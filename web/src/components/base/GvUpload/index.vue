<template>
  <div class="gv-upload">
    <el-upload
      ref="uploadRef"
      v-bind="$attrs"
      :auto-upload="autoUpload"
      :on-change="handleChange"
      :on-remove="handleRemove"
      :on-success="handleSuccess"
      :on-error="handleError"
      :on-progress="handleProgress"
      :before-upload="beforeUpload"
      :file-list="fileList"
      :disabled="disabled"
      :limit="limit"
      :accept="accept"
      class="gv-upload-wrapper"
      :class="{ 'gv-upload--disabled': disabled }"
    >
      <template #trigger>
        <GvButton
          :variant="variant"
          :size="size"
          :disabled="disabled"
          :loading="uploading"
        >
          <template v-if="!uploading" #icon>
            <el-icon><Upload /></el-icon>
          </template>
          {{ uploading ? uploadingText : buttonText }}
        </GvButton>
      </template>
      <template #tip>
        <div v-if="tip" class="gv-upload-tip">
          {{ tip }}
        </div>
      </template>
    </el-upload>

    <!-- 文件列表预览 -->
    <div v-if="showFileList && fileList.length > 0" class="gv-upload-list">
      <div
        v-for="(file, index) in fileList"
        :key="index"
        class="gv-upload-item"
      >
        <div class="gv-upload-item-info">
          <el-icon class="gv-upload-item-icon">
            <Document />
          </el-icon>
          <span class="gv-upload-item-name">{{ file.name }}</span>
          <span class="gv-upload-item-size">{{ formatSize(file.size) }}</span>
        </div>
        <div v-if="file.status === 'uploading'" class="gv-upload-item-progress">
          <el-progress :percentage="file.percentage || 0" :stroke-width="4" />
        </div>
        <GvButton
          v-if="!disabled && file.status !== 'uploading'"
          variant="text"
          size="small"
          color="error"
          @click="handleRemoveFile(file)"
        >
          <template #icon>
            <el-icon><Delete /></el-icon>
          </template>
        </GvButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Upload, Document, Delete } from '@element-plus/icons-vue'
import type { UploadFile, UploadFiles } from 'element-plus'
import GvButton from '../GvButton/index.vue'

export interface GvUploadProps {
  modelValue?: UploadFile[]
  autoUpload?: boolean
  disabled?: boolean
  limit?: number
  accept?: string
  buttonText?: string
  uploadingText?: string
  tip?: string
  variant?: 'filled' | 'outlined' | 'text'
  size?: 'small' | 'medium' | 'large'
  showFileList?: boolean
  maxSize?: number
  beforeUpload?: (file: File) => boolean | Promise<boolean>
}

const props = withDefaults(defineProps<GvUploadProps>(), {
  modelValue: () => [],
  autoUpload: false,
  disabled: false,
  limit: 1,
  accept: undefined,
  buttonText: '选择文件',
  uploadingText: '上传中...',
  tip: '',
  variant: 'outlined',
  size: 'medium',
  showFileList: true,
  maxSize: undefined,
  beforeUpload: undefined
})

const emit = defineEmits<{
  'update:modelValue': [files: UploadFile[]]
  'change': [file: UploadFile, fileList: UploadFiles]
  'remove': [file: UploadFile, fileList: UploadFiles]
  'success': [response: any, file: UploadFile, fileList: UploadFiles]
  'error': [error: Error, file: UploadFile, fileList: UploadFiles]
  'progress': [event: any, file: UploadFile, fileList: UploadFiles]
}>()

const uploadRef = ref()
const uploading = ref(false)
const fileList = ref<UploadFile[]>(props.modelValue || [])

watch(
  () => props.modelValue,
  (newVal) => {
    fileList.value = newVal || []
  },
  { deep: true }
)

watch(fileList, (newVal) => {
  emit('update:modelValue', newVal)
  uploading.value = newVal.some((file) => file.status === 'uploading')
}, { deep: true })

function handleChange(file: UploadFile, files: UploadFiles) {
  fileList.value = files
  emit('change', file, files)
}

function handleRemove(file: UploadFile, files: UploadFiles) {
  fileList.value = files
  emit('remove', file, files)
}

function handleSuccess(response: any, file: UploadFile, files: UploadFiles) {
  fileList.value = files
  emit('success', response, file, files)
}

function handleError(error: Error, file: UploadFile, files: UploadFiles) {
  fileList.value = files
  emit('error', error, file, files)
}

function handleProgress(event: any, file: UploadFile, files: UploadFiles) {
  fileList.value = files
  emit('progress', event, file, files)
}

function handleRemoveFile(file: UploadFile) {
  const index = fileList.value.findIndex((f) => f.uid === file.uid)
  if (index > -1) {
    fileList.value.splice(index, 1)
    emit('remove', file, fileList.value)
  }
}

function beforeUpload(file: File): boolean | Promise<boolean> {
  if (props.maxSize && file.size > props.maxSize) {
    ElMessage.error(`文件大小不能超过 ${formatSize(props.maxSize)}`)
    return false
  }

  if (props.beforeUpload) {
    return props.beforeUpload(file)
  }

  return true
}

function formatSize(size: number): string {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`
  if (size < 1024 * 1024 * 1024) return `${(size / (1024 * 1024)).toFixed(2)} MB`
  return `${(size / (1024 * 1024 * 1024)).toFixed(2)} GB`
}

defineExpose({
  submit: () => uploadRef.value?.submit(),
  clearFiles: () => uploadRef.value?.clearFiles(),
  abort: () => uploadRef.value?.abort(),
  fileList: computed(() => fileList.value)
})
</script>

<style scoped lang="scss">
.gv-upload {
  width: 100%;
}

.gv-upload-wrapper {
  width: 100%;
}

.gv-upload--disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.gv-upload-tip {
  margin-top: 8px;
  font-size: 12px;
  color: var(--text-tertiary);
}

.gv-upload-list {
  margin-top: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.gv-upload-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  background: var(--bg-secondary);
  border-radius: 8px;
  transition: all 0.2s;
}

.gv-upload-item:hover {
  background: var(--bg-tertiary);
}

.gv-upload-item-info {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.gv-upload-item-icon {
  flex-shrink: 0;
  color: var(--text-secondary);
  font-size: 16px;
}

.gv-upload-item-name {
  flex: 1;
  font-size: 14px;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.gv-upload-item-size {
  flex-shrink: 0;
  font-size: 12px;
  color: var(--text-tertiary);
  margin-left: auto;
}

.gv-upload-item-progress {
  flex: 1;
  min-width: 0;
}
</style>
