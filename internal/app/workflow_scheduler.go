package app

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

// WorkflowScheduler 工作流调度器
type WorkflowScheduler struct {
	scheduler gocron.Scheduler
	repo      port.Repository
	engine    port.WorkflowEngine
	jobs      map[uuid.UUID]gocron.Job
	jobsMu    sync.RWMutex
}

// NewWorkflowScheduler 创建工作流调度器
func NewWorkflowScheduler(repo port.Repository, engine port.WorkflowEngine) (*WorkflowScheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("create scheduler: %w", err)
	}

	return &WorkflowScheduler{
		scheduler: s,
		repo:      repo,
		engine:    engine,
		jobs:      make(map[uuid.UUID]gocron.Job),
	}, nil
}

// Start 启动调度器
func (s *WorkflowScheduler) Start(ctx context.Context) error {
	s.scheduler.Start()
	return s.loadAndScheduleWorkflows(ctx)
}

// Stop 停止调度器
func (s *WorkflowScheduler) Stop() error {
	return s.scheduler.Shutdown()
}

// loadAndScheduleWorkflows 加载并调度所有启用的工作流
func (s *WorkflowScheduler) loadAndScheduleWorkflows(ctx context.Context) error {
	workflows, err := s.repo.ListEnabledWorkflows(ctx)
	if err != nil {
		return fmt.Errorf("list enabled workflows: %w", err)
	}

	for _, workflow := range workflows {
		if workflow.TriggerType == domain.TriggerTypeSchedule {
			if err := s.ScheduleWorkflow(ctx, workflow); err != nil {
				continue
			}
		}
	}

	return nil
}

// ScheduleWorkflow 调度工作流
func (s *WorkflowScheduler) ScheduleWorkflow(ctx context.Context, workflow *domain.Workflow) error {
	s.jobsMu.Lock()
	defer s.jobsMu.Unlock()

	if _, exists := s.jobs[workflow.ID]; exists {
		return nil
	}

	var triggerConf domain.TriggerConfig
	if workflow.TriggerConf != nil && len(workflow.TriggerConf) > 0 {
		if err := json.Unmarshal(workflow.TriggerConf, &triggerConf); err != nil {
			return fmt.Errorf("parse trigger config: %w", err)
		}
	}

	job, err := s.createJob(workflow, &triggerConf)
	if err != nil {
		return fmt.Errorf("create job: %w", err)
	}

	s.jobs[workflow.ID] = job
	return nil
}

// UnscheduleWorkflow 取消调度工作流
func (s *WorkflowScheduler) UnscheduleWorkflow(workflowID uuid.UUID) error {
	s.jobsMu.Lock()
	defer s.jobsMu.Unlock()

	job, exists := s.jobs[workflowID]
	if !exists {
		return nil
	}

	if err := s.scheduler.RemoveJob(job.ID()); err != nil {
		return fmt.Errorf("remove job: %w", err)
	}

	delete(s.jobs, workflowID)
	return nil
}

// createJob 创建调度任务
func (s *WorkflowScheduler) createJob(workflow *domain.Workflow, triggerConf *domain.TriggerConfig) (gocron.Job, error) {
	if triggerConf.Schedule != "" {
		return s.createCronJob(workflow, triggerConf)
	}

	if triggerConf.IntervalSec > 0 {
		return s.createIntervalJob(workflow, triggerConf)
	}

	return nil, fmt.Errorf("invalid trigger config: no schedule or interval specified")
}

// createIntervalJob 创建间隔任务
func (s *WorkflowScheduler) createIntervalJob(workflow *domain.Workflow, triggerConf *domain.TriggerConfig) (gocron.Job, error) {
	duration := time.Duration(triggerConf.IntervalSec) * time.Second

	job, err := s.scheduler.NewJob(
		gocron.DurationJob(duration),
		gocron.NewTask(s.runWorkflow, workflow.ID),
	)
	if err != nil {
		return nil, fmt.Errorf("create interval job: %w", err)
	}

	return job, nil
}

// createCronJob 创建 Cron 任务
func (s *WorkflowScheduler) createCronJob(workflow *domain.Workflow, triggerConf *domain.TriggerConfig) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.CronJob(triggerConf.Schedule, false),
		gocron.NewTask(s.runWorkflow, workflow.ID),
	)
	if err != nil {
		return nil, fmt.Errorf("create cron job: %w", err)
	}

	return job, nil
}

// runWorkflow 执行工作流
func (s *WorkflowScheduler) runWorkflow(workflowID uuid.UUID) {
	ctx := context.Background()

	workflow, err := s.repo.GetWorkflowWithNodes(ctx, workflowID)
	if err != nil {
		return
	}

	if !workflow.IsEnabled() {
		s.UnscheduleWorkflow(workflowID)
		return
	}

	task := &domain.Task{
		WorkflowID: workflow.ID,
		Status:     domain.TaskStatusPending,
		Progress:   0,
	}

	if err := s.repo.CreateTask(ctx, task); err != nil {
		return
	}

	go func() {
		if err := s.engine.Execute(context.Background(), workflow, task); err != nil {
			now := time.Now()
			task.Status = domain.TaskStatusFailed
			task.Error = err.Error()
			task.CompletedAt = &now
			s.repo.UpdateTask(context.Background(), task)
		}
	}()
}

// TriggerWorkflow 手动触发工作流
func (s *WorkflowScheduler) TriggerWorkflow(ctx context.Context, workflowID uuid.UUID, assetID *uuid.UUID) (*domain.Task, error) {
	workflow, err := s.repo.GetWorkflowWithNodes(ctx, workflowID)
	if err != nil {
		return nil, fmt.Errorf("get workflow: %w", err)
	}

	if !workflow.IsEnabled() {
		return nil, fmt.Errorf("workflow is not enabled")
	}

	var inputParams map[string]interface{}
	if assetID != nil {
		inputParams = map[string]interface{}{
			"asset_id": assetID.String(),
		}
	}

	inputParamsJSON, _ := json.Marshal(inputParams)

	task := &domain.Task{
		WorkflowID:  workflow.ID,
		AssetID:     assetID,
		Status:      domain.TaskStatusPending,
		Progress:    0,
		InputParams: inputParamsJSON,
	}

	if err := s.repo.CreateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	go func() {
		if err := s.engine.Execute(context.Background(), workflow, task); err != nil {
			now := time.Now()
			task.Status = domain.TaskStatusFailed
			task.Error = err.Error()
			task.CompletedAt = &now
			s.repo.UpdateTask(context.Background(), task)
		}
	}()

	return task, nil
}
