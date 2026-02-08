package app

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"goyavision/internal/app/event"
	"goyavision/internal/domain/workflow"
	"goyavision/internal/port"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

// WorkflowScheduler 工作流调度器
type WorkflowScheduler struct {
	scheduler gocron.Scheduler
	repo      port.Repository
	engine    port.WorkflowEngine
	eventBus  port.EventBus
	jobs      map[uuid.UUID]gocron.Job
	jobsMu    sync.RWMutex
}

// NewWorkflowScheduler 创建工作流调度器。eventBus 可选，为 nil 时不启用事件触发。
func NewWorkflowScheduler(repo port.Repository, engine port.WorkflowEngine, eventBus port.EventBus) (*WorkflowScheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("create scheduler: %w", err)
	}

	ws := &WorkflowScheduler{
		scheduler: s,
		repo:      repo,
		engine:    engine,
		eventBus:  eventBus,
		jobs:      make(map[uuid.UUID]gocron.Job),
	}
	return ws, nil
}

// Start 启动调度器
func (s *WorkflowScheduler) Start(ctx context.Context) error {
	s.scheduler.Start()
	if err := s.loadAndScheduleWorkflows(ctx); err != nil {
		return err
	}
	if s.eventBus != nil {
		s.eventBus.Subscribe(event.EventTypeAssetNew, s.handleAssetNew)
		s.eventBus.Subscribe(event.EventTypeAssetDone, s.handleAssetDone)
	}
	return nil
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

	for _, wf := range workflows {
		if wf.TriggerType == workflow.TriggerTypeSchedule {
			if err := s.ScheduleWorkflow(ctx, wf); err != nil {
				continue
			}
		}
	}

	return nil
}

// ScheduleWorkflow 调度工作流
func (s *WorkflowScheduler) ScheduleWorkflow(ctx context.Context, wf *workflow.Workflow) error {
	s.jobsMu.Lock()
	defer s.jobsMu.Unlock()

	if _, exists := s.jobs[wf.ID]; exists {
		return nil
	}

	var triggerConf *workflow.TriggerConfig
	if wf.TriggerConf != nil {
		triggerConf = wf.TriggerConf
	}

	job, err := s.createJob(wf, triggerConf)
	if err != nil {
		return fmt.Errorf("create job: %w", err)
	}

	s.jobs[wf.ID] = job
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
func (s *WorkflowScheduler) createJob(workflow *workflow.Workflow, triggerConf *workflow.TriggerConfig) (gocron.Job, error) {
	if triggerConf.Schedule != "" {
		return s.createCronJob(workflow, triggerConf)
	}

	if triggerConf.IntervalSec > 0 {
		return s.createIntervalJob(workflow, triggerConf)
	}

	return nil, fmt.Errorf("invalid trigger config: no schedule or interval specified")
}

// createIntervalJob 创建间隔任务
func (s *WorkflowScheduler) createIntervalJob(wf *workflow.Workflow, triggerConf *workflow.TriggerConfig) (gocron.Job, error) {
	duration := time.Duration(triggerConf.IntervalSec) * time.Second

	job, err := s.scheduler.NewJob(
		gocron.DurationJob(duration),
		gocron.NewTask(s.runWorkflow, wf.ID),
	)
	if err != nil {
		return nil, fmt.Errorf("create interval job: %w", err)
	}

	return job, nil
}

// createCronJob 创建 Cron 任务
func (s *WorkflowScheduler) createCronJob(wf *workflow.Workflow, triggerConf *workflow.TriggerConfig) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.CronJob(triggerConf.Schedule, false),
		gocron.NewTask(s.runWorkflow, wf.ID),
	)
	if err != nil {
		return nil, fmt.Errorf("create cron job: %w", err)
	}

	return job, nil
}

// runWorkflow 执行工作流（由定时任务调用，无请求 context）
func (s *WorkflowScheduler) runWorkflow(workflowID uuid.UUID) {
	ctx := context.Background()

	wf, err := s.repo.GetWorkflowWithNodes(ctx, workflowID)
	if err != nil {
		log.Printf("[WorkflowScheduler] runWorkflow: get workflow %s: %v", workflowID, err)
		return
	}

	if !wf.IsEnabled() {
		s.UnscheduleWorkflow(workflowID)
		return
	}

	task := &workflow.Task{
		WorkflowID: wf.ID,
		Status:     workflow.TaskStatusPending,
		Progress:   0,
	}

	if err := s.repo.CreateTask(ctx, task); err != nil {
		log.Printf("[WorkflowScheduler] runWorkflow: create task for workflow %s: %v", workflowID, err)
		return
	}

	go func() {
		runCtx := context.Background()
		if err := s.engine.Execute(runCtx, wf, task); err != nil {
			log.Printf("[WorkflowScheduler] runWorkflow: execute failed workflow=%s task=%s: %v", workflowID, task.ID, err)
			now := time.Now()
			task.Status = workflow.TaskStatusFailed
			task.Error = err.Error()
			task.CompletedAt = &now
			if updateErr := s.repo.UpdateTask(runCtx, task); updateErr != nil {
				log.Printf("[WorkflowScheduler] runWorkflow: update task status failed task=%s: %v", task.ID, updateErr)
			}
		}
	}()
}

