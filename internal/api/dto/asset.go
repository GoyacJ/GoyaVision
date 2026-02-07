package dto

import (
	"strings"
	"time"

	"goyavision/internal/domain/media"

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
	Path       string                 `json:"path"`
	Duration   *float64               `json:"duration,omitempty"`
	Size       int64                  `json:"size"`
	Format     string                 `json:"format,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	Status         string                 `json:"status,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Visibility     *int                   `json:"visibility,omitempty"`
	VisibleRoleIDs []string               `json:"visible_role_ids,omitempty"`
}

// AssetUpdateReq 更新资产请求
type AssetUpdateReq struct {
	Name           *string                `json:"name,omitempty"`
	Status         *string                `json:"status,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Visibility     *int                   `json:"visibility,omitempty"`
	VisibleRoleIDs []string               `json:"visible_role_ids,omitempty"`
}

// AssetResponse 资产响应
type AssetResponse struct {
	ID             uuid.UUID              `json:"id"`
	Type           string                 `json:"type"`
	SourceType     string                 `json:"source_type"`
	SourceID       *uuid.UUID             `json:"source_id,omitempty"`
	ParentID       *uuid.UUID             `json:"parent_id,omitempty"`
	Name           string                 `json:"name"`
	Path           string                 `json:"path"`
	Duration       *float64               `json:"duration,omitempty"`
	Size           int64                  `json:"size"`
	Format         string                 `json:"format,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	Status         string                 `json:"status"`
	Tags           []string               `json:"tags,omitempty"`
	Visibility     int                    `json:"visibility"`
	VisibleRoleIDs []string               `json:"visible_role_ids,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// AssetListResponse 资产列表响应
type AssetListResponse struct {
	Items []*AssetResponse `json:"items"`
	Total int64            `json:"total"`
}

// AssetToResponse 转换为响应
// minioEndpoint: MinIO 端点地址（如 "39.105.2.5:14250"）
// minioBucket: MinIO 存储桶名称
// minioPublicBase: MinIO 公共访问基址（如 "https://vision.ysmjjsy.com/minio"）
// minioUseSSL: 是否使用 SSL
func AssetToResponse(a *media.Asset, minioEndpoint, minioBucket, minioPublicBase string, minioUseSSL bool) *AssetResponse {
	if a == nil {
		return nil
	}

	metadata := a.Metadata
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	tags := a.Tags
	if tags == nil {
		tags = []string{}
	}

	// 生成完整的文件 URL
	path := a.Path
	if a.SourceType == media.AssetSourceUpload || a.SourceType == media.AssetSourceGenerated ||
		a.SourceType == media.AssetSourceOperatorOutput {
		// 如果是上传或生成的文件，且 path 不是完整 URL，则生成 MinIO 完整 URL
		if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
			if minioPublicBase != "" {
				base := strings.TrimRight(minioPublicBase, "/")
				path = base + "/" + minioBucket + "/" + strings.TrimPrefix(path, "/")
			} else {
				protocol := "http"
				if minioUseSSL {
					protocol = "https"
				}
				path = strings.TrimSuffix(protocol+"://"+minioEndpoint, "/") + "/" + minioBucket + "/" + strings.TrimPrefix(path, "/")
			}
		}
	}

	return &AssetResponse{
		ID:             a.ID,
		Type:           string(a.Type),
		SourceType:     string(a.SourceType),
		SourceID:       a.SourceID,
		ParentID:       a.ParentID,
		Name:           a.Name,
		Path:           path,
		Duration:       a.Duration,
		Size:           a.Size,
		Format:         a.Format,
		Metadata:       metadata,
		Status:         string(a.Status),
		Tags:           tags,
		Visibility:     int(a.Visibility),
		VisibleRoleIDs: a.VisibleRoleIDs,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
}

// AssetsToResponse 转换为响应列表
func AssetsToResponse(assets []*media.Asset, minioEndpoint, minioBucket, minioPublicBase string, minioUseSSL bool) []*AssetResponse {
	result := make([]*AssetResponse, len(assets))
	for i, a := range assets {
		result[i] = AssetToResponse(a, minioEndpoint, minioBucket, minioPublicBase, minioUseSSL)
	}
	return result
}
