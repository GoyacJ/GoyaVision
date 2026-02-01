# GoyaVision 架构文档

## 概述

GoyaVision 采用分层架构（Clean Architecture / Hexagonal Architecture），确保业务逻辑与基础设施解耦，提高可维护性和可测试性。

## 架构层次

### 1. Domain Layer（领域层）

**位置**: `internal/domain/`

**职责**:
- 定义核心业务实体
- 包含纯业务逻辑
- 无外部依赖（不导入 adapter、api、config 等）

**实体**:
- `Stream`: 视频流实体
- `Algorithm`: AI 算法配置
- `AlgorithmBinding`: 流与算法的绑定关系
- `RecordSession`: 录制会话
- `InferenceResult`: 推理结果
- `User`: 用户实体
- `Role`: 角色实体
- `Permission`: 权限实体
- `Menu`: 菜单实体

**原则**:
- 实体为纯 Go struct，包含业务字段
- 可包含领域方法（如验证逻辑）
- 不依赖任何基础设施

### 2. Port Layer（端口层）

**位置**: `internal/port/`

**职责**:
- 定义应用边界接口
- 定义领域需要的服务接口
- 可被 app、api 层依赖

**接口**:
- `Repository`: 数据持久化接口
- `Inference`: AI 推理服务接口

**原则**:
- 接口定义在 port 层
- 实现放在 adapter 层
- 通过依赖注入使用

### 3. App Layer（应用层）

**位置**: `internal/app/`

**职责**:
- 编排 domain 与 port
- 实现业务用例
- 不直接依赖 adapter 具体实现

**服务**:
- `StreamService`: 流管理用例（CRUD、状态管理）
- `AlgorithmService`: 算法管理用例（CRUD）
- `AlgorithmBindingService`: 算法绑定管理用例（CRUD、验证）
- `RecordService`: 录制用例（启停、会话管理、任务监控）
- `InferenceService`: 推理结果查询用例（过滤、分页）
- `PreviewService`: 预览用例（启停、HLS URL 管理）
- `Scheduler`: 调度器（gocron，管理推理任务）
- `AuthService`: 认证服务（登录、登出、Token 刷新、密码修改）
- `UserService`: 用户管理（CRUD、角色分配、密码重置）
- `RoleService`: 角色管理（CRUD、权限分配、菜单分配）
- `MenuService`: 菜单管理（CRUD、树形结构）

**原则**:
- 通过 port 接口操作，而非直接调用 adapter
- 一个服务对应一个业务用例
- 处理业务规则和流程编排

### 4. Adapter Layer（适配器层）

**位置**: `internal/adapter/`

**职责**:
- 实现 port 定义的接口
- 处理基础设施细节
- 可依赖 domain

**适配器**:
- `persistence`: 实现 `port.Repository`（GORM + PostgreSQL）
- `ai`: 实现 `port.Inference`（HTTP 客户端，支持超时和重试）

**原则**:
- 实现 port 接口
- 处理技术细节（数据库、HTTP、进程等）
- 将外部模型转换为 domain 模型

### 5. API Layer（接口层）

**位置**: `internal/api/`

**职责**:
- HTTP 路由定义
- 请求/响应处理
- DTO 转换
- 中间件

**组件**:
- `router.go`: 路由注册（认证路由、受保护路由、管理路由）
- `handler/`: 请求处理器（auth、user、role、menu、stream 等）
- `dto/`: 数据传输对象
- `middleware/auth.go`: 认证中间件（JWT 验证、权限校验）

**原则**:
- Handler 调用 app 服务或 port 接口
- 不直接操作数据库（通过 Repository）
- DTO 与 domain 模型分离

## 依赖关系

```
┌─────────────┐
│   API       │  HTTP 层
└──────┬──────┘
       │ 依赖
┌──────▼──────┐
│    App      │  应用服务层
└──────┬──────┘
       │ 依赖
┌──────▼──────┐      ┌─────────────┐
│   Port      │◄─────┤   Domain    │  接口与领域
└──────┬──────┘      └─────────────┘
       │ 实现
┌──────▼──────┐
│  Adapter    │  适配器层
└─────────────┘
```

### 依赖规则

1. **Domain** 不依赖任何其他层
2. **Port** 可依赖 Domain
3. **App** 可依赖 Domain 和 Port
4. **Adapter** 可依赖 Domain 和 Port（实现接口）
5. **API** 可依赖 App、Port、Domain

