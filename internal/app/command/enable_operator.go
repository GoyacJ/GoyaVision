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

type EnableOperatorHandler struct {
	uow port.UnitOfWork
}

func NewEnableOperatorHandler(uow port.UnitOfWork) *EnableOperatorHandler {
	return &EnableOperatorHandler{uow: uow}
}

func (h *EnableOperatorHandler) Handle(ctx context.Context, cmd dto.EnableOperatorCommand) (*operator.Operator, error) {
	var result *operator.Operator
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		op, err := repos.Operators.Get(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		if cmd.Enabled {
			op.Status = operator.StatusEnabled
		} else {
			op.Status = operator.StatusDisabled
		}

		if err := repos.Operators.Update(ctx, op); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to update operator status")
		}

		result = op
		return nil
	})

	return result, err
}
