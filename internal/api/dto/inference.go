package dto

import (
	"encoding/json"
	"time"

	"goyavision/internal/domain"

	"github.com/google/uuid"
)

type InferenceResultListQuery struct {
	StreamID  *string `query:"stream_id"`
	BindingID *string `query:"binding_id"`
	From      *int64  `query:"from"`
	To        *int64  `query:"to"`
	Limit     int     `query:"limit"`
	Offset    int     `query:"offset"`
}

type InferenceResultResponse struct {
	ID                 uuid.UUID       `json:"id"`
	AlgorithmBindingID uuid.UUID       `json:"algorithm_binding_id"`
	StreamID           uuid.UUID       `json:"stream_id"`
	Ts                 time.Time       `json:"ts"`
	FrameRef           string          `json:"frame_ref"`
	Output             json.RawMessage `json:"output"`
	LatencyMs          *int            `json:"latency_ms,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
}

type InferenceResultListResponse struct {
	Items []*InferenceResultResponse `json:"items"`
	Total int64                      `json:"total"`
}

func InferenceResultToResponse(r *domain.InferenceResult) *InferenceResultResponse {
	if r == nil {
		return nil
	}
	return &InferenceResultResponse{
		ID:                 r.ID,
		AlgorithmBindingID: r.AlgorithmBindingID,
		StreamID:           r.StreamID,
		Ts:                 r.Ts,
		FrameRef:           r.FrameRef,
		Output:             r.Output,
		LatencyMs:          r.LatencyMs,
		CreatedAt:          r.CreatedAt,
	}
}

func InferenceResultsToResponse(results []*domain.InferenceResult) []*InferenceResultResponse {
	responses := make([]*InferenceResultResponse, len(results))
	for i, r := range results {
		responses[i] = InferenceResultToResponse(r)
	}
	return responses
}
