import apiClient from './client'

export interface Tenant {
  id: string
  name: string
  code: string
  status: number
  created_at: string
  updated_at: string
}

export interface TenantCreateReq {
  name: string
  code: string
  status?: number
}

export interface TenantUpdateReq {
  name?: string
  status?: number
}

export const tenantApi = {
  list() {
    return apiClient.get<Tenant[]>('/tenants')
  },
  get(id: string) {
    return apiClient.get<Tenant>(`/tenants/${id}`)
  },
  create(data: TenantCreateReq) {
    return apiClient.post<Tenant>('/tenants', data)
  },
  update(id: string, data: TenantUpdateReq) {
    return apiClient.put<Tenant>(`/tenants/${id}`, data)
  },
  delete(id: string) {
    return apiClient.delete(`/tenants/${id}`)
  }
}
