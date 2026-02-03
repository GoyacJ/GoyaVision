# GoyaVision V1.0 需求文档

## 1. 项目定义

**GoyaVision** 是一个基于 AI 的智能媒体处理平台，提供媒体资产管理、算子编排、工作流自动化能力。

### 1.1 产品定位

- **智能媒体处理平台**：支持视频、图片、音频的 AI 分析与处理
- **算子编排引擎**：可插拔的算子体系，支持自定义算子
- **工作流自动化**：通过 DAG 工作流编排算子，实现复杂业务场景
- **多媒体资产管理**：统一管理媒体资产的全生命周期

### 1.2 产品范围

| 模块 | 功能 |
|------|------|
| **资产库** | 媒体源管理、资产管理、录制管理、存储配置 |
| **算子中心** | 算子市场、算子配置、算子监控 |
| **任务中心** | 工作流编排、任务调度、产物管理 |
| **控制台** | 用户管理、权限管理、系统监控 |

## 2. 核心概念

### 2.1 MediaAsset（媒体资产）

媒体资产是系统中所有可被追踪、处理、复用的媒体实体。

**属性**：
- `id`：唯一标识
- `type`：资产类型（video、image、audio）
- `source_type`：来源类型（live、vod、upload、generated）
- `source_id`：关联的媒体源 ID
- `parent_id`：派生自哪个资产（可选）
- `name`：资产名称
- `path`：存储路径
- `duration`：时长（秒，视频/音频）
- `size`：文件大小（字节）
- `format`：格式（mp4、jpg、wav 等）
- `metadata`：扩展元数据（分辨率、帧率、编码等）
- `status`：状态（pending、ready、failed）
- `tags`：标签（分类、场景等）
- `created_at`：创建时间
- `updated_at`：更新时间

**资产来源**：
- **live**：实时流录制或抽帧
- **vod**：点播视频
- **upload**：用户上传
- **generated**：算子生成（剪辑、转码、生成）

### 2.2 MediaSource（媒体源）

媒体源是媒体资产的来源，包括视频流、文件上传等。

**属性**：
- `id`：唯一标识
- `name`：源名称
- `type`：源类型（pull、push、upload）
- `url`：源地址（拉流地址，可选）
- `protocol`：协议类型（rtsp、rtmp、hls、webrtc、file）
- `enabled`：是否启用
- `status`：实时状态（ready、online、offline）
- `metadata`：扩展元数据
- `created_at`：创建时间
- `updated_at`：更新时间

**源类型**：
- **pull**：拉流（从外部地址拉取）
- **push**：推流（等待外部推送）
- **upload**：文件上传

### 2.3 Operator（算子）

算子是 AI/媒体处理的能力单元，具有标准化的输入输出协议。

**属性**：
- `id`：唯一标识
- `code`：算子编码（唯一）
- `name`：算子名称
- `category`：算子分类（analyze、edit、generate、transform）
- `version`：版本号
- `input_spec`：输入规格（支持的媒体类型、参数）
- `output_spec`：输出规格（产物类型、结构）
- `endpoint`：服务端点（HTTP）
- `config`：默认配置
- `status`：状态（enabled、disabled）
- `is_builtin`：是否内置（内置 vs 自定义）
- `description`：描述
- `icon`：图标
- `created_at`：创建时间
- `updated_at`：更新时间

**算子分类**：
- **analyze**（分析）：检测、识别、分类、追踪、OCR、ASR
- **edit**（编辑）：剪辑、裁剪、打码、去水印、字幕、水印
- **generate**（生成）：TTS、配音、扩展、摘要、高光
- **transform**（转换）：转码、压缩、分辨率调整、格式转换、增强

**标准化 I/O 协议**：

输入：
```json
{
  "asset_id": "资产 ID",
  "params": {
    "key": "value"
  }
}
```

