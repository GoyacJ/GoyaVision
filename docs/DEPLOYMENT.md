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
