package command

import (
	"context"
	"errors"
	"strings"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	port2 "goyavision/internal/port"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type PublishOperatorHandler struct {
	uow       port.UnitOfWork
	mcpClient port2.MCPClient
	validator port.SchemaValidator
}

func NewPublishOperatorHandler(uow port.UnitOfWork, mcpClient port2.MCPClient, validator port.SchemaValidator) *PublishOperatorHandler {
	return &PublishOperatorHandler{uow: uow, mcpClient: mcpClient, validator: validator}
}

func (h *PublishOperatorHandler) Handle(ctx context.Context, cmd dto.PublishOperatorCommand) (*operator.Operator, error) {
	var result *operator.Operator
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		op, err := repos.Operators.GetWithActiveVersion(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		if op.ActiveVersion == nil || op.ActiveVersionID == nil {
			return apperr.InvalidInput("publish requires active version")
		}

		if h.validator != nil {
			if op.ActiveVersion.InputSchema != nil {
				if err := h.validator.IsValidJSONSchema(ctx, op.ActiveVersion.InputSchema); err != nil {
					return apperr.Wrap(err, apperr.CodeInvalidInput, "invalid active version input schema")
				}
			}
			if op.ActiveVersion.OutputSpec != nil {
				if err := h.validator.IsValidJSONSchema(ctx, op.ActiveVersion.OutputSpec); err != nil {
					return apperr.Wrap(err, apperr.CodeInvalidInput, "invalid active version output spec")
				}
			}
		}

		if ok, unmet, err := repos.OperatorDependencies.CheckDependenciesSatisfied(ctx, op.ID); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to check operator dependencies")
		} else if !ok {
			return apperr.InvalidInput("publish blocked by unmet dependencies: " + strings.Join(unmet, ","))
		}

		if op.ActiveVersion.ExecMode == operator.ExecModeMCP {
			if h.mcpClient == nil {
				return apperr.ServiceUnavailable("mcp client is not configured")
			}
			if op.ActiveVersion.ExecConfig == nil || op.ActiveVersion.ExecConfig.MCP == nil {
				return apperr.InvalidInput("mcp exec config is required for mcp operator")
			}

			mcpCfg := op.ActiveVersion.ExecConfig.MCP
			if mcpCfg.ServerID == "" || mcpCfg.ToolName == "" {
				return apperr.InvalidInput("mcp server_id and tool_name are required")
			}

			if err := h.mcpClient.HealthCheck(ctx, mcpCfg.ServerID); err != nil {
				return apperr.Wrap(err, apperr.CodeServiceUnavailable, "mcp server health check failed")
			}

			tools, err := h.mcpClient.ListTools(ctx, mcpCfg.ServerID)
			if err != nil {
				return apperr.Wrap(err, apperr.CodeServiceUnavailable, "failed to list mcp tools")
			}
			found := false
			for i := range tools {
				if tools[i].Name == mcpCfg.ToolName {
					found = true
					break
				}
			}
			if !found {
				return apperr.NotFound("mcp tool", mcpCfg.ToolName)
			}
		}

		op.Status = operator.StatusPublished
		if err := repos.Operators.Update(ctx, op); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to publish operator")
		}

		result = op
		return nil
	})

	return result, err
}
