/**
 * GvModal 组件类型定义
 * Material Design 3 模态框规范
 */

/**
 * 模态框尺寸
 */
export type ModalSize = 'small' | 'medium' | 'large' | 'full'

/**
 * 模态框 Props
 */
export interface ModalProps {
  /**
   * 是否显示
   */
  modelValue: boolean
  
  /**
   * 标题
   */
  title?: string
  
  /**
   * 模态框尺寸
   * @default 'medium'
   */
  size?: ModalSize
  
  /**
   * 是否显示关闭按钮
   * @default true
   */
  showClose?: boolean
  
  /**
   * 是否可以通过点击遮罩关闭
   * @default true
   */
  closeOnClickModal?: boolean
  
  /**
   * 是否可以通过按下 ESC 关闭
   * @default true
   */
  closeOnPressEscape?: boolean
  
  /**
   * 是否在关闭时销毁子元素
   * @default false
   */
  destroyOnClose?: boolean
  
  /**
   * 是否显示底部操作区域
   * @default true
   */
  showFooter?: boolean
  
  /**
   * 确认按钮文字
   * @default '确定'
   */
  confirmText?: string
  
  /**
   * 取消按钮文字
   * @default '取消'
   */
  cancelText?: string
  
  /**
   * 是否显示确认按钮
   * @default true
   */
  showConfirm?: boolean
  
  /**
   * 是否显示取消按钮
   * @default true
   */
  showCancel?: boolean
  
  /**
   * 确认按钮加载状态
   * @default false
   */
  confirmLoading?: boolean
  
  /**
   * 是否居中显示
   * @default false
   */
  center?: boolean
  
  /**
   * 自定义类名
   */
  customClass?: string
  
  /**
   * z-index
   * @default 1000
   */
  zIndex?: number
}

/**
 * 模态框 Emits
 */
export interface ModalEmits {
  /**
   * 更新显示状态
   */
  (e: 'update:modelValue', value: boolean): void
  
  /**
   * 打开事件
   */
  (e: 'open'): void
  
  /**
   * 打开动画结束事件
   */
  (e: 'opened'): void
  
  /**
   * 关闭事件
   */
  (e: 'close'): void
  
  /**
   * 关闭动画结束事件
   */
  (e: 'closed'): void
  
  /**
   * 确认事件
   */
  (e: 'confirm'): void
  
  /**
   * 取消事件
   */
  (e: 'cancel'): void
}
