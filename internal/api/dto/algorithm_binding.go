package dto

import (
	"encoding/json"
	"time"

	"goyavision/internal/domain"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AlgorithmBindingCreateReq struct {
	AlgorithmID     uuid.UUID       `json:"algorithm_id" validate:"required"`
	Enabled         *bool           `json:"enabled"`
	IntervalSec     int             `json:"interval_sec" validate:"required,min=1"`
	InitialDelaySec int             `json:"initial_delay_sec"`
	Schedule        json.RawMessage `json:"schedule"`
	Config          json.RawMessage `json:"config"`
}

type AlgorithmBindingUpdateReq struct {
	AlgorithmID     *uuid.UUID      `json:"algorithm_id"`
	Enabled         *bool           `json:"enabled"`
	IntervalSec     *int            `json:"interval_sec"`
	InitialDelaySec *int            `json:"initial_delay_sec"`
	Schedule        json.RawMessage `json:"schedule"`
	Config          json.RawMessage `json:"config"`
}

type AlgorithmBindingResponse struct {
	ID              uuid.UUID       `json:"id"`
	StreamID        uuid.UUID       `json:"stream_id"`
	AlgorithmID     uuid.UUID       `json:"algorithm_id"`
	Enabled         bool            `json:"enabled"`
	IntervalSec     int             `json:"interval_sec"`
	InitialDelaySec int             `json:"initial_delay_sec"`
	Schedule        json.RawMessage `json:"schedule"`
	Config          json.RawMessage `json:"config"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func AlgorithmBindingToResponse(b *domain.AlgorithmBinding) *AlgorithmBindingResponse {
	if b == nil {
		return nil
	}
	return &AlgorithmBindingResponse{
		ID:              b.ID,
		StreamID:        b.StreamID,
		AlgorithmID:     b.AlgorithmID,
		Enabled:         b.Enabled,
		IntervalSec:     b.IntervalSec,
		InitialDelaySec: b.InitialDelaySec,
		Schedule:        b.Schedule,
		Config:          b.Config,
		CreatedAt:       b.CreatedAt,
		UpdatedAt:       b.UpdatedAt,
	}
}

func AlgorithmBindingsToResponse(bindings []*domain.AlgorithmBinding) []*AlgorithmBindingResponse {
	result := make([]*AlgorithmBindingResponse, len(bindings))
	for i, b := range bindings {
		result[i] = AlgorithmBindingToResponse(b)
	}
	return result
}
