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
- `StreamService`: 流管理用例
- `RecordService`: 录制用例
- `InferenceService`: 推理用例
- `PreviewService`: 预览用例

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
- `persistence`: 实现 `port.Repository`（GORM）
- `ffmpeg`: FFmpeg 进程管理（规划中）
- `preview`: 预览服务（规划中）
- `ai`: AI 推理服务调用（规划中）

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
- `router.go`: 路由注册
- `handler/`: 请求处理器
- `dto/`: 数据传输对象
- `middleware.go`: 中间件

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

**位置**: `pkg/ffmpeg/`

- FFmpeg 进程池
- 资源限流（最大录制数、最大抽帧数）
- 进程生命周期管理

## 扩展点

### 添加新的基础设施适配器

1. 在 `port/` 定义接口
2. 在 `adapter/` 实现接口
3. 在 `app/` 或 `api/` 中注入使用

### 添加新的业务用例

1. 在 `app/` 创建新的 Service
2. 编排 domain 和 port
3. 在 `api/handler/` 中调用

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

- 输入验证（API 层）
- SQL 注入防护（使用 GORM 参数化查询）
- 认证与鉴权（规划中）
- 敏感信息加密（规划中）

---

**注意**: 本文档会随着项目演进持续更新。
