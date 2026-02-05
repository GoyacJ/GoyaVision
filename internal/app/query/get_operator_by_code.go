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

type GetOperatorByCodeHandler struct {
	uow port.UnitOfWork
}

func NewGetOperatorByCodeHandler(uow port.UnitOfWork) *GetOperatorByCodeHandler {
	return &GetOperatorByCodeHandler{uow: uow}
}

func (h *GetOperatorByCodeHandler) Handle(ctx context.Context, query dto.GetOperatorByCodeQuery) (*operator.Operator, error) {
	var result *operator.Operator
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		op, err := repos.Operators.GetByCode(ctx, query.Code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", query.Code)
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator by code")
		}
		result = op
		return nil
	})

	return result, err
}
