package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	"goyavision/pkg/apperr"
)

type ListOperatorsHandler struct {
	uow port.UnitOfWork
}

func NewListOperatorsHandler(uow port.UnitOfWork) *ListOperatorsHandler {
	return &ListOperatorsHandler{uow: uow}
}

func (h *ListOperatorsHandler) Handle(ctx context.Context, query dto.ListOperatorsQuery) (*dto.PagedResult[*operator.Operator], error) {
	query.Pagination.Normalize()

	filter := operator.Filter{
		Category:  query.Category,
		Type:      query.Type,
		Status:    query.Status,
		IsBuiltin: query.IsBuiltin,
		Tags:      query.Tags,
		Keyword:   query.Keyword,
		Limit:     query.Pagination.Limit,
		Offset:    query.Pagination.Offset,
	}

	var items []*operator.Operator
	var total int64
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		items, total, err = repos.Operators.List(ctx, filter)
		return err
	})

	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeDBError, "failed to list operators")
	}

	return &dto.PagedResult[*operator.Operator]{
		Items:  items,
		Total:  total,
		Limit:  query.Pagination.Limit,
		Offset: query.Pagination.Offset,
	}, nil
}
