/**
 * 字体系统 - Material Design 3 字阶系统
 */

/**
 * 字体家族
 */
export const fontFamily = {
  sans: [
    'Inter',
    '-apple-system',
    'BlinkMacSystemFont',
    'Segoe UI',
    'Roboto',
    'Helvetica Neue',
    'Arial',
    'Noto Sans',
    'sans-serif',
    'Apple Color Emoji',
    'Segoe UI Emoji',
    'Segoe UI Symbol',
    'Noto Color Emoji'
  ].join(', '),
  
  mono: [
    'Fira Code',
    'Consolas',
    'Monaco',
    'Courier New',
    'monospace'
  ].join(', ')
} as const

/**
 * 字号和行高
 * Material Design 3 推荐的 9 档字阶
 */
export const fontSize = {
  xs: { size: '0.75rem', lineHeight: '1rem' },      // 12px / 16px
  sm: { size: '0.875rem', lineHeight: '1.25rem' },  // 14px / 20px
  base: { size: '1rem', lineHeight: '1.5rem' },     // 16px / 24px
  lg: { size: '1.125rem', lineHeight: '1.75rem' },  // 18px / 28px
  xl: { size: '1.25rem', lineHeight: '1.75rem' },   // 20px / 28px
  '2xl': { size: '1.5rem', lineHeight: '2rem' },    // 24px / 32px
  '3xl': { size: '1.875rem', lineHeight: '2.25rem' }, // 30px / 36px
  '4xl': { size: '2.25rem', lineHeight: '2.5rem' },   // 36px / 40px
  '5xl': { size: '3rem', lineHeight: '1' }            // 48px
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
 * Material Design 3 排版比例
 * 用于标题、正文、标签等文本样式
 */
export const typographyScale = {
  // Display - 超大标题
  displayLarge: {
    fontSize: fontSize['5xl'].size,
    lineHeight: fontSize['5xl'].lineHeight,
    fontWeight: fontWeight.bold
  },
  displayMedium: {
    fontSize: fontSize['4xl'].size,
    lineHeight: fontSize['4xl'].lineHeight,
    fontWeight: fontWeight.bold
  },
  displaySmall: {
    fontSize: fontSize['3xl'].size,
    lineHeight: fontSize['3xl'].lineHeight,
    fontWeight: fontWeight.bold
  },
  
  // Headline - 标题
  headlineLarge: {
    fontSize: fontSize['2xl'].size,
    lineHeight: fontSize['2xl'].lineHeight,
    fontWeight: fontWeight.semibold
  },
  headlineMedium: {
    fontSize: fontSize.xl.size,
    lineHeight: fontSize.xl.lineHeight,
    fontWeight: fontWeight.semibold
  },
  headlineSmall: {
    fontSize: fontSize.lg.size,
    lineHeight: fontSize.lg.lineHeight,
    fontWeight: fontWeight.semibold
  },
  
  // Title - 副标题
  titleLarge: {
    fontSize: fontSize.lg.size,
    lineHeight: fontSize.lg.lineHeight,
    fontWeight: fontWeight.medium
  },
  titleMedium: {
    fontSize: fontSize.base.size,
    lineHeight: fontSize.base.lineHeight,
    fontWeight: fontWeight.medium
  },
  titleSmall: {
    fontSize: fontSize.sm.size,
    lineHeight: fontSize.sm.lineHeight,
    fontWeight: fontWeight.medium
  },
  
  // Body - 正文
  bodyLarge: {
    fontSize: fontSize.base.size,
    lineHeight: fontSize.base.lineHeight,
    fontWeight: fontWeight.normal
  },
  bodyMedium: {
    fontSize: fontSize.sm.size,
    lineHeight: fontSize.sm.lineHeight,
    fontWeight: fontWeight.normal
  },
  bodySmall: {
    fontSize: fontSize.xs.size,
    lineHeight: fontSize.xs.lineHeight,
    fontWeight: fontWeight.normal
  },
  
  // Label - 标签
  labelLarge: {
    fontSize: fontSize.sm.size,
    lineHeight: fontSize.sm.lineHeight,
    fontWeight: fontWeight.medium
  },
  labelMedium: {
    fontSize: fontSize.xs.size,
    lineHeight: fontSize.xs.lineHeight,
    fontWeight: fontWeight.medium
  },
  labelSmall: {
    fontSize: fontSize.xs.size,
    lineHeight: fontSize.xs.lineHeight,
    fontWeight: fontWeight.normal
  }
} as const

export type FontSize = keyof typeof fontSize
export type FontWeight = keyof typeof fontWeight
export type TypographyScale = keyof typeof typographyScale
