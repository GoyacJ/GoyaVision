/**
 * Design Tokens - GoyaVision 设计令牌系统
 * 基于 Material Design 3 规范
 */

export * from './colors'
export * from './spacing'
export * from './typography'
export * from './shadows'
export * from './radius'

/**
 * 动画曲线 - Material Design 3 运动规范
 */
export const easing = {
  // 标准缓动
  standard: 'cubic-bezier(0.2, 0, 0, 1)',
  
  // 强调缓动（适用于重要动作）
  emphasized: 'cubic-bezier(0.2, 0, 0, 1)',
  emphasizedDecelerate: 'cubic-bezier(0.05, 0.7, 0.1, 1)',
  emphasizedAccelerate: 'cubic-bezier(0.3, 0, 0.8, 0.15)',
  
  // 简单缓动
  linear: 'linear',
  easeIn: 'cubic-bezier(0.4, 0, 1, 1)',
  easeOut: 'cubic-bezier(0, 0, 0.2, 1)',
  easeInOut: 'cubic-bezier(0.4, 0, 0.2, 1)'
} as const

/**
 * 动画时长 - Material Design 3 推荐时长
 */
export const duration = {
  // 短时长（快速反馈）
  short1: 50,
  short2: 100,
  short3: 150,
  short4: 200,
  
  // 中等时长（默认动画）
  medium1: 250,
  medium2: 300,
  medium3: 350,
  medium4: 400,
  
  // 长时长（复杂动画）
  long1: 450,
  long2: 500,
  long3: 550,
  long4: 600,
  
  // 超长时长（页面过渡）
  extraLong1: 700,
  extraLong2: 800,
  extraLong3: 900,
  extraLong4: 1000
} as const

/**
 * 断点系统 - 响应式设计断点
 */
export const breakpoints = {
  xs: 0,
  sm: 640,
  md: 768,
  lg: 1024,
  xl: 1280,
  '2xl': 1536
} as const

/**
 * Z-Index 层级
 */
export const zIndex = {
  base: 0,
  dropdown: 1000,
  sticky: 1020,
  fixed: 1030,
  modalBackdrop: 1040,
  modal: 1050,
  popover: 1060,
  tooltip: 1070,
  notification: 1080
} as const

export type Easing = keyof typeof easing
export type Duration = keyof typeof duration
export type Breakpoint = keyof typeof breakpoints
export type ZIndex = keyof typeof zIndex
