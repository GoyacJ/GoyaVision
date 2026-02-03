# GvTable - 表格组件

Material Design 3 风格的表格组件，基于 Element Plus Table 封装，统一样式和 API。

## 基本用法

```vue
<template>
  <GvTable :data="tableData" :columns="columns" />
</template>

<script setup>
import { ref } from 'vue'
import { GvTable } from '@/components'

const tableData = ref([
  { id: 1, name: '张三', age: 28, address: '北京市' },
  { id: 2, name: '李四', age: 32, address: '上海市' }
])

const columns = [
  { prop: 'id', label: 'ID', width: '80' },
  { prop: 'name', label: '姓名', width: '120' },
  { prop: 'age', label: '年龄', width: '80' },
  { prop: 'address', label: '地址' }
]
</script>
```

## 带边框和斑马纹

```vue
<GvTable :data="tableData" :columns="columns" border stripe />
```

## 固定列

```vue
<script setup>
const columns = [
  { prop: 'id', label: 'ID', width: '80', fixed: 'left' },
  { prop: 'name', label: '姓名', width: '120', fixed: 'left' },
  { prop: 'age', label: '年龄', width: '80' },
  { prop: 'address', label: '地址', width: '300' },
  { prop: 'actions', label: '操作', width: '150', fixed: 'right' }
]
</script>
```

## 可排序

```vue
<script setup>
const columns = [
  { prop: 'id', label: 'ID', sortable: true },
  { prop: 'name', label: '姓名' },
  { prop: 'age', label: '年龄', sortable: true },
  { prop: 'address', label: '地址' }
]
</script>
```

## 自定义列渲染

```vue
<template>
  <GvTable :data="tableData" :columns="columns">
    <!-- 自定义状态列 -->
    <template #status="{ row }">
      <GvBadge
        :variant="row.status === 'active' ? 'success' : 'error'"
        :text="row.status === 'active' ? '激活' : '禁用'"
      />
    </template>
    
    <!-- 操作列 -->
    <template #actions="{ row }">
      <GvSpace size="xs">
        <GvButton size="small" @click="handleEdit(row)">编辑</GvButton>
        <GvButton size="small" variant="text" @click="handleDelete(row)">
          删除
        </GvButton>
      </GvSpace>
    </template>
  </GvTable>
</template>

<script setup>
const columns = [
  { prop: 'name', label: '姓名' },
  { prop: 'status', label: '状态' },
  { prop: 'actions', label: '操作', width: '200' }
]
</script>
```

## 分页

```vue
<template>
  <GvTable
    :data="tableData"
    :columns="columns"
    pagination
    :pagination-config="paginationConfig"
    @current-change="handlePageChange"
    @size-change="handleSizeChange"
  />
</template>

<script setup>
const paginationConfig = ref({
  currentPage: 1,
  pageSize: 10,
  total: 100,
  pageSizes: [10, 20, 50, 100]
})

const handlePageChange = (page: number) => {
  paginationConfig.value.currentPage = page
  fetchData()
}

const handleSizeChange = (size: number) => {
  paginationConfig.value.pageSize = size
  fetchData()
}
</script>
```

## 可选择行

```vue
<template>
  <GvTable
    :data="tableData"
    :columns="columns"
    selectable
    @selection-change="handleSelectionChange"
  />
</template>

<script setup>
const handleSelectionChange = (selection: any[]) => {
  console.log('选中的行:', selection)
}
</script>
```

## 加载状态

```vue
<GvTable :data="tableData" :columns="columns" :loading="loading" />
```

## 高度和最大高度

