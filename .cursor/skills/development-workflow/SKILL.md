---
name: development-workflow
description: 开发工作流管理 - 开始开发前查阅文档，完成后更新文档并规范提交
---

# GoyaVision 开发工作流

管理 GoyaVision 项目的完整开发生命周期，确保每次开发都遵循统一流程：开始前熟悉项目状态，完成后同步文档和规范提交。

## 何时使用

✅ **推荐场景**：
- 开始新功能开发或 Bug 修复前（触发词："开始开发"、"新需求"）
- 完成功能开发或修复后（触发词："完成开发"、"提交前检查"、"收尾"）
- 需要一次性查看开发前/后的完整检查清单

❌ **不适用场景**：
- 只需要查看项目架构（使用 `/goyavision-context`）
- 只是阅读代码或文档（使用 Read 工具）

## 工作流阶段

```
📋 开始开发 → 💻 实现功能 → ✅ 完成开发
   (查阅文档)   (遵循规范)   (更新文档+提交)
```

---

## 阶段一：开始开发（新需求前必做）

在开始任何新需求或功能开发前，按顺序完成以下步骤。

### 1. 必读文档（按顺序查阅）

| 序号 | 文档 | 路径 | 用途 | 关注重点 |
|------|------|------|------|----------|
| 1️⃣ | 开发进度 | `docs/development-progress.md` | 当前功能状态 | ✅/🚧/⏸️ 标记，避免重复工作 |
| 2️⃣ | 变更日志 | `CHANGELOG.md` | 最近变更 | **[未发布]** 章节，了解最新修改 |
| 3️⃣ | 需求文档 | `docs/requirements.md` | 功能需求 | 验收标准、业务规则 |
| 4️⃣ | 架构文档 | `docs/architecture.md` | 系统设计 | 分层架构、依赖规则、数据流 |
| 5️⃣ | API 文档 | `docs/api.md` | 接口约定 | 已有端点、请求/响应格式 |

**为什么按这个顺序？**
- 开发进度 → 避免重复或冲突
- 变更日志 → 了解最新改动
- 需求/架构/API → 理解实现方式

### 2. 必遵循的规则

加载并遵循以下 Rules（位于 `.cursor/rules/`）：

- ✅ **`goyavision.mdc`** - 核心规则
  - 分层架构与依赖规则
  - 算子标准协议
  - 错误处理约定
  - 代码风格（无行尾注释）

- ✅ **`development-workflow.mdc`** - 工作流规则
  - 开发前/后要求
  - 文档更新规范
  - Git 提交格式（Conventional Commits）

- ✅ **`frontend-components.mdc`** - 前端规则（如涉及前端）
  - 组件结构与命名
  - Tailwind + Element Plus 使用
  - Pinia 状态管理

### 3. 引用项目上下文

需要时使用 `/goyavision-context` 获取：
- 核心概念（MediaSource → MediaAsset → Operator → Workflow → Task → Artifact）
- 分层架构（Domain/Port/App/Adapter/API）
- API 端点完整列表
- 算子标准协议
- 开发状态（✅/🚧/⏸️）

### 4. 开发中的依据

**架构遵循**：
- 实现必须符合 Clean Architecture 分层
- Domain 层无外部依赖
- App 层通过 Port 接口调用，禁止直接依赖 Adapter
- API 层使用 DTO，禁止暴露 Domain 实体

**API 设计**：
- 遵循 RESTful 约定
- 前缀统一为 `/api/v1`
- 实现后立即更新 `docs/api.md`

**代码风格**：
- Go: `gofmt` + `goimports`，无行尾注释
- 前端: TypeScript strict，Vue 3 Composition API
- 错误处理: 使用 `internal/api/errors.go` 定义的错误类型

### 5. 快速参考

**默认凭证**: admin / admin123
**服务端口**: 8080 (API), 5432 (DB), 8554 (RTSP), 8888 (HLS)
**构建命令**: `make build` (后端), `make build-web` (前端), `make build-all` (全部)
**配置文件**: `configs/config.<env>.yaml`，环境变量前缀 `GOYAVISION_*`

---

## 阶段二：完成开发（功能/修复完成后必做）

完成功能开发或问题修复后，**按顺序**完成以下 4 个步骤。

### 步骤 1️⃣：更新开发进度

**文件**: `docs/development-progress.md`

**操作清单**:
- [ ] 将本次涉及的功能状态更新为 ✅ 已完成、🚧 进行中或 ⏸️ 待开始
- [ ] 更新"当前迭代"或"说明"列的描述
- [ ] 如有技术债务或阻塞问题，在文档中注明
- [ ] 更新完成百分比（如有跟踪）

