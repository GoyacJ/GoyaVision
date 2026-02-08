# GoyaVision 整体重构方案

> 版本: 1.0
> 日期: 2026-02-04
> 基于: 全部 69 个 Go 源文件（~10,500 LOC）逐行分析

---

## 目录

- [A. 现状诊断](#a-现状诊断)
- [B. 目标架构蓝图](#b-目标架构蓝图)
- [C. 最终目录结构](#c-最终目录结构)
- [D. 核心模块边界](#d-核心模块边界)
- [E. 契约设计](#e-契约设计)
- [F. 数据与持久化设计](#f-数据与持久化设计)
- [G. 并发与可靠性](#g-并发与可靠性)
- [J. 代码骨架](#j-代码骨架)
- [K. 交付与切换方案](#k-交付与切换方案)
- [L. 工作拆分清单](#l-工作拆分清单)
- [附录: 关键假设列表](#附录-关键假设列表)

---

## A. 现状诊断

基于对全部 69 个 Go 文件的逐行分析，识别出以下结构性问题，按影响范围/修复收益/风险排序：

| # | 问题 | 影响范围 | 修复收益 | 风险 |
|---|------|---------|---------|------|
| 1 | **巨型 Repository 接口 (125 方法)**：`port/repository.go` 一个接口覆盖 User/Role/Permission/Menu/Asset/Source/Operator/Workflow/Task/Artifact/File 全部 11 个领域，严重违反 ISP | 全局 | 极高 | 中 |
| 2 | **Domain 实体被 GORM 标签污染**：`domain/*.go` 所有实体直接携带 `gorm:"type:uuid;primaryKey"` 等 ORM 标签，域层耦合基础设施 | 全局 | 高 | 中 |
| 3 | **无事务支持**：`WorkflowService.Create` 分步创建 workflow + nodes + edges 无原子保证；`UserService.Create` 创建用户 + 分配角色非原子；`MediaSourceService.Create` 跨 MediaMTX + DB 只有单向补偿 | 核心业务 | 极高 | 高 |
| 4 | **贫血域模型**：Domain 实体只有 `IsXxx()` 判断方法，无业务规则、无状态机、无不变量校验 | 核心业务 | 高 | 低 |
| 5 | **Deps 结构体引用具体类型**：`handler/deps.go` 直接依赖 `*mediamtx.Client`、`*storage.MinIOClient`，违反 DIP | API 层 | 高 | 低 |
| 6 | **Handler 每次注册时创建 Service 实例**：`RegisterUser(g, d)` 内部 `app.NewUserService(d.Repo)`，服务生命周期管理混乱 | API 层 | 中 | 低 |
| 7 | **错误处理不一致**：部分用 `errors.Is()`、部分用 `==`、部分用 `strings.Contains()`；App 层返回 `errors.New("string")`，丢失错误类型和上下文 | 全局 | 高 | 低 |
| 8 | **无结构化日志**：全局使用 `log.Printf`，无 request ID、无 trace ID、无级别控制 | 运维/调试 | 高 | 低 |
| 9 | **验证逻辑重复散落**：DTO 定义了 `validate:"required"` 但从未启用验证器；Limit 校验在 handler 和 service 重复 10+ 次 | 全局 | 中 | 低 |
| 10 | **Filter 结构体放在 Domain 层**：`domain.MediaAssetFilter`、`domain.TaskFilter` 等查询关注点泄漏到领域 | 架构边界 | 中 | 低 |
| 11 | **权限每次请求 N+1 查询**：`LoadUserPermissions` 中间件每个请求都执行 `GetUserWithRoles` + `GetPermissionsByRoleIDs` | 性能 | 中 | 低 |
| 12 | **异步任务 Context 误用**：`WorkflowScheduler.runWorkflow` 使用 `context.Background()` 启动 goroutine，丢失取消信号和 trace 传播 | 可靠性 | 高 | 中 |
| 13 | **单文件 Repository 实现 1277 行**：`adapter/persistence/repository.go` 所有表的 CRUD 塞在一个文件 | 可维护性 | 中 | 低 |
| 14 | **无统一 API 响应信封**：列表接口返回 `{items, total}`，错误返回 `{error, message}`，缺少 `code`/`request_id`/`trace_id` | API 质量 | 中 | 低 |
| 15 | **无测试基础设施**：仅发现 `media_source_test.go` 一个测试文件，无 mock、无 fixture、无集成测试 | 质量保障 | 极高 | 低 |
| 16 | ~~**工作流引擎顺序执行**~~：`SimpleWorkflowEngine` 已删除，由 `DAGWorkflowEngine` 完全替代（支持拓扑排序、并行执行、条件分支） | 功能完整性 | ✅ 已解决 | — |

---

## B. 目标架构蓝图

### 推荐架构：Clean Architecture + DDD-lite

**选型理由：**

1. **现有代码已有 Clean Architecture 雏形**（domain/port/app/adapter/api 分层），重构是矫正而非推倒，成本最低。
2. **项目规模 (~10K LOC) 不适合完整 DDD**（Aggregate Root、Domain Event Sourcing 过重），DDD-lite 取其精华（聚合边界、值对象、领域服务、领域事件）即可。
3. **媒体处理是 IO 密集型**，Hexagonal 在这个场景下与 Clean Architecture 差异不大，后者生态和 Go 社区实践更成熟。

### 核心设计原则

```
依赖规则（严格单向）：

  Infrastructure --> Application --> Domain
       ^                              ^
    Interface                      (无依赖)
    Adapters
```

**关键差异（对比现状）：**

- Domain 实体**不携带任何 ORM 标签**
- Repository 接口**按聚合拆分**（不再是一个巨型接口）
- 引入 **UnitOfWork** 解决事务问题
- 引入 **Domain Event** 替代直接 goroutine 调度
- 引入 **Application Error** 类型体系替代 `errors.New("string")`

---

## C. 最终目录结构

```
goyavision/
├── cmd/
│   ├── server/
│   │   └── main.go                    # 应用启动、DI 组装
│   └── migrate/
│       └── main.go                    # 数据库迁移工具
│
├── internal/
│   ├── domain/                        # 【域层】纯业务逻辑，零外部依赖
│   │   ├── media/                     # 媒体聚合
│   │   │   ├── source.go              # MediaSource 实体 + 值对象 + 业务方法
│   │   │   ├── asset.go               # MediaAsset 实体
│   │   │   └── repository.go          # SourceRepository, AssetRepository 接口
│   │   ├── operator/                  # 算子聚合
│   │   │   ├── operator.go            # Operator 实体 + 标准协议
│   │   │   └── repository.go          # OperatorRepository 接口
│   │   ├── workflow/                  # 工作流聚合
│   │   │   ├── workflow.go            # Workflow 聚合根 (含 Node/Edge)
│   │   │   ├── task.go                # Task 实体
│   │   │   ├── artifact.go            # Artifact 实体
│   │   │   ├── engine.go              # WorkflowEngine, OperatorExecutor 接口
│   │   │   └── repository.go          # WorkflowRepository, TaskRepository 接口
│   │   ├── identity/                  # 身份认证聚合
│   │   │   ├── user.go                # User 实体
│   │   │   ├── role.go                # Role + Permission + Menu
│   │   │   └── repository.go          # UserRepository, RoleRepository 接口
│   │   ├── storage/                   # 文件存储聚合
│   │   │   ├── file.go                # File 实体
│   │   │   └── repository.go          # FileRepository 接口
│   │   └── event.go                   # 领域事件定义
│   │
│   ├── app/                           # 【应用层】用例编排，不含业务规则
│   │   ├── command/                   # 写操作（Command）
│   │   │   ├── create_source.go
│   │   │   ├── create_asset.go
│   │   │   ├── create_workflow.go
│   │   │   ├── trigger_workflow.go
│   │   │   ├── create_operator.go
│   │   │   ├── upload_file.go
│   │   │   ├── create_user.go
│   │   │   └── auth_login.go
│   │   ├── query/                     # 读操作（Query）
│   │   │   ├── list_sources.go
│   │   │   ├── list_assets.go
│   │   │   ├── list_workflows.go
│   │   │   ├── list_tasks.go
│   │   │   ├── get_profile.go
│   │   │   └── list_operators.go
│   │   ├── dto/                       # 应用层输入/输出 DTO
│   │   │   ├── command.go             # Command 输入类型
│   │   │   ├── query.go               # Query 输入类型 + Filter
│   │   │   └── result.go              # 用例返回类型
│   │   ├── port/                      # 应用层出站端口（基础设施抽象）
│   │   │   ├── unit_of_work.go        # UnitOfWork 接口
│   │   │   ├── media_gateway.go       # MediaMTX 网关接口
│   │   │   ├── object_storage.go      # 对象存储接口
│   │   │   ├── token_service.go       # JWT 令牌接口
│   │   │   └── event_bus.go           # 事件总线接口
│   │   └── scheduler.go              # 工作流调度器
│   │
│   ├── infra/                         # 【基础设施层】外部系统对接
│   │   ├── persistence/               # 数据持久化
│   │   │   ├── model/                 # GORM 模型（携带 ORM 标签）
│   │   │   │   ├── media_source.go
│   │   │   │   ├── media_asset.go
│   │   │   │   ├── operator.go
│   │   │   │   ├── workflow.go
│   │   │   │   ├── task.go
│   │   │   │   ├── artifact.go
│   │   │   │   ├── user.go
│   │   │   │   └── file.go
│   │   │   ├── mapper/                # Model <-> Domain 映射
│   │   │   │   ├── media.go
│   │   │   │   ├── workflow.go
│   │   │   │   ├── identity.go
│   │   │   │   └── file.go
│   │   │   ├── repo/                  # Repository 实现（按聚合分文件）
│   │   │   │   ├── media_source.go
│   │   │   │   ├── media_asset.go
│   │   │   │   ├── operator.go
│   │   │   │   ├── workflow.go
│   │   │   │   ├── task.go
│   │   │   │   ├── identity.go
│   │   │   │   └── file.go
│   │   │   ├── uow.go                # UnitOfWork GORM 实现
│   │   │   ├── migrate.go            # 数据库迁移
│   │   │   └── seed.go               # 初始数据
│   │   ├── mediamtx/                  # MediaMTX HTTP 客户端
│   │   │   ├── client.go
│   │   │   ├── models.go
│   │   │   └── gateway.go            # 实现 app/port/media_gateway.go
│   │   ├── engine/                    # 工作流引擎
│   │   │   ├── dag_engine.go          # DAG 执行引擎（含拓扑排序）
│   │   │   └── http_executor.go       # HTTP 算子执行器
│   │   ├── minio/                     # MinIO 对象存储
│   │   │   └── client.go             # 实现 app/port/object_storage.go
│   │   ├── auth/                      # JWT 实现
│   │   │   └── jwt.go                # 实现 app/port/token_service.go
│   │   └── eventbus/                  # 事件总线
│   │       └── local.go              # 本地内存事件总线
│   │
│   └── api/                           # 【接口适配器层】HTTP 入站
│       ├── handler/                   # HTTP 处理器
│       │   ├── auth.go
│       │   ├── media_source.go
│       │   ├── media_asset.go
│       │   ├── operator.go
│       │   ├── workflow.go
│       │   ├── task.go
│       │   ├── artifact.go
│       │   ├── file.go
│       │   └── admin.go               # User/Role/Menu/Permission
│       ├── request/                   # HTTP 请求 DTO
│       │   ├── auth.go
│       │   ├── media.go
│       │   ├── workflow.go
│       │   └── admin.go
│       ├── response/                  # HTTP 响应 DTO
│       │   ├── envelope.go            # 统一响应信封
│       │   ├── media.go
│       │   ├── workflow.go
│       │   └── admin.go
│       ├── middleware/
│       │   ├── auth.go                # JWT 认证
│       │   ├── permission.go          # RBAC 权限
│       │   ├── requestid.go           # Request ID
│       │   ├── logging.go             # 结构化日志
│       │   └── validator.go           # 请求验证
│       ├── router.go                  # 路由注册
│       └── errors.go                  # HTTP 错误映射
│
├── pkg/                               # 共享库（可被外部项目引用）
│   ├── apperr/                        # 应用错误类型体系
│   │   └── errors.go
│   ├── pagination/                    # 分页工具
│   │   └── pagination.go
│   └── logger/                        # 结构化日志
│       └── logger.go
│
├── config/
│   └── config.go                      # 配置加载
├── configs/
│   ├── config.dev.yaml
│   ├── config.prod.yaml
│   ├── config.example.yaml
│   ├── .env.example
│   └── mediamtx.yml
├── web/                               # Vue 3 前端
├── embed.go                           # 前端嵌入
├── Makefile
├── docker-compose.yml
├── go.mod
└── go.sum
```

### 各层职责与依赖规则

| 层 | 包路径 | 职责 | 允许依赖 | 禁止依赖 |
|----|--------|------|---------|---------|
| **Domain** | `internal/domain/*` | 实体、值对象、聚合、领域服务、Repository 接口 | `pkg/apperr` | App, Infra, API, 任何外部库 |
| **Application** | `internal/app/*` | 用例编排、事务协调、DTO 转换 | Domain, `pkg/*` | Infra, API |
| **Infrastructure** | `internal/infra/*` | GORM, HTTP 客户端, MinIO, JWT | Domain, App(port), `pkg/*`, 外部库 | API |
| **API** | `internal/api/*` | HTTP 路由、请求/响应 DTO、中间件 | App, Domain, `pkg/*` | Infra |
| **pkg** | `pkg/*` | 通用工具，无业务语义 | 标准库 | internal/* |

---

## D. 核心模块边界

### 聚合清单

### 1. Media 聚合（media/）

```
输入:
  - CreateSourceCommand{Name, Type, URL, Protocol}
  - CreateAssetCommand{Type, SourceType, Name, Path, Tags, Metadata}
  - ListAssetsQuery{Filter, Pagination}

输出:
  - SourceResult{ID, Name, PathName, ...}
  - AssetResult{ID, Name, Type, Tags, ...}
  - PagedResult[AssetResult]{Items, Total, Limit, Offset}

领域事件:
  - SourceCreated{SourceID, PathName}
  - SourceDeleted{SourceID, PathName}
  - AssetCreated{AssetID, Type, SourceID}
```

### 2. Workflow 聚合（workflow/）

```
输入:
  - CreateWorkflowCommand{Code, Name, TriggerType, Nodes[], Edges[]}
  - TriggerWorkflowCommand{WorkflowID, AssetID}
  - CancelTaskCommand{TaskID}

输出:
  - WorkflowResult{ID, Code, Name, Nodes[], Edges[]}
  - TaskResult{ID, Status, Progress, Artifacts[]}

领域事件:
  - WorkflowCreated{WorkflowID}
  - TaskStarted{TaskID, WorkflowID}
  - TaskCompleted{TaskID, Status}
  - TaskFailed{TaskID, Error}
```

### 3. Operator 聚合（operator/）

```
输入:
  - CreateOperatorCommand{Code, Name, Category, Endpoint, ...}
  - OperatorInput{AssetID, Params}  // 标准协议

输出:
  - OperatorResult{ID, Code, Status, ...}
  - OperatorOutput{OutputAssets[], Results[], Timeline[], Diagnostics}

领域事件:
  - OperatorRegistered{OperatorID, Code}
```

### 4. Identity 聚合（identity/）

```
输入:
  - LoginCommand{Username, Password}
  - CreateUserCommand{Username, Password, Nickname, RoleIDs[]}
  - AssignRolesCommand{UserID, RoleIDs[]}

输出:
  - TokenPair{AccessToken, RefreshToken, ExpiresIn}
  - UserProfile{ID, Username, Roles[], Permissions[], Menus[]}
```

### 5. Storage 聚合（storage/）

```
输入:
  - UploadFileCommand{File, OriginalName, UploaderID}

输出:
  - FileResult{ID, Name, Path, Size, MimeType, URL}
```

---

## E. 契约设计

### 1. 统一 API 响应信封

**成功响应：**

```json
{
  "code": 0,
  "message": "ok",
  "data": { "..." },
  "request_id": "req-a1b2c3d4",
  "timestamp": "2026-02-04T12:00:00Z"
}
```

**列表响应：**

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [ "..." ],
    "total": 150,
    "limit": 20,
    "offset": 0
  },
  "request_id": "req-a1b2c3d4",
  "timestamp": "2026-02-04T12:00:00Z"
}
```

**错误响应：**

```json
{
  "code": 40401,
  "message": "media source not found",
  "details": [
    {"field": "id", "reason": "no record with id=xxx"}
  ],
  "request_id": "req-a1b2c3d4",
  "timestamp": "2026-02-04T12:00:00Z"
}
```

### 2. 错误码规范

格式: `HHHDD`（HTTP 状态码 3 位 + 业务子码 2 位）

| 错误码 | 含义 |
|--------|------|
| 40001 | 请求参数校验失败 |
| 40002 | 请求体解析失败 |
| 40101 | Token 缺失 |
| 40102 | Token 过期 |
| 40103 | Token 无效 |
| 40301 | 权限不足 |
| 40401 | 媒体源不存在 |
| 40402 | 媒体资产不存在 |
| 40403 | 算子不存在 |
| 40404 | 工作流不存在 |
| 40405 | 任务不存在 |
| 40901 | 资源已存在（唯一约束冲突） |
| 40902 | 媒体源有关联资产，不可删除 |
| 42201 | 工作流无节点，不可启用 |
| 42202 | 工作流存在环路 |
| 50001 | 数据库错误 |
| 50301 | MediaMTX 不可用 |
| 50302 | MinIO 不可用 |

### 3. 分页规范

```
请求: GET /api/v1/assets?limit=20&offset=0&type=video&tags=ai,detection

响应 data 字段:
{
  "items": [...],
  "total": 150,       // 符合过滤条件的总数
  "limit": 20,        // 本次请求的 limit（已校正: 1 <= limit <= 100）
  "offset": 0         // 本次请求的 offset
}

默认值: limit=20, offset=0
上限: limit <= 100（管理端 <= 1000）
```

### 4. 日志/Trace 字段规范

```json
{
  "level": "info",
  "ts": "2026-02-04T12:00:00.000Z",
  "caller": "handler/asset.go:45",
  "msg": "asset created",
  "request_id": "req-a1b2c3d4",
  "trace_id": "abc123def456",
  "user_id": "uuid-...",
  "method": "POST",
  "path": "/api/v1/assets",
  "status": 201,
  "latency_ms": 42,
  "asset_id": "uuid-...",
  "error": null
}
```

### 5. 幂等设计

写操作通过 `Idempotency-Key` 请求头实现幂等：

- 客户端在 POST 请求中携带 `Idempotency-Key: <uuid>`
- 服务端在 Redis/内存 map 中缓存 key -> response，TTL 24h
- 重复请求直接返回缓存结果
- 适用于：创建媒体源、触发工作流、上传文件

---

## F. 数据与持久化设计

### 1. Domain Entity vs GORM Model 分离

**纯领域实体（无 ORM 标签）：**

```go
// internal/domain/media/source.go
type Source struct {
    ID            uuid.UUID
    Name          string
    PathName      string
    Type          SourceType
    URL           string
    Protocol      Protocol
    Enabled       bool
    RecordEnabled bool
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

// 领域业务方法
func (s *Source) Enable()  { s.Enabled = true }
func (s *Source) Disable() { s.Enabled = false }
func (s *Source) Validate() error {
    if s.Name == "" { return apperr.InvalidInput("name is required") }
    if s.Type == SourceTypePull && s.URL == "" {
        return apperr.InvalidInput("url is required for pull source")
    }
    return nil
}
```

**GORM 模型（基础设施层）：**

```go
// internal/infra/persistence/model/media_source.go
type MediaSourceModel struct {
    ID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
    Name          string     `gorm:"type:varchar(255);not null"`
    PathName      string     `gorm:"type:varchar(255);not null;uniqueIndex"`
    Type          string     `gorm:"type:varchar(20);not null"`
    URL           string     `gorm:"type:varchar(1024)"`
    Protocol      string     `gorm:"type:varchar(20)"`
    Enabled       bool       `gorm:"not null;default:true"`
    RecordEnabled bool       `gorm:"not null;default:false"`
    CreatedAt     time.Time  `gorm:"autoCreateTime"`
    UpdatedAt     time.Time  `gorm:"autoUpdateTime"`
}
func (MediaSourceModel) TableName() string { return "media_sources" }
```

**映射函数：**

```go
// internal/infra/persistence/mapper/media.go
func SourceToModel(s *media.Source) *model.MediaSourceModel { ... }
func SourceFromModel(m *model.MediaSourceModel) *media.Source { ... }
```

### 2. Repository 按聚合拆分

```go
// internal/domain/media/repository.go
type SourceRepository interface {
    Create(ctx context.Context, s *Source) error
    Get(ctx context.Context, id uuid.UUID) (*Source, error)
    GetByPathName(ctx context.Context, pathName string) (*Source, error)
    List(ctx context.Context, filter SourceFilter) ([]*Source, int64, error)
    Update(ctx context.Context, s *Source) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type AssetRepository interface {
    Create(ctx context.Context, a *Asset) error
    Get(ctx context.Context, id uuid.UUID) (*Asset, error)
    List(ctx context.Context, filter AssetFilter) ([]*Asset, int64, error)
    Update(ctx context.Context, a *Asset) error
    Delete(ctx context.Context, id uuid.UUID) error
    ListBySource(ctx context.Context, sourceID uuid.UUID) ([]*Asset, error)
    ListByParent(ctx context.Context, parentID uuid.UUID) ([]*Asset, error)
    GetAllTags(ctx context.Context) ([]string, error)
}
```

### 3. UnitOfWork 事务边界

```go
// internal/app/port/unit_of_work.go
type UnitOfWork interface {
    // Do 在事务内执行 fn；fn 返回 error 则回滚，否则提交
    Do(ctx context.Context, fn func(tx UnitOfWork) error) error

    // 每个聚合的 Repository 访问器
    Sources() media.SourceRepository
    Assets() media.AssetRepository
    Operators() operator.OperatorRepository
    Workflows() workflow.WorkflowRepository
    Tasks() workflow.TaskRepository
    Artifacts() workflow.ArtifactRepository
    Users() identity.UserRepository
    Roles() identity.RoleRepository
    Files() storage.FileRepository
}
```

**事务使用示例：**

```go
// internal/app/command/create_workflow.go
func (h *CreateWorkflowHandler) Handle(ctx context.Context, cmd CreateWorkflowCommand) error {
    return h.uow.Do(ctx, func(tx port.UnitOfWork) error {
        wf := workflow.NewWorkflow(cmd.Code, cmd.Name, cmd.TriggerType)
        if err := tx.Workflows().Create(ctx, wf); err != nil {
            return err
        }
        for _, n := range cmd.Nodes {
            node := wf.AddNode(n.NodeKey, n.NodeType, n.OperatorID)
            if err := tx.Workflows().CreateNode(ctx, node); err != nil {
                return err  // 自动回滚
            }
        }
        return nil  // 自动提交
    })
}
```

### 4. Filter/Query 结构体归属

Filter 从 domain 移到 `app/dto/query.go`：

```go
// internal/app/dto/query.go
type AssetFilter struct {
    Type       *string
    SourceType *string
    SourceID   *uuid.UUID
    Status     *string
    Tags       []string
    From       *time.Time
    To         *time.Time
}

type Pagination struct {
    Limit  int
    Offset int
}

func (p *Pagination) Normalize() {
    if p.Limit <= 0 { p.Limit = 20 }
    if p.Limit > 100 { p.Limit = 100 }
    if p.Offset < 0 { p.Offset = 0 }
}
```

---

## G. 并发与可靠性

### 1. 超时控制

```go
// 算子执行超时：通过 context 传递
func (e *HTTPExecutor) Execute(ctx context.Context, op *operator.Operator, input *OperatorInput) (*OperatorOutput, error) {
    timeout := time.Duration(op.TimeoutSeconds) * time.Second
    if timeout == 0 {
        timeout = 5 * time.Minute // 默认
    }
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()
    // ... HTTP 调用
}
```

### 2. 重试（算子执行）

```go
// internal/infra/engine/http_executor.go
func (e *HTTPExecutor) ExecuteWithRetry(ctx context.Context, op *Operator, input *OperatorInput, maxRetries int) (*OperatorOutput, error) {
    var lastErr error
    for attempt := 0; attempt <= maxRetries; attempt++ {
        if attempt > 0 {
            backoff := time.Duration(attempt*attempt) * time.Second
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            case <-time.After(backoff):
            }
        }
        output, err := e.Execute(ctx, op, input)
        if err == nil {
            return output, nil
        }
        lastErr = err
    }
    return nil, fmt.Errorf("after %d retries: %w", maxRetries, lastErr)
}
```

### 3. 任务失败补偿

```go
// 工作流调度器中的异步任务，使用带 recover 的 goroutine
func (s *Scheduler) runWorkflowAsync(ctx context.Context, wfID uuid.UUID, task *workflow.Task) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                s.logger.Error("workflow panic",
                    "workflow_id", wfID,
                    "task_id", task.ID,
                    "panic", r,
                )
                s.failTask(context.Background(), task, fmt.Errorf("panic: %v", r))
            }
        }()

        if err := s.engine.Execute(ctx, wfID, task); err != nil {
            s.failTask(context.Background(), task, err)
        }
    }()
}

func (s *Scheduler) failTask(ctx context.Context, task *workflow.Task, err error) {
    task.Fail(err.Error())
    _ = s.uow.Tasks().Update(ctx, task)
    s.eventBus.Publish(ctx, workflow.TaskFailed{TaskID: task.ID, Error: err.Error()})
}
```

### 4. 限流（API 层）

```go
// internal/api/middleware/ratelimit.go
// 使用 Echo 内置 + golang.org/x/time/rate
func RateLimit(rps int) echo.MiddlewareFunc {
    limiter := rate.NewLimiter(rate.Limit(rps), rps*2)
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if !limiter.Allow() {
                return c.JSON(http.StatusTooManyRequests, response.Error(42901, "rate limit exceeded"))
            }
            return next(c)
        }
    }
}
```

---

## J. 代码骨架

### 骨架 1：应用错误类型体系

```go
// pkg/apperr/errors.go
package apperr

import "fmt"

type Code int

const (
    CodeInvalidInput    Code = 40001
    CodeUnauthorized    Code = 40101
    CodeForbidden       Code = 40301
    CodeNotFound        Code = 40401
    CodeConflict        Code = 40901
    CodeUnprocessable   Code = 42201
    CodeInternal        Code = 50001
    CodeServiceUnavail  Code = 50301
)

type Error struct {
    code    Code
    message string
    cause   error
    details []Detail
}

type Detail struct {
    Field  string `json:"field"`
    Reason string `json:"reason"`
}

func (e *Error) Error() string {
    if e.cause != nil {
        return fmt.Sprintf("[%d] %s: %v", e.code, e.message, e.cause)
    }
    return fmt.Sprintf("[%d] %s", e.code, e.message)
}

func (e *Error) Code() Code      { return e.code }
func (e *Error) Unwrap() error   { return e.cause }
func (e *Error) Details() []Detail { return e.details }

func NotFound(msg string) *Error {
    return &Error{code: CodeNotFound, message: msg}
}

func InvalidInput(msg string) *Error {
    return &Error{code: CodeInvalidInput, message: msg}
}

func Conflict(msg string) *Error {
    return &Error{code: CodeConflict, message: msg}
}

func Internal(msg string, cause error) *Error {
    return &Error{code: CodeInternal, message: msg, cause: cause}
}

func Wrap(code Code, msg string, cause error) *Error {
    return &Error{code: code, message: msg, cause: cause}
}

func WithDetails(err *Error, details ...Detail) *Error {
    err.details = append(err.details, details...)
    return err
}
```

### 骨架 2：统一响应信封

```go
// internal/api/response/envelope.go
package response

import (
    "time"

    "github.com/labstack/echo/v4"
)

type Envelope struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    Details   interface{} `json:"details,omitempty"`
    RequestID string      `json:"request_id"`
    Timestamp time.Time   `json:"timestamp"`
}

type PagedData struct {
    Items  interface{} `json:"items"`
    Total  int64       `json:"total"`
    Limit  int         `json:"limit"`
    Offset int         `json:"offset"`
}

func OK(c echo.Context, data interface{}) error {
    return c.JSON(200, Envelope{
        Code:      0,
        Message:   "ok",
        Data:      data,
        RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
        Timestamp: time.Now(),
    })
}

func Created(c echo.Context, data interface{}) error {
    return c.JSON(201, Envelope{
        Code:      0,
        Message:   "created",
        Data:      data,
        RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
        Timestamp: time.Now(),
    })
}

func Paged(c echo.Context, items interface{}, total int64, limit, offset int) error {
    return OK(c, PagedData{
        Items:  items,
        Total:  total,
        Limit:  limit,
        Offset: offset,
    })
}

func Error(code int, message string) Envelope {
    return Envelope{
        Code:      code,
        Message:   message,
        Timestamp: time.Now(),
    }
}
```

### 骨架 3：UnitOfWork GORM 实现

```go
// internal/infra/persistence/uow.go
package persistence

import (
    "context"

    "goyavision/internal/app/port"
    "goyavision/internal/domain/media"
    "goyavision/internal/domain/workflow"
    "goyavision/internal/domain/identity"
    "goyavision/internal/domain/storage"
    "goyavision/internal/domain/operator"
    "goyavision/internal/infra/persistence/repo"

    "gorm.io/gorm"
)

type gormUoW struct {
    db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) port.UnitOfWork {
    return &gormUoW{db: db}
}

func (u *gormUoW) Do(ctx context.Context, fn func(tx port.UnitOfWork) error) error {
    return u.db.WithContext(ctx).Transaction(func(gormTx *gorm.DB) error {
        txUoW := &gormUoW{db: gormTx}
        return fn(txUoW)
    })
}

func (u *gormUoW) Sources() media.SourceRepository {
    return repo.NewSourceRepo(u.db)
}

func (u *gormUoW) Assets() media.AssetRepository {
    return repo.NewAssetRepo(u.db)
}

func (u *gormUoW) Operators() operator.OperatorRepository {
    return repo.NewOperatorRepo(u.db)
}

func (u *gormUoW) Workflows() workflow.WorkflowRepository {
    return repo.NewWorkflowRepo(u.db)
}

func (u *gormUoW) Tasks() workflow.TaskRepository {
    return repo.NewTaskRepo(u.db)
}

func (u *gormUoW) Artifacts() workflow.ArtifactRepository {
    return repo.NewArtifactRepo(u.db)
}

func (u *gormUoW) Users() identity.UserRepository {
    return repo.NewUserRepo(u.db)
}

func (u *gormUoW) Roles() identity.RoleRepository {
    return repo.NewRoleRepo(u.db)
}

func (u *gormUoW) Files() storage.FileRepository {
    return repo.NewFileRepo(u.db)
}
```

### 骨架 4：Command Handler（用例）

```go
// internal/app/command/create_source.go
package command

import (
    "context"

    "goyavision/internal/app/port"
    "goyavision/internal/domain/media"
    "goyavision/pkg/apperr"
    "goyavision/pkg/logger"
)

type CreateSourceCommand struct {
    Name     string
    Type     string
    URL      string
    Protocol string
    Enabled  bool
}

type CreateSourceHandler struct {
    uow     port.UnitOfWork
    gateway port.MediaGateway
    log     *logger.Logger
}

func NewCreateSourceHandler(uow port.UnitOfWork, gw port.MediaGateway, log *logger.Logger) *CreateSourceHandler {
    return &CreateSourceHandler{uow: uow, gateway: gw, log: log}
}

func (h *CreateSourceHandler) Handle(ctx context.Context, cmd CreateSourceCommand) (*media.Source, error) {
    src, err := media.NewSource(cmd.Name, media.SourceType(cmd.Type), cmd.URL, media.Protocol(cmd.Protocol))
    if err != nil {
        return nil, err
    }
    src.Enabled = cmd.Enabled

    // 1. 先在 MediaMTX 创建路径
    mtxSource := cmd.URL
    if src.Type == media.SourceTypePush {
        mtxSource = "publisher"
    }
    if err := h.gateway.AddPath(ctx, src.PathName, mtxSource); err != nil {
        return nil, apperr.Wrap(apperr.CodeServiceUnavail, "mediamtx add path failed", err)
    }

    // 2. 写入数据库（失败则补偿删除 MediaMTX 路径）
    if err := h.uow.Sources().Create(ctx, src); err != nil {
        _ = h.gateway.DeletePath(ctx, src.PathName)
        return nil, apperr.Internal("create source in db", err)
    }

    h.log.Info("source created", "source_id", src.ID, "path_name", src.PathName)
    return src, nil
}
```

### 骨架 5：HTTP Handler

```go
// internal/api/handler/media_source.go
package handler

import (
    "net/http"

    "goyavision/internal/api/request"
    "goyavision/internal/api/response"
    "goyavision/internal/app/command"
    "goyavision/internal/app/query"
    "goyavision/pkg/apperr"

    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
)

type SourceHandler struct {
    createCmd *command.CreateSourceHandler
    listQry   *query.ListSourcesHandler
}

func NewSourceHandler(c *command.CreateSourceHandler, q *query.ListSourcesHandler) *SourceHandler {
    return &SourceHandler{createCmd: c, listQry: q}
}

func (h *SourceHandler) Register(g *echo.Group) {
    g.GET("/sources", h.List)
    g.POST("/sources", h.Create)
    g.GET("/sources/:id", h.Get)
}

func (h *SourceHandler) Create(c echo.Context) error {
    var req request.CreateSourceReq
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, response.Error(40002, "invalid request body"))
    }
    if err := c.Validate(req); err != nil {
        return err // 由 validator middleware 处理
    }

    src, err := h.createCmd.Handle(c.Request().Context(), command.CreateSourceCommand{
        Name:     req.Name,
        Type:     req.Type,
        URL:      req.URL,
        Protocol: req.Protocol,
        Enabled:  req.Enabled,
    })
    if err != nil {
        return err // 由 error handler 映射
    }

    return response.Created(c, response.SourceFromDomain(src))
}

func (h *SourceHandler) List(c echo.Context) error {
    var req request.ListSourcesReq
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, response.Error(40002, "invalid query"))
    }
    req.Pagination.Normalize()

    result, err := h.listQry.Handle(c.Request().Context(), query.ListSourcesQuery{
        Type:       req.Type,
        Pagination: req.Pagination,
    })
    if err != nil {
        return err
    }

    items := make([]response.SourceResponse, len(result.Items))
    for i, s := range result.Items {
        items[i] = *response.SourceFromDomain(s)
    }
    return response.Paged(c, items, result.Total, req.Limit, req.Offset)
}

func (h *SourceHandler) Get(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return apperr.InvalidInput("invalid id format")
    }
    // ... 调用 query handler
    _ = id
    return nil
}
```

### 骨架 6：HTTP 错误映射中间件

```go
// internal/api/errors.go
package api

import (
    "errors"
    "net/http"

    "goyavision/internal/api/response"
    "goyavision/pkg/apperr"

    "github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
    if c.Response().Committed {
        return
    }

    var appErr *apperr.Error
    if errors.As(err, &appErr) {
        status := codeToHTTPStatus(appErr.Code())
        env := response.Envelope{
            Code:    int(appErr.Code()),
            Message: appErr.Error(),
        }
        if len(appErr.Details()) > 0 {
            env.Details = appErr.Details()
        }
        env.RequestID = c.Response().Header().Get(echo.HeaderXRequestID)
        c.JSON(status, env)
        return
    }

    // Echo 内置错误
    var echoErr *echo.HTTPError
    if errors.As(err, &echoErr) {
        msg, _ := echoErr.Message.(string)
        c.JSON(echoErr.Code, response.Error(echoErr.Code*100+1, msg))
        return
    }

    // 未知错误
    c.JSON(http.StatusInternalServerError, response.Error(50001, "internal server error"))
}

func codeToHTTPStatus(code apperr.Code) int {
    switch {
    case code >= 50000:
        return http.StatusInternalServerError
    case code >= 42200:
        return http.StatusUnprocessableEntity
    case code >= 40900:
        return http.StatusConflict
    case code >= 40400:
        return http.StatusNotFound
    case code >= 40300:
        return http.StatusForbidden
    case code >= 40100:
        return http.StatusUnauthorized
    case code >= 40000:
        return http.StatusBadRequest
    default:
        return http.StatusInternalServerError
    }
}
```

---

## K. 交付与切换方案

### 策略：全量切换 + API 版本化

由于项目处于 V1.0 开发阶段（核心功能仍在进行中），没有线上生产流量，推荐直接切换而非蓝绿/双写。

### 具体方案

1. **在 `refactor` 分支上完成全部重构**，不影响 `master` 上的当前开发
2. **保持 API 路径 `/api/v1/` 不变**，请求/响应格式升级为新的信封格式
3. **数据库表结构不变**（只是 Go 代码侧分离了 model/entity），无需数据迁移
4. **前端兼容层**：如果前端尚未适配新响应格式，在 handler 中提供临时的 legacy response wrapper

### 回滚方案

- `master` 分支始终保持可用
- 重构分支完成后通过 PR merge 到 `master`
- 如果发现严重问题，`git revert` 整个 PR（单次 revert 因为是 squash merge）

### 风险点

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|---------|
| 前端响应格式不兼容 | 中 | 中 | 提供 legacy wrapper middleware，逐步迁移 |
| GORM model 映射遗漏字段 | 低 | 高 | mapper 层单元测试全覆盖 |
| UnitOfWork 事务死锁 | 低 | 高 | 控制事务粒度，添加超时 |
| 重构期间 master 有新功能合入 | 高 | 中 | 定期 rebase master，保持同步 |

---

## L. 工作拆分清单

### Phase 1：基础设施层

| # | 任务 | 产出物 | 验收标准 |
|---|------|--------|---------|
| 1.1 | 建立 `pkg/apperr` 错误类型体系 | `pkg/apperr/errors.go` | `go test ./pkg/apperr/...` 通过；覆盖全部错误码 |
| 1.2 | 建立 `pkg/pagination` 分页工具 | `pkg/pagination/pagination.go` | 单元测试覆盖边界值 |
| 1.3 | 建立 `pkg/logger` 结构化日志 | `pkg/logger/logger.go` (封装 `log/slog`) | 支持 JSON 输出、level 控制、context 提取 request_id |
| 1.4 | 创建 `internal/api/response` 统一信封 | `envelope.go`, `OK/Created/Paged/Error` 函数 | 单元测试验证 JSON 序列化格式 |

### Phase 2：Domain 层重构

| # | 任务 | 产出物 | 验收标准 |
|---|------|--------|---------|
| 2.1 | 拆分 Domain 到子包 | `domain/media/`, `domain/operator/`, `domain/workflow/`, `domain/identity/`, `domain/storage/` | 编译通过；无外部依赖 (`go vet`) |
| 2.2 | 域实体去除 GORM 标签 | 每个实体文件 | 实体文件不 import `gorm.io/*` |
| 2.3 | 添加域业务方法 | `Validate()`, 状态机方法, `NewXxx()` 构造函数 | 单元测试覆盖核心业务规则 |
| 2.4 | 定义按聚合拆分的 Repository 接口 | 每个子包内 `repository.go` | 接口方法数 <= 15/聚合 |
| 2.5 | 定义领域事件 | `domain/event.go` | 事件结构体定义 + 文档注释 |

### Phase 3：基础设施持久化层

| # | 任务 | 产出物 | 验收标准 |
|---|------|--------|---------|
| 3.1 | 创建 GORM Model 层 | `infra/persistence/model/*.go` | 所有模型携带正确 GORM 标签；`AutoMigrate` 通过 |
| 3.2 | 实现 Mapper 层 | `infra/persistence/mapper/*.go` | 单元测试：Domain->Model->Domain 往返转换无损 |
| 3.3 | 按聚合实现 Repository | `infra/persistence/repo/*.go` (7 个文件) | 集成测试连接 PostgreSQL 验证 CRUD |
| 3.4 | 实现 UnitOfWork | `infra/persistence/uow.go` | 事务回滚测试：fn 返回 error 时数据不写入 |
| 3.5 | 迁移 seed 数据 | `infra/persistence/seed.go` | 幂等初始化测试 |

### Phase 4：Application 层

| # | 任务 | 产出物 | 验收标准 |
|---|------|--------|---------|
| 4.1 | 定义出站端口 | `app/port/*.go` (UnitOfWork, MediaGateway, ObjectStorage, TokenService, EventBus) | 接口编译通过 |
| 4.2 | 实现 Command Handlers | `app/command/*.go` (8 个用例) | 单元测试 mock UoW 验证业务流程 |
| 4.3 | 实现 Query Handlers | `app/query/*.go` (6 个用例) | 单元测试 mock Repository 验证查询逻辑 |
| 4.4 | 重构 WorkflowScheduler | `app/scheduler.go` | 正确使用 context；panic recovery；失败补偿 |
| 4.5 | 实现 EventBus | `infra/eventbus/local.go` | 发布/订阅测试 |

### Phase 5：基础设施适配器

| # | 任务 | 产出物 | 验收标准 |
|---|------|--------|---------|
| 5.1 | MediaMTX Gateway 实现 | `infra/mediamtx/gateway.go` | 实现 `app/port/MediaGateway`；Ping 测试 |
| 5.2 | MinIO ObjectStorage 实现 | `infra/minio/client.go` | 实现 `app/port/ObjectStorage` |
| 5.3 | JWT TokenService 实现 | `infra/auth/jwt.go` | 实现 `app/port/TokenService`；Token 生成/验证测试 |
| 5.4 | 升级 DAG 引擎 | `infra/engine/dag_engine.go` | 拓扑排序测试；环路检测测试 |

### Phase 6：API 层

| # | 任务 | 产出物 | 验收标准 |
|---|------|--------|---------|
| 6.1 | Request/Response DTO | `api/request/*.go`, `api/response/*.go` | 编译通过；JSON 标签正确 |
| 6.2 | Handler 重写 | `api/handler/*.go` | 注入 Command/Query Handler 而非直接创建 Service |
| 6.3 | Middleware 重写 | `api/middleware/*.go` (requestid, logging, validator, auth, permission) | Request ID 传播测试；结构化日志验证 |
| 6.4 | 错误映射 | `api/errors.go` | `apperr.Error` -> HTTP 状态码映射测试 |
| 6.5 | Router 注册 | `api/router.go` | 全部路由注册；中间件链正确 |

### Phase 7：组装与集成

| # | 任务 | 产出物 | 验收标准 |
|---|------|--------|---------|
| 7.1 | DI 组装 | `cmd/server/main.go` | 启动成功；所有依赖正确注入 |
| 7.2 | 全链路集成测试 | `tests/integration/*.go` | docker-compose up -> API 调用 -> 数据库验证 |
| 7.3 | 前端兼容验证 | 手动测试 | 前端页面功能正常 |
| 7.4 | 文档更新 | `docs/architecture.md`, `CHANGELOG.md`, `CLAUDE.md` | 文档与代码一致 |

---

## 附录: 关键假设列表

| # | 假设 | 依赖此假设的结论 | 如果假设不成立 |
|---|------|----------------|---------------|
| 1 | 项目当前无生产流量 | K 节切换方案采用全量切换而非蓝绿 | 需增加 API 版本化 `/api/v2/` + 双路由 |
| 2 | PostgreSQL 是唯一数据库 | UnitOfWork 只需 GORM 实现 | 需抽象 DB driver 层 |
| 3 | 不需要 gRPC 支持 | 只设计 HTTP 契约 | 需增加 `internal/api/grpc/` 层 |
| 4 | 算子服务全部通过 HTTP 调用 | HTTPExecutor 是唯一执行器 | 需增加 gRPC/SDK 执行器适配 |
| 5 | 重构期间不新增领域实体 | WBS 以现有 11 个实体为准 | 新实体按同样模式添加 |
| 6 | 前端可以适配新的响应信封格式 | 不需要 legacy wrapper | 需在 Phase 6 增加兼容中间件 |
