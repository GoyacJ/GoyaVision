# 算子管理模块全面重设计方案（最终收敛版）

> **版本**: 1.1  
> **日期**: 2026-02-06  
> **范围**: Operator 模块从"基础 CRUD"升级为"多执行模式（HTTP/CLI/MCP）+ 版本化 + 生态市场"  
> **基于**: 当前 `internal/domain/operator/` 及相关层逐文件分析，并纳入 MCP 支持要求

---

## 目录

- [一、现状分析与设计动机](#一现状分析与设计动机)
- [二、Domain 层重设计](#二domain-层重设计)
- [三、Port 层重设计](#三port-层重设计)
- [四、App 层 CQRS 重设计](#四app-层-cqrs-重设计)
- [五、Adapter 层实现](#五adapter-层实现)
- [六、API 层重设计](#六api-层重设计)
- [七、前端重设计](#七前端重设计)
- [八、分阶段实施路线](#八分阶段实施路线)
- [九、文件变更清单](#九文件变更清单)
- [十、新增依赖](#十新增依赖)

---

## 一、现状分析与设计动机

### 1.1 当前状态

当前算子模块（`internal/domain/operator/operator.go`）为基础 CRUD 模型，元数据、执行配置、Schema 和版本信息均混在单一实体中。

### 1.2 核心局限

| # | 问题 | 影响 | 改进方向 |
|---|------|------|---------|
| 1 | 执行模式单一（HTTP） | 无法接入本地命令与 MCP 工具生态 | 多态 ExecConfig + 执行器注册表 |
| 2 | 版本语义缺失 | 无法灰度、回滚、稳定运营 | `Operator` 与 `OperatorVersion` 解耦 |
| 3 | Schema 缺少强校验 | 运行时错误暴露晚，编排风险高 | JSON Schema 全流程校验 |
| 4 | 生态能力弱 | 模板复用与共享能力不足 | 模板市场 + MCP 生态接入 |
| 5 | 生命周期过于简化 | 缺少测试/发布/弃用治理 | 完整状态机与发布门禁 |

### 1.3 设计目标

```text
当前：Operator = 元数据 + 执行配置 + Schema（扁平结构）
目标：Operator = 元数据 + Origin + 生命周期
      └── OperatorVersion = 执行配置 + Schema + 版本状态（多版本结构）
      └── OperatorTemplate = 市场模板（生态体系）
      └── OperatorDependency = 依赖关系（依赖管理）
```

---

## 二、Domain 层重设计

### 2.1 Operator 实体重构

**文件**: `internal/domain/operator/operator.go`

**变更要点**:
- 移除字段：`Endpoint`, `Method`, `InputSchema`, `OutputSpec`, `Config`, `Version`（下沉到 `OperatorVersion`）
- `IsBuiltin` 升级为 `Origin` 枚举
- 增加 `ActiveVersionID` 与 `ActiveVersion`
- 生命周期状态扩展为 `draft/testing/published/deprecated`

```go
type ExecMode string

const (
    ExecModeHTTP ExecMode = "http"
    ExecModeCLI  ExecMode = "cli"
    ExecModeMCP  ExecMode = "mcp"
)

type Origin string

const (
    OriginBuiltin     Origin = "builtin"
    OriginCustom      Origin = "custom"
    OriginMarketplace Origin = "marketplace"
    OriginMCP         Origin = "mcp"
)

type Status string

const (
    StatusDraft      Status = "draft"
    StatusTesting    Status = "testing"
    StatusPublished  Status = "published"
    StatusDeprecated Status = "deprecated"
)

type Operator struct {
    ID              uuid.UUID
    Code            string
    Name            string
    Description     string
    Category        Category
    Type            Type
    Origin          Origin
    ActiveVersionID *uuid.UUID
    Status          Status
    Tags            []string
    CreatedAt       time.Time
    UpdatedAt       time.Time

    ActiveVersion   *OperatorVersion
}
```

**Filter 扩展**:

```go
type Filter struct {
    Category *Category
    Type     *Type
    Status   *Status
    Origin   *Origin
    ExecMode *ExecMode
    Tags     []string
    Keyword  string
    Limit    int
    Offset   int
}
```

### 2.2 OperatorVersion 实体（新增）

**新文件**: `internal/domain/operator/version.go`

```go
type VersionStatus string

const (
    VersionStatusDraft    VersionStatus = "draft"
    VersionStatusTesting  VersionStatus = "testing"
    VersionStatusActive   VersionStatus = "active"
    VersionStatusArchived VersionStatus = "archived"
)

type OperatorVersion struct {
    ID          uuid.UUID
    OperatorID  uuid.UUID
    Version     string                 // semver
    ExecMode    ExecMode
    ExecConfig  *ExecConfig
    InputSchema map[string]interface{}
    OutputSpec  map[string]interface{}
    Config      map[string]interface{}
    Changelog   string
    Status      VersionStatus
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### 2.3 ExecConfig 值对象（新增）

**新文件**: `internal/domain/operator/exec_config.go`

```go
type ExecConfig struct {
    HTTP *HTTPExecConfig `json:"http,omitempty"`
    CLI  *CLIExecConfig  `json:"cli,omitempty"`
    MCP  *MCPExecConfig  `json:"mcp,omitempty"`
}

type HTTPExecConfig struct {
    Endpoint   string            `json:"endpoint"`
    Method     string            `json:"method"`
    Headers    map[string]string `json:"headers,omitempty"`
    TimeoutSec int               `json:"timeout_sec,omitempty"`
    AuthType   string            `json:"auth_type,omitempty"`
    AuthConfig map[string]string `json:"auth_config,omitempty"`
}

type CLIExecConfig struct {
    Command    string            `json:"command"`
    Args       []string          `json:"args"`
    WorkDir    string            `json:"work_dir,omitempty"`
    Env        map[string]string `json:"env,omitempty"`
    TimeoutSec int               `json:"timeout_sec,omitempty"`
}

type MCPExecConfig struct {
    ServerID      string                 `json:"server_id"`
    ToolName      string                 `json:"tool_name"`
    ToolVersion   string                 `json:"tool_version,omitempty"`
    TimeoutSec    int                    `json:"timeout_sec,omitempty"`
    InputMapping  map[string]interface{} `json:"input_mapping,omitempty"`
    OutputMapping map[string]interface{} `json:"output_mapping,omitempty"`
}
```

### 2.4 OperatorTemplate 实体（新增）

**新文件**: `internal/domain/operator/template.go`

保持模板市场模型不变，但 `ExecMode` 限定为 `http/cli/mcp`。

### 2.5 OperatorDependency 实体（新增）

**新文件**: `internal/domain/operator/dependency.go`

保持依赖模型不变，用于发布前校验。

---

## 三、Port 层重设计

### 3.1 Repository 接口扩展

**文件**: `internal/domain/operator/repository.go`

- `ListEnabled` → `ListPublished`
- 新增 `GetWithActiveVersion`
- 新增 `VersionRepository` / `TemplateRepository` / `DependencyRepository`

### 3.2 OperatorExecutor + ExecutorRegistry

**文件**: `internal/port/engine.go`

```go
type OperatorExecutor interface {
    Execute(ctx context.Context, version *operator.OperatorVersion, input *operator.Input) (*operator.Output, error)
    Mode() operator.ExecMode
    HealthCheck(ctx context.Context, version *operator.OperatorVersion) error
}

type ExecutorRegistry interface {
    Register(mode operator.ExecMode, executor OperatorExecutor)
    Get(mode operator.ExecMode) (OperatorExecutor, error)
}
```

### 3.3 SchemaValidator

**新文件**: `internal/app/port/schema_validator.go`

保持原设计：`ValidateInput/ValidateOutput/ValidateConnection/IsValidJSONSchema`。

### 3.4 MCP Port（新增）

**新文件建议**: `internal/port/mcp.go`

```go
type MCPClient interface {
    ListTools(ctx context.Context, serverID string) ([]MCPTool, error)
    CallTool(ctx context.Context, serverID, toolName string, args map[string]interface{}) (map[string]interface{}, error)
    HealthCheck(ctx context.Context, serverID string) error
}

type MCPRegistry interface {
    ListServers(ctx context.Context) ([]MCPServer, error)
    GetServer(ctx context.Context, serverID string) (*MCPServer, error)
}
```

### 3.5 UnitOfWork Repositories 扩展

`internal/app/port/unit_of_work.go` 增加：

```go
type Repositories struct {
    OperatorVersions     operator.VersionRepository
    OperatorTemplates    operator.TemplateRepository
    OperatorDependencies operator.DependencyRepository
}
```

---

## 四、App 层 CQRS 重设计

### 4.1 Command 变更

保留原先版本化命令，并新增 MCP 生态命令：

| 文件 | 操作 | 说明 |
|------|------|------|
| `publish_operator.go` | 新增 | 发布算子（含门禁校验） |
| `deprecate_operator.go` | 新增 | 弃用算子 |
| `create_operator_version.go` | 新增 | 创建新版本 |
| `activate_version.go` | 新增 | 激活版本（事务） |
| `rollback_version.go` | 新增 | 回滚版本 |
| `archive_version.go` | 新增 | 归档版本 |
| `test_operator.go` | 新增 | 连通性/试运行测试 |
| `install_template.go` | 新增 | 从模板安装 |
| `set_operator_dependencies.go` | 新增 | 设置依赖 |
| `sync_mcp_templates.go` | 新增 | 从 MCP 同步模板 |
| `install_mcp_operator.go` | 新增 | 从 MCP Tool 安装算子 |

### 4.2 Query 变更

保留原先 Query，并新增：

| 文件 | 操作 | 说明 |
|------|------|------|
| `list_mcp_servers.go` | 新增 | MCP Server 列表 |
| `list_mcp_tools.go` | 新增 | MCP Tool 列表 |
| `preview_mcp_tool.go` | 新增 | MCP Tool 预览（Schema/元信息） |

### 4.3 强约束业务规则

1. `publish` 前必须存在 `active version`。
2. `publish` 前必须通过依赖校验。
3. `exec_mode=mcp` 时，`publish` 前必须通过：`server health check + tool exists + schema valid`。
4. `activate_version` 必须同事务更新：新版本 active、旧版本 archived、operator.active_version_id。

---

## 五、Adapter 层实现

### 5.1 执行器实现

**目录**: `internal/adapter/engine/`

| 文件 | 操作 | 说明 |
|------|------|------|
| `executor_registry.go` | 新增 | 执行器注册表 |
| `http_executor.go` | 修改 | 从 `version.ExecConfig.HTTP` 读取配置 |
| `cli_executor.go` | 新增 | 本地命令执行 |
| `mcp_executor.go` | 新增 | 调用 MCP Tool |
| `simple_engine.go` | 修改 | 使用 `ExecutorRegistry` 按 `ExecMode` 路由 |

> 本期不实现：`docker_executor.go`、`grpc_executor.go`。

### 5.2 Schema 校验器（新增）

**新文件**: `internal/adapter/schema/json_schema_validator.go`

基于 `github.com/santhosh-tekuri/jsonschema/v5` 实现 Schema 校验与兼容性检查。

### 5.3 MCP 适配器（新增）

| 文件 | 说明 |
|------|------|
| `internal/adapter/mcp/client.go` | MCP 客户端适配 |
| `internal/adapter/mcp/template_sync.go` | MCP Tool → OperatorTemplate 映射同步 |

### 5.4 数据库模型变更

- 修改 `internal/infra/persistence/model/operator.go`：新增 `origin/active_version_id`，移除旧执行字段。
- 新增 `operator_versions/operator_templates/operator_dependencies`。
- 增加对应 mapper/repo 文件。

### 5.5 数据迁移策略

1. 现有 `operators` 每条生成一个 `operator_versions`，`exec_mode = http`。
2. `operators.active_version_id` 指向新版本。
3. `is_builtin` 映射到 `origin`。
4. `enabled → published`，`disabled/draft → draft`。
5. 删除废弃列。

---

## 六、API 层重设计

### 6.1 算子、版本、依赖、模板（沿用并扩展）

保持 `/api/v1/operators` 体系，新增和调整如下：

- CRUD + `origin/exec_mode` 筛选
- 生命周期：`/publish`、`/deprecate`、`/test`
- 版本管理：`/versions`、`/activate`、`/rollback`
- 依赖管理：`/dependencies`
- 模板市场：`/templates` + `install`
- Schema 校验：`/validate-schema`、`/validate-connection`

### 6.2 MCP API（新增）

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/operators/mcp/servers` | MCP Server 列表 |
| `GET` | `/operators/mcp/servers/:id/tools` | MCP Tool 列表 |
| `GET` | `/operators/mcp/servers/:id/tools/:tool/preview` | Tool 预览 |
| `POST` | `/operators/mcp/install` | 从 MCP Tool 安装算子 |
| `POST` | `/operators/mcp/sync-templates` | 同步 MCP 模板 |

### 6.3 DTO 变更

- `OperatorResponse` 增加 `origin/active_version/exec_mode`。
- `OperatorCreateReq` 使用 `exec_mode + exec_config`。
- 新增 `MCPInstallReq/MCPToolPreviewResp/MCPServerResp` 等 DTO。

---

## 七、前端重设计

### 7.1 API 客户端

**文件**: `web/src/api/operator.ts`

- 更新模型：仅保留 `http/cli/mcp`。
- 增加 MCP 接口：`listMCPServers/listMCPTools/previewMCPTool/installMCPTool/syncMCPTemplates`。

### 7.2 算子管理页

**文件**: `web/src/views/operator/index.vue`

- 筛选增加 `Origin` 与 `ExecMode`。
- 详情展示 ActiveVersion 与执行配置。

### 7.3 组件变更

**目录**: `web/src/views/operator/components/`

- `ExecConfigForm.vue`：仅三面板 `HTTP / CLI / MCP`。
- 其他组件维持：`OperatorForm/VersionList/VersionForm/SchemaEditor/DependencyManager/TemplateCard`。

### 7.4 模板市场页面

保留 `web/src/views/operator-marketplace/index.vue`，新增 MCP 来源标签与安装入口。

### 7.5 useJsonSchema

保留 `web/src/composables/useJsonSchema.ts`。

---

## 八、分阶段实施路线

### Phase A：基础重构（Domain + 版本化数据模型）

目标：完成模型重构与数据迁移，确保 CRUD/工作流不回归。

### Phase B：多执行模式（HTTP + CLI + MCP）

| 步骤 | 内容 |
|------|------|
| B.1 | 实现 ExecutorRegistry |
| B.2 | 实现 CLIExecutor |
| B.3 | 实现 MCPExecutor + MCPClient |
| B.4 | 前端 ExecConfigForm（三模式） |
| B.5 | 验证：CLI 执行 ffmpeg 截帧 + MCP Tool 调用通过 |

### Phase C：版本管理

目标：实现发布、弃用、激活、回滚全流程。

### Phase D：JSON Schema 集成

目标：创建、执行、连接三个环节全量校验。

### Phase E：模板市场 + MCP 生态

目标：模板浏览、安装、MCP 同步一体化。

### Phase F：依赖管理

目标：依赖配置、发布前门禁、阻断机制。

---

## 九、文件变更清单

### 需修改的现有文件（核心）

- `internal/domain/operator/operator.go`
- `internal/domain/operator/repository.go`
- `internal/port/engine.go`
- `internal/app/port/unit_of_work.go`
- `internal/app/command/create_operator.go`
- `internal/app/command/update_operator.go`
- `internal/app/command/delete_operator.go`
- `internal/app/query/get_operator.go`
- `internal/app/query/list_operators.go`
- `internal/adapter/engine/http_executor.go`
- `internal/adapter/engine/simple_engine.go`
- `internal/infra/persistence/model/operator.go`
- `internal/infra/persistence/mapper/operator.go`
- `internal/infra/persistence/repo/operator.go`
- `internal/api/handler/operator.go`
- `internal/api/handler/handlers.go`
- `internal/api/dto/operator.go`
- `internal/api/router.go`
- `web/src/api/operator.ts`
- `web/src/views/operator/index.vue`

### 需新增的文件（核心）

#### Domain
- `internal/domain/operator/version.go`
- `internal/domain/operator/exec_config.go`
- `internal/domain/operator/template.go`
- `internal/domain/operator/dependency.go`

#### Port / App
- `internal/port/mcp.go`
- `internal/app/port/schema_validator.go`
- `internal/app/command/publish_operator.go`
- `internal/app/command/deprecate_operator.go`
- `internal/app/command/create_operator_version.go`
- `internal/app/command/activate_version.go`
- `internal/app/command/rollback_version.go`
- `internal/app/command/archive_version.go`
- `internal/app/command/test_operator.go`
- `internal/app/command/install_template.go`
- `internal/app/command/set_operator_dependencies.go`
- `internal/app/command/sync_mcp_templates.go`
- `internal/app/command/install_mcp_operator.go`
- `internal/app/query/list_operator_versions.go`
- `internal/app/query/get_operator_version.go`
- `internal/app/query/list_templates.go`
- `internal/app/query/get_template.go`
- `internal/app/query/list_operator_dependencies.go`
- `internal/app/query/check_dependencies.go`
- `internal/app/query/validate_schema.go`
- `internal/app/query/validate_connection.go`
- `internal/app/query/list_mcp_servers.go`
- `internal/app/query/list_mcp_tools.go`
- `internal/app/query/preview_mcp_tool.go`

#### Adapter / Infra
- `internal/adapter/engine/executor_registry.go`
- `internal/adapter/engine/cli_executor.go`
- `internal/adapter/engine/mcp_executor.go`
- `internal/adapter/schema/json_schema_validator.go`
- `internal/adapter/mcp/client.go`
- `internal/adapter/mcp/template_sync.go`
- `internal/infra/persistence/model/operator_version.go`
- `internal/infra/persistence/model/operator_template.go`
- `internal/infra/persistence/model/operator_dependency.go`
- `internal/infra/persistence/mapper/operator_version.go`
- `internal/infra/persistence/mapper/operator_template.go`
- `internal/infra/persistence/mapper/operator_dependency.go`
- `internal/infra/persistence/repo/operator_version.go`
- `internal/infra/persistence/repo/operator_template.go`
- `internal/infra/persistence/repo/operator_dependency.go`

#### Frontend
- `web/src/views/operator/components/OperatorForm.vue`
- `web/src/views/operator/components/ExecConfigForm.vue`
- `web/src/views/operator/components/VersionList.vue`
- `web/src/views/operator/components/VersionForm.vue`
- `web/src/views/operator/components/SchemaEditor.vue`
- `web/src/views/operator/components/DependencyManager.vue`
- `web/src/views/operator/components/TemplateCard.vue`
- `web/src/views/operator-marketplace/index.vue`
- `web/src/composables/useJsonSchema.ts`

> 本清单已移除 Docker/gRPC 相关文件。

---

## 十、新增依赖

| 依赖 | 用途 | 引入阶段 |
|------|------|---------|
| `github.com/santhosh-tekuri/jsonschema/v5` | JSON Schema 校验 | Phase D |

> 本期不引入：`github.com/docker/docker/client`、`google.golang.org/grpc`。

---

## 文档结构概览（最终）

1. 现状分析与设计动机（含 MCP 场景）
2. Domain 重构（Operator + Version + ExecConfig + Template + Dependency）
3. Port 重构（Repository、ExecutorRegistry、SchemaValidator、MCP Port）
4. App CQRS（版本化命令查询 + MCP 同步/安装）
5. Adapter/Infra（HTTP/CLI/MCP 执行器、Schema、迁移）
6. API（算子/版本/依赖/模板 + MCP 专属端点）
7. 前端（三执行模式 UI 与模板市场）
8. 分阶段实施（A~F）
9. 文件变更清单（Docker/gRPC 已移除）
10. 新增依赖（仅保留本期必需）