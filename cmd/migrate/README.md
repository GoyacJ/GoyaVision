# GoyaVision 数据迁移工具

## 概述

此工具用于将 GoyaVision 从旧架构迁移到 V1.0 新架构。

## 迁移内容

### 1. Streams → MediaAssets
- 将所有视频流转换为媒体资产
- 保留源 ID 关联
- 状态映射：`enabled=true` → `ready`, `enabled=false` → `pending`

### 2. Algorithms → Operators
- 将所有算法转换为算子
- 分类统一设为 `analysis`
- 类型添加 `legacy_` 前缀
- 状态设为 `published`

### 3. 清理旧表
- 删除 `algorithm_bindings` 表
- 删除 `inference_results` 表

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
1. 迁移 streams → media_assets（作为媒体源）
2. 迁移 algorithms → operators
3. 清理废弃表（algorithm_bindings、inference_results）

是否继续？ [y/N]: y

开始迁移...

[1/3] 迁移 Streams → MediaAssets
找到 5 个流
  ✓ 迁移流: Stream1 → 资产 ID: xxx-xxx-xxx
  ✓ 迁移流: Stream2 → 资产 ID: xxx-xxx-xxx
✅ 成功迁移 5/5 个流

[2/3] 迁移 Algorithms → Operators
找到 3 个算法
  ✓ 迁移算法: ObjectDetection → 算子 ID: xxx-xxx-xxx
  ✓ 迁移算法: FaceRecognition → 算子 ID: xxx-xxx-xxx
✅ 成功迁移 3/3 个算法

[3/3] 清理废弃表
  删除表: algorithm_bindings
  ✓ 已删除: algorithm_bindings
  删除表: inference_results
  ✓ 已删除: inference_results
✅ 清理完成

✅ 迁移完成！
```

## 迁移后验证

```bash
# 启动服务
./bin/goyavision

# 检查媒体资产
curl http://localhost:8080/api/v1/assets

# 检查算子
curl http://localhost:8080/api/v1/operators
```

## 回滚

如果迁移出现问题，使用备份恢复：

```bash
psql goyavision < backup_YYYYMMDD_HHMMSS.sql
```
