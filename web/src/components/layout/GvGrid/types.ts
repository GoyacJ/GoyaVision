/**
 * GvGrid 组件类型定义
 * 基于 CSS Grid 的网格布局
 */

/**
 * 网格列数
 */
export type GridCols = 1 | 2 | 3 | 4 | 5 | 6 | 12

/**
 * 间距大小
 */
export type GapSize = 'none' | 'xs' | 'sm' | 'md' | 'lg' | 'xl'

/**
 * GvGrid Props
 */
export interface GridProps {
  /**
   * 列数
   * @default 3
   */
  cols?: GridCols
  
  /**
   * 响应式列数配置
   * @example { xs: 1, sm: 2, md: 3, lg: 4 }
   */
  responsive?: {
    xs?: GridCols
    sm?: GridCols
    md?: GridCols
    lg?: GridCols
    xl?: GridCols
    '2xl'?: GridCols
  }
  
  /**
   * 间距大小
   * @default 'md'
   */
  gap?: GapSize
  
  /**
   * 水平间距（覆盖 gap）
   */
  gapX?: GapSize
  
  /**
   * 垂直间距（覆盖 gap）
   */
  gapY?: GapSize
  
  /**
   * 是否自动填充列
   * @default false
   */
  autoFill?: boolean
  
  /**
   * 自动填充时的最小列宽
   * @default '200px'
   */
  minColWidth?: string
  
  /**
   * 对齐方式
   */
  align?: 'start' | 'center' | 'end' | 'stretch'
  
  /**
   * 垂直对齐方式
   */
  justify?: 'start' | 'center' | 'end' | 'between' | 'around' | 'evenly'
}
