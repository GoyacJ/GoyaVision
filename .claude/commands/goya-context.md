# context

获取 GoyaVision 项目的完整上下文，包括架构设计、核心概念、API 端点和开发状态。

## 核心概念

### 数据流

```
MediaSource → MediaAsset → Operator → Workflow → Task → Artifact
   (媒体源)    (媒体资产)    (算子)     (工作流)   (任务)   (产物)
```

### 关键实体

| 实体 | 作用 | 关键属性 |
|------|------|----------|
| **MediaSource** | 媒体来源（流/上传） | type(pull/push/upload), protocol(rtsp/rtmp/hls/webrtc/file) |
| **MediaAsset** | 媒体资产管理 | type(video/image/audio), source_type, parent_id(派生追踪), tags |
| **Operator** | AI/媒体处理单元 | category(analyze/edit/generate/transform), endpoint, input_spec, output_spec |
| **Workflow** | DAG 编排 | trigger(manual/schedule/event), nodes, edges |
| **Task** | 工作流执行实例 | status(pending/running/completed/failed), progress, current_node |
| **Artifact** | 算子输出产物 | type(asset/result/timeline/diagnostic), data |

## 分层架构（Clean Architecture）

```
internal/
├── domain/      # 核心实体（无外部依赖）
├── port/        # 接口定义（契约）
├── app/         # 业务服务（用例编排）
├── adapter/     # 基础设施实现（persistence, mediamtx, engine, ai）
└── api/         # HTTP 表现层（handler, dto, middleware, router）
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

### 输入格式
```json
{
  "asset_id": "资产 UUID",
  "params": {"key": "value"}
}
```

### 输出格式
```json
{
  "output_assets": [{"type": "video|image|audio", "path": "...", "format": "...", "metadata": {}}],
  "results": [{"type": "detection|classification|...", "data": {}, "confidence": 0.95}],
  "timeline": [{"start": 0.0, "end": 5.0, "event_type": "...", "confidence": 0.95, "data": {}}],
  "diagnostics": {"latency_ms": 150, "model_version": "v1.0", "device": "gpu"}
}
```

## API 端点（前缀：/api/v1）

**认证**: `/auth/login`, `/auth/refresh`, `/auth/profile`, `/auth/password`, `/auth/logout`
**媒体源**: `/sources` (CRUD), `/sources/:id/enable|disable|status|preview`
**录制**: `/sources/:id/record/start|stop|status|sessions|files`
**点播**: `/sources/:id/playback`, `/sources/:id/playback/segments`
**媒体资产**: `/assets` (CRUD, 支持过滤), `/assets/:id/children`
**算子**: `/operators` (CRUD, 支持过滤), `/operators/:id/enable|disable|test`
**工作流**: `/workflows` (CRUD), `/workflows/:id/activate|pause|validate`
**任务**: `/tasks` (CRUD, 支持过滤), `/tasks/:id/cancel|retry|logs`
**产物**: `/artifacts` (列表, 支持过滤), `/artifacts/:id/download`
**用户管理**: `/users`, `/roles`, `/menus`, `/permissions`
**文件管理**: `/files` (CRUD, 上传, 下载)

## 开发状态（V1.0）

### ✅ 已完成
- MediaMTX 集成、媒体源管理、录制与点播
- JWT 认证（双 Token 机制）
- RBAC 权限模型
- 分层架构、Docker Compose 部署

### 🚧 进行中
- 媒体资产管理、算子管理
- 简化工作流（Phase 1：单算子任务）
- 任务调度与执行、产物管理
- 前端页面

### ⏸️ 待开始
- 可视化工作流设计器
- 复杂工作流（DAG 编排、并行、条件分支）
- 更多内置算子、自定义算子
- 多租户支持、监控与告警

## 配置快速参考

**默认凭证**: admin / admin123
**服务端口**: 8080 (API), 5432 (DB), 8554 (RTSP), 1935 (RTMP), 8888 (HLS), 8889 (WebRTC), 9997 (MediaMTX API)
**构建命令**: `make build`, `make build-web`, `make build-all`
**配置文件**: `configs/config.yaml`
**环境变量**: `GOYAVISION_*` 前缀（如 `GOYAVISION_DB_DSN`）

## 技术栈

**后端**: Go 1.22+, Echo v4, GORM, PostgreSQL, Viper, JWT
**流媒体**: MediaMTX, FFmpeg
**前端**: Vue 3, TypeScript, Vite, Element Plus, Tailwind CSS
**部署**: Docker, Docker Compose
