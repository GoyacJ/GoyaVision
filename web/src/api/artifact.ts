import apiClient from './client'
import type { Task } from './task'
import type { MediaAsset } from './asset'

export interface Artifact {
  id: string
  task_id: string
  type: 'asset' | 'result' | 'timeline' | 'report'
  asset_id?: string
  data?: Record<string, any>
  created_at: string
  updated_at: string
  task?: Task
  asset?: MediaAsset
}

export interface ArtifactListQuery {
  task_id?: string
  node_key?: string
  type?: 'asset' | 'result' | 'timeline' | 'report'
  asset_id?: string
  page?: number
  page_size?: number
}

export interface ArtifactCreateReq {
  task_id: string
  type: 'asset' | 'result' | 'timeline' | 'report'
  asset_id?: string
  data?: Record<string, any>
}

export interface ArtifactListResponse {
  items: Artifact[]
  total: number
  page: number
  page_size: number
}

export const artifactApi = {
  list(params?: ArtifactListQuery) {
    return apiClient.get<ArtifactListResponse>('/artifacts', { params })
  },

  get(id: string) {
    return apiClient.get<Artifact>(`/artifacts/${id}`)
  },

  create(data: ArtifactCreateReq) {
    return apiClient.post<Artifact>('/artifacts', data)
  },

  delete(id: string) {
    return apiClient.delete(`/artifacts/${id}`)
  },

  listByTask(taskId: string, type?: string) {
    const params: any = {}
    if (type) params.type = type
    return apiClient.get<ArtifactListResponse>(`/tasks/${taskId}/artifacts`, { params })
  }
}
