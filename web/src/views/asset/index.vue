<template>
  <GvContainer max-width="full">
    <!-- 页面头部 -->
    <PageHeader
      title="媒体资产库"
      description="管理所有媒体资产，支持视频、图片、音频等多种格式"
    >
      <template #actions>
        <GvSpace>
          <SearchBar
            v-model="searchName"
            placeholder="搜索资产名称"
            class="w-80"
            immediate
            @search="loadAssets"
          />
          <GvButton @click="showUploadDialog = true">
            <template #icon>
              <el-icon><Upload /></el-icon>
            </template>
            上传资产
          </GvButton>
        </GvSpace>
      </template>
    </PageHeader>

    <!-- 筛选栏 -->
    <FilterBar
      v-model="filters"
      :fields="filterFields"
      :loading="loading"
      @filter="loadAssets"
      @reset="handleResetFilter"
    />

    <!-- 数据表格 -->
    <GvTable
      :data="assets"
      :columns="columns"
      :loading="loading"
      border
      stripe
      pagination
      :pagination-config="paginationConfig"
      @current-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <template #type="{ row }">
        <GvTag :color="getTypeColor(row.type)" size="small">
          {{ getTypeLabel(row.type) }}
        </GvTag>
      </template>
      
      <template #source_type="{ row }">
        <GvTag color="info" size="small" variant="tonal">
          {{ getSourceTypeLabel(row.source_type) }}
        </GvTag>
      </template>
      
      <template #size="{ row }">
        {{ formatSize(row.size) }}
      </template>
      
      <template #duration="{ row }">
        {{ row.duration ? formatDuration(row.duration) : '-' }}
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
          <GvButton size="small" variant="text" @click="handleDelete(row)">
            删除
          </GvButton>
        </GvSpace>
      </template>
    </GvTable>

    <!-- 上传对话框 -->
    <GvModal
      v-model="showUploadDialog"
      title="上传资产"
      :confirm-loading="uploading"
      @confirm="handleUpload"
      @cancel="showUploadDialog = false"
    >
      <el-form ref="uploadFormRef" :model="uploadForm" :rules="uploadRules" label-width="100px">
        <el-form-item label="资产名称" prop="name">
          <GvInput v-model="uploadForm.name" placeholder="请输入资产名称" />
        </el-form-item>
        <el-form-item label="资产类型" prop="type">
          <GvSelect
            v-model="uploadForm.type"
            :options="typeOptions"
            placeholder="请选择类型"
          />
        </el-form-item>
        <el-form-item label="文件路径" prop="path">
          <GvInput v-model="uploadForm.path" placeholder="请输入文件路径（相对或绝对路径）" />
        </el-form-item>
        <el-form-item label="文件大小" prop="size">
          <GvInput v-model.number="uploadForm.size" type="number" placeholder="文件大小（字节）" />
        </el-form-item>
        <el-form-item label="格式" prop="format">
          <GvInput v-model="uploadForm.format" placeholder="如：mp4, jpg, mp3" />
        </el-form-item>
      </el-form>
    </GvModal>

    <!-- 编辑对话框 -->
    <GvModal
      v-model="showEditDialog"
      title="编辑资产"
      @confirm="handleUpdate"
      @cancel="showEditDialog = false"
    >
      <el-form ref="editFormRef" :model="editForm" :rules="editRules" label-width="100px">
        <el-form-item label="资产名称" prop="name">
          <GvInput v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <GvSelect
            v-model="editForm.status"
            :options="statusOptions"
          />
        </el-form-item>
      </el-form>
    </GvModal>

    <!-- 详情对话框 -->
    <GvModal
      v-model="showViewDialog"
      title="资产详情"
      size="large"
      :show-confirm="false"
      cancel-text="关闭"
    >
      <el-descriptions v-if="currentAsset" :column="1" border>
        <el-descriptions-item label="ID">{{ currentAsset.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ currentAsset.name }}</el-descriptions-item>
        <el-descriptions-item label="类型">
          <GvTag :color="getTypeColor(currentAsset.type)" size="small">
            {{ getTypeLabel(currentAsset.type) }}
          </GvTag>
        </el-descriptions-item>
        <el-descriptions-item label="来源">
          <GvTag color="info" size="small" variant="tonal">
            {{ getSourceTypeLabel(currentAsset.source_type) }}
          </GvTag>
        </el-descriptions-item>
        <el-descriptions-item label="路径">{{ currentAsset.path }}</el-descriptions-item>
        <el-descriptions-item label="格式">{{ currentAsset.format || '-' }}</el-descriptions-item>
        <el-descriptions-item label="大小">{{ formatSize(currentAsset.size) }}</el-descriptions-item>
        <el-descriptions-item label="时长">{{ currentAsset.duration ? formatDuration(currentAsset.duration) : '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <StatusBadge :status="mapStatus(currentAsset.status)" />
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentAsset.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentAsset.updated_at) }}</el-descriptions-item>
      </el-descriptions>
    </GvModal>
  </GvContainer>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import { assetApi, type MediaAsset, type AssetCreateReq, type AssetUpdateReq } from '@/api/asset'
import {
  GvContainer,
  GvTable,
  GvModal,
  GvButton,
  GvSpace,
  GvTag,
  GvInput,
  GvSelect,
  PageHeader,
  FilterBar,
  SearchBar,
  StatusBadge,
  type TableColumn,
  type FilterField
} from '@/components'

const loading = ref(false)
const uploading = ref(false)
const assets = ref<MediaAsset[]>([])
const showUploadDialog = ref(false)
const showEditDialog = ref(false)
const showViewDialog = ref(false)
const currentAsset = ref<MediaAsset | null>(null)
const uploadFormRef = ref<FormInstance>()
const editFormRef = ref<FormInstance>()

