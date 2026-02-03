<template>
  <GvContainer max-width="full" class="h-full">
    <div class="flex h-full gap-4">
      <!-- 左侧：类型和标签筛选 -->
      <aside class="w-64 flex-shrink-0">
        <GvCard shadow="sm" padding="md" class="sticky top-4">
          <!-- 媒体类型筛选 -->
          <div class="mb-6">
            <h3 class="text-sm font-semibold text-text-primary mb-3">媒体类型</h3>
            <div class="space-y-2">
              <div
                v-for="type in mediaTypes"
                :key="type.value"
                :class="[
                  'flex items-center justify-between px-3 py-2 rounded-lg cursor-pointer transition-all',
                  selectedType === type.value
                    ? 'bg-primary-50 text-primary-600 font-medium'
                    : 'hover:bg-neutral-50 text-text-secondary'
                ]"
                @click="handleTypeChange(type.value)"
              >
                <div class="flex items-center gap-2">
                  <el-icon :size="16">
                    <component :is="type.icon" />
                  </el-icon>
                  <span class="text-sm">{{ type.label }}</span>
                </div>
                <GvBadge
                  v-if="type.count !== undefined"
                  :color="selectedType === type.value ? 'primary' : 'neutral'"
                  size="small"
                  variant="tonal"
                >
                  {{ type.count }}
                </GvBadge>
              </div>
            </div>
          </div>

          <!-- 标签筛选 -->
          <div>
            <div class="flex items-center justify-between mb-3">
              <h3 class="text-sm font-semibold text-text-primary">标签</h3>
              <GvButton variant="text" size="small" @click="loadTags">
                <template #icon>
                  <el-icon><Refresh /></el-icon>
                </template>
              </GvButton>
            </div>
            <div v-if="tagsLoading" class="text-center py-4">
              <GvLoading size="small" />
            </div>
            <div v-else-if="tags.length === 0" class="text-center py-4 text-text-tertiary text-sm">
              暂无标签
            </div>
            <div v-else class="space-y-2 max-h-64 overflow-y-auto">
              <div
                v-for="tag in tags"
                :key="tag"
                :class="[
                  'px-3 py-2 rounded-lg cursor-pointer transition-all text-sm',
                  selectedTag === tag
                    ? 'bg-primary-50 text-primary-600 font-medium'
                    : 'hover:bg-neutral-50 text-text-secondary'
                ]"
                @click="handleTagChange(tag)"
              >
                {{ tag }}
              </div>
            </div>
          </div>
        </GvCard>
      </aside>

      <!-- 右侧：资产列表 -->
      <main class="flex-1 min-w-0">
        <!-- 页面头部 -->
        <PageHeader
          title="媒体资产库"
          description="管理所有媒体资产，支持视频、图片、音频、流媒体等多种格式"
        >
          <template #actions>
            <GvSpace>
              <SearchBar
                v-model="searchName"
                placeholder="搜索资产名称"
                class="w-80"
                immediate
                :show-button="false"
                @search="loadAssets"
              />
              <div class="view-switch-group">
                <button
                  :class="['view-switch-btn', { active: viewMode === 'grid' }]"
                  @click="viewMode = 'grid'"
                  title="网格视图"
                >
                  <el-icon :size="18"><Grid /></el-icon>
                </button>
                <button
                  :class="['view-switch-btn', { active: viewMode === 'list' }]"
                  @click="viewMode = 'list'"
                  title="列表视图"
                >
                  <el-icon :size="18"><List /></el-icon>
                </button>
              </div>
              <GvButton @click="showUploadDialog = true">
                <template #icon>
                  <el-icon><Upload /></el-icon>
                </template>
                添加资产
              </GvButton>
            </GvSpace>
          </template>
        </PageHeader>

        <!-- 资产展示 -->
        <div v-if="loading" class="flex justify-center items-center py-20">
          <GvLoading />
        </div>
        <div v-else-if="assets.length === 0" class="text-center py-20">
          <el-icon :size="64" class="text-neutral-300 mb-4">
            <FolderOpened />
          </el-icon>
          <p class="text-text-tertiary">暂无资产</p>
        </div>
        <div v-else>
          <!-- 网格视图 -->
          <div v-if="viewMode === 'grid'" class="grid gap-4 mb-6" :class="gridClass">
            <AssetCard
              v-for="asset in assets"
              :key="asset.id"
              :asset="asset"
              @view="handleView"
              @edit="handleEdit"
              @delete="handleDelete"
            />
          </div>

          <!-- 列表视图 -->
          <GvTable
            v-else
            :data="assets"
            :columns="tableColumns"
            :loading="loading"
            class="mb-6"
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
            <template #tags="{ row }">
              <GvSpace v-if="row.tags && row.tags.length > 0" size="xs" wrap>
                <GvTag v-for="tag in row.tags.slice(0, 3)" :key="tag" size="small" color="primary" variant="tonal">
                  {{ tag }}
                </GvTag>
                <GvTag v-if="row.tags.length > 3" size="small" color="neutral" variant="tonal">
                  +{{ row.tags.length - 3 }}
                </GvTag>
              </GvSpace>
              <span v-else class="text-text-tertiary text-sm">-</span>
            </template>
            <template #created_at="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
            <template #actions="{ row }">
              <GvSpace size="xs">
                <GvButton variant="text" size="small" @click="handleView(row)">查看</GvButton>
                <GvButton variant="text" size="small" @click="handleEdit(row)">编辑</GvButton>
                <GvButton variant="text" size="small" color="error" @click="handleDelete(row)">删除</GvButton>
              </GvSpace>
            </template>
          </GvTable>

          <!-- 分页 -->
          <div class="flex justify-end">
            <el-pagination
              v-model:current-page="pagination.page"
              v-model:page-size="pagination.page_size"
              :page-sizes="[12, 24, 48, 96]"
              :total="pagination.total"
              layout="total, sizes, prev, pager, next, jumper"
              @current-change="handlePageChange"
              @size-change="handleSizeChange"
            />
          </div>
        </div>

        <!-- 添加资产对话框 -->
        <GvModal
          v-model="showUploadDialog"
          title="添加资产"
          size="large"
          :confirm-loading="uploading"
          @confirm="handleUpload"
          @cancel="showUploadDialog = false"
        >
          <el-tabs v-model="uploadType" class="mb-4">
            <el-tab-pane label="URL 地址" name="url" />
            <el-tab-pane label="文件上传" name="file" />
          </el-tabs>

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

            <!-- URL 模式 -->
            <template v-if="uploadType === 'url'">
              <el-form-item label="资源地址" prop="path">
                <GvInput v-model="uploadForm.path" placeholder="请输入资源 URL" />
              </el-form-item>
            </template>

            <!-- 文件上传模式 -->
            <template v-else>
              <el-form-item label="选择文件" prop="file">
                <el-upload
                  ref="uploadRef"
                  :auto-upload="false"
                  :limit="1"
                  :on-change="handleFileChange"
                  :on-remove="handleFileRemove"
                >
                  <template #trigger>
                    <GvButton variant="outlined">选择文件</GvButton>
                  </template>
                  <template #tip>
                    <div class="text-sm text-text-tertiary mt-2">
                      支持视频、图片、音频文件
                    </div>
                  </template>
                </el-upload>
              </el-form-item>
            </template>

            <el-form-item label="标签" prop="tags">
              <el-select
                v-model="uploadForm.tags"
                multiple
                filterable
                allow-create
                placeholder="输入标签并回车"
                class="w-full"
              >
                <el-option
                  v-for="tag in tags"
                  :key="tag"
                  :label="tag"
                  :value="tag"
                />
              </el-select>
            </el-form-item>
          </el-form>
        </GvModal>

        <!-- 编辑资产对话框 -->
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
            <el-form-item label="标签" prop="tags">
              <el-select
                v-model="editForm.tags"
                multiple
                filterable
                allow-create
                placeholder="输入标签并回车"
                class="w-full"
              >
                <el-option
                  v-for="tag in tags"
                  :key="tag"
                  :label="tag"
                  :value="tag"
                />
              </el-select>
            </el-form-item>
          </el-form>
        </GvModal>

        <!-- 资产详情对话框 -->
        <GvModal
          v-model="showViewDialog"
          title="资产详情"
          size="large"
          :show-confirm="false"
          cancel-text="关闭"
        >
          <el-descriptions v-if="currentAsset" :column="2" border>
            <el-descriptions-item label="ID" :span="2">{{ currentAsset.id }}</el-descriptions-item>
            <el-descriptions-item label="名称" :span="2">{{ currentAsset.name }}</el-descriptions-item>
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
            <el-descriptions-item label="路径" :span="2">{{ currentAsset.path }}</el-descriptions-item>
            <el-descriptions-item label="格式">{{ currentAsset.format || '-' }}</el-descriptions-item>
            <el-descriptions-item label="大小">{{ formatSize(currentAsset.size) }}</el-descriptions-item>
            <el-descriptions-item label="时长">{{ currentAsset.duration ? formatDuration(currentAsset.duration) : '-' }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <StatusBadge :status="mapStatus(currentAsset.status)" />
            </el-descriptions-item>
            <el-descriptions-item label="标签" :span="2">
              <GvSpace v-if="currentAsset.tags && currentAsset.tags.length > 0" size="xs" wrap>
                <GvTag v-for="tag in currentAsset.tags" :key="tag" size="small" color="primary" variant="tonal">
                  {{ tag }}
                </GvTag>
              </GvSpace>
              <span v-else class="text-text-tertiary">-</span>
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatDate(currentAsset.created_at) }}</el-descriptions-item>
            <el-descriptions-item label="更新时间">{{ formatDate(currentAsset.updated_at) }}</el-descriptions-item>
          </el-descriptions>
        </GvModal>
      </main>
    </div>
  </GvContainer>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadFile } from 'element-plus'
