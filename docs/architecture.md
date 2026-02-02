# GoyaVision V1.0 架构文档

## 概述

GoyaVision 采用**分层架构**（Clean Architecture / Hexagonal Architecture），确保业务逻辑与基础设施解耦，提高可维护性和可测试性。V1.0 版本引入全新的核心概念体系，围绕 **MediaAsset**、**Operator**、**Workflow** 构建智能媒体处理平台。

## 核心设计理念

### 1. 业务 = 配置，能力 = 插件，执行 = 引擎

- **配置**：工作流通过 JSON/YAML 配置定义，无需编码
- **插件**：算子作为可插拔的能力单元，支持自定义扩展
- **引擎**：统一的任务执行引擎，负责工作流编排和调度

### 2. 标准化 I/O 协议

所有算子遵循统一的输入输出协议，确保算子间可互操作：

```
输入：MediaAsset + Params
输出：Artifact（新资产、结构化结果、时间轴片段）
```

### 3. 资产驱动

媒体资产是系统的核心"事实"，所有处理围绕资产展开：

```
MediaSource → MediaAsset → Operator → Artifact → MediaAsset
```

### 4. 工作流编排

复杂业务通过工作流编排实现，不修改底层能力：

```
工作流 = DAG（算子节点 + 数据流转 + 条件分支）
```

## 架构层次

### 1. Domain Layer（领域层）

**位置**: `internal/domain/`

**职责**:
- 定义核心业务实体
- 包含纯业务逻辑
- 无外部依赖

**核心实体**:

#### 媒体资产类
- **MediaSource**：媒体源（流、上传）
- **MediaAsset**：媒体资产（视频、图片、音频）

#### 算子与工作流类
- **Operator**：算子定义
- **Workflow**：工作流定义
- **WorkflowNode**：工作流节点
- **WorkflowEdge**：工作流边
- **Task**：任务实例
- **Artifact**：产物

#### 认证授权类
- **User**：用户
- **Role**：角色
- **Permission**：权限
- **Menu**：菜单

**原则**:
- 实体为纯 Go struct，包含业务字段和验证逻辑
- 不依赖任何基础设施
- 可包含领域方法（如 `MediaAsset.IsVideo()`）

### 2. Port Layer（端口层）

**位置**: `internal/port/`

**职责**:
- 定义应用边界接口
- 定义领域需要的服务接口

**核心接口**:

```go
type Repository interface {
    MediaSource
    MediaAsset
    Operator
    Workflow
    Task
    Artifact
    User
    Role
    Permission
    Menu
}

type MediaSourceRepository interface {
    Create(ctx context.Context, source *domain.MediaSource) error
    GetByID(ctx context.Context, id string) (*domain.MediaSource, error)
    List(ctx context.Context, filter SourceFilter) ([]*domain.MediaSource, error)
    Update(ctx context.Context, source *domain.MediaSource) error
    Delete(ctx context.Context, id string) error
}

type MediaAssetRepository interface {
    Create(ctx context.Context, asset *domain.MediaAsset) error
    GetByID(ctx context.Context, id string) (*domain.MediaAsset, error)
    List(ctx context.Context, filter AssetFilter) ([]*domain.MediaAsset, int64, error)
    Update(ctx context.Context, asset *domain.MediaAsset) error
    Delete(ctx context.Context, id string) error
    ListBySource(ctx context.Context, sourceID string) ([]*domain.MediaAsset, error)
}

type OperatorPort interface {
    Execute(ctx context.Context, input OperatorInput) (*OperatorOutput, error)
}

type WorkflowEngine interface {
    Execute(ctx context.Context, task *domain.Task) error
    GetStatus(ctx context.Context, taskID string) (*TaskStatus, error)
    Cancel(ctx context.Context, taskID string) error
}
```

**原则**:
- 接口定义在 port 层
- 实现放在 adapter 层
- 通过依赖注入使用

### 3. App Layer（应用层）

