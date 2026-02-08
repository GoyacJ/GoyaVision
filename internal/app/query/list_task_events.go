package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/agent"
	"goyavision/pkg/apperr"
)

type ListTaskEventsHandler struct {
	uow port.UnitOfWork
}

func NewListTaskEventsHandler(uow port.UnitOfWork) *ListTaskEventsHandler {
	return &ListTaskEventsHandler{uow: uow}
}

func (h *ListTaskEventsHandler) Handle(ctx context.Context, query dto.ListTaskEventsQuery) (*dto.PagedResult[*agent.RunEvent], error) {
	query.Pagination.Normalize()

	var items []*agent.RunEvent
	var total int64
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.RunEvents == nil {
			return apperr.Internal("run event repository is not configured", nil)
		}
		var err error
		items, total, err = repos.RunEvents.List(ctx, agent.EventFilter{
			TaskID:  &query.TaskID,
			Limit:   query.Pagination.Limit,
			Offset:  query.Pagination.Offset,
			Source:  query.Source,
			NodeKey: query.NodeKey,
		})
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to list task events")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dto.PagedResult[*agent.RunEvent]{
		Items:  items,
		Total:  total,
		Limit:  query.Pagination.Limit,
		Offset: query.Pagination.Offset,
	}, nil
}