**示例**:
```markdown
| 功能 | 状态 | 说明 |
|------|------|------|
| 媒体资产管理 | ✅ | 已完成 CRUD、搜索、派生追踪功能 |
| 算子管理 | 🚧 | 进行中 - CRUD 已完成，版本管理开发中 |
```

### 步骤 2️⃣：更新变更日志

**文件**: `CHANGELOG.md`

**操作清单**:
- [ ] 在 `[未发布]` 章节下添加条目
- [ ] 按类型分类：**新增**/变更/修复/弃用/移除/安全
- [ ] 描述清晰、面向用户，便于他人理解
- [ ] 包含 Issue/PR 编号（如适用）

**示例**:
```markdown
## [未发布]

### 新增
- **媒体资产管理**: 实现 CRUD、标签系统、派生追踪功能
- **资产搜索**: 支持按类型、来源、标签过滤

### 修复
- 修复工作流 DAG 验证时的死循环问题 (#123)
- 修复媒体源状态查询超时问题

### 变更
- 优化媒体资产列表加载性能（从 2s 降至 300ms）
```

### 步骤 3️⃣：按需更新相关文档

根据变更类型，判断是否需要更新以下文档：

| 变更类型 | 需更新的文档 | 更新内容 |
|----------|--------------|----------|
| **新增或修改 API** | `docs/api.md` | 端点文档、请求/响应示例、参数说明、错误码 |
| **架构变更** | `docs/architecture.md` | 新增层/组件、依赖图、设计决策 |
| **需求变更** | `docs/requirements.md` | 功能规格、验收标准、用例 |
| **用户影响** | `README.md` | 新功能说明、使用方式 |
| **部署/配置** | `docs/DEPLOYMENT.md` | 配置项、环境变量、安装步骤 |

### 步骤 4️⃣：Git 提交

