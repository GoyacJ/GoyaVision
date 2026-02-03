import apiClient from './client'

export interface MediaAsset {
  id: string
  type: 'video' | 'image' | 'audio' | 'stream' // ✅ 添加 stream 类型
  source_type: 'upload' | 'stream_capture' | 'operator_output' | 'live' | 'vod' | 'generated'
  source_id?: string
  parent_id?: string
  name: string
  path: string
  duration?: number
  size: number
  format?: string
  metadata?: Record<string, any>
  status: 'pending' | 'ready' | 'processing' | 'error'
  tags?: string[]
  created_at: string
  updated_at: string
}

export interface AssetListQuery {
  type?: 'video' | 'image' | 'audio' | 'stream' // ✅ 添加 stream 类型
  source_type?: 'upload' | 'stream_capture' | 'operator_output' | 'live' | 'vod' | 'generated'
  source_id?: string
  parent_id?: string
  status?: 'pending' | 'ready' | 'processing' | 'error'
  tags?: string
  name?: string
  page?: number
  page_size?: number
}

export interface AssetCreateReq {
  type: 'video' | 'image' | 'audio' | 'stream' // ✅ 添加 stream 类型
  source_type: 'upload' | 'stream_capture' | 'operator_output' | 'live' | 'vod' | 'generated'
  source_id?: string
  parent_id?: string
  name: string
  path: string
  duration?: number
  size: number
  format?: string
  metadata?: Record<string, any>
  tags?: string[]
}

export interface AssetUpdateReq {
  name?: string
  metadata?: Record<string, any>
  status?: 'pending' | 'ready' | 'processing' | 'error'
  tags?: string[]
}

export interface AssetListResponse {
  items: MediaAsset[]
  total: number
  page: number
  page_size: number
}

export const assetApi = {
  list(params?: AssetListQuery) {
    return apiClient.get<AssetListResponse>('/assets', { params })
  },

  get(id: string) {
    return apiClient.get<MediaAsset>(`/assets/${id}`)
  },

  create(data: AssetCreateReq) {
    return apiClient.post<MediaAsset>('/assets', data)
  },

  update(id: string, data: AssetUpdateReq) {
    return apiClient.put<MediaAsset>(`/assets/${id}`, data)
  },

  delete(id: string) {
    return apiClient.delete(`/assets/${id}`)
  },

  listChildren(id: string, params?: AssetListQuery) {
    return apiClient.get<AssetListResponse>(`/assets/${id}/children`, { params })
  },

  getTags() {
    return apiClient.get<{ tags: string[] }>('/assets/tags')
  },

  // 文件上传
  upload(file: File, type: string, name?: string, tags?: string[]) {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('type', type)
    if (name) {
      formData.append('name', name)
    }
    if (tags && tags.length > 0) {
      formData.append('tags', JSON.stringify(tags))
    }

    return apiClient.post<MediaAsset>('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
}
