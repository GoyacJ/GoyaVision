# frontend-component

前端组件开发流程。用于创建或修改 Vue 组件/页面的标准流程。

## 1. 明确组件范围

在开始开发前，请明确：
- 组件用途和所在页面
- 数据来源（API/composables/store）
- 期望的交互行为
- 是否需要复用现有组件

## 2. 检查现有组件

搜索 `web/src/components` 与 `web/src/views`，确认是否已有可复用组件：

```bash
# 查看组件目录
ls web/src/components

# 查看业务组件
ls web/src/components/business

# 查看基础组件
ls web/src/components/base
```

## 3. 创建/修改组件

### 组件结构

```vue
<template>
  <!-- 使用 Gv 组件和 Tailwind -->
  <GvCard>
    <GvButton variant="filled">操作</GvButton>
  </GvCard>
</template>

<script setup lang="ts">
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
/* 优先使用 Tailwind，必要时使用 scoped CSS */
</style>
```

### 开发规范

- 使用 Vue 3 Composition API
- 优先使用 Gv 组件库与 Tailwind 工具类
- 设计令牌来源：`web/src/design-tokens`
- 组件名使用 PascalCase，文件名使用 kebab-case

## 4. 校验与提示

### 检查清单

- [ ] Props 和 Emits 已定义类型
- [ ] 使用了 Gv 组件或 Element Plus（按需）
- [ ] 样式使用 Tailwind 和 design tokens
- [ ] 响应式设计已考虑（移动端适配）
- [ ] 无硬编码颜色值

### API 集成

若新增 API 调用：
- 更新 `web/src/api` 中的 API 客户端
- 检查类型定义是否完整
- 使用统一的错误处理

## 5. 测试

- 组件功能测试
- 响应式布局测试
- 交互行为测试

## 相关资源

- 组件库：`web/src/components`
- 设计令牌：`web/src/design-tokens`
- API 客户端：`web/src/api`
- 组合式函数：`web/src/composables`
