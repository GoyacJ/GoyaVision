# GoyaVision 多租户与数据隔离重构设计方案

## 1. 总体目标
实现系统的多租户架构改造，确保数据安全隔离，并为核心业务资产（媒体、算子、工作流等）提供细粒度的权限控制（私有/角色可见/公开）。

## 2. 核心设计决策

### 2.1 多租户模式 (Multi-tenancy Model)
采用 **共享数据库、共享 Schema、字段隔离** 的模式。
- **低成本**：无需为每个租户创建独立数据库或 Schema，便于运维和资源共享。
- **实现**：所有业务表增加 `tenant_id` 字段。
- **隔离机制**：通过 GORM Scope 和 Middleware 强制注入 `tenant_id` 过滤条件，防止跨租户访问。

### 2.2 租户识别 (Tenant Resolution)
采用 **JWT Token 携带租户信息** 的方式。
- 用户登录时，Token Payload 中包含 `tenant_id`。
- 后端 Middleware 从 Token 解析 `tenant_id` 并存入 `context`。
- Repository 层从 `context` 获取 `tenant_id` 进行数据操作。

### 2.3 资源可见性模型 (Resource Visibility)
为所有核心资源（MediaAsset, MediaSource, Operator, AIModel, Workflow）引入统一的可见性控制：
- **范围定义**：
  - `Private` (0): 仅 **拥有者(Owner)** 及 **租户管理员** 可见。
  - `Role` (1): 拥有者 + 指定角色 (如 "Dev", "Ops") 的用户可见。
  - `Public` (2): 该 **租户下** 所有用户可见。
- **所有权**：资源创建者自动成为 **Owner**。

---

## 3. 数据库模型变更 (Schema Changes)

### 3.1 新增实体
新增 `tenants` 表：
```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(64) UNIQUE NOT NULL, -- 租户唯一标识
    status INT DEFAULT 1,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### 3.2 现有表变更
对以下表增加字段：
1.  **Users 表**：
    *   `tenant_id UUID` (索引)
2.  **资源表** (`media_assets`, `media_sources`, `operators`, `ai_models`, `workflows`)：
    *   `tenant_id UUID NOT NULL` (索引)
    *   `owner_id UUID` (索引)
    *   `visibility INT DEFAULT 0`
    *   `visible_role_ids JSONB` (存储可见的角色ID列表)
3.  **Tasks 表**：
    *   `tenant_id UUID NOT NULL`
    *   `triggered_by_user_id UUID` (用于“我的任务”筛选)

---

## 4. 业务逻辑与架构改造

### 4.1 上下文传播 (Context Propagation)
*   **Middleware**: 创建 `TenantMiddleware`，解析 JWT 中的 `tenant_id`。
*   **Context Key**: 定义专用 Key (e.g., `ctxKeyTenantID`, `ctxKeyUserRoles`)。

### 4.2 持久层改造 (Persistence Layer)
*   **Query Scopes**: 封装 GORM Scopes。
    *   `ScopeTenant(ctx)`: `WHERE tenant_id = ?`
    *   `ScopeVisibility(ctx, user_id, role_ids)`:
      ```sql
      AND (
          owner_id = @current_user_id
          OR visibility = 2 -- Public
          OR (visibility = 1 AND EXISTS (SELECT 1 FROM jsonb_array_elements_text(visible_role_ids) as r WHERE r IN (@current_user_roles)))
      )
      ```

### 4.3 业务服务层 (Application Layer)
*   **创建逻辑**：注入 `tenant_id` 和 `owner_id`。
*   **查询逻辑**：应用 `ScopeTenant` 和 `ScopeVisibility`。

---

## 5. 接口与前端适配

### 5.1 API DTO 变更
资源创建/详情 DTO 增加：
```go
type VisibilityConfig struct {
    Visibility     int      `json:"visibility"`
    VisibleRoleIDs []string `json:"visible_role_ids"`
}
```

### 5.2 前端改造
*   **表单**：增加“可见范围”设置组件。
*   **列表**：增加可见性状态展示。

---

## 6. 迁移策略 (Migration)
1.  **数据初始化**：
    *   创建默认租户。
    *   将现有数据归属到默认租户。
    *   初始化 `owner_id` 和 `visibility`。
