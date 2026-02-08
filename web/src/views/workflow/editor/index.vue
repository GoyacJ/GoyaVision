<template>
  <div class="h-full flex flex-col">
    <WorkflowToolbar />
    <div class="flex-1 flex overflow-hidden">
      <OperatorPalette class="w-60 border-r bg-white flex-shrink-0" />
      <div class="flex-1 relative">
        <WorkflowCanvas />
      </div>
      <div class="w-80 border-l bg-white flex-shrink-0">
        <NodeInspector v-if="selectedNode" :node="selectedNode" />
        <EdgeInspector v-else-if="selectedEdge" :edge="selectedEdge" />
        <div v-else class="p-4 text-gray-400 text-center mt-10">
          请选择节点或连线进行配置
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useVueFlow } from '@vue-flow/core'
import WorkflowToolbar from './components/WorkflowToolbar.vue'
import OperatorPalette from './components/OperatorPalette.vue'
import WorkflowCanvas from './WorkflowCanvas.vue'
import NodeInspector from './components/NodeInspector.vue'
import EdgeInspector from './components/EdgeInspector.vue'

const { getSelectedNodes, getSelectedEdges } = useVueFlow({ id: 'workflow-editor' })

const selectedNode = computed(() => {
  const selected = getSelectedNodes.value
  return selected.length === 1 ? selected[0] : null
})

const selectedEdge = computed(() => {
  const selected = getSelectedEdges.value
  return selected.length === 1 ? selected[0] : null
})
</script>
