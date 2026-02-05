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

type UpdateTaskHandler struct {
	uow port.UnitOfWork
}

func NewUpdateTaskHandler(uow port.UnitOfWork) *UpdateTaskHandler {
	return &UpdateTaskHandler{uow: uow}
}

func (h *UpdateTaskHandler) Handle(ctx context.Context, cmd dto.UpdateTaskCommand) (*workflow.Task, error) {
	var result *workflow.Task
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		task, err := repos.Tasks.Get(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("task", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get task")
		}

		if cmd.Status != nil {
			if *cmd.Status == workflow.TaskStatusRunning && task.StartedAt == nil {
				now := time.Now()
				task.StartedAt = &now
			}
			if (*cmd.Status == workflow.TaskStatusSuccess || *cmd.Status == workflow.TaskStatusFailed || *cmd.Status == workflow.TaskStatusCancelled) && task.CompletedAt == nil {
				now := time.Now()
				task.CompletedAt = &now
			}
			task.Status = *cmd.Status
		}

		if cmd.Progress != nil {
			progress := *cmd.Progress
			if progress < 0 {
				progress = 0
			}
			if progress > 100 {
				progress = 100
			}
			task.Progress = progress
		}

		if cmd.CurrentNode != nil {
			task.CurrentNode = *cmd.CurrentNode
		}

		if cmd.Error != nil {
			task.Error = *cmd.Error
		}

		if err := repos.Tasks.Update(ctx, task); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to update task")
		}

		result = task
		return nil
	})

	return result, err
}
