# Infrastructure Layer

基础设施层实现了 `internal/app/port` 中定义的出站端口接口，提供具体的技术实现。

## 目录结构

```
internal/infra/
├── minio/          # MinIO 对象存储实现
│   ├── client.go
│   └── verify_interface.go
├── auth/           # JWT 令牌服务实现
│   ├── jwt.go
│   └── verify_interface.go
├── eventbus/       # 本地事件总线实现
│   ├── local.go
│   └── verify_interface.go
├── mediamtx/       # MediaMTX 媒体网关实现
│   └── gateway.go
└── persistence/    # GORM 数据持久化实现
    ├── model/      # GORM 模型
    ├── mapper/     # Domain <-> Model 映射器
    ├── repo/       # Repository 实现
    └── uow.go      # Unit of Work 实现
```

## 适配器实现

### 1. ObjectStorage - MinIO 客户端

**文件**: `internal/infra/minio/client.go`

**职责**:
- 文件上传和下载
- 预签名 URL 生成
- 文件元数据管理

**初始化示例**:

```go
import (
    "goyavision/config"
    "goyavision/internal/infra/minio"
)

cfg, _ := config.Load()
storage, err := minio.NewClient(&cfg.MinIO)
if err != nil {
    // 处理错误
}

// 上传文件
result, err := storage.Upload(ctx, &port.UploadRequest{
    Bucket:      "goya-vision",
    ObjectName:  "media/video.mp4",
    Reader:      fileReader,
    Size:        fileSize,
    ContentType: "video/mp4",
    Metadata:    map[string]string{"source": "upload"},
})

// 下载文件
reader, err := storage.Download(ctx, "goya-vision", "media/video.mp4")
defer reader.Close()

// 生成预签名 URL（15 分钟有效）
url, err := storage.GetPresignedURL(ctx, "goya-vision", "media/video.mp4", 15*time.Minute)
```

**配置** (configs/config.<env>.yaml):

```yaml
minio:
  endpoint: "39.105.2.5:14250"
  access_key: "goyaminio"
  secret_key: "goyaminio"
  bucket_name: "goya-vision"
  use_ssl: false
```

**依赖**: `github.com/minio/minio-go/v7`

---

### 2. TokenService - JWT 服务

**文件**: `internal/infra/auth/jwt.go`

**职责**:
- 生成 Access Token (2 小时) 和 Refresh Token (7 天)
- 验证和解析 Token
- Token 刷新机制

**初始化示例**:

```go
import (
    "goyavision/config"
    "goyavision/internal/infra/auth"
)

cfg, _ := config.Load()
tokenService, err := auth.NewJWTService(&cfg.JWT)
if err != nil {
    // 处理错误
}

// 生成 Token 对
tokenPair, err := tokenService.GenerateTokenPair(userID, username)
// tokenPair.AccessToken  - 用于 API 请求
// tokenPair.RefreshToken - 用于刷新 Access Token
// tokenPair.ExpiresIn    - Access Token 过期时间（秒）
// tokenPair.ExpiresAt    - Access Token 过期时刻

// 验证 Access Token
claims, isExpired, err := tokenService.ValidateAccessToken(accessToken)
if err != nil {
    if isExpired {
        // Token 已过期，需要刷新
    } else {
        // Token 无效
    }
}

// 刷新 Token
newTokenPair, err := tokenService.RefreshTokenPair(refreshToken)
```

**配置** (configs/config.<env>.yaml):

```yaml
jwt:
  secret: "goyavision-secret-change-in-production"  # ⚠️ 生产环境必须更改
  expire: 2h                                        # Access Token 有效期
  refresh_exp: 168h                                 # Refresh Token 有效期（7天）
  issuer: "goyavision"
```

**依赖**: `github.com/golang-jwt/jwt/v5`

**Token 结构**:

```json
{
  "user_id": "uuid",
  "username": "admin",
  "type": "access",
  "iss": "goyavision",
  "sub": "uuid",
  "iat": 1234567890,
  "exp": 1234574890,
  "nbf": 1234567890
}
```

---

### 3. EventBus - 本地事件总线

**文件**: `internal/infra/eventbus/local.go`

**职责**:
- 发布领域事件
- 订阅事件处理器
- 异步事件处理

**初始化示例**:

```go
import "goyavision/internal/infra/eventbus"

// 创建事件总线（缓冲区大小 100）
eventBus := eventbus.NewLocalEventBus(100)

// 订阅事件
eventBus.Subscribe("media.source.created", func(ctx context.Context, event port.Event) error {
    // 处理事件
    log.Printf("received event: %s", event.EventType())
    return nil
})

// 发布事件
err := eventBus.Publish(ctx, &MediaSourceCreatedEvent{
    SourceID:   uuid.New(),
    SourceName: "RTSP Camera 01",
    OccurredAt: time.Now().Unix(),
})

// 取消订阅
eventBus.Unsubscribe("media.source.created", handler)
```

**特性**:
- ✅ 并发安全（使用 `sync.RWMutex`）
- ✅ 异步处理（每个 handler 在独立 goroutine 中执行）
- ✅ Panic 恢复（handler panic 不会影响其他 handler）
- ✅ 错误日志（handler 错误会被记录但不会中断执行）

**典型使用场景**:

