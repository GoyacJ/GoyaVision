import apiClient from './client'

export interface User {
  id: string
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  status: number
  roles: { id: string; code: string; name: string }[]
  created_at: string
  updated_at: string
}

export interface UserListResponse {
  items: User[]
  total: number
}

export interface CreateUserRequest {
  username: string
  password: string
  nickname?: string
  email?: string
  phone?: string
  status?: number
  role_ids?: string[]
}

export interface UpdateUserRequest {
  nickname?: string
  email?: string
  phone?: string
  status?: number
  password?: string
  role_ids?: string[]
}

export const userApi = {
  list(params?: { status?: number; limit?: number; offset?: number }) {
    return apiClient.get<UserListResponse>('/users', { params })
  },

  get(id: string) {
    return apiClient.get<User>(`/users/${id}`)
  },

  create(data: CreateUserRequest) {
    return apiClient.post<User>('/users', data)
  },

  update(id: string, data: UpdateUserRequest) {
    return apiClient.put<User>(`/users/${id}`, data)
  },

  delete(id: string) {
    return apiClient.delete(`/users/${id}`)
  },

  resetPassword(id: string, newPassword: string) {
    return apiClient.post(`/users/${id}/reset-password`, { new_password: newPassword })
  }
}
