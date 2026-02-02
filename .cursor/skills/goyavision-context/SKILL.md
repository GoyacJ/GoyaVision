---
name: goyavision-context
description: GoyaVision V1.0 项目结构、核心概念、API 约定与开发状态。在实现或评审 GoyaVision 功能时使用，以便遵循既定分层、数据模型和标准协议。
---

# GoyaVision V1.0 项目上下文

## 何时使用

- 在 GoyaVision 仓库中实现新功能、修改 handler/app/domain/adapter 时
- 需要确认实体、API 路径、配置项或开发状态时
- 需要了解已实现的功能和代码结构时
- 需要遵循算子标准协议或工作流编排规范时

## 版本说明

**当前版本**：V1.0（架构重构版本）

**核心变更**：
- 引入全新核心概念：MediaAsset、MediaSource、Operator、Workflow、Task、Artifact
- 废弃：AlgorithmBinding、InferenceResult
- 模块重命名：资产库、算子中心、任务中心、控制台
- **不向后兼容**

## 项目结构（核心）

```
cmd/server/          入口；config、GORM、Echo、Router、Scheduler、embed、初始化数据
config/              配置加载（Viper + YAML）
configs/             配置文件（config.yaml、mediamtx.yml）
internal/
  domain/            MediaSource, MediaAsset, Operator, Workflow, WorkflowNode, WorkflowEdge,
                     Task, Artifact, User, Role, Permission, Menu
  port/              Repository, OperatorPort, WorkflowEngine, MediaMTXClient
  app/               MediaSourceService, MediaAssetService, OperatorService,
                     WorkflowService, TaskService, ArtifactService, Scheduler,
                     RecordService, PlaybackService,
                     AuthService, UserService, RoleService, MenuService
  adapter/
    persistence/     Repository 实现（GORM）、初始化数据（init_data.go）
    mediamtx/        MediaMTX HTTP API 客户端
    workflow/        WorkflowEngine 实现（DAG 执行）
    ai/              OperatorPort 实现（HTTP 客户端）
  api/
    handler/         source, asset, operator, workflow, task, artifact,
                     record, playback, auth, user, role, menu
    dto/             source, asset, operator, workflow, task, artifact,
                     record, playback, auth, user, role, menu
    middleware/      auth.go（JWT 认证、权限校验）
    errors.go        统一错误处理
    static.go        前端静态文件服务（embed）
    router.go        路由注册（公开路由、认证路由、管理路由）
pkg/
  ffmpeg/            Pool（进程池）、Manager（抽帧，用于 AI 推理）
  storage/           Manager（文件管理）、Lifecycle（生命周期）
web/                 Vue 3 前端（src/, dist/）
  src/store/         Pinia 状态管理（用户、权限）
  src/views/login/   登录页面
  src/views/asset/   资产库页面（源、资产、录制、点播）
  src/views/operator/ 算子中心页面（算子市场、配置）
  src/views/workflow/ 任务中心页面（工作流、任务、产物）
  src/views/system/  系统管理页面（用户、角色、菜单）
  src/layout/        动态菜单布局
  src/directives/    权限指令（v-permission）
  src/router/guard.ts 路由守卫
docs/               需求、开发进度、架构文档、API 文档、部署指南
```

## 核心概念（V1.0）

### 资产类

#### MediaSource（媒体源）
- **作用**：媒体的来源（流、上传）
- **类型**：
  - `pull`：拉流（从外部地址拉取）
  - `push`：推流（等待外部推送）
  - `upload`：文件上传
- **协议**：rtsp、rtmp、hls、webrtc、file
- **状态**：ready、online、offline

#### MediaAsset（媒体资产）
- **作用**：统一管理视频、图片、音频资产
- **类型**：video、image、audio
- **来源类型**：
  - `live`：实时流录制或抽帧
  - `vod`：点播视频
  - `upload`：用户上传
  - `generated`：算子生成
