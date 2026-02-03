<template>
  <GvContainer max-width="full">
    <!-- 页面头部 -->
    <PageHeader
      title="工作流管理"
      description="管理 AI 工作流，配置触发方式和执行流程"
    >
      <template #actions>
        <GvSpace>
          <SearchBar
            v-model="searchKeyword"
            placeholder="搜索工作流"
            class="w-80"
            immediate
            @search="loadWorkflows"
          />
          <GvButton @click="showCreateDialog = true">
            <template #icon>
              <el-icon><Plus /></el-icon>
            </template>
            创建工作流
          </GvButton>
        </GvSpace>
      </template>
    </PageHeader>

    <!-- 筛选栏 -->
    <FilterBar
      v-model="filters"
      :fields="filterFields"
      :columns="2"
      :loading="loading"
      @filter="loadWorkflows"
      @reset="handleResetFilter"
    />

    <!-- 数据表格 -->
    <GvTable
      :data="workflows"
      :columns="columns"
      :loading="loading"
      border
      stripe
      pagination
      :pagination-config="paginationConfig"
      @current-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <template #trigger_type="{ row }">
        <GvTag :color="getTriggerTypeColor(row.trigger_type)" size="small">
          {{ getTriggerTypeLabel(row.trigger_type) }}
        </GvTag>
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
          <GvButton size="small" @click="handleEdit(row)">
            编辑
          </GvButton>
          <GvButton size="small" @click="handleTrigger(row)">
            触发
          </GvButton>
          <GvButton
            v-if="row.status === 'published'"
            size="small"
            variant="text"
            @click="handleDisable(row)"
          >
            禁用
          </GvButton>
          <GvButton
            v-else
            size="small"
            variant="text"
            @click="handleEnable(row)"
          >
            启用
          </GvButton>
          <GvButton
            size="small"
            variant="text"
            @click="handleDelete(row)"
          >
            删除
          </GvButton>
        </GvSpace>
      </template>
    </GvTable>

    <!-- 创建对话框 -->
    <GvModal
      v-model="showCreateDialog"
      title="创建工作流"
      size="large"
      :confirm-loading="creating"
      @confirm="handleCreate"
      @cancel="showCreateDialog = false"
    >
      <el-form ref="createFormRef" :model="createForm" :rules="createRules" label-width="100px">
        <el-form-item label="工作流代码" prop="code">
          <GvInput v-model="createForm.code" placeholder="唯一标识" />
        </el-form-item>
        <el-form-item label="工作流名称" prop="name">
          <GvInput v-model="createForm.name" placeholder="工作流显示名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <GvInput v-model="createForm.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="触发方式" prop="trigger_type">
          <GvSelect
            v-model="createForm.trigger_type"
            :options="triggerTypeOptions"
          />
        </el-form-item>
      </el-form>
    </GvModal>

    <!-- 详情对话框 -->
    <GvModal
      v-model="showViewDialog"
      title="工作流详情"
      size="large"
      :show-confirm="false"
      cancel-text="关闭"
    >
      <el-descriptions v-if="currentWorkflow" :column="2" border>
        <el-descriptions-item label="ID" :span="2">{{ currentWorkflow.id }}</el-descriptions-item>
        <el-descriptions-item label="代码">{{ currentWorkflow.code }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ currentWorkflow.name }}</el-descriptions-item>
        <el-descriptions-item label="触发方式">
          <GvTag :color="getTriggerTypeColor(currentWorkflow.trigger_type)" size="small">
            {{ getTriggerTypeLabel(currentWorkflow.trigger_type) }}
          </GvTag>
        </el-descriptions-item>
        <el-descriptions-item label="版本">{{ currentWorkflow.version }}</el-descriptions-item>
        <el-descriptions-item label="状态" :span="2">
          <StatusBadge :status="mapStatus(currentWorkflow.status)" />
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentWorkflow.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ formatDate(currentWorkflow.created_at) }}</el-descriptions-item>
      </el-descriptions>
    </GvModal>

    <!-- 触发对话框 -->
    <GvModal
      v-model="showTriggerDialog"
      title="触发工作流"
      :confirm-loading="triggering"
      @confirm="handleTriggerConfirm"
      @cancel="showTriggerDialog = false"
    >
      <el-form label-width="100px">
        <el-form-item label="关联资产">
          <GvInput v-model="triggerForm.asset_id" placeholder="可选：输入资产 ID" />
        </el-form-item>
      </el-form>
    </GvModal>
  </GvContainer>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { workflowApi, type Workflow, type WorkflowCreateReq } from '@/api/workflow'
import { useRouter } from 'vue-router'
import GvContainer from '@/components/layout/GvContainer/index.vue'
import GvTable from '@/components/base/GvTable/index.vue'
import GvModal from '@/components/base/GvModal/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvSpace from '@/components/layout/GvSpace/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'
import GvInput from '@/components/base/GvInput/index.vue'
import GvSelect from '@/components/base/GvSelect/index.vue'
import PageHeader from '@/components/business/PageHeader/index.vue'
import FilterBar from '@/components/business/FilterBar/index.vue'
import SearchBar from '@/components/business/SearchBar/index.vue'
import StatusBadge from '@/components/business/StatusBadge/index.vue'
import type { TableColumn } from '@/components/base/GvTable/types'
import type { FilterField } from '@/components/business/FilterBar/types'

