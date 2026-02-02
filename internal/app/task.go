package app

import (
	"context"
	"errors"
	"time"

	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	WorkflowID  uuid.UUID              `json:"workflow_id"`
	AssetID     *uuid.UUID             `json:"asset_id,omitempty"`
	InputParams map[string]interface{} `json:"input_params,omitempty"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Status      *domain.TaskStatus `json:"status,omitempty"`
	Progress    *int               `json:"progress,omitempty"`
	CurrentNode *string            `json:"current_node,omitempty"`
	Error       *string            `json:"error,omitempty"`
}

// ListTasksRequest 列出任务请求
type ListTasksRequest struct {
	WorkflowID *uuid.UUID
	AssetID    *uuid.UUID
	Status     *domain.TaskStatus
	From       *time.Time
	To         *time.Time
	Limit      int
	Offset     int
}

type TaskService struct {
	repo port.Repository
}

func NewTaskService(repo port.Repository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

// Create 创建任务
func (s *TaskService) Create(ctx context.Context, req *CreateTaskRequest) (*domain.Task, error) {
	if req.WorkflowID == uuid.Nil {
		return nil, errors.New("workflow_id is required")
	}

	workflow, err := s.repo.GetWorkflow(ctx, req.WorkflowID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("workflow not found")
		}
		return nil, err
	}

	if !workflow.IsEnabled() {
		return nil, errors.New("workflow is not enabled")
	}

	if req.AssetID != nil {
		if _, err := s.repo.GetMediaAsset(ctx, *req.AssetID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("asset not found")
			}
			return nil, err
		}
	}

	task := &domain.Task{
		WorkflowID: req.WorkflowID,
		AssetID:    req.AssetID,
		Status:     domain.TaskStatusPending,
		Progress:   0,
	}

	if err := s.repo.CreateTask(ctx, task); err != nil {
		return nil, err
	}

	return s.repo.GetTaskWithRelations(ctx, task.ID)
}

// Get 获取任务
func (s *TaskService) Get(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	task, err := s.repo.GetTask(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return task, nil
}

// GetWithRelations 获取任务及其关联数据
func (s *TaskService) GetWithRelations(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	task, err := s.repo.GetTaskWithRelations(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return task, nil
}

// List 列出任务
func (s *TaskService) List(ctx context.Context, req *ListTasksRequest) ([]*domain.Task, int64, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 1000 {
		req.Limit = 1000
	}

	filter := domain.TaskFilter{
		WorkflowID: req.WorkflowID,
		AssetID:    req.AssetID,
		Status:     req.Status,
		From:       req.From,
		To:         req.To,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}

	return s.repo.ListTasks(ctx, filter)
}

// Update 更新任务
func (s *TaskService) Update(ctx context.Context, id uuid.UUID, req *UpdateTaskRequest) (*domain.Task, error) {
	task, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		if *req.Status == domain.TaskStatusRunning && task.StartedAt == nil {
			now := time.Now()
			task.StartedAt = &now
		}
		if (*req.Status == domain.TaskStatusSuccess || *req.Status == domain.TaskStatusFailed || *req.Status == domain.TaskStatusCancelled) && task.CompletedAt == nil {
			now := time.Now()
			task.CompletedAt = &now
		}
		task.Status = *req.Status
	}

	if req.Progress != nil {
		if *req.Progress < 0 {
			*req.Progress = 0
		}
		if *req.Progress > 100 {
			*req.Progress = 100
		}
		task.Progress = *req.Progress
	}

	if req.CurrentNode != nil {
		task.CurrentNode = *req.CurrentNode
	}

	if req.Error != nil {
		task.Error = *req.Error
	}

	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// Delete 删除任务
func (s *TaskService) Delete(ctx context.Context, id uuid.UUID) error {
	task, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if task.IsRunning() {
		return errors.New("cannot delete running task")
	}

	return s.repo.DeleteTask(ctx, id)
}

// Start 启动任务
func (s *TaskService) Start(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	task, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !task.IsPending() {
		return nil, errors.New("task is not pending")
	}

	now := time.Now()
	task.Status = domain.TaskStatusRunning
	task.StartedAt = &now
	task.Progress = 0

	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// Complete 完成任务
func (s *TaskService) Complete(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	task, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !task.IsRunning() {
		return nil, errors.New("task is not running")
	}

	now := time.Now()
	task.Status = domain.TaskStatusSuccess
	task.CompletedAt = &now
	task.Progress = 100

	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// Fail 任务失败
func (s *TaskService) Fail(ctx context.Context, id uuid.UUID, errorMsg string) (*domain.Task, error) {
	task, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !task.IsRunning() {
		return nil, errors.New("task is not running")
	}

	now := time.Now()
	task.Status = domain.TaskStatusFailed
	task.CompletedAt = &now
	task.Error = errorMsg

	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// Cancel 取消任务
func (s *TaskService) Cancel(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	task, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if task.IsCompleted() {
		return nil, errors.New("task is already completed")
	}

	now := time.Now()
	task.Status = domain.TaskStatusCancelled
	if task.CompletedAt == nil {
		task.CompletedAt = &now
	}

	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// GetStats 获取任务统计
func (s *TaskService) GetStats(ctx context.Context, workflowID *uuid.UUID) (*domain.TaskStats, error) {
	return s.repo.GetTaskStats(ctx, workflowID)
}

// ListRunning 列出所有运行中的任务
func (s *TaskService) ListRunning(ctx context.Context) ([]*domain.Task, error) {
	return s.repo.ListRunningTasks(ctx)
}
