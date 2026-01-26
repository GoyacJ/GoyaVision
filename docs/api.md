# GoyaVision API 文档

## 概述

GoyaVision API 遵循 RESTful 设计原则，所有 API 端点以 `/api/v1` 为前缀。

## 认证

当前版本暂未实现认证，生产环境需要添加 API Key 或 JWT 认证。

## 错误响应

所有错误响应遵循统一格式：

```json
{
  "error": "错误类型",
  "message": "错误详情"
}
```

HTTP 状态码：
- `200`：成功
- `201`：创建成功
- `204`：删除成功（无响应体）
- `400`：请求错误（Bad Request）
- `404`：资源不存在（Not Found）
- `409`：资源冲突（Conflict）
- `500`：服务器内部错误

## API 端点

### 视频流（Streams）

#### 列出流

```http
GET /api/v1/streams?enabled=true
```

**查询参数**：
- `enabled`（可选）：过滤启用/禁用的流（true/false）

**响应**：
```json
[
  {
    "id": "uuid",
    "url": "rtsp://example.com/stream",
    "name": "流名称",
    "enabled": true,
    "created_at": "2025-01-26T10:00:00Z",
    "updated_at": "2025-01-26T10:00:00Z"
  }
]
```

#### 创建流

```http
POST /api/v1/streams
Content-Type: application/json

{
  "url": "rtsp://example.com/stream",
  "name": "流名称",
  "enabled": true
}
```

#### 获取流

```http
GET /api/v1/streams/:id
```

#### 更新流

```http
PUT /api/v1/streams/:id
Content-Type: application/json

{
  "url": "rtsp://example.com/stream2",
  "name": "新名称",
  "enabled": false
}
```

#### 删除流

```http
DELETE /api/v1/streams/:id
```

### 算法（Algorithms）

#### 列出算法

```http
GET /api/v1/algorithms
```

#### 创建算法

```http
POST /api/v1/algorithms
Content-Type: application/json

{
  "name": "算法名称",
  "endpoint": "http://ai-service:8080/inference",
  "input_spec": {},
  "output_spec": {}
}
```

#### 获取算法

```http
GET /api/v1/algorithms/:id
```

#### 更新算法

```http
PUT /api/v1/algorithms/:id
Content-Type: application/json

{
  "name": "新名称",
  "endpoint": "http://ai-service:8080/inference2"
}
```

#### 删除算法

```http
DELETE /api/v1/algorithms/:id
```

### 算法绑定（Algorithm Bindings）

#### 列出流的算法绑定

```http
GET /api/v1/streams/:id/algorithm-bindings
```

#### 创建算法绑定

```http
POST /api/v1/streams/:id/algorithm-bindings
Content-Type: application/json

{
  "algorithm_id": "uuid",
  "interval_sec": 30,
  "initial_delay_sec": 0,
  "enabled": true,
  "schedule": {
    "start": "09:00:00",
    "end": "18:00:00",
    "days_of_week": [1, 2, 3, 4, 5]
  }
}
```

#### 获取算法绑定

```http
GET /api/v1/streams/:id/algorithm-bindings/:bid
```

#### 更新算法绑定

```http
PUT /api/v1/streams/:id/algorithm-bindings/:bid
Content-Type: application/json

{
  "interval_sec": 60,
  "enabled": false
}
```

#### 删除算法绑定

```http
DELETE /api/v1/streams/:id/algorithm-bindings/:bid
```

### 录制（Recording）

#### 启动录制

```http
POST /api/v1/streams/:id/record/start
```

**响应**：
```json
{
  "session_id": "uuid"
}
```

#### 停止录制

```http
POST /api/v1/streams/:id/record/stop
```

#### 查询录制会话

```http
GET /api/v1/streams/:id/record/sessions
```

**响应**：
```json
[
  {
    "id": "uuid",
    "stream_id": "uuid",
    "status": "running",
    "base_path": "./data/recordings/uuid",
    "started_at": "2025-01-26T10:00:00Z",
    "stopped_at": null
  }
]
```

### 推理结果（Inference Results）

#### 查询推理结果

```http
GET /api/v1/inference_results?stream_id=uuid&binding_id=uuid&from=1706256000&to=1706342400&limit=50&offset=0
```

**查询参数**：
- `stream_id`（可选）：流 ID
- `binding_id`（可选）：算法绑定 ID
- `from`（可选）：开始时间戳（Unix 秒）
- `to`（可选）：结束时间戳（Unix 秒）
- `limit`（可选）：每页数量，默认 50，最大 1000
- `offset`（可选）：偏移量，默认 0

**响应**：
```json
{
  "items": [
    {
      "id": "uuid",
      "algorithm_binding_id": "uuid",
      "stream_id": "uuid",
      "ts": "2025-01-26T10:00:00Z",
      "frame_ref": "./data/frames/uuid/frame_xxx.jpg",
      "output": {},
      "latency_ms": 150,
      "created_at": "2025-01-26T10:00:00Z"
    }
  ],
  "total": 100
}
```

### 预览（Preview）

#### 启动预览

```http
GET /api/v1/streams/:id/preview/start
```

**响应**：
```json
{
  "hls_url": "/live/uuid/index.m3u8"
}
```

#### 停止预览

```http
POST /api/v1/streams/:id/preview/stop
```

## HLS 文件服务

HLS 文件通过 `/live/*` 路径提供：

```http
GET /live/:stream_id/index.m3u8
GET /live/:stream_id/segment_001.ts
```

## 前端静态文件

前端 SPA 通过根路径 `/` 访问，所有非 API 路由返回 `index.html`。

---

**注意**：本文档会随着 API 演进持续更新。
