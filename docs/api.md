# GoyaVision V1.0 API 文档

## 概述

GoyaVision API 遵循 RESTful 设计原则，所有 API 端点以 `/api/v1` 为前缀。

## 认证

使用 JWT（JSON Web Token）进行认证。除登录和刷新 Token 接口外，所有 API 请求需要在请求头中携带 Access Token：

```http
Authorization: Bearer <access_token>
```

### 认证流程

1. 调用登录接口获取 `access_token` 和 `refresh_token`
2. 使用 `access_token` 访问其他 API
3. 当 `access_token` 过期时，使用 `refresh_token` 获取新的 Token
4. 当 `refresh_token` 也过期时，需要重新登录
5. Access Token 内部包含 `token_type=access`，Refresh Token 包含 `token_type=refresh`

## 通用约定

### 分页参数

```
limit: 每页数量（默认 20，最大 1000）
offset: 偏移量（默认 0）
```

### 分页响应

```json
{
  "items": [...],
  "total": 100
}
```

### 错误响应

API 统一返回错误信封，包含错误码、消息与请求追踪信息：

```json
{
  "code": 40000,
  "message": "错误详情",
  "details": {},
  "request_id": "req-xxx",
  "timestamp": 1700000000
}
```

参数校验失败等基础错误可能返回简化格式：

```json
{
  "error": "Bad Request",
  "message": "invalid request body"
}
```

HTTP 状态码：
- `200`：成功
- `201`：创建成功
- `204`：删除成功（无响应体）
- `400`：请求错误（Bad Request）
- `401`：未认证（Unauthorized）
- `403`：无权限（Forbidden）
- `404`：资源不存在（Not Found）
- `409`：资源冲突（Conflict）
- `500`：服务器内部错误

## API 端点

### 认证（Auth）

#### 登录

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

**响应**：
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 7200,
  "token_type": "access",
  "user": {
    "id": "uuid",
    "username": "admin",
    "nickname": "管理员",
    "email": "",
    "phone": "",
    "avatar": "",
    "roles": ["super_admin"],
    "permissions": ["*"],
    "menus": [...]
  }
}
```

#### 刷新 Token

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### 获取当前用户信息

```http
GET /api/v1/auth/profile
Authorization: Bearer <access_token>
```

#### 修改密码

```http
PUT /api/v1/auth/password
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "old_password": "admin123",
  "new_password": "newpassword"
}
```

#### 登出

```http
POST /api/v1/auth/logout
Authorization: Bearer <access_token>
```

---

### 媒体源（Sources）

媒体源与 MediaMTX path 一一对应，用于流媒体资产接入。设计详见 `docs/stream-asset-mediamtx-design.md`。

#### 已实现的端点

##### 列出媒体源

```http
GET /api/v1/sources?type=pull&limit=20&offset=0
```

**查询参数**：
- `type`（可选）：源类型过滤（pull、push）
- `limit`（可选）：每页数量，默认 20
- `offset`（可选）：偏移量，默认 0

**响应**：
```json
{
  "items": [
    {
      "id": "uuid",
      "name": "camera1",
      "path_name": "live/camera1-a1b2c3d4",
      "type": "pull",
      "url": "rtsp://192.168.1.100:554/stream",
      "protocol": "rtsp",
      "enabled": true,
      "record_enabled": false,
      "created_at": "2025-02-01T10:00:00Z",
      "updated_at": "2025-02-01T10:00:00Z"
    }
  ],
  "total": 1
}
```

##### 创建媒体源

拉流（type=pull 时必填 url）；推流（type=push 时 url 可为空，创建后通过预览接口获取推流地址）。

```http
POST /api/v1/sources
Content-Type: application/json

{
  "name": "camera1",
  "type": "pull",
  "url": "rtsp://192.168.1.100:554/stream",
  "protocol": "rtsp",
  "enabled": true
}
```

##### 获取媒体源详情

```http
GET /api/v1/sources/:id
```

##### 更新媒体源

```http
PUT /api/v1/sources/:id
Content-Type: application/json

