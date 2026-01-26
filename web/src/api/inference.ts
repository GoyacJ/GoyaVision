import apiClient from './client'

export interface InferenceResult {
  id: string
  algorithm_binding_id: string
  stream_id: string
  ts: string
  frame_ref: string
  output: any
  latency_ms?: number
  created_at: string
}

export interface InferenceResultListQuery {
  stream_id?: string
  binding_id?: string
  from?: number
  to?: number
  limit?: number
  offset?: number
}

export interface InferenceResultListResponse {
  items: InferenceResult[]
  total: number
}

export const inferenceApi = {
  list: (params?: InferenceResultListQuery) => apiClient.get<InferenceResultListResponse>('/inference_results', { params })
}
