package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/port"
	"goyavision/pkg/apperr"
)

type ListMCPToolsHandler struct {
	client port.MCPClient
}

func NewListMCPToolsHandler(client port.MCPClient) *ListMCPToolsHandler {
	return &ListMCPToolsHandler{client: client}
}

func (h *ListMCPToolsHandler) Handle(ctx context.Context, query dto.ListMCPToolsQuery) ([]port.MCPTool, error) {
	if query.ServerID == "" {
		return nil, apperr.InvalidInput("server_id is required")
	}
	if h.client == nil {
		return nil, apperr.ServiceUnavailable("mcp client is not configured")
	}

	tools, err := h.client.ListTools(ctx, query.ServerID)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeServiceUnavailable, "failed to list mcp tools")
	}

	return tools, nil
}
