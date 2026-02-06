package engine

import (
	"context"
	"encoding/json"
	"fmt"

	"goyavision/internal/domain/operator"
	"goyavision/internal/port"
)

var _ port.OperatorExecutor = (*MCPOperatorExecutor)(nil)

// MCPOperatorExecutor MCP 算子执行器
type MCPOperatorExecutor struct {
	client port.MCPClient
}

func NewMCPOperatorExecutor(client port.MCPClient) *MCPOperatorExecutor {
	return &MCPOperatorExecutor{client: client}
}

func (e *MCPOperatorExecutor) Execute(ctx context.Context, version *operator.OperatorVersion, input *operator.Input) (*operator.Output, error) {
	if version == nil {
		return nil, fmt.Errorf("operator version is nil")
	}
	if version.ExecMode != operator.ExecModeMCP {
		return nil, fmt.Errorf("mcp executor does not support exec mode: %s", version.ExecMode)
	}
	if e.client == nil {
		return nil, fmt.Errorf("mcp client is not configured")
	}
	if version.ExecConfig == nil || version.ExecConfig.MCP == nil {
		return nil, fmt.Errorf("mcp exec config is required")
	}

	mcpCfg := version.ExecConfig.MCP
	if mcpCfg.ServerID == "" || mcpCfg.ToolName == "" {
		return nil, fmt.Errorf("mcp server_id and tool_name are required")
	}

	args := map[string]interface{}{}
	if input != nil {
		for k, v := range input.Params {
			args[k] = v
		}
		if input.AssetID.String() != "00000000-0000-0000-0000-000000000000" {
			args["asset_id"] = input.AssetID.String()
		}
	}

	result, err := e.client.CallTool(ctx, mcpCfg.ServerID, mcpCfg.ToolName, args)
	if err != nil {
		return nil, fmt.Errorf("failed to call mcp tool: %w", err)
	}

	if len(result) == 0 {
		return &operator.Output{}, nil
	}

	// 默认按标准输出结构反序列化
	raw, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal mcp result: %w", err)
	}

	var output operator.Output
	if err := json.Unmarshal(raw, &output); err == nil {
		return &output, nil
	}

	// 兼容非标准输出：归档到 diagnostics
	return &operator.Output{Diagnostics: result}, nil
}

func (e *MCPOperatorExecutor) Mode() operator.ExecMode {
	return operator.ExecModeMCP
}

func (e *MCPOperatorExecutor) HealthCheck(ctx context.Context, version *operator.OperatorVersion) error {
	if version == nil {
		return fmt.Errorf("operator version is nil")
	}
	if e.client == nil {
		return fmt.Errorf("mcp client is not configured")
	}
	if version.ExecConfig == nil || version.ExecConfig.MCP == nil {
		return fmt.Errorf("mcp exec config is required")
	}
	mcpCfg := version.ExecConfig.MCP
	if mcpCfg.ServerID == "" || mcpCfg.ToolName == "" {
		return fmt.Errorf("mcp server_id and tool_name are required")
	}
	if err := e.client.HealthCheck(ctx, mcpCfg.ServerID); err != nil {
		return fmt.Errorf("mcp server health check failed: %w", err)
	}
	tools, err := e.client.ListTools(ctx, mcpCfg.ServerID)
	if err != nil {
		return fmt.Errorf("failed to list mcp tools: %w", err)
	}
	for i := range tools {
		if tools[i].Name == mcpCfg.ToolName {
			return nil
		}
	}
	return fmt.Errorf("mcp tool %s not found on server %s", mcpCfg.ToolName, mcpCfg.ServerID)
}
