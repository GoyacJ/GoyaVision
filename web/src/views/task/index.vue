<template>
  <GvContainer max-width="full">
    <!-- 页面头部 -->
    <PageHeader
      title="任务管理"
      description="查看和管理所有工作流任务的执行状态"
    >
      <template #actions>
        <GvButton @click="handleRefresh">
          <template #icon>
            <el-icon><Refresh /></el-icon>
          </template>
          刷新
        </GvButton>
      </template>
      
      <template #extra>
        <GvSpace wrap>
          <StatusBadge status="pending" :text="`${stats.pending} 个待执行`" :animated="false" />
          <StatusBadge status="running" :text="`${stats.running} 个执行中`" />
          <StatusBadge status="success" :text="`${stats.success} 个已成功`" :animated="false" />
          <StatusBadge status="failed" :text="`${stats.failed} 个已失败`" :animated="false" />
        </GvSpace>
      </template>
    </PageHeader>

    <!-- 统计卡片 -->
    <GvGrid :cols="6" gap="lg" class="mb-6">
      <GvCard class="text-center">
        <div class="text-text-tertiary text-sm mb-2">总任务数</div>
        <div class="text-3xl font-bold text-text-primary">{{ stats.total }}</div>
      </GvCard>
      <GvCard class="text-center border-l-4 border-neutral-400">
        <div class="text-text-tertiary text-sm mb-2">待执行</div>
        <div class="text-3xl font-bold text-neutral-600">{{ stats.pending }}</div>
      </GvCard>
      <GvCard class="text-center border-l-4 border-primary-500">
        <div class="text-text-tertiary text-sm mb-2">执行中</div>
        <div class="text-3xl font-bold text-primary-600">{{ stats.running }}</div>
      </GvCard>
      <GvCard class="text-center border-l-4 border-success-500">
        <div class="text-text-tertiary text-sm mb-2">已成功</div>
        <div class="text-3xl font-bold text-success-600">{{ stats.success }}</div>
      </GvCard>
      <GvCard class="text-center border-l-4 border-error-500">
        <div class="text-text-tertiary text-sm mb-2">已失败</div>
        <div class="text-3xl font-bold text-error-600">{{ stats.failed }}</div>
      </GvCard>
      <GvCard class="text-center border-l-4 border-warning-500">
        <div class="text-text-tertiary text-sm mb-2">已取消</div>
        <div class="text-3xl font-bold text-warning-600">{{ stats.cancelled }}</div>
      </GvCard>
    </GvGrid>

    <!-- 筛选栏 -->
    <FilterBar
      v-model="filters"
      :fields="filterFields"
      :columns="1"
      :loading="loading"
      @filter="() => { pagination.page = 1 }"
      @reset="handleResetFilter"
    />

    <ErrorState
      v-if="error && !loading"
      :error="error"
      title="加载失败"
      @retry="handleRefresh"
    />

    <EmptyState
      v-else-if="!loading && tasks.length === 0"
      icon="📋"
      title="还没有任务"
      description="触发工作流后将在这里显示执行任务"
    />

    <!-- 数据表格 -->
    <GvTable
      v-else
      :data="tasks"
      :columns="columns"
      :loading="loading"
      border
      stripe
      pagination
      :pagination-config="paginationConfig"
      @current-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <template #workflow="{ row }">
        {{ row.workflow?.name || '-' }}
      </template>
      
      <template #status="{ row }">
        <StatusBadge :status="mapStatus(row.status)" />
      </template>
      
      <template #progress="{ row }">
        <el-progress
          :percentage="row.progress"
          :status="getProgressStatus(row.status)"
        />
      </template>
      
      <template #started_at="{ row }">
        {{ row.started_at ? formatDate(row.started_at) : '-' }}
      </template>
      
      <template #completed_at="{ row }">
        {{ row.completed_at ? formatDate(row.completed_at) : '-' }}
      </template>
      
      <template #duration="{ row }">
        {{ calculateDuration(row) }}
      </template>
      
      <template #actions="{ row }">
        <GvSpace size="xs">
          <GvButton size="small" variant="tonal" @click="handleView(row)">
            查看
          </GvButton>
          <GvButton size="small" variant="tonal" @click="handleViewArtifacts(row)">
            产物
          </GvButton>
          <GvButton
            v-if="row.status === 'running'"
            size="small"
            variant="text"
            @click="handleCancel(row)"
          >
            取消
          </GvButton>
          <GvButton
            v-if="row.status !== 'running' && row.status !== 'pending'"
            size="small"
            variant="text"
            @click="handleDelete(row)"
          >
            删除
          </GvButton>
        </GvSpace>
      </template>
    </GvTable>

    <!-- 详情对话框 -->
    <GvModal
      v-model="showViewDialog"
      title="任务详情"
      size="large"
      :show-confirm="false"
      cancel-text="关闭"
    >
      <el-descriptions v-if="currentTask" :column="2" border>
        <el-descriptions-item label="任务 ID" :span="2">{{ currentTask.id }}</el-descriptions-item>
        <el-descriptions-item label="工作流">{{ currentTask.workflow?.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <StatusBadge :status="mapStatus(currentTask.status)" />
        </el-descriptions-item>
        <el-descriptions-item label="进度">{{ currentTask.progress }}%</el-descriptions-item>
        <el-descriptions-item label="当前节点">{{ currentTask.current_node || '-' }}</el-descriptions-item>
        <el-descriptions-item label="开始时间" :span="2">{{ currentTask.started_at ? formatDate(currentTask.started_at) : '-' }}</el-descriptions-item>
        <el-descriptions-item label="完成时间" :span="2">{{ currentTask.completed_at ? formatDate(currentTask.completed_at) : '-' }}</el-descriptions-item>
        <el-descriptions-item label="耗时" :span="2">{{ calculateDuration(currentTask) }}</el-descriptions-item>
        <el-descriptions-item v-if="currentTask.error" label="错误信息" :span="2">
          <GvAlert type="error" :closable="false">{{ currentTask.error }}</GvAlert>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ formatDate(currentTask.created_at) }}</el-descriptions-item>
      </el-descriptions>
    </GvModal>
  </GvContainer>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { taskApi, type Task, type TaskStats } from '@/api/task'