import { Upload, VideoCamera, Picture, Headset, Connection, Refresh, FolderOpened, Grid, List } from '@element-plus/icons-vue'
import { assetApi, type MediaAsset, type AssetCreateReq, type AssetUpdateReq } from '@/api/asset'
import {
  GvContainer,
  GvCard,
  GvModal,
  GvButton,
  GvSpace,
  GvTag,
  GvBadge,
  GvInput,
  GvSelect,
  GvLoading,
  GvTable,
  PageHeader,
  SearchBar,
  StatusBadge,
  AssetCard
} from '@/components'

const loading = ref(false)
const uploading = ref(false)
const tagsLoading = ref(false)
const assets = ref<MediaAsset[]>([])
const tags = ref<string[]>([])
const showUploadDialog = ref(false)
const showEditDialog = ref(false)
const showViewDialog = ref(false)
const currentAsset = ref<MediaAsset | null>(null)
const uploadFormRef = ref<FormInstance>()
const editFormRef = ref<FormInstance>()
const uploadRef = ref()
const uploadType = ref<'url' | 'file'>('url')
const selectedFile = ref<UploadFile | null>(null)
const viewMode = ref<'grid' | 'list'>('grid')

const searchName = ref('')
const selectedType = ref<string | null>(null)
const selectedTag = ref<string | null>(null)