**位置**: `internal/app/`

**职责**:
- 编排 domain 与 port
- 实现业务用例
- 不直接依赖 adapter 具体实现

**核心服务**:

#### 资产库（Asset Library）
- **MediaSourceService**：媒体源管理
  - CRUD 操作
  - 流状态查询（集成 MediaMTX）
  - 多协议地址生成
- **MediaAssetService**：媒体资产管理
  - CRUD 操作
  - 资产搜索与过滤
  - 资产关联管理（parent-child）
- **RecordService**：录制管理
  - 启停录制（集成 MediaMTX API）
  - 录制会话管理
  - 录制文件索引
- **PlaybackService**：点播服务
  - 录制段查询
  - 点播 URL 生成

#### 算子中心（Operator Hub）
- **OperatorService**：算子管理
  - CRUD 操作
  - 算子分类与搜索
  - 算子版本管理

#### 任务中心（Task Center）
- **WorkflowService**：工作流管理
  - CRUD 操作
  - DAG 验证
  - 工作流版本管理
- **TaskService**：任务管理
  - 任务创建与执行
  - 任务状态查询
  - 任务控制（取消、重试）
- **Scheduler**：调度器
  - 定时任务调度（gocron）
  - 事件驱动调度
  - 任务优先级队列
- **WorkflowEngine**：工作流引擎
  - DAG 执行
  - 节点间数据流转
  - 错误处理与重试
- **ArtifactService**：产物管理
  - 产物查询
  - 产物关联

#### 控制台（Console）
- **AuthService**：认证服务
  - 登录、登出
  - Token 刷新
  - 密码管理
- **UserService**：用户管理
  - CRUD 操作
  - 角色分配
- **RoleService**：角色管理
  - CRUD 操作
  - 权限分配
  - 菜单分配
- **MenuService**：菜单管理
  - CRUD 操作
  - 树形结构

**原则**:
- 通过 port 接口操作，不直接调用 adapter
- 一个服务对应一个业务用例
- 处理业务规则和流程编排

### 4. Adapter Layer（适配器层）

**位置**: `internal/adapter/`

**职责**:
- 实现 port 定义的接口
- 处理基础设施细节

**适配器**:

#### 持久化（Persistence）
- **位置**: `internal/adapter/persistence/`
- **职责**: 实现 `port.Repository`
- **技术**: GORM + PostgreSQL
- **组件**:
  - `repository.go`：统一 Repository 实现
  - `media_source.go`：媒体源持久化
  - `media_asset.go`：媒体资产持久化
  - `operator.go`：算子持久化
  - `workflow.go`：工作流持久化
  - `task.go`：任务持久化
  - `artifact.go`：产物持久化
  - `user.go`：用户持久化
  - `init_data.go`：初始化数据

#### 流媒体（MediaMTX）
- **位置**: `internal/adapter/mediamtx/`
- **职责**: MediaMTX HTTP API 客户端
- **功能**:
  - 路径配置管理
  - 录制管理
  - 流状态查询
  - 点播文件索引

#### AI 推理（AI）
- **位置**: `internal/adapter/ai/`
- **职责**: 实现 `port.OperatorPort`
- **功能**:
  - HTTP + JSON 调用算子服务
  - 超时与重试
  - 输入输出转换

#### 工作流引擎（Workflow）
- **位置**: `internal/adapter/workflow/`
- **职责**: 实现 `port.WorkflowEngine`
- **功能**:
  - DAG 解析与验证
  - 节点执行编排
  - 数据流转
  - 错误处理

**原则**:
- 实现 port 接口
- 处理技术细节（数据库、HTTP、进程等）
- 将外部模型转换为 domain 模型

### 5. API Layer（接口层）

**位置**: `internal/api/`

**职责**:
- HTTP 路由定义
- 请求/响应处理
- DTO 转换
- 中间件

**组件**:

