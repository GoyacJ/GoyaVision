<template>
  <div :class="containerClasses">
    <div :class="spinnerClasses">
      <div class="spinner-ring" />
    </div>
    <p v-if="message" class="loading-message">{{ message }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { LoadingStateProps } from './types'

const props = withDefaults(defineProps<LoadingStateProps>(), {
  size: 'medium',
  fullscreen: false
})

// 容器类名
const containerClasses = computed(() => {
  const base = [
    'loading-state',
    'flex flex-col items-center justify-center gap-4'
  ]

  const fullscreenClass = props.fullscreen
    ? 'fixed inset-0 bg-surface-dim z-50'
    : 'min-h-[400px]'

  return [base, fullscreenClass].flat()
})

// 加载指示器类名
const spinnerClasses = computed(() => {
  const sizeClasses = {
    small: 'w-8 h-8',
    medium: 'w-12 h-12',
    large: 'w-16 h-16'
  }

  return ['spinner', sizeClasses[props.size]]
})
</script>

<style scoped>
/* 加载指示器 */
.spinner {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.spinner-ring {
  width: 100%;
  height: 100%;
  border: 3px solid #E5E5E5;  /* neutral.200 */
  border-top-color: #4F5B93;  /* primary.600 */
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* 提示文本 */
.loading-message {
  color: #737373;  /* neutral.500 */
  font-size: 0.875rem;  /* 14px */
  line-height: 1.5;
}
</style>