## 数据流

### 创建流的流程示例

```
1. HTTP Request → API Handler
2. Handler 转换 DTO → Domain 模型
3. Handler 调用 App Service
4. App Service 通过 Port.Repository 保存
5. Adapter.Persistence 实现 Repository，使用 GORM
6. 返回 Domain 模型
7. App Service 返回给 Handler
8. Handler 转换 Domain → DTO
9. HTTP Response
```

## 配置管理

**位置**: `config/`

- 使用 Viper 加载 YAML 配置
- 支持环境变量覆盖（`GOYAVISION_*` 前缀）
- 配置结构体定义在 `config.go`

## 进程管理

**位置**: `pkg/ffmpeg/`、`pkg/preview/`

- **FFmpeg Pool**：进程池与限流（最大录制数、最大抽帧数）
- **FFmpegManager**：录制任务（RTSP -> 分段 MP4）、抽帧任务（单帧提取、连续抽帧）
- **Preview Pool**：预览资源池（最大预览数）
- **PreviewManager**：预览任务管理（MediaMTX/FFmpeg HLS）
- 进程生命周期管理（启动、停止、监控）

## 调度器

**位置**: `internal/app/scheduler.go`

- 使用 gocron 管理定时任务
- 支持固定间隔（`interval_sec`）
- 支持定时调度（`schedule`：start、end、days_of_week）
- 支持首次延迟（`initial_delay_sec`）
- 启动时自动加载启用的算法绑定
- 任务管理（创建、删除、监控）

## 前端集成

**位置**: `web/`、`internal/api/static.go`

- Vue 3 + TypeScript + Vite + Element Plus + video.js
- 构建产物嵌入到 Go 二进制（`embed.FS`）
- SPA 路由处理（所有非 API 路由返回 index.html）
- HLS 文件服务（`/live/*`）

## 扩展点

### 添加新的基础设施适配器

1. 在 `port/` 定义接口
2. 在 `adapter/` 实现接口
3. 在 `app/` 或 `api/` 中注入使用

### 添加新的业务用例

1. 在 `app/` 创建新的 Service
2. 编排 domain 和 port
3. 在 `api/handler/` 中调用
4. 创建对应的 DTO

### 添加新的进程管理

1. 在 `pkg/` 创建新的管理器
2. 实现进程池和限流
3. 在对应的 Service 中集成

## 测试策略

- **单元测试**: Domain、App 层逻辑
- **集成测试**: Adapter 实现（数据库、HTTP 客户端）
- **端到端测试**: API 层完整流程

## 性能考虑

- 进程池限制并发数
- 数据库连接池
- 异步处理（推理、录制）
- 缓存策略（规划中）

## 安全考虑

- 输入验证（API 层和 Service 层）
- SQL 注入防护（使用 GORM 参数化查询）
- 错误处理不泄露敏感信息
- **JWT 认证**：Access Token + Refresh Token 双 Token 机制
- **RBAC 权限模型**：用户-角色-权限三级授权
- **密码加密**：使用 bcrypt 算法加密存储
- **权限中间件**：API 级别权限校验

## 已实现的关键组件

### 错误处理
- `internal/api/errors.go`：统一错误响应格式
- 区分业务错误（4xx）与基础设施错误（5xx）
- Echo HTTPErrorHandler 集成

### 数据库约束
- RecordSession 唯一约束（部分唯一索引，确保一个流只有一个 running 状态）
- InferenceResult 查询索引（stream_id + ts, algorithm_binding_id + ts）
- User、Role、Permission、Menu 唯一索引（username、code）

### 任务管理
- RecordService：内存中存储活跃录制任务，后台监控
- PreviewManager：内存中存储活跃预览任务
- Scheduler：内存中存储活跃调度任务

### 认证授权
- **JWT 认证**：`internal/api/middleware/auth.go`
  - Token 生成和解析（HS256 签名）
  - JWTAuth 中间件（Authorization: Bearer xxx）
  - 支持 Access Token 和 Refresh Token
- **权限校验**：RequirePermission 中间件
  - 基于角色-权限关系校验
  - 超级管理员（super_admin）跳过校验
- **初始化数据**：`internal/adapter/persistence/init_data.go`
  - 默认权限（30+ API 权限）
  - 默认菜单（系统管理、视频流管理等）
  - 超级管理员角色
  - 默认管理员账号（admin/admin123）

---

**注意**: 本文档会随着项目演进持续更新。
