---
name: api-doc
description: 更新或新增 API 文档时使用。提供端点文档模板、请求/响应示例与错误码说明要求。
---

# API 文档更新指南

用于新增或修改 API 端点时，确保 `docs/api.md` 与实现一致。

## 何时使用

✅ **推荐场景**：
- 新增 API 端点
- 修改请求/响应结构或参数
- 变更错误码/响应封装

❌ **不适用场景**：
- 只需要查看现有 API（使用 Read 工具）
- 前端组件开发（使用 frontend-components skill）

## 操作步骤

1. **确认路由前缀**：`/api/v1` 与 RESTful 设计
2. **在 `docs/api.md` 中补充**：
   - 端点说明与用途
   - 请求参数（Path/Query/Body）
   - 请求示例
   - 响应结构与示例
   - 错误码与说明
3. **如涉及 DTO/响应信封变更**，更新 `internal/api/response` 说明

## 文档模板

```markdown
### 任务（Tasks）

#### 创建任务
POST /api/v1/tasks

**描述**: 创建新的工作流执行任务

**认证**: 需要 JWT Token

**权限**: `task:create`

**请求体**:
```json
{
  "workflow_id": "uuid",
  "input_assets": ["uuid"],
  "trigger_type": "manual"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "uuid",
    "status": "pending",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

**错误响应**:
| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 工作流不存在 |
| 500 | 服务器错误 |
```

## 通用约定

- **路径前缀**：所有 API 使用 `/api/v1` 前缀
- **认证**：需要 JWT Token（Bearer Token）
- **响应格式**：统一使用 `{code, message, data}` 格式
- **错误码**：200=成功，400=参数错误，401=未认证，403=无权限，404=不存在，500=服务器错误

## 更新检查清单

- [ ] 端点路径和方法已记录
- [ ] 请求参数（Path/Query/Body）已说明
- [ ] 请求示例已提供
- [ ] 响应结构和示例已提供
- [ ] 错误码和说明已列出
- [ ] 认证和权限要求已说明
- [ ] 文档已同步到 `CHANGELOG.md`
