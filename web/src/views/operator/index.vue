<template>
  <GvContainer max-width="full">
    <!-- 页面头部 -->
    <PageHeader
      title="算子中心"
      description="管理 AI 算子，包括内置和自定义算子"
    >
      <template #actions>
        <GvSpace wrap>
          <SearchBar
            v-model="searchKeyword"
            placeholder="搜索算子"
            class="w-full sm:w-80"
            immediate
            :show-button="false"
            @search="() => { pagination.page = 1 }"
          />
          <GvButton @click="openCreateDialog">
            <template #icon>
              <el-icon><Plus /></el-icon>
            </template>
            添加算子
          </GvButton>
          <GvButton variant="tonal" @click="router.push('/operator-marketplace')">
            模板市场
          </GvButton>
        </GvSpace>
      </template>
    </PageHeader>

    <!-- 筛选栏 -->
    <FilterBar
      v-model="filters"
      :fields="filterFields"
      :loading="loading"
      @filter="() => { pagination.page = 1 }"
      @reset="handleResetFilter"
    />

    <ErrorState
      v-if="error && !loading"
      :error="error"
      title="加载失败"
      @retry="refreshTable"
    />

    <EmptyState
      v-else-if="!loading && operators.length === 0"
      icon="⚙️"
      title="还没有算子"
      description="添加 AI 算子以处理您的媒体资产"
      action-text="添加算子"
      show-action
      @action="openCreateDialog"
    />

    <!-- 数据表格 -->
    <GvTable
      v-else
      :data="operators"
      :columns="columns"
      :loading="loading"
      border
      stripe
      pagination
      :pagination-config="paginationConfig"
      @current-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <template #category="{ row }">
        <GvTag :color="getCategoryColor(row.category)" size="small">
          {{ getCategoryLabel(row.category) }}
        </GvTag>
      </template>

      <template #origin="{ row }">
        {{ getOriginLabel(row.origin, row.is_builtin) }}
      </template>

      <template #exec_mode="{ row }">
        {{ getExecModeLabel(row.exec_mode || row.active_version?.exec_mode) }}
      </template>
      
      <template #is_builtin="{ row }">
        <GvTag v-if="row.is_builtin" color="info" size="small" variant="tonal">
          内置
        </GvTag>
        <span v-else class="text-text-tertiary">-</span>
      </template>
      
      <template #status="{ row }">
        <StatusBadge :status="mapStatus(row.status)" />
      </template>
      
      <template #created_at="{ row }">
        {{ formatDate(row.created_at) }}
      </template>
      
      <template #actions="{ row }">
        <GvSpace size="xs">
          <GvButton size="small" variant="tonal" @click="handleView(row)">
            查看
          </GvButton>
          <GvButton size="small" variant="text" @click="handleOpenVersions(row)">
            版本
          </GvButton>
          <GvButton
            size="small"
            @click="handleEdit(row)"
            :disabled="row.is_builtin"
          >
            编辑
          </GvButton>
          <GvButton
            v-if="row.status === 'published'"
            size="small"
            variant="text"
            @click="handleDeprecate(row)"
          >
            弃用
          </GvButton>
          <GvButton
            v-else
            size="small"
            variant="text"
            @click="handlePublish(row)"
          >
            发布
          </GvButton>
          <GvButton
            size="small"
            variant="text"
            @click="handleTest(row)"
          >
            测试
          </GvButton>
          <GvButton
            size="small"
            variant="text"
            @click="handleDelete(row)"
            :disabled="row.is_builtin"
          >
            删除
          </GvButton>
        </GvSpace>
      </template>
    </GvTable>

    <!-- 创建对话框 -->
    <GvModal
      v-model="showCreateDialog"
      title="添加算子"
      size="large"
      confirm-text="保存"
      :confirm-loading="creating"
      @confirm="handleCreateConfirm"
      @cancel="showCreateDialog = false"
    >
      <OperatorForm
        ref="createFormRef"
        :model-value="operatorFormModel"
        :loading="creating"
        @submit="handleCreateSubmit"
        @cancel="showCreateDialog = false"
      />
    </GvModal>

    <GvModal
      v-model="showEditDialog"
      title="编辑算子"
      size="large"
      confirm-text="保存"
      :confirm-loading="editing"
      @confirm="handleEditConfirm"
      @cancel="showEditDialog = false"
    >
      <OperatorForm
        ref="editFormRef"
        :model-value="operatorFormModel"
        :loading="editing"
        @submit="handleEditSubmit"
        @cancel="showEditDialog = false"
      />
    </GvModal>

    <!-- 详情对话框 -->
    <GvModal
      v-model="showViewDialog"
      title="算子详情"
      size="large"
      :show-confirm="false"
      cancel-text="关闭"
    >
      <el-descriptions v-if="currentOperator" :column="2" border>
        <el-descriptions-item label="ID" :span="2">{{ currentOperator.id }}</el-descriptions-item>
        <el-descriptions-item label="代码">{{ currentOperator.code }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ currentOperator.name }}</el-descriptions-item>
        <el-descriptions-item label="分类">
          <GvTag :color="getCategoryColor(currentOperator.category)" size="small">
            {{ getCategoryLabel(currentOperator.category) }}
          </GvTag>
        </el-descriptions-item>
        <el-descriptions-item label="类型">{{ currentOperator.type }}</el-descriptions-item>
        <el-descriptions-item label="来源">{{ getOriginLabel(currentOperator.origin, currentOperator.is_builtin) }}</el-descriptions-item>
        <el-descriptions-item label="执行模式">{{ getExecModeLabel(currentOperator.exec_mode || currentOperator.active_version?.exec_mode) }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ currentOperator.version }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <StatusBadge :status="mapStatus(currentOperator.status)" />
        </el-descriptions-item>
        <el-descriptions-item label="激活版本ID" :span="2">{{ currentOperator.active_version_id || '-' }}</el-descriptions-item>
        <el-descriptions-item label="执行配置" :span="2">
          <pre class="max-h-64 overflow-auto rounded bg-neutral-50 p-2 text-xs">{{ formatJson(currentOperator.active_version?.exec_config || {}) }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="内置">
          <GvTag v-if="currentOperator.is_builtin" color="info" size="small" variant="tonal">
            是
          </GvTag>
          <span v-else>否</span>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentOperator.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ formatDate(currentOperator.created_at) }}</el-descriptions-item>
      </el-descriptions>
    </GvModal>

    <GvModal
      v-model="showVersionDialog"
      :title="`版本与依赖管理 - ${currentOperator?.name || ''}`"
      size="large"
      :show-confirm="false"
      cancel-text="关闭"
    >
      <el-tabs>
        <el-tab-pane label="版本列表">
          <VersionList
            :versions="versionList"
            :loading="versionLoading"
            @activate="handleActivateVersion"
            @rollback="handleRollbackVersion"
            @archive="handleArchiveVersion"
          />
        </el-tab-pane>
        <el-tab-pane label="创建版本">
          <VersionForm :loading="versionSubmitting" @submit="handleCreateVersion" />
        </el-tab-pane>
        <el-tab-pane label="依赖管理">
          <div class="mb-3">
            <GvButton size="small" variant="tonal" :loading="dependencyChecking" @click="handleCheckDependencies">
              检查依赖满足性
            </GvButton>
          </div>
          <DependencyManager
            :dependencies="dependencyList"
            :loading="dependencyLoading"
            @save="handleSaveDependencies"
          />
        </el-tab-pane>
      </el-tabs>
    </GvModal>
  </GvContainer>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { operatorApi, type Operator, type OperatorCreateReq } from '@/api/operator'
