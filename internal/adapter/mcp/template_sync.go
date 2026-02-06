package mcp

import (
	"fmt"
	"strings"

	"goyavision/internal/domain/operator"
	"goyavision/internal/port"
)

// ToolToTemplate 将 MCP Tool 映射为 OperatorTemplate（最小适配层）。
// 说明：当前项目仍使用 StaticClient，此映射用于统一同步逻辑与后续真实 MCP 适配替换。
func ToolToTemplate(serverID string, tool port.MCPTool) *operator.OperatorTemplate {
	code := fmt.Sprintf("mcp_%s_%s", sanitize(serverID), sanitize(tool.Name))
	return &operator.OperatorTemplate{
		Code:        code,
		Name:        tool.Name,
		Description: tool.Description,
		Category:    operator.CategoryUtility,
		Type:        operator.TypeTranscode,
		ExecMode:    operator.ExecModeMCP,
		ExecConfig: &operator.ExecConfig{MCP: &operator.MCPExecConfig{
			ServerID:    serverID,
			ToolName:    tool.Name,
			ToolVersion: tool.Version,
		}},
		InputSchema: tool.InputSchema,
		OutputSpec:  tool.OutputSchema,
		Config:      map[string]interface{}{},
		Author:      "mcp",
		Tags:        []string{"mcp", serverID},
	}
}

func sanitize(v string) string {
	v = strings.TrimSpace(strings.ToLower(v))
	v = strings.ReplaceAll(v, " ", "_")
	v = strings.ReplaceAll(v, "/", "_")
	v = strings.ReplaceAll(v, "-", "_")
	return v
}
