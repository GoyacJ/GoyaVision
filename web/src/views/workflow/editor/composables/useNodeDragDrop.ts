import { useVueFlow } from '@vue-flow/core'

function generateId() {
  return Math.random().toString(36).slice(2, 10)
}

export function useNodeDragDrop() {
  const { addNodes, project, vueFlowRef } = useVueFlow({ id: 'workflow-editor' })

  function onDragOver(event: DragEvent) {
    event.preventDefault()
    if (event.dataTransfer) {
      event.dataTransfer.dropEffect = 'copy'
    }
  }

  function onDrop(event: DragEvent) {
    event.preventDefault()

    // Try both standard and custom types for maximum compatibility
    const data = event.dataTransfer?.getData('text') || event.dataTransfer?.getData('application/vueflow')
    if (!data) {
      console.warn('Drop event triggered but no data found in dataTransfer')
      return
    }

    let nodeData
    try {
      nodeData = JSON.parse(data)
    } catch (e) {
      console.error('Failed to parse drag data', e)
      return
    }

    // Calculate position
    // Since vueFlowRef is now bound to the container div, vueFlowRef.value is the element
    const target = (vueFlowRef.value as any)?.$el || vueFlowRef.value
    if (!target || typeof target.getBoundingClientRect !== 'function') {
      console.warn('VueFlow target element not found or invalid for drop')
      return
    }

    const { left, top } = target.getBoundingClientRect()
    const position = project({
      x: event.clientX - left,
      y: event.clientY - top,
    })

    const newNode = {
      id: `node-${generateId()}`,
      type: nodeData.type === 'trigger' ? 'trigger' : 'operator',
      position,
      data: {
        operatorId: nodeData.operatorId,
        operatorCode: nodeData.operatorCode,
        operatorName: nodeData.operatorName,
        config: {
          params: {},
          retry_count: 0,
          timeout_seconds: 0
        }
      },
    }

    addNodes([newNode])
  }

  return {
    onDragOver,
    onDrop
  }
}