- **关键属性**：
  - `source_id`：关联的媒体源
  - `parent_id`：派生自哪个资产（资产派生追踪）
  - `tags`：标签数组
  - `metadata`：扩展元数据（分辨率、帧率、时长等）

### 算子与工作流类

#### Operator（算子）
- **作用**：AI/媒体处理的能力单元
- **分类**：
  - `analyze`（分析）：检测、识别、分类、追踪、OCR、ASR
  - `edit`（编辑）：剪辑、裁剪、打码、去水印、字幕、水印
  - `generate`（生成）：TTS、配音、摘要、高光
  - `transform`（转换）：转码、压缩、分辨率调整、增强
- **标准化协议**：统一的输入输出格式（见下文）
- **关键属性**：
  - `code`：唯一编码
  - `version`：版本号
  - `input_spec`：输入规格
  - `output_spec`：输出规格
  - `endpoint`：HTTP 服务端点
  - `is_builtin`：内置 vs 自定义

#### Workflow（工作流）
- **作用**：通过 DAG 编排算子，实现复杂业务流程
- **触发器**：
  - `manual`：手动触发
  - `schedule`：定时触发（cron 表达式）
  - `event`：事件触发（新资产、录制完成、流上线）
- **组成**：
  - `nodes`：工作流节点（operator_id、params、retry、timeout）
  - `edges`：节点连接（from、to、condition）

#### Task（任务）
- **作用**：工作流的执行实例
- **状态**：pending、running、completed、failed、cancelled
- **关键属性**：
  - `workflow_id`：关联的工作流
  - `input_assets`：输入资产列表
  - `progress`：进度（0-100）
  - `current_node`：当前执行节点

#### Artifact（产物）
- **作用**：算子/工作流的输出结果
- **类型**：
  - `asset`：新生成的媒体资产
  - `result`：结构化结果（检测框、标签、文本）
  - `timeline`：时间轴片段（事件、高光、镜头切分）
  - `diagnostic`：诊断信息（性能指标、模型版本）
- **关联**：task_id、node_id、operator_id、asset_id

### 废弃概念（不再使用）
- ❌ `Stream`：升级为 MediaSource
- ❌ `Algorithm`：升级为 Operator
- ❌ `AlgorithmBinding`：由 Workflow 替代
- ❌ `InferenceResult`：由 Artifact 替代

## 算子标准协议

所有算子必须遵循统一的输入输出协议。

### 输入格式

```json
{
  "asset_id": "资产 ID",
  "params": {
    "key": "value"
  }
}
```

### 输出格式

```json
{
  "output_assets": [
    {
      "type": "video|image|audio",
      "path": "存储路径",
      "format": "格式",
      "metadata": {}
    }
  ],
  "results": [
    {
      "type": "detection|classification|ocr|...",
      "data": {},
      "confidence": 0.95
    }
  ],
  "timeline": [
    {
      "start": 0.0,
      "end": 5.0,
      "event_type": "事件类型",
      "confidence": 0.95,
      "data": {}
    }
  ],
  "diagnostics": {
    "latency_ms": 150,
    "model_version": "v1.0",
    "device": "gpu"
  }
}
```

## 已实现功能（V1.0 开发状态）

### ✅ 已完成（从旧版本保留）

#### 流媒体基础
- MediaMTX 集成（多协议支持：RTSP/RTMP/HLS/WebRTC）
- 媒体源管理（拉流/推流）
- 实时状态查询
- 多协议预览
- 录制与点播（集成 MediaMTX）
- 录制文件索引

#### 认证授权
- JWT 认证（Access Token + Refresh Token 双 Token 机制）
- RBAC 权限模型
- 用户管理、角色管理、菜单管理
- 权限中间件
- 前端：Pinia 状态管理、登录页面、路由守卫、权限指令、动态菜单

#### 基础设施
- 分层架构（Domain、Port、App、Adapter、API）
- 配置管理（Viper + YAML）
- 数据库持久化（GORM + PostgreSQL）
- 统一错误处理
- FFmpeg 抽帧管理
- Docker Compose 部署

