import apiClient from './client'

export interface Operator {
  id: string
  code: string
  name: string
  description?: string
  category: 'analysis' | 'processing' | 'generation'
  type: string
  version: string
  endpoint: string
  method: string
  input_schema?: Record<string, any>
  output_spec?: Record<string, any>
  config?: Record<string, any>
  status: 'draft' | 'testing' | 'published' | 'deprecated'
  is_builtin: boolean
  tags?: string[]
  created_at: string
  updated_at: string
}

export interface OperatorListQuery {
  category?: 'analysis' | 'processing' | 'generation'
  type?: string
  status?: 'draft' | 'testing' | 'published' | 'deprecated'
  is_builtin?: boolean
  tags?: string[]
  keyword?: string
  page?: number
  page_size?: number
}

export interface OperatorCreateReq {
  code: string
  name: string
  description?: string
  category: 'analysis' | 'processing' | 'generation'
  type: string
  version?: string
  endpoint: string
  method?: string
  input_schema?: Record<string, any>
  output_spec?: Record<string, any>
  config?: Record<string, any>
  tags?: string[]
}

export interface OperatorUpdateReq {
  name?: string
  description?: string
  category?: 'analysis' | 'processing' | 'generation'
  type?: string
  version?: string
  endpoint?: string
  method?: string
  input_schema?: Record<string, any>
  output_spec?: Record<string, any>
  config?: Record<string, any>
  tags?: string[]
}

export interface OperatorListResponse {
  items: Operator[]
  total: number
  page: number
  page_size: number
}

export const operatorApi = {
  list(params?: OperatorListQuery) {
    return apiClient.get<OperatorListResponse>('/operators', { params })
  },

  get(id: string) {
    return apiClient.get<Operator>(`/operators/${id}`)
  },

  create(data: OperatorCreateReq) {
    return apiClient.post<Operator>('/operators', data)
  },

  update(id: string, data: OperatorUpdateReq) {
    return apiClient.put<Operator>(`/operators/${id}`, data)
  },

  delete(id: string) {
    return apiClient.delete(`/operators/${id}`)
  },

  enable(id: string) {
    return apiClient.post<Operator>(`/operators/${id}/enable`)
  },

  disable(id: string) {
    return apiClient.post<Operator>(`/operators/${id}/disable`)
  },

  listByCategory(category: string, params?: OperatorListQuery) {
    return apiClient.get<OperatorListResponse>(`/operators/category/${category}`, { params })
  }
}
