import apiClient from './client'

export interface AlgorithmBinding {
  id: string
  stream_id: string
  algorithm_id: string
  enabled: boolean
  interval_sec: number
  initial_delay_sec: number
  schedule: any
  config: any
  created_at: string
  updated_at: string
}

export interface AlgorithmBindingCreateReq {
  algorithm_id: string
  enabled?: boolean
  interval_sec: number
  initial_delay_sec?: number
  schedule?: any
  config?: any
}

export interface AlgorithmBindingUpdateReq {
  algorithm_id?: string
  enabled?: boolean
  interval_sec?: number
  initial_delay_sec?: number
  schedule?: any
  config?: any
}

export const algorithmBindingApi = {
  list: (streamId: string) => apiClient.get<AlgorithmBinding[]>(`/streams/${streamId}/algorithm-bindings`),
  get: (streamId: string, bindingId: string) => apiClient.get<AlgorithmBinding>(`/streams/${streamId}/algorithm-bindings/${bindingId}`),
  create: (streamId: string, data: AlgorithmBindingCreateReq) => apiClient.post<AlgorithmBinding>(`/streams/${streamId}/algorithm-bindings`, data),
  update: (streamId: string, bindingId: string, data: AlgorithmBindingUpdateReq) => apiClient.put<AlgorithmBinding>(`/streams/${streamId}/algorithm-bindings/${bindingId}`, data),
  delete: (streamId: string, bindingId: string) => apiClient.delete(`/streams/${streamId}/algorithm-bindings/${bindingId}`)
}
