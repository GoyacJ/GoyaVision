# GoyaVision Git 工作流规范

本手册定义了 GoyaVision 项目的 Git 分支管理与提交规范，旨在提高协作效率并保证代码质量。

## 1. 分支管理规范

项目采用简化版的 **Git Flow** 模型。

### 1.1 分支分类

| 分支类型 | 命名规范 | 来源分支 | 合并回 | 说明 |
| :--- | :--- | :--- | :--- | :--- |
| **主分支** | `main` | - | - | 生产环境分支，对应线上稳定版本 |
| **开发分支** | `develop` | `main` | `main` | 核心开发分支，用于日常集成 |
| **功能分支** | `feature/*` | `develop` | `develop` | 新功能开发，完成后删除 |
| **修复分支** | `fix/*` | `develop` | `develop` | 普通 BUG 修复，完成后删除 |
| **紧急修复** | `hotfix/*` | `main` | `main`, `develop` | 线上紧急故障修复，完成后删除 |
| **重构分支** | `refactor/*` | `develop` | `develop` | 技术重构或架构调整 |
| **发布分支** | `release/*` | `develop` | `main`, `develop` | 版本发布前的预发布测试与准备 |

### 1.2 工作流操作指引

1.  **开发新功能**：
    *   从 `develop` 分支拉取：`git checkout -b feature/user-auth develop`
    *   开发、测试并提交代码。
    *   发起 PR 合并回 `develop` 分支。
2.  **修复线上 BUG**：
    *   从 `main` 分支拉取：`git checkout -b hotfix/critical-leak main`
    *   修复并验证。
    *   发起 PR 同时合并回 `main` 和 `develop` 分支。
3.  **合并策略**：
    *   功能分支合并建议使用 **Squash Merge**，保持 `develop` 分支历史整洁。
    *   预发布分支合并建议使用 **No-FF**，保留版本发布的历史节点。

---

## 2. 提交规范 (Commit Specs)

遵循 **Conventional Commits** (约定式提交) 标准。

### 2.1 提交格式

```text
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

### 2.2 类型 (Type)

| 类型 | 说明 |
| :--- | :--- |
| **feat** | 新增功能 |
| **fix** | 修复缺陷 |
| **docs** | 文档更新 |
| **style** | 代码格式调整（不影响逻辑，如缩进、空格） |
| **refactor** | 代码重构 |
| **perf** | 性能优化 |
| **test** | 测试用例编写 |
| **chore** | 构建流程或辅助工具变更 |

### 2.3 范围 (Scope)

常用范围如下：
*   `asset`: 资产库
*   `operator`: 算子中心
*   `workflow`: 工作流
*   `task`: 任务管理
*   `auth`: 认证授权
*   `api`: API 接口层
*   `ui`: 前端界面
*   `infra`: 基础设施适配器

### 2.4 描述要求

*   **Subject**：必须使用中文描述，简明扼要。
*   **Body**：对于复杂的变更，需在正文中说明实现方案、技术背景或影响点。
*   **Footer**：
    *   修复 BUG 应引用：`Closes #IssueID`。
    *   破坏性变更应说明：`BREAKING CHANGE: <description>`。

---

## 3. 开发流程约束

1.  **提交即更新**：代码变更必须与相关的 `development-progress.md` 和 `CHANGELOG.md` 更新一并提交。
2.  **原子化提交**：一次提交应仅解决一个特定任务或 BUG。
3.  **提交前自检**：
    *   代码已通过本地编译。
    *   格式已对齐（`go fmt` / `npm run lint`）。
    *   无任何调试代码（如 `fmt.Println`、`console.log`）。
    *   涉及 API 变更已同步更新 `docs/api.md`。

---

最后更新：2026-02-08
