/**
 * LoadingState 组件类型定义
 * 用于显示加载中状态
 */

export type LoadingSize = 'small' | 'medium' | 'large'

export interface LoadingStateProps {
  /**
   * 加载指示器大小
   * @default 'medium'
   */
  size?: LoadingSize

  /**
   * 加载提示文本
   */
  message?: string

  /**
   * 是否全屏显示
   * @default false
   */
  fullscreen?: boolean
}
