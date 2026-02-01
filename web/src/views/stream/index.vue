<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>视频流管理</span>
          <el-button type="primary" @click="showCreateDialog = true">
            <el-icon><Plus /></el-icon>
            添加流
          </el-button>
        </div>
      </template>

      <el-table :data="streams" v-loading="loading" stripe>
        <el-table-column prop="name" label="名称" width="160" />
        <el-table-column prop="type" label="类型" width="80">
          <template #default="{ row }">
            <el-tag :type="row.type === 'push' ? 'warning' : 'primary'" size="small">
              {{ row.type === 'push' ? '推流' : '拉流' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="url" label="源地址" min-width="200" show-overflow-tooltip />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
            <el-tag v-if="row.status?.ready" type="success" size="small" class="ml-1">在线</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="读取者" width="80" v-if="showStatus">
          <template #default="{ row }">
            {{ row.status?.reader_count || 0 }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="320" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="primary" size="small" @click="handlePreview(row)">预览</el-button>
            <el-button link type="primary" size="small" @click="handleRecord(row)">录制</el-button>
            <el-button link type="primary" size="small" @click="handlePlayback(row)">点播</el-button>
            <el-button link type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建流对话框 -->
    <el-dialog v-model="showCreateDialog" title="添加流" width="500px">
      <el-form ref="createFormRef" :model="createForm" :rules="rules" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入流名称" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="createForm.type">
            <el-radio value="pull">拉流（从源地址拉取）</el-radio>
            <el-radio value="push">推流（等待推送）</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="源地址" prop="url" v-if="createForm.type === 'pull'">
          <el-input v-model="createForm.url" placeholder="rtsp://..." />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="createForm.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>

    <!-- 编辑流对话框 -->
    <el-dialog v-model="showEditDialog" title="编辑流" width="500px">
      <el-form ref="editFormRef" :model="editForm" :rules="editRules" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="源地址" prop="url" v-if="editingStream?.type === 'pull'">
          <el-input v-model="editForm.url" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="editForm.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleUpdate">确定</el-button>
      </template>
    </el-dialog>

    <!-- 预览对话框 -->
    <el-dialog v-model="showPreviewDialog" title="视频预览" width="850px" destroy-on-close>
      <el-tabs v-model="previewTab">
        <el-tab-pane label="HLS 播放" name="hls">
          <HLSPreview v-if="previewUrls?.hls_url" :hls-url="previewUrls.hls_url" :width="780" :height="439" />
        </el-tab-pane>
        <el-tab-pane label="协议地址" name="urls">
          <el-descriptions :column="1" border>
            <el-descriptions-item label="RTSP">
              <el-input :value="previewUrls?.rtsp_url" readonly>
                <template #append>
                  <el-button @click="copyUrl(previewUrls?.rtsp_url)">复制</el-button>
                </template>
              </el-input>
            </el-descriptions-item>
            <el-descriptions-item label="RTMP">
              <el-input :value="previewUrls?.rtmp_url" readonly>
                <template #append>
                  <el-button @click="copyUrl(previewUrls?.rtmp_url)">复制</el-button>
                </template>
              </el-input>
            </el-descriptions-item>
            <el-descriptions-item label="HLS">
              <el-input :value="previewUrls?.hls_url" readonly>
                <template #append>
                  <el-button @click="copyUrl(previewUrls?.hls_url)">复制</el-button>
                </template>
              </el-input>
            </el-descriptions-item>
            <el-descriptions-item label="WebRTC">
              <el-input :value="previewUrls?.webrtc_url" readonly>
                <template #append>
                  <el-button @click="copyUrl(previewUrls?.webrtc_url)">复制</el-button>
                </template>
              </el-input>
            </el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>

    <!-- 点播对话框 -->
    <el-dialog v-model="showPlaybackDialog" title="录制点播" width="850px" destroy-on-close>
      <el-table :data="playbackSegments" v-loading="loadingSegments" stripe max-height="300">
        <el-table-column prop="start" label="开始时间" width="200">
          <template #default="{ row }">
            {{ formatDate(row.start) }}
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="{ row }">
            <el-button link type="primary" @click="playSegment(row)">播放</el-button>
            <el-button link type="primary" @click="copyUrl(row.playback_url)">复制链接</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="currentPlaybackUrl" class="mt-4">
        <h4>正在播放</h4>
        <HLSPreview :hls-url="currentPlaybackUrl" :width="780" :height="439" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { 
  streamApi, 
  type StreamWithStatus, 
  type StreamCreateReq, 
  type StreamUpdateReq,
  type PreviewURLs,
  type PlaybackSegment
} from '../../api/stream'
import HLSPreview from '../../components/HLSPreview.vue'

const streams = ref<StreamWithStatus[]>([])
const loading = ref(false)
const showStatus = ref(true)
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const showPreviewDialog = ref(false)
const showPlaybackDialog = ref(false)
const previewTab = ref('hls')
const previewUrls = ref<PreviewURLs | null>(null)
const playbackSegments = ref<PlaybackSegment[]>([])
const loadingSegments = ref(false)
const currentPlaybackUrl = ref('')
const currentEditId = ref('')
const editingStream = ref<StreamWithStatus | null>(null)
const createFormRef = ref<FormInstance>()
const editFormRef = ref<FormInstance>()

const createForm = ref<StreamCreateReq>({
  name: '',
  url: '',
  type: 'pull',
  enabled: true
})

const editForm = ref<StreamUpdateReq>({
  name: '',
  url: '',
  enabled: true
})

const rules = computed<FormRules>(() => ({
  name: [{ required: true, message: '请输入流名称', trigger: 'blur' }],
  url: createForm.value.type === 'pull' 
    ? [{ required: true, message: '请输入源地址', trigger: 'blur' }]
    : []
}))

const editRules: FormRules = {
  name: [{ required: true, message: '请输入流名称', trigger: 'blur' }]
}

const loadStreams = async () => {
  loading.value = true
  try {
    const res = await streamApi.list(undefined, showStatus.value)
    streams.value = res.data
  } catch (error: any) {
    ElMessage.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

const handleCreate = async () => {
  try {
    await createFormRef.value?.validate()
    await streamApi.create(createForm.value)
    ElMessage.success('创建成功')
    showCreateDialog.value = false
    createForm.value = { name: '', url: '', type: 'pull', enabled: true }
    loadStreams()
  } catch (error: any) {
    if (error !== 'cancel' && error?.message) {
      ElMessage.error(error.message || '创建失败')
    }
  }
}

const handleEdit = (row: StreamWithStatus) => {
  currentEditId.value = row.id
  editingStream.value = row
  editForm.value = {
    name: row.name,
    url: row.url,
    enabled: row.enabled
  }
  showEditDialog.value = true
}

const handleUpdate = async () => {
  try {
    await editFormRef.value?.validate()
    await streamApi.update(currentEditId.value, editForm.value)
    ElMessage.success('更新成功')
    showEditDialog.value = false
    loadStreams()
  } catch (error: any) {
    if (error !== 'cancel' && error?.message) {
      ElMessage.error(error.message || '更新失败')
    }
  }
}

const handleDelete = async (row: StreamWithStatus) => {
  try {
    await ElMessageBox.confirm('确定要删除这个流吗？', '确认删除', {
      type: 'warning'
    })
    await streamApi.delete(row.id)
    ElMessage.success('删除成功')
    loadStreams()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const handlePreview = async (row: StreamWithStatus) => {
  try {
    const res = await streamApi.startPreview(row.id)
    previewUrls.value = res.data
    previewTab.value = 'hls'
    showPreviewDialog.value = true
  } catch (error: any) {
    ElMessage.error(error.message || '获取预览地址失败')
  }
}

const handleRecord = async (row: StreamWithStatus) => {
  try {
    const statusRes = await streamApi.isRecording(row.id)
    if (statusRes.data.recording) {
      await ElMessageBox.confirm('当前正在录制，是否停止？', '录制控制', {
        confirmButtonText: '停止录制',
        cancelButtonText: '取消',
        type: 'warning'
      })
      await streamApi.stopRecord(row.id)
      ElMessage.success('录制已停止')
    } else {
      await streamApi.startRecord(row.id)
      ElMessage.success('录制已启动')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '录制操作失败')
    }
  }
}

const handlePlayback = async (row: StreamWithStatus) => {
  loadingSegments.value = true
  currentPlaybackUrl.value = ''
  showPlaybackDialog.value = true
  try {
    const res = await streamApi.listPlaybackSegments(row.id)
    playbackSegments.value = res.data
  } catch (error: any) {
    ElMessage.error(error.message || '获取录制列表失败')
  } finally {
    loadingSegments.value = false
  }
}

const playSegment = (segment: PlaybackSegment) => {
  currentPlaybackUrl.value = segment.playback_url
}

const copyUrl = (url?: string) => {
  if (!url) return
  navigator.clipboard.writeText(url)
  ElMessage.success('已复制到剪贴板')
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleString('zh-CN')
}

onMounted(() => {
  loadStreams()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.ml-1 {
  margin-left: 4px;
}

.mt-4 {
  margin-top: 16px;
}
</style>
