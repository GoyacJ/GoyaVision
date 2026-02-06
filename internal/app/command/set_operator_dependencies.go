package command

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	"goyavision/pkg/apperr"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SetOperatorDependenciesHandler struct {
	uow port.UnitOfWork
}

func NewSetOperatorDependenciesHandler(uow port.UnitOfWork) *SetOperatorDependenciesHandler {
	return &SetOperatorDependenciesHandler{uow: uow}
}

func (h *SetOperatorDependenciesHandler) Handle(ctx context.Context, cmd dto.SetOperatorDependenciesCommand) error {
	if cmd.OperatorID == uuid.Nil {
		return apperr.InvalidInput("operator_id is required")
	}

	return h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if _, err := repos.Operators.Get(ctx, cmd.OperatorID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", cmd.OperatorID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		if err := repos.OperatorDependencies.DeleteByOperator(ctx, cmd.OperatorID); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to clear operator dependencies")
		}

		for i := range cmd.Dependencies {
			dep := cmd.Dependencies[i]
			if dep.DependsOnID == uuid.Nil {
				return apperr.InvalidInput("depends_on_id is required")
			}
			if dep.DependsOnID == cmd.OperatorID {
				return apperr.InvalidInput("operator cannot depend on itself")
			}

			if _, err := repos.Operators.Get(ctx, dep.DependsOnID); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return apperr.NotFound("depends_on operator", dep.DependsOnID.String())
				}
				return apperr.Wrap(err, apperr.CodeDBError, "failed to validate dependency operator")
			}

			item := &operator.OperatorDependency{
				ID:          uuid.New(),
				OperatorID:  cmd.OperatorID,
				DependsOnID: dep.DependsOnID,
				MinVersion:  dep.MinVersion,
				IsOptional:  dep.IsOptional,
			}
			if err := repos.OperatorDependencies.Create(ctx, item); err != nil {
				return apperr.Wrap(err, apperr.CodeDBError, "failed to create operator dependency")
			}
		}

		return nil
	})
}
