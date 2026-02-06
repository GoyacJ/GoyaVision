package port

import "context"

// MCPTool 描述 MCP Server 暴露的工具元信息。
type MCPTool struct {
	Name        string
	Description string
	Version     string
	InputSchema map[string]interface{}
	OutputSchema map[string]interface{}
}

// MCPServer 描述可用的 MCP 服务端。
type MCPServer struct {
	ID          string
	Name        string
	Description string
	Status      string
}

// MCPClient 定义 MCP 工具调用能力。
type MCPClient interface {
	ListTools(ctx context.Context, serverID string) ([]MCPTool, error)
	CallTool(ctx context.Context, serverID, toolName string, args map[string]interface{}) (map[string]interface{}, error)
	HealthCheck(ctx context.Context, serverID string) error
}

// MCPRegistry 定义 MCP Server 管理能力。
type MCPRegistry interface {
	ListServers(ctx context.Context) ([]MCPServer, error)
	GetServer(ctx context.Context, serverID string) (*MCPServer, error)
}
