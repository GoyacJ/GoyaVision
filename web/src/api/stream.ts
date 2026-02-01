import apiClient from './client'

export interface Stream {
  id: string
  url: string
  name: string
  type: 'pull' | 'push'
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface StreamStatus {
  path_name: string
  ready: boolean
  online: boolean
  tracks: string[]
  bytes_received: number
  bytes_sent: number
  reader_count: number
  rtsp_url: string
  rtmp_url: string
  hls_url: string
  webrtc_url: string
}

export interface StreamWithStatus extends Stream {
  status?: StreamStatus
}

export interface StreamCreateReq {
  url?: string
  name: string
  type?: 'pull' | 'push'
  enabled?: boolean
}

export interface StreamUpdateReq {
  url?: string
  name?: string
  enabled?: boolean
}

export interface PreviewURLs {
  hls_url: string
  rtsp_url: string
  rtmp_url: string
  webrtc_url: string
}

export interface PlaybackURLs {
  hls_url: string
  mp4_url: string
}

export interface PlaybackSegment {
  start: string
  playback_url: string
}

export const streamApi = {
  list: (enabled?: boolean, withStatus?: boolean) => {
    const params: any = {}
    if (enabled !== undefined) params.enabled = enabled
    if (withStatus) params.with_status = true
    return apiClient.get<StreamWithStatus[]>('/streams', { params })
  },
  get: (id: string, withStatus?: boolean) => {
    const params = withStatus ? { with_status: true } : {}
    return apiClient.get<StreamWithStatus>(`/streams/${id}`, { params })
  },
  create: (data: StreamCreateReq) => apiClient.post<Stream>('/streams', data),
  update: (id: string, data: StreamUpdateReq) => apiClient.put<Stream>(`/streams/${id}`, data),
  delete: (id: string) => apiClient.delete(`/streams/${id}`),
  enable: (id: string) => apiClient.post(`/streams/${id}/enable`),
  disable: (id: string) => apiClient.post(`/streams/${id}/disable`),
  getStatus: (id: string) => apiClient.get<StreamStatus>(`/streams/${id}/status`),

  getPreviewURLs: (id: string) => apiClient.get<PreviewURLs>(`/streams/${id}/preview`),
  startPreview: (id: string) => apiClient.get<PreviewURLs>(`/streams/${id}/preview/start`),
  isReady: (id: string) => apiClient.get<{ ready: boolean }>(`/streams/${id}/preview/ready`),

  startRecord: (id: string) => apiClient.post<{ session_id: string }>(`/streams/${id}/record/start`),
  stopRecord: (id: string) => apiClient.post(`/streams/${id}/record/stop`),
  listRecordSessions: (id: string) => apiClient.get(`/streams/${id}/record/sessions`),
  getRecordings: (id: string) => apiClient.get(`/streams/${id}/record/files`),
  isRecording: (id: string) => apiClient.get<{ recording: boolean }>(`/streams/${id}/record/status`),

  getPlaybackURL: (id: string, start: string) => 
    apiClient.get<PlaybackURLs>(`/streams/${id}/playback`, { params: { start } }),
  listPlaybackSegments: (id: string) => 
    apiClient.get<PlaybackSegment[]>(`/streams/${id}/playback/segments`)
}
