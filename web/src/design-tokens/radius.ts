/**
 * 圆角系统 - Material Design 3 圆角规范
 */

export const radius = {
  none: '0',
  sm: '0.25rem',     // 4px  - 小元素（标签、徽章）
  DEFAULT: '0.5rem', // 8px  - 基准圆角
  md: '0.75rem',     // 12px - 卡片、按钮
  lg: '1rem',        // 16px - 大卡片
  xl: '1.5rem',      // 24px - 模态框
  '2xl': '2rem',     // 32px - 大型容器
  '3xl': '3rem',     // 48px - 特大容器
  full: '9999px'     // 完全圆形
} as const

/**
 * 按组件类型推荐的圆角大小
 */
export const radiusByComponent = {
  button: radius.md,
  buttonSmall: radius.sm,
  buttonLarge: radius.lg,
  card: radius.lg,
  modal: radius.xl,
  drawer: radius.xl,
  input: radius.md,
  select: radius.md,
  badge: radius.sm,
  tag: radius.sm,
  avatar: radius.full,
  chip: radius.full,
  tooltip: radius.md,
  popover: radius.md,
  dropdown: radius.md
} as const

export type RadiusSize = keyof typeof radius
export type ComponentRadius = keyof typeof radiusByComponent
