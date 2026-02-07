package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ArtifactModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TenantID  uuid.UUID      `gorm:"type:uuid;not null;index:idx_artifacts_tenant_id"`
	TaskID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_artifacts_task_id"`
	Type      string         `gorm:"type:varchar(50);not null;index:idx_artifacts_type"`
	AssetID   *uuid.UUID     `gorm:"type:uuid;index:idx_artifacts_asset_id"`
	Data      datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt time.Time      `gorm:"autoCreateTime;index:idx_artifacts_created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`

	Task  *TaskModel       `gorm:"foreignKey:TaskID"`
	Asset *MediaAssetModel `gorm:"foreignKey:AssetID"`
}

func (ArtifactModel) TableName() string { return "artifacts" }
