package port

import (
	"context"
	"io"
	"time"
)

// ObjectStorage 对象存储接口
//
// 职责：
//  1. 文件上传和下载
//  2. 生成预签名 URL
//  3. 文件元数据管理
//
// 实现：
//  - infra/minio/client.go (MinIO 实现)
//  - 未来可扩展：S3, OSS, COS 等
type ObjectStorage interface {
	// Upload 上传文件
	Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error)

	// Download 下载文件（返回 Reader）
	Download(ctx context.Context, bucket, objectName string) (io.ReadCloser, error)

	// Delete 删除文件
	Delete(ctx context.Context, bucket, objectName string) error

	// GetPresignedURL 生成预签名 URL（用于临时访问）
	GetPresignedURL(ctx context.Context, bucket, objectName string, expires time.Duration) (string, error)

	// Exists 检查文件是否存在
	Exists(ctx context.Context, bucket, objectName string) (bool, error)

	// GetMetadata 获取文件元数据
	GetMetadata(ctx context.Context, bucket, objectName string) (*ObjectMetadata, error)
}

// UploadRequest 上传请求
type UploadRequest struct {
	Bucket      string            // 存储桶名称
	ObjectName  string            // 对象名称（文件路径）
	Reader      io.Reader         // 文件数据流
	Size        int64             // 文件大小（字节）
	ContentType string            // MIME 类型
	Metadata    map[string]string // 自定义元数据
}

// UploadResult 上传结果
type UploadResult struct {
	Bucket     string
	ObjectName string
	ETag       string
	Size       int64
	URL        string // 访问 URL
}

// ObjectMetadata 对象元数据
type ObjectMetadata struct {
	Bucket       string
	ObjectName   string
	Size         int64
	ContentType  string
	ETag         string
	LastModified time.Time
	Metadata     map[string]string
}
