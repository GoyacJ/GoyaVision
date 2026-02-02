<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>算子中心</span>
          <div class="header-actions">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索算子"
              clearable
              style="width: 200px; margin-right: 10px"
              @change="loadOperators"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-button type="primary" @click="showCreateDialog = true">
              <el-icon><Plus /></el-icon>
              添加算子
            </el-button>
          </div>
        </div>
      </template>

      <div class="filter-bar">
        <el-space wrap>
          <el-select v-model="filterCategory" placeholder="分类" clearable @change="loadOperators" style="width: 120px">
            <el-option label="分析" value="analysis" />
            <el-option label="处理" value="processing" />
            <el-option label="生成" value="generation" />
          </el-select>
          <el-select v-model="filterStatus" placeholder="状态" clearable @change="loadOperators" style="width: 120px">
            <el-option label="草稿" value="draft" />
            <el-option label="测试中" value="testing" />
            <el-option label="已发布" value="published" />
            <el-option label="已废弃" value="deprecated" />
          </el-select>
          <el-checkbox v-model="showBuiltin" @change="loadOperators">仅显示内置算子</el-checkbox>
        </el-space>
      </div>

      <el-table :data="operators" v-loading="loading" stripe>
        <el-table-column prop="name" label="名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="code" label="代码" width="150" />
        <el-table-column prop="category" label="分类" width="100">
          <template #default="{ row }">
            <el-tag :type="getCategoryColor(row.category)" size="small">
              {{ getCategoryLabel(row.category) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="120" />
        <el-table-column prop="version" label="版本" width="80" />
        <el-table-column label="内置" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.is_builtin" type="info" size="small">内置</el-tag>
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
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleView(row)">查看</el-button>
            <el-button link type="primary" size="small" @click="handleEdit(row)" :disabled="row.is_builtin">编辑</el-button>
            <el-button
              v-if="row.status === 'published'"
              link
              type="warning"
              size="small"
              @click="handleDisable(row)"
            >
              禁用
            </el-button>
            <el-button
              v-else
              link
              type="success"
              size="small"
              @click="handleEnable(row)"
            >
              启用
            </el-button>
            <el-button link type="danger" size="small" @click="handleDelete(row)" :disabled="row.is_builtin">删除</el-button>
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
          @size-change="loadOperators"
          @current-change="loadOperators"
        />
      </div>
    </el-card>

    <el-dialog v-model="showCreateDialog" title="添加算子" width="600px">
      <el-form ref="createFormRef" :model="createForm" :rules="createRules" label-width="100px">
        <el-form-item label="算子代码" prop="code">
          <el-input v-model="createForm.code" placeholder="唯一标识，如：frame_extract" />
        </el-form-item>
        <el-form-item label="算子名称" prop="name">
          <el-input v-model="createForm.name" placeholder="算子显示名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="createForm.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="createForm.category" style="width: 100%">
            <el-option label="分析" value="analysis" />
            <el-option label="处理" value="processing" />
            <el-option label="生成" value="generation" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-input v-model="createForm.type" placeholder="如：object_detection, frame_extract" />
        </el-form-item>
        <el-form-item label="端点地址" prop="endpoint">
          <el-input v-model="createForm.endpoint" placeholder="http://..." />
        </el-form-item>
        <el-form-item label="HTTP方法" prop="method">
          <el-select v-model="createForm.method" style="width: 100%">
            <el-option label="POST" value="POST" />
            <el-option label="GET" value="GET" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showViewDialog" title="算子详情" width="700px">
      <el-descriptions v-if="currentOperator" :column="2" border>
        <el-descriptions-item label="ID" :span="2">{{ currentOperator.id }}</el-descriptions-item>
        <el-descriptions-item label="代码">{{ currentOperator.code }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ currentOperator.name }}</el-descriptions-item>
        <el-descriptions-item label="分类">
          <el-tag :type="getCategoryColor(currentOperator.category)" size="small">
            {{ getCategoryLabel(currentOperator.category) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="类型">{{ currentOperator.type }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ currentOperator.version }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusColor(currentOperator.status)" size="small">
            {{ getStatusLabel(currentOperator.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="端点地址" :span="2">{{ currentOperator.endpoint }}</el-descriptions-item>
        <el-descriptions-item label="HTTP方法">{{ currentOperator.method }}</el-descriptions-item>
        <el-descriptions-item label="内置">
          <el-tag v-if="currentOperator.is_builtin" type="info" size="small">是</el-tag>
          <span v-else>否</span>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentOperator.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ formatDate(currentOperator.created_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { operatorApi, type Operator, type OperatorCreateReq } from '@/api/operator'

const loading = ref(false)
const creating = ref(false)
const operators = ref<Operator[]>([])
const showCreateDialog = ref(false)
const showViewDialog = ref(false)
const currentOperator = ref<Operator | null>(null)
const createFormRef = ref<FormInstance>()

const searchKeyword = ref('')
const filterCategory = ref('')
const filterStatus = ref('')
const showBuiltin = ref(false)

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const createForm = reactive<OperatorCreateReq>({
  code: '',
  name: '',
  description: '',
  category: 'analysis',
  type: '',
  endpoint: '',
  method: 'POST'
})

const createRules: FormRules = {
  code: [{ required: true, message: '请输入算子代码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入算子名称', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  type: [{ required: true, message: '请输入类型', trigger: 'blur' }],
  endpoint: [{ required: true, message: '请输入端点地址', trigger: 'blur' }]
}

onMounted(() => {
  loadOperators()
})

async function loadOperators() {
  loading.value = true
  try {
    const response = await operatorApi.list({
      keyword: searchKeyword.value || undefined,
      category: filterCategory.value as any,
      status: filterStatus.value as any,
      is_builtin: showBuiltin.value || undefined,
      page: pagination.page,
      page_size: pagination.page_size
    })
    operators.value = response.data.items
    pagination.total = response.data.total
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function handleCreate() {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (!valid) return
    creating.value = true
    try {
      await operatorApi.create(createForm)
      ElMessage.success('创建成功')
      showCreateDialog.value = false
      loadOperators()
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || '创建失败')
    } finally {
      creating.value = false
    }
  })
}

function handleView(row: Operator) {
  currentOperator.value = row
  showViewDialog.value = true
}

function handleEdit(row: Operator) {
  ElMessage.info('编辑功能开发中')
}

async function handleEnable(row: Operator) {
  try {
    await operatorApi.enable(row.id)
    ElMessage.success('启用成功')
    loadOperators()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '启用失败')
  }
}

async function handleDisable(row: Operator) {
  try {
    await operatorApi.disable(row.id)
    ElMessage.success('禁用成功')
    loadOperators()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '禁用失败')
  }
}

async function handleDelete(row: Operator) {
  try {
    await ElMessageBox.confirm('确定要删除此算子吗？', '提示', {
      type: 'warning'
    })
    await operatorApi.delete(row.id)
    ElMessage.success('删除成功')
    loadOperators()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

function getCategoryLabel(category: string) {
  const map: Record<string, string> = {
    analysis: '分析',
    processing: '处理',
    generation: '生成'
  }
  return map[category] || category
}

function getCategoryColor(category: string) {
  const map: Record<string, any> = {
    analysis: 'primary',
    processing: 'success',
    generation: 'warning'
  }
  return map[category] || ''
}

function getStatusLabel(status: string) {
  const map: Record<string, string> = {
    draft: '草稿',
    testing: '测试中',
    published: '已发布',
    deprecated: '已废弃'
  }
  return map[status] || status
}

function getStatusColor(status: string) {
  const map: Record<string, any> = {
    draft: 'info',
    testing: 'warning',
    published: 'success',
    deprecated: 'danger'
  }
  return map[status] || ''
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
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  align-items: center;
}

.filter-bar {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
