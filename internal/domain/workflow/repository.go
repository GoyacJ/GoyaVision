package workflow

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, w *Workflow) error
	Get(ctx context.Context, id uuid.UUID) (*Workflow, error)
	GetByCode(ctx context.Context, code string) (*Workflow, error)
	GetWithNodes(ctx context.Context, id uuid.UUID) (*Workflow, error)
	List(ctx context.Context, filter Filter) ([]*Workflow, int64, error)
	Update(ctx context.Context, w *Workflow) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListEnabled(ctx context.Context) ([]*Workflow, error)

	CreateNode(ctx context.Context, n *Node) error
	ListNodes(ctx context.Context, workflowID uuid.UUID) ([]*Node, error)
	DeleteNodes(ctx context.Context, workflowID uuid.UUID) error

	CreateEdge(ctx context.Context, e *Edge) error
	ListEdges(ctx context.Context, workflowID uuid.UUID) ([]*Edge, error)
	DeleteEdges(ctx context.Context, workflowID uuid.UUID) error
}

type TaskRepository interface {
	Create(ctx context.Context, t *Task) error
	Get(ctx context.Context, id uuid.UUID) (*Task, error)
	GetWithRelations(ctx context.Context, id uuid.UUID) (*Task, error)
	List(ctx context.Context, filter TaskFilter) ([]*Task, int64, error)
	Update(ctx context.Context, t *Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetStats(ctx context.Context, workflowID *uuid.UUID) (*TaskStats, error)
	ListRunning(ctx context.Context) ([]*Task, error)
}

type ArtifactRepository interface {
	Create(ctx context.Context, a *Artifact) error
	Get(ctx context.Context, id uuid.UUID) (*Artifact, error)
	List(ctx context.Context, filter ArtifactFilter) ([]*Artifact, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByTask(ctx context.Context, taskID uuid.UUID) ([]*Artifact, error)
	ListByType(ctx context.Context, taskID uuid.UUID, artifactType ArtifactType) ([]*Artifact, error)
}
