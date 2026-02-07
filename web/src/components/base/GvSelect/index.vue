<template>
  <div :class="containerClasses">
    <!-- 标签 -->
    <label v-if="label" :class="labelClasses">
      {{ label }}
      <span v-if="required" class="text-error-600 ml-0.5">*</span>
    </label>
    
    <!-- 选择器 -->
    <el-select
      :model-value="modelValue"
      :class="selectClasses"
      :size="size"
      :placeholder="placeholder"
      :disabled="disabled"
      :clearable="clearable"
      :multiple="multiple"
      :multiple-limit="multipleLimit"
      :filterable="filterable"
      :allow-create="allowCreate"
      :filter-method="filterMethod"
      :remote="remote"
      :remote-method="remoteMethod"
      :loading="loading"
      :loading-text="loadingText"
      :no-data-text="noDataText"
      :no-match-text="noMatchText"
      :popper-class="cn('gv-select-dropdown', popperClass)"
      @update:model-value="handleUpdate"
      @change="handleChange"
      @visible-change="handleVisibleChange"
      @remove-tag="handleRemoveTag"
      @clear="handleClear"
      @focus="handleFocus"
      @blur="handleBlur"
    >
      <template #suffix>
        <div class="flex items-center gap-1">
          <el-icon 
            v-if="clearable && hasValue" 
            class="cursor-pointer text-text-tertiary hover:text-text-secondary transition-colors"
            @click.stop="handleCustomClear"
          >
            <CircleClose />
          </el-icon>
          <el-icon class="text-text-tertiary"><ArrowDown /></el-icon>
        </div>
      </template>
      <el-option
        v-for="option in options"
        :key="option[valueKey]"
        :label="option[labelKey]"
        :value="option[valueKey]"
        :disabled="option.disabled"
      />
    </el-select>
    
    <!-- 验证状态图标 -->
    <div v-if="status" :class="statusIconClasses">
      <el-icon>
        <SuccessFilled v-if="status === 'success'" />
        <CircleCloseFilled v-if="status === 'error'" />
        <WarningFilled v-if="status === 'warning'" />
      </el-icon>
    </div>
    
    <!-- 错误提示 -->
    <div v-if="errorMessage && status === 'error'" class="gv-select__error">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'
import type { SelectProps, SelectEmits } from './types'

const props = withDefaults(defineProps<SelectProps>(), {
  size: 'medium',
  disabled: false,
  clearable: true,
  multiple: false,
  filterable: false,
  allowCreate: false,
  remote: false,
  loading: false,
  required: false,
  labelKey: 'label',
  valueKey: 'value',
  options: () => []
})

const emit = defineEmits<SelectEmits>()

// 容器类名
const containerClasses = computed(() => {
  return cn('gv-select', 'w-full', 'relative')
})

// 标签类名
const labelClasses = computed(() => {
  return cn(
    'gv-select__label',
    'block mb-1.5 text-sm font-medium text-text-primary'
  )
})

// 选择器类名
const selectClasses = computed(() => {
  const base = ['gv-select__wrapper', 'w-full']
  
  // 状态
  const stateClasses = {
    success: 'gv-select--success',
    error: 'gv-select--error',
    warning: 'gv-select--warning'
  }
  
  const stateClass = props.status ? stateClasses[props.status] : ''
  
  return cn(base, stateClass)
})

// 状态图标类名
const statusIconClasses = computed(() => {
  const base = [
    'gv-select__status',
    'absolute right-10 top-1/2 -translate-y-1/2',
    'pointer-events-none'
  ]
  
  const colorClasses = {
    success: 'text-success-600',
    error: 'text-error-600',
    warning: 'text-warning-600'
  }
  
  const colorClass = props.status ? colorClasses[props.status] : ''
  
  // 如果有标签，需要调整位置
  const topClass = props.label ? 'top-[calc(50%+12px)]' : 'top-1/2'
  
  return cn(base, colorClass, topClass)
})

// 更新事件
const handleUpdate = (value: string | number | Array<string | number> | undefined) => {
  emit('update:modelValue', value)
}

// 变化事件
const handleChange = (value: string | number | Array<string | number> | undefined) => {
  emit('change', value)
}

// 下拉框显示/隐藏
const handleVisibleChange = (visible: boolean) => {
  emit('visible-change', visible)
}

