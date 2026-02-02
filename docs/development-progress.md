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

**任务清单**:

- [ ] **实体层（Domain）**
  - [ ] MediaAsset 实体定义
  - [ ] Operator 实体定义（重构 Algorithm）
  - [ ] Workflow 实体定义（替代 AlgorithmBinding）
  - [ ] Task 实体定义
  - [ ] Artifact 实体定义
  - [ ] WorkflowNode、WorkflowEdge 定义

- [ ] **端口层（Port）**
  - [ ] MediaAssetRepository 接口
  - [ ] OperatorRepository 接口
  - [ ] WorkflowRepository 接口
  - [ ] TaskRepository 接口
  - [ ] ArtifactRepository 接口
  - [ ] OperatorPort 接口（算子执行）
  - [ ] WorkflowEngine 接口（工作流引擎）

- [ ] **应用层（App）**
  - [ ] MediaAssetService 实现
  - [ ] OperatorService 实现
  - [ ] WorkflowService 实现
  - [ ] TaskService 实现
  - [ ] ArtifactService 实现
  - [ ] Scheduler 重构（适配新架构）

- [ ] **适配器层（Adapter）**
  - [ ] MediaAssetRepository 实现（GORM）
  - [ ] OperatorRepository 实现（GORM）
  - [ ] WorkflowRepository 实现（GORM）
  - [ ] TaskRepository 实现（GORM）
  - [ ] ArtifactRepository 实现（GORM）
  - [ ] SimpleWorkflowEngine 实现（单算子）

- [ ] **API 层（API）**
  - [ ] MediaAsset Handler + DTO
  - [ ] Operator Handler + DTO
  - [ ] Workflow Handler + DTO
  - [ ] Task Handler + DTO
  - [ ] Artifact Handler + DTO
  - [ ] 路由注册

- [ ] **数据库迁移**
  - [ ] 创建新表（media_assets、operators、workflows、tasks、artifacts）
  - [ ] 数据迁移脚本（streams → media_sources、algorithms → operators）
  - [ ] 删除旧表（algorithm_bindings、inference_results）

### 迭代 2：前端适配

**目标**: 前端适配新 API 和概念

**任务清单**:

- [ ] **页面重构**
  - [ ] 媒体源页面（保持现有功能）
  - [ ] 媒体资产页面（新增）
  - [ ] 算子中心页面（替代算法管理）
  - [ ] 工作流页面（替代算法绑定）
  - [ ] 任务页面（新增）
  - [ ] 产物页面（替代推理结果）

- [ ] **API 客户端**
  - [ ] 更新 API 定义（TypeScript）
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
