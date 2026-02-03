/**
 * GvUpload 组件类型定义
 */

import type { UploadFile, UploadFiles } from 'element-plus'

export interface GvUploadProps {
  /**
   * 文件列表（v-model）
   */
  modelValue?: UploadFile[]

  /**
   * 是否自动上传
   * @default false
   */
  autoUpload?: boolean

  /**
   * 是否禁用
   * @default false
   */
  disabled?: boolean

  /**
   * 最大上传数量
   * @default 1
   */
  limit?: number

  /**
   * 接受的文件类型
   * @example "image/*" | ".jpg,.png" | "video/*"
   */
  accept?: string

  /**
   * 按钮文字
   * @default "选择文件"
   */
  buttonText?: string

  /**
   * 上传中按钮文字
   * @default "上传中..."
   */
  uploadingText?: string

  /**
   * 提示文字
   */
  tip?: string

  /**
   * 按钮样式变体
   * @default "outlined"
   */
  variant?: 'filled' | 'outlined' | 'text'

  /**
   * 按钮大小
   * @default "medium"
   */
  size?: 'small' | 'medium' | 'large'

  /**
   * 是否显示文件列表
   * @default true
   */
  showFileList?: boolean

  /**
   * 最大文件大小（字节）
   */
  maxSize?: number

  /**
   * 上传前的钩子
   */
  beforeUpload?: (file: File) => boolean | Promise<boolean>
}

export interface GvUploadEmits {
  /**
   * 文件列表更新
   */
  'update:modelValue': (files: UploadFile[]) => void

  /**
   * 文件状态改变
   */
  'change': (file: UploadFile, fileList: UploadFiles) => void

  /**
   * 文件移除
   */
  'remove': (file: UploadFile, fileList: UploadFiles) => void

  /**
   * 上传成功
   */
  'success': (response: any, file: UploadFile, fileList: UploadFiles) => void

  /**
   * 上传失败
   */
  'error': (error: Error, file: UploadFile, fileList: UploadFiles) => void

  /**
   * 上传进度
   */
  'progress': (event: any, file: UploadFile, fileList: UploadFiles) => void
}

export interface GvUploadExpose {
  /**
   * 手动上传文件列表
   */
  submit: () => void

  /**
   * 清空文件列表
   */
  clearFiles: () => void

  /**
   * 取消上传
   */
  abort: () => void

  /**
   * 当前文件列表
   */
  fileList: UploadFile[]
}
