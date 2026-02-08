import apiClient from './client'

export type AlgorithmStatus = 'draft' | 'published' | 'deprecated'
export type AlgorithmVersionStatus = 'draft' | 'tested' | 'published' | 'archived'
export type AlgorithmImplementationType = 'operator_version' | 'mcp_tool' | 'ai_chain'
export type AlgorithmSelectionPolicy = 'stable' | 'high_quality' | 'low_cost'

export interface AlgorithmImplementation {
  id?: string
  name?: string
  type?: AlgorithmImplementationType
  binding_ref: string
  config?: Record<string, any>
  latency_ms?: number
  cost_score?: number
  quality_score?: number
  tier?: string
  is_default?: boolean
}

export interface AlgorithmEvaluation {
  id?: string
  dataset_ref?: string
  metrics?: Record<string, number>
  report_artifact_id?: string
  summary?: string
}

export interface AlgorithmVersion {
  id?: string
  version: string
  status?: AlgorithmVersionStatus
  selection_policy?: AlgorithmSelectionPolicy
  default_implementation?: string
  implementations: AlgorithmImplementation[]
  evaluations?: AlgorithmEvaluation[]
}

export interface Algorithm {
  id: string
  code: string
  name: string
  description?: string
  scenario?: string
  status: AlgorithmStatus
  tags?: string[]
  visibility?: number
  visible_role_ids?: string[]
  versions?: AlgorithmVersion[]
  created_at: string
  updated_at: string
}

export interface AlgorithmListResponse {
  items: Algorithm[]
  total: number
}

export interface AlgorithmListQuery {
  status?: AlgorithmStatus
  scenario?: string
  tags?: string[]
  keyword?: string
  page?: number
  page_size?: number
}

export interface AlgorithmCreateReq {
  code: string
  name: string
  description?: string
  scenario?: string
  status?: AlgorithmStatus
  tags?: string[]
  visibility?: number
  visible_role_ids?: string[]
  initial_version?: AlgorithmVersion
}

export interface AlgorithmUpdateReq {
  name?: string
  description?: string
  scenario?: string
  status?: AlgorithmStatus
  tags?: string[]
  visibility?: number
  visible_role_ids?: string[]
}

export interface CreateAlgorithmVersionReq {
  version: string
  status?: AlgorithmVersionStatus
  selection_policy?: AlgorithmSelectionPolicy
  implementations: AlgorithmImplementation[]
  evaluations?: AlgorithmEvaluation[]
}

function toListParams(params?: { page?: number; page_size?: number; [key: string]: any }) {
  if (!params) return undefined
  const { page, page_size, ...rest } = params
  const normalized: Record<string, any> = { ...rest }
  if (Array.isArray(normalized.tags)) {
    normalized.tags = normalized.tags.join(',')
  }
  if (typeof page_size === 'number' && page_size > 0) {
    normalized.limit = page_size
  }
  if (typeof page === 'number' && page > 0 && typeof page_size === 'number' && page_size > 0) {
    normalized.offset = (page - 1) * page_size
  }
  return normalized
}

export const algorithmApi = {
  list(params?: AlgorithmListQuery) {
    return apiClient.get<AlgorithmListResponse>('/algorithms', { params: toListParams(params) })
  },

  get(id: string) {
    return apiClient.get<Algorithm>(`/algorithms/${id}`)
  },

  create(data: AlgorithmCreateReq) {
    return apiClient.post<Algorithm>('/algorithms', data)
  },

  update(id: string, data: AlgorithmUpdateReq) {
    return apiClient.put<Algorithm>(`/algorithms/${id}`, data)
  },

  delete(id: string) {
    return apiClient.delete(`/algorithms/${id}`)
  },

  createVersion(id: string, data: CreateAlgorithmVersionReq) {
    return apiClient.post<AlgorithmVersion>(`/algorithms/${id}/versions`, data)
  },

  publishVersion(id: string, versionID: string) {
    return apiClient.post<AlgorithmVersion>(`/algorithms/${id}/versions/${versionID}/publish`)
  },
}
