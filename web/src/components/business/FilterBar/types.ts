/**
 * FilterBar 业务组件类型定义
 * 筛选栏组件
 */

/**
 * 筛选字段类型
 */
export type FilterFieldType = 'input' | 'select' | 'daterange' | 'date' | 'datetime'

/**
 * 筛选字段配置
 */
export interface FilterField {
  /**
   * 字段 key
   */
  key: string
  
  /**
   * 字段标签
   */
  label: string
  
  /**
   * 字段类型
   * @default 'input'
   */
  type?: FilterFieldType
  
  /**
   * 占位文本
   */
  placeholder?: string
  
  /**
   * 选项（type 为 select 时）
   */
  options?: Array<{ label: string; value: any }>
  
  /**
   * 默认值
   */
  defaultValue?: any
  
  /**
   * 起始占位文本（type 为 daterange 时）
   */
  startPlaceholder?: string
  
  /**
   * 结束占位文本（type 为 daterange 时）
   */
  endPlaceholder?: string
}

/**
 * FilterBar Props
 */
export interface FilterBarProps {
  /**
   * 筛选字段配置
   */
  fields: FilterField[]
  
  /**
   * 筛选值
   */
  modelValue: Record<string, any>
  
  /**
   * 是否显示重置按钮
   * @default true
   */
  showReset?: boolean
  
  /**
   * 重置按钮文本
   * @default '重置'
   */
  resetText?: string
  
  /**
   * 筛选按钮文本
   * @default '筛选'
   */
  filterText?: string
  
  /**
   * 是否正在筛选
   * @default false
   */
  loading?: boolean
  
  /**
   * 是否可折叠
   * @default false
   */
  collapsible?: boolean
  
  /**
   * 默认是否展开（collapsible 为 true 时有效）
   * @default true
   */
  defaultExpanded?: boolean
  
  /**
   * 每行显示的字段数
   * @default 3
   */
  columns?: number
}

/**
 * FilterBar Emits
 */
export interface FilterBarEmits {
  /**
   * 值更新事件
   */
  (e: 'update:modelValue', value: Record<string, any>): void
  
  /**
   * 筛选事件
   */
  (e: 'filter', value: Record<string, any>): void
  
  /**
   * 重置事件
   */
  (e: 'reset'): void
}
