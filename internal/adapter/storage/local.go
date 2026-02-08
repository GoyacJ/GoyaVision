package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"goyavision/config"
	"goyavision/internal/app/port"
)

var _ port.FileStorage = (*LocalAdapter)(nil)

type LocalAdapter struct {
	basePath string
	baseURL  string
}

func NewLocalAdapter(cfg *config.LocalStorage) (*LocalAdapter, error) {
	basePath := strings.TrimRight(cfg.BasePath, string(os.PathSeparator))
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("create base path: %w", err)
	}
	return &LocalAdapter{basePath: basePath, baseURL: strings.TrimRight(cfg.BaseURL, "/")}, nil
}

func (l *LocalAdapter) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	objectName = strings.TrimPrefix(objectName, "/")
	fullPath := filepath.Join(l.basePath, filepath.FromSlash(objectName))
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create dir: %w", err)
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	if err != nil {
		_ = os.Remove(fullPath)
		return "", fmt.Errorf("write file: %w", err)
	}
	return objectName, nil
}

func (l *LocalAdapter) Delete(ctx context.Context, objectName string) error {
	objectName = strings.TrimPrefix(objectName, "/")
	fullPath := filepath.Join(l.basePath, filepath.FromSlash(objectName))
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete file: %w", err)
	}
	return nil
}

func (l *LocalAdapter) GetPublicURL(objectName string) string {
	objectName = strings.TrimPrefix(objectName, "/")
	return l.baseURL + "/" + objectName
}

func (l *LocalAdapter) URLConfig() port.StorageURLConfig {
	return port.StorageURLConfig{
		Endpoint:   "",
		BucketName: "",
		PublicBase: l.baseURL,
		UseSSL:     strings.HasPrefix(l.baseURL, "https://"),
	}
}
