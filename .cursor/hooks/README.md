# GoyaVision Hooks

本目录包含 GoyaVision 项目的 Cursor Hooks 脚本。

## 当前 Hooks

### finish-dev-reminder.sh

**用途**：在 Agent 任务结束（stop）时自动触发完成开发检查清单提醒。

**实现规范**：
- ✅ 符合 Cursor Hooks 官方规范
- ✅ 从 stdin 读取 JSON 输入
- ✅ 输出 JSON 格式到 stdout
- ✅ 使用 `followup_message` 自动触发后续消息
- ✅ 检查 `loop_count` 限制（防止无限循环）
- ✅ 支持 jq 和 fallback 解析（提高兼容性）

**输入格式**：
```json
{
  "status": "completed" | "aborted" | "error",
  "loop_count": 0,
  "conversation_id": "string",
  "generation_id": "string",
  "model": "string",
  ...
}
```

**输出格式**：
```json
{
  "followup_message": "检查清单消息"
}
```

当 `loop_count >= 5` 时，输出空 JSON `{}`（系统限制）。

## 测试方法

### 测试脚本输出格式

```bash
# 测试正常情况（loop_count < 5）
echo '{"status":"completed","loop_count":0,"conversation_id":"test"}' | \
  bash .cursor/hooks/finish-dev-reminder.sh | jq .

# 测试超过限制（loop_count >= 5）
echo '{"status":"completed","loop_count":6,"conversation_id":"test"}' | \
  bash .cursor/hooks/finish-dev-reminder.sh | jq .

# 测试语法
bash -n .cursor/hooks/finish-dev-reminder.sh
```

### 在 Cursor 中测试

1. 启动 Cursor Agent（Cmd+K 或 Agent Chat）
2. 完成一个任务
3. Agent 任务结束时，会自动触发 `stop` hook
4. 检查清单消息会自动显示在聊天中

## 配置

Hooks 配置在 `.cursor/hooks.json`：

```json
{
  "version": 1,
  "hooks": {
    "stop": [
      {
        "command": ".cursor/hooks/finish-dev-reminder.sh",
        "timeout": 10,
        "loop_limit": 5
      }
    ]
  }
}
```

**配置说明**：
- `command`: 脚本路径（项目级 hooks 使用项目根目录相对路径）
- `timeout`: 执行超时时间（秒），默认 10 秒
- `loop_limit`: stop hook 的单脚本循环上限，默认 5

## 开发新 Hook

创建新 hook 时，请遵循以下规范：

1. **脚本格式**：
   ```bash
   #!/usr/bin/env bash
   set -euo pipefail
   
   # 从 stdin 读取 JSON 输入
   input_json=$(cat)
   
   # 处理逻辑...
   
   # 输出 JSON 到 stdout
   echo '{"key": "value"}'
   exit 0
   ```

2. **退出码**：
   - `0`: Hook 执行成功
   - `2`: 阻止操作（等同于返回 `permission: "deny"`）
   - 其他: Hook 失败，操作继续（fail-open）

3. **JSON 处理**：
   - 优先使用 `jq` 解析和生成 JSON
   - 提供 fallback 方案（使用 grep/sed）以提高兼容性

4. **添加到 hooks.json**：
   ```json
   {
     "hooks": {
       "hookName": [
         {
           "command": ".cursor/hooks/your-script.sh",
           "timeout": 10
         }
       ]
     }
   }
   ```

5. **设置可执行权限**：
   ```bash
   chmod +x .cursor/hooks/your-script.sh
   ```

## 参考文档

- [Cursor Hooks 官方文档](https://cursor.com/cn/docs/agent/hooks)
- [第三方 Hooks 兼容性](https://cursor.com/cn/docs/agent/third-party-hooks)
- [项目 Hooks 配置](../hooks.json)