```go
// 1. 媒体源创建后，触发资产索引
eventBus.Subscribe("media.source.created", func(ctx context.Context, event port.Event) error {
    e := event.(*MediaSourceCreatedEvent)
    return indexService.IndexMediaSource(ctx, e.SourceID)
})

// 2. 任务完成后，发送通知
eventBus.Subscribe("task.completed", func(ctx context.Context, event port.Event) error {
    e := event.(*TaskCompletedEvent)
    return notificationService.NotifyTaskResult(ctx, e.TaskID)
})

// 3. 工作流状态变更，记录审计日志
eventBus.Subscribe("workflow.status.changed", func(ctx context.Context, event port.Event) error {
    e := event.(*WorkflowStatusChangedEvent)
    return auditService.LogWorkflowChange(ctx, e)
})
```

**扩展性**:
- 当前实现为本地内存版本（适合单机部署）
- 未来可扩展为分布式版本：
  - Redis Pub/Sub
  - Apache Kafka
  - RabbitMQ

---

## 依赖注入示例

在应用启动时初始化所有基础设施组件：

```go
// cmd/server/main.go

func main() {
    // 加载配置
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    // 初始化基础设施
    storage, err := minio.NewClient(&cfg.MinIO)
    if err != nil {
        log.Fatal(err)
    }

    tokenService, err := auth.NewJWTService(&cfg.JWT)
    if err != nil {
        log.Fatal(err)
    }

    eventBus := eventbus.NewLocalEventBus(100)

    // 初始化 GORM
    db, err := gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 初始化 Unit of Work
    uow := persistence.NewUnitOfWork(db)

    // 初始化 MediaMTX Gateway
    mediaGateway := mediamtx.NewGateway(
        cfg.MediaMTX.APIAddress,
        cfg.MediaMTX.Username, cfg.MediaMTX.Password,
        cfg.MediaMTX.RecordPath, cfg.MediaMTX.RecordFormat, cfg.MediaMTX.SegmentDuration,
    )

    // 初始化 Application Services
    mediaSourceService := app.NewMediaSourceService(uow, mediaGateway, eventBus)
    authService := app.NewAuthService(uow, tokenService)
    // ...

    // 启动 HTTP 服务器
    router := api.NewRouter(mediaSourceService, authService, ...)
    router.Start(cfg.Server.Addr())
}
```

## 测试

每个适配器都有对应的接口验证文件（`verify_interface.go`），确保编译时接口实现的正确性。

**单元测试示例**:

```go
// internal/infra/eventbus/local_test.go

func TestLocalEventBus_PublishAndSubscribe(t *testing.T) {
    bus := NewLocalEventBus(10)

    received := false
    handler := func(ctx context.Context, event port.Event) error {
        received = true
        return nil
    }

    bus.Subscribe("test.event", handler)

    err := bus.Publish(context.Background(), &testEvent{eventType: "test.event"})
    assert.NoError(t, err)

    time.Sleep(100 * time.Millisecond) // 等待异步处理
    assert.True(t, received)
}
```

## 错误处理

所有适配器使用统一的错误处理机制（`pkg/apperr`）:

```go
// 输入验证错误
return apperr.InvalidInput("object name is required")

// 资源不存在
return apperr.NotFound("object", objectName)

// 内部错误（包装原始错误）
return apperr.Wrap(err, apperr.CodeInternal, "failed to upload object")

// 未授权
return apperr.Unauthorized("invalid token")
```

## 日志

所有适配器使用统一的日志包（`pkg/logger`）:

```go
logger.Info("object uploaded", "bucket", bucket, "object", objectName)
logger.Error("upload failed", "error", err)
logger.Debug("handler subscribed", "event_type", eventType)
logger.Warn("invalid configuration", "field", "endpoint")
```

## 性能注意事项

1. **MinIO Client**: 使用连接池，支持并发上传/下载
2. **JWT Service**: Token 生成和验证是 CPU 密集型操作，考虑缓存已验证的 Token
3. **EventBus**: 异步处理避免阻塞发布者，但需注意 goroutine 泄漏
4. **Database Connection**: GORM 内置连接池，注意设置合理的 `MaxOpenConns` 和 `MaxIdleConns`

## 安全注意事项

1. **JWT Secret**: 生产环境必须使用强密码，建议通过环境变量配置
2. **MinIO Credentials**: 不要在代码中硬编码，使用环境变量或密钥管理服务
3. **Token 过期**: Access Token 应设置较短的有效期（2小时），Refresh Token 设置较长的有效期（7天）
4. **预签名 URL**: 设置合理的过期时间，避免长期有效的 URL 泄露

## 配置最佳实践

```yaml
# 开发环境
minio:
  endpoint: "localhost:9000"
  access_key: "minioadmin"
  secret_key: "minioadmin"
  bucket_name: "goyavision-dev"
  use_ssl: false

jwt:
  secret: "dev-secret-key"
  expire: 2h
  refresh_exp: 168h
  issuer: "goyavision"

# 生产环境（使用环境变量覆盖）
# GOYAVISION_MINIO_ENDPOINT=minio.example.com:9000
# GOYAVISION_MINIO_ACCESS_KEY=${MINIO_ACCESS_KEY}
# GOYAVISION_MINIO_SECRET_KEY=${MINIO_SECRET_KEY}
# GOYAVISION_MINIO_USE_SSL=true
# GOYAVISION_JWT_SECRET=${JWT_SECRET}
```
