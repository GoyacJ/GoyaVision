<template>
  <div :class="dividerClasses" role="separator">
    <!-- 水平分割线 -->
    <template v-if="direction === 'horizontal'">
      <span v-if="$slots.default" :class="contentClasses">
        <slot />
      </span>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'
import type { DividerProps } from './types'

const props = withDefaults(defineProps<DividerProps>(), {
  direction: 'horizontal',
  contentPosition: 'center',
  dashed: false
})

// 分割线类名
const dividerClasses = computed(() => {
  const base = ['gv-divider', 'relative']
  
  // 方向
  if (props.direction === 'horizontal') {
    const lineStyle = props.dashed
      ? 'border-t border-dashed border-neutral-200'
      : 'border-t border-neutral-200'
    
    return cn(
      base,
      'w-full my-4',
      lineStyle,
      props.$slots.default && 'flex items-center',
      props.customClass
    )
  } else {
    const lineStyle = props.dashed
      ? 'border-l border-dashed border-neutral-200'
      : 'border-l border-neutral-200'
    
    return cn(
      base,
      'h-full mx-4 inline-block',
      lineStyle,
      props.customClass
    )
  }
})

// 内容类名
const contentClasses = computed(() => {
  const base = [
    'gv-divider__content',
    'px-4 bg-surface text-text-secondary text-sm'
  ]
  
  const positionClasses = {
    left: 'mr-auto',
    center: 'mx-auto',
    right: 'ml-auto'
  }
  
  return cn(base, positionClasses[props.contentPosition])
})
</script>

<style scoped>
.dark .gv-divider {
  @apply border-neutral-700;
}

.dark .gv-divider__content {
  @apply bg-surface-dark text-neutral-400;
}
</style>