{
  "name": "camera1_updated",
  "url": "rtsp://192.168.1.101:554/stream",
  "enabled": true
}
```

##### 删除媒体源

仅当无流媒体资产关联该源时允许；有关联时返回 409。

```http
DELETE /api/v1/sources/:id
```

##### 获取预览 URL

type=push 时响应中同时包含 `push_url`（推流地址，供 OBS 等配置）。

```http
GET /api/v1/sources/:id/preview
```

**响应**：
```json
{
  "path_name": "live/camera1-a1b2c3d4",
  "hls_url": "http://localhost:8888/live/camera1-a1b2c3d4/index.m3u8",
  "rtsp_url": "rtsp://localhost:8554/live/camera1-a1b2c3d4",
  "rtmp_url": "rtmp://localhost:1935/live/camera1-a1b2c3d4",
  "webrtc_url": "http://localhost:8889/live/camera1-a1b2c3d4/whep",
  "push_url": "rtmp://localhost:1935/live/camera1-a1b2c3d4"
}
```

（`push_url` 仅当 type=push 时存在。）

#### 计划实现的端点（当前未实现）

以下按设计文档 2.3 节规划，后续迭代实现：列表/详情 with_status、enable/disable、status、preview/ready、录制与点播相关端点。

---

### 媒体资产（Assets）

#### 列出媒体资产

```http
GET /api/v1/assets?type=video&source_type=upload&tags=tag1,tag2&limit=20&offset=0
```

**查询参数**：
- `type`（可选）：资产类型（video、image、audio）
- `source_type`（可选）：来源类型（upload、generated、operator_output）
- `source_id`（可选）：媒体源 ID
- `parent_id`（可选）：父资产 ID
- `tags`（可选）：标签列表（逗号分隔）
- `from`（可选）：开始时间戳
- `to`（可选）：结束时间戳
- `limit`（可选）：每页数量
- `offset`（可选）：偏移量

**响应**：
```json
{
  "items": [
    {
      "id": "uuid",
      "type": "video",
      "source_type": "upload",
      "source_id": null,
      "parent_id": null,
      "name": "uploaded_video.mp4",
      "path": "http://minio:9000/goyavision/uploads/uploaded_video.mp4",
      "duration": 3600,
      "size": 104857600,
      "format": "mp4",
      "metadata": {
        "resolution": "1920x1080",
        "fps": 30,
        "codec": "h264"
      },
      "status": "ready",
      "tags": ["demo"],
      "created_at": "2025-02-01T10:00:00Z",
      "updated_at": "2025-02-01T10:00:00Z"
    }
  ],
  "total": 100
}
```

#### 创建媒体资产

```http
POST /api/v1/assets
Content-Type: application/json

{
  "type": "video",
  "source_type": "upload",
  "name": "uploaded_video.mp4",
  "path": "uploads/uploaded_video.mp4",
  "duration": 120,
  "size": 10485760,
  "format": "mp4",
  "metadata": {
    "resolution": "1920x1080",
    "fps": 30
  },
  "tags": ["upload", "demo"]
}
```

**字段说明**：
- `type`（必填）：资产类型，可选值 `video`、`image`、`audio`
- `source_type`（必填）：来源类型，可选值 `upload`（用户上传）、`generated`（系统生成）、`operator_output`（算子输出）
- `path`（必填）：资源路径。对于 `upload`/`generated`/`operator_output` 类型，可传 MinIO 相对路径，响应中会返回完整 URL

> **注意**：流媒体接入功能已迁移至媒体源（Sources）模块，资产模块不再支持 `type=stream` 和 `source_type=live/vod`。

#### 获取媒体资产详情

```http
GET /api/v1/assets/:id
```

#### 更新媒体资产

```http
PUT /api/v1/assets/:id
Content-Type: application/json