const router = useRouter()
const loading = ref(false)
const creating = ref(false)
const triggering = ref(false)
const workflows = ref<Workflow[]>([])
const showCreateDialog = ref(false)
const showViewDialog = ref(false)
const showTriggerDialog = ref(false)
const currentWorkflow = ref<Workflow | null>(null)
const createFormRef = ref<FormInstance>()

const searchKeyword = ref('')

const filters = ref({
  trigger_type: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const createForm = reactive<WorkflowCreateReq>({
  code: '',
  name: '',
  description: '',
  trigger_type: 'manual'
})

const triggerForm = reactive({
  asset_id: ''
})

const createRules: FormRules = {
  code: [{ required: true, message: '请输入工作流代码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入工作流名称', trigger: 'blur' }],
  trigger_type: [{ required: true, message: '请选择触发方式', trigger: 'change' }]
}

const triggerTypeOptions = [
  { label: '手动触发', value: 'manual' },
  { label: '定时调度', value: 'schedule' },
  { label: '事件触发', value: 'event' }
]

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '测试中', value: 'testing' },
  { label: '已发布', value: 'published' },
  { label: '已归档', value: 'archived' }
]

const filterFields: FilterField[] = [
  {
    key: 'trigger_type',
    label: '触发方式',
    type: 'select',
    placeholder: '选择触发方式',
    options: triggerTypeOptions
  },
  {
    key: 'status',
    label: '状态',
    type: 'select',
    placeholder: '选择状态',
    options: statusOptions
  }
]

const columns: TableColumn[] = [
  { prop: 'name', label: '名称', minWidth: '150', showOverflowTooltip: true },
  { prop: 'code', label: '代码', width: '150' },
  { prop: 'trigger_type', label: '触发方式', width: '100' },
  { prop: 'version', label: '版本', width: '80' },
  { prop: 'status', label: '状态', width: '120' },
  { prop: 'created_at', label: '创建时间', width: '160' },
  { prop: 'actions', label: '操作', width: '360', fixed: 'right' }
]

const paginationConfig = computed(() => ({
  currentPage: pagination.page,
  pageSize: pagination.page_size,
  total: pagination.total
}))

onMounted(() => {
  loadWorkflows()
})

async function loadWorkflows() {
  loading.value = true
  try {
    const response = await workflowApi.list({
      keyword: searchKeyword.value || undefined,
      trigger_type: filters.value.trigger_type as any,
      status: filters.value.status as any,
      page: pagination.page,
      page_size: pagination.page_size
    })
    workflows.value = response.data.items
    pagination.total = response.data.total
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function handleCreate() {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (!valid) return
    creating.value = true
    try {
      await workflowApi.create(createForm)
      ElMessage.success('创建成功')
      showCreateDialog.value = false
      loadWorkflows()
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || '创建失败')
    } finally {
      creating.value = false
    }
  })
}

function handleView(row: Workflow) {
  currentWorkflow.value = row
  showViewDialog.value = true
}

function handleEdit(row: Workflow) {
  ElMessage.info('编辑功能开发中')
}

function handleTrigger(row: Workflow) {
  currentWorkflow.value = row
  triggerForm.asset_id = ''
  showTriggerDialog.value = true
}

async function handleTriggerConfirm() {
  if (!currentWorkflow.value) return
  triggering.value = true
  try {
    const assetId = triggerForm.asset_id || undefined
    await workflowApi.trigger(currentWorkflow.value.id, assetId)
    ElMessage.success('工作流已触发，任务已创建')
    showTriggerDialog.value = false
    router.push('/tasks')
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '触发失败')
  } finally {
    triggering.value = false
  }
}

async function handleEnable(row: Workflow) {
  try {
    await workflowApi.enable(row.id)
    ElMessage.success('启用成功')
    loadWorkflows()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '启用失败')
  }
}

async function handleDisable(row: Workflow) {
  try {
    await workflowApi.disable(row.id)
    ElMessage.success('禁用成功')
    loadWorkflows()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '禁用失败')
  }
}

async function handleDelete(row: Workflow) {
  try {
    await ElMessageBox.confirm('确定要删除此工作流吗？', '提示', {
      type: 'warning'
    })
    await workflowApi.delete(row.id)
    ElMessage.success('删除成功')
    loadWorkflows()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

function handlePageChange(page: number) {
  pagination.page = page
  loadWorkflows()
}

function handleSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  loadWorkflows()
}

function handleResetFilter() {
  searchKeyword.value = ''
  loadWorkflows()
}

function getTriggerTypeLabel(type: string) {
  const map: Record<string, string> = {
    manual: '手动',
    schedule: '定时',
    event: '事件'
  }
  return map[type] || type
}

function getTriggerTypeColor(type: string) {
  const map: Record<string, string> = {
    manual: 'primary',
    schedule: 'success',
    event: 'warning'
  }
  return map[type] || 'neutral'
}

function mapStatus(status: string): any {
  const map: Record<string, string> = {
    draft: 'pending',
    testing: 'processing',
    published: 'active',
    archived: 'disabled'
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
