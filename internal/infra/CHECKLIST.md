# Infrastructure Adapters Implementation Checklist

## ✅ ObjectStorage 实现 (MinIO Client)

**文件**: `internal/infra/minio/client.go`

- ✅ 实现 `port.ObjectStorage` 接口
- ✅ 使用 MinIO Go SDK (`github.com/minio/minio-go/v7`)
- ✅ 实现 `Upload` 方法
  - ✅ 支持自定义 bucket
  - ✅ 支持 ContentType 设置
  - ✅ 支持自定义元数据
  - ✅ 返回完整的上传结果（ETag, Size, URL）
- ✅ 实现 `Download` 方法
  - ✅ 返回 `io.ReadCloser`
  - ✅ 检查对象是否存在
  - ✅ 处理 NoSuchKey 错误
- ✅ 实现 `Delete` 方法
  - ✅ 完整的错误处理
- ✅ 实现 `GetPresignedURL` 方法
  - ✅ 支持自定义过期时间
  - ✅ 默认 15 分钟有效期
- ✅ 实现 `Exists` 方法
  - ✅ 返回布尔值，不抛出错误（对象不存在时）
- ✅ 实现 `GetMetadata` 方法
  - ✅ 返回完整的元数据信息
- ✅ 配置验证
  - ✅ Endpoint 必填
  - ✅ AccessKey 必填
  - ✅ SecretKey 必填
  - ✅ BucketName 必填
- ✅ 自动创建 bucket（如不存在）
- ✅ 使用统一错误处理 (`pkg/apperr`)
- ✅ 构建对象访问 URL（支持 HTTP/HTTPS）
- ✅ 接口验证文件 (`verify_interface.go`)

---

## ✅ TokenService 实现 (JWT Service)

**文件**: `internal/infra/auth/jwt.go`

- ✅ 实现 `port.TokenService` 接口
- ✅ 使用 `github.com/golang-jwt/jwt/v5`
- ✅ 实现 `GenerateTokenPair` 方法
  - ✅ 生成 Access Token（2 小时有效）
  - ✅ 生成 Refresh Token（7 天有效）
  - ✅ 返回 ExpiresIn 和 ExpiresAt
  - ✅ 输入验证（userID 和 username）
- ✅ 实现 `ValidateAccessToken` 方法
  - ✅ 验证签名
  - ✅ 验证 Token 类型（必须是 "access"）
  - ✅ 返回 Claims、是否过期、错误
  - ✅ 处理过期 Token
- ✅ 实现 `ValidateRefreshToken` 方法
  - ✅ 验证签名
  - ✅ 验证 Token 类型（必须是 "refresh"）
  - ✅ 返回 Claims、是否过期、错误
- ✅ 实现 `RefreshTokenPair` 方法
  - ✅ 验证 Refresh Token
  - ✅ 生成新的 Token 对
  - ✅ 处理过期 Refresh Token
- ✅ 自定义 Claims 结构
  - ✅ UserID（UUID 字符串）
  - ✅ Username
  - ✅ Type（"access" 或 "refresh"）
  - ✅ 标准 Claims（iss, sub, iat, exp, nbf）
- ✅ HS256 签名算法
- ✅ 配置验证
  - ✅ Secret 必填
  - ✅ Expire 默认 2 小时
  - ✅ RefreshExp 默认 7 天
  - ✅ Issuer 默认 "goyavision"
- ✅ 使用统一错误处理 (`pkg/apperr`)
  - ✅ `CodeTokenExpired` - Token 过期
  - ✅ `CodeTokenInvalid` - Token 无效
  - ✅ `CodeUnauthorized` - 未授权
- ✅ 接口验证文件 (`verify_interface.go`)

---

## ✅ EventBus 实现 (Local EventBus)

**文件**: `internal/infra/eventbus/local.go`

- ✅ 实现 `port.EventBus` 接口
- ✅ 本地内存实现（使用 map + channel）
- ✅ 实现 `Publish` 方法
  - ✅ 验证 event 不为 nil
  - ✅ 验证 eventType 不为空
  - ✅ 异步调用所有 handler（goroutine）
  - ✅ 复制 handler 列表（避免并发修改）
  - ✅ 无 handler 时记录 Debug 日志
- ✅ 实现 `Subscribe` 方法
  - ✅ 验证 eventType 不为空
  - ✅ 验证 handler 不为 nil
  - ✅ 分配唯一 handler ID
  - ✅ 记录订阅日志
- ✅ 实现 `Unsubscribe` 方法
  - ✅ 验证 eventType 不为空
  - ✅ 验证 handler 不为 nil
  - ✅ 删除 handler
  - ✅ 清理空的 eventType 映射
  - ✅ 记录取消订阅日志
- ✅ 并发安全
  - ✅ 使用 `sync.RWMutex` 保护 handlers map
  - ✅ 使用独立 mutex 保护 handler ID 计数器
  - ✅ 复制 handler 列表后释放锁
