package command

import (
	"context"
	"errors"
	"time"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type CancelTaskHandler struct {
	uow port.UnitOfWork
}

func NewCancelTaskHandler(uow port.UnitOfWork) *CancelTaskHandler {
	return &CancelTaskHandler{uow: uow}
}

func (h *CancelTaskHandler) Handle(ctx context.Context, cmd dto.CancelTaskCommand) (*workflow.Task, error) {
	var result *workflow.Task
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		task, err := repos.Tasks.Get(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("task", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get task")
		}

		if task.IsCompleted() {
			return apperr.InvalidInput("task is already completed")
		}

		now := time.Now()
		task.Status = workflow.TaskStatusCancelled
		if task.CompletedAt == nil {
			task.CompletedAt = &now
		}

		if err := repos.Tasks.Update(ctx, task); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to cancel task")
		}

		result = task
		return nil
	})

	return result, err
}
