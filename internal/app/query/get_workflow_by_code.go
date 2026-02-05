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

type GetWorkflowByCodeHandler struct {
	uow port.UnitOfWork
}

func NewGetWorkflowByCodeHandler(uow port.UnitOfWork) *GetWorkflowByCodeHandler {
	return &GetWorkflowByCodeHandler{uow: uow}
}

func (h *GetWorkflowByCodeHandler) Handle(ctx context.Context, query dto.GetWorkflowByCodeQuery) (*workflow.Workflow, error) {
	var result *workflow.Workflow
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		wf, err := repos.Workflows.GetByCode(ctx, query.Code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("workflow", query.Code)
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get workflow by code")
		}
		result = wf
		return nil
	})

	return result, err
}
