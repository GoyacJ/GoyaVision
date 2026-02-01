import apiClient from './client'

export interface Permission {
  id: string
  code: string
  name: string
  method: string
  path: string
  description: string
}

export interface Role {
  id: string
  code: string
  name: string
  description: string
  status: number
  permissions: { id: string; code: string; name: string }[]
  menus: { id: string; code: string; name: string }[]
  created_at: string
  updated_at: string
}

export interface CreateRoleRequest {
  code: string
  name: string
  description?: string
  status?: number
  permission_ids?: string[]
  menu_ids?: string[]
}

export interface UpdateRoleRequest {
  name?: string
  description?: string
  status?: number
  permission_ids?: string[]
  menu_ids?: string[]
}

export const roleApi = {
  list(params?: { status?: number }) {
    return apiClient.get<Role[]>('/roles', { params })
  },

  get(id: string) {
    return apiClient.get<Role>(`/roles/${id}`)
  },

  create(data: CreateRoleRequest) {
    return apiClient.post<Role>('/roles', data)
  },

  update(id: string, data: UpdateRoleRequest) {
    return apiClient.put<Role>(`/roles/${id}`, data)
  },

  delete(id: string) {
    return apiClient.delete(`/roles/${id}`)
  }
}

export const permissionApi = {
  list() {
    return apiClient.get<Permission[]>('/permissions')
  }
}
