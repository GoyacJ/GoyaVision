import apiClient from './client'

export interface WorkflowNode {
  id: string
  workflow_id: string
  node_key: string
  name: string
  operator_id?: string
  config?: Record<string, any>
  position_x: number
  position_y: number
  created_at: string
  updated_at: string
}

export interface WorkflowEdge {
  id: string
  workflow_id: string
  source_node: string
  target_node: string
  condition?: Record<string, any>
  created_at: string
  updated_at: string
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

export interface WorkflowCreateReq {
  code: string
  name: string
  description?: string
  version?: string
  trigger_type: 'manual' | 'schedule' | 'event'
  trigger_conf?: Record<string, any>
  tags?: string[]
  visibility?: number
  nodes?: Omit<WorkflowNode, 'id' | 'workflow_id' | 'created_at' | 'updated_at'>[]
  edges?: Omit<WorkflowEdge, 'id' | 'workflow_id' | 'created_at' | 'updated_at'>[]
}

export interface WorkflowUpdateReq {
  name?: string
  description?: string
  version?: string
  trigger_type?: 'manual' | 'schedule' | 'event'
  trigger_conf?: Record<string, any>
  tags?: string[]
  visibility?: number
  nodes?: Omit<WorkflowNode, 'id' | 'workflow_id' | 'created_at' | 'updated_at'>[]
  edges?: Omit<WorkflowEdge, 'id' | 'workflow_id' | 'created_at' | 'updated_at'>[]
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
