---
name: goyavision-context
description: GoyaVision 项目结构、技术方案与文档索引。在实现或评审 GoyaVision 功能时使用，以便遵循既定分层、数据模型和 API 约定。
---

# GoyaVision 项目上下文

## 何时使用

- 在 GoyaVision 仓库中实现新功能、修改 handler/app/domain/adapter 时；
- 需要确认实体、API 路径、配置项或阶段划分时；
- 需要了解已实现的功能和代码结构时。

## 项目结构（核心）

```
cmd/server/          入口；config、GORM、Echo、Router、Scheduler、embed、初始化数据
config/              配置加载（Viper + YAML）
configs/             配置文件（config.yaml）
internal/
  domain/            Stream, Algorithm, AlgorithmBinding, RecordSession, InferenceResult,
                     User, Role, Permission, Menu
  port/              Repository, Inference
  app/               StreamService, AlgorithmService, AlgorithmBindingService,
                     RecordService, InferenceService, PreviewService, Scheduler,
                     AuthService, UserService, RoleService, MenuService
  adapter/
    persistence/     Repository 实现（GORM）、初始化数据（init_data.go）
    ai/              Inference 实现（HTTP 客户端）
  api/
    handler/         stream, algorithm, algorithm_binding, record, inference, preview,
                     auth, user, role, menu
    dto/             stream, algorithm, algorithm_binding, record, inference, preview,
                     auth, user, role, menu
    middleware/      auth.go（JWT 认证、权限校验）
    errors.go        统一错误处理
    static.go        前端静态文件服务（embed）
    router.go        路由注册（公开路由、认证路由、管理路由）
pkg/
  ffmpeg/            Pool（进程池）、Manager（录制和抽帧）
  preview/           Pool（预览池）、Manager（预览任务）
web/                 Vue 3 前端（src/, dist/）
  src/store/         Pinia 状态管理（用户、权限）
  src/views/login/   登录页面
  src/views/system/  系统管理页面（用户、角色、菜单）
  src/layout/        动态菜单布局
  src/directives/    权限指令（v-permission）
  src/router/guard.ts 路由守卫
docs/               需求、开发进度、架构文档、API 文档
```

## 已实现功能

### 核心功能（阶段 1-8 已完成）

1. **基础 CRUD**（阶段 2）
   - Stream、Algorithm、AlgorithmBinding 的完整 CRUD
   - 统一错误处理和 DTO 转换
   - 数据库索引和约束

2. **FFmpeg 与池**（阶段 3）
   - FFmpeg Pool（录制和抽帧限流）
   - FFmpegManager（录制、单帧提取、连续抽帧）

3. **录制功能**（阶段 4）
   - RecordService（启停、会话查询）
   - Record Handler（start、stop、sessions）
   - 任务监控和自动状态更新

4. **抽帧与推理**（阶段 5）
   - Scheduler（gocron 调度器）
   - AI 推理适配器（HTTP + JSON）
   - InferenceService（结果查询）
   - 支持 interval_sec、schedule、initial_delay_sec

5. **预览功能**（阶段 6）
   - PreviewManager（MediaMTX/FFmpeg）
   - Preview Pool（预览限流）
   - PreviewService（启停）
   - HLS 文件服务（/live）

6. **前端界面**（阶段 7）
   - Vue 3 + TypeScript + Element Plus
   - 流列表、算法管理、推理结果查询页面
   - HLS 预览组件（video.js）
   - Go embed 集成

7. **认证授权**（阶段 8）
   - RBAC 权限模型（User、Role、Permission、Menu）
   - JWT 认证（Access Token + Refresh Token）
   - 认证中间件和权限校验中间件
   - 登录/登出/Token 刷新/修改密码
   - 用户/角色/菜单管理 API
   - 前端：Pinia 状态管理、登录页面、路由守卫、权限指令、动态菜单

## 数据与 API 要点

- **无 Task**：Stream → AlgorithmBinding → Algorithm；`AlgorithmBinding` 含 `interval_sec`、`schedule`、`initial_delay_sec`、`enabled`。
- **schedule**：`{start,end,days_of_week}` JSON 格式；`start`/`end` 为时间字符串（"HH:MM:SS"），`days_of_week` 为星期数组（0-6）。
- **API 前缀**：`/api/v1`
  - **认证（公开）**：
    - `POST /auth/login`：登录
    - `POST /auth/refresh`：刷新 Token
  - **认证（需登录）**：
    - `GET /auth/profile`：获取当前用户信息
    - `PUT /auth/password`：修改密码
    - `POST /auth/logout`：登出
  - **用户管理**：`GET/POST/PUT/DELETE /users`、`POST /users/:id/reset-password`
  - **角色管理**：`GET/POST/PUT/DELETE /roles`
  - **菜单管理**：`GET/POST/PUT/DELETE /menus`、`GET /menus/tree`
  - **权限列表**：`GET /permissions`
  - **流**：`GET/POST/PUT/DELETE /streams`、`GET/POST/PUT/DELETE /streams/:id`
  - **算法**：`GET/POST/PUT/DELETE /algorithms`、`GET/POST/PUT/DELETE /algorithms/:id`
  - **绑定**：`GET/POST/PUT/DELETE /streams/:id/algorithm-bindings`、`GET/POST/PUT/DELETE /streams/:id/algorithm-bindings/:bid`
  - **录制**：`POST /streams/:id/record/start`、`POST /streams/:id/record/stop`、`GET /streams/:id/record/sessions`
  - **推理**：`GET /inference_results`（支持 stream_id、binding_id、from、to、limit、offset）
  - **预览**：`GET /streams/:id/preview/start`、`POST /streams/:id/preview/stop`
- **静态文件**：`/live/*`（HLS 文件）、`/*`（前端 SPA）
- **认证**：所有业务 API 需要在请求头中携带 `Authorization: Bearer <access_token>`

## 配置项

- **数据库**：`db.dsn`（PostgreSQL 连接字符串）
- **FFmpeg**：`ffmpeg.bin`（可执行文件路径）、`ffmpeg.max_record`、`ffmpeg.max_frame`
- **预览**：`preview.provider`（"mediamtx" 或 "ffmpeg"）、`preview.max_preview`、`preview.hls_base`
- **录制**：`record.base_path`、`record.segment_sec`
- **AI**：`ai.timeout`、`ai.retry`
- **JWT**：`jwt.secret`（签名密钥）、`jwt.expire`（Access Token 过期时间，默认 2h）、`jwt.refresh_exp`（Refresh Token 过期时间，默认 168h）、`jwt.issuer`

## 文档

- 需求：`docs/requirements.md`
- 进度与阶段：`docs/development-progress.md`
- 架构：`docs/architecture.md`
- 技术方案：`.cursor/plans/` 下 `goyavision_技术实现方案_*.plan.md`（若存在）

## 开发状态

- **阶段 1-8**：已完成
- **阶段 9**：联调与优化（待开始）

## 默认账号

- **用户名**：admin
- **密码**：admin123
- **角色**：超级管理员（拥有所有权限）
