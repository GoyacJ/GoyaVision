/**
 * 阴影系统 - Material Design 3 阴影层级
 */

export const shadows = {
  // 层级 1 - 微小提升（卡片）
  sm: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
  
  // 层级 2 - 小提升（悬停卡片）
  DEFAULT: '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1)',
  
  // 层级 3 - 中等提升（下拉菜单）
  md: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1)',
  
  // 层级 4 - 较大提升（模态框）
  lg: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1)',
  
  // 层级 5 - 最大提升（抽屉、浮动按钮）
  xl: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1)',
  
  // 特殊阴影
  '2xl': '0 25px 50px -12px rgba(0, 0, 0, 0.25)',
  inner: 'inset 0 2px 4px 0 rgba(0, 0, 0, 0.05)',
  none: 'none',
  
  // 彩色阴影（品牌色）
  primary: '0 8px 16px -4px rgba(102, 126, 234, 0.3)',
  secondary: '0 8px 16px -4px rgba(118, 75, 162, 0.3)',
  success: '0 8px 16px -4px rgba(16, 185, 129, 0.3)',
  error: '0 8px 16px -4px rgba(239, 68, 68, 0.3)',
  warning: '0 8px 16px -4px rgba(245, 158, 11, 0.3)',
  info: '0 8px 16px -4px rgba(59, 130, 246, 0.3)'
} as const

/**
 * Material Design 3 阴影层级语义化命名
 * 根据组件类型推荐使用的阴影
 */
export const shadowsByComponent = {
  card: shadows.sm,
  cardHover: shadows.DEFAULT,
  dropdown: shadows.md,
  modal: shadows.lg,
  drawer: shadows.xl,
  tooltip: shadows.md,
  popover: shadows.md,
  button: shadows.sm,
  buttonHover: shadows.md,
  fab: shadows.lg,          // Floating Action Button
  fabHover: shadows.xl
} as const

export type ShadowSize = keyof typeof shadows
export type ComponentShadow = keyof typeof shadowsByComponent
