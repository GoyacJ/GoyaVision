package engine

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"goyavision/internal/domain/operator"
	"goyavision/internal/domain/workflow"
	"goyavision/internal/port"

	"github.com/google/uuid"
)

var _ port.WorkflowEngine = (*SimpleWorkflowEngine)(nil)

// SimpleWorkflowEngine 简单工作流引擎（支持单算子执行）
type SimpleWorkflowEngine struct {
	repo     port.Repository
	registry port.ExecutorRegistry
	tasks    map[uuid.UUID]*taskExecution
	mu       sync.RWMutex
}

type taskExecution struct {
	ctx            context.Context
	cancel         context.CancelFunc
	progress       int
	currentNode    string
	nodeExecutions map[string]*workflow.NodeExecution
}

// NewSimpleWorkflowEngine 创建简单工作流引擎
func NewSimpleWorkflowEngine(repo port.Repository, executor port.OperatorExecutor) *SimpleWorkflowEngine {
	registry := NewExecutorRegistry()
	registry.Register(executor.Mode(), executor)

	return &SimpleWorkflowEngine{
		repo:     repo,
		registry: registry,
		tasks:    make(map[uuid.UUID]*taskExecution),
	}
}

// Execute 执行工作流
func (e *SimpleWorkflowEngine) Execute(ctx context.Context, wf *workflow.Workflow, task *workflow.Task) error {
	if len(wf.Nodes) == 0 {
		return errors.New("workflow has no nodes")
	}

	execCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	nodeExecs := make(map[string]*workflow.NodeExecution)
	for _, node := range wf.Nodes {
		nodeExecs[node.NodeKey] = &workflow.NodeExecution{
			NodeKey: node.NodeKey,
			Status:  workflow.NodeExecPending,
		}
	}

	e.mu.Lock()
	e.tasks[task.ID] = &taskExecution{
		ctx:            execCtx,
		cancel:         cancel,
		progress:       0,
		nodeExecutions: nodeExecs,
	}
	e.mu.Unlock()

	defer func() {
		e.mu.Lock()
		delete(e.tasks, task.ID)
		e.mu.Unlock()
	}()

	now := time.Now()
	task.Status = workflow.TaskStatusRunning
	task.StartedAt = &now
	task.Progress = 0
	task.NodeExecutions = e.collectNodeExecutions(task.ID)
	if err := e.repo.UpdateTask(ctx, task); err != nil {
		return err
	}

	totalNodes := len(wf.Nodes)
	for i, node := range wf.Nodes {
		select {
		case <-execCtx.Done():
			return execCtx.Err()
		default:
		}

		if !e.shouldExecuteNode(&node, wf.Edges, task.ID) {
			e.setNodeStatus(task.ID, node.NodeKey, workflow.NodeExecSkipped, "")
			task.NodeExecutions = e.collectNodeExecutions(task.ID)
			if err := e.repo.UpdateTask(ctx, task); err != nil {
				return err
			}
			continue
		}

		e.mu.Lock()
		if exec, ok := e.tasks[task.ID]; ok {
			exec.currentNode = node.NodeKey
			exec.progress = int((float64(i) / float64(totalNodes)) * 100)
		}
		e.mu.Unlock()

		task.CurrentNode = node.NodeKey
		task.Progress = int((float64(i) / float64(totalNodes)) * 100)

		startedAt := time.Now()
		e.mu.Lock()
		if exec, ok := e.tasks[task.ID]; ok {
			if ne, ok := exec.nodeExecutions[node.NodeKey]; ok {
				ne.Status = workflow.NodeExecRunning
				ne.StartedAt = &startedAt
			}
		}
		e.mu.Unlock()

		task.NodeExecutions = e.collectNodeExecutions(task.ID)
		if err := e.repo.UpdateTask(ctx, task); err != nil {
			return err
		}

		if node.OperatorID == nil {
			e.setNodeStatus(task.ID, node.NodeKey, workflow.NodeExecSuccess, "")
			completedAt := time.Now()
			e.mu.Lock()
			if exec, ok := e.tasks[task.ID]; ok {
				if ne, ok := exec.nodeExecutions[node.NodeKey]; ok {
					ne.CompletedAt = &completedAt
				}
			}
			e.mu.Unlock()
			task.NodeExecutions = e.collectNodeExecutions(task.ID)
			if err := e.repo.UpdateTask(ctx, task); err != nil {
				return err
			}
			continue
		}

		op, err := e.repo.GetOperator(ctx, *node.OperatorID)
		if err != nil {
			e.setNodeStatus(task.ID, node.NodeKey, workflow.NodeExecFailed, fmt.Sprintf("failed to get operator: %v", err))
			task.NodeExecutions = e.collectNodeExecutions(task.ID)
			_ = e.repo.UpdateTask(ctx, task)
			return fmt.Errorf("failed to get operator: %w", err)
		}

		version := op.ActiveVersion
		if version == nil {
			errMsg := fmt.Sprintf("operator %s has no active version", op.Code)
			e.setNodeStatus(task.ID, node.NodeKey, workflow.NodeExecFailed, errMsg)
			task.NodeExecutions = e.collectNodeExecutions(task.ID)
			_ = e.repo.UpdateTask(ctx, task)
			return fmt.Errorf("operator %s has no active version", op.Code)
		}

		executor, err := e.registry.Get(version.ExecMode)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get executor for mode %s: %v", version.ExecMode, err)
			e.setNodeStatus(task.ID, node.NodeKey, workflow.NodeExecFailed, errMsg)
			task.NodeExecutions = e.collectNodeExecutions(task.ID)
			_ = e.repo.UpdateTask(ctx, task)
			return fmt.Errorf("failed to get executor for mode %s: %w", version.ExecMode, err)
		}

		inputParams := task.InputParams
		if inputParams == nil {
			inputParams = make(map[string]interface{})
		}
		if node.Config != nil && node.Config.Params != nil {
			for k, v := range node.Config.Params {
				inputParams[k] = v
			}
		}

		var assetID uuid.UUID
		if task.AssetID != nil {
			assetID = *task.AssetID
		}
		input := &operator.Input{
			AssetID: assetID,
			Params:  inputParams,
		}

		output, err := executor.Execute(execCtx, version, input)
		if err != nil {
			errMsg := fmt.Sprintf("failed to execute operator %s: %v", op.Code, err)
			e.setNodeStatus(task.ID, node.NodeKey, workflow.NodeExecFailed, errMsg)
			task.NodeExecutions = e.collectNodeExecutions(task.ID)
			_ = e.repo.UpdateTask(ctx, task)
			return fmt.Errorf("failed to execute operator %s: %w", op.Code, err)
		}

		artifactIDs, err := e.saveArtifacts(ctx, task.ID, node.NodeKey, output)
		if err != nil {
			errMsg := fmt.Sprintf("failed to save artifacts: %v", err)
			e.setNodeStatus(task.ID, node.NodeKey, workflow.NodeExecFailed, errMsg)
			task.NodeExecutions = e.collectNodeExecutions(task.ID)
			_ = e.repo.UpdateTask(ctx, task)
			return fmt.Errorf("failed to save artifacts: %w", err)
		}

		completedAt := time.Now()
		e.mu.Lock()
		if exec, ok := e.tasks[task.ID]; ok {
			if ne, ok := exec.nodeExecutions[node.NodeKey]; ok {
				ne.Status = workflow.NodeExecSuccess
				ne.CompletedAt = &completedAt
				ne.ArtifactIDs = artifactIDs
			}
		}
		e.mu.Unlock()

		task.NodeExecutions = e.collectNodeExecutions(task.ID)
		if err := e.repo.UpdateTask(ctx, task); err != nil {
			return err
		}
	}

	task.Status = workflow.TaskStatusSuccess
	task.Progress = 100
	completedAt := time.Now()
	task.CompletedAt = &completedAt
	task.NodeExecutions = e.collectNodeExecutions(task.ID)
	if err := e.repo.UpdateTask(ctx, task); err != nil {
		return err
	}

	return nil
}