// TriggerWorkflow 手动触发工作流
func (s *WorkflowScheduler) TriggerWorkflow(ctx context.Context, workflowID uuid.UUID, assetID *uuid.UUID) (*workflow.Task, error) {
	wf, err := s.repo.GetWorkflowWithNodes(ctx, workflowID)
	if err != nil {
		return nil, fmt.Errorf("get workflow: %w", err)
	}

	if !wf.IsEnabled() {
		return nil, fmt.Errorf("workflow is not enabled")
	}

	var inputParams map[string]interface{}
	if assetID != nil {
		inputParams = map[string]interface{}{
			"asset_id": assetID.String(),
		}
	}

	task := &workflow.Task{
		WorkflowID:  wf.ID,
		AssetID:     assetID,
		Status:      workflow.TaskStatusPending,
		Progress:    0,
		InputParams: inputParams,
	}

	if err := s.repo.CreateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	goCtx := context.WithoutCancel(ctx)
	go func() {
		if err := s.engine.Execute(goCtx, wf, task); err != nil {
			log.Printf("[WorkflowScheduler] TriggerWorkflow: execute failed workflow=%s task=%s: %v", workflowID, task.ID, err)
			now := time.Now()
			task.Status = workflow.TaskStatusFailed
			task.Error = err.Error()
			task.CompletedAt = &now
			if updateErr := s.repo.UpdateTask(goCtx, task); updateErr != nil {
				log.Printf("[WorkflowScheduler] TriggerWorkflow: update task status failed task=%s: %v", task.ID, updateErr)
			}
		}
	}()

	return task, nil
}

func (s *WorkflowScheduler) handleAssetNew(ctx context.Context, ev port.Event) error {
	e, ok := ev.(*event.AssetCreatedEvent)
	if !ok {
		return nil
	}
	return s.triggerWorkflowsByEvent(ctx, e.AssetID, workflow.TriggerTypeAssetNew)
}

func (s *WorkflowScheduler) handleAssetDone(ctx context.Context, ev port.Event) error {
	e, ok := ev.(*event.AssetDoneEvent)
	if !ok {
		return nil
	}
	return s.triggerWorkflowsByEvent(ctx, e.AssetID, workflow.TriggerTypeAssetDone)
}

func (s *WorkflowScheduler) triggerWorkflowsByEvent(ctx context.Context, assetID uuid.UUID, triggerType workflow.TriggerType) error {
	workflows, err := s.repo.ListEnabledWorkflows(ctx)
	if err != nil {
		log.Printf("[WorkflowScheduler] triggerWorkflowsByEvent: list workflows: %v", err)
		return err
	}
	assetIDPtr := &assetID
	for _, wf := range workflows {
		if wf.TriggerType != triggerType {
			continue
		}
		wfWithNodes, err := s.repo.GetWorkflowWithNodes(ctx, wf.ID)
		if err != nil {
			log.Printf("[WorkflowScheduler] triggerWorkflowsByEvent: get workflow %s: %v", wf.ID, err)
			continue
		}
		task := &workflow.Task{
			WorkflowID:  wfWithNodes.ID,
			AssetID:     assetIDPtr,
			Status:      workflow.TaskStatusPending,
			Progress:    0,
			InputParams: map[string]interface{}{"asset_id": assetID.String()},
		}
		if err := s.repo.CreateTask(ctx, task); err != nil {
			log.Printf("[WorkflowScheduler] triggerWorkflowsByEvent: create task workflow=%s: %v", wf.ID, err)
			continue
		}
		goCtx := context.WithoutCancel(ctx)
		go func(w *workflow.Workflow, t *workflow.Task) {
			if err := s.engine.Execute(goCtx, w, t); err != nil {
				log.Printf("[WorkflowScheduler] triggerWorkflowsByEvent: execute failed workflow=%s task=%s: %v", w.ID, t.ID, err)
				now := time.Now()
				t.Status = workflow.TaskStatusFailed
				t.Error = err.Error()
				t.CompletedAt = &now
				if updateErr := s.repo.UpdateTask(goCtx, t); updateErr != nil {
					log.Printf("[WorkflowScheduler] triggerWorkflowsByEvent: update task failed task=%s: %v", t.ID, updateErr)
				}
			}
		}(wfWithNodes, task)
	}
	return nil
}
