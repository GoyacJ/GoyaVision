package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"
)

type ListRunningTasksHandler struct {
	uow port.UnitOfWork
}

func NewListRunningTasksHandler(uow port.UnitOfWork) *ListRunningTasksHandler {
	return &ListRunningTasksHandler{uow: uow}
}

func (h *ListRunningTasksHandler) Handle(ctx context.Context, query dto.ListRunningTasksQuery) ([]*workflow.Task, error) {
	var result []*workflow.Task
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		tasks, err := repos.Tasks.ListRunning(ctx)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to list running tasks")
		}
		result = tasks
		return nil
	})

	return result, err
}
