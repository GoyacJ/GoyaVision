package command

import (
	"context"
	"errors"
	"time"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateTaskContextSnapshotHandler struct {
	uow port.UnitOfWork
}

func NewCreateTaskContextSnapshotHandler(uow port.UnitOfWork) *CreateTaskContextSnapshotHandler {
	return &CreateTaskContextSnapshotHandler{uow: uow}
}

func (h *CreateTaskContextSnapshotHandler) Handle(ctx context.Context, cmd dto.CreateTaskContextSnapshotCommand) (*workflow.TaskContextSnapshot, error) {
	if cmd.TaskID == uuid.Nil {
		return nil, apperr.InvalidInput("task_id is required")
	}

	trigger := cmd.Trigger
	if trigger == "" {
		trigger = "manual"
	}

	var result *workflow.TaskContextSnapshot
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.Contexts == nil {
			return apperr.Internal("context repository is not configured", nil)
		}
		if _, err := repos.Tasks.Get(ctx, cmd.TaskID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("task", cmd.TaskID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get task")
		}

		state, err := repos.Contexts.GetState(ctx, cmd.TaskID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("task_context", cmd.TaskID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get task context")
		}

		snapshot := &workflow.TaskContextSnapshot{
			ID:        uuid.New(),
			TaskID:    cmd.TaskID,
			Version:   state.Version,
			Data:      state.Data,
			Trigger:   trigger,
			CreatedAt: time.Now().UTC(),
		}
		if err := repos.Contexts.CreateSnapshot(ctx, snapshot); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create context snapshot")
		}
		result = snapshot
		return nil
	})
	return result, err
}
