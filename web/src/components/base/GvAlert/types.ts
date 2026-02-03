/**
 * GvAlert 组件类型定义
 * Material Design 3 警告框规范
 */

/**
 * 警告框类型
 */
export type AlertType = 'success' | 'info' | 'warning' | 'error'

/**
 * 警告框 Props
 */
export interface AlertProps {
  /**
   * 警告框类型
   * @default 'info'
   */
  type?: AlertType
  
  /**
   * 标题
   */
  title?: string
  
  /**
   * 描述文本
   */
  description?: string
  
  /**
   * 是否可关闭
   * @default false
   */
  closable?: boolean
  
  /**
   * 是否显示图标
   * @default true
   */
  showIcon?: boolean
  
  /**
   * 是否居中显示
   * @default false
   */
  center?: boolean
  
  /**
   * 关闭按钮文本
   */
  closeText?: string
}

/**
 * 警告框 Emits
 */
export interface AlertEmits {
  /**
   * 关闭事件
   */
  (e: 'close'): void
}
