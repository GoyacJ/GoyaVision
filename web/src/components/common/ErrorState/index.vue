<template>
  <div class="error-state">
    <div class="error-icon">
      <svg
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <circle cx="12" cy="12" r="10" />
        <line x1="12" y1="8" x2="12" y2="12" />
        <line x1="12" y1="16" x2="12.01" y2="16" />
      </svg>
    </div>

    <h3 class="error-title">{{ displayTitle }}</h3>

    <p v-if="displayMessage" class="error-message">
      {{ displayMessage }}
    </p>

    <button
      v-if="showRetry"
      class="retry-button"
      @click="handleRetry"
    >
      {{ retryText }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ErrorStateProps, ErrorStateEmits } from './types'

const props = withDefaults(defineProps<ErrorStateProps>(), {
  title: '加载失败',
  retryText: '重试',
  showRetry: true
})

const emit = defineEmits<ErrorStateEmits>()

// 显示的标题
const displayTitle = computed(() => props.title)

// 显示的消息
const displayMessage = computed(() => {
  if (props.message) {
    return props.message
  }
  if (props.error) {
    return props.error.message || '发生未知错误'
  }
  return '请稍后重试'
})

// 重试处理
const handleRetry = () => {
  emit('retry')
}
</script>

<style scoped>
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  padding: 32px;
  text-align: center;
}

/* 错误图标 */
.error-icon {
  width: 64px;
  height: 64px;
  margin-bottom: 24px;
  color: #EF4444;  /* error.500 */
}

.error-icon svg {
  width: 100%;
  height: 100%;
}

/* 错误标题 */
.error-title {
  font-size: 1.125rem;  /* 18px */
  font-weight: 600;
  color: #262626;  /* neutral.800 */
  margin-bottom: 8px;
  line-height: 1.4;
}

/* 错误消息 */
.error-message {
  font-size: 0.875rem;  /* 14px */
  color: #737373;  /* neutral.500 */
  margin-bottom: 24px;
  max-width: 400px;
  line-height: 1.5;
}

/* 重试按钮 */
.retry-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 8px 24px;
  background-color: #4F5B93;  /* primary.600 */
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 0.875rem;  /* 14px */
  font-weight: 500;
  cursor: pointer;
  transition: background-color 150ms;
}

.retry-button:hover {
  background-color: #3E4A7A;  /* primary.700 */
}

.retry-button:active {
  background-color: #2F3A61;  /* primary.800 */
}

.retry-button:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(79, 91, 147, 0.2);
}
</style>
