<template>
  <el-form :model="form" label-width="100px">
    <el-form-item label="版本号">
      <el-input v-model="form.version" placeholder="如 v1.0.1" />
    </el-form-item>

    <el-form-item label="执行模式">
      <el-select v-model="form.exec_mode" style="width: 100%">
        <el-option label="HTTP" value="http" />
        <el-option label="CLI" value="cli" />
        <el-option label="MCP" value="mcp" />
      </el-select>
    </el-form-item>

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

    <div class="mt-4">
      <el-button type="primary" :loading="loading" :disabled="!canSubmit" @click="handleSubmit">创建版本</el-button>
    </div>
  </el-form>
</template>

<script setup lang="ts">
import { computed, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import ExecConfigForm from './ExecConfigForm.vue'
import SchemaEditor from './SchemaEditor.vue'

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
