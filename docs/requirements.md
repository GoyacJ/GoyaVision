# GoyaVision 需求文档

## 1. 项目与范围

- **项目名**：GoyaVision（AI 视频流分析处理平台）
- **范围**：视频流接入、抽帧、录制、AI 模型调用、分析任务配置、管理前端。

## 2. 功能需求

### 2.1 视频流

- 支持 **RTSP** 接入。
- 单实例支持 **10+ 路**并发流。
- 流的 CRUD；支持按 `enabled` 过滤。

### 2.2 抽帧与录制

- **抽帧**：从 RTSP 按设定频率取帧，用于 AI 推理；与录制解耦。
- **录制**：与分析独立开关；按流启停，分段落盘（如 5 分钟/段）；`-c copy` 不重编码。

### 2.3 AI 模型调用

- 通过 **HTTP + JSON** 调用推理服务；算法配置中维护 `endpoint`、`input_spec`、`output_spec`。
- 输入形态：v1 以 **base64 图片** 为主。
- 推理结果**持久化**到 DB，支持按流、算法、时间查询。

### 2.4 分析任务配置（无「任务」概念）

- **一流多算法**：一个流可绑定多个算法（Stream → AlgorithmBinding → Algorithm），无中间「Task」。
- **处理频率**：`AlgorithmBinding.interval_sec` 表示抽帧 + 模型调用周期（每 N 秒 1 帧 1 次推理）。
- **定时**：`schedule`（JSON）：`{start, end, days_of_week}`，仅在指定时段内运行；为空则始终按 `interval_sec` 运行。
- **延迟**：`initial_delay_sec`：绑定生效后，首次推理前等待秒数。

### 2.5 预览

- 需要**实时预览**画面。
- 前端通过 HLS 播放（video.js）；后端经 MediaMTX 或 FFmpeg 将 RTSP 转 HLS，并提供 URL。

### 2.6 前端

- **Vue 3 + TypeScript + Vite + Element Plus + video.js**。
- 构建产物**内嵌**到 Go 单体（`embed`），由 Echo 提供 SPA 与 `/live`。
- 页面：流列表（启停、预览）、流的算法绑定（`interval_sec`、`schedule`、`initial_delay_sec`、`enabled`）、录制启停、推理结果查询、HLS 预览。

## 3. 非功能需求

### 3.1 资源与规模

- 单实例 10+ 路 RTSP；需 FFmpeg / 预览进程池与上限，避免无界拉起。

### 3.2 部署

- 单二进制（Go + embed 前端）；PostgreSQL、FFmpeg、MediaMTX 需预装或可配置路径。

### 3.3 运维（规划）

- 健康检查：`/health`、`/ready`；指标：`/metrics`（Prometheus）；结构化日志。
- 认证与鉴权：生产需 API Key 或 JWT（具体方案待定）。

## 4. 数据模型概要

- **Stream**：id, url, name, enabled, created_at, updated_at
- **Algorithm**：id, name, endpoint, input_spec, output_spec, created_at, updated_at
- **AlgorithmBinding**：id, stream_id, algorithm_id, enabled, interval_sec, initial_delay_sec, schedule, config, created_at, updated_at
- **RecordSession**：id, stream_id, status, base_path, started_at, stopped_at
- **InferenceResult**：id, algorithm_binding_id, stream_id, ts, frame_ref, output, latency_ms, created_at

## 5. API 概要

| 资源            | 方法 | 路径 |
|-----------------|------|------|
| Stream          | CRUD | /api/v1/streams |
| Algorithm       | CRUD | /api/v1/algorithms |
| AlgorithmBinding| CRUD | /api/v1/streams/:id/algorithm-bindings |
| Record          | POST | /api/v1/streams/:id/record/start, .../stop |
| Record          | GET  | /api/v1/streams/:id/record/sessions |
| InferenceResult | GET  | /api/v1/inference_results |
| Preview         | GET  | /api/v1/streams/:id/preview/start |
| Preview         | POST | /api/v1/streams/:id/preview/stop |

## 6. 变更与维护

- 本文档随迭代更新；重大变更需在「变更记录」中注明日期与内容。
- 变更记录（示例）：
  - 2025-01：初稿；去掉 Task，一流多算法，算法支持 schedule、initial_delay_sec。
