<template>
  <div class="space-y-4">
    <GvInput v-model="form.code" label="算子代码" placeholder="唯一标识，如 frame_extract" />
    <GvInput v-model="form.name" label="算子名称" placeholder="请输入算子名称" />
    <GvInput v-model="form.description" label="描述" type="textarea" :rows="2" />
    
    <div class="grid grid-cols-2 gap-4">
      <GvSelect
        v-model="form.category"
        label="分类"
        :options="categoryOptions"
        class="w-full"
      />
      <GvInput v-model="form.type" label="类型" placeholder="如 object_detection / transcode" />
    </div>

    <div class="grid grid-cols-2 gap-4">
      <GvSelect
        v-model="form.origin"
        label="来源"
        :options="originOptions"
        class="w-full"
      />
      <GvSelect
        v-model="form.exec_mode"
        label="执行模式"
        :options="execModeOptions"
        class="w-full"
      />
    </div>

    <GvSelect
      v-model="form.ai_model_id"
      label="AI 模型"
      :options="aiModelOptions"
      placeholder="选择关联的 AI 模型 (可选)"
      clearable
      class="w-full"
    />

    <ExecConfigForm v-model="form.exec_config" :exec-mode="form.exec_mode" />
  </div>
</template>

<script setup lang="ts">
import { reactive, watch, ref } from 'vue'
import { aiModelApi } from '@/api/ai-model'
import ExecConfigForm from './ExecConfigForm.vue'
import GvInput from '@/components/base/GvInput/index.vue'
import GvSelect from '@/components/base/GvSelect/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'

type OperatorFormModel = {
  code: string
  name: string
  description?: string
  category: 'analysis' | 'processing' | 'generation' | 'utility'
  type: string
  origin: 'custom' | 'builtin' | 'marketplace' | 'mcp'
  exec_mode: 'http' | 'cli' | 'mcp'
  exec_config?: Record<string, any>
  ai_model_id?: string
}

const props = defineProps<{
  modelValue?: Partial<OperatorFormModel>
  loading?: boolean
}>()

const emit = defineEmits<{
  submit: [payload: OperatorFormModel]
  cancel: []
}>()

const form = reactive<OperatorFormModel>({
  code: '',
  name: '',
  description: '',
  category: 'analysis',
  type: '',
  origin: 'custom',
  exec_mode: 'http',
  exec_config: {},
  ai_model_id: ''
})

const aiModelOptions = ref<Array<{ label: string; value: string }>>([])
const loadAIModels = async () => {
  try {
    const res = await aiModelApi.list({ page: 1, page_size: 100 })
    aiModelOptions.value = (res.data?.items || []).map((m: any) => ({
      label: m.name,
      value: m.id
    }))
  } catch (e) {
    console.error(e)
  }
}
loadAIModels()

const categoryOptions = [
  { label: '分析', value: 'analysis' },
  { label: '处理', value: 'processing' },
  { label: '生成', value: 'generation' },
  { label: '工具', value: 'utility' }
]

const originOptions = [
  { label: '自定义', value: 'custom' },
  { label: '内置', value: 'builtin' },
  { label: '市场', value: 'marketplace' },
  { label: 'MCP', value: 'mcp' }
]

const execModeOptions = [
  { label: 'HTTP', value: 'http' },
  { label: 'CLI', value: 'cli' },
  { label: 'MCP', value: 'mcp' }
]

watch(
  () => props.modelValue,
  (value) => {
    if (!value) return
    Object.assign(form, value)
  },
  { immediate: true, deep: true }
)

const submit = () => {
  emit('submit', form)
}

defineExpose({
  submit
})
</script>
