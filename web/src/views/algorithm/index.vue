<template>
  <GvContainer max-width="full">
    <PageHeader
      title="算法库"
      description="管理算法意图层资产，维护版本、实现绑定与发布状态"
    >
      <template #actions>
        <GvSpace wrap>
          <SearchBar
            v-model="searchKeyword"
            placeholder="搜索算法名称/编码"
            class="w-full sm:w-80"
            immediate
            :show-button="false"
            @search="() => { pagination.page = 1 }"
          />
          <GvButton @click="openCreateDialog">
            <template #icon>
              <el-icon><Plus /></el-icon>
            </template>
            新建算法
          </GvButton>
        </GvSpace>
      </template>
    </PageHeader>

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
      v-else-if="!loading && algorithms.length === 0"
      icon="🧠"
      title="还没有算法"
      description="创建算法后可在工作流节点通过 algorithm_ref 直接引用"
      action-text="新建算法"
      show-action
      @action="openCreateDialog"
    />

    <GvTable
      v-else
      :data="algorithms"
      :columns="columns"
      :loading="loading"
      border
      stripe
      pagination
      :pagination-config="paginationConfig"
      @current-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <template #status="{ row }">
        <StatusBadge :status="mapStatus(row.status)" />
      </template>

      <template #updated_at="{ row }">
        {{ formatDate(row.updated_at) }}
      </template>

      <template #tags="{ row }">
        <GvSpace size="xs" wrap>
          <GvTag v-for="tag in (row.tags || []).slice(0, 3)" :key="`${row.id}-${tag}`" size="small" variant="tonal">
            {{ tag }}
          </GvTag>
          <span v-if="(row.tags || []).length > 3" class="text-xs text-text-tertiary">
            +{{ (row.tags || []).length - 3 }}
          </span>
          <span v-if="!row.tags || row.tags.length === 0" class="text-text-tertiary">-</span>
        </GvSpace>
      </template>

      <template #actions="{ row }">
        <GvSpace size="xs">
          <GvButton size="small" variant="tonal" @click="handleView(row)">查看</GvButton>
          <GvButton size="small" @click="handleEdit(row)">编辑</GvButton>
          <GvButton size="small" variant="text" @click="handleManageVersions(row)">版本</GvButton>
          <GvButton size="small" variant="text" @click="handleDelete(row)">删除</GvButton>
        </GvSpace>
      </template>
    </GvTable>

    <GvModal
      v-model="showCreateDialog"
      title="新建算法"
      size="large"
      confirm-text="保存"
      :confirm-loading="creating"
      @confirm="handleCreateConfirm"
      @cancel="showCreateDialog = false"
    >
      <el-form ref="createFormRef" :model="algorithmForm" :rules="algorithmFormRules" label-width="96px">
        <el-form-item label="算法编码" prop="code">
          <el-input v-model="algorithmForm.code" placeholder="例如: forest-inspection" />
        </el-form-item>
        <el-form-item label="算法名称" prop="name">
          <el-input v-model="algorithmForm.name" placeholder="例如: 森林巡检" />
        </el-form-item>
        <el-form-item label="应用场景" prop="scenario">
          <el-input v-model="algorithmForm.scenario" placeholder="例如: forestry" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="algorithmForm.status" class="w-full">
            <el-option label="草稿" value="draft" />
            <el-option label="已发布" value="published" />
            <el-option label="已弃用" value="deprecated" />
          </el-select>
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="algorithmForm.tagsText" placeholder="多个标签请用逗号分隔" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="algorithmForm.description" type="textarea" :rows="4" placeholder="填写算法用途、约束和说明" />
        </el-form-item>
      </el-form>
    </GvModal>

    <GvModal
      v-model="showEditDialog"
      title="编辑算法"
      size="large"
      confirm-text="保存"
      :confirm-loading="editing"
      @confirm="handleEditConfirm"
      @cancel="showEditDialog = false"
    >
      <el-form ref="editFormRef" :model="algorithmForm" :rules="algorithmFormRules" label-width="96px">
        <el-form-item label="算法编码">
          <el-input v-model="algorithmForm.code" disabled />
        </el-form-item>
        <el-form-item label="算法名称" prop="name">
          <el-input v-model="algorithmForm.name" />
        </el-form-item>
        <el-form-item label="应用场景" prop="scenario">
          <el-input v-model="algorithmForm.scenario" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="algorithmForm.status" class="w-full">
            <el-option label="草稿" value="draft" />
            <el-option label="已发布" value="published" />
            <el-option label="已弃用" value="deprecated" />
          </el-select>
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="algorithmForm.tagsText" placeholder="多个标签请用逗号分隔" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="algorithmForm.description" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
    </GvModal>

    <GvModal
      v-model="showViewDialog"
      title="算法详情"
      size="large"
      :show-confirm="false"
      cancel-text="关闭"
    >
      <el-descriptions v-if="currentAlgorithm" :column="2" border>
        <el-descriptions-item label="ID" :span="2">{{ currentAlgorithm.id }}</el-descriptions-item>
        <el-descriptions-item label="编码">{{ currentAlgorithm.code }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ currentAlgorithm.name }}</el-descriptions-item>
        <el-descriptions-item label="场景">{{ currentAlgorithm.scenario || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <StatusBadge :status="mapStatus(currentAlgorithm.status)" />
        </el-descriptions-item>
        <el-descriptions-item label="标签" :span="2">
          <GvSpace size="xs" wrap>
            <GvTag v-for="tag in currentAlgorithm.tags || []" :key="`detail-${tag}`" size="small" variant="tonal">{{ tag }}</GvTag>
            <span v-if="!currentAlgorithm.tags || currentAlgorithm.tags.length === 0">-</span>
          </GvSpace>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentAlgorithm.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="版本数量">{{ currentAlgorithm.versions?.length || 0 }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentAlgorithm.updated_at) }}</el-descriptions-item>
      </el-descriptions>
    </GvModal>

    <GvModal
      v-model="showVersionDialog"
      title="算法版本管理"
      size="large"
      :show-confirm="false"
      cancel-text="关闭"
    >
      <div v-if="currentAlgorithm" class="space-y-4">
        <div class="text-sm text-text-tertiary">
          当前算法：<span class="text-text-primary font-medium">{{ currentAlgorithm.name }}</span>
          <span class="ml-2 font-mono">{{ currentAlgorithm.code }}</span>
        </div>

        <el-table :data="currentAlgorithm.versions || []" border stripe>
          <el-table-column prop="version" label="版本" min-width="120" />
          <el-table-column prop="status" label="状态" min-width="110">
            <template #default="{ row }">
              <StatusBadge :status="mapVersionStatus(row.status)" />
            </template>
          </el-table-column>
          <el-table-column prop="selection_policy" label="选择策略" min-width="130" />
          <el-table-column label="实现数" min-width="90">
            <template #default="{ row }">
              {{ row.implementations?.length || 0 }}
            </template>
          </el-table-column>
          <el-table-column label="评测数" min-width="90">
            <template #default="{ row }">
              {{ row.evaluations?.length || 0 }}
            </template>
          </el-table-column>
          <el-table-column label="操作" min-width="120" fixed="right">
            <template #default="{ row }">
              <GvButton
                size="small"
                variant="text"
                :loading="publishingVersionID === row.id"
                :disabled="row.status === 'published'"
                @click="handlePublishVersion(row.id)"
              >
                发布
              </GvButton>
            </template>
          </el-table-column>
        </el-table>

        <el-divider content-position="left">新增版本</el-divider>

        <el-form ref="versionFormRef" :model="versionForm" :rules="versionFormRules" label-width="96px">
          <el-form-item label="版本号" prop="version">
            <el-input v-model="versionForm.version" placeholder="例如: v1.0.0" />
          </el-form-item>
          <el-form-item label="版本状态">
            <el-select v-model="versionForm.status" class="w-full">
              <el-option label="草稿" value="draft" />
              <el-option label="已测试" value="tested" />
              <el-option label="已发布" value="published" />
              <el-option label="已归档" value="archived" />
            </el-select>
          </el-form-item>
          <el-form-item label="选择策略">
            <el-select v-model="versionForm.selection_policy" class="w-full">
              <el-option label="稳定优先" value="stable" />
              <el-option label="高精度优先" value="high_quality" />
              <el-option label="低成本优先" value="low_cost" />
            </el-select>
          </el-form-item>

          <el-divider content-position="left">默认实现</el-divider>

          <el-form-item label="实现名称">
            <el-input v-model="versionForm.implName" placeholder="例如: stable-op" />
          </el-form-item>
          <el-form-item label="实现类型">
            <el-select v-model="versionForm.implType" class="w-full">
              <el-option label="operator_version" value="operator_version" />
              <el-option label="mcp_tool" value="mcp_tool" />
              <el-option label="ai_chain" value="ai_chain" />
            </el-select>
          </el-form-item>
          <el-form-item label="绑定引用" prop="bindingRef">
            <el-input v-model="versionForm.bindingRef" placeholder="operator_version UUID 或工具引用" />
          </el-form-item>
          <el-form-item label="Tier">
            <el-input v-model="versionForm.tier" placeholder="stable / high_quality / low_cost" />
          </el-form-item>
          <el-form-item label="延迟(ms)">
            <el-input-number v-model="versionForm.latencyMS" :min="0" :step="10" class="w-full" />
          </el-form-item>
          <el-form-item label="成本分">
            <el-input-number v-model="versionForm.costScore" :min="0" :step="0.1" :precision="2" class="w-full" />
          </el-form-item>
          <el-form-item label="质量分">
            <el-input-number v-model="versionForm.qualityScore" :min="0" :step="0.1" :precision="2" class="w-full" />
          </el-form-item>
          <el-form-item>
            <GvButton :loading="versionSubmitting" @click="handleCreateVersion">创建版本</GvButton>
          </el-form-item>
        </el-form>
      </div>
    </GvModal>
  </GvContainer>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  algorithmApi,
  type Algorithm,
  type AlgorithmCreateReq,
  type AlgorithmStatus,
  type CreateAlgorithmVersionReq,
} from '@/api/algorithm'
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
import { EmptyState, ErrorState } from '@/components/common'
import type { TableColumn } from '@/components/base/GvTable/types'
import type { FilterField } from '@/components/business/FilterBar/types'

