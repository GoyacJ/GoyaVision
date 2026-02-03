# GvFlex - Flex 布局组件

基于 CSS Flexbox 的弹性布局组件，提供简洁的 API。

## 基本用法

```vue
<template>
  <GvFlex>
    <div>项目 1</div>
    <div>项目 2</div>
    <div>项目 3</div>
  </GvFlex>
</template>

<script setup>
import { GvFlex } from '@/components'
</script>
```

## 方向

```vue
<GvFlex direction="row">水平布局</GvFlex>
<GvFlex direction="col">垂直布局</GvFlex>

<!-- 快捷方式 -->
<GvFlex vertical>垂直布局</GvFlex>
```

## 对齐方式

### 主轴对齐

```vue
<GvFlex justify="start">起点对齐</GvFlex>
<GvFlex justify="center">居中对齐</GvFlex>
<GvFlex justify="end">终点对齐</GvFlex>
<GvFlex justify="between">两端对齐</GvFlex>
<GvFlex justify="around">环绕对齐</GvFlex>
<GvFlex justify="evenly">均匀对齐</GvFlex>
```

### 交叉轴对齐

```vue
<GvFlex align="start">起点对齐</GvFlex>
<GvFlex align="center">居中对齐</GvFlex>
<GvFlex align="end">终点对齐</GvFlex>
<GvFlex align="stretch">拉伸</GvFlex>
<GvFlex align="baseline">基线对齐</GvFlex>
```

## 间距

```vue
<GvFlex gap="xs">2px 间距</GvFlex>
<GvFlex gap="sm">12px 间距</GvFlex>
<GvFlex gap="md">16px 间距</GvFlex>
<GvFlex gap="lg">24px 间距</GvFlex>
<GvFlex gap="xl">32px 间距</GvFlex>
```

## 换行

```vue
<GvFlex wrap="wrap">
  <div v-for="i in 20" :key="i" class="w-32">
    项目 {{ i }}
  </div>
</GvFlex>
```

## 内联 Flex

```vue
<GvFlex inline>
  <span>内联 Flex</span>
</GvFlex>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| direction | `'row' \| 'row-reverse' \| 'col' \| 'col-reverse'` | `'row'` | Flex 方向 |
| wrap | `'nowrap' \| 'wrap' \| 'wrap-reverse'` | `'nowrap'` | 是否换行 |
| justify | `'start' \| 'center' \| 'end' \| 'between' \| 'around' \| 'evenly'` | `'start'` | 主轴对齐 |
| align | `'start' \| 'center' \| 'end' \| 'stretch' \| 'baseline'` | `'stretch'` | 交叉轴对齐 |
| gap | `'none' \| 'xs' \| 'sm' \| 'md' \| 'lg' \| 'xl'` | `'md'` | 间距大小 |
| inline | `boolean` | `false` | 是否内联 |
| vertical | `boolean` | `false` | 是否垂直布局 |

## 使用场景

### 页面头部

```vue
<GvFlex justify="between" align="center">
  <h1 class="text-2xl font-bold">页面标题</h1>
  <GvButton icon="Plus">新建</GvButton>
</GvFlex>
```

### 操作按钮组

```vue
<GvFlex gap="sm">
  <GvButton variant="filled">保存</GvButton>
  <GvButton variant="tonal">取消</GvButton>
  <GvButton variant="text">重置</GvButton>
</GvFlex>
```

### 垂直列表

```vue
<GvFlex vertical gap="lg">
  <GvCard v-for="item in list" :key="item.id">
    {{ item.name }}
  </GvCard>
</GvFlex>
```

### 标签列表

```vue
<GvFlex wrap="wrap" gap="sm">
  <GvTag v-for="tag in tags" :key="tag" closable>
    {{ tag }}
  </GvTag>
</GvFlex>
```

### 居中内容

```vue
<GvFlex justify="center" align="center" class="min-h-screen">
  <div class="text-center">
    <h1>居中内容</h1>
  </div>
</GvFlex>
```

## 配合 Tailwind

```vue
<!-- Flex 项控制 -->
<GvFlex>
  <div class="flex-1">自动增长</div>
  <div class="flex-shrink-0">不缩小</div>
  <div class="flex-grow-0">不增长</div>
</GvFlex>

<!-- 顺序控制 -->
<GvFlex>
  <div class="order-2">第二个</div>
  <div class="order-1">第一个</div>
  <div class="order-3">第三个</div>
</GvFlex>
```

## 最佳实践

1. **选择 Flex vs Grid**：
   - 一维布局（行或列）→ Flex
   - 二维布局（行和列）→ Grid

2. **响应式**：
   - 使用 Tailwind 的响应式类名
   - 配合 direction 实现不同屏幕布局

3. **间距**：
   - 使用统一的 gap 属性
   - 避免给子元素添加 margin

4. **对齐**：
   - 垂直居中常用 align="center"
   - 两端对齐常用 justify="between"

5. **性能**：
   - Flex 比 Grid 性能更好
   - 简单布局优先使用 Flex
