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
    'font-medium transition-all duration-300',
    'focus:outline-none focus:ring-2 focus:ring-offset-2',
    'disabled:cursor-not-allowed disabled:opacity-50'
  ]
  
  // 尺寸
  const sizeClasses = {
    small: 'h-8 px-3 text-sm gap-1.5 rounded-lg',
    medium: 'h-10 px-4 text-base gap-2 rounded-xl',
    large: 'h-12 px-6 text-lg gap-2.5 rounded-2xl'
  }
  
  // 变体和颜色组合
  const variantColorClasses = {
    // Filled 变体（填充按钮）
    'filled-primary': 'bg-primary-600 text-white hover:bg-primary-700 active:bg-primary-800 focus:ring-primary-500 shadow-primary',
    'filled-secondary': 'bg-secondary-600 text-white hover:bg-secondary-700 active:bg-secondary-800 focus:ring-secondary-500 shadow-secondary',
    'filled-success': 'bg-success-600 text-white hover:bg-success-700 active:bg-success-800 focus:ring-success-500 shadow-success',
    'filled-error': 'bg-error-600 text-white hover:bg-error-700 active:bg-error-800 focus:ring-error-500 shadow-error',
    'filled-warning': 'bg-warning-600 text-white hover:bg-warning-700 active:bg-warning-800 focus:ring-warning-500 shadow-warning',
    'filled-info': 'bg-info-600 text-white hover:bg-info-700 active:bg-info-800 focus:ring-info-500 shadow-info',
    
    // Tonal 变体（色调按钮）
    'tonal-primary': 'bg-primary-100 text-primary-700 hover:bg-primary-200 active:bg-primary-300 focus:ring-primary-500',
    'tonal-secondary': 'bg-secondary-100 text-secondary-700 hover:bg-secondary-200 active:bg-secondary-300 focus:ring-secondary-500',
    'tonal-success': 'bg-success-100 text-success-700 hover:bg-success-200 active:bg-success-300 focus:ring-success-500',
    'tonal-error': 'bg-error-100 text-error-700 hover:bg-error-200 active:bg-error-300 focus:ring-error-500',
    'tonal-warning': 'bg-warning-100 text-warning-700 hover:bg-warning-200 active:bg-warning-300 focus:ring-warning-500',
    'tonal-info': 'bg-info-100 text-info-700 hover:bg-info-200 active:bg-info-300 focus:ring-info-500',
    
    // Outlined 变体（边框按钮）
    'outlined-primary': 'border-2 border-primary-600 text-primary-600 hover:bg-primary-50 active:bg-primary-100 focus:ring-primary-500',
    'outlined-secondary': 'border-2 border-secondary-600 text-secondary-600 hover:bg-secondary-50 active:bg-secondary-100 focus:ring-secondary-500',
    'outlined-success': 'border-2 border-success-600 text-success-600 hover:bg-success-50 active:bg-success-100 focus:ring-success-500',
    'outlined-error': 'border-2 border-error-600 text-error-600 hover:bg-error-50 active:bg-error-100 focus:ring-error-500',
    'outlined-warning': 'border-2 border-warning-600 text-warning-600 hover:bg-warning-50 active:bg-warning-100 focus:ring-warning-500',
    'outlined-info': 'border-2 border-info-600 text-info-600 hover:bg-info-50 active:bg-info-100 focus:ring-info-500',
    
    // Text 变体（文本按钮）
    'text-primary': 'text-primary-600 hover:bg-primary-50 active:bg-primary-100 focus:ring-primary-500',
    'text-secondary': 'text-secondary-600 hover:bg-secondary-50 active:bg-secondary-100 focus:ring-secondary-500',
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

.gv-button:active {
  transform: scale(0.98);
}

.gv-button:disabled {
  transform: none;
}
</style>
