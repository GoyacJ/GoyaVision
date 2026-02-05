import { reactive, computed, type ComputedRef } from 'vue'

/**
 * 分页配置选项
 */
export interface UsePaginationOptions {
  /** 初始页码 */
  initialPage?: number
  /** 初始每页条数 */
  initialPageSize?: number
  /** 每页条数选项 */
  pageSizeOptions?: number[]
  /** 总数 */
  total?: number
}

/**
 * 分页状态
 */
export interface PaginationState {
  /** 当前页码 */
  page: number
  /** 每页条数 */
  pageSize: number
  /** 总数 */
  total: number
}

/**
 * 分页返回值
 */
export interface UsePaginationReturn {
  /** 分页状态（响应式） */
  pagination: PaginationState
  /** 总页数 */
  totalPages: ComputedRef<number>
  /** 是否有上一页 */
  hasPrevPage: ComputedRef<boolean>
  /** 是否有下一页 */
  hasNextPage: ComputedRef<boolean>
  /** 当前页的起始索引（从 0 开始） */
  startIndex: ComputedRef<number>
  /** 当前页的结束索引（从 0 开始） */
  endIndex: ComputedRef<number>
  /** 跳转到指定页 */
  goToPage: (page: number) => void
  /** 上一页 */
  prevPage: () => void
  /** 下一页 */
  nextPage: () => void
  /** 更改每页条数 */
  changePageSize: (size: number) => void
  /** 更新总数 */
  setTotal: (total: number) => void
  /** 重置分页 */
  reset: () => void
}

/**
 * 使用分页
 * @param options 分页选项
 * @returns 分页状态和方法
 *
 * @example
 * ```ts
 * const { pagination, goToPage, changePageSize } = usePagination({
 *   initialPage: 1,
 *   initialPageSize: 20
 * })
 *
 * // 使用
 * goToPage(2)
 * changePageSize(50)
 * ```
 */
export function usePagination(
  options: UsePaginationOptions = {}
): UsePaginationReturn {
  const {
    initialPage = 1,
    initialPageSize = 20,
    total = 0
  } = options

  // 分页状态
  const pagination = reactive<PaginationState>({
    page: initialPage,
    pageSize: initialPageSize,
    total
  })

  // 总页数
  const totalPages = computed(() => {
    return Math.ceil(pagination.total / pagination.pageSize) || 1
  })

  // 是否有上一页
  const hasPrevPage = computed(() => {
    return pagination.page > 1
  })

  // 是否有下一页
  const hasNextPage = computed(() => {
    return pagination.page < totalPages.value
  })

  // 当前页的起始索引（从 0 开始）
  const startIndex = computed(() => {
    return (pagination.page - 1) * pagination.pageSize
  })

  // 当前页的结束索引（从 0 开始）
  const endIndex = computed(() => {
    return Math.min(
      startIndex.value + pagination.pageSize - 1,
      pagination.total - 1
    )
  })

  /**
   * 跳转到指定页
   */
  const goToPage = (page: number) => {
    if (page < 1 || page > totalPages.value) {
      console.warn(`Invalid page number: ${page}. Valid range: 1-${totalPages.value}`)
      return
    }
    pagination.page = page
  }

  /**
   * 上一页
   */
  const prevPage = () => {
    if (hasPrevPage.value) {
      pagination.page--
    }
  }

  /**
   * 下一页
   */
  const nextPage = () => {
    if (hasNextPage.value) {
      pagination.page++
    }
  }

  /**
   * 更改每页条数
   */
  const changePageSize = (size: number) => {
    if (size <= 0) {
      console.warn(`Invalid page size: ${size}`)
      return
    }

    pagination.pageSize = size

    // 重新计算当前页，确保不超出范围
    const newTotalPages = Math.ceil(pagination.total / size) || 1
    if (pagination.page > newTotalPages) {
      pagination.page = newTotalPages
    }
  }

  /**
   * 更新总数
   */
  const setTotal = (total: number) => {
    if (total < 0) {
      console.warn(`Invalid total: ${total}`)
      return
    }

    pagination.total = total

    // 重新计算当前页，确保不超出范围
    if (pagination.page > totalPages.value) {
      pagination.page = totalPages.value || 1
    }
  }

  /**
   * 重置分页
   */
  const reset = () => {
    pagination.page = initialPage
    pagination.pageSize = initialPageSize
    pagination.total = total
  }

  return {
    pagination,
    totalPages,
    hasPrevPage,
    hasNextPage,
    startIndex,
    endIndex,
    goToPage,
    prevPage,
    nextPage,
    changePageSize,
    setTotal,
    reset
  }
}
