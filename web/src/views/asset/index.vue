<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>媒体资产库</span>
          <div class="header-actions">
            <el-input
              v-model="searchName"
              placeholder="搜索资产名称"
              clearable
              style="width: 200px; margin-right: 10px"
              @change="loadAssets"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-button type="primary" @click="showUploadDialog = true">
              <el-icon><Upload /></el-icon>
              上传资产
            </el-button>
          </div>
        </div>
      </template>

      <div class="filter-bar">
        <el-space wrap>
          <el-select v-model="filterType" placeholder="资产类型" clearable @change="loadAssets" style="width: 120px">
            <el-option label="视频" value="video" />
            <el-option label="图片" value="image" />
            <el-option label="音频" value="audio" />
          </el-select>
          <el-select v-model="filterSourceType" placeholder="来源类型" clearable @change="loadAssets" style="width: 150px">
            <el-option label="上传" value="upload" />
            <el-option label="流捕获" value="stream_capture" />
            <el-option label="算子输出" value="operator_output" />
          </el-select>
          <el-select v-model="filterStatus" placeholder="状态" clearable @change="loadAssets" style="width: 120px">
            <el-option label="就绪" value="ready" />
            <el-option label="处理中" value="processing" />
            <el-option label="待处理" value="pending" />
            <el-option label="错误" value="error" />
          </el-select>
        </el-space>
      </div>

      <el-table :data="assets" v-loading="loading" stripe>
        <el-table-column prop="name" label="名称" min-width="180" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="90">
          <template #default="{ row }">
            <el-tag :type="getTypeColor(row.type)" size="small">
              {{ getTypeLabel(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="source_type" label="来源" width="120">
          <template #default="{ row }">
            <el-tag type="info" size="small">
              {{ getSourceTypeLabel(row.source_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="format" label="格式" width="90" />
        <el-table-column label="大小" width="100">
          <template #default="{ row }">
            {{ formatSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column label="时长" width="100">
          <template #default="{ row }">
            {{ row.duration ? formatDuration(row.duration) : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusColor(row.status)" size="small">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleView(row)">查看</el-button>
            <el-button link type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadAssets"
          @current-change="loadAssets"
        />
      </div>
    </el-card>

    <el-dialog v-model="showUploadDialog" title="上传资产" width="500px">
      <el-form ref="uploadFormRef" :model="uploadForm" :rules="uploadRules" label-width="100px">
        <el-form-item label="资产名称" prop="name">
          <el-input v-model="uploadForm.name" placeholder="请输入资产名称" />
        </el-form-item>
        <el-form-item label="资产类型" prop="type">
          <el-select v-model="uploadForm.type" placeholder="请选择类型" style="width: 100%">
            <el-option label="视频" value="video" />
            <el-option label="图片" value="image" />
            <el-option label="音频" value="audio" />
          </el-select>
        </el-form-item>
        <el-form-item label="文件路径" prop="path">
          <el-input v-model="uploadForm.path" placeholder="请输入文件路径（相对或绝对路径）" />
        </el-form-item>
        <el-form-item label="文件大小" prop="size">
          <el-input v-model.number="uploadForm.size" type="number" placeholder="文件大小（字节）" />
        </el-form-item>
        <el-form-item label="格式" prop="format">
          <el-input v-model="uploadForm.format" placeholder="如：mp4, jpg, mp3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showUploadDialog = false">取消</el-button>
        <el-button type="primary" :loading="uploading" @click="handleUpload">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showEditDialog" title="编辑资产" width="500px">
      <el-form ref="editFormRef" :model="editForm" :rules="editRules" label-width="100px">
        <el-form-item label="资产名称" prop="name">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="editForm.status" style="width: 100%">
            <el-option label="就绪" value="ready" />
            <el-option label="处理中" value="processing" />
            <el-option label="待处理" value="pending" />
            <el-option label="错误" value="error" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleUpdate">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showViewDialog" title="资产详情" width="600px">
      <el-descriptions v-if="currentAsset" :column="1" border>
        <el-descriptions-item label="ID">{{ currentAsset.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ currentAsset.name }}</el-descriptions-item>
        <el-descriptions-item label="类型">
          <el-tag :type="getTypeColor(currentAsset.type)" size="small">
            {{ getTypeLabel(currentAsset.type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="来源">
          <el-tag type="info" size="small">
            {{ getSourceTypeLabel(currentAsset.source_type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="路径">{{ currentAsset.path }}</el-descriptions-item>
        <el-descriptions-item label="格式">{{ currentAsset.format || '-' }}</el-descriptions-item>
        <el-descriptions-item label="大小">{{ formatSize(currentAsset.size) }}</el-descriptions-item>
        <el-descriptions-item label="时长">{{ currentAsset.duration ? formatDuration(currentAsset.duration) : '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusColor(currentAsset.status)" size="small">
            {{ getStatusLabel(currentAsset.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentAsset.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentAsset.updated_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { assetApi, type MediaAsset, type AssetCreateReq, type AssetUpdateReq } from '@/api/asset'

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
const filterType = ref('')
const filterSourceType = ref('')
const filterStatus = ref('')

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

onMounted(() => {
  loadAssets()
})

async function loadAssets() {
  loading.value = true
  try {
    const response = await assetApi.list({
      name: searchName.value || undefined,
      type: filterType.value as any,
      source_type: filterSourceType.value as any,
      status: filterStatus.value as any,
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
  const map: Record<string, any> = {
    video: 'primary',
    image: 'success',
    audio: 'warning'
  }
  return map[type] || ''
}

function getSourceTypeLabel(type: string) {
  const map: Record<string, string> = {
    upload: '上传',
    stream_capture: '流捕获',
    operator_output: '算子输出'
  }
  return map[type] || type
}

function getStatusLabel(status: string) {
  const map: Record<string, string> = {
    ready: '就绪',
    processing: '处理中',
    pending: '待处理',
    error: '错误'
  }
  return map[status] || status
}

function getStatusColor(status: string) {
  const map: Record<string, any> = {
    ready: 'success',
    processing: 'warning',
    pending: 'info',
    error: 'danger'
  }
  return map[status] || ''
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
.page-container {
  padding: 0;
  animation: fadeIn 0.4s ease-out;
}

:deep(.el-card) {
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(31, 38, 135, 0.12);
  border: 1px solid rgba(102, 126, 234, 0.1);
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  overflow: hidden;
}

:deep(.el-card__header) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
  padding: 20px 24px;
}

:deep(.el-card__body) {
  padding: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header > span {
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

:deep(.header-actions .el-input__wrapper) {
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  transition: all 0.3s;
}

:deep(.header-actions .el-input__wrapper:hover) {
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
}

:deep(.header-actions .el-button--primary) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 10px;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  transition: all 0.3s;
}

:deep(.header-actions .el-button--primary:hover) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
}

.filter-bar {
  margin-bottom: 20px;
  padding: 16px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.03) 0%, rgba(118, 75, 162, 0.03) 100%);
  border-radius: 12px;
  border: 1px solid rgba(102, 126, 234, 0.08);
}

:deep(.filter-bar .el-select .el-input__wrapper),
:deep(.filter-bar .el-input__wrapper) {
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  transition: all 0.3s;
}

:deep(.filter-bar .el-select .el-input__wrapper:hover),
:deep(.filter-bar .el-input__wrapper:hover) {
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.12);
}

:deep(.el-table) {
  border-radius: 12px;
  overflow: hidden;
}

:deep(.el-table__header-wrapper) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
}

:deep(.el-table thead) {
  color: #333;
  font-weight: 600;
}

:deep(.el-table th.el-table__cell) {
  background: transparent;
  border-bottom: 2px solid rgba(102, 126, 234, 0.15);
}

:deep(.el-table .el-table__row) {
  transition: all 0.3s;
}

:deep(.el-table .el-table__row:hover > td) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.03) 0%, rgba(118, 75, 162, 0.03) 100%);
}

:deep(.el-tag) {
  border-radius: 6px;
  padding: 4px 12px;
  font-weight: 500;
  border: none;
}

:deep(.el-tag--primary) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.15) 0%, rgba(118, 75, 162, 0.15) 100%);
  color: #667eea;
}

:deep(.el-tag--success) {
  background: linear-gradient(135deg, rgba(103, 194, 58, 0.15) 0%, rgba(56, 249, 215, 0.15) 100%);
  color: #67c23a;
}

:deep(.el-tag--warning) {
  background: linear-gradient(135deg, rgba(230, 162, 60, 0.15) 0%, rgba(245, 87, 108, 0.15) 100%);
  color: #e6a23c;
}

:deep(.el-tag--danger) {
  background: linear-gradient(135deg, rgba(245, 108, 108, 0.15) 0%, rgba(245, 87, 108, 0.15) 100%);
  color: #f56c6c;
}

:deep(.el-tag--info) {
  background: linear-gradient(135deg, rgba(144, 147, 153, 0.15) 0%, rgba(144, 147, 153, 0.15) 100%);
  color: #909399;
}

:deep(.el-button--text),
:deep(.el-button--link) {
  font-weight: 500;
  transition: all 0.3s;
}

:deep(.el-button--text:hover),
:deep(.el-button--link:hover) {
  transform: translateX(2px);
}

.pagination {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
  padding: 16px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.03) 0%, rgba(118, 75, 162, 0.03) 100%);
  border-radius: 12px;
}

:deep(.el-pagination) {
  font-weight: 500;
}

:deep(.el-pagination .el-pager li) {
  border-radius: 8px;
  transition: all 0.3s;
}

:deep(.el-pagination .el-pager li.is-active) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

:deep(.el-pagination .btn-prev),
:deep(.el-pagination .btn-next) {
  border-radius: 8px;
  transition: all 0.3s;
}

:deep(.el-pagination .btn-prev:hover),
:deep(.el-pagination .btn-next:hover) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  color: #667eea;
}

:deep(.el-dialog) {
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(31, 38, 135, 0.2);
  border: 1px solid rgba(102, 126, 234, 0.1);
}

:deep(.el-dialog__header) {
  padding: 24px 24px 16px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
  border-radius: 16px 16px 0 0;
}

:deep(.el-dialog__title) {
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-dialog__footer) {
  padding: 16px 24px 24px;
  border-top: 1px solid rgba(102, 126, 234, 0.1);
}

:deep(.el-form-item__label) {
  font-weight: 600;
  color: #333;
}

:deep(.el-input__wrapper) {
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  transition: all 0.3s;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.12);
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.25);
}

:deep(.el-select .el-input__wrapper) {
  border-radius: 10px;
}

:deep(.el-descriptions) {
  border-radius: 12px;
  overflow: hidden;
}

:deep(.el-descriptions__label) {
  font-weight: 600;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
