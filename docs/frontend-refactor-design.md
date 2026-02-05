# GoyaVision 前端重构设计方案

> **设计标准**：真实可上线产品级 UI/UX 设计
> **设计风格**：克制、现代、内容优先（参考 Medium / Apple 官网 / 简洁视频平台）
> **核心理念**：干净、克制、强调内容而非装饰

---

## 目录

1. [现状分析](#1-现状分析)
2. [设计系统 (Design System)](#2-设计系统-design-system)
3. [技术架构重构](#3-技术架构重构)
4. [页面设计规范](#4-页面设计规范)
5. [组件设计原则](#5-组件设计原则)
6. [API 层设计](#6-api-层设计)
7. [状态管理策略](#7-状态管理策略)
8. [路由与导航](#8-路由与导航)
9. [性能优化策略](#9-性能优化策略)
10. [实施路线图](#10-实施路线图)

---

## 1. 现状分析

### 1.1 优势

✅ **已有基础**
- Design Tokens 系统完善（颜色、间距、阴影、圆角、动画）
- 基础组件库初步建立（Gv* 系列组件）
- TypeScript + Vue 3 Composition API
- Tailwind CSS + Element Plus 混合使用
- API 层已分模块（auth, asset, source, operator, workflow, task）
- Pinia 状态管理已接入

### 1.2 问题与改进点

❌ **设计层面**
- 视觉风格不够统一（渐变色使用过多，装饰性强）
- 色彩对比度过高（`#667eea` → `#764ba2` 渐变过于鲜艳）
- 阴影和毛玻璃效果偏重（backdrop-filter, box-shadow 过度使用）
- 缺乏视觉层级的克制感

❌ **技术层面**
- Element Plus 与自定义组件混用，风格不统一
- 部分组件逻辑过重（如 `asset/index.vue` 近 1200 行）
- API 调用未统一 Loading/Error 处理
- 缺乏全局错误边界
- 类型定义不够完善（部分 `any` 使用）

---

## 2. 设计系统 (Design System)

### 2.1 色彩系统重构

**主色调：克制的蓝灰色系**

```typescript
// 新色彩系统 - 更低调、更专业
const colors = {
  // 主色：冷静的蓝灰色（降低饱和度）
  primary: {
    DEFAULT: '#4F5B93',  // 从 #667eea 降低饱和度
    50: '#F5F6FA',
    100: '#EBEDF5',
    200: '#D4D8E8',
    300: '#B3BAD5',
    400: '#8A94B8',
    500: '#4F5B93',
    600: '#3E4A7A',
    700: '#2F3A61',
    800: '#232D4B',
    900: '#1A2238',
  },

  // 辅助色：极简灰色系（内容优先）
  neutral: {
    50: '#FAFAFA',   // 背景色
    100: '#F5F5F5',  // 容器背景
    200: '#E5E5E5',  // 边框
    300: '#D4D4D4',  // 禁用状态
    400: '#A3A3A3',  // 占位符
    500: '#737373',  // 次要文本
    600: '#525252',  // 主要文本
    700: '#404040',  // 标题
    800: '#262626',  // 深色文本
    900: '#171717',  // 强调文本
  },

  // 语义色：功能性色彩
  success: '#10B981',  // 绿色
  warning: '#F59E0B',  // 橙色
  error: '#EF4444',    // 红色
  info: '#3B82F6',     // 蓝色
}
```

**弃用渐变色背景**
- ❌ 移除所有 `linear-gradient` 装饰性背景
- ✅ 仅在强调元素（CTA 按钮、重要标签）使用单色

---

### 2.2 排版系统

**字体层级（参考 Apple/Medium）**

```typescript
const typography = {
  // 页面标题
  h1: {
    size: '32px',
    weight: 700,
    lineHeight: 1.25,
    letterSpacing: '-0.02em',  // 紧凑字距
  },

  // 区块标题
  h2: {
    size: '24px',
    weight: 600,
    lineHeight: 1.3,
    letterSpacing: '-0.01em',
  },

  // 卡片标题
  h3: {
    size: '18px',
    weight: 600,
    lineHeight: 1.4,
    letterSpacing: '0',
  },

  // 正文
  body: {
    size: '15px',      // 比标准 14px 稍大，增强可读性
    weight: 400,
    lineHeight: 1.6,
    letterSpacing: '0',
  },

  // 次要文本
  caption: {
    size: '13px',
    weight: 400,
    lineHeight: 1.5,
    color: 'neutral.500',
  },

  // 标签/微型文本
  label: {
    size: '12px',
    weight: 500,
    lineHeight: 1.4,
    textTransform: 'uppercase',
    letterSpacing: '0.05em',
  },
}
```

---

### 2.3 间距系统（8px 基准）

```typescript
const spacing = {
  0: '0',
  1: '4px',    // 0.25rem
  2: '8px',    // 0.5rem  - 基准单位
  3: '12px',   // 0.75rem
  4: '16px',   // 1rem
  5: '20px',   // 1.25rem
  6: '24px',   // 1.5rem
  8: '32px',   // 2rem
  10: '40px',  // 2.5rem
  12: '48px',  // 3rem
  16: '64px',  // 4rem
  20: '80px',  // 5rem
}
```

**布局规则**
- 页面容器左右 Padding：`32px`（桌面）/ `16px`（移动）
- 卡片内边距：`24px`
- 元素间距：`16px`（紧密）/ `24px`（标准）/ `40px`（宽松）

---

### 2.4 视觉效果克制化

**阴影系统：从装饰转向功能**

```typescript
const shadows = {
  // 弃用：彩色阴影（primary-shadow, secondary-shadow）

  // 新方案：极简功能性阴影
  none: 'none',
  sm: '0 1px 2px rgba(0, 0, 0, 0.04)',           // 轻微层级
  DEFAULT: '0 1px 3px rgba(0, 0, 0, 0.06)',      // 标准卡片
  md: '0 4px 6px rgba(0, 0, 0, 0.07)',           // 浮动元素
  lg: '0 10px 15px rgba(0, 0, 0, 0.08)',         // 模态框
  xl: '0 20px 25px rgba(0, 0, 0, 0.10)',         // 抽屉/弹出层
}
```

**圆角：适度减小**

```typescript
const radius = {
  sm: '4px',    // 小元素（标签、徽章）
  DEFAULT: '6px',   // 按钮、输入框
  md: '8px',    // 卡片
  lg: '12px',   // 容器、模态框
  xl: '16px',   // 大型容器
  full: '9999px',  // 圆形头像
}
```

**弃用效果**
- ❌ 毛玻璃背景（`backdrop-filter: blur()`）
- ❌ 过度的 hover 动画（`transform: scale()`, `translateY()`）
- ❌ 装饰性渐变背景
- ✅ 仅保留必要的状态反馈（hover 透明度变化、focus 边框）

---

### 2.5 交互动画

**动画原则：迅速、自然、不干扰**

```typescript
const transitions = {
  fast: '150ms cubic-bezier(0.4, 0, 0.2, 1)',    // 快速反馈（hover）
  normal: '200ms cubic-bezier(0.4, 0, 0.2, 1)',  // 标准过渡
  slow: '300ms cubic-bezier(0.4, 0, 0.2, 1)',    // 页面切换
}
```

**禁止**
- ❌ 页面级动画（`fade-enter-active`）
- ❌ 列表动画（`transition-group`）
- ❌ 骨架屏过度使用

**允许**
- ✅ 按钮 hover/active 状态
- ✅ 模态框淡入淡出
- ✅ 加载指示器旋转

---

## 3. 技术架构重构

### 3.1 目录结构

```
web/src/
├── api/                    # API 层（统一封装）
│   ├── client.ts          # Axios 实例配置
│   ├── interceptors.ts    # 请求/响应拦截器
│   ├── types.ts           # API 通用类型
│   └── modules/           # 按业务模块分组
│       ├── asset.ts
│       ├── source.ts
│       ├── operator.ts
│       ├── workflow.ts
│       └── ...
├── components/            # 组件库
│   ├── base/             # 基础组件（替换 Element Plus）
│   │   ├── GvButton/
│   │   ├── GvInput/
│   │   ├── GvCard/
│   │   └── ...
│   ├── business/         # 业务组件
│   │   ├── AssetCard/
│   │   ├── TaskList/
│   │   └── ...
│   └── layout/           # 布局组件
│       ├── GvContainer/
│       ├── GvGrid/
│       └── ...
├── composables/          # 组合式函数
│   ├── useAsyncData.ts  # 统一数据加载
│   ├── usePagination.ts # 分页逻辑
│   ├── useTable.ts      # 表格逻辑
│   └── ...
├── layouts/              # 页面布局
│   ├── DefaultLayout.vue
│   ├── AuthLayout.vue
│   └── ...
├── views/                # 页面视图（仅组合逻辑）
│   ├── assets/
│   │   ├── index.vue    # 资产列表
│   │   └── detail.vue   # 资产详情
│   ├── sources/
│   ├── operators/
│   └── ...
├── stores/               # Pinia 状态管理
│   ├── user.ts
│   ├── app.ts           # 全局应用状态
│   └── ...
├── router/
│   ├── index.ts
│   ├── guards.ts        # 路由守卫
│   └── routes.ts
├── types/                # TypeScript 类型定义
│   ├── api.ts           # API 响应类型
│   ├── components.ts    # 组件 Props 类型
│   └── models.ts        # 业务模型
├── utils/                # 工具函数
│   ├── format.ts        # 格式化函数
│   ├── validate.ts      # 表单验证
│   └── ...
└── design-tokens/        # 设计令牌（保留）
```

---

### 3.2 API 层统一封装

#### 3.2.1 Axios 实例配置

**选择 Axios**（理由：拦截器、请求取消、类型友好）

```typescript
// api/client.ts
import axios, { type AxiosInstance, type AxiosRequestConfig } from 'axios'
import { useUserStore } from '@/stores/user'

const client: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
client.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// 响应拦截器
client.interceptors.response.use(
  (response) => response.data, // 直接返回 data
  async (error) => {
    const userStore = useUserStore()

    // Token 过期自动刷新
    if (error.response?.status === 401 && !error.config._retry) {
      error.config._retry = true
      try {
        await userStore.refreshAccessToken()
        return client(error.config)
      } catch {
        userStore.logout()
        return Promise.reject(error)
      }
    }

    return Promise.reject(error)
  }
)

export default client
```

#### 3.2.2 API 响应类型

```typescript
// api/types.ts

// 统一响应结构
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 分页响应
export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

// 错误响应
export interface ApiError {
  code: string
  message: string
  details?: Record<string, any>
}

// 请求状态
export type RequestStatus = 'idle' | 'pending' | 'success' | 'error'
```

#### 3.2.3 API 模块示例

```typescript
// api/modules/asset.ts
import client from '../client'
import type { ApiResponse, PaginatedResponse } from '../types'

export interface MediaAsset {
  id: string
  name: string
  type: 'video' | 'image' | 'audio' | 'stream'
  source_type: 'upload' | 'live' | 'vod' | 'generated'
  path: string
  format: string
  size: number
  duration?: number
  status: 'ready' | 'processing' | 'pending' | 'error'
  tags: string[]
  created_at: string
  updated_at: string
}

export interface AssetListParams {
  name?: string
  type?: MediaAsset['type']
  tags?: string
  page?: number
  page_size?: number
}

export interface AssetCreateData {
  name: string
  type: MediaAsset['type']
  source_type: MediaAsset['source_type']
  path?: string
  stream_url?: string
  source_id?: string
  tags?: string[]
}

export interface AssetUpdateData {
  name?: string
  status?: MediaAsset['status']
  tags?: string[]
}

export const assetApi = {
  // 获取资产列表
  list(params: AssetListParams) {
    return client.get<ApiResponse<PaginatedResponse<MediaAsset>>>('/assets', { params })
  },

  // 获取资产详情
  get(id: string) {
    return client.get<ApiResponse<MediaAsset>>(`/assets/${id}`)
  },

  // 创建资产
  create(data: AssetCreateData) {
    return client.post<ApiResponse<MediaAsset>>('/assets', data)
  },

  // 上传文件
  upload(file: File, type: string, name: string, tags: string[] = []) {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('type', type)
    formData.append('name', name)
    formData.append('tags', JSON.stringify(tags))

    return client.post<ApiResponse<MediaAsset>>('/assets/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },

  // 更新资产
  update(id: string, data: AssetUpdateData) {
    return client.put<ApiResponse<MediaAsset>>(`/assets/${id}`, data)
  },

  // 删除资产
  delete(id: string) {
    return client.delete<ApiResponse<void>>(`/assets/${id}`)
  },

  // 获取标签列表
  getTags() {
    return client.get<ApiResponse<{ tags: string[] }>>('/assets/tags')
  },
}
```

---

### 3.3 统一 Loading/Error 处理

#### 3.3.1 Composable: useAsyncData

```typescript
// composables/useAsyncData.ts
import { ref, type Ref } from 'vue'
import type { RequestStatus } from '@/api/types'

export interface UseAsyncDataOptions<T> {
  immediate?: boolean           // 是否立即执行
  initialData?: T              // 初始数据
  onSuccess?: (data: T) => void
  onError?: (error: Error) => void
}

export function useAsyncData<T>(
  asyncFn: (...args: any[]) => Promise<T>,
  options: UseAsyncDataOptions<T> = {}
) {
  const { immediate = false, initialData, onSuccess, onError } = options

  const data = ref<T | undefined>(initialData) as Ref<T | undefined>
  const error = ref<Error | null>(null)
  const status = ref<RequestStatus>('idle')
  const isLoading = computed(() => status.value === 'pending')

  async function execute(...args: any[]) {
    status.value = 'pending'
    error.value = null

    try {
      const result = await asyncFn(...args)
      data.value = result
      status.value = 'success'
      onSuccess?.(result)
      return result
    } catch (err) {
      error.value = err instanceof Error ? err : new Error(String(err))
      status.value = 'error'
      onError?.(error.value)
      throw error.value
    }
  }

  function reset() {
    data.value = initialData
    error.value = null
    status.value = 'idle'
  }

  if (immediate) {
    execute()
  }

  return {
    data,
    error,
    status,
    isLoading,
    execute,
    reset,
  }
}
```

#### 3.3.2 使用示例

```vue
<!-- views/assets/index.vue -->
<script setup lang="ts">
import { computed } from 'vue'
import { assetApi, type AssetListParams } from '@/api/modules/asset'
import { useAsyncData } from '@/composables/useAsyncData'
import AssetCard from '@/components/business/AssetCard/index.vue'
import LoadingState from '@/components/common/LoadingState.vue'
import ErrorState from '@/components/common/ErrorState.vue'
import EmptyState from '@/components/common/EmptyState.vue'

const params = ref<AssetListParams>({ page: 1, page_size: 12 })

const {
  data: assetsData,
  error,
  isLoading,
  execute: loadAssets,
} = useAsyncData(
  () => assetApi.list(params.value),
  { immediate: true }
)

const assets = computed(() => assetsData.value?.data.items ?? [])
const total = computed(() => assetsData.value?.data.total ?? 0)
</script>

<template>
  <div class="assets-page">
    <!-- Loading 状态 -->
    <LoadingState v-if="isLoading" />

    <!-- Error 状态 -->
    <ErrorState v-else-if="error" :error="error" @retry="loadAssets" />

    <!-- Empty 状态 -->
    <EmptyState v-else-if="assets.length === 0" />

    <!-- 正常内容 -->
    <div v-else class="assets-grid">
      <AssetCard v-for="asset in assets" :key="asset.id" :asset="asset" />
    </div>
  </div>
</template>
```

---

### 3.4 TypeScript 类型规范

#### 3.4.1 禁止 any

```typescript
// ❌ 禁止
function handleData(data: any) { }

// ✅ 使用 unknown + 类型守卫
function handleData(data: unknown) {
  if (isMediaAsset(data)) {
    // 类型安全
  }
}

// 类型守卫
function isMediaAsset(value: unknown): value is MediaAsset {
  return (
    typeof value === 'object' &&
    value !== null &&
    'id' in value &&
    'name' in value
  )
}
```

#### 3.4.2 组件 Props 类型

```typescript
// components/business/AssetCard/types.ts
import type { MediaAsset } from '@/api/modules/asset'

export interface AssetCardProps {
  asset: MediaAsset
  selectable?: boolean
  selected?: boolean
}

export interface AssetCardEmits {
  (e: 'view', asset: MediaAsset): void
  (e: 'edit', asset: MediaAsset): void
  (e: 'delete', asset: MediaAsset): void
  (e: 'select', selected: boolean): void
}
```

```vue
<!-- components/business/AssetCard/index.vue -->
<script setup lang="ts">
import type { AssetCardProps, AssetCardEmits } from './types'

const props = withDefaults(defineProps<AssetCardProps>(), {
  selectable: false,
  selected: false,
})

const emit = defineEmits<AssetCardEmits>()
</script>
```

---

## 4. 页面设计规范

### 4.1 页面布局模式

#### 模式 1: 双栏布局（资产管理、算子中心）

```
┌─────────────────────────────────────────────────────────┐
│  Sidebar (240px)  │  Main Content (flex-1)              │
│  ┌──────────────┐ │  ┌────────────────────────────────┐ │
│  │ Filters      │ │  │  Toolbar (Search + Actions)    │ │
│  │ - Type       │ │  └────────────────────────────────┘ │
│  │ - Tags       │ │                                      │
│  │ - Status     │ │  ┌────────────────────────────────┐ │
│  └──────────────┘ │  │  Grid / List View              │ │
│                   │  │  ┌────┐ ┌────┐ ┌────┐          │ │
│                   │  │  │Card│ │Card│ │Card│          │ │
│                   │  │  └────┘ └────┘ └────┘          │ │
│                   │  └────────────────────────────────┘ │
│                   │                                      │
│                   │  [Pagination]                        │
└─────────────────────────────────────────────────────────┘
```

#### 模式 2: 单栏布局（任务中心、工作流）

```
┌─────────────────────────────────────────────────────────┐
│  Page Header                                            │
│  ┌────────────────────────────────────────────────────┐ │
│  │  Title + Description    [Actions]                  │ │
│  └────────────────────────────────────────────────────┘ │
│                                                          │
│  Filters Bar                                            │
│  ┌────────────────────────────────────────────────────┐ │
│  │  [Status] [Type] [Date Range]  [Search]           │ │
│  └────────────────────────────────────────────────────┘ │
│                                                          │
│  Content (Table / Cards)                                │
│  ┌────────────────────────────────────────────────────┐ │
│  │  [Table with sorting/filtering]                    │ │
│  │  ...                                               │ │
│  └────────────────────────────────────────────────────┘ │
│                                                          │
│  [Pagination]                                           │
└─────────────────────────────────────────────────────────┘
```

### 4.2 视觉层级

```
1. 页面背景     : #FAFAFA (neutral.50)
2. 容器背景     : #FFFFFF (白色卡片)
3. 边框颜色     : #E5E5E5 (neutral.200)
4. 主要文本     : #262626 (neutral.800)
5. 次要文本     : #737373 (neutral.500)
6. 占位符文本   : #A3A3A3 (neutral.400)
```

### 4.3 状态组件

#### Loading State

```vue
<!-- components/common/LoadingState.vue -->
<template>
  <div class="loading-state">
    <div class="spinner" />
    <p v-if="message" class="message">{{ message }}</p>
  </div>
</template>

<style scoped>
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #E5E5E5;
  border-top-color: #4F5B93;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.message {
  margin-top: 16px;
  color: #737373;
  font-size: 14px;
}
</style>
```

#### Error State

```vue
<!-- components/common/ErrorState.vue -->
<template>
  <div class="error-state">
    <div class="icon">⚠️</div>
    <h3 class="title">加载失败</h3>
    <p class="message">{{ error?.message || '发生未知错误' }}</p>
    <button class="retry-btn" @click="$emit('retry')">
      重试
    </button>
  </div>
</template>

<style scoped>
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

.icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.title {
  font-size: 18px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 8px;
}

.message {
  font-size: 14px;
  color: #737373;
  margin-bottom: 24px;
}

.retry-btn {
  padding: 8px 16px;
  background: #4F5B93;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: opacity 0.15s;
}

.retry-btn:hover {
  opacity: 0.9;
}
</style>
```

#### Empty State

```vue
<!-- components/common/EmptyState.vue -->
<template>
  <div class="empty-state">
    <div class="icon">📭</div>
    <h3 class="title">{{ title || '暂无数据' }}</h3>
    <p v-if="description" class="description">{{ description }}</p>
    <button v-if="actionText" class="action-btn" @click="$emit('action')">
      {{ actionText }}
    </button>
  </div>
</template>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

.icon {
  font-size: 64px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.title {
  font-size: 16px;
  font-weight: 500;
  color: #525252;
  margin-bottom: 8px;
}

.description {
  font-size: 14px;
  color: #737373;
  margin-bottom: 24px;
  text-align: center;
}

.action-btn {
  padding: 8px 16px;
  background: #4F5B93;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: opacity 0.15s;
}
</style>
```

---

## 5. 组件设计原则

### 5.1 职责单一原则

**❌ 反例：逻辑堆砌**

```vue
<!-- 1200 行的巨型组件 -->
<script setup lang="ts">
// 混合了：数据获取 + 表单验证 + 上传逻辑 + UI 状态管理
const { data, loading, error } = useAsyncData(...)
const uploadForm = reactive({ ... })
const validateForm = () => { ... }
const handleUpload = () => { ... }
// ... 更多逻辑
</script>
```

**✅ 正例：拆分关注点**

```vue
<!-- views/assets/index.vue - 页面层（组合逻辑） -->
<script setup lang="ts">
import { useAssetList } from './composables/useAssetList'
import { useAssetFilters } from './composables/useAssetFilters'
import AssetList from './components/AssetList.vue'
import AssetFilters from './components/AssetFilters.vue'

const { assets, loading, error, refresh } = useAssetList()
const { filters, updateFilter } = useAssetFilters()
</script>

<template>
  <div class="assets-page">
    <AssetFilters :filters="filters" @update="updateFilter" />
    <AssetList
      :assets="assets"
      :loading="loading"
      :error="error"
      @refresh="refresh"
    />
  </div>
</template>
```

```typescript
// composables/useAssetList.ts - 业务逻辑
export function useAssetList() {
  const params = ref({ page: 1, page_size: 12 })

  const { data, error, isLoading, execute } = useAsyncData(
    () => assetApi.list(params.value),
    { immediate: true }
  )

  const assets = computed(() => data.value?.data.items ?? [])

  return { assets, loading: isLoading, error, refresh: execute }
}
```

### 5.2 组件分类

#### 5.2.1 基础组件（Base Components）

**特点**
- 无业务逻辑
- 接受 Props，发出 Events
- 可在任意项目复用

```vue
<!-- components/base/GvButton/index.vue -->
<script setup lang="ts">
import type { ButtonProps } from './types'

const props = withDefaults(defineProps<ButtonProps>(), {
  variant: 'solid',
  size: 'medium',
  disabled: false,
  loading: false,
})
</script>

<template>
  <button
    :class="[
      'gv-button',
      `gv-button--${variant}`,
      `gv-button--${size}`,
      { 'is-disabled': disabled, 'is-loading': loading }
    ]"
    :disabled="disabled || loading"
  >
    <span v-if="loading" class="gv-button__spinner" />
    <slot />
  </button>
</template>
```

#### 5.2.2 业务组件（Business Components）

**特点**
- 包含 GoyaVision 特定业务逻辑
- 可直接调用 API
- 不可跨项目复用

```vue
<!-- components/business/AssetCard/index.vue -->
<script setup lang="ts">
import type { MediaAsset } from '@/api/modules/asset'
import { formatFileSize, formatDuration } from '@/utils/format'

interface Props {
  asset: MediaAsset
}

const props = defineProps<Props>()
const emit = defineEmits<{
  view: [asset: MediaAsset]
  edit: [asset: MediaAsset]
  delete: [asset: MediaAsset]
}>()
</script>

<template>
  <div class="asset-card">
    <div class="asset-card__preview">
      <img v-if="asset.type === 'image'" :src="asset.path" :alt="asset.name" />
      <div v-else class="asset-card__placeholder">
        {{ asset.type }}
      </div>
    </div>

    <div class="asset-card__content">
      <h3 class="asset-card__title">{{ asset.name }}</h3>
      <div class="asset-card__meta">
        <span>{{ formatFileSize(asset.size) }}</span>
        <span v-if="asset.duration">{{ formatDuration(asset.duration) }}</span>
      </div>
    </div>

    <div class="asset-card__actions">
      <button @click="emit('view', asset)">查看</button>
      <button @click="emit('edit', asset)">编辑</button>
      <button @click="emit('delete', asset)">删除</button>
    </div>
  </div>
</template>

<style scoped>
.asset-card {
  background: white;
  border: 1px solid #E5E5E5;
  border-radius: 8px;
  overflow: hidden;
  transition: box-shadow 0.15s;
}

.asset-card:hover {
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.07);
}

.asset-card__preview {
  aspect-ratio: 16/9;
  background: #F5F5F5;
  overflow: hidden;
}

.asset-card__content {
  padding: 16px;
}

.asset-card__title {
  font-size: 15px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.asset-card__meta {
  display: flex;
  gap: 12px;
  font-size: 13px;
  color: #737373;
}

.asset-card__actions {
  padding: 12px 16px;
  border-top: 1px solid #E5E5E5;
  display: flex;
  gap: 8px;
}
</style>
```

---

## 6. API 层设计

### 6.1 统一错误处理

```typescript
// api/interceptors.ts
import { ElMessage } from 'element-plus'
import type { AxiosError } from 'axios'

export function handleApiError(error: AxiosError) {
  const status = error.response?.status
  const message = error.response?.data?.message || '请求失败'

  switch (status) {
    case 400:
      ElMessage.error(`请求错误：${message}`)
      break
    case 401:
      // 已在拦截器处理 Token 刷新
      break
    case 403:
      ElMessage.error('无权限访问')
      break
    case 404:
      ElMessage.error('资源不存在')
      break
    case 500:
      ElMessage.error('服务器错误')
      break
    default:
      ElMessage.error(message)
  }
}
```

### 6.2 请求取消

```typescript
// composables/useAsyncData.ts 增强版
import { ref, onBeforeUnmount } from 'vue'
import axios, { type CancelTokenSource } from 'axios'

export function useAsyncData<T>(asyncFn: () => Promise<T>) {
  const cancelToken = ref<CancelTokenSource>()

  async function execute() {
    // 取消之前的请求
    cancelToken.value?.cancel('New request started')

    // 创建新的 cancel token
    cancelToken.value = axios.CancelToken.source()

    try {
      const result = await asyncFn()
      return result
    } catch (err) {
      if (!axios.isCancel(err)) {
        throw err
      }
    }
  }

  // 组件卸载时取消请求
  onBeforeUnmount(() => {
    cancelToken.value?.cancel('Component unmounted')
  })

  return { execute }
}
```

---

## 7. 状态管理策略

### 7.1 全局状态（Pinia）

**仅存储以下内容**
- 用户信息（token, profile, permissions）
- 应用配置（theme, locale, sidebar collapsed）
- 跨页面共享数据（websocket connection）

```typescript
// stores/app.ts
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', () => {
  const theme = ref<'light' | 'dark'>('light')
  const sidebarCollapsed = ref(false)
  const locale = ref('zh-CN')

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setTheme(newTheme: 'light' | 'dark') {
    theme.value = newTheme
    document.documentElement.classList.toggle('dark', newTheme === 'dark')
  }

  return { theme, sidebarCollapsed, locale, toggleSidebar, setTheme }
})
```

### 7.2 局部状态

**使用 Composables 管理页面状态**

```typescript
// views/assets/composables/useAssetList.ts
export function useAssetList() {
  const filters = ref({ type: null, tags: null })
  const pagination = ref({ page: 1, page_size: 12 })

  const { data, execute } = useAsyncData(
    () => assetApi.list({ ...filters.value, ...pagination.value })
  )

  function updateFilters(newFilters: Partial<typeof filters.value>) {
    filters.value = { ...filters.value, ...newFilters }
    pagination.value.page = 1
    execute()
  }

  return { filters, pagination, data, updateFilters }
}
```

---

## 8. 路由与导航

### 8.1 路由配置

```typescript
// router/routes.ts
import type { RouteRecordRaw } from 'vue-router'

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    redirect: '/assets',
    children: [
      {
        path: 'assets',
        name: 'Assets',
        component: () => import('@/views/assets/index.vue'),
        meta: { title: '媒体资产', icon: 'film' },
      },
      {
        path: 'assets/:id',
        name: 'AssetDetail',
        component: () => import('@/views/assets/detail.vue'),
        meta: { title: '资产详情', hidden: true },
      },
      // ...
    ],
  },
]
```

### 8.2 路由守卫

```typescript
// router/guards.ts
import { useUserStore } from '@/stores/user'

export function setupRouterGuards(router: Router) {
  router.beforeEach(async (to, from, next) => {
    const userStore = useUserStore()

    if (to.path === '/login') {
      return next()
    }

    if (!userStore.isLoggedIn) {
      return next('/login')
    }

    // 加载用户信息
    if (!userStore.userInfo) {
      try {
        await userStore.getProfile()
      } catch {
        return next('/login')
      }
    }

    // 权限检查
    if (to.meta.permission && !userStore.hasPermission(to.meta.permission)) {
      return next('/403')
    }

    next()
  })
}
```

---

## 9. 性能优化策略

### 9.1 代码分割

```typescript
// 路由级代码分割（已实现）
const AssetList = () => import('@/views/assets/index.vue')

// 组件级按需加载
const HeavyChart = defineAsyncComponent(() => import('./HeavyChart.vue'))
```

### 9.2 图片优化

```vue
<template>
  <!-- 懒加载 + 占位符 -->
  <img
    v-lazy="asset.thumbnail"
    :alt="asset.name"
    class="asset-thumbnail"
    loading="lazy"
  />
</template>
```

### 9.3 虚拟滚动（大列表）

```vue
<script setup lang="ts">
import { useVirtualList } from '@vueuse/core'

const { list, containerProps, wrapperProps } = useVirtualList(
  largeAssetList,
  { itemHeight: 200 }
)
</script>

<template>
  <div v-bind="containerProps" class="asset-list">
    <div v-bind="wrapperProps">
      <AssetCard v-for="item in list" :key="item.data.id" :asset="item.data" />
    </div>
  </div>
</template>
```

---

## 10. 实施路线图

### Phase 1: 设计系统重构（2 周）

**Week 1**
- [ ] 更新 Design Tokens（色彩、排版、间距、阴影）
- [ ] 重构基础组件（GvButton, GvInput, GvCard, GvModal）
- [ ] 创建状态组件（LoadingState, ErrorState, EmptyState）
- [ ] 编写 Storybook 文档

**Week 2**
- [ ] 重构 Layout 组件（移除毛玻璃、渐变色）
- [ ] 统一组件样式（圆角、阴影、动画）
- [ ] 创建布局组件（PageHeader, PageContainer, Sidebar）

### Phase 2: API 层与 Composables（1 周）

- [ ] 统一 Axios 配置和拦截器
- [ ] 重写所有 API 模块（完善类型定义）
- [ ] 实现 useAsyncData, usePagination, useTable
- [ ] 全局错误处理

### Phase 3: 页面重构（3 周）

**优先级排序**
1. 登录页（AuthLayout）
2. 媒体资产管理（双栏布局）
3. 媒体源管理
4. 算子中心
5. 工作流管理
6. 任务中心

**每个页面必须**
- 拆分 Composables 提取逻辑
- 使用统一的 Loading/Error/Empty 状态
- 类型定义完整（无 any）
- 组件职责单一

### Phase 4: 测试与优化（1 周）

- [ ] 单元测试（Composables）
- [ ] 组件测试（Vitest + Testing Library）
- [ ] E2E 测试（Playwright）
- [ ] 性能优化（Lighthouse 评分 > 90）
- [ ] 无障碍性检查（WCAG AA 标准）

---

## 附录 A: 禁止事项清单

### 视觉设计

- ❌ 彩色渐变背景（linear-gradient）
- ❌ 毛玻璃效果（backdrop-filter: blur）
- ❌ 彩色阴影（box-shadow 带颜色）
- ❌ 过度动画（transition > 300ms）
- ❌ 装饰性图标（非功能性）
- ❌ 过小的字体（< 12px）

### 技术实现

- ❌ 使用 `any` 类型
- ❌ 组件超过 500 行
- ❌ 直接在组件内调用 API（应使用 Composables）
- ❌ 在 `<template>` 中写复杂逻辑
- ❌ 忽略 Loading/Error 状态
- ❌ 混用 Element Plus 和自定义组件样式

---

## 附录 B: 检查清单

每完成一个页面/组件，必须通过以下检查：

### 代码质量

- [ ] TypeScript 严格模式无错误
- [ ] ESLint 无警告
- [ ] 所有 Props 有类型定义
- [ ] 所有 Emits 有类型定义
- [ ] 无 `any` 类型使用

### 功能完整性

- [ ] Loading 状态正常显示
- [ ] Error 状态可重试
- [ ] Empty 状态有提示
- [ ] 所有表单有验证
- [ ] 所有操作有反馈

### 视觉一致性

- [ ] 使用 Design Tokens 变量
- [ ] 圆角符合规范（4/6/8/12px）
- [ ] 阴影符合规范（sm/md/lg）
- [ ] 间距使用 8px 基准
- [ ] 颜色符合色彩系统

### 性能

- [ ] 图片使用懒加载
- [ ] 大列表使用虚拟滚动
- [ ] 路由使用懒加载
- [ ] 无不必要的重渲染

---

## 总结

本重构方案以**内容优先**为核心，通过：

1. **克制的设计系统**：降低视觉噪音，聚焦内容
2. **严格的类型系统**：TypeScript 覆盖率 100%
3. **清晰的分层架构**：API → Composables → Components → Views
4. **统一的状态处理**：Loading/Error/Empty 标准化
5. **职责单一的组件**：每个组件不超过 500 行

打造一个**专业、可维护、可上线**的前端应用。

---

**版本**：v1.0
**作者**：GoyaVision 前端架构团队
**更新日期**：2026-02-05
