package port

import (
	"context"
	"io"
)

// FileStorage 文件存储接口，供应用层上传/删除及生成展示 URL 使用。
// 实现：adapter/storage（MinIO、S3、Local）。
type FileStorage interface {
	Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
	Delete(ctx context.Context, objectName string) error
	GetPublicURL(objectName string) string
}

// StorageURLConfig 用于在 DTO 中构建资产/文件/产物等展示 URL 的配置。
// 由 main 根据当前存储类型注入到 Handler。
type StorageURLConfig struct {
	Endpoint   string
	BucketName string
	PublicBase string
	UseSSL     bool
}
