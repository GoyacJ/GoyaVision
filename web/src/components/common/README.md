# 通用状态组件

用于统一展示 Loading、Error、Empty 三种常见状态的组件库。

## 组件列表

- **LoadingState** - 加载中状态
- **ErrorState** - 错误状态
- **EmptyState** - 空状态

---

## LoadingState

显示加载中状态，带有旋转的加载指示器。

### Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `size` | `'small' \| 'medium' \| 'large'` | `'medium'` | 加载指示器大小 |
| `message` | `string` | - | 加载提示文本 |
| `fullscreen` | `boolean` | `false` | 是否全屏显示 |

### 使用示例

```vue
<script setup lang="ts">
import { LoadingState } from '@/components/common'
</script>

<template>
  <!-- 基础用法 -->
  <LoadingState />

  <!-- 带提示文本 -->
  <LoadingState message="加载中..." />

  <!-- 小尺寸 -->
  <LoadingState size="small" message="加载中..." />

  <!-- 全屏加载 -->
  <LoadingState fullscreen message="正在处理..." />
</template>
```

### 配合 useAsyncData 使用

```vue
<script setup lang="ts">
import { LoadingState } from '@/components/common'
import { useAsyncData } from '@/composables/useAsyncData'
import { assetApi } from '@/api/modules/asset'

const { data, isLoading } = useAsyncData(
  () => assetApi.list(),
  { immediate: true }
)
</script>

<template>
  <LoadingState v-if="isLoading" message="加载资产列表..." />
  <div v-else>
    <!-- 数据展示 -->
  </div>
</template>
```

---

## ErrorState

显示错误状态，包含错误图标、提示信息和重试按钮。

### Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `error` | `Error \| null` | - | 错误对象 |
| `title` | `string` | `'加载失败'` | 错误标题 |
| `message` | `string` | - | 错误描述（未提供时使用 error.message） |
| `retryText` | `string` | `'重试'` | 重试按钮文本 |
| `showRetry` | `boolean` | `true` | 是否显示重试按钮 |

### Emits

| 事件 | 参数 | 说明 |
|------|------|------|
| `retry` | - | 点击重试按钮时触发 |

### 使用示例

```vue
<script setup lang="ts">
import { ErrorState } from '@/components/common'
import { useAsyncData } from '@/composables/useAsyncData'
import { assetApi } from '@/api/modules/asset'

const { data, error, isLoading, execute } = useAsyncData(
  () => assetApi.list(),
  { immediate: true }
)
</script>

<template>
  <LoadingState v-if="isLoading" />
  <ErrorState
    v-else-if="error"
    :error="error"
    @retry="execute"
  />
  <div v-else>
    <!-- 数据展示 -->
  </div>
</template>
```

### 自定义错误信息

```vue
<template>
  <ErrorState
    title="网络连接失败"
    message="请检查网络连接后重试"
    retry-text="重新加载"
    @retry="handleRetry"
  />
</template>
```

### 隐藏重试按钮

```vue
<template>
  <ErrorState
    title="权限不足"
    message="您没有访问此资源的权限"
    :show-retry="false"
  />
</template>
```

---

## EmptyState

显示空状态，包含图标、提示信息和可选的操作按钮。

### Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `icon` | `string` | `'📭'` | 空状态图标（emoji） |
| `title` | `string` | `'暂无数据'` | 空状态标题 |
| `description` | `string` | - | 空状态描述 |
| `actionText` | `string` | - | 操作按钮文本 |
| `showAction` | `boolean` | `false` | 是否显示操作按钮 |

### Emits

| 事件 | 参数 | 说明 |
|------|------|------|
| `action` | - | 点击操作按钮时触发 |

### 使用示例

```vue
<script setup lang="ts">
import { EmptyState } from '@/components/common'
import { computed } from 'vue'

const assets = ref([])
const isEmpty = computed(() => assets.value.length === 0)
</script>

<template>
  <EmptyState v-if="isEmpty" />
  <div v-else>
    <!-- 数据展示 -->
  </div>
</template>
```