import { useTable } from '@/composables'
import GvContainer from '@/components/layout/GvContainer/index.vue'
import GvTable from '@/components/base/GvTable/index.vue'
import GvModal from '@/components/base/GvModal/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvSpace from '@/components/layout/GvSpace/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'
import PageHeader from '@/components/business/PageHeader/index.vue'
import FilterBar from '@/components/business/FilterBar/index.vue'
import SearchBar from '@/components/business/SearchBar/index.vue'
import StatusBadge from '@/components/business/StatusBadge/index.vue'
import { ErrorState, EmptyState } from '@/components/common'
import type { TableColumn } from '@/components/base/GvTable/types'
import type { FilterField } from '@/components/business/FilterBar/types'
import VersionList from './components/VersionList.vue'
import VersionForm from './components/VersionForm.vue'
import DependencyManager from './components/DependencyManager.vue'
import OperatorForm from './components/OperatorForm.vue'

// UI 状态
const creating = ref(false)
const editing = ref(false)
const router = useRouter()
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const showViewDialog = ref(false)
const showVersionDialog = ref(false)
const currentOperator = ref<Operator | null>(null)
const editOperatorID = ref('')
const versionList = ref<Array<any>>([])
const dependencyList = ref<Array<any>>([])
const versionLoading = ref(false)
const versionSubmitting = ref(false)
const dependencyLoading = ref(false)
const dependencyChecking = ref(false)
const createFormRef = ref<InstanceType<typeof OperatorForm>>()
const editFormRef = ref<InstanceType<typeof OperatorForm>>()

