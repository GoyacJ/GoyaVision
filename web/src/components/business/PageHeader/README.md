# PageHeader - 页面头部组件

业务组件，统一的页面头部，包含标题、描述、面包屑导航和操作区域。

## 基本用法

```vue
<template>
  <PageHeader title="资产管理" />
</template>

<script setup>
import { PageHeader } from '@/components'
</script>
```

## 带描述

```vue
<PageHeader
  title="资产管理"
  description="管理所有视频资产，包括上传、编辑、删除等操作"
/>
```

## 带操作按钮

```vue
<template>
  <PageHeader title="资产管理">
    <template #actions>
      <GvSpace>
        <GvButton icon="Plus" @click="handleAdd">新建</GvButton>
        <GvButton variant="tonal" @click="handleExport">导出</GvButton>
      </GvSpace>
    </template>
  </PageHeader>
</template>
```

## 带面包屑导航

```vue
<PageHeader
  title="编辑资产"
  :breadcrumb="[
    { label: '首页', to: '/' },
    { label: '资产管理', to: '/assets' },
    { label: '编辑资产' }
  ]"
/>
```

## 带返回按钮

```vue
<template>
  <PageHeader
    title="资产详情"
    show-back
    @back="handleBack"
  />
</template>

<script setup>
const handleBack = () => {
  router.back()
}
</script>
```

## 带额外内容

```vue
<template>
  <PageHeader title="资产管理">
    <template #actions>
      <GvButton icon="Plus" @click="handleAdd">新建</GvButton>
    </template>
    
    <template #extra>
      <GvSpace wrap>
        <StatusBadge status="running" text="3 个任务运行中" />
        <StatusBadge status="pending" text="5 个任务等待中" />
      </GvSpace>
    </template>
  </PageHeader>
</template>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| title | `string` | - | 页面标题（必填） |
| description | `string` | - | 页面描述 |
| breadcrumb | `BreadcrumbItem[]` | - | 面包屑导航 |
| showBack | `boolean` | `false` | 是否显示返回按钮 |
| backText | `string` | `'返回'` | 返回按钮文本 |

## BreadcrumbItem

| 属性 | 类型 | 说明 |
|------|------|------|
| label | `string` | 文本 |
| to | `string` | 跳转路径 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| back | `()` | 返回按钮点击 |

## Slots

| 插槽名 | 说明 |
|--------|------|
| actions | 操作区域（右侧） |
| extra | 额外内容区域（底部） |

## 使用场景

### 列表页

```vue
<template>
  <GvContainer>
    <PageHeader
      title="资产管理"
      description="管理所有视频资产"
    >
      <template #actions>
        <GvSpace>
          <SearchBar v-model="searchText" class="w-80" />
          <GvButton icon="Plus" @click="handleAdd">新建</GvButton>
        </GvSpace>
      </template>
    </PageHeader>
    
    <GvTable :data="tableData" :columns="columns" />
  </GvContainer>
</template>
```

### 详情页

```vue
<template>
  <GvContainer>
    <PageHeader
      title="资产详情"
      :breadcrumb="[
        { label: '首页', to: '/' },
        { label: '资产管理', to: '/assets' },
        { label: asset.name }
      ]"
      show-back
      @back="router.back()"
    >
      <template #actions>
        <GvSpace>
          <GvButton @click="handleEdit">编辑</GvButton>
          <GvButton variant="text" @click="handleDelete">删除</GvButton>
        </GvSpace>
      </template>
      
      <template #extra>
        <GvSpace>
          <StatusBadge :status="asset.status" />
          <span class="text-text-tertiary">
            创建于 {{ asset.createdAt }}
          </span>
        </GvSpace>
      </template>
    </PageHeader>
    
    <GvCard>
      <!-- 详情内容 -->
    </GvCard>
  </GvContainer>
</template>
```

### 表单页

```vue
<template>
  <GvContainer>
    <PageHeader
      title="新建资产"
      show-back
      @back="router.back()"
    >
      <template #actions>
        <GvSpace>
          <GvButton variant="tonal" @click="handleCancel">取消</GvButton>
          <GvButton @click="handleSave">保存</GvButton>
        </GvSpace>
      </template>
    </PageHeader>
    
    <GvCard>
      <el-form :model="form">
        <!-- 表单内容 -->
      </el-form>
    </GvCard>
  </GvContainer>
</template>
```

## 最佳实践

1. **标题**：
   - 清晰描述页面功能
   - 不超过 20 个字符

2. **描述**：
   - 可选，补充说明
   - 不超过 50 个字符

3. **面包屑**：
   - 列表页通常不需要
   - 详情页、编辑页需要
   - 最后一项不设置 to

4. **返回按钮**：
   - 详情页、编辑页使用
   - 列表页不需要

5. **操作按钮**：
   - 主要操作放在右侧
   - 使用 GvSpace 保持间距
   - 不超过 3 个按钮

6. **额外内容**：
   - 展示页面状态信息
   - 展示统计数据
