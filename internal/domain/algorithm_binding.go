package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AlgorithmBinding struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey"`
	StreamID        uuid.UUID      `gorm:"type:uuid;not null;index"`
	AlgorithmID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	Enabled         bool           `gorm:"default:true"`
	IntervalSec     int            `gorm:"not null"`
	InitialDelaySec int            `gorm:"default:0"`
	Schedule        datatypes.JSON `gorm:"type:jsonb"`
	Config          datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
}

func (AlgorithmBinding) TableName() string { return "algorithm_bindings" }