### 🚧 进行中（V1.0 核心功能）

#### 资产库
- [ ] 媒体资产管理（CRUD、搜索、派生追踪）
- [ ] 存储配置（生命周期管理）

#### 算子中心
- [ ] 算子管理（CRUD、分类、版本管理）
- [ ] 内置算子（抽帧、目标检测 - 需要重构为算子）
- [ ] 算子监控（调用统计、性能指标）

#### 任务中心
- [ ] 工作流管理（CRUD、DAG 验证）
- [ ] 简化工作流（Phase 1：单算子任务）
- [ ] 任务管理（创建、执行、查询、控制）
- [ ] 任务调度（定时调度、事件触发）
- [ ] 产物管理（查询、关联）

#### 前端
- [ ] 媒体资产页面
- [ ] 算子中心页面
- [ ] 工作流编排页面
- [ ] 任务列表页面
- [ ] 产物列表页面

### ⏸️ 待开始（V1.0 后续）

- 可视化工作流设计器
- 更多内置算子（编辑、生成、转换类）
- 复杂工作流（DAG 编排、并行执行、条件分支）
- 自定义算子（Docker 镜像上传）
- 多租户支持
- 监控与告警（Prometheus + Grafana）

## API 端点（V1.0）

### 基础
- **前缀**：`/api/v1`
- **认证**：所有业务 API 需要 `Authorization: Bearer <access_token>`

### 认证（Auth）
- `POST /auth/login`：登录
- `POST /auth/refresh`：刷新 Token
- `GET /auth/profile`：获取当前用户信息
- `PUT /auth/password`：修改密码
- `POST /auth/logout`：登出

### 媒体源（Sources）
- `GET/POST /sources`：列表、创建
- `GET/PUT/DELETE /sources/:id`：详情、更新、删除
- `POST /sources/:id/enable`：启用
- `POST /sources/:id/disable`：禁用
- `GET /sources/:id/status`：获取实时状态
- `GET /sources/:id/preview`：获取预览 URL
- `GET /sources/:id/preview/ready`：检查流就绪

### 录制（Record）
- `POST /sources/:id/record/start`：启动录制
- `POST /sources/:id/record/stop`：停止录制
- `GET /sources/:id/record/status`：获取录制状态
- `GET /sources/:id/record/sessions`：列出录制会话
- `GET /sources/:id/record/files`：列出录制文件

### 点播（Playback）
- `GET /sources/:id/playback?start=<timestamp>`：获取点播 URL
- `GET /sources/:id/playback/segments`：列出录制段

### 媒体资产（Assets）
- `GET/POST /assets`：列表、创建（支持过滤：type、source_type、source_id、tags）
- `GET/PUT/DELETE /assets/:id`：详情、更新、删除
- `GET /assets/:id/children`：列出子资产（派生资产）

### 算子（Operators）
- `GET/POST /operators`：列表、创建（支持过滤：category、status、is_builtin）
- `GET/PUT/DELETE /operators/:id`：详情、更新、删除
- `POST /operators/:id/enable`：启用
- `POST /operators/:id/disable`：禁用
- `POST /operators/:id/test`：测试算子

### 工作流（Workflows）
- `GET/POST /workflows`：列表、创建（支持过滤：status）
- `GET/PUT/DELETE /workflows/:id`：详情、更新、删除
- `POST /workflows/:id/activate`：启用工作流
- `POST /workflows/:id/pause`：暂停工作流
- `POST /workflows/:id/validate`：验证工作流

### 任务（Tasks）
- `GET/POST /tasks`：列表、创建（支持过滤：workflow_id、status、trigger_type）
- `GET /tasks/:id`：详情
- `POST /tasks/:id/cancel`：取消任务
- `POST /tasks/:id/retry`：重试任务
- `GET /tasks/:id/logs`：获取任务日志

