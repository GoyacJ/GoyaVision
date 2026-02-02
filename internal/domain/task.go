package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusSuccess   TaskStatus = "success"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
)

// Task 任务实体（工作流执行实例）
type Task struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	WorkflowID     uuid.UUID      `gorm:"type:uuid;not null;index:idx_tasks_workflow_id"`
	AssetID        *uuid.UUID     `gorm:"type:uuid;index:idx_tasks_asset_id"`
	Status         TaskStatus     `gorm:"type:varchar(20);not null;default:'pending';index:idx_tasks_status"`
	Progress       int            `gorm:"not null;default:0"`
	CurrentNode    string         `gorm:"type:varchar(100)"`
	InputParams    datatypes.JSON `gorm:"type:jsonb"`
	Error          string         `gorm:"type:text"`
	StartedAt      *time.Time
	CompletedAt    *time.Time
	CreatedAt      time.Time      `gorm:"autoCreateTime;index:idx_tasks_created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`

	Workflow  *Workflow     `gorm:"foreignKey:WorkflowID"`
	Asset     *MediaAsset   `gorm:"foreignKey:AssetID"`
	Artifacts []Artifact    `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE"`
}

func (Task) TableName() string { return "tasks" }

// IsPending 判断任务是否待执行
func (t *Task) IsPending() bool {
	return t.Status == TaskStatusPending
}

// IsRunning 判断任务是否运行中
func (t *Task) IsRunning() bool {
	return t.Status == TaskStatusRunning
}

// IsSuccess 判断任务是否成功
func (t *Task) IsSuccess() bool {
	return t.Status == TaskStatusSuccess
}

// IsFailed 判断任务是否失败
func (t *Task) IsFailed() bool {
	return t.Status == TaskStatusFailed
}

// IsCancelled 判断任务是否已取消
func (t *Task) IsCancelled() bool {
	return t.Status == TaskStatusCancelled
}

// IsCompleted 判断任务是否已完成（成功、失败或取消）
func (t *Task) IsCompleted() bool {
	return t.Status == TaskStatusSuccess || t.Status == TaskStatusFailed || t.Status == TaskStatusCancelled
}

// Duration 计算任务执行时长（秒）
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

// TaskFilter 任务过滤器
type TaskFilter struct {
	WorkflowID *uuid.UUID
	AssetID    *uuid.UUID
	Status     *TaskStatus
	From       *time.Time
	To         *time.Time
	Limit      int
	Offset     int
}

// TaskStats 任务统计
type TaskStats struct {
	Total     int64
	Pending   int64
	Running   int64
	Success   int64
	Failed    int64
	Cancelled int64
}

// TaskInputParams 任务输入参数
type TaskInputParams struct {
	AssetID uuid.UUID              `json:"asset_id"`
	Params  map[string]interface{} `json:"params,omitempty"`
}
