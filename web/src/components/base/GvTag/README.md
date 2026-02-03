# GvTag - 标签组件

Material Design 3 风格的标签组件，用于分类、标记、状态展示等。

## 基本用法

```vue
<GvTag>标签文字</GvTag>
```

## 颜色主题

```vue
<GvTag color="primary">Primary</GvTag>
<GvTag color="success">Success</GvTag>
<GvTag color="error">Error</GvTag>
<GvTag color="warning">Warning</GvTag>
<GvTag color="neutral">Neutral</GvTag>
```

## 变体

```vue
<GvTag variant="filled">Filled</GvTag>
<GvTag variant="tonal">Tonal</GvTag>
<GvTag variant="outlined">Outlined</GvTag>
```

## 尺寸

```vue
<GvTag size="small">Small</GvTag>
<GvTag size="medium">Medium</GvTag>
<GvTag size="large">Large</GvTag>
```

## 可关闭标签

```vue
<GvTag closable @close="handleClose">
  可关闭标签
</GvTag>
```

## 带图标

```vue
<GvTag icon="Check" color="success">已完成</GvTag>
<GvTag icon="Clock" color="warning">进行中</GvTag>
<GvTag icon="Close" color="error">失败</GvTag>
```

## 圆形标签

```vue
<GvTag rounded>圆形标签</GvTag>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| color | `'primary' \| 'secondary' \| 'success' \| 'error' \| 'warning' \| 'info' \| 'neutral'` | `'neutral'` | 颜色主题 |
| size | `'small' \| 'medium' \| 'large'` | `'medium'` | 尺寸 |
| variant | `'filled' \| 'tonal' \| 'outlined'` | `'tonal'` | 变体 |
| closable | `boolean` | `false` | 是否可关闭 |
| rounded | `boolean` | `false` | 是否圆形 |
| icon | `string` | - | 前置图标 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| click | `(event: MouseEvent)` | 点击事件 |
| close | `(event: MouseEvent)` | 关闭事件 |

## Slots

| 插槽名 | 说明 |
|--------|------|
| default | 标签内容 |

## 使用场景

### 状态标签

```vue
<GvTag color="success" icon="Check">已完成</GvTag>
<GvTag color="warning" icon="Clock">处理中</GvTag>
<GvTag color="error" icon="Close">失败</GvTag>
```

### 分类标签

```vue
<GvTag variant="tonal" closable>视频</GvTag>
<GvTag variant="tonal" closable>图片</GvTag>
<GvTag variant="tonal" closable>音频</GvTag>
```

### 标签组

```vue
<div class="flex flex-wrap gap-2">
  <GvTag v-for="tag in tags" :key="tag" closable @close="removeTag(tag)">
    {{ tag }}
  </GvTag>
</div>
```
