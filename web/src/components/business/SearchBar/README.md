# SearchBar - 搜索栏组件

业务组件，统一的搜索栏，基于 `GvInput` 和 `GvButton` 封装。

## 基本用法

```vue
<template>
  <SearchBar
    v-model="searchText"
    @search="handleSearch"
  />
</template>

<script setup>
import { ref } from 'vue'
import { SearchBar } from '@/components'

const searchText = ref('')

const handleSearch = (value: string) => {
  console.log('搜索:', value)
}
</script>
```

## 不显示按钮

```vue
<SearchBar
  v-model="searchText"
  :show-button="false"
  @search="handleSearch"
/>
```

## 立即搜索

```vue
<!-- 输入时自动搜索（防抖 300ms） -->
<SearchBar
  v-model="searchText"
  immediate
  @search="handleSearch"
/>

<!-- 自定义防抖时间 -->
<SearchBar
  v-model="searchText"
  immediate
  :debounce="500"
  @search="handleSearch"
/>
```

## 加载状态

```vue
<template>
  <SearchBar
    v-model="searchText"
    :loading="isSearching"
    @search="handleSearch"
  />
</template>

<script setup>
const isSearching = ref(false)

const handleSearch = async (value: string) => {
  isSearching.value = true
  try {
    await searchAPI(value)
  } finally {
    isSearching.value = false
  }
}
</script>
```

## 尺寸

```vue
<SearchBar v-model="searchText" size="small" />
<SearchBar v-model="searchText" size="default" />
<SearchBar v-model="searchText" size="large" />
```

## 自定义占位文本和按钮

```vue
<SearchBar
  v-model="searchText"
  placeholder="输入资产名称"
  button-text="查询"
  @search="handleSearch"
/>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| modelValue | `string` | - | 搜索值（必填） |
| placeholder | `string` | `'搜索...'` | 占位文本 |
| showButton | `boolean` | `true` | 是否显示搜索按钮 |
| buttonText | `string` | `'搜索'` | 搜索按钮文本 |
| loading | `boolean` | `false` | 是否正在搜索 |
| disabled | `boolean` | `false` | 是否禁用 |
| size | `'small' \| 'default' \| 'large'` | `'default'` | 尺寸 |
| clearable | `boolean` | `true` | 是否可清空 |
| immediate | `boolean` | `false` | 是否立即搜索 |
| debounce | `number` | `300` | 防抖延迟（毫秒） |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| update:modelValue | `(value: string)` | 值更新 |
| search | `(value: string)` | 搜索事件 |
| clear | `()` | 清空事件 |

## 使用场景

### 列表搜索

```vue
<template>
  <GvCard>
    <SearchBar
      v-model="searchText"
      placeholder="搜索资产名称"
      immediate
      @search="handleSearch"
    />
    
    <GvTable
      :data="filteredData"
      :columns="columns"
      :loading="loading"
    />
  </GvCard>
</template>

<script setup>
const searchText = ref('')
const loading = ref(false)
const tableData = ref([])

const filteredData = computed(() => {
  if (!searchText.value) return tableData.value
  return tableData.value.filter(item => 
    item.name.includes(searchText.value)
  )
})

const handleSearch = async (value: string) => {
  loading.value = true
  try {
    tableData.value = await fetchData({ search: value })
  } finally {
    loading.value = false
  }
}
</script>
```

### 页面头部搜索

```vue
<template>
  <GvFlex justify="between" align="center" class="mb-6">
    <h1 class="text-2xl font-bold">资产管理</h1>
    
    <SearchBar
      v-model="searchText"
      class="w-80"
      @search="handleSearch"
    />
  </GvFlex>
</template>
```

### 筛选栏搜索

```vue
<template>
  <GvCard class="mb-6">
    <GvGrid :cols="3" gap="lg">
      <SearchBar
        v-model="filters.search"
        placeholder="搜索名称"
        :show-button="false"
        immediate
        @search="handleFilter"
      />
      
      <GvSelect
        v-model="filters.type"
        :options="typeOptions"
        placeholder="选择类型"
      />
      
      <GvDatePicker
        v-model="filters.date"
        type="daterange"
      />
    </GvGrid>
  </GvCard>
</template>
```

## 最佳实践

1. **立即搜索 vs 按钮搜索**：
   - 本地筛选使用 immediate
   - 远程搜索使用按钮
   - 大数据量避免立即搜索

2. **防抖时间**：
   - 本地筛选: 100-200ms
   - 远程搜索: 300-500ms
   - 慢接口: 500-1000ms

3. **加载状态**：
   - 远程搜索显示 loading
   - 提升用户体验

4. **清空处理**：
   - 清空时恢复全部数据
   - immediate 模式自动触发搜索

5. **占位文本**：
   - 提示搜索的内容类型
   - 简洁明了
