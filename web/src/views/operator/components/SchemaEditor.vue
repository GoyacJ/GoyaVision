<template>
  <div class="schema-editor">
    <div class="mb-2 flex items-center gap-2 text-sm font-medium text-text-primary">
      <span>{{ title }}</span>
      <span v-if="validating" class="text-xs text-text-tertiary">校验中...</span>
    </div>
    <el-input
      v-model="text"
      type="textarea"
      :rows="rows"
      :placeholder="placeholder"
      @blur="handleBlur"
    />
    <div v-if="error" class="mt-1 text-xs text-danger-500">{{ error }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useJsonSchema } from '@/composables'

const props = withDefaults(defineProps<{
  modelValue?: Record<string, any>
  title?: string
  rows?: number
  placeholder?: string
}>(), {
  title: 'JSON Schema',
  rows: 8,
  placeholder: '请输入合法 JSON'
})

const emit = defineEmits<{
  'update:modelValue': [value: Record<string, any>]
  validate: [payload: { valid: boolean; message?: string; value?: Record<string, any> }]
}>()

const text = ref('')
const error = ref('')
const validating = ref(false)
const { parseJsonObject, validateSchema } = useJsonSchema()

watch(
  () => props.modelValue,
  (value) => {
    text.value = value ? JSON.stringify(value, null, 2) : '{}'
  },
  { immediate: true }
)

async function handleBlur() {
  const parsedResult = parseJsonObject(text.value || '{}')
  if (!parsedResult.valid || !parsedResult.data) {
    error.value = parsedResult.message || 'JSON 格式不合法'
    emit('validate', { valid: false, message: error.value })
    return
  }

  emit('update:modelValue', parsedResult.data)
  validating.value = true
  const result = await validateSchema(parsedResult.data)
  validating.value = false

  if (!result.valid) {
    error.value = result.message || 'Schema 校验失败'
    emit('validate', { valid: false, message: error.value, value: parsedResult.data })
    return
  }

  error.value = ''
  emit('validate', { valid: true, value: parsedResult.data })
}
</script>
