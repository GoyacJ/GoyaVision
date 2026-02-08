<template>
  <div class="h-full flex flex-col bg-gray-50">
    <!-- Header -->
    <div class="h-14 border-b bg-white flex items-center px-4 justify-between flex-shrink-0">
      <div class="flex items-center gap-4">
        <el-button @click="$router.back()">
          <el-icon><Back /></el-icon>
        </el-button>
        <div class="flex items-center gap-3">
          <span class="font-medium text-lg">任务详情</span>
          <span class="text-gray-400 text-sm">#{{ task?.id }}</span>
          <StatusBadge v-if="task" :status="mapStatus(task.status)" />
        </div>
      </div>
      
      <div class="flex items-center gap-4" v-if="task">
        <div class="flex flex-col items-end">
          <div class="text-sm text-gray-500">进度</div>
          <div class="w-32">
            <el-progress :percentage="task.progress" :status="progressStatus" :show-text="false" />
          </div>
        </div>
      </div>
    </div>

    <div class="flex-1 flex overflow-hidden">
      <!-- DAG Canvas (Read-only) -->
      <div class="flex-1 relative border-r bg-white">
        <VueFlow
          v-if="elements.length > 0"
          id="task-detail"
          v-model="elements"
          :default-zoom="1.0"
          :min-zoom="0.2"
          :max-zoom="4"
          :nodes-draggable="false"
          :nodes-connectable="false"
          :elements-selectable="true"
          :fit-view-on-init="true"
          @node-click="onNodeClick"
        >
          <Background pattern-color="#aaa" :gap="8" />
          <Controls :show-interactive="false" />
          <MiniMap />

          <template #node-operator="props">
            <div 
              class="min-w-[180px] bg-white rounded-lg border-2 shadow-sm transition-all relative overflow-hidden"
              :class="getNodeClass(props.id)"
            >
              <!-- Progress Bar Overlay -->
              <div 
                v-if="getNodeStatus(props.id) === 'running'"
                class="absolute bottom-0 left-0 h-1 bg-blue-500 animate-pulse w-full"
              ></div>

              <Handle type="target" position="top" class="w-3 h-3 !bg-gray-300" />
              
              <div class="p-3">
                <div class="flex items-center gap-2 mb-2">
                  <div class="w-8 h-8 rounded bg-gray-50 flex items-center justify-center text-lg">
                    📦
                  </div>
                  <div class="flex-1 min-w-0">
                    <div class="font-medium text-sm truncate">{{ props.data.operatorName }}</div>
                    <div class="text-xs text-gray-400 font-mono truncate">{{ props.data.operatorCode }}</div>
                  </div>
                </div>
                
                <div class="flex justify-between items-center text-xs">
                  <span 
                    class="px-1.5 py-0.5 rounded text-white"
                    :class="getStatusBadgeClass(getNodeStatus(props.id))"
                  >
                    {{ getStatusLabel(getNodeStatus(props.id)) }}
                  </span>
                  <span v-if="getNodeDuration(props.id)" class="text-gray-400">
                    {{ getNodeDuration(props.id) }}
                  </span>
                </div>
              </div>

              <Handle type="source" position="bottom" class="w-3 h-3 !bg-gray-300" />
            </div>
          </template>

          <template #node-trigger="props">
            <div class="min-w-[120px] bg-white rounded-lg border-2 border-green-500 shadow-sm p-3">
              <div class="flex items-center gap-2">
                <div class="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center text-lg text-green-600">▶</div>
                <div class="font-bold text-gray-700">开始</div>
              </div>
              <Handle type="source" position="bottom" class="w-3 h-3 !bg-green-500" />
            </div>
          </template>
        </VueFlow>
      </div>

      <!-- Node Detail Panel -->
      <div class="w-96 bg-white flex-shrink-0 flex flex-col" v-if="selectedNode">
        <div class="p-4 border-b flex justify-between items-center">
          <h3 class="font-medium">节点详情</h3>
          <el-button link @click="selectedNode = null">
            <el-icon><Close /></el-icon>
          </el-button>
        </div>
        
        <div class="flex-1 overflow-y-auto p-4 space-y-6">
          <!-- Status -->
          <div class="space-y-2">
            <div class="text-sm text-gray-500">执行状态</div>
            <div class="flex items-center gap-2">
              <span 
                class="px-2 py-1 rounded text-sm font-medium text-white"
                :class="getStatusBadgeClass(selectedNodeExec?.status)"
              >
                {{ getStatusLabel(selectedNodeExec?.status) }}
              </span>
              <span class="text-sm text-gray-400" v-if="selectedNodeExec?.completed_at">
                完成于 {{ formatDate(selectedNodeExec.completed_at) }}
              </span>
            </div>
            <div v-if="selectedNodeExec?.error" class="bg-red-50 text-red-600 p-3 rounded text-sm mt-2">
              {{ selectedNodeExec.error }}
            </div>
          </div>

          <!-- Artifacts -->
          <div class="space-y-3" v-if="nodeArtifacts.length > 0">
            <div class="text-sm text-gray-500">产物列表</div>
            <div class="space-y-2">
              <div 
                v-for="art in nodeArtifacts" 
                :key="art.id"
                class="border rounded p-3 hover:shadow-sm transition-all"
              >
                <div class="flex justify-between items-start mb-2">
                  <div class="font-medium text-sm">{{ getArtifactTypeLabel(art.type) }}</div>
                  <div class="text-xs text-gray-400">{{ formatDate(art.created_at) }}</div>
                </div>
                
                <!-- Asset Preview -->
                <div v-if="art.type === 'asset' && art.data?.asset_info" class="flex gap-2 text-sm text-gray-600">
                  <div class="w-16 h-16 bg-gray-100 rounded flex items-center justify-center flex-shrink-0">
                    <span v-if="art.data.asset_info.type === 'image'">🖼️</span>
                    <span v-else-if="art.data.asset_info.type === 'video'">🎬</span>
                    <span v-else>📄</span>
                  </div>
                  <div class="min-w-0 flex-1">
                    <div class="truncate">{{ art.data.asset_info.path }}</div>
                    <div class="text-xs text-gray-400 mt-1">
                      {{ art.data.asset_info.format }} / {{ formatSize(art.data.asset_info.size) }}
                    </div>
                  </div>
                </div>

                <!-- Result Preview -->
                <div v-else-if="art.type === 'result'" class="bg-gray-50 p-2 rounded text-xs font-mono overflow-x-auto">
                  {{ JSON.stringify(art.data?.results || {}, null, 2).slice(0, 100) }}...
                </div>
              </div>
            </div>
          </div>
          
          <div v-else-if="selectedNodeExec?.status === 'success'" class="text-center text-gray-400 py-4 text-sm">
            该节点无产出
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { Back, Close } from '@element-plus/icons-vue'
import { VueFlow, Handle } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import { MiniMap } from '@vue-flow/minimap'
import { taskApi } from '@/api/task'
import { workflowApi } from '@/api/workflow'
import { artifactApi } from '@/api/artifact'
import { useTaskProgress } from '@/composables/useTaskProgress'
import StatusBadge from '@/components/business/StatusBadge/index.vue'
import dagre from 'dagre'

