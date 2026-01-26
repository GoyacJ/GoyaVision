# GoyaVision

AI 视频流分析处理平台。支持 RTSP 流接入、抽帧、录制、AI 模型调用，以及按流配置多算法与定时、延迟等策略。

## 功能概览

- **视频流**：RTSP 接入，单实例 10+ 路
- **抽帧与推理**：按算法绑定的 `interval_sec`、`schedule`、`initial_delay_sec` 抽帧并调用 HTTP 推理服务，结果落库
- **录制**：与分析独立开关，按流启停，分段落盘
- **预览**：HLS 实时预览（MediaMTX 或 FFmpeg）
- **前端**：Vue 3 内嵌，流 / 算法绑定 / 录制 / 结果 CRUD 与 HLS 播放

## 技术栈

| 类别     | 选型                          |
|----------|-------------------------------|
| 后端     | Go 1.22, Echo v4, Viper, GORM |
| 数据库   | PostgreSQL                    |
| 调度     | gocron v2                     |
| 视频     | FFmpeg CLI, MediaMTX          |
| 前端     | Vue 3, TypeScript, Vite, Element Plus, video.js |

## 项目结构

```
goyavision/
├── cmd/server/          # 入口
├── config/              # 配置加载
├── configs/             # config.yaml
├── internal/
│   ├── domain/          # 领域实体
│   ├── port/             # 端口（Repository、Inference）
│   ├── app/              # 应用服务（编排）
│   ├── adapter/          # 基础设施（persistence、ffmpeg、preview、ai）
│   └── api/              # HTTP 路由、中间件、handler、dto
├── pkg/ffmpeg/           # FFmpeg 进程池
├── web/                  # Vue 前端（待建）
├── migrations/
└── docs/                 # 需求、开发进度
```

## 运行

### 依赖

- Go 1.22+
- PostgreSQL（库 `goyavision`，用户 `goyavision` / 密码 `goyavision` 或通过 `GOYAVISION_DB_DSN` 覆盖）
- FFmpeg（ PATH 或 `config.ffmpeg.bin`）
- 可选：MediaMTX（预览）

### 配置

- 主配置：`configs/config.yaml`
- 环境变量：`GOYAVISION_` 前缀覆盖，如 `GOYAVISION_DB_DSN`、`GOYAVISION_SERVER_PORT`

### 启动

```bash
# 拉依赖
go mod download

# 建库（需先有 PostgreSQL）
# createdb goyavision
# 或: psql -c "CREATE DATABASE goyavision;"

# 启动
go run ./cmd/server
```

默认 HTTP：`http://localhost:8080`。API 以 `/api/v1` 为前缀。

## 文档

- [需求文档](docs/requirements.md)
- [开发进度](docs/development-progress.md)
- 技术方案：`.cursor/plans/goyavision_技术实现方案_*.plan.md`（若存在）

## 开发与维护

- 规则与约定：`.cursor/rules/`
- 项目技能：`.cursor/skills/goyavision-context/`
- 开发时维护 `docs/development-progress.md` 中的阶段与状态。

## 许可证

未定。
