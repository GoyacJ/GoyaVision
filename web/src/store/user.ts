import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, type UserInfo, type MenuInfo, type LoginRequest, type LoginOAuthRequest } from '../api/auth'
import { getToken, setToken, getRefreshToken, setRefreshToken, clearTokens } from '../utils/auth'
import { buildRoutesFromMenus, hasRouteComponent } from '../utils/dynamic-routes'

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(getToken())
  const refreshToken = ref<string | null>(getRefreshToken())
  const userInfo = ref<UserInfo | null>(null)
  const routesLoaded = ref(false)

  const isLoggedIn = computed(() => !!token.value)
  const username = computed(() => userInfo.value?.username || '')
  const nickname = computed(() => userInfo.value?.nickname || userInfo.value?.username || '')
  const avatar = computed(() => userInfo.value?.avatar || '')
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

    routesLoaded.value = false

    if (data.user && data.user.menus) {
      registerDynamicRoutes()
    }

    return data
  }

  async function loginOAuth(loginData: LoginOAuthRequest) {
    const response = await authApi.loginOAuth(loginData)
    const data = response.data

    token.value = data.access_token
    refreshToken.value = data.refresh_token
    userInfo.value = data.user

    setToken(data.access_token)
    setRefreshToken(data.refresh_token)

    routesLoaded.value = false

    if (data.user && data.user.menus) {
      registerDynamicRoutes()
    }

    return data
  }

  async function getProfile() {
    if (!token.value) {
      throw new Error('No token')
    }

    const response = await authApi.getProfile()
    userInfo.value = response.data
    if (!routesLoaded.value) {
      registerDynamicRoutes()
    }
    return response.data
  }

  async function registerDynamicRoutes() {
    const { default: router } = await import('../router')
    const rootRouteName = 'Root'
    const menus = userInfo.value?.menus || []
    
    if (menus.length === 0) {
      console.warn('[Router] No menus found, skipping route registration')
      routesLoaded.value = true
      return
    }
    
    const menuRoutes = buildRoutesFromMenus(menus)
    let registeredCount = 0
    
    menuRoutes.forEach((route) => {
      if (!hasRouteComponent(route) && (!route.children || route.children.length === 0)) {
        return
      }
      router.addRoute(rootRouteName, route)
      registeredCount++
    })
    
    console.log(`[Router] Registered ${registeredCount} dynamic routes from ${menus.length} menus`)
    routesLoaded.value = true
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
    routesLoaded.value = false
    clearTokens()

    const { default: router } = await import('../router')
    router.push('/login')
  }

  function resetState() {
    token.value = null
    refreshToken.value = null
    userInfo.value = null
    routesLoaded.value = false
    clearTokens()
  }

  return {
    token,
    refreshToken,
    userInfo,
    isLoggedIn,
    username,
    nickname,
    avatar,
    roles,
    permissions,
    menus,
    routesLoaded,
    hasPermission,
    hasAnyPermission,
    hasRole,
    login,
    loginOAuth,
    getProfile,
    registerDynamicRoutes,
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
