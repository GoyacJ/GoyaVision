# GoyaVision Claude Code Commands

这些命令帮助你在 GoyaVision 项目中遵循标准化的开发工作流。

## 可用命令

### 🚀 开发工作流

| 命令 | 描述 | 使用时机 |
|------|------|----------|
| `/goya-dev-start` | 开始开发前的检查清单 | 开始新功能或 Bug 修复前 |
| `/goya-dev-done` | 完成开发后的检查清单 | 功能/修复完成，提交前 |
| `/goya-commit` | 创建规范的 Git 提交 | 准备提交代码时 |

### 📚 项目上下文

| 命令 | 描述 | 使用时机 |
|------|------|----------|
| `/goya-context` | 获取项目完整上下文 | 需要了解架构、API、实体定义时 |

### 📝 文档更新

| 命令 | 描述 | 使用时机 |
|------|------|----------|
| `/goya-api-doc` | 更新 API 文档指南 | 新增或修改 API 端点后 |
| `/goya-progress` | 更新开发进度指南 | 完成功能或重要里程碑后 |

## 典型工作流

### 开始新功能

```bash
1. /goya-dev-start   # 查看开发前检查清单
2. /goya-context     # 了解项目架构和 API
3. 实现功能...
4. /goya-dev-done    # 完成后检查清单
5. /goya-api-doc     # 更新 API 文档（如适用）
6. /goya-progress    # 更新开发进度
7. /goya-commit      # 创建规范提交
```

### 修复 Bug

```bash
1. /goya-dev-start   # 查看开发前检查清单
2. 修复问题...
3. /goya-dev-done    # 完成后检查清单
4. /goya-commit      # 创建规范提交（fix 类型）
```

### 查询项目信息

```bash
/goya-context        # 查看架构、API、实体、开发状态
```

## 命令详解

### `/goya-dev-start` - 开始开发

**执行内容**:
- 按顺序阅读 5 个关键文档（开发进度 → 变更日志 → 需求 → 架构 → API）
- 了解核心架构规则（分层依赖、算子协议、代码风格）
- 快速参考（默认凭证、服务端口、构建命令）

**为什么按这个顺序**:
1. 开发进度 → 避免重复或冲突的工作
2. 变更日志 → 了解最新改动和影响
3. 需求/架构/API → 理解正确的实现方式

### `/goya-dev-done` - 完成开发

**4 个必做步骤**:
1. ✅ 更新 `docs/development-progress.md` - 标记功能状态（✅/🚧/⏸️）
2. ✅ 更新 `CHANGELOG.md` - 在 `[未发布]` 下添加条目
3. ✅ 按需更新相关文档（API/架构/需求/部署）
4. ✅ 提交前自检（测试、格式化、无调试代码）

### `/goya-commit` - 规范提交

**Conventional Commits 格式**:
```
<type>(<scope>): <subject>

[optional body]
```

**常用类型**:
- `feat` - 新功能
- `fix` - Bug 修复
- `docs` - 文档变更
- `refactor` - 重构（不改变功能）
- `test`, `chore`, `perf`, `style`

**常用范围**:
- `asset`, `operator`, `workflow`, `task`, `auth`, `api`, `ui`, `source`

### `/goya-context` - 项目上下文

**提供内容**:
- 核心概念和数据流
- Clean Architecture 分层和依赖规则
- 算子标准 I/O 协议
- 完整 API 端点列表（/api/v1/*）
- 开发状态概览（✅ 已完成 / 🚧 进行中 / ⏸️ 待开始）
- 技术栈和配置快速参考

### `/goya-api-doc` - API 文档更新

**提供**:
- 端点文档模板（路径、方法、参数、请求/响应示例）
- 标准响应格式规范
- 常用错误码说明
- 更新步骤指导

### `/goya-progress` - 开发进度更新

**状态标记**:
- ✅ 已完成
- 🚧 进行中
- ⏸️ 待开始
- ⚠️ 阻塞
- 🔄 重构中

**更新指导**:
- 何时更新（必须 vs 可选）
- 如何描述（已完成子功能、技术债务、阻塞原因）
- 变更示例

## 命令设计原则

1. **渐进式** - 从简单到复杂，逐步引导
2. **清单式** - 提供明确的检查项，减少遗漏
3. **规范化** - 统一代码提交、文档更新的格式
4. **上下文感知** - 提供项目特定的架构和约定
5. **统一前缀** - 所有命令以 `goya-` 开头，避免命名冲突

## 与 Cursor Skills 的关系

这些 Claude Code 命令是基于 `.cursor/skills/` 中的 skills 设计的：

- `development-workflow` skill → `/goya-dev-start`, `/goya-dev-done`, `/goya-commit`
- `goyavision-context` skill → `/goya-context`

它们遵循相同的开发规范，但针对 Claude Code 的命令接口进行了优化。

## 命令文件结构

```
.claude/commands/
├── README.md               # 本文件
├── goya-dev-start.md       # 开发前检查清单
├── goya-dev-done.md        # 完成后检查清单
├── goya-commit.md          # Git 提交规范
├── goya-context.md         # 项目上下文
├── goya-api-doc.md         # API 文档更新
└── goya-progress.md        # 开发进度更新
```

## 相关文档

- `CLAUDE.md` - Claude Code 使用指南（包含完整命令说明）
- `docs/development-progress.md` - 开发进度跟踪
- `docs/architecture.md` - 架构设计详细说明
- `docs/api.md` - RESTful API 参考文档
- `CHANGELOG.md` - 版本变更历史
- `.cursor/skills/` - Cursor IDE Skills（相同规范）

## 快速开始

**首次使用建议**:
1. 运行 `/goya-context` 了解项目全貌
2. 运行 `/goya-dev-start` 查看开发前必读文档
3. 开始编码...
4. 完成后运行 `/goya-dev-done` 确保文档同步
5. 使用 `/goya-commit` 创建规范提交

## 反馈与改进

如有问题或建议，请在项目仓库中提 Issue 或直接修改命令文件。
