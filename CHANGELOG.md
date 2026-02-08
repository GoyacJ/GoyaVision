# 变更日志

本文档记录项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [未发布]

### 新增
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
