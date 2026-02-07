package ai_model

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, d *AIModel) error
	Get(ctx context.Context, id uuid.UUID) (*AIModel, error)
	Update(ctx context.Context, d *AIModel) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter Filter) ([]*AIModel, int64, error)
}
