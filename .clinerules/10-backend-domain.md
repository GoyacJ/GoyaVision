---
paths:
  - "internal/domain/**"
  - "internal/port/**"
---

# 后端领域与端口规则

## 分层边界
- Domain 层仅包含纯业务实体与验证逻辑，不依赖任何外部包。
- Port 层只定义接口与类型，不引入 Adapter 具体实现。
- 禁止在 Domain/Port 引入 `internal/adapter` 或 `internal/api`。

## 领域建模
- 实体名称与字段遵循现有命名：MediaSource、MediaAsset、Operator、Workflow、Task、Artifact、User、Role、Menu。
- 领域状态与枚举保持与现有常量一致，避免新增未经讨论的状态值。
- 涉及资产派生必须维护 parent_id 关系。

## 质量要求
- 所有错误必须返回或记录，不吞错误。
- 所有可取消操作使用 context.Context。
