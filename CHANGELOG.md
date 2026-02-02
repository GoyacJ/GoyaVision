# 变更日志

本文档记录项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [未发布]

### 新增
- **MediaAsset 完整功能**（V1.0 迭代 1）
  - 添加 MediaAsset 实体（internal/domain/media_asset.go）
    - 支持视频、图片、音频三种类型
    - 支持四种来源类型（live、vod、upload、generated）
    - 支持资产派生追踪（parent_id）
    - 支持标签系统和元数据存储
  - 添加 MediaAssetRepository 接口和实现
    - 完整的 CRUD 操作
    - 支持复杂过滤和分页
  - 添加 MediaAssetService（internal/app/media_asset.go）
    - 完整的业务逻辑和验证
    - 防止删除有子资产的资产
  - 添加 MediaAsset API（internal/api/handler/asset.go）
    - GET /api/v1/assets（列表，支持过滤）
    - POST /api/v1/assets（创建）
    - GET /api/v1/assets/:id（详情）
    - PUT /api/v1/assets/:id（更新）
    - DELETE /api/v1/assets/:id（删除）
    - GET /api/v1/assets/:id/children（子资产列表）
  - 数据库迁移：自动创建 media_assets 表

- **Operator 完整功能**（V1.0 迭代 1）
  - 添加 Operator 实体（internal/domain/operator.go）
    - 支持四种分类（analysis、processing、generation、utility）
    - 支持 15+ 种算子类型（检测、OCR、ASR、剪辑等）
    - 支持版本管理和状态控制（enabled、disabled、draft）
    - 支持内置算子标识
    - 定义标准输入输出协议（OperatorInput、OperatorOutput）
  - 添加 OperatorRepository 接口和实现
    - 完整的 CRUD 操作
    - 支持复杂过滤（分类、类型、状态、内置标识、关键词搜索）
    - 支持分页查询
  - 添加 OperatorService（internal/app/operator.go）
    - Create、Get、GetByCode、List、Update、Delete
    - Enable、Disable、ListEnabled、ListByCategory
    - 完整的业务验证逻辑
    - 防止修改/删除内置算子
    - 代码唯一性检查
  - 添加 Operator API（internal/api/handler/operator.go）
    - GET /api/v1/operators（列表，支持过滤）
    - POST /api/v1/operators（创建）
    - GET /api/v1/operators/:id（详情）
    - PUT /api/v1/operators/:id（更新）
    - DELETE /api/v1/operators/:id（删除）
    - POST /api/v1/operators/:id/enable（启用）
    - POST /api/v1/operators/:id/disable（禁用）
    - GET /api/v1/operators/category/:category（按分类列出）
  - 数据库迁移：自动创建 operators 表

- **项目规范**
  - 添加文档更新强制要求（每次功能开发或修改后必须更新文档）
  - 添加 Git 提交规范（遵循 Conventional Commits）
  - 提供详细的提交检查清单和示例

### 变更
- **文档更新**
  - 更新所有 V1.0 项目文档（requirements.md、architecture.md、api.md、development-progress.md）
  - 更新 README.md 反映新架构
  - 重写 CHANGELOG.md 包含 V1.0 变更
  - 更新 .cursor/rules/goyavision.mdc（项目规则）
  - 更新 .cursor/skills/goyavision-context/SKILL.md（项目上下文）

### 计划中（V1.0 开发中）

**当前迭代重点**：
- [ ] 实现核心实体（MediaAsset、Operator、Workflow、Task、Artifact）
- [ ] 实现 Repository 和 Service 层
- [ ] 实现简化版 WorkflowEngine（单算子任务）
- [ ] API 层适配新架构
- [ ] 前端页面重构
- [ ] 数据迁移方案

**后续计划**：
- 可视化工作流设计器
- 更多内置算子（编辑、生成、转换类）
- 复杂工作流（DAG 编排）
- 自定义算子支持
- 多租户支持
- 监控与告警（Prometheus + Grafana）

## [1.0.0] - 2025-02（架构重构版本）

### 🚨 破坏性变更（不向后兼容）

此版本为架构重构版本，引入全新核心概念体系，不兼容旧版本数据和 API。

#### 核心概念重定义

