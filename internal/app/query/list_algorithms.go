package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/algorithm"
	"goyavision/pkg/apperr"
)

type ListAlgorithmsHandler struct {
	uow port.UnitOfWork
}

func NewListAlgorithmsHandler(uow port.UnitOfWork) *ListAlgorithmsHandler {
	return &ListAlgorithmsHandler{uow: uow}
}

func (h *ListAlgorithmsHandler) Handle(ctx context.Context, query dto.ListAlgorithmsQuery) (*dto.PagedResult[*algorithm.Algorithm], error) {
	query.Pagination.Normalize()
	filter := algorithm.Filter{
		Status:   query.Status,
		Scenario: query.Scenario,
		Tags:     query.Tags,
		Keyword:  query.Keyword,
		Limit:    query.Pagination.Limit,
		Offset:   query.Pagination.Offset,
	}

	var items []*algorithm.Algorithm
	var total int64
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		items, total, err = repos.Algorithms.List(ctx, filter)
		return err
	})
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeDBError, "failed to list algorithms")
	}

	return &dto.PagedResult[*algorithm.Algorithm]{
		Items:  items,
		Total:  total,
		Limit:  query.Pagination.Limit,
		Offset: query.Pagination.Offset,
	}, nil
}
