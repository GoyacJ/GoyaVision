package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/port"
	"goyavision/pkg/apperr"
)

type PreviewMCPToolHandler struct {
	client port.MCPClient
}

func NewPreviewMCPToolHandler(client port.MCPClient) *PreviewMCPToolHandler {
	return &PreviewMCPToolHandler{client: client}
}

func (h *PreviewMCPToolHandler) Handle(ctx context.Context, query dto.PreviewMCPToolQuery) (*port.MCPTool, error) {
	if query.ServerID == "" {
		return nil, apperr.InvalidInput("server_id is required")
	}
	if query.ToolName == "" {
		return nil, apperr.InvalidInput("tool_name is required")
	}
	if h.client == nil {
		return nil, apperr.ServiceUnavailable("mcp client is not configured")
	}

	tools, err := h.client.ListTools(ctx, query.ServerID)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeServiceUnavailable, "failed to list mcp tools")
	}

	for i := range tools {
		if tools[i].Name == query.ToolName {
			tool := tools[i]
			return &tool, nil
		}
	}

	return nil, apperr.NotFound("mcp tool", query.ToolName)
}
