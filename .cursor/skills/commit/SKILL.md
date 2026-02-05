---
name: commit
description: Git 提交规范指导（Conventional Commits）。用于提交前检查与撰写提交信息。
---

# Git 提交规范

用于按 Conventional Commits 规范完成提交。

## 何时使用

✅ **推荐场景**：
- 功能/修复完成准备提交
- 需要检查提交信息格式

❌ **不适用场景**：
- 只需要查看 Git 历史（使用 git log）
- 需要合并分支（使用 git merge）

## 提交格式

```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

> **提交信息需使用中文描述**

## Type 类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档变更
- `refactor`: 代码重构（不改变功能）
- `test`: 测试相关
- `chore`: 构建、配置、依赖更新
- `perf`: 性能优化
- `style`: 代码格式（不影响逻辑）

## Scope 范围（可选）

- `asset`: 资产库
- `operator`: 算子中心
- `workflow`: 工作流
- `task`: 任务管理
- `auth`: 认证授权
- `api`: API 层
- `ui`: 前端界面
- `cursor`: Cursor 配置
- `config`: 配置与运维

## 提交示例

### ✅ 新功能
```bash
git commit -m "feat(asset): 实现媒体资产管理功能

- 添加 MediaAsset 实体和 Repository
- 实现 CQRS 命令和查询（create/update/delete/get/list）
- 添加 Asset Handler 和 DTO
- 更新 API 文档和开发进度"
```

### ✅ Bug 修复
```bash
git commit -m "fix(workflow): 修复 DAG 验证时的死循环问题

在验证包含自环的 DAG 时会立即返回错误，
避免进入无限循环。

Closes #123"
```

### ✅ 文档更新
```bash
git commit -m "docs: 更新开发进度、CHANGELOG 与 API 文档"
```

## 提交前检查

- [ ] **代码已测试** - 单元测试、集成测试、手动测试已通过
- [ ] **文档已同步** - 开发进度、CHANGELOG、API 文档已更新
- [ ] **代码已格式化** - Go: `go fmt ./...`，前端: `pnpm run format`
- [ ] **无调试代码** - 移除 console.log、fmt.Println、注释代码块
- [ ] **提交信息符合规范** - 有正确的 type 和 scope，Subject 清晰简洁

## 好提交 vs 坏提交

### ❌ 坏提交（不要这样做）
```bash
git commit -m "update"
git commit -m "fix bug"
git commit -m "完成功能"
```

### ✅ 好提交（推荐做法）
```bash
git commit -m "feat(asset): 实现媒体资产管理功能"
git commit -m "fix(workflow): 修复 DAG 验证死循环"
git commit -m "docs: 更新 V1.0 架构文档"
```
