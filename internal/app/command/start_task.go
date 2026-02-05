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

type StartTaskHandler struct {
	uow port.UnitOfWork
}

func NewStartTaskHandler(uow port.UnitOfWork) *StartTaskHandler {
	return &StartTaskHandler{uow: uow}
}

func (h *StartTaskHandler) Handle(ctx context.Context, cmd dto.StartTaskCommand) (*workflow.Task, error) {
	var result *workflow.Task
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		task, err := repos.Tasks.Get(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("task", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get task")
		}

		if !task.IsPending() {
			return apperr.InvalidInput("task is not pending")
		}

		now := time.Now()
		task.Status = workflow.TaskStatusRunning
		task.StartedAt = &now
		task.Progress = 0

		if err := repos.Tasks.Update(ctx, task); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to start task")
		}

		result = task
		return nil
	})

	return result, err
}
