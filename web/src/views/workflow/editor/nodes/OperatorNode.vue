<template>
  <div 
    class="min-w-[180px] bg-white rounded-lg border shadow-sm transition-all"
    :class="{ 'border-blue-500 ring-2 ring-blue-100': selected, 'border-gray-200': !selected }"
  >
    <Handle type="target" position="top" class="w-3 h-3 !bg-gray-400" />
    
    <div class="p-3">
      <div class="flex items-center gap-2 mb-2">
        <div class="w-8 h-8 rounded bg-blue-50 flex items-center justify-center text-lg">
          {{ icon }}
        </div>
        <div class="flex-1 min-w-0">
          <div class="font-medium text-sm truncate" :title="data.operatorName">
            {{ data.operatorName || '未命名算子' }}
          </div>
          <div class="text-xs text-gray-400 truncate font-mono">
            {{ data.operatorCode }}
          </div>
        </div>
      </div>
      
      <div class="text-xs text-gray-500 flex gap-2">
        <span class="bg-gray-100 px-1.5 py-0.5 rounded text-gray-600">
          v{{ version }}
        </span>
        <span class="bg-gray-100 px-1.5 py-0.5 rounded text-gray-600">
          {{ execMode }}
        </span>
      </div>
    </div>

    <Handle type="source" position="bottom" class="w-3 h-3 !bg-blue-500" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Handle } from '@vue-flow/core'
import { operatorApi } from '@/api/operator'
import { useAsyncData } from '@/composables/useAsyncData'

const props = defineProps<{
  id: string
  selected: boolean
  data: {
    operatorId: string
    operatorCode: string
    operatorName: string
    config: any
  }
}>()

// Fetch additional info if needed, or rely on passed data
const { data: operatorData } = useAsyncData(
  () => operatorApi.get(props.data.operatorId),
  { immediate: true } // Assuming operatorId is present
)

const operator = computed(() => operatorData.value?.data)

const icon = computed(() => {
  const cat = operator.value?.category || 'other'
  const map: Record<string, string> = {
    analyze: '🔍',
    edit: '✂️',
    generate: '✨',
    transform: '🔄',
    tool: '🛠️',
    other: '📦'
  }
  return map[cat] || '📦'
})

const version = computed(() => operator.value?.active_version?.version || 'latest')
const execMode = computed(() => operator.value?.active_version?.exec_mode || 'HTTP')

</script>
