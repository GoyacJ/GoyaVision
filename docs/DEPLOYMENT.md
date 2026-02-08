# GoyaVision 部署指南

## 系统要求

### 最低要求

- **操作系统**：Linux、macOS、Windows
- **Go**：1.22 或更高版本
- **PostgreSQL**：12 或更高版本
- **FFmpeg**：5.0 或更高版本
- **内存**：至少 2GB RAM
- **磁盘**：至少 10GB 可用空间（用于录制和抽帧文件）

### 推荐配置

- **CPU**：4 核或更多
- **内存**：8GB 或更多
- **磁盘**：SSD，至少 100GB（根据录制需求调整）
- **网络**：稳定的网络连接（用于 RTSP 流接入）

## 构建

### 完整构建（包含前端）

```bash
# 构建前端
cd web
npm install
npm run build

# 构建后端（会自动嵌入前端）
cd ..
go build -o bin/goyavision ./cmd/server
```

或使用 Makefile：

```bash
make build-all
```

### 仅构建后端

```bash
go build -o bin/goyavision ./cmd/server
```

### Docker 构建

```bash
docker-compose build
```
或
```bash
make docker-build
```

## 数据库准备

### PostgreSQL 安装

**Ubuntu/Debian**:
```bash
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib
```

**macOS**:
```bash
brew install postgresql
```

**CentOS/RHEL**:
```bash
sudo yum install postgresql-server postgresql-contrib
```

### 创建数据库和用户

```bash
# 创建数据库
createdb goyavision

# 或使用 psql
psql -c "CREATE DATABASE goyavision;"
psql -c "CREATE USER goyavision WITH PASSWORD 'goyavision';"
psql -c "GRANT ALL PRIVILEGES ON DATABASE goyavision TO goyavision;"
```

### 多数据库支持

除 PostgreSQL 外，支持 **MySQL** 和 **SQLite**，通过 `db.driver` 与 `db.dsn` 配置：

| 驱动 | 说明 | DSN 示例 |
|------|------|----------|
| `postgres` | 默认，生产推荐 | `host=localhost user=goyavision password=xxx dbname=goyavision port=5432 sslmode=disable` |
| `mysql` | 需 MySQL 5.7+（支持 JSON） | `user:pass@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=True` |
| `sqlite3` | 单文件，无需独立服务 | `file:./data/goyavision.db?mode=rwc` |

环境变量：`GOYAVISION_DB_DRIVER`、`GOYAVISION_DB_DSN`。未配置 `db.dsn` 时跳过数据库连接（如仅运行前端静态服务）。

### 多文件存储支持

文件与资产存储支持三种后端，通过 `storage.type` 选择：

| 类型 | 说明 | 配置 |
|------|------|------|
| `minio` | 默认，MinIO 或兼容 S3 的服务 | `minio.*`（endpoint、access_key、secret_key、bucket_name 等） |
| `s3` | AWS S3 或兼容端点 | `storage.s3.*`（region、bucket、endpoint、access_key、secret_key、public_base 等） |
| `local` | 本地目录 | `storage.local.base_path`（存储路径）、`storage.local.base_url`（对外访问 base URL） |

未配置或 `storage.type` 为空时默认为 `minio`，行为与仅配置 `minio` 段一致。使用 `local` 时需确保应用能读写 `base_path`，且 `base_url` 与对外提供文件访问的地址一致（如通过反向代理或静态路由提供下载）。

## 配置

### 配置结构

配置按环境隔离，使用 `GOYAVISION_ENV` 选择加载文件：

```
configs/
  ├── config.dev.yaml      # 开发环境默认
  ├── config.prod.yaml     # 生产环境（建议配合环境变量）
  ├── config.example.yaml  # 配置模板（可用于初始化）
  └── .env.example         # 环境变量示例
```

### 环境变量

通过环境变量覆盖（`GOYAVISION_` 前缀）：

```bash
export GOYAVISION_ENV=dev
export GOYAVISION_DB_DSN="host=localhost user=goyavision password=goyavision dbname=goyavision port=5432 sslmode=disable"
export GOYAVISION_JWT_SECRET="replace-with-secure-secret"
export GOYAVISION_MEDIAMTX_API_ADDRESS="http://localhost:9997"
export GOYAVISION_MINIO_ENDPOINT="localhost:9000"
export GOYAVISION_MEDIAMTX_RECORD_PATH="./data/recordings/%path/%Y-%m-%d_%H-%M-%S"
export GOYAVISION_MINIO_USE_SSL=false
export GOYAVISION_MINIO_PUBLIC_BASE="https://vision.ysmjjsy.com/minio"
```

生产环境建议使用 `.env`（参考 `configs/.env.example`），并在配置文件中引用环境变量占位符（如 `config.prod.yaml`）。

> 说明：启动时会优先加载 `configs/.env`（最高优先级），其变量会覆盖系统已有环境变量与 `config.<env>.yaml` 的值。
> 环境变量键名支持 `GOYAVISION_SERVER_PORT` 这类下划线格式（对应 `server.port`）。

## 运行

### 开发模式

```bash
go run ./cmd/server
```

### 生产模式

```bash
./bin/goyavision
```

### Docker 运行

使用 Docker Compose 启动所有服务（包含数据库、MediaMTX、MinIO）：

```bash
docker-compose up -d
```

### 后台运行（Linux）

使用 systemd：

```ini
# /etc/systemd/system/goyavision.service
[Unit]
Description=GoyaVision AI Video Stream Analysis Platform
After=network.target postgresql.service

[Service]
Type=simple
User=goyavision
WorkingDirectory=/opt/goyavision
ExecStart=/opt/goyavision/bin/goyavision
Restart=always
RestartSec=5
Environment="GOYAVISION_DB_DSN=host=localhost user=goyavision password=xxx dbname=goyavision port=5432 sslmode=disable"

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable goyavision
sudo systemctl start goyavision
```

## 目录结构

确保以下目录存在且有写权限：

```
./data/
├── recordings/    # 录制文件
├── frames/        # 抽帧文件
└── uploads/       # 用户上传文件
```

## 验证部署

1. **检查服务**：
   ```bash
   curl http://localhost:8080/api/v1/sources
   ```

2. **访问 Web 界面**：
   打开浏览器访问 `http://localhost:8080`

3. **检查日志**：
   查看应用日志确认无错误

## 性能调优

### PostgreSQL

```sql
-- 调整连接池
ALTER SYSTEM SET max_connections = 200;
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
```

### FFmpeg 限制

根据服务器性能调整 `config.<env>.yaml` 中的限制：

```yaml
ffmpeg:
  max_record: 16  # 最大并发录制数
  max_frame: 16   # 最大并发抽帧数

preview:
  max_preview: 10  # 最大并发预览数
```

## 监控

### 健康检查（规划中）

- `/health`：服务健康状态
- `/ready`：服务就绪状态
- `/metrics`：Prometheus 指标

## 故障排查

### 常见问题

1. **数据库连接失败**
   - 检查 PostgreSQL 是否运行
   - 验证连接字符串是否正确
   - 检查防火墙设置

2. **FFmpeg 未找到**
   - 确认 FFmpeg 已安装并在 PATH 中
   - 或通过 `GOYAVISION_FFMPEG_BIN` 指定路径

3. **端口被占用**
   - 检查端口是否被其他服务占用
   - 修改 `server.port` 配置

4. **权限错误**
   - 确保数据目录有写权限
   - 检查文件系统权限

---

**注意**：本文档会随着部署方式演进持续更新。
