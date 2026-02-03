# GvInput - 输入框组件

Material Design 3 风格的输入框组件，支持多种类型和验证状态。

## 基本用法

```vue
<template>
  <GvInput v-model="value" placeholder="请输入内容" />
</template>

<script setup>
import { ref } from 'vue'
import { GvInput } from '@/components'

const value = ref('')
</script>
```

## 输入框类型

```vue
<GvInput type="text" placeholder="文本输入" />
<GvInput type="password" placeholder="密码输入" />
<GvInput type="number" placeholder="数字输入" />
<GvInput type="email" placeholder="邮箱输入" />
<GvInput type="tel" placeholder="电话输入" />
<GvInput type="url" placeholder="网址输入" />
<GvInput type="search" placeholder="搜索输入" />
```

## 尺寸

```vue
<GvInput size="small" placeholder="小尺寸" />
<GvInput size="medium" placeholder="中等尺寸" />
<GvInput size="large" placeholder="大尺寸" />
```

## 带标签

```vue
<GvInput label="用户名" placeholder="请输入用户名" />
<GvInput label="密码" required placeholder="请输入密码" />
```

## 可清除

```vue
<GvInput v-model="value" clearable placeholder="可清除的输入框" />
```

## 密码框

```vue
<GvInput
  v-model="password"
  type="password"
  show-password
  placeholder="请输入密码"
/>
```

## 图标

### 前置图标

```vue
<GvInput prefix-icon="User" placeholder="请输入用户名" />
<GvInput prefix-icon="Lock" placeholder="请输入密码" />
```

### 后置图标

```vue
<GvInput suffix-icon="Calendar" placeholder="选择日期" />
<GvInput suffix-icon="Location" placeholder="选择地点" />
```

### 自定义插槽

```vue
<GvInput placeholder="搜索">
  <template #prefix>
    <el-icon><Search /></el-icon>
  </template>
  <template #suffix>
    <el-button text>搜索</el-button>
  </template>
</GvInput>
```

## 字数限制

```vue
<GvInput
  v-model="value"
  :maxlength="100"
  show-count
  placeholder="最多输入 100 字"
/>
```

## 验证状态

```vue
<!-- 成功状态 -->
<GvInput
  v-model="value"
  status="success"
  placeholder="验证成功"
/>

<!-- 错误状态 -->
<GvInput
  v-model="value"
  status="error"
  error-message="用户名不能为空"
  placeholder="验证失败"
/>

<!-- 警告状态 -->
<GvInput
  v-model="value"
  status="warning"
  placeholder="警告提示"
/>
```

## 禁用和只读

```vue
<GvInput disabled placeholder="禁用状态" />
<GvInput readonly value="只读状态" />
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| modelValue | `string \| number` | - | 绑定值 |
| type | `'text' \| 'password' \| 'number' \| 'email' \| 'tel' \| 'url' \| 'search'` | `'text'` | 输入框类型 |
| size | `'small' \| 'medium' \| 'large'` | `'medium'` | 尺寸 |
| placeholder | `string` | - | 占位文本 |
| disabled | `boolean` | `false` | 是否禁用 |
| readonly | `boolean` | `false` | 是否只读 |
| clearable | `boolean` | `false` | 是否显示清除按钮 |
| showPassword | `boolean` | `false` | 是否显示密码切换 |
| prefixIcon | `string \| Component` | - | 前置图标 |
| suffixIcon | `string \| Component` | - | 后置图标 |
| maxlength | `number` | - | 最大输入长度 |
| showCount | `boolean` | `false` | 是否显示字数统计 |
| status | `'success' \| 'error' \| 'warning'` | - | 验证状态 |
| errorMessage | `string` | - | 错误提示信息 |
| autofocus | `boolean` | `false` | 是否自动聚焦 |
| autocomplete | `string` | - | 原生 autocomplete |
| name | `string` | - | 原生 name |
| form | `string` | - | 原生 form |
| required | `boolean` | `false` | 是否必填 |
| label | `string` | - | 标签文本 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| update:modelValue | `(value: string \| number)` | 值更新 |
| input | `(value: string \| number)` | 输入事件 |
| change | `(value: string \| number)` | 变化事件 |
| focus | `(event: FocusEvent)` | 获得焦点 |
| blur | `(event: FocusEvent)` | 失去焦点 |
| clear | `()` | 清除事件 |
| keydown | `(event: KeyboardEvent)` | 按键事件 |
| keyup | `(event: KeyboardEvent)` | 按键抬起 |

## Slots

| 插槽名 | 说明 |
|--------|------|
| prefix | 前置内容 |
| suffix | 后置内容 |

## Exposes

| 方法名 | 说明 |
|--------|------|
| focus() | 使输入框获得焦点 |
| blur() | 使输入框失去焦点 |
| select() | 选中输入框文本 |

## 使用场景

### 表单输入

```vue
<GvInput
  v-model="form.username"
  label="用户名"
  required
  placeholder="请输入用户名"
  prefix-icon="User"
  clearable
/>

<GvInput
  v-model="form.password"
  type="password"
  label="密码"
  required
  show-password
  placeholder="请输入密码"
  prefix-icon="Lock"
/>

<GvInput
  v-model="form.email"
  type="email"
  label="邮箱"
  placeholder="请输入邮箱"
  prefix-icon="Message"
/>
```

### 搜索框

```vue
<GvInput
  v-model="keyword"
  type="search"
  size="large"
  placeholder="搜索资产..."
  prefix-icon="Search"
  clearable
  @keyup.enter="handleSearch"
/>
```

### 带验证的输入框

```vue
<GvInput
  v-model="form.name"
  label="资产名称"
  required
  :status="nameError ? 'error' : undefined"
  :error-message="nameError"
  placeholder="请输入资产名称"
  @blur="validateName"
/>
```

## 最佳实践

1. **标签使用**：重要输入框建议添加 label
2. **必填标识**：必填字段使用 required 属性显示红色星号
3. **验证反馈**：使用 status 和 errorMessage 提供及时反馈
4. **清除按钮**：长文本输入建议添加 clearable
5. **密码输入**：密码框建议使用 showPassword
6. **字数限制**：有长度限制时使用 maxlength + showCount
7. **图标使用**：使用语义化图标提升识别度
8. **自动聚焦**：关键输入框可使用 autofocus
