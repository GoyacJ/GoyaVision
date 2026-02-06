package command

import (
	"context"
	"errors"

	adaptermcp "goyavision/internal/adapter/mcp"
	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	port2 "goyavision/internal/port"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
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
			tplPayload := adaptermcp.ToolToTemplate(cmd.ServerID, tools[i])
			code := tplPayload.Code

			existing, getErr := repos.OperatorTemplates.GetByCode(ctx, code)
			if getErr != nil && !errors.Is(getErr, gorm.ErrRecordNotFound) {
				return apperr.Wrap(getErr, apperr.CodeDBError, "failed to query mcp template by code")
			}
			if getErr == nil && existing != nil {
				existing.Name = tplPayload.Name
				existing.Description = tplPayload.Description
				existing.Category = tplPayload.Category
				existing.Type = tplPayload.Type
				existing.ExecMode = tplPayload.ExecMode
				existing.ExecConfig = tplPayload.ExecConfig
				existing.InputSchema = tplPayload.InputSchema
				existing.OutputSpec = tplPayload.OutputSpec
				existing.Config = tplPayload.Config
				existing.Author = tplPayload.Author
				existing.Tags = tplPayload.Tags
				if err := repos.OperatorTemplates.Update(ctx, existing); err != nil {
					return apperr.Wrap(err, apperr.CodeDBError, "failed to update mcp template")
				}
				result.Updated++
				continue
			}

			if err := repos.OperatorTemplates.Create(ctx, tplPayload); err != nil {
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
