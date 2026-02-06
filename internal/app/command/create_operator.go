package command

import (
	"context"
	"errors"
	"fmt"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	"goyavision/pkg/apperr"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateOperatorHandler struct {
	uow port.UnitOfWork
}

func NewCreateOperatorHandler(uow port.UnitOfWork) *CreateOperatorHandler {
	return &CreateOperatorHandler{uow: uow}
}

func (h *CreateOperatorHandler) Handle(ctx context.Context, cmd dto.CreateOperatorCommand) (*operator.Operator, error) {
	if cmd.Code == "" {
		return nil, apperr.InvalidInput("code is required")
	}
	if cmd.Name == "" {
		return nil, apperr.InvalidInput("name is required")
	}
	if cmd.Category == "" {
		return nil, apperr.InvalidInput("category is required")
	}
	if cmd.Type == "" {
		return nil, apperr.InvalidInput("type is required")
	}
	if cmd.Category != operator.CategoryAnalysis &&
		cmd.Category != operator.CategoryProcessing &&
		cmd.Category != operator.CategoryGeneration &&
		cmd.Category != operator.CategoryUtility {
		return nil, apperr.InvalidInput("invalid category")
	}

	version := "1.0.0"
	if cmd.Version != "" {
		version = cmd.Version
	}

	method := "POST"
	if cmd.Method != "" {
		method = cmd.Method
	}

	status := operator.StatusDraft
	if cmd.Status != "" {
		status = cmd.Status
	}

	origin := cmd.Origin
	if origin == "" {
		if cmd.IsBuiltin {
			origin = operator.OriginBuiltin
		} else {
			origin = operator.OriginCustom
		}
	}

	execMode := cmd.ExecMode
	if execMode == "" {
		execMode = operator.ExecModeHTTP
	}

	var result *operator.Operator
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if _, err := repos.Operators.GetByCode(ctx, cmd.Code); err == nil {
			return apperr.Conflict(fmt.Sprintf("operator with code %s already exists", cmd.Code))
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to check operator code uniqueness")
		}

		op := &operator.Operator{
			Code:        cmd.Code,
			Name:        cmd.Name,
			Description: cmd.Description,
			Category:    cmd.Category,
			Type:        cmd.Type,
			Origin:      origin,
			Version:     version,
			Endpoint:    cmd.Endpoint,
			Method:      method,
			InputSchema: cmd.InputSchema,
			OutputSpec:  cmd.OutputSpec,
			Config:      cmd.Config,
			Status:      status,
			IsBuiltin:   cmd.IsBuiltin,
			Tags:        cmd.Tags,
		}

		if err := repos.Operators.Create(ctx, op); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create operator")
		}

		// Phase A：创建首个版本并绑定为激活版本（兼容旧字段）
		execConfig := cmd.ExecConfig
		if execConfig == nil && (cmd.Endpoint != "" || method != "") {
			execConfig = &operator.ExecConfig{
				HTTP: &operator.HTTPExecConfig{
					Endpoint: cmd.Endpoint,
					Method:   method,
				},
			}
		}

		ov := &operator.OperatorVersion{
			ID:          uuid.New(),
			OperatorID:  op.ID,
			Version:     version,
			ExecMode:    execMode,
			ExecConfig:  execConfig,
			InputSchema: cmd.InputSchema,
			OutputSpec:  cmd.OutputSpec,
			Config:      cmd.Config,
			Status:      operator.VersionStatusActive,
		}

		if err := repos.OperatorVersions.Create(ctx, ov); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create operator initial version")
		}

		op.ActiveVersionID = &ov.ID
		op.ActiveVersion = ov
		if err := repos.Operators.Update(ctx, op); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to bind operator active version")
		}

		result = op
		return nil
	})

	return result, err
}
