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

type FailTaskHandler struct {
	uow port.UnitOfWork
}

func NewFailTaskHandler(uow port.UnitOfWork) *FailTaskHandler {
	return &FailTaskHandler{uow: uow}
}

func (h *FailTaskHandler) Handle(ctx context.Context, cmd dto.FailTaskCommand) (*workflow.Task, error) {
	var result *workflow.Task
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		task, err := repos.Tasks.Get(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("task", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get task")
		}

		if !task.IsRunning() {
			return apperr.InvalidInput("task is not running")
		}

		now := time.Now()
		task.Status = workflow.TaskStatusFailed
		task.CompletedAt = &now
		task.Error = cmd.ErrorMsg

		if err := repos.Tasks.Update(ctx, task); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to mark task as failed")
		}

		result = task
		return nil
	})

	return result, err
}
