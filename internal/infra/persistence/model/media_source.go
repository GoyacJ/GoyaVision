package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type MediaSourceModel struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TenantID       uuid.UUID      `gorm:"type:uuid;not null;index:idx_media_sources_tenant_id"`
	OwnerID        uuid.UUID      `gorm:"type:uuid;index:idx_media_sources_owner_id"`
	Visibility     int            `gorm:"default:0;index:idx_media_sources_visibility"`
	VisibleRoleIDs datatypes.JSON `gorm:"serializer:json"`
	Name           string         `gorm:"type:varchar(255);not null"`
	PathName      string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_media_sources_path_name"`
	Type          string    `gorm:"type:varchar(20);not null;index:idx_media_sources_type"`
	URL           string    `gorm:"type:varchar(1024)"`
	Protocol      string    `gorm:"type:varchar(20)"`
	Enabled       bool      `gorm:"not null;default:true"`
	RecordEnabled bool      `gorm:"not null;default:false"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (MediaSourceModel) TableName() string { return "media_sources" }
