---
name: goyavision-context
description: 获取 GoyaVision 项目架构、核心概念、API 端点和开发状态的完整上下文
---

# GoyaVision 项目上下文

提供 GoyaVision V1.0 智能媒体处理平台的完整项目上下文，包括架构设计、核心概念、API 约定和开发状态。

## 何时使用

✅ **推荐场景**：
- 开始实现新功能前，需要了解项目整体架构
- 修改 handler/app/domain/adapter 层代码时
- 需要查询 API 端点、实体定义或配置项
- 遵循算子标准协议或工作流编排规范时
- 新团队成员快速了解项目

❌ **不适用场景**：
- 只需要查看单个文件内容（使用 Read 工具）
- 执行具体开发任务（使用 development-workflow skill）

## 核心概念

### 数据流

```
MediaSource → MediaAsset → Operator → Workflow → Task → Artifact
   (媒体源)    (媒体资产)    (算子)     (工作流)   (任务)   (产物)
```

### 关键实体

| 实体 | 作用 | 属性示例 |
|------|------|----------|
| **MediaSource** | 媒体来源（流/上传） | type(pull/push/upload), protocol(rtsp/rtmp/hls/webrtc/file) |
| **MediaAsset** | 媒体资产管理 | type(video/image/audio), source_type, parent_id(派生追踪), tags |
| **Operator** | AI/媒体处理单元 | category(analysis/processing/generation/utility), type(frame_extract/object_detection/ocr/...), endpoint, input_schema, output_spec, status(enabled/disabled/draft) |
| **Workflow** | DAG 编排 | trigger(manual/schedule/event), nodes, edges |
| **Task** | 工作流执行实例 | status(pending/running/success/failed/cancelled), progress, current_node, asset_id |
| **Artifact** | 算子输出产物 | type(asset/result/timeline/diagnostic), data |

### 废弃概念（V1.0 不再使用）

- ❌ Stream → 升级为 MediaSource
- ❌ Algorithm → 升级为 Operator
- ❌ AlgorithmBinding → 由 Workflow 替代
- ❌ InferenceResult → 由 Artifact 替代

## 分层架构（Clean Architecture）

```
internal/
├── domain/      # 核心实体（无外部依赖）
│   └── 实体：MediaSource, MediaAsset, Operator, Workflow, Task, Artifact, User, Role, Menu
│
├── port/        # 接口定义（契约）
│   └── 接口：Repository, OperatorPort, WorkflowEngine, MediaMTXClient
│
├── app/         # 业务服务（CQRS 模式）
│   ├── command/    # 命令（写操作：创建、更新、删除）
│   ├── query/      # 查询（读操作：查询、列表）
│   ├── dto/        # 数据传输对象
│   ├── port/       # 应用端口接口（MediaGateway, ObjectStorage, TokenService, EventBus, UnitOfWork）
│   ├── artifact.go          # 产物管理
│   ├── file.go              # 文件管理服务
│   ├── user_management.go   # 用户管理服务
│   └── workflow_scheduler.go # 工作流调度器
│
├── adapter/     # 基础设施实现
│   ├── persistence/   # GORM + PostgreSQL
│   ├── mediamtx/      # MediaMTX HTTP 客户端
│   ├── engine/        # DAG 工作流引擎
│   └── ai/            # 算子 HTTP 客户端
│
└── api/         # HTTP 表现层
    ├── handler/       # 请求处理器
    ├── dto/           # 数据传输对象
    ├── middleware/    # 中间件（JWT 认证、权限校验）
    └── router.go      # 路由注册
```

### 依赖规则（严格遵守）

```
✅ 正确：App → Port Interface → Adapter Implementation
❌ 错误：App → Adapter directly

依赖流：
- Domain: 不依赖任何层
- Port: 可依赖 Domain
- App: 可依赖 Domain + Port（禁止依赖 Adapter）
- Adapter: 实现 Port，可依赖 Domain
- API: 可依赖 App + Port + Domain（禁止依赖 Adapter）
```

## 算子标准协议

所有算子必须遵循统一的 I/O 协议，确保互操作性。

### 输入格式

