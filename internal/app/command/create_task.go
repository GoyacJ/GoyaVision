package command

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateTaskHandler struct {
	uow port.UnitOfWork
}

func NewCreateTaskHandler(uow port.UnitOfWork) *CreateTaskHandler {
	return &CreateTaskHandler{uow: uow}
}

func (h *CreateTaskHandler) Handle(ctx context.Context, cmd dto.CreateTaskCommand) (*workflow.Task, error) {
	if cmd.WorkflowID == uuid.Nil {
		return nil, apperr.InvalidInput("workflow_id is required")
	}

	var result *workflow.Task
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		wf, err := repos.Workflows.Get(ctx, cmd.WorkflowID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("workflow", cmd.WorkflowID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get workflow")
		}

		if !wf.IsEnabled() {
			return apperr.InvalidInput("workflow is not enabled")
		}
		if err := ensureWorkflowHasActiveRevision(wf); err != nil {
			return err
		}

		if cmd.AssetID != nil {
			if _, err := repos.Assets.Get(ctx, *cmd.AssetID); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return apperr.NotFound("asset", cmd.AssetID.String())
				}
				return apperr.Wrap(err, apperr.CodeDBError, "failed to get asset")
			}
		}

		task := &workflow.Task{
			WorkflowID:         cmd.WorkflowID,
			WorkflowRevisionID: wf.CurrentRevisionID,
			WorkflowRevision:   wf.CurrentRevision,
			AssetID:            cmd.AssetID,
			Status:             workflow.TaskStatusPending,
			Progress:           0,
			InputParams:        cmd.InputParams,
		}

		if err := repos.Tasks.Create(ctx, task); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create task")
		}

		taskWithRelations, err := repos.Tasks.GetWithRelations(ctx, task.ID)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get task with relations")
		}
		result = taskWithRelations
		return nil
	})

	return result, err
}
