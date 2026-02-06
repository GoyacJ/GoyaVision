package command

import (
	"context"
	"fmt"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	port2 "goyavision/internal/port"
	"goyavision/pkg/apperr"
)

type SyncMCPTemplatesHandler struct {
	uow       port.UnitOfWork
	mcpClient port2.MCPClient
}

func NewSyncMCPTemplatesHandler(uow port.UnitOfWork, mcpClient port2.MCPClient) *SyncMCPTemplatesHandler {
	return &SyncMCPTemplatesHandler{uow: uow, mcpClient: mcpClient}
}

func (h *SyncMCPTemplatesHandler) Handle(ctx context.Context, cmd dto.SyncMCPTemplatesCommand) (*dto.SyncMCPTemplatesResult, error) {
	if cmd.ServerID == "" {
		return nil, apperr.InvalidInput("server_id is required")
	}
	if h.mcpClient == nil {
		return nil, apperr.ServiceUnavailable("mcp client is not configured")
	}

	tools, err := h.mcpClient.ListTools(ctx, cmd.ServerID)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeServiceUnavailable, "failed to list mcp tools")
	}

	result := &dto.SyncMCPTemplatesResult{ServerID: cmd.ServerID, Total: len(tools)}

	err = h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		for i := range tools {
			code := fmt.Sprintf("mcp_%s_%s", cmd.ServerID, tools[i].Name)

			existing, getErr := repos.OperatorTemplates.GetByCode(ctx, code)
			if getErr == nil && existing != nil {
				existing.Name = tools[i].Name
				existing.Description = tools[i].Description
				existing.ExecMode = operator.ExecModeMCP
				existing.ExecConfig = &operator.ExecConfig{MCP: &operator.MCPExecConfig{ServerID: cmd.ServerID, ToolName: tools[i].Name, ToolVersion: tools[i].Version}}
				existing.InputSchema = tools[i].InputSchema
				existing.OutputSpec = tools[i].OutputSchema
				if err := repos.OperatorTemplates.Update(ctx, existing); err != nil {
					return apperr.Wrap(err, apperr.CodeDBError, "failed to update mcp template")
				}
				result.Updated++
				continue
			}

			tpl := &operator.OperatorTemplate{
				Code:        code,
				Name:        tools[i].Name,
				Description: tools[i].Description,
				Category:    operator.CategoryUtility,
				Type:        operator.TypeTranscode,
				ExecMode:    operator.ExecModeMCP,
				ExecConfig: &operator.ExecConfig{MCP: &operator.MCPExecConfig{
					ServerID:    cmd.ServerID,
					ToolName:    tools[i].Name,
					ToolVersion: tools[i].Version,
				}},
				InputSchema: tools[i].InputSchema,
				OutputSpec:  tools[i].OutputSchema,
				Config:      map[string]interface{}{},
				Author:      "mcp",
				Tags:        []string{"mcp", cmd.ServerID},
			}
			if err := repos.OperatorTemplates.Create(ctx, tpl); err != nil {
				return apperr.Wrap(err, apperr.CodeDBError, "failed to create mcp template")
			}
			result.Created++
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
