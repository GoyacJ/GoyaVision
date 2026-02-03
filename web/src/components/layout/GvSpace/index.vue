<template>
  <div :class="spaceClasses" :style="spaceStyle">
    <template v-for="(child, index) in children" :key="index">
      <div
        :class="itemClasses"
        :style="itemStyle"
      >
        <component :is="child" />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, useSlots } from 'vue'
import { cn } from '@/utils/cn'
import type { SpaceProps, SpaceSize } from './types'

const props = withDefaults(defineProps<SpaceProps>(), {
  size: 'md',
  direction: 'horizontal',
  align: 'center',
  wrap: false,
  fill: false
})

const slots = useSlots()

// 获取子元素
const children = computed(() => {
  const defaultSlot = slots.default?.()
  return defaultSlot?.filter(child => child.type !== Comment) || []
})

// 间距大小映射
const sizeMap = {
  xs: '0.5rem',    // 8px
  sm: '0.75rem',   // 12px
  md: '1rem',      // 16px
  lg: '1.5rem',    // 24px
  xl: '2rem'       // 32px
}

// 获取间距值
const getGapValue = (size: SpaceSize): string => {
  if (typeof size === 'number') {
    return `${size}px`
  }
  return sizeMap[size] || sizeMap.md
}

// 计算间距
const gap = computed(() => {
  if (Array.isArray(props.size)) {
    const [horizontal, vertical] = props.size
    return {
      columnGap: getGapValue(horizontal),
      rowGap: getGapValue(vertical)
    }
  }
  const gapValue = getGapValue(props.size)
  return {
    columnGap: gapValue,
    rowGap: gapValue
  }
})

// Space 类名
const spaceClasses = computed(() => {
  const base = ['gv-space', 'inline-flex']
  
  // 方向
  const directionClass = props.direction === 'vertical' ? 'flex-col' : 'flex-row'
  
  // 对齐方式
  const alignMap = {
    start: 'items-start',
    center: 'items-center',
    end: 'items-end',
    baseline: 'items-baseline'
  }
  
  // 换行
  const wrapClass = props.wrap && props.direction === 'horizontal' ? 'flex-wrap' : ''
  
  // 填充
  const fillClass = props.fill ? 'w-full' : ''
  
  return cn(
    base,
    directionClass,
    alignMap[props.align],
    wrapClass,
    fillClass
  )
})

// Space 样式
const spaceStyle = computed(() => {
  return {
    ...gap.value
  }
})

// 子项类名
const itemClasses = computed(() => {
  return cn('gv-space__item', props.fill && 'flex-1')
})

// 子项样式
const itemStyle = computed(() => {
  if (props.fill && props.fillRatio) {
    return {
      flex: `${props.fillRatio} ${props.fillRatio} 0`
    }
  }
  return {}
})
</script>

<style scoped>
.gv-space__item {
  display: inline-flex;
}
</style>
