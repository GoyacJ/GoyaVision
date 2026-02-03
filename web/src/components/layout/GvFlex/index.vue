<template>
  <div :class="flexClasses">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'
import type { FlexProps } from './types'

const props = withDefaults(defineProps<FlexProps>(), {
  direction: 'row',
  wrap: 'nowrap',
  justify: 'start',
  align: 'stretch',
  gap: 'md',
  inline: false,
  vertical: false
})

// 间距映射
const gapMap = {
  none: 'gap-0',
  xs: 'gap-2',
  sm: 'gap-3',
  md: 'gap-4',
  lg: 'gap-6',
  xl: 'gap-8'
}

// 方向映射
const directionMap = {
  row: 'flex-row',
  'row-reverse': 'flex-row-reverse',
  col: 'flex-col',
  'col-reverse': 'flex-col-reverse'
}

// 换行映射
const wrapMap = {
  nowrap: 'flex-nowrap',
  wrap: 'flex-wrap',
  'wrap-reverse': 'flex-wrap-reverse'
}

// 主轴对齐映射
const justifyMap = {
  start: 'justify-start',
  center: 'justify-center',
  end: 'justify-end',
  between: 'justify-between',
  around: 'justify-around',
  evenly: 'justify-evenly'
}

// 交叉轴对齐映射
const alignMap = {
  start: 'items-start',
  center: 'items-center',
  end: 'items-end',
  stretch: 'items-stretch',
  baseline: 'items-baseline'
}

// Flex 类名
const flexClasses = computed(() => {
  const base = ['gv-flex', props.inline ? 'inline-flex' : 'flex']
  
  // 如果使用 vertical 快捷方式
  const actualDirection = props.vertical ? 'col' : props.direction
  
  return cn(
    base,
    directionMap[actualDirection],
    wrapMap[props.wrap],
    justifyMap[props.justify],
    alignMap[props.align],
    gapMap[props.gap]
  )
})
</script>
