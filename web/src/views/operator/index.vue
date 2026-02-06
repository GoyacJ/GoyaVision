<template>
  <GvContainer max-width="full">
    <!-- 页面头部 -->
    <PageHeader
      title="算子中心"
      description="管理 AI 算子，包括内置和自定义算子"
    >
      <template #actions>
        <GvSpace>
          <SearchBar
            v-model="searchKeyword"
            placeholder="搜索算子"
            class="w-80"
            immediate
            :show-button="false"
            @search="() => { pagination.page = 1 }"
          />
          <GvButton @click="showCreateDialog = true">
            <template #icon>
              <el-icon><Plus /></el-icon>
            </template>
            添加算子
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
      @action="showCreateDialog = true"
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
      :confirm-loading="creating"
      @confirm="handleCreate"
      @cancel="showCreateDialog = false"
    >
      <el-form ref="createFormRef" :model="createForm" :rules="createRules" label-width="100px">
        <el-form-item label="算子代码" prop="code">
          <GvInput v-model="createForm.code" placeholder="唯一标识，如：frame_extract" />
        </el-form-item>
        <el-form-item label="算子名称" prop="name">
          <GvInput v-model="createForm.name" placeholder="算子显示名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <GvInput v-model="createForm.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <GvSelect
            v-model="createForm.category"
            :options="categoryOptions"
          />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <GvInput v-model="createForm.type" placeholder="如：object_detection, frame_extract" />
        </el-form-item>
        <el-form-item label="端点地址" prop="endpoint">
          <GvInput v-model="createForm.endpoint" placeholder="http://..." />
        </el-form-item>
        <el-form-item label="HTTP方法" prop="method">
          <GvSelect
            v-model="createForm.method"
            :options="methodOptions"
          />
        </el-form-item>
      </el-form>
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
        <el-descriptions-item label="版本">{{ currentOperator.version }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <StatusBadge :status="mapStatus(currentOperator.status)" />
        </el-descriptions-item>
        <el-descriptions-item label="端点地址" :span="2">{{ currentOperator.endpoint }}</el-descriptions-item>
        <el-descriptions-item label="HTTP方法">{{ currentOperator.method }}</el-descriptions-item>
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
  </GvContainer>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { operatorApi, type Operator, type OperatorCreateReq } from '@/api/operator'
import { useTable } from '@/composables'
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
import { ErrorState, EmptyState } from '@/components/common'
import type { TableColumn } from '@/components/base/GvTable/types'
import type { FilterField } from '@/components/business/FilterBar/types'

// UI 状态
const creating = ref(false)
const showCreateDialog = ref(false)
const showViewDialog = ref(false)
const currentOperator = ref<Operator | null>(null)
const createFormRef = ref<FormInstance>()

const searchKeyword = ref('')

const filters = ref({
  category: '',
  status: '',
  is_builtin: ''
})

// 计算筛选参数
const filterParams = computed(() => ({
  keyword: searchKeyword.value || undefined,
  category: filters.value.category || undefined,
  status: filters.value.status || undefined,
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

const createForm = reactive<OperatorCreateReq>({
  code: '',
  name: '',
  description: '',
  category: 'analysis',
  type: '',
  endpoint: '',
  method: 'POST'
})

const createRules: FormRules = {
  code: [{ required: true, message: '请输入算子代码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入算子名称', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  type: [{ required: true, message: '请输入类型', trigger: 'blur' }],
  endpoint: [{ required: true, message: '请输入端点地址', trigger: 'blur' }]
}

const categoryOptions = [
  { label: '分析', value: 'analysis' },
  { label: '处理', value: 'processing' },
  { label: '生成', value: 'generation' }
]

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '测试中', value: 'testing' },
  { label: '已发布', value: 'published' },
  { label: '已废弃', value: 'deprecated' }
]

const methodOptions = [
  { label: 'POST', value: 'POST' },
  { label: 'GET', value: 'GET' }
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
  { prop: 'type', label: '类型', width: '120' },
  { prop: 'version', label: '版本', width: '80' },
  { prop: 'is_builtin', label: '内置', width: '80' },
  { prop: 'status', label: '状态', width: '120' },
  { prop: 'created_at', label: '创建时间', width: '160' },
  { prop: 'actions', label: '操作', width: '320', fixed: 'right' }
]

const paginationConfig = computed(() => ({
  currentPage: pagination.page,
  pageSize: pagination.pageSize,
  total: pagination.total
}))


async function handleCreate() {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (!valid) return
    creating.value = true
    try {
      await operatorApi.create(createForm)
      ElMessage.success('创建成功')
      showCreateDialog.value = false
      refreshTable()
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || '创建失败')
    } finally {
      creating.value = false
    }
  })
}

function handleView(row: Operator) {
  currentOperator.value = row
  showViewDialog.value = true
}

function handleEdit(row: Operator) {
  ElMessage.info('编辑功能开发中')
}

async function handleEnable(row: Operator) {
  try {
    await operatorApi.enable(row.id)
    ElMessage.success('启用成功')
    refreshTable()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '启用失败')
  }
}

async function handleDisable(row: Operator) {
  try {
    await operatorApi.disable(row.id)
    ElMessage.success('禁用成功')
    refreshTable()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '禁用失败')
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
    generation: '生成'
  }
  return map[category] || category
}

function getCategoryColor(category: string) {
  const map: Record<string, string> = {
    analysis: 'primary',
    processing: 'success',
    generation: 'warning'
  }
  return map[category] || 'neutral'
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
