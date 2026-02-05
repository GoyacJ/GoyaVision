---
paths:
  - "internal/adapter/**"
  - "internal/api/**"
---

# 后端适配器与 API 层规则

## Adapter 层
- 仅实现 Port 接口，禁止直接被 App/Domain 依赖。
- 处理外部系统细节（DB、HTTP、流媒体、引擎），需转换为 Domain 模型。
- 使用 context.Context 传播取消与超时。

## API 层
- Handler 只做参数校验、DTO 转换与调用 App/Port，不直接访问数据库。
- 响应必须使用统一错误与响应封装（internal/api/errors.go / response）。
- DTO 与 Domain 分离，禁止直接暴露 Domain。
- 保持 RESTful 路径与 `/api/v1` 前缀一致。

## 安全
- JWT 与 RBAC 必须通过 middleware 校验。
- 避免返回敏感信息。
