package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type MediaAssetModel struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TenantID       uuid.UUID      `gorm:"type:uuid;not null;index:idx_assets_tenant_id"`
	OwnerID        uuid.UUID      `gorm:"type:uuid;index:idx_assets_owner_id"`
	Visibility     int            `gorm:"default:0;index:idx_assets_visibility"` // 0: Private, 1: Role, 2: Public
	VisibleRoleIDs datatypes.JSON `gorm:"serializer:json"`
	Type           string         `gorm:"type:varchar(20);not null;index:idx_assets_type"`
	SourceType     string         `gorm:"type:varchar(20);not null;index:idx_assets_source_type"`
	SourceID       *uuid.UUID     `gorm:"type:uuid;index:idx_assets_source_id"`
	ParentID       *uuid.UUID     `gorm:"type:uuid;index:idx_assets_parent_id"`
	Name           string         `gorm:"type:varchar(255);not null"`
	Path           string         `gorm:"type:varchar(1024);not null"`
	Duration       *float64       `gorm:"type:float8"`
	Size           int64          `gorm:"not null;default:0"`
	Format         string         `gorm:"type:varchar(50)"`
	Metadata       datatypes.JSON `gorm:"serializer:json"`
	Status         string         `gorm:"type:varchar(20);not null;default:'pending';index:idx_assets_status"`
	Tags           datatypes.JSON `gorm:"serializer:json"`
	CreatedAt      time.Time      `gorm:"autoCreateTime;index:idx_assets_created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
}

func (MediaAssetModel) TableName() string { return "media_assets" }
