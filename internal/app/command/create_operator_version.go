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

type CreateOperatorVersionHandler struct {
	uow             port.UnitOfWork
	schemaValidator port.SchemaValidator
}

func NewCreateOperatorVersionHandler(uow port.UnitOfWork, schemaValidator port.SchemaValidator) *CreateOperatorVersionHandler {
	return &CreateOperatorVersionHandler{uow: uow, schemaValidator: schemaValidator}
}

func (h *CreateOperatorVersionHandler) Handle(ctx context.Context, cmd dto.CreateOperatorVersionCommand) (*operator.OperatorVersion, error) {
	if cmd.OperatorID == uuid.Nil {
		return nil, apperr.InvalidInput("operator_id is required")
	}
	if cmd.Version == "" {
		return nil, apperr.InvalidInput("version is required")
	}
	if cmd.ExecMode == "" {
		return nil, apperr.InvalidInput("exec_mode is required")
	}
	if err := validateSemver(cmd.Version); err != nil {
		return nil, err
	}
	if err := validateExecMode(cmd.ExecMode); err != nil {
		return nil, err
	}

	status := cmd.Status
	if status == "" {
		status = operator.VersionStatusDraft
	}

	if h.schemaValidator != nil {
		if cmd.InputSchema != nil {
			if err := h.schemaValidator.IsValidJSONSchema(ctx, cmd.InputSchema); err != nil {
				return nil, err
			}
		}
		if cmd.OutputSpec != nil {
			if err := h.schemaValidator.IsValidJSONSchema(ctx, cmd.OutputSpec); err != nil {
				return nil, err
			}
		}
	}

	var result *operator.OperatorVersion
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if _, err := repos.Operators.Get(ctx, cmd.OperatorID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", cmd.OperatorID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		if _, err := repos.OperatorVersions.GetByOperatorAndVersion(ctx, cmd.OperatorID, cmd.Version); err == nil {
			return apperr.Conflict("operator version already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to check operator version uniqueness")
		}

		v := &operator.OperatorVersion{
			ID:          uuid.New(),
			OperatorID:  cmd.OperatorID,
			Version:     cmd.Version,
			ExecMode:    cmd.ExecMode,
			ExecConfig:  cmd.ExecConfig,
			InputSchema: cmd.InputSchema,
			OutputSpec:  cmd.OutputSpec,
			Config:      cmd.Config,
			Changelog:   cmd.Changelog,
			Status:      status,
		}

		if err := repos.OperatorVersions.Create(ctx, v); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create operator version")
		}

		result = v
		return nil
	})

	return result, err
}
