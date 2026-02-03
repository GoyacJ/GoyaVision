<template>
  <div :class="cardClasses" @click="handleClick">
    <!-- 卡片头部 -->
    <div v-if="$slots.header" :class="headerClasses">
      <slot name="header" />
    </div>
    
    <!-- 卡片主体 -->
    <div :class="bodyClasses">
      <slot />
    </div>
    
    <!-- 卡片底部 -->
    <div v-if="$slots.footer" :class="footerClasses">
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'
import type { CardProps, CardEmits } from './types'

const props = withDefaults(defineProps<CardProps>(), {
  shadow: 'md',
  padding: 'md',
  hoverable: false,
  bordered: false,
  background: 'default'
})

const emit = defineEmits<CardEmits>()

// 卡片容器类名
const cardClasses = computed(() => {
  const base = [
    'gv-card',
    'rounded-xl overflow-hidden',
    'transition-all duration-300'
  ]
  
  // 阴影
  const shadowClasses = {
    none: '',
    sm: 'shadow-sm',
    md: 'shadow-md',
    lg: 'shadow-lg',
    xl: 'shadow-xl'
  }
  
  // 背景色
  const backgroundClasses = {
    default: 'bg-white',
    container: 'bg-surface-container',
    transparent: 'bg-transparent'
  }
  
  // 边框
  const borderClass = props.bordered ? 'border border-neutral-200' : ''
  
  // 悬停效果
  const hoverClass = props.hoverable
    ? 'hover:shadow-lg hover:-translate-y-1 cursor-pointer'
    : ''
  
  return cn(
    base,
    shadowClasses[props.shadow],
    backgroundClasses[props.background],
    borderClass,
    hoverClass
  )
})

// 头部类名
const headerClasses = computed(() => {
  const paddingClasses = {
    none: '',
    sm: 'px-4 py-3',
    md: 'px-6 py-4',
    lg: 'px-8 py-5'
  }
  
  return cn(
    'gv-card__header',
    'border-b border-neutral-100',
    'bg-gradient-to-r from-primary-50/30 to-secondary-50/30',
    paddingClasses[props.padding]
  )
})

// 主体类名
const bodyClasses = computed(() => {
  const paddingClasses = {
    none: '',
    sm: 'p-4',
    md: 'p-6',
    lg: 'p-8'
  }
  
  return cn('gv-card__body', paddingClasses[props.padding])
})

// 底部类名
const footerClasses = computed(() => {
  const paddingClasses = {
    none: '',
    sm: 'px-4 py-3',
    md: 'px-6 py-4',
    lg: 'px-8 py-5'
  }
  
  return cn(
    'gv-card__footer',
    'border-t border-neutral-100',
    'bg-neutral-50/50',
    paddingClasses[props.padding]
  )
})

// 点击事件处理
const handleClick = (event: MouseEvent) => {
  if (props.hoverable) {
    emit('click', event)
  }
}
</script>

<style scoped>
.gv-card {
  backdrop-filter: blur(10px);
}

.dark .gv-card {
  @apply bg-surface-dark border-neutral-700;
}

.dark .gv-card__header {
  @apply border-neutral-700 bg-gradient-to-r from-primary-950/30 to-secondary-950/30;
}

.dark .gv-card__footer {
  @apply border-neutral-700 bg-neutral-800/50;
}
</style>
