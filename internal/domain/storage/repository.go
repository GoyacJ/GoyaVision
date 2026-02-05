package storage

import (
	"context"

	"github.com/google/uuid"
)

type FileRepository interface {
	Create(ctx context.Context, f *File) error
	Get(ctx context.Context, id uuid.UUID) (*File, error)
	GetByPath(ctx context.Context, path string) (*File, error)
	GetByHash(ctx context.Context, hash string) (*File, error)
	List(ctx context.Context, filter FileFilter) ([]*File, int64, error)
	Update(ctx context.Context, f *File) error
	Delete(ctx context.Context, id uuid.UUID) error
}
