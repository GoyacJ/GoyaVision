---
name: development-workflow
description: 管理 GoyaVision 开发全流程（开始前查阅文档、完成后更新文档与提交）。用于开始开发、完成开发、提交前检查。
---

# GoyaVision 开发工作流

统一执行开发前与开发后的检查流程，确保遵循项目文档、架构与提交规范。

## 何时使用
- 开始新功能或修复 Bug 前："开始开发"、"新需求"。
- 完成开发后："完成开发"、"收尾"、"提交前检查"。

## 阶段一：开始开发（必做）

### 1. 必读文档（按顺序）
1. `docs/development-progress.md`（状态 ✅/🚧/⏸️）
2. `CHANGELOG.md`（关注 [未发布]）
3. `docs/requirements.md`
4. `docs/architecture.md`
5. `docs/api.md`

### 2. 必遵循规则
- `.clinerules/00-universal.md`
- 相关条件规则（后端/前端/测试/文档）

### 3. 引用项目上下文
需要时启用 `/goyavision-context` 获取实体、分层、API 端点与状态。

## 阶段二：完成开发（必做）

### 1. 更新开发进度
- 编辑 `docs/development-progress.md` 更新状态与说明。

### 2. 更新变更日志
- 在 `CHANGELOG.md` 的 [未发布] 下按类型记录。

### 3. 同步其他文档
- API 变更：更新 `docs/api.md`。
- 架构/需求变更：更新 `docs/architecture.md` / `docs/requirements.md`。

### 4. 规范提交
- Conventional Commits：`<type>(<scope>): <subject>`。
- 提交前自检：格式化、无调试代码。
