<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>工作流管理</span>
          <div class="header-actions">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索工作流"
              clearable
              style="width: 200px; margin-right: 10px"
              @change="loadWorkflows"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-button type="primary" @click="showCreateDialog = true">
              <el-icon><Plus /></el-icon>
              创建工作流
            </el-button>
          </div>
        </div>
      </template>

      <div class="filter-bar">
        <el-space wrap>
          <el-select v-model="filterTriggerType" placeholder="触发方式" clearable @change="loadWorkflows" style="width: 120px">
            <el-option label="手动" value="manual" />
            <el-option label="定时" value="schedule" />
            <el-option label="事件" value="event" />
          </el-select>
          <el-select v-model="filterStatus" placeholder="状态" clearable @change="loadWorkflows" style="width: 120px">
            <el-option label="草稿" value="draft" />
            <el-option label="测试中" value="testing" />
            <el-option label="已发布" value="published" />
            <el-option label="已归档" value="archived" />
          </el-select>
        </el-space>
      </div>

      <el-table :data="workflows" v-loading="loading" stripe>
        <el-table-column prop="name" label="名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="code" label="代码" width="150" />
        <el-table-column prop="trigger_type" label="触发方式" width="100">
          <template #default="{ row }">
            <el-tag :type="getTriggerTypeColor(row.trigger_type)" size="small">
              {{ getTriggerTypeLabel(row.trigger_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="版本" width="80" />
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
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleView(row)">查看</el-button>
            <el-button link type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="success" size="small" @click="handleTrigger(row)">触发</el-button>
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
          @size-change="loadWorkflows"
          @current-change="loadWorkflows"
        />
      </div>
    </el-card>

    <el-dialog v-model="showCreateDialog" title="创建工作流" width="600px">
      <el-form ref="createFormRef" :model="createForm" :rules="createRules" label-width="100px">
        <el-form-item label="工作流代码" prop="code">
          <el-input v-model="createForm.code" placeholder="唯一标识" />
        </el-form-item>
        <el-form-item label="工作流名称" prop="name">
          <el-input v-model="createForm.name" placeholder="工作流显示名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="createForm.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="触发方式" prop="trigger_type">
          <el-select v-model="createForm.trigger_type" style="width: 100%">
            <el-option label="手动触发" value="manual" />
            <el-option label="定时调度" value="schedule" />
            <el-option label="事件触发" value="event" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showViewDialog" title="工作流详情" width="700px">
      <el-descriptions v-if="currentWorkflow" :column="2" border>
        <el-descriptions-item label="ID" :span="2">{{ currentWorkflow.id }}</el-descriptions-item>
        <el-descriptions-item label="代码">{{ currentWorkflow.code }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ currentWorkflow.name }}</el-descriptions-item>
        <el-descriptions-item label="触发方式">
          <el-tag :type="getTriggerTypeColor(currentWorkflow.trigger_type)" size="small">
            {{ getTriggerTypeLabel(currentWorkflow.trigger_type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="版本">{{ currentWorkflow.version }}</el-descriptions-item>
        <el-descriptions-item label="状态" :span="2">
          <el-tag :type="getStatusColor(currentWorkflow.status)" size="small">
            {{ getStatusLabel(currentWorkflow.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentWorkflow.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ formatDate(currentWorkflow.created_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <el-dialog v-model="showTriggerDialog" title="触发工作流" width="500px">
      <el-form label-width="100px">
        <el-form-item label="关联资产">
          <el-input v-model="triggerForm.asset_id" placeholder="可选：输入资产 ID" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showTriggerDialog = false">取消</el-button>
        <el-button type="primary" :loading="triggering" @click="handleTriggerConfirm">确定触发</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { workflowApi, type Workflow, type WorkflowCreateReq } from '@/api/workflow'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const creating = ref(false)
const triggering = ref(false)
const workflows = ref<Workflow[]>([])
const showCreateDialog = ref(false)
const showViewDialog = ref(false)
const showTriggerDialog = ref(false)
const currentWorkflow = ref<Workflow | null>(null)
const createFormRef = ref<FormInstance>()

const searchKeyword = ref('')
const filterTriggerType = ref('')
const filterStatus = ref('')

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const createForm = reactive<WorkflowCreateReq>({
  code: '',
  name: '',
  description: '',
  trigger_type: 'manual'
})

const triggerForm = reactive({
  asset_id: ''
})

const createRules: FormRules = {
  code: [{ required: true, message: '请输入工作流代码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入工作流名称', trigger: 'blur' }],
  trigger_type: [{ required: true, message: '请选择触发方式', trigger: 'change' }]
}

onMounted(() => {
  loadWorkflows()
})

async function loadWorkflows() {
  loading.value = true
  try {
    const response = await workflowApi.list({
      keyword: searchKeyword.value || undefined,
      trigger_type: filterTriggerType.value as any,
      status: filterStatus.value as any,
      page: pagination.page,
      page_size: pagination.page_size
    })
    workflows.value = response.data.items
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
      await workflowApi.create(createForm)
      ElMessage.success('创建成功')
      showCreateDialog.value = false
      loadWorkflows()
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || '创建失败')
    } finally {
      creating.value = false
    }
  })
}

function handleView(row: Workflow) {
  currentWorkflow.value = row
  showViewDialog.value = true
}

function handleEdit(row: Workflow) {
  ElMessage.info('编辑功能开发中')
}

function handleTrigger(row: Workflow) {
  currentWorkflow.value = row
  triggerForm.asset_id = ''
  showTriggerDialog.value = true
}

async function handleTriggerConfirm() {
  if (!currentWorkflow.value) return
  triggering.value = true
  try {
    const assetId = triggerForm.asset_id || undefined
    await workflowApi.trigger(currentWorkflow.value.id, assetId)
    ElMessage.success('工作流已触发，任务已创建')
    showTriggerDialog.value = false
    router.push('/tasks')
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '触发失败')
  } finally {
    triggering.value = false
  }
}

async function handleEnable(row: Workflow) {
  try {
    await workflowApi.enable(row.id)
    ElMessage.success('启用成功')
    loadWorkflows()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '启用失败')
  }
}

async function handleDisable(row: Workflow) {
  try {
    await workflowApi.disable(row.id)
    ElMessage.success('禁用成功')
    loadWorkflows()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '禁用失败')
  }
}

async function handleDelete(row: Workflow) {
  try {
    await ElMessageBox.confirm('确定要删除此工作流吗？', '提示', {
      type: 'warning'
    })
    await workflowApi.delete(row.id)
    ElMessage.success('删除成功')
    loadWorkflows()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

function getTriggerTypeLabel(type: string) {
  const map: Record<string, string> = {
    manual: '手动',
    schedule: '定时',
    event: '事件'
  }
  return map[type] || type
}

function getTriggerTypeColor(type: string) {
  const map: Record<string, any> = {
    manual: 'primary',
    schedule: 'success',
    event: 'warning'
  }
  return map[type] || ''
}

function getStatusLabel(status: string) {
  const map: Record<string, string> = {
    draft: '草稿',
    testing: '测试中',
    published: '已发布',
    archived: '已归档'
  }
  return map[status] || status
}

function getStatusColor(status: string) {
  const map: Record<string, any> = {
    draft: 'info',
    testing: 'warning',
    published: 'success',
    archived: 'danger'
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
