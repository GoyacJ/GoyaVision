<template>
  <GvContainer max-width="full" class="h-full">
    <div class="flex h-full gap-4">
      <!-- 左侧：类型和标签筛选 -->
      <aside class="w-64 flex-shrink-0">
        <!-- 页面标题 -->
        <div class="mb-4">
          <h1 class="text-2xl font-bold text-text-primary">媒体资产库</h1>
        </div>

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
        <!-- 操作栏 -->
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center gap-3">
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
          </div>
          <GvButton @click="showUploadDialog = true">
            <template #icon>
              <el-icon><Upload /></el-icon>
            </template>
            添加资产
          </GvButton>
        </div>

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
            <el-tab-pane label="流媒体接入" name="stream" />
          </el-tabs>

          <el-form ref="uploadFormRef" :model="uploadForm" :rules="uploadRules" label-width="100px">
            <el-form-item label="资产名称" prop="name">
              <GvInput v-model="uploadForm.name" placeholder="请输入资产名称" />
            </el-form-item>
            <el-form-item v-if="uploadType !== 'stream'" label="资产类型" prop="type">
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

            <!-- 流媒体接入模式 -->
            <template v-else-if="uploadType === 'stream'">
              <el-form-item label="创建方式">
                <el-radio-group v-model="streamCreateMode">
                  <el-radio value="url">输入流地址（新建媒体源并创建资产）</el-radio>
                  <el-radio value="from_source">从已有媒体源创建资产</el-radio>
                </el-radio-group>
              </el-form-item>
              <el-form-item v-if="streamCreateMode === 'url'" label="流地址" prop="path">
                <GvInput
                  v-model="uploadForm.path"
                  placeholder="请输入流地址，如 rtsp://...、rtmp://...、https://.../live.m3u8"
                  type="textarea"
                  :rows="2"
                />
                <div class="text-text-tertiary text-xs mt-1">
                  支持 RTSP、RTMP、HLS 等协议，将自动接入 MediaMTX 并创建媒体源与资产
                </div>
              </el-form-item>
              <el-form-item v-else label="媒体源" prop="source_id">
                <el-select
                  v-model="uploadForm.source_id"
                  placeholder="请选择已有媒体源"
                  filterable
                  class="w-full"
                  :loading="sourcesLoading"
                >
                  <el-option
                    v-for="s in sources"
                    :key="s.id"
                    :label="`${s.name} (${s.path_name})`"
                    :value="s.id"
                  />
                </el-select>
                <div class="text-text-tertiary text-xs mt-1">
                  在媒体源管理页可先创建拉流/推流源，再在此处选源创建资产
                </div>
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
          <div v-if="currentAsset" class="asset-detail-container">
            <!-- 左侧：资产信息 -->
            <div class="asset-detail-info">
              <div class="info-section">
                <div class="info-item">
                  <span class="info-label">名称</span>
                  <span class="info-value">{{ currentAsset.name }}</span>
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
                  <StatusBadge :status="mapStatus(currentAsset.status)" />
                </div>
                <div v-if="currentAsset.tags && currentAsset.tags.length > 0" class="info-item">
                  <span class="info-label">标签</span>
                  <GvSpace size="xs" wrap>
                    <GvTag v-for="tag in currentAsset.tags" :key="tag" size="small" color="primary" variant="tonal">
                      {{ tag }}
                    </GvTag>
                  </GvSpace>
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
            </div>

            <!-- 右侧：资源预览 -->
            <div class="asset-detail-preview">
              <!-- 视频预览 -->
              <div v-if="currentAsset.type === 'video'" class="preview-container">
                <video
                  :src="currentAsset.path"
                  controls
                  class="preview-media"
                >
                  您的浏览器不支持视频播放
                </video>
              </div>

              <!-- 图片预览 -->
              <div v-else-if="currentAsset.type === 'image'" class="preview-container">
                <img
                  :src="currentAsset.path"
                  :alt="currentAsset.name"
                  class="preview-media"
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

              <!-- 流媒体预览 -->
              <div v-else-if="currentAsset.type === 'stream'" class="preview-container stream-preview">
                <div class="stream-info">
                  <el-icon :size="80" class="text-info-500 mb-4">
                    <Connection />
                  </el-icon>
                  <p class="text-text-secondary mb-2">流媒体地址</p>
                  <p class="text-sm text-text-primary font-mono bg-neutral-50 px-3 py-2 rounded break-all">
                    {{ currentAsset.path }}
                  </p>
                </div>
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
          </div>
        </GvModal>
      </main>
    </div>
  </GvContainer>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadFile, type UploadFiles } from 'element-plus'
import { Upload, VideoCamera, Picture, Headset, Connection, Refresh, FolderOpened, Grid, List } from '@element-plus/icons-vue'
import { assetApi, type MediaAsset, type AssetCreateReq, type AssetUpdateReq } from '@/api/asset'
import { sourceApi } from '@/api/source'
import GvContainer from '@/components/layout/GvContainer/index.vue'
import GvCard from '@/components/base/GvCard/index.vue'
import GvModal from '@/components/base/GvModal/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvSpace from '@/components/layout/GvSpace/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'
import GvBadge from '@/components/base/GvBadge/index.vue'
import GvInput from '@/components/base/GvInput/index.vue'
import GvSelect from '@/components/base/GvSelect/index.vue'
import GvLoading from '@/components/base/GvLoading/index.vue'
import GvTable from '@/components/base/GvTable/index.vue'
import GvUpload from '@/components/base/GvUpload/index.vue'
import PageHeader from '@/components/business/PageHeader/index.vue'
import SearchBar from '@/components/business/SearchBar/index.vue'
import StatusBadge from '@/components/business/StatusBadge/index.vue'
import AssetCard from '@/components/business/AssetCard/index.vue'

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
const uploadType = ref<'url' | 'file' | 'stream'>('url')
const streamCreateMode = ref<'url' | 'from_source'>('url')
const sources = ref<import('@/api/source').MediaSource[]>([])
const sourcesLoading = ref(false)
const uploadFileList = ref<UploadFile[]>([])
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
  source_id: undefined,
  tags: []
})

