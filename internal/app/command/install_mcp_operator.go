package command

import (
	"context"
	"errors"
	"fmt"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	port2 "goyavision/internal/port"
	"goyavision/pkg/apperr"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InstallMCPOperatorHandler struct {
	uow       port.UnitOfWork
	mcpClient port2.MCPClient
}

func NewInstallMCPOperatorHandler(uow port.UnitOfWork, mcpClient port2.MCPClient) *InstallMCPOperatorHandler {
	return &InstallMCPOperatorHandler{uow: uow, mcpClient: mcpClient}
}

func (h *InstallMCPOperatorHandler) Handle(ctx context.Context, cmd dto.InstallMCPOperatorCommand) (*operator.Operator, error) {
	if cmd.ServerID == "" || cmd.ToolName == "" {
		return nil, apperr.InvalidInput("server_id and tool_name are required")
	}
	if cmd.OperatorCode == "" || cmd.OperatorName == "" {
		return nil, apperr.InvalidInput("operator_code and operator_name are required")
	}

	if h.mcpClient == nil {
		return nil, apperr.ServiceUnavailable("mcp client is not configured")
	}

	tools, err := h.mcpClient.ListTools(ctx, cmd.ServerID)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeServiceUnavailable, "failed to list mcp tools")
	}

	found := false
	for i := range tools {
		if tools[i].Name == cmd.ToolName {
			found = true
			break
		}
	}
	if !found {
		return nil, apperr.NotFound("mcp tool", cmd.ToolName)
	}

	category := operator.CategoryUtility
	if cmd.Category != nil {
		category = *cmd.Category
	}
	typeVal := operator.TypeTranscode
	if cmd.Type != nil {
		typeVal = *cmd.Type
	}

	var result *operator.Operator
	err = h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if _, err := repos.Operators.GetByCode(ctx, cmd.OperatorCode); err == nil {
			return apperr.Conflict(fmt.Sprintf("operator with code %s already exists", cmd.OperatorCode))
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to check operator code uniqueness")
		}

		op := &operator.Operator{
			Code:        cmd.OperatorCode,
			Name:        cmd.OperatorName,
			Description: fmt.Sprintf("Installed from MCP tool %s/%s", cmd.ServerID, cmd.ToolName),
			Category:    category,
			Type:        typeVal,
			Origin:      operator.OriginMCP,
			Status:      operator.StatusDraft,
			Tags:        cmd.Tags,
		}
		if err := repos.Operators.Create(ctx, op); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create mcp operator")
		}

		execConfig := &operator.ExecConfig{MCP: &operator.MCPExecConfig{
			ServerID:   cmd.ServerID,
			ToolName:   cmd.ToolName,
			TimeoutSec: cmd.TimeoutSec,
		}}
		ov := &operator.OperatorVersion{
			ID:          uuid.New(),
			OperatorID:  op.ID,
			Version:     "1.0.0",
			ExecMode:    operator.ExecModeMCP,
			ExecConfig:  execConfig,
			InputSchema: map[string]interface{}{},
			OutputSpec:  map[string]interface{}{},
			Config:      map[string]interface{}{},
			Status:      operator.VersionStatusActive,
		}
		if err := repos.OperatorVersions.Create(ctx, ov); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create mcp operator initial version")
		}

		op.ActiveVersionID = &ov.ID
		op.ActiveVersion = ov
		syncOperatorCompatFieldsFromVersion(op, ov)
		if err := repos.Operators.Update(ctx, op); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to bind operator active version")
		}

		result = op
		return nil
	})

	return result, err
}