```json
{
  "asset_id": "资产 UUID",
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

### 产物类型说明

- **output_assets**: 新生成的媒体资产（剪辑视频、检测结果图片）
- **results**: 结构化结果（检测框、分类标签、OCR 文本）
- **timeline**: 时间轴片段（事件、高光、镜头切分）
- **diagnostics**: 诊断信息（性能指标、模型版本）

## API 端点（前缀：/api/v1）

### 认证（Auth）
- `POST /auth/login` - 登录
- `POST /auth/refresh` - 刷新 Token
- `GET /auth/profile` - 获取当前用户
- `PUT /auth/password` - 修改密码
- `POST /auth/logout` - 登出

### 媒体源（Sources）
- `GET|POST /sources` - 列表、创建
- `GET|PUT|DELETE /sources/:id` - 详情、更新、删除
- `POST /sources/:id/enable|disable` - 启用、禁用
- `GET /sources/:id/status` - 获取实时状态
- `GET /sources/:id/preview` - 获取预览 URL

### 录制（Record）
- `POST /sources/:id/record/start|stop` - 启动、停止录制
- `GET /sources/:id/record/status|sessions|files` - 状态、会话、文件列表

### 点播（Playback）
- `GET /sources/:id/playback?start=<timestamp>` - 获取点播 URL
- `GET /sources/:id/playback/segments` - 列出录制段

### 媒体资产（Assets）
- `GET|POST /assets` - 列表、创建（支持过滤：type, source_type, tags）
- `GET|PUT|DELETE /assets/:id` - 详情、更新、删除
- `GET /assets/:id/children` - 列出派生资产

### 算子（Operators）
- `GET|POST /operators` - 列表、创建（支持过滤：category, status, is_builtin）
- `GET|PUT|DELETE /operators/:id` - 详情、更新、删除
- `POST /operators/:id/enable|disable|test` - 启用、禁用、测试

### 工作流（Workflows）
- `GET|POST /workflows` - 列表、创建
- `GET|PUT|DELETE /workflows/:id` - 详情、更新、删除
- `POST /workflows/:id/activate|pause|validate` - 启用、暂停、验证

### 任务（Tasks）
- `GET|POST /tasks` - 列表、创建（支持过滤：workflow_id, status, trigger_type）
- `GET /tasks/:id` - 详情
- `POST /tasks/:id/cancel|retry` - 取消、重试
- `GET /tasks/:id/logs` - 获取日志

### 产物（Artifacts）
- `GET /artifacts` - 列表（支持过滤：task_id, node_id, operator_id, type）
- `GET|DELETE /artifacts/:id` - 详情、删除
- `GET /artifacts/:id/download` - 下载

### 用户管理（Users/Roles/Menus）
- `GET|POST /users` - 用户管理
- `GET|POST /roles` - 角色管理
- `GET|POST /menus` - 菜单管理
- `GET /menus/tree` - 菜单树
- `GET /permissions` - 权限列表

### 文件管理（Files）
- `GET|POST /files` - 列表、上传
- `GET|PUT|DELETE /files/:id` - 详情、更新、删除
- `GET /files/:id/download` - 下载

## 配置管理

### 主配置文件：`configs/config.<env>.yaml`

```yaml
server:
  port: 8080

db:
  dsn: "host=localhost user=goyavision password=goyavision dbname=goyavision port=5432 sslmode=disable"

jwt:
  secret: "your-secret-key-change-in-production"
  expire: 2h
  refresh_exp: 168h

mediamtx:
  api_address: "http://localhost:9997"
  rtsp_address: "rtsp://localhost:8554"
  rtmp_address: "rtmp://localhost:1935"
  hls_address: "http://localhost:8888"
  webrtc_address: "http://localhost:8889"
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

所有配置项可通过环境变量覆盖，前缀为 `GOYAVISION_`：

```bash
export GOYAVISION_DB_DSN="host=localhost ..."
export GOYAVISION_JWT_SECRET="your-production-secret"
export GOYAVISION_MEDIAMTX_API_ADDRESS="http://mediamtx:9997"
```

## 开发状态（V1.0）

