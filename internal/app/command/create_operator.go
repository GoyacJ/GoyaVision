package command

import (
	"context"
	"fmt"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	"goyavision/pkg/apperr"
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
	if cmd.Endpoint == "" {
		return nil, apperr.InvalidInput("endpoint is required")
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

	var result *operator.Operator
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if _, err := repos.Operators.GetByCode(ctx, cmd.Code); err == nil {
			return apperr.Conflict(fmt.Sprintf("operator with code %s already exists", cmd.Code))
		}

		op := &operator.Operator{
			Code:        cmd.Code,
			Name:        cmd.Name,
			Description: cmd.Description,
			Category:    cmd.Category,
			Type:        cmd.Type,
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

		result = op
		return nil
	})

	return result, err
}
