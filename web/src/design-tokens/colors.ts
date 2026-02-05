/**
 * 颜色系统 - GoyaVision 克制设计系统
 * 设计原则：内容优先、极简、专业
 * 参考：Medium / Apple 官网风格
 */

export const colors = {
  // 主色：克制的蓝灰色（降低饱和度，更专业）
  primary: {
    DEFAULT: '#4F5B93',
    50: '#F5F6FA',
    100: '#EBEDF5',
    200: '#D4D8E8',
    300: '#B3BAD5',
    400: '#8A94B8',
    500: '#4F5B93',
    600: '#3E4A7A',
    700: '#2F3A61',
    800: '#232D4B',
    900: '#1A2238',
    950: '#131A2B'
  },
  
  // 功能色 - 成功
  success: {
    DEFAULT: '#10b981',
    50: '#f0fdf4',
    100: '#dcfce7',
    200: '#bbf7d0',
    300: '#86efac',
    400: '#4ade80',
    500: '#22c55e',
    600: '#10b981',
    700: '#059669',
    800: '#047857',
    900: '#065f46',
    950: '#064e3b'
  },
  
  // 功能色 - 错误
  error: {
    DEFAULT: '#ef4444',
    50: '#fef2f2',
    100: '#fee2e2',
    200: '#fecaca',
    300: '#fca5a5',
    400: '#f87171',
    500: '#ef4444',
    600: '#dc2626',
    700: '#b91c1c',
    800: '#991b1b',
    900: '#7f1d1d',
    950: '#450a0a'
  },
  
  // 功能色 - 警告
  warning: {
    DEFAULT: '#f59e0b',
    50: '#fffbeb',
    100: '#fef3c7',
    200: '#fde68a',
    300: '#fcd34d',
    400: '#fbbf24',
    500: '#f59e0b',
    600: '#d97706',
    700: '#b45309',
    800: '#92400e',
    900: '#78350f',
    950: '#451a03'
  },
  
  // 功能色 - 信息
  info: {
    DEFAULT: '#3b82f6',
    50: '#eff6ff',
    100: '#dbeafe',
    200: '#bfdbfe',
    300: '#93c5fd',
    400: '#60a5fa',
    500: '#3b82f6',
    600: '#2563eb',
    700: '#1d4ed8',
    800: '#1e40af',
    900: '#1e3a8a',
    950: '#172554'
  },
  
  // 中性色：极简灰度系统（9 级，适合内容优先设计）
  neutral: {
    DEFAULT: '#737373',
    50: '#FAFAFA',    // 页面背景
    100: '#F5F5F5',   // 容器背景
    200: '#E5E5E5',   // 边框
    300: '#D4D4D4',   // 禁用状态
    400: '#A3A3A3',   // 占位符
    500: '#737373',   // 次要文本
    600: '#525252',   // 主要文本
    700: '#404040',   // 标题
    800: '#262626',   // 深色文本
    900: '#171717',   // 强调文本
    950: '#0A0A0A'    // 最深（暗黑模式背景）
  },

  // 表面色（背景、卡片）- 简化版
  surface: {
    DEFAULT: '#FFFFFF',        // 白色卡片
    dim: '#FAFAFA',           // 浅灰背景
    container: '#F5F5F5',     // 容器背景
    containerHigh: '#E5E5E5', // 高层级容器
    dark: '#262626',          // 暗黑模式表面
    darkContainer: '#171717'  // 暗黑模式容器
  },

  // 文字色（语义化命名）
  text: {
    primary: '#262626',    // 主要文本（neutral.800）
    secondary: '#525252',  // 次要文本（neutral.600）
    tertiary: '#737373',   // 三级文本（neutral.500）
    placeholder: '#A3A3A3',// 占位符（neutral.400）
    disabled: '#D4D4D4',   // 禁用文本（neutral.300）
    inverse: '#FFFFFF'     // 反色文本（用于深色背景）
  }
} as const

/**
 * 设计说明：
 *
 * 1. 主色（Primary）：从鲜艳的 #667eea 降低饱和度至 #4F5B93
 *    - 更专业、更克制，适合长期内容消费
 *    - 减少视觉疲劳
 *
 * 2. 中性色（Neutral）：采用纯灰度（无色相偏向）
 *    - 不使用蓝灰色（#64748b），改用纯灰色（#737373）
 *    - 9 级灰度满足所有层级需求
 *
 * 3. 移除 Secondary 色系：
 *    - 不再使用渐变色（#667eea → #764ba2）
 *    - 避免过度装饰
 *
 * 4. 功能色保持不变：
 *    - Success: #10b981（绿色）
 *    - Warning: #f59e0b（橙色）
 *    - Error: #ef4444（红色）
 *    - Info: #3b82f6（蓝色）
 */

export type ColorPalette = typeof colors
export type ColorName = keyof typeof colors
export type ColorShade = keyof typeof colors.primary
