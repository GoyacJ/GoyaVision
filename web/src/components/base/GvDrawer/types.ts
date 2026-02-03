/**
 * GvDrawer 组件类型定义
 * Material Design 3 抽屉规范
 */

/**
 * 抽屉方向
 */
export type DrawerDirection = 'left' | 'right' | 'top' | 'bottom'

/**
 * 抽屉尺寸
 */
export type DrawerSize = 'small' | 'medium' | 'large' | 'full'

/**
 * 抽屉 Props
 */
export interface DrawerProps {
  /**
   * 是否显示
   */
  modelValue: boolean
  
  /**
   * 标题
   */
  title?: string
  
  /**
   * 抽屉方向
   * @default 'right'
   */
  direction?: DrawerDirection
  
  /**
   * 抽屉尺寸
   * @default 'medium'
   */
  size?: DrawerSize
  
  /**
   * 自定义宽度（direction 为 left/right 时）
   */
  width?: string
  
  /**
   * 自定义高度（direction 为 top/bottom 时）
   */
  height?: string
  
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
 * 抽屉 Emits
 */
export interface DrawerEmits {
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
