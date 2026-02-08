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

type CreateWorkflowRevisionHandler struct {
	uow port.UnitOfWork
}

func NewCreateWorkflowRevisionHandler(uow port.UnitOfWork) *CreateWorkflowRevisionHandler {
	return &CreateWorkflowRevisionHandler{uow: uow}
}

func (h *CreateWorkflowRevisionHandler) Handle(ctx context.Context, cmd dto.CreateWorkflowRevisionCommand) (*workflow.WorkflowRevision, error) {
	if cmd.WorkflowID == uuid.Nil {
		return nil, apperr.InvalidInput("workflow_id is required")
	}

	var result *workflow.WorkflowRevision
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		wf, err := repos.Workflows.GetWithNodes(ctx, cmd.WorkflowID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("workflow", cmd.WorkflowID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get workflow")
		}

		nextRevision := nextWorkflowRevision(wf.CurrentRevision)
		rev := &workflow.WorkflowRevision{
			ID:         uuid.New(),
			WorkflowID: wf.ID,
			Revision:   nextRevision,
			Status:     workflow.RevisionStatusDraft,
			Definition: workflow.BuildDefinitionFromWorkflow(wf),
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}
		if cmd.Activate {
			rev.Status = workflow.RevisionStatusActive
		}

		if err := repos.Workflows.CreateRevision(ctx, rev); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create workflow revision")
		}
		if cmd.Activate {
			if err := repos.Workflows.ActivateRevision(ctx, wf.ID, rev.ID); err != nil {
				return apperr.Wrap(err, apperr.CodeDBError, "failed to activate workflow revision")
			}
		}
		result = rev
		return nil
	})
	return result, err
}