- **MediaSource**（媒体源）：替代旧的 `Stream`，支持拉流、推流、上传
- **MediaAsset**（媒体资产）：新增，统一管理视频、图片、音频资产
- **Operator**（算子）：替代旧的 `Algorithm`，算子是 AI/媒体处理的能力单元
- **Workflow**（工作流）：新增，通过 DAG 编排算子
- **Task**（任务）：新增，工作流的执行实例
- **Artifact**（产物）：替代旧的 `InferenceResult`，统一管理算子输出

#### 废弃的概念

- ❌ **AlgorithmBinding**：由 Workflow 替代
- ❌ **InferenceResult**：由 Artifact 替代
- ❌ 旧的 `Stream` 概念：升级为 MediaSource
- ❌ 旧的 `Algorithm` 概念：升级为 Operator

#### 模块重命名

| 旧模块 | 新模块 | 说明 |
|--------|--------|------|
| 视频流管理 | **资产库**（Asset Library） | 媒体源、资产、录制、存储 |
| 算法管理 | **算子中心**（Operator Hub） | 算子市场、配置、监控 |
| 算法绑定 | **任务中心**（Task Center） | 工作流、任务、产物 |
| 系统管理 | **控制台**（Console） | 用户、角色、菜单、监控 |

### 新增

#### 核心能力

- **媒体资产管理**
  - 统一管理视频、图片、音频资产
  - 资产派生追踪（parent-child 关系）
  - 标签系统
  - 搜索与过滤
  - 多媒体类型支持

- **算子体系**
  - 标准化 I/O 协议（统一输入输出格式）
  - 算子分类（analyze、edit、generate、transform）
  - 内置算子（抽帧、目标检测、OCR、ASR、剪辑、转码等）
  - 算子监控（调用统计、性能指标）
  - 自定义算子支持（规划中）

- **工作流引擎**
  - DAG 工作流编排
  - 多种触发器（手动、定时、事件）
  - 节点执行与数据流转
  - 错误处理与重试
  - 简化版实现（Phase 1：单算子任务）

- **任务管理**
  - 任务创建与执行
  - 任务状态查询（实时进度）
  - 任务控制（取消、重试）
  - 任务日志

- **产物管理**
  - 统一管理算子输出
  - 产物类型：asset、result、timeline、diagnostic
  - 产物关联（任务、节点、算子、资产）
  - 产物下载导出

#### 架构改进

- **标准化协议**：算子统一的输入输出协议，确保互操作性
- **资产驱动**：以媒体资产为中心的设计理念
- **插件化**：算子作为可插拔的能力单元
- **配置化**：业务流程通过工作流配置定义

### 变更

#### API 变更

- 所有 API 端点根据新模块重新设计
- 新增端点：
  - `/api/v1/sources`（媒体源，替代 `/api/v1/streams`）
  - `/api/v1/assets`（媒体资产）
  - `/api/v1/operators`（算子，替代 `/api/v1/algorithms`）
  - `/api/v1/workflows`（工作流）
  - `/api/v1/tasks`（任务）
  - `/api/v1/artifacts`（产物，替代 `/api/v1/inference_results`）
- 废弃端点：
  - `/api/v1/streams/:id/algorithm-bindings`（由工作流替代）

#### 数据模型变更

- 新增表：
  - `media_sources`（替代 `streams`）
  - `media_assets`（新增）
  - `operators`（替代 `algorithms`）
  - `workflows`（新增）
  - `workflow_nodes`（新增）
  - `workflow_edges`（新增）
  - `tasks`（新增）
  - `artifacts`（替代 `inference_results`）
- 删除表：
  - `algorithm_bindings`
  - `inference_results`

#### 前端变更

- 模块重构：
  - 视频流管理 → 资产库
  - 算法管理 → 算子中心
  - 推理结果 → 任务中心/产物管理
- 新增页面：
  - 媒体资产管理
  - 工作流编排
  - 任务列表
  - 产物列表

### 保留（从旧版本）

#### 流媒体基础
- ✅ MediaMTX 集成（多协议支持）
- ✅ 流管理（拉流/推流）
- ✅ 实时状态查询
- ✅ 多协议预览（HLS/RTSP/RTMP/WebRTC）
- ✅ 录制与点播
- ✅ 录制文件索引

#### 认证授权
- ✅ JWT 认证（双 Token 机制）
- ✅ RBAC 权限模型
- ✅ 用户管理
- ✅ 角色管理
- ✅ 菜单管理
- ✅ 权限中间件

