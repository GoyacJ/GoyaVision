import apiClient from './client'

export interface AIModel {
  id: string
  name: string
  provider: 'openai' | 'anthropic' | 'ollama' | 'local' | 'custom'
  endpoint: string
  model_name: string
  config: Record<string, any>
  status: 'active' | 'disabled'
  created_at: string
  updated_at: string
}

export interface AIModelCreateReq {
  name: string
  provider: string
  endpoint: string
  api_key: string
  model_name: string
  config: Record<string, any>
}

export interface AIModelUpdateReq {
  name?: string
  provider?: string
  endpoint?: string
  api_key?: string
  model_name?: string
  config?: Record<string, any>
  status?: string
}

export interface AIModelListQuery {
  keyword?: string
  limit?: number
  offset?: number
  page?: number
  page_size?: number
}

export interface AIModelListRes {
  items: AIModel[]
  total: number
}

export const aiModelApi = {
  list(params?: AIModelListQuery) {
    const query = { ...params }
    if (query.page && query.page_size) {
      query.limit = query.page_size
      query.offset = (query.page - 1) * query.page_size
      delete query.page
      delete query.page_size
    }
    return apiClient.get<AIModelListRes>('/ai-models', { params: query })
  },
  
  create(data: AIModelCreateReq) {
    return apiClient.post<AIModel>('/ai-models', data)
  },

  get(id: string) {
    return apiClient.get<AIModel>(`/ai-models/${id}`)
  },

  update(id: string, data: AIModelUpdateReq) {
    return apiClient.put<AIModel>(`/ai-models/${id}`, data)
  },

  delete(id: string) {
    return apiClient.delete(`/ai-models/${id}`)
  }
}
