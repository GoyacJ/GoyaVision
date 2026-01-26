# 变更日志

本文档记录项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [未发布]

### 计划中
- 健康检查端点（/health、/ready）
- Prometheus 指标（/metrics）
- 认证与鉴权
- Docker 支持
- 算法绑定管理页面完善
- 录制会话管理页面

## [0.2.0] - 2025-01-26

### 新增
- **前端界面**（阶段 7）
  - Vue 3 + TypeScript + Vite + Element Plus + video.js
  - 流列表页面（CRUD、预览、录制）
  - 算法管理页面
  - 推理结果查询页面
  - HLS 预览组件
  - Go embed 集成（单二进制部署）
- **预览功能**（阶段 6）
  - PreviewManager（MediaMTX/FFmpeg HLS）
  - 预览池限流
  - HLS 文件服务（/live）
- **抽帧与推理**（阶段 5）
  - Scheduler（gocron 调度器）
  - AI 推理适配器（HTTP + JSON）
  - 支持 interval_sec、schedule、initial_delay_sec
  - 推理结果查询（过滤、分页）
- **录制功能**（阶段 4）
  - RecordService（启停、会话管理）
  - 任务监控和自动状态更新
- **FFmpeg 与池**（阶段 3）
  - FFmpeg Pool（进程池与限流）
  - FFmpegManager（录制、单帧提取、连续抽帧）
- **基础与持久化**（阶段 2）
  - Stream、Algorithm、AlgorithmBinding 完整 CRUD
  - 统一错误处理机制
  - 数据库索引和约束

## [0.1.0] - 2025-01-26

### 新增
- 项目初始化和骨架搭建
- 分层架构设计（domain/port/app/adapter/api）
- 配置管理（Viper + YAML）
- 数据库模型定义（Stream, Algorithm, AlgorithmBinding, RecordSession, InferenceResult）
- HTTP API 路由框架（Echo）
- 项目文档（需求文档、开发进度、架构文档）

### 变更
- 项目从 Maas 重命名为 GoyaVision

---

## 版本说明

- **[未发布]**: 开发中，尚未发布的功能
- **[主版本.次版本.修订版本]**: 已发布的版本

### 变更类型

- **新增**: 新功能
- **变更**: 现有功能的变更
- **弃用**: 即将移除的功能
- **移除**: 已移除的功能
- **修复**: Bug 修复
- **安全**: 安全相关的修复