**提交格式**: 遵循 [Conventional Commits](https://www.conventionalcommits.org/)

```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

**Type 类型**:
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档变更
- `refactor`: 代码重构（不改变功能）
- `test`: 测试相关
- `chore`: 构建、配置、依赖更新
- `perf`: 性能优化
- `style`: 代码格式（不影响逻辑）

**Scope 范围**（可选）:
- `asset`, `operator`, `workflow`, `task`, `auth`, `api`, `ui`

**提交示例**:

```bash
# ✅ 新功能
git commit -m "feat(asset): 实现媒体资产管理功能

- 添加 MediaAsset 实体和 Repository
- 实现 CQRS 命令和查询（create/update/delete/get/list）
- 添加 Asset Handler 和 DTO
- 更新 API 文档和开发进度"

# ✅ Bug 修复
git commit -m "fix(workflow): 修复 DAG 验证时的死循环问题

在验证包含自环的 DAG 时会立即返回错误，
避免进入无限循环。

Closes #123"

# ✅ 文档更新
git commit -m "docs: 更新开发进度、CHANGELOG 与 API 文档"

# ✅ 重构
git commit -m "refactor(operator): 重构算子标准协议实现

提取公共接口，简化算子实现代码。"
```

**提交前自检清单**:

- [ ] **代码已测试**
  - 单元测试（Domain/App 层）
  - 集成测试（Adapter/API 层）
  - 手动测试已通过

- [ ] **文档已更新**
  - 开发进度状态已更新
  - CHANGELOG 已添加条目
  - API/架构文档已同步（如适用）

- [ ] **代码已格式化**
  - Go: 运行 `go fmt ./...` 和 `go vet ./...`
  - 前端: 运行 `pnpm run format`（如配置）
  - 无 linter 错误

- [ ] **无调试代码**
  - 移除 `console.log`、`fmt.Println`（非日志）
  - 移除注释代码块
  - 移除临时文件

- [ ] **提交信息符合规范**
  - 有正确的 type 和 scope
  - Subject 清晰简洁
  - Body 解释"为什么"而非"做了什么"（如需要）

### 提交后操作

- [ ] 推送到远程: `git push origin <branch>`
- [ ] 创建 Pull Request（如团队协作）
- [ ] 关联相关 Issues
- [ ] 请求代码审查（如需要）

---

## 好提交 vs 坏提交

### ❌ 坏提交（不要这样做）

```bash
git commit -m "update"
git commit -m "fix bug"
git commit -m "完成功能"
git commit -m "WIP"
git commit -m "asdfasdf"
git commit -m "修改了一些东西"
```

**问题**:
- 没有类型标识
- 描述不清晰
- 无法从历史中理解改动

### ✅ 好提交（推荐做法）

```bash
git commit -m "feat(asset): 实现媒体资产管理功能"
git commit -m "fix(workflow): 修复 DAG 验证死循环"
git commit -m "docs: 更新 V1.0 架构文档"
git commit -m "refactor(ui): 优化媒体资产列表加载性能"
git commit -m "test(operator): 添加算子标准协议测试用例"
```

**优点**:
- 类型清晰（feat/fix/docs/refactor/test）
- 范围明确（asset/workflow/ui/operator）
- 描述具体，便于理解

---

## 完整流程示例

### 场景：实现媒体资产管理功能

**1. 开始开发** 🚀

```bash
# 触发 development-workflow skill
"开始开发"

# 执行检查清单：
✅ 阅读 docs/development-goya-progress.md - 确认"媒体资产管理"状态为 🚧
✅ 阅读 CHANGELOG.md - 了解最近无相关改动
✅ 阅读 docs/requirements.md - 确认需求：CRUD、搜索、派生追踪
✅ 阅读 docs/architecture.md - 确认分层架构
✅ 阅读 docs/api.md - 确认需新增 /api/v1/assets 端点
✅ 遵循 .cursor/rules/goyavision.mdc
✅ 引用 /goyavision-context 了解 MediaAsset 实体定义
```

**2. 实现功能** 💻

```
实现步骤（CQRS 模式）：
1. internal/domain/media/asset.go - 定义 MediaAsset 实体
2. internal/port/repository.go - 定义 MediaAssetRepository 接口
3. internal/adapter/persistence/media/asset.go - 实现 Repository
4. internal/app/command/create_asset.go - 创建命令
5. internal/app/command/update_asset.go - 更新命令
6. internal/app/command/delete_asset.go - 删除命令
7. internal/app/query/get_asset.go - 查询单个资产
8. internal/app/query/list_assets.go - 列表查询
9. internal/api/dto/asset.go - 定义 DTO
10. internal/api/handler/asset.go - 实现 Handler
11. internal/api/router.go - 注册路由
```

**3. 完成开发** ✅

```bash
# 触发 development-workflow skill
"完成开发"

# 步骤 1：更新开发进度
编辑 docs/development-goya-progress.md:
- 媒体资产管理: 🚧 → ✅

# 步骤 2：更新变更日志
编辑 CHANGELOG.md [未发布]:
### 新增
- **媒体资产管理**: 实现 CRUD、标签系统、派生追踪功能

# 步骤 3：更新 API 文档
编辑 docs/api.md:
添加 GET|POST /assets 端点文档

# 步骤 4：Git 提交
git add internal/ docs/
git commit -m "feat(asset): 实现媒体资产管理功能

- 添加 MediaAsset 实体和 Repository
- 实现 MediaAssetService（CRUD、搜索、派生追踪）
- 添加 Asset Handler 和 DTO
- 更新 API 文档和开发进度"

git push origin feature/asset-management
```

---

## 与其他工具的配合

### Cursor Hooks

`.cursor/hooks.json` 中配置的 `stop` hook 会在 Cursor 任务结束时触发提醒脚本，输出"完成开发"检查清单。

### Skills 组合使用

```bash
# 开始开发前
/goyavision-context  # 了解项目架构
/development-workflow "开始开发"  # 执行开发前检查

# 开发中
/goyavision-context  # 需要时查询实体/API

# 完成后
/development-workflow "完成开发"  # 执行完成检查
/review-architecture  # （可选）架构合规审查
```

---

## 常见问题

**Q: 如果只是修改一个小 Bug，也需要完整流程吗？**
A: 是的。即使是小改动，也需要：
- 更新 CHANGELOG.md（修复章节）
- 使用规范的提交格式 `fix(scope): description`
- 开发进度文档通常不需要更新（除非是重要 Bug）

**Q: 提交信息的 Body 是必须的吗？**
A: 不是必须，但推荐在以下情况添加：
- 变更原因不明显
- 涉及多个改动点
- 有特殊的技术决策
- 关联 Issue/PR

**Q: 文档更新和代码提交是分开还是一起？**
A: **推荐一起提交**。这样保证文档和代码同步，Git 历史更清晰。

**Q: 如何判断是 `feat` 还是 `refactor`？**
A:
- `feat`: 新增功能或用户可见的改进
- `refactor`: 代码重构，不改变外部行为

---

## 相关 Skills

- `/goyavision-context` - 项目架构和核心概念
- `/create-entity` - 创建新领域实体（自动化开发模式）
- `/create-operator` - 创建新算子
- `/review-architecture` - 架构合规性审查

---

## 检查清单速查

### 开始开发 ✅
- [ ] 阅读开发进度（避免重复）
- [ ] 阅读 CHANGELOG（了解最新）
- [ ] 阅读需求/架构/API
- [ ] 遵循 Rules
- [ ] 引用项目上下文

### 完成开发 ✅
- [ ] 更新开发进度
- [ ] 更新 CHANGELOG
- [ ] 更新相关文档（API/架构/需求）
- [ ] Git 提交（Conventional Commits）
- [ ] 提交前自检（测试、格式、文档）
