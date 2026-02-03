---
name: development-workflow
description: GoyaVision 开发工作流：开始开发前查阅文档与进度，完成开发后更新文档并 Git 提交。可通过 @development-workflow 或「开始开发」「完成开发」触发。
---

# GoyaVision 开发工作流

## 何时使用

- 用户说「开始开发」「新需求」「接一个新需求」时：执行「开始开发」清单。
- 用户说「完成开发」「收尾」「提交前检查」时：执行「完成开发」清单。
- 需要一次性对照「开始前必读」或「完成后必做」时：引用本 Skill 并按步骤执行。

---

## 一、开始开发（新需求前必做）

在开始任何新需求或功能开发前，按顺序完成以下步骤以熟悉项目。

### 1. 必读文档（按顺序阅读或引用）

| 文档 | 路径 | 用途 |
|------|------|------|
| 开发进度 | `docs/development-progress.md` | 当前功能状态、路线图、进行中/待开始项 |
| 变更日志 | `CHANGELOG.md` | 最近变更、未发布内容，避免重复或冲突 |
| 需求文档 | `docs/requirements.md` | 功能需求与验收标准 |
| 架构文档 | `docs/architecture.md` | 分层、依赖、领域概念、数据模型 |
| API 文档 | `docs/api.md` | 已有接口、请求/响应约定 |

### 2. 必用 Rules 与 Skills

- **Rules**：确保已加载并遵循 `.cursor/rules/goyavision.mdc`、`.cursor/rules/development-workflow.mdc`；涉及前端时遵循 `.cursor/rules/frontend-components.mdc`。
- **Skills**：实现或评审功能时使用 `@goyavision-context`（`.cursor/skills/goyavision-context/SKILL.md`）以获取项目结构、实体、API 与开发状态。

### 3. 开发中依据

- 实现必须符合 `docs/architecture.md` 的分层与依赖。
- API 设计须与 `docs/api.md` 一致，并在实现后更新该文档。
- 代码风格、错误处理、Git 提交格式见 `goyavision.mdc`。

---

## 二、完成开发（功能/修复完成后必做）

每次完成功能开发或问题修复后，按顺序完成以下步骤，最后执行 Git 提交。

### 1. 更新开发进度

- **文件**：`docs/development-progress.md`
- **操作**：
  - 将本次涉及的功能/模块状态更新为 ✅ 已完成、🚧 进行中或 ⏸️ 待开始。
  - 更新「当前迭代」或「说明」列中与本需求相关的描述。
  - 如有技术债务或阻塞问题，在文档中注明。

### 2. 更新变更日志

- **文件**：`CHANGELOG.md`
- **操作**：
  - 在 `[未发布]` 下按类型添加条目：新增、变更、修复、弃用、移除、安全。
  - 描述清晰，便于他人理解；若有 Issue/PR 可注明编号。

### 3. 按需更新其他项目文档

| 变更类型 | 需更新的文档 |
|----------|--------------|
| 新增或修改 API | `docs/api.md`（请求/响应示例、参数、错误码） |
| 需求或验收标准变更 | `docs/requirements.md` |
| 架构或分层变更 | `docs/architecture.md` |
| 影响用户使用、安装或配置 | `README.md`、`docs/DEPLOYMENT.md` |

### 4. Git 提交

- **范围**：将本次所有代码与文档变更一并纳入同一提交（或按逻辑拆分为多个符合规范的提交）。
- **格式**：遵循 Conventional Commits，例如：
  - `feat(asset): 实现媒体资产管理功能`
  - `fix(workflow): 修复 DAG 验证死循环`
  - `docs: 更新开发进度与 CHANGELOG`
- **提交前自检**：
  - [ ] 代码已测试（单元/集成或手动）
  - [ ] 相关文档已更新（开发进度、CHANGELOG、必要时 API/架构/README）
  - [ ] 代码已格式化（gofmt / goimports，前端按项目规范）
  - [ ] 无调试代码或临时注释
  - [ ] Commit message 符合规范

---

## 三、与 Cursor 的配合

- **Rules**：`.cursor/rules/development-workflow.mdc` 已对上述流程做概括，并设为 alwaysApply，Agent 会默认遵循。
- **Hooks**：任务结束（`stop`）时，`.cursor/hooks.json` 可触发提醒脚本，输出本「完成开发」清单，便于人工或 Agent 核对。
- **本 Skill**：通过 @development-workflow 或自然语言「开始开发」「完成开发」引用，可一次性获得完整步骤与文件路径，便于逐项执行。
