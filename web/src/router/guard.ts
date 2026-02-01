import router from './index'
import { getToken } from '../utils/auth'
import { ElMessage } from 'element-plus'

const whiteList = ['/login']

router.beforeEach(async (to, _from, next) => {
  const token = getToken()

  if (token) {
    if (to.path === '/login') {
      next({ path: '/' })
    } else {
      const { useUserStore } = await import('../store/user')
      const userStore = useUserStore()
      
      if (userStore.userInfo) {
        next()
      } else {
        try {
          await userStore.getProfile()
          next({ ...to, replace: true })
        } catch (error) {
          userStore.resetState()
          ElMessage.error('登录已过期，请重新登录')
          next(`/login?redirect=${to.path}`)
        }
      }
    }
  } else {
    if (whiteList.includes(to.path)) {
      next()
    } else {
      next(`/login?redirect=${to.path}`)
    }
  }
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title} - GoyaVision` : 'GoyaVision'
})