const searchKeyword = ref('')

const filters = ref({
  category: '',
  status: '',
  origin: '',
  exec_mode: '',
  is_builtin: ''
})

const operatorFormModel = reactive<OperatorCreateReq>({
  code: '',
  name: '',
  description: '',
  category: 'analysis',
  type: '',
  origin: 'custom',
  exec_mode: 'http',
  exec_config: {}
})

// 计算筛选参数
const filterParams = computed(() => ({
  keyword: searchKeyword.value || undefined,
  category: filters.value.category || undefined,
  status: filters.value.status || undefined,
  origin: filters.value.origin || undefined,
  exec_mode: filters.value.exec_mode || undefined,
  is_builtin: filters.value.is_builtin ? filters.value.is_builtin === 'true' : undefined
}))

// 使用 useTable 管理算子列表
const {
  items: operators,
  isLoading: loading,
  error,
  pagination,
  goToPage,
  changePageSize,
  refreshTable
} = useTable(
  async (params) => {
    const res = await operatorApi.list(params)
    return { items: res.data?.items ?? [], total: res.data?.total ?? 0 }
  },
  {
    immediate: true,
    initialPageSize: 20,
    extraParams: filterParams
  }
)

const categoryOptions = [
  { label: '分析', value: 'analysis' },
  { label: '处理', value: 'processing' },
  { label: '生成', value: 'generation' },
  { label: '工具', value: 'utility' }
]

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '测试中', value: 'testing' },
  { label: '已发布', value: 'published' },
  { label: '已废弃', value: 'deprecated' }
]

const filterFields: FilterField[] = [
  {
    key: 'category',
    label: '分类',
    type: 'select',
    placeholder: '选择分类',
    options: categoryOptions
  },
  {
    key: 'status',
    label: '状态',
    type: 'select',
    placeholder: '选择状态',
    options: statusOptions
  },
  {
    key: 'origin',
    label: '来源',
    type: 'select',
    placeholder: '选择来源',
    options: [
      { label: '内置', value: 'builtin' },
      { label: '自定义', value: 'custom' },
      { label: '模板市场', value: 'marketplace' },
      { label: 'MCP', value: 'mcp' }
    ]
  },
  {
    key: 'exec_mode',
    label: '执行模式',
    type: 'select',
    placeholder: '选择执行模式',
    options: [
      { label: 'HTTP', value: 'http' },
      { label: 'CLI', value: 'cli' },
      { label: 'MCP', value: 'mcp' }
    ]
  },
  {
    key: 'is_builtin',
    label: '内置',
    type: 'select',
    placeholder: '是否内置',
    options: [
      { label: '内置算子', value: 'true' },
      { label: '自定义算子', value: 'false' }
    ]
  }
]

const columns: TableColumn[] = [
  { prop: 'name', label: '名称', minWidth: '150', showOverflowTooltip: true },
  { prop: 'code', label: '代码', width: '150' },
  { prop: 'category', label: '分类', width: '100' },
  { prop: 'origin', label: '来源', width: '100' },
  { prop: 'exec_mode', label: '执行模式', width: '100' },
  { prop: 'type', label: '类型', width: '120' },
  { prop: 'version', label: '版本', width: '80' },
  { prop: 'is_builtin', label: '内置', width: '80' },
  { prop: 'status', label: '状态', width: '120' },
  { prop: 'created_at', label: '创建时间', width: '160' },
  { prop: 'actions', label: '操作', width: '420', fixed: 'right' }
]