// shouldExecuteNode 检查节点是否应该执行（基于入边条件）
func (e *SimpleWorkflowEngine) shouldExecuteNode(node *workflow.Node, edges []workflow.Edge, taskID uuid.UUID) bool {
	e.mu.RLock()
	exec, ok := e.tasks[taskID]
	e.mu.RUnlock()
	if !ok {
		return true
	}

	hasInEdge := false
	for _, edge := range edges {
		if edge.TargetKey != node.NodeKey {
			continue
		}
		hasInEdge = true

		sourceExec, ok := exec.nodeExecutions[edge.SourceKey]
		if !ok {
			continue
		}

		condType := "always"
		if edge.Condition != nil && edge.Condition.Type != "" {
			condType = edge.Condition.Type
		}

		switch condType {
		case "always":
			continue
		case "on_success":
			if sourceExec.Status != workflow.NodeExecSuccess {
				return false
			}
		case "on_failure":
			if sourceExec.Status != workflow.NodeExecFailed {
				return false
			}
		default:
			continue
		}
	}

	if !hasInEdge {
		return true
	}

	return true
}

// setNodeStatus 设置节点执行状态
func (e *SimpleWorkflowEngine) setNodeStatus(taskID uuid.UUID, nodeKey string, status workflow.NodeExecutionStatus, errMsg string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	exec, ok := e.tasks[taskID]
	if !ok {
		return
	}
	ne, ok := exec.nodeExecutions[nodeKey]
	if !ok {
		return
	}
	ne.Status = status
	if errMsg != "" {
		ne.Error = errMsg
	}
	if status == workflow.NodeExecSuccess || status == workflow.NodeExecFailed || status == workflow.NodeExecSkipped {
		now := time.Now()
		ne.CompletedAt = &now
	}
}

