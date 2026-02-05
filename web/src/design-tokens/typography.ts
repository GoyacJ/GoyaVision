/**
 * 字体系统 - GoyaVision 克制设计系统
 * 设计原则：增强可读性、紧凑字距、内容优先
 *
 * 改动说明：
 * 1. 正文字号从 16px 改为 15px（参考 Medium）
 * 2. 增加字距（letterSpacing）
 * 3. 标题使用紧凑字距（负值）
 */

/**
 * 字体家族
 * 优先使用系统字体，确保跨平台一致性
 */
export const fontFamily = {
  sans: [
    '-apple-system',
    'BlinkMacSystemFont',
    'Segoe UI',
    'Roboto',
    'Helvetica Neue',
    'Arial',
    'Noto Sans SC',           // 中文优先
    'sans-serif',
    'Apple Color Emoji',
    'Segoe UI Emoji',
    'Segoe UI Symbol'
  ].join(', '),

  mono: [
    'SF Mono',
    'Menlo',
    'Monaco',
    'Consolas',
    'Liberation Mono',
    'Courier New',
    'monospace'
  ].join(', ')
} as const

/**
 * 字号和行高（参考 Apple / Medium 设计规范）
 * 正文字号从 16px 改为 15px，增强可读性
 */
export const fontSize = {
  xs: { size: '0.75rem', lineHeight: '1rem' },       // 12px / 16px
  sm: { size: '0.8125rem', lineHeight: '1.25rem' },  // 13px / 20px
  base: { size: '0.9375rem', lineHeight: '1.5rem' }, // 15px / 24px - 正文字号
  lg: { size: '1.125rem', lineHeight: '1.75rem' },   // 18px / 28px
  xl: { size: '1.25rem', lineHeight: '1.75rem' },    // 20px / 28px
  '2xl': { size: '1.5rem', lineHeight: '2rem' },     // 24px / 32px
  '3xl': { size: '1.875rem', lineHeight: '2.25rem' },// 30px / 36px
  '4xl': { size: '2rem', lineHeight: '2.5rem' },     // 32px / 40px
  '5xl': { size: '2.5rem', lineHeight: '1.2' }       // 40px
} as const

/**
 * 字重
 */
export const fontWeight = {
  light: 300,
  normal: 400,
  medium: 500,
  semibold: 600,
  bold: 700,
  extrabold: 800
} as const

/**
 * 字距（Letter Spacing）
 * 标题使用紧凑字距，正文使用标准字距
 */
export const letterSpacing = {
  tighter: '-0.02em',  // 超紧凑（大标题）
  tight: '-0.01em',    // 紧凑（标题）
  normal: '0',         // 标准（正文）
  wide: '0.025em',     // 宽松（标签）
  wider: '0.05em'      // 超宽松（全大写标签）
} as const

/**
 * 排版比例（Typography Scale）
 * 用于标题、正文、标签等文本样式
 * 增加字距配置，提升可读性
 */
export const typographyScale = {
  // Page Title - 页面标题
  h1: {
    fontSize: fontSize['4xl'].size,
    lineHeight: '1.25',
    fontWeight: fontWeight.bold,
    letterSpacing: letterSpacing.tighter
  },

  // Section Title - 区块标题
  h2: {
    fontSize: fontSize['2xl'].size,
    lineHeight: '1.3',
    fontWeight: fontWeight.semibold,
    letterSpacing: letterSpacing.tight
  },

  // Card Title - 卡片标题
  h3: {
    fontSize: fontSize.lg.size,
    lineHeight: '1.4',
    fontWeight: fontWeight.semibold,
    letterSpacing: letterSpacing.normal
  },

  // Small Heading - 小标题
  h4: {
    fontSize: fontSize.base.size,
    lineHeight: '1.5',
    fontWeight: fontWeight.semibold,
    letterSpacing: letterSpacing.normal
  },

  // Body - 正文（15px，增强可读性）
  body: {
    fontSize: fontSize.base.size,
    lineHeight: '1.6',
    fontWeight: fontWeight.normal,
    letterSpacing: letterSpacing.normal
  },

  // Body Small - 小正文
  bodySmall: {
    fontSize: fontSize.sm.size,
    lineHeight: '1.5',
    fontWeight: fontWeight.normal,
    letterSpacing: letterSpacing.normal
  },

  // Caption - 次要文本
  caption: {
    fontSize: fontSize.sm.size,
    lineHeight: '1.5',
    fontWeight: fontWeight.normal,
    letterSpacing: letterSpacing.normal,
    color: '#737373' // neutral.500
  },

  // Label - 标签/按钮文本
  label: {
    fontSize: fontSize.sm.size,
    lineHeight: '1.4',
    fontWeight: fontWeight.medium,
    letterSpacing: letterSpacing.normal
  },

  // Label Small - 小标签
  labelSmall: {
    fontSize: fontSize.xs.size,
    lineHeight: '1.4',
    fontWeight: fontWeight.medium,
    letterSpacing: letterSpacing.wide
  },

  // Overline - 全大写标签（如状态标签）
  overline: {
    fontSize: fontSize.xs.size,
    lineHeight: '1.4',
    fontWeight: fontWeight.medium,
    letterSpacing: letterSpacing.wider,
    textTransform: 'uppercase' as const
  }
} as const

export type FontSize = keyof typeof fontSize
export type FontWeight = keyof typeof fontWeight
export type LetterSpacing = keyof typeof letterSpacing
export type TypographyScale = keyof typeof typographyScale
