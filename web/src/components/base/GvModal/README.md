# GvModal - 模态框组件

Material Design 3 风格的模态框组件，用于显示重要信息或需要用户交互的内容。

## 基本用法

```vue
<template>
  <div>
    <GvButton @click="visible = true">打开模态框</GvButton>
    
    <GvModal v-model="visible" title="模态框标题">
      <p>这是模态框的内容</p>
    </GvModal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { GvModal, GvButton } from '@/components'

const visible = ref(false)
</script>
```

## 尺寸

```vue
<GvModal v-model="visible" size="small" title="小尺寸">
  <p>内容</p>
</GvModal>

<GvModal v-model="visible" size="medium" title="中等尺寸">
  <p>内容</p>
</GvModal>

<GvModal v-model="visible" size="large" title="大尺寸">
  <p>内容</p>
</GvModal>

<GvModal v-model="visible" size="full" title="全屏">
  <p>内容</p>
</GvModal>
```

## 自定义头部和底部

```vue
<GvModal v-model="visible">
  <template #header>
    <div class="flex items-center gap-2">
      <el-icon><InfoFilled /></el-icon>
      <h3>自定义标题</h3>
    </div>
  </template>
  
  <p>模态框内容</p>
  
  <template #footer>
    <div class="flex justify-between">
      <GvButton variant="text">帮助</GvButton>
      <div class="flex gap-2">
        <GvButton variant="tonal" @click="visible = false">取消</GvButton>
        <GvButton variant="filled" @click="handleSave">保存</GvButton>
      </div>
    </div>
  </template>
</GvModal>
```

## 无底部按钮

```vue
<GvModal v-model="visible" title="信息提示" :show-footer="false">
  <p>这是一条信息提示</p>
</GvModal>
```

## 确认框

```vue
<GvModal
  v-model="visible"
  title="确认删除"
  confirm-text="删除"
  cancel-text="取消"
  :confirm-loading="loading"
  @confirm="handleDelete"
  @cancel="visible = false"
>
  <p>确定要删除这条记录吗？此操作不可撤销。</p>
</GvModal>
```

## 禁用遮罩点击关闭

```vue
<GvModal
  v-model="visible"
  title="重要操作"
  :close-on-click-modal="false"
>
  <p>必须完成操作才能关闭</p>
</GvModal>
```

## 禁用 ESC 键关闭

```vue
<GvModal
  v-model="visible"
  title="强制确认"
  :close-on-press-escape="false"
>
  <p>必须点击按钮才能关闭</p>
</GvModal>
```

## 居中显示

```vue
<GvModal v-model="visible" title="提示" center>
  <p>内容居中显示</p>
</GvModal>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| modelValue | `boolean` | - | 是否显示（必填） |
| title | `string` | - | 标题 |
| size | `'small' \| 'medium' \| 'large' \| 'full'` | `'medium'` | 尺寸 |
| showClose | `boolean` | `true` | 是否显示关闭按钮 |
| closeOnClickModal | `boolean` | `true` | 点击遮罩是否关闭 |
| closeOnPressEscape | `boolean` | `true` | 按 ESC 是否关闭 |
| destroyOnClose | `boolean` | `false` | 关闭时销毁子元素 |
| showFooter | `boolean` | `true` | 是否显示底部 |
| confirmText | `string` | `'确定'` | 确认按钮文字 |
| cancelText | `string` | `'取消'` | 取消按钮文字 |
| showConfirm | `boolean` | `true` | 是否显示确认按钮 |
| showCancel | `boolean` | `true` | 是否显示取消按钮 |
| confirmLoading | `boolean` | `false` | 确认按钮加载状态 |
| center | `boolean` | `false` | 是否居中显示 |
| customClass | `string` | - | 自定义类名 |
| zIndex | `number` | `1000` | z-index |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| update:modelValue | `(value: boolean)` | 更新显示状态 |
| open | `()` | 打开事件 |
| opened | `()` | 打开动画结束 |
| close | `()` | 关闭事件 |
| closed | `()` | 关闭动画结束 |
| confirm | `()` | 确认事件 |
| cancel | `()` | 取消事件 |

## Slots

| 插槽名 | 说明 |
|--------|------|
| header | 头部内容 |
| default | 主体内容 |
| footer | 底部内容 |

## 使用场景

### 表单提交

```vue
<GvModal
  v-model="visible"
  title="新建资产"
  :confirm-loading="loading"
  @confirm="handleSubmit"
