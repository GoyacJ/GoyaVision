package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/port"
	"goyavision/pkg/apperr"
)

type ListMCPServersHandler struct {
	registry port.MCPRegistry
}

func NewListMCPServersHandler(registry port.MCPRegistry) *ListMCPServersHandler {
	return &ListMCPServersHandler{registry: registry}
}

func (h *ListMCPServersHandler) Handle(ctx context.Context, _ dto.ListMCPServersQuery) ([]port.MCPServer, error) {
	if h.registry == nil {
		return nil, apperr.ServiceUnavailable("mcp registry is not configured")
	}

	servers, err := h.registry.ListServers(ctx)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeServiceUnavailable, "failed to list mcp servers")
	}

	return servers, nil
}
