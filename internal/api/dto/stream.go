package dto

import (
	"time"

	"goyavision/internal/domain"

	"github.com/google/uuid"
)

type StreamListQuery struct {
	Enabled *bool `query:"enabled"`
}

type StreamCreateReq struct {
	URL     string `json:"url" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Enabled *bool  `json:"enabled"`
}

type StreamUpdateReq struct {
	URL     *string `json:"url"`
	Name    *string `json:"name"`
	Enabled *bool   `json:"enabled"`
}

type StreamResponse struct {
	ID        uuid.UUID `json:"id"`
	URL       string    `json:"url"`
	Name      string    `json:"name"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func StreamToResponse(s *domain.Stream) *StreamResponse {
	if s == nil {
		return nil
	}
	return &StreamResponse{
		ID:        s.ID,
		URL:       s.URL,
		Name:      s.Name,
		Enabled:   s.Enabled,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func StreamsToResponse(streams []*domain.Stream) []*StreamResponse {
	result := make([]*StreamResponse, len(streams))
	for i, s := range streams {
		result[i] = StreamToResponse(s)
	}
	return result
}
