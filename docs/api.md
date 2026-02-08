# GoyaVision API 文档（V1.0）

## 概述

- 基础前缀：`/api/v1`
- 认证方式：JWT Bearer Token（`Authorization: Bearer <token>`）
- 公开与受保护接口混合存在：部分查询接口可匿名访问，写操作需登录

## 通用约定

- 主要资源支持 `limit/offset` 分页参数。
- 多租户核心字段：`tenant_id`, `owner_id`, `visibility`。
- 运行时相关接口统一使用 `task_id` / `session_id` 追踪。

## 认证与用户

- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/profile`
- `PUT /api/v1/auth/password`

## 系统配置

- `GET /api/v1/public/config`
- `PUT /api/v1/system/config`

## 资产与媒体源

### Media Sources
- `GET /api/v1/sources`
- `GET /api/v1/sources/:id`
- `POST /api/v1/sources`
- `PUT /api/v1/sources/:id`
- `DELETE /api/v1/sources/:id`

### Media Assets
- `GET /api/v1/assets`
- `GET /api/v1/assets/:id`
- `GET /api/v1/assets/:id/children`
- `GET /api/v1/assets/tags`
- `POST /api/v1/assets`
- `PUT /api/v1/assets/:id`
- `DELETE /api/v1/assets/:id`

## 算子中心

### Operators
- `GET /api/v1/operators`
- `GET /api/v1/operators/:id`
- `POST /api/v1/operators`
- `PUT /api/v1/operators/:id`
- `DELETE /api/v1/operators/:id`
- `POST /api/v1/operators/:id/publish`
- `POST /api/v1/operators/:id/deprecate`
- `POST /api/v1/operators/:id/test`
- `POST /api/v1/operators/validate-schema`
- `POST /api/v1/operators/validate-connection`

### Operator Versions / Templates / Dependencies
- `GET /api/v1/operators/:id/versions`
- `GET /api/v1/operators/:id/versions/:version_id`
- `POST /api/v1/operators/:id/versions`
- `POST /api/v1/operators/:id/versions/activate`
- `POST /api/v1/operators/:id/versions/rollback`
- `POST /api/v1/operators/:id/versions/archive`
- `GET /api/v1/operators/templates`
- `GET /api/v1/operators/templates/:template_id`
- `POST /api/v1/operators/templates/install`
- `GET /api/v1/operators/:id/dependencies`
- `GET /api/v1/operators/:id/dependencies/check`
- `PUT /api/v1/operators/:id/dependencies`

### MCP
- `GET /api/v1/operators/mcp/servers`
- `GET /api/v1/operators/mcp/servers/:id/tools`
- `GET /api/v1/operators/mcp/servers/:id/tools/:tool/preview`
- `POST /api/v1/operators/mcp/install`
- `POST /api/v1/operators/mcp/sync-templates`

## 算法库

- `GET /api/v1/algorithms`
- `GET /api/v1/algorithms/:id`
- `POST /api/v1/algorithms`
- `PUT /api/v1/algorithms/:id`
- `DELETE /api/v1/algorithms/:id`
- `POST /api/v1/algorithms/:id/versions`
- `POST /api/v1/algorithms/:id/versions/:version_id/publish`

## AI 模型

- `GET /api/v1/ai-models`
- `GET /api/v1/ai-models/:id`
- `POST /api/v1/ai-models`
- `PUT /api/v1/ai-models/:id`
- `DELETE /api/v1/ai-models/:id`
- `POST /api/v1/ai-models/:id/test-connection`

## 工作流与任务

### Workflows
- `GET /api/v1/workflows`
- `GET /api/v1/workflows/:id`
- `POST /api/v1/workflows`
- `PUT /api/v1/workflows/:id`
- `DELETE /api/v1/workflows/:id`
- `POST /api/v1/workflows/:id/enable`
- `POST /api/v1/workflows/:id/disable`
- `POST /api/v1/workflows/:id/trigger`

### Workflow Revisions
- `GET /api/v1/workflows/:id/revisions`
- `GET /api/v1/workflows/:id/revisions/:revision`
- `POST /api/v1/workflows/:id/revisions`

### Tasks
- `GET /api/v1/tasks`
- `GET /api/v1/tasks/stats`
- `GET /api/v1/tasks/:id`
- `POST /api/v1/tasks`
- `PUT /api/v1/tasks/:id`
- `DELETE /api/v1/tasks/:id`
- `POST /api/v1/tasks/:id/start`
- `POST /api/v1/tasks/:id/complete`
- `POST /api/v1/tasks/:id/fail`
- `POST /api/v1/tasks/:id/cancel`
- `GET /api/v1/tasks/:id/progress/stream`（SSE）

### Task Context / Events
- `GET /api/v1/tasks/:id/context`
- `GET /api/v1/tasks/:id/context/patches`
- `POST /api/v1/tasks/:id/context/snapshot`
- `GET /api/v1/tasks/:id/events`

## Agent Run Loop

- `GET /api/v1/agent/sessions`
- `GET /api/v1/agent/sessions/:id`
- `GET /api/v1/agent/sessions/:id/events`
- `POST /api/v1/agent/sessions`
- `POST /api/v1/agent/sessions/:id/run`
- `POST /api/v1/agent/sessions/:id/stop`

## 产物与文件

### Artifacts
- `GET /api/v1/artifacts`
- `GET /api/v1/artifacts/:id`
- `POST /api/v1/artifacts`
- `DELETE /api/v1/artifacts/:id`
- `GET /api/v1/tasks/:task_id/artifacts`

### Files
- `GET /api/v1/files`
- `POST /api/v1/files`
- `GET /api/v1/files/:id`
- `PUT /api/v1/files/:id`
- `DELETE /api/v1/files/:id`
- `GET /api/v1/files/:id/download`

## 管理接口

- Users: `GET/POST/PUT/DELETE /api/v1/users...`
- Roles: `GET/POST/PUT/DELETE /api/v1/roles...`
- Menus: `GET/POST/PUT/DELETE /api/v1/menus...`
- Tenants: `GET/POST/PUT/DELETE /api/v1/tenants...`

最后更新：2026-02-08
