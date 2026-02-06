<template>
  <el-form :model="form" label-width="100px">
    <el-form-item label="算子代码">
      <el-input v-model="form.code" placeholder="唯一标识，如 frame_extract" />
    </el-form-item>
    <el-form-item label="算子名称">
      <el-input v-model="form.name" placeholder="请输入算子名称" />
    </el-form-item>
    <el-form-item label="描述">
      <el-input v-model="form.description" type="textarea" :rows="2" />
    </el-form-item>
    <el-form-item label="分类">
      <el-select v-model="form.category" style="width: 100%">
        <el-option label="分析" value="analysis" />
        <el-option label="处理" value="processing" />
        <el-option label="生成" value="generation" />
        <el-option label="工具" value="utility" />
      </el-select>
    </el-form-item>
    <el-form-item label="类型">
      <el-input v-model="form.type" placeholder="如 object_detection / transcode" />
    </el-form-item>
    <el-form-item label="来源">
      <el-select v-model="form.origin" style="width: 100%">
        <el-option label="自定义" value="custom" />
        <el-option label="内置" value="builtin" />
        <el-option label="市场" value="marketplace" />
        <el-option label="MCP" value="mcp" />
      </el-select>
    </el-form-item>
    <el-form-item label="执行模式">
      <el-select v-model="form.exec_mode" style="width: 100%">
        <el-option label="HTTP" value="http" />
        <el-option label="CLI" value="cli" />
        <el-option label="MCP" value="mcp" />
      </el-select>
    </el-form-item>

    <ExecConfigForm v-model="form.exec_config" :exec-mode="form.exec_mode" />

    <div class="mt-4 flex gap-2">
      <el-button type="primary" :loading="loading" @click="emit('submit', form)">保存</el-button>
      <el-button @click="emit('cancel')">取消</el-button>
    </div>
  </el-form>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'
import ExecConfigForm from './ExecConfigForm.vue'

type OperatorFormModel = {
  code: string
  name: string
  description?: string
  category: 'analysis' | 'processing' | 'generation' | 'utility'
  type: string
  origin: 'custom' | 'builtin' | 'marketplace' | 'mcp'
  exec_mode: 'http' | 'cli' | 'mcp'
  exec_config?: Record<string, any>
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
  exec_config: {}
})

watch(
  () => props.modelValue,
  (value) => {
    if (!value) return
    Object.assign(form, value)
  },
  { immediate: true, deep: true }
)
</script>
