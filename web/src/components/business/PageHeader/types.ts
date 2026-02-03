/**
 * PageHeader 业务组件类型定义
 * 页面头部组件
 */

/**
 * 面包屑项
 */
export interface BreadcrumbItem {
  label: string
  to?: string
}

/**
 * PageHeader Props
 */
export interface PageHeaderProps {
  /**
   * 页面标题
   */
  title: string
  
  /**
   * 页面描述
   */
  description?: string
  
  /**
   * 面包屑导航
   */
  breadcrumb?: BreadcrumbItem[]
  
  /**
   * 是否显示返回按钮
   * @default false
   */
  showBack?: boolean
  
  /**
   * 返回按钮文本
   * @default '返回'
   */
  backText?: string
}

/**
 * PageHeader Emits
 */
export interface PageHeaderEmits {
  /**
   * 返回按钮点击事件
   */
  (e: 'back'): void
}
