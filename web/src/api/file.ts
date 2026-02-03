import apiClient from './client'

export interface File {
  id: string
  name: string
  original_name: string
  path: string
  url: string
  size: number
  mime_type: string
  type: 'image' | 'video' | 'audio' | 'document' | 'archive' | 'other'
  extension: string
  status: 'uploading' | 'completed' | 'failed' | 'deleted'
  hash: string
  uploader_id?: string
  metadata?: Record<string, any>
  created_at: string
  updated_at: string
}

export interface FileListQuery {
  type?: 'image' | 'video' | 'audio' | 'document' | 'archive' | 'other'
  status?: 'uploading' | 'completed' | 'failed' | 'deleted'
  uploader_id?: string
  search?: string
  from?: number
  to?: number
  page?: number
  page_size?: number
}

export interface FileListResponse {
  items: File[]
  total: number
}

export interface FileUpdateReq {
  name?: string
  status?: 'uploading' | 'completed' | 'failed' | 'deleted'
  metadata?: Record<string, any>
}

export const fileApi = {
  /**
   * 上传文件
   */
  upload(file: File, onProgress?: (progress: number) => void) {
    const formData = new FormData()
    formData.append('file', file)

    return apiClient.post<File>('/files', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          onProgress(progress)
        }
      }
    })
  },

  /**
   * 获取文件列表
   */
  list(params?: FileListQuery) {
    return apiClient.get<FileListResponse>('/files', { params })
  },

  /**
   * 获取文件详情
   */
  get(id: string) {
    return apiClient.get<File>(`/files/${id}`)
  },

  /**
   * 更新文件
   */
  update(id: string, data: FileUpdateReq) {
    return apiClient.put<File>(`/files/${id}`, data)
  },

  /**
   * 删除文件
   */
  delete(id: string) {
    return apiClient.delete(`/files/${id}`)
  },

  /**
   * 下载文件（返回文件 URL）
   */
  getDownloadUrl(id: string): string {
    return `${apiClient.defaults.baseURL}/files/${id}/download`
  }
}
