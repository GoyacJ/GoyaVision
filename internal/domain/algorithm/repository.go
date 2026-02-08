package algorithm

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, a *Algorithm) error
	Get(ctx context.Context, id uuid.UUID) (*Algorithm, error)
	GetByCode(ctx context.Context, code string) (*Algorithm, error)
	GetWithRelations(ctx context.Context, id uuid.UUID) (*Algorithm, error)
	List(ctx context.Context, filter Filter) ([]*Algorithm, int64, error)
	Update(ctx context.Context, a *Algorithm) error
	Delete(ctx context.Context, id uuid.UUID) error

	CreateVersion(ctx context.Context, v *Version) error
	GetVersion(ctx context.Context, id uuid.UUID) (*Version, error)
	GetVersionByName(ctx context.Context, algorithmID uuid.UUID, version string) (*Version, error)
	ListVersions(ctx context.Context, algorithmID uuid.UUID) ([]*Version, error)
	UpdateVersion(ctx context.Context, v *Version) error

	ReplaceImplementations(ctx context.Context, versionID uuid.UUID, impls []Implementation) error
	ListImplementations(ctx context.Context, versionID uuid.UUID) ([]*Implementation, error)

	ReplaceEvaluations(ctx context.Context, versionID uuid.UUID, profiles []EvaluationProfile) error
	ListEvaluations(ctx context.Context, versionID uuid.UUID) ([]*EvaluationProfile, error)
}
