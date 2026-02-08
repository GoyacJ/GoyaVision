package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OperatorModel struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TenantID       uuid.UUID      `gorm:"type:uuid;not null;index:idx_operators_tenant_id"`
	OwnerID        uuid.UUID      `gorm:"type:uuid;index:idx_operators_owner_id"`
	Visibility     int            `gorm:"default:0;index:idx_operators_visibility"`
	VisibleRoleIDs datatypes.JSON `gorm:"serializer:json"`
	Code           string         `gorm:"type:varchar(100);not null;uniqueIndex"`
	Name           string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Category    string         `gorm:"type:varchar(50);not null;index:idx_operators_category"`
	Type        string         `gorm:"type:varchar(50);not null;index:idx_operators_type"`
	Origin      string         `gorm:"type:varchar(20);not null;default:'custom';index:idx_operators_origin"`

	ActiveVersionID *uuid.UUID `gorm:"type:uuid;index:idx_operators_active_version_id"`
	Status      string         `gorm:"type:varchar(20);not null;default:'draft';index:idx_operators_status"`
	Tags        datatypes.JSON `gorm:"serializer:json"`
	CreatedAt   time.Time      `gorm:"autoCreateTime;index:idx_operators_created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`

	ActiveVersion *OperatorVersionModel `gorm:"foreignKey:ActiveVersionID"`
}

func (OperatorModel) TableName() string { return "operators" }