const creating = ref(false)
const editing = ref(false)
const versionSubmitting = ref(false)
const publishingVersionID = ref('')

const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const showViewDialog = ref(false)
const showVersionDialog = ref(false)

const createFormRef = ref<FormInstance>()
const editFormRef = ref<FormInstance>()
const versionFormRef = ref<FormInstance>()

const currentAlgorithm = ref<Algorithm | null>(null)
const editAlgorithmID = ref('')
const searchKeyword = ref('')

const filters = ref({
  status: '',
  scenario: '',
})

const algorithmForm = reactive({
  code: '',
  name: '',
  scenario: '',
  description: '',
  status: 'draft' as AlgorithmStatus,
  tagsText: '',
})

const versionForm = reactive({
  version: '',
  status: 'draft' as 'draft' | 'tested' | 'published' | 'archived',
  selection_policy: 'stable' as 'stable' | 'high_quality' | 'low_cost',
  implName: 'default-impl',
  implType: 'operator_version' as 'operator_version' | 'mcp_tool' | 'ai_chain',
  bindingRef: '',
  tier: 'stable',
  latencyMS: 0,
  costScore: 0,
  qualityScore: 0,
})

const algorithmFormRules: FormRules = {
  code: [{ required: true, message: '请输入算法编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入算法名称', trigger: 'blur' }],
}

