# 算子重设计进度审计报告（2026-02-06，二次深度核验）

## 1. 结论摘要

基于 `docs/operator-redesign.md` 的最新二次核验结果：

- 当前实现已从“后端主干 + 前端缺失”阶段，进入“前后端主路径基本贯通、剩余治理与收口问题”阶段。
- 设计文件清单覆盖率：**100%（72/72）**。
- 综合完成度（含前后端）：**约 83%**。
- 后端完成度：**约 88%**。
- 前端完成度：**约 78%**。

---

## 2. 核验范围与方法

- 设计基线：`docs/operator-redesign.md`
- 核验范围：
  - `internal/domain/operator`
  - `internal/port` / `internal/app/port`
  - `internal/app/command` / `internal/app/query`
  - `internal/adapter` / `internal/infra/persistence`
  - `internal/api`
  - `web/src/api` / `web/src/views/operator` / `web/src/views/operator-marketplace` / `web/src/composables`

核验方式：

1. 文件清单逐项对账（设计清单 vs 实际存在）。
2. 关键业务规则逐条核验（发布门禁、版本激活事务、Schema 校验、依赖门禁、MCP链路）。
3. 可运行性验证：
   - `go test ./...` 通过
   - `npm --prefix web run build` 通过（仅 chunk size 警告）

---

## 3. 总体进度

### 3.1 清单覆盖率

- redesign 文档中的核心清单路径：**72**
- 已存在：**72**
- 缺失：**0**

### 3.2 分阶段完成度（A~F）

- Phase A（基础重构）：**89%**
- Phase B（多执行模式）：**82%**
- Phase C（版本管理）：**88%**
- Phase D（Schema 集成）：**84%**
- Phase E（模板市场 + MCP 生态）：**76%**
- Phase F（依赖管理）：**81%**

---

## 4. 从整体到局部的实现核验

### 4.1 Domain / 模型层（完成度：高）

已完成：

- `Operator` 已具备 `Origin`、`ActiveVersionID`、`ActiveVersion`。
- `OperatorVersion`、`ExecConfig`、`OperatorTemplate`、`OperatorDependency` 均已落地。
- `ExecMode` 已收敛为 `http/cli/mcp`。

现状与差距：

- 仍处于兼容迁移态：`Operator` 和 `operators` 表保留旧执行字段（`version/endpoint/method/input_schema/output_spec/config/is_builtin`），尚未完全下沉到 `OperatorVersion`。
- 迁移脚本尚未落地，`migrations` 目录仍为空占位。

### 4.2 Port / UoW / 执行抽象（完成度：高）

已完成：

- `OperatorExecutor` + `ExecutorRegistry` 完整落地。
- `MCPClient/MCPRegistry` 与 `SchemaValidator` Port 已落地。
- UnitOfWork 仓储聚合已扩展版本、模板、依赖仓储。

评价：

- 抽象层设计与文档目标基本一致，后续主要是实现深度问题（特别是 MCP 真接入）。

### 4.3 App/CQRS（完成度：中高）

已完成：

- 生命周期命令：`publish/deprecate/test`。
- 版本命令：`create/activate/rollback/archive`。
- 模板与依赖命令：`install_template/set_operator_dependencies`。
- MCP 命令：`install_mcp_operator/sync_mcp_templates`。
- 对应 Query 与 API 已完整对接。

关键进展：

- `publish` 已实现门禁：`active version`、依赖检查、MCP 健康检查/工具存在校验、Schema 有效性校验。
- `activate_version` 已实现事务内三件事：旧版本归档、新版本激活、`operator.active_version_id` 更新。
- `test_operator` 已从占位实现升级为真实执行器健康检查 + 试运行。
- workflow 创建/更新已接入连接 Schema 校验门禁。

### 4.4 Adapter / Infra（完成度：中高）

已完成：

- HTTP/CLI/MCP 三类执行器 + 注册表 + 路由执行器完成。
- JSON Schema 校验器基于 `jsonschema/v5` 已落地。
- `operator_version/operator_template/operator_dependency` model/mapper/repo 全部落地。

关键差距：

- MCP 仍为 `StaticClient`（内存桩），尚未接入真实 MCP 协议/远程 Server。
- `sync_mcp_templates` 逻辑与 `template_sync` 映射函数未统一复用，存在重复实现与命名规范不一致风险。

### 4.5 API（完成度：高）

已完成：

- operator、versions、dependencies、templates、schema、mcp 全套路由已挂载。
- operator 列表已经 `Preload("ActiveVersion")`，`exec_mode/active_version` 可回传。

评价：

- API 层与 redesign 基本对齐，剩余问题主要在参数约束与数据一致性细节。

### 4.6 Frontend（完成度：中高）

已完成：

- `web/src/api/operator.ts` 已接入版本、依赖、模板、MCP 全量接口。
- `operator/index.vue` 已支持 origin/exec_mode 过滤、发布/弃用/测试、版本与依赖管理。
- redesign 要求组件已落地：`OperatorForm/ExecConfigForm/VersionList/VersionForm/SchemaEditor/DependencyManager/TemplateCard`。
- `operator-marketplace` 页面与 `useJsonSchema` 已落地。

差距：

- `ExecConfigForm` 仍以 JSON 编辑为主，三模式“结构化表单深度”仍可加强。
- 类型定义细节未完全对齐（如 `OperatorVersion.status` 缺少 `testing`）。

---

## 5. 现阶段主要风险与未收口项

1. 兼容字段未收口（Domain/DB 双轨）
- 风险：长期维持双轨字段会放大一致性成本。

2. 迁移脚本缺失
- 风险：跨环境升级时难以保证数据可预期迁移。

3. MCP 仍为静态桩
- 风险：生态能力看似可用，但真实联通能力尚未闭环。

4. 安装模板/MCP 算子后的兼容字段同步不充分
- 现象：`install_template/install_mcp_operator` 绑定 active version 后未调用兼容字段同步辅助函数。
- 风险：依赖版本检查仍读取 operator 兼容字段，可能出现版本字段不一致。

5. Create 命令的约束仍偏弱
- 现象：`create_operator/create_operator_version` 未严格校验 `exec_mode` 枚举与版本格式。
- 风险：无效数据可入库，错误延后到运行期暴露。

6. Schema 校验尚未覆盖运行时 I/O 全链路
- 现象：当前主要覆盖“创建/发布/连接校验”，执行时未统一校验 `input/output` 载荷。

---

## 6. 与上版报告相比的关键变化

已修复或显著推进：

- redesign 清单文件已全部补齐（由缺失状态变为 100% 覆盖）。
- 前端版本管理、依赖管理、模板市场、MCP 交互基础页面已落地。
- `test_operator` 从占位实现升级为真实执行测试。
- workflow 创建/更新已接入连接校验。
- 依赖校验已加入 `min_version` 版本比较逻辑。
- operator 列表已支持 ActiveVersion 预加载。

---

## 7. 验证记录

- 后端：`go test ./...`（通过）
- 前端：`npm --prefix web run build`（通过）

---

## 8. 下一步优先级建议

1. 先收口兼容层：补迁移脚本并规划旧字段退役路径。  
2. 修复安装路径的一致性：在模板安装与 MCP 安装后同步兼容字段。  
3. 强化输入约束：补 `exec_mode` 枚举校验与 version（semver）校验。  
4. 推进 MCP 真接入：替换 `StaticClient`，打通真实服务发现与调用链。  
5. 前端深化：将 `ExecConfigForm` 从 JSON 编辑器升级为模式化结构表单。
