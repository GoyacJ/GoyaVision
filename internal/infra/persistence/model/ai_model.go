package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AIModelModel struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TenantID       uuid.UUID      `gorm:"type:uuid;not null;index:idx_ai_models_tenant_id"`
	OwnerID        uuid.UUID      `gorm:"type:uuid;index:idx_ai_models_owner_id"`
	Visibility     int            `gorm:"default:0;index:idx_ai_models_visibility"`
	VisibleRoleIDs datatypes.JSON `gorm:"type:jsonb"`
	Name           string         `gorm:"type:varchar(100);not null"`
	Description string         `gorm:"type:text"`
	Provider    string         `gorm:"type:varchar(50);not null"`
	Endpoint  string         `gorm:"type:varchar(255)"`
	APIKey    string         `gorm:"type:varchar(255)"`
	ModelName string         `gorm:"type:varchar(100)"`
	Config    datatypes.JSON `gorm:"type:jsonb"`
	Status    string         `gorm:"type:varchar(20);not null;default:'active'"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
}

func (AIModelModel) TableName() string { return "ai_models" }
