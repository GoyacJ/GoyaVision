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

type CompleteTaskHandler struct {
	uow port.UnitOfWork
}

func NewCompleteTaskHandler(uow port.UnitOfWork) *CompleteTaskHandler {
	return &CompleteTaskHandler{uow: uow}
}

func (h *CompleteTaskHandler) Handle(ctx context.Context, cmd dto.CompleteTaskCommand) (*workflow.Task, error) {
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
		task.Status = workflow.TaskStatusSuccess
		task.CompletedAt = &now
		task.Progress = 100

		if err := repos.Tasks.Update(ctx, task); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to complete task")
		}

		result = task
		return nil
	})

	return result, err
}