{
  "name": "new_name.mp4",
  "tags": ["tag1", "tag2"]
}
```

**权限要求**：
- 需要 `asset:update` 权限
- 后端会基于 JWT + RBAC 强校验，前端权限控制仅用于交互层隐藏/禁用

**无权限响应**：
```json
{
  "error": "Forbidden",
  "message": "无编辑权限"
}
```

#### 删除媒体资产

```http
DELETE /api/v1/assets/:id
```

#### 列出子资产（派生资产）

```http
GET /api/v1/assets/:id/children
```

**响应**：
```json
[
  {
    "id": "uuid",
    "type": "image",
    "source_type": "generated",
    "parent_id": "uuid",
    "name": "frame_001.jpg",
    "path": "./data/frames/uuid/frame_001.jpg",
    "format": "jpg",
    "status": "ready",
    "created_at": "2025-02-01T10:01:00Z"
  }
]
```

---

### 算子（Operators）

#### 列出算子

```http
GET /api/v1/operators?category=analyze&status=enabled
```

**查询参数**：
- `category`（可选）：算子分类（analyze、edit、generate、transform）
- `status`（可选）：状态（enabled、disabled）
- `is_builtin`（可选）：是否内置（true、false）

**响应**：
```json
[
  {
    "id": "uuid",
    "code": "frame_extract",
    "name": "抽帧",
    "category": "analyze",
    "version": "1.0.0",
    "input_spec": {
      "asset_types": ["video"],
      "params": {
        "interval": {
          "type": "integer",
          "description": "抽帧间隔（秒）",
          "default": 30
        }
      }
    },
    "output_spec": {
      "assets": ["image"],
      "results": [],
      "timeline": []
    },
    "endpoint": "http://localhost:8080/internal/operators/frame_extract",
    "status": "enabled",
    "is_builtin": true,
    "description": "从视频中按间隔提取帧",
    "icon": "frame",
    "created_at": "2025-02-01T10:00:00Z",
    "updated_at": "2025-02-01T10:00:00Z"
  }
]
```

#### 创建算子

```http
POST /api/v1/operators
Content-Type: application/json

{
  "code": "custom_detector",
  "name": "自定义检测器",
  "category": "analyze",
  "version": "1.0.0",
  "input_spec": {
    "asset_types": ["image"],
    "params": {
      "threshold": {
        "type": "float",
        "description": "检测阈值",
        "default": 0.5
      }
    }
  },
  "output_spec": {
    "assets": [],
    "results": ["detection"],
    "timeline": []
  },
  "endpoint": "http://ai-service:8080/detect",
  "config": {
    "timeout": 10,
    "retry": 2
  },
  "status": "enabled",
  "is_builtin": false,
  "description": "自定义目标检测算子"
}
```

#### 获取算子详情

```http
GET /api/v1/operators/:id
```

#### 更新算子

```http
PUT /api/v1/operators/:id
Content-Type: application/json

{
  "name": "自定义检测器 V2",
  "version": "2.0.0",
  "endpoint": "http://ai-service:8080/detect-v2"
}
```

#### 删除算子

```http
DELETE /api/v1/operators/:id
```

#### 发布算子

```http
POST /api/v1/operators/:id/publish
```

#### 弃用算子

```http
POST /api/v1/operators/:id/deprecate
```

#### 测试算子

```http
POST /api/v1/operators/:id/test
Content-Type: application/json

{
  "asset_id": "uuid",
  "params": {
    "threshold": 0.7
  }
}
```

**响应**：
```json
{
  "success": true,
  "message": "ok",
  "diagnostics": {
    "operator_id": "uuid",
    "status": "published",
    "checked": true
  }
}
```

#### MCP Server 列表

```http
GET /api/v1/operators/mcp/servers
```

**响应**：
```json
[
  {
    "id": "default",
    "name": "默认 MCP",
    "description": "MCP server",
    "status": "online"
  }
]
```

#### MCP Tool 列表

```http
GET /api/v1/operators/mcp/servers/:id/tools
```

**响应**：
```json
[
  {
    "name": "tool_name",
    "description": "工具描述",
    "version": "1.0.0",
    "input_schema": {},
    "output_schema": {}
  }
]
```

#### MCP Tool 预览

```http
GET /api/v1/operators/mcp/servers/:id/tools/:tool/preview
```

**响应**：
```json
{
  "name": "tool_name",
  "description": "工具描述",
  "version": "1.0.0",
  "input_schema": {},
  "output_schema": {}
}
```

#### 从 MCP Tool 安装算子

```http
POST /api/v1/operators/mcp/install
Content-Type: application/json

{
  "server_id": "default",
  "tool_name": "tool_name",
  "operator_code": "mcp_tool_name",
  "operator_name": "MCP Tool Name",
  "category": "utility",
  "type": "transcode",
  "timeout_sec": 30,
  "tags": ["mcp"]
}
```

**响应**：返回标准 `OperatorResponse`。

#### 同步 MCP 模板

```http
POST /api/v1/operators/mcp/sync-templates
Content-Type: application/json

