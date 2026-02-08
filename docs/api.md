# GoyaVision V1.0 API 文档

## 概述

GoyaVision API 遵循 RESTful 设计原则，所有 API 端点以 `/api/v1` 为前缀。

## 认证

使用 JWT (JSON Web Token) 进行认证。Access Token 需携带在 `Authorization` 头中：

```http
Authorization: Bearer <access_token>
```

- **Access Token**: `token_type=access`，有效期 2 小时。
- **Refresh Token**: `token_type=refresh`，有效期 7 天。

## 通用约定

### 多租户与权限
所有核心资源 DTO 包含以下公共字段：
- `tenant_id`: 归属租户 ID。
- `owner_id`: 创建者 ID。
- `visibility`: 可见性级别（0: 私有, 1: 租户公开, 2: 系统公开）。

### 错误响应
API 统一返回错误信封：

```json
{
  "code": 40000,
  "message": "错误详情",
  "request_id": "req-xxx",
  "timestamp": 1700000000
}
```

## API 端点

### 认证 (Auth)
- `POST /auth/login`: 登录获取双 Token。
- `POST /auth/refresh`: 使用 Refresh Token 刷新。
- `GET /auth/profile`: 获取当前用户信息、角色、权限与菜单。
- `GET /auth/oauth/login`: OAuth 三方登录入口。

### 媒体源 (Sources)
- `GET /sources`: 列出媒体源（集成 MediaMTX 状态）。
- `POST /sources`: 创建媒体源（支持 pull/push 类型）。
- `GET /sources/:id/preview`: 获取多协议预览地址（HLS/RTSP/RTMP/WebRTC）。

### 媒体资产 (Assets)
- `GET /assets`: 搜索与列出资产（支持 type/source_type/tags/visibility 过滤）。
- `POST /assets`: 接入新资产（支持文件上传与 URL 接入）。
- `GET /assets/:id/children`: 获取派生资产列表。

### 算子管理 (Operators)
- `POST /operators`: 创建算子基础信息。
- `POST /operators/:id/publish`: 发布算子（触发 Schema 与依赖门禁校验）。
- `POST /operators/:id/test`: 连通性测试（真实执行器试运行）。

#### 版本管理
- `GET /operators/:id/versions`: 列出所有版本。
- `POST /operators/:id/versions`: 创建新版本（定义 ExecMode 与 ExecConfig）。
- `POST /operators/:id/versions/activate`: 激活指定版本为生产版本。

#### MCP 生态集成
- `GET /operators/mcp/servers`: 列出已连接的 MCP 服务。
- `GET /operators/mcp/servers/:id/tools`: 浏览工具列表。
- `POST /operators/mcp/install`: 将 MCP Tool 直接安装为系统算子。
- `POST /operators/mcp/sync-templates`: 从 MCP 同步市场模板。

### 工作流与任务 (Workflows & Tasks)
- `POST /workflows`: 创建 DAG 工作流（触发连接兼容性校验）。
- `POST /workflows/:id/trigger`: 手动触发工作流执行。
- `GET /tasks`: 任务列表与统计。
- `GET /tasks/:id`: 任务详情（含 **NodeExecutions** 节点状态追踪）。
- `GET /tasks/:id/progress/stream`: **SSE** 实时进度推送。

### 系统配置 (System Config)
- `GET /system/configs`: 按分类获取系统配置。
- `PUT /system/configs`: 批量更新系统参数。

---

最后更新：2026-02-08
