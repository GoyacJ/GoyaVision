package media

import (
	"context"

	"github.com/google/uuid"
)

type SourceRepository interface {
	Create(ctx context.Context, s *Source) error
	Get(ctx context.Context, id uuid.UUID) (*Source, error)
	GetByPathName(ctx context.Context, pathName string) (*Source, error)
	List(ctx context.Context, filter SourceFilter) ([]*Source, int64, error)
	Update(ctx context.Context, s *Source) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type AssetRepository interface {
	Create(ctx context.Context, a *Asset) error
	Get(ctx context.Context, id uuid.UUID) (*Asset, error)
	List(ctx context.Context, filter AssetFilter) ([]*Asset, int64, error)
	Update(ctx context.Context, a *Asset) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListBySource(ctx context.Context, sourceID uuid.UUID) ([]*Asset, error)
	ListByParent(ctx context.Context, parentID uuid.UUID) ([]*Asset, error)
	GetAllTags(ctx context.Context) ([]string, error)
}