{
  "server_id": "default"
}
```

**响应**：
```json
{
  "server_id": "default",
  "total": 10,
  "created": 8,
  "updated": 2
}
```

---

### 工作流（Workflows）

#### 列出工作流

```http
GET /api/v1/workflows?status=active
```

**查询参数**：
- `status`（可选）：状态（draft、active、paused）

**响应**：
```json
[
  {
    "id": "uuid",
    "name": "实时视频分析",
    "description": "对实时视频流进行目标检测和追踪",
    "trigger": {
      "type": "schedule",
      "config": {
        "interval": "30s"
      }
    },
    "nodes": [
      {
        "id": "node1",
        "operator_id": "uuid",
        "params": {
          "interval": 30
        },
        "retry": 3,
        "timeout": 300
      }
    ],
    "edges": [],
    "status": "active",
    "created_at": "2025-02-01T10:00:00Z",
    "updated_at": "2025-02-01T10:00:00Z"
  }
]
```

#### 创建工作流

```http
POST /api/v1/workflows
Content-Type: application/json

{
  "name": "实时视频分析",
  "description": "对实时视频流进行目标检测和追踪",
  "trigger": {
    "type": "schedule",
    "config": {
      "interval": "30s"
    }
  },
  "nodes": [
    {
      "id": "node1",
      "operator_id": "uuid",
      "params": {
        "interval": 30
      },
      "retry": 3,
      "timeout": 300
    },
    {
      "id": "node2",
      "operator_id": "uuid2",
      "params": {
        "threshold": 0.7
      },
      "retry": 2,
      "timeout": 60
    }
  ],
  "edges": [
    {
      "from": "node1",
      "to": "node2"
    }
  ],
  "status": "draft"
}
```

#### 获取工作流详情

```http
GET /api/v1/workflows/:id
```

#### 更新工作流

```http
PUT /api/v1/workflows/:id
Content-Type: application/json

{
  "name": "实时视频分析 V2",
  "status": "active"
}
```

#### 删除工作流

```http
DELETE /api/v1/workflows/:id
```

#### 启用工作流

```http
POST /api/v1/workflows/:id/enable
```

#### 禁用工作流

```http
POST /api/v1/workflows/:id/disable
```

#### 手动触发工作流

```http
POST /api/v1/workflows/:id/trigger
Content-Type: application/json

{
  "asset_id": "uuid"
}
```

---

### 任务（Tasks）

#### 列出任务

```http
GET /api/v1/tasks?workflow_id=uuid&status=running&limit=20&offset=0
```

**查询参数**：
- `workflow_id`（可选）：工作流 ID
- `status`（可选）：状态（pending、running、success、failed、cancelled）
- `from`（可选）：开始时间戳
- `to`（可选）：结束时间戳
- `limit`（可选）：每页数量
- `offset`（可选）：偏移量

**响应**：
```json
{
  "items": [
    {
      "id": "uuid",
      "workflow_id": "uuid",
      "asset_id": "uuid",
      "status": "running",
      "progress": 50,
      "current_node": "node2",
      "started_at": "2025-02-01T10:00:00Z",
      "completed_at": null,
      "error": null,
      "created_at": "2025-02-01T10:00:00Z"
    }
  ],
  "total": 100
}
```

#### 创建任务

```http
POST /api/v1/tasks
Content-Type: application/json

