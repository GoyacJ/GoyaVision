package command

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type DeprecateOperatorHandler struct {
	uow port.UnitOfWork
}

func NewDeprecateOperatorHandler(uow port.UnitOfWork) *DeprecateOperatorHandler {
	return &DeprecateOperatorHandler{uow: uow}
}

func (h *DeprecateOperatorHandler) Handle(ctx context.Context, cmd dto.DeprecateOperatorCommand) (*operator.Operator, error) {
	var result *operator.Operator
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		op, err := repos.Operators.Get(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		op.Status = operator.StatusDeprecated
		if err := repos.Operators.Update(ctx, op); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to deprecate operator")
		}

		result = op
		return nil
	})

	return result, err
}
