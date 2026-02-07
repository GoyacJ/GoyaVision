<template>
  <GvContainer max-width="full">
    <PageHeader
      title="AI 模型管理"
      description="管理大模型配置，支持本地与远程模型接入"
    >
      <template #actions>
        <GvSpace>
          <GvSelect
            v-model="filterProvider"
            placeholder="提供商"
            :options="[{ label: '全部提供商', value: '' }, ...providerOptions]"
            class="w-44"
            @change="refreshTable"
          />
          <GvSelect
            v-model="filterStatus"
            placeholder="状态"
            :options="[{ label: '全部状态', value: '' }, ...statusOptions]"
            class="w-32"
            @change="refreshTable"
          />
          <SearchBar
            v-model="searchKeyword"
            placeholder="搜索模型"
            class="w-80"
            immediate
            :show-button="false"
            @search="refreshTable"
          />
          <GvButton @click="openCreateDialog">
            <template #icon>
              <el-icon><Plus /></el-icon>
            </template>
            添加模型
          </GvButton>
        </GvSpace>
      </template>
    </PageHeader>

    <GvTable
      :data="aiModels"
      :columns="columns"
      :loading="loading"
      border
      stripe
      pagination
      :pagination-config="paginationConfig"
      @current-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <template #provider="{ row }">
        <GvTag :color="getProviderColor(row.provider)" size="small">
          {{ row.provider }}
        </GvTag>
      </template>

      <template #has_api_key="{ row }">
        <GvTag :color="row.has_api_key ? 'success' : 'neutral'" size="small">
          {{ row.has_api_key ? '已配置' : '未配置' }}
        </GvTag>
      </template>

      <template #status="{ row }">
        <StatusBadge :status="row.status === 'active' ? 'active' : 'disabled'" />
      </template>

      <template #created_at="{ row }">
        {{ formatDate(row.created_at) }}
      </template>

      <template #actions="{ row }">
        <GvSpace size="xs">
          <GvButton
            size="small"
            variant="text"
            :loading="testingId === row.id"
            @click="handleTestConnection(row)"
          >
            测试连接
          </GvButton>
          <GvButton size="small" variant="text" @click="handleEdit(row)">
            编辑
          </GvButton>
          <GvButton size="small" variant="text" color="error" @click="handleDelete(row)">
            删除
          </GvButton>
        </GvSpace>
      </template>
    </GvTable>

    <!-- 编辑/创建对话框 -->
    <GvModal
      v-model="showDialog"
      :title="isEdit ? '编辑模型' : '添加模型'"
      size="medium"
      confirm-text="保存"
      :confirm-loading="submitting"
      @confirm="handleSubmit"
      @cancel="showDialog = false"
    >
      <div class="space-y-4">
        <GvInput v-model="form.name" label="名称" placeholder="如 GPT-4, Llama3" />
        <GvInput v-model="form.description" label="描述" type="textarea" :rows="2" placeholder="模型用途说明 (可选)" />
        <GvSelect
          v-model="form.provider"
          label="提供商"
          :options="providerOptions"
          class="w-full"
        />
        <GvInput v-model="form.endpoint" label="Endpoint" placeholder="API 地址 (可选)" />
        <div class="space-y-1.5">
          <label class="block text-sm font-medium text-text-primary">API Key</label>
          <div class="flex items-center gap-2">
            <GvInput
              v-model="form.api_key"
              type="password"
              show-password
              placeholder="密钥 (可选)"
              class="flex-1"
            />
            <GvTag v-if="isEdit" :color="editHasApiKey ? 'success' : 'neutral'" size="small">
              {{ editHasApiKey ? '已配置' : '未配置' }}
            </GvTag>
          </div>
          <p v-if="isEdit" class="text-xs text-text-tertiary">留空则不更新密钥</p>
        </div>
        <GvInput v-model="form.model_name" label="模型标识" placeholder="如 gpt-4-turbo" />

        <div class="space-y-1.5">
          <label class="block text-sm font-medium text-text-primary">额外配置 (JSON)</label>
          <SchemaEditor
            v-model="form.config"
            :rows="5"
          />
        </div>

        <GvSelect
          v-model="form.visibility"
          label="可见范围"
          :options="VISIBILITY_OPTIONS"
          class="w-full"
        />

        <GvSelect
          v-if="isEdit"
          v-model="form.status"
          label="状态"
          :options="statusOptions"
          class="w-full"
        />

      </div>
    </GvModal>
  </GvContainer>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { aiModelApi, type AIModel } from '@/api/ai-model'
import { useTable } from '@/composables'
import GvContainer from '@/components/layout/GvContainer/index.vue'
import GvTable from '@/components/base/GvTable/index.vue'
import GvModal from '@/components/base/GvModal/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvInput from '@/components/base/GvInput/index.vue'
import GvSelect from '@/components/base/GvSelect/index.vue'
import GvSpace from '@/components/layout/GvSpace/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'
import PageHeader from '@/components/business/PageHeader/index.vue'
import SearchBar from '@/components/business/SearchBar/index.vue'
import StatusBadge from '@/components/business/StatusBadge/index.vue'
import SchemaEditor from '@/views/operator/components/SchemaEditor.vue'
import type { TableColumn } from '@/components/base/GvTable/types'
import { VISIBILITY_OPTIONS } from '@/constants/visibility'