输出：
```json
{
  "output_assets": [
    {
      "type": "video|image|audio",
      "path": "存储路径",
      "format": "格式",
      "metadata": {}
    }
  ],
  "results": [
    {
      "type": "detection|classification|ocr|...",
      "data": {},
      "confidence": 0.95
    }
  ],
  "timeline": [
    {
      "start": 0.0,
      "end": 5.0,
      "event_type": "事件类型",
      "confidence": 0.95,
      "data": {}
    }
  ],
  "diagnostics": {
    "latency_ms": 150,
    "model_version": "v1.0",
    "device": "gpu"
  }
}
```

### 2.4 Workflow（工作流）

工作流是算子的编排组合，通过 DAG（有向无环图）定义算子的执行顺序和数据流转。

**属性**：
- `id`：唯一标识
- `name`：工作流名称
- `description`：描述
- `trigger`：触发配置
- `nodes`：DAG 节点列表
- `edges`：节点连接关系
- `status`：状态（draft、active、paused）
- `created_at`：创建时间
- `updated_at`：更新时间

**触发器类型**：
- **manual**：手动触发
- **schedule**：定时触发（cron 表达式）
- **event**：事件触发（新资产、录制完成、流上线）

**DAG 节点**：
```json
{
  "id": "节点 ID",
  "operator_id": "算子 ID",
  "params": {
    "key": "value"
  },
  "retry": 3,
  "timeout": 300
}
```

**DAG 边**：
```json
{
  "from": "源节点 ID",
  "to": "目标节点 ID",
  "condition": "条件表达式（可选）"
}
```

### 2.5 Task（任务）

任务是工作流的一次执行实例。

**属性**：
- `id`：唯一标识
- `workflow_id`：工作流 ID
- `trigger_type`：触发类型（manual、schedule、event）
- `input_assets`：输入资产 ID 列表
- `status`：状态（pending、running、completed、failed、cancelled）
- `progress`：进度（0-100）
- `current_node`：当前执行节点
- `started_at`：开始时间
- `completed_at`：完成时间
- `error`：错误信息（失败时）
- `created_at`：创建时间

### 2.6 Artifact（产物）

产物是算子/工作流的输出结果，包括新生成的媒体资产、结构化数据、时间轴片段等。

**属性**：
- `id`：唯一标识
- `task_id`：任务 ID
- `node_id`：节点 ID
- `operator_id`：算子 ID
- `type`：产物类型（asset、result、timeline、diagnostic）
- `asset_id`：关联的资产 ID（如果是资产类型）
- `data`：产物数据（JSON）
- `created_at`：创建时间

**产物类型**：
- **asset**：新生成的媒体资产
- **result**：结构化结果（检测框、分类标签、OCR 文本）
- **timeline**：时间轴片段（事件、高光、镜头切分）
- **diagnostic**：诊断信息（性能指标、模型版本）

## 3. 功能需求

### 3.1 资产库（Asset Library）

#### 3.1.1 媒体源管理

- **CRUD**：创建、查询、更新、删除媒体源
- **类型支持**：拉流（pull）、推流（push）、上传（upload）
- **协议支持**：RTSP、RTMP、HLS、WebRTC、SRT、File
- **状态查询**：实时状态（ready、online、offline、readers、bitrate）
- **多协议分发**：单一源自动转换为多协议
- **启停控制**：启用/禁用媒体源

#### 3.1.2 媒体资产管理

- **CRUD**：创建、查询、更新、删除媒体资产
- **类型支持**：视频（video）、图片（image）、音频（audio）、流媒体（stream）
- **来源追踪**：记录资产来源（source_id、parent_id）
- **元数据管理**：分辨率、帧率、编码、时长、大小
- **标签系统**：支持多标签分类
- **搜索过滤**：按类型、来源、标签、时间范围搜索
- **预览**：支持多协议预览（HLS、RTSP、RTMP、WebRTC）
- **下载导出**：支持资产下载

