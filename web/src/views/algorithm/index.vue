<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>算法管理</span>
          <el-button type="primary" @click="showCreateDialog = true">
            <el-icon><Plus /></el-icon>
            添加算法
          </el-button>
        </div>
      </template>

      <el-table :data="algorithms" v-loading="loading" stripe>
        <el-table-column prop="name" label="名称" width="200" />
        <el-table-column prop="endpoint" label="端点" min-width="300" />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="showCreateDialog" title="添加算法" width="500px">
      <el-form ref="createFormRef" :model="createForm" :rules="rules" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入算法名称" />
        </el-form-item>
        <el-form-item label="端点" prop="endpoint">
          <el-input v-model="createForm.endpoint" placeholder="http://..." />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showEditDialog" title="编辑算法" width="500px">
      <el-form ref="editFormRef" :model="editForm" :rules="rules" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="端点" prop="endpoint">
          <el-input v-model="editForm.endpoint" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleUpdate">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { algorithmApi, type Algorithm, type AlgorithmCreateReq, type AlgorithmUpdateReq } from '../../api/algorithm'

const algorithms = ref<Algorithm[]>([])
const loading = ref(false)
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const currentEditId = ref('')
const createFormRef = ref<FormInstance>()
const editFormRef = ref<FormInstance>()

const createForm = ref<AlgorithmCreateReq>({
  name: '',
  endpoint: ''
})

const editForm = ref<AlgorithmUpdateReq>({
  name: '',
  endpoint: ''
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入算法名称', trigger: 'blur' }],
  endpoint: [{ required: true, message: '请输入端点地址', trigger: 'blur' }]
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString('zh-CN')
}

const loadAlgorithms = async () => {
  loading.value = true
  try {
    const res = await algorithmApi.list()
    algorithms.value = res.data
  } catch (error: any) {
    ElMessage.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

const handleCreate = async () => {
  try {
    await createFormRef.value?.validate()
    await algorithmApi.create(createForm.value)
    ElMessage.success('创建成功')
    showCreateDialog.value = false
    createForm.value = { name: '', endpoint: '' }
    loadAlgorithms()
  } catch (error: any) {
    if (error !== 'cancel' && error?.message) {
      ElMessage.error(error.message || '创建失败')
    }
  }
}

const handleEdit = (row: Algorithm) => {
  currentEditId.value = row.id
  editForm.value = {
    name: row.name,
    endpoint: row.endpoint
  }
  showEditDialog.value = true
}

const handleUpdate = async () => {
  try {
    await editFormRef.value?.validate()
    await algorithmApi.update(currentEditId.value, editForm.value)
    ElMessage.success('更新成功')
    showEditDialog.value = false
    loadAlgorithms()
  } catch (error: any) {
    if (error !== 'cancel' && error?.message) {
      ElMessage.error(error.message || '更新失败')
    }
  }
}

const handleDelete = async (row: Algorithm) => {
  try {
    await ElMessageBox.confirm('确定要删除这个算法吗？', '确认删除', {
      type: 'warning'
    })
    await algorithmApi.delete(row.id)
    ElMessage.success('删除成功')
    loadAlgorithms()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

onMounted(() => {
  loadAlgorithms()
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
