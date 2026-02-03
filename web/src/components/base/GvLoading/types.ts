/**
 * GvLoading 组件类型定义
 * Material Design 3 加载规范
 */

/**
 * 加载组件尺寸
 */
export type LoadingSize = 'small' | 'medium' | 'large'

/**
 * 加载组件类型
 */
export type LoadingType = 'spinner' | 'circular' | 'dots' | 'bars'

/**
 * 加载组件 Props
 */
export interface LoadingProps {
  /**
   * 是否显示加载
   * @default true
   */
  loading?: boolean
  
  /**
   * 加载组件类型
   * @default 'circular'
   */
  type?: LoadingType
  
  /**
   * 加载组件尺寸
   * @default 'medium'
   */
  size?: LoadingSize
  
  /**
   * 加载文本
   */
  text?: string
  
  /**
   * 是否全屏
   * @default false
   */
  fullscreen?: boolean
  
  /**
   * 是否锁定屏幕滚动（fullscreen 时有效）
   * @default true
   */
  lock?: boolean
  
  /**
   * 背景色
   */
  background?: string
  
  /**
   * 自定义类名
   */
  customClass?: string
  
  /**
   * 颜色
   * @default 'primary'
   */
  color?: 'primary' | 'secondary' | 'success' | 'error' | 'warning' | 'info' | 'white'
}
