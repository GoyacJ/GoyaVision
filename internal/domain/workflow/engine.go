package workflow

import (
	"context"

	"goyavision/internal/domain/operator"

	"github.com/google/uuid"
)

type Engine interface {
	Execute(ctx context.Context, workflow *Workflow, task *Task) error
	Cancel(ctx context.Context, taskID uuid.UUID) error
	GetProgress(ctx context.Context, taskID uuid.UUID) (int, error)
}

type OperatorExecutor interface {
	Execute(ctx context.Context, op *operator.Operator, input *operator.Input) (*operator.Output, error)
}
