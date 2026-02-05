package dto

import (
	"encoding/json"
	"time"

	"goyavision/internal/domain/workflow"

	"github.com/google/uuid"
)

// ArtifactListQuery 列出产物查询参数
type ArtifactListQuery struct {
	TaskID  *uuid.UUID `query:"task_id"`
	Type    *string    `query:"type"`
	AssetID *uuid.UUID `query:"asset_id"`
	From    *int64     `query:"from"`
	To      *int64     `query:"to"`
	Limit   int        `query:"limit"`
	Offset  int        `query:"offset"`
}

// ArtifactCreateReq 创建产物请求
type ArtifactCreateReq struct {
	TaskID  uuid.UUID              `json:"task_id" validate:"required"`
	Type    string                 `json:"type" validate:"required"`
	AssetID *uuid.UUID             `json:"asset_id,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// ArtifactResponse 产物响应
type ArtifactResponse struct {
	ID        uuid.UUID              `json:"id"`
	TaskID    uuid.UUID              `json:"task_id"`
	Type      string                 `json:"type"`
	AssetID   *uuid.UUID             `json:"asset_id,omitempty"`
	Asset     *AssetResponse         `json:"asset,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// ArtifactListResponse 产物列表响应
type ArtifactListResponse struct {
	Items []*ArtifactResponse `json:"items"`
	Total int64               `json:"total"`
}

// ArtifactToResponse 转换为响应
func ArtifactToResponse(a *workflow.Artifact, minioEndpoint, minioBucket string, minioUseSSL bool) *ArtifactResponse {
	if a == nil {
		return nil
	}

	var data map[string]interface{}
	if a.Data != nil {
		dataBytes, err := json.Marshal(a.Data)
		if err == nil {
			json.Unmarshal(dataBytes, &data)
		}
	}

	resp := &ArtifactResponse{
		ID:        a.ID,
		TaskID:    a.TaskID,
		Type:      string(a.Type),
		AssetID:   a.AssetID,
		Data:      data,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}

	return resp
}

// ArtifactsToResponse 转换为响应列表
func ArtifactsToResponse(artifacts []*workflow.Artifact, minioEndpoint, minioBucket string, minioUseSSL bool) []*ArtifactResponse {
	result := make([]*ArtifactResponse, len(artifacts))
	for i, a := range artifacts {
		result[i] = ArtifactToResponse(a, minioEndpoint, minioBucket, minioUseSSL)
	}
	return result
}
