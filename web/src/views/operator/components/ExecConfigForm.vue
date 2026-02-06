<template>
  <div>
    <el-alert
      :title="`当前执行模式：${execMode || '未设置'}`"
      type="info"
      :closable="false"
      class="mb-3"
    />
    <div class="rounded border border-border-secondary p-3">
      <template v-if="normalizedMode === 'http'">
        <el-form-item label="Endpoint" label-width="100px" class="mb-3">
          <el-input v-model="httpConfig.endpoint" placeholder="http://localhost:8080/run" @input="emitValue" />
        </el-form-item>
        <el-form-item label="Method" label-width="100px" class="mb-3">
          <el-select v-model="httpConfig.method" style="width: 100%" @change="emitValue">
            <el-option label="POST" value="POST" />
            <el-option label="GET" value="GET" />
            <el-option label="PUT" value="PUT" />
            <el-option label="PATCH" value="PATCH" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item label="超时(秒)" label-width="100px" class="mb-0">
          <el-input-number v-model="httpConfig.timeout_sec" :min="1" :max="3600" @change="emitValue" />
        </el-form-item>
        <el-form-item label="Headers" label-width="100px" class="mt-3 mb-0">
          <el-input
            v-model="httpHeadersText"
            type="textarea"
            :rows="3"
            placeholder='JSON，例如 {"Authorization":"Bearer xxx"}'
            @input="emitValue"
          />
        </el-form-item>
        <el-form-item label="AuthType" label-width="100px" class="mt-3 mb-0">
          <el-input v-model="httpConfig.auth_type" placeholder="bearer/basic/apikey" @input="emitValue" />
        </el-form-item>
        <el-form-item label="Auth配置" label-width="100px" class="mt-3 mb-0">
          <el-input
            v-model="httpAuthConfigText"
            type="textarea"
            :rows="3"
            placeholder='JSON，例如 {"token":"xxx"}'
            @input="emitValue"
          />
        </el-form-item>
      </template>

      <template v-else-if="normalizedMode === 'cli'">
        <el-form-item label="Command" label-width="100px" class="mb-3">
          <el-input v-model="cliConfig.command" placeholder="python" @input="emitValue" />
        </el-form-item>
        <el-form-item label="Args" label-width="100px" class="mb-3">
          <el-input v-model="cliArgsText" placeholder='每行一个参数，例如 main.py' type="textarea" :rows="3" @input="emitCliArgs" />
        </el-form-item>
        <el-form-item label="超时(秒)" label-width="100px" class="mb-0">
          <el-input-number v-model="cliConfig.timeout_sec" :min="1" :max="3600" @change="emitValue" />
        </el-form-item>
        <el-form-item label="工作目录" label-width="100px" class="mt-3 mb-0">
          <el-input v-model="cliConfig.work_dir" placeholder="/usr/local/bin" @input="emitValue" />
        </el-form-item>
        <el-form-item label="环境变量" label-width="100px" class="mt-3 mb-0">
          <el-input
            v-model="cliEnvText"
            type="textarea"
            :rows="3"
            placeholder='JSON，例如 {"PYTHONPATH":"/app"}'
            @input="emitValue"
          />
        </el-form-item>
      </template>

      <template v-else>
        <el-form-item label="Server ID" label-width="100px" class="mb-3">
          <el-input v-model="mcpConfig.server_id" placeholder="default" @input="emitValue" />
        </el-form-item>
        <el-form-item label="Tool Name" label-width="100px" class="mb-3">
          <el-input v-model="mcpConfig.tool_name" placeholder="echo" @input="emitValue" />
        </el-form-item>
        <el-form-item label="超时(秒)" label-width="100px" class="mb-0">
          <el-input-number v-model="mcpConfig.timeout_sec" :min="1" :max="3600" @change="emitValue" />
        </el-form-item>
        <el-form-item label="Tool版本" label-width="100px" class="mt-3 mb-0">
          <el-input v-model="mcpConfig.tool_version" placeholder="1.0.0" @input="emitValue" />
        </el-form-item>
        <el-form-item label="输入映射" label-width="100px" class="mt-3 mb-0">
          <el-input
            v-model="mcpInputMappingText"
            type="textarea"
            :rows="3"
            placeholder='JSON，例如 {"text":"$.params.prompt"}'
            @input="emitValue"
          />
        </el-form-item>
        <el-form-item label="输出映射" label-width="100px" class="mt-3 mb-0">
          <el-input
            v-model="mcpOutputMappingText"
            type="textarea"
            :rows="3"
            placeholder='JSON，例如 {"result":"$.data.output"}'
            @input="emitValue"
          />
        </el-form-item>
      </template>
    </div>

    <div class="mt-2 flex gap-2">
      <el-button size="small" text @click="applyTemplateConfig">重置 {{ execModeLabel }} 默认模板</el-button>
      <span class="text-xs text-text-tertiary">已提供结构化表单，便于按模式编辑配置</span>
    </div>

    <SchemaEditor
      class="mt-3"
      :model-value="modelValue"
      title="执行配置 JSON 预览"
      :rows="5"
      @update:model-value="(v) => applyJsonEditorValue(v)"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, watch } from 'vue'
