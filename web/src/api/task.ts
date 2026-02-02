import apiClient from './client'
import type { Workflow } from './workflow'
import type { MediaAsset } from './asset'

export interface Task {
  id: string
  workflow_id: string
  asset_id?: string
  status: 'pending' | 'running' | 'success' | 'failed' | 'cancelled'
  progress: number
  current_node?: string
  input_params?: Record<string, any>
  error?: string
  started_at?: string
  completed_at?: string
  created_at: string
  updated_at: string
  workflow?: Workflow
  asset?: MediaAsset
}

export interface TaskListQuery {
  workflow_id?: string
  asset_id?: string
  status?: 'pending' | 'running' | 'success' | 'failed' | 'cancelled'
  page?: number
  page_size?: number
}

export interface TaskCreateReq {
  workflow_id: string
  asset_id?: string
  input_params?: Record<string, any>
}

export interface TaskUpdateReq {
  status?: 'pending' | 'running' | 'success' | 'failed' | 'cancelled'
  progress?: number
  current_node?: string
  error?: string
}

export interface TaskStats {
  total: number
  pending: number
  running: number
  success: number
  failed: number
  cancelled: number
}

export interface TaskListResponse {
  items: Task[]
  total: number
  page: number
  page_size: number
}

export const taskApi = {
  list(params?: TaskListQuery) {
    return apiClient.get<TaskListResponse>('/tasks', { params })
  },

  get(id: string, withRelations: boolean = false) {
    return apiClient.get<Task>(`/tasks/${id}`, { params: { with_relations: withRelations } })
  },

  create(data: TaskCreateReq) {
    return apiClient.post<Task>('/tasks', data)
  },

  update(id: string, data: TaskUpdateReq) {
    return apiClient.put<Task>(`/tasks/${id}`, data)
  },

  delete(id: string) {
    return apiClient.delete(`/tasks/${id}`)
  },

  start(id: string) {
    return apiClient.post<Task>(`/tasks/${id}/start`)
  },

  complete(id: string) {
    return apiClient.post<Task>(`/tasks/${id}/complete`)
  },

  fail(id: string, error: string) {
    return apiClient.post<Task>(`/tasks/${id}/fail`, { error })
  },

  cancel(id: string) {
    return apiClient.post<Task>(`/tasks/${id}/cancel`)
  },

  getStats() {
    return apiClient.get<TaskStats>('/tasks/stats')
  }
}
