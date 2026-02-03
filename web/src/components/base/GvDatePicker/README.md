# GvDatePicker - 日期选择器组件

Material Design 3 风格的日期选择器组件，基于 Element Plus DatePicker 封装，统一样式和 API。

## 基本用法

```vue
<template>
  <GvDatePicker
    v-model="date"
    placeholder="请选择日期"
  />
</template>

<script setup>
import { ref } from 'vue'
import { GvDatePicker } from '@/components'

const date = ref('')
</script>
```

## 日期选择器类型

```vue
<!-- 日期 -->
<GvDatePicker v-model="date" type="date" placeholder="选择日期" />

<!-- 日期时间 -->
<GvDatePicker v-model="datetime" type="datetime" placeholder="选择日期时间" />

<!-- 月份 -->
<GvDatePicker v-model="month" type="month" placeholder="选择月份" />

<!-- 年份 -->
<GvDatePicker v-model="year" type="year" placeholder="选择年份" />

<!-- 日期范围 -->
<GvDatePicker
  v-model="daterange"
  type="daterange"
  start-placeholder="开始日期"
  end-placeholder="结束日期"
/>

<!-- 日期时间范围 -->
<GvDatePicker
  v-model="datetimerange"
  type="datetimerange"
  start-placeholder="开始日期时间"
  end-placeholder="结束日期时间"
/>
```

## 尺寸

```vue
<GvDatePicker v-model="date" size="small" placeholder="小尺寸" />
<GvDatePicker v-model="date" size="default" placeholder="默认尺寸" />
<GvDatePicker v-model="date" size="large" placeholder="大尺寸" />
```

## 带标签

```vue
<GvDatePicker
  v-model="date"
  label="选择日期"
  required
  placeholder="请选择日期"
/>
```

## 日期格式

```vue
<!-- 显示格式 -->
<GvDatePicker
  v-model="date"
  format="YYYY/MM/DD"
  placeholder="YYYY/MM/DD"
/>

<!-- 值格式 -->
<GvDatePicker
  v-model="date"
  value-format="YYYY-MM-DD"
  placeholder="值格式为 YYYY-MM-DD"
/>
```

## 禁用日期

```vue
<template>
  <GvDatePicker
    v-model="date"
    :disabled-date="disabledDate"
    placeholder="禁用过去的日期"
  />
</template>

<script setup>
const disabledDate = (date: Date) => {
  return date < new Date()
}
</script>
```

## 快捷选项

```vue
<template>
  <GvDatePicker
    v-model="date"
    :shortcuts="shortcuts"
    placeholder="带快捷选项"
  />
</template>

<script setup>
const shortcuts = [
  {
    text: '今天',
    value: new Date()
  },
  {
    text: '昨天',
    value: () => {
      const date = new Date()
      date.setDate(date.getDate() - 1)
      return date
    }
  },
  {
    text: '一周前',
    value: () => {
      const date = new Date()
      date.setDate(date.getDate() - 7)
      return date
    }
  }
]
</script>
```

## 范围快捷选项

```vue
<template>
  <GvDatePicker
    v-model="daterange"
    type="daterange"
    :shortcuts="rangeShortcuts"
    start-placeholder="开始日期"
    end-placeholder="结束日期"
  />
</template>

<script setup>
const rangeShortcuts = [
  {
    text: '最近一周',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setDate(start.getDate() - 7)
      return [start, end]
    }
  },
  {
    text: '最近一个月',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setMonth(start.getMonth() - 1)
      return [start, end]
    }
  },
  {
    text: '最近三个月',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setMonth(start.getMonth() - 3)
      return [start, end]
    }
  }
]
</script>
```

## 验证状态

```vue
<!-- 成功状态 -->
<GvDatePicker
  v-model="date"
  status="success"
/>

<!-- 错误状态 -->
<GvDatePicker
  v-model="date"
  status="error"
  error-message="请选择一个有效日期"
/>

<!-- 警告状态 -->
<GvDatePicker
  v-model="date"
  status="warning"
/>
```

## 禁用和清空

