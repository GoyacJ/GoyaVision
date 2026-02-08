<template>
  <GvContainer max-width="full">
    <PageHeader
      title="Agent 会话"
      description="查看会话状态，手动执行 Run Step，并聚合会话决策与跨任务 DAG 运行链路"
    >
      <template #actions>
        <GvButton @click="refreshTable">
          <template #icon>
            <el-icon><Refresh /></el-icon>
          </template>
          刷新
        </GvButton>
      </template>
    </PageHeader>

    <FilterBar
      v-model="filters"
      :fields="filterFields"
      :loading="loading"
      @filter="() => { pagination.page = 1 }"
      @reset="handleResetFilter"
    />

    <ErrorState
      v-if="error && !loading"
      :error="error"
      title="加载失败"
      @retry="refreshTable"
    />

    <EmptyState
      v-else-if="!loading && sessions.length === 0"
      icon="🤖"
      title="暂无会话"
      description="创建 Agent Session 后，这里会显示会话运行状态"
    />

    <GvTable
      v-else
      :data="sessions"
      :columns="columns"
      :loading="loading"
      border
      stripe
      pagination
      :pagination-config="paginationConfig"
      @current-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <template #status="{ row }">
        <StatusBadge :status="mapStatus(row.status)" />
      </template>

      <template #started_at="{ row }">
        {{ formatDateTime(row.started_at) }}
      </template>

      <template #ended_at="{ row }">
        {{ formatDateTime(row.ended_at) }}
      </template>

      <template #actions="{ row }">
        <GvSpace size="xs" wrap>
          <GvButton size="small" variant="tonal" @click="handleView(row)">链路</GvButton>
          <GvButton
            size="small"
            :loading="runningSessionID === row.id"
            :disabled="row.status !== 'running'"
            @click="handleRunStep(row.id, 1)"
          >
            Run x1
          </GvButton>
          <GvButton
            size="small"
            variant="tonal"
            :loading="runningSessionID === row.id"
            :disabled="row.status !== 'running'"
            @click="handleRunStep(row.id, 3)"
          >
            Run x3
          </GvButton>
          <GvButton
            size="small"
            variant="text"
            :loading="stoppingSessionID === row.id"
            :disabled="row.status !== 'running'"
            @click="handleStopSession(row.id)"
          >
            停止
          </GvButton>
        </GvSpace>
      </template>
    </GvTable>

    <GvModal
      v-model="showDetailDialog"
      title="会话恢复链路"
      size="large"
      :show-confirm="false"
      cancel-text="关闭"
    >
      <template v-if="selectedSession">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="Session ID" :span="2">{{ selectedSession.id }}</el-descriptions-item>
          <el-descriptions-item label="Task ID">{{ selectedSession.task_id }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <StatusBadge :status="mapStatus(selectedSession.status)" />
          </el-descriptions-item>
          <el-descriptions-item label="Step Count">{{ selectedSession.step_count }}</el-descriptions-item>
          <el-descriptions-item label="开始时间">{{ formatDateTime(selectedSession.started_at) }}</el-descriptions-item>
          <el-descriptions-item label="结束时间">{{ formatDateTime(selectedSession.ended_at) }}</el-descriptions-item>
        </el-descriptions>

        <div class="mt-4 rounded border border-neutral-200 p-3">
          <div class="mb-3 flex items-center gap-3">
            <span class="text-sm text-text-tertiary">手动执行</span>
            <el-input-number v-model="detailRunMaxActions" :min="1" :max="10" size="small" />
            <GvButton
              size="small"
              :loading="runningSessionID === selectedSession.id"
              :disabled="selectedSession.status !== 'running'"
              @click="handleRunStep(selectedSession.id, detailRunMaxActions, true)"
            >
              执行 Run Step
            </GvButton>
            <GvButton
              size="small"
              variant="tonal"
              :loading="stoppingSessionID === selectedSession.id"
              :disabled="selectedSession.status !== 'running'"
              @click="handleStopSession(selectedSession.id, true)"
            >
              停止会话
            </GvButton>
            <GvButton size="small" variant="text" :loading="detailLoading" @click="loadDetail(selectedSession.id)">刷新链路</GvButton>
          </div>

          <el-timeline v-if="chainEvents.length > 0" class="max-h-[440px] overflow-y-auto pr-2">
            <el-timeline-item
              v-for="event in chainEvents"
              :key="event.id"
              :timestamp="formatDateTime(event.created_at)"
              :type="eventTimelineType(event.event_type)"
              placement="top"
            >
              <div class="rounded border border-neutral-200 p-3">
                <div class="flex items-center justify-between gap-2">
                  <div class="font-medium text-sm">{{ event.event_type }}</div>
                  <div class="text-xs text-text-tertiary">
                    {{ event.source }}
                    <span class="ml-2 rounded bg-neutral-100 px-1.5 py-0.5">{{ eventSourceLabel(event.source) }}</span>
                  </div>
                </div>
                <div class="mt-1 text-xs text-text-tertiary">
                  Task: {{ event.task_id }}
                  <span v-if="event.node_key"> | Node: {{ event.node_key }}</span>
                  <span v-if="event.tool_name"> | Tool: {{ event.tool_name }}</span>
                </div>
                <pre class="mt-2 max-h-44 overflow-auto rounded bg-gray-50 p-2 text-xs leading-5">{{ formatPrettyJSON(event.payload || {}) }}</pre>
              </div>
            </el-timeline-item>
          </el-timeline>
          <div v-else class="py-10 text-center text-sm text-text-tertiary">暂无会话事件</div>
        </div>
      </template>
    </GvModal>
  </GvContainer>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { agentApi, type AgentSession, type AgentSessionEvent, type AgentSessionStatus } from '@/api/agent'
import { taskApi, type TaskRunEvent } from '@/api/task'
import { useTable } from '@/composables'
import GvContainer from '@/components/layout/GvContainer/index.vue'
import GvTable from '@/components/base/GvTable/index.vue'
import GvModal from '@/components/base/GvModal/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvSpace from '@/components/layout/GvSpace/index.vue'
import PageHeader from '@/components/business/PageHeader/index.vue'
import FilterBar from '@/components/business/FilterBar/index.vue'
import StatusBadge from '@/components/business/StatusBadge/index.vue'
import { ErrorState, EmptyState } from '@/components/common'
import type { TableColumn } from '@/components/base/GvTable/types'
import type { FilterField } from '@/components/business/FilterBar/types'

const filters = ref({
  status: '',
  task_id: '',
})

const runningSessionID = ref('')
const stoppingSessionID = ref('')

const showDetailDialog = ref(false)
const detailLoading = ref(false)
const selectedSession = ref<AgentSession | null>(null)
const sessionEvents = ref<AgentSessionEvent[]>([])
const dagTaskEvents = ref<TaskRunEvent[]>([])
const detailRunMaxActions = ref(1)

const chainEvents = computed(() => {
  const merged = [...sessionEvents.value, ...dagTaskEvents.value]
  const deduped = new Map<string, AgentSessionEvent | TaskRunEvent>()
  for (const event of merged) {
    if (!event?.id) continue
    if (!deduped.has(event.id)) {
      deduped.set(event.id, event)
    }
  }
  return Array.from(deduped.values()).sort((a, b) => {
    const at = new Date(a.created_at || '').getTime()
    const bt = new Date(b.created_at || '').getTime()
    return at - bt
  })
})

const filterParams = computed(() => ({
  status: (filters.value.status || undefined) as AgentSessionStatus | undefined,
  task_id: filters.value.task_id?.trim() || undefined,
}))

const {
  items: sessions,
  isLoading: loading,
  error,
  pagination,
  goToPage,
  changePageSize,
  refreshTable,
} = useTable(
  async (params) => {
    const res = await agentApi.listSessions(params)
    const payload = unwrapPayload<any>(res)
    return {
      items: payload?.items || [],
      total: payload?.total || 0,
    }
  },
  {
    immediate: true,
    initialPageSize: 20,
    extraParams: filterParams,
  }
)

const columns: TableColumn[] = [
  { prop: 'id', label: 'Session ID', minWidth: '230', showOverflowTooltip: true },
  { prop: 'task_id', label: 'Task ID', minWidth: '230', showOverflowTooltip: true },
  { prop: 'status', label: '状态', width: '120' },
  { prop: 'step_count', label: 'Step', width: '90' },
  { prop: 'started_at', label: '开始时间', width: '180' },
  { prop: 'ended_at', label: '结束时间', width: '180' },
  { prop: 'actions', label: '操作', minWidth: '290', fixed: 'right' },
]

const filterFields: FilterField[] = [
  {
    key: 'status',
    label: '状态',
    type: 'select',
    placeholder: '全部状态',
    options: [
      { label: '运行中', value: 'running' },
      { label: '成功', value: 'succeeded' },
      { label: '失败', value: 'failed' },
      { label: '已取消', value: 'cancelled' },
    ],
  },
  {
    key: 'task_id',
    label: 'Task ID',
    type: 'input',
    placeholder: '按 Task ID 过滤',
  },
]

const paginationConfig = computed(() => ({
  currentPage: pagination.page,
  pageSize: pagination.pageSize,
  total: pagination.total,
}))

async function handleView(row: AgentSession) {
  showDetailDialog.value = true
  await loadDetail(row.id)
}

async function loadDetail(sessionID: string) {
  detailLoading.value = true
  sessionEvents.value = []
  dagTaskEvents.value = []
  try {
    const [sessionRes, eventRes] = await Promise.all([
      agentApi.getSession(sessionID),
      agentApi.listSessionEvents(sessionID, { page: 1, page_size: 200 }),
    ])

    const session = unwrapPayload<AgentSession>(sessionRes)
    const eventPayload = unwrapPayload<any>(eventRes)
    const sessionEventItems = eventPayload?.items || []

    selectedSession.value = session
    sessionEvents.value = sessionEventItems

    const relatedTaskIDs = collectRelatedTaskIDs(session, sessionEventItems)
    if (relatedTaskIDs.length === 0) {
      dagTaskEvents.value = []
      return
    }

    const taskEventResponses = await Promise.all(
      relatedTaskIDs.map((taskID) =>
        taskApi.listEvents(taskID, { source: 'dag_engine', limit: 200, offset: 0 }).catch(() => null)
      )
    )

    const mergedTaskEvents: TaskRunEvent[] = []
    for (const response of taskEventResponses) {
      if (!response) continue
      const payload = unwrapPayload<any>(response)
      const items = payload?.items || []
      for (const event of items) {
        mergedTaskEvents.push(event)
      }
    }
    dagTaskEvents.value = dedupeEventsByID(mergedTaskEvents)
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '加载会话链路失败')
  } finally {
    detailLoading.value = false
  }
}

