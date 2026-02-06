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

type ListOperatorDependenciesHandler struct {
	uow port.UnitOfWork
}

func NewListOperatorDependenciesHandler(uow port.UnitOfWork) *ListOperatorDependenciesHandler {
	return &ListOperatorDependenciesHandler{uow: uow}
}

func (h *ListOperatorDependenciesHandler) Handle(ctx context.Context, q dto.ListOperatorDependenciesQuery) ([]*operator.OperatorDependency, error) {
	var result []*operator.OperatorDependency
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if _, err := repos.Operators.Get(ctx, q.OperatorID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", q.OperatorID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		deps, err := repos.OperatorDependencies.ListByOperator(ctx, q.OperatorID)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to list operator dependencies")
		}
		result = deps
		return nil
	})

	return result, err
}
