<template>
  <GvContainer max-width="full" class="h-full">
    <div class="flex flex-col lg:flex-row h-full gap-4">
      <!-- 左侧：类型和标签筛选 -->
      <aside class="w-full lg:w-64 flex-shrink-0">
        <div class="mb-4">
          <h1 class="text-2xl font-bold text-text-primary">媒体资产库</h1>
        </div>

        <GvCard shadow="sm" padding="md" class="lg:sticky lg:top-4">
          <div class="mb-6">
            <h3 class="text-sm font-semibold text-text-primary mb-3">媒体类型</h3>
            <div class="space-y-2">
              <div
                v-for="type in mediaTypes"
                :key="String(type.value)"
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
              </div>
            </div>
          </div>

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
            <div v-else class="flex flex-wrap gap-2 max-h-64 overflow-y-auto">
              <div
                v-for="tag in tags"
                :key="tag"
                :class="[
                  'px-3 py-1.5 rounded-lg cursor-pointer transition-all text-sm whitespace-nowrap',
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

      <main class="flex-1 min-w-0">
        <!-- 操作栏 -->
        <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between mb-6 gap-4">
          <div class="flex items-center gap-3 w-full sm:w-auto">
            <SearchBar
              v-model="searchName"
              placeholder="搜索资产名称"
              class="w-full sm:w-80"
              immediate
              :show-button="false"
              @search="() => { pagination.page = 1 }"
            />
            <div class="view-switch-group flex-shrink-0">
              <button
                :class="['view-switch-btn', { active: viewMode === 'grid' }]"
                @click="viewMode = 'grid'"
                title="网格视图"
              >
                <el-icon :size="18"><Grid /></el-icon>
              </button>
              <button
                v-if="!isMobile"
                :class="['view-switch-btn', { active: viewMode === 'list' }]"
                @click="viewMode = 'list'"
                title="列表视图"
              >
                <el-icon :size="18"><List /></el-icon>
              </button>
            </div>
          </div>
          <GvButton @click="showUploadDialog = true" class="w-full sm:w-auto">
            <template #icon>
              <el-icon><Upload /></el-icon>
            </template>
            添加资产
          </GvButton>
        </div>

        <!-- 资产展示 -->
        <LoadingState v-if="loading" message="加载资产列表..." />

        <ErrorState
          v-else-if="error"
          :error="error"
          title="加载失败"
          @retry="refreshTable"
        />

        <EmptyState
          v-else-if="assets.length === 0"
          icon="🎬"
          title="还没有媒体资产"
          description="开始上传您的第一个视频、图片或音频文件"
          action-text="添加资产"
          show-action
          @action="showUploadDialog = true"
        />

        <div v-else>
          <!-- 网格视图 -->
          <div v-if="viewMode === 'grid'" class="grid gap-4 mb-6" :class="gridClass">
            <AssetCard
              v-for="asset in assets"
              :key="asset.id"
              :asset="asset"
              :can-edit="canEditPermission"
              @click="handleDetail"
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
              <GvTag :color="getTypeColor(row.type)" size="small" variant="tonal">
                <span class="inline-flex items-center gap-1">
                  <el-icon :size="14">
                    <component :is="getTypeIcon(row.type)" />
                  </el-icon>
                  {{ getTypeLabel(row.type) }}
                </span>
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
              <div class="flex flex-col gap-1">
                <StatusBadge :status="mapStatus(row.status)" />
                <GvTag v-if="row.visibility !== undefined" :color="row.visibility === 2 ? 'success' : (row.visibility === 1 ? 'warning' : 'neutral')" size="xs" variant="tonal">
                  {{ row.visibility === 2 ? '公开' : (row.visibility === 1 ? '角色' : '私有') }}
                </GvTag>
              </div>
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
                <GvButton variant="text" size="small" @click="handleDetail(row)">详情</GvButton>
                <GvButton v-if="canEditPermission" variant="text" size="small" color="error" @click="handleDelete(row)">删除</GvButton>
              </GvSpace>
            </template>
          </GvTable>

          <!-- 分页 -->
          <div class="flex justify-end">
            <el-pagination
              v-model:current-page="pagination.page"
              v-model:page-size="pagination.pageSize"
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
                <GvUpload
                  ref="uploadRef"
                  v-model="uploadFileList"
                  :auto-upload="false"
                  :limit="1"
                  button-text="选择文件"
                  tip="支持视频、图片、音频文件"
                  @change="handleFileChange"
                  @remove="handleFileRemove"
                />
              </el-form-item>
            </template>

            <el-form-item label="可见范围">
              <GvSelect
                v-model="uploadForm.visibility"
                :options="VISIBILITY_OPTIONS"
                placeholder="请选择可见范围"
              />
            </el-form-item>

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

        <!-- 资产详情抽屉（查看 + 编辑一体化） -->
        <GvDrawer
          v-model="showDetailDrawer"
          title="资产详情"
          direction="right"
          size="large"
          :show-footer="false"
        >
          <div v-if="currentAsset" class="asset-detail-panel">
            <div class="asset-detail-toolbar">
              <GvButton size="small" variant="tonal" @click="handleCopyLink">复制链接</GvButton>
              <GvButton size="small" variant="tonal" @click="handleDownload">下载</GvButton>
            </div>

            <div class="asset-preview-card">
              <!-- 视频预览 -->
              <div v-if="currentAsset.type === 'video'" class="preview-container">
                <video
                  :src="currentAsset.path"
                  controls
                  class="preview-media preview-media--zoomable"
                  @dblclick="openVideoPreview"
                >
                  您的浏览器不支持视频播放
                </video>
              </div>

              <!-- 图片预览 -->
              <div v-else-if="currentAsset.type === 'image'" class="preview-container">
                <img
                  :src="currentAsset.path"
                  :alt="currentAsset.name"
                  class="preview-media preview-media--zoomable"
                  @dblclick="openImagePreview"
                />
              </div>

              <!-- 音频预览 -->
              <div v-else-if="currentAsset.type === 'audio'" class="preview-container audio-preview">
                <div class="audio-icon">
                  <el-icon :size="80" class="text-primary-500">
                    <Headset />
                  </el-icon>
                </div>
                <audio
                  :src="currentAsset.path"
                  controls
                  class="audio-player"
                >
                  您的浏览器不支持音频播放
                </audio>
              </div>

              <!-- 未知类型 -->
              <div v-else class="preview-container">
                <div class="text-center text-text-tertiary">
                  <el-icon :size="80" class="mb-4">
                    <FolderOpened />
                  </el-icon>
                  <p>暂无预览</p>
                </div>
              </div>
            </div>

            <div class="asset-form-grid">
              <div class="info-item info-item--full">
                <span class="info-label">名称</span>
                <template v-if="canEditPermission">
                  <GvInput v-model="editForm.name" />
                </template>
                <span v-else class="info-value">{{ currentAsset.name }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">类型</span>
                <GvTag :color="getTypeColor(currentAsset.type)" size="small">
                  {{ getTypeLabel(currentAsset.type) }}
                </GvTag>
              </div>
              <div class="info-item">
                <span class="info-label">来源</span>
                <GvTag color="info" size="small" variant="tonal">
                  {{ getSourceTypeLabel(currentAsset.source_type) }}
                </GvTag>
              </div>
              <div class="info-item">
                <span class="info-label">格式</span>
                <span class="info-value">{{ currentAsset.format || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">大小</span>
                <span class="info-value">{{ formatSize(currentAsset.size) }}</span>
              </div>
              <div v-if="currentAsset.duration" class="info-item">
                <span class="info-label">时长</span>
                <span class="info-value">{{ formatDuration(currentAsset.duration) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">状态</span>
                <template v-if="canEditPermission">
                  <GvSelect v-model="editForm.status" :options="statusOptions" />
                </template>
                <StatusBadge v-else :status="mapStatus(currentAsset.status)" />
              </div>
              <div class="info-item">
                <span class="info-label">可见范围</span>
                <template v-if="canEditPermission">
                  <GvSelect
                    v-model="editForm.visibility"
                    :options="VISIBILITY_OPTIONS"
                    placeholder="请选择可见范围"
                  />
                </template>
                <span v-else class="info-value">{{ (currentAsset.visibility === 2 ? '公开' : (currentAsset.visibility === 1 ? '角色可见' : '私有')) }}</span>
              </div>
              <div class="info-item info-item--full">
                <span class="info-label">标签</span>
                <template v-if="canEditPermission">
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
                </template>
                <GvSpace v-else-if="currentAsset.tags && currentAsset.tags.length > 0" size="xs" wrap>
                  <GvTag v-for="tag in currentAsset.tags" :key="tag" size="small" color="primary" variant="tonal">{{ tag }}</GvTag>
                </GvSpace>
                <span v-else class="info-value">-</span>
              </div>
              <div class="info-item">
                <span class="info-label">创建时间</span>
                <span class="info-value text-xs">{{ formatDate(currentAsset.created_at) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">ID</span>
                <span class="info-value text-xs text-text-tertiary">{{ currentAsset.id }}</span>
              </div>
            </div>

            <div v-if="canEditPermission" class="asset-detail-actions">
              <GvButton size="small" variant="filled" :loading="savingSection === 'all'" @click="saveAll">保存</GvButton>
            </div>
            <div v-else class="pt-2">
              <GvButton size="small" variant="tonal" @click="showDetailDrawer = false">关闭</GvButton>
            </div>
          </div>
        </GvDrawer>

        <el-dialog v-model="showImagePreview" title="图片预览" width="70%" append-to-body>
          <div class="preview-dialog-content">
            <img v-if="currentAsset" :src="currentAsset.path" :alt="currentAsset.name" class="preview-dialog-image" />
          </div>
        </el-dialog>

        <el-dialog v-model="showVideoPreview" title="视频预览" width="80%" append-to-body>
          <div class="preview-dialog-content">
            <video v-if="currentAsset" :src="currentAsset.path" controls autoplay class="preview-dialog-video">
              您的浏览器不支持视频播放
            </video>
          </div>
        </el-dialog>
      </main>
    </div>
  </GvContainer>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadFile, type UploadFiles } from 'element-plus'
import { Upload, VideoCamera, Picture, Headset, Refresh, FolderOpened, Grid, List } from '@element-plus/icons-vue'
import { assetApi, type MediaAsset, type AssetCreateReq, type AssetUpdateReq } from '@/api/asset'
import { roleApi } from '@/api/role'
import { useTable, useAsyncData } from '@/composables'
import { useBreakpoint } from '@/composables/useBreakpoint'
import GvContainer from '@/components/layout/GvContainer/index.vue'
import GvCard from '@/components/base/GvCard/index.vue'
import GvModal from '@/components/base/GvModal/index.vue'
import GvDrawer from '@/components/base/GvDrawer/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvSpace from '@/components/layout/GvSpace/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'
import GvInput from '@/components/base/GvInput/index.vue'
import GvSelect from '@/components/base/GvSelect/index.vue'
import GvLoading from '@/components/base/GvLoading/index.vue'
import GvTable from '@/components/base/GvTable/index.vue'
import GvUpload from '@/components/base/GvUpload/index.vue'
import SearchBar from '@/components/business/SearchBar/index.vue'
import StatusBadge from '@/components/business/StatusBadge/index.vue'
import AssetCard from '@/components/business/AssetCard/index.vue'
import { LoadingState, ErrorState, EmptyState } from '@/components/common'
import { useUserStore } from '@/store/user'
import { VISIBILITY_OPTIONS } from '@/constants/visibility'

// UI 状态
const uploading = ref(false)
const showUploadDialog = ref(false)
const showDetailDrawer = ref(false)
const savingSection = ref<'all' | ''>('')
const currentAsset = ref<MediaAsset | null>(null)
const uploadFormRef = ref<FormInstance>()
const uploadRef = ref()
const uploadType = ref<'url' | 'file'>('url')
const uploadFileList = ref<UploadFile[]>([])
const selectedFile = ref<UploadFile | null>(null)
const viewMode = ref<'grid' | 'list'>('grid')
const showImagePreview = ref(false)
const showVideoPreview = ref(false)
const { isMobile } = useBreakpoint()

watch(isMobile, (val) => {
  if (val) {
    viewMode.value = 'grid'
  }
}, { immediate: true })

// 筛选参数
const searchName = ref('')
const selectedType = ref<string | null>(null)
const selectedTag = ref<string | null>(null)
const userStore = useUserStore()
const canEditPermission = computed(() => userStore.hasPermission('asset:update'))

// 计算筛选参数
const filterParams = computed(() => ({
  name: searchName.value || undefined,
  type: selectedType.value || undefined,
  tags: selectedTag.value || undefined
}))

// 使用 useTable 管理资产列表
const {
  items: assets,
  isLoading: loading,
  error,
  pagination,
  goToPage,
  changePageSize,
  refreshTable
} = useTable(
  async (params) => {
    const res = await assetApi.list(params)
    return { items: res.data?.items ?? [], total: res.data?.total ?? 0 }
  },
  {
    immediate: true,
    initialPageSize: 12,
    extraParams: filterParams
  }
)

// 使用 useAsyncData 管理标签加载
const {
  data: tagsData,
  isLoading: tagsLoading,
  execute: loadTags
} = useAsyncData(
  () => assetApi.getTags(),
  { immediate: true }
)

const tags = computed(() => tagsData.value?.data.tags || [])

const uploadForm = reactive<any>({
  type: 'video',
  source_type: 'upload',
  name: '',
  path: '',
  size: 0,
  format: '',
  source_id: undefined,
  tags: [],
  visibility: 0
})

const editForm = reactive<any>({
  name: '',
  status: 'ready',
  tags: [],
  visibility: 0
})

const uploadRules: FormRules = {
  name: [{ required: true, message: '请输入资产名称', trigger: 'blur' }],
  type: [
    {
      required: true,
      message: '请选择资产类型',
      trigger: 'change'
    }
  ],
  path: [
    {
      required: true,
      message: '请输入资源地址',
      trigger: 'blur',
      validator: (_rule: unknown, value: string, callback: (e?: Error) => void) => {
        if (uploadType.value === 'url') {
          if (!value || !value.trim()) {
            callback(new Error('请输入资源地址'))
          } else {
            callback()
          }
        } else {
          callback()
        }
      }
    }
  ]
}

const mediaTypes = computed(() => [
  { label: '全部', value: null, icon: FolderOpened },
  { label: '视频', value: 'video', icon: VideoCamera },
  { label: '图片', value: 'image', icon: Picture },
  { label: '音频', value: 'audio', icon: Headset }
])

const typeOptions = [
  { label: '视频', value: 'video' },
  { label: '图片', value: 'image' },
  { label: '音频', value: 'audio' }
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

function handleTypeChange(type: string | null) {
  selectedType.value = type
}

function handleTagChange(tag: string) {
  selectedTag.value = selectedTag.value === tag ? null : tag
}

// 直接使用 useTable 提供的方法
const handlePageChange = goToPage
const handleSizeChange = changePageSize

function handleFileChange(file: UploadFile, fileList: UploadFiles) {
  if (fileList.length > 0 && file.raw) {
    selectedFile.value = file
    uploadForm.name = file.name.split('.')[0]
    uploadForm.size = file.size || 0
    uploadForm.format = file.name.split('.').pop() || ''
    const detectedType = detectAssetType(file.name, file.raw.type)
    if (detectedType) {
      uploadForm.type = detectedType
    }
  } else {
    selectedFile.value = null
    uploadForm.size = 0
    uploadForm.format = ''
  }
}

function handleFileRemove(file: UploadFile, fileList: UploadFiles) {
  if (fileList.length === 0) {
    selectedFile.value = null
    uploadForm.size = 0
    uploadForm.format = ''
  }
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
        await assetApi.upload(
          selectedFile.value.raw,
          uploadForm.type,
          uploadForm.name,
          uploadForm.tags || [],
          uploadForm.visibility ?? 0
        )
      } else {
        const createData: any = {
          type: uploadForm.type,
          source_type: uploadForm.source_type,
          name: uploadForm.name,
          path: uploadForm.path,
          size: uploadForm.size || 0,
          format: uploadForm.format || '',
          tags: uploadForm.tags || [],
          visibility: uploadForm.visibility ?? 0
        }
        await assetApi.create(createData)
      }
      ElMessage.success('添加成功')
      showUploadDialog.value = false
      resetUploadForm()
      refreshTable()
      loadTags()
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || '添加失败')
    } finally {
      uploading.value = false
    }
  })
}

function resetUploadForm() {
  uploadType.value = 'url'
  uploadForm.type = 'video'
  uploadForm.source_type = 'upload'
  uploadForm.name = ''
  uploadForm.path = ''
  uploadForm.size = 0
  uploadForm.format = ''
  uploadForm.source_id = undefined
  uploadForm.tags = []
  uploadForm.visibility = 0
  selectedFile.value = null
  uploadFileList.value = []
  uploadRef.value?.clearFiles()
  uploadFormRef.value?.resetFields()
}

function detectAssetType(nameOrPath: string, mimeType?: string): 'video' | 'image' | 'audio' | null {
  const source = (nameOrPath || '').toLowerCase()
  const mime = (mimeType || '').toLowerCase()

  if (mime.startsWith('image/')) return 'image'
  if (mime.startsWith('video/')) return 'video'
  if (mime.startsWith('audio/')) return 'audio'

  const imageExt = /\.(jpg|jpeg|png|gif|webp|bmp|svg|heic)(\?.*)?$/i
  const videoExt = /\.(mp4|mov|mkv|avi|webm|m4v)(\?.*)?$/i
  const audioExt = /\.(mp3|wav|aac|m4a|flac|ogg)(\?.*)?$/i

  if (imageExt.test(source)) return 'image'
  if (videoExt.test(source)) return 'video'
  if (audioExt.test(source)) return 'audio'

  return null
}

watch(
  () => uploadForm.path,
  (val) => {
    if (uploadType.value !== 'url') return
    const detectedType = detectAssetType(val)
    if (detectedType) {
      uploadForm.type = detectedType
    }
  }
)

function handleDetail(asset: MediaAsset) {
  currentAsset.value = asset
  editForm.name = asset.name
  editForm.status = asset.status
  editForm.tags = asset.tags || []
  editForm.visibility = asset.visibility ?? 0
  showDetailDrawer.value = true
}

function resetEditForm() {
  if (!currentAsset.value) return
  editForm.name = currentAsset.value.name
  editForm.status = currentAsset.value.status
  editForm.tags = currentAsset.value.tags || []
  editForm.visibility = currentAsset.value.visibility ?? 0
}

async function doUpdate(payload: any) {
  if (!currentAsset.value) return

  savingSection.value = 'all'
  try {
    const res = await assetApi.update(currentAsset.value.id, payload)
    const updated = res.data as MediaAsset
    currentAsset.value = updated
    resetEditForm()
    ElMessage.success('保存成功')
    refreshTable()
    loadTags()
  } catch (error: any) {
    if (error?.response?.status === 403) {
      ElMessage.error('无编辑权限')
      return
    }
    ElMessage.error(error.response?.data?.message || '更新失败')
  } finally {
    savingSection.value = ''
  }
}

async function saveAll() {
  await doUpdate({
    name: editForm.name,
    status: editForm.status,
    tags: editForm.tags || [],
    visibility: editForm.visibility ?? 0
  })
}

async function handleCopyLink() {
  if (!currentAsset.value?.path) return
  try {
    await navigator.clipboard.writeText(currentAsset.value.path)
    ElMessage.success('链接已复制')
  } catch {
    ElMessage.error('复制失败，请手动复制')
  }
}

function handleDownload() {
  if (!currentAsset.value?.path) return
  window.open(currentAsset.value.path, '_blank')
}

function openImagePreview() {
  if (currentAsset.value?.type !== 'image') return
  showImagePreview.value = true
}

function openVideoPreview() {
  if (currentAsset.value?.type !== 'video') return
  showVideoPreview.value = true
}

async function handleDelete(asset: MediaAsset) {
  try {
    await ElMessageBox.confirm('确定要删除此资产吗？', '提示', {
      type: 'warning'
    })
    await assetApi.delete(asset.id)
    ElMessage.success('删除成功')
    refreshTable()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

function getTypeIcon(type: string) {
  const iconMap: Record<string, any> = {
    video: VideoCamera,
    image: Picture,
    audio: Headset
  }
  return iconMap[type] || Picture
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
    generated: '生成',
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

.asset-detail-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.asset-detail-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.asset-preview-card {
  background: #f9fafb;
  border-radius: 8px;
  overflow: hidden;
  min-height: 260px;
}

.asset-form-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.info-item--full {
  grid-column: 1 / -1;
}

.asset-detail-actions {
  display: flex;
  justify-content: flex-end;
  position: sticky;
  bottom: 0;
  padding-top: 8px;
  background: linear-gradient(to top, rgba(255, 255, 255, 0.96), rgba(255, 255, 255, 0.65));
  backdrop-filter: blur(4px);
  z-index: 2;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 12px;
  color: #6b7280;
  font-weight: 500;
}

.info-value {
  font-size: 14px;
  color: #111827;
  word-break: break-all;
}

.preview-container {
  position: relative;
  width: 100%;
  min-height: 260px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.preview-dialog-content {
  display: flex;
  justify-content: center;
  align-items: center;
}

.preview-dialog-image,
.preview-dialog-video {
  max-width: 100%;
  max-height: 70vh;
  border-radius: 8px;
}

.preview-media {
  max-width: 100%;
  max-height: 500px;
  width: auto;
  height: auto;
  border-radius: 8px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.preview-media--zoomable {
  cursor: zoom-in;
}

/* 音频预览 */
.audio-preview {
  flex-direction: column;
  gap: 24px;
}

.audio-icon {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

.audio-player {
  width: 100%;
  max-width: 400px;
}

/* 深色模式 */
.dark .info-value {
  color: #f3f4f6;
}

.dark .asset-preview-card {
  background: #1f2937;
}

.dark .asset-detail-actions {
  background: linear-gradient(to top, rgba(17, 24, 39, 0.96), rgba(17, 24, 39, 0.65));
}

@media (max-width: 960px) {
  .asset-form-grid {
    grid-template-columns: 1fr;
  }
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
</style>
