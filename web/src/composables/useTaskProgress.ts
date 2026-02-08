import { ref, onMounted, onUnmounted, type Ref } from 'vue'

export interface NodeExecution {
  node_key: string
  status: 'pending' | 'running' | 'success' | 'failed' | 'skipped'
  error?: string
  started_at?: string
  completed_at?: string
  artifact_ids?: string[]
}

export function useTaskProgress(taskId: Ref<string>) {
  const nodeExecutions = ref<NodeExecution[]>([])
  const status = ref<string>('')
  const progress = ref(0)
  const error = ref<string>('')
  let eventSource: EventSource | null = null

  function connect() {
    if (!taskId.value) return

    // Close existing connection if any
    disconnect()

    eventSource = new EventSource(`/api/v1/tasks/${taskId.value}/progress/stream`)
    
    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        nodeExecutions.value = data.node_executions || []
        status.value = data.status
        progress.value = data.progress
        error.value = data.error || ''

        if (['success', 'failed', 'cancelled'].includes(data.status)) {
          disconnect()
        }
      } catch (e) {
        console.error('Failed to parse SSE data', e)
      }
    }

    eventSource.onerror = (e) => {
      console.error('SSE Error', e)
      // Attempt reconnection logic could go here, or handled by browser default
      // For now, if error occurs and task is not done, we might want to just let it try reconnecting
      if (eventSource && eventSource.readyState === EventSource.CLOSED) {
         // Connection closed
      }
    }
  }

  function disconnect() {
    if (eventSource) {
      eventSource.close()
      eventSource = null
    }
  }

  onMounted(connect)
  onUnmounted(disconnect)

  return { nodeExecutions, status, progress, error, connect, disconnect }
}
