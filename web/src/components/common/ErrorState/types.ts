/**
 * ErrorState 组件类型定义
 * 用于显示错误状态
 */

export interface ErrorStateProps {
  /**
   * 错误对象
   */
  error?: Error | null

  /**
   * 错误标题
   * @default '加载失败'
   */
  title?: string

  /**
   * 错误描述
   * 如果未提供且 error 存在，则使用 error.message
   */
  message?: string

  /**
   * 重试按钮文本
   * @default '重试'
   */
  retryText?: string

  /**
   * 是否显示重试按钮
   * @default true
   */
  showRetry?: boolean
}

export interface ErrorStateEmits {
  /**
   * 点击重试按钮时触发
   */
  (e: 'retry'): void
}
