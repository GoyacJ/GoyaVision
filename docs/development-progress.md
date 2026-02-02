# GoyaVision V1.0 开发进度

## 版本说明

**当前版本**: V1.0（架构重构版本）

**核心变更**:
- 引入全新核心概念：MediaAsset、MediaSource、Operator、Workflow、Task、Artifact
- 废弃 AlgorithmBinding，由 Workflow 替代
- 模块重命名：资产库、算子中心、任务中心、控制台
- 不向后兼容，全新架构

## 开发路线

### Phase 1：核心闭环（V1.0）

基础能力建设，实现最小可用系统。

| 模块 | 功能 | 状态 | 说明 |
|------|------|------|------|
| **资产库** | | | |
| 媒体源管理 | CRUD、状态查询 | ✅ 已完成 | 基于 MediaMTX，支持拉流/推流 |
| 媒体资产管理 | CRUD、搜索过滤 | 🚧 进行中 | 需要实现 MediaAsset 实体和服务 |
| 录制管理 | 启停录制、文件索引 | ✅ 已完成 | 集成 MediaMTX 录制 API |
| 点播服务 | 录制段查询、URL 生成 | ✅ 已完成 | 集成 MediaMTX Playback |
| 存储配置 | 存储路径配置 | ⏸️ 待开始 | 支持本地存储 |
| **算子中心** | | | |
| 算子管理 | CRUD、分类 | 🚧 进行中 | 需要实现 Operator 实体和服务 |
| 内置算子 | 抽帧、目标检测 | ✅ 部分完成 | 已有抽帧和推理，需要重构为算子 |
| 算子监控 | 调用统计、性能 | ⏸️ 待开始 | |
| **任务中心** | | | |
| 工作流管理 | CRUD | 🚧 进行中 | 需要实现 Workflow 实体和服务 |
| 简化工作流 | 单算子任务 | 🚧 进行中 | 先支持单算子，不支持 DAG |
| 任务管理 | 创建、查询、控制 | 🚧 进行中 | 需要实现 Task 实体和服务 |
| 任务调度 | 定时调度、事件触发 | 🚧 进行中 | 基于 gocron，需要适配新架构 |
| 产物管理 | 查询、关联 | 🚧 进行中 | 需要实现 Artifact 实体和服务 |
| **控制台** | | | |
| 认证服务 | 登录、Token 刷新 | ✅ 已完成 | JWT 双 Token 机制 |
| 用户管理 | CRUD、角色分配 | ✅ 已完成 | RBAC 权限模型 |
| 角色管理 | CRUD、权限分配 | ✅ 已完成 | |
| 菜单管理 | CRUD、树形结构 | ✅ 已完成 | 动态菜单 |
| 仪表盘 | 系统概览 | ⏸️ 待开始 | |
| 审计日志 | 操作日志 | ⏸️ 待开始 | |
| **前端** | | | |
| 媒体源页面 | 流管理、预览 | ✅ 已完成 | 支持多协议预览 |
| 媒体资产页面 | 资产列表、搜索 | ⏸️ 待开始 | |
| 算子中心页面 | 算子市场 | ⏸️ 待开始 | |
| 工作流页面 | 工作流列表 | ⏸️ 待开始 | |
| 任务页面 | 任务列表、详情 | ⏸️ 待开始 | |
| 产物页面 | 产物列表 | ⏸️ 待开始 | |
| 系统管理页面 | 用户、角色、菜单 | ✅ 已完成 | |

### Phase 2：能力扩展

扩展媒体类型和算子能力。

| 模块 | 功能 | 状态 | 说明 |
|------|------|------|------|
| **资产库** | | | |
| 图片资产 | 图片上传、管理 | ⏸️ 待开始 | |
| 音频资产 | 音频上传、管理 | ⏸️ 待开始 | |
| 资产标签 | 标签系统 | ⏸️ 待开始 | |
| 生命周期管理 | 自动清理策略 | ⏸️ 待开始 | |
| **算子中心** | | | |
| 编辑类算子 | 剪辑、打码、水印 | ⏸️ 待开始 | |
| 生成类算子 | TTS、高光摘要 | ⏸️ 待开始 | |
| 转换类算子 | 转码、压缩、增强 | ⏸️ 待开始 | |
| 算子版本 | 多版本管理 | ⏸️ 待开始 | |
| **任务中心** | | | |
| 复杂工作流 | DAG 编排 | ⏸️ 待开始 | 支持并行、条件分支 |
| 可视化设计器 | 拖拽式 DAG 设计 | ⏸️ 待开始 | |
| 工作流模板 | 预定义模板 | ⏸️ 待开始 | |
| 任务优先级 | 优先级队列 | ⏸️ 待开始 | |

