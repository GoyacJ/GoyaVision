# FilterBar - 筛选栏组件

业务组件，统一的筛选栏，支持多种筛选字段类型和灵活配置。

## 基本用法

```vue
<template>
  <FilterBar
    v-model="filters"
    :fields="filterFields"
    @filter="handleFilter"
  />
</template>

<script setup>
import { ref } from 'vue'
import { FilterBar } from '@/components'

const filters = ref({
  name: '',
  type: '',
  dateRange: []
})

const filterFields = [
  {
    key: 'name',
    label: '名称',
    type: 'input',
    placeholder: '搜索名称'
  },
  {
    key: 'type',
    label: '类型',
    type: 'select',
    placeholder: '选择类型',
    options: [
      { label: '类型 A', value: 'a' },
      { label: '类型 B', value: 'b' }
    ]
  },
  {
    key: 'dateRange',
    label: '日期范围',
    type: 'daterange',
    startPlaceholder: '开始日期',
    endPlaceholder: '结束日期'
  }
]

const handleFilter = (values) => {
  console.log('筛选条件:', values)
  // 调用 API 进行筛选
}
</script>
```

## 支持的字段类型

### 输入框（input）

```javascript
{
  key: 'keyword',
  label: '关键词',
  type: 'input',
  placeholder: '输入关键词'
}
```

### 选择器（select）

```javascript
{
  key: 'status',
  label: '状态',
  type: 'select',
  placeholder: '选择状态',
  options: [
    { label: '运行中', value: 'running' },
    { label: '已停止', value: 'stopped' }
  ]
}
```

### 日期范围（daterange）

```javascript
{
  key: 'dateRange',
  label: '日期范围',
  type: 'daterange',
  startPlaceholder: '开始日期',
  endPlaceholder: '结束日期'
}
```

### 日期（date）

```javascript
{
  key: 'createDate',
  label: '创建日期',
  type: 'date',
  placeholder: '选择日期'
}
```

### 日期时间（datetime）

```javascript
{
  key: 'createTime',
  label: '创建时间',
  type: 'datetime',
  placeholder: '选择日期时间'
}
```

## 可折叠

```vue
<FilterBar
  v-model="filters"
  :fields="filterFields"
  collapsible
  :default-expanded="false"
  @filter="handleFilter"
/>
```

## 自定义列数

```vue
<!-- 2 列布局 -->
<FilterBar
  v-model="filters"
  :fields="filterFields"
  :columns="2"
  @filter="handleFilter"
/>

<!-- 4 列布局 -->
<FilterBar
  v-model="filters"
  :fields="filterFields"
  :columns="4"
  @filter="handleFilter"
/>
```

## 加载状态

```vue
<template>
  <FilterBar
    v-model="filters"
    :fields="filterFields"
    :loading="isLoading"
    @filter="handleFilter"
  />
</template>

<script setup>
const isLoading = ref(false)

const handleFilter = async (values) => {
  isLoading.value = true
  try {
    await fetchData(values)
  } finally {
    isLoading.value = false
  }
}
</script>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| fields | `FilterField[]` | - | 筛选字段配置（必填） |
| modelValue | `Record<string, any>` | - | 筛选值（必填） |
| showReset | `boolean` | `true` | 是否显示重置按钮 |
| resetText | `string` | `'重置'` | 重置按钮文本 |
| filterText | `string` | `'筛选'` | 筛选按钮文本 |
| loading | `boolean` | `false` | 是否正在筛选 |
| collapsible | `boolean` | `false` | 是否可折叠 |
| defaultExpanded | `boolean` | `true` | 默认是否展开 |
| columns | `number` | `3` | 每行显示的字段数 |

## FilterField

| 属性 | 类型 | 说明 |
|------|------|------|
| key | `string` | 字段 key（必填） |
| label | `string` | 字段标签（必填） |
| type | `FilterFieldType` | 字段类型 |
| placeholder | `string` | 占位文本 |
| options | `Array` | 选项（select 类型） |
| defaultValue | `any` | 默认值 |
| startPlaceholder | `string` | 起始占位（daterange） |
| endPlaceholder | `string` | 结束占位（daterange） |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| update:modelValue | `(value)` | 值更新 |
| filter | `(value)` | 筛选事件 |
| reset | `()` | 重置事件 |

## 使用场景

### 资产列表筛选

```vue
<template>
  <GvContainer>
    <PageHeader title="资产管理">
      <template #actions>
        <GvButton icon="Plus" @click="handleAdd">新建</GvButton>
      </template>
    </PageHeader>
    
    <FilterBar
      v-model="filters"
      :fields="filterFields"
      :loading="isLoading"
      @filter="handleFilter"
      @reset="handleReset"
    />
    
    <GvTable
      :data="tableData"
      :columns="columns"
      :loading="isLoading"
    />
  </GvContainer>
</template>

<script setup>
const filters = ref({
  name: '',
  type: '',
  status: '',
  dateRange: []
})

const filterFields = [
  {
    key: 'name',
    label: '名称',
    type: 'input',
    placeholder: '搜索名称'
  },
  {
    key: 'type',
    label: '类型',
    type: 'select',
    placeholder: '选择类型',
    options: typeOptions
  },
  {
    key: 'status',
    label: '状态',
    type: 'select',
    placeholder: '选择状态',
    options: statusOptions
  },
  {
    key: 'dateRange',
    label: '创建时间',
    type: 'daterange'
  }
]

const handleFilter = async (values) => {
  isLoading.value = true
  try {
    tableData.value = await fetchAssets(values)
  } finally {
    isLoading.value = false
  }
}

const handleReset = () => {
  handleFilter(filters.value)
}
</script>
```

## 最佳实践

1. **字段数量**：
   - 常用筛选: 3-4 个字段
   - 高级筛选: 6-8 个字段
   - 过多字段使用折叠

2. **字段类型**：
   - 文本搜索使用 input
   - 固定选项使用 select
   - 时间范围使用 daterange

3. **默认值**：
   - 设置合理的默认值
   - 提升用户体验

4. **布局**：
   - 3 列适合大多数场景
   - 2 列适合字段较少
   - 4 列适合宽屏

5. **加载状态**：
   - 显示 loading 提升体验
   - 防止重复提交

6. **重置功能**：
   - 提供重置按钮
   - 重置后自动筛选
