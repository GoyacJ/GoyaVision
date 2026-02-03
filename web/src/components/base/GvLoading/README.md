# GvLoading - 加载组件

Material Design 3 风格的加载组件，支持多种加载样式和场景。

## 基本用法

```vue
<template>
  <GvLoading :loading="true" />
</template>

<script setup>
import { GvLoading } from '@/components'
</script>
```

## 包裹内容

```vue
<GvLoading :loading="isLoading">
  <div class="content">
    <p>这里是内容区域</p>
    <p>加载时会显示遮罩</p>
  </div>
</GvLoading>
```

## 加载类型

```vue
<GvLoading type="circular" />  <!-- 圆形加载器（默认） -->
<GvLoading type="spinner" />   <!-- 旋转加载器 -->
<GvLoading type="dots" />      <!-- 点状加载器 -->
<GvLoading type="bars" />      <!-- 条形加载器 -->
```

## 尺寸

```vue
<GvLoading size="small" />
<GvLoading size="medium" />
<GvLoading size="large" />
```

## 带文本

```vue
<GvLoading text="加载中..." />
<GvLoading text="正在处理，请稍候..." />
```

## 颜色

```vue
<GvLoading color="primary" />
<GvLoading color="secondary" />
<GvLoading color="success" />
<GvLoading color="white" />  <!-- 深色背景上使用 -->
```

## 全屏加载

```vue
<GvLoading
  :loading="isLoading"
  fullscreen
  text="加载中..."
/>
```

## 自定义背景

```vue
<GvLoading
  fullscreen
  background="rgba(0, 0, 0, 0.8)"
  color="white"
  text="加载中..."
/>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| loading | `boolean` | `true` | 是否显示加载 |
| type | `'spinner' \| 'circular' \| 'dots' \| 'bars'` | `'circular'` | 加载器类型 |
| size | `'small' \| 'medium' \| 'large'` | `'medium'` | 尺寸 |
| text | `string` | - | 加载文本 |
| fullscreen | `boolean` | `false` | 是否全屏 |
| lock | `boolean` | `true` | 是否锁定滚动 |
| background | `string` | - | 背景色 |
| customClass | `string` | - | 自定义类名 |
| color | `'primary' \| ...` | `'primary'` | 颜色 |

## Slots

| 插槽名 | 说明 |
|--------|------|
| default | 被加载内容 |

## 使用场景

### 局部加载（卡片）

```vue
<GvCard>
  <GvLoading :loading="isLoading">
    <div class="content">
      <p>卡片内容</p>
      <p>数据加载中...</p>
    </div>
  </GvLoading>
</GvCard>
```

### 全屏加载（页面初始化）

```vue
<template>
  <div>
    <GvLoading
      :loading="isInitializing"
      fullscreen
      text="正在初始化..."
    />
    
    <PageContent />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const isInitializing = ref(true)

onMounted(async () => {
  await fetchData()
  isInitializing.value = false
})
</script>
```

### 按钮加载（已由 GvButton 内置）

```vue
<GvButton :loading="isSubmitting" @click="handleSubmit">
  提交
</GvButton>
```

### 表格加载

```vue
<GvCard>
  <GvLoading :loading="isTableLoading" type="bars">
    <el-table :data="tableData">
      <!-- 表格列 -->
    </el-table>
  </GvLoading>
</GvCard>
```

### 异步操作

```vue
<template>
  <div>
    <GvButton @click="handleRefresh">刷新数据</GvButton>
    
    <GvLoading :loading="isRefreshing" text="刷新中...">
      <div class="data-list">
        <DataItem v-for="item in data" :key="item.id" :data="item" />
      </div>
    </GvLoading>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const isRefreshing = ref(false)
const data = ref([])

const handleRefresh = async () => {
  isRefreshing.value = true
  try {
    data.value = await fetchData()
  } finally {
    isRefreshing.value = false
  }
}
</script>
```

### 分步加载

```vue
<GvLoading
  :loading="isLoading"
  fullscreen
  type="circular"
  :text="loadingText"
/>

<script setup>
const loadingSteps = ['加载配置...', '加载数据...', '渲染页面...']
const currentStep = ref(0)

const loadingText = computed(() => loadingSteps[currentStep.value])
</script>
```

## 最佳实践

1. **类型选择**：
   - circular: 默认推荐，优雅简洁
   - spinner: 轻量快速
   - dots: 轻松活泼
   - bars: 数据处理场景

2. **全屏加载**：
   - 页面初始化使用 fullscreen
   - 关键操作使用 fullscreen + lock
   - 提供清晰的加载文本

3. **局部加载**：
   - 卡片、表格等容器内使用局部加载
   - 避免过多的全屏加载影响体验

4. **加载文本**：
   - 长时间加载建议添加文本说明
   - 分步加载显示当前步骤

5. **性能考虑**：
   - 避免过多的同时加载动画
   - 短时间操作（< 300ms）可不显示加载

6. **用户体验**：
   - 提供准确的加载时长预期
   - 避免无限加载，设置超时处理
   - 加载失败时提供错误提示

## 组合使用

### 配合 GvButton

```vue
<GvButton :loading="isSubmitting" @click="handleSubmit">
  提交
</GvButton>
```

### 配合 GvCard

```vue
<GvCard>
  <GvLoading :loading="isLoading">
    <CardContent />
  </GvLoading>
</GvCard>
```

### 配合 GvModal

```vue
<GvModal v-model="visible" :confirm-loading="isSubmitting">
  <GvLoading :loading="isDataLoading">
    <FormContent />
  </GvLoading>
</GvModal>
```
