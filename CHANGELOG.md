# 变更日志

本文档记录项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [未发布]

### 计划中
- 健康检查端点（/health、/ready）
- Prometheus 指标（/metrics）
- Docker 支持
- 算法绑定管理页面完善
- 录制会话管理页面
- 登录失败次数限制
- 审计日志

## [0.3.0] - 2025-01-26

### 新增
- **RBAC 认证授权**（阶段 8）
  - User/Role/Permission/Menu 领域实体
  - JWT 认证（Access Token + Refresh Token）
  - 认证中间件和权限校验中间件
  - 登录/登出/刷新 Token/修改密码 API
  - 用户管理 API（CRUD、角色分配、重置密码）
  - 角色管理 API（CRUD、权限分配、菜单分配）
  - 菜单管理 API（CRUD、树形结构）
  - 权限列表 API
  - 初始化数据（默认权限、菜单、超级管理员角色、admin 账号）
- **前端认证集成**
  - Pinia 状态管理（用户、Token、权限）
  - 登录页面
  - 路由守卫（未登录跳转登录页）
  - 权限指令（v-permission）
  - 动态菜单布局
  - 系统管理页面（用户、角色、菜单管理）

### 变更
- 所有业务 API 现在需要认证才能访问
- 前端布局改为动态菜单侧边栏
- 添加 @element-plus/icons-vue 依赖

### 依赖
- 新增 golang-jwt/jwt/v5
- 新增 golang.org/x/crypto（bcrypt）
- 新增 pinia、pinia-plugin-persistedstate

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
