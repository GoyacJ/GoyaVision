/**
 * GvFlex 组件类型定义
 * 基于 CSS Flexbox 的弹性布局
 */

import type { GapSize } from '../GvGrid/types'

export type { GapSize }

/**
 * Flex 方向
 */
export type FlexDirection = 'row' | 'row-reverse' | 'col' | 'col-reverse'

/**
 * 对齐方式
 */
export type AlignItems = 'start' | 'center' | 'end' | 'stretch' | 'baseline'

/**
 * 主轴对齐
 */
export type JustifyContent = 'start' | 'center' | 'end' | 'between' | 'around' | 'evenly'

/**
 * 换行方式
 */
export type FlexWrap = 'nowrap' | 'wrap' | 'wrap-reverse'

/**
 * GvFlex Props
 */
export interface FlexProps {
  /**
   * Flex 方向
   * @default 'row'
   */
  direction?: FlexDirection
  
  /**
   * 是否换行
   * @default 'nowrap'
   */
  wrap?: FlexWrap
  
  /**
   * 主轴对齐方式
   * @default 'start'
   */
  justify?: JustifyContent
  
  /**
   * 交叉轴对齐方式
   * @default 'stretch'
   */
  align?: AlignItems
  
  /**
   * 间距大小
   * @default 'md'
   */
  gap?: GapSize
  
  /**
   * 是否内联 flex
   * @default false
   */
  inline?: boolean
  
  /**
   * 是否垂直布局（direction="col" 的快捷方式）
   * @default false
   */
  vertical?: boolean
}
