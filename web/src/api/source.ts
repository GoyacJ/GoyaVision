import apiClient from './client'

export interface MediaSource {
  id: string
  name: string
  path_name: string
  type: 'pull' | 'push'
  url?: string
  protocol?: string
  enabled: boolean
  record_enabled: boolean
  created_at: string
  updated_at: string
}

export interface SourceListQuery {
  type?: 'pull' | 'push'
  limit?: number
  offset?: number
}

export interface SourceCreateReq {
  name: string
  type: 'pull' | 'push'
  url?: string
  protocol?: string
  enabled?: boolean
}

export interface SourceUpdateReq {
  name?: string
  url?: string
  protocol?: string
  enabled?: boolean
}

export interface SourcePreviewResponse {
  path_name: string
  hls_url: string
  rtsp_url: string
  rtmp_url: string
  webrtc_url?: string
  push_url?: string
}

export interface SourceListResponse {
  items: MediaSource[]
  total: number
}

export const sourceApi = {
  list(params?: SourceListQuery) {
    return apiClient.get<SourceListResponse>('/sources', { params })
  },

  get(id: string) {
    return apiClient.get<MediaSource>(`/sources/${id}`)
  },

  create(data: SourceCreateReq) {
    return apiClient.post<MediaSource>('/sources', data)
  },

  update(id: string, data: SourceUpdateReq) {
    return apiClient.put<MediaSource>(`/sources/${id}`, data)
  },

  remove(id: string) {
    return apiClient.delete(`/sources/${id}`)
  },

  getPreview(id: string) {
    return apiClient.get<SourcePreviewResponse>(`/sources/${id}/preview`)
  }
}
