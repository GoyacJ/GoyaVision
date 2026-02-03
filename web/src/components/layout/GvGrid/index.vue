<template>
  <div :class="gridClasses" :style="gridStyle">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'
import type { GridProps } from './types'

const props = withDefaults(defineProps<GridProps>(), {
  cols: 3,
  gap: 'md',
  autoFill: false,
  minColWidth: '200px'
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

const gapXMap = {
  none: 'gap-x-0',
  xs: 'gap-x-2',
  sm: 'gap-x-3',
  md: 'gap-x-4',
  lg: 'gap-x-6',
  xl: 'gap-x-8'
}

const gapYMap = {
  none: 'gap-y-0',
  xs: 'gap-y-2',
  sm: 'gap-y-3',
  md: 'gap-y-4',
  lg: 'gap-y-6',
  xl: 'gap-y-8'
}

// 对齐方式映射
const alignMap = {
  start: 'items-start',
  center: 'items-center',
  end: 'items-end',
  stretch: 'items-stretch'
}

const justifyMap = {
  start: 'justify-items-start',
  center: 'justify-items-center',
  end: 'justify-items-end',
  between: 'justify-items-between',
  around: 'justify-items-around',
  evenly: 'justify-items-evenly'
}

// 网格类名
const gridClasses = computed(() => {
  const base = ['gv-grid', 'grid']
  
  // 自动填充模式
  if (props.autoFill) {
    return cn(
      base,
      props.gapX ? gapXMap[props.gapX] : props.gapY ? gapYMap[props.gapY] : gapMap[props.gap],
      props.align && alignMap[props.align],
      props.justify && justifyMap[props.justify]
    )
  }
  
  // 固定列数模式
  const colsClasses = []
  
  if (props.responsive) {
    // 响应式列数
    if (props.responsive.xs) colsClasses.push(`grid-cols-${props.responsive.xs}`)
    if (props.responsive.sm) colsClasses.push(`sm:grid-cols-${props.responsive.sm}`)
    if (props.responsive.md) colsClasses.push(`md:grid-cols-${props.responsive.md}`)
    if (props.responsive.lg) colsClasses.push(`lg:grid-cols-${props.responsive.lg}`)
    if (props.responsive.xl) colsClasses.push(`xl:grid-cols-${props.responsive.xl}`)
    if (props.responsive['2xl']) colsClasses.push(`2xl:grid-cols-${props.responsive['2xl']}`)
  } else {
    // 固定列数
    colsClasses.push(`grid-cols-${props.cols}`)
  }
  
  // 间距
  const gapClasses = []
  if (props.gapX) {
    gapClasses.push(gapXMap[props.gapX])
  }
  if (props.gapY) {
    gapClasses.push(gapYMap[props.gapY])
  }
  if (!props.gapX && !props.gapY) {
    gapClasses.push(gapMap[props.gap])
  }
  
  return cn(
    base,
    colsClasses,
    gapClasses,
    props.align && alignMap[props.align],
    props.justify && justifyMap[props.justify]
  )
})

// 网格样式（自动填充模式）
const gridStyle = computed(() => {
  if (props.autoFill) {
    return {
      gridTemplateColumns: `repeat(auto-fill, minmax(${props.minColWidth}, 1fr))`
    }
  }
  return {}
})
</script>
