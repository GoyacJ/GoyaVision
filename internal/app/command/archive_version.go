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

type ArchiveVersionHandler struct {
	uow port.UnitOfWork
}

func NewArchiveVersionHandler(uow port.UnitOfWork) *ArchiveVersionHandler {
	return &ArchiveVersionHandler{uow: uow}
}

func (h *ArchiveVersionHandler) Handle(ctx context.Context, cmd dto.ArchiveVersionCommand) (*operator.OperatorVersion, error) {
	var result *operator.OperatorVersion
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		op, err := repos.Operators.Get(ctx, cmd.OperatorID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", cmd.OperatorID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		version, err := repos.OperatorVersions.Get(ctx, cmd.VersionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator_version", cmd.VersionID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator version")
		}
		if version.OperatorID != op.ID {
			return apperr.InvalidInput("version does not belong to operator")
		}

		if op.ActiveVersionID != nil && *op.ActiveVersionID == version.ID {
			return apperr.InvalidInput("cannot archive active version")
		}

		version.Status = operator.VersionStatusArchived
		if err := repos.OperatorVersions.Update(ctx, version); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to archive operator version")
		}

		result = version
		return nil
	})

	return result, err
}