**添加资产 - 流媒体接入**（设计，结合 MediaMTX，详见 `docs/stream-asset-mediamtx-design.md`、`docs/asset-stream-ingestion.md`）：

- **原则**：本平台与 MediaMTX 深度集成，作为 MediaMTX 客户端；**创建流媒体资产时必须接入 MediaMTX**，不提供“仅登记流地址、不接入 MediaMTX”的模式。**建 MediaSource 表**，与 MediaMTX path 一一对应；流媒体资产通过 `source_id` 关联 MediaSource。
- **方式一：新建流并创建资产**（主入口）
  - 用户输入流地址、资产名称、标签；系统生成 path name，调用 MediaMTX AddPath(source=流地址)，创建 MediaSource(path_name, url, type=pull, …)，再创建 MediaAsset(type=stream, source_type=live, source_id=MediaSource.ID, path=path_name)。具备完整预览、状态、录制、点播能力。
- **方式二：从已有媒体源创建资产**
  - 用户选择已有 MediaSource、填写资产名称、标签；系统创建 MediaAsset(type=stream, source_type=live, source_id=MediaSource.ID, path=MediaSource.path_name)。同一路流可对应多个资产（不同标签/用途）。

#### 3.1.3 录制管理

- **录制控制**：启动/停止录制
- **录制格式**：fMP4（推荐）、MPEG-TS
- **分段录制**：可配置段时长（默认 1 小时）
- **录制索引**：自动索引录制段
- **点播播放**：支持 HLS/MP4 点播，支持 Seek 跳转
- **录制列表**：查询录制文件列表
- **录制状态**：实时录制状态查询

#### 3.1.4 存储配置

- **存储路径**：配置录制、抽帧、上传的存储路径
- **生命周期**：设置资产自动清理策略
- **存储限额**：设置存储空间限额和告警

### 3.2 算子中心（Operator Hub）

#### 3.2.1 算子市场

- **算子列表**：展示所有可用算子
- **算子分类**：按分类（分析、编辑、生成、转换）筛选
- **算子搜索**：按名称、分类、标签搜索
- **算子详情**：查看算子详细信息、输入输出规格、示例
- **算子启用**：启用/禁用算子

#### 3.2.2 内置算子

**分析类（Analyze）**：
- **FrameExtract**：抽帧
- **ObjectDetection**：目标检测
- **FaceRecognition**：人脸识别
- **OCR**：文字识别
- **ObjectTracking**：目标追踪
- **ImageClassification**：图像分类
- **SceneDetection**：场景检测
- **ASR**：语音转文字

**编辑类（Edit）**：
- **VideoClip**：视频剪辑
- **VideoCrop**：视频裁剪
- **VideoMosaic**：视频打码
- **RemoveWatermark**：去水印
- **AddWatermark**：加水印
- **GenerateSubtitle**：生成字幕
- **AudioMix**：音频混合

**生成类（Generate）**：
- **TTS**：文本转语音
- **VideoHighlight**：高光摘要
- **VideoExtend**：视频扩展
- **ImageGenerate**：图像生成

**转换类（Transform）**：
- **VideoTranscode**：视频转码
- **VideoCompress**：视频压缩
- **VideoResize**：分辨率调整
- **AudioTranscode**：音频转码
- **ImageResize**：图片缩放
- **VideoEnhance**：视频增强

#### 3.2.3 自定义算子

- **算子上传**：上传自定义算子（Docker 镜像）
- **算子配置**：配置算子端点、输入输出规格、默认参数
- **算子测试**：测试算子功能
- **版本管理**：支持算子多版本
- **算子分享**：分享算子给其他用户（规划）

#### 3.2.4 算子监控

- **调用统计**：调用次数、成功率、失败率
- **性能指标**：平均延迟、P95/P99 延迟
- **资源使用**：CPU、内存、GPU 使用率
- **成本分析**：按算子统计计算成本

### 3.3 任务中心（Task Center）

#### 3.3.1 工作流编排

