package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"
)

type ListTasksHandler struct {
	uow port.UnitOfWork
}

func NewListTasksHandler(uow port.UnitOfWork) *ListTasksHandler {
	return &ListTasksHandler{uow: uow}
}

func (h *ListTasksHandler) Handle(ctx context.Context, query dto.ListTasksQuery) (*dto.PagedResult[*workflow.Task], error) {
	query.Pagination.Normalize()

	filter := workflow.TaskFilter{
		WorkflowID:        query.WorkflowID,
		AssetID:           query.AssetID,
		Status:            query.Status,
		TriggeredByUserID: query.TriggeredByUserID,
		From:              query.From,
		To:                query.To,
		Limit:             query.Pagination.Limit,
		Offset:            query.Pagination.Offset,
	}

	var items []*workflow.Task
	var total int64
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		items, total, err = repos.Tasks.List(ctx, filter)
		return err
	})

	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeDBError, "failed to list tasks")
	}

	return &dto.PagedResult[*workflow.Task]{
		Items:  items,
		Total:  total,
		Limit:  query.Pagination.Limit,
		Offset: query.Pagination.Offset,
	}, nil
}
