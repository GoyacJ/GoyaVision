<template>
  <div class="space-y-4">
    <GvInput v-model="form.version" label="版本号" placeholder="如 v1.0.1" />

    <GvSelect
      v-model="form.exec_mode"
      label="执行模式"
      :options="execModeOptions"
      class="w-full"
    />

    <ExecConfigForm v-model="form.exec_config" :exec-mode="form.exec_mode" />

    <div class="mt-4 grid grid-cols-1 gap-3 lg:grid-cols-2">
      <SchemaEditor
        v-model="form.input_schema"
        title="输入 Schema"
        @validate="(p) => { schemaValid.input = p.valid }"
      />
      <SchemaEditor
        v-model="form.output_spec"
        title="输出 Spec"
        @validate="(p) => { schemaValid.output = p.valid }"
      />
    </div>

    <div class="mt-4">
      <SchemaEditor
        v-model="form.config"
        title="兼容配置（config）"
        :rows="5"
        @validate="(p) => { schemaValid.config = p.valid }"
      />
    </div>

    <div class="mt-4 flex justify-end">
      <GvButton variant="filled" color="primary" :loading="loading" :disabled="!canSubmit" @click="handleSubmit">
        创建版本
      </GvButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import ExecConfigForm from './ExecConfigForm.vue'
import SchemaEditor from './SchemaEditor.vue'
import GvInput from '@/components/base/GvInput/index.vue'
import GvSelect from '@/components/base/GvSelect/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'

type VersionFormModel = {
  version: string
  exec_mode: 'http' | 'cli' | 'mcp'
  exec_config?: Record<string, any>
  input_schema?: Record<string, any>
  output_spec?: Record<string, any>
  config?: Record<string, any>
  changelog?: string
}

const props = defineProps<{
  modelValue?: VersionFormModel
  loading?: boolean
}>()

const emit = defineEmits<{
  submit: [payload: VersionFormModel]
}>()

const form = reactive<VersionFormModel>({
  version: '',
  exec_mode: 'http',
  exec_config: {},
  input_schema: {},
  output_spec: {},
  config: {}
})

const schemaValid = reactive({
  input: true,
  output: true,
  config: true
})

const canSubmit = computed(() => {
  return !!form.version && !!form.exec_mode && schemaValid.input && schemaValid.output && schemaValid.config
})

const semverPattern = /^v?\d+\.\d+\.\d+$/

const execModeOptions = [
  { label: 'HTTP', value: 'http' },
  { label: 'CLI', value: 'cli' },
  { label: 'MCP', value: 'mcp' }
]

function handleSubmit() {
  if (!semverPattern.test(form.version)) {
    ElMessage.warning('版本号格式不正确，请使用 x.y.z 或 vx.y.z')
    return
  }
  emit('submit', form)
}

watch(
  () => props.modelValue,
  (value) => {
    if (!value) return
    Object.assign(form, value)
  },
  { immediate: true, deep: true }
)
</script>
