<template>
  <div>
    <GvAlert
      :title="`当前执行模式：${execModeLabel}`"
      type="info"
      :closable="false"
      class="mb-3"
    />
    <div class="rounded-xl border border-border-secondary p-4 space-y-4">
      <template v-if="normalizedMode === 'http'">
        <GvInput v-model="httpConfig.endpoint" label="Endpoint" placeholder="http://localhost:8080/run" @input="emitValue" />
        <GvSelect
          v-model="httpConfig.method"
          label="Method"
          :options="methodOptions"
          class="w-full"
          @change="emitValue"
        />
        <GvInput
          v-model="httpConfig.timeout_sec"
          label="超时(秒)"
          type="number"
          placeholder="30"
          @input="emitValue"
        />
        <GvInput
          v-model="httpHeadersText"
          label="Headers"
          type="textarea"
          :rows="3"
          placeholder='JSON，例如 {"Authorization":"Bearer xxx"}'
          @input="emitValue"
        />
        <GvInput v-model="httpConfig.auth_type" label="AuthType" placeholder="bearer/basic/api_key" @input="emitValue" />
        <GvInput
          v-model="httpAuthConfigText"
          label="Auth配置"
          type="textarea"
          :rows="3"
          placeholder='JSON，例如 {"token":"xxx"}'
          @input="emitValue"
        />
      </template>

      <template v-else-if="normalizedMode === 'cli'">
        <GvInput v-model="cliConfig.command" label="Command" placeholder="python" @input="emitValue" />
        <GvInput
          v-model="cliArgsText"
          label="Args"
          type="textarea"
          :rows="3"
          placeholder='每行一个参数，例如 main.py'
          @input="emitCliArgs"
        />
        <GvInput
          v-model="cliConfig.timeout_sec"
          label="超时(秒)"
          type="number"
          placeholder="60"
          @input="emitValue"
        />
        <GvInput v-model="cliConfig.work_dir" label="工作目录" placeholder="/usr/local/bin" @input="emitValue" />
        <GvInput
          v-model="cliEnvText"
          label="环境变量"
          type="textarea"
          :rows="3"
          placeholder='JSON，例如 {"PYTHONPATH":"/app"}'
          @input="emitValue"
        />
      </template>

      <template v-else-if="normalizedMode === 'ai_model'">
        <GvSelect
          v-model="aiModelConfig.model_id"
          label="AI 模型"
          :options="aiModelOptions"
          placeholder="选择 AI 模型"
          class="w-full"
          @change="emitValue"
        />
        <GvSelect
          v-model="aiModelConfig.interaction_mode"
          label="交互模式"
          :options="interactionModeOptions"
          class="w-full"
          @change="emitValue"
        />
        <GvInput
          v-model="aiModelConfig.system_prompt"
          label="System Prompt"
          type="textarea"
          :rows="3"
          placeholder="系统提示词，如：你是一个图像分析专家..."
          @input="emitValue"
        />
        <GvInput
          v-model="aiModelConfig.user_prompt_template"
          label="User Prompt 模板"
          type="textarea"
          :rows="3"
          placeholder="支持变量：{{asset_id}}, {{params.xxx}}"
          @input="emitValue"
        />
        <div class="grid grid-cols-3 gap-4">
          <GvInput
            v-model="aiModelConfig.temperature"
            label="Temperature"
            type="number"
            placeholder="0.7"
            @input="emitValue"
          />
          <GvInput
            v-model="aiModelConfig.max_tokens"
            label="Max Tokens"
            type="number"
            placeholder="4096"
            @input="emitValue"
          />
          <GvInput
            v-model="aiModelConfig.top_p"
            label="Top P"
            type="number"
            placeholder="1.0"
            @input="emitValue"
          />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <GvSelect
            v-model="aiModelConfig.response_format"
            label="响应格式"
            :options="responseFormatOptions"
            class="w-full"
            @change="emitValue"
          />
          <GvInput
            v-model="aiModelConfig.timeout_sec"
            label="超时(秒)"
            type="number"
            placeholder="60"
            @input="emitValue"
          />
        </div>
        <GvInput
          v-model="aiModelOutputMappingText"
          label="输出映射 (JSON)"
          type="textarea"
          :rows="3"
          placeholder='JSON，定义 AI 响应到算子输出的映射'
          @input="emitValue"
        />
      </template>

      <template v-else>
        <GvInput v-model="mcpConfig.server_id" label="Server ID" placeholder="default" @input="emitValue" />
        <GvInput v-model="mcpConfig.tool_name" label="Tool Name" placeholder="echo" @input="emitValue" />
        <GvInput
          v-model="mcpConfig.timeout_sec"
          label="超时(秒)"
          type="number"
          placeholder="30"
          @input="emitValue"
        />
        <GvInput v-model="mcpConfig.tool_version" label="Tool版本" placeholder="1.0.0" @input="emitValue" />
        <GvInput
          v-model="mcpInputMappingText"
          label="输入映射"
          type="textarea"
          :rows="3"
          placeholder='JSON，例如 {"text":"$.params.prompt"}'
          @input="emitValue"
        />
        <GvInput
          v-model="mcpOutputMappingText"
          label="输出映射"
          type="textarea"
          :rows="3"
          placeholder='JSON，例如 {"result":"$.data.output"}'
          @input="emitValue"
        />
      </template>
    </div>

    <div class="mt-2 flex gap-2 items-center">
      <GvButton size="small" variant="text" @click="applyTemplateConfig">重置 {{ execModeLabel }} 默认模板</GvButton>
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
import { computed, reactive, ref, watch, onMounted } from 'vue'
import { aiModelApi } from '@/api/ai-model'
import SchemaEditor from './SchemaEditor.vue'
import GvInput from '@/components/base/GvInput/index.vue'
import GvSelect from '@/components/base/GvSelect/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvAlert from '@/components/base/GvAlert/index.vue'

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

