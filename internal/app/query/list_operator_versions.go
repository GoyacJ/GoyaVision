package query

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type ListOperatorVersionsHandler struct {
	uow port.UnitOfWork
}

func NewListOperatorVersionsHandler(uow port.UnitOfWork) *ListOperatorVersionsHandler {
	return &ListOperatorVersionsHandler{uow: uow}
}

func (h *ListOperatorVersionsHandler) Handle(ctx context.Context, q dto.ListOperatorVersionsQuery) (*dto.PagedResult[*operator.OperatorVersion], error) {
	q.Pagination.Normalize()

	var items []*operator.OperatorVersion
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if _, err := repos.Operators.Get(ctx, q.OperatorID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", q.OperatorID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		versions, err := repos.OperatorVersions.ListByOperator(ctx, q.OperatorID)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to list operator versions")
		}
		items = versions
		return nil
	})
	if err != nil {
		return nil, err
	}

	total := int64(len(items))
	start := q.Pagination.Offset
	if start > len(items) {
		start = len(items)
	}
	end := start + q.Pagination.Limit
	if end > len(items) {
		end = len(items)
	}

	return &dto.PagedResult[*operator.OperatorVersion]{
		Items:  items[start:end],
		Total:  total,
		Limit:  q.Pagination.Limit,
		Offset: q.Pagination.Offset,
	}, nil
}
