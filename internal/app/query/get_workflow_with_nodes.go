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

type GetWorkflowWithNodesHandler struct {
	uow port.UnitOfWork
}

func NewGetWorkflowWithNodesHandler(uow port.UnitOfWork) *GetWorkflowWithNodesHandler {
	return &GetWorkflowWithNodesHandler{uow: uow}
}

func (h *GetWorkflowWithNodesHandler) Handle(ctx context.Context, query dto.GetWorkflowWithNodesQuery) (*workflow.Workflow, error) {
	var result *workflow.Workflow
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		wf, err := repos.Workflows.GetWithNodes(ctx, query.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("workflow", query.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get workflow with nodes")
		}
		result = wf
		return nil
	})

	return result, err
}