const pagination = reactive({
  page: 1,
  page_size: 12,
  total: 0
})

const uploadForm = reactive<AssetCreateReq>({
  type: 'video',
  source_type: 'upload',
  name: '',
  path: '',
  size: 0,
  format: '',
  tags: []
})

const editForm = reactive<AssetUpdateReq>({
  name: '',
  status: 'ready',
  tags: []
})

const uploadRules: FormRules = {
  name: [{ required: true, message: '请输入资产名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择资产类型', trigger: 'change' }],
  path: [{ required: uploadType.value === 'url', message: '请输入资源地址', trigger: 'blur' }]
}

const editRules: FormRules = {
  name: [{ required: true, message: '请输入资产名称', trigger: 'blur' }]
}

const mediaTypes = computed(() => [
  { label: '全部', value: null, icon: FolderOpened },
  { label: '视频', value: 'video', icon: VideoCamera },
  { label: '图片', value: 'image', icon: Picture },
  { label: '音频', value: 'audio', icon: Headset },
  { label: '流媒体', value: 'stream', icon: Connection }
])

const typeOptions = [
  { label: '视频', value: 'video' },
  { label: '图片', value: 'image' },
  { label: '音频', value: 'audio' },
  { label: '流媒体', value: 'stream' }
]

const statusOptions = [
  { label: '就绪', value: 'ready' },
  { label: '处理中', value: 'processing' },
  { label: '待处理', value: 'pending' },
  { label: '错误', value: 'error' }
]

const tableColumns = [
  { prop: 'name', label: '名称', minWidth: 200 },
  { prop: 'type', label: '类型', width: 100 },
  { prop: 'source_type', label: '来源', width: 120 },
  { prop: 'format', label: '格式', width: 80 },
  { prop: 'size', label: '大小', width: 100 },
  { prop: 'duration', label: '时长', width: 100 },
  { prop: 'status', label: '状态', width: 100 },
  { prop: 'tags', label: '标签', width: 200 },
  { prop: 'created_at', label: '创建时间', width: 180 },
  { prop: 'actions', label: '操作', width: 200, fixed: 'right' }
]

