/**
 * GvDivider 组件类型定义
 */

/**
 * 分割线方向
 */
export type DividerDirection = 'horizontal' | 'vertical'

/**
 * 文本位置
 */
export type ContentPosition = 'left' | 'center' | 'right'

/**
 * GvDivider Props
 */
export interface DividerProps {
  /**
   * 方向
   * @default 'horizontal'
   */
  direction?: DividerDirection
  
  /**
   * 文本位置
   * @default 'center'
   */
  contentPosition?: ContentPosition
  
  /**
   * 是否为虚线
   * @default false
   */
  dashed?: boolean
  
  /**
   * 自定义类名
   */
  customClass?: string
}
