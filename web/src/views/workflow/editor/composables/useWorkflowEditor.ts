import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useVueFlow } from '@vue-flow/core'
import { workflowApi } from '@/api/workflow'
import { ElMessage } from 'element-plus'

export function useWorkflowEditor() {
  const route = useRoute()
  const router = useRouter()
  const workflowId = route.params.id as string
  const workflow = ref<any>(null)
  
  const { nodes, edges, setNodes, setEdges } = useVueFlow({ id: 'workflow-editor' })

  // Load workflow data
  async function loadWorkflow() {
    if (!workflowId) return

    try {
      const res = await workflowApi.get(workflowId, true)
      const wf = res.data
      workflow.value = wf

      // Map server nodes to VueFlow nodes
      const flowNodes = (wf.nodes || []).map((node: any) => ({
        id: node.node_key,
        type: node.node_type === 'trigger' ? 'trigger' : 'operator',
        position: node.position || { x: 0, y: 0 },
        data: {
          operatorId: node.operator_id,
          operatorName: node.operator?.name,
          operatorCode: node.operator?.code,
          config: node.config || {},
        },
      }))

      // Map server edges to VueFlow edges
      const flowEdges = (wf.edges || []).map((edge: any) => ({
        id: `e-${edge.source_key}-${edge.target_key}`,
        source: edge.source_key,
        target: edge.target_key,
        type: 'conditional',
        data: {
          condition: edge.condition || { type: 'always' }
        }
      }))

      setNodes(flowNodes)
      setEdges(flowEdges)
    } catch (error) {
      console.error('Failed to load workflow:', error)
      ElMessage.error('加载工作流失败')
    }
  }

  // Save workflow data
  async function saveWorkflow() {
    if (!workflowId) return

    try {
      // Map VueFlow nodes to server nodes
      const serverNodes = nodes.value.map((node) => ({
        node_key: node.id,
        node_type: node.type === 'trigger' ? 'trigger' : 'operator',
        operator_id: node.data.operatorId,
        position: node.position,
        config: node.data.config
      }))

      // Map VueFlow edges to server edges
      const serverEdges = edges.value.map((edge) => ({
        source_key: edge.source,
        target_key: edge.target,
        condition: edge.data?.condition || { type: 'always' }
      }))

      await workflowApi.update(workflowId, {
        nodes: serverNodes,
        edges: serverEdges
      })
      ElMessage.success('保存成功')
    } catch (error) {
      console.error('Failed to save workflow:', error)
      ElMessage.error('保存失败')
    }
  }

  async function runWorkflow() {
    if (!workflowId) return
    try {
      await saveWorkflow() // Save before run
      const res = await workflowApi.trigger(workflowId)
      const task = res.data
      ElMessage.success('任务已启动')
      // Note: trigger returns the task ID or task object? Assuming object with id
      router.push(`/tasks/${task.id}`)
    } catch (error) {
      console.error('Failed to run workflow:', error)
      ElMessage.error('启动任务失败')
    }
  }

  // Initialize
  loadWorkflow()

  return {
    workflow,
    saveWorkflow,
    runWorkflow
  }
}