### 自定义内容和操作

```vue
<template>
  <EmptyState
    icon="🎬"
    title="还没有媒体资产"
    description="开始上传您的第一个视频、图片或音频文件"
    action-text="上传资产"
    show-action
    @action="handleUpload"
  />
</template>
```

### 不同场景的图标

```vue
<template>
  <!-- 搜索无结果 -->
  <EmptyState
    icon="🔍"
    title="未找到相关内容"
    description="尝试使用其他关键词搜索"
  />

  <!-- 筛选无结果 -->
  <EmptyState
    icon="🎯"
    title="没有符合条件的项目"
    description="调整筛选条件后再试"
  />

  <!-- 历史记录为空 -->
  <EmptyState
    icon="📝"
    title="暂无历史记录"
    description="您的操作历史将显示在这里"
  />

  <!-- 收藏为空 -->
  <EmptyState
    icon="⭐"
    title="还没有收藏"
    description="收藏您喜欢的内容以便快速访问"
  />
</template>
```

---

## 完整示例：列表页面

结合三个状态组件的完整示例：

```vue
<script setup lang="ts">
import { computed } from 'vue'
import { LoadingState, ErrorState, EmptyState } from '@/components/common'
import { useAsyncData } from '@/composables/useAsyncData'
import { assetApi } from '@/api/modules/asset'

const {
  data: assetsData,
  error,
  isLoading,
  execute: loadAssets
} = useAsyncData(
  () => assetApi.list({ page: 1, page_size: 12 }),
  { immediate: true }
)

const assets = computed(() => assetsData.value?.data.items ?? [])
const isEmpty = computed(() => !isLoading.value && !error.value && assets.value.length === 0)
</script>

<template>
  <div class="assets-page">
    <!-- Loading 状态 -->
    <LoadingState v-if="isLoading" message="加载资产列表..." />

    <!-- Error 状态 -->
    <ErrorState
      v-else-if="error"
      :error="error"
      title="加载失败"
      @retry="loadAssets"
    />

    <!-- Empty 状态 -->
    <EmptyState
      v-else-if="isEmpty"
      icon="🎬"
      title="还没有媒体资产"
      description="开始上传您的第一个视频、图片或音频文件"
      action-text="上传资产"
      show-action
      @action="handleUpload"
    />

    <!-- 正常内容 -->
    <div v-else class="assets-grid">
      <AssetCard
        v-for="asset in assets"
        :key="asset.id"
        :asset="asset"
      />
    </div>
  </div>
</template>
```

---

## 设计原则

这些状态组件遵循 GoyaVision 克制设计系统：

### 1. 色彩
- 使用新的主色 `#4F5B93`（primary.600）
- 中性灰色系 `#525252`、`#737373`
- 错误色 `#EF4444`

### 2. 排版
- 标题：16-18px，font-weight: 500-600
- 描述：14px，font-weight: 400
- 使用 Design Tokens 中的字距和行高

### 3. 间距
- 组件内边距：32px
- 元素间距：8px、24px
- 最小高度：400px

### 4. 动画
- 加载指示器：0.8s 线性旋转
- 按钮过渡：150ms
- 无过度动画

### 5. 无障碍
- 语义化 HTML
- 键盘可访问（按钮支持 focus）
- 清晰的视觉反馈

---

## TypeScript 支持

所有组件都提供完整的 TypeScript 类型定义：

```typescript
import type {
  LoadingStateProps,
  ErrorStateProps,
  ErrorStateEmits,
  EmptyStateProps,
  EmptyStateEmits
} from '@/components/common'
```

---

## 浏览器兼容性

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+

---

## 更新日志

### v1.0.0 (2026-02-05)
- ✅ 初始版本
- ✅ LoadingState 组件
- ✅ ErrorState 组件
- ✅ EmptyState 组件
- ✅ 完整的 TypeScript 类型定义
- ✅ 遵循克制设计系统
