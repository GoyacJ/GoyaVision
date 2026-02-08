package model

import (
	"time"

	"gorm.io/datatypes"
)

// SystemConfigModel GORM model for system_configs table
type SystemConfigModel struct {
	Key         string         `gorm:"primaryKey;type:varchar(255)"`
	Value       datatypes.JSON `gorm:"type:jsonb;not null"`
	Description string         `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (SystemConfigModel) TableName() string {
	return "system_configs"
}
