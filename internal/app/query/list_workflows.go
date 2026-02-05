package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"
)

type ListWorkflowsHandler struct {
	uow port.UnitOfWork
}

func NewListWorkflowsHandler(uow port.UnitOfWork) *ListWorkflowsHandler {
	return &ListWorkflowsHandler{uow: uow}
}

func (h *ListWorkflowsHandler) Handle(ctx context.Context, query dto.ListWorkflowsQuery) (*dto.PagedResult[*workflow.Workflow], error) {
	query.Pagination.Normalize()

	filter := workflow.Filter{
		Status:      query.Status,
		TriggerType: query.TriggerType,
		Tags:        query.Tags,
		Keyword:     query.Keyword,
		Limit:       query.Pagination.Limit,
		Offset:      query.Pagination.Offset,
	}

	var items []*workflow.Workflow
	var total int64
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		items, total, err = repos.Workflows.List(ctx, filter)
		return err
	})

	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeDBError, "failed to list workflows")
	}

	return &dto.PagedResult[*workflow.Workflow]{
		Items:  items,
		Total:  total,
		Limit:  query.Pagination.Limit,
		Offset: query.Pagination.Offset,
	}, nil
}