### 产物（Artifacts）
- `GET /artifacts`：列表（支持过滤：task_id、node_id、operator_id、type）
- `GET /artifacts/:id`：详情
- `DELETE /artifacts/:id`：删除
- `GET /artifacts/:id/download`：下载产物

### 用户管理（Users）
- `GET/POST /users`：列表、创建
- `GET/PUT/DELETE /users/:id`：详情、更新、删除
- `POST /users/:id/reset-password`：重置密码

### 角色管理（Roles）
- `GET/POST /roles`：列表、创建
- `GET/PUT/DELETE /roles/:id`：详情、更新、删除

### 菜单管理（Menus）
- `GET/POST /menus`：列表、创建
- `GET/PUT/DELETE /menus/:id`：详情、更新、删除
- `GET /menus/tree`：获取菜单树

### 权限（Permissions）
- `GET /permissions`：列出所有权限

### 静态文件
- `/live/*`：HLS 文件服务（已废弃，使用 MediaMTX）
- `/*`：前端 SPA

## 配置项（V1.0）

### 主配置文件：`configs/config.yaml`

```yaml
server:
  port: 8080

db:
  dsn: "host=localhost user=goyavision password=goyavision dbname=goyavision port=5432 sslmode=disable"

ffmpeg:
  bin: "ffmpeg"
  max_frame: 16

ai:
  timeout: 10s
  retry: 2

jwt:
  secret: "your-secret-key-change-in-production"
  expire: 2h
  refresh_exp: 168h
  issuer: "goyavision"

mediamtx:
  api_address: "http://localhost:9997"
  rtsp_address: "rtsp://localhost:8554"
  rtmp_address: "rtmp://localhost:1935"
  hls_address: "http://localhost:8888"
  webrtc_address: "http://localhost:8889"
  playback_address: "http://localhost:9996"
  record_path: "./data/recordings/%path/%Y-%m-%d_%H-%M-%S"
  record_format: "fmp4"
  segment_duration: "1h"

storage:
  base_path: "./data"
  recordings_path: "./data/recordings"
  frames_path: "./data/frames"
  uploads_path: "./data/uploads"
```

### 环境变量覆盖
所有配置项都可以通过环境变量覆盖，前缀为 `GOYAVISION_`。

例如：
```bash
export GOYAVISION_DB_DSN="host=localhost ..."
export GOYAVISION_JWT_SECRET="your-production-secret"
```

## 数据模型概要（V1.0）

### 资产库
- `media_sources`：id, name, type, url, protocol, enabled, created_at, updated_at
- `media_assets`：id, type, source_type, source_id, parent_id, name, path, duration, size, format, metadata, status, tags, created_at, updated_at

### 算子与工作流
- `operators`：id, code, name, category, version, input_spec, output_spec, endpoint, config, status, is_builtin, description, icon, created_at, updated_at
- `workflows`：id, name, description, trigger, nodes, edges, status, created_at, updated_at
- `tasks`：id, workflow_id, trigger_type, input_assets, status, progress, current_node, started_at, completed_at, error, created_at
- `artifacts`：id, task_id, node_id, operator_id, type, asset_id, data, created_at

### 认证授权
- `users`：id, username, password, nickname, email, phone, avatar, status, created_at, updated_at
- `roles`：id, code, name, description, status, created_at, updated_at
- `permissions`：id, code, name, method, path, description
- `menus`：id, parent_id, code, name, type, path, icon, component, permission, sort, visible, status, created_at, updated_at
- `user_roles`：user_id, role_id
- `role_permissions`：role_id, permission_id
- `role_menus`：role_id, menu_id

## 文档

- **需求文档**：`docs/requirements.md`
- **架构文档**：`docs/architecture.md`
- **开发进度**：`docs/development-progress.md`
- **API 文档**：`docs/api.md`
- **部署指南**：`docs/DEPLOYMENT.md`
- **变更日志**：`CHANGELOG.md`

## 开发路线

