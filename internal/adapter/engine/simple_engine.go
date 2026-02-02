package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
)

var _ port.WorkflowEngine = (*SimpleWorkflowEngine)(nil)

// SimpleWorkflowEngine 简单工作流引擎（支持单算子执行）
type SimpleWorkflowEngine struct {
	repo     port.Repository
	executor port.OperatorExecutor
	tasks    map[uuid.UUID]*taskExecution
	mu       sync.RWMutex
}

type taskExecution struct {
	ctx        context.Context
	cancel     context.CancelFunc
	progress   int
	currentNode string
}

// NewSimpleWorkflowEngine 创建简单工作流引擎
func NewSimpleWorkflowEngine(repo port.Repository, executor port.OperatorExecutor) *SimpleWorkflowEngine {
	return &SimpleWorkflowEngine{
		repo:     repo,
		executor: executor,
		tasks:    make(map[uuid.UUID]*taskExecution),
	}
}

// Execute 执行工作流
func (e *SimpleWorkflowEngine) Execute(ctx context.Context, workflow *domain.Workflow, task *domain.Task) error {
	if len(workflow.Nodes) == 0 {
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
	task.Status = domain.TaskStatusRunning
	task.StartedAt = &now
	task.Progress = 0
	if err := e.repo.UpdateTask(ctx, task); err != nil {
		return err
	}

	totalNodes := len(workflow.Nodes)
	for i, node := range workflow.Nodes {
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

		operator, err := e.repo.GetOperator(ctx, *node.OperatorID)
		if err != nil {
			return fmt.Errorf("failed to get operator: %w", err)
		}

		var inputParams map[string]interface{}
		if task.InputParams != nil && len(task.InputParams) > 0 {
			if err := json.Unmarshal(task.InputParams, &inputParams); err != nil {
				return fmt.Errorf("failed to unmarshal input params: %w", err)
			}
		}

		input := &domain.OperatorInput{
			Params: inputParams,
		}

		if task.AssetID != nil {
			input.AssetID = *task.AssetID
		}

		output, err := e.executor.Execute(execCtx, operator, input)
		if err != nil {
			return fmt.Errorf("failed to execute operator %s: %w", operator.Code, err)
		}

		if err := e.saveArtifacts(ctx, task.ID, output); err != nil {
			return fmt.Errorf("failed to save artifacts: %w", err)
		}
	}

	task.Status = domain.TaskStatusSuccess
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
func (e *SimpleWorkflowEngine) saveArtifacts(ctx context.Context, taskID uuid.UUID, output *domain.OperatorOutput) error {
	if output == nil {
		return nil
	}

	for _, asset := range output.OutputAssets {
		data := map[string]interface{}{
			"type":     asset.Type,
			"path":     asset.Path,
			"format":   asset.Format,
			"metadata": asset.Metadata,
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}

		artifact := &domain.Artifact{
			TaskID: taskID,
			Type:   domain.ArtifactTypeAsset,
		}
		artifact.Data = dataBytes

		if err := e.repo.CreateArtifact(ctx, artifact); err != nil {
			return err
		}
	}

	if len(output.Results) > 0 {
		data := map[string]interface{}{
			"results": output.Results,
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}

		artifact := &domain.Artifact{
			TaskID: taskID,
			Type:   domain.ArtifactTypeResult,
		}
		artifact.Data = dataBytes

		if err := e.repo.CreateArtifact(ctx, artifact); err != nil {
			return err
		}
	}

	if len(output.Timeline) > 0 {
		data := map[string]interface{}{
			"timeline": output.Timeline,
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}

		artifact := &domain.Artifact{
			TaskID: taskID,
			Type:   domain.ArtifactTypeTimeline,
		}
		artifact.Data = dataBytes

		if err := e.repo.CreateArtifact(ctx, artifact); err != nil {
			return err
		}
	}

	return nil
}
