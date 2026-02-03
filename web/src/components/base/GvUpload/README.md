# GvUpload 文件上传组件

通用的文件上传组件，基于 Element Plus Upload 封装，提供统一的样式和交互体验。

## 基础用法

```vue
<template>
  <GvUpload
    v-model="fileList"
    :action="uploadUrl"
    @success="handleSuccess"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { UploadFile } from 'element-plus'
import GvUpload from '@/components/base/GvUpload/index.vue'

const fileList = ref<UploadFile[]>([])
const uploadUrl = '/api/v1/upload'

function handleSuccess(response: any, file: UploadFile) {
  console.log('上传成功', response)
}
</script>
```

## 自动上传

```vue
<GvUpload
  v-model="fileList"
  :action="uploadUrl"
  :auto-upload="true"
  @success="handleSuccess"
/>
```

## 限制文件类型和大小

```vue
<GvUpload
  v-model="fileList"
  :action="uploadUrl"
  accept="image/*"
  :max-size="5 * 1024 * 1024"
  tip="支持 JPG、PNG 格式，大小不超过 5MB"
/>
```

## 多文件上传

```vue
<GvUpload
  v-model="fileList"
  :action="uploadUrl"
  :limit="5"
  tip="最多上传 5 个文件"
/>
```

## 自定义按钮样式

```vue
<GvUpload
  v-model="fileList"
  :action="uploadUrl"
  button-text="上传图片"
  variant="filled"
  size="large"
/>
```

## Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| `modelValue` | 文件列表（v-model） | `UploadFile[]` | `[]` |
| `autoUpload` | 是否自动上传 | `boolean` | `false` |
| `disabled` | 是否禁用 | `boolean` | `false` |
| `limit` | 最大上传数量 | `number` | `1` |
| `accept` | 接受的文件类型 | `string` | - |
| `buttonText` | 按钮文字 | `string` | `"选择文件"` |
| `uploadingText` | 上传中按钮文字 | `string` | `"上传中..."` |
| `tip` | 提示文字 | `string` | - |
| `variant` | 按钮样式变体 | `'filled' \| 'outlined' \| 'text'` | `"outlined"` |
| `size` | 按钮大小 | `'small' \| 'medium' \| 'large'` | `"medium"` |
| `showFileList` | 是否显示文件列表 | `boolean` | `true` |
| `maxSize` | 最大文件大小（字节） | `number` | - |
| `beforeUpload` | 上传前的钩子 | `(file: File) => boolean \| Promise<boolean>` | - |

## Events

| 事件名 | 说明 | 回调参数 |
|--------|------|----------|
| `update:modelValue` | 文件列表更新 | `(files: UploadFile[])` |
| `change` | 文件状态改变 | `(file: UploadFile, fileList: UploadFiles)` |
| `remove` | 文件移除 | `(file: UploadFile, fileList: UploadFiles)` |
| `success` | 上传成功 | `(response: any, file: UploadFile, fileList: UploadFiles)` |
| `error` | 上传失败 | `(error: Error, file: UploadFile, fileList: UploadFiles)` |
| `progress` | 上传进度 | `(event: any, file: UploadFile, fileList: UploadFiles)` |

## Methods

通过 ref 可以调用以下方法：

| 方法名 | 说明 | 参数 |
|--------|------|------|
| `submit` | 手动上传文件列表 | - |
| `clearFiles` | 清空文件列表 | - |
| `abort` | 取消上传 | - |

## 完整示例

```vue
<template>
  <GvUpload
    ref="uploadRef"
    v-model="fileList"
    :action="uploadUrl"
    :auto-upload="false"
    :limit="3"
    accept="image/*"
    :max-size="10 * 1024 * 1024"
    button-text="选择图片"
    tip="支持 JPG、PNG、GIF 格式，单个文件不超过 10MB，最多上传 3 个"
    @change="handleChange"
    @success="handleSuccess"
    @error="handleError"
  />

  <GvButton @click="handleUpload">开始上传</GvButton>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { UploadFile, UploadFiles } from 'element-plus'
import GvUpload from '@/components/base/GvUpload/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'

const uploadRef = ref()
const fileList = ref<UploadFile[]>([])
const uploadUrl = '/api/v1/upload'

function handleChange(file: UploadFile, files: UploadFiles) {
  console.log('文件改变', file, files)
}

function handleSuccess(response: any, file: UploadFile) {
  console.log('上传成功', response)
  ElMessage.success('上传成功')
}

function handleError(error: Error, file: UploadFile) {
  console.error('上传失败', error)
  ElMessage.error('上传失败')
}

function handleUpload() {
  uploadRef.value?.submit()
}
</script>
```
