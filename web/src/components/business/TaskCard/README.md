# TaskCard - 任务卡片组件

业务组件，用于展示任务信息的卡片，基于 `GvCard` 封装。

## 基本用法

```vue
<template>
  <TaskCard
    :task="task"
    @click="handleClick"
  />
</template>

<script setup>
import { TaskCard } from '@/components'

const task = {
  id: 1,
  name: '视频转码任务',
  type: '转码',
  status: 'processing',
  description: '将视频转码为 MP4 格式',
  progress: 65,
  startTime: '2024-01-01 12:00',
  endTime: '2024-01-01 13:00'
}

const handleClick = (task) => {
  console.log('点击任务:', task)
}
</script>
```

## 网格布局

```vue
<GvGrid :cols="2" gap="lg">
  <TaskCard
    v-for="task in tasks"
    :key="task.id"
    :task="task"
    @click="handleClick"
    @cancel="handleCancel"
  />
</GvGrid>
```

## 垂直列表

```vue
<GvSpace vertical size="lg">
  <TaskCard
    v-for="task in tasks"
    :key="task.id"
    :task="task"
    @view="handleView"
    @cancel="handleCancel"
  />
</GvSpace>
```

## 不显示进度条

```vue
<TaskCard
  :task="task"
  :show-progress="false"
/>
```

## 自定义操作按钮

```vue
<TaskCard
  :task="task"
  :show-cancel="false"
  @view="handleView"
/>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| task | `Task` | - | 任务数据（必填） |
| showProgress | `boolean` | `true` | 是否显示进度条 |
| showActions | `boolean` | `true` | 是否显示操作按钮 |
| showView | `boolean` | `true` | 是否显示查看按钮 |
| showCancel | `boolean` | `true` | 是否显示取消按钮 |

## Task

| 属性 | 类型 | 说明 |
|------|------|------|
| id | `string \| number` | ID（必填） |
| name | `string` | 名称（必填） |
| status | `string` | 状态（必填） |
| type | `string` | 类型 |
| description | `string` | 描述 |
| progress | `number` | 进度（0-100） |
| startTime | `string` | 开始时间 |
| endTime | `string` | 结束时间 |
| createdAt | `string` | 创建时间 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| click | `(task)` | 点击卡片 |
| view | `(task)` | 点击查看 |
| cancel | `(task)` | 点击取消 |

## 使用场景

### 任务列表

```vue
<template>
  <GvContainer>
    <PageHeader title="任务管理">
      <template #extra>
        <GvSpace>
          <StatusBadge status="processing" text="3 个任务运行中" />
          <StatusBadge status="pending" text="5 个任务等待中" />
        </GvSpace>
      </template>
    </PageHeader>
    
    <GvGrid :cols="2" gap="lg">
      <TaskCard
        v-for="task in tasks"
        :key="task.id"
        :task="task"
        @click="handleView"
        @cancel="handleCancel"
      />
    </GvGrid>
  </GvContainer>
</template>
```

### 按状态分组

```vue
<template>
  <div>
    <!-- 进行中 -->
    <h3 class="text-lg font-semibold mb-4">进行中</h3>
    <GvSpace vertical size="lg" class="mb-6">
      <TaskCard
        v-for="task in processingTasks"
        :key="task.id"
        :task="task"
        @cancel="handleCancel"
      />
    </GvSpace>
    
    <GvDivider />
    
    <!-- 已完成 -->
    <h3 class="text-lg font-semibold mb-4">已完成</h3>
    <GvSpace vertical size="lg">
      <TaskCard
        v-for="task in completedTasks"
        :key="task.id"
        :task="task"
        :show-cancel="false"
      />
    </GvSpace>
  </div>
</template>

<script setup>
const processingTasks = computed(() => 
  tasks.value.filter(t => t.status === 'processing')
)

const completedTasks = computed(() => 
  tasks.value.filter(t => t.status === 'success')
)
</script>
```

### 侧边栏任务

```vue
<template>
  <GvDrawer v-model="visible" title="任务列表" size="small">
    <GvSpace vertical size="md">
      <TaskCard
        v-for="task in recentTasks"
        :key="task.id"
        :task="task"
        :show-actions="false"
        @click="handleViewTask"
      />
    </GvSpace>
  </GvDrawer>
</template>
```

## 最佳实践

1. **状态显示**：
   - 使用 StatusBadge 统一状态
   - 不同状态不同颜色

2. **进度显示**：
   - 处理中显示进度
   - 已完成不显示进度

3. **取消按钮**：
   - 仅在可取消状态显示
   - 确认后再取消

4. **时间显示**：
   - 显示开始时间
   - 已完成显示结束时间

5. **布局选择**：
   - 网格布局: 2-3 列
   - 垂直列表: 详细信息多时使用
