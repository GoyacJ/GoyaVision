# GoyaVision API 文档

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

**响应**：
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 7200
}
```

#### 获取当前用户信息（需认证）

```http
GET /api/v1/auth/profile
Authorization: Bearer <access_token>
```

#### 修改密码（需认证）

```http
PUT /api/v1/auth/password
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "old_password": "admin123",
  "new_password": "newpassword"
}
```

#### 登出（需认证）

```http
POST /api/v1/auth/logout
Authorization: Bearer <access_token>
```

### 用户管理（Users）- 需认证

#### 列出用户

```http
GET /api/v1/users?status=1&limit=20&offset=0
```

**查询参数**：
- `status`（可选）：状态过滤（1=启用，0=禁用）
- `limit`（可选）：每页数量，默认 20
- `offset`（可选）：偏移量

**响应**：
```json
{
  "items": [
    {
      "id": "uuid",
      "username": "admin",
      "nickname": "管理员",
      "email": "",
      "phone": "",
      "status": 1,
      "roles": [{"id": "uuid", "code": "super_admin", "name": "超级管理员"}],
      "created_at": "2025-01-26T10:00:00Z",
      "updated_at": "2025-01-26T10:00:00Z"
    }
  ],
  "total": 1
}
```

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

### 角色管理（Roles）- 需认证

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

### 菜单管理（Menus）- 需认证

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
  "code": "system:user",
  "name": "用户管理",
  "type": 2,
  "path": "/system/user",
  "icon": "User",
  "component": "system/user/index",
  "permission": "user:list",
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

### 权限列表（Permissions）- 需认证

#### 列出所有权限

```http
GET /api/v1/permissions
```

**响应**：
```json
[
  {
    "id": "uuid",
    "code": "stream:list",
    "name": "查看视频流列表",
    "method": "GET",
    "path": "/api/v1/streams",
    "description": ""
  }
]
```

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
