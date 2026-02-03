# GvButton - 按钮组件

Material Design 3 风格的按钮组件，支持多种变体和颜色主题。

## 基本用法

```vue
<template>
  <GvButton>默认按钮</GvButton>
</template>

<script setup>
import { GvButton } from '@/components'
</script>
```

## 变体（Variants）

按钮提供 4 种变体，从高到低强调级别：

### Filled（填充按钮）- 最高强调

用于最重要的操作。

```vue
<GvButton variant="filled" color="primary">主要操作</GvButton>
```

### Tonal（色调按钮）- 中等强调

用于次要但仍然重要的操作。

```vue
<GvButton variant="tonal" color="primary">次要操作</GvButton>
```

### Outlined（边框按钮）- 中等强调

用于与 Filled 按钮搭配使用的中等强调操作。

```vue
<GvButton variant="outlined" color="primary">其他操作</GvButton>
```

### Text（文本按钮）- 最低强调

用于最不重要的操作。

```vue
<GvButton variant="text" color="primary">辅助操作</GvButton>
```

## 颜色主题

支持 6 种颜色主题：

```vue
<GvButton color="primary">主要</GvButton>
<GvButton color="secondary">辅助</GvButton>
<GvButton color="success">成功</GvButton>
<GvButton color="error">错误</GvButton>
<GvButton color="warning">警告</GvButton>
<GvButton color="info">信息</GvButton>
```

## 尺寸

提供 3 种尺寸：

```vue
<GvButton size="small">小按钮</GvButton>
<GvButton size="medium">中按钮</GvButton>
<GvButton size="large">大按钮</GvButton>
```

## 图标

### 左侧图标

```vue
<GvButton icon="Plus" icon-position="left">新建</GvButton>
```

### 右侧图标

```vue
<GvButton icon="ArrowRight" icon-position="right">下一步</GvButton>
```

### 纯图标按钮

```vue
<GvButton icon="Search" rounded />
```

## 状态

### 加载状态

```vue
<GvButton :loading="isLoading">保存中...</GvButton>
```

### 禁用状态

```vue
<GvButton disabled>禁用按钮</GvButton>
```

## 块级按钮

```vue
<GvButton block>块级按钮</GvButton>
```

## 圆形按钮

```vue
<GvButton rounded>圆形</GvButton>
```

## 链接按钮

```vue
<GvButton href="https://example.com" target="_blank">
  外部链接
</GvButton>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| variant | `'filled' \| 'tonal' \| 'outlined' \| 'text'` | `'filled'` | 按钮变体 |
| color | `'primary' \| 'secondary' \| 'success' \| 'error' \| 'warning' \| 'info'` | `'primary'` | 颜色主题 |
| size | `'small' \| 'medium' \| 'large'` | `'medium'` | 按钮尺寸 |
| disabled | `boolean` | `false` | 是否禁用 |
| loading | `boolean` | `false` | 是否加载中 |
| icon | `string \| Component` | - | 图标 |
| iconPosition | `'left' \| 'right'` | `'left'` | 图标位置 |
| rounded | `boolean` | `false` | 是否圆形 |
| block | `boolean` | `false` | 是否块级 |
| type | `'button' \| 'submit' \| 'reset'` | `'button'` | HTML 类型 |
| href | `string` | - | 链接地址 |
| target | `string` | - | 链接打开方式 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| click | `(event: MouseEvent)` | 点击事件 |

## Slots

| 插槽名 | 说明 |
|--------|------|
| default | 按钮内容 |

## 使用场景

### 表单提交

```vue
<GvButton variant="filled" color="primary" type="submit">
  提交
</GvButton>
```

### 取消操作

```vue
<GvButton variant="tonal" color="secondary" @click="handleCancel">
  取消
</GvButton>
```

### 删除操作

```vue
<GvButton variant="outlined" color="error" @click="handleDelete">
  删除
</GvButton>
```

### 查看详情

```vue
<GvButton variant="text" size="small" @click="handleView">
  查看详情
</GvButton>
```

## 最佳实践

1. **强调层次**：在同一界面中，使用不同变体区分操作重要性
2. **颜色语义**：错误操作使用 error 颜色，成功操作使用 success 颜色
3. **图标使用**：重要操作可以添加图标增强识别度
4. **加载状态**：异步操作时显示 loading 状态提升用户体验
5. **响应式**：移动端可以使用 block 属性让按钮占满宽度
