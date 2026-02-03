# GvSelect - 选择器组件

Material Design 3 风格的选择器组件，基于 Element Plus Select 封装，统一样式和 API。

## 基本用法

```vue
<template>
  <GvSelect
    v-model="value"
    :options="options"
    placeholder="请选择"
  />
</template>

<script setup>
import { ref } from 'vue'
import { GvSelect } from '@/components'

const value = ref('')
const options = [
  { label: '选项1', value: '1' },
  { label: '选项2', value: '2' },
  { label: '选项3', value: '3' }
]
</script>
```

## 尺寸

```vue
<GvSelect :options="options" size="small" placeholder="小尺寸" />
<GvSelect :options="options" size="medium" placeholder="中等尺寸" />
<GvSelect :options="options" size="large" placeholder="大尺寸" />
```

## 带标签

```vue
<GvSelect
  v-model="value"
  :options="options"
  label="选择器"
  required
  placeholder="请选择"
/>
```

## 可清空

```vue
<GvSelect
  v-model="value"
  :options="options"
  clearable
  placeholder="可清空"
/>
```

## 禁用

```vue
<GvSelect
  v-model="value"
  :options="options"
  disabled
  placeholder="禁用状态"
/>

<!-- 禁用某个选项 -->
<GvSelect
  v-model="value"
  :options="[
    { label: '选项1', value: '1' },
    { label: '选项2', value: '2', disabled: true },
    { label: '选项3', value: '3' }
  ]"
/>
```

## 多选

```vue
<GvSelect
  v-model="values"
  :options="options"
  multiple
  placeholder="请选择多个选项"
/>

<!-- 限制多选数量 -->
<GvSelect
  v-model="values"
  :options="options"
  multiple
  :multiple-limit="3"
  placeholder="最多选择3个"
/>
```

## 可搜索

```vue
<GvSelect
  v-model="value"
  :options="options"
  filterable
  placeholder="可搜索"
/>
```

## 允许创建新条目

```vue
<GvSelect
  v-model="value"
  :options="options"
  filterable
  allow-create
  placeholder="输入创建新条目"
/>
```

## 远程搜索

```vue
<template>
  <GvSelect
    v-model="value"
    :options="remoteOptions"
    remote
    filterable
    :loading="loading"
    :remote-method="remoteSearch"
    placeholder="输入关键词搜索"
  />
</template>

<script setup>
import { ref } from 'vue'

const value = ref('')
const loading = ref(false)
const remoteOptions = ref([])

const remoteSearch = async (query) => {
  if (query) {
    loading.value = true
    // 模拟远程搜索
    setTimeout(() => {
      remoteOptions.value = [
        { label: `${query} 结果1`, value: '1' },
        { label: `${query} 结果2`, value: '2' }
      ]
      loading.value = false
    }, 500)
  } else {
    remoteOptions.value = []
  }
}
</script>
```

## 验证状态

```vue
<!-- 成功状态 -->
<GvSelect
  v-model="value"
  :options="options"
  status="success"
/>

<!-- 错误状态 -->
<GvSelect
  v-model="value"
  :options="options"
  status="error"
  error-message="请选择一个选项"
/>

<!-- 警告状态 -->
<GvSelect
  v-model="value"
  :options="options"
  status="warning"
/>
```

## 自定义字段名

```vue
<GvSelect
  v-model="value"
  :options="customOptions"
  label-key="name"
  value-key="id"
/>

<script setup>
const customOptions = [
  { name: '选项1', id: 1 },
  { name: '选项2', id: 2 }
]
</script>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| modelValue | `string \| number \| Array` | - | 绑定值 |
| options | `SelectOption[]` | `[]` | 选项数据 |
| size | `'small' \| 'medium' \| 'large'` | `'medium'` | 尺寸 |
| placeholder | `string` | - | 占位文本 |
| disabled | `boolean` | `false` | 是否禁用 |
| clearable | `boolean` | `false` | 是否可清空 |
| multiple | `boolean` | `false` | 是否多选 |
| multipleLimit | `number` | - | 最多可选数量 |
| filterable | `boolean` | `false` | 是否可搜索 |
| allowCreate | `boolean` | `false` | 是否允许创建新条目 |
| filterMethod | `Function` | - | 自定义搜索方法 |
| remote | `boolean` | `false` | 是否远程搜索 |
| remoteMethod | `Function` | - | 远程搜索方法 |
| loading | `boolean` | `false` | 是否正在加载 |
| loadingText | `string` | - | 加载时显示的文字 |
| noDataText | `string` | - | 无数据时显示的文字 |
| noMatchText | `string` | - | 无匹配时显示的文字 |
| popperClass | `string` | - | 下拉框的类名 |
| status | `'success' \| 'error' \| 'warning'` | - | 验证状态 |
| errorMessage | `string` | - | 错误提示信息 |
| required | `boolean` | `false` | 是否必填 |
| label | `string` | - | 标签文本 |
| labelKey | `string` | `'label'` | 选项的标签字段名 |
| valueKey | `string` | `'value'` | 选项的值字段名 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| update:modelValue | `(value)` | 值更新 |
| change | `(value)` | 选中值变化 |
| visible-change | `(visible: boolean)` | 下拉框显示/隐藏 |
| remove-tag | `(value)` | 移除标签 |
| clear | `()` | 清空 |
| focus | `(event: FocusEvent)` | 获得焦点 |
| blur | `(event: FocusEvent)` | 失去焦点 |

## 类型定义

### SelectOption

```typescript
interface SelectOption {
  label: string
  value: string | number
  disabled?: boolean
  [key: string]: any
}
```

## 使用场景

### 表单选择

```vue
<GvSelect
  v-model="form.type"
  :options="typeOptions"
  label="资产类型"
  required
  placeholder="请选择资产类型"
/>

<GvSelect
  v-model="form.category"
  :options="categoryOptions"
  label="分类"
  clearable
  placeholder="请选择分类"
/>
```

### 多选标签

```vue
<GvSelect
  v-model="form.tags"
  :options="tagOptions"
  label="标签"
  multiple
  filterable
  allow-create
  placeholder="选择或创建标签"
/>
```

### 级联选择

```vue
<div class="space-y-4">
  <GvSelect
    v-model="province"
    :options="provinceOptions"
    label="省份"
    @change="handleProvinceChange"
  />
  
  <GvSelect
    v-model="city"
    :options="cityOptions"
    label="城市"
    :disabled="!province"
  />
</div>
```

### 远程搜索用户

```vue
<GvSelect
  v-model="userId"
  :options="userOptions"
  label="选择用户"
  remote
  filterable
  :loading="loading"
  :remote-method="searchUsers"
  placeholder="输入用户名搜索"
/>
```

## 最佳实践

1. **选项数据**：
   - 数据量小（< 100）：直接使用本地数据
   - 数据量大：使用 remote + filterable

2. **搜索功能**：
   - 选项较多时建议启用 filterable
   - 远程数据使用 remote + remoteMethod

3. **多选限制**：
   - 使用 multipleLimit 限制选择数量
   - 避免过多选项影响性能

4. **验证状态**：
   - 使用 status 提供即时反馈
   - 配合 errorMessage 显示错误

5. **标签使用**：
   - 重要选择器添加 label
   - 必填字段使用 required

6. **禁用选项**：
   - 不可用选项设置 disabled: true
   - 提供清晰的禁用原因

7. **占位文本**：
   - 使用清晰的 placeholder 引导用户
   - 避免过长的占位文本