#### 路由（Router）
- **位置**: `internal/api/router.go`
- **职责**: 注册所有 HTTP 路由
- **分组**:
  - 公开路由（无需认证）：`/api/v1/auth/login`、`/api/v1/auth/refresh`
  - 受保护路由（需认证）：所有业务 API
  - 管理路由（需权限）：用户、角色、菜单管理

#### 处理器（Handler）
- **位置**: `internal/api/handler/`
- **职责**: 处理 HTTP 请求
- **组件**:
  - `auth.go`：认证相关
  - `user.go`：用户管理
  - `role.go`：角色管理
  - `menu.go`：菜单管理
  - `source.go`：媒体源管理
  - `asset.go`：媒体资产管理
  - `record.go`：录制管理
  - `playback.go`：点播服务
  - `operator.go`：算子管理
  - `workflow.go`：工作流管理
  - `task.go`：任务管理
  - `artifact.go`：产物管理

#### DTO（数据传输对象）
- **位置**: `internal/api/dto/`
- **职责**: API 请求/响应结构定义
- **原则**:
  - DTO 与 domain 模型分离
  - 提供转换函数（`ToDTO()`、`FromDTO()`）
  - 包含验证标签

#### 中间件（Middleware）
- **位置**: `internal/api/middleware/`
- **组件**:
  - `auth.go`：JWT 认证、权限校验
  - `cors.go`：跨域配置
  - `logger.go`：请求日志
  - `recovery.go`：错误恢复

**原则**:
- Handler 调用 app 服务或 port 接口
- 不直接操作数据库
- DTO 与 domain 模型分离

## 依赖关系

```
┌─────────────┐
│   API       │  HTTP 层
└──────┬──────┘
       │ 依赖
┌──────▼──────┐
│    App      │  应用服务层
└──────┬──────┘
       │ 依赖
┌──────▼──────┐      ┌─────────────┐
│   Port      │◄─────┤   Domain    │  接口与领域
└──────┬──────┘      └─────────────┘
       │ 实现
┌──────▼──────┐
│  Adapter    │  适配器层
└─────────────┘
```

### 依赖规则

1. **Domain** 不依赖任何其他层
2. **Port** 可依赖 Domain
3. **App** 可依赖 Domain 和 Port
4. **Adapter** 可依赖 Domain 和 Port（实现接口）
5. **API** 可依赖 App、Port、Domain

## 数据流

### 示例：创建媒体资产并执行工作流

```
1. HTTP Request → API Handler (asset.go)
2. Handler 转换 DTO → Domain 模型 (MediaAsset)
3. Handler 调用 MediaAssetService.Create()
4. MediaAssetService 调用 Repository.Create()
5. Adapter.Persistence 保存到数据库
6. 返回 MediaAsset
7. 如果配置了自动工作流，触发 WorkflowService.TriggerEvent()
8. WorkflowService 创建 Task
9. Scheduler 调度任务执行
10. WorkflowEngine 执行 DAG
11. 调用 Operator 处理资产
12. 保存 Artifact
13. 返回执行结果
```

## 核心模块设计

### 1. 媒体资产管理

```
MediaSource (媒体源)
  ├── type: pull, push, upload
  ├── protocol: rtsp, rtmp, hls, webrtc, file
  └── status: ready, online, offline

MediaAsset (媒体资产)
  ├── type: video, image, audio
  ├── source_type: live, vod, upload, generated
  ├── source_id: 关联 MediaSource
  ├── parent_id: 派生关系
  └── metadata: 分辨率、帧率、时长等
```

**设计原则**：
- 资产是不可变的（immutable），修改产生新资产
- 通过 parent_id 追踪资产派生关系
- 支持多级派生（原始视频 → 抽帧图片 → 检测结果图片）

### 2. 算子系统

```
Operator (算子)
  ├── category: analyze, edit, generate, transform
  ├── input_spec: 输入规格（支持的资产类型、参数）
  ├── output_spec: 输出规格（产物类型、结构）
  ├── endpoint: HTTP 服务端点
  └── is_builtin: 内置 vs 自定义
```

