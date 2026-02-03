/**
 * GvTable 组件类型定义
 * Material Design 3 表格规范
 */

/**
 * 表格尺寸
 */
export type TableSize = 'small' | 'default' | 'large'

/**
 * 列定义
 */
export interface TableColumn {
  /**
   * 列的 key（对应数据字段名）
   */
  prop: string
  
  /**
   * 列标题
   */
  label: string
  
  /**
   * 列宽度
   */
  width?: string | number
  
  /**
   * 最小宽度
   */
  minWidth?: string | number
  
  /**
   * 是否固定列
   */
  fixed?: boolean | 'left' | 'right'
  
  /**
   * 是否可排序
   */
  sortable?: boolean | 'custom'
  
  /**
   * 对齐方式
   */
  align?: 'left' | 'center' | 'right'
  
  /**
   * 自定义渲染函数
   */
  formatter?: (row: any, column: TableColumn, cellValue: any, index: number) => string
  
  /**
   * 是否显示溢出提示
   */
  showOverflowTooltip?: boolean
}

/**
 * 分页配置
 */
export interface PaginationConfig {
  /**
   * 当前页码
   */
  currentPage: number
  
  /**
   * 每页显示条数
   */
  pageSize: number
  
  /**
   * 总条数
   */
  total: number
  
  /**
   * 每页显示条数选择器的选项
   */
  pageSizes?: number[]
  
  /**
   * 分页布局
   */
  layout?: string
}

/**
 * GvTable Props
 */
export interface TableProps {
  /**
   * 表格数据
   */
  data: any[]
  
  /**
   * 列配置
   */
  columns: TableColumn[]
  
  /**
   * 表格尺寸
   * @default 'default'
   */
  size?: TableSize
  
  /**
   * 是否显示边框
   * @default false
   */
  border?: boolean
  
  /**
   * 是否显示斑马纹
   * @default true
   */
  stripe?: boolean
  
  /**
   * 是否显示表头
   * @default true
   */
  showHeader?: boolean
  
  /**
   * 是否高亮当前行
   * @default false
   */
  highlightCurrentRow?: boolean
  
  /**
   * 是否显示加载状态
   * @default false
   */
  loading?: boolean
  
  /**
   * 空数据时显示的文本
   * @default '暂无数据'
   */
  emptyText?: string
  
  /**
   * 行的 key
   * @default 'id'
   */
  rowKey?: string
  
  /**
   * 是否启用分页
   * @default false
   */
  pagination?: boolean
  
  /**
   * 分页配置
   */
  paginationConfig?: PaginationConfig
  
  /**
   * 表格高度
   */
  height?: string | number
  
  /**
   * 表格最大高度
   */
  maxHeight?: string | number
  
  /**
   * 是否可选择行
   * @default false
   */
  selectable?: boolean
  
  /**
   * 默认选中的行
   */
  defaultSelection?: any[]
}

/**
 * GvTable Emits
 */
export interface TableEmits {
  /**
   * 当选择项发生变化时触发
   */
  (e: 'selection-change', selection: any[]): void
  
  /**
   * 当某一行被点击时触发
   */
  (e: 'row-click', row: any, column: TableColumn, event: Event): void
  
  /**
   * 当某一行被双击时触发
   */
  (e: 'row-dblclick', row: any, column: TableColumn, event: Event): void
  
  /**
   * 当某个单元格被点击时触发
   */
  (e: 'cell-click', row: any, column: TableColumn, cell: any, event: Event): void
  
  /**
   * 当排序条件发生变化时触发
   */
  (e: 'sort-change', data: { column: TableColumn; prop: string; order: string }): void
  
  /**
   * 当前页码改变时触发
   */
  (e: 'current-change', page: number): void
  
  /**
   * 每页条数改变时触发
   */
  (e: 'size-change', size: number): void
}
