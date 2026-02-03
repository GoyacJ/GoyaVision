<template>
  <div :class="containerClasses">
    <GvInput
      :model-value="modelValue"
      :placeholder="placeholder"
      :size="size"
      :disabled="disabled"
      :clearable="clearable"
      :class="inputClasses"
      @update:model-value="handleInput"
      @clear="handleClear"
      @keyup.enter="handleSearch"
    >
      <template #prefix>
        <el-icon>
          <Search />
        </el-icon>
      </template>
    </GvInput>
    
    <GvButton
      v-if="showButton"
      :size="size"
      :loading="loading"
      :disabled="disabled || !modelValue"
      @click="handleSearch"
    >
      {{ buttonText }}
    </GvButton>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import { cn } from '@/utils/cn'
import { GvInput, GvButton } from '@/components'
import { Search } from '@element-plus/icons-vue'
import type { SearchBarProps, SearchBarEmits } from './types'
import { useDebounceFn } from '@vueuse/core'

const props = withDefaults(defineProps<SearchBarProps>(), {
  placeholder: '搜索...',
  showButton: true,
  buttonText: '搜索',
  loading: false,
  disabled: false,
  size: 'default',
  clearable: true,
  immediate: false,
  debounce: 300
})

const emit = defineEmits<SearchBarEmits>()

// 容器类名
const containerClasses = computed(() => {
  return cn('search-bar', 'flex gap-3 w-full')
})

// 输入框类名
const inputClasses = computed(() => {
  return cn('search-bar__input', props.showButton ? 'flex-1' : 'w-full')
})

// 防抖搜索
const debouncedSearch = useDebounceFn((value: string) => {
  if (props.immediate) {
    emit('search', value)
  }
}, props.debounce)

// 输入处理
const handleInput = (value: string) => {
  emit('update:modelValue', value)
  
  // 如果启用立即搜索，触发防抖搜索
  if (props.immediate) {
    debouncedSearch(value)
  }
}

// 搜索处理
const handleSearch = () => {
  emit('search', props.modelValue)
}

// 清空处理
const handleClear = () => {
  emit('update:modelValue', '')
  emit('clear')
  
  // 如果启用立即搜索，清空时也触发搜索
  if (props.immediate) {
    emit('search', '')
  }
}
</script>
