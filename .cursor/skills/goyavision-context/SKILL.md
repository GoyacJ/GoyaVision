---
name: goyavision-context
description: GoyaVision 项目结构、技术方案与文档索引。在实现或评审 GoyaVision 功能时使用，以便遵循既定分层、数据模型和 API 约定。
---

# GoyaVision 项目上下文

## 何时使用

- 在 GoyaVision 仓库中实现新功能、修改 handler/app/domain/adapter 时；
- 需要确认实体、API 路径、配置项或阶段划分时。

## 项目结构（核心）

```
cmd/server/      入口；config、GORM、Echo、Router、embed（后续）
config/          配置加载（Viper）
internal/
  domain/        Stream, Algorithm, AlgorithmBinding, RecordSession, InferenceResult
  port/          Repository, Inference
  app/           StreamService, RecordService, InferenceService, PreviewService
  adapter/       persistence（Repository）, ffmpeg, preview, ai（后续）
  api/           Deps, Router, Middleware, handler/*, dto/
pkg/ffmpeg/      进程池与限流
web/             Vue 3 前端（后续）
docs/            requirements.md, development-progress.md
```

## 数据与 API 要点

- **无 Task**：Stream → AlgorithmBinding → Algorithm；`AlgorithmBinding` 含 `interval_sec`、`schedule`、`initial_delay_sec`、`enabled`。
- **schedule**：`{start,end,days_of_week}` 或 null；`initial_delay_sec`：首次推理前延迟。
- **API 前缀**：`/api/v1`；流 `/streams`，算法 `/algorithms`，绑定 `/streams/:id/algorithm-bindings`，录制 `/streams/:id/record/*`，推理 `/inference_results`，预览 `/streams/:id/preview/start|stop`。

## 文档

- 需求：`docs/requirements.md`
- 进度与阶段：`docs/development-progress.md`
- 技术方案：`.cursor/plans/` 下 `goyavision_技术实现方案_*.plan.md`（若存在）
