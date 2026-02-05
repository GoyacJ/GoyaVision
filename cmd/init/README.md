# GoyaVision 数据库初始化工具

## 概述

此工具用于初始化 GoyaVision 数据库，包括创建表结构和初始化默认数据。支持多次执行（幂等性），已存在的数据会被跳过。

## 功能

### 1. 创建数据库表结构
- 自动创建 V1.0 所需的所有表
- 如果表已存在，会自动更新结构（添加缺失字段）
- 创建的表包括：
  - 认证授权：users, roles, permissions, menus
  - 媒体管理：media_sources, media_assets
  - 算子与工作流：operators, workflows, workflow_nodes, workflow_edges
  - 任务与产物：tasks, artifacts
  - 文件管理：files

### 2. 初始化权限数据
- 创建 45 个系统权限（asset、source、operator、workflow、task、artifact、user、role、menu、file）
- 已存在的权限会被跳过（除非使用 `--force`）

### 3. 初始化菜单数据
- 创建 10 个系统菜单（媒体资产、媒体源、算子管理、工作流、任务管理、系统管理及其子菜单）
- 已存在的菜单会被跳过（除非使用 `--force`）

### 4. 初始化角色数据
- 创建超级管理员角色（super_admin）
- 分配所有权限和菜单
- 已存在的角色会被跳过（除非使用 `--force`）

### 5. 初始化管理员用户
- 创建默认管理员用户
  - 用户名：`admin`
  - 密码：`admin123`
- 分配超级管理员角色
- 已存在的用户会被跳过（除非使用 `--force`）

## 使用方法

### 1. 模拟运行（推荐首次使用）

```bash
go run cmd/init/main.go --dry-run
```

这将显示所有将要执行的操作，但不会修改数据库。

### 2. 正常初始化

```bash
go run cmd/init/main.go
```

执行时会要求确认，输入 `y` 继续，输入 `N` 取消。

### 3. 强制重新初始化

```bash
go run cmd/init/main.go --force
```

⚠️ **警告**：此操作会删除现有数据后重新创建，请谨慎使用！

## 参数说明

- `--dry-run`：模拟运行模式，只显示操作不实际执行
- `--force`：强制模式，删除现有数据后重新创建

## 幂等性保证

- **默认模式**：已存在的数据会被跳过，不会重复创建
- **强制模式**：删除现有数据后重新创建

## 配置要求

- 需要正确配置 `configs/config.<env>.yaml` 中的数据库连接
- 或通过环境变量设置 `GOYAVISION_DB_DSN`

## 输出示例

```
GoyaVision 数据库初始化工具 v1.0
====================================

📊 初始化计划:
1. 创建数据库表结构
2. 初始化权限数据
3. 初始化菜单数据
4. 初始化角色数据（超级管理员）
5. 初始化管理员用户

是否继续？ [y/N]: y

开始初始化...

[1/5] 创建数据库表结构
  创建 V1.0 表结构...
  ✓ 已创建/更新以下表:
    - users, roles, permissions, menus
    - media_sources, media_assets
    - operators
    - workflows, workflow_nodes, workflow_edges
    - tasks, artifacts
    - files
✅ 数据库表结构创建完成

[2/5] 初始化权限数据
  ✓ 创建权限: 查看媒体资产列表
  ✓ 创建权限: 创建媒体资产
  ...
  ✓ 新增权限: 45 个
✅ 权限数据初始化完成

[3/5] 初始化菜单数据
  ✓ 创建菜单: 媒体资产
  ✓ 创建菜单: 媒体源
  ...
  ✓ 新增菜单: 10 个
✅ 菜单数据初始化完成

[4/5] 初始化角色数据
  创建超级管理员角色...
  ✓ 已创建超级管理员角色
  分配权限和菜单...
  ✓ 已分配 45 个权限
  ✓ 已分配 10 个菜单
✅ 角色数据初始化完成

[5/5] 初始化管理员用户
  创建管理员用户...
  ✓ 已创建管理员用户
✅ 管理员用户初始化完成

✅ 数据库初始化完成！

默认管理员账号:
  用户名: admin
  密码: admin123
  ⚠️  生产环境请立即修改密码！
```

## 多次执行

### 正常执行（幂等性）

```bash
# 第一次执行
go run cmd/init/main.go
# 输出：创建所有数据

# 第二次执行（相同命令）
go run cmd/init/main.go
# 输出：跳过已存在的数据，不重复创建
```

### 强制重新初始化

```bash
go run cmd/init/main.go --force
# 输出：删除现有数据后重新创建
```

## 与迁移工具的区别

| 功能 | 初始化工具 (`cmd/init`) | 迁移工具 (`cmd/migrate`) |
|------|------------------------|-------------------------|
| **用途** | 新项目初始化 | 旧版本数据迁移 |
| **表创建** | ✅ 创建所有表 | ✅ 创建所有表 |
| **数据初始化** | ✅ 初始化默认数据 | ❌ 不初始化默认数据 |
| **数据迁移** | ❌ 不迁移旧数据 | ✅ 迁移旧数据 |
| **清理旧表** | ❌ 不清理 | ✅ 清理废弃表 |

## 注意事项

⚠️ **重要提示**

1. **生产环境安全**
   - 默认管理员密码为 `admin123`，生产环境必须立即修改
   - 建议使用强密码策略

2. **数据备份**
   - 使用 `--force` 前请备份数据库
   - 强制模式会删除现有数据

3. **多次执行**
   - 默认模式支持多次执行，不会重复创建数据
   - 使用 `--force` 会删除现有数据后重新创建

4. **配置要求**
   - 需要正确配置数据库连接
   - 确保数据库用户有创建表和插入数据的权限

## 验证初始化结果

```bash
# 启动服务
./bin/goyavision

# 使用默认账号登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 检查菜单
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/menus/tree

# 检查权限
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/permissions
```

## 常见问题

**Q: 初始化后如何修改管理员密码？**

A: 可以通过 API 或数据库直接修改：
```bash
# 通过 API（需要先登录）
curl -X PUT http://localhost:8080/api/v1/auth/password \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"old_password":"admin123","new_password":"新密码"}'
```

**Q: 初始化失败怎么办？**

A: 检查：
1. 数据库连接配置是否正确
2. 数据库用户是否有足够权限
3. 查看错误日志定位问题

**Q: 可以只初始化表结构不初始化数据吗？**

A: 可以，但需要修改代码。当前工具会同时初始化表和数据。如果只需要表结构，可以使用迁移工具的 `--dry-run` 模式查看表创建步骤。
