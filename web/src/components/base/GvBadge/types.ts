/**
 * GvBadge 组件类型定义
 * Material Design 3 徽章规范
 */

/**
 * 徽章颜色
 */
export type BadgeColor = 'primary' | 'secondary' | 'success' | 'error' | 'warning' | 'info' | 'neutral'

/**
 * 徽章尺寸
 */
export type BadgeSize = 'small' | 'medium' | 'large'

/**
 * 徽章变体
 */
export type BadgeVariant = 'filled' | 'tonal' | 'outlined'

/**
 * 徽章 Props
 */
export interface BadgeProps {
  /**
   * 徽章颜色
   * @default 'primary'
   */
  color?: BadgeColor
  
  /**
   * 徽章尺寸
   * @default 'medium'
   */
  size?: BadgeSize
  
  /**
   * 徽章变体
   * @default 'filled'
   */
  variant?: BadgeVariant
  
  /**
   * 显示的数值
   */
  value?: number | string
  
  /**
   * 最大值，超过显示为 max+
   * @default 99
   */
  max?: number
  
  /**
   * 是否显示为点状徽章
   * @default false
   */
  dot?: boolean
  
  /**
   * 是否隐藏徽章
   * @default false
   */
  hidden?: boolean
  
  /**
   * 自定义偏移量（用于角标模式）
   */
  offset?: [number, number]
}
