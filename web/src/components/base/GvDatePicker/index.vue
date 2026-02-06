<template>
  <div :class="containerClasses">
    <!-- 标签 -->
    <label v-if="label" :class="labelClasses">
      {{ label }}
      <span v-if="required" class="text-error-600 ml-0.5">*</span>
    </label>
    
    <!-- 日期选择器 -->
    <el-date-picker
      :model-value="modelValue"
      :class="datePickerClasses"
      :type="type"
      :size="size"
      :placeholder="placeholder"
      :start-placeholder="startPlaceholder"
      :end-placeholder="endPlaceholder"
      :disabled="disabled"
      :clearable="clearable"
      :format="format"
      :value-format="valueFormat"
      :disabled-date="disabledDate"
      :shortcuts="shortcuts"
      :range-separator="rangeSeparator"
      :default-time="defaultTime"
      @update:model-value="handleUpdate"
      @change="handleChange"
      @clear="handleClear"
      @focus="handleFocus"
      @blur="handleBlur"
    />
    
    <!-- 验证状态图标 -->
    <div v-if="status" :class="statusIconClasses">
      <el-icon>
        <SuccessFilled v-if="status === 'success'" />
        <CircleCloseFilled v-if="status === 'error'" />
        <WarningFilled v-if="status === 'warning'" />
      </el-icon>
    </div>
    
    <!-- 错误提示 -->
    <div v-if="errorMessage && status === 'error'" class="gv-date-picker__error">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'
import type { DatePickerProps, DatePickerEmits } from './types'

const props = withDefaults(defineProps<DatePickerProps>(), {
  type: 'date',
  size: 'default',
  disabled: false,
  clearable: true,
  required: false,
  rangeSeparator: '-'
})

const emit = defineEmits<DatePickerEmits>()

// 容器类名
const containerClasses = computed(() => {
  return cn('gv-date-picker', 'w-full', 'relative')
})

// 标签类名
const labelClasses = computed(() => {
  return cn(
    'gv-date-picker__label',
    'block mb-1.5 text-sm font-medium text-text-primary'
  )
})

// 日期选择器类名
const datePickerClasses = computed(() => {
  const base = ['gv-date-picker__wrapper', 'w-full']
  
  // 状态
  const stateClasses = {
    success: 'gv-date-picker--success',
    error: 'gv-date-picker--error',
    warning: 'gv-date-picker--warning'
  }
  
  const stateClass = props.status ? stateClasses[props.status] : ''
  
  return cn(base, stateClass)
})

// 状态图标类名
const statusIconClasses = computed(() => {
  const base = [
    'gv-date-picker__status',
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
const handleUpdate = (value: any) => {
  emit('update:modelValue', value)
}

// 变化事件
const handleChange = (value: any) => {
  emit('change', value)
}

// 清空事件
const handleClear = () => {
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
/* 自定义 Element Plus DatePicker 样式 */
.gv-date-picker__wrapper .el-input__wrapper {
  @apply bg-white border border-neutral-300 rounded-xl;
  box-shadow: none !important;
  padding: 0 12px;
}

.gv-date-picker__wrapper.el-date-editor--small .el-input__wrapper {
  @apply h-8 text-sm rounded-lg;
}

.gv-date-picker__wrapper.el-date-editor--default .el-input__wrapper {
  @apply h-10 text-base rounded-xl;
}

.gv-date-picker__wrapper.el-date-editor--large .el-input__wrapper {
  @apply h-12 text-lg rounded-2xl;
}

.gv-date-picker__wrapper .el-input__wrapper.is-focus {
  @apply border-neutral-300;
  box-shadow: none !important;
}

.gv-date-picker__wrapper.gv-date-picker--success .el-input__wrapper,
.gv-date-picker__wrapper.gv-date-picker--success .el-input__wrapper.is-focus {
  @apply border-success-600;
}

.gv-date-picker__wrapper.gv-date-picker--error .el-input__wrapper,
.gv-date-picker__wrapper.gv-date-picker--error .el-input__wrapper.is-focus {
  @apply border-error-600;
}

.gv-date-picker__wrapper.gv-date-picker--warning .el-input__wrapper,
.gv-date-picker__wrapper.gv-date-picker--warning .el-input__wrapper.is-focus {
  @apply border-warning-600;
}

.gv-date-picker__wrapper .el-input__wrapper.is-disabled {
  @apply bg-neutral-100 border-neutral-200 cursor-not-allowed;
}

.gv-date-picker__wrapper .el-input__inner {
  @apply text-text-primary;
}

.gv-date-picker__wrapper .el-input__prefix,
.gv-date-picker__wrapper .el-input__suffix {
  @apply text-text-tertiary;
}

/* 下拉面板样式 */
.el-picker-panel {
  @apply rounded-xl border border-neutral-200 shadow-lg;
}

.el-picker-panel__body {
  @apply p-4;
}

.el-date-picker__header {
  @apply mb-3;
}

.el-date-table td.available:hover {
  @apply bg-primary-50;
}

.el-date-table td.current:not(.disabled) {
  @apply bg-primary-600 text-white;
}

.el-date-table td.today .el-date-table-cell__text {
  @apply text-primary-600 font-semibold;
}

/* 快捷选项样式 */
.el-picker-panel__sidebar {
  @apply border-r border-neutral-200;
}

.el-picker-panel__shortcut {
  @apply px-4 py-2 text-sm text-text-primary;
  @apply transition-colors duration-150;
  @apply hover:bg-primary-50 hover:text-primary-600;
}

/* 错误提示 */
.gv-date-picker__error {
  margin-top: 4px;
  font-size: 0.75rem;
  color: rgb(239 68 68);
  line-height: 1.4;
}

/* 深色模式 */
.dark .gv-date-picker__wrapper .el-input__wrapper {
  @apply bg-surface-dark border-neutral-700;
}

.dark .gv-date-picker__wrapper .el-input__inner {
  @apply text-text-inverse;
}

.dark .gv-date-picker__label {
  @apply text-text-inverse;
}

.dark .el-picker-panel {
  @apply bg-surface-dark border-neutral-700;
}

.dark .el-date-picker__header {
  @apply text-text-inverse;
}

.dark .el-date-table td.available:hover {
  @apply bg-primary-950/30;
}

.dark .el-date-table td.current:not(.disabled) {
  @apply bg-primary-700 text-white;
}

.dark .el-date-table td.today .el-date-table-cell__text {
  @apply text-primary-400;
}

.dark .el-picker-panel__sidebar {
  @apply border-neutral-700;
}

.dark .el-picker-panel__shortcut {
  @apply text-text-inverse;
  @apply hover:bg-primary-950/30 hover:text-primary-300;
}
</style>
