import apiClient from './client'
import type { MenuInfo } from './auth'

export interface GetPublicConfigResponse {
  home_path: string
  public_menus: MenuInfo[]
}

export const systemApi = {
  getPublicConfig() {
    return apiClient.get<GetPublicConfigResponse>('/public/config')
  },

  updateConfig(data: Record<string, any>) {
    return apiClient.put('/system/config', data)
  }
}
