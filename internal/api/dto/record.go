package dto

import (
	"time"

	"goyavision/internal/domain"

	"github.com/google/uuid"
)

type RecordSessionResponse struct {
	ID        uuid.UUID  `json:"id"`
	StreamID  uuid.UUID  `json:"stream_id"`
	Status    string     `json:"status"`
	BasePath  string     `json:"base_path"`
	StartedAt time.Time  `json:"started_at"`
	StoppedAt *time.Time `json:"stopped_at,omitempty"`
}

type RecordStartResponse struct {
	SessionID uuid.UUID `json:"session_id"`
}

func RecordSessionToResponse(s *domain.RecordSession) *RecordSessionResponse {
	if s == nil {
		return nil
	}
	return &RecordSessionResponse{
		ID:        s.ID,
		StreamID:  s.StreamID,
		Status:    s.Status,
		BasePath:  s.BasePath,
		StartedAt: s.StartedAt,
		StoppedAt: s.StoppedAt,
	}
}

func RecordSessionsToResponse(sessions []*domain.RecordSession) []*RecordSessionResponse {
	result := make([]*RecordSessionResponse, len(sessions))
	for i, s := range sessions {
		result[i] = RecordSessionToResponse(s)
	}
	return result
}
