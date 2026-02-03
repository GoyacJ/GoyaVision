/**
 * GvInput 组件类型定义
 * Material Design 3 输入框规范
 */

import type { Component } from 'vue'

/**
 * 输入框尺寸
 */
export type InputSize = 'small' | 'medium' | 'large'

/**
 * 输入框类型
 */
export type InputType = 'text' | 'password' | 'number' | 'email' | 'tel' | 'url' | 'search'

/**
 * 输入框 Props
 */
export interface InputProps {
  /**
   * 绑定值
   */
  modelValue?: string | number
  
  /**
   * 输入框类型
   * @default 'text'
   */
  type?: InputType
  
  /**
   * 输入框尺寸
   * @default 'medium'
   */
  size?: InputSize
  
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
   * 是否只读
   * @default false
   */
  readonly?: boolean
  
  /**
   * 是否显示清除按钮
   * @default false
   */
  clearable?: boolean
  
  /**
   * 是否显示密码切换按钮
   * @default false
   */
  showPassword?: boolean
  
  /**
   * 前置图标
   */
  prefixIcon?: string | Component
  
  /**
   * 后置图标
   */
  suffixIcon?: string | Component
  
  /**
   * 最大输入长度
   */
  maxlength?: number
  
  /**
   * 是否显示字数统计
   * @default false
   */
  showCount?: boolean
  
  /**
   * 验证状态
   */
  status?: 'success' | 'error' | 'warning'
  
  /**
   * 错误提示信息
   */
  errorMessage?: string
  
  /**
   * 是否自动聚焦
   * @default false
   */
  autofocus?: boolean
  
  /**
   * 原生 autocomplete 属性
   */
  autocomplete?: string
  
  /**
   * 原生 name 属性
   */
  name?: string
  
  /**
   * 原生 form 属性
   */
  form?: string
  
  /**
   * 是否必填（显示红色星号）
   * @default false
   */
  required?: boolean
  
  /**
   * 标签文本
   */
  label?: string
}

/**
 * 输入框 Emits
 */
export interface InputEmits {
  /**
   * 值更新事件
   */
  (e: 'update:modelValue', value: string | number): void
  
  /**
   * 输入事件
   */
  (e: 'input', value: string | number): void
  
  /**
   * 变化事件
   */
  (e: 'change', value: string | number): void
  
  /**
   * 获得焦点事件
   */
  (e: 'focus', event: FocusEvent): void
  
  /**
   * 失去焦点事件
   */
  (e: 'blur', event: FocusEvent): void
  
  /**
   * 清除事件
   */
  (e: 'clear'): void
  
  /**
   * 按键事件
   */
  (e: 'keydown', event: KeyboardEvent): void
  
  /**
   * 按键抬起事件
   */
  (e: 'keyup', event: KeyboardEvent): void
}
