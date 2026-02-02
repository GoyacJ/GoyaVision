package dto

import (
	"time"

	"goyavision/internal/domain"

	"github.com/google/uuid"
)

// AssetListQuery 列出资产查询参数
type AssetListQuery struct {
	Type       *string    `query:"type"`
	SourceType *string    `query:"source_type"`
	SourceID   *uuid.UUID `query:"source_id"`
	ParentID   *uuid.UUID `query:"parent_id"`
	Status     *string    `query:"status"`
	Tags       *string    `query:"tags"`
	From       *int64     `query:"from"`
	To         *int64     `query:"to"`
	Limit      int        `query:"limit"`
	Offset     int        `query:"offset"`
}

// AssetCreateReq 创建资产请求
type AssetCreateReq struct {
	Type       string                 `json:"type" validate:"required"`
	SourceType string                 `json:"source_type" validate:"required"`
	SourceID   *uuid.UUID             `json:"source_id,omitempty"`
	ParentID   *uuid.UUID             `json:"parent_id,omitempty"`
	Name       string                 `json:"name" validate:"required"`
	Path       string                 `json:"path" validate:"required"`
	Duration   *float64               `json:"duration,omitempty"`
	Size       int64                  `json:"size"`
	Format     string                 `json:"format,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Status     string                 `json:"status,omitempty"`
	Tags       []string               `json:"tags,omitempty"`
}

// AssetUpdateReq 更新资产请求
type AssetUpdateReq struct {
	Name     *string                `json:"name,omitempty"`
	Status   *string                `json:"status,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Tags     []string               `json:"tags,omitempty"`
}

// AssetResponse 资产响应
type AssetResponse struct {
	ID         uuid.UUID              `json:"id"`
	Type       string                 `json:"type"`
	SourceType string                 `json:"source_type"`
	SourceID   *uuid.UUID             `json:"source_id,omitempty"`
	ParentID   *uuid.UUID             `json:"parent_id,omitempty"`
	Name       string                 `json:"name"`
	Path       string                 `json:"path"`
	Duration   *float64               `json:"duration,omitempty"`
	Size       int64                  `json:"size"`
	Format     string                 `json:"format,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Status     string                 `json:"status"`
	Tags       []string               `json:"tags,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// AssetListResponse 资产列表响应
type AssetListResponse struct {
	Items []*AssetResponse `json:"items"`
	Total int64            `json:"total"`
}

// AssetToResponse 转换为响应
func AssetToResponse(a *domain.MediaAsset) *AssetResponse {
	if a == nil {
		return nil
	}

	var metadata map[string]interface{}
	if a.Metadata != nil {
		if err := a.Metadata.Unmarshal(&metadata); err == nil {
		}
	}

	var tags []string
	if a.Tags != nil {
		if err := a.Tags.Unmarshal(&tags); err == nil {
		}
	}

	return &AssetResponse{
		ID:         a.ID,
		Type:       string(a.Type),
		SourceType: string(a.SourceType),
		SourceID:   a.SourceID,
		ParentID:   a.ParentID,
		Name:       a.Name,
		Path:       a.Path,
		Duration:   a.Duration,
		Size:       a.Size,
		Format:     a.Format,
		Metadata:   metadata,
		Status:     string(a.Status),
		Tags:       tags,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}

// AssetsToResponse 转换为响应列表
func AssetsToResponse(assets []*domain.MediaAsset) []*AssetResponse {
	result := make([]*AssetResponse, len(assets))
	for i, a := range assets {
		result[i] = AssetToResponse(a)
	}
	return result
}
