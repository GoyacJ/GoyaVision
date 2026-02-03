<template>
  <GvBadge
    :variant="badgeVariant"
    :color="badgeColor"
    :class="badgeClasses"
  >
    <GvFlex align="center" gap="xs">
      <!-- 状态图标/动画 -->
      <span v-if="showIcon" :class="iconClasses">
        <span v-if="shouldAnimate" class="status-badge__dot"></span>
        <el-icon v-else>
          <component :is="statusIcon" />
        </el-icon>
      </span>
      
      <!-- 状态文本 -->
      <span>{{ displayText }}</span>
    </GvFlex>
  </GvBadge>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'
import { GvBadge, GvFlex } from '@/components'
import type { StatusBadgeProps } from './types'
import {
  CircleCheck,
  CircleClose,
  Warning,
  VideoPlay,
  VideoPause,
  Clock,
  Loading
} from '@element-plus/icons-vue'

const props = withDefaults(defineProps<StatusBadgeProps>(), {
  animated: true,
  showIcon: true
})

// 状态配置映射
const statusConfig = {
  running: {
    text: '运行中',
    color: 'success' as const,
    variant: 'filled' as const,
    icon: VideoPlay,
    animated: true
  },
  stopped: {
    text: '已停止',
    color: 'neutral' as const,
    variant: 'tonal' as const,
    icon: VideoPause,
    animated: false
  },
  pending: {
    text: '待处理',
    color: 'warning' as const,
    variant: 'tonal' as const,
    icon: Clock,
    animated: false
  },
  processing: {
    text: '处理中',
    color: 'info' as const,
    variant: 'filled' as const,
    icon: Loading,
    animated: true
  },
  success: {
    text: '成功',
    color: 'success' as const,
    variant: 'tonal' as const,
    icon: CircleCheck,
    animated: false
  },
  failed: {
    text: '失败',
    color: 'error' as const,
    variant: 'tonal' as const,
    icon: CircleClose,
    animated: false
  },
  error: {
    text: '错误',
    color: 'error' as const,
    variant: 'filled' as const,
    icon: CircleClose,
    animated: false
  },
  warning: {
    text: '警告',
    color: 'warning' as const,
    variant: 'filled' as const,
    icon: Warning,
    animated: false
  },
  active: {
    text: '激活',
    color: 'success' as const,
    variant: 'filled' as const,
    icon: CircleCheck,
    animated: true
  },
  inactive: {
    text: '未激活',
    color: 'neutral' as const,
    variant: 'tonal' as const,
    icon: CircleClose,
    animated: false
  },
  online: {
    text: '在线',
    color: 'success' as const,
    variant: 'filled' as const,
    icon: CircleCheck,
    animated: true
  },
  offline: {
    text: '离线',
    color: 'neutral' as const,
    variant: 'tonal' as const,
    icon: CircleClose,
    animated: false
  },
  enabled: {
    text: '启用',
    color: 'success' as const,
    variant: 'tonal' as const,
    icon: CircleCheck,
    animated: false
  },
  disabled: {
    text: '禁用',
    color: 'neutral' as const,
    variant: 'tonal' as const,
    icon: CircleClose,
    animated: false
  }
}

// 当前状态配置
const currentConfig = computed(() => statusConfig[props.status])

// 显示文本
const displayText = computed(() => props.text || currentConfig.value.text)

// 徽章颜色
const badgeColor = computed(() => currentConfig.value.color)

// 徽章变体
const badgeVariant = computed(() => currentConfig.value.variant)

// 状态图标
const statusIcon = computed(() => currentConfig.value.icon)

// 是否应该显示动画
const shouldAnimate = computed(() => {
  return props.animated && currentConfig.value.animated
})

// 徽章类名
const badgeClasses = computed(() => {
  return cn('status-badge', shouldAnimate.value && 'status-badge--animated')
})

// 图标类名
const iconClasses = computed(() => {
  return cn(
    'status-badge__icon',
    'inline-flex items-center justify-center'
  )
})
</script>

<style scoped>
/* 跳动的点动画 */
.status-badge__dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: currentColor;
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

/* 旋转动画（用于 Loading 图标） */
.status-badge--animated :deep(.el-icon) {
  animation: rotate 1s linear infinite;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
