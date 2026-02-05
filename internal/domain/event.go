package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Type      string
	Timestamp time.Time
	Payload   interface{}
}

type SourceCreated struct {
	SourceID uuid.UUID
	Name     string
	Type     string
}

type SourceDeleted struct {
	SourceID uuid.UUID
}

type AssetCreated struct {
	AssetID    uuid.UUID
	Type       string
	SourceType string
	SourceID   *uuid.UUID
}

type TaskStarted struct {
	TaskID     uuid.UUID
	WorkflowID uuid.UUID
}

type TaskCompleted struct {
	TaskID     uuid.UUID
	WorkflowID uuid.UUID
	Status     string
}

type TaskFailed struct {
	TaskID     uuid.UUID
	WorkflowID uuid.UUID
	Error      string
}
