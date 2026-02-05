# GoyaVision Cursor 配置

本目录包含 GoyaVision 项目的 Cursor IDE 配置，包括 Rules、Skills、Commands 和 Hooks。

## 目录结构

```
.cursor/
├── rules/              # 项目规则（Project Rules）
│   ├── goyavision.mdc              # 核心项目规则（始终应用）
│   ├── development-workflow.mdc    # 开发工作流规则（始终应用）
│   ├── backend-domain.mdc          # 后端领域与端口层规则（Domain/Port 文件）
│   ├── backend-app.mdc             # 后端应用层规则（App 文件）
│   ├── backend-adapter-api.mdc     # 后端适配器与 API 层规则（Adapter/API 文件）
│   ├── frontend-components.mdc     # 前端组件规范（前端文件）
│   ├── testing.mdc                 # 测试规则（测试文件）
│   ├── docs.mdc                    # 文档规则（文档文件）
│   └── config-ops.mdc              # 配置与运维规则（配置文件）
├── skills/             # Agent Skills
│   ├── development-workflow/       # 开发工作流技能
│   ├── goyavision-context/         # 项目上下文技能
│   ├── frontend-components/        # 前端组件开发技能
│   ├── api-doc/                    # API 文档更新技能
│   ├── commit/                     # Git 提交规范技能
│   └── progress/                    # 开发进度更新技能
├── commands/           # 自定义命令（Slash Commands）
│   ├── dev-start.md    # 开始开发检查清单
│   ├── dev-done.md     # 完成开发检查清单
│   ├── commit.md       # Git 提交规范
│   ├── context.md      # 项目上下文
│   ├── api-doc.md      # API 文档更新指南
│   ├── progress.md     # 开发进度更新指南
│   └── frontend-component.md       # 前端组件开发流程
├── hooks/              # Hooks 脚本
│   ├── finish-dev-reminder.sh      # 完成开发提醒脚本
│   ├── pre-tool-use.sh            # 工具使用前检查
│   ├── post-tool-use.sh           # 工具使用后检查
│   └── before-submit-prompt.sh    # 提交 prompt 前注入上下文
├── hooks.json          # Hooks 配置
└── README.md           # 本文件
```

## Rules（规则）

Rules 提供系统级指令，指导 AI Agent 如何理解和编写代码。

### 规则类型

| 规则文件 | 应用方式 | 说明 |
|---------|---------|------|
| `goyavision.mdc` | Always Apply | 核心项目规则，每个会话都应用 |
| `development-workflow.mdc` | Always Apply | 开发工作流规范，每个会话都应用 |
| `backend-domain.mdc` | Apply to Specific Files | Domain/Port 层文件（`internal/domain/**`, `internal/port/**`） |
| `backend-app.mdc` | Apply to Specific Files | App 层文件（`internal/app/**`） |
| `backend-adapter-api.mdc` | Apply to Specific Files | Adapter/API 层文件（`internal/adapter/**`, `internal/api/**`） |
| `frontend-components.mdc` | Apply to Specific Files | 前端文件（`web/**/*.vue`, `web/**/*.ts`） |
| `testing.mdc` | Apply to Specific Files | 测试文件（`**/*_test.go`, `**/*.test.ts`） |
| `docs.mdc` | Apply to Specific Files | 文档文件（`docs/**`, `**/*.md`） |
| `config-ops.mdc` | Apply to Specific Files | 配置文件（`configs/**`, `Makefile`, `docker-compose.yml`） |

### 规则格式

每个规则文件使用 Markdown 格式，包含 YAML frontmatter：

```markdown
---
description: "规则描述"
alwaysApply: true | false
globs:
  - "pattern/**/*.ext"
---
```

## Skills（技能）

Skills 是可移植的知识包，Agent 可以根据上下文自动调用。

### 可用技能

| 技能 | 描述 | 使用时机 |
|------|------|----------|
| `development-workflow` | 开发工作流管理 | 开始/完成开发时 |
| `goyavision-context` | 项目架构和上下文 | 需要了解项目结构、API、实体时 |
| `frontend-components` | 前端组件开发指南 | 新增/修改 Vue 组件、页面、样式时 |
| `api-doc` | API 文档更新指南 | 新增或修改 API 端点时 |
| `commit` | Git 提交规范指导 | 提交前检查与撰写提交信息时 |
| `progress` | 开发进度更新指南 | 更新开发进度文档时 |

### 技能格式

每个技能位于独立目录，包含 `SKILL.md` 文件：

```markdown
---
name: skill-name
description: "技能描述"
---
```

## Commands（命令）

Commands 是可通过 `/` 前缀触发的自定义工作流。

### 可用命令

| 命令 | 描述 | 使用示例 |
|------|------|----------|
| `/dev-start` | 开始开发前检查清单 | `/dev-start` |
| `/dev-done` | 完成开发后检查清单 | `/dev-done` |
| `/commit` | Git 提交规范指南 | `/commit` |
| `/context` | 项目上下文信息 | `/context` |
| `/api-doc` | API 文档更新指南 | `/api-doc` |
| `/progress` | 开发进度更新指南 | `/progress` |
| `/frontend-component` | 前端组件开发流程 | `/frontend-component` |

### 命令格式

命令文件使用 Markdown 格式，无需 frontmatter：

```markdown
# command-name

命令描述和使用说明...
```