import SchemaEditor from './SchemaEditor.vue'

const props = defineProps<{
  modelValue?: Record<string, any>
  execMode?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: Record<string, any>]
}>()

const httpConfig = reactive({
  endpoint: '',
  method: 'POST',
  timeout_sec: 30,
  headers: {} as Record<string, string>,
  auth_type: '',
  auth_config: {} as Record<string, string>
})

const cliConfig = reactive({
  command: 'python',
  args: [] as string[],
  timeout_sec: 60,
  work_dir: '',
  env: {} as Record<string, string>
})

const mcpConfig = reactive({
  server_id: 'default',
  tool_name: 'echo',
  timeout_sec: 30,
  tool_version: '',
  input_mapping: {} as Record<string, any>,
  output_mapping: {} as Record<string, any>
})

const httpHeadersText = computed({
  get: () => formatObject(httpConfig.headers),
  set: (value: string) => {
    httpConfig.headers = parseObject(value) as Record<string, string>
  }
})

const cliEnvText = computed({
  get: () => formatObject(cliConfig.env),
  set: (value: string) => {
    cliConfig.env = parseObject(value) as Record<string, string>
  }
})

const httpAuthConfigText = computed({
  get: () => formatObject(httpConfig.auth_config),
  set: (value: string) => {
    httpConfig.auth_config = parseObject(value) as Record<string, string>
  }
})

const mcpInputMappingText = computed({
  get: () => formatObject(mcpConfig.input_mapping),
  set: (value: string) => {
    mcpConfig.input_mapping = parseObject(value)
  }
})

const mcpOutputMappingText = computed({
  get: () => formatObject(mcpConfig.output_mapping),
  set: (value: string) => {
    mcpConfig.output_mapping = parseObject(value)
  }
})

const normalizedMode = computed(() => {
  if (props.execMode === 'cli' || props.execMode === 'mcp') {
    return props.execMode
  }
  return 'http'
})

const cliArgsText = computed({
  get: () => cliConfig.args.join('\n'),
  set: (value: string) => {
    cliConfig.args = value
      .split('\n')
      .map((s) => s.trim())
      .filter(Boolean)
  }
})

const execModeLabel = computed(() => {
  const mode = props.execMode || 'http'
  const map: Record<string, string> = { http: 'HTTP', cli: 'CLI', mcp: 'MCP' }
  return map[mode] || mode.toUpperCase()
})

function applyTemplateConfig() {
  if (normalizedMode.value === 'cli') {
    Object.assign(cliConfig, {
      command: 'python',
      args: ['main.py'],
      timeout_sec: 60,
      work_dir: '',
      env: {}
    })
    emitValue()
    return
  }
  if (normalizedMode.value === 'mcp') {
    Object.assign(mcpConfig, {
      server_id: 'default',
      tool_name: 'echo',
      timeout_sec: 30,
      tool_version: '',
      input_mapping: {},
      output_mapping: {}
    })
    emitValue()
    return
  }
  Object.assign(httpConfig, {
    endpoint: 'http://localhost:8080/run',
    method: 'POST',
    timeout_sec: 30,
    headers: {},
    auth_type: '',
    auth_config: {}
  })
  emitValue()
}

