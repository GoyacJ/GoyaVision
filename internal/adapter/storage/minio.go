package storage

import (
	"context"
	"io"
	"strings"

	"goyavision/config"
	"goyavision/internal/app/port"
	"goyavision/pkg/storage"
)

var _ port.FileStorage = (*MinIOAdapter)(nil)

type MinIOAdapter struct {
	client   *storage.MinIOClient
	urlCfg   port.StorageURLConfig
}

func NewMinIOAdapter(cfg *config.MinIO) (*MinIOAdapter, error) {
	client, err := storage.NewMinIOClient(
		cfg.Endpoint,
		cfg.AccessKey,
		cfg.SecretKey,
		cfg.BucketName,
		cfg.UseSSL,
	)
	if err != nil {
		return nil, err
	}
	urlCfg := port.StorageURLConfig{
		Endpoint:   cfg.Endpoint,
		BucketName: cfg.BucketName,
		PublicBase: cfg.PublicBase,
		UseSSL:     cfg.UseSSL,
	}
	return &MinIOAdapter{client: client, urlCfg: urlCfg}, nil
}

func (m *MinIOAdapter) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	return m.client.Upload(ctx, objectName, reader, size, contentType)
}

func (m *MinIOAdapter) Delete(ctx context.Context, objectName string) error {
	return m.client.Delete(ctx, objectName)
}

func (m *MinIOAdapter) GetPublicURL(objectName string) string {
	objectName = strings.TrimPrefix(objectName, "/")
	if m.urlCfg.PublicBase != "" {
		base := strings.TrimRight(m.urlCfg.PublicBase, "/")
		return base + "/" + m.urlCfg.BucketName + "/" + objectName
	}
	return m.client.GetPublicURL(m.urlCfg.Endpoint, m.urlCfg.BucketName, objectName, m.urlCfg.UseSSL)
}

func (m *MinIOAdapter) URLConfig() port.StorageURLConfig {
	return m.urlCfg
}
