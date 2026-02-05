# GoyaVision 数据迁移工具

## 概述

此工具用于将 GoyaVision 从旧架构迁移到 V1.0 新架构。

## 迁移内容

### 0. 创建数据库表结构
- 自动创建 V1.0 所需的所有表
- 如果表已存在，会自动更新结构（添加缺失字段）
- 创建的表包括：
  - 认证授权：users, roles, permissions, menus
  - 媒体管理：media_sources, media_assets
  - 算子与工作流：operators, workflows, workflow_nodes, workflow_edges
  - 任务与产物：tasks, artifacts
  - 文件管理：files

### 1. 更新菜单和权限（V1.0 新功能）
- 清理旧菜单和权限（stream、algorithm、inference 相关）
- 添加新菜单（媒体资产、媒体源、算子管理、工作流、任务管理、系统管理）
- 添加新权限（asset、source、operator、workflow、task、artifact、user、role、menu、file）
- 更新超级管理员角色权限

### 2. Streams → MediaSources
- 将所有视频流转换为媒体源（MediaSource）
- 自动识别协议类型（rtsp/rtmp/hls）
- 生成 PathName（用于 MediaMTX）
- 保留原始 ID 和启用状态

### 3. Streams → MediaAssets
- 将所有视频流转换为媒体资产（MediaAsset）
- 类型设为 `stream`
- 来源类型设为 `live`
- 关联到对应的媒体源（SourceID）
- 状态映射：`enabled=true` → `ready`, `enabled=false` → `pending`

### 4. Algorithms → Operators
- 将所有算法转换为算子（Operator）
- 分类统一设为 `analysis`
- 类型添加 `legacy_` 前缀
- 状态设为 `enabled`
- 保留 InputSpec、OutputSpec、Config（JSON 格式）

### 5. 清理废弃表
- 删除 `algorithm_bindings` 表
- 删除 `inference_results` 表
- 删除 `streams` 表（已迁移到 media_sources 和 media_assets）
- 删除 `record_sessions` 表（如果存在）

## 使用方法

### 1. 模拟运行（推荐首次使用）

```bash
go run cmd/migrate/main.go --dry-run
```

这将显示所有将要执行的操作，但不会修改数据库。

### 2. 正式迁移

```bash
go run cmd/migrate/main.go
```

执行时会要求确认，输入 `y` 继续，输入 `N` 取消。

## 注意事项

⚠️ **重要提示**

1. **备份数据库**
   ```bash
   pg_dump goyavision > backup_$(date +%Y%m%d_%H%M%S).sql
   ```

2. **停止服务**
   - 迁移前请停止 GoyaVision 服务
   - 确保没有正在运行的任务

3. **不可逆操作**
   - 旧表会被永久删除
   - 请确保已备份重要数据

4. **配置要求**
  - 需要正确配置 `configs/config.<env>.yaml` 中的数据库连接
   - 或通过环境变量设置 `GOYAVISION_DB_DSN`

## 输出示例

```
GoyaVision 数据迁移工具 v1.0
================================

📊 数据迁移计划:
0. 创建数据库表结构（如果不存在）
1. 更新菜单和权限（V1.0 新功能）
2. 迁移 streams → media_sources（媒体源）
3. 迁移 streams → media_assets（媒体资产）
4. 迁移 algorithms → operators（算子）
5. 清理废弃表（algorithm_bindings、inference_results、streams、record_sessions）

是否继续？ [y/N]: y

开始迁移...

[0/5] 创建数据库表结构
  创建 V1.0 表结构...
  ✓ 已创建/更新以下表:
    - users, roles, permissions, menus
    - media_sources, media_assets
    - operators
    - workflows, workflow_nodes, workflow_edges
    - tasks, artifacts
    - files
✅ 数据库表结构创建完成

[1/5] 更新菜单和权限
  清理旧菜单...
  ✓ 删除旧菜单: stream
  ✓ 删除旧菜单: algorithm
  清理旧权限...
  ✓ 删除旧权限: stream:list
  ✓ 删除旧权限: algorithm:list
  添加新菜单...
  ✓ 创建新菜单: 媒体资产
  ✓ 创建新菜单: 媒体源
  ✓ 创建新菜单: 算子管理
  ✓ 创建新菜单: 工作流
  ✓ 创建新菜单: 任务管理
  ✓ 新增菜单: 10 个
  添加新权限...
  ✓ 新增权限: 45 个
  更新超级管理员角色权限...
  ✓ 已更新超级管理员权限
✅ 菜单和权限更新完成

[2/5] 迁移 Streams → MediaSources
找到 5 个流
  ✓ 迁移流: Stream1 → 媒体源 ID: xxx-xxx-xxx
  ✓ 迁移流: Stream2 → 媒体源 ID: xxx-xxx-xxx
✅ 成功迁移 5/5 个流到媒体源

[3/5] 迁移 Streams → MediaAssets
找到 5 个流
  ✓ 迁移流: Stream1 → 资产 ID: xxx-xxx-xxx
  ✓ 迁移流: Stream2 → 资产 ID: xxx-xxx-xxx
✅ 成功迁移 5/5 个流到媒体资产

[4/5] 迁移 Algorithms → Operators
找到 3 个算法
  ✓ 迁移算法: ObjectDetection → 算子 ID: xxx-xxx-xxx
  ✓ 迁移算法: FaceRecognition → 算子 ID: xxx-xxx-xxx
✅ 成功迁移 3/3 个算法

[5/5] 清理废弃表
  删除表: algorithm_bindings
  ✓ 已删除: algorithm_bindings
  删除表: inference_results
  ✓ 已删除: inference_results
  删除表: streams
  ✓ 已删除: streams
✅ 清理完成

✅ 迁移完成！
```

## 迁移后验证

```bash
# 启动服务
./bin/goyavision

# 检查媒体源
curl http://localhost:8080/api/v1/sources

# 检查媒体资产
curl http://localhost:8080/api/v1/assets

# 检查算子
curl http://localhost:8080/api/v1/operators

# 检查菜单（需要认证）
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/menus/tree
```

## 回滚

如果迁移出现问题，使用备份恢复：

```bash
psql goyavision < backup_YYYYMMDD_HHMMSS.sql
```
