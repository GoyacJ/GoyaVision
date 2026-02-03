# AssetCard - 资产卡片组件

业务组件，用于展示资产信息的卡片，基于 `GvCard` 封装。

## 基本用法

```vue
<template>
  <AssetCard
    :asset="asset"
    @click="handleClick"
  />
</template>

<script setup>
import { AssetCard } from '@/components'

const asset = {
  id: 1,
  name: '视频资产 1',
  type: 'MP4',
  status: 'running',
  thumbnail: 'https://example.com/thumb.jpg',
  description: '这是一个视频资产的描述信息',
  size: '128 MB',
  duration: '05:30',
  createdAt: '2024-01-01 12:00'
}

const handleClick = (asset) => {
  console.log('点击资产:', asset)
}
</script>
```

## 网格布局

```vue
<GvGrid :cols="3" gap="lg">
  <AssetCard
    v-for="asset in assets"
    :key="asset.id"
    :asset="asset"
    @click="handleClick"
    @edit="handleEdit"
    @delete="handleDelete"
  />
</GvGrid>
```

## 可选择

```vue
<template>
  <GvGrid :cols="3" gap="lg">
    <AssetCard
      v-for="asset in assets"
      :key="asset.id"
      :asset="asset"
      selectable
      :selected="selectedIds.includes(asset.id)"
      @select="handleSelect"
    />
  </GvGrid>
</template>

<script setup>
const selectedIds = ref([])

const handleSelect = (selected, asset) => {
  if (selected) {
    selectedIds.value.push(asset.id)
  } else {
    selectedIds.value = selectedIds.value.filter(id => id !== asset.id)
  }
}
</script>
```

## 自定义操作按钮

```vue
<AssetCard
  :asset="asset"
  :show-edit="false"
  :show-delete="false"
  @detail="handleDetail"
/>
```

## 无缩略图

```vue
<AssetCard
  :asset="{
    id: 1,
    name: '资产名称',
    type: 'MP4',
    status: 'stopped'
  }"
/>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| asset | `Asset` | - | 资产数据（必填） |
| selectable | `boolean` | `false` | 是否可选择 |
| selected | `boolean` | `false` | 是否选中 |
| showActions | `boolean` | `true` | 是否显示操作按钮 |
| showDetail | `boolean` | `true` | 是否显示详情按钮 |
| showEdit | `boolean` | `true` | 是否显示编辑按钮 |
| showDelete | `boolean` | `true` | 是否显示删除按钮 |

## Asset

| 属性 | 类型 | 说明 |
|------|------|------|
| id | `string \| number` | ID（必填） |
| name | `string` | 名称（必填） |
| type | `string` | 类型（必填） |
| status | `string` | 状态（必填） |
| thumbnail | `string` | 缩略图 URL |
| description | `string` | 描述 |
| size | `string` | 大小 |
| duration | `string` | 时长 |
| createdAt | `string` | 创建时间 |
| updatedAt | `string` | 更新时间 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| click | `(asset)` | 点击卡片 |
| select | `(selected, asset)` | 选择状态变化 |
| detail | `(asset)` | 点击详情 |
| edit | `(asset)` | 点击编辑 |
| delete | `(asset)` | 点击删除 |

## 使用场景

### 资产列表

```vue
<template>
  <GvContainer>
    <PageHeader title="资产管理">
      <template #actions>
        <GvSpace>
          <SearchBar v-model="searchText" class="w-80" />
          <GvButton icon="Plus" @click="handleAdd">新建</GvButton>
        </GvSpace>
      </template>
    </PageHeader>
    
    <FilterBar
      v-model="filters"
      :fields="filterFields"
      @filter="handleFilter"
    />
    
    <GvGrid :cols="3" gap="lg">
      <AssetCard
        v-for="asset in assets"
        :key="asset.id"
        :asset="asset"
        @click="handleView"
        @edit="handleEdit"
        @delete="handleDelete"
      />
    </GvGrid>
  </GvContainer>
</template>
```

### 批量选择

```vue
<template>
  <div>
    <GvFlex justify="between" class="mb-4">
      <span>已选 {{ selectedAssets.length }} 项</span>
      <GvSpace>
        <GvButton
          :disabled="selectedAssets.length === 0"
          @click="handleBatchDelete"
        >
          批量删除
        </GvButton>
        <GvButton
          :disabled="selectedAssets.length === 0"
          @click="handleBatchExport"
        >
          批量导出
        </GvButton>
      </GvSpace>
    </GvFlex>
    
    <GvGrid :cols="3" gap="lg">
      <AssetCard
        v-for="asset in assets"
        :key="asset.id"
        :asset="asset"
        selectable
        :selected="isSelected(asset.id)"
        @select="handleSelect"
      />
    </GvGrid>
  </div>
</template>
```

## 最佳实践

1. **缩略图**：
   - 提供默认占位图
   - 使用懒加载
   - 统一宽高比

2. **状态显示**：
   - 使用 StatusBadge 统一状态
   - 状态显示在右上角

3. **操作按钮**：
   - 悬停时显示
   - 最多 3 个按钮
   - 危险操作放最后

4. **选择功能**：
   - 批量操作时启用
   - 显示选中数量
   - 提供全选功能

5. **响应式**：
   - 使用 GvGrid 响应式布局
   - 移动端 1 列
   - 平板 2 列
   - 桌面 3-4 列
