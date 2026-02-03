/**
 * GvContainer 组件类型定义
 */

/**
 * 容器最大宽度
 */
export type ContainerMaxWidth = 'sm' | 'md' | 'lg' | 'xl' | '2xl' | 'full'

/**
 * 容器 Props
 */
export interface ContainerProps {
  /**
   * 最大宽度
   * @default 'xl'
   */
  maxWidth?: ContainerMaxWidth
  
  /**
   * 是否添加水平内边距
   * @default true
   */
  padding?: boolean
  
  /**
   * 是否居中对齐
   * @default true
   */
  centered?: boolean
}
