<template>
  <transition name="gv-alert-fade">
    <div v-show="visible" :class="alertClasses">
      <!-- 图标 -->
      <div v-if="showIcon" class="gv-alert__icon">
        <el-icon :size="18">
          <SuccessFilled v-if="type === 'success'" />
          <InfoFilled v-if="type === 'info'" />
          <WarningFilled v-if="type === 'warning'" />
          <CircleCloseFilled v-if="type === 'error'" />
        </el-icon>
      </div>
      
      <!-- 内容 -->
      <div class="gv-alert__content">
        <!-- 标题 -->
        <div v-if="title || $slots.title" class="gv-alert__title">
          <slot name="title">{{ title }}</slot>
        </div>
        
        <!-- 描述 -->
        <div v-if="description || $slots.default" class="gv-alert__description">
          <slot>{{ description }}</slot>
        </div>
      </div>
      
      <!-- 关闭按钮 -->
      <div v-if="closable" class="gv-alert__close" @click="handleClose">
        <span v-if="closeText" class="gv-alert__close-text">
          {{ closeText }}
        </span>
        <el-icon v-else :size="16">
          <Close />
        </el-icon>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { cn } from '@/utils/cn'
import type { AlertProps, AlertEmits } from './types'

const props = withDefaults(defineProps<AlertProps>(), {
  type: 'info',
  closable: false,
  showIcon: true,
  center: false
})

const emit = defineEmits<AlertEmits>()

const visible = ref(true)

// 警告框类名
const alertClasses = computed(() => {
  const base = [
    'gv-alert',
    'relative flex items-start gap-3',
    'px-4 py-3 rounded-xl',
    'transition-all duration-300'
  ]
  
  // 类型样式
  const typeClasses = {
    success: 'bg-success-50 border border-success-200 text-success-800',
    info: 'bg-info-50 border border-info-200 text-info-800',
    warning: 'bg-warning-50 border border-warning-200 text-warning-800',
    error: 'bg-error-50 border border-error-200 text-error-800'
  }
  
  // 居中
  const centerClass = props.center ? 'justify-center text-center' : ''
  
  return cn(base, typeClasses[props.type], centerClass)
})

// 关闭事件
const handleClose = () => {
  visible.value = false
  emit('close')
}
</script>

<style scoped>
.gv-alert__icon {
  flex-shrink: 0;
  margin-top: 1px;
}

.gv-alert__icon .el-icon.success {
  color: rgb(16 185 129);
}

.gv-alert__icon .el-icon.info {
  color: rgb(59 130 246);
}

.gv-alert__icon .el-icon.warning {
  color: rgb(245 158 11);
}

.gv-alert__icon .el-icon.error {
  color: rgb(239 68 68);
}

.gv-alert__content {
  flex: 1;
  min-width: 0;
}

.gv-alert__title {
  font-size: 0.875rem;
  font-weight: 600;
  line-height: 1.5;
  margin-bottom: 4px;
}

.gv-alert__description {
  font-size: 0.875rem;
  line-height: 1.5;
  opacity: 0.9;
}

.gv-alert__close {
  flex-shrink: 0;
  margin-left: auto;
  cursor: pointer;
  opacity: 0.7;
  transition: opacity 0.2s;
}

.gv-alert__close:hover {
  opacity: 1;
}

.gv-alert__close-text {
  font-size: 0.875rem;
  font-weight: 500;
}

/* 过渡动画 */
.gv-alert-fade-enter-active,
.gv-alert-fade-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.gv-alert-fade-enter-from {
  opacity: 0;
  transform: translateY(-8px);
}

.gv-alert-fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* 深色模式 */
.dark .gv-alert[class*="success"] {
  @apply bg-success-950/30 border-success-800/50 text-success-200;
}

.dark .gv-alert[class*="info"] {
  @apply bg-info-950/30 border-info-800/50 text-info-200;
}

.dark .gv-alert[class*="warning"] {
  @apply bg-warning-950/30 border-warning-800/50 text-warning-200;
}

.dark .gv-alert[class*="error"] {
  @apply bg-error-950/30 border-error-800/50 text-error-200;
}
</style>
