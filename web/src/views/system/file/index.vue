<template>
  <GvContainer max-width="full" class="h-full">
    <div class="flex h-full gap-4">
      <!-- 左侧：类型筛选 -->
      <aside class="w-64 flex-shrink-0">
        <div class="mb-4">
          <h1 class="text-2xl font-bold text-text-primary">文件管理</h1>
        </div>

        <GvCard shadow="sm" padding="md" class="sticky top-4">
          <!-- 文件类型筛选 -->
          <div class="mb-6">
            <h3 class="text-sm font-semibold text-text-primary mb-3">文件类型</h3>
            <div class="space-y-2">
              <div
                v-for="type in fileTypes"
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
              </div>
            </div>
          </div>

          <!-- 状态筛选 -->
          <div>
            <h3 class="text-sm font-semibold text-text-primary mb-3">状态</h3>
            <div class="space-y-2">
              <div
                v-for="status in statusOptions"
                :key="status.value"
                :class="[
                  'px-3 py-2 rounded-lg cursor-pointer transition-all text-sm',
                  selectedStatus === status.value
                    ? 'bg-primary-50 text-primary-600 font-medium'
                    : 'hover:bg-neutral-50 text-text-secondary'
                ]"
                @click="handleStatusChange(status.value)"
              >
                {{ status.label }}
              </div>
            </div>
          </div>
        </GvCard>
      </aside>

      <!-- 右侧：文件列表 -->
      <main class="flex-1 min-w-0">
        <!-- 操作栏 -->
        <div class="flex items-center justify-between mb-6">
          <SearchBar
            v-model="searchName"
            placeholder="搜索文件名称"
            class="w-80"
            immediate
            :show-button="false"
            @search="loadFiles"
          />
          <GvButton @click="showUploadDialog = true" v-permission="'file:create'">
            <template #icon>
              <el-icon><Upload /></el-icon>
            </template>
            上传文件
          </GvButton>
        </div>

        <!-- 文件列表 -->
        <div v-if="loading" class="flex justify-center items-center py-20">
          <GvLoading />
        </div>
        <div v-else-if="files.length === 0" class="text-center py-20">
          <el-icon :size="64" class="text-neutral-300 mb-4">
            <FolderOpened />
          </el-icon>
          <p class="text-text-tertiary">暂无文件</p>
        </div>
        <GvTable
          v-else
          :data="files"
          :columns="tableColumns"
          :loading="loading"
          class="mb-6"
        >
          <template #type="{ row }">
            <GvTag :color="getTypeColor(row.type)" size="small" variant="tonal">
              {{ getTypeLabel(row.type) }}
            </GvTag>
          </template>
          <template #size="{ row }">
            {{ formatSize(row.size) }}
          </template>
          <template #status="{ row }">
            <StatusBadge :status="mapStatus(row.status)" />
          </template>
          <template #created_at="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
          <template #actions="{ row }">
            <GvSpace size="xs">
              <GvButton variant="text" size="small" @click="handleView(row)">查看</GvButton>
              <GvButton variant="text" size="small" @click="handleDownload(row)" v-permission="'file:download'">下载</GvButton>
              <GvButton variant="text" size="small" color="error" @click="handleDelete(row)" v-permission="'file:delete'">删除</GvButton>
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

        <!-- 上传文件对话框 -->
        <GvModal
          v-model="showUploadDialog"
          title="上传文件"
          :confirm-loading="uploading"
          @confirm="handleUpload"
          @cancel="showUploadDialog = false"
        >
          <GvUpload
            ref="uploadRef"
            v-model="uploadFileList"
            :auto-upload="false"
            :limit="10"
            button-text="选择文件"
            tip="支持所有文件类型，单个文件不超过 100MB"
            :max-size="100 * 1024 * 1024"
            @change="handleFileChange"
            @remove="handleFileRemove"
          />
        </GvModal>

        <!-- 文件详情对话框 -->
        <GvModal
          v-model="showViewDialog"
          title="文件详情"
          size="large"
          :show-confirm="false"
          cancel-text="关闭"
        >
          <div v-if="currentFile" class="file-detail-container">
            <div class="grid grid-cols-2 gap-6">
              <div class="space-y-4">
                <div class="info-item">
                  <span class="info-label">文件名</span>
                  <span class="info-value">{{ currentFile.name }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">原始文件名</span>
                  <span class="info-value">{{ currentFile.original_name }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">类型</span>
                  <GvTag :color="getTypeColor(currentFile.type)" size="small">
                    {{ getTypeLabel(currentFile.type) }}
                  </GvTag>
                </div>
                <div class="info-item">
                  <span class="info-label">大小</span>
                  <span class="info-value">{{ formatSize(currentFile.size) }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">MIME 类型</span>
                  <span class="info-value">{{ currentFile.mime_type }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">状态</span>
                  <StatusBadge :status="mapStatus(currentFile.status)" />
                </div>
                <div class="info-item">
                  <span class="info-label">创建时间</span>
                  <span class="info-value text-xs">{{ formatDate(currentFile.created_at) }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">文件 URL</span>
                  <span class="info-value text-xs text-text-tertiary break-all">{{ currentFile.url }}</span>
                </div>
              </div>
              <div class="preview-container">
                <!-- 图片预览 -->
                <img
                  v-if="currentFile.type === 'image'"
                  :src="currentFile.url"
                  :alt="currentFile.name"
                  class="w-full h-auto rounded-lg"
                />
                <!-- 视频预览 -->
                <video
                  v-else-if="currentFile.type === 'video'"
                  :src="currentFile.url"
                  controls
                  class="w-full rounded-lg"
                >
                  您的浏览器不支持视频播放
                </video>
                <!-- 音频预览 -->
                <div v-else-if="currentFile.type === 'audio'" class="flex flex-col items-center justify-center h-full">
                  <el-icon :size="80" class="text-primary-500 mb-4">
                    <Headset />
                  </el-icon>
                  <audio :src="currentFile.url" controls class="w-full">
                    您的浏览器不支持音频播放
                  </audio>
                </div>
                <!-- 其他类型 -->
                <div v-else class="flex flex-col items-center justify-center h-full text-text-tertiary">
                  <el-icon :size="80" class="mb-4">
                    <Document />
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
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type UploadFile, type UploadFiles } from 'element-plus'
import {
  Picture,
  VideoCamera,
  Headset,
  Document,
  FolderOpened,
  Upload,
  Refresh
} from '@element-plus/icons-vue'
import { fileApi, type File, type FileListQuery } from '@/api/file'
import {
  GvContainer,
  GvCard,
  GvModal,
  GvButton,
  GvSpace,
  GvTag,
  GvLoading,
  GvTable,
  GvUpload,
  SearchBar,
  StatusBadge
} from '@/components'

const loading = ref(false)
const uploading = ref(false)
const files = ref<File[]>([])
const showUploadDialog = ref(false)
const showViewDialog = ref(false)
const currentFile = ref<File | null>(null)
const uploadRef = ref()
const uploadFileList = ref<UploadFile[]>([])

const searchName = ref('')
const selectedType = ref<string | null>(null)
const selectedStatus = ref<string | null>(null)

const pagination = reactive({
  page: 1,
  page_size: 12,
  total: 0
})

const fileTypes = computed(() => [
  { label: '全部', value: null, icon: FolderOpened },
  { label: '图片', value: 'image', icon: Picture },
  { label: '视频', value: 'video', icon: VideoCamera },
  { label: '音频', value: 'audio', icon: Headset },
  { label: '文档', value: 'document', icon: Document },
  { label: '其他', value: 'other', icon: Document }
])

const statusOptions = [
  { label: '全部', value: null },
  { label: '上传中', value: 'uploading' },
  { label: '已完成', value: 'completed' },
  { label: '失败', value: 'failed' },
  { label: '已删除', value: 'deleted' }
]

const tableColumns = [
  { prop: 'name', label: '文件名', minWidth: 200 },
  { prop: 'type', label: '类型', width: 100 },
  { prop: 'size', label: '大小', width: 100 },
  { prop: 'mime_type', label: 'MIME 类型', width: 150 },
  { prop: 'status', label: '状态', width: 100 },
  { prop: 'created_at', label: '创建时间', width: 180 },
  { prop: 'actions', label: '操作', width: 200, fixed: 'right' }
]

onMounted(() => {
  loadFiles()
})

async function loadFiles() {
  loading.value = true
  try {
    const params: FileListQuery = {
      type: selectedType.value as any,
      status: selectedStatus.value as any,
      search: searchName.value || undefined,
      page: pagination.page,
      page_size: pagination.page_size
    }
    const response = await fileApi.list(params)
    files.value = response.data.items
    pagination.total = response.data.total
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '加载失败')
  } finally {
    loading.value = false
  }
}

function handleTypeChange(type: string | null) {
  selectedType.value = type
  pagination.page = 1
  loadFiles()
}

function handleStatusChange(status: string | null) {
  selectedStatus.value = status
  pagination.page = 1
  loadFiles()
}

function handlePageChange(page: number) {
  pagination.page = page
  loadFiles()
}

function handleSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  loadFiles()
}

function handleFileChange(file: UploadFile, fileList: UploadFiles) {
  // 文件选择变化
}

function handleFileRemove(file: UploadFile, fileList: UploadFiles) {
  // 文件移除
}

async function handleUpload() {
  if (uploadFileList.value.length === 0) {
    ElMessage.warning('请选择文件')
    return
  }

  uploading.value = true
  try {
    for (const fileItem of uploadFileList.value) {
      if (fileItem.raw) {
        await fileApi.upload(fileItem.raw, (progress) => {
          // 更新进度
          if (fileItem.uid) {
            const index = uploadFileList.value.findIndex((f) => f.uid === fileItem.uid)
            if (index > -1) {
              uploadFileList.value[index].percentage = progress
            }
          }
        })
      }
    }
    ElMessage.success('上传成功')
    showUploadDialog.value = false
    uploadFileList.value = []
    uploadRef.value?.clearFiles()
    loadFiles()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '上传失败')
  } finally {
    uploading.value = false
  }
}

function handleView(file: File) {
  currentFile.value = file
  showViewDialog.value = true
}

function handleDownload(file: File) {
  window.open(file.url, '_blank')
}

async function handleDelete(file: File) {
  try {
    await ElMessageBox.confirm('确定要删除此文件吗？', '提示', {
      type: 'warning'
    })
    await fileApi.delete(file.id)
    ElMessage.success('删除成功')
    loadFiles()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

function getTypeLabel(type: string) {
  const map: Record<string, string> = {
    image: '图片',
    video: '视频',
    audio: '音频',
    document: '文档',
    archive: '压缩包',
    other: '其他'
  }
  return map[type] || type
}

function getTypeColor(type: string) {
  const map: Record<string, string> = {
    image: 'success',
    video: 'primary',
    audio: 'warning',
    document: 'info',
    archive: 'neutral',
    other: 'neutral'
  }
  return map[type] || 'neutral'
}

function mapStatus(status: string): any {
  const map: Record<string, string> = {
    completed: 'success',
    uploading: 'processing',
    failed: 'error',
    deleted: 'inactive'
  }
  return map[status] || 'inactive'
}

function formatSize(size: number): string {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`
  if (size < 1024 * 1024 * 1024) return `${(size / (1024 * 1024)).toFixed(2)} MB`
  return `${(size / (1024 * 1024 * 1024)).toFixed(2)} GB`
}

function formatDate(date: string): string {
  return new Date(date).toLocaleString('zh-CN')
}
</script>

<style scoped>
.view-switch-group {
  @apply flex items-center gap-1 bg-neutral-100 rounded-lg p-1;
}

.view-switch-btn {
  @apply px-3 py-1.5 rounded transition-all text-text-secondary;
}

.view-switch-btn:hover {
  @apply bg-white text-text-primary;
}

.view-switch-btn.active {
  @apply bg-white text-primary-600 shadow-sm;
}

.file-detail-container {
  @apply p-4;
}

.info-item {
  @apply flex items-start gap-4 py-2 border-b border-neutral-100 last:border-0;
}

.info-label {
  @apply text-sm font-medium text-text-secondary w-24 flex-shrink-0;
}

.info-value {
  @apply text-sm text-text-primary flex-1;
}

.preview-container {
  @apply bg-neutral-50 rounded-lg p-4 min-h-[300px] flex items-center justify-center;
}

.dark .preview-container {
  @apply bg-neutral-800;
}
</style>