**标准化 I/O 协议**：

```go
type OperatorInput struct {
    AssetID string                 `json:"asset_id"`
    Params  map[string]interface{} `json:"params"`
}

type OperatorOutput struct {
    OutputAssets []OutputAsset     `json:"output_assets"`
    Results      []Result           `json:"results"`
    Timeline     []TimelineEvent    `json:"timeline"`
    Diagnostics  *Diagnostics       `json:"diagnostics"`
}

type OutputAsset struct {
    Type     string                 `json:"type"`
    Path     string                 `json:"path"`
    Format   string                 `json:"format"`
    Metadata map[string]interface{} `json:"metadata"`
}

type Result struct {
    Type       string                 `json:"type"`
    Data       map[string]interface{} `json:"data"`
    Confidence float64                `json:"confidence"`
}

type TimelineEvent struct {
    Start      float64                `json:"start"`
    End        float64                `json:"end"`
    EventType  string                 `json:"event_type"`
    Confidence float64                `json:"confidence"`
    Data       map[string]interface{} `json:"data"`
}
```

**设计原则**：
- 算子无状态，幂等执行
- 统一输入输出协议，确保互操作性
- 支持超时、重试、取消

### 3. 工作流引擎

```
Workflow (工作流)
  ├── trigger: manual, schedule, event
  ├── nodes: DAG 节点列表
  └── edges: 节点连接关系

Task (任务)
  ├── workflow_id: 关联工作流
  ├── status: pending, running, completed, failed
  ├── progress: 0-100
  └── current_node: 当前执行节点

WorkflowNode (工作流节点)
  ├── operator_id: 关联算子
  ├── params: 节点参数
  ├── retry: 重试次数
  └── timeout: 超时时间

WorkflowEdge (工作流边)
  ├── from: 源节点 ID
  ├── to: 目标节点 ID
  └── condition: 条件表达式（可选）
```

**执行流程**：

```
1. 触发器激活 → 创建 Task
2. Scheduler 调度任务 → WorkflowEngine 执行
3. 解析 DAG → 拓扑排序
4. 按顺序执行节点：
   a. 获取输入资产
   b. 调用 Operator
   c. 保存 Artifact
   d. 传递给下游节点
5. 所有节点完成 → 任务成功
6. 任何节点失败 → 重试或标记失败
```

**设计原则**：
- DAG 验证：无环、连通
- 支持并行执行：多个独立分支并行
- 支持条件分支：根据上游结果选择分支
- 数据流转：上游输出作为下游输入
- 错误处理：节点失败不影响其他分支

### 4. 产物管理

```
Artifact (产物)
  ├── type: asset, result, timeline, diagnostic
  ├── task_id: 关联任务
  ├── node_id: 关联节点
  ├── operator_id: 关联算子
  ├── asset_id: 关联资产（如果是 asset 类型）
  └── data: 产物数据（JSON）
```

**产物类型**：
- **asset**：新生成的媒体资产（剪辑后的视频、检测结果图片）
- **result**：结构化结果（检测框、分类标签、OCR 文本）
- **timeline**：时间轴片段（事件、高光、镜头切分）
- **diagnostic**：诊断信息（性能指标、模型版本）

**设计原则**：
- 产物与任务、节点、算子关联
- 产物可关联到新资产（资产派生）
- 产物支持查询、下载、导出

## 配置管理

**位置**: `config/`、`configs/`

- 使用 Viper 加载 YAML 配置
- 支持环境变量覆盖（`GOYAVISION_*` 前缀）
- 配置结构体定义在 `config/config.go`

**核心配置**：

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
  secret: "your-secret-key"
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

## 进程管理

**位置**: `pkg/`

### FFmpeg 管理
- **位置**: `pkg/ffmpeg/`
- **职责**: 抽帧任务管理（AI 推理用）
- **组件**:
  - `pool.go`：进程池与限流（max_frame）
  - `manager.go`：抽帧任务管理

