package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Algorithm struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name       string         `gorm:"not null"`
	Endpoint   string         `gorm:"not null"`
	InputSpec  datatypes.JSON `gorm:"type:jsonb"`
	OutputSpec datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
}

func (Algorithm) TableName() string { return "algorithms" }
