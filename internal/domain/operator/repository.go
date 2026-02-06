package operator

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, o *Operator) error
	Get(ctx context.Context, id uuid.UUID) (*Operator, error)
	GetWithActiveVersion(ctx context.Context, id uuid.UUID) (*Operator, error)
	GetByCode(ctx context.Context, code string) (*Operator, error)
	List(ctx context.Context, filter Filter) ([]*Operator, int64, error)
	Update(ctx context.Context, o *Operator) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListPublished(ctx context.Context) ([]*Operator, error)
	ListByCategory(ctx context.Context, category Category) ([]*Operator, error)
}

type VersionRepository interface {
	Create(ctx context.Context, version *OperatorVersion) error
	Get(ctx context.Context, id uuid.UUID) (*OperatorVersion, error)
	ListByOperator(ctx context.Context, operatorID uuid.UUID) ([]*OperatorVersion, error)
	GetByOperatorAndVersion(ctx context.Context, operatorID uuid.UUID, version string) (*OperatorVersion, error)
	Update(ctx context.Context, version *OperatorVersion) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TemplateFilter struct {
	Category *Category
	Type     *Type
	ExecMode *ExecMode
	Keyword  string
	Tags     []string
	Limit    int
	Offset   int
}

type TemplateRepository interface {
	Create(ctx context.Context, tpl *OperatorTemplate) error
	Get(ctx context.Context, id uuid.UUID) (*OperatorTemplate, error)
	GetByCode(ctx context.Context, code string) (*OperatorTemplate, error)
	List(ctx context.Context, filter TemplateFilter) ([]*OperatorTemplate, int64, error)
	Update(ctx context.Context, tpl *OperatorTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementDownloads(ctx context.Context, id uuid.UUID) error
}

type DependencyRepository interface {
	Create(ctx context.Context, dep *OperatorDependency) error
	ListByOperator(ctx context.Context, operatorID uuid.UUID) ([]*OperatorDependency, error)
	DeleteByOperator(ctx context.Context, operatorID uuid.UUID) error
	CheckDependenciesSatisfied(ctx context.Context, operatorID uuid.UUID) (bool, []string, error)
}
