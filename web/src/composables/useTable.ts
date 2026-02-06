import { watch, computed, type Ref, type ComputedRef } from 'vue'
import { useAsyncData, type UseAsyncDataReturn } from './useAsyncData'
import { usePagination, type UsePaginationReturn } from './usePagination'

/**
 * 表格数据接口
 */
export interface TableData<T> {
  /** 数据列表 */
  items: T[]
  /** 总数 */
  total: number
}

/**
 * 表格查询参数接口
 */
export interface TableQueryParams {
  /** 页码 */
  page: number
  /** 每页条数 */
  page_size: number
  /** 其他查询参数 */
  [key: string]: any
}

/**
 * 表格配置选项
 */
export interface UseTableOptions<T> {
  /** 是否立即加载 */
  immediate?: boolean
  /** 初始页码 */
  initialPage?: number
  /** 初始每页条数 */
  initialPageSize?: number
  /** 额外的查询参数（响应式） */
  extraParams?: Ref<Record<string, any>> | Record<string, any>
  /** 成功回调 */
  onSuccess?: (data: TableData<T>) => void
  /** 错误回调 */
  onError?: (error: Error) => void
}

/**
 * 表格返回值
 */
export interface UseTableReturn<T> extends Omit<UseAsyncDataReturn<TableData<T>>, 'refresh' | 'reset'>, UsePaginationReturn {
  /** 表格数据列表 */
  items: ComputedRef<T[]>
  /** 加载表格数据 */
  loadTable: () => Promise<void>
  /** 刷新当前页 */
  refreshTable: () => Promise<void>
  /** 重置并重新加载 */
  resetTable: () => Promise<void>
}

/**
 * 使用表格
 * @param fetchFn 获取表格数据的函数
 * @param options 表格选项
 * @returns 表格状态和方法
 *
 * @example
 * ```ts
 * const searchParams = ref({ keyword: '' })
 *
 * const {
 *   items,
 *   isLoading,
 *   error,
 *   pagination,
 *   goToPage,
 *   changePageSize,
 *   refreshTable
 * } = useTable(
 *   (params) => assetApi.list(params),
 *   {
 *     immediate: true,
 *     initialPageSize: 20,
 *     extraParams: searchParams
 *   }
 * )
 *
 * // 搜索时自动重新加载
 * watch(searchParams, () => {
 *   pagination.page = 1
 *   loadTable()
 * })
 * ```
 */
export function useTable<T>(
  fetchFn: (params: TableQueryParams) => Promise<TableData<T>>,
  options: UseTableOptions<T> = {}
): UseTableReturn<T> {
  const {
    immediate = false,
    initialPage = 1,
    initialPageSize = 20,
    extraParams,
    onSuccess,
    onError
  } = options

  // 初始化分页
  const paginationHook = usePagination({
    initialPage,
    initialPageSize,
    total: 0
  })

  const { pagination, setTotal } = paginationHook

  // 初始化异步数据
  const asyncDataHook = useAsyncData<TableData<T>>(
    async () => {
      // 构建查询参数
      const params: TableQueryParams = {
        page: pagination.page,
        page_size: pagination.pageSize
      }

      // 合并额外参数
      if (extraParams) {
        const extra = 'value' in extraParams ? extraParams.value : extraParams
        Object.assign(params, extra)
      }

      // 调用 API
      const result = await fetchFn(params)

      // 更新总数
      setTotal(result.total)

      // 调用成功回调
      if (onSuccess) {
        onSuccess(result)
      }

      return result
    },
    {
      immediate,
      onError
    }
  )

  const { data, isLoading, error, execute, reset: resetAsyncData } = asyncDataHook

  // 表格数据列表
  const items = computed(() => {
    return data.value?.items || []
  })

  /**
   * 加载表格数据
   */
  const loadTable = async (): Promise<void> => {
    await execute()
  }

  /**
   * 刷新当前页
   */
  const refreshTable = async (): Promise<void> => {
    await execute()
  }

  /**
   * 重置并重新加载
   */
  const resetTable = async (): Promise<void> => {
    paginationHook.reset()
    resetAsyncData()
    await execute()
  }

  // 监听分页变化，自动重新加载
  watch(
    () => [pagination.page, pagination.pageSize],
    () => {
      loadTable()
    }
  )

  // 监听额外参数变化，自动重置到第一页并重新加载
  if (extraParams && 'value' in extraParams) {
    watch(
      () => extraParams.value,
      () => {
        if (pagination.page === 1) {
          loadTable()
        } else {
          pagination.page = 1
        }
      },
      { deep: true }
    )
  }

  return {
    // 从 useAsyncData 继承
    data,
    isLoading,
    error,
    execute,

    // 从 usePagination 继承
    pagination,
    totalPages: paginationHook.totalPages,
    hasPrevPage: paginationHook.hasPrevPage,
    hasNextPage: paginationHook.hasNextPage,
    startIndex: paginationHook.startIndex,
    endIndex: paginationHook.endIndex,
    goToPage: paginationHook.goToPage,
    prevPage: paginationHook.prevPage,
    nextPage: paginationHook.nextPage,
    changePageSize: paginationHook.changePageSize,
    setTotal: paginationHook.setTotal,
    reset: paginationHook.reset,

    // 表格特定方法
    items,
    loadTable,
    refreshTable,
    resetTable
  }
}
