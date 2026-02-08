import router from './index'
import { getToken } from '../utils/auth'
import { ElMessage } from 'element-plus'

const whiteList = ['/login']

router.beforeEach(async (to, _from, next) => {
  const { useAppStore } = await import('../store/app')
  const appStore = useAppStore()
  
  // Ensure config is loaded
  if (!appStore.configLoaded) {
    await appStore.fetchPublicConfig()
    // Restart navigation to ensure dynamic routes are picked up
    next({ ...to, replace: true })
    return
  }

  const token = getToken()

  if (token) {
    if (to.path === '/login') {
      next({ path: '/' })
    } else {
      const { useUserStore } = await import('../store/user')
      const userStore = useUserStore()

      if (userStore.userInfo && userStore.routesLoaded) {
        if (to.path === '/') {
          next({ path: appStore.defaultHome, replace: true })
        } else {
          next()
        }
      } else {
        try {
          await userStore.getProfile()

          if (to.path === '/') {
            next({ path: appStore.defaultHome, replace: true })
          } else {
            next({ ...to, replace: true })
          }
        } catch (error) {
          userStore.resetState()
          ElMessage.error('登录已过期，请重新登录')
          next(`/login?redirect=${to.path}`)
        }
      }
    }
  } else {
    // Not logged in
    if (whiteList.includes(to.path)) {
      next()
      return
    }
    
    // Redirect root to default home
    if (to.path === '/') {
       next({ path: appStore.defaultHome, replace: true })
       return
    }
    
    // Check if path is in public menus
    const { flattenMenus } = await import('../store/user')
    const flatPublic = flattenMenus(appStore.publicMenus)
    // Simple check: exact match
    const isPublic = flatPublic.some(m => m.path === to.path)
    
    if (isPublic) {
      next()
    } else {
      next(`/login?redirect=${to.path}`)
    }
  }
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title} - GoyaVision` : 'GoyaVision'
})