function emitCliArgs() {
  cliConfig.args = cliArgsText.value
    .split('\n')
    .map((s) => s.trim())
    .filter(Boolean)
  emitValue()
}

function emitValue() {
  if (normalizedMode.value === 'cli') {
    emit('update:modelValue', {
      cli: {
        command: cliConfig.command,
        args: cliConfig.args,
        timeout_sec: cliConfig.timeout_sec,
        work_dir: cliConfig.work_dir,
        env: cliConfig.env
      }
    })
    return
  }

  if (normalizedMode.value === 'mcp') {
    emit('update:modelValue', {
      mcp: {
        server_id: mcpConfig.server_id,
        tool_name: mcpConfig.tool_name,
        timeout_sec: mcpConfig.timeout_sec,
        tool_version: mcpConfig.tool_version,
        input_mapping: mcpConfig.input_mapping,
        output_mapping: mcpConfig.output_mapping
      }
    })
    return
  }

  emit('update:modelValue', {
    http: {
      endpoint: httpConfig.endpoint,
      method: httpConfig.method,
      timeout_sec: httpConfig.timeout_sec,
      headers: httpConfig.headers,
      auth_type: httpConfig.auth_type,
      auth_config: httpConfig.auth_config
    }
  })
}

function syncFromModel(v?: Record<string, any>) {
  if (!v || typeof v !== 'object') {
    return
  }

  const http = v.http
  const cli = v.cli
  const mcp = v.mcp
  if (http && typeof http === 'object') {
    httpConfig.endpoint = String(http.endpoint || '')
    httpConfig.method = String(http.method || 'POST').toUpperCase()
    httpConfig.timeout_sec = Number(http.timeout_sec || 30)
    httpConfig.headers = isObject(http.headers) ? (http.headers as Record<string, string>) : {}
    httpConfig.auth_type = String(http.auth_type || '')
    httpConfig.auth_config = isObject(http.auth_config) ? (http.auth_config as Record<string, string>) : {}
  }
  if (cli && typeof cli === 'object') {
    cliConfig.command = String(cli.command || 'python')
    cliConfig.args = Array.isArray(cli.args) ? cli.args.map((i) => String(i)) : []
    cliConfig.timeout_sec = Number(cli.timeout_sec || 60)
    cliConfig.work_dir = String(cli.work_dir || '')
    cliConfig.env = isObject(cli.env) ? (cli.env as Record<string, string>) : {}
  }
  if (mcp && typeof mcp === 'object') {
    mcpConfig.server_id = String(mcp.server_id || 'default')
    mcpConfig.tool_name = String(mcp.tool_name || 'echo')
    mcpConfig.timeout_sec = Number(mcp.timeout_sec || 30)
    mcpConfig.tool_version = String(mcp.tool_version || '')
    mcpConfig.input_mapping = isObject(mcp.input_mapping) ? mcp.input_mapping : {}
    mcpConfig.output_mapping = isObject(mcp.output_mapping) ? mcp.output_mapping : {}
  }
}

function applyJsonEditorValue(v?: Record<string, any>) {
  syncFromModel(v)
  emit('update:modelValue', v)
}

watch(
  () => props.modelValue,
  (v) => syncFromModel(v),
  { immediate: true, deep: true }
)

function parseObject(raw: string): Record<string, any> {
  const v = String(raw || '').trim()
  if (!v) return {}
  try {
    const parsed = JSON.parse(v)
    return isObject(parsed) ? parsed : {}
  } catch {
    return {}
  }
}

function formatObject(obj: unknown): string {
  if (!isObject(obj) || Object.keys(obj).length === 0) {
    return ''
  }
  return JSON.stringify(obj, null, 2)
}

function isObject(v: unknown): v is Record<string, any> {
  return !!v && typeof v === 'object' && !Array.isArray(v)
}
</script>
