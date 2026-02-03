# GvDivider - 分割线组件

简洁的分割线组件，用于分隔内容区域。

## 基本用法

```vue
<template>
  <div>
    <p>内容 1</p>
    <GvDivider />
    <p>内容 2</p>
  </div>
</template>

<script setup>
import { GvDivider } from '@/components'
</script>
```

## 带文字

```vue
<GvDivider>分割线文字</GvDivider>

<GvDivider content-position="left">左侧文字</GvDivider>
<GvDivider content-position="center">居中文字</GvDivider>
<GvDivider content-position="right">右侧文字</GvDivider>
```

## 虚线

```vue
<GvDivider dashed />
<GvDivider dashed>虚线分割线</GvDivider>
```

## 垂直分割线

```vue
<div class="flex items-center">
  <span>选项 1</span>
  <GvDivider direction="vertical" />
  <span>选项 2</span>
  <GvDivider direction="vertical" />
  <span>选项 3</span>
</div>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| direction | `'horizontal' \| 'vertical'` | `'horizontal'` | 方向 |
| contentPosition | `'left' \| 'center' \| 'right'` | `'center'` | 文本位置 |
| dashed | `boolean` | `false` | 是否虚线 |
| customClass | `string` | - | 自定义类名 |

## Slots

| 插槽名 | 说明 |
|--------|------|
| default | 分割线文字 |

## 使用场景

### 内容分隔

```vue
<div>
  <section>
    <h2>章节 1</h2>
    <p>内容...</p>
  </section>
  
  <GvDivider />
  
  <section>
    <h2>章节 2</h2>
    <p>内容...</p>
  </section>
</div>
```

### 按钮组分隔

```vue
<GvFlex align="center">
  <GvButton>操作 1</GvButton>
  <GvDivider direction="vertical" />
  <GvButton>操作 2</GvButton>
  <GvDivider direction="vertical" />
  <GvButton>操作 3</GvButton>
</GvFlex>
```

### 列表分组

```vue
<div>
  <GvDivider content-position="left">今天</GvDivider>
  <div class="space-y-2">
    <div>项目 1</div>
    <div>项目 2</div>
  </div>
  
  <GvDivider content-position="left">昨天</GvDivider>
  <div class="space-y-2">
    <div>项目 3</div>
    <div>项目 4</div>
  </div>
</div>
```

## 最佳实践

1. **方向选择**：
   - horizontal: 垂直分隔内容
   - vertical: 水平分隔元素

2. **文字位置**：
   - center: 标题分隔
   - left: 列表分组
   - right: 较少使用

3. **虚线使用**：
   - 弱分隔使用虚线
   - 强分隔使用实线

4. **配合布局**：
   - 配合 GvFlex 实现垂直分割
   - 配合 section 实现章节分隔
