<template>
  <path :id="id" class="vue-flow__edge-path" :d="path[0]" :marker-end="markerEnd" />
  
  <div v-if="label" class="vue-flow__edge-label-container" :style="{ transform: `translate(-50%, -50%) translate(${path[1]}px,${path[2]}px)` }">
    <div 
      class="px-2 py-1 rounded text-xs font-bold border"
      :class="labelClass"
    >
      {{ label }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { getBezierPath, useVueFlow } from '@vue-flow/core'

const props = defineProps<{
  id: string
  sourceX: number
  sourceY: number
  targetX: number
  targetY: number
  sourcePosition: any
  targetPosition: any
  markerEnd?: string
  data?: any
}>()

const path = computed(() => getBezierPath(props))

const label = computed(() => {
  const type = props.data?.condition?.type
  switch (type) {
    case 'on_success': return '成功'
    case 'on_failure': return '失败'
    default: return ''
  }
})

const labelClass = computed(() => {
  const type = props.data?.condition?.type
  switch (type) {
    case 'on_success': return 'bg-green-100 text-green-700 border-green-200'
    case 'on_failure': return 'bg-red-100 text-red-700 border-red-200'
    default: return 'bg-white text-gray-500 border-gray-200'
  }
})
</script>

<style>
.vue-flow__edge-path {
  stroke: #b1b1b7;
  stroke-width: 2;
}
.vue-flow__edge:hover .vue-flow__edge-path {
  stroke: #555;
}
.vue-flow__edge.selected .vue-flow__edge-path {
  stroke: #3b82f6;
}
</style>
