import apiClient from './client'

export interface Stream {
  id: string
  url: string
  name: string
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface StreamCreateReq {
  url: string
  name: string
  enabled?: boolean
}

export interface StreamUpdateReq {
  url?: string
  name?: string
  enabled?: boolean
}

export const streamApi = {
  list: (enabled?: boolean) => {
    const params = enabled !== undefined ? { enabled } : {}
    return apiClient.get<Stream[]>('/streams', { params })
  },
  get: (id: string) => apiClient.get<Stream>(`/streams/${id}`),
  create: (data: StreamCreateReq) => apiClient.post<Stream>('/streams', data),
  update: (id: string, data: StreamUpdateReq) => apiClient.put<Stream>(`/streams/${id}`, data),
  delete: (id: string) => apiClient.delete(`/streams/${id}`),
  startPreview: (id: string) => apiClient.get<{ hls_url: string }>(`/streams/${id}/preview/start`),
  stopPreview: (id: string) => apiClient.post(`/streams/${id}/preview/stop`),
  startRecord: (id: string) => apiClient.post<{ session_id: string }>(`/streams/${id}/record/start`),
  stopRecord: (id: string) => apiClient.post(`/streams/${id}/record/stop`),
  listRecordSessions: (id: string) => apiClient.get(`/streams/${id}/record/sessions`)
}
