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

type ActivateVersionHandler struct {
	uow port.UnitOfWork
}

func NewActivateVersionHandler(uow port.UnitOfWork) *ActivateVersionHandler {
	return &ActivateVersionHandler{uow: uow}
}

func (h *ActivateVersionHandler) Handle(ctx context.Context, cmd dto.ActivateVersionCommand) (*operator.Operator, error) {
	var result *operator.Operator
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		op, err := repos.Operators.GetWithActiveVersion(ctx, cmd.OperatorID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", cmd.OperatorID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		target, err := repos.OperatorVersions.Get(ctx, cmd.VersionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator_version", cmd.VersionID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get target version")
		}
		if target.OperatorID != op.ID {
			return apperr.InvalidInput("version does not belong to operator")
		}

		if op.ActiveVersionID != nil && *op.ActiveVersionID != target.ID {
			current, err := repos.OperatorVersions.Get(ctx, *op.ActiveVersionID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return apperr.NotFound("operator_version", op.ActiveVersionID.String())
				}
				return apperr.Wrap(err, apperr.CodeDBError, "failed to get current active version")
			}
			current.Status = operator.VersionStatusArchived
			if err := repos.OperatorVersions.Update(ctx, current); err != nil {
				return apperr.Wrap(err, apperr.CodeDBError, "failed to archive current active version")
			}
		}

		target.Status = operator.VersionStatusActive
		if err := repos.OperatorVersions.Update(ctx, target); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to activate target version")
		}

		op.ActiveVersionID = &target.ID
		op.ActiveVersion = target
		syncOperatorCompatFieldsFromVersion(op, target)

		// 每次切换版本后，将算子状态重置为草稿，强制重新发布以进行完整校验
		if op.Status == operator.StatusPublished || op.Status == operator.StatusDeprecated {
			op.Status = operator.StatusDraft
		}

		if err := repos.Operators.Update(ctx, op); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to update operator active version")
		}

		result = op
		return nil
	})

	return result, err
}