{
  "workflow_id": "uuid",
  "asset_id": "uuid",
  "input_params": {}
}
```

#### 获取任务详情

```http
GET /api/v1/tasks/:id
```

**响应**：
```json
{
  "id": "uuid",
  "workflow_id": "uuid",
  "workflow": {...},
  "asset_id": "uuid",
  "input_params": {},
  "status": "success",
  "progress": 100,
  "current_node": null,
  "started_at": "2025-02-01T10:00:00Z",
  "completed_at": "2025-02-01T10:05:00Z",
  "error": null,
  "artifacts": [...],
  "created_at": "2025-02-01T10:00:00Z"
}
```

#### 取消任务

```http
POST /api/v1/tasks/:id/cancel
```

#### 获取任务统计

```http
GET /api/v1/tasks/stats
```

---

### 产物（Artifacts）

#### 列出产物

```http
GET /api/v1/artifacts?task_id=uuid&type=result&limit=20&offset=0
```

**查询参数**：
- `task_id`（可选）：任务 ID
- `node_id`（可选）：节点 ID
- `operator_id`（可选）：算子 ID
- `type`（可选）：产物类型（asset、result、timeline、diagnostic）
- `from`（可选）：开始时间戳
- `to`（可选）：结束时间戳
- `limit`（可选）：每页数量
- `offset`（可选）：偏移量

**响应**：
```json
{
  "items": [
    {
      "id": "uuid",
      "task_id": "uuid",
      "node_id": "node1",
      "operator_id": "uuid",
      "type": "result",
      "asset_id": null,
      "data": {
        "type": "detection",
        "detections": [
          {
            "class": "person",
            "confidence": 0.95,
            "bbox": [100, 200, 300, 400]
          }
        ]
      },
      "created_at": "2025-02-01T10:01:00Z"
    }
  ],
  "total": 100
}
```

#### 获取产物详情

```http
GET /api/v1/artifacts/:id
```

#### 删除产物

```http
DELETE /api/v1/artifacts/:id
```

#### 下载产物

```http
GET /api/v1/artifacts/:id/download
```

---

### 用户管理（Users）

#### 列出用户

```http
GET /api/v1/users?status=1&limit=20&offset=0
```

**查询参数**：
- `status`（可选）：状态（1=启用，0=禁用）
- `limit`（可选）：每页数量
- `offset`（可选）：偏移量

#### 创建用户

```http
POST /api/v1/users
Content-Type: application/json

{
  "username": "user1",
  "password": "password123",
  "nickname": "用户1",
  "email": "user1@example.com",
  "phone": "13800138000",
  "status": 1,
  "role_ids": ["uuid"]
}
```

#### 更新用户

```http
PUT /api/v1/users/:id
Content-Type: application/json

{
  "nickname": "新昵称",
  "status": 1,
  "role_ids": ["uuid1", "uuid2"]
}
```

#### 删除用户

```http
DELETE /api/v1/users/:id
```

#### 重置密码

```http
POST /api/v1/users/:id/reset-password
Content-Type: application/json

{
  "new_password": "newpassword"
}
```

---

### 角色管理（Roles）

#### 列出角色

```http
GET /api/v1/roles?status=1
```

#### 创建角色

```http
POST /api/v1/roles
Content-Type: application/json

{
  "code": "operator",
  "name": "操作员",
  "description": "普通操作员",
  "status": 1,
  "permission_ids": ["uuid1", "uuid2"],
  "menu_ids": ["uuid1", "uuid2"]
}
```

#### 更新角色

```http
PUT /api/v1/roles/:id
Content-Type: application/json

{
  "name": "新名称",
  "permission_ids": ["uuid1", "uuid2"],
  "menu_ids": ["uuid1", "uuid2"]
}
```

#### 删除角色

```http
DELETE /api/v1/roles/:id
```

---

### 权限管理（Permissions）

#### 列出所有权限

```http
GET /api/v1/permissions
```

**响应**：
```json
[
  {
    "id": "uuid",
    "code": "asset:list",
    "name": "查看媒体资产列表",
    "method": "GET",
    "path": "/api/v1/assets",
    "description": ""
  }
]
```

---

### 菜单管理（Menus）

#### 列出菜单

```http
GET /api/v1/menus?status=1
```

#### 列出菜单树

```http
GET /api/v1/menus/tree?status=1
```

#### 创建菜单

```http
POST /api/v1/menus
Content-Type: application/json

{
  "parent_id": "uuid",
  "code": "asset:list",
  "name": "媒体资产",
  "type": 2,
  "path": "/asset/list",
  "icon": "VideoCamera",
  "component": "asset/index",
  "permission": "asset:list",
  "sort": 1,
  "visible": true,
  "status": 1
}
```

菜单类型：
- `1`：目录
- `2`：菜单
- `3`：按钮

#### 更新菜单

```http
PUT /api/v1/menus/:id
Content-Type: application/json

{
  "name": "新名称",
  "sort": 2
}
```

#### 删除菜单

```http
DELETE /api/v1/menus/:id
```

---

## 算子标准协议

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

---

## 前端静态文件

前端 SPA 通过根路径 `/` 访问，所有非 API 路由返回 `index.html`。

---

**注意**：本文档会随着 API 演进持续更新。
