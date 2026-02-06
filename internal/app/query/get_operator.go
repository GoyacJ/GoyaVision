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

type GetOperatorHandler struct {
	uow port.UnitOfWork
}

func NewGetOperatorHandler(uow port.UnitOfWork) *GetOperatorHandler {
	return &GetOperatorHandler{uow: uow}
}

func (h *GetOperatorHandler) Handle(ctx context.Context, query dto.GetOperatorQuery) (*operator.Operator, error) {
	var result *operator.Operator
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		op, err := repos.Operators.GetWithActiveVersion(ctx, query.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", query.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}
		result = op
		return nil
	})

	return result, err
}