// Import CSS
import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'
import '@vue-flow/controls/dist/style.css'
import '@vue-flow/minimap/dist/style.css'

const route = useRoute()
const taskId = ref(route.params.id as string)
const task = ref<any>(null)
const elements = ref<any[]>([])
const selectedNode = ref<any>(null)
const nodeArtifacts = ref<any[]>([])

const { nodeExecutions, status, progress, connect } = useTaskProgress(taskId)

// Layout graph
function layoutGraph(nodes: any[], edges: any[]) {
  const dagreGraph = new dagre.graphlib.Graph()
  dagreGraph.setDefaultEdgeLabel(() => ({}))
  dagreGraph.setGraph({ rankdir: 'TB' })

  nodes.forEach((node) => {
    dagreGraph.setNode(node.id, { width: 180, height: 100 })
  })

  edges.forEach((edge) => {
    dagreGraph.setEdge(edge.source, edge.target)
  })

  dagre.layout(dagreGraph)

  return nodes.map((node) => {
    const nodeWithPosition = dagreGraph.node(node.id)
    return {
      ...node,
      position: {
        x: nodeWithPosition.x - 90,
        y: nodeWithPosition.y - 50,
      },
    }
  })
}

// Initial Data Load
onMounted(async () => {
  if (!taskId.value) return

  // 1. Get Task
  const taskRes = await taskApi.get(taskId.value, true)
  task.value = taskRes
  
  // 2. Get Workflow Structure
  if (taskRes.workflow_id) {
    const wf = await workflowApi.get(taskRes.workflow_id, true)
    
    const flowNodes = (wf.nodes || []).map((node: any) => ({
      id: node.node_key,
      type: node.node_type === 'trigger' ? 'trigger' : 'operator',
      data: {
        operatorName: node.operator?.name,
        operatorCode: node.operator?.code,
      }
    }))

    const flowEdges = (wf.edges || []).map((edge: any) => ({
      id: `e-${edge.source_key}-${edge.target_key}`,
      source: edge.source_key,
      target: edge.target_key,
      type: 'default', // Simple edges for read-only view
      animated: true,
    }))

    const layoutedNodes = layoutGraph(flowNodes, flowEdges)
    elements.value = [...layoutedNodes, ...flowEdges]
  }

  // 3. Connect SSE
  connect()
})

