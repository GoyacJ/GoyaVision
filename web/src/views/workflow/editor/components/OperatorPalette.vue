<template>
  <div class="flex flex-col h-full">
    <div class="p-3 border-b bg-gray-50 font-bold text-gray-700">
      算子库
    </div>
    <div class="flex-1 overflow-y-auto p-2">
      <div v-for="(group, category) in groupedOperators" :key="category" class="mb-4">
        <div class="text-xs font-bold text-gray-500 uppercase mb-2 px-1">
          {{ getCategoryLabel(category) }}
        </div>
        <div class="space-y-2">
          <div
            v-for="op in group"
            :key="op.id"
            class="p-2 bg-white border rounded shadow-sm hover:border-blue-400 hover:shadow cursor-grab active:cursor-grabbing transition-all flex items-center gap-2"
            draggable="true"
            @dragstart="(event) => onDragStart(event, op)"
          >
            <div class="w-8 h-8 rounded bg-gray-100 flex items-center justify-center text-lg">
              {{ getIcon(op.category) }}
            </div>
            <div class="flex-1 min-w-0">
              <div class="text-sm font-medium truncate">{{ op.name }}</div>
              <div class="text-xs text-gray-400 truncate">{{ op.code }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAsyncData } from '@/composables/useAsyncData'
import { operatorApi } from '@/api/operator'
import type { Operator } from '@/api/operator'

const { data: operatorList } = useAsyncData(async () => {
  const res = await operatorApi.list({ page_size: 100 })
  return res.data
}, { immediate: true })

const operators = computed(() => operatorList.value?.items || [])

const groupedOperators = computed(() => {
  const groups: Record<string, Operator[]> = {}
  if (!operators.value) return groups

  operators.value.forEach((op) => {
    const cat = op.category || 'other'
    if (!groups[cat]) groups[cat] = []
    groups[cat].push(op)
  })
  return groups
})

function getCategoryLabel(category: string | number) {
  const map: Record<string, string> = {
    analyze: '分析',
    edit: '编辑',
    generate: '生成',
    transform: '转换',
    tool: '工具',
    other: '其他'
  }
  return map[String(category)] || String(category)
}

function getIcon(category: string) {
  const map: Record<string, string> = {
    analyze: '🔍',
    edit: '✂️',
    generate: '✨',
    transform: '🔄',
    tool: '🛠️',
    other: '📦'
  }
  return map[category] || '📦'
}

function onDragStart(event: DragEvent, op: Operator) {
  if (event.dataTransfer) {
    const data = {
      type: 'operator',
      operatorId: op.id,
      operatorCode: op.code,
      operatorName: op.name
    }
    event.dataTransfer.setData('text', JSON.stringify(data))
    event.dataTransfer.effectAllowed = 'copy'
  }
}
</script>
