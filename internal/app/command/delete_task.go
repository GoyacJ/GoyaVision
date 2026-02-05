package command

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type DeleteTaskHandler struct {
	uow port.UnitOfWork
}

func NewDeleteTaskHandler(uow port.UnitOfWork) *DeleteTaskHandler {
	return &DeleteTaskHandler{uow: uow}
}

func (h *DeleteTaskHandler) Handle(ctx context.Context, cmd dto.DeleteTaskCommand) error {
	return h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		task, err := repos.Tasks.Get(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("task", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get task")
		}

		if task.IsRunning() {
			return apperr.InvalidInput("cannot delete running task")
		}

		if err := repos.Tasks.Delete(ctx, cmd.ID); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to delete task")
		}

		return nil
	})
}
