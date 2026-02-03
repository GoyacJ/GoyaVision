package dto

import (
	"time"

	"goyavision/internal/domain"

	"github.com/google/uuid"
)

// SourceListQuery 列出媒体源查询参数
type SourceListQuery struct {
	Type   *string `query:"type"`
	Limit  int     `query:"limit"`
	Offset int     `query:"offset"`
}

// SourceCreateReq 创建媒体源请求
type SourceCreateReq struct {
	Name  string `json:"name" validate:"required"`
	Type  string `json:"type" validate:"required"`
	URL   string `json:"url,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Enabled bool  `json:"enabled,omitempty"`
}

// SourceUpdateReq 更新媒体源请求
type SourceUpdateReq struct {
	Name     *string `json:"name,omitempty"`
	URL      *string `json:"url,omitempty"`
	Protocol *string `json:"protocol,omitempty"`
	Enabled  *bool   `json:"enabled,omitempty"`
}

// SourceResponse 媒体源响应
type SourceResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string   `json:"name"`
	PathName      string   `json:"path_name"`
	Type          string   `json:"type"`
	URL           string   `json:"url,omitempty"`
	Protocol      string   `json:"protocol,omitempty"`
	Enabled       bool     `json:"enabled"`
	RecordEnabled bool     `json:"record_enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// SourceListResponse 媒体源列表响应
type SourceListResponse struct {
	Items []*SourceResponse `json:"items"`
	Total int64             `json:"total"`
}

// SourcePreviewResponse 媒体源预览 URL 响应（含各协议及 push 时推流地址）
type SourcePreviewResponse struct {
	PathName  string `json:"path_name"`
	HLSURL    string `json:"hls_url"`
	RTSPURL   string `json:"rtsp_url"`
	RTMPURL   string `json:"rtmp_url"`
	WebRTCURL string `json:"webrtc_url,omitempty"`
	PushURL   string `json:"push_url,omitempty"`
}

// SourceToResponse 转换为响应
func SourceToResponse(s *domain.MediaSource) *SourceResponse {
	if s == nil {
		return nil
	}
	return &SourceResponse{
		ID:            s.ID,
		Name:          s.Name,
		PathName:      s.PathName,
		Type:          string(s.Type),
		URL:           s.URL,
		Protocol:      s.Protocol,
		Enabled:       s.Enabled,
		RecordEnabled: s.RecordEnabled,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

// SourcesToResponse 转换为响应列表
func SourcesToResponse(list []*domain.MediaSource) []*SourceResponse {
	out := make([]*SourceResponse, len(list))
	for i, s := range list {
		out[i] = SourceToResponse(s)
	}
	return out
}
