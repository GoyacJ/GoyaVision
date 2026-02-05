package model

import (
	"time"

	"github.com/google/uuid"
)

type MediaSourceModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name          string    `gorm:"type:varchar(255);not null"`
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