### Phase 1：核心闭环（当前 V1.0）
- 媒体源管理（✅ 已完成）
- 媒体资产管理（🚧 进行中）
- 内置算子（抽帧、目标检测）（🚧 重构中）
- 简化工作流（单算子任务）（🚧 进行中）
- 任务调度与执行（🚧 进行中）
- 产物管理（🚧 进行中）

### Phase 2：能力扩展
- 多媒体类型（图片、音频）
- 更多内置算子（编辑、生成、转换类）
- 复杂工作流（DAG 编排、并行、条件分支）
- 可视化工作流设计器
- 工作流模板市场

### Phase 3：平台化
- 自定义算子（Docker 镜像）
- 算子市场（第三方算子）
- 多租户支持
- 开放 API 与 SDK
- 监控与告警（Prometheus + Grafana）

## 默认账号

- **用户名**：admin
- **密码**：admin123
- **角色**：超级管理员（拥有所有权限）

## 常见开发模式

### 创建新实体

1. 在 `internal/domain/` 定义实体
2. 在 `internal/port/` 定义 Repository 接口
3. 在 `internal/adapter/persistence/` 实现 Repository
4. 在 `internal/app/` 实现 Service
5. 在 `internal/api/dto/` 定义 DTO
6. 在 `internal/api/handler/` 实现 Handler
7. 在 `internal/api/router.go` 注册路由

### 实现新算子

1. 实现算子 HTTP 服务（符合标准 I/O 协议）
2. 在算子中心注册算子（code、category、version、endpoint、input_spec、output_spec）
3. 在工作流中使用算子

### 创建工作流

1. 定义工作流（name、description、trigger）
2. 添加节点（operator_id、params、retry、timeout）
3. 添加边（from、to、condition）
4. 验证 DAG（无环、连通）
5. 启用工作流

## 注意事项

1. **V1.0 不向后兼容**：旧版本数据和 API 需要手动迁移
2. **资产派生追踪**：使用 `parent_id` 追踪资产派生关系（原始视频 → 抽帧图片 → 检测结果图片）
3. **算子幂等性**：算子应设计为无状态、幂等执行
4. **工作流验证**：DAG 必须无环、连通
5. **错误传播**：节点失败不影响其他独立分支
6. **产物关联**：产物可关联新资产（通过 asset_id）

## 开发规范

### 文档更新要求（强制）

**每次完成功能开发或修改后，必须同步更新相关文档：**

1. **必须更新**：
   - `docs/development-progress.md`：更新功能状态、迭代进度
   - `docs/api.md`：新增或修改 API 时更新
   - `CHANGELOG.md`：在 `[未发布]` 章节记录变更

2. **可能需要更新**：
   - `docs/requirements.md`：功能需求变更时
   - `docs/architecture.md`：架构设计变更时
   - `README.md`：影响用户使用时

### Git 提交规范（强制）

**每次完成功能开发或修改后，必须进行 Git 提交：**

#### Commit Message 格式

遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>(<scope>): <subject>
```

**Type 类型**：
- `feat`：新功能
- `fix`：Bug 修复
- `docs`：文档变更
- `refactor`：代码重构
- `test`：测试相关
- `chore`：构建、配置、依赖等

**Scope 范围**（可选）：
- `asset`、`operator`、`workflow`、`task`、`auth`、`api`、`ui`

**示例**：
```bash
feat(asset): 实现媒体资产管理功能
fix(workflow): 修复 DAG 验证死循环
docs: 更新 V1.0 架构文档
```

#### 提交检查清单

- [ ] 代码已测试
- [ ] 相关文档已更新
- [ ] 代码已格式化（gofmt / goimports）
- [ ] Commit message 符合规范

## 技术债务

- AlgorithmBinding 迁移到 Workflow（高优先级）
- InferenceResult 迁移到 Artifact（高优先级）
- FFmpeg Pool 优化（中优先级）
- 数据库索引优化（中优先级）
- 前端性能优化（低优先级）
