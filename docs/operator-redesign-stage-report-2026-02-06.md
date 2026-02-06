# 算子重设计阶段性报告（2026-02-06）

## 一、背景与目标

本阶段基于 `docs/operator-redesign.md` 推进算子模块重设计，目标是先完成可落地的后端基础能力：

1. 完成 Phase A 的核心版本化数据模型（Operator + OperatorVersion）。
2. 打通 v1.1 的 MCP 最小业务闭环（发现工具、安装算子、同步模板）。
3. 在不破坏现有工作流主链路的前提下，逐步切换执行语义到 ActiveVersion。

---

## 二、本阶段已完成内容

## 2.1 Domain 与持久化基础（Phase A 主体）

已完成新增实体与持久化结构：

- Domain：
  - `internal/domain/operator/version.go`
  - `internal/domain/operator/exec_config.go`
  - `internal/domain/operator/template.go`
  - `internal/domain/operator/dependency.go`
- Repository 接口扩展：
  - `VersionRepository`
  - `TemplateRepository`
  - `DependencyRepository`
  - `GetWithActiveVersion` / `ListPublished`
- Infra 持久化新增：
  - model / mapper / repo 的 `operator_version`、`operator_template`、`operator_dependency`

同时，`create_operator`、`get_operator`、`list_operators`、`delete_operator` 已按版本化方向做了适配（首版本创建、ActiveVersion 读取、筛选扩展、级联清理）。

## 2.2 执行链路改造（部分完成）

- `internal/port/engine.go` 已升级为以 `OperatorVersion` 执行为核心。
- 新增 `ExecutorRegistry` 接口与实现（`internal/adapter/engine/executor_registry.go`）。
- `simple_engine` / `http_executor` 已按 ActiveVersion 和 ExecConfig 语义调整。

## 2.3 MCP 最小闭环（v1.1）

### 查询能力

- `GET /api/v1/operators/mcp/servers`
- `GET /api/v1/operators/mcp/servers/:id/tools`
- `GET /api/v1/operators/mcp/servers/:id/tools/:tool/preview`

### 写入能力

- `POST /api/v1/operators/mcp/install`
  - 从 MCP Tool 安装为 Operator（origin=mcp）
  - 自动创建首个 MCP 执行版本
- `POST /api/v1/operators/mcp/sync-templates`
  - 批量同步 MCP Tool 到 operator_templates
  - 返回 total/created/updated 统计

### 对应应用层与 DTO

- Command：
  - `InstallMCPOperatorHandler`
  - `SyncMCPTemplatesHandler`
- DTO：
  - `MCPInstallReq`
  - `SyncMCPTemplatesReq`
  - `SyncMCPTemplatesResponse`

---

## 三、当前完成度评估（对照 redesign 文档）

- **Phase A（基础重构）**：约 **80%**
  - 核心实体、Repo、持久化与主流程改造已完成。
  - 迁移脚本与端到端回归仍待补强。
- **Phase B（多执行模式）**：约 **20%**
  - Registry 已有；CLI/Docker/gRPC 执行器尚未落地。
- **Phase C（版本管理）**：约 **35%**
  - 发布/弃用/测试已落地；Activate/Rollback/Archive 未完成。
- **Phase D（Schema 集成）**：约 **10%**
  - 尚未落地 SchemaValidator 端口与 adapter 实现。
- **Phase E（模板市场）**：约 **30%**
  - 模板实体与持久化具备，完整模板 API/前端市场未完成。
- **Phase F（依赖管理）**：约 **20%**
  - 依赖实体与持久化具备，发布前依赖校验与 API/UI 未完成。

综合完成度（后端+文档）：约 **40%~45%**。

---

## 四、已识别风险与待办

1. MCP 客户端注入需收口到统一依赖装配，避免运行期不可用。
2. 版本状态机仍缺关键命令（activate/rollback/archive），线上回滚闭环未完成。
3. 缺少 JSON Schema 质量门禁，工作流连接仍存在运行时风险。

---

## 五、下一阶段建议（按优先级）

1. **先完成版本管理闭环（C）**：CreateVersion + Activate + Rollback + Archive + 对应 API。
2. **补 Schema 门禁（D）**：SchemaValidator 端口、jsonschema adapter、创建/执行/连接三处校验。
3. **推进执行器扩展（B）**：优先 CLI，再 Docker，再 gRPC，并同步安全约束。

以上完成后，再进入模板市场 UI 与依赖发布治理（E/F）。