async function handleRunStep(sessionID: string, maxActions: number, refreshDetail = false) {
  runningSessionID.value = sessionID
  try {
    await agentApi.runSession(sessionID, { max_actions: maxActions })
    ElMessage.success(`Run Step 执行完成（max_actions=${maxActions}）`)
    await refreshTable()
    if (refreshDetail || selectedSession.value?.id === sessionID) {
      await loadDetail(sessionID)
    }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || 'Run Step 执行失败')
  } finally {
    runningSessionID.value = ''
  }
}

async function handleStopSession(sessionID: string, refreshDetail = false) {
  stoppingSessionID.value = sessionID
  try {
    await agentApi.stopSession(sessionID)
    ElMessage.success('会话已停止')
    await refreshTable()
    if (refreshDetail || selectedSession.value?.id === sessionID) {
      await loadDetail(sessionID)
    }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '停止会话失败')
  } finally {
    stoppingSessionID.value = ''
  }
}

const handlePageChange = goToPage
const handleSizeChange = changePageSize

function handleResetFilter() {
  filters.value.status = ''
  filters.value.task_id = ''
  pagination.page = 1
}

function mapStatus(status: string): string {
  const map: Record<string, string> = {
    running: 'processing',
    succeeded: 'active',
    failed: 'inactive',
    cancelled: 'neutral',
  }
  return map[status] || 'neutral'
}

