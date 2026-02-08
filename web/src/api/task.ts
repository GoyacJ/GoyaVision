import apiClient from './client'
import type { Workflow } from './workflow'
import type { MediaAsset } from './asset'

export interface Task {
  id: string
  workflow_id: string
  workflow_revision_id?: string
  workflow_revision?: number
  asset_id?: string
  status: 'pending' | 'running' | 'success' | 'failed' | 'cancelled'
  progress: number
  current_node?: string
  context_version?: number
  input_params?: Record<string, any>
  error?: string
  started_at?: string
  completed_at?: string
  created_at: string
  updated_at: string
  workflow?: Workflow
  asset?: MediaAsset
  node_executions?: any[]
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

export interface TaskContext {
  task_id: string
  version: number
  data?: Record<string, any>
  updated_at: string
}

export interface TaskContextPatch {
  id: string
  task_id: string
  writer_node_key: string
  before_version: number
  after_version: number
  diff: {
    set?: Record<string, any>
    unset?: string[]
  }
  created_at: string
}

export interface TaskContextPatchListResponse {
  items: TaskContextPatch[]
  total: number
}

export interface TaskRunEvent {
  id: string
  task_id: string
  session_id?: string
  event_type: string
  source: string
  node_key?: string
  tool_name?: string
  payload?: Record<string, any>
  created_at: string
}

export interface TaskRunEventListResponse {
  items: TaskRunEvent[]
  total: number
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
  },

  getContext(id: string) {
    return apiClient.get<TaskContext>(`/tasks/${id}/context`)
  },

  listContextPatches(id: string, params?: { limit?: number; offset?: number }) {
    return apiClient.get<TaskContextPatchListResponse>(`/tasks/${id}/context/patches`, { params })
  },

  createContextSnapshot(id: string, trigger: string = 'manual') {
    return apiClient.post(`/tasks/${id}/context/snapshot`, { trigger })
  },

  listEvents(id: string, params?: { source?: string; node_key?: string; limit?: number; offset?: number }) {
    return apiClient.get<TaskRunEventListResponse>(`/tasks/${id}/events`, { params })
  }
}
