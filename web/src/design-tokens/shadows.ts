/**
 * 阴影系统 - GoyaVision 克制设计系统
 * 设计原则：极简、功能性、不装饰
 *
 * 改动说明：
 * 1. 移除所有彩色阴影（primary, secondary, success, error）
 * 2. 降低阴影透明度（0.05 → 0.04, 0.1 → 0.06~0.08）
 * 3. 仅保留必要的层级区分
 */

export const shadows = {
  none: 'none',

  // 层级 1 - 轻微层级（卡片静态状态）
  sm: '0 1px 2px rgba(0, 0, 0, 0.04)',

  // 层级 2 - 标准卡片（默认）
  DEFAULT: '0 1px 3px rgba(0, 0, 0, 0.06)',

  // 层级 3 - 浮动元素（hover 卡片、下拉菜单）
  md: '0 4px 6px rgba(0, 0, 0, 0.07)',

  // 层级 4 - 模态框、抽屉
  lg: '0 10px 15px rgba(0, 0, 0, 0.08)',

  // 层级 5 - 最高层级（全屏遮罩、重要弹窗）
  xl: '0 20px 25px rgba(0, 0, 0, 0.10)',

  // 特殊阴影
  '2xl': '0 25px 50px rgba(0, 0, 0, 0.12)',
  inner: 'inset 0 2px 4px rgba(0, 0, 0, 0.04)'
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
