/**
 * StatusBadge 业务组件类型定义
 * 统一的状态徽章组件
 */

/**
 * 状态类型
 */
export type StatusType = 
  | 'running'      // 运行中
  | 'stopped'      // 已停止
  | 'pending'      // 待处理
  | 'processing'   // 处理中
  | 'success'      // 成功
  | 'failed'       // 失败
  | 'error'        // 错误
  | 'warning'      // 警告
  | 'active'       // 激活
  | 'inactive'     // 未激活
  | 'online'       // 在线
  | 'offline'      // 离线
  | 'enabled'      // 启用
  | 'disabled'     // 禁用
  | 'neutral'      // 未知/中性

/**
 * StatusBadge Props
 */
export interface StatusBadgeProps {
  /**
   * 状态类型
   */
  status: StatusType
  
  /**
   * 自定义显示文本（不传则使用默认文本）
   */
  text?: string
  
  /**
   * 是否显示动画（运行中、处理中等状态）
   * @default true
   */
  animated?: boolean
  
  /**
   * 是否显示图标
   * @default true
   */
  showIcon?: boolean
}
