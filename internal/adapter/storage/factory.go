package storage

import (
	"fmt"

	"goyavision/config"
	"goyavision/internal/app/port"
)

func NewFileStorageFromConfig(cfg *config.Config) (port.FileStorage, port.StorageURLConfig, error) {
	stype := cfg.Storage.Type
	if stype == "" {
		stype = "minio"
	}
	switch stype {
	case "minio":
		a, err := NewMinIOAdapter(&cfg.MinIO)
		if err != nil {
			return nil, port.StorageURLConfig{}, err
		}
		return a, a.URLConfig(), nil
	case "s3":
		a, err := NewS3Adapter(&cfg.Storage.S3)
		if err != nil {
			return nil, port.StorageURLConfig{}, err
		}
		return a, a.URLConfig(), nil
	case "local":
		a, err := NewLocalAdapter(&cfg.Storage.Local)
		if err != nil {
			return nil, port.StorageURLConfig{}, err
		}
		return a, a.URLConfig(), nil
	default:
		return nil, port.StorageURLConfig{}, fmt.Errorf("unsupported storage type: %s", stype)
	}
}
