<template>
  <div class="algorithm-list">
    <div class="header">
      <h2>算法管理</h2>
      <el-button type="primary" @click="showCreateDialog = true">添加算法</el-button>
    </div>

    <el-table :data="algorithms" v-loading="loading" style="width: 100%">
      <el-table-column prop="name" label="名称" width="200" />
      <el-table-column prop="endpoint" label="端点" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showCreateDialog" title="添加算法" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="名称" required>
          <el-input v-model="createForm.name" />
        </el-form-item>
        <el-form-item label="端点" required>
          <el-input v-model="createForm.endpoint" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showEditDialog" title="编辑算法" width="500px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="端点">
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { algorithmApi, type Algorithm, type AlgorithmCreateReq, type AlgorithmUpdateReq } from '../api/algorithm'

const algorithms = ref<Algorithm[]>([])
const loading = ref(false)
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const currentEditId = ref('')

const createForm = ref<AlgorithmCreateReq>({
  name: '',
  endpoint: ''
})

const editForm = ref<AlgorithmUpdateReq>({
  name: '',
  endpoint: ''
})

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
    await algorithmApi.create(createForm.value)
    ElMessage.success('创建成功')
    showCreateDialog.value = false
    createForm.value = { name: '', endpoint: '' }
    loadAlgorithms()
  } catch (error: any) {
    ElMessage.error(error.message || '创建失败')
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
    await algorithmApi.update(currentEditId.value, editForm.value)
    ElMessage.success('更新成功')
    showEditDialog.value = false
    loadAlgorithms()
  } catch (error: any) {
    ElMessage.error(error.message || '更新失败')
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
.algorithm-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h2 {
  margin: 0;
}
</style>