## Hooks（钩子）

Hooks 允许在 Agent 循环的特定阶段执行自定义脚本，观察、控制和扩展 agent 行为。

### 当前配置

- **preToolUse hook**: 在工具使用前检查 Domain 层依赖规则
- **postToolUse hook**: 在工具使用后检查性能（执行时间超过 5 秒时提醒）
- **beforeSubmitPrompt hook**: 在提交 prompt 前根据内容注入相关上下文提醒
- **stop hook**: 在 Agent 任务结束时触发 `finish-dev-reminder.sh`，自动显示完成开发检查清单

### Hooks 配置

`hooks.json` 文件定义 hooks：

```json
{
  "version": 1,
  "hooks": {
    "preToolUse": [
      {
        "command": ".cursor/hooks/pre-tool-use.sh",
        "timeout": 5,
        "matcher": "Write"
      }
    ],
    "postToolUse": [
      {
        "command": ".cursor/hooks/post-tool-use.sh",
        "timeout": 5
      }
    ],
    "beforeSubmitPrompt": [
      {
        "command": ".cursor/hooks/before-submit-prompt.sh",
        "timeout": 5
      }
    ],
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
- `loop_limit`: stop hook 的单脚本循环上限，默认 5（防止无限循环）

### Hook 实现规范

所有 hook 脚本必须符合官方规范：

1. **输入格式**：从 stdin 读取 JSON 输入
2. **输出格式**：输出 JSON 到 stdout
3. **退出码**：
   - `0`: Hook 执行成功
   - `2`: 阻止操作（等同于返回 `permission: "deny"`）
   - 其他: Hook 失败，操作继续（fail-open）

**stop hook 示例**：
```bash
#!/usr/bin/env bash
# 从 stdin 读取 JSON 输入
input_json=$(cat)

# 解析 loop_count
loop_count=$(echo "$input_json" | jq -r '.loop_count // 0')

# 输出 JSON 响应
echo '{"followup_message": "检查清单消息"}'
exit 0
```

### 支持的 Hook 事件

**Agent（Cmd+K/Agent Chat）**：
- `sessionStart` / `sessionEnd` - 会话生命周期
- `preToolUse` / `postToolUse` / `postToolUseFailure` - 工具使用前后
- `subagentStart` / `subagentStop` - Subagent 生命周期
- `beforeShellExecution` / `afterShellExecution` - Shell 命令执行前后
- `beforeMCPExecution` / `afterMCPExecution` - MCP 工具执行前后
- `beforeReadFile` / `afterFileEdit` - 文件读取/编辑前后
- `beforeSubmitPrompt` - 提交前校验 prompt
- `preCompact` - 上下文窗口压缩前
- `stop` - Agent 任务结束
- `afterAgentResponse` / `afterAgentThought` - Agent 响应跟踪

**Tab（行内补全）**：
- `beforeTabFileRead` - Tab 读取文件前
- `afterTabFileEdit` - Tab 编辑文件后

## 使用指南

### 开始新功能开发

1. 在聊天中输入 `/dev-start` 查看开发前检查清单
2. 按顺序阅读文档（开发进度 → 变更日志 → 需求 → 架构 → API）
3. 遵循 Rules 中的架构约束和代码规范
4. 实现功能...

### 完成功能开发

1. 在聊天中输入 `/dev-done` 查看完成后检查清单
2. 更新开发进度文档
3. 更新 CHANGELOG.md
4. 按需更新 API/架构文档
5. 使用 `/commit` 创建规范的 Git 提交

### 查询项目信息

- 使用 `/context` 获取项目架构、API 端点、实体定义
- 使用 `/api-doc` 查看 API 文档更新指南
- 使用 `/progress` 查看开发进度更新指南

## 与 Claude Code 的兼容性

本项目同时支持 Cursor 和 Claude Code：

- **Cursor**: 使用 `.cursor/` 目录下的配置
- **Claude Code**: 使用 `.claude/commands/` 目录下的命令（兼容性保留）

两个工具共享相同的开发规范和流程。

## 参考文档

- [Cursor Rules 文档](https://cursor.com/cn/docs/context/rules)
- [Cursor Skills 文档](https://cursor.com/cn/docs/context/skills)
- [Cursor Commands 文档](https://cursor.com/cn/docs/context/commands)
- [Cursor Hooks 文档](https://cursor.com/cn/docs/agent/hooks)
- [项目开发规范](../.cursor/rules/development-workflow.mdc)

## 更新日志

- **2024-01-XX**: 初始配置，符合 Cursor 官方规范
  - 添加 Rules（goyavision, development-workflow, frontend-components）
  - 添加 Skills（development-workflow, goyavision-context）
  - 添加 Commands（dev-start, dev-done, commit, context, api-doc, progress）
  - 配置 Hooks（stop hook 用于完成开发提醒）

- **2024-02-06**: 完善配置，参考 .clinerules/ 和 .cline/ 补充内容
  - 新增 Rules：backend-domain, backend-app, backend-adapter-api, testing, docs, config-ops
  - 新增 Skills：frontend-components, api-doc, commit, progress
  - 新增 Hooks：preToolUse, postToolUse, beforeSubmitPrompt
  - 新增 Commands：frontend-component
  - 更新 goyavision.mdc：添加信息完整性与提问规范、通用代码质量
  - 更新 development-workflow.mdc：引用新增的规则文件
