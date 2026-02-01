package domain

import (
	"time"

	"github.com/google/uuid"
)

// Permission 权限实体（API 资源）
type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Code        string    `gorm:"uniqueIndex;not null;size:64"`
	Name        string    `gorm:"not null;size:64"`
	Method      string    `gorm:"size:16"`
	Path        string    `gorm:"size:256"`
	Description string    `gorm:"size:256"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (Permission) TableName() string {
	return "permissions"
}
