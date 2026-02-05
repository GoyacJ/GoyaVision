#!/usr/bin/env bash
# BeforeSubmitPrompt Hook - 在提交 prompt 前注入上下文
# 符合 Cursor Hooks 官方规范

set -euo pipefail

# 从 stdin 读取 JSON 输入
input_json=$(cat)

# 解析 prompt 内容
prompt=$(echo "$input_json" | jq -r '.prompt // ""')

if [[ -z "$prompt" ]]; then
    echo '{"continue": true}'
    exit 0
fi

# 根据 prompt 内容注入相关上下文提醒
context=""

if echo "$prompt" | grep -qiE "API|接口|endpoint"; then
    context="修改/新增 API 时需更新 docs/api.md，并保持 /api/v1 前缀与统一错误响应。"
elif echo "$prompt" | grep -qiE "前端|Vue|组件|UI"; then
    context="前端使用 Vue 3 + TS + Tailwind，优先复用 Gv 组件库与 design tokens。"
elif echo "$prompt" | grep -qiE "测试|test"; then
    context="测试需覆盖核心逻辑，命名清晰，并避免依赖外部资源。"
fi

# 输出响应（beforeSubmitPrompt 不支持 contextModification，只能阻止或允许）
if [[ -n "$context" ]]; then
    # 注意：beforeSubmitPrompt 的输出格式是 {continue: boolean, user_message?: string}
    # 这里我们允许继续，但可以通过 user_message 提醒
    echo "{\"continue\": true, \"user_message\": \"$context\"}"
else
    echo '{"continue": true}'
fi

exit 0
