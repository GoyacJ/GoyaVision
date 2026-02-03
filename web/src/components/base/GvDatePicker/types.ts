/**
 * GvDatePicker 组件类型定义
 * Material Design 3 日期选择器规范
 */

/**
 * 日期选择器类型
 */
export type DatePickerType = 'date' | 'datetime' | 'daterange' | 'datetimerange' | 'month' | 'year'

/**
 * 日期选择器尺寸
 */
export type DatePickerSize = 'small' | 'default' | 'large'

/**
 * GvDatePicker Props
 */
export interface DatePickerProps {
  /**
   * 绑定值
   */
  modelValue?: Date | string | number | [Date, Date] | [string, string] | [number, number]
  
  /**
   * 选择器类型
   * @default 'date'
   */
  type?: DatePickerType
  
  /**
   * 尺寸
   * @default 'default'
   */
  size?: DatePickerSize
  
  /**
   * 占位文本
   */
  placeholder?: string
  
  /**
   * 范围选择时的占位文本
   */
  startPlaceholder?: string
  
  /**
   * 范围选择时的占位文本
   */
  endPlaceholder?: string
  
  /**
   * 是否禁用
   * @default false
   */
  disabled?: boolean
  
  /**
   * 是否可清空
   * @default true
   */
  clearable?: boolean
  
  /**
   * 日期格式
   */
  format?: string
  
  /**
   * 值的格式
   */
  valueFormat?: string
  
  /**
   * 禁用日期函数
   */
  disabledDate?: (date: Date) => boolean
  
  /**
   * 可选择的时间范围
   */
  shortcuts?: Array<{
    text: string
    value: Date | (() => Date)
  }>
  
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
   * 范围选择时的分隔符
   * @default '-'
   */
  rangeSeparator?: string
  
  /**
   * 是否默认显示当前时间
   * @default false
   */
  defaultTime?: Date | [Date, Date]
}

/**
 * GvDatePicker Emits
 */
export interface DatePickerEmits {
  /**
   * 值更新事件
   */
  (e: 'update:modelValue', value: any): void
  
  /**
   * 值发生变化时触发
   */
  (e: 'change', value: any): void
  
  /**
   * 用户清空日期时触发
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