// Helpers
const selectedNodeExec = computed(() => {
  if (!selectedNode.value) return null
  return nodeExecutions.value.find(e => e.node_key === selectedNode.value.id)
})

watch(selectedNode, async (node) => {
  if (!node) {
    nodeArtifacts.value = []
    return
  }
  // Fetch artifacts for this node
  try {
    const res = await artifactApi.list({
      task_id: taskId.value,
      node_key: node.id,
      page_size: 100
    })
    nodeArtifacts.value = res.items
  } catch (e) {
    console.error('Failed to load artifacts', e)
  }
})

function onNodeClick(e: any) {
  selectedNode.value = e.node
}

function getNodeStatus(nodeKey: string) {
  const exec = nodeExecutions.value.find(e => e.node_key === nodeKey)
  return exec?.status || 'pending'
}

function getNodeDuration(nodeKey: string) {
  const exec = nodeExecutions.value.find(e => e.node_key === nodeKey)
  if (!exec?.started_at || !exec?.completed_at) return ''
  const start = new Date(exec.started_at).getTime()
  const end = new Date(exec.completed_at).getTime()
  return ((end - start) / 1000).toFixed(1) + 's'
}

function getNodeClass(nodeKey: string) {
  const status = getNodeStatus(nodeKey)
  const isSelected = selectedNode.value?.id === nodeKey
  
  let baseClass = ''
  if (isSelected) baseClass += ' ring-2 ring-blue-400'

  switch (status) {
    case 'success': return baseClass + ' border-green-500 bg-green-50'
    case 'failed': return baseClass + ' border-red-500 bg-red-50'
    case 'running': return baseClass + ' border-blue-500 bg-white'
    case 'skipped': return baseClass + ' border-gray-300 bg-gray-100 opacity-60 border-dashed'
    default: return baseClass + ' border-gray-200'
  }
}

function getStatusBadgeClass(status?: string) {
  switch (status) {
    case 'success': return 'bg-green-500'
    case 'failed': return 'bg-red-500'
    case 'running': return 'bg-blue-500'
    case 'skipped': return 'bg-gray-400'
    default: return 'bg-gray-300'
  }
}

function getStatusLabel(status?: string) {
  const map: Record<string, string> = {
    pending: '等待中',
    running: '运行中',
    success: '成功',
    failed: '失败',
    skipped: '跳过'
  }
  return map[status || 'pending'] || status
}

function getArtifactTypeLabel(type: string) {
  const map: Record<string, string> = {
    asset: '媒体资产',
    result: '分析结果',
    timeline: '时间轴',
    report: '报告'
  }
  return map[type] || type
}

function formatSize(bytes?: number) {
  if (!bytes) return '-'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

function formatDate(dateStr?: string) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleTimeString()
}

function mapStatus(status: string) {
  // reusing existing status mapping if available or simple map
  const map: Record<string, any> = {
    pending: 'pending',
    running: 'processing',
    success: 'active',
    failed: 'inactive',
    cancelled: 'neutral'
  }
  return map[status] || 'neutral'
}

const progressStatus = computed(() => {
  if (task.value?.status === 'failed') return 'exception'
  if (task.value?.status === 'success') return 'success'
  return ''
})
</script>
