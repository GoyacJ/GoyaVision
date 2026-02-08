# 变更日志

本文档记录项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [未发布]

### 新增
- **事件触发与 EventBus 集成**：WorkflowScheduler 支持事件驱动触发。
  - 注入 EventBus，启动时订阅 `asset_new`、`asset_done`；收到事件后按触发器类型筛选启用的工作流并自动创建任务执行。
  - 新增 `internal/app/event` 包，定义 `AssetCreatedEvent`、`AssetDoneEvent`，与领域触发器类型对齐。
  - 资产创建（CreateAssetHandler）成功后发布 `asset_new` 事件，便于「资产上传后自动触发工作流」场景。

### 修复
- **事件结构体字段与接口方法同名**：`AssetCreatedEvent`/`AssetDoneEvent` 的 `OccurredAt` 字段改为 `At`，`OccurredAt()` 方法返回 `e.At`，消除与 `port.Event` 接口同名导致的编译错误。
- **WorkflowScheduler 端口引用**：`EventBus`/`Event` 改为使用 `internal/app/port`（appport）而非 `internal/port`，修复 undefined 编译错误。
- **Scheduler goroutine 错误处理**：`runWorkflow` 与 `TriggerWorkflow` 中引擎执行失败或 `UpdateTask` 失败时不再静默丢弃，统一打日志 `[WorkflowScheduler]`。
- **Scheduler 使用可追溯 context**：手动触发（`TriggerWorkflow`）时 goroutine 内改用 `context.WithoutCancel(ctx)`，保留调用方 TenantID/追踪等信息；定时触发的 `runWorkflow` 仍使用 `Background` 并补充注释说明。
- **前端 Artifact 类型与后端一致**：`web/src/api/artifact.ts` 中产物类型由 `diagnostic` 改为 `report`，与领域 `ArtifactTypeReport` 一致，避免按类型过滤不匹配。
- **ArtifactService 错误风格统一**：`internal/app/artifact.go` 全部改为使用 `pkg/apperr`（InvalidInput、NotFound、Internal），与其余 Command/Query 一致，便于 API 层统一映射。
- **DAG 层内并行错误汇总**：`internal/infra/engine/dag_engine.go` 中层内多节点并行执行时不再只返回首个错误，改为收集全部错误后通过 `errors.Join` 返回，便于排查多节点同时失败。

### 变更
- **WorkflowScheduler 构造函数**：`NewWorkflowScheduler(repo, engine, eventBus)` 增加可选参数 `eventBus`，为 nil 时不启用事件触发。
- **CreateAssetHandler 构造函数**：`NewCreateAssetHandler(uow, eventBus)` 增加可选参数 `eventBus`，用于资产创建后发布 `asset_new` 事件。
- **API/Handler 装配**：`api.NewHandlers`、`handler.NewHandlers` 增加 `eventBus` 参数；`cmd/server/main.go` 创建 `LocalEventBus` 并注入 Scheduler 与 Handlers。

## [1.0.1] - 2026-02-08

### 移除
- **删除 SimpleWorkflowEngine**：移除 `internal/adapter/engine/simple_engine.go`，该引擎功能已被 `DAGWorkflowEngine` 完全覆盖。

### 新增
- **个人信息中心深度重构**：
    - **视觉升级**：重构顶部 Banner 区域，采用简约现代的渐变与毛玻璃动效设计，优化移动端布局。
    - **组件规范化**：个人资料页全面集成 `GvInput`, `GvButton`, `GvCard`, `GvTag` 等自研 UI 组件。
    - **资产管理模块**：新增支付管理、积分管理、订阅管理、使用记录统计四大核心模块。
- **后端资产管理功能落地**：
    - 实现 `UserAsset` 领域实体、Repository 与 CQRS Handlers。
    - **支付系统集成**：引入 `gopay` 库，完整对接 **支付宝、微信、银联** 真实支付流程（统一下单、异步通知验签），支持配置化管理。
    - **积分系统**：实现每日签到获取积分逻辑及变动记录，对接前端交互。
    - **订阅系统**：支持多级订阅计划变更与状态同步。
    - **统计系统**：对接真实 API 调用统计数据，展示算子/模型消耗分布。
- **前端个人中心功能全量对接**：
    - 改造 `PaymentManager`, `PointsManager`, `SubscriptionManager`, `UsageStats` 组件。
    - 接入真实 API，实现充值跳转、扫码支付、签到领积分等核心交互逻辑。
    - 修复了数据绑定问题，解决了列表加载状态不消失的 bug。

