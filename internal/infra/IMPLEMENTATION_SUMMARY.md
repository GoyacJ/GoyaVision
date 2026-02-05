# Infrastructure Adapters Implementation Summary

## 已完成的实现

### 1. ObjectStorage - MinIO Client

**文件**: `internal/infra/minio/client.go`
**接口**: `internal/app/port/object_storage.go`
**状态**: ✅ 已完成

**实现的方法**:
- ✅ `Upload(ctx, req) (*UploadResult, error)` - 上传文件到对象存储
- ✅ `Download(ctx, bucket, objectName) (io.ReadCloser, error)` - 下载文件
- ✅ `Delete(ctx, bucket, objectName) error` - 删除文件
- ✅ `GetPresignedURL(ctx, bucket, objectName, expires) (string, error)` - 生成预签名 URL
- ✅ `Exists(ctx, bucket, objectName) (bool, error)` - 检查文件是否存在
- ✅ `GetMetadata(ctx, bucket, objectName) (*ObjectMetadata, error)` - 获取文件元数据

**依赖**:
- `github.com/minio/minio-go/v7`

**关键特性**:
- 自动创建 bucket（如不存在）
- 支持自定义元数据
- 完整的错误处理（使用 `pkg/apperr`）
- 生成对象访问 URL（HTTP/HTTPS）

**配置项**:
```yaml
minio:
  endpoint: "39.105.2.5:14250"
  access_key: "goyaminio"
  secret_key: "goyaminio"
  bucket_name: "goya-vision"
  use_ssl: false
```

---

### 2. TokenService - JWT Service

**文件**: `internal/infra/auth/jwt.go`
**接口**: `internal/app/port/token_service.go`
**状态**: ✅ 已完成

**实现的方法**:
- ✅ `GenerateTokenPair(userID, username) (*TokenPair, error)` - 生成 Token 对
- ✅ `ValidateAccessToken(token) (*TokenClaims, bool, error)` - 验证 Access Token
- ✅ `ValidateRefreshToken(token) (*TokenClaims, bool, error)` - 验证 Refresh Token
- ✅ `RefreshTokenPair(refreshToken) (*TokenPair, error)` - 刷新 Token 对

**依赖**:
- `github.com/golang-jwt/jwt/v5`
- `github.com/google/uuid`

**关键特性**:
- 双 Token 机制（Access + Refresh）
- Access Token 有效期：2 小时
- Refresh Token 有效期：7 天
- HS256 签名算法
- 完整的 Token 验证（包括类型检查）
- 过期状态检测

**Token Claims 结构**:
```json
{
  "user_id": "uuid",
  "username": "admin",
  "type": "access",  // or "refresh"
  "iss": "goyavision",
  "sub": "uuid",
  "iat": 1234567890,
  "exp": 1234574890,
  "nbf": 1234567890
}
```

**配置项**:
```yaml
jwt:
  secret: "goyavision-secret-change-in-production"
  expire: 2h
  refresh_exp: 168h
  issuer: "goyavision"
```

---

### 3. EventBus - Local EventBus

**文件**: `internal/infra/eventbus/local.go`
**接口**: `internal/app/port/event_bus.go`
**状态**: ✅ 已完成

**实现的方法**:
- ✅ `Publish(ctx, event) error` - 发布事件
- ✅ `Subscribe(eventType, handler)` - 订阅事件
- ✅ `Unsubscribe(eventType, handler)` - 取消订阅

**额外方法**（用于测试和监控）:
- `GetSubscriberCount(eventType) int` - 获取订阅者数量
- `Clear()` - 清空所有订阅

**关键特性**:
- 本地内存实现（适合单机部署）
- 并发安全（`sync.RWMutex`）
- 异步处理（每个 handler 在独立 goroutine 中执行）
- Panic 恢复（handler panic 不会影响其他 handler）
- 错误日志（handler 错误会被记录但不会中断执行）
- 可配置缓冲区大小（默认 100）

**使用场景**:
1. 媒体源创建后，触发资产索引
2. 任务完成后，发送通知
3. 工作流状态变更，记录审计日志

**扩展性**:
- 当前：本地内存版本
- 未来：Redis Pub/Sub, Kafka, RabbitMQ

---

## 接口验证

每个实现都包含编译时接口验证文件：

1. `internal/infra/minio/verify_interface.go`
   ```go
   var _ port.ObjectStorage = (*Client)(nil)
   ```

2. `internal/infra/auth/verify_interface.go`
   ```go
   var _ port.TokenService = (*JWTService)(nil)
   ```

3. `internal/infra/eventbus/verify_interface.go`
   ```go
   var _ port.EventBus = (*LocalEventBus)(nil)
   ```

---

## 错误处理统一规范