// 移除标签
const handleRemoveTag = (value: string | number) => {
  emit('remove-tag', value)
}

// 清空
const handleClear = () => {
  emit('clear')
}

const hasValue = computed(() => {
  if (Array.isArray(props.modelValue)) {
    return props.modelValue.length > 0
  }
  return props.modelValue !== undefined && props.modelValue !== null && props.modelValue !== ''
})

const handleCustomClear = () => {
  const emptyValue = props.multiple ? [] : undefined
  emit('update:modelValue', emptyValue)
  emit('change', emptyValue)
  emit('clear')
}

// 聚焦
const handleFocus = (event: FocusEvent) => {
  emit('focus', event)
}

// 失焦
const handleBlur = (event: FocusEvent) => {
  emit('blur', event)
}
</script>

<style>
/* 自定义 Element Plus Select 样式 */
.gv-select__wrapper .el-select__wrapper {
  @apply bg-white border border-neutral-300 rounded-xl transition-all;
  @apply hover:border-neutral-400;
  box-shadow: none !important;
  padding: 0 12px;
}

.gv-select__wrapper.el-select--small .el-select__wrapper {
  @apply h-8 text-sm rounded-lg;
}

.gv-select__wrapper.el-select--default .el-select__wrapper {
  @apply h-10 text-base rounded-xl;
}

.gv-select__wrapper.el-select--large .el-select__wrapper {
  @apply h-12 text-lg rounded-2xl;
}

.gv-select__wrapper .el-select__wrapper.is-focused {
  @apply border-primary-600 ring-4 ring-primary-50;
  box-shadow: none !important;
}

.gv-select__wrapper.gv-select--success .el-select__wrapper,
.gv-select__wrapper.gv-select--success .el-select__wrapper.is-focused {
  @apply border-success-600;
}

.gv-select__wrapper.gv-select--error .el-select__wrapper,
.gv-select__wrapper.gv-select--error .el-select__wrapper.is-focused {
  @apply border-error-600;
}

.gv-select__wrapper.gv-select--warning .el-select__wrapper,
.gv-select__wrapper.gv-select--warning .el-select__wrapper.is-focused {
  @apply border-warning-600;
}

.gv-select__wrapper .el-select__wrapper.is-disabled {
  @apply bg-neutral-100 border-neutral-200 cursor-not-allowed;
}

.gv-select__wrapper .el-select__placeholder {
  @apply text-text-tertiary;
}

.gv-select__wrapper .el-select__input {
  @apply text-text-primary;
}

/* 下拉框样式 */
.gv-select-dropdown {
  @apply rounded-xl shadow-lg border-none;
  margin-top: 4px;
}

.gv-select-dropdown .el-select-dropdown__item {
  @apply px-4 py-2 text-text-primary flex items-center justify-start;
  @apply transition-colors duration-150;
  height: auto;
  min-height: 34px;
  line-height: normal;
}

.gv-select-dropdown .el-select-dropdown__item:hover {
  @apply bg-primary-50;
}

.gv-select-dropdown .el-select-dropdown__item.is-selected {
  @apply bg-primary-100 text-primary-700 font-medium;
}

.gv-select-dropdown .el-select-dropdown__item.is-disabled {
  @apply text-text-disabled cursor-not-allowed;
  @apply hover:bg-transparent;
}

/* 错误提示 */
.gv-select__error {
  margin-top: 4px;
  font-size: 0.75rem;
  color: rgb(239 68 68);
  line-height: 1.4;
}

/* 深色模式 */
.dark .gv-select__wrapper .el-select__wrapper {
  @apply bg-surface-dark border-neutral-700;
}

.dark .gv-select__wrapper .el-select__input {
  @apply text-text-inverse;
}

.dark .gv-select__label {
  @apply text-text-inverse;
}

.dark .gv-select-dropdown {
  @apply bg-surface-dark border-neutral-700;
}

.dark .gv-select-dropdown .el-select-dropdown__item {
  @apply text-text-inverse;
}

.dark .gv-select-dropdown .el-select-dropdown__item:hover {
  @apply bg-primary-950/30;
}

.dark .gv-select-dropdown .el-select-dropdown__item.is-selected {
  @apply bg-primary-900/50 text-primary-300;
}
</style>
