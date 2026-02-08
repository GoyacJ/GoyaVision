package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"goyavision/config"
	"goyavision/internal/app/port"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var _ port.FileStorage = (*S3Adapter)(nil)

type S3Adapter struct {
	client *minio.Client
	bucket string
	urlCfg port.StorageURLConfig
}

func NewS3Adapter(cfg *config.S3) (*S3Adapter, error) {
	endpoint := cfg.Endpoint
	if endpoint == "" {
		if cfg.Region != "" {
			endpoint = fmt.Sprintf("s3.%s.amazonaws.com", cfg.Region)
		} else {
			endpoint = "s3.amazonaws.com"
		}
	}
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("create s3 client: %w", err)
	}
	urlCfg := port.StorageURLConfig{
		Endpoint:   endpoint,
		BucketName: cfg.Bucket,
		PublicBase: cfg.PublicBase,
		UseSSL:     cfg.UseSSL,
	}
	return &S3Adapter{client: client, bucket: cfg.Bucket, urlCfg: urlCfg}, nil
}

func (s *S3Adapter) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	objectName = strings.TrimPrefix(objectName, "/")
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("put object: %w", err)
	}
	return objectName, nil
}

func (s *S3Adapter) Delete(ctx context.Context, objectName string) error {
	objectName = strings.TrimPrefix(objectName, "/")
	if err := s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("remove object: %w", err)
	}
	return nil
}

func (s *S3Adapter) GetPublicURL(objectName string) string {
	objectName = strings.TrimPrefix(objectName, "/")
	if s.urlCfg.PublicBase != "" {
		base := strings.TrimRight(s.urlCfg.PublicBase, "/")
		return base + "/" + s.bucket + "/" + objectName
	}
	protocol := "http"
	if s.urlCfg.UseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, s.urlCfg.Endpoint, s.bucket, objectName)
}

func (s *S3Adapter) URLConfig() port.StorageURLConfig {
	return s.urlCfg
}
