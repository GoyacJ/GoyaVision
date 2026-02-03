<template>
  <GvCard
    class="task-card cursor-pointer"
    hoverable
    @click="handleClick"
  >
    <!-- 头部：标题和状态 -->
    <GvFlex justify="between" align="start" class="mb-3">
      <div class="flex-1 min-w-0">
        <h3 class="task-card__title">{{ task.name }}</h3>
        <GvTag v-if="task.type" size="small" variant="tonal" class="mt-1">
          {{ task.type }}
        </GvTag>
      </div>
      
      <StatusBadge :status="task.status as any" />
    </GvFlex>
    
    <!-- 描述 -->
    <p v-if="task.description" class="task-card__description">
      {{ task.description }}
    </p>
    
    <!-- 进度条 -->
    <div v-if="showProgress && task.progress !== undefined" class="task-card__progress">
      <GvFlex justify="between" align="center" class="mb-2">
        <span class="text-sm text-text-secondary">进度</span>
        <span class="text-sm font-medium text-text-primary">
          {{ task.progress }}%
        </span>
      </GvFlex>
      <el-progress
        :percentage="task.progress"
        :status="getProgressStatus(task.progress)"
        :show-text="false"
      />
    </div>
    
    <!-- 时间信息 -->
    <GvFlex class="task-card__time" wrap gap="md">
      <div v-if="task.startTime" class="flex items-center gap-1">
        <el-icon class="text-text-tertiary"><Clock /></el-icon>
        <span class="text-xs text-text-tertiary">{{ task.startTime }}</span>
      </div>
      <div v-if="task.endTime" class="flex items-center gap-1">
        <el-icon class="text-text-tertiary"><Finished /></el-icon>
        <span class="text-xs text-text-tertiary">{{ task.endTime }}</span>
      </div>
    </GvFlex>
    
    <!-- 操作按钮 -->
    <GvDivider v-if="showActions" class="my-3" />
    <GvFlex v-if="showActions" justify="end">
      <GvSpace size="xs">
        <GvButton
          v-if="showView"
          size="small"
          variant="tonal"
          @click.stop="handleView"
        >
          查看
        </GvButton>
        <GvButton
          v-if="showCancel && canCancel"
          size="small"
          variant="text"
          @click.stop="handleCancel"
        >
          取消
        </GvButton>
      </GvSpace>
    </GvFlex>
  </GvCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  GvCard,
  GvFlex,
  GvTag,
  GvButton,
  GvSpace,
  GvDivider,
  StatusBadge
} from '@/components'
import { Clock, Finished } from '@element-plus/icons-vue'
import type { TaskCardProps, TaskCardEmits } from './types'

const props = withDefaults(defineProps<TaskCardProps>(), {
  showProgress: true,
  showActions: true,
  showView: true,
  showCancel: true
})

const emit = defineEmits<TaskCardEmits>()

// 是否可以取消
const canCancel = computed(() => {
  return ['pending', 'processing', 'running'].includes(props.task.status)
})

// 进度条状态
const getProgressStatus = (progress: number) => {
  if (progress === 100) return 'success'
  if (progress >= 80) return undefined
  if (progress >= 50) return undefined
  return undefined
}

// 点击卡片
const handleClick = () => {
  emit('click', props.task)
}

// 查看任务
const handleView = () => {
  emit('view', props.task)
}

// 取消任务
const handleCancel = () => {
  emit('cancel', props.task)
}
</script>

<style scoped>
.task-card__title {
  @apply text-base font-semibold text-text-primary;
  @apply m-0 truncate;
}

.task-card__description {
  @apply text-sm text-text-secondary;
  @apply m-0 mb-3 line-clamp-2;
}

.task-card__progress {
  @apply mb-3;
}

.task-card__time {
  @apply mt-3;
}

/* 深色模式 */
.dark .task-card__title {
  @apply text-text-inverse;
}

.dark .task-card__description {
  @apply text-neutral-400;
}
</style>
