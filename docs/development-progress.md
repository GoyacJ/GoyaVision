# GoyaVision 开发进度

## 阶段与状态

| 阶段 | 状态 | 说明 |
|------|------|------|
| 1. 项目初始化、骨架搭建 | 已完成 | 目录、go.mod、cmd/server、config、domain/port/adapter/app/api 骨架、README、需求、本文档 |
| 2. 基础与持久化 | 待开始 | GORM 迁表、Repository 实现与接入、Stream/Algorithm/AlgorithmBinding CRUD |
| 3. FFmpeg 与池 | 待开始 | FFmpegManager（录制 segment、抽帧 fps+image2）、进程池与限流 |
| 4. 录制 | 待开始 | RecordService、RecordSession 启停；/record/start、/record/stop、/record/sessions |
| 5. 抽帧与推理 | 待开始 | gocron Job（interval_sec、schedule、initial_delay_sec）；取帧→AI→InferenceResult；/inference_results |
| 6. 预览 | 待开始 | PreviewManager（MediaMTX/FFmpeg HLS）；/preview/start、/preview/stop；/live 或代理 |
| 7. 前端 | 待开始 | Vue 3、流/算法绑定/录制/结果 CRUD、video.js 预览、构建与 embed |
| 8. 联调与优化 | 待开始 | 10+ 路压测、FFmpeg/预览上限、DB 索引与分页 |

---

## 1. 项目初始化、骨架搭建（已完成）

### 完成项

- [x] 项目目录按方案搭建：`cmd/server`、`config`、`configs`、`internal/domain|port|app|adapter|api`、`pkg/ffmpeg`、`migrations`、`docs`
- [x] `go.mod`：Echo、Viper、GORM、PostgreSQL、gocron、uuid、gorm.io/datatypes
- [x] `config`：Viper 加载 `configs/config.yaml` 与 `GOYAVISION_*` 环境变量
- [x] `cmd/server/main.go`：加载配置、GORM 连接与 AutoMigrate、Echo、Router、优雅退出
- [x] `internal/domain`：Stream、Algorithm、AlgorithmBinding、RecordSession、InferenceResult 实体
- [x] `internal/port`：Repository、Inference 接口
- [x] `internal/adapter/persistence`：Repository 实现与 AutoMigrate
- [x] `internal/api`：Deps、Router、Middleware、Handler 占位（stream、algorithm、algorithm_binding、record、inference、preview）、dto（stream）
- [x] `internal/app`：StreamService、RecordService、InferenceService、PreviewService 占位
- [x] `pkg/ffmpeg/pool.go`：Pool 占位
- [x] `README.md`：功能概览、技术栈、结构、运行方式、文档索引
- [x] `docs/requirements.md`：需求文档
- [x] `docs/development-progress.md`：本文档

### 待后续

- 健康检查：`/health`、`/ready`
- 指标：`/metrics`（Prometheus）

---

## 2. 基础与持久化（待开始）

- [ ] Handler 中接入 app/Repo，实现 Stream、Algorithm、AlgorithmBinding 的完整 CRUD 与校验
- [ ] 完善 dto 与错误处理
- [ ] 索引与约束按方案补充（如 RecordSession 唯一 running、InferenceResult 查询索引）

---

## 3–8. 后续阶段

- 在对应阶段开发时，将上表中「待开始」改为「进行中 / 已完成」，并在下方增补该阶段的「完成项」与「待后续」。

---

## 维护说明

- 每完成一个阶段或重要功能，更新本表及对应阶段的清单。
- 需求变更时，同步更新 `docs/requirements.md`。
