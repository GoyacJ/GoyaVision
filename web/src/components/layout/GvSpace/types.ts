/**
 * GvSpace 组件类型定义
 */

/**
 * 间距大小
 */
export type SpaceSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl' | number

/**
 * 间距方向
 */
export type SpaceDirection = 'horizontal' | 'vertical'

/**
 * 对齐方式
 */
export type SpaceAlign = 'start' | 'center' | 'end' | 'baseline'

/**
 * GvSpace Props
 */
export interface SpaceProps {
  /**
   * 间距大小
   * @default 'md'
   */
  size?: SpaceSize | [SpaceSize, SpaceSize]
  
  /**
   * 间距方向
   * @default 'horizontal'
   */
  direction?: SpaceDirection
  
  /**
   * 对齐方式
   * @default 'center'
   */
  align?: SpaceAlign
  
  /**
   * 是否自动换行（仅在 horizontal 时有效）
   * @default false
   */
  wrap?: boolean
  
  /**
   * 是否填充整个父容器
   * @default false
   */
  fill?: boolean
  
  /**
   * 填充比例，仅在 fill 为 true 时有效
   */
  fillRatio?: number
}
