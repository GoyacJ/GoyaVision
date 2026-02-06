package query

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type CheckDependenciesHandler struct {
	uow port.UnitOfWork
}

func NewCheckDependenciesHandler(uow port.UnitOfWork) *CheckDependenciesHandler {
	return &CheckDependenciesHandler{uow: uow}
}

func (h *CheckDependenciesHandler) Handle(ctx context.Context, q dto.CheckDependenciesQuery) (bool, []string, error) {
	var ok bool
	var unmet []string
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if _, err := repos.Operators.Get(ctx, q.OperatorID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", q.OperatorID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		status, missing, err := repos.OperatorDependencies.CheckDependenciesSatisfied(ctx, q.OperatorID)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to check dependencies")
		}
		ok = status
		unmet = missing
		return nil
	})

	return ok, unmet, err
}
