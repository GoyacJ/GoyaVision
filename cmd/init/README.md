# GoyaVision 数据库初始化工具

## 概述

`cmd/init/main.go` 用于在开发环境快速完成数据库建表与基础数据初始化。

工具特性：
- 幂等执行：默认模式下已存在数据会跳过。
- 支持强制刷新：`--force` 会重建初始化数据（权限、菜单、角色关联等）。
- 支持演练：`--dry-run` 仅输出流程，不写数据库。

## 初始化内容

当前版本共 7 个步骤：

1. 创建数据库表结构
2. 初始化租户数据
3. 初始化权限数据
4. 初始化菜单数据
5. 初始化角色数据
6. 初始化管理员用户
7. 初始化系统配置

### 1) 表结构（AutoMigrate）

已覆盖核心模块：
- 身份与权限：`tenants`, `users`, `roles`, `permissions`, `menus`, `user_identities`
- 资产与媒体源：`media_assets`, `media_sources`
- 算子中心：`operators`, `operator_versions`, `operator_templates`, `operator_dependencies`
- 算法库：`algorithms`, `algorithm_versions`, `algorithm_implementations`, `algorithm_evaluation_profiles`
- 工作流与上下文：`workflows`, `workflow_revisions`, `workflow_nodes`, `workflow_edges`, `tasks`, `task_context_state`, `task_context_patches`, `task_context_snapshots`
- Agent 运行：`agent_sessions`, `run_events`, `tool_policies`
- 其他：`artifacts`, `files`, `ai_models`, `system_configs`, 个人资产相关表

### 2) 默认种子数据

- 默认租户：`default`
- 权限：104 个（含算法库、工作流版本、任务上下文、Agent 会话）
- 菜单：17 个（含“算法库”“Agent会话”）
- 角色：`super_admin`、`user`
- 管理员账号：`admin / admin123`
- 系统配置：
  - `system.home_path = /assets`
  - `system.public_menus = []`

## 使用方式

### 演练模式

```bash
go run cmd/init/main.go --dry-run
```

### 正常初始化

```bash
go run cmd/init/main.go
```

### 强制刷新初始化数据

```bash
go run cmd/init/main.go --force
```

## 参数

- `--dry-run`：只打印操作，不执行写入。
- `--force`：更新/重建初始化数据（会清理权限、菜单及其关联）。

## 注意事项

- 默认管理员密码仅用于开发联调，部署后必须立即修改。
- 使用 `--force` 前建议备份数据库。
- 本工具目标是“开发初始化”，不负责历史版本迁移。
