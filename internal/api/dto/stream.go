package dto

import (
	"time"

	"goyavision/internal/domain"

	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type StreamListQuery struct {
	Enabled    *bool `query:"enabled"`
	WithStatus bool  `query:"with_status"`
}

type StreamCreateReq struct {
	URL     string `json:"url"`
	Name    string `json:"name" validate:"required"`
	Type    string `json:"type"`
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
	Type      string    `json:"type"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StreamStatusResponse struct {
	PathName      string   `json:"path_name"`
	Ready         bool     `json:"ready"`
	Online        bool     `json:"online"`
	Tracks        []string `json:"tracks"`
	BytesReceived uint64   `json:"bytes_received"`
	BytesSent     uint64   `json:"bytes_sent"`
	ReaderCount   int      `json:"reader_count"`
	RTSPUrl       string   `json:"rtsp_url"`
	RTMPUrl       string   `json:"rtmp_url"`
	HLSUrl        string   `json:"hls_url"`
	WebRTCUrl     string   `json:"webrtc_url"`
}

type StreamWithStatusResponse struct {
	ID        uuid.UUID             `json:"id"`
	URL       string                `json:"url"`
	Name      string                `json:"name"`
	Type      string                `json:"type"`
	Enabled   bool                  `json:"enabled"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	Status    *StreamStatusResponse `json:"status,omitempty"`
}

func StreamToResponse(s *domain.Stream) *StreamResponse {
	if s == nil {
		return nil
	}
	return &StreamResponse{
		ID:        s.ID,
		URL:       s.URL,
		Name:      s.Name,
		Type:      string(s.Type),
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

func StreamStatusToResponse(s *domain.StreamStatus) *StreamStatusResponse {
	if s == nil {
		return nil
	}
	return &StreamStatusResponse{
		PathName:      s.PathName,
		Ready:         s.Ready,
		Online:        s.Online,
		Tracks:        s.Tracks,
		BytesReceived: s.BytesReceived,
		BytesSent:     s.BytesSent,
		ReaderCount:   s.ReaderCount,
		RTSPUrl:       s.RTSPUrl,
		RTMPUrl:       s.RTMPUrl,
		HLSUrl:        s.HLSUrl,
		WebRTCUrl:     s.WebRTCUrl,
	}
}

func StreamWithStatusToResponse(s *domain.StreamWithStatus) *StreamWithStatusResponse {
	if s == nil {
		return nil
	}
	return &StreamWithStatusResponse{
		ID:        s.Stream.ID,
		URL:       s.Stream.URL,
		Name:      s.Stream.Name,
		Type:      string(s.Stream.Type),
		Enabled:   s.Stream.Enabled,
		CreatedAt: s.Stream.CreatedAt,
		UpdatedAt: s.Stream.UpdatedAt,
		Status:    StreamStatusToResponse(s.Status),
	}
}

func StreamsWithStatusToResponse(streams []*domain.StreamWithStatus) []*StreamWithStatusResponse {
	result := make([]*StreamWithStatusResponse, len(streams))
	for i, s := range streams {
		result[i] = StreamWithStatusToResponse(s)
	}
	return result
}