function eventTimelineType(eventType: string): 'primary' | 'success' | 'warning' | 'danger' | 'info' {
  const v = (eventType || '').toLowerCase()
  if (v.includes('failed') || v.includes('escalation')) return 'danger'
  if (v.includes('recover') || v.includes('decision')) return 'warning'
  if (v.includes('action')) return 'primary'
  if (v.includes('succeeded')) return 'success'
  return 'info'
}

function eventSourceLabel(source: string): string {
  if (source === 'agent_run_loop') return 'Agent'
  if (source === 'dag_engine') return 'DAG'
  return 'Other'
}

function formatDateTime(value?: string) {
  if (!value) return '-'
  return new Date(value).toLocaleString()
}

function formatPrettyJSON(value: any) {
  try {
    return JSON.stringify(value || {}, null, 2)
  } catch {
    return '{}'
  }
}

function unwrapPayload<T>(response: any): T {
  if (!response || typeof response !== 'object') return response as T
  if ('data' in response) {
    const level1 = response.data
    if (level1 && typeof level1 === 'object' && 'data' in level1) {
      return level1.data as T
    }
    return level1 as T
  }
  return response as T
}

function collectRelatedTaskIDs(session: AgentSession, events: AgentSessionEvent[]): string[] {
  const ids = new Set<string>()
  appendIfUUID(ids, session?.task_id)

  for (const event of events || []) {
    appendIfUUID(ids, event.task_id)
    extractTaskIDsFromPayload(ids, event.payload)
  }

  return Array.from(ids)
}

function extractTaskIDsFromPayload(target: Set<string>, payload: Record<string, any> | undefined) {
  if (!payload) return

  appendTaskIDValue(target, payload.task_id)
  appendTaskIDValue(target, payload.from_task_id)
  appendTaskIDValue(target, payload.retry_task_id)
  appendTaskIDValue(target, payload.original_task_id)
  appendTaskIDValue(target, payload.current_task_id)
}

function appendTaskIDValue(target: Set<string>, value: any) {
  if (Array.isArray(value)) {
    for (const item of value) {
      appendIfUUID(target, item)
    }
    return
  }
  appendIfUUID(target, value)
}

function appendIfUUID(target: Set<string>, value: any) {
  if (typeof value !== 'string') return
  const text = value.trim()
  if (isUUID(text)) {
    target.add(text)
  }
}

function isUUID(value: string): boolean {
  return /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i.test(value)
}

function dedupeEventsByID<T extends { id: string }>(events: T[]): T[] {
  const seen = new Set<string>()
  const out: T[] = []
  for (const event of events) {
    if (!event?.id || seen.has(event.id)) continue
    seen.add(event.id)
    out.push(event)
  }
  return out
}
</script>
