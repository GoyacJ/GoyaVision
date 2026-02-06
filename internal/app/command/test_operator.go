package command

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type TestOperatorHandler struct {
	uow port.UnitOfWork
}

func NewTestOperatorHandler(uow port.UnitOfWork) *TestOperatorHandler {
	return &TestOperatorHandler{uow: uow}
}

func (h *TestOperatorHandler) Handle(ctx context.Context, cmd dto.TestOperatorCommand) (*dto.TestOperatorResult, error) {
	result := &dto.TestOperatorResult{Success: true, Message: "ok"}

	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		op, err := repos.Operators.GetWithActiveVersion(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		if op.ActiveVersion == nil && op.ActiveVersionID == nil {
			return apperr.InvalidInput("operator has no active version")
		}

		result.Diagnostics = map[string]interface{}{
			"operator_id": op.ID.String(),
			"status":      string(op.Status),
			"checked":     true,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
