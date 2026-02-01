import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, type UserInfo, type MenuInfo, type LoginRequest } from '../api/auth'
import { getToken, setToken, getRefreshToken, setRefreshToken, clearTokens } from '../utils/auth'
import router from '../router'

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(getToken())
  const refreshToken = ref<string | null>(getRefreshToken())
  const userInfo = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const username = computed(() => userInfo.value?.username || '')
  const nickname = computed(() => userInfo.value?.nickname || userInfo.value?.username || '')
  const roles = computed(() => userInfo.value?.roles || [])
  const permissions = computed(() => userInfo.value?.permissions || [])
  const menus = computed(() => userInfo.value?.menus || [])

  function hasPermission(permission: string): boolean {
    if (permissions.value.includes('*')) {
      return true
    }
    return permissions.value.includes(permission)
  }

  function hasAnyPermission(permissionList: string[]): boolean {
    if (permissions.value.includes('*')) {
      return true
    }
    return permissionList.some(p => permissions.value.includes(p))
  }

  function hasRole(role: string): boolean {
    return roles.value.includes(role)
  }

  async function login(loginData: LoginRequest) {
    const response = await authApi.login(loginData)
    const data = response.data

    token.value = data.access_token
    refreshToken.value = data.refresh_token
    userInfo.value = data.user

    setToken(data.access_token)
    setRefreshToken(data.refresh_token)

    return data
  }

  async function getProfile() {
    if (!token.value) {
      throw new Error('No token')
    }

    const response = await authApi.getProfile()
    userInfo.value = response.data
    return response.data
  }

  async function refreshAccessToken() {
    if (!refreshToken.value) {
      throw new Error('No refresh token')
    }

    try {
      const response = await authApi.refreshToken({ refresh_token: refreshToken.value })
      const data = response.data

      token.value = data.access_token
      refreshToken.value = data.refresh_token

      setToken(data.access_token)
      setRefreshToken(data.refresh_token)

      return data
    } catch {
      logout()
      throw new Error('Refresh token expired')
    }
  }

  async function logout() {
    try {
      if (token.value) {
        await authApi.logout()
      }
    } catch {
      // ignore
    }

    token.value = null
    refreshToken.value = null
    userInfo.value = null
    clearTokens()

    router.push('/login')
  }

  function resetState() {
    token.value = null
    refreshToken.value = null
    userInfo.value = null
    clearTokens()
  }

  return {
    token,
    refreshToken,
    userInfo,
    isLoggedIn,
    username,
    nickname,
    roles,
    permissions,
    menus,
    hasPermission,
    hasAnyPermission,
    hasRole,
    login,
    getProfile,
    refreshAccessToken,
    logout,
    resetState
  }
})

export function flattenMenus(menus: MenuInfo[]): MenuInfo[] {
  const result: MenuInfo[] = []
  function traverse(items: MenuInfo[]) {
    for (const item of items) {
      result.push(item)
      if (item.children && item.children.length > 0) {
        traverse(item.children)
      }
    }
  }
  traverse(menus)
  return result
}