- ✅ 异步处理
  - ✅ 每个 handler 在独立 goroutine 中执行
  - ✅ Panic 恢复（不影响其他 handler）
  - ✅ 错误日志（不中断执行）
- ✅ 可配置缓冲区大小（默认 100）
- ✅ 额外方法（用于测试和监控）
  - ✅ `GetSubscriberCount` - 获取订阅者数量
  - ✅ `Clear` - 清空所有订阅
- ✅ 使用统一日志 (`pkg/logger`)
  - ✅ Debug 日志（订阅、取消订阅、无 handler）
  - ✅ Error 日志（handler panic、handler error）
  - ✅ Warn 日志（无效输入）
- ✅ 接口验证文件 (`verify_interface.go`)

---

## 📋 通用要求检查

### 代码质量
- ✅ 所有方法完整实现（无占位符）
- ✅ 完善的错误处理
- ✅ 适当的代码注释
- ✅ 遵循 Go 命名规范（CamelCase）
- ✅ 无行尾注释（用户偏好）

### 依赖管理
- ✅ MinIO: `github.com/minio/minio-go/v7`
- ✅ JWT: `github.com/golang-jwt/jwt/v5`
- ✅ UUID: `github.com/google/uuid`

### 配置管理
- ✅ MinIO 配置（endpoint, access_key, secret_key, bucket_name, use_ssl）
- ✅ JWT 配置（secret, expire, refresh_exp, issuer）
- ✅ EventBus 配置（bufferSize）

### 错误处理
- ✅ 使用 `pkg/apperr` 统一错误包
- ✅ 输入验证错误：`apperr.InvalidInput`
- ✅ 资源不存在：`apperr.NotFound`
- ✅ 内部错误：`apperr.Wrap`
- ✅ 未授权：`apperr.Unauthorized`

### 日志记录
- ✅ 使用 `pkg/logger` 统一日志包
- ✅ 结构化日志（key-value pairs）
- ✅ 适当的日志级别（Info, Error, Debug, Warn）

### 接口验证
- ✅ 每个实现都有 `verify_interface.go`
- ✅ 编译时接口类型检查

---

## 📄 文档

- ✅ 实现总结文档 (`IMPLEMENTATION_SUMMARY.md`)
- ✅ 详细使用指南 (`README.md`)
- ✅ 实现检查清单 (`CHECKLIST.md`)

---

## 🧪 测试建议

### MinIO Client
- [ ] 上传文件测试
- [ ] 下载文件测试
- [ ] 删除文件测试
- [ ] 预签名 URL 测试
- [ ] 文件存在性检查测试
- [ ] 元数据获取测试
- [ ] 错误处理测试（对象不存在、网络错误等）

### JWT Service
- [ ] 生成 Token 对测试
- [ ] 验证 Access Token 测试
- [ ] 验证 Refresh Token 测试
- [ ] 刷新 Token 测试
- [ ] Token 过期测试
- [ ] Token 类型验证测试
- [ ] 无效 Token 测试

### EventBus
- [ ] 发布和订阅测试
- [ ] 多个 handler 测试
- [ ] 取消订阅测试
- [ ] 并发发布测试
- [ ] Handler panic 恢复测试
- [ ] Handler 错误处理测试
- [ ] 订阅者计数测试

---

## 🚀 部署注意事项

### 生产环境配置
- ⚠️ 修改 JWT Secret（使用强密码）
- ⚠️ 使用环境变量覆盖敏感配置
- ⚠️ 启用 MinIO SSL（use_ssl: true）
- ⚠️ 配置合理的 Token 过期时间

### 环境变量示例
```bash
export GOYAVISION_MINIO_ENDPOINT=minio.example.com:9000
export GOYAVISION_MINIO_ACCESS_KEY=your-access-key
export GOYAVISION_MINIO_SECRET_KEY=your-secret-key
export GOYAVISION_MINIO_USE_SSL=true
export GOYAVISION_JWT_SECRET=$(openssl rand -base64 32)
```

---

## 📊 总结

### 完成情况
- ✅ ObjectStorage 实现 (100%)
- ✅ TokenService 实现 (100%)
- ✅ EventBus 实现 (100%)
- ⏳ MediaGateway 实现（进行中）

### 代码统计
- **MinIO Client**: 242 行代码，7 个方法
- **JWT Service**: 181 行代码，7 个方法（4 个公开，3 个私有）
- **EventBus**: 164 行代码，5 个方法（3 个接口方法，2 个辅助方法）
- **总计**: 587 行代码

### 特性亮点
1. 完全实现 Port 接口，遵循 Clean Architecture
2. 完善的错误处理和日志记录
3. 并发安全（EventBus）
4. 异步处理（EventBus）
5. 配置验证和默认值
6. 接口编译时验证
7. 详细的文档和使用示例

### 下一步
1. 完成 MediaGateway 实现
2. 编写单元测试和集成测试
3. 性能测试和优化
4. 在 Application 层中使用这些适配器
5. 更新 API 层以使用新的服务
