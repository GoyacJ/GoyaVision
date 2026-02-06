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

type DeleteOperatorHandler struct {
	uow port.UnitOfWork
}

func NewDeleteOperatorHandler(uow port.UnitOfWork) *DeleteOperatorHandler {
	return &DeleteOperatorHandler{uow: uow}
}

func (h *DeleteOperatorHandler) Handle(ctx context.Context, cmd dto.DeleteOperatorCommand) error {
	return h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		op, err := repos.Operators.Get(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		if op.Origin == operator.OriginBuiltin {
			return apperr.InvalidInput("cannot delete builtin operator")
		}

		versions, err := repos.OperatorVersions.ListByOperator(ctx, cmd.ID)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to list operator versions")
		}
		for _, v := range versions {
			if err := repos.OperatorVersions.Delete(ctx, v.ID); err != nil {
				return apperr.Wrap(err, apperr.CodeDBError, "failed to delete operator version")
			}
		}

		if err := repos.OperatorDependencies.DeleteByOperator(ctx, cmd.ID); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to delete operator dependencies")
		}

		if err := repos.Operators.Delete(ctx, cmd.ID); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to delete operator")
		}

		return nil
	})
}
