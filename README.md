# GoyaVision

<div align="center">

**智能媒体处理平台**

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/GoyacJ/GoyaVision)](https://goreportcard.com/report/github.com/GoyacJ/GoyaVision)

[功能特性](#-功能特性) • [系统架构](#-系统架构) • [快速开始](#-快速开始) • [配置说明](#-配置说明) • [API 文档](#-api-文档) • [贡献](#-贡献)

</div>

---

**GoyaVision** 是一个基于 Go + Vue 3 的智能媒体处理平台，采用 Clean Architecture 架构设计。核心理念是**"业务=配置、能力=插件、执行=引擎"**。平台通过版本化算子、多执行模式（MCP/HTTP/CLI/AIModel）、多租户隔离和高性能 DAG 工作流引擎，实现复杂的 AI 媒体处理场景自动化。

## ✨ 功能特性

### 📦 资产库（Asset Library）

| 功能 | 说明 |
|------|------|
| **媒体源管理** | 深度集成 MediaMTX，支持拉流（RTSP/RTMP）、推流、实时状态监控 |
| **媒体资产** | 统一管理视频、图片、音频资产，支持租户可见性、多标签搜索与派生追踪 |
| **录制与点播** | 集成 MediaMTX 录制 API，支持 fMP4/MPEG-TS 分段录制、HLS/MP4 点播 |
| **存储配置** | 支持 MinIO、S3、本地文件系统三种后端（`storage.type`），集中管理存储与展示 URL |

### 🧩 算子中心（Operator Hub）

| 功能 | 说明 |
|------|------|
| **版本化管理** | 支持算子多版本共存，提供创建、激活（灰度）、回滚与归档全生命周期管理 |
| **多执行模式** | 统一路由 HTTP、CLI、MCP (Model Context Protocol) 与 AIModel 四种执行器 |
| **标准化协议** | 基于 JSON Schema 的输入输出强校验，确保节点间连接的 Schema 兼容性 |
| **模板市场** | 支持算子模板浏览、一键安装及从远程 MCP Server 自动同步工具模板 |
| **AI 模型集成** | 集中管理 OpenAI、Anthropic、Ollama 及国产大模型厂商连接 |

#### 内置算子分类

- **分析类（Analysis）**：抽帧、目标检测、人脸识别、OCR、目标追踪、场景检测、ASR
- **编辑类（Processing）**：视频剪辑、裁剪、打码、去水印、加水印、字幕生成、音频混合
- **生成类（Generation）**：TTS、高光摘要、视频扩展、图像生成
- **转换类（Transform）**：视频转码、压缩、分辨率调整、图片缩放、视频增强
- **实用类（Utility）**：格式转换、元数据提取、MCP 工具适配

### 🔄 任务中心（Task Center）

| 功能 | 说明 |
|------|------|
| **工作流编排** | 基于 Vue Flow 实现的可视化编辑器，支持拖拽编排、参数配置与自动布局 |
| **高性能引擎** | 支持 DAG 拓扑执行、并行节点处理、重试机制（指数退避）与超时控制 |
| **节点级追踪** | 详细记录每个节点的执行状态、耗时、输入输出与产生的产物（Artifact） |
| **实时监控** | 采用 SSE (Server-Sent Events) 实现任务详情页的实时进度推送与状态着色 |

### 🎛️ 控制台（Console）

| 功能 | 说明 |
|------|------|
| **多租户体系** | 全局资源支持租户级隔离（tenant_id）与细粒度的可见性策略配置 |
| **系统分类配置** | 在线管理存储（MinIO）、媒体服务器、网络及 AI 服务等模块参数 |
| **RBAC 增强** | 基于角色所有权的权限模型，支持基于条件的自动角色分配与动态菜单 |

## 🏗️ 系统架构

### 核心概念

```text
┌─────────────┐
│ MediaSource │  媒体源（MediaMTX 路径映射）
└──────┬──────┘
       │ 接入
       ▼
┌─────────────┐
│ MediaAsset  │  媒体资产（不可变实体，支持派生追踪）
└──────┬──────┘
       │ 输入
       ▼
┌─────────────┐      ┌───────────────────────────────────┐
│  Operator   │──────┤ OperatorVersion (HTTP/CLI/MCP...) │
└──────┬──────┘      └───────────────────────────────────┘
       │ 编排
       ▼
┌─────────────┐
│  Workflow   │  工作流（DAG 有向无环图）
└──────┬──────┘
       │ 执行
       ▼
┌─────────────┐      ┌───────────────────────────────────┐
│    Task     │──────┤ NodeExecution (节点级状态与日志)   │
└──────┬──────┘      └───────────────────────────────────┘
       │ 产出
       ▼
┌─────────────┐
│  Artifact   │  产物（结构化结果、新资产、时间轴、报告）
└─────────────┘
```

### 分层架构

GoyaVision 遵循 **Clean Architecture** 规范，并采用 **CQRS** 模式分离业务读写逻辑：

- **API Layer**: RESTful 路由、DTO 映射、SSE 进度推送。
- **App Layer**: 命令与查询处理器（Command/Query Handlers）、WorkflowScheduler。
- **Port Layer**: 定义仓储、执行器与校验器接口。
- **Domain Layer**: 纯粹的业务实体与核心规则，零外部库依赖。
- **Adapter Layer**: GORM (PostgreSQL)、MediaMTX、FFmpeg 及各模式执行器实现。

## 🚀 快速开始

### 方式一：Docker Compose（推荐）

```bash
# 克隆仓库
git clone https://github.com/GoyacJ/GoyaVision.git
cd GoyaVision

# 一键启动
docker-compose up -d
```

服务默认地址：
- **Web 界面**: http://localhost:8080
- **API 基址**: http://localhost:8080/api/v1
- **默认账号**: `admin` / `admin123`

### 方式二：本地开发

1. **环境准备**: Go 1.22+, Node.js 20+，数据库（PostgreSQL/MySQL/SQLite 任选），FFmpeg, MediaMTX。
2. **初始化配置**: 复制 `configs/config.example.yaml` 为 `configs/config.dev.yaml`。
3. **启动后端**: `go build -o bin/server ./cmd/server && ./bin/server`
4. **启动前端**: `cd web && pnpm install && pnpm dev`

## ⚙️ 配置说明

系统支持按环境加载 YAML 配置文件，并可使用环境变量进行覆盖（前缀 `GOYAVISION_`）。

核心配置项包括：
- `db`: 数据库驱动（`db.driver`: postgres/mysql/sqlite3）与连接串（`db.dsn`），详见 [部署指南](docs/DEPLOYMENT.md#多数据库支持)。
- `storage`: 文件存储类型（`storage.type`: minio/s3/local）及对应 `minio`/`storage.s3`/`storage.local` 配置，详见 [部署指南](docs/DEPLOYMENT.md#多文件存储支持)。
- `mediamtx`: API 与各协议访问基址。
- `jwt`: 认证密钥与过期策略。
- `mcp`: 远程 MCP Server 注册清单。

## 📖 文档

- [文档总览](docs/README.md) - 文档索引与维护约定
- [需求文档](docs/requirements.md) - V1.0 核心业务定义
- [架构文档](docs/architecture.md) - 技术细节与依赖规则
- [API 文档](docs/api.md) - 完整端点与 DTO 说明
- [开发进度](docs/development-progress.md) - 里程碑与迭代记录

---

<div align="center">

**GoyaVision V1.0** - 资产驱动的智能媒体编排引擎

**业务 = 配置 · 能力 = 插件 · 执行 = 引擎**

[⭐ Star us on GitHub](https://github.com/GoyacJ/GoyaVision) • [📖 阅读文档](docs/README.md) • [🤝 参与贡献](CONTRIBUTING.md)

</div>
