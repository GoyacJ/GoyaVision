# GvGrid - 网格布局组件

基于 CSS Grid 的网格布局组件，提供简洁的 API 和响应式支持。

## 基本用法

```vue
<template>
  <GvGrid :cols="3">
    <GvCard>项目 1</GvCard>
    <GvCard>项目 2</GvCard>
    <GvCard>项目 3</GvCard>
  </GvGrid>
</template>

<script setup>
import { GvGrid, GvCard } from '@/components'
</script>
```

## 列数

```vue
<GvGrid :cols="2">
  <div>2 列布局</div>
</GvGrid>

<GvGrid :cols="3">
  <div>3 列布局</div>
</GvGrid>

<GvGrid :cols="4">
  <div>4 列布局</div>
</GvGrid>

<GvGrid :cols="6">
  <div>6 列布局</div>
</GvGrid>

<GvGrid :cols="12">
  <div>12 列栅格</div>
</GvGrid>
```

## 响应式布局

```vue
<GvGrid
  :responsive="{
    xs: 1,
    sm: 2,
    md: 3,
    lg: 4,
    xl: 6
  }"
>
  <GvCard v-for="i in 12" :key="i">
    项目 {{ i }}
  </GvCard>
</GvGrid>
```

## 间距

```vue
<GvGrid gap="none">无间距</GvGrid>
<GvGrid gap="xs">2 像素间距</GvGrid>
<GvGrid gap="sm">12 像素间距</GvGrid>
<GvGrid gap="md">16 像素间距</GvGrid>
<GvGrid gap="lg">24 像素间距</GvGrid>
<GvGrid gap="xl">32 像素间距</GvGrid>
```

## 不同方向间距

```vue
<GvGrid gap-x="lg" gap-y="sm">
  <div>水平间距大，垂直间距小</div>
</GvGrid>
```

## 自动填充

```vue
<!-- 自动填充，每列最小宽度 200px -->
<GvGrid auto-fill>
  <GvCard v-for="i in 20" :key="i">
    项目 {{ i }}
  </GvCard>
</GvGrid>

<!-- 自定义最小列宽 -->
<GvGrid auto-fill min-col-width="250px">
  <GvCard v-for="i in 20" :key="i">
    项目 {{ i }}
  </GvCard>
</GvGrid>
```

## 对齐方式

```vue
<!-- 垂直对齐 -->
<GvGrid align="start">顶部对齐</GvGrid>
<GvGrid align="center">垂直居中</GvGrid>
<GvGrid align="end">底部对齐</GvGrid>
<GvGrid align="stretch">拉伸填充</GvGrid>

<!-- 水平对齐 -->
<GvGrid justify="start">左对齐</GvGrid>
<GvGrid justify="center">水平居中</GvGrid>
<GvGrid justify="end">右对齐</GvGrid>
<GvGrid justify="between">两端对齐</GvGrid>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| cols | `1 \| 2 \| 3 \| 4 \| 5 \| 6 \| 12` | `3` | 列数 |
| responsive | `Object` | - | 响应式列数配置 |
| gap | `'none' \| 'xs' \| 'sm' \| 'md' \| 'lg' \| 'xl'` | `'md'` | 间距大小 |
| gapX | `GapSize` | - | 水平间距 |
| gapY | `GapSize` | - | 垂直间距 |
| autoFill | `boolean` | `false` | 是否自动填充 |
| minColWidth | `string` | `'200px'` | 自动填充时的最小列宽 |
| align | `'start' \| 'center' \| 'end' \| 'stretch'` | - | 垂直对齐 |
| justify | `'start' \| 'center' \| 'end' \| 'between' \| 'around' \| 'evenly'` | - | 水平对齐 |

## 使用场景

### 卡片网格

```vue
<GvGrid
  :responsive="{ xs: 1, sm: 2, md: 3, lg: 4 }"
  gap="lg"
>
  <GvCard v-for="asset in assets" :key="asset.id" hoverable>
    <template #header>
      <h3>{{ asset.name }}</h3>
    </template>
    <p>{{ asset.description }}</p>
  </GvCard>
</GvGrid>
```

### 图片画廊

```vue
<GvGrid auto-fill min-col-width="200px" gap="md">
  <img
    v-for="image in images"
    :key="image.id"
    :src="image.url"
    class="w-full h-auto rounded-lg"
  />
</GvGrid>
```

### 仪表板

```vue
<GvGrid :cols="12" gap="lg">
  <!-- 占 8 列 -->
  <div class="col-span-8">
    <GvCard>主要内容</GvCard>
  </div>
  
  <!-- 占 4 列 -->
  <div class="col-span-4">
    <GvCard>侧边栏</GvCard>
  </div>
</GvGrid>
```

### 表单布局

```vue
<GvGrid :cols="2" gap="lg">
  <GvInput label="名称" />
  <GvInput label="类型" />
  <GvInput label="描述" class="col-span-2" />
  <GvSelect label="分类" :options="[]" />
  <GvSelect label="状态" :options="[]" />
</GvGrid>
```

## 最佳实践

1. **列数选择**：
   - 2-3 列：常规卡片列表
   - 4-6 列：图标、小卡片
   - 12 列：复杂页面布局

2. **响应式**：
   - 移动端使用 1 列
   - 平板使用 2-3 列
   - 桌面使用 3-4 列

3. **自动填充**：
   - 图片画廊使用 autoFill
   - 不确定项目数量时使用

4. **间距**：
   - 卡片列表使用 lg
   - 紧凑列表使用 sm
   - 密集网格使用 xs

5. **配合 Tailwind**：
   - 使用 col-span-* 控制跨列
   - 使用 row-span-* 控制跨行