import { useTable, useAsyncData } from '@/composables'
import GvContainer from '@/components/layout/GvContainer/index.vue'
import GvGrid from '@/components/layout/GvGrid/index.vue'
import GvCard from '@/components/base/GvCard/index.vue'
import GvTable from '@/components/base/GvTable/index.vue'
import GvModal from '@/components/base/GvModal/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvSpace from '@/components/layout/GvSpace/index.vue'
import GvAlert from '@/components/base/GvAlert/index.vue'
import PageHeader from '@/components/business/PageHeader/index.vue'
import FilterBar from '@/components/business/FilterBar/index.vue'
import StatusBadge from '@/components/business/StatusBadge/index.vue'
import { ErrorState, EmptyState } from '@/components/common'
import type { TableColumn } from '@/components/base/GvTable/types'
import type { FilterField } from '@/components/business/FilterBar/types'

// UI 状态
const showViewDialog = ref(false)
const currentTask = ref<Task | null>(null)

const filters = ref({
  status: ''
})

// 计算筛选参数
const filterParams = computed(() => ({
  status: filters.value.status || undefined
}))

// 使用 useTable 管理任务列表
const {
  items: tasks,
  isLoading: loading,
  error,
  pagination,
  goToPage,
  changePageSize,
  refreshTable
} = useTable(
  (params) => taskApi.list(params),
  {
    immediate: true,
    initialPageSize: 20,
    extraParams: filterParams
  }
)

// 使用 useAsyncData 管理统计数据
const {
  data: statsData,
  execute: loadStats
} = useAsyncData(
  () => taskApi.getStats(),
  { immediate: true }
)

