/**
 * GvTag 组件类型定义
 * Material Design 3 标签规范
 */

/**
 * 标签颜色
 */
export type TagColor = 'primary' | 'secondary' | 'success' | 'error' | 'warning' | 'info' | 'neutral'

/**
 * 标签尺寸
 */
export type TagSize = 'small' | 'medium' | 'large'

/**
 * 标签变体
 */
export type TagVariant = 'filled' | 'tonal' | 'outlined'

/**
 * 标签 Props
 */
export interface TagProps {
  /**
   * 标签颜色
   * @default 'neutral'
   */
  color?: TagColor
  
  /**
   * 标签尺寸
   * @default 'medium'
   */
  size?: TagSize
  
  /**
   * 标签变体
   * @default 'tonal'
   */
  variant?: TagVariant
  
  /**
   * 是否可关闭
   * @default false
   */
  closable?: boolean
  
  /**
   * 是否圆形标签
   * @default false
   */
  rounded?: boolean
  
  /**
   * 前置图标
   */
  icon?: string
}

/**
 * 标签 Emits
 */
export interface TagEmits {
  /**
   * 点击事件
   */
  (e: 'click', event: MouseEvent): void
  
  /**
   * 关闭事件
   */
  (e: 'close', event: MouseEvent): void
}
