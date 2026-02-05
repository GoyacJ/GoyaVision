import axios, { type AxiosError, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import { getToken, getRefreshToken, setToken, setRefreshToken, clearTokens } from '../utils/auth'
import router from '../router'

/**
 * API 响应数据格式
 */
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

/**
 * API 错误响应格式
 */
export interface ApiError {
  code: number
  message: string
  details?: any
}

/**
 * 创建 Axios 实例
 */
const apiClient = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

const refreshClient = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

let isRefreshing = false
const pendingRequests: Array<(token: string | null) => void> = []

function subscribeTokenRefresh(callback: (token: string | null) => void) {
  pendingRequests.push(callback)
}

function onTokenRefreshed(token: string | null) {
  pendingRequests.splice(0).forEach((callback) => callback(token))
}

/**
 * 请求拦截器
 */
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 添加 JWT Token
    const token = getToken()
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }

    return config
  },
  (error: AxiosError) => {
    console.error('[Request Error]', error)
    return Promise.reject(error)
  }
)

/**
 * 响应拦截器
 */
apiClient.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    // 直接返回响应（不提取 data，让调用者决定如何处理）
    return response
  },
  async (error: AxiosError<ApiError>) => {
    const originalRequest = error.config as (InternalAxiosRequestConfig & { _retry?: boolean }) | undefined

    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      originalRequest._retry = true
      const refreshToken = getRefreshToken()

      if (!refreshToken) {
        handleApiError(error)
        return Promise.reject(error)
      }

      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          subscribeTokenRefresh((newToken) => {
            if (!newToken) {
              reject(error)
              return
            }
            if (originalRequest.headers) {
              originalRequest.headers.Authorization = `Bearer ${newToken}`
            }
            resolve(apiClient(originalRequest))
          })
        })
      }

      isRefreshing = true
      try {
        const response = await refreshClient.post<{
          access_token: string
          refresh_token: string
        }>('/auth/refresh', { refresh_token: refreshToken })
        const data = response.data
        setToken(data.access_token)
        setRefreshToken(data.refresh_token)
        onTokenRefreshed(data.access_token)

        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${data.access_token}`
        }
        return apiClient(originalRequest)
      } catch (refreshError) {
        onTokenRefreshed(null)
        handleApiError(error)
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    // 处理各种错误情况
    handleApiError(error)
    return Promise.reject(error)
  }
)

/**
 * 统一错误处理
 */
function handleApiError(error: AxiosError<ApiError>): void {
  if (!error.response) {
    // 网络错误（无响应）
    if (error.code === 'ECONNABORTED') {
      console.error('[Network Timeout]', error.message)
    } else if (error.code === 'ERR_NETWORK') {
      console.error('[Network Error]', error.message)
    } else {
      console.error('[Unknown Error]', error.message)
    }
    return
  }

  const { status, data } = error.response

  switch (status) {
    case 400:
      // 请求参数错误
      console.error('[Bad Request]', data?.message || '请求参数错误')
      break

    case 401:
      // 未授权 - 清除 Token 并跳转登录
      console.warn('[Unauthorized]', '登录已过期，请重新登录')
      clearTokens()

      // 避免重复跳转
      if (router.currentRoute.value.path !== '/login') {
        router.push('/login')
      }
      break

    case 403:
      // 禁止访问
      console.error('[Forbidden]', data?.message || '没有权限访问此资源')
      ElMessage.error(data?.message || '没有权限访问此资源')
      break

    case 404:
      // 资源不存在
      console.error('[Not Found]', data?.message || '请求的资源不存在')
      break

    case 409:
      // 冲突（如资源已存在）
      console.error('[Conflict]', data?.message || '资源冲突')
      break

    case 422:
      // 验证失败
      console.error('[Validation Error]', data?.message || '数据验证失败')
      break

    case 429:
      // 请求过于频繁
      console.warn('[Too Many Requests]', '请求过于频繁，请稍后再试')
      ElMessage.warning('请求过于频繁，请稍后再试')
      break

    case 500:
      // 服务器错误
      console.error('[Server Error]', data?.message || '服务器内部错误')
      ElMessage.error(data?.message || '服务器内部错误')
      break

    case 502:
    case 503:
    case 504:
      // 网关错误或服务不可用
      console.error('[Service Unavailable]', '服务暂时不可用，请稍后再试')
      ElMessage.error('服务暂时不可用，请稍后再试')
      break

    default:
      console.error(`[HTTP ${status}]`, data?.message || '未知错误')
      break
  }
}

export default apiClient
