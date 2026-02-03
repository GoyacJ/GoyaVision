/**
 * GvSelect 组件类型定义
 * Material Design 3 选择器规范
 */

/**
 * 选择器尺寸
 */
export type SelectSize = 'small' | 'medium' | 'large'

/**
 * 选项数据
 */
export interface SelectOption {
  label: string
  value: string | number
  disabled?: boolean
  [key: string]: any
}

/**
 * 选择器 Props
 */
export interface SelectProps {
  /**
   * 绑定值
   */
  modelValue?: string | number | Array<string | number>
  
  /**
   * 选项数据
   */
  options?: SelectOption[]
  
  /**
   * 选择器尺寸
   * @default 'medium'
   */
  size?: SelectSize
  
  /**
   * 占位文本
   */
  placeholder?: string
  
  /**
   * 是否禁用
   * @default false
   */
  disabled?: boolean
  
  /**
   * 是否可清空
   * @default false
   */
  clearable?: boolean
  
  /**
   * 是否多选
   * @default false
   */
  multiple?: boolean
  
  /**
   * 多选时用户最多可以选择的项目数
   */
  multipleLimit?: number
  
  /**
   * 是否可搜索
   * @default false
   */
  filterable?: boolean
  
  /**
   * 是否允许创建新条目
   * @default false
   */
  allowCreate?: boolean
  
  /**
   * 自定义搜索方法
   */
  filterMethod?: (query: string, option: SelectOption) => boolean
  
  /**
   * 是否远程搜索
   * @default false
   */
  remote?: boolean
  
  /**
   * 远程搜索方法
   */
  remoteMethod?: (query: string) => void
  
  /**
   * 是否正在从远程获取数据
   * @default false
   */
  loading?: boolean
  
  /**
   * 远程加载时显示的文字
   */
  loadingText?: string
  
  /**
   * 无数据时显示的文字
   */
  noDataText?: string
  
  /**
   * 无匹配数据时显示的文字
   */
  noMatchText?: string
  
  /**
   * 下拉框的类名
   */
  popperClass?: string
  
  /**
   * 验证状态
   */
  status?: 'success' | 'error' | 'warning'
  
  /**
   * 错误提示信息
   */
  errorMessage?: string
  
  /**
   * 是否必填（显示红色星号）
   * @default false
   */
  required?: boolean
  
  /**
   * 标签文本
   */
  label?: string
  
  /**
   * 选项的标签字段名
   * @default 'label'
   */
  labelKey?: string
  
  /**
   * 选项的值字段名
   * @default 'value'
   */
  valueKey?: string
}

/**
 * 选择器 Emits
 */
export interface SelectEmits {
  /**
   * 值更新事件
   */
  (e: 'update:modelValue', value: string | number | Array<string | number> | undefined): void
  
  /**
   * 选中值发生变化时触发
   */
  (e: 'change', value: string | number | Array<string | number> | undefined): void
  
  /**
   * 下拉框出现/隐藏时触发
   */
  (e: 'visible-change', visible: boolean): void
  
  /**
   * 多选模式下移除 tag 时触发
   */
  (e: 'remove-tag', value: string | number): void
  
  /**
   * 可清空的单选模式下用户点击清空按钮时触发
   */
  (e: 'clear'): void
  
  /**
   * 获得焦点时触发
   */
  (e: 'focus', event: FocusEvent): void
  
  /**
   * 失去焦点时触发
   */
  (e: 'blur', event: FocusEvent): void
}