#### 基础设施
- ✅ 分层架构
- ✅ 配置管理（Viper）
- ✅ 数据库持久化（GORM + PostgreSQL）
- ✅ 统一错误处理
- ✅ FFmpeg 抽帧管理
- ✅ Docker Compose 部署

### 文档更新

- 完全重写需求文档（`docs/requirements.md`）
- 完全重写架构文档（`docs/architecture.md`）
- 完全重写 API 文档（`docs/api.md`）
- 更新开发进度文档（`docs/development-progress.md`）
- 更新 README.md

### 迁移指南

由于 V1.0 是架构重构版本，不提供自动迁移路径。如果您正在使用旧版本，建议：

1. **导出重要数据**：导出流配置、算法配置、推理结果
2. **全新部署 V1.0**：使用新的 Docker Compose 或手动部署
3. **手动迁移配置**：
   - 流配置 → 媒体源
   - 算法配置 → 算子
   - 算法绑定 → 工作流（需要重新配置）
4. **历史数据**：推理结果需要转换为产物格式（提供转换脚本）

---

## [0.3.0] - 2025-01-26

### 新增
- **RBAC 认证授权**（阶段 8）
  - User/Role/Permission/Menu 领域实体
  - JWT 认证（Access Token + Refresh Token）
  - 认证中间件和权限校验中间件
  - 登录/登出/刷新 Token/修改密码 API
  - 用户管理 API（CRUD、角色分配、重置密码）
  - 角色管理 API（CRUD、权限分配、菜单分配）
  - 菜单管理 API（CRUD、树形结构）
  - 权限列表 API
  - 初始化数据（默认权限、菜单、超级管理员角色、admin 账号）
- **前端认证集成**
  - Pinia 状态管理（用户、Token、权限）
  - 登录页面
  - 路由守卫（未登录跳转登录页）
  - 权限指令（v-permission）
  - 动态菜单布局
  - 系统管理页面（用户、角色、菜单管理）

### 变更
- 所有业务 API 现在需要认证才能访问
- 前端布局改为动态菜单侧边栏
- 添加 @element-plus/icons-vue 依赖

### 依赖
- 新增 golang-jwt/jwt/v5
- 新增 golang.org/x/crypto（bcrypt）
- 新增 pinia、pinia-plugin-persistedstate

## [0.2.0] - 2025-01-26

### 新增
- **前端界面**（阶段 7）
  - Vue 3 + TypeScript + Vite + Element Plus + video.js
  - 流列表页面（CRUD、预览、录制）
  - 算法管理页面
  - 推理结果查询页面
  - HLS 预览组件
  - Go embed 集成（单二进制部署）
- **预览功能**（阶段 6）
  - PreviewManager（MediaMTX/FFmpeg HLS）
  - 预览池限流
  - HLS 文件服务（/live）
- **抽帧与推理**（阶段 5）
  - Scheduler（gocron 调度器）
  - AI 推理适配器（HTTP + JSON）
  - 支持 interval_sec、schedule、initial_delay_sec
  - 推理结果查询（过滤、分页）
- **录制功能**（阶段 4）
  - RecordService（启停、会话管理）
  - 任务监控和自动状态更新
- **FFmpeg 与池**（阶段 3）
  - FFmpeg Pool（进程池与限流）
  - FFmpegManager（录制、单帧提取、连续抽帧）
- **基础与持久化**（阶段 2）
  - Stream、Algorithm、AlgorithmBinding 完整 CRUD
  - 统一错误处理机制
  - 数据库索引和约束

## [0.1.0] - 2025-01-26

### 新增
- 项目初始化和骨架搭建
- 分层架构设计（domain/port/app/adapter/api）
- 配置管理（Viper + YAML）
- 数据库模型定义（Stream, Algorithm, AlgorithmBinding, RecordSession, InferenceResult）
- HTTP API 路由框架（Echo）
- 项目文档（需求文档、开发进度、架构文档）

### 变更
- 项目从 Maas 重命名为 GoyaVision

---

## 版本说明

- **[未发布]**: 开发中，尚未发布的功能
- **[主版本.次版本.修订版本]**: 已发布的版本

### 变更类型

- **新增**: 新功能
- **变更**: 现有功能的变更
- **弃用**: 即将移除的功能
- **移除**: 已移除的功能
- **修复**: Bug 修复
- **安全**: 安全相关的修复
- **破坏性变更**: 不向后兼容的变更
