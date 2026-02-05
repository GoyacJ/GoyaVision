import { ref, type Ref } from 'vue'

/**
 * 异步数据加载 Composable
 * 统一处理 Loading、Error、Data 三种状态
 */

export interface UseAsyncDataOptions<T> {
  /** 是否立即执行 */
  immediate?: boolean
  /** 初始数据 */
  initialData?: T
  /** 成功回调 */
  onSuccess?: (data: T) => void
  /** 错误回调 */
  onError?: (error: Error) => void
  /** 重置错误状态的延迟时间（ms），0 表示不自动重置 */
  resetErrorDelay?: number
}

export interface UseAsyncDataReturn<T> {
  /** 数据 */
  data: Ref<T | null>
  /** 加载状态 */
  isLoading: Ref<boolean>
  /** 错误 */
  error: Ref<Error | null>
  /** 执行函数 */
  execute: (...args: any[]) => Promise<T | null>
  /** 重置状态 */
  reset: () => void
  /** 刷新（重新执行上次的请求） */
  refresh: () => Promise<T | null>
}

/**
 * 使用异步数据
 * @param fn 异步函数
 * @param options 选项
 * @returns 异步数据状态和方法
 *
 * @example
 * ```ts
 * const { data, isLoading, error, execute } = useAsyncData(
 *   () => assetApi.list({ page: 1 }),
 *   { immediate: true }
 * )
 * ```
 */
export function useAsyncData<T>(
  fn: (...args: any[]) => Promise<T>,
  options: UseAsyncDataOptions<T> = {}
): UseAsyncDataReturn<T> {
  const {
    immediate = false,
    initialData = null,
    onSuccess,
    onError,
    resetErrorDelay = 0
  } = options

  const data = ref<T | null>(initialData) as Ref<T | null>
  const isLoading = ref(false)
  const error = ref<Error | null>(null)

  // 存储上次执行的参数，用于 refresh
  let lastArgs: any[] = []

  // 错误重置定时器
  let errorResetTimer: number | null = null

  /**
   * 执行异步函数
   */
  const execute = async (...args: any[]): Promise<T | null> => {
    // 保存参数供 refresh 使用
    lastArgs = args

    // 清除之前的错误重置定时器
    if (errorResetTimer !== null) {
      clearTimeout(errorResetTimer)
      errorResetTimer = null
    }

    isLoading.value = true
    error.value = null

    try {
      const result = await fn(...args)
      data.value = result

      // 调用成功回调
      if (onSuccess) {
        onSuccess(result)
      }

      return result
    } catch (err) {
      const errorObj = err instanceof Error ? err : new Error(String(err))
      error.value = errorObj

      // 调用错误回调
      if (onError) {
        onError(errorObj)
      }

      // 设置错误自动重置
      if (resetErrorDelay > 0) {
        errorResetTimer = window.setTimeout(() => {
          error.value = null
        }, resetErrorDelay)
      }

      return null
    } finally {
      isLoading.value = false
    }
  }

  /**
   * 刷新（使用上次的参数重新执行）
   */
  const refresh = (): Promise<T | null> => {
    return execute(...lastArgs)
  }

  /**
   * 重置所有状态
   */
  const reset = () => {
    data.value = initialData
    isLoading.value = false
    error.value = null
    lastArgs = []

    if (errorResetTimer !== null) {
      clearTimeout(errorResetTimer)
      errorResetTimer = null
    }
  }

  // 立即执行
  if (immediate) {
    execute()
  }

  return {
    data,
    isLoading,
    error,
    execute,
    refresh,
    reset
  }
}
