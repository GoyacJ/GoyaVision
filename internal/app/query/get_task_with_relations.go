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

type GetTaskWithRelationsHandler struct {
	uow port.UnitOfWork
}

func NewGetTaskWithRelationsHandler(uow port.UnitOfWork) *GetTaskWithRelationsHandler {
	return &GetTaskWithRelationsHandler{uow: uow}
}

func (h *GetTaskWithRelationsHandler) Handle(ctx context.Context, query dto.GetTaskWithRelationsQuery) (*workflow.Task, error) {
	var result *workflow.Task
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		task, err := repos.Tasks.GetWithRelations(ctx, query.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("task", query.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get task with relations")
		}
		result = task
		return nil
	})

	return result, err
}
