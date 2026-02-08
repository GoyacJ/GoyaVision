# GoyaVision V1.0 架构文档

## 概述

GoyaVision 采用**分层架构**（Clean Architecture / Hexagonal Architecture），通过严格的依赖规则确保业务逻辑与基础设施解耦。V1.0 版本引入了版本化算子与基于 DAG 的工作流引擎，构建资产驱动的智能媒体处理平台。

## 核心设计理念

### 1. 业务 = 配置，能力 = 插件，执行 = 引擎

- **配置**：工作流通过 JSON 配置定义算子节点关系，无需二次开发。
- **插件**：算子作为独立的能力单元，支持 HTTP、CLI 与 MCP 多种接入模式。
- **引擎**：高性能任务引擎负责 DAG 拓扑排序、并行执行与状态管理。

### 2. 资产驱动与不可变性

媒体资产（MediaAsset）是系统的核心。算子处理总是输入资产并产生新的资产或结构化产物（Artifact），确保处理链路的可追溯性。

### 3. 多租户与可见性隔离

全局支持 `tenant_id` 隔离，通过数据层 Scopes 自动实现所有权（Owner）与可见性（Visibility）的过滤。

## 架构层次

### 1. Domain Layer（领域层）

**职责**: 定义核心业务实体与纯业务逻辑，零外部依赖。

**核心实体**:
- **MediaSource/MediaAsset**：媒体源与多模态资产（视频、图片、音频）。
- **Operator/OperatorVersion**：版本化算子元数据与多执行模式配置。
- **Workflow/WorkflowNode/WorkflowEdge**：DAG 工作流结构定义。
- **Task/NodeExecution**：任务实例与节点级执行状态追踪。
- **Artifact**：多维产物（资产、结果、时间轴、报告）。
- **Identity (User/Role/Menu)**：权限与菜单实体。
- **SystemConfig**：分类系统配置。

### 2. Port Layer（端口层）

**职责**: 定义应用边界接口。

**关键接口**:
- `Repository`：持久化抽象。
- `ExecutorRegistry/OperatorExecutor`：算子执行器路由与实现接口。
- `SchemaValidator`：基于 JSON Schema 的 I/O 校验接口。
- `MCPClient/MCPRegistry`：MCP 协议工具发现与调用接口。
- `UnitOfWork`：跨 Repository 事务管理。

### 3. App Layer（应用层）

**职责**: 实现业务用例，采用 **CQRS** 模式分离读写操作。

**核心组件**:
- **Command Handlers**：处理写操作（如 `CreateWorkflow`、`PublishOperator`），通过 UnitOfWork 保证原子性。
- **Query Handlers**：处理读操作（如 `ListAssets`、`GetTaskStats`），支持复杂的过滤与预加载。
- **WorkflowScheduler**：基于定时（Cron/Interval）或事件（AssetNew）触发任务执行。

### 4. Adapter Layer（适配器层）

**职责**: 实现 Port 定义的接口，处理基础设施细节。

**主要适配器**:
- **Persistence**：GORM + PostgreSQL 实现，集成租户/可见性隔离 Scopes。
- **Engine**：基于 Kahn 算法的 DAG 引擎，支持并行 goroutine 执行。
- **Executors**：
  - `HTTPOperatorExecutor`：Restful 调用。
  - `CLIOperatorExecutor`：本地命令执行。
  - `MCPOperatorExecutor`：MCP 协议会话。
  - `AIModelExecutor`：大模型直连。
- **Schema**：基于 `jsonschema/v5` 的编译与校验器。

### 5. API Layer（接口层）

**职责**: HTTP 路由定义、DTO 转换、认证（JWT）与错误处理映射。

## 依赖关系规则

```text
API Layer  ──►  App Layer  ──►  Port Layer  ◄──  Domain Layer
                                     ▲
                                     │ (实现)
                                Adapter Layer
```

- **严格单向依赖**：禁止 Domain 依赖任何层，禁止 App 直接依赖 Adapter。
- **DIP 原则**：App 层仅通过 Port 接口操作基础设施。

## 数据流向：任务执行示例

1. **触发阶段**：`WorkflowScheduler` 检测到 Cron 或 Event 信号。
2. **创建阶段**：调用 `CreateTaskHandler` 初始化 Task 与待执行节点。
3. **编排阶段**：`WorkflowEngine` 对 DAG 进行拓扑排序。
4. **执行阶段**：
   - 引擎按顺序拉取 `NodeExecution`。
   - 通过 `ExecutorRegistry` 获取对应模式（HTTP/CLI/MCP）的执行器。
   - `SchemaValidator` 校验输入 Payload。
   - 执行器返回结果，持久化为 `Artifact`。
5. **结束阶段**：更新 Task 状态，通过 SSE 向前端推送进度。

---

最后更新：2026-02-08