// 响应式网格类名
const gridClass = computed(() => {
  return 'grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6'
})

onMounted(() => {
  loadAssets()
  loadTags()
})

async function loadAssets() {
  loading.value = true
  try {
    const response = await assetApi.list({
      name: searchName.value || undefined,
      type: selectedType.value as any,
      tags: selectedTag.value || undefined,
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

async function loadTags() {
  tagsLoading.value = true
  try {
    const response = await assetApi.getTags()
    tags.value = response.data.tags || []
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '加载标签失败')
  } finally {
    tagsLoading.value = false
  }
}

function handleTypeChange(type: string | null) {
  selectedType.value = type
  pagination.page = 1
  loadAssets()
}

function handleTagChange(tag: string) {
  selectedTag.value = selectedTag.value === tag ? null : tag
  pagination.page = 1
  loadAssets()
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

function handleFileChange(file: UploadFile) {
  selectedFile.value = file
  uploadForm.name = file.name.split('.')[0]
  uploadForm.size = file.size || 0
  uploadForm.format = file.name.split('.').pop() || ''
  uploadForm.path = file.name
}

function handleFileRemove() {
  selectedFile.value = null
  uploadForm.path = ''
  uploadForm.size = 0
  uploadForm.format = ''
}

async function handleUpload() {
  if (!uploadFormRef.value) return

  await uploadFormRef.value.validate(async (valid) => {
    if (!valid) return

    if (uploadType.value === 'file' && !selectedFile.value) {
      ElMessage.warning('请选择文件')
      return
    }

    uploading.value = true
    try {
      if (uploadType.value === 'file' && selectedFile.value?.raw) {
        // 文件上传模式
        await assetApi.upload(
          selectedFile.value.raw,
          uploadForm.type,
          uploadForm.name,
          uploadForm.tags
        )
      } else {
        // URL 模式
        await assetApi.create(uploadForm)
      }
      ElMessage.success('添加成功')
      showUploadDialog.value = false
      resetUploadForm()
      loadAssets()
      loadTags()
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || '添加失败')
    } finally {
      uploading.value = false
    }
  })
}

function resetUploadForm() {
  uploadForm.type = 'video'
  uploadForm.source_type = 'upload'
  uploadForm.name = ''
  uploadForm.path = ''
  uploadForm.size = 0
  uploadForm.format = ''
  uploadForm.tags = []
  selectedFile.value = null
  uploadFormRef.value?.resetFields()
}

function handleView(asset: MediaAsset) {
  currentAsset.value = asset
  showViewDialog.value = true
}

function handleEdit(asset: MediaAsset) {
  currentAsset.value = asset
  editForm.name = asset.name
  editForm.status = asset.status
  editForm.tags = asset.tags || []
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
      loadTags()
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || '更新失败')
    }
  })
}

async function handleDelete(asset: MediaAsset) {
  try {
    await ElMessageBox.confirm('确定要删除此资产吗？', '提示', {
      type: 'warning'
    })
    await assetApi.delete(asset.id)
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
    audio: '音频',
    stream: '流媒体'
  }
  return map[type] || type
}

function getTypeColor(type: string) {
  const map: Record<string, string> = {
    video: 'primary',
    image: 'success',
    audio: 'warning',
    stream: 'info'
  }
  return map[type] || 'neutral'
}

function getSourceTypeLabel(type: string) {
  const map: Record<string, string> = {
    upload: '上传',
    live: '直播',
    vod: '点播',
    generated: '生成',
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
/* 视图切换按钮组 */
.view-switch-group {
  display: inline-flex;
  background: #f5f7fa;
  border-radius: 8px;
  padding: 4px;
  gap: 4px;
}

.view-switch-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: #606266;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.view-switch-btn::before {
  content: '';
  position: absolute;
  inset: 0;
  background: currentColor;
  opacity: 0;
  transition: opacity 0.2s;
}

.view-switch-btn:hover::before {
  opacity: 0.08;
}

.view-switch-btn:active {
  transform: scale(0.95);
}

.view-switch-btn.active {
  background: #ffffff;
  color: #409eff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
}

.view-switch-btn.active:hover::before {
  opacity: 0;
}

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
