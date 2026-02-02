<template>
  <div class="page-container">
    <el-row :gutter="20" class="stats-row">
      <el-col :span="4">
        <el-card class="stats-card">
          <div class="stat-item">
            <div class="stat-label">总任务数</div>
            <div class="stat-value">{{ stats.total }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stats-card pending">
          <div class="stat-item">
            <div class="stat-label">待执行</div>
            <div class="stat-value">{{ stats.pending }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stats-card running">
          <div class="stat-item">
            <div class="stat-label">执行中</div>
            <div class="stat-value">{{ stats.running }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stats-card success">
          <div class="stat-item">
            <div class="stat-label">已成功</div>
            <div class="stat-value">{{ stats.success }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stats-card failed">
          <div class="stat-item">
            <div class="stat-label">已失败</div>
            <div class="stat-value">{{ stats.failed }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stats-card cancelled">
          <div class="stat-item">
            <div class="stat-label">已取消</div>
            <div class="stat-value">{{ stats.cancelled }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="mt-4">
      <template #header>
        <div class="card-header">
          <span>任务列表</span>
          <el-button type="primary" @click="loadTasks">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <div class="filter-bar">
        <el-space wrap>
          <el-select v-model="filterStatus" placeholder="状态" clearable @change="loadTasks" style="width: 120px">
            <el-option label="待执行" value="pending" />
            <el-option label="执行中" value="running" />
            <el-option label="已成功" value="success" />
            <el-option label="已失败" value="failed" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
        </el-space>
      </div>

      <el-table :data="tasks" v-loading="loading" stripe>
        <el-table-column prop="id" label="任务 ID" width="280" show-overflow-tooltip />
        <el-table-column label="工作流" min-width="150">
          <template #default="{ row }">
            {{ row.workflow?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusColor(row.status)" size="small">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="进度" width="120">
          <template #default="{ row }">
            <el-progress :percentage="row.progress" :status="getProgressStatus(row.status)" />
          </template>
        </el-table-column>
        <el-table-column prop="current_node" label="当前节点" width="120" show-overflow-tooltip />
        <el-table-column label="开始时间" width="160">
          <template #default="{ row }">
            {{ row.started_at ? formatDate(row.started_at) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="完成时间" width="160">
          <template #default="{ row }">
            {{ row.completed_at ? formatDate(row.completed_at) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="耗时" width="100">
          <template #default="{ row }">
            {{ calculateDuration(row) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleView(row)">查看</el-button>
            <el-button link type="success" size="small" @click="handleViewArtifacts(row)">产物</el-button>
            <el-button
              v-if="row.status === 'running'"
              link
              type="warning"
              size="small"
              @click="handleCancel(row)"
            >
              取消
            </el-button>
            <el-button
              v-if="row.status !== 'running' && row.status !== 'pending'"
              link
              type="danger"
              size="small"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
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
          @size-change="loadTasks"
          @current-change="loadTasks"
        />
      </div>
    </el-card>

    <el-dialog v-model="showViewDialog" title="任务详情" width="700px">
      <el-descriptions v-if="currentTask" :column="2" border>
        <el-descriptions-item label="任务 ID" :span="2">{{ currentTask.id }}</el-descriptions-item>
        <el-descriptions-item label="工作流">{{ currentTask.workflow?.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusColor(currentTask.status)" size="small">
            {{ getStatusLabel(currentTask.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="进度">{{ currentTask.progress }}%</el-descriptions-item>
        <el-descriptions-item label="当前节点">{{ currentTask.current_node || '-' }}</el-descriptions-item>
        <el-descriptions-item label="开始时间" :span="2">{{ currentTask.started_at ? formatDate(currentTask.started_at) : '-' }}</el-descriptions-item>
        <el-descriptions-item label="完成时间" :span="2">{{ currentTask.completed_at ? formatDate(currentTask.completed_at) : '-' }}</el-descriptions-item>
        <el-descriptions-item label="耗时" :span="2">{{ calculateDuration(currentTask) }}</el-descriptions-item>
        <el-descriptions-item v-if="currentTask.error" label="错误信息" :span="2">
          <el-alert type="error" :closable="false">{{ currentTask.error }}</el-alert>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ formatDate(currentTask.created_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { taskApi, type Task, type TaskStats } from '@/api/task'

const loading = ref(false)
const tasks = ref<Task[]>([])
const showViewDialog = ref(false)
const currentTask = ref<Task | null>(null)

const filterStatus = ref('')

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const stats = reactive<TaskStats>({
  total: 0,
  pending: 0,
  running: 0,
  success: 0,
  failed: 0,
  cancelled: 0
})

onMounted(() => {
  loadTasks()
  loadStats()
})

async function loadTasks() {
  loading.value = true
  try {
    const response = await taskApi.list({
      status: filterStatus.value as any,
      page: pagination.page,
      page_size: pagination.page_size
    })
    tasks.value = response.data.items
    pagination.total = response.data.total
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function loadStats() {
  try {
    const response = await taskApi.getStats()
    Object.assign(stats, response.data)
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '加载统计失败')
  }
}

function handleView(row: Task) {
  currentTask.value = row
  showViewDialog.value = true
}

function handleViewArtifacts(row: Task) {
  ElMessage.info('产物查看功能开发中')
}

async function handleCancel(row: Task) {
  try {
    await ElMessageBox.confirm('确定要取消此任务吗？', '提示', {
      type: 'warning'
    })
    await taskApi.cancel(row.id)
    ElMessage.success('取消成功')
    loadTasks()
    loadStats()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '取消失败')
    }
  }
}

async function handleDelete(row: Task) {
  try {
    await ElMessageBox.confirm('确定要删除此任务吗？', '提示', {
      type: 'warning'
    })
    await taskApi.delete(row.id)
    ElMessage.success('删除成功')
    loadTasks()
    loadStats()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

function getStatusLabel(status: string) {
  const map: Record<string, string> = {
    pending: '待执行',
    running: '执行中',
    success: '已成功',
    failed: '已失败',
    cancelled: '已取消'
  }
  return map[status] || status
}

function getStatusColor(status: string) {
  const map: Record<string, any> = {
    pending: 'info',
    running: 'primary',
    success: 'success',
    failed: 'danger',
    cancelled: 'warning'
  }
  return map[status] || ''
}

function getProgressStatus(status: string) {
  if (status === 'success') return 'success'
  if (status === 'failed') return 'exception'
  return undefined
}

function calculateDuration(task: Task): string {
  if (!task.started_at) return '-'
  const start = new Date(task.started_at).getTime()
  const end = task.completed_at ? new Date(task.completed_at).getTime() : Date.now()
  const duration = Math.floor((end - start) / 1000)
  
  if (duration < 60) return `${duration}秒`
  if (duration < 3600) return `${Math.floor(duration / 60)}分${duration % 60}秒`
  return `${Math.floor(duration / 3600)}时${Math.floor((duration % 3600) / 60)}分`
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

.stats-row {
  margin-bottom: 20px;
}

.stats-card {
  text-align: center;
  border-left: 4px solid #409EFF;
}

.stats-card.pending {
  border-left-color: #909399;
}

.stats-card.running {
  border-left-color: #409EFF;
}

.stats-card.success {
  border-left-color: #67C23A;
}

.stats-card.failed {
  border-left-color: #F56C6C;
}

.stats-card.cancelled {
  border-left-color: #E6A23C;
}

.stat-item {
  padding: 10px 0;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.card-header {
  display: flex;
  justify-content: space-between;
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

.mt-4 {
  margin-top: 20px;
}
</style>
