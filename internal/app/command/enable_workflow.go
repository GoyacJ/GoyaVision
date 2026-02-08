package command

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type EnableWorkflowHandler struct {
	uow port.UnitOfWork
}

func NewEnableWorkflowHandler(uow port.UnitOfWork) *EnableWorkflowHandler {
	return &EnableWorkflowHandler{uow: uow}
}

func (h *EnableWorkflowHandler) Handle(ctx context.Context, cmd dto.EnableWorkflowCommand) (*workflow.Workflow, error) {
	var result *workflow.Workflow
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		wf, err := repos.Workflows.GetWithNodes(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("workflow", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get workflow")
		}

		if cmd.Enabled {
			if len(wf.Nodes) == 0 {
				return apperr.InvalidInput("workflow must have at least one node to enable")
			}
			wf.Status = workflow.StatusEnabled
			if wf.CurrentRevisionID == nil || wf.CurrentRevision == 0 {
				if _, err := persistAndActivateWorkflowRevision(ctx, repos, wf, 1); err != nil {
					return err
				}
			}
		} else {
			wf.Status = workflow.StatusDisabled
		}

		if err := repos.Workflows.Update(ctx, wf); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to update workflow status")
		}
		reloaded, err := repos.Workflows.GetWithNodes(ctx, wf.ID)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to reload workflow")
		}
		result = reloaded
		return nil
	})

	return result, err
}
