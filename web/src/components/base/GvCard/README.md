# GvCard - 卡片组件

Material Design 3 风格的卡片容器组件。

## 基本用法

```vue
<template>
  <GvCard>
    <p>卡片内容</p>
  </GvCard>
</template>

<script setup>
import { GvCard } from '@/components'
</script>
```

## 带头部和底部

```vue
<GvCard>
  <template #header>
    <h3>卡片标题</h3>
  </template>
  
  <p>卡片内容区域</p>
  
  <template #footer>
    <div class="flex justify-end gap-2">
      <GvButton variant="text">取消</GvButton>
      <GvButton variant="filled">确定</GvButton>
    </div>
  </template>
</GvCard>
```

## 阴影大小

```vue
<GvCard shadow="none">无阴影</GvCard>
<GvCard shadow="sm">小阴影</GvCard>
<GvCard shadow="md">中等阴影</GvCard>
<GvCard shadow="lg">大阴影</GvCard>
<GvCard shadow="xl">超大阴影</GvCard>
```

## 内边距

```vue
<GvCard padding="none">无内边距</GvCard>
<GvCard padding="sm">小内边距</GvCard>
<GvCard padding="md">中等内边距</GvCard>
<GvCard padding="lg">大内边距</GvCard>
```

## 可悬停卡片

```vue
<GvCard hoverable @click="handleCardClick">
  <p>鼠标悬停时会有动画效果</p>
</GvCard>
```

## 带边框

```vue
<GvCard bordered>
  <p>带边框的卡片</p>
</GvCard>
```

## 自定义背景

```vue
<GvCard background="container">
  <p>容器背景色</p>
</GvCard>

<GvCard background="transparent">
  <p>透明背景</p>
</GvCard>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| shadow | `'none' \| 'sm' \| 'md' \| 'lg' \| 'xl'` | `'md'` | 阴影大小 |
| padding | `'none' \| 'sm' \| 'md' \| 'lg'` | `'md'` | 内边距 |
| hoverable | `boolean` | `false` | 是否支持悬停效果 |
| bordered | `boolean` | `false` | 是否显示边框 |
| background | `'default' \| 'container' \| 'transparent'` | `'default'` | 背景色 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| click | `(event: MouseEvent)` | 点击事件（仅在 hoverable 为 true 时触发） |

## Slots

| 插槽名 | 说明 |
|--------|------|
| header | 卡片头部 |
| default | 卡片主体内容 |
| footer | 卡片底部 |

## 使用场景

### 信息展示卡片

```vue
<GvCard shadow="md" padding="lg">
  <h3 class="text-lg font-semibold mb-2">资产名称</h3>
  <p class="text-sm text-text-secondary">资产描述信息</p>
</GvCard>
```

### 可点击卡片

```vue
<GvCard hoverable @click="handleView">
  <div class="flex items-center gap-4">
    <img src="thumbnail.jpg" class="w-16 h-16 rounded" />
    <div>
      <h4 class="font-medium">资产标题</h4>
      <p class="text-sm text-text-secondary">点击查看详情</p>
    </div>
  </div>
</GvCard>
```

### 表单卡片

```vue
<GvCard>
  <template #header>
    <h2 class="text-xl font-semibold">新建资产</h2>
  </template>
  
  <el-form :model="form">
    <el-form-item label="名称">
      <el-input v-model="form.name" />
    </el-form-item>
  </el-form>
  
  <template #footer>
    <div class="flex justify-end gap-2">
      <GvButton variant="tonal">取消</GvButton>
      <GvButton variant="filled">提交</GvButton>
    </div>
  </template>
</GvCard>
```

## 最佳实践

1. **阴影使用**：页面主要内容使用 md，浮动元素使用 lg/xl
2. **内边距**：根据内容密度选择合适的 padding
3. **悬停效果**：列表页卡片建议使用 hoverable
4. **背景色**：深色背景上使用 container 背景色增强对比度
5. **深色模式**：组件自动适配深色模式，无需额外配置