const versionFormRules: FormRules = {
  version: [{ required: true, message: '请输入版本号', trigger: 'blur' }],
  bindingRef: [{ required: true, message: '请输入绑定引用', trigger: 'blur' }],
}

const filterParams = computed(() => ({
  keyword: searchKeyword.value || undefined,
  status: filters.value.status || undefined,
  scenario: filters.value.scenario || undefined,
}))

const {
  items: algorithms,
  isLoading: loading,
  error,
  pagination,
  goToPage,
  changePageSize,
  refreshTable,
} = useTable(
  async (params) => {
    const res = await algorithmApi.list(params)
    return { items: res.data?.items || [], total: res.data?.total || 0 }
  },
  {
    immediate: true,
    initialPageSize: 20,
    extraParams: filterParams,
  }
)

const columns: TableColumn[] = [
  { prop: 'name', label: '名称', minWidth: '160', showOverflowTooltip: true },
  { prop: 'code', label: '编码', minWidth: '150' },
  { prop: 'scenario', label: '场景', minWidth: '120' },
  { prop: 'tags', label: '标签', minWidth: '170' },
  { prop: 'status', label: '状态', width: '120' },
  { prop: 'updated_at', label: '更新时间', width: '180' },
  { prop: 'actions', label: '操作', width: '300', fixed: 'right' },
]

