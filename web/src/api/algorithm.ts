import apiClient from './client'

export interface Algorithm {
  id: string
  name: string
  endpoint: string
  input_spec: any
  output_spec: any
  created_at: string
  updated_at: string
}

export interface AlgorithmCreateReq {
  name: string
  endpoint: string
  input_spec?: any
  output_spec?: any
}

export interface AlgorithmUpdateReq {
  name?: string
  endpoint?: string
  input_spec?: any
  output_spec?: any
}

export const algorithmApi = {
  list: () => apiClient.get<Algorithm[]>('/algorithms'),
  get: (id: string) => apiClient.get<Algorithm>(`/algorithms/${id}`),
  create: (data: AlgorithmCreateReq) => apiClient.post<Algorithm>('/algorithms', data),
  update: (id: string, data: AlgorithmUpdateReq) => apiClient.put<Algorithm>(`/algorithms/${id}`, data),
  delete: (id: string) => apiClient.delete(`/algorithms/${id}`)
}
