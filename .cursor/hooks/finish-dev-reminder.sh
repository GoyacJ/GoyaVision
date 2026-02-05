#!/usr/bin/env bash
# GoyaVision 完成开发提醒 Hook
# 在 Cursor Agent 任务结束（stop）时输出检查清单提示
# 路径：.cursor/hooks/finish-dev-reminder.sh
#
# 符合 Cursor Hooks 官方规范：
# - 从 stdin 读取 JSON 输入
# - 输出 JSON 格式到 stdout
# - 只在有未提交变更时输出提示信息（不自动发送请求）

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

# 检查是否有未提交的变更（避免无限循环）
# 获取项目根目录（假设脚本在 .cursor/hooks/ 目录下）
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
project_root="$(cd "$script_dir/../.." && pwd)"

# 检查 git 状态（如果不在 git 仓库中，跳过检查）
if [ -d "$project_root/.git" ]; then
    cd "$project_root"
    # 检查是否有未暂存或未提交的变更
    if ! git diff --quiet --exit-code 2>/dev/null && ! git diff --cached --quiet --exit-code 2>/dev/null; then
        has_changes=true
    else
        has_changes=false
    fi
else
    # 不在 git 仓库中，默认不触发（避免频繁提醒）
    has_changes=false
fi

# 只有在有未提交变更时才输出提示
if [ "$has_changes" = false ]; then
    # 没有变更，不输出提示
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

# 输出提示信息到 stderr（这样可以在日志中看到，但不会自动发送请求）
echo "$checklist_message" >&2

# 输出空 JSON（不包含 followup_message，因此不会自动发送请求）
echo '{}'

exit 0