const filterFields: FilterField[] = [
  {
    key: 'status',
    label: '状态',
    type: 'select',
    placeholder: '选择状态',
    options: [
      { label: '草稿', value: 'draft' },
      { label: '已发布', value: 'published' },
      { label: '已弃用', value: 'deprecated' },
    ],
  },
  {
    key: 'scenario',
    label: '场景',
    type: 'input',
    placeholder: '按场景筛选',
  },
]

const paginationConfig = computed(() => ({
  currentPage: pagination.page,
  pageSize: pagination.pageSize,
  total: pagination.total,
}))

function resetAlgorithmForm() {
  Object.assign(algorithmForm, {
    code: '',
    name: '',
    scenario: '',
    description: '',
    status: 'draft',
    tagsText: '',
  })
}

function resetVersionForm() {
  Object.assign(versionForm, {
    version: '',
    status: 'draft',
    selection_policy: 'stable',
    implName: 'default-impl',
    implType: 'operator_version',
    bindingRef: '',
    tier: 'stable',
    latencyMS: 0,
    costScore: 0,
    qualityScore: 0,
  })
}

function toTags(text: string): string[] {
  return text
    .split(',')
    .map((item) => item.trim())
    .filter(Boolean)
}

function toAlgorithmPayload(): AlgorithmCreateReq {
  return {
    code: algorithmForm.code.trim(),
    name: algorithmForm.name.trim(),
    description: algorithmForm.description.trim() || undefined,
    scenario: algorithmForm.scenario.trim() || undefined,
    status: algorithmForm.status,
    tags: toTags(algorithmForm.tagsText),
  }
}

function openCreateDialog() {
  resetAlgorithmForm()
  showCreateDialog.value = true
}

async function handleCreateConfirm() {
  if (!createFormRef.value) return
  const valid = await createFormRef.value.validate().catch(() => false)
  if (!valid) return

  creating.value = true
  try {
    await algorithmApi.create(toAlgorithmPayload())
    ElMessage.success('算法创建成功')
    showCreateDialog.value = false
    refreshTable()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '算法创建失败')
  } finally {
    creating.value = false
  }
}

async function loadAlgorithmDetail(id: string): Promise<Algorithm> {
  const res = await algorithmApi.get(id)
  return res.data
}

async function handleView(row: Algorithm) {
  try {
    currentAlgorithm.value = await loadAlgorithmDetail(row.id)
    showViewDialog.value = true
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '加载算法详情失败')
  }
}