### Phase 3：平台化

开放能力，支持自定义扩展。

| 模块 | 功能 | 状态 | 说明 |
|------|------|------|------|
| **算子中心** | | | |
| 自定义算子 | Docker 镜像上传 | ⏸️ 待开始 | |
| 算子市场 | 第三方算子 | ⏸️ 待开始 | |
| 算子沙箱 | 隔离执行 | ⏸️ 待开始 | |
| **开放平台** | | | |
| API 文档 | OpenAPI 规范 | ⏸️ 待开始 | |
| SDK | Go/Python/JS SDK | ⏸️ 待开始 | |
| Webhook | 事件通知 | ⏸️ 待开始 | |
| **多租户** | | | |
| 租户隔离 | tenant_id 隔离 | ⏸️ 待开始 | |
| 资源配额 | 存储、计算配额 | ⏸️ 待开始 | |
| **监控告警** | | | |
| Prometheus | 指标暴露 | ⏸️ 待开始 | |
| Grafana | 可视化看板 | ⏸️ 待开始 | |
| 告警规则 | 告警配置 | ⏸️ 待开始 | |

## 当前迭代重点（V1.0）

### 迭代 0：文档与规范（已完成）

**目标**: 更新所有项目文档，建立开发规范

**已完成**:
- [x] 更新需求文档（`docs/requirements.md`）
- [x] 更新架构文档（`docs/architecture.md`）
- [x] 更新 API 文档（`docs/api.md`）
- [x] 更新开发进度文档（`docs/development-progress.md`）
- [x] 更新 README.md
- [x] 更新 CHANGELOG.md
- [x] 更新项目规则（`.cursor/rules/goyavision.mdc`）
- [x] 更新项目技能（`.cursor/skills/goyavision-context/SKILL.md`）
- [x] 建立文档更新规范
- [x] 建立 Git 提交规范（Conventional Commits）

### 迭代 1：核心实体与服务（当前）

**目标**: 实现新架构的核心实体和服务

**已完成（全部 5 个核心实体）**:

- [x] **实体层（Domain）**
  - [x] MediaAsset 实体定义（media_asset.go）
    - 支持视频、图片、音频三种类型
    - 支持四种来源类型（live、vod、upload、generated）
    - 支持资产派生追踪（parent_id）
    - 支持标签系统（tags）
    - 支持元数据存储（metadata）
  - [x] Operator 实体定义（operator.go）
    - 支持四种分类（analysis、processing、generation、utility）
    - 支持 15+ 种算子类型（检测、OCR、ASR、剪辑等）
    - 支持版本管理和状态控制
    - 支持内置算子标识
    - 定义标准输入输出协议（OperatorInput、OperatorOutput）
  - [x] Workflow 实体定义（workflow.go）
    - 支持五种触发类型（manual、schedule、event、asset_new、asset_done）
    - 支持 DAG 工作流定义（WorkflowNode、WorkflowEdge）
    - 支持节点配置和位置信息
    - 支持边条件和路由
    - 支持版本管理和状态控制（enabled、disabled、draft）
  - [x] Task 实体定义（task.go）
    - 支持五种状态（pending、running、success、failed、cancelled）
    - 关联工作流和资产
    - 支持进度跟踪（0-100%）
    - 记录当前执行节点
    - 记录执行时间（开始、完成）
    - 支持错误信息记录
  - [x] Artifact 实体定义（artifact.go）
    - 支持四种类型（asset、result、timeline、report）
    - 关联任务和资产
    - 支持 JSONB 数据存储
    - 定义标准数据结构（AssetInfo、TimelineSegment、AnalysisResult）

