# GvBadge - 徽章组件

Material Design 3 风格的徽章组件，用于显示状态、通知数量等信息。

## 基本用法

### 独立徽章

```vue
<GvBadge>NEW</GvBadge>
```

### 角标徽章

```vue
<GvBadge :value="5">
  <GvButton icon="Notification" />
</GvBadge>
```

## 颜色主题

```vue
<GvBadge color="primary">Primary</GvBadge>
<GvBadge color="success">Success</GvBadge>
<GvBadge color="error">Error</GvBadge>
<GvBadge color="warning">Warning</GvBadge>
```

## 变体

```vue
<GvBadge variant="filled">Filled</GvBadge>
<GvBadge variant="tonal">Tonal</GvBadge>
<GvBadge variant="outlined">Outlined</GvBadge>
```

## 数字徽章

```vue
<GvBadge :value="5">
  <el-icon><Bell /></el-icon>
</GvBadge>

<GvBadge :value="100" :max="99">
  <el-icon><Message /></el-icon>
</GvBadge>
```

## 点状徽章

```vue
<GvBadge dot>
  <el-icon><Notification /></el-icon>
</GvBadge>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| color | `'primary' \| 'secondary' \| 'success' \| 'error' \| 'warning' \| 'info' \| 'neutral'` | `'primary'` | 颜色主题 |
| size | `'small' \| 'medium' \| 'large'` | `'medium'` | 尺寸 |
| variant | `'filled' \| 'tonal' \| 'outlined'` | `'filled'` | 变体 |
| value | `number \| string` | - | 显示的值 |
| max | `number` | `99` | 最大值 |
| dot | `boolean` | `false` | 是否为点状 |
| hidden | `boolean` | `false` | 是否隐藏 |
| offset | `[number, number]` | - | 偏移量 |

## Slots

| 插槽名 | 说明 |
|--------|------|
| default | 徽章附着的元素 |

## 使用场景

### 状态标识

```vue
<GvBadge color="success">运行中</GvBadge>
<GvBadge color="warning">待处理</GvBadge>
<GvBadge color="error">失败</GvBadge>
```

### 通知数量

```vue
<GvBadge :value="unreadCount">
  <el-icon><Bell /></el-icon>
</GvBadge>
```

### 在线状态

```vue
<GvBadge dot color="success">
  <el-avatar :src="avatar" />
</GvBadge>
```