>
  <el-form :model="form" label-width="80px">
    <el-form-item label="名称">
      <el-input v-model="form.name" />
    </el-form-item>
    <el-form-item label="类型">
      <el-select v-model="form.type" class="w-full">
        <el-option label="视频" value="video" />
        <el-option label="图片" value="image" />
      </el-select>
    </el-form-item>
  </el-form>
</GvModal>
```

### 确认对话框

```vue
<GvModal
  v-model="deleteVisible"
  title="确认删除"
  size="small"
  @confirm="handleDelete"
>
  <div class="text-center">
    <el-icon :size="48" color="#f59e0b" class="mb-4">
      <WarningFilled />
    </el-icon>
    <p class="text-base mb-2">确定要删除这条记录吗？</p>
    <p class="text-sm text-text-secondary">此操作不可撤销</p>
  </div>
</GvModal>
```

### 详情查看

```vue
<GvModal
  v-model="detailVisible"
  title="资产详情"
  size="large"
  :show-confirm="false"
  cancel-text="关闭"
>
  <div class="space-y-4">
    <div>
      <h4 class="font-semibold mb-2">基本信息</h4>
      <p>名称：{{ asset.name }}</p>
      <p>类型：{{ asset.type }}</p>
    </div>
    <div>
      <h4 class="font-semibold mb-2">详细信息</h4>
      <p>创建时间：{{ asset.createdAt }}</p>
      <p>文件大小：{{ asset.size }}</p>
    </div>
  </div>
</GvModal>
```

### 多步骤流程

```vue
<GvModal
  v-model="wizardVisible"
  title="创建工作流"
  size="large"
  :close-on-click-modal="false"
>
  <div class="min-h-[400px]">
    <!-- 步骤指示器 -->
    <el-steps :active="step" class="mb-6">
      <el-step title="选择算子" />
      <el-step title="配置参数" />
      <el-step title="确认创建" />
    </el-steps>
    
    <!-- 步骤内容 -->
    <div v-show="step === 0">步骤 1 内容</div>
    <div v-show="step === 1">步骤 2 内容</div>
    <div v-show="step === 2">步骤 3 内容</div>
  </div>
  
  <template #footer>
    <div class="flex justify-between">
      <GvButton
        v-if="step > 0"
        variant="text"
        @click="step--"
      >
        上一步
      </GvButton>
      <div class="ml-auto flex gap-2">
        <GvButton variant="tonal" @click="wizardVisible = false">
          取消
        </GvButton>
        <GvButton
          v-if="step < 2"
          variant="filled"
          @click="step++"
        >
          下一步
        </GvButton>
        <GvButton
          v-else
          variant="filled"
          :loading="loading"
          @click="handleCreate"
        >
          创建
        </GvButton>
      </div>
    </div>
  </template>
</GvModal>
```

## 最佳实践

1. **尺寸选择**：
   - small: 简单确认框
   - medium: 普通表单
   - large: 复杂内容
   - full: 需要大量空间的内容

2. **关闭控制**：
   - 重要操作禁用遮罩点击关闭
   - 多步骤流程禁用 ESC 关闭

3. **加载状态**：
   - 异步操作时使用 confirmLoading

4. **事件处理**：
   - 使用 @confirm 处理确认逻辑
   - 使用 @cancel 处理取消逻辑

5. **用户体验**：
   - 提供清晰的标题和说明
   - 使用合适的按钮文本
   - 重要操作添加二次确认

6. **性能优化**：
   - 复杂内容使用 destroyOnClose
   - 避免在模态框中进行大量计算
