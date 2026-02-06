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

type GetOperatorVersionHandler struct {
	uow port.UnitOfWork
}

func NewGetOperatorVersionHandler(uow port.UnitOfWork) *GetOperatorVersionHandler {
	return &GetOperatorVersionHandler{uow: uow}
}

func (h *GetOperatorVersionHandler) Handle(ctx context.Context, q dto.GetOperatorVersionQuery) (*operator.OperatorVersion, error) {
	var result *operator.OperatorVersion
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if _, err := repos.Operators.Get(ctx, q.OperatorID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", q.OperatorID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		v, err := repos.OperatorVersions.Get(ctx, q.VersionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator_version", q.VersionID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator version")
		}
		if v.OperatorID != q.OperatorID {
			return apperr.InvalidInput("version does not belong to operator")
		}

		result = v
		return nil
	})

	return result, err
}
