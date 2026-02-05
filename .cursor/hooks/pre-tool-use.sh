#!/usr/bin/env bash
# PreToolUse Hook - 在工具使用前检查 Domain 层依赖
# 符合 Cursor Hooks 官方规范

set -euo pipefail

# 从 stdin 读取 JSON 输入
input_json=$(cat)

# 解析工具名称和输入（根据 Cursor 官方格式）
tool_name=$(echo "$input_json" | jq -r '.tool_name // ""')
tool_input=$(echo "$input_json" | jq -r '.tool_input // {}')

# 只检查 Write 工具（文件写入）
if [[ "$tool_name" != "Write" ]]; then
    echo '{"decision": "allow"}'
    exit 0
fi

# 检查文件路径（Write 工具的输入格式）
file_path=$(echo "$tool_input" | jq -r '.path // .file_path // ""')

if [[ -z "$file_path" ]]; then
    echo '{"decision": "allow"}'
    exit 0
fi

# 如果写入 Domain 层文件，添加提醒（通过 reason 字段）
if [[ "$file_path" == internal/domain/* ]]; then
    reason="Domain 层不得依赖 adapter/api，确保纯业务模型。"
    echo "{\"decision\": \"allow\", \"reason\": \"$reason\"}"
    exit 0
fi

echo '{"decision": "allow"}'
exit 0
