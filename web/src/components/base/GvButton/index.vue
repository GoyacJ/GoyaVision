<template>
  <component
    :is="tag"
    :class="buttonClasses"
    :disabled="disabled || loading"
    :type="href ? undefined : type"
    :href="href"
    :target="target"
    @click="handleClick"
  >
    <!-- 加载图标 -->
    <span v-if="loading" class="gv-button__icon gv-button__icon--loading">
      <el-icon class="is-loading">
        <Loading />
      </el-icon>
    </span>
    
    <!-- 左侧图标 -->
    <span v-else-if="icon && iconPosition === 'left'" class="gv-button__icon">
      <el-icon v-if="typeof icon === 'string'">
        <component :is="icon" />
      </el-icon>
      <component v-else :is="icon" />
    </span>
    
    <!-- 文字内容 -->
    <span v-if="$slots.default" class="gv-button__text">
      <slot />
    </span>
    
    <!-- 右侧图标 -->
    <span v-if="icon && iconPosition === 'right' && !loading" class="gv-button__icon">
      <el-icon v-if="typeof icon === 'string'">
        <component :is="icon" />
      </el-icon>
      <component v-else :is="icon" />
    </span>
  </component>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'
import type { ButtonProps, ButtonEmits } from './types'

const props = withDefaults(defineProps<ButtonProps>(), {
  variant: 'filled',
  color: 'primary',
  size: 'medium',
  iconPosition: 'left',
  type: 'button',
  disabled: false,
  loading: false,
  rounded: false,
  block: false
})

const emit = defineEmits<ButtonEmits>()

// 根据是否有 href 决定渲染为 button 还是 a 标签
const tag = computed(() => (props.href ? 'a' : 'button'))

// 按钮类名
const buttonClasses = computed(() => {
  const base = [
    'gv-button',
    'relative inline-flex items-center justify-center',
    'font-medium transition-all duration-fast',
    'focus:outline-none focus:ring-2 focus:ring-offset-2',
    'disabled:cursor-not-allowed disabled:opacity-50'
  ]

  // 尺寸 - 使用更克制的圆角
  const sizeClasses = {
    small: 'h-8 px-3 text-sm gap-1.5 rounded',      // 6px
    medium: 'h-10 px-4 text-base gap-2 rounded-md',  // 8px
    large: 'h-12 px-6 text-lg gap-2.5 rounded-lg'    // 12px
  }
  
  // 变体和颜色组合（移除彩色阴影，secondary 改为 neutral）
  const variantColorClasses = {
    // Filled 变体（填充按钮）- 使用极简阴影
    'filled-primary': 'bg-primary-600 text-white hover:bg-primary-700 active:bg-primary-800 focus:ring-primary-500 shadow-sm',
    'filled-secondary': 'bg-neutral-600 text-white hover:bg-neutral-700 active:bg-neutral-800 focus:ring-neutral-500 shadow-sm',
    'filled-success': 'bg-success-600 text-white hover:bg-success-700 active:bg-success-800 focus:ring-success-500 shadow-sm',
    'filled-error': 'bg-error-600 text-white hover:bg-error-700 active:bg-error-800 focus:ring-error-500 shadow-sm',
    'filled-warning': 'bg-warning-600 text-white hover:bg-warning-700 active:bg-warning-800 focus:ring-warning-500 shadow-sm',
    'filled-info': 'bg-info-600 text-white hover:bg-info-700 active:bg-info-800 focus:ring-info-500 shadow-sm',
    
    // Tonal 变体（色调按钮）- secondary 改为 neutral
    'tonal-primary': 'bg-primary-100 text-primary-700 hover:bg-primary-200 active:bg-primary-300 focus:ring-primary-500',
    'tonal-secondary': 'bg-neutral-100 text-neutral-700 hover:bg-neutral-200 active:bg-neutral-300 focus:ring-neutral-500',
    'tonal-success': 'bg-success-100 text-success-700 hover:bg-success-200 active:bg-success-300 focus:ring-success-500',
    'tonal-error': 'bg-error-100 text-error-700 hover:bg-error-200 active:bg-error-300 focus:ring-error-500',
    'tonal-warning': 'bg-warning-100 text-warning-700 hover:bg-warning-200 active:bg-warning-300 focus:ring-warning-500',
    'tonal-info': 'bg-info-100 text-info-700 hover:bg-info-200 active:bg-info-300 focus:ring-info-500',

    // Outlined 变体（边框按钮）- secondary 改为 neutral
    'outlined-primary': 'border-2 border-primary-600 text-primary-600 hover:bg-primary-50 active:bg-primary-100 focus:ring-primary-500',
    'outlined-secondary': 'border-2 border-neutral-600 text-neutral-600 hover:bg-neutral-50 active:bg-neutral-100 focus:ring-neutral-500',
    'outlined-success': 'border-2 border-success-600 text-success-600 hover:bg-success-50 active:bg-success-100 focus:ring-success-500',
    'outlined-error': 'border-2 border-error-600 text-error-600 hover:bg-error-50 active:bg-error-100 focus:ring-error-500',
    'outlined-warning': 'border-2 border-warning-600 text-warning-600 hover:bg-warning-50 active:bg-warning-100 focus:ring-warning-500',
    'outlined-info': 'border-2 border-info-600 text-info-600 hover:bg-info-50 active:bg-info-100 focus:ring-info-500',

    // Text 变体（文本按钮）- secondary 改为 neutral
    'text-primary': 'text-primary-600 hover:bg-primary-50 active:bg-primary-100 focus:ring-primary-500',
    'text-secondary': 'text-neutral-600 hover:bg-neutral-50 active:bg-neutral-100 focus:ring-neutral-500',
    'text-success': 'text-success-600 hover:bg-success-50 active:bg-success-100 focus:ring-success-500',
    'text-error': 'text-error-600 hover:bg-error-50 active:bg-error-100 focus:ring-error-500',
    'text-warning': 'text-warning-600 hover:bg-warning-50 active:bg-warning-100 focus:ring-warning-500',
    'text-info': 'text-info-600 hover:bg-info-50 active:bg-info-100 focus:ring-info-500'
  }
  
  // 圆形按钮
  const roundedClass = props.rounded ? 'rounded-full' : ''
  
  // 块级按钮
  const blockClass = props.block ? 'w-full' : ''
  
  const variantColorKey = `${props.variant}-${props.color}` as keyof typeof variantColorClasses
  
  return cn(
    base,
    sizeClasses[props.size],
    variantColorClasses[variantColorKey],
    roundedClass,
    blockClass
  )
})

// 点击事件处理
const handleClick = (event: MouseEvent) => {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
}
</script>

<style scoped>
.gv-button__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.gv-button__icon--loading .el-icon {
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

/* 移除过度的缩放动画，保持克制 */
</style>
