<template>
  <span :class="tagClasses" @click="handleClick">
    <!-- 前置图标 -->
    <el-icon v-if="icon" class="gv-tag__icon">
      <component :is="icon" />
    </el-icon>
    
    <!-- 标签内容 -->
    <span class="gv-tag__content">
      <slot />
    </span>
    
    <!-- 关闭按钮 -->
    <el-icon
      v-if="closable"
      class="gv-tag__close"
      @click.stop="handleClose"
    >
      <Close />
    </el-icon>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'
import type { TagProps, TagEmits } from './types'

const props = withDefaults(defineProps<TagProps>(), {
  color: 'neutral',
  size: 'medium',
  variant: 'tonal',
  closable: false,
  rounded: false
})

const emit = defineEmits<TagEmits>()

// 标签类名
const tagClasses = computed(() => {
  const base = [
    'gv-tag',
    'inline-flex items-center justify-center gap-1',
    'font-medium transition-all duration-200',
    'cursor-default'
  ]
  
  // 尺寸
  const sizeClasses = {
    small: 'h-5 px-2 text-xs rounded',
    medium: 'h-6 px-2.5 text-sm rounded-md',
    large: 'h-7 px-3 text-base rounded-lg'
  }
  
  // 变体和颜色组合
  const variantColorClasses = {
    // Filled 变体
    'filled-primary': 'bg-primary-600 text-white',
    'filled-secondary': 'bg-neutral-600 text-white',
    'filled-success': 'bg-success-600 text-white',
    'filled-error': 'bg-error-600 text-white',
    'filled-warning': 'bg-warning-600 text-white',
    'filled-info': 'bg-info-600 text-white',
    'filled-neutral': 'bg-neutral-600 text-white',
    
    // Tonal 变体（默认推荐）
    'tonal-primary': 'bg-primary-100 text-primary-700',
    'tonal-secondary': 'bg-neutral-100 text-neutral-700',
    'tonal-success': 'bg-success-100 text-success-700',
    'tonal-error': 'bg-error-100 text-error-700',
    'tonal-warning': 'bg-warning-100 text-warning-700',
    'tonal-info': 'bg-info-100 text-info-700',
    'tonal-neutral': 'bg-neutral-100 text-neutral-700',
    
    // Outlined 变体
    'outlined-primary': 'border border-primary-600 text-primary-600 bg-white',
    'outlined-secondary': 'border border-neutral-600 text-neutral-600 bg-white',
    'outlined-success': 'border border-success-600 text-success-600 bg-white',
    'outlined-error': 'border border-error-600 text-error-600 bg-white',
    'outlined-warning': 'border border-warning-600 text-warning-600 bg-white',
    'outlined-info': 'border border-info-600 text-info-600 bg-white',
    'outlined-neutral': 'border border-neutral-600 text-neutral-600 bg-white'
  }
  
  const roundedClass = props.rounded ? 'rounded-full' : ''
  
  const variantColorKey = `${props.variant}-${props.color}` as keyof typeof variantColorClasses
  
  return cn(
    base,
    sizeClasses[props.size],
    variantColorClasses[variantColorKey],
    roundedClass
  )
})

const handleClick = (event: MouseEvent) => {
  emit('click', event)
}

const handleClose = (event: MouseEvent) => {
  emit('close', event)
}
</script>

<style scoped>
.gv-tag__icon {
  font-size: 0.875em;
}

.gv-tag__close {
  font-size: 0.875em;
  margin-left: 4px;
  cursor: pointer;
  opacity: 0.7;
  transition: opacity 0.2s;
}

.gv-tag__close:hover {
  opacity: 1;
}

.dark .gv-tag[class*="outlined"] {
  @apply bg-surface-dark;
}
</style>