```vue
<!-- 固定高度 -->
<GvTable :data="tableData" :columns="columns" height="400" />

<!-- 最大高度 -->
<GvTable :data="tableData" :columns="columns" max-height="600" />
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| data | `any[]` | `[]` | 表格数据（必填） |
| columns | `TableColumn[]` | `[]` | 列配置（必填） |
| size | `'small' \| 'default' \| 'large'` | `'default'` | 表格尺寸 |
| border | `boolean` | `false` | 是否显示边框 |
| stripe | `boolean` | `true` | 是否显示斑马纹 |
| showHeader | `boolean` | `true` | 是否显示表头 |
| highlightCurrentRow | `boolean` | `false` | 是否高亮当前行 |
| loading | `boolean` | `false` | 是否显示加载 |
| emptyText | `string` | `'暂无数据'` | 空数据文本 |
| rowKey | `string` | `'id'` | 行的 key |
| pagination | `boolean` | `false` | 是否启用分页 |
| paginationConfig | `PaginationConfig` | - | 分页配置 |
| height | `string \| number` | - | 表格高度 |
| maxHeight | `string \| number` | - | 表格最大高度 |
| selectable | `boolean` | `false` | 是否可选择行 |
| defaultSelection | `any[]` | - | 默认选中的行 |

## TableColumn

| 属性 | 类型 | 说明 |
|------|------|------|
| prop | `string` | 列的 key（必填） |
| label | `string` | 列标题（必填） |
| width | `string \| number` | 列宽度 |
| minWidth | `string \| number` | 最小宽度 |
| fixed | `boolean \| 'left' \| 'right'` | 是否固定列 |
| sortable | `boolean \| 'custom'` | 是否可排序 |
| align | `'left' \| 'center' \| 'right'` | 对齐方式 |
| formatter | `Function` | 自定义渲染函数 |
| showOverflowTooltip | `boolean` | 是否显示溢出提示 |

## PaginationConfig

| 属性 | 类型 | 说明 |
|------|------|------|
| currentPage | `number` | 当前页码 |
| pageSize | `number` | 每页条数 |
| total | `number` | 总条数 |
| pageSizes | `number[]` | 每页条数选项 |
| layout | `string` | 分页布局 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| selection-change | `(selection: any[])` | 选择项变化 |
| row-click | `(row, column, event)` | 行点击 |
| row-dblclick | `(row, column, event)` | 行双击 |
| cell-click | `(row, column, cell, event)` | 单元格点击 |
| sort-change | `(data)` | 排序变化 |
| current-change | `(page: number)` | 页码变化 |
| size-change | `(size: number)` | 每页条数变化 |

## Slots

| 插槽名 | 参数 | 说明 |
|--------|------|------|
| [prop] | `{ row, column, $index }` | 自定义列内容 |
| actions | - | 操作列 |

## Methods

| 方法名 | 参数 | 说明 |
|--------|------|------|
| clearSelection | `()` | 清空选择 |
| toggleRowSelection | `(row, selected?)` | 切换行选择 |
| setCurrentRow | `(row)` | 设置当前行 |

## 使用场景

### 数据列表

```vue
<GvTable
  :data="assetList"
  :columns="columns"
  :loading="loading"
  pagination
  :pagination-config="paginationConfig"
>
  <template #type="{ row }">
    <GvTag :color="getTypeColor(row.type)">
      {{ row.type }}
    </GvTag>
  </template>
  
  <template #status="{ row }">
    <GvBadge
      :variant="row.status === 'active' ? 'success' : 'error'"
      :text="row.status"
    />
  </template>
  
  <template #actions="{ row }">
    <GvSpace size="xs">
      <GvButton size="small" @click="handleView(row)">查看</GvButton>
      <GvButton size="small" @click="handleEdit(row)">编辑</GvButton>
      <GvButton size="small" variant="text" @click="handleDelete(row)">
        删除
      </GvButton>
    </GvSpace>
  </template>
</GvTable>
```

### 批量操作

```vue
<template>
  <div>
    <GvFlex justify="between" class="mb-4">
      <GvSpace>
        <GvButton
          :disabled="selectedRows.length === 0"
          @click="handleBatchDelete"
        >
          批量删除
        </GvButton>
        <GvButton
          :disabled="selectedRows.length === 0"
          @click="handleBatchExport"
        >
          批量导出
        </GvButton>
      </GvSpace>
      
      <span class="text-text-secondary">
        已选 {{ selectedRows.length }} 项
      </span>
    </GvFlex>
    
    <GvTable
      :data="tableData"
      :columns="columns"
      selectable
      @selection-change="handleSelectionChange"
    />
  </div>
</template>

<script setup>
const selectedRows = ref([])

const handleSelectionChange = (selection: any[]) => {
  selectedRows.value = selection
}
</script>
```

### 固定表头

```vue
<GvTable
  :data="tableData"
  :columns="columns"
  max-height="500"
/>
```

## 最佳实践

1. **列配置**：
   - 设置合理的列宽
   - 长文本使用 showOverflowTooltip
   - 操作列固定在右侧

2. **分页**：
   - 大数据量必须使用分页
   - 提供合理的 pageSizes
   - 保存用户的分页偏好

3. **性能优化**：
   - 使用 rowKey 提升性能
   - 避免在表格中进行复杂计算
   - 自定义列使用插槽而非 formatter

4. **用户体验**：
   - 加载时显示 loading
   - 空数据提供友好提示
   - 重要操作添加确认

5. **响应式**：
   - 移动端考虑使用卡片列表
   - 固定重要列
   - 隐藏次要列