## [1.0.0] - 2026-02-08

### 新增
- **个人中心 UI 优化与 Bug 修复**：
  - 优化个人中心顶部样式为浅色透明设计。
  - 修复后端获取 Profile 时 CreatedAt/UpdatedAt 丢失导致显示错误的问题。
  - 修复顶部导航栏头像不显示实际头像的问题。
- **用户注册与个人中心**：
    - 后端新增 `POST /api/v1/auth/register` 用户自主注册接口。
    - 后端新增 `PUT /api/v1/auth/profile` 个人资料更新接口。
    - 前端登录页新增注册模式切换功能，移除调试期的默认账号提示。
    - 新增个人中心页面 (`/profile`)，支持用户修改昵称、邮箱、手机号及头像。
    - 个人中心页面集成密码修改功能，并支持实时查看用户角色与账户状态。
- **工作流与任务管理重构**：实现基于 DAG 的可视化工作流编排与实时任务监控。
  - **后端增强**：
    - 修复节点配置 (Config)、位置 (Position) 及边条件 (Condition) 的持久化问题。
    - 任务实体增加 `NodeExecutions` 追踪每个节点的执行状态、时间与产物关联。
    - DAG 引擎支持并行执行、数据流传递与条件分支评估 (`always`/`on_success`/`on_failure`)。
    - 新增 SSE (Server-Sent Events) 端点 `GET /api/v1/tasks/:id/progress/stream` 实时推送任务进度。
    - 产物查询支持按 `node_key` 过滤，方便前端按节点展示产物。
  - **可视化编辑器**：
    - 基于 `@vue-flow/core` 与 `dagre` 实现拖拽式 DAG 设计器。
    - 支持算子库拖拽、连线、参数动态表单配置、连线 Schema 兼容性校验。
    - 支持自动布局 (TB/LR) 与全屏编辑。
  - **任务监控中心**：
    - 任务详情页支持只读 DAG 视图展示。
    - 接入 SSE 实时更新节点执行状态着色（等待/运行/成功/失败/跳过）。
    - 集成节点级产物查看器，点击节点即可查看产生的资产、结果或报告。
- **系统配置管理**：新增系统基础配置模块，支持存储、媒体处理、网络等参数的在线管理。
  - 后端：实现 `SystemConfig` 领域实体、Repository 及 API 端点。
  - 前端：新增系统配置管理页面，支持按分类编辑配置项。
- **AI 模型厂商扩展**：新增千问(Qwen)、豆包(Doubao)、智谱(Zhipu)、vLLM 四种提供商支持。
- **多租户架构重构 (Phase 1 & 2)**：
  - **基础设施**：新增 `Tenant` 实体与表结构，JWT 增强 `tenant_id` 载荷支持。
  - **数据模型**：所有核心业务表增加租户隔离与可见性字段。
  - **API适配**：更新核心资源 DTO 和 Handler，支持可见性设置。
- **前端响应式重构**：全面优化移动端适配（全局导航、资产库、GvTable、算子中心）。

### 修复
- **工作流编辑器拖拽修复**：修复从算子库拖拽算子到画布无法落位的问题。
  - 修正事件冒泡拦截逻辑，确保 `drop` 事件能被正确捕获。
  - 优化 `vueFlowRef` 绑定位置至容器 DOM，提高坐标投影计算准确性。
  - 移除导致无限更新循环的冗余 `v-model` 绑定。
  - 激活连线校验监听，修复节点间无法连线的问题。
  - 统一全局 `useVueFlow` ID，解决多组件状态不同步及操作失效问题。
- **文件管理页 500 修复**：修复系统管理-文件管理页打开报错 `column "visibility" does not exist`。
- **可见性参数传递修复**：修复所有涉及页面可见性设置参数传递始终为 0 的问题。

### 变更
- 资产管理与预览流程优化，移除废弃的算法绑定（AlgorithmBinding）相关代码，全面转向工作流驱动。
- 全面更新项目核心文档（README, requirements, architecture, api），对齐算子版本化、DAG 工作流、MCP 集成与多租户隔离等最新设计实现。
- 完善 Git 规范，新增 `docs/git-workflow.md`，并同步更新 `.clinerules`, `.cursor`, `.cline`, `.claude` 配置文件。
