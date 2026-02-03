# GvSpace - 间距组件

简洁的间距组件，用于给子元素添加统一间距。

## 基本用法

```vue
<template>
  <GvSpace>
    <GvButton>按钮 1</GvButton>
    <GvButton>按钮 2</GvButton>
    <GvButton>按钮 3</GvButton>
  </GvSpace>
</template>

<script setup>
import { GvSpace, GvButton } from '@/components'
</script>
```

## 间距大小

```vue
<GvSpace size="xs">8px 间距</GvSpace>
<GvSpace size="sm">12px 间距</GvSpace>
<GvSpace size="md">16px 间距</GvSpace>
<GvSpace size="lg">24px 间距</GvSpace>
<GvSpace size="xl">32px 间距</GvSpace>

<!-- 自定义数值 -->
<GvSpace :size="20">20px 间距</GvSpace>
```

## 不同方向间距

```vue
<!-- 水平 12px，垂直 24px -->
<GvSpace :size="['sm', 'lg']" wrap>
  <GvButton v-for="i in 10" :key="i">
    按钮 {{ i }}
  </GvButton>
</GvSpace>
```

## 垂直间距

```vue
<GvSpace direction="vertical">
  <GvCard>卡片 1</GvCard>
  <GvCard>卡片 2</GvCard>
  <GvCard>卡片 3</GvCard>
</GvSpace>
```

## 对齐方式

```vue
<GvSpace align="start">顶部对齐</GvSpace>
<GvSpace align="center">居中对齐</GvSpace>
<GvSpace align="end">底部对齐</GvSpace>
<GvSpace align="baseline">基线对齐</GvSpace>
```

## 自动换行

```vue
<GvSpace wrap>
  <GvTag v-for="i in 20" :key="i">
    标签 {{ i }}
  </GvTag>
</GvSpace>
```

## 填充父容器

```vue
<GvSpace fill>
  <GvButton>按钮 1</GvButton>
  <GvButton>按钮 2</GvButton>
  <GvButton>按钮 3</GvButton>
</GvSpace>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| size | `SpaceSize \| [SpaceSize, SpaceSize]` | `'md'` | 间距大小 |
| direction | `'horizontal' \| 'vertical'` | `'horizontal'` | 方向 |
| align | `'start' \| 'center' \| 'end' \| 'baseline'` | `'center'` | 对齐方式 |
| wrap | `boolean` | `false` | 是否自动换行 |
| fill | `boolean` | `false` | 是否填充父容器 |
| fillRatio | `number` | - | 填充比例 |

## Slots

| 插槽名 | 说明 |
|--------|------|
| default | 子元素 |

## 使用场景

### 按钮组

```vue
<GvSpace>
  <GvButton variant="filled">保存</GvButton>
  <GvButton variant="tonal">取消</GvButton>
  <GvButton variant="text">重置</GvButton>
</GvSpace>
```

### 标签列表

```vue
<GvSpace wrap size="sm">
  <GvTag v-for="tag in tags" :key="tag" closable>
    {{ tag }}
  </GvTag>
</GvSpace>
```

### 表单操作栏

```vue
<GvSpace fill>
  <div class="flex-1"></div>
  <GvButton variant="tonal">取消</GvButton>
  <GvButton variant="filled">提交</GvButton>
</GvSpace>
```

### 垂直列表

```vue
<GvSpace direction="vertical" size="lg">
  <GvCard v-for="item in list" :key="item.id">
    {{ item.title }}
  </GvCard>
</GvSpace>
```

## 最佳实践

1. **vs GvFlex**：
   - GvSpace: 简单间距，无需对齐控制
   - GvFlex: 复杂布局，需要精确控制

2. **间距选择**：
   - xs/sm: 紧凑布局
   - md: 默认推荐
   - lg/xl: 宽松布局

3. **方向**：
   - horizontal: 按钮组、标签列表
   - vertical: 卡片列表、表单项

4. **换行**：
   - 标签列表建议 wrap
   - 按钮组通常不换行
