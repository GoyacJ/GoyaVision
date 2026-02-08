package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type TaskModel struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TenantID          uuid.UUID      `gorm:"type:uuid;not null;index:idx_tasks_tenant_id"`
	TriggeredByUserID *uuid.UUID     `gorm:"type:uuid;index:idx_tasks_triggered_by"`
	WorkflowID        uuid.UUID      `gorm:"type:uuid;not null;index:idx_tasks_workflow_id"`
	AssetID           *uuid.UUID     `gorm:"type:uuid;index:idx_tasks_asset_id"`
	Status            string         `gorm:"type:varchar(20);not null;default:'pending';index:idx_tasks_status"`
	Progress          int            `gorm:"not null;default:0"`
	CurrentNode       string         `gorm:"type:varchar(100)"`
	InputParams       datatypes.JSON `gorm:"serializer:json"`
	Error             string         `gorm:"type:text"`
	NodeExecutions    datatypes.JSON `gorm:"serializer:json"`
	StartedAt         *time.Time
	CompletedAt       *time.Time
	CreatedAt         time.Time `gorm:"autoCreateTime;index:idx_tasks_created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`

	Workflow  *WorkflowModel   `gorm:"foreignKey:WorkflowID"`
	Asset     *MediaAssetModel `gorm:"foreignKey:AssetID"`
	Artifacts []ArtifactModel  `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE"`
}

func (TaskModel) TableName() string { return "tasks" }
