/**
 * 响应式断点 Composable
 * 用于在 JavaScript 中判断当前屏幕尺寸
 */

import { computed } from 'vue'
import { useWindowSize } from '@vueuse/core'
import { breakpoints } from '@/design-tokens'

export function useBreakpoint() {
  const { width } = useWindowSize()
  
  /**
   * 当前是否为移动端
   */
  const isMobile = computed(() => width.value < breakpoints.md)
  
  /**
   * 当前是否为平板
   */
  const isTablet = computed(() => width.value >= breakpoints.md && width.value < breakpoints.lg)
  
  /**
   * 当前是否为桌面端
   */
  const isDesktop = computed(() => width.value >= breakpoints.lg)
  
  /**
   * 当前是否为小屏幕
   */
  const isSmall = computed(() => width.value < breakpoints.sm)
  
  /**
   * 当前是否为大屏幕
   */
  const isLarge = computed(() => width.value >= breakpoints.xl)
  
  /**
   * 当前断点名称
   */
  const current = computed(() => {
    if (width.value < breakpoints.sm) return 'xs'
    if (width.value < breakpoints.md) return 'sm'
    if (width.value < breakpoints.lg) return 'md'
    if (width.value < breakpoints.xl) return 'lg'
    if (width.value < breakpoints['2xl']) return 'xl'
    return '2xl'
  })
  
  /**
   * 是否大于等于指定断点
   */
  const isGreaterOrEqual = (breakpoint: keyof typeof breakpoints) => {
    return computed(() => width.value >= breakpoints[breakpoint])
  }
  
  return {
    width,
    isMobile,
    isTablet,
    isDesktop,
    isSmall,
    isLarge,
    current,
    isGreaterOrEqual
  }
}