- **可视化编排**：拖拽式 DAG 工作流设计器
- **节点配置**：配置算子参数、重试策略、超时
- **条件分支**：支持 if/else、switch 条件路由
- **并行执行**：支持并行节点
- **数据传递**：节点间数据流转
- **工作流模板**：预定义常用工作流模板
- **工作流版本**：支持工作流版本管理

#### 3.3.2 任务调度

- **触发器配置**：
  - 手动触发：手动执行工作流
  - 定时触发：cron 表达式定时执行
  - 事件触发：新资产、录制完成、流上线等事件
- **调度策略**：
  - 固定间隔：每 N 秒执行一次
  - 定时调度：指定时间段内执行
  - 首次延迟：支持延迟首次执行
- **并发控制**：限制同时执行的任务数
- **优先级**：支持任务优先级

#### 3.3.3 任务管理

- **任务列表**：查看所有任务
- **任务详情**：查看任务执行详情、日志、产物
- **任务控制**：启动、停止、重试任务
- **任务搜索**：按状态、工作流、时间范围搜索
- **任务监控**：实时查看任务执行进度

#### 3.3.4 产物管理

- **产物列表**：查看所有产物
- **产物详情**：查看产物详细信息
- **产物下载**：下载产物（资产、结果文件）
- **产物搜索**：按任务、算子、类型搜索
- **产物关联**：查看产物关联的资产、任务

### 3.4 控制台（Console）

#### 3.4.1 仪表盘

- **系统概览**：资产数、任务数、算子数、存储使用
- **实时监控**：正在运行的任务、流状态、系统资源
- **统计报表**：任务执行趋势、算子使用排行、成本分析

#### 3.4.2 用户管理

- **CRUD**：创建、查询、更新、删除用户
- **角色分配**：为用户分配角色
- **状态控制**：启用/禁用用户
- **密码管理**：修改密码、重置密码

#### 3.4.3 角色管理

- **CRUD**：创建、查询、更新、删除角色
- **权限分配**：为角色分配权限
- **菜单分配**：为角色分配菜单
- **超级管理员**：super_admin 角色拥有所有权限

#### 3.4.4 权限管理

- **权限列表**：查看所有权限
- **RBAC 模型**：用户-角色-权限三级授权
- **按钮级权限**：支持前端按钮级权限控制

#### 3.4.5 系统配置

- **配置管理**：系统参数配置
- **集成配置**：MediaMTX、AI 服务等外部系统集成

#### 3.4.6 审计日志

- **操作日志**：记录用户操作
- **系统日志**：记录系统事件
- **日志查询**：按用户、操作、时间查询

## 4. 非功能需求

### 4.1 性能

- **并发处理**：支持 10+ 路流同时处理
- **任务并发**：支持 100+ 任务并发执行
- **响应时间**：API 响应时间 < 500ms
- **吞吐量**：支持 1000+ 次/秒 API 调用

### 4.2 可扩展性

- **算子插件化**：支持自定义算子
- **工作流灵活性**：支持任意复杂度工作流
- **多租户支持**：架构上支持多租户（tenant_id）
- **水平扩展**：支持服务水平扩展

### 4.3 可靠性

- **高可用**：支持集群部署
- **容错**：支持任务失败重试
- **数据持久化**：所有关键数据持久化到数据库
- **优雅关闭**：支持优雅关闭，确保任务不丢失

### 4.4 安全

- **认证**：JWT Token 认证
- **授权**：RBAC 权限模型
- **密码加密**：bcrypt 加密
- **输入验证**：所有输入进行验证
- **审计日志**：记录所有操作

### 4.5 可维护性

- **分层架构**：清晰的分层架构
- **日志**：结构化日志
- **监控**：Prometheus 指标暴露
- **健康检查**：/health、/ready 端点
- **文档**：完善的 API 文档

### 4.6 部署

