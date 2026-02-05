/**
 * EmptyState 组件类型定义
 * 用于显示空状态
 */

export interface EmptyStateProps {
  /**
   * 空状态图标（可以是 emoji 或图标组件）
   * @default '📭'
   */
  icon?: string

  /**
   * 空状态标题
   * @default '暂无数据'
   */
  title?: string

  /**
   * 空状态描述
   */
  description?: string

  /**
   * 操作按钮文本
   */
  actionText?: string

  /**
   * 是否显示操作按钮
   * @default false
   */
  showAction?: boolean
}

export interface EmptyStateEmits {
  /**
   * 点击操作按钮时触发
   */
  (e: 'action'): void
}
