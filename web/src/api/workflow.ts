import apiClient from './client'

export interface WorkflowNode {
  id: string
  node_key: string
  node_type: string
  operator_id?: string
  config?: Record<string, any>
  position?: Record<string, any>
}

export interface WorkflowEdge {
  id: string
  source_key: string
  target_key: string
  condition?: Record<string, any>
}

export interface Workflow {
  id: string
  code: string
  name: string
  description?: string
  version: string
  trigger_type: 'manual' | 'schedule' | 'event'
  trigger_conf?: Record<string, any>
  status: 'draft' | 'testing' | 'published' | 'archived'
  tags?: string[]
  visibility?: number
  visible_role_ids?: string[]
  created_at: string
  updated_at: string
  nodes?: WorkflowNode[]
  edges?: WorkflowEdge[]
}

export interface WorkflowListQuery {
  trigger_type?: 'manual' | 'schedule' | 'event'
  status?: 'draft' | 'testing' | 'published' | 'archived'
  tags?: string[]
  keyword?: string
  page?: number
  page_size?: number
}

export interface WorkflowNodeInput {
  node_key: string
  node_type: string
  operator_id?: string
  config?: Record<string, any>
  position?: Record<string, any>
}

export interface WorkflowEdgeInput {
  source_key: string
  target_key: string
  condition?: Record<string, any>
}

export interface WorkflowCreateReq {
  code: string
  name: string
  description?: string
  version?: string
  trigger_type: 'manual' | 'schedule' | 'event'
  trigger_conf?: Record<string, any>
  status?: string
  tags?: string[]
  visibility?: number
  nodes?: WorkflowNodeInput[]
  edges?: WorkflowEdgeInput[]
}

export interface WorkflowUpdateReq {
  name?: string
  description?: string
  version?: string
  trigger_type?: 'manual' | 'schedule' | 'event'
  trigger_conf?: Record<string, any>
  status?: string
  tags?: string[]
  visibility?: number
  nodes?: WorkflowNodeInput[]
  edges?: WorkflowEdgeInput[]
}

export interface WorkflowListResponse {
  items: Workflow[]
  total: number
  page: number
  page_size: number
}

export const workflowApi = {
  list(params?: WorkflowListQuery) {
    return apiClient.get<WorkflowListResponse>('/workflows', { params })
  },

  get(id: string, withNodes: boolean = false) {
    return apiClient.get<Workflow>(`/workflows/${id}`, { params: { with_nodes: withNodes } })
  },

  create(data: WorkflowCreateReq) {
    return apiClient.post<Workflow>('/workflows', data)
  },

  update(id: string, data: WorkflowUpdateReq) {
    return apiClient.put<Workflow>(`/workflows/${id}`, data)
  },

  delete(id: string) {
    return apiClient.delete(`/workflows/${id}`)
  },

  enable(id: string) {
    return apiClient.post<Workflow>(`/workflows/${id}/enable`)
  },

  disable(id: string) {
    return apiClient.post<Workflow>(`/workflows/${id}/disable`)
  },

  trigger(id: string, assetId?: string) {
    return apiClient.post(`/workflows/${id}/trigger`, { asset_id: assetId })
  }
}