所有实现使用 `pkg/apperr` 包进行错误处理：

```go
// 输入验证错误
apperr.InvalidInput("message")

// 资源不存在
apperr.NotFound("entity", id)

// 内部错误（包装原始错误）
apperr.Wrap(err, apperr.CodeInternal, "message")

// 未授权
apperr.Unauthorized("message")
```

---

## 日志规范

所有实现使用 `pkg/logger` 包进行日志记录：

```go
logger.Info("message", "key", value)
logger.Error("message", "error", err)
logger.Debug("message", "key", value)
logger.Warn("message", "key", value)
```

---

## 初始化示例

```go
package main

import (
    "goyavision/config"
    "goyavision/internal/infra/minio"
    "goyavision/internal/infra/auth"
    "goyavision/internal/infra/eventbus"
)

func main() {
    // 加载配置
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    // 初始化 MinIO 客户端
    storage, err := minio.NewClient(&cfg.MinIO)
    if err != nil {
        log.Fatal(err)
    }

    // 初始化 JWT 服务
    tokenService, err := auth.NewJWTService(&cfg.JWT)
    if err != nil {
        log.Fatal(err)
    }

    // 初始化事件总线
    eventBus := eventbus.NewLocalEventBus(100)

    // 使用适配器...
}
```

---

## 依赖包清单

需要在 `go.mod` 中添加以下依赖：

```
require (
    github.com/minio/minio-go/v7 v7.0.66
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/google/uuid v1.5.0
)
```

---

## 测试建议

### MinIO Client 测试
```bash
# 使用 Docker 启动 MinIO 测试实例
docker run -p 9000:9000 -p 9001:9001 \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  minio/minio server /data --console-address ":9001"
```

### JWT Service 测试
```go
func TestJWTService_GenerateAndValidate(t *testing.T) {
    cfg := &config.JWT{
        Secret:     "test-secret",
        Expire:     2 * time.Hour,
        RefreshExp: 7 * 24 * time.Hour,
        Issuer:     "test",
    }

    service, _ := auth.NewJWTService(cfg)

    userID := uuid.New()
    tokenPair, err := service.GenerateTokenPair(userID, "testuser")
    assert.NoError(t, err)

    claims, isExpired, err := service.ValidateAccessToken(tokenPair.AccessToken)
    assert.NoError(t, err)
    assert.False(t, isExpired)
    assert.Equal(t, userID, claims.UserID)
}
```

### EventBus 测试
```go
func TestEventBus_PublishAndSubscribe(t *testing.T) {
    bus := eventbus.NewLocalEventBus(10)

    received := false
    handler := func(ctx context.Context, event port.Event) error {
        received = true
        return nil
    }

    bus.Subscribe("test.event", handler)
    bus.Publish(context.Background(), &testEvent{eventType: "test.event"})

    time.Sleep(100 * time.Millisecond)
    assert.True(t, received)
}
```

---

## 性能优化建议

### MinIO Client
- 使用连接池（SDK 内置）
- 大文件上传使用分片上传（multipart upload）
- 考虑使用 CDN 加速对象访问

### JWT Service
- 考虑缓存已验证的 Token（使用 Redis）
- 使用 Token 黑名单机制（用于注销）
- 定期轮换 JWT Secret

### EventBus
- 避免 goroutine 泄漏（确保 handler 能够正常退出）
- 使用带超时的 context（避免 handler 长时间阻塞）
- 监控事件处理延迟和失败率

---

## 安全注意事项

1. **JWT Secret**: 生产环境必须使用强密码
   ```bash
   export GOYAVISION_JWT_SECRET=$(openssl rand -base64 32)
   ```

2. **MinIO Credentials**: 不要在代码中硬编码
   ```bash
   export GOYAVISION_MINIO_ACCESS_KEY=your-access-key
   export GOYAVISION_MINIO_SECRET_KEY=your-secret-key
   ```

3. **Token 过期策略**:
   - Access Token: 2 小时（短期）
   - Refresh Token: 7 天（长期）
   - 考虑实现 Token 黑名单（用于主动注销）

4. **预签名 URL**: 设置合理的过期时间（建议 15 分钟）

---

## 后续工作

1. ✅ ObjectStorage 实现（MinIO）
2. ✅ TokenService 实现（JWT）
3. ✅ EventBus 实现（本地内存）
4. ⏳ MediaGateway 实现（MediaMTX）- 进行中
5. ⏳ CQRS 模式拆分 Application 层
6. ⏳ 编写单元测试和集成测试
7. ⏳ 性能测试和优化

---

## 文档

- 详细使用指南: `internal/infra/README.md`
- Port 接口定义: `internal/app/port/`
- 配置文件: `configs/config.yaml`
- 项目架构: `docs/architecture.md`
