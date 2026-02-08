package query

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type GetTaskContextHandler struct {
	uow port.UnitOfWork
}

func NewGetTaskContextHandler(uow port.UnitOfWork) *GetTaskContextHandler {
	return &GetTaskContextHandler{uow: uow}
}

func (h *GetTaskContextHandler) Handle(ctx context.Context, query dto.GetTaskContextQuery) (*workflow.TaskContextState, error) {
	var result *workflow.TaskContextState
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.Contexts == nil {
			return apperr.Internal("context repository is not configured", nil)
		}
		state, err := repos.Contexts.GetState(ctx, query.TaskID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("task_context", query.TaskID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get task context")
		}
		result = state
		return nil
	})
	return result, err
}
