<template>
  <div :class="tableContainerClasses">
    <!-- 表格 -->
    <el-table
      ref="tableRef"
      :data="data"
      :class="tableClasses"
      :size="size"
      :border="border"
      :stripe="stripe"
      :show-header="showHeader"
      :highlight-current-row="highlightCurrentRow"
      :height="height"
      :max-height="maxHeight"
      :row-key="rowKey"
      :empty-text="emptyText"
      v-loading="loading"
      @selection-change="handleSelectionChange"
      @row-click="handleRowClick"
      @row-dblclick="handleRowDblclick"
      @cell-click="handleCellClick"
      @sort-change="handleSortChange"
    >
      <!-- 选择列 -->
      <el-table-column
        v-if="selectable"
        type="selection"
        width="55"
        fixed="left"
      />
      
      <!-- 数据列 -->
      <el-table-column
        v-for="column in columns"
        :key="column.prop"
        :prop="column.prop"
        :label="column.label"
        :width="column.width"
        :min-width="column.minWidth"
        :fixed="column.fixed"
        :sortable="column.sortable"
        :align="column.align || 'left'"
        :formatter="column.formatter"
        :show-overflow-tooltip="column.showOverflowTooltip !== false"
      >
        <!-- 自定义列插槽 -->
        <template #default="scope">
          <slot
            :name="column.prop"
            :row="scope.row"
            :column="column"
            :$index="scope.$index"
          >
            {{ scope.row[column.prop] }}
          </slot>
        </template>
      </el-table-column>
      
      <!-- 操作列插槽 -->
      <slot name="actions" />
    </el-table>
    
    <!-- 分页 -->
    <div v-if="pagination && paginationConfig" class="gv-table__pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="currentPageSize"
        :page-sizes="paginationConfig.pageSizes || [10, 20, 50, 100]"
        :layout="paginationConfig.layout || 'total, sizes, prev, pager, next, jumper'"
        :total="paginationConfig.total"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { cn } from '@/utils/cn'
import type { TableProps, TableEmits, TableColumn } from './types'
import type { ElTable } from 'element-plus'

const props = withDefaults(defineProps<TableProps>(), {
  size: 'default',
  border: false,
  stripe: true,
  showHeader: true,
  highlightCurrentRow: false,
  loading: false,
  emptyText: '暂无数据',
  rowKey: 'id',
  pagination: false,
  selectable: false
})

const emit = defineEmits<TableEmits>()

// 表格引用
const tableRef = ref<InstanceType<typeof ElTable>>()

// 分页状态
const currentPage = ref(props.paginationConfig?.currentPage || 1)
const currentPageSize = ref(props.paginationConfig?.pageSize || 10)

// 监听分页配置变化
watch(
  () => props.paginationConfig,
  (newConfig) => {
    if (newConfig) {
      currentPage.value = newConfig.currentPage
      currentPageSize.value = newConfig.pageSize
    }
  },
  { deep: true }
)

// 表格容器类名
const tableContainerClasses = computed(() => {
  return cn('gv-table', 'w-full')
})

// 表格类名
const tableClasses = computed(() => {
  return cn('gv-table__table')
})

// 选择变化
const handleSelectionChange = (selection: any[]) => {
  emit('selection-change', selection)
}

// 行点击
const handleRowClick = (row: any, column: any, event: Event) => {
  emit('row-click', row, column, event)
}

// 行双击
const handleRowDblclick = (row: any, column: any, event: Event) => {
  emit('row-dblclick', row, column, event)
}

// 单元格点击
const handleCellClick = (row: any, column: any, cell: any, event: Event) => {
  emit('cell-click', row, column, cell, event)
}

// 排序变化
const handleSortChange = (data: { column: any; prop: string; order: string }) => {
  emit('sort-change', data as any)
}

// 当前页变化
const handleCurrentChange = (page: number) => {
  emit('current-change', page)
}

// 每页条数变化
const handleSizeChange = (size: number) => {
  emit('size-change', size)
}

// 暴露方法
const clearSelection = () => {
  tableRef.value?.clearSelection()
}

const toggleRowSelection = (row: any, selected?: boolean) => {
  tableRef.value?.toggleRowSelection(row, selected)
}

const setCurrentRow = (row: any) => {
  tableRef.value?.setCurrentRow(row)
}

defineExpose({
  clearSelection,
  toggleRowSelection,
  setCurrentRow
})
</script>

<style>
/* 自定义 Element Plus Table 样式 */
.gv-table__table {
  @apply overflow-hidden;
}

.gv-table__table .el-table__header-wrapper {
  @apply bg-neutral-50;
}

.gv-table__table .el-table__header th {
  @apply bg-neutral-50 text-text-primary font-semibold;
  border-bottom: 2px solid rgb(226 232 240);
}

.gv-table__table .el-table__body tr {
  @apply transition-colors duration-150;
}

.gv-table__table .el-table__body tr:hover {
  @apply bg-primary-50;
}

.gv-table__table .el-table__body td {
  @apply text-text-primary;
}

.gv-table__table.el-table--border {
  @apply border border-neutral-200;
}

.gv-table__table.el-table--border th,
.gv-table__table.el-table--border td {
  @apply border-neutral-200;
}

.gv-table__table .el-table__empty-text {
  @apply text-text-tertiary;
}

/* 分页样式 */
.gv-table__pagination {
  @apply mt-4 flex justify-end;
}

.gv-table__pagination .el-pagination {
  @apply flex items-center gap-2;
}

.gv-table__pagination .el-pagination button,
.gv-table__pagination .el-pagination .el-pager li {
  @apply min-w-8 h-8 rounded-lg;
  @apply transition-colors duration-150;
}

.gv-table__pagination .el-pagination button:hover,
.gv-table__pagination .el-pagination .el-pager li:hover {
  @apply bg-primary-50 text-primary-600;
}

.gv-table__pagination .el-pagination .el-pager li.is-active {
  @apply bg-primary-600 text-white;
}

.gv-table__pagination .el-pagination .el-select {
  @apply rounded-lg;
}

/* 深色模式 */
.dark .gv-table__table .el-table__header th {
  @apply bg-neutral-800 text-text-inverse;
  border-bottom-color: rgb(64 64 64);
}

.dark .gv-table__table .el-table__body tr {
  @apply bg-surface-dark;
}

.dark .gv-table__table .el-table__body tr:hover {
  @apply bg-primary-950/30;
}

.dark .gv-table__table .el-table__body td {
  @apply text-text-inverse border-neutral-700;
}

.dark .gv-table__table.el-table--border {
  @apply border-neutral-700;
}

.dark .gv-table__table.el-table--border th,
.dark .gv-table__table.el-table--border td {
  @apply border-neutral-700;
}

.dark .gv-table__table .el-table__empty-text {
  @apply text-neutral-500;
}

.dark .gv-table__pagination .el-pagination button,
.dark .gv-table__pagination .el-pagination .el-pager li {
  @apply bg-neutral-800 text-text-inverse;
}

.dark .gv-table__pagination .el-pagination button:hover,
.dark .gv-table__pagination .el-pagination .el-pager li:hover {
  @apply bg-primary-950/50 text-primary-300;
}

.dark .gv-table__pagination .el-pagination .el-pager li.is-active {
  @apply bg-primary-700 text-white;
}
</style>
