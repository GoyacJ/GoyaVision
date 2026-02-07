<template>
  <div :class="containerClasses">
    <!-- 标签 -->
    <label v-if="label" :class="labelClasses">
      {{ label }}
      <span v-if="required" class="text-error-600 ml-0.5">*</span>
    </label>
    
    <!-- 输入框容器 -->
    <div :class="wrapperClasses">
      <!-- 前置图标 -->
      <span v-if="prefixIcon || $slots.prefix" class="gv-input__prefix">
        <slot name="prefix">
          <el-icon v-if="prefixIcon">
            <component :is="prefixIcon" />
          </el-icon>
        </slot>
      </span>
      
      <!-- 输入框 -->
      <component
        :is="type === 'textarea' ? 'textarea' : 'input'"
        ref="inputRef"
        :class="[inputClasses, type === 'textarea' ? 'h-auto py-2 resize-none' : '']"
        :type="type === 'textarea' ? undefined : computedType"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        :maxlength="maxlength"
        :autocomplete="autocomplete"
        :name="name"
        :form="form"
        :autofocus="autofocus"
        :rows="rows"
        @input="handleInput"
        @change="handleChange"
        @focus="handleFocus"
        @blur="handleBlur"
        @keydown="handleKeydown"
        @keyup="handleKeyup"
      />
      
      <!-- 字数统计 -->
      <span v-if="showCount && maxlength" class="gv-input__count">
        {{ String(modelValue || '').length }}/{{ maxlength }}
      </span>
      
      <!-- 清除按钮 -->
      <span
        v-if="clearable && modelValue && !disabled && !readonly"
        class="gv-input__clear"
        @click="handleClear"
      >
        <el-icon><CircleClose /></el-icon>
      </span>
      
      <!-- 密码切换 -->
      <span
        v-if="showPassword && type === 'password'"
        class="gv-input__password"
        @click="togglePasswordVisible"
      >
        <el-icon>
          <component :is="passwordVisible ? 'View' : 'Hide'" />
        </el-icon>
      </span>
      
      <!-- 后置图标 -->
      <span v-if="suffixIcon || $slots.suffix" class="gv-input__suffix">
        <slot name="suffix">
          <el-icon v-if="suffixIcon">
            <component :is="suffixIcon" />
          </el-icon>
        </slot>
      </span>
      
      <!-- 验证状态图标 -->
      <span v-if="status && !suffixIcon && !$slots.suffix" class="gv-input__status">
        <el-icon>
          <SuccessFilled v-if="status === 'success'" />
          <CircleCloseFilled v-if="status === 'error'" />
          <WarningFilled v-if="status === 'warning'" />
        </el-icon>
      </span>
    </div>
    
    <!-- 错误提示 -->
    <div v-if="errorMessage && status === 'error'" class="gv-input__error">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { cn } from '@/utils/cn'
import type { InputProps, InputEmits } from './types'

const props = withDefaults(defineProps<InputProps>(), {
  type: 'text',
  size: 'medium',
  disabled: false,
  readonly: false,
  clearable: false,
  showPassword: false,
  showCount: false,
  autofocus: false,
  required: false,
  rows: 2
})

const emit = defineEmits<InputEmits>()

const inputRef = ref<HTMLInputElement>()
const passwordVisible = ref(false)
const isFocused = ref(false)

// 计算的输入框类型
const computedType = computed(() => {
  if (props.type === 'password' && props.showPassword) {
    return passwordVisible.value ? 'text' : 'password'
  }
  return props.type
})

// 容器类名
const containerClasses = computed(() => {
  return cn('gv-input', 'w-full')
})

// 标签类名
const labelClasses = computed(() => {
  return cn(
    'gv-input__label',
    'block mb-1.5 text-sm font-medium text-text-primary'
  )
})

