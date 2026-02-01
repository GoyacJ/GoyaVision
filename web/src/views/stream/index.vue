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
        <el-table-column prop="name" label="名称" width="200" />
        <el-table-column prop="url" label="RTSP URL" min-width="300" />
        <el-table-column prop="enabled" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="primary" size="small" @click="handlePreview(row)">预览</el-button>
            <el-button link type="primary" size="small" @click="handleRecord(row)">录制</el-button>
            <el-button link type="primary" size="small" @click="handleBindings(row)">绑定</el-button>
            <el-button link type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="showCreateDialog" title="添加流" width="500px">
      <el-form ref="createFormRef" :model="createForm" :rules="rules" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入流名称" />
        </el-form-item>
        <el-form-item label="RTSP URL" prop="url">
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

    <el-dialog v-model="showEditDialog" title="编辑流" width="500px">
      <el-form ref="editFormRef" :model="editForm" :rules="rules" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="RTSP URL" prop="url">
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

    <el-dialog v-model="showPreviewDialog" title="视频预览" width="800px" destroy-on-close>
      <HLSPreview v-if="previewUrl" :hls-url="previewUrl" :width="720" :height="405" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { streamApi, type Stream, type StreamCreateReq, type StreamUpdateReq } from '../../api/stream'
import HLSPreview from '../../components/HLSPreview.vue'

const streams = ref<Stream[]>([])
const loading = ref(false)
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const showPreviewDialog = ref(false)
const previewUrl = ref('')
const currentEditId = ref('')
const createFormRef = ref<FormInstance>()
const editFormRef = ref<FormInstance>()

const createForm = ref<StreamCreateReq>({
  name: '',
  url: '',
  enabled: true
})

const editForm = ref<StreamUpdateReq>({
  name: '',
  url: '',
  enabled: true
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入流名称', trigger: 'blur' }],
  url: [{ required: true, message: '请输入RTSP URL', trigger: 'blur' }]
}

const loadStreams = async () => {
  loading.value = true
  try {
    const res = await streamApi.list()
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
    createForm.value = { name: '', url: '', enabled: true }
    loadStreams()
  } catch (error: any) {
    if (error !== 'cancel' && error?.message) {
      ElMessage.error(error.message || '创建失败')
    }
  }
}

const handleEdit = (row: Stream) => {
  currentEditId.value = row.id
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

const handleDelete = async (row: Stream) => {
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

const handlePreview = async (row: Stream) => {
  try {
    const res = await streamApi.startPreview(row.id)
    previewUrl.value = res.data.hls_url
    showPreviewDialog.value = true
  } catch (error: any) {
    ElMessage.error(error.message || '启动预览失败')
  }
}

const handleRecord = async (row: Stream) => {
  try {
    await streamApi.startRecord(row.id)
    ElMessage.success('录制已启动')
  } catch (error: any) {
    ElMessage.error(error.message || '启动录制失败')
  }
}

const handleBindings = (row: Stream) => {
  ElMessage.info('算法绑定功能开发中')
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
</style>
