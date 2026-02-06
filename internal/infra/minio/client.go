package minio

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"goyavision/config"
	"goyavision/internal/app/port"
	"goyavision/pkg/apperr"
)

// Client MinIO 对象存储客户端实现
type Client struct {
	client     *minio.Client
	bucketName string
	endpoint   string
	useSSL     bool
}

// NewClient 创建 MinIO 客户端实例
func NewClient(cfg *config.MinIO) (*Client, error) {
	if cfg.Endpoint == "" {
		return nil, apperr.InvalidInput("MinIO endpoint is required")
	}
	if cfg.AccessKey == "" {
		return nil, apperr.InvalidInput("MinIO access key is required")
	}
	if cfg.SecretKey == "" {
		return nil, apperr.InvalidInput("MinIO secret key is required")
	}
	if cfg.BucketName == "" {
		return nil, apperr.InvalidInput("MinIO bucket name is required")
	}

	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeInternal, "failed to initialize MinIO client")
	}

	client := &Client{
		client:     minioClient,
		bucketName: cfg.BucketName,
		endpoint:   cfg.Endpoint,
		useSSL:     cfg.UseSSL,
	}

	if err := client.ensureBucket(context.Background()); err != nil {
		return nil, err
	}

	return client, nil
}

// ensureBucket 确保存储桶存在，如不存在则创建
func (c *Client) ensureBucket(ctx context.Context) error {
	exists, err := c.client.BucketExists(ctx, c.bucketName)
	if err != nil {
		return apperr.Wrap(err, apperr.CodeInternal, "failed to check bucket existence")
	}

	if !exists {
		err = c.client.MakeBucket(ctx, c.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return apperr.Wrap(err, apperr.CodeInternal, "failed to create bucket")
		}
	}

	return nil
}

// Upload 上传文件到对象存储
func (c *Client) Upload(ctx context.Context, req *port.UploadRequest) (*port.UploadResult, error) {
	if req == nil {
		return nil, apperr.InvalidInput("upload request is required")
	}
	if req.Bucket == "" {
		req.Bucket = c.bucketName
	}
	if req.ObjectName == "" {
		return nil, apperr.InvalidInput("object name is required")
	}
	if req.Reader == nil {
		return nil, apperr.InvalidInput("reader is required")
	}

	opts := minio.PutObjectOptions{
		ContentType:  req.ContentType,
		UserMetadata: req.Metadata,
	}

	info, err := c.client.PutObject(ctx, req.Bucket, req.ObjectName, req.Reader, req.Size, opts)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeInternal, "failed to upload object")
	}

	objectURL := c.buildObjectURL(req.Bucket, req.ObjectName)

	return &port.UploadResult{
		Bucket:     req.Bucket,
		ObjectName: req.ObjectName,
		ETag:       info.ETag,
		Size:       info.Size,
		URL:        objectURL,
	}, nil
}

// Download 下载文件
func (c *Client) Download(ctx context.Context, bucket, objectName string) (io.ReadCloser, error) {
	if bucket == "" {
		bucket = c.bucketName
	}
	if objectName == "" {
		return nil, apperr.InvalidInput("object name is required")
	}

	object, err := c.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeInternal, "failed to download object")
	}

	_, err = object.Stat()
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return nil, apperr.NotFound("object", objectName)
		}
		return nil, apperr.Wrap(err, apperr.CodeInternal, "failed to stat object")
	}

	return object, nil
}

// Delete 删除文件
func (c *Client) Delete(ctx context.Context, bucket, objectName string) error {
	if bucket == "" {
		bucket = c.bucketName
	}
	if objectName == "" {
		return apperr.InvalidInput("object name is required")
	}

	err := c.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return apperr.Wrap(err, apperr.CodeInternal, "failed to delete object")
	}

	return nil
}

// GetPresignedURL 生成预签名 URL（用于临时访问）
func (c *Client) GetPresignedURL(ctx context.Context, bucket, objectName string, expires time.Duration) (string, error) {
	if bucket == "" {
		bucket = c.bucketName
	}
	if objectName == "" {
		return "", apperr.InvalidInput("object name is required")
	}
	if expires <= 0 {
		expires = 15 * time.Minute
	}

	presignedURL, err := c.client.PresignedGetObject(ctx, bucket, objectName, expires, nil)
	if err != nil {
		return "", apperr.Wrap(err, apperr.CodeInternal, "failed to generate presigned URL")
	}

	return presignedURL.String(), nil
}

// Exists 检查文件是否存在
func (c *Client) Exists(ctx context.Context, bucket, objectName string) (bool, error) {
	if bucket == "" {
		bucket = c.bucketName
	}
	if objectName == "" {
		return false, apperr.InvalidInput("object name is required")
	}

	_, err := c.client.StatObject(ctx, bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, apperr.Wrap(err, apperr.CodeInternal, "failed to check object existence")
	}

	return true, nil
}

// GetMetadata 获取文件元数据
func (c *Client) GetMetadata(ctx context.Context, bucket, objectName string) (*port.ObjectMetadata, error) {
	if bucket == "" {
		bucket = c.bucketName
	}
	if objectName == "" {
		return nil, apperr.InvalidInput("object name is required")
	}

	info, err := c.client.StatObject(ctx, bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return nil, apperr.NotFound("object", objectName)
		}
		return nil, apperr.Wrap(err, apperr.CodeInternal, "failed to get object metadata")
	}

	return &port.ObjectMetadata{
		Bucket:       bucket,
		ObjectName:   objectName,
		Size:         info.Size,
		ContentType:  info.ContentType,
		ETag:         info.ETag,
		LastModified: info.LastModified,
		Metadata:     info.UserMetadata,
	}, nil
}

// buildObjectURL 构建对象访问 URL
func (c *Client) buildObjectURL(bucket, objectName string) string {
	scheme := "http"
	if c.useSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, c.endpoint, bucket, url.PathEscape(objectName))
}