const editForm = reactive<AssetUpdateReq>({
  name: '',
  status: 'ready',
  tags: []
})

const uploadRules: FormRules = {
  name: [{ required: true, message: '请输入资产名称', trigger: 'blur' }],
  type: [
    {
      required: true,
      message: '请选择资产类型',
      trigger: 'change',
      validator: (_rule: unknown, value: string, callback: (e?: Error) => void) => {
        if (uploadType.value === 'stream') {
          callback()
        } else if (!value) {
          callback(new Error('请选择资产类型'))
        } else {
          callback()
        }
      }
    }
  ],
  path: [
    {
      required: true,
      message: '请输入资源地址或流地址',
      trigger: 'blur',
      validator: (_rule: unknown, value: string, callback: (e?: Error) => void) => {
        if (uploadType.value === 'url') {
          if (!value || !value.trim()) {
            callback(new Error('请输入资源地址'))
          } else {
            callback()
          }
        } else if (uploadType.value === 'stream' && streamCreateMode.value === 'url') {
          if (!value || !value.trim()) {
            callback(new Error('请输入流地址'))
          } else {
            callback()
          }
        } else {
          callback()
        }
      }
    }
  ],
  source_id: [
    {
      validator: (_rule: unknown, _value: string, callback: (e?: Error) => void) => {
        if (uploadType.value === 'stream' && streamCreateMode.value === 'from_source') {
          if (!uploadForm.source_id) {
            callback(new Error('请选择媒体源'))
          } else {
            callback()
          }
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ]
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

watch(uploadType, (t) => {
  if (t === 'stream') {
    uploadForm.type = 'stream'
    uploadForm.source_type = 'live'
    if (streamCreateMode.value === 'from_source') {
      loadSources()
    }
  }
})

watch([uploadType, streamCreateMode], () => {
  if (uploadType.value === 'stream' && streamCreateMode.value === 'from_source') {
    loadSources()
  }
})

async function loadSources() {
  sourcesLoading.value = true
  try {
    const res = await sourceApi.list({ limit: 200 })
    sources.value = res.data?.items ?? []
  } catch {
    sources.value = []
  } finally {
    sourcesLoading.value = false
  }
}

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

function handleFileChange(file: UploadFile, fileList: UploadFiles) {
  if (fileList.length > 0 && file.raw) {
    selectedFile.value = file
    uploadForm.name = file.name.split('.')[0]
    uploadForm.size = file.size || 0
    uploadForm.format = file.name.split('.').pop() || ''
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
          uploadForm.tags || []
        )
      } else if (uploadType.value === 'stream') {
        if (streamCreateMode.value === 'url') {
          await assetApi.create({
            type: 'stream',
            source_type: 'live',
            name: uploadForm.name,
            stream_url: uploadForm.path!.trim(),
            size: 0,
            format: '',
            tags: uploadForm.tags || []
          })
        } else {
          await assetApi.create({
            type: 'stream',
            source_type: 'live',
            name: uploadForm.name,
            source_id: uploadForm.source_id!,
            size: 0,
            format: '',
            tags: uploadForm.tags || []
          })
        }
      } else {
        const createData = {
          type: uploadForm.type,
          source_type: uploadForm.source_type,
          name: uploadForm.name,
          path: uploadForm.path,
          size: uploadForm.size || 0,
          format: uploadForm.format || '',
          tags: uploadForm.tags || []
        }
        await assetApi.create(createData)
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
  uploadType.value = 'url'
  streamCreateMode.value = 'url'
  uploadForm.type = 'video'
  uploadForm.source_type = 'upload'
  uploadForm.name = ''
  uploadForm.path = ''
  uploadForm.size = 0
  uploadForm.format = ''
  uploadForm.source_id = undefined
  uploadForm.tags = []
  selectedFile.value = null
  uploadFileList.value = []
  uploadRef.value?.clearFiles()
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

function getTypeIcon(type: string) {
  const iconMap: Record<string, any> = {
    video: VideoCamera,
    image: Picture,
    audio: Headset,
    stream: Connection
  }
  return iconMap[type] || Picture
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

/* 资产详情两栏布局 */
.asset-detail-container {
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 24px;
  min-height: 400px;
}

/* 左侧信息区域 */
.asset-detail-info {
  border-right: 1px solid #e5e7eb;
  padding-right: 24px;
}

.info-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
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

/* 右侧预览区域 */
.asset-detail-preview {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f9fafb;
  border-radius: 8px;
  overflow: hidden;
}

.preview-container {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.preview-media {
  max-width: 100%;
  max-height: 500px;
  width: auto;
  height: auto;
  border-radius: 8px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
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

/* 流媒体预览 */
.stream-preview {
  padding: 40px;
}

.stream-info {
  text-align: center;
  max-width: 500px;
}

/* 深色模式 */
.dark .asset-detail-info {
  border-right-color: #374151;
}

.dark .info-value {
  color: #f3f4f6;
}

.dark .asset-detail-preview {
  background: #1f2937;
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