### 存储管理
- **位置**: `pkg/storage/`
- **职责**: 文件存储管理
- **组件**:
  - `manager.go`：文件上传、下载、删除
  - `lifecycle.go`：生命周期管理（自动清理）

## 调度器

**位置**: `internal/app/scheduler.go`

- 使用 gocron/v2 管理定时任务
- 支持触发器类型：
  - **manual**：手动触发
  - **schedule**：定时触发（cron 表达式）
  - **event**：事件触发（新资产、录制完成、流上线）
- 支持并发控制、优先级队列
- 启动时自动加载启用的工作流

## 前端集成

**位置**: `web/`、`internal/api/static.go`

- Vue 3 + TypeScript + Vite + Element Plus + video.js
- 构建产物嵌入到 Go 二进制（`embed.FS`）
- SPA 路由处理（所有非 API 路由返回 index.html）
- HLS/WebRTC 播放器集成

## 扩展点

### 添加新算子

1. 实现算子 HTTP 服务（符合标准 I/O 协议）
2. 在算子中心注册算子
3. 配置算子端点、输入输出规格
4. 在工作流中使用算子

### 添加新工作流

1. 在任务中心创建工作流
2. 配置 DAG 节点和边
3. 配置触发器
4. 启用工作流

### 添加新存储适配器

1. 在 `port/` 定义存储接口
2. 在 `adapter/storage/` 实现接口（如 OSS、S3）
3. 在配置中选择存储后端

## 安全考虑

- **输入验证**：API 层和 Service 层双重验证
- **SQL 注入防护**：使用 GORM 参数化查询
- **XSS 防护**：前端输入转义
- **CSRF 防护**：SameSite Cookie
- **JWT 认证**：Access Token + Refresh Token 双 Token 机制
- **RBAC 权限模型**：用户-角色-权限三级授权
- **密码加密**：bcrypt 算法
- **审计日志**：记录所有关键操作

## 性能优化

- **数据库**：
  - 连接池
  - 索引优化（stream_id、task_id、created_at）
  - 分页查询
- **缓存**：
  - 用户权限缓存
  - 算子配置缓存
- **并发**：
  - 进程池限流
  - 任务并发控制
  - 读写锁保护共享资源
- **异步处理**：
  - 工作流异步执行
  - 产物异步生成

## 监控与可观测性

### 日志
- 结构化日志（JSON 格式）
- 日志级别：DEBUG、INFO、WARN、ERROR
- 请求 ID 追踪

### 指标（规划）
- Prometheus 指标暴露（`/metrics`）
- 关键指标：
  - 任务执行数、成功率、失败率
  - 算子调用数、延迟
  - 系统资源使用率

### 健康检查
- `/health`：应用健康状态
- `/ready`：就绪探针（数据库、MediaMTX 连通性）

## 测试策略

- **单元测试**：Domain、App 层逻辑
- **集成测试**：Adapter 实现（数据库、HTTP 客户端）
- **端到端测试**：API 层完整流程
- **工作流测试**：DAG 执行、数据流转、错误处理

## 部署架构

### 单机部署

```
┌─────────────────────────────────────┐
│  GoyaVision + MediaMTX + PostgreSQL │
│  Docker Compose                     │
└─────────────────────────────────────┘
```

### 集群部署（规划）

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ GoyaVision  │     │ GoyaVision  │     │ GoyaVision  │
│  Instance 1 │     │  Instance 2 │     │  Instance 3 │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
       ┌───────────────────┴───────────────────┐
       │                                       │
       ▼                                       ▼
┌─────────────┐                         ┌─────────────┐
│  MediaMTX   │                         │ PostgreSQL  │
│   Cluster   │                         │   Cluster   │
└─────────────┘                         └─────────────┘
```

---

**注意**: 本文档会随着项目演进持续更新。
