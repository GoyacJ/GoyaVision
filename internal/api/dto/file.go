package dto

import (
	"strings"
	"time"

	"goyavision/internal/domain/storage"

	"github.com/google/uuid"
)

// FileListQuery 列出文件查询参数
type FileListQuery struct {
	Type       *string    `query:"type"`
	Status     *string    `query:"status"`
	UploaderID *uuid.UUID `query:"uploader_id"`
	Search     string     `query:"search"`
	From       *int64     `query:"from"`
	To         *int64     `query:"to"`
	Limit      int        `query:"limit"`
	Offset     int        `query:"offset"`
}

// FileCreateReq 创建文件请求（用于手动创建文件记录）
type FileCreateReq struct {
	Name        string                 `json:"name" validate:"required"`
	OriginalName string                `json:"original_name" validate:"required"`
	Path        string                 `json:"path" validate:"required"`
	Size        int64                  `json:"size"`
	MimeType    string                 `json:"mime_type"`
	Type        string                 `json:"type"`
	Extension   string                 `json:"extension"`
	Hash        string                 `json:"hash"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// FileUpdateReq 更新文件请求
type FileUpdateReq struct {
	Name     *string                `json:"name,omitempty"`
	Status   *string                `json:"status,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// FileResponse 文件响应
type FileResponse struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	OriginalName string                `json:"original_name"`
	Path        string                 `json:"path"`
	URL         string                 `json:"url"`
	Size        int64                  `json:"size"`
	MimeType    string                 `json:"mime_type"`
	Type        string                 `json:"type"`
	Extension   string                 `json:"extension"`
	Status      string                 `json:"status"`
	Hash        string                 `json:"hash"`
	UploaderID  *uuid.UUID             `json:"uploader_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// FileListResponse 文件列表响应
type FileListResponse struct {
	Items []*FileResponse `json:"items"`
	Total int64           `json:"total"`
}

// FileToResponse 转换为响应
func FileToResponse(f *storage.File, minioEndpoint, minioBucket string, minioUseSSL bool) *FileResponse {
	if f == nil {
		return nil
	}

	metadata := f.Metadata
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	// 生成完整的文件 URL
	url := f.Path
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		protocol := "http"
		if minioUseSSL {
			protocol = "https"
		}
		url = strings.TrimSuffix(protocol+"://"+minioEndpoint, "/") + "/" + minioBucket + "/" + strings.TrimPrefix(f.Path, "/")
	}

	return &FileResponse{
		ID:          f.ID,
		Name:        f.Name,
		OriginalName: f.OriginalName,
		Path:        f.Path,
		URL:         url,
		Size:        f.Size,
		MimeType:    f.MimeType,
		Type:        string(f.Type),
		Extension:   f.Extension,
		Status:      string(f.Status),
		Hash:        f.Hash,
		UploaderID:  f.UploaderID,
		Metadata:    metadata,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}
}

// FilesToResponse 转换为响应列表
func FilesToResponse(files []*storage.File, minioEndpoint, minioBucket string, minioUseSSL bool) []*FileResponse {
	result := make([]*FileResponse, len(files))
	for i, f := range files {
		result[i] = FileToResponse(f, minioEndpoint, minioBucket, minioUseSSL)
	}
	return result
}
