package query

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type GetWorkflowRevisionHandler struct {
	uow port.UnitOfWork
}

func NewGetWorkflowRevisionHandler(uow port.UnitOfWork) *GetWorkflowRevisionHandler {
	return &GetWorkflowRevisionHandler{uow: uow}
}

func (h *GetWorkflowRevisionHandler) Handle(ctx context.Context, query dto.GetWorkflowRevisionQuery) (*workflow.WorkflowRevision, error) {
	var result *workflow.WorkflowRevision
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		rev, err := repos.Workflows.GetRevisionByNumber(ctx, query.WorkflowID, query.Revision)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("workflow_revision", query.WorkflowID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get workflow revision")
		}
		result = rev
		return nil
	})
	return result, err
}
