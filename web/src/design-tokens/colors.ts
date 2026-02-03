/**
 * 颜色系统 - Material Design 3 动态配色方案
 * 基于 GoyaVision 品牌色生成完整色阶
 */

export const colors = {
  // 主色（品牌色 #667eea）
  primary: {
    DEFAULT: '#667eea',
    50: '#f5f7ff',
    100: '#ebedff',
    200: '#d6dcff',
    300: '#b3bdff',
    400: '#8d9eff',
    500: '#667eea',
    600: '#5568d3',
    700: '#4553bd',
    800: '#3640a6',
    900: '#2a3290',
    950: '#1f2673'
  },
  
  // 辅助色
  secondary: {
    DEFAULT: '#764ba2',
    50: '#f9f5fc',
    100: '#f3ebf9',
    200: '#e7d7f3',
    300: '#d4b8e9',
    400: '#b88ed9',
    500: '#9d64c9',
    600: '#764ba2',
    700: '#65408b',
    800: '#543574',
    900: '#432a5d',
    950: '#321f46'
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
  
  // 中性色（灰阶）
  neutral: {
    DEFAULT: '#64748b',
    50: '#f8fafc',
    100: '#f1f5f9',
    200: '#e2e8f0',
    300: '#cbd5e1',
    400: '#94a3b8',
    500: '#64748b',
    600: '#475569',
    700: '#334155',
    800: '#1e293b',
    900: '#0f172a',
    950: '#020617'
  },
  
  // 表面色（背景、卡片）
  surface: {
    DEFAULT: '#ffffff',
    dark: '#1e293b',
    dim: '#f8fafc',
    bright: '#ffffff',
    container: '#f1f5f9',
    containerHigh: '#e2e8f0',
    containerHighest: '#cbd5e1'
  },
  
  // 文字色
  text: {
    primary: '#0f172a',
    secondary: '#475569',
    tertiary: '#64748b',
    disabled: '#cbd5e1',
    inverse: '#ffffff'
  }
} as const

export type ColorPalette = typeof colors
export type ColorName = keyof typeof colors
export type ColorShade = keyof typeof colors.primary