// 输入框容器类名
const wrapperClasses = computed(() => {
  const base = [
    'gv-input__wrapper',
    'relative flex',
    'bg-white border rounded-xl'
  ]

  if (props.type === 'textarea') {
    base.push('items-start')
  } else {
    base.push('items-center')
  }

  // 尺寸
  const sizeClasses = {
    small: 'h-8',
    medium: 'h-10',
    large: 'h-12'
  }

  const sizeClass = props.type === 'textarea' ? 'h-auto' : sizeClasses[props.size]

  // 状态
  const stateClasses = {
    default: 'border-neutral-300 hover:border-neutral-400 transition-colors',
    focused: 'border-primary-600 ring-4 ring-primary-50 transition-all',
    disabled: 'bg-neutral-100 border-neutral-200 cursor-not-allowed',
    readonly: 'bg-neutral-50 border-neutral-200',
    success: 'border-success-600',
    error: 'border-error-600',
    warning: 'border-warning-600'
  }
  
  let stateClass = stateClasses.default
  if (props.disabled) {
    stateClass = stateClasses.disabled
  } else if (props.readonly) {
    stateClass = stateClasses.readonly
  } else if (props.status) {
    stateClass = stateClasses[props.status]
  } else if (isFocused.value) {
    stateClass = stateClasses.focused
  }
  
  return cn(base, sizeClass, stateClass)
})

// 输入框类名
const inputClasses = computed(() => {
  const base = [
    'gv-input__inner',
    'flex-1 px-3 bg-transparent',
    'text-text-primary placeholder:text-text-tertiary',
    'outline-none border-none',
    'disabled:cursor-not-allowed disabled:text-text-disabled'
  ]
  
  // 尺寸
  const sizeClasses = {
    small: 'text-sm',
    medium: 'text-base',
    large: 'text-lg'
  }
  
  return cn(base, sizeClasses[props.size])
})

// 输入事件
const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  const value = target.value
  emit('update:modelValue', value)
  emit('input', value)
}

// 变化事件
const handleChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  const value = target.value
  emit('change', value)
}

// 聚焦事件
const handleFocus = (event: FocusEvent) => {
  isFocused.value = true
  emit('focus', event)
}

// 失焦事件
const handleBlur = (event: FocusEvent) => {
  isFocused.value = false
  emit('blur', event)
}

// 清除事件
const handleClear = () => {
  emit('update:modelValue', '')
  emit('clear')
  inputRef.value?.focus()
}

// 切换密码可见性
const togglePasswordVisible = () => {
  passwordVisible.value = !passwordVisible.value
}

// 按键事件
const handleKeydown = (event: KeyboardEvent) => {
  emit('keydown', event)
}

const handleKeyup = (event: KeyboardEvent) => {
  emit('keyup', event)
}

// 暴露方法
defineExpose({
  focus: () => inputRef.value?.focus(),
  blur: () => inputRef.value?.blur(),
  select: () => inputRef.value?.select()
})
</script>

<style scoped>
.gv-input__prefix,
.gv-input__suffix,
.gv-input__clear,
.gv-input__password,
.gv-input__status,
.gv-input__count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: var(--el-text-color-secondary);
  padding: 0 8px;
}

.gv-input__clear,
.gv-input__password {
  cursor: pointer;
  transition: color 0.2s;
}

.gv-input__clear:hover,
.gv-input__password:hover {
  color: var(--el-text-color-primary);
}

.gv-input__count {
  font-size: 0.75rem;
  color: var(--el-text-color-placeholder);
}

.gv-input__error {
  margin-top: 4px;
  font-size: 0.75rem;
  color: rgb(239 68 68);
  line-height: 1.4;
}

.gv-input__status .el-icon {
  font-size: 1.125rem;
}

.gv-input__status .el-icon.success {
  color: rgb(16 185 129);
}

.gv-input__status .el-icon.error {
  color: rgb(239 68 68);
}

.gv-input__status .el-icon.warning {
  color: rgb(245 158 11);
}

.dark .gv-input__wrapper {
  @apply bg-surface-dark border-neutral-700;
}

.dark .gv-input__inner {
  @apply text-text-inverse;
}

.dark .gv-input__label {
  @apply text-text-inverse;
}

/* 强制移除默认 outline 和 border */
.gv-input__inner {
  outline: none !important;
  border: none !important;
  box-shadow: none !important;
}

.gv-input__inner:focus {
  outline: none !important;
  border: none !important;
  box-shadow: none !important;
}
</style>