- [x] **端口层（Port）**
  - [x] MediaAssetRepository 接口（7个方法）
    - Create、Get、List、Update、Delete
    - ListBySource、ListByParent
  - [x] OperatorRepository 接口（8个方法）
    - Create、Get、GetByCode、List、Update、Delete
    - ListEnabled、ListByCategory
  - [x] WorkflowRepository 接口（8个方法）
    - Create、Get、GetByCode、GetWithNodes、List、Update、Delete
    - ListEnabled
  - [x] WorkflowNode/Edge Repository 接口（6个方法）
    - CreateNode、ListNodes、DeleteNodes
    - CreateEdge、ListEdges、DeleteEdges
  - [x] TaskRepository 接口（8个方法）
    - Create、Get、GetWithRelations、List、Update、Delete
    - GetStats、ListRunning
  - [x] ArtifactRepository 接口（6个方法）
    - Create、Get、List、Delete
    - ListByTask、ListByType

- [x] **适配器层（Adapter）**
  - [x] MediaAssetRepository 实现（GORM + PostgreSQL）
    - 完整的 CRUD 实现
    - 支持复杂过滤（类型、来源、状态、标签、时间范围）
    - 支持分页查询
    - AutoMigrate 集成
  - [x] OperatorRepository 实现（GORM + PostgreSQL）
    - 完整的 CRUD 实现
    - 支持复杂过滤（分类、类型、状态、内置标识、关键词搜索）
    - 支持分页查询
    - AutoMigrate 集成
  - [x] WorkflowRepository 实现（GORM + PostgreSQL）
    - 完整的 CRUD 实现
    - 支持复杂过滤（状态、触发类型、标签、关键词搜索）
    - 支持预加载节点和边（Preload）
    - 级联删除支持（CASCADE）
    - AutoMigrate 集成
  - [x] TaskRepository 实现（GORM + PostgreSQL）
    - 完整的 CRUD 实现
    - 支持复杂过滤（工作流、资产、状态、时间范围）
    - 支持预加载关联数据（Workflow、Asset、Artifacts）
    - 支持统计查询（按状态分组）
    - AutoMigrate 集成
  - [x] ArtifactRepository 实现（GORM + PostgreSQL）
    - 完整的 CRUD 实现
    - 支持复杂过滤（任务、类型、资产、时间范围）
    - 支持预加载关联数据（Task、Asset）
    - 支持按任务和类型查询
    - AutoMigrate 集成

- [x] **应用层（App）**
  - [x] MediaAssetService 实现（media_asset.go）
    - Create、Get、List、Update、Delete
    - ListBySource、ListChildren
    - 完整的业务验证逻辑
    - 防止删除有子资产的资产
  - [x] OperatorService 实现（operator.go）
    - Create、Get、GetByCode、List、Update、Delete
    - Enable、Disable、ListEnabled、ListByCategory
    - 完整的业务验证逻辑
    - 防止修改/删除内置算子
    - 代码唯一性检查
  - [x] WorkflowService 实现（workflow.go）
    - Create、Get、GetWithNodes、GetByCode、List、Update、Delete
    - Enable、Disable、ListEnabled
    - 完整的业务验证逻辑
    - 节点和边的级联管理
    - 启用前验证工作流完整性
    - 代码唯一性检查
  - [x] TaskService 实现（task.go）
    - Create、Get、GetWithRelations、List、Update、Delete
    - Start、Complete、Fail、Cancel
    - GetStats、ListRunning
    - 完整的业务验证逻辑
    - 状态转换管理（自动记录开始/完成时间）
    - 进度范围验证（0-100%）
    - 防止删除运行中的任务
  - [x] ArtifactService 实现（artifact.go）
    - Create、Get、List、Delete
    - ListByTask、ListByType
    - 完整的业务验证逻辑
    - 验证关联的任务和资产存在性

