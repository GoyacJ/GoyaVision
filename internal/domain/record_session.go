package domain

import (
	"time"

	"github.com/google/uuid"
)

type RecordSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	StreamID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Status    string    `gorm:"not null;size:32"`
	BasePath  string    `gorm:"not null"`
	StartedAt time.Time `gorm:"not null"`
	StoppedAt *time.Time
}

func (RecordSession) TableName() string { return "record_sessions" }

const (
	RecordStatusRunning = "running"
	RecordStatusStopped = "stopped"
)
