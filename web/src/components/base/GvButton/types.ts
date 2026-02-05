/**
 * GvButton 组件类型定义
 * Material Design 3 按钮规范
 */

import type { Component } from 'vue'

/**
 * 按钮变体
 * - filled: 填充按钮（默认，最高强调）
 * - tonal: 色调按钮（中等强调）
 * - outlined: 边框按钮（中等强调）
 * - text: 文本按钮（最低强调）
 */
export type ButtonVariant = 'filled' | 'tonal' | 'outlined' | 'text'

/**
 * 按钮颜色
 * @deprecated secondary 已弃用，自动映射为 neutral
 */
export type ButtonColor = 'primary' | 'secondary' | 'success' | 'error' | 'warning' | 'info'

/**
 * 按钮尺寸
 */
export type ButtonSize = 'small' | 'medium' | 'large'

/**
 * 按钮 Props
 */
export interface ButtonProps {
  /**
   * Material Design 3 变体
   * @default 'filled'
   */
  variant?: ButtonVariant
  
  /**
   * 颜色主题
   * @default 'primary'
   */
  color?: ButtonColor
  
  /**
   * 按钮尺寸
   * @default 'medium'
   */
  size?: ButtonSize
  
  /**
   * 是否禁用
   * @default false
   */
  disabled?: boolean
  
  /**
   * 是否加载中
   * @default false
   */
  loading?: boolean
  
  /**
   * 图标（Element Plus 图标名或组件）
   */
  icon?: string | Component
  
  /**
   * 图标位置
   * @default 'left'
   */
  iconPosition?: 'left' | 'right'
  
  /**
   * 是否圆形按钮
   * @default false
   */
  rounded?: boolean
  
  /**
   * 是否块级按钮（宽度100%）
   * @default false
   */
  block?: boolean
  
  /**
   * HTML 按钮类型
   * @default 'button'
   */
  type?: 'button' | 'submit' | 'reset'
  
  /**
   * 链接地址（作为 a 标签）
   */
  href?: string
  
  /**
   * 链接打开方式
   */
  target?: '_blank' | '_self' | '_parent' | '_top'
}

/**
 * 按钮 Emits
 */
export interface ButtonEmits {
  /**
   * 点击事件
   */
  (e: 'click', event: MouseEvent): void
}
