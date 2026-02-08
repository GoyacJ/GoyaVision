package workflow

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusSuccess   TaskStatus = "success"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
)

type NodeExecutionStatus string

const (
	NodeExecPending NodeExecutionStatus = "pending"
	NodeExecRunning NodeExecutionStatus = "running"
	NodeExecSuccess NodeExecutionStatus = "success"
	NodeExecFailed  NodeExecutionStatus = "failed"
	NodeExecSkipped NodeExecutionStatus = "skipped"
)

type NodeExecution struct {
	NodeKey     string              `json:"node_key"`
	Status      NodeExecutionStatus `json:"status"`
	Error       string              `json:"error,omitempty"`
	StartedAt   *time.Time          `json:"started_at,omitempty"`
	CompletedAt *time.Time          `json:"completed_at,omitempty"`
	ArtifactIDs []uuid.UUID         `json:"artifact_ids,omitempty"`
}

type Task struct {
	ID                 uuid.UUID
	TenantID           uuid.UUID
	TriggeredByUserID  *uuid.UUID
	WorkflowID         uuid.UUID
	WorkflowRevisionID *uuid.UUID
	WorkflowRevision   int64
	AssetID            *uuid.UUID
	Status             TaskStatus
	Progress           int
	CurrentNode        string
	ContextVersion     int64
	InputParams        map[string]interface{}
	Error              string
	NodeExecutions     []NodeExecution
	StartedAt          *time.Time
	CompletedAt        *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (t *Task) IsPending() bool {
	return t.Status == TaskStatusPending
}

func (t *Task) IsRunning() bool {
	return t.Status == TaskStatusRunning
}

func (t *Task) IsSuccess() bool {
	return t.Status == TaskStatusSuccess
}

func (t *Task) IsFailed() bool {
	return t.Status == TaskStatusFailed
}

func (t *Task) IsCancelled() bool {
	return t.Status == TaskStatusCancelled
}

func (t *Task) IsCompleted() bool {
	return t.Status == TaskStatusSuccess || t.Status == TaskStatusFailed || t.Status == TaskStatusCancelled
}

func (t *Task) Start() {
	now := time.Now()
	t.Status = TaskStatusRunning
	t.StartedAt = &now
}

func (t *Task) Complete() {
	now := time.Now()
	t.Status = TaskStatusSuccess
	t.Progress = 100
	t.CompletedAt = &now
}

func (t *Task) Fail(errMsg string) {
	now := time.Now()
	t.Status = TaskStatusFailed
	t.Error = errMsg
	t.CompletedAt = &now
}

func (t *Task) Cancel() {
	now := time.Now()
	t.Status = TaskStatusCancelled
	t.CompletedAt = &now
}

func (t *Task) Duration() float64 {
	if t.StartedAt == nil {
		return 0
	}
	endTime := time.Now()
	if t.CompletedAt != nil {
		endTime = *t.CompletedAt
	}
	return endTime.Sub(*t.StartedAt).Seconds()
}

type TaskFilter struct {
	WorkflowID        *uuid.UUID
	AssetID           *uuid.UUID
	TriggeredByUserID *uuid.UUID
	Status            *TaskStatus
	From              *time.Time
	To                *time.Time
	Limit             int
	Offset            int
}

type TaskStats struct {
	Total     int64
	Pending   int64
	Running   int64
	Success   int64
	Failed    int64
	Cancelled int64
}
