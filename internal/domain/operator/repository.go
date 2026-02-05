package operator

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, o *Operator) error
	Get(ctx context.Context, id uuid.UUID) (*Operator, error)
	GetByCode(ctx context.Context, code string) (*Operator, error)
	List(ctx context.Context, filter Filter) ([]*Operator, int64, error)
	Update(ctx context.Context, o *Operator) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListEnabled(ctx context.Context) ([]*Operator, error)
	ListByCategory(ctx context.Context, category Category) ([]*Operator, error)
}
