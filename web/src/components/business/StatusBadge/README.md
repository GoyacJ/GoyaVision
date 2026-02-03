# StatusBadge - 状态徽章组件

业务组件，用于显示统一的状态徽章，基于 `GvBadge` 封装。

## 基本用法

```vue
<template>
  <StatusBadge status="running" />
  <StatusBadge status="stopped" />
  <StatusBadge status="success" />
  <StatusBadge status="failed" />
</template>

<script setup>
import { StatusBadge } from '@/components'
</script>
```

## 支持的状态

### 运行状态

```vue
<StatusBadge status="running" />    <!-- 运行中 -->
<StatusBadge status="stopped" />    <!-- 已停止 -->
```

### 任务状态

```vue
<StatusBadge status="pending" />     <!-- 待处理 -->
<StatusBadge status="processing" />  <!-- 处理中 -->
<StatusBadge status="success" />     <!-- 成功 -->
<StatusBadge status="failed" />      <!-- 失败 -->
```

### 系统状态

```vue
<StatusBadge status="error" />      <!-- 错误 -->
<StatusBadge status="warning" />    <!-- 警告 -->
```

### 激活状态

```vue
<StatusBadge status="active" />     <!-- 激活 -->
<StatusBadge status="inactive" />   <!-- 未激活 -->
```

### 在线状态

```vue
<StatusBadge status="online" />     <!-- 在线 -->
<StatusBadge status="offline" />    <!-- 离线 -->
```

### 启用状态

```vue
<StatusBadge status="enabled" />    <!-- 启用 -->
<StatusBadge status="disabled" />   <!-- 禁用 -->
```

## 自定义文本

```vue
<StatusBadge status="running" text="正在录制" />
<StatusBadge status="success" text="已完成" />
```

## 禁用动画

```vue
<StatusBadge status="running" :animated="false" />
```

## 隐藏图标

```vue
<StatusBadge status="success" :show-icon="false" />
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| status | `StatusType` | - | 状态类型（必填） |
| text | `string` | - | 自定义显示文本 |
| animated | `boolean` | `true` | 是否显示动画 |
| showIcon | `boolean` | `true` | 是否显示图标 |

## StatusType

| 状态 | 默认文本 | 颜色 | 动画 |
|------|----------|------|------|
| running | 运行中 | success | ✅ |
| stopped | 已停止 | neutral | - |
| pending | 待处理 | warning | - |
| processing | 处理中 | info | ✅ |
| success | 成功 | success | - |
| failed | 失败 | error | - |
| error | 错误 | error | - |
| warning | 警告 | warning | - |
| active | 激活 | success | ✅ |
| inactive | 未激活 | neutral | - |
| online | 在线 | success | ✅ |
| offline | 离线 | neutral | - |
| enabled | 启用 | success | - |
| disabled | 禁用 | neutral | - |

## 使用场景

### 流状态

```vue
<template>
  <GvTable :data="streams" :columns="columns">
    <template #status="{ row }">
      <StatusBadge :status="row.status" />
    </template>
  </GvTable>
</template>
```

### 任务列表

```vue
<template>
  <GvSpace vertical>
    <GvCard v-for="task in tasks" :key="task.id">
      <GvFlex justify="between" align="center">
        <h3>{{ task.name }}</h3>
        <StatusBadge :status="task.status" />
      </GvFlex>
    </GvCard>
  </GvSpace>
</template>
```

### 用户在线状态

```vue
<template>
  <GvFlex align="center" gap="sm">
    <img :src="user.avatar" class="w-10 h-10 rounded-full" />
    <div>
      <div>{{ user.name }}</div>
      <StatusBadge :status="user.onlineStatus" />
    </div>
  </GvFlex>
</template>
```

## 最佳实践

1. **使用默认文本**：
   - 优先使用默认文本保持一致性
   - 特殊场景才自定义文本

2. **动画使用**：
   - 运行中、处理中等状态保持动画
   - 静态状态可禁用动画

3. **状态映射**：
   - 后端状态统一映射到标准状态
   - 避免过多自定义状态

4. **颜色语义**：
   - success: 正常、成功、在线
   - error: 失败、错误
   - warning: 警告、待处理
   - neutral: 停止、禁用、离线
