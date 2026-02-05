<template>
  <span v-if="!$slots.default" :class="standaloneBadgeClasses">
    <span v-if="!dot" class="gv-badge__content">
      {{ displayValue }}
    </span>
  </span>
  
  <span v-else class="gv-badge gv-badge--wrapper relative inline-flex">
    <slot />
    <span
      v-if="!hidden && (dot || value !== undefined)"
      :class="badgeClasses"
      :style="badgeStyle"
    >
      <span v-if="!dot" class="gv-badge__content">
        {{ displayValue }}
      </span>
    </span>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'
import type { BadgeProps } from './types'

const props = withDefaults(defineProps<BadgeProps>(), {
  color: 'primary',
  size: 'medium',
  variant: 'filled',
  max: 99,
  dot: false,
  hidden: false
})

// 显示的值
const displayValue = computed(() => {
  if (typeof props.value === 'number' && props.value > props.max) {
    return `${props.max}+`
  }
  return props.value
})

// 独立徽章类名（无子元素时）
const standaloneBadgeClasses = computed(() => {
  const base = [
    'gv-badge',
    'gv-badge--standalone',
    'inline-flex items-center justify-center',
    'font-medium rounded-full',
    'transition-all duration-200'
  ]
  
  // 尺寸
  const sizeClasses = {
    small: 'h-5 px-2 text-xs min-w-[20px]',
    medium: 'h-6 px-2.5 text-sm min-w-[24px]',
    large: 'h-7 px-3 text-base min-w-[28px]'
  }
  
  // 点状徽章尺寸
  const dotSizeClasses = {
    small: 'w-2 h-2',
    medium: 'w-2.5 h-2.5',
    large: 'w-3 h-3'
  }
  
  return cn(
    base,
    props.dot ? dotSizeClasses[props.size] : sizeClasses[props.size],
    variantColorClasses.value
  )
})

// 角标徽章类名
const badgeClasses = computed(() => {
  const base = [
    'gv-badge__indicator',
    'absolute inline-flex items-center justify-center',
    'font-medium rounded-full',
    'border-2 border-white',
    'transition-all duration-200'
  ]
  
  // 尺寸
  const sizeClasses = {
    small: props.dot ? 'w-2 h-2' : 'h-4 px-1.5 text-xs min-w-[16px]',
    medium: props.dot ? 'w-2.5 h-2.5' : 'h-5 px-2 text-xs min-w-[20px]',
    large: props.dot ? 'w-3 h-3' : 'h-6 px-2.5 text-sm min-w-[24px]'
  }
  
  // 默认位置（右上角）
  const positionClass = 'top-0 right-0 transform translate-x-1/2 -translate-y-1/2'
  
  return cn(
    base,
    sizeClasses[props.size],
    positionClass,
    variantColorClasses.value
  )
})

// 变体和颜色组合类名
const variantColorClasses = computed(() => {
  const classes = {
    // Filled 变体
    'filled-primary': 'bg-primary-600 text-white',
    'filled-secondary': 'bg-neutral-600 text-white',
    'filled-success': 'bg-success-600 text-white',
    'filled-error': 'bg-error-600 text-white',
    'filled-warning': 'bg-warning-600 text-white',
    'filled-info': 'bg-info-600 text-white',
    'filled-neutral': 'bg-neutral-600 text-white',
    
    // Tonal 变体
    'tonal-primary': 'bg-primary-100 text-primary-700',
    'tonal-secondary': 'bg-neutral-100 text-neutral-700',
    'tonal-success': 'bg-success-100 text-success-700',
    'tonal-error': 'bg-error-100 text-error-700',
    'tonal-warning': 'bg-warning-100 text-warning-700',
    'tonal-info': 'bg-info-100 text-info-700',
    'tonal-neutral': 'bg-neutral-100 text-neutral-700',
    
    // Outlined 变体
    'outlined-primary': 'border-2 border-primary-600 text-primary-600 bg-white',
    'outlined-secondary': 'border-2 border-neutral-600 text-neutral-600 bg-white',
    'outlined-success': 'border-2 border-success-600 text-success-600 bg-white',
    'outlined-error': 'border-2 border-error-600 text-error-600 bg-white',
    'outlined-warning': 'border-2 border-warning-600 text-warning-600 bg-white',
    'outlined-info': 'border-2 border-info-600 text-info-600 bg-white',
    'outlined-neutral': 'border-2 border-neutral-600 text-neutral-600 bg-white'
  }
  
  const key = `${props.variant}-${props.color}` as keyof typeof classes
  return classes[key]
})

// 自定义偏移样式
const badgeStyle = computed(() => {
  if (!props.offset) return {}
  
  const [x, y] = props.offset
  return {
    transform: `translate(${x}px, ${y}px)`
  }
})
</script>

<style scoped>
.gv-badge__content {
  line-height: 1;
}

.dark .gv-badge__indicator {
  @apply border-neutral-800;
}
</style>