- **容器化**：Docker 镜像
- **编排**：Docker Compose / Kubernetes
- **配置管理**：环境变量覆盖
- **一键部署**：支持一键启动

## 5. 数据模型

### 5.1 核心实体

| 实体 | 说明 |
|------|------|
| MediaSource | 媒体源 |
| MediaAsset | 媒体资产 |
| Operator | 算子 |
| Workflow | 工作流 |
| Task | 任务 |
| Artifact | 产物 |
| User | 用户 |
| Role | 角色 |
| Permission | 权限 |
| Menu | 菜单 |

### 5.2 关系

```
MediaSource 1:N MediaAsset
MediaAsset 1:N MediaAsset (parent-child)
Workflow 1:N Task
Task 1:N Artifact
Artifact N:1 MediaAsset (可选)

User N:M Role
Role N:M Permission
Role N:M Menu
```

## 6. API 设计

### 6.1 API 原则

- RESTful 风格
- 统一 `/api/v1` 前缀
- JWT 认证
- 统一错误响应格式
- 支持分页、过滤、排序

### 6.2 API 端点

| 模块 | 端点前缀 |
|------|----------|
| 认证 | `/api/v1/auth` |
| 用户 | `/api/v1/users` |
| 角色 | `/api/v1/roles` |
| 权限 | `/api/v1/permissions` |
| 菜单 | `/api/v1/menus` |
| 媒体源 | `/api/v1/sources` |
| 媒体资产 | `/api/v1/assets` |
| 录制 | `/api/v1/sources/:id/record` |
| 算子 | `/api/v1/operators` |
| 工作流 | `/api/v1/workflows` |
| 任务 | `/api/v1/tasks` |
| 产物 | `/api/v1/artifacts` |

详细 API 设计见 [API 文档](api.md)。

## 7. 技术栈

### 7.1 后端

- **语言**：Go 1.22+
- **框架**：Echo v4
- **ORM**：GORM
- **配置**：Viper
- **调度**：gocron/v2
- **认证**：golang-jwt/jwt/v5
- **数据库**：PostgreSQL 12+

### 7.2 流媒体

- **流媒体服务器**：MediaMTX
- **协议**：RTSP、RTMP、HLS、WebRTC、SRT
- **视频处理**：FFmpeg

### 7.3 前端

- **框架**：Vue 3
- **语言**：TypeScript
- **构建**：Vite
- **UI**：Element Plus
- **播放器**：video.js
- **状态管理**：Pinia
- **路由**：Vue Router

### 7.4 部署

- **容器化**：Docker
- **编排**：Docker Compose / Kubernetes
- **监控**：Prometheus + Grafana（规划）

## 8. 开发路线

### 8.1 Phase 1：核心闭环（当前 V1.0）

- [x] 媒体源管理（拉流、推流）
- [x] 实时流录制与点播
- [x] 内置算子（抽帧、目标检测）
- [ ] 媒体资产管理
- [ ] 简化工作流（单算子任务）
- [ ] 任务调度与执行
- [ ] 产物管理

### 8.2 Phase 2：能力扩展

- [ ] 多媒体类型（图片、音频）
- [ ] 更多内置算子（编辑、生成、转换）
- [ ] 复杂工作流（DAG 编排）
- [ ] 可视化工作流设计器
- [ ] 工作流模板市场

### 8.3 Phase 3：平台化

- [ ] 自定义算子（Docker 镜像）
- [ ] 算子市场（第三方算子）
- [ ] 多租户支持
- [ ] 开放 API 与 SDK
- [ ] 监控与告警

## 9. 变更记录

| 日期 | 版本 | 变更内容 |
|------|------|----------|
| 2025-02 | V1.0 | 架构重构：引入 MediaAsset、Operator、Workflow、Task、Artifact 核心概念；废弃 AlgorithmBinding；模块重命名为资产库、算子中心、任务中心、控制台 |

---

**注意**：本文档会随着项目演进持续更新。
