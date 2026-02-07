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

type UpdateOperatorHandler struct {
	uow port.UnitOfWork
}

func NewUpdateOperatorHandler(uow port.UnitOfWork) *UpdateOperatorHandler {
	return &UpdateOperatorHandler{uow: uow}
}

func (h *UpdateOperatorHandler) Handle(ctx context.Context, cmd dto.UpdateOperatorCommand) (*operator.Operator, error) {
	var result *operator.Operator
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		op, err := repos.Operators.Get(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		if op.Origin == operator.OriginBuiltin {
			return apperr.InvalidInput("cannot update builtin operator")
		}

		if cmd.Name != nil {
			op.Name = *cmd.Name
		}
		if cmd.Description != nil {
			op.Description = *cmd.Description
		}
		if cmd.Category != nil {
			op.Category = *cmd.Category
		}
		if len(cmd.Tags) > 0 {
			op.Tags = cmd.Tags
		}

		if err := repos.Operators.Update(ctx, op); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to update operator")
		}

		result = op
		return nil
	})

	return result, err
}
