import apiClient from './client'

export interface LoginRequest {
  username: string
  password: string
}

export interface MenuInfo {
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
  children?: MenuInfo[]
}

export interface UserInfo {
  id: string
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  roles: string[]
  permissions: string[]
  menus: MenuInfo[]
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: UserInfo
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

export interface LoginOAuthRequest {
  provider: string
  code: string
  state?: string
}

export const authApi = {
  login(data: LoginRequest) {
    return apiClient.post<LoginResponse>('/auth/login', data)
  },

  loginOAuth(data: LoginOAuthRequest) {
    return apiClient.post<LoginResponse>('/auth/oauth/login', data)
  },

  refreshToken(data: RefreshTokenRequest) {
    return apiClient.post<LoginResponse>('/auth/refresh', data)
  },

  getProfile() {
    return apiClient.get<UserInfo>('/auth/profile')
  },

  changePassword(data: ChangePasswordRequest) {
    return apiClient.put('/auth/password', data)
  },

  logout() {
    return apiClient.post('/auth/logout')
  }
}
