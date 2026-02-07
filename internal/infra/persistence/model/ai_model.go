package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AIModelModel struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name        string         `gorm:"type:varchar(100);not null"`
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
