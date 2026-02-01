import apiClient from './client'

export interface Menu {
  id: string
  parent_id?: string
  code: string
  name: string
  type: number
  path: string
  icon: string
  component: string
  permission: string
  sort: number
  visible: boolean
  status: number
  children?: Menu[]
  created_at: string
  updated_at: string
}

export interface CreateMenuRequest {
  parent_id?: string
  code: string
  name: string
  type: number
  path?: string
  icon?: string
  component?: string
  permission?: string
  sort?: number
  visible?: boolean
  status?: number
}

export interface UpdateMenuRequest {
  parent_id?: string
  name?: string
  type?: number
  path?: string
  icon?: string
  component?: string
  permission?: string
  sort?: number
  visible?: boolean
  status?: number
}

export const menuApi = {
  list(params?: { status?: number }) {
    return apiClient.get<Menu[]>('/menus', { params })
  },

  listTree(params?: { status?: number }) {
    return apiClient.get<Menu[]>('/menus/tree', { params })
  },

  get(id: string) {
    return apiClient.get<Menu>(`/menus/${id}`)
  },

  create(data: CreateMenuRequest) {
    return apiClient.post<Menu>('/menus', data)
  },

  update(id: string, data: UpdateMenuRequest) {
    return apiClient.put<Menu>(`/menus/${id}`, data)
  },

  delete(id: string) {
    return apiClient.delete(`/menus/${id}`)
  }
}