const searchName = ref('')

const filters = ref({
  name: '',
  type: '',
  source_type: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const uploadForm = reactive<AssetCreateReq>({
  type: 'video',
  source_type: 'upload',
  name: '',
  path: '',
  size: 0,
  format: ''
})

const editForm = reactive<AssetUpdateReq>({
  name: '',
  status: 'ready'
})

const uploadRules: FormRules = {
  name: [{ required: true, message: '请输入资产名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择资产类型', trigger: 'change' }],
  path: [{ required: true, message: '请输入文件路径', trigger: 'blur' }],
  size: [{ required: true, message: '请输入文件大小', trigger: 'blur' }]
}

const editRules: FormRules = {
  name: [{ required: true, message: '请输入资产名称', trigger: 'blur' }]
}

const typeOptions = [
  { label: '视频', value: 'video' },
  { label: '图片', value: 'image' },
  { label: '音频', value: 'audio' }
]

const sourceTypeOptions = [
  { label: '上传', value: 'upload' },
  { label: '流捕获', value: 'stream_capture' },
  { label: '算子输出', value: 'operator_output' }
]

const statusOptions = [
  { label: '就绪', value: 'ready' },
  { label: '处理中', value: 'processing' },
  { label: '待处理', value: 'pending' },
  { label: '错误', value: 'error' }
]

const filterFields: FilterField[] = [
  {
    key: 'name',
    label: '名称',
    type: 'input',
    placeholder: '搜索资产名称'
  },
  {
    key: 'type',
    label: '资产类型',
    type: 'select',
    placeholder: '选择资产类型',
    options: typeOptions
  },
  {
    key: 'source_type',
    label: '来源类型',
    type: 'select',
    placeholder: '选择来源类型',
    options: sourceTypeOptions
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
  { prop: 'name', label: '名称', minWidth: '180', showOverflowTooltip: true },
  { prop: 'type', label: '类型', width: '90' },
  { prop: 'source_type', label: '来源', width: '120' },
  { prop: 'format', label: '格式', width: '90' },
  { prop: 'size', label: '大小', width: '100' },
  { prop: 'duration', label: '时长', width: '100' },
  { prop: 'status', label: '状态', width: '120' },
  { prop: 'created_at', label: '创建时间', width: '160' },
  { prop: 'actions', label: '操作', width: '200', fixed: 'right' }
]

const paginationConfig = computed(() => ({
  currentPage: pagination.page,
  pageSize: pagination.page_size,
  total: pagination.total
}))

onMounted(() => {
  loadAssets()
})

async function loadAssets() {
  loading.value = true
  try {
    const response = await assetApi.list({
      name: filters.value.name || searchName.value || undefined,
      type: filters.value.type as any,
      source_type: filters.value.source_type as any,
      status: filters.value.status as any,
      page: pagination.page,
      page_size: pagination.page_size
    })
    assets.value = response.data.items
    pagination.total = response.data.total
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '加载失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  pagination.page = page
  loadAssets()
}

function handleSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  loadAssets()
}

function handleResetFilter() {
  searchName.value = ''
  loadAssets()
}

async function handleUpload() {
  if (!uploadFormRef.value) return
  await uploadFormRef.value.validate(async (valid) => {
    if (!valid) return
    uploading.value = true
    try {
      await assetApi.create(uploadForm)
      ElMessage.success('创建成功')
      showUploadDialog.value = false
      loadAssets()
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || '创建失败')
    } finally {
      uploading.value = false
    }
  })
}

function handleView(row: MediaAsset) {
  currentAsset.value = row
  showViewDialog.value = true
}

function handleEdit(row: MediaAsset) {
  currentAsset.value = row
  editForm.name = row.name
  editForm.status = row.status
  showEditDialog.value = true
}

async function handleUpdate() {
  if (!editFormRef.value || !currentAsset.value) return
  await editFormRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      await assetApi.update(currentAsset.value!.id, editForm)
      ElMessage.success('更新成功')
      showEditDialog.value = false
      loadAssets()
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || '更新失败')
    }
  })
}

async function handleDelete(row: MediaAsset) {
  try {
    await ElMessageBox.confirm('确定要删除此资产吗？', '提示', {
      type: 'warning'
    })
    await assetApi.delete(row.id)
    ElMessage.success('删除成功')
    loadAssets()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

function getTypeLabel(type: string) {
  const map: Record<string, string> = {
    video: '视频',
    image: '图片',
    audio: '音频'
  }
  return map[type] || type
}

function getTypeColor(type: string) {
  const map: Record<string, string> = {
    video: 'primary',
    image: 'success',
    audio: 'warning'
  }
  return map[type] || 'neutral'
}

function getSourceTypeLabel(type: string) {
  const map: Record<string, string> = {
    upload: '上传',
    stream_capture: '流捕获',
    operator_output: '算子输出'
  }
  return map[type] || type
}

function mapStatus(status: string): any {
  const map: Record<string, string> = {
    ready: 'success',
    processing: 'processing',
    pending: 'pending',
    error: 'error'
  }
  return map[status] || 'inactive'
}

function formatSize(size: number): string {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`
  if (size < 1024 * 1024 * 1024) return `${(size / (1024 * 1024)).toFixed(2)} MB`
  return `${(size / (1024 * 1024 * 1024)).toFixed(2)} GB`
}

function formatDuration(seconds: number): string {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = Math.floor(seconds % 60)
  if (h > 0) return `${h}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
  return `${m}:${s.toString().padStart(2, '0')}`
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
