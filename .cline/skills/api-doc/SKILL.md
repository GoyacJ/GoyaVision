---
name: api-doc
description: 更新或新增 API 文档时使用。提供端点文档模板、请求/响应示例与错误码说明要求。
---

# API 文档更新指南

用于新增或修改 API 端点时，确保 `docs/api.md` 与实现一致。

## 何时使用
- 新增 API 端点
- 修改请求/响应结构或参数
- 变更错误码/响应封装

## 操作步骤
1. 确认路由前缀 `/api/v1` 与 RESTful 设计。
2. 在 `docs/api.md` 中补充：
   - 端点说明与用途
   - 请求参数（Path/Query/Body）
   - 请求示例
   - 响应结构与示例
   - 错误码与说明
3. 如涉及 DTO/响应信封变更，更新 `internal/api/response` 说明。

## 文档模板（示例）
```markdown
### 任务（Tasks）

#### 创建任务
POST /api/v1/tasks

**请求体**
```json
{
  "workflow_id": "uuid",
  "input_assets": ["uuid"],
  "trigger_type": "manual"
}
```

**响应**
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": "uuid",
    "status": "pending"
  }
}
```
```
