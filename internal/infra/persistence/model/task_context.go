package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type TaskContextStateModel struct {
	TaskID    uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Version   int64          `gorm:"not null;default:1"`
	Data      datatypes.JSON `gorm:"serializer:json"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
}

func (TaskContextStateModel) TableName() string { return "task_context_state" }

type TaskContextPatchModel struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TaskID        uuid.UUID      `gorm:"type:uuid;not null;index:idx_task_context_patches_task_id"`
	WriterNodeKey string         `gorm:"type:varchar(100);not null;index:idx_task_context_patches_writer_node_key"`
	BeforeVersion int64          `gorm:"not null"`
	AfterVersion  int64          `gorm:"not null;index:idx_task_context_patches_after_version"`
	Diff          datatypes.JSON `gorm:"serializer:json"`
	CreatedAt     time.Time      `gorm:"autoCreateTime;index:idx_task_context_patches_created_at"`
}

func (TaskContextPatchModel) TableName() string { return "task_context_patches" }

type TaskContextSnapshotModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TaskID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_task_context_snapshots_task_id"`
	Version   int64          `gorm:"not null;index:idx_task_context_snapshots_version"`
	Data      datatypes.JSON `gorm:"serializer:json"`
	Trigger   string         `gorm:"type:varchar(50);not null;default:'periodic'"`
	CreatedAt time.Time      `gorm:"autoCreateTime;index:idx_task_context_snapshots_created_at"`
}

func (TaskContextSnapshotModel) TableName() string { return "task_context_snapshots" }
