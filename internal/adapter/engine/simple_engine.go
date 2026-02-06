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
	repo      port.Repository
	registry  port.ExecutorRegistry
	tasks     map[uuid.UUID]*taskExecution
	mu        sync.RWMutex
}

type taskExecution struct {
	ctx         context.Context
	cancel      context.CancelFunc
	progress    int
	currentNode string
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

	e.mu.Lock()
	e.tasks[task.ID] = &taskExecution{
		ctx:      execCtx,
		cancel:   cancel,
		progress: 0,
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

		e.mu.Lock()
		if exec, ok := e.tasks[task.ID]; ok {
			exec.currentNode = node.NodeKey
			exec.progress = int((float64(i) / float64(totalNodes)) * 100)
		}
		e.mu.Unlock()

		task.CurrentNode = node.NodeKey
		task.Progress = int((float64(i) / float64(totalNodes)) * 100)
		if err := e.repo.UpdateTask(ctx, task); err != nil {
			return err
		}

		if node.OperatorID == nil {
			continue
		}

		op, err := e.repo.GetOperator(ctx, *node.OperatorID)
		if err != nil {
			return fmt.Errorf("failed to get operator: %w", err)
		}

		version := op.ActiveVersion
		if version == nil {
			return fmt.Errorf("operator %s has no active version", op.Code)
		}

		executor, err := e.registry.Get(version.ExecMode)
		if err != nil {
			return fmt.Errorf("failed to get executor for mode %s: %w", version.ExecMode, err)
		}

		inputParams := task.InputParams
		if inputParams == nil {
			inputParams = make(map[string]interface{})
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
			return fmt.Errorf("failed to execute operator %s: %w", op.Code, err)
		}

		if err := e.saveArtifacts(ctx, task.ID, output); err != nil {
			return fmt.Errorf("failed to save artifacts: %w", err)
		}
	}

	task.Status = workflow.TaskStatusSuccess
	task.Progress = 100
	completedAt := time.Now()
	task.CompletedAt = &completedAt
	if err := e.repo.UpdateTask(ctx, task); err != nil {
		return err
	}

	return nil
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

// saveArtifacts 保存产物
func (e *SimpleWorkflowEngine) saveArtifacts(ctx context.Context, taskID uuid.UUID, output *operator.Output) error {
	if output == nil {
		return nil
	}

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
			},
		}

		if err := e.repo.CreateArtifact(ctx, artifact); err != nil {
			return err
		}
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
			},
		}

		if err := e.repo.CreateArtifact(ctx, artifact); err != nil {
			return err
		}
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
			},
		}

		if err := e.repo.CreateArtifact(ctx, artifact); err != nil {
			return err
		}
	}

	return nil
}
