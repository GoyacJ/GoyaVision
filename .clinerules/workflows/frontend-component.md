# 前端组件开发流程

用于创建或修改 Vue 组件/页面的标准流程。

## 1. 明确组件范围

```xml
<ask_followup_question>
  <question>请描述组件用途、所在页面、数据来源（API/composables/store）以及期望的交互。</question>
</ask_followup_question>
```

## 2. 检查现有组件

我将搜索 `web/src/components` 与 `web/src/views`，确认是否已有可复用组件。

```bash
ls web/src/components
```

## 3. 创建/修改组件

- 使用 Vue 3 Composition API。
- 优先使用 Gv 组件库与 Tailwind 工具类。
- 设计令牌来源：`web/src/design-tokens`。

## 4. 校验与提示

若新增 API 调用，请同步更新 `web/src/api` 并检查类型定义。
