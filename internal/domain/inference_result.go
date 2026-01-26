package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type InferenceResult struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey"`
	AlgorithmBindingID uuid.UUID      `gorm:"type:uuid;not null;index"`
	StreamID           uuid.UUID      `gorm:"type:uuid;not null;index"`
	Ts                 time.Time      `gorm:"not null;index"`
	FrameRef           string         `gorm:"size:512"`
	Output             datatypes.JSON `gorm:"type:jsonb"`
	LatencyMs          *int
	CreatedAt          time.Time `gorm:"autoCreateTime"`
}

func (InferenceResult) TableName() string { return "inference_results" }