```vue
<!-- 禁用 -->
<GvDatePicker v-model="date" disabled />

<!-- 不可清空 -->
<GvDatePicker v-model="date" :clearable="false" />
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| modelValue | `Date \| string \| number \| [Date, Date]` | - | 绑定值 |
| type | `'date' \| 'datetime' \| 'daterange' \| 'datetimerange' \| 'month' \| 'year'` | `'date'` | 选择器类型 |
| size | `'small' \| 'default' \| 'large'` | `'default'` | 尺寸 |
| placeholder | `string` | - | 占位文本 |
| startPlaceholder | `string` | - | 范围选择起始占位 |
| endPlaceholder | `string` | - | 范围选择结束占位 |
| disabled | `boolean` | `false` | 是否禁用 |
| clearable | `boolean` | `true` | 是否可清空 |
| format | `string` | - | 日期显示格式 |
| valueFormat | `string` | - | 绑定值格式 |
| disabledDate | `Function` | - | 禁用日期函数 |
| shortcuts | `Array` | - | 快捷选项 |
| status | `'success' \| 'error' \| 'warning'` | - | 验证状态 |
| errorMessage | `string` | - | 错误提示信息 |
| required | `boolean` | `false` | 是否必填 |
| label | `string` | - | 标签文本 |
| rangeSeparator | `string` | `'-'` | 范围分隔符 |
| defaultTime | `Date \| [Date, Date]` | - | 默认时间 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| update:modelValue | `(value)` | 值更新 |
| change | `(value)` | 值变化 |
| clear | `()` | 清空 |
| focus | `(event: FocusEvent)` | 获得焦点 |
| blur | `(event: FocusEvent)` | 失去焦点 |

## 使用场景

### 表单日期选择

```vue
<el-form :model="form">
  <el-form-item>
    <GvDatePicker
      v-model="form.startDate"
      label="开始日期"
      required
      placeholder="请选择开始日期"
    />
  </el-form-item>
  
  <el-form-item>
    <GvDatePicker
      v-model="form.endDate"
      label="结束日期"
      required
      :disabled-date="disabledEndDate"
      placeholder="请选择结束日期"
    />
  </el-form-item>
</el-form>
```

### 日期范围筛选

```vue
<GvFlex gap="md">
  <GvDatePicker
    v-model="filters.dateRange"
    type="daterange"
    :shortcuts="shortcuts"
    start-placeholder="开始日期"
    end-placeholder="结束日期"
  />
  
  <GvButton @click="handleFilter">筛选</GvButton>
  <GvButton variant="text" @click="handleReset">重置</GvButton>
</GvFlex>
```

### 生日选择

```vue
<GvDatePicker
  v-model="birthday"
  type="date"
  label="出生日期"
  :disabled-date="disabledFutureDate"
  placeholder="请选择出生日期"
/>

<script setup>
const disabledFutureDate = (date: Date) => {
  return date > new Date()
}
</script>
```

### 预约时间

```vue
<GvDatePicker
  v-model="appointmentTime"
  type="datetime"
  label="预约时间"
  :disabled-date="disabledPastDate"
  :shortcuts="appointmentShortcuts"
  placeholder="请选择预约时间"
/>

<script setup>
const disabledPastDate = (date: Date) => {
  return date < new Date()
}

const appointmentShortcuts = [
  { text: '明天上午 9:00', value: getTomorrow(9, 0) },
  { text: '明天下午 14:00', value: getTomorrow(14, 0) },
  { text: '后天上午 9:00', value: getDayAfterTomorrow(9, 0) }
]
</script>
```

## 最佳实践

1. **类型选择**：
   - date: 日常日期选择
   - datetime: 需要具体时间
   - daterange: 筛选、统计
   - month/year: 报表、归档

2. **快捷选项**：
   - 常用场景提供快捷选项
   - 减少用户操作步骤
   - 提升用户体验

3. **禁用日期**：
   - 生日禁用未来日期
   - 预约禁用过去日期
   - 结束日期不能早于开始日期

4. **日期格式**：
   - 显示格式友好易读
   - 值格式标准统一
   - 根据业务需求选择

5. **验证**：
   - 使用 status 提供即时反馈
   - 配合 errorMessage 显示错误
   - required 标识必填字段

6. **范围选择**：
   - 提供清晰的占位文本
   - 使用快捷选项
   - 合理设置分隔符