const aiModelConfig = reactive({
  model_id: '',
  interaction_mode: 'chat',
  system_prompt: '',
  user_prompt_template: '',
  temperature: '' as string | number,
  max_tokens: '' as string | number,
  top_p: '' as string | number,
  response_format: 'text',
  timeout_sec: 60,
  output_mapping: {} as Record<string, any>
})

const aiModelOptions = ref<Array<{ label: string; value: string }>>([])

const loadAIModels = async () => {
  try {
    const res = await aiModelApi.list({ page: 1, page_size: 100 })
    aiModelOptions.value = (res.data?.items || []).map((m: any) => ({
      label: `${m.name} (${m.provider}/${m.model_name})`,
      value: m.id
    }))
  } catch (e) {
    // silently fail
  }
}

onMounted(() => {
  if (props.execMode === 'ai_model') {
    loadAIModels()
  }
})

watch(() => props.execMode, (mode) => {
  if (mode === 'ai_model') {
    loadAIModels()
  }
})

const methodOptions = [
  { label: 'POST', value: 'POST' },
  { label: 'GET', value: 'GET' },
  { label: 'PUT', value: 'PUT' },
  { label: 'PATCH', value: 'PATCH' },
  { label: 'DELETE', value: 'DELETE' }
]

const interactionModeOptions = [
  { label: 'Chat (文本对话)', value: 'chat' },
  { label: 'Vision (图像理解)', value: 'vision' }
]

const responseFormatOptions = [
  { label: '文本', value: 'text' },
  { label: 'JSON', value: 'json' }
]

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

const aiModelOutputMappingText = computed({
  get: () => formatObject(aiModelConfig.output_mapping),
  set: (value: string) => {
    aiModelConfig.output_mapping = parseObject(value)
  }
})

const normalizedMode = computed(() => {
  if (props.execMode === 'cli' || props.execMode === 'mcp' || props.execMode === 'ai_model') {
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
  const map: Record<string, string> = { http: 'HTTP', cli: 'CLI', mcp: 'MCP', ai_model: 'AI 模型' }
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
  if (normalizedMode.value === 'ai_model') {
    Object.assign(aiModelConfig, {
      model_id: '',
      interaction_mode: 'chat',
      system_prompt: '',
      user_prompt_template: '请分析资产 {{asset_id}}',
      temperature: 0.7,
      max_tokens: 4096,
      top_p: 1.0,
      response_format: 'text',
      timeout_sec: 60,
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
        timeout_sec: Number(cliConfig.timeout_sec),
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
        timeout_sec: Number(mcpConfig.timeout_sec),
        tool_version: mcpConfig.tool_version,
        input_mapping: mcpConfig.input_mapping,
        output_mapping: mcpConfig.output_mapping
      }
    })
    return
  }

  if (normalizedMode.value === 'ai_model') {
    const cfg: Record<string, any> = {
      model_id: aiModelConfig.model_id,
      interaction_mode: aiModelConfig.interaction_mode,
      system_prompt: aiModelConfig.system_prompt,
      user_prompt_template: aiModelConfig.user_prompt_template,
      response_format: aiModelConfig.response_format,
      timeout_sec: Number(aiModelConfig.timeout_sec),
      output_mapping: aiModelConfig.output_mapping
    }
    if (aiModelConfig.temperature !== '' && aiModelConfig.temperature !== undefined) {
      cfg.temperature = Number(aiModelConfig.temperature)
    }
    if (aiModelConfig.max_tokens !== '' && aiModelConfig.max_tokens !== undefined) {
      cfg.max_tokens = Number(aiModelConfig.max_tokens)
    }
    if (aiModelConfig.top_p !== '' && aiModelConfig.top_p !== undefined) {
      cfg.top_p = Number(aiModelConfig.top_p)
    }
    emit('update:modelValue', { ai_model: cfg })
    return
  }

  emit('update:modelValue', {
    http: {
      endpoint: httpConfig.endpoint,
      method: httpConfig.method,
      timeout_sec: Number(httpConfig.timeout_sec),
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
  const aiModel = v.ai_model
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
  if (aiModel && typeof aiModel === 'object') {
    aiModelConfig.model_id = String(aiModel.model_id || '')
    aiModelConfig.interaction_mode = String(aiModel.interaction_mode || 'chat')
    aiModelConfig.system_prompt = String(aiModel.system_prompt || '')
    aiModelConfig.user_prompt_template = String(aiModel.user_prompt_template || '')
    aiModelConfig.temperature = aiModel.temperature ?? ''
    aiModelConfig.max_tokens = aiModel.max_tokens ?? ''
    aiModelConfig.top_p = aiModel.top_p ?? ''
    aiModelConfig.response_format = String(aiModel.response_format || 'text')
    aiModelConfig.timeout_sec = Number(aiModel.timeout_sec || 60)
    aiModelConfig.output_mapping = isObject(aiModel.output_mapping) ? aiModel.output_mapping : {}
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
