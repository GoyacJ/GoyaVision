---
paths:
  - "web/src/**"
---

# 前端开发规则

## 技术栈
- Vue 3 + TypeScript + Vite。
- UI 组件优先使用 GoyaVision 组件库（`web/src/components`），复杂组件使用 Element Plus。
- 样式优先使用 Tailwind 与设计令牌（`web/src/design-tokens`）。

## 组件与结构
- 组件名使用 PascalCase，文件名使用 kebab-case。
- Composition API 优先，不使用 Options API。
- API 调用统一放在 `web/src/api`，页面逻辑放在 `web/src/views`。
- 业务组件优先复用 `components/business`。

## 类型与质量
- 所有 Props/Emits 必须定义类型。
- TypeScript 严格模式下不得忽略类型错误。
- 避免在组件内硬编码颜色，使用 design tokens。
