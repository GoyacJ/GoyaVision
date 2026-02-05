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

type GetWorkflowHandler struct {
	uow port.UnitOfWork
}

func NewGetWorkflowHandler(uow port.UnitOfWork) *GetWorkflowHandler {
	return &GetWorkflowHandler{uow: uow}
}

func (h *GetWorkflowHandler) Handle(ctx context.Context, query dto.GetWorkflowQuery) (*workflow.Workflow, error) {
	var result *workflow.Workflow
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		wf, err := repos.Workflows.Get(ctx, query.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("workflow", query.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get workflow")
		}
		result = wf
		return nil
	})

	return result, err
}