function handleEdit(row: Algorithm) {
  editAlgorithmID.value = row.id
  Object.assign(algorithmForm, {
    code: row.code,
    name: row.name,
    scenario: row.scenario || '',
    description: row.description || '',
    status: row.status,
    tagsText: (row.tags || []).join(', '),
  })
  showEditDialog.value = true
}

async function handleEditConfirm() {
  if (!editAlgorithmID.value || !editFormRef.value) return
  const valid = await editFormRef.value.validate().catch(() => false)
  if (!valid) return

  editing.value = true
  try {
    await algorithmApi.update(editAlgorithmID.value, {
      name: algorithmForm.name.trim(),
      description: algorithmForm.description.trim() || undefined,
      scenario: algorithmForm.scenario.trim() || undefined,
      status: algorithmForm.status,
      tags: toTags(algorithmForm.tagsText),
    })
    ElMessage.success('算法更新成功')
    showEditDialog.value = false
    refreshTable()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '算法更新失败')
  } finally {
    editing.value = false
  }
}

async function handleDelete(row: Algorithm) {
  try {
    await ElMessageBox.confirm('确定删除该算法吗？该操作不可撤销。', '删除确认', { type: 'warning' })
    await algorithmApi.delete(row.id)
    ElMessage.success('算法已删除')
    refreshTable()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

async function handleManageVersions(row: Algorithm) {
  try {
    currentAlgorithm.value = await loadAlgorithmDetail(row.id)
    resetVersionForm()
    showVersionDialog.value = true
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '加载版本数据失败')
  }
}

async function handleCreateVersion() {
  if (!currentAlgorithm.value || !versionFormRef.value) return
  const valid = await versionFormRef.value.validate().catch(() => false)
  if (!valid) return

  const payload: CreateAlgorithmVersionReq = {
    version: versionForm.version.trim(),
    status: versionForm.status,
    selection_policy: versionForm.selection_policy,
    implementations: [
      {
        name: versionForm.implName.trim() || 'default-impl',
        type: versionForm.implType,
        binding_ref: versionForm.bindingRef.trim(),
        tier: versionForm.tier.trim() || 'stable',
        latency_ms: Number(versionForm.latencyMS || 0),
        cost_score: Number(versionForm.costScore || 0),
        quality_score: Number(versionForm.qualityScore || 0),
        is_default: true,
      },
    ],
  }

  versionSubmitting.value = true
  try {
    await algorithmApi.createVersion(currentAlgorithm.value.id, payload)
    ElMessage.success('版本创建成功')
    currentAlgorithm.value = await loadAlgorithmDetail(currentAlgorithm.value.id)
    resetVersionForm()
    refreshTable()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '版本创建失败')
  } finally {
    versionSubmitting.value = false
  }
}

async function handlePublishVersion(versionID: string) {
  if (!currentAlgorithm.value) return

  publishingVersionID.value = versionID
  try {
    await algorithmApi.publishVersion(currentAlgorithm.value.id, versionID)
    ElMessage.success('版本发布成功')
    currentAlgorithm.value = await loadAlgorithmDetail(currentAlgorithm.value.id)
    refreshTable()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '版本发布失败')
  } finally {
    publishingVersionID.value = ''
  }
}

const handlePageChange = goToPage
const handleSizeChange = changePageSize

function handleResetFilter() {
  searchKeyword.value = ''
  filters.value.status = ''
  filters.value.scenario = ''
  pagination.page = 1
}

function mapStatus(status: string): string {
  const map: Record<string, string> = {
    draft: 'pending',
    published: 'active',
    deprecated: 'disabled',
  }
  return map[status] || 'neutral'
}

function mapVersionStatus(status: string): string {
  const map: Record<string, string> = {
    draft: 'pending',
    tested: 'processing',
    published: 'active',
    archived: 'disabled',
  }
  return map[status] || 'neutral'
}

function formatDate(value?: string) {
  if (!value) return '-'
  return new Date(value).toLocaleString()
}
</script>
