<template>
  <GvCard class="filter-bar" shadow="sm" padding="md">
    <!-- 筛选表单 - 横向布局 -->
    <GvFlex align="center" gap="md" wrap>
      <!-- 动态渲染筛选字段 -->
      <template v-for="field in fields" :key="field.key">
        <!-- 输入框 -->
        <div v-if="field.type === 'input' || !field.type" class="filter-item">
          <GvInput
            :model-value="localValue[field.key]"
            :placeholder="field.placeholder || field.label"
            size="small"
            @update:model-value="handleFieldChange(field.key, $event)"
          />
        </div>

        <!-- 选择器 -->
        <div v-else-if="field.type === 'select'" class="filter-item">
          <GvSelect
            :model-value="localValue[field.key]"
            :placeholder="field.placeholder || field.label"
            :options="field.options || []"
            size="small"
            @update:model-value="handleFieldChange(field.key, $event)"
          />
        </div>

        <!-- 日期范围 -->
        <div v-else-if="field.type === 'daterange'" class="filter-item-wide">
          <GvDatePicker
            :model-value="localValue[field.key]"
            type="daterange"
            :start-placeholder="field.startPlaceholder || '开始日期'"
            :end-placeholder="field.endPlaceholder || '结束日期'"
            size="small"
            @update:model-value="handleFieldChange(field.key, $event)"
          />
        </div>

        <!-- 日期 -->
        <div v-else-if="field.type === 'date'" class="filter-item">
          <GvDatePicker
            :model-value="localValue[field.key]"
            type="date"
            :placeholder="field.placeholder || field.label"
            size="small"
            @update:model-value="handleFieldChange(field.key, $event)"
          />
        </div>

        <!-- 日期时间 -->
        <div v-else-if="field.type === 'datetime'" class="filter-item">
          <GvDatePicker
            :model-value="localValue[field.key]"
            type="datetime"
            :placeholder="field.placeholder || field.label"
            size="small"
            @update:model-value="handleFieldChange(field.key, $event)"
          />
        </div>
      </template>

      <!-- 操作按钮 -->
      <div class="filter-actions">
        <GvSpace size="xs">
          <GvButton
            size="small"
            :loading="loading"
            @click="handleFilter"
          >
            {{ filterText }}
          </GvButton>
          <GvButton
            v-if="showReset"
            variant="text"
            size="small"
            @click="handleReset"
          >
            {{ resetText }}
          </GvButton>
        </GvSpace>
      </div>
    </GvFlex>
  </GvCard>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  GvCard,
  GvFlex,
  GvGrid,
  GvButton,
  GvSpace,
  GvInput,
  GvSelect,
  GvDatePicker
} from '@/components'
import { ArrowUp, ArrowDown } from '@element-plus/icons-vue'
import type { FilterBarProps, FilterBarEmits } from './types'

const props = withDefaults(defineProps<FilterBarProps>(), {
  showReset: true,
  resetText: '重置',
  filterText: '筛选',
  loading: false,
  collapsible: false,
  defaultExpanded: true,
  columns: 3
})

const emit = defineEmits<FilterBarEmits>()

// 本地值
const localValue = ref<Record<string, any>>({ ...props.modelValue })

// 展开状态
const expanded = ref(props.defaultExpanded)

// 监听外部值变化
watch(
  () => props.modelValue,
  (newValue) => {
    localValue.value = { ...newValue }
  },
  { deep: true }
)

// 字段变化处理
const handleFieldChange = (key: string, value: any) => {
  localValue.value[key] = value
  emit('update:modelValue', localValue.value)
}

// 筛选处理
const handleFilter = () => {
  emit('filter', localValue.value)
}

// 重置处理
const handleReset = () => {
  // 重置为默认值
  const resetValue: Record<string, any> = {}
  props.fields.forEach(field => {
    resetValue[field.key] = field.defaultValue ?? ''
  })
  
  localValue.value = resetValue
  emit('update:modelValue', resetValue)
  emit('reset')
}

// 切换展开/收起
const toggleExpanded = () => {
  expanded.value = !expanded.value
}
</script>

<style scoped>
.filter-bar {
  @apply mb-4;
}

.filter-item {
  @apply min-w-[180px];
}

.filter-item-wide {
  @apply min-w-[280px];
}

.filter-actions {
  @apply ml-auto;
}
</style>
