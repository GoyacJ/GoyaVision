---
paths:
  - "internal/app/**"
---

# 后端应用层规则

## 依赖与编排
- 仅依赖 Domain + Port，不直接依赖 Adapter。
- 通过接口注入 Repository/Engine/Client 等端口。
- 一个服务对应一个业务用例，避免过度聚合。

## 业务逻辑
- 业务验证放在 App 层完成，API 层只做输入校验与 DTO 转换。
- 任务状态流转需遵循既有状态机定义。
- 触发工作流或任务时，必须维护 Artifact 与 Asset 关联。

## 错误处理
- 统一返回项目错误类型（参考 pkg/apperr 或 api/errors）。
- 不吞错误，必要时记录结构化日志。
