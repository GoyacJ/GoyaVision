---
name: frontend-components
description: 前端组件开发指南。用于新增/修改 Vue 组件、页面、样式与设计令牌。
---

# 前端组件开发指南

适用于 Vue 3 + TypeScript 前端开发的组件与页面实现。

## 何时使用

✅ **推荐场景**：
- 新增/修改页面（views）
- 新增/修改组件（components）
- 调整样式与设计令牌（design-tokens / tailwind）

❌ **不适用场景**：
- 只需要查看现有组件（使用 Read 工具）
- 后端 API 开发（使用 backend-adapter-api rule）

## 基本规范

- 使用 Vue 3 Composition API
- UI 组件优先使用 `web/src/components` 中的 Gv 组件
- Element Plus 复杂组件按需使用，但不要从 `element-plus` 直接导入
- 样式优先使用 Tailwind + design tokens

## 目录约定

- **页面**：`web/src/views`
- **组件**：`web/src/components`（base/business/layout/common）
- **API**：`web/src/api`
- **组合式函数**：`web/src/composables`
- **工具函数**：`web/src/utils`
- **状态管理**：`web/src/stores`（Pinia）

## 组件结构

```vue
<template>
  <!-- 模板内容 -->
</template>

<script setup lang="ts">
// 导入
import { ref, computed } from 'vue'
import type { ComponentProps } from './types'

// Props 定义
interface Props {
  // ...
}
const props = defineProps<Props>()

// Emits 定义
const emit = defineEmits<{
  // ...
}>()

// 逻辑
</script>

<style scoped>
/* 样式（优先使用 Tailwind） */
</style>
```

## 类型与质量

- 所有 Props/Emits 必须定义类型
- TypeScript 严格模式下不忽略类型错误
- 禁止硬编码颜色值，使用 design tokens
- 组件名使用 PascalCase，文件名使用 kebab-case

## 样式规范

- 优先使用 Tailwind 工具类
- 使用 design tokens（`web/src/design-tokens`）
- 避免内联样式，使用 scoped CSS
- 响应式设计使用 Tailwind 断点（sm/md/lg/xl）

## 常见模式

### 列表页面
- 使用 `FilterBar` 组件进行筛选
- 使用 `DataTable` 或 `GvGrid` 展示数据
- 使用 `GvPagination` 进行分页

### 表单页面
- 使用 `el-form` 进行表单验证
- 使用 `GvButton` 进行操作
- 使用 `GvInput`、`GvSelect` 等表单组件

### 详情页面
- 使用 `GvCard` 展示信息
- 使用 `GvBadge` 展示状态
- 使用 `GvTag` 展示标签

## 相关资源

- 组件库：`web/src/components`
- 设计令牌：`web/src/design-tokens`
- API 客户端：`web/src/api`
- 路由配置：`web/src/router`
