import { defineStore } from 'pinia'
import { ref } from 'vue'
import { systemApi } from '../api/system'
import type { MenuInfo } from '../api/auth'
import router from '../router'
import { buildRoutesFromMenus, hasRouteComponent } from '../utils/dynamic-routes'

export const useAppStore = defineStore('app', () => {
  const defaultHome = ref<string>('/assets')
  const publicMenus = ref<MenuInfo[]>([])
  const configLoaded = ref(false)
  const routesRegistered = ref(false)

  async function fetchPublicConfig() {
    try {
      const res = await systemApi.getPublicConfig()
      defaultHome.value = res.data.home_path || '/assets'
      publicMenus.value = res.data.public_menus || []
      configLoaded.value = true
      
      registerPublicRoutes()
    } catch (error) {
      console.error('Failed to load public config', error)
      configLoaded.value = true
    }
  }

  function registerPublicRoutes() {
    if (routesRegistered.value || publicMenus.value.length === 0) return

    const rootRouteName = 'Root'
    const menuRoutes = buildRoutesFromMenus(publicMenus.value)
    
    menuRoutes.forEach((route) => {
      if (!hasRouteComponent(route) && (!route.children || route.children.length === 0)) {
        return
      }
      router.addRoute(rootRouteName, route)
    })
    
    routesRegistered.value = true
    console.log('[App] Registered public routes')
  }

  return {
    defaultHome,
    publicMenus,
    configLoaded,
    fetchPublicConfig
  }
})
