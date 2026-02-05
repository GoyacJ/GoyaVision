package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/media"
	"goyavision/pkg/apperr"
)

type ListSourcesHandler struct {
	uow port.UnitOfWork
}

func NewListSourcesHandler(uow port.UnitOfWork) *ListSourcesHandler {
	return &ListSourcesHandler{uow: uow}
}

func (h *ListSourcesHandler) Handle(ctx context.Context, query dto.ListSourcesQuery) (*dto.PagedResult[*media.Source], error) {
	query.Pagination.Normalize()

	filter := media.SourceFilter{
		Type:   query.Type,
		Limit:  query.Pagination.Limit,
		Offset: query.Pagination.Offset,
	}

	var sources []*media.Source
	var total int64
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		sources, total, err = repos.Sources.List(ctx, filter)
		return err
	})

	if err != nil {
		return nil, apperr.Internal("list sources", err)
	}

	return &dto.PagedResult[*media.Source]{
		Items:  sources,
		Total:  total,
		Limit:  query.Pagination.Limit,
		Offset: query.Pagination.Offset,
	}, nil
}
