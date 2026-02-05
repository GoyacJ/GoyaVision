---
name: frontend-components
description: 前端组件开发指南。用于新增/修改 Vue 组件、页面、样式与设计令牌。
---

# 前端组件开发指南

适用于 Vue 3 + TypeScript 前端开发的组件与页面实现。

## 何时使用
- 新增/修改页面（views）
- 新增/修改组件（components）
- 调整样式与设计令牌（design-tokens / tailwind）

## 基本规范
- 使用 Vue 3 Composition API。
- UI 组件优先使用 `web/src/components` 中的 Gv 组件。
- Element Plus 复杂组件按需使用，但不要从 `element-plus` 直接导入。
- 样式优先使用 Tailwind + design tokens。

## 目录约定
- 页面：`web/src/views`
- 组件：`web/src/components`（base/business/layout/common）
- API：`web/src/api`
- 组合式函数：`web/src/composables`

## 类型与质量
- 所有 Props/Emits 定义类型。
- TypeScript 严格模式下不忽略类型错误。
- 禁止硬编码颜色值。
