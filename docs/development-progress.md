# GoyaVision 开发进度

## 阶段与状态

| 阶段 | 状态 | 说明 |
|------|------|------|
| 1. 项目初始化、骨架搭建 | 已完成 | 目录、go.mod、cmd/server、config、domain/port/adapter/app/api 骨架、README、需求、本文档 |
| 2. 基础与持久化 | 已完成 | GORM 迁表、Repository 实现与接入、Stream/Algorithm/AlgorithmBinding CRUD |
| 3. FFmpeg 与池 | 已完成 | FFmpegManager（录制 segment、抽帧 fps+image2）、进程池与限流 |
| 4. 录制 | 已完成 | RecordService、RecordSession 启停；/record/start、/record/stop、/record/sessions |
| 5. 抽帧与推理 | 已完成 | gocron Job（interval_sec、schedule、initial_delay_sec）；取帧→AI→InferenceResult；/inference_results |
| 6. 预览 | 已完成 | PreviewManager（MediaMTX/FFmpeg HLS）；/preview/start、/preview/stop；/live 或代理 |
| 7. 前端 | 已完成 | Vue 3、流/算法绑定/录制/结果 CRUD、video.js 预览、构建与 embed |
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

## 2. 基础与持久化（已完成）

### 完成项

- [x] 实现 StreamService（App 层）- Stream CRUD 业务逻辑
- [x] 实现 AlgorithmService（App 层）- Algorithm CRUD 业务逻辑
- [x] 实现 AlgorithmBindingService（App 层）- AlgorithmBinding CRUD 业务逻辑
- [x] 完善 Stream DTO（请求与响应 DTO）
- [x] 完善 Algorithm DTO（请求与响应 DTO）
- [x] 完善 AlgorithmBinding DTO（请求与响应 DTO）
- [x] 实现 Stream Handler 完整 CRUD 与错误处理
- [x] 实现 Algorithm Handler 完整 CRUD 与错误处理
- [x] 实现 AlgorithmBinding Handler 完整 CRUD 与错误处理
- [x] 统一错误处理机制（`internal/api/errors.go`）
- [x] 数据库索引与约束：
  - RecordSession 唯一约束（一个流只能有一个 running 状态的录制会话）
  - InferenceResult 查询索引（stream_id + ts, algorithm_binding_id + ts）

### 技术实现

- **错误处理**：统一的错误响应格式，区分业务错误（4xx）与基础设施错误（5xx）
- **输入验证**：在 Service 层进行业务规则验证
- **DTO 转换**：Domain 模型与 API DTO 分离，提供转换函数
- **数据库约束**：使用 PostgreSQL 部分唯一索引确保 RecordSession 唯一性

### 待后续

- 输入验证可以使用 validator 库增强
- 分页支持（当前 List 接口返回全部数据）

---

## 3. FFmpeg 与池（已完成）

### 完成项

- [x] 实现 FFmpeg Pool（进程池与限流）
  - 支持最大录制数限制（max_record）
  - 支持最大抽帧数限制（max_frame）
  - 使用互斥锁保证线程安全
  - 支持上下文取消自动释放资源
- [x] 实现 FFmpegManager
  - 录制功能：RTSP -> 分段 MP4（使用 `-c copy` 不重编码）
  - 单帧提取：RTSP -> 单张图片（用于推理）
  - 连续抽帧：RTSP -> 按间隔抽帧（fps + image2）
  - 进程生命周期管理（启动、停止、监控）
- [x] 进程管理
  - RecordTask：录制任务管理
  - FrameTask：抽帧任务管理
  - 支持优雅停止和强制终止

### 技术实现

- **进程池**：使用互斥锁和计数器实现资源限流
- **录制**：使用 FFmpeg segment muxer，支持分段落盘
- **抽帧**：支持单帧提取和连续抽帧两种模式
- **错误处理**：完善的错误处理和资源清理

### 待后续

- 进程健康检查与自动重启
- 录制文件完整性验证
- 抽帧任务支持 schedule 和 initial_delay_sec

---

## 4. 录制（已完成）

### 完成项

- [x] 实现 RecordService（App 层）
  - Start：启动录制，创建 RecordSession，启动 FFmpeg 录制任务
  - Stop：停止录制，更新 RecordSession 状态
  - ListSessions：列出流的录制会话历史
- [x] 集成 FFmpegManager
  - 使用 FFmpeg Pool 进行资源限流
  - 启动录制任务（RTSP -> 分段 MP4）
  - 任务生命周期管理
- [x] 实现 Record Handler
  - POST `/streams/:id/record/start`：启动录制
  - POST `/streams/:id/record/stop`：停止录制
  - GET `/streams/:id/record/sessions`：查询录制会话列表
- [x] 创建 Record DTO
  - RecordSessionResponse：录制会话响应
  - RecordStartResponse：启动录制响应
- [x] 录制任务管理
  - 内存中存储活跃任务
  - 任务监控（自动检测进程退出）
  - 线程安全的任务管理

### 技术实现

- **业务逻辑**：检查流状态、确保唯一运行中的录制会话
- **资源管理**：使用 FFmpeg Pool 限制并发录制数
- **任务监控**：后台监控任务状态，自动更新数据库
- **错误处理**：完善的错误处理和资源清理

### 待后续

- 录制文件完整性验证
- 录制文件大小和时长统计
- 支持录制质量配置

---

## 5. 抽帧与推理（已完成）

### 完成项

- [x] 实现 AI 推理适配器（adapter/ai）
  - HTTP + JSON 调用推理服务
  - 支持超时和重试机制
  - 错误处理和响应解析
