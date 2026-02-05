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

type GetTaskHandler struct {
	uow port.UnitOfWork
}

func NewGetTaskHandler(uow port.UnitOfWork) *GetTaskHandler {
	return &GetTaskHandler{uow: uow}
}

func (h *GetTaskHandler) Handle(ctx context.Context, query dto.GetTaskQuery) (*workflow.Task, error) {
	var result *workflow.Task
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		task, err := repos.Tasks.Get(ctx, query.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("task", query.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get task")
		}
		result = task
		return nil
	})

	return result, err
}
