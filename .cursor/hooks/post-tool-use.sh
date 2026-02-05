#!/usr/bin/env bash
# PostToolUse Hook - 在工具使用后检查性能
# 符合 Cursor Hooks 官方规范

set -euo pipefail

# 从 stdin 读取 JSON 输入
input_json=$(cat)

# 解析工具信息（根据 Cursor 官方格式）
tool_name=$(echo "$input_json" | jq -r '.tool_name // ""')
duration=$(echo "$input_json" | jq -r '.duration // 0')

# postToolUse 的输出格式：可以返回 updated_mcp_tool_output（仅针对 MCP 工具）
# 对于性能监控，这里只返回空对象（不影响执行）
if (( duration > 5000 )); then
    # 性能提醒可以通过日志记录，但不影响工具执行
    echo '{}'
    exit 0
fi

echo '{}'
exit 0