- [x] **API 层（API）**
  - [x] MediaAsset DTO（asset.go）
    - Request：AssetCreateReq、AssetUpdateReq、AssetListQuery
    - Response：AssetResponse、AssetListResponse
    - 转换函数：AssetToResponse、AssetsToResponse
  - [x] MediaAsset Handler（asset.go）
    - GET /assets（列表，支持过滤）
    - POST /assets（创建）
    - GET /assets/:id（详情）
    - PUT /assets/:id（更新）
    - DELETE /assets/:id（删除）
    - GET /assets/:id/children（子资产列表）
  - [x] Operator DTO（operator.go）
    - Request：OperatorCreateReq、OperatorUpdateReq、OperatorListQuery
    - Response：OperatorResponse、OperatorListResponse
    - 转换函数：OperatorToResponse、OperatorsToResponse
  - [x] Operator Handler（operator.go）
    - GET /operators（列表，支持过滤）
    - POST /operators（创建）
    - GET /operators/:id（详情）
    - PUT /operators/:id（更新）
    - DELETE /operators/:id（删除）
    - POST /operators/:id/enable（启用）
    - POST /operators/:id/disable（禁用）
    - GET /operators/category/:category（按分类列出）
  - [x] Workflow DTO（workflow.go）
    - Request：WorkflowCreateReq、WorkflowUpdateReq、WorkflowListQuery
    - Response：WorkflowResponse、WorkflowWithNodesResponse、WorkflowNodeResponse、WorkflowEdgeResponse
    - 转换函数：WorkflowToResponse、WorkflowToResponseWithNodes、WorkflowsToResponse
  - [x] Workflow Handler（workflow.go）
    - GET /workflows（列表，支持过滤）
    - POST /workflows（创建）
    - GET /workflows/:id（详情，支持 with_nodes 参数）
    - PUT /workflows/:id（更新）
    - DELETE /workflows/:id（删除）
    - POST /workflows/:id/enable（启用）
    - POST /workflows/:id/disable（禁用）
  - [x] Task DTO（task.go）
    - Request：TaskCreateReq、TaskUpdateReq、TaskListQuery
    - Response：TaskResponse、TaskWithRelationsResponse、TaskStatsResponse
    - 转换函数：TaskToResponse、TaskToResponseWithRelations、TasksToResponse、TaskStatsToResponse
  - [x] Task Handler（task.go）
    - GET /tasks（列表，支持过滤）
    - POST /tasks（创建）
    - GET /tasks/:id（详情，支持 with_relations 参数）
    - PUT /tasks/:id（更新）
    - DELETE /tasks/:id（删除）
    - POST /tasks/:id/start（启动）
    - POST /tasks/:id/complete（完成）
    - POST /tasks/:id/fail（失败）
    - POST /tasks/:id/cancel（取消）
    - GET /tasks/stats（统计）
  - [x] Artifact DTO（artifact.go）
    - Request：ArtifactCreateReq、ArtifactListQuery
    - Response：ArtifactResponse、ArtifactListResponse
    - 转换函数：ArtifactToResponse、ArtifactsToResponse
  - [x] Artifact Handler（artifact.go）
    - GET /artifacts（列表，支持过滤）
    - POST /artifacts（创建）
    - GET /artifacts/:id（详情）
    - DELETE /artifacts/:id（删除）
    - GET /tasks/:task_id/artifacts（列出任务的产物，支持类型过滤）
  - [x] 路由注册（router.go）

## 迭代 1 总结

**✅ 核心实体层（5/5 完成 - 100%）**

全部 5 个核心实体已完成实现！

---

### 迭代 2：工作流引擎与调度器（当前）

**目标**: 实现工作流执行引擎和任务调度系统

**已完成**:

- [x] **端口层（Port）**
  - [x] OperatorExecutor 接口（engine.go）
    - Execute：执行算子
  - [x] WorkflowEngine 接口（engine.go）
    - Execute：执行工作流
    - Cancel：取消执行
    - GetProgress：获取进度

- [x] **适配器层（Adapter）**
  - [x] HTTPOperatorExecutor 实现（engine/http_executor.go）
    - 通过 HTTP 调用外部算子服务
    - 支持自定义 HTTP 方法（POST/GET）
    - 支持超时控制（5 分钟）
    - 标准化输入输出协议
    - 完整的错误处理
  - [x] SimpleWorkflowEngine 实现（engine/simple_engine.go）
    - 支持单算子顺序执行
    - 支持进度跟踪（按节点数计算）
    - 支持取消执行（Context 取消）
    - 自动保存产物（Assets、Results、Timeline）
    - 完整的任务状态管理
    - 并发安全（sync.RWMutex）

- [x] **应用层（App）**
  - [x] WorkflowScheduler 实现（workflow_scheduler.go）
    - 支持定时调度（Cron、Interval）
    - 支持手动触发（TriggerWorkflow）
    - 自动加载启用的工作流
    - 支持取消调度
    - 异步执行工作流（goroutine）
    - 完整的错误处理

- [x] **集成与 API**
  - [x] 更新 main.go 集成工作流引擎和调度器
  - [x] 更新 handler.Deps 传递 WorkflowScheduler
  - [x] 更新 Router 签名适配新 Deps
  - [x] 添加手动触发 API
    - POST /api/v1/workflows/:id/trigger（手动触发工作流）

---

### 迭代 4：数据迁移与清理（当前）

