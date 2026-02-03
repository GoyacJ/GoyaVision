<template>
  <GvCard class="filter-bar">
    <!-- 折叠按钮 -->
    <GvFlex v-if="collapsible" justify="between" align="center" class="mb-4">
      <h3 class="text-base font-semibold">筛选条件</h3>
      <GvButton variant="text" size="small" @click="toggleExpanded">
        {{ expanded ? '收起' : '展开' }}
        <template #suffix>
          <el-icon>
            <ArrowUp v-if="expanded" />
            <ArrowDown v-else />
          </el-icon>
        </template>
      </GvButton>
    </GvFlex>
    
    <!-- 筛选表单 -->
    <div v-show="!collapsible || expanded">
      <GvGrid :cols="columns" gap="lg" class="mb-4">
        <!-- 动态渲染筛选字段 -->
        <template v-for="field in fields" :key="field.key">
          <!-- 输入框 -->
          <GvInput
            v-if="field.type === 'input' || !field.type"
            :model-value="localValue[field.key]"
            :label="field.label"
            :placeholder="field.placeholder"
            @update:model-value="handleFieldChange(field.key, $event)"
          />
          
          <!-- 选择器 -->
          <GvSelect
            v-else-if="field.type === 'select'"
            :model-value="localValue[field.key]"
            :label="field.label"
            :placeholder="field.placeholder"
            :options="field.options || []"
            @update:model-value="handleFieldChange(field.key, $event)"
          />
          
          <!-- 日期范围 -->
          <GvDatePicker
            v-else-if="field.type === 'daterange'"
            :model-value="localValue[field.key]"
            type="daterange"
            :label="field.label"
            :start-placeholder="field.startPlaceholder || '开始日期'"
            :end-placeholder="field.endPlaceholder || '结束日期'"
            @update:model-value="handleFieldChange(field.key, $event)"
          />
          
          <!-- 日期 -->
          <GvDatePicker
            v-else-if="field.type === 'date'"
            :model-value="localValue[field.key]"
            type="date"
            :label="field.label"
            :placeholder="field.placeholder"
            @update:model-value="handleFieldChange(field.key, $event)"
          />
          
          <!-- 日期时间 -->
          <GvDatePicker
            v-else-if="field.type === 'datetime'"
            :model-value="localValue[field.key]"
            type="datetime"
            :label="field.label"
            :placeholder="field.placeholder"
            @update:model-value="handleFieldChange(field.key, $event)"
          />
        </template>
      </GvGrid>
      
      <!-- 操作按钮 -->
      <GvFlex justify="end">
        <GvSpace>
          <GvButton
            v-if="showReset"
            variant="text"
            @click="handleReset"
          >
            {{ resetText }}
          </GvButton>
          <GvButton
            :loading="loading"
            @click="handleFilter"
          >
            {{ filterText }}
          </GvButton>
        </GvSpace>
      </GvFlex>
    </div>
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
