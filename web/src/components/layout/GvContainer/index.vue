<template>
  <div :class="containerClasses">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'
import type { ContainerProps } from './types'

const props = withDefaults(defineProps<ContainerProps>(), {
  maxWidth: 'xl',
  padding: true,
  centered: true
})

const containerClasses = computed(() => {
  const base = ['gv-container', 'w-full']
  
  // 最大宽度
  const maxWidthClasses = {
    sm: 'max-w-screen-sm',     // 640px
    md: 'max-w-screen-md',     // 768px
    lg: 'max-w-screen-lg',     // 1024px
    xl: 'max-w-screen-xl',     // 1280px
    '2xl': 'max-w-screen-2xl', // 1536px
    full: 'max-w-full'
  }
  
  // 水平内边距
  const paddingClass = props.padding ? 'px-4 md:px-6 lg:px-8' : ''
  
  // 居中对齐
  const centeredClass = props.centered ? 'mx-auto' : ''
  
  return cn(
    base,
    maxWidthClasses[props.maxWidth],
    paddingClass,
    centeredClass
  )
})
</script>
