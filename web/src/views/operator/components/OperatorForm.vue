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

    <div class="grid grid-cols-2 gap-4">
      <GvSelect
        v-model="form.visibility"
        label="可见范围"
        :options="visibilityOptions"
        class="w-full"
      />
      <div v-if="form.visibility === 1">
        <label class="block text-sm font-medium text-text-secondary mb-1">可见角色</label>
        <el-select
          v-model="form.visible_role_ids"
          multiple
          placeholder="请选择可见角色"
          class="w-full"
        >
          <el-option
            v-for="role in roleOptions"
            :key="role.value"
            :label="role.label"
            :value="role.value"
          />
        </el-select>
      </div>
    </div>

    <ExecConfigForm v-model="form.exec_config" :exec-mode="form.exec_mode" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { roleApi } from '@/api/role'
import ExecConfigForm from './ExecConfigForm.vue'
import GvInput from '@/components/base/GvInput/index.vue'
import GvSelect from '@/components/base/GvSelect/index.vue'

type OperatorFormModel = {
  code: string
  name: string
  description?: string
  category: 'analysis' | 'processing' | 'generation' | 'utility'
  type: string
  origin: 'custom' | 'builtin' | 'marketplace' | 'mcp'
  exec_mode: 'http' | 'cli' | 'mcp' | 'ai_model'
  exec_config?: Record<string, any>
  visibility?: number
  visible_role_ids?: string[]
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
  visibility: 0,
  visible_role_ids: []
})

const roleOptions = ref<{ label: string; value: string }[]>([])
const loadRoles = async () => {
  try {
    const res = await roleApi.list()
    roleOptions.value = (res.data || []).map((r: any) => ({ label: r.name, value: r.id }))
  } catch (e) {
    console.error('Failed to load roles', e)
  }
}
loadRoles()

const visibilityOptions = [
  { label: '私有', value: 0 },
  { label: '角色可见', value: 1 },
  { label: '公开', value: 2 }
]

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
  { label: 'MCP', value: 'mcp' },
  { label: 'AI 模型', value: 'ai_model' }
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
