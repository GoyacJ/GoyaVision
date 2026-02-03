/**
 * TaskCard 业务组件类型定义
 * 任务卡片组件
 */

/**
 * 任务数据
 */
export interface Task {
  id: string | number
  name: string
  status: string
  type?: string
  description?: string
  progress?: number
  startTime?: string
  endTime?: string
  createdAt?: string
  [key: string]: any
}

/**
 * TaskCard Props
 */
export interface TaskCardProps {
  /**
   * 任务数据
   */
  task: Task
  
  /**
   * 是否显示进度条
   * @default true
   */
  showProgress?: boolean
  
  /**
   * 是否显示操作按钮
   * @default true
   */
  showActions?: boolean
  
  /**
   * 是否显示查看按钮
   * @default true
   */
  showView?: boolean
  
  /**
   * 是否显示取消按钮
   * @default true
   */
  showCancel?: boolean
}

/**
 * TaskCard Emits
 */
export interface TaskCardEmits {
  /**
   * 点击卡片
   */
  (e: 'click', task: Task): void
  
  /**
   * 点击查看
   */
  (e: 'view', task: Task): void
  
  /**
   * 点击取消
   */
  (e: 'cancel', task: Task): void
}
