package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type FileModel struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name         string         `gorm:"type:varchar(255);not null;index:idx_files_name"`
	OriginalName string         `gorm:"type:varchar(255);not null"`
	Path         string         `gorm:"type:varchar(1024);not null;index:idx_files_path"`
	Size         int64          `gorm:"not null;default:0;index:idx_files_size"`
	MimeType     string         `gorm:"type:varchar(100);not null"`
	Type         string         `gorm:"type:varchar(20);not null;index:idx_files_type"`
	Extension    string         `gorm:"type:varchar(20)"`
	Status       string         `gorm:"type:varchar(20);not null;default:'completed';index:idx_files_status"`
	Hash         string         `gorm:"type:varchar(64);index:idx_files_hash"`
	UploaderID   *uuid.UUID     `gorm:"type:uuid;index:idx_files_uploader"`
	Metadata     datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt    time.Time      `gorm:"autoCreateTime;index:idx_files_created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    *time.Time     `gorm:"index:idx_files_deleted_at"`
}

func (FileModel) TableName() string { return "files" }