const stats = computed(() => statsData.value?.data || {
  total: 0,
  pending: 0,
  running: 0,
  success: 0,
  failed: 0,
  cancelled: 0
})

const statusOptions = [
  { label: '待执行', value: 'pending' },
  { label: '执行中', value: 'running' },
  { label: '已成功', value: 'success' },
  { label: '已失败', value: 'failed' },
  { label: '已取消', value: 'cancelled' }
]

const filterFields: FilterField[] = [
  {
    key: 'status',
    label: '任务状态',
    type: 'select',
    placeholder: '选择任务状态',
    options: statusOptions
  }
]

const columns: TableColumn[] = [
  { prop: 'id', label: '任务 ID', width: '280', showOverflowTooltip: true },
  { prop: 'workflow', label: '工作流', minWidth: '150' },
  { prop: 'status', label: '状态', width: '120' },
  { prop: 'progress', label: '进度', width: '120' },
  { prop: 'current_node', label: '当前节点', width: '120', showOverflowTooltip: true },
  { prop: 'started_at', label: '开始时间', width: '160' },
  { prop: 'completed_at', label: '完成时间', width: '160' },
  { prop: 'duration', label: '耗时', width: '100' },
  { prop: 'actions', label: '操作', width: '280', fixed: 'right' }
]

const paginationConfig = computed(() => ({
  currentPage: pagination.page,
  pageSize: pagination.pageSize,
  total: pagination.total
}))

function handleRefresh() {
  refreshTable()
  loadStats()
}

// 直接使用 useTable 提供的方法
const handlePageChange = goToPage
const handleSizeChange = changePageSize

function handleResetFilter() {
  pagination.page = 1
  // useTable 监听 pagination.page 变化会自动重新加载
}

function handleView(row: Task) {
  currentTask.value = row
  showViewDialog.value = true
}

function handleViewArtifacts(row: Task) {
  ElMessage.info('产物查看功能开发中')
}

async function handleCancel(row: Task) {
  try {
    await ElMessageBox.confirm('确定要取消此任务吗？', '提示', {
      type: 'warning'
    })
    await taskApi.cancel(row.id)
    ElMessage.success('取消成功')
    refreshTable()
    loadStats()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '取消失败')
    }
  }
}

async function handleDelete(row: Task) {
  try {
    await ElMessageBox.confirm('确定要删除此任务吗？', '提示', {
      type: 'warning'
    })
    await taskApi.delete(row.id)
    ElMessage.success('删除成功')
    refreshTable()
    loadStats()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

function mapStatus(status: string): any {
  const map: Record<string, string> = {
    pending: 'pending',
    running: 'running',
    success: 'success',
    failed: 'failed',
    cancelled: 'stopped'
  }
  return map[status] || 'inactive'
}

function getProgressStatus(status: string) {
  if (status === 'success') return 'success'
  if (status === 'failed') return 'exception'
  return undefined
}

function calculateDuration(task: Task): string {
  if (!task.started_at) return '-'
  const start = new Date(task.started_at).getTime()
  const end = task.completed_at ? new Date(task.completed_at).getTime() : Date.now()
  const duration = Math.floor((end - start) / 1000)
  
  if (duration < 60) return `${duration}秒`
  if (duration < 3600) return `${Math.floor(duration / 60)}分${duration % 60}秒`
  return `${Math.floor(duration / 3600)}时${Math.floor((duration % 3600) / 60)}分`
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}
</script>

<style scoped>
:deep(.el-descriptions) {
  @apply rounded-lg overflow-hidden;
}

:deep(.el-descriptions__label) {
  @apply font-semibold bg-neutral-50;
}

:deep(.el-descriptions__content) {
  @apply text-text-primary;
}

.dark :deep(.el-descriptions__label) {
  @apply bg-neutral-800 text-text-inverse;
}

.dark :deep(.el-descriptions__content) {
  @apply bg-surface-dark text-text-inverse;
}

.dark .text-text-tertiary {
  @apply text-neutral-400;
}
</style>