**目标**: 清理废弃代码，创建数据迁移工具

**已完成**:

- [x] **数据迁移工具**
  - [x] 创建迁移命令（cmd/migrate/main.go）
    - 支持 dry-run 模式
    - Streams → MediaAssets 迁移（作为媒体源）
    - Algorithms → Operators 迁移
    - 清理旧表（algorithm_bindings、inference_results）
    - 确认提示和详细日志

- [x] **删除废弃代码**
  - [x] Domain 层（3 个文件）
    - algorithm.go
    - algorithm_binding.go
    - inference_result.go
  - [x] Handler 层（3 个文件）
    - algorithm.go
    - algorithm_binding.go
    - inference.go
  - [x] App 层（4 个文件）
    - algorithm.go
    - algorithm_binding.go
    - inference.go
    - scheduler.go（旧调度器）
  - [x] DTO 层（3 个文件）
    - algorithm.go
    - algorithm_binding.go
    - inference.go
  - [x] Adapter 层（1 个文件）
    - ai/inference.go
  - [x] Port 层（1 个文件）
    - inference.go

- [x] **更新核心文件**
  - [x] internal/port/repository.go（删除 13 个旧方法）
  - [x] internal/adapter/persistence/repository.go（删除实现，更新 AutoMigrate）
  - [x] internal/api/router.go（删除 3 个旧路由）
  - [x] cmd/server/main.go（移除旧 Scheduler，简化导入）

**待实现**:

---

### 迭代 3：前端适配（当前）

**目标**: 前端适配新 API 和概念，升级为顶部菜单栏布局

**已完成**:

- [x] **布局改造**
  - [x] 将侧边栏布局改为顶部菜单栏布局（layout/index.vue）
    - 移除侧边栏（el-aside）
    - Logo 移至顶部左侧
    - 菜单横向显示（mode="horizontal"）
    - 现代化视觉设计（渐变 Logo、悬停效果）
    - 保留用户下拉菜单功能

- [x] **API 客户端（TypeScript）**
  - [x] asset.ts（媒体资产 API）
    - 类型定义：MediaAsset、AssetListQuery、AssetCreateReq、AssetUpdateReq
    - 6 个 API 方法：list、get、create、update、delete、listChildren
  - [x] operator.ts（算子 API）
    - 类型定义：Operator、OperatorListQuery、OperatorCreateReq、OperatorUpdateReq
    - 8 个 API 方法：list、get、create、update、delete、enable、disable、listByCategory
  - [x] workflow.ts（工作流 API）
    - 类型定义：Workflow、WorkflowNode、WorkflowEdge、WorkflowListQuery、WorkflowCreateReq、WorkflowUpdateReq
    - 8 个 API 方法：list、get、create、update、delete、enable、disable、trigger
  - [x] task.ts（任务 API）
    - 类型定义：Task、TaskListQuery、TaskCreateReq、TaskUpdateReq、TaskStats
    - 9 个 API 方法：list、get、create、update、delete、start、complete、fail、cancel、getStats
  - [x] artifact.ts（产物 API）
    - 类型定义：Artifact、ArtifactListQuery、ArtifactCreateReq
    - 5 个 API 方法：list、get、create、delete、listByTask

- [x] **页面实现**
  - [x] 媒体资产页面（views/asset/index.vue）
    - 列表展示（类型、来源、格式、大小、时长、状态）
    - 搜索过滤（名称、类型、来源类型、状态）
    - CRUD 操作（创建、查看、编辑、删除）
    - 分页支持
  - [x] 算子中心页面（views/operator/index.vue）
    - 列表展示（代码、名称、分类、类型、版本、状态、内置标识）
    - 搜索过滤（关键词、分类、状态、内置算子）
    - CRUD 操作（创建、查看、编辑、删除）
    - 启用/禁用功能
    - 保护内置算子（不可编辑/删除）
  - [x] 工作流页面（views/workflow/index.vue）
    - 列表展示（代码、名称、触发方式、版本、状态）
    - 搜索过滤（关键词、触发方式、状态）
    - CRUD 操作（创建、查看、编辑、删除）
    - 启用/禁用功能
    - 手动触发功能（支持指定资产）
  - [x] 任务中心页面（views/task/index.vue）
    - 统计卡片（总数、待执行、执行中、已成功、已失败、已取消）
    - 列表展示（任务 ID、工作流、状态、进度、当前节点、时间、耗时）
    - 状态过滤
    - 查看任务详情
    - 取消运行中的任务
    - 删除已完成/失败的任务
    - 查看任务产物（入口）