- [x] 实现 Scheduler（调度器）
  - 使用 gocron 管理定时任务
  - 支持 interval_sec（固定间隔）
  - 支持 schedule（定时调度：start、end、days_of_week）
  - 支持 initial_delay_sec（首次延迟）
  - 自动加载启用的绑定并创建任务
- [x] 实现 InferenceService（App 层）
  - ListResults：查询推理结果（支持过滤、分页）
  - 支持按流、绑定、时间范围查询
- [x] 实现推理流程
  - 抽帧：使用 FFmpegManager 提取单帧
  - 编码：将图片编码为 base64
  - 调用：通过 Inference 接口调用 AI 服务
  - 落库：保存 InferenceResult 到数据库
- [x] 实现 Inference Handler
  - GET `/api/v1/inference_results`：查询推理结果列表
- [x] 创建 Inference DTO
  - InferenceResultListQuery：查询参数
  - InferenceResultResponse：响应格式
  - InferenceResultListResponse：列表响应（含总数）
- [x] 集成调度器到主程序
  - 启动时自动加载并调度所有启用的绑定
  - 优雅关闭时停止调度器

### 技术实现

- **调度策略**：
  - 固定间隔：按 `interval_sec` 定期执行
  - 定时调度：仅在 `schedule` 指定的时间段内执行
  - 首次延迟：支持 `initial_delay_sec` 延迟首次执行
- **抽帧**：使用 FFmpeg 提取单帧，保存为图片文件
- **AI 调用**：HTTP POST 请求，支持 base64 图片输入
- **结果持久化**：保存推理结果、延迟时间、帧引用等信息

### 待后续

- 支持更灵活的 schedule 配置（cron 表达式）
- 推理结果缓存机制
- 批量推理优化
- 推理失败重试策略优化

---

## 6. 预览（已完成）

### 完成项

- [x] 实现 PreviewManager
  - 支持 MediaMTX 和 FFmpeg 两种提供者
  - MediaMTX：使用 MediaMTX 服务（需外部运行）
  - FFmpeg：使用 FFmpeg 将 RTSP 转 HLS
  - HLS 输出管理（index.m3u8、segment.ts）
- [x] 实现预览池（Preview Pool）
  - 支持最大预览数限制（max_preview）
  - 使用互斥锁保证线程安全
  - 支持上下文取消自动释放资源
- [x] 实现 PreviewService（App 层）
  - Start：启动预览，返回 HLS URL
  - Stop：停止预览
  - 检查流状态和启用状态
- [x] 实现 Preview Handler
  - GET `/api/v1/streams/:id/preview/start`：启动预览
  - POST `/api/v1/streams/:id/preview/stop`：停止预览
- [x] 创建 Preview DTO
  - PreviewStartResponse：启动预览响应（包含 HLS URL）
- [x] 预览任务管理
  - 内存中存储活跃任务
  - 任务生命周期管理
  - 线程安全的任务管理

### 技术实现

- **提供者支持**：
  - MediaMTX：使用外部 MediaMTX 服务
  - FFmpeg：使用 FFmpeg 直接转换 RTSP 到 HLS
- **HLS 配置**：
  - 段时长：2 秒
  - 播放列表大小：3 段
  - 自动删除旧段
- **资源管理**：使用预览池限制并发预览数

### 待后续

- HLS 文件服务（/live 路由）
- MediaMTX 集成优化
- 预览质量配置
- 预览状态持久化

---

## 7. 前端（已完成）

### 完成项

- [x] 创建 Vue 3 项目基础结构
  - package.json：依赖管理（Vue 3、TypeScript、Vite、Element Plus、video.js）
  - vite.config.ts：Vite 构建配置
  - tsconfig.json：TypeScript 配置
- [x] 配置 Element Plus 和 video.js
  - Element Plus UI 组件库
  - video.js HLS 播放器
- [x] 实现流列表页面（StreamList）
  - 流的 CRUD 操作
  - 流状态显示（启用/禁用）
  - 预览功能（启动预览、HLS 播放）
  - 录制功能（启动录制）
  - 算法绑定入口
- [x] 实现算法管理页面（AlgorithmList）
  - 算法的 CRUD 操作
- [x] 实现推理结果查询页面（InferenceResultList）
  - 按流、绑定、时间范围查询
  - 分页支持
  - 推理输出查看
- [x] 实现 HLS 预览组件（HLSPreview）
  - 使用 video.js 播放 HLS 流
  - 支持动态 URL 切换
- [x] API 客户端封装
  - axios 封装
  - 统一错误处理
  - TypeScript 类型定义
- [x] 路由配置
  - Vue Router 配置
  - 页面路由（流列表、算法管理、推理结果）
- [x] Go embed 集成
  - 使用 embed.FS 嵌入前端构建产物
  - SPA 路由处理（所有非 API 路由返回 index.html）
  - 静态文件服务
- [x] 构建脚本
  - Makefile 支持（build-web、build-all）
  - 前端构建流程

### 技术实现

- **前端技术栈**：
  - Vue 3 Composition API
  - TypeScript
  - Vite 构建工具
  - Element Plus UI 组件
  - video.js HLS 播放
- **API 集成**：使用 axios 封装 API 调用，支持类型安全
- **SPA 路由**：Vue Router 管理前端路由
- **嵌入方式**：Go embed 将构建产物嵌入二进制文件

### 待后续

- 算法绑定管理页面（完整的 CRUD 和 schedule 配置）
- 录制会话管理页面
- 更完善的错误处理和用户提示
- 响应式设计优化
- 国际化支持

---

## 8. 联调与优化（待开始）

- 10+ 路压测
- FFmpeg/预览上限验证
- DB 索引与分页优化

---

## 维护说明

- 每完成一个阶段或重要功能，更新本表及对应阶段的清单。
- 需求变更时，同步更新 `docs/requirements.md`。
