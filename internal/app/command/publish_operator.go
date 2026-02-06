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

type PublishOperatorHandler struct {
	uow port.UnitOfWork
}

func NewPublishOperatorHandler(uow port.UnitOfWork) *PublishOperatorHandler {
	return &PublishOperatorHandler{uow: uow}
}

func (h *PublishOperatorHandler) Handle(ctx context.Context, cmd dto.PublishOperatorCommand) (*operator.Operator, error) {
	var result *operator.Operator
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		op, err := repos.Operators.GetWithActiveVersion(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		if op.ActiveVersion == nil && op.ActiveVersionID == nil {
			return apperr.InvalidInput("publish requires active version")
		}

		op.Status = operator.StatusPublished
		if err := repos.Operators.Update(ctx, op); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to publish operator")
		}

		result = op
		return nil
	})

	return result, err
}