### ✅ 已完成
- MediaMTX 集成（RTSP/RTMP/HLS/WebRTC）
- 媒体源管理（拉流/推流）
- 录制与点播（集成 MediaMTX）
- JWT 认证（Access Token + Refresh Token 双 Token）
- RBAC 权限模型（用户、角色、菜单）
- 分层架构（Domain/Port/App/Adapter/API）
- Docker Compose 部署

### 🚧 进行中
- 媒体资产管理（CRUD、搜索、派生追踪）
- 算子管理（CRUD、分类、版本管理）
- 简化工作流（Phase 1：单算子任务）
- 任务调度与执行
- 产物管理
- 前端页面（资产、算子、工作流、任务）

### ⏸️ 待开始
- 可视化工作流设计器
- 更多内置算子（编辑、生成、转换类）
- 复杂工作流（DAG 编排、并行、条件分支）
- 自定义算子（Docker 镜像上传）
- 多租户支持
- 监控与告警（Prometheus + Grafana）

## 关键文档

| 文档 | 路径 | 用途 |
|------|------|------|
| 需求文档 | `docs/requirements.md` | 功能规格与验收标准 |
| 架构文档 | `docs/architecture.md` | 系统设计详细说明 |
| 开发进度 | `docs/development-progress.md` | 实现状态追踪 |
| API 文档 | `docs/api.md` | RESTful API 参考 |
| 变更日志 | `CHANGELOG.md` | 版本历史（关注 [未发布] 章节） |
| 部署指南 | `docs/DEPLOYMENT.md` | 部署与运维 |
| Claude 指南 | `CLAUDE.md` | Claude Code 使用指南 |

## 常见开发模式

### 添加新实体

1. 在 `internal/domain/` 定义实体
2. 在 `internal/port/` 定义 Repository 接口
3. 在 `internal/adapter/persistence/` 实现 Repository
4. 在 `internal/app/` 实现 Service
5. 在 `internal/api/dto/` 定义 DTO
6. 在 `internal/api/handler/` 实现 Handler
7. 在 `internal/api/router.go` 注册路由

### 实现新算子

1. 实现算子 HTTP 服务（符合标准 I/O 协议）
2. 在算子中心注册（code、category、version、endpoint、input_spec、output_spec）
3. 在工作流中使用算子

### 创建工作流

1. 定义工作流（name、description、trigger）
2. 添加节点（operator_id、params、retry、timeout）
3. 添加边（from、to、condition）
4. 验证 DAG（无环、连通）
5. 启用工作流

## 重要注意事项

⚠️ **关键约束**：
1. V1.0 不向后兼容旧版本
2. 资产派生追踪使用 `parent_id`（原始视频 → 抽帧图片 → 检测结果图片）
3. 算子必须无状态、幂等执行
4. 工作流 DAG 必须无环、连通
5. 节点失败不影响其他独立分支
6. 产物可通过 `asset_id` 关联新资产

## 默认凭证

- **用户名**: admin
- **密码**: admin123
- **角色**: 超级管理员（拥有所有权限）

⚠️ **安全警告**: 生产环境必须立即修改默认密码！

## 快速参考

**服务端口**:
- 8080: GoyaVision (Web UI + API)
- 5432: PostgreSQL
- 8554: MediaMTX RTSP
- 1935: MediaMTX RTMP
- 8888: MediaMTX HLS
- 8889: MediaMTX WebRTC
- 9997: MediaMTX API

**构建命令**:
- `make build` - 构建后端
- `make build-web` - 构建前端
- `make build-all` - 构建全部
- `make clean` - 清理构建产物

**技术栈**:
- 后端: Go 1.22+, Echo v4, GORM, PostgreSQL, Viper, JWT
- 流媒体: MediaMTX, FFmpeg
- 前端: Vue 3, TypeScript, Vite, Element Plus, Tailwind CSS
- 部署: Docker, Docker Compose

## 使用示例

```bash
# 在开始实现媒体资产管理功能前
/goyavision-context

# 快速查看：
# - MediaAsset 实体定义和属性
# - API 端点 GET|POST /assets
# - 分层架构中的位置
# - 已实现功能状态（🚧 进行中）
```

## 相关 Skills

- `/development-workflow` - 开发工作流（开始/完成开发）
- `/create-entity` - 创建新领域实体
- `/create-operator` - 创建新算子
- `/review-architecture` - 架构合规性审查
