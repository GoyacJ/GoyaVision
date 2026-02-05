<template>
  <div class="empty-state">
    <div class="empty-icon">
      {{ icon }}
    </div>

    <h3 class="empty-title">{{ title }}</h3>

    <p v-if="description" class="empty-description">
      {{ description }}
    </p>

    <button
      v-if="showAction && actionText"
      class="action-button"
      @click="handleAction"
    >
      {{ actionText }}
    </button>
  </div>
</template>

<script setup lang="ts">
import type { EmptyStateProps, EmptyStateEmits } from './types'

const props = withDefaults(defineProps<EmptyStateProps>(), {
  icon: '📭',
  title: '暂无数据',
  showAction: false
})

const emit = defineEmits<EmptyStateEmits>()

// 操作处理
const handleAction = () => {
  emit('action')
}
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  padding: 32px;
  text-align: center;
}

/* 空状态图标 */
.empty-icon {
  font-size: 64px;
  line-height: 1;
  margin-bottom: 24px;
  opacity: 0.5;
}

/* 空状态标题 */
.empty-title {
  font-size: 1rem;  /* 16px */
  font-weight: 500;
  color: #525252;  /* neutral.600 */
  margin-bottom: 8px;
  line-height: 1.5;
}

/* 空状态描述 */
.empty-description {
  font-size: 0.875rem;  /* 14px */
  color: #737373;  /* neutral.500 */
  margin-bottom: 24px;
  max-width: 400px;
  line-height: 1.5;
}

/* 操作按钮 */
.action-button {
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

.action-button:hover {
  background-color: #3E4A7A;  /* primary.700 */
}

.action-button:active {
  background-color: #2F3A61;  /* primary.800 */
}

.action-button:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(79, 91, 147, 0.2);
}
</style>
