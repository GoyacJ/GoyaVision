/**
 * 圆角系统 - GoyaVision 克制设计系统
 * 设计原则：适度圆角，避免过度圆润
 *
 * 改动说明：
 * - DEFAULT 从 8px 改为 6px（更克制）
 * - 整体圆角尺寸适度减小
 * - 移除 3xl（48px 过大）
 */

export const radius = {
  none: '0',
  sm: '0.25rem',      // 4px  - 小元素（标签、徽章）
  DEFAULT: '0.375rem',// 6px  - 基准圆角（按钮、输入框）
  md: '0.5rem',       // 8px  - 卡片
  lg: '0.75rem',      // 12px - 大卡片、容器
  xl: '1rem',         // 16px - 模态框、抽屉
  '2xl': '1.5rem',    // 24px - 大型容器
  full: '9999px'      // 完全圆形（头像）
} as const

/**
 * 按组件类型推荐的圆角大小
 */
export const radiusByComponent = {
  button: radius.DEFAULT,      // 6px - 更克制
  buttonSmall: radius.sm,      // 4px
  buttonLarge: radius.md,      // 8px
  card: radius.md,             // 8px - 从 12px 减小
  modal: radius.lg,            // 12px - 从 16px 减小
  drawer: radius.lg,           // 12px
  input: radius.DEFAULT,       // 6px - 更克制
  select: radius.DEFAULT,      // 6px
  badge: radius.sm,            // 4px
  tag: radius.sm,              // 4px
  avatar: radius.full,         // 圆形
  chip: radius.full,           // 圆形
  tooltip: radius.DEFAULT,     // 6px
  popover: radius.md,          // 8px
  dropdown: radius.md          // 8px
} as const

export type RadiusSize = keyof typeof radius
export type ComponentRadius = keyof typeof radiusByComponent
