#!/usr/bin/env bash
# GoyaVision 完成开发提醒 Hook
# 在 Cursor Agent 任务结束（stop）时自动触发检查清单提醒
# 路径：.cursor/hooks/finish-dev-reminder.sh
#
# 符合 Cursor Hooks 官方规范：
# - 从 stdin 读取 JSON 输入
# - 输出 JSON 格式到 stdout
# - 使用 followup_message 自动触发后续消息

set -euo pipefail

# 从 stdin 读取 JSON 输入
input_json=$(cat)

# 检查是否有 jq 命令
if command -v jq >/dev/null 2>&1; then
    # 使用 jq 解析 loop_count
    loop_count=$(echo "$input_json" | jq -r '.loop_count // 0' 2>/dev/null || echo "0")
else
    # 如果没有 jq，使用 grep 和 sed 简单解析（fallback）
    loop_count=$(echo "$input_json" | grep -o '"loop_count"[[:space:]]*:[[:space:]]*[0-9]*' | grep -o '[0-9]*' || echo "0")
    if [ -z "$loop_count" ]; then
        loop_count="0"
    fi
fi

# 检查 loop_count 是否超过限制（系统限制为 5）
if [ "$loop_count" -ge 5 ]; then
    # 超过限制，不输出 followup_message
    echo '{}'
    exit 0
fi

# 构建检查清单消息
checklist_message="请完成以下开发后检查清单：

1. **更新开发进度**
   - 文件: docs/development-progress.md
   - 操作: 更新功能状态（✅/🚧/⏸️）与说明

2. **更新变更日志**
   - 文件: CHANGELOG.md
   - 操作: 在 [未发布] 下按类型添加条目

3. **按需更新其他文档**
   - API 变更 -> docs/api.md
   - 需求/架构变更 -> docs/requirements.md, docs/architecture.md
   - 用户/部署影响 -> README.md, docs/DEPLOYMENT.md

4. **Git 提交**
   - 格式: <type>(<scope>): <subject>
   - 示例: feat(asset): 实现媒体资产管理
   - 自检: 已测试、已格式化、文档已更新

详细步骤见: .cursor/skills/development-workflow/SKILL.md
规则说明: .cursor/rules/development-workflow.mdc"

# 输出 JSON 格式的响应
if command -v jq >/dev/null 2>&1; then
    # 使用 jq 确保 JSON 格式正确，并转义特殊字符
    response=$(jq -n \
        --arg msg "$checklist_message" \
        '{followup_message: $msg}')
    echo "$response"
else
    # 如果没有 jq，手动构建 JSON（转义特殊字符）
    # 转义双引号、反斜杠和换行符
    escaped_msg=$(echo "$checklist_message" | sed 's/\\/\\\\/g' | sed 's/"/\\"/g' | sed ':a;N;$!ba;s/\n/\\n/g')
    echo "{\"followup_message\":\"$escaped_msg\"}"
fi

exit 0
