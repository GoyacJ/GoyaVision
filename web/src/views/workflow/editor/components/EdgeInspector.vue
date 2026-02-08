<template>
  <div class="flex flex-col h-full bg-white">
    <div class="p-4 border-b">
      <h3 class="font-medium text-gray-900">连线配置</h3>
      <div class="text-xs text-gray-500 mt-1">ID: {{ edge.id }}</div>
    </div>

    <div class="flex-1 p-4 space-y-6">
      <div class="space-y-3">
        <h4 class="text-sm font-medium text-gray-700">条件设置</h4>
        <div class="space-y-4">
          <div class="space-y-1">
            <label class="block text-sm text-gray-600">执行条件</label>
            <el-select v-model="conditionType" class="w-full" size="small">
              <el-option label="始终执行" value="always" />
              <el-option label="成功时执行" value="on_success" />
              <el-option label="失败时执行" value="on_failure" />
            </el-select>
            <div class="text-xs text-gray-400 mt-1">
              控制下游节点在何时执行
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="p-4 border-t bg-gray-50">
      <el-button type="danger" plain class="w-full" @click="removeEdge">
        删除连线
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useVueFlow } from '@vue-flow/core'

const props = defineProps<{
  edge: any
}>()

const { removeEdges } = useVueFlow({ id: 'workflow-editor' })

const conditionType = computed({
  get: () => {
    return props.edge.data?.condition?.type || 'always'
  },
  set: (val) => {
    if (!props.edge.data) props.edge.data = {}
    if (!props.edge.data.condition) props.edge.data.condition = {}
    props.edge.data.condition.type = val
    // Update label or style if needed via edge.label or class
    props.edge.label = getLabel(val)
  }
})

function getLabel(type: string) {
  switch (type) {
    case 'on_success': return '成功时'
    case 'on_failure': return '失败时'
    default: return ''
  }
}

function removeEdge() {
  removeEdges([props.edge.id])
}
</script>