const searchKeyword = ref('')
const filterProvider = ref('')
const filterStatus = ref('')
const showDialog = ref(false)
const submitting = ref(false)
const isEdit = ref(false)
const editId = ref('')
const editHasApiKey = ref(false)
const testingId = ref('')


const form = reactive<any>({
  name: '',
  description: '',
  provider: 'openai',
  endpoint: '',
  api_key: '',
  model_name: '',
  config: {},
  status: 'active',
  visibility: 0
})

const {
  items: aiModels,
  isLoading: loading,
  pagination,
  goToPage,
  changePageSize,
  refreshTable
} = useTable(
  async (params) => {
    const res = await aiModelApi.list({
      ...params,
      keyword: searchKeyword.value || undefined,
      provider: filterProvider.value || undefined,
      status: filterStatus.value || undefined
    })
    return { items: res.data?.items ?? [], total: res.data?.total ?? 0 }
  },
  {
    immediate: true,
    initialPageSize: 20
  }
)

const providerOptions = [
  { label: 'OpenAI', value: 'openai' },
  { label: 'Anthropic', value: 'anthropic' },
  { label: '千问 (Qwen)', value: 'qwen' },
  { label: '豆包 (Doubao)', value: 'doubao' },
  { label: '智谱 (Zhipu)', value: 'zhipu' },
  { label: 'Ollama', value: 'ollama' },
  { label: 'vLLM', value: 'vllm' },
  { label: 'Local', value: 'local' },
  { label: 'Custom', value: 'custom' }
]

const statusOptions = [
  { label: '启用', value: 'active' },
  { label: '禁用', value: 'disabled' }
]

const columns: TableColumn[] = [
  { prop: 'name', label: '名称', minWidth: '140' },
  { prop: 'description', label: '描述', minWidth: '160', showOverflowTooltip: true },
  { prop: 'provider', label: '提供商', width: '120' },
  { prop: 'model_name', label: '模型标识', minWidth: '140' },
  { prop: 'endpoint', label: 'Endpoint', minWidth: '180', showOverflowTooltip: true },
  { prop: 'has_api_key', label: 'API Key', width: '100' },
  { prop: 'status', label: '状态', width: '100' },
  { prop: 'created_at', label: '创建时间', width: '160' },
  { prop: 'actions', label: '操作', width: '200', fixed: 'right' }
]

const paginationConfig = computed(() => ({
  currentPage: pagination.page,
  pageSize: pagination.pageSize,
  total: pagination.total
}))

const handlePageChange = goToPage
const handleSizeChange = changePageSize

function getProviderColor(provider: string): any {
  const map: Record<string, string> = {
    openai: 'success',
    anthropic: 'warning',
    qwen: 'primary',
    doubao: 'info',
    zhipu: 'success',
    ollama: 'primary',
    vllm: 'info',
    local: 'info',
    custom: 'neutral'
  }
  return map[provider] || 'neutral'
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

function resetForm() {
  form.name = ''
  form.description = ''
  form.provider = 'openai'
  form.endpoint = ''
  form.api_key = ''
  form.model_name = ''
  form.config = {}
  form.status = 'active'
  form.visibility = 0
}

function openCreateDialog() {
  isEdit.value = false
  editId.value = ''
  editHasApiKey.value = false
  resetForm()
  showDialog.value = true
}

function handleEdit(row: AIModel) {
  isEdit.value = true
  editId.value = row.id
  editHasApiKey.value = row.has_api_key
  form.name = row.name
  form.description = row.description || ''
  form.provider = row.provider
  form.endpoint = row.endpoint
  form.api_key = ''
  form.model_name = row.model_name
  form.config = row.config || {}
  form.status = row.status
  form.visibility = row.visibility ?? 0
  showDialog.value = true
}

async function handleTestConnection(row: AIModel) {
  testingId.value = row.id
  try {
    const res = await aiModelApi.testConnection(row.id)
    if (res.data?.success) {
      ElMessage.success(res.data.message || '连接成功')
    } else {
      ElMessage.error(res.data?.message || '连接失败')
    }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '测试连接失败')
  } finally {
    testingId.value = ''
  }
}

async function handleSubmit() {
  if (!form.name || !form.provider) {
    ElMessage.warning('请填写名称和提供商')
    return
  }

  submitting.value = true
  try {
    const data: any = {
      name: form.name,
      description: form.description,
      provider: form.provider,
      endpoint: form.endpoint,
      model_name: form.model_name,
      config: form.config,
      visibility: form.visibility ?? 0
    }
    if (form.api_key) {
      data.api_key = form.api_key
    }

    if (isEdit.value) {
      data.status = form.status
      await aiModelApi.update(editId.value, data)
      ElMessage.success('更新成功')
    } else {
      await aiModelApi.create(data)
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    refreshTable()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(row: AIModel) {
  try {
    await ElMessageBox.confirm('确定要删除此模型吗？', '提示', {
      type: 'warning'
    })
    await aiModelApi.delete(row.id)
    ElMessage.success('删除成功')
    refreshTable()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}
</script>
