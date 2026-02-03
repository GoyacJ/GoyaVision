# GvAlert - 警告框组件

Material Design 3 风格的警告框组件，用于显示重要的提示信息。

## 基本用法

```vue
<template>
  <GvAlert
    title="这是一条提示信息"
    description="这是详细的描述文本"
  />
</template>

<script setup>
import { GvAlert } from '@/components'
</script>
```

## 警告框类型

```vue
<GvAlert type="success" title="成功提示" description="操作成功完成" />
<GvAlert type="info" title="信息提示" description="这是一条普通信息" />
<GvAlert type="warning" title="警告提示" description="请注意相关风险" />
<GvAlert type="error" title="错误提示" description="操作失败，请重试" />
```

## 可关闭

```vue
<GvAlert
  type="warning"
  title="可关闭的警告框"
  closable
  @close="handleClose"
/>
```

## 自定义关闭文本

```vue
<GvAlert
  type="info"
  title="自定义关闭按钮"
  closable
  close-text="知道了"
/>
```

## 不显示图标

```vue
<GvAlert
  title="不显示图标"
  :show-icon="false"
/>
```

## 居中显示

```vue
<GvAlert
  type="success"
  title="居中显示的警告框"
  center
/>
```

## 使用插槽

```vue
<GvAlert type="info">
  <template #title>
    <strong>自定义标题</strong>
  </template>
  <div>
    <p>这是自定义的内容</p>
    <ul>
      <li>列表项 1</li>
      <li>列表项 2</li>
    </ul>
  </div>
</GvAlert>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| type | `'success' \| 'info' \| 'warning' \| 'error'` | `'info'` | 警告框类型 |
| title | `string` | - | 标题 |
| description | `string` | - | 描述文本 |
| closable | `boolean` | `false` | 是否可关闭 |
| showIcon | `boolean` | `true` | 是否显示图标 |
| center | `boolean` | `false` | 是否居中显示 |
| closeText | `string` | - | 关闭按钮文本 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| close | `()` | 关闭事件 |

## Slots

| 插槽名 | 说明 |
|--------|------|
| title | 标题内容 |
| default | 描述内容 |

## 使用场景

### 操作成功提示

```vue
<GvAlert
  type="success"
  title="操作成功"
  description="您的更改已保存"
  closable
/>
```

### 重要提醒

```vue
<GvAlert
  type="warning"
  title="注意"
  description="此操作不可撤销，请谨慎操作"
  show-icon
/>
```

### 错误提示

```vue
<GvAlert
  type="error"
  title="保存失败"
  description="网络连接失败，请检查您的网络设置"
  closable
/>
```

### 通知信息

```vue
<GvAlert
  type="info"
  title="系统维护通知"
  description="系统将于今晚 22:00 进行例行维护，预计 1 小时"
/>
```

### 带操作的警告框

```vue
<GvAlert type="warning" closable>
  <template #title>
    <strong>版本更新</strong>
  </template>
  <div>
    <p class="mb-2">发现新版本 v2.0.0，建议立即更新</p>
    <GvButton size="small" variant="tonal">立即更新</GvButton>
  </div>
</GvAlert>
```

## 最佳实践

1. **类型选择**：根据信息重要性选择合适的类型
   - success: 操作成功
   - info: 普通提示
   - warning: 需要注意
   - error: 错误或失败

2. **可关闭性**：
   - 临时提示建议可关闭
   - 重要警告建议不可关闭

3. **图标使用**：
   - 默认显示图标增强识别度
   - 简单提示可隐藏图标

4. **文案**：
   - 标题简洁明了
   - 描述提供必要细节

5. **位置**：
   - 页面级提示放在内容顶部
   - 局部提示放在相关模块内