// collectNodeExecutions 收集节点执行状态快照
func (e *SimpleWorkflowEngine) collectNodeExecutions(taskID uuid.UUID) []workflow.NodeExecution {
	e.mu.RLock()
	defer e.mu.RUnlock()

	exec, ok := e.tasks[taskID]
	if !ok {
		return nil
	}

	result := make([]workflow.NodeExecution, 0, len(exec.nodeExecutions))
	for _, ne := range exec.nodeExecutions {
		result = append(result, *ne)
	}
	return result
}

// Cancel 取消工作流执行
func (e *SimpleWorkflowEngine) Cancel(ctx context.Context, taskID uuid.UUID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	exec, ok := e.tasks[taskID]
	if !ok {
		return errors.New("task is not running")
	}

	exec.cancel()
	return nil
}

// GetProgress 获取工作流执行进度
func (e *SimpleWorkflowEngine) GetProgress(ctx context.Context, taskID uuid.UUID) (int, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	exec, ok := e.tasks[taskID]
	if !ok {
		task, err := e.repo.GetTask(ctx, taskID)
		if err != nil {
			return 0, err
		}
		return task.Progress, nil
	}

	return exec.progress, nil
}

// saveArtifacts 保存产物并返回产物 ID 列表
func (e *SimpleWorkflowEngine) saveArtifacts(ctx context.Context, taskID uuid.UUID, nodeKey string, output *operator.Output) ([]uuid.UUID, error) {
	if output == nil {
		return nil, nil
	}

	var artifactIDs []uuid.UUID

	for _, asset := range output.OutputAssets {
		artifact := &workflow.Artifact{
			TaskID: taskID,
			Type:   workflow.ArtifactTypeAsset,
			Data: &workflow.ArtifactData{
				AssetInfo: &workflow.AssetInfo{
					Type:     string(asset.Type),
					Path:     asset.Path,
					Format:   asset.Format,
					Metadata: asset.Metadata,
				},
				Metadata: map[string]interface{}{
					"node_key": nodeKey,
				},
			},
		}

		if err := e.repo.CreateArtifact(ctx, artifact); err != nil {
			return artifactIDs, err
		}
		artifactIDs = append(artifactIDs, artifact.ID)
	}

	if len(output.Results) > 0 {
		results := make([]map[string]interface{}, len(output.Results))
		for i, r := range output.Results {
			results[i] = map[string]interface{}{
				"type":       r.Type,
				"data":       r.Data,
				"confidence": r.Confidence,
			}
		}

		artifact := &workflow.Artifact{
			TaskID: taskID,
			Type:   workflow.ArtifactTypeResult,
			Data: &workflow.ArtifactData{
				Results: results,
				Metadata: map[string]interface{}{
					"node_key": nodeKey,
				},
			},
		}

		if err := e.repo.CreateArtifact(ctx, artifact); err != nil {
			return artifactIDs, err
		}
		artifactIDs = append(artifactIDs, artifact.ID)
	}

	if len(output.Timeline) > 0 {
		timeline := make([]workflow.TimelineSegment, len(output.Timeline))
		for i, t := range output.Timeline {
			timeline[i] = workflow.TimelineSegment{
				Start:      t.Start,
				End:        t.End,
				EventType:  t.EventType,
				Confidence: t.Confidence,
				Data:       t.Data,
			}
		}

		artifact := &workflow.Artifact{
			TaskID: taskID,
			Type:   workflow.ArtifactTypeTimeline,
			Data: &workflow.ArtifactData{
				Timeline: timeline,
				Metadata: map[string]interface{}{
					"node_key": nodeKey,
				},
			},
		}

		if err := e.repo.CreateArtifact(ctx, artifact); err != nil {
			return artifactIDs, err
		}
		artifactIDs = append(artifactIDs, artifact.ID)
	}

	return artifactIDs, nil
}
