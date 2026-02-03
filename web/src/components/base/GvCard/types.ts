/**
 * GvCard 组件类型定义
 * Material Design 3 卡片规范
 */

/**
 * 卡片阴影大小
 */
export type CardShadow = 'none' | 'sm' | 'md' | 'lg' | 'xl'

/**
 * 卡片内边距
 */
export type CardPadding = 'none' | 'sm' | 'md' | 'lg'

/**
 * 卡片 Props
 */
export interface CardProps {
  /**
   * 阴影大小
   * @default 'md'
   */
  shadow?: CardShadow
  
  /**
   * 内边距
   * @default 'md'
   */
  padding?: CardPadding
  
  /**
   * 是否支持悬停效果
   * @default false
   */
  hoverable?: boolean
  
  /**
   * 是否显示边框
   * @default false
   */
  bordered?: boolean
  
  /**
   * 背景色
   */
  background?: 'default' | 'container' | 'transparent'
}

/**
 * 卡片 Emits
 */
export interface CardEmits {
  /**
   * 点击事件（仅在 hoverable 为 true 时有效）
   */
  (e: 'click', event: MouseEvent): void
}
