<template>
  <div ref="vueFlowRef" class="w-full h-full relative" @dragover="onDragOver" @drop="onDrop">
    <VueFlow
      id="workflow-editor"
      :default-zoom="1.5"
      :min-zoom="0.2"
      :max-zoom="4"
      :fit-view-on-init="true"
    >
      <Background pattern-color="#aaa" :gap="8" />
      <Controls />
      <MiniMap />

      <template #node-operator="props">
        <OperatorNode v-bind="props" />
      </template>

      <template #node-trigger="props">
        <TriggerNode v-bind="props" />
      </template>

      <template #edge-conditional="props">
        <ConditionalEdge v-bind="props" />
      </template>
    </VueFlow>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { VueFlow, useVueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import { MiniMap } from '@vue-flow/minimap'
import OperatorNode from './nodes/OperatorNode.vue'
import TriggerNode from './nodes/TriggerNode.vue'
import ConditionalEdge from './edges/ConditionalEdge.vue'
import { useNodeDragDrop } from './composables/useNodeDragDrop'
import { useConnectionValidation } from './composables/useConnectionValidation'

// Import CSS
import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'
import '@vue-flow/controls/dist/style.css'
import '@vue-flow/minimap/dist/style.css'

const { nodes, edges, vueFlowRef } = useVueFlow({ id: 'workflow-editor' })

const { onDragOver, onDrop } = useNodeDragDrop()
useConnectionValidation()
</script>

<style>
/* Adjust Vue Flow Styles if needed */
.vue-flow__node {
  padding: 0;
  border-radius: 8px;
  border: 1px solid #ddd;
  background: white;
}
.vue-flow__node.selected {
  border-color: var(--el-color-primary);
  box-shadow: 0 0 0 2px var(--el-color-primary-light-8);
}
</style>