- [x] **路由配置**
  - [x] 更新路由定义（router/index.ts）
    - 注册新页面：/assets、/operators、/workflows、/tasks
    - 保留旧页面（标记为"旧"）：/streams、/algorithms、/inference-results
    - 默认重定向到 /assets

**待实现**:
  - [ ] 适配新端点

- [ ] **路由与菜单**
  - [ ] 更新路由配置
  - [ ] 更新菜单配置
  - [ ] 更新权限配置

### 迭代 3：测试与优化

**目标**: 确保新架构稳定可用

**任务清单**:

- [ ] **单元测试**
  - [ ] Domain 层测试
  - [ ] App 层测试

- [ ] **集成测试**
  - [ ] Adapter 层测试
  - [ ] API 层测试

- [ ] **端到端测试**
  - [ ] 创建媒体源 → 录制 → 创建资产
  - [ ] 创建工作流 → 触发任务 → 生成产物
  - [ ] 完整业务流程测试

- [ ] **文档更新**
  - [ ] API 文档
  - [ ] 用户手册
  - [ ] 部署文档

## 技术债务

| 问题 | 优先级 | 状态 | 说明 |
|------|--------|------|------|
| AlgorithmBinding 迁移 | 高 | 待处理 | 需要迁移到 Workflow |
| InferenceResult 迁移 | 高 | 待处理 | 需要迁移到 Artifact |
| FFmpeg Pool 优化 | 中 | 待处理 | 资源泄漏检查 |
| 数据库索引优化 | 中 | 待处理 | 添加缺失索引 |
| 前端性能优化 | 低 | 待处理 | 大列表虚拟滚动 |

## 已完成功能（从旧版本保留）

### 流媒体基础
- ✅ MediaMTX 集成（多协议支持）
- ✅ 流管理（拉流/推流）
- ✅ 实时状态查询
- ✅ 多协议预览（HLS/RTSP/RTMP/WebRTC）
- ✅ 录制与点播
- ✅ 录制文件索引

### 认证授权
- ✅ JWT 认证（双 Token 机制）
- ✅ RBAC 权限模型
- ✅ 用户管理
- ✅ 角色管理
- ✅ 菜单管理
- ✅ 权限中间件
- ✅ 前端权限指令

### 基础设施
- ✅ 分层架构
- ✅ 配置管理（Viper）
- ✅ 数据库持久化（GORM + PostgreSQL）
- ✅ 统一错误处理
- ✅ FFmpeg 抽帧管理
- ✅ Docker Compose 部署

## 风险与阻塞

| 风险 | 影响 | 应对措施 | 状态 |
|------|------|----------|------|
| 数据迁移复杂性 | 高 | 编写迁移脚本，充分测试 | 待处理 |
| 前端重构工作量 | 中 | 分阶段迭代，保持核心功能可用 | 待处理 |
| 工作流引擎复杂度 | 中 | Phase 1 先实现简化版（单算子） | 待处理 |
| 算子标准化 | 中 | 定义清晰的 I/O 协议文档 | 进行中 |

## 下一步行动

### 本周（Week 1）

1. 完成核心实体定义（MediaAsset、Operator、Workflow、Task、Artifact）
2. 实现 Repository 接口和 GORM 持久化
3. 数据库迁移方案设计

### 下周（Week 2）

1. 实现 App 层服务（MediaAssetService、OperatorService、WorkflowService）
2. 实现简化版 WorkflowEngine（单算子任务）
3. 实现 API Handler 和 DTO

### 两周后（Week 3）

1. 前端页面重构
2. 集成测试
3. 文档更新

### 一个月后（Week 4）

1. 端到端测试
2. 性能优化
3. V1.0 版本发布

## 变更记录

| 日期 | 版本 | 变更内容 |
|------|------|----------|
| 2025-02 | V1.0 | 架构重构：引入 MediaAsset、Operator、Workflow、Task、Artifact；废弃 AlgorithmBinding；模块重命名；不向后兼容 |
| 2025-01 | V0.9 | MediaMTX 集成、录制重构、点播服务、认证授权完成 |
| 2024-12 | V0.1 | 项目初始化、基础骨架搭建 |

---

**注意**: 本文档会随着项目演进持续更新。每周更新迭代进度。