const paginationConfig = computed(() => ({
  currentPage: pagination.page,
  pageSize: pagination.pageSize,
  total: pagination.total
}))


function resetOperatorForm() {
  Object.assign(operatorFormModel, {
    code: '',
    name: '',
    description: '',
    category: 'analysis',
    type: '',
    origin: 'custom',
    exec_mode: 'http',
    exec_config: {}
  })
}

function openCreateDialog() {
  resetOperatorForm()
  showCreateDialog.value = true
}

function handleCreateConfirm() {
  createFormRef.value?.submit()
}

async function handleCreateSubmit(payload: OperatorCreateReq) {
  if (!payload.code || !payload.name || !payload.type || !payload.category) {
    ElMessage.warning('请完整填写算子代码、名称、分类与类型')
    return
  }
  creating.value = true
  try {
    await operatorApi.create(payload)
    ElMessage.success('创建成功')
    showCreateDialog.value = false
    resetOperatorForm()
    refreshTable()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '创建失败')
  } finally {
    creating.value = false
  }
}

function handleView(row: Operator) {
  currentOperator.value = row
  showViewDialog.value = true
}

function handleEdit(row: Operator) {
  editOperatorID.value = row.id
  Object.assign(operatorFormModel, {
    code: row.code,
    name: row.name,
    description: row.description || '',
    category: row.category,
    type: row.type,
    origin: row.origin || (row.is_builtin ? 'builtin' : 'custom'),
    exec_mode: row.exec_mode || row.active_version?.exec_mode || 'http',
    exec_config: row.active_version?.exec_config || {}
  })
  showEditDialog.value = true
}

function handleEditConfirm() {
  editFormRef.value?.submit()
}

async function handleEditSubmit(payload: OperatorCreateReq) {
  if (!editOperatorID.value) return
  editing.value = true
  try {
    await operatorApi.update(editOperatorID.value, {
      name: payload.name,
      description: payload.description,
      category: payload.category,
      tags: payload.tags
    })
    ElMessage.success('更新成功')
    ElMessage.warning('执行模式/执行配置变更需通过“创建版本”完成')
    showEditDialog.value = false
    refreshTable()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '更新失败')
  } finally {
    editing.value = false
  }
}

async function handlePublish(row: Operator) {
  try {
    await operatorApi.publish(row.id)
    ElMessage.success('发布成功')
    refreshTable()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '发布失败')
  }
}

async function handleDeprecate(row: Operator) {
  try {
    await operatorApi.deprecate(row.id)
    ElMessage.success('弃用成功')
    refreshTable()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '弃用失败')
  }
}

async function handleTest(row: Operator) {
  try {
    const res = await operatorApi.test(row.id)
    const message = res.data?.message || '测试完成'
    ElMessage.success(message)
    if (res.data?.diagnostics) {
      await ElMessageBox.alert(`<pre style="max-height:360px;overflow:auto">${formatJson(res.data.diagnostics)}</pre>`, '测试诊断信息', {
        dangerouslyUseHTMLString: true,
        confirmButtonText: '关闭'
      })
    }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '测试失败')
  }
}

async function handleOpenVersions(row: Operator) {
  currentOperator.value = row
  showVersionDialog.value = true
  await Promise.all([loadVersions(row.id), loadDependencies(row.id)])
}

async function loadVersions(operatorId: string) {
  versionLoading.value = true
  try {
    const res = await operatorApi.listVersions(operatorId, { page: 1, page_size: 100 })
    versionList.value = res.data?.items || []
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '加载版本失败')
  } finally {
    versionLoading.value = false
  }
}

async function loadDependencies(operatorId: string) {
  dependencyLoading.value = true
  try {
    const res = await operatorApi.listDependencies(operatorId)
    dependencyList.value = res.data || []
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '加载依赖失败')
  } finally {
    dependencyLoading.value = false
  }
}

