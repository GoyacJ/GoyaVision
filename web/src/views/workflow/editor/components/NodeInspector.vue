<template>
  <div class="flex flex-col h-full bg-white">
    <div class="p-4 border-b">
      <h3 class="font-medium text-gray-900">节点配置</h3>
      <div class="text-xs text-gray-500 mt-1">ID: {{ node.id }}</div>
    </div>

    <div class="flex-1 overflow-y-auto p-4 space-y-6">
      <!-- Operator Info -->
      <div class="space-y-3">
        <h4 class="text-sm font-medium text-gray-700">算子信息</h4>
        <div class="bg-gray-50 rounded p-3 text-sm space-y-2">
          <div class="flex justify-between">
            <span class="text-gray-500">名称</span>
            <span class="font-medium">{{ node.data.operatorName }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-500">编码</span>
            <span class="font-mono text-xs">{{ node.data.operatorCode }}</span>
          </div>
        </div>
      </div>

      <!-- Parameters -->
      <div class="space-y-3" v-if="schema && Object.keys(schema).length > 0">
        <h4 class="text-sm font-medium text-gray-700">参数配置</h4>
        <div class="space-y-4">
          <template v-for="(prop, key) in schema" :key="key">
            <div class="space-y-1">
              <label class="block text-sm text-gray-600">
                {{ prop.title || key }}
                <span v-if="prop.required" class="text-red-500">*</span>
              </label>
              
              <!-- String/Select -->
              <el-input
                v-if="prop.type === 'string' && !prop.enum"
                v-model="node.data.config.params[key]"
                :placeholder="prop.description"
                size="small"
              />
              <el-select
                v-else-if="prop.type === 'string' && prop.enum"
                v-model="node.data.config.params[key]"
                :placeholder="prop.description"
                size="small"
                class="w-full"
              >
                <el-option
                  v-for="opt in prop.enum"
                  :key="opt"
                  :label="opt"
                  :value="opt"
                />
              </el-select>

              <!-- Number -->
              <el-input-number
                v-else-if="prop.type === 'number' || prop.type === 'integer'"
                v-model="node.data.config.params[key]"
                :min="prop.minimum"
                :max="prop.maximum"
                size="small"
                class="w-full"
              />

              <!-- Boolean -->
              <el-switch
                v-else-if="prop.type === 'boolean'"
                v-model="node.data.config.params[key]"
              />

              <div v-if="prop.description" class="text-xs text-gray-400 mt-1">
                {{ prop.description }}
              </div>
            </div>
          </template>
        </div>
      </div>

      <!-- Execution Config -->
      <div class="space-y-3">
        <h4 class="text-sm font-medium text-gray-700">运行设置</h4>
        <div class="space-y-4">
          <div class="space-y-1">
            <label class="block text-sm text-gray-600">重试次数</label>
            <el-input-number
              v-model="node.data.config.retry_count"
              :min="0"
              :max="10"
              size="small"
              class="w-full"
            />
          </div>
          <div class="space-y-1">
            <label class="block text-sm text-gray-600">超时时间 (秒)</label>
            <el-input-number
              v-model="node.data.config.timeout_seconds"
              :min="0"
              :step="10"
              size="small"
              class="w-full"
            />
          </div>
        </div>
      </div>
    </div>

    <div class="p-4 border-t bg-gray-50">
      <el-button type="danger" plain class="w-full" @click="removeNode">
        删除节点
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import { useVueFlow } from '@vue-flow/core'
import { operatorApi } from '@/api/operator'
import { useAsyncData } from '@/composables/useAsyncData'

const props = defineProps<{
  node: any
}>()

const { removeNodes } = useVueFlow({ id: 'workflow-editor' })

// Ensure config structure exists
watch(() => props.node, (newNode) => {
  if (!newNode.data.config) {
    newNode.data.config = { params: {}, retry_count: 0, timeout_seconds: 0 }
  }
  if (!newNode.data.config.params) {
    newNode.data.config.params = {}
  }
}, { immediate: true })

// Fetch full operator details to get input schema
const { data: operatorData, execute } = useAsyncData(
  async () => {
    const res = await operatorApi.get(props.node.data.operatorId)
    return res.data
  },
  { immediate: false }
)

const operator = computed(() => operatorData.value)

watch(() => props.node.data.operatorId, (newId) => {
  if (newId) execute()
}, { immediate: true })

const schema = computed(() => {
  const inputSchema = operator.value?.active_version?.input_schema
  if (!inputSchema || !inputSchema.properties) return {}
  return inputSchema.properties
})

function removeNode() {
  removeNodes([props.node.id])
}
</script>
