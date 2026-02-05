# API 文档更新流程

用于新增或修改 API 端点后的文档同步。

## 1. 确认端点变更

```xml
<ask_followup_question>
  <question>请提供变更的 API 端点（方法 + 路径）与主要请求/响应字段。</question>
</ask_followup_question>
```

## 2. 更新 docs/api.md

在对应模块中补充：
- 端点说明与用途
- 请求参数（Path/Query/Body）
- 请求示例
- 响应结构与示例
- 错误码说明