async function handleCreateVersion(payload: any) {
  if (!currentOperator.value) return
  versionSubmitting.value = true
  try {
    await operatorApi.createVersion(currentOperator.value.id, payload)
    ElMessage.success('创建版本成功')
    await loadVersions(currentOperator.value.id)
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '创建版本失败')
  } finally {
    versionSubmitting.value = false
  }
}

async function handleActivateVersion(versionId: string) {
  if (!currentOperator.value) return
  try {
    await operatorApi.activateVersion(currentOperator.value.id, { version_id: versionId })
    ElMessage.success('激活成功')
    await loadVersions(currentOperator.value.id)
    refreshTable()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '激活失败')
  }
}

async function handleRollbackVersion(versionId: string) {
  if (!currentOperator.value) return
  try {
    await operatorApi.rollbackVersion(currentOperator.value.id, { version_id: versionId })
    ElMessage.success('回滚成功')
    await loadVersions(currentOperator.value.id)
    refreshTable()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '回滚失败')
  }
}

async function handleArchiveVersion(versionId: string) {
  if (!currentOperator.value) return
  try {
    await operatorApi.archiveVersion(currentOperator.value.id, { version_id: versionId })
    ElMessage.success('归档成功')
    await loadVersions(currentOperator.value.id)
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '归档失败')
  }
}

async function handleSaveDependencies(deps: Array<{ depends_on_id: string; min_version?: string; is_optional?: boolean }>) {
  if (!currentOperator.value) return
  try {
    await operatorApi.setDependencies(currentOperator.value.id, {
      dependencies: deps.filter((d) => d.depends_on_id)
    })
    ElMessage.success('依赖保存成功')
    await loadDependencies(currentOperator.value.id)
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '依赖保存失败')
  }
}

async function handleCheckDependencies() {
  if (!currentOperator.value) return
  dependencyChecking.value = true
  try {
    const res = await operatorApi.checkDependencies(currentOperator.value.id)
    if (res.data?.satisfied) {
      ElMessage.success('依赖检查通过')
      return
    }

    const unmet = res.data?.unmet || []
    await ElMessageBox.alert(
      `<div>以下依赖未满足：</div><pre style="max-height:280px;overflow:auto">${unmet.join('\n') || '(空)'}</pre>`,
      '依赖检查未通过',
      { dangerouslyUseHTMLString: true, confirmButtonText: '关闭' }
    )
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '依赖检查失败')
  } finally {
    dependencyChecking.value = false
  }
}

async function handleDelete(row: Operator) {
  try {
    await ElMessageBox.confirm('确定要删除此算子吗？', '提示', {
      type: 'warning'
    })
    await operatorApi.delete(row.id)
    ElMessage.success('删除成功')
    refreshTable()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

// 直接使用 useTable 提供的方法
const handlePageChange = goToPage
const handleSizeChange = changePageSize

function handleResetFilter() {
  searchKeyword.value = ''
  pagination.page = 1
  // useTable 监听 pagination.page 变化会自动重新加载
}

function getCategoryLabel(category: string) {
  const map: Record<string, string> = {
    analysis: '分析',
    processing: '处理',
    generation: '生成',
    utility: '工具'
  }
  return map[category] || category
}

function getCategoryColor(category: string) {
  const map: Record<string, string> = {
    analysis: 'primary',
    processing: 'success',
    generation: 'warning',
    utility: 'info'
  }
  return map[category] || 'neutral'
}

function getOriginLabel(origin?: string, isBuiltin?: boolean) {
  const resolved = origin || (isBuiltin ? 'builtin' : 'custom')
  const map: Record<string, string> = {
    builtin: '内置',
    custom: '自定义',
    marketplace: '市场',
    mcp: 'MCP'
  }
  return map[resolved] || resolved
}

function getExecModeLabel(execMode?: string) {
  const map: Record<string, string> = {
    http: 'HTTP',
    cli: 'CLI',
    mcp: 'MCP'
  }
  return map[execMode || ''] || '-'
}

function formatJson(value: any) {
  if (!value || typeof value !== 'object') return '{}'
  return JSON.stringify(value, null, 2)
}

function mapStatus(status: string): any {
  const map: Record<string, string> = {
    draft: 'pending',
    testing: 'processing',
    published: 'active',
    deprecated: 'disabled'
  }
  return map[status] || 'inactive'
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
</style>
