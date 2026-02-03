/**
 * 间距系统 - 基于 8px 网格
 * Material Design 3 推荐的间距规范
 */

export const spacing = {
  0: '0',
  0.5: '0.125rem',  // 2px
  1: '0.25rem',     // 4px
  1.5: '0.375rem',  // 6px
  2: '0.5rem',      // 8px  ← 基准单位
  3: '0.75rem',     // 12px
  4: '1rem',        // 16px
  5: '1.25rem',     // 20px
  6: '1.5rem',      // 24px
  8: '2rem',        // 32px
  10: '2.5rem',     // 40px
  12: '3rem',       // 48px
  16: '4rem',       // 64px
  20: '5rem',       // 80px
  24: '6rem',       // 96px
  32: '8rem'        // 128px
} as const

/**
 * 组件内边距规范
 * 用于快速应用统一的内边距
 */
export const componentPadding = {
  xs: { x: spacing[2], y: spacing[1] },      // 8px × 4px
  sm: { x: spacing[3], y: spacing[1.5] },    // 12px × 6px
  md: { x: spacing[4], y: spacing[2] },      // 16px × 8px
  lg: { x: spacing[6], y: spacing[3] },      // 24px × 12px
  xl: { x: spacing[8], y: spacing[4] }       // 32px × 16px
} as const

/**
 * 容器最大宽度
 */
export const containerMaxWidth = {
  sm: '640px',
  md: '768px',
  lg: '1024px',
  xl: '1280px',
  '2xl': '1536px',
  full: '100%'
} as const

export type SpacingSize = keyof typeof spacing
export type ComponentPaddingSize = keyof typeof componentPadding
export type ContainerMaxWidth = keyof typeof containerMaxWidth
