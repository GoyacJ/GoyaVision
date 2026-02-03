/**
 * SearchBar 业务组件类型定义
 * 搜索栏组件
 */

/**
 * SearchBar Props
 */
export interface SearchBarProps {
  /**
   * 搜索值
   */
  modelValue: string
  
  /**
   * 占位文本
   * @default '搜索...'
   */
  placeholder?: string
  
  /**
   * 是否显示搜索按钮
   * @default true
   */
  showButton?: boolean
  
  /**
   * 搜索按钮文本
   * @default '搜索'
   */
  buttonText?: string
  
  /**
   * 是否正在搜索
   * @default false
   */
  loading?: boolean
  
  /**
   * 是否禁用
   * @default false
   */
  disabled?: boolean
  
  /**
   * 尺寸
   * @default 'default'
   */
  size?: 'small' | 'default' | 'large'
  
  /**
   * 是否可清空
   * @default true
   */
  clearable?: boolean
  
  /**
   * 是否在输入时立即搜索
   * @default false
   */
  immediate?: boolean
  
  /**
   * 输入防抖延迟（毫秒）
   * @default 300
   */
  debounce?: number
}

/**
 * SearchBar Emits
 */
export interface SearchBarEmits {
  /**
   * 值更新事件
   */
  (e: 'update:modelValue', value: string): void
  
  /**
   * 搜索事件
   */
  (e: 'search', value: string): void
  
  /**
   * 清空事件
   */
  (e: 'clear'): void
}
