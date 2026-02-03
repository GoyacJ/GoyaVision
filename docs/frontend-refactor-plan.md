# GoyaVision 前端重构方案

> 参考 Google Material Design 3，构建现代化、标准化、可维护的前端架构

---

## 📋 目录

1. [项目背景](#项目背景)
2. [当前状态分析](#当前状态分析)
3. [重构目标](#重构目标)
4. [技术选型](#技术选型)
5. [架构设计](#架构设计)
6. [组件体系](#组件体系)
7. [设计系统](#设计系统)
8. [实施计划](#实施计划)
9. [规范文档](#规范文档)
10. [风险评估](#风险评估)

---

## 项目背景

### 当前技术栈

```json
{
  "框架": "Vue 3.4 + TypeScript 5.3",
  "构建工具": "Vite 5.0",
  "UI 库": "Element Plus 2.5",
  "状态管理": "Pinia 2.1",
  "路由": "Vue Router 4.2",
  "HTTP 客户端": "Axios 1.6",
  "样式方案": "Scoped CSS + CSS Variables"
}
```

### 重构驱动因素

1. **缺乏统一的设计系统** - 样式散落在各个组件中，难以维护
2. **组件复用性低** - 缺少公共组件库，大量重复代码
3. **样式方案原始** - 手写 CSS，效率低，一致性差
4. **缺少规范文档** - AI 辅助开发时无法准确调用组件
5. **设计风格不统一** - 需要对齐 Google Material Design 3 标准

---

## 当前状态分析

### ✅ 优势

- Vue 3 Composition API 使用规范
- TypeScript 类型安全
- 路径别名配置完善（`@/`）
- 基础权限指令已实现
- Pinia 状态管理清晰

### ❌ 不足

- **样式管理混乱**
  - 每个组件独立写样式
  - 缺少统一的 Design Token
  - 大量重复的样式代码

- **组件复用性差**
  - 只有 1 个公共组件（HLSPreview）
  - 缺少 Button、Card、Form 等基础组件封装
  - 布局组件功能单一

- **缺少工程化工具**
  - 无 CSS 预处理器或工具类框架
  - 无组件文档生成工具
  - 无自动化测试

- **设计系统缺失**
  - 无统一的配色、间距、圆角规范
  - 无组件使用规范文档
  - AI 无法准确理解组件用法

---

## 重构目标

### 1. 技术架构目标

- ✅ 引入 **Tailwind CSS** 作为工具类 CSS 框架
- ✅ 保留 **Element Plus** 作为基础组件库
- ✅ 建立 **Material Design 3** 设计系统
- ✅ 封装 **公共组件库**（30+ 组件）
- ✅ 配置 **Storybook** 组件文档

### 2. 代码质量目标

- ✅ 组件复用率 > 70%
- ✅ 样式一致性 > 95%
- ✅ TypeScript 类型覆盖率 100%
- ✅ 组件文档覆盖率 100%
- ✅ 性能优化（减少 30% CSS 体积）

### 3. 开发体验目标

- ✅ AI 可精准调用组件（通过规范文档）
- ✅ 新页面开发效率提升 50%
- ✅ 样式调整时间减少 70%
- ✅ 代码可维护性提升 80%

---

## 技术选型

### 1. CSS 框架方案对比

| 方案 | 优势 | 劣势 | 推荐度 |
|------|------|------|--------|
| **Tailwind CSS** | 🔥 工具类丰富<br>🔥 按需加载<br>🔥 生态成熟<br>🔥 可自定义主题 | ⚠️ HTML 类名多 | ⭐⭐⭐⭐⭐ |
| UnoCSS | 性能极高<br>灵活性强 | 生态较新<br>社区较小 | ⭐⭐⭐⭐ |
| Windi CSS | 按需编译快 | 项目已停更 | ⭐⭐ |
| 手写 CSS | 灵活度最高 | 效率低<br>维护难 | ⭐ |

**最终选择：Tailwind CSS v3.4**

**理由：**
- 生态最成熟，插件丰富
- 官方 Material Design 3 主题支持
- 与 Vue 3 集成良好
- AI 训练数据充足，辅助开发友好

### 2. 组件库策略

**双层架构：Element Plus + 自定义组件**

```
┌─────────────────────────────────────┐
│     GoyaVision 业务组件层           │  ← 高度定制的业务组件
├─────────────────────────────────────┤
│     GoyaVision 基础组件层           │  ← 封装 Element Plus + 自定义
├─────────────────────────────────────┤
│     Element Plus 原子组件层         │  ← 第三方 UI 库
└─────────────────────────────────────┘
```

**保留 Element Plus 的原因：**
- ✅ 复杂组件（Table、Form、DatePicker）成熟稳定
- ✅ 国际化、主题、可访问性开箱即用
- ✅ 社区活跃，持续更新
- ✅ 减少重复造轮子

**自定义组件的范围：**
- 基础组件：Button、Card、Badge、Tag（Tailwind 样式）
- 布局组件：Container、Grid、Flex、Divider
- 业务组件：AssetCard、TaskCard、WorkflowNode

### 3. 设计系统方案

**参考：Google Material Design 3**

| 特性 | Material Design 3 | GoyaVision 适配 |
|------|-------------------|-----------------|
| **配色** | Dynamic Color（动态配色） | 基于主题色生成完整色板 |
| **形状** | Rounded Corners（圆角） | 统一 4 档圆角规范 |
| **排版** | Type Scale（字阶） | 9 档字号 + 4 档行高 |
| **间距** | 8px Grid（8 像素网格） | 统一间距系统 |
| **动效** | Motion（运动） | 4 种缓动曲线 |
| **状态** | State Layers（状态层） | 悬停、激活、禁用状态 |

---

## 架构设计

### 1. 项目结构重构

```
web/
├── src/
│   ├── assets/                    # 静态资源
│   │   ├── icons/                 # SVG 图标
│   │   ├── images/                # 图片资源
│   │   └── fonts/                 # 字体文件
│   │
│   ├── components/                # 公共组件库 ⭐ 核心重构区域
│   │   ├── base/                  # 基础组件（30+）
│   │   │   ├── GvButton/          # 按钮组件
│   │   │   │   ├── index.vue
│   │   │   │   ├── types.ts
│   │   │   │   └── README.md      # 组件文档
│   │   │   ├── GvCard/            # 卡片组件
│   │   │   ├── GvBadge/           # 徽章组件
│   │   │   ├── GvTag/             # 标签组件
│   │   │   ├── GvInput/           # 输入框组件
│   │   │   ├── GvSelect/          # 选择器组件
│   │   │   ├── GvModal/           # 模态框组件
│   │   │   ├── GvDrawer/          # 抽屉组件
│   │   │   ├── GvTooltip/         # 提示框组件
│   │   │   ├── GvAlert/           # 警告框组件
│   │   │   └── ...                # 30+ 基础组件
│   │   │
│   │   ├── layout/                # 布局组件
│   │   │   ├── GvContainer/       # 容器组件
│   │   │   ├── GvGrid/            # 栅格组件
│   │   │   ├── GvFlex/            # 弹性布局
│   │   │   ├── GvDivider/         # 分割线
│   │   │   ├── GvHeader/          # 头部组件
│   │   │   ├── GvSidebar/         # 侧边栏组件
│   │   │   ├── GvFooter/          # 底部组件
│   │   │   └── GvPageLayout/      # 页面布局
│   │   │
│   │   ├── business/              # 业务组件
│   │   │   ├── AssetCard/         # 资产卡片
│   │   │   ├── TaskCard/          # 任务卡片
│   │   │   ├── WorkflowNode/      # 工作流节点
│   │   │   ├── UserAvatar/        # 用户头像
│   │   │   ├── StatusBadge/       # 状态徽章
│   │   │   └── ...                # 业务组件
│   │   │
│   │   ├── feedback/              # 反馈组件
│   │   │   ├── GvLoading/         # 加载组件
│   │   │   ├── GvEmpty/           # 空状态组件
│   │   │   ├── GvResult/          # 结果页组件
│   │   │   └── GvSkeleton/        # 骨架屏
│   │   │
│   │   └── index.ts               # 组件统一导出
│   │
│   ├── composables/               # 组合式函数 ⭐ 新增
│   │   ├── useTheme.ts            # 主题切换
│   │   ├── useBreakpoint.ts       # 响应式断点
│   │   ├── useModal.ts            # 模态框管理
│   │   ├── useToast.ts            # 消息提示
│   │   ├── useClipboard.ts        # 剪贴板
│   │   └── ...                    # 更多 hooks
│   │
│   ├── styles/                    # 样式系统 ⭐ 核心重构区域
│   │   ├── tailwind.css           # Tailwind 入口
│   │   ├── variables.css          # CSS 变量（Design Token）
│   │   ├── transitions.css        # 动画过渡
│   │   ├── utilities.css          # 自定义工具类
│   │   └── themes/                # 主题文件
│   │       ├── light.css          # 浅色主题
│   │       ├── dark.css           # 深色主题
│   │       └── default.css        # 默认主题
│   │
│   ├── design-tokens/             # 设计令牌 ⭐ 新增
│   │   ├── colors.ts              # 颜色系统
│   │   ├── spacing.ts             # 间距系统
│   │   ├── typography.ts          # 字体系统
│   │   ├── shadows.ts             # 阴影系统
│   │   ├── radius.ts              # 圆角系统
│   │   └── index.ts               # 统一导出
│   │
│   ├── api/                       # API 接口（保持不变）
│   ├── router/                    # 路由配置（保持不变）
│   ├── store/                     # 状态管理（保持不变）
│   ├── directives/                # 自定义指令（保持不变）
│   ├── utils/                     # 工具函数（保持不变）
│   ├── views/                     # 页面视图（渐进式重构）
│   ├── App.vue                    # 根组件
│   └── main.ts                    # 入口文件
│
├── .storybook/                    # Storybook 配置 ⭐ 新增
│   ├── main.ts                    # 主配置
│   ├── preview.ts                 # 预览配置
│   └── manager.ts                 # 管理器配置
│
├── tailwind.config.js             # Tailwind 配置 ⭐ 新增
├── postcss.config.js              # PostCSS 配置 ⭐ 新增
├── package.json                   # 依赖配置（更新）
└── vite.config.ts                 # Vite 配置（更新）
```

### 2. 组件命名规范

**前缀规范：**
- `Gv` = GoyaVision（自定义组件）
- `El` = Element Plus（第三方组件）

**命名模式：**
```typescript
// ✅ 好的命名
GvButton.vue          // 基础按钮
GvCardHeader.vue      // 卡片头部
AssetCard.vue         // 资产卡片（业务组件）

// ❌ 避免的命名
Button.vue            // 太通用
MyButton.vue          // 无意义前缀
asset-card.vue        // 小写命名
```

---

## 组件体系

### 1. 基础组件清单（30+）

#### 数据展示类（8 个）

| 组件 | 描述 | 优先级 |
|------|------|--------|
| **GvCard** | 卡片容器 | P0 |
| **GvTable** | 表格（封装 ElTable） | P0 |
| **GvList** | 列表 | P1 |
| **GvDescriptions** | 描述列表 | P1 |
| **GvTimeline** | 时间轴 | P2 |
| **GvTree** | 树形控件 | P2 |
| **GvTag** | 标签 | P0 |
| **GvBadge** | 徽章 | P0 |

#### 表单类（10 个）

| 组件 | 描述 | 优先级 |
|------|------|--------|
| **GvButton** | 按钮 | P0 |
| **GvInput** | 输入框 | P0 |
| **GvTextarea** | 文本域 | P0 |
| **GvSelect** | 选择器 | P0 |
| **GvCheckbox** | 复选框 | P1 |
| **GvRadio** | 单选框 | P1 |
| **GvSwitch** | 开关 | P1 |
| **GvSlider** | 滑块 | P2 |
| **GvDatePicker** | 日期选择（封装 ElDatePicker） | P0 |
| **GvUpload** | 上传（封装 ElUpload） | P1 |

#### 反馈类（7 个）

| 组件 | 描述 | 优先级 |
|------|------|--------|
| **GvModal** | 模态框 | P0 |
| **GvDrawer** | 抽屉 | P0 |
| **GvAlert** | 警告提示 | P0 |
| **GvMessage** | 消息提示（函数式） | P0 |
| **GvLoading** | 加载状态 | P0 |
| **GvEmpty** | 空状态 | P1 |
| **GvSkeleton** | 骨架屏 | P2 |

#### 布局类（5 个）

| 组件 | 描述 | 优先级 |
|------|------|--------|
| **GvContainer** | 容器 | P0 |
| **GvGrid** | 栅格 | P0 |
| **GvFlex** | 弹性布局 | P0 |
| **GvDivider** | 分割线 | P1 |
| **GvSpace** | 间距 | P1 |

### 2. 业务组件清单（10+）

| 组件 | 描述 | 依赖基础组件 |
|------|------|-------------|
| **AssetCard** | 资产卡片 | GvCard, GvBadge, GvButton |
| **TaskCard** | 任务卡片 | GvCard, GvBadge, GvTag |
| **WorkflowNode** | 工作流节点 | GvCard, GvBadge |
| **UserAvatar** | 用户头像 | GvBadge |
| **StatusBadge** | 状态徽章 | GvBadge |
| **FilterBar** | 筛选栏 | GvSelect, GvInput, GvButton |
| **DataTable** | 数据表格 | GvTable, GvPagination |
| **FormBuilder** | 表单构建器 | GvInput, GvSelect, GvButton |
| **SearchBar** | 搜索栏 | GvInput, GvButton |
| **PageHeader** | 页面头部 | GvButton, GvBreadcrumb |

### 3. 组件 API 设计规范

#### 示例：GvButton 组件

```typescript
// components/base/GvButton/types.ts
export type ButtonType = 'primary' | 'secondary' | 'tertiary' | 'text' | 'outlined'
export type ButtonSize = 'small' | 'medium' | 'large'
export type ButtonVariant = 'filled' | 'tonal' | 'elevated' | 'outlined' | 'text'

export interface ButtonProps {
  // Material Design 3 变体
  variant?: ButtonVariant
  
  // 颜色（支持主题色）
  color?: 'primary' | 'secondary' | 'tertiary' | 'error' | 'success'
  
  // 尺寸
  size?: ButtonSize
  
  // 状态
  disabled?: boolean
  loading?: boolean
  
  // 图标
  icon?: string | Component
  iconPosition?: 'left' | 'right'
  
  // 样式
  rounded?: boolean
  block?: boolean
  
  // HTML 属性
  type?: 'button' | 'submit' | 'reset'
  href?: string
  target?: string
}
```

```vue
<!-- components/base/GvButton/index.vue -->
<template>
  <button
    :class="buttonClasses"
    :disabled="disabled || loading"
    :type="type"
    @click="handleClick"
  >
    <GvIcon v-if="loading" name="loading" class="animate-spin" />
    <GvIcon v-else-if="icon && iconPosition === 'left'" :name="icon" />
    <span v-if="$slots.default" class="gv-button__text">
      <slot />
    </span>
    <GvIcon v-if="icon && iconPosition === 'right'" :name="icon" />
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ButtonProps } from './types'

const props = withDefaults(defineProps<ButtonProps>(), {
  variant: 'filled',
  color: 'primary',
  size: 'medium',
  iconPosition: 'left',
  type: 'button'
})

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

// Material Design 3 风格类名
const buttonClasses = computed(() => [
  'gv-button',
  `gv-button--${props.variant}`,
  `gv-button--${props.color}`,
  `gv-button--${props.size}`,
  {
    'gv-button--rounded': props.rounded,
    'gv-button--block': props.block,
    'gv-button--disabled': props.disabled,
    'gv-button--loading': props.loading
  }
])

const handleClick = (event: MouseEvent) => {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
}
</script>

<style scoped>
/* Material Design 3 按钮样式 */
.gv-button {
  @apply relative inline-flex items-center justify-center;
  @apply font-medium transition-all duration-200;
  @apply focus:outline-none focus:ring-2 focus:ring-offset-2;
  @apply disabled:cursor-not-allowed disabled:opacity-50;
}

/* 尺寸 */
.gv-button--small {
  @apply h-8 px-3 text-sm gap-1.5 rounded-lg;
}

.gv-button--medium {
  @apply h-10 px-4 text-base gap-2 rounded-xl;
}

.gv-button--large {
  @apply h-12 px-6 text-lg gap-2.5 rounded-2xl;
}

/* Filled 变体 - 主要按钮 */
.gv-button--filled.gv-button--primary {
  @apply bg-primary-600 text-white;
  @apply hover:bg-primary-700 active:bg-primary-800;
  @apply focus:ring-primary-500;
}

/* Tonal 变体 - 次要按钮 */
.gv-button--tonal.gv-button--primary {
  @apply bg-primary-100 text-primary-700;
  @apply hover:bg-primary-200 active:bg-primary-300;
  @apply focus:ring-primary-500;
}

/* Outlined 变体 - 边框按钮 */
.gv-button--outlined.gv-button--primary {
  @apply border-2 border-primary-600 text-primary-600;
  @apply hover:bg-primary-50 active:bg-primary-100;
  @apply focus:ring-primary-500;
}

/* Text 变体 - 文本按钮 */
.gv-button--text.gv-button--primary {
  @apply text-primary-600;
  @apply hover:bg-primary-50 active:bg-primary-100;
  @apply focus:ring-primary-500;
}

/* 圆形按钮 */
.gv-button--rounded {
  @apply rounded-full;
}

/* 块级按钮 */
.gv-button--block {
  @apply w-full;
}
</style>
```

---

## 设计系统

### 1. Design Tokens（设计令牌）

#### 颜色系统

```typescript
// src/design-tokens/colors.ts

// Material Design 3 动态配色方案
export const colors = {
  // 主色（从品牌色 #667eea 生成）
  primary: {
    DEFAULT: '#667eea',
    50: '#f5f7ff',
    100: '#ebedff',
    200: '#d6dcff',
    300: '#b3bdff',
    400: '#8d9eff',
    500: '#667eea',      // 基准色
    600: '#5568d3',
    700: '#4553bd',
    800: '#3640a6',
    900: '#2a3290',
    950: '#1f2673'
  },
  
  // 辅助色（从 #764ba2 生成）
  secondary: {
    DEFAULT: '#764ba2',
    50: '#f9f5fc',
    100: '#f3ebf9',
    200: '#e7d7f3',
    300: '#d4b8e9',
    400: '#b88ed9',
    500: '#9d64c9',
    600: '#764ba2',     // 基准色
    700: '#65408b',
    800: '#543574',
    900: '#432a5d',
    950: '#321f46'
  },
  
  // 功能色
  success: {
    DEFAULT: '#10b981',
    50: '#f0fdf4',
    // ... 完整色阶
  },
  
  error: {
    DEFAULT: '#ef4444',
    50: '#fef2f2',
    // ... 完整色阶
  },
  
  warning: {
    DEFAULT: '#f59e0b',
    50: '#fffbeb',
    // ... 完整色阶
  },
  
  info: {
    DEFAULT: '#3b82f6',
    50: '#eff6ff',
    // ... 完整色阶
  },
  
  // 中性色（灰阶）
  neutral: {
    DEFAULT: '#64748b',
    50: '#f8fafc',
    100: '#f1f5f9',
    200: '#e2e8f0',
    300: '#cbd5e1',
    400: '#94a3b8',
    500: '#64748b',
    600: '#475569',
    700: '#334155',
    800: '#1e293b',
    900: '#0f172a',
    950: '#020617'
  },
  
  // 表面色（背景、卡片）
  surface: {
    DEFAULT: '#ffffff',
    dark: '#1e293b',
    dim: '#f8fafc',
    bright: '#ffffff',
    container: '#f1f5f9',
    'container-high': '#e2e8f0',
    'container-highest': '#cbd5e1'
  },
  
  // 文字色
  text: {
    primary: '#0f172a',
    secondary: '#475569',
    tertiary: '#64748b',
    disabled: '#cbd5e1',
    inverse: '#ffffff'
  }
}
```

#### 间距系统

```typescript
// src/design-tokens/spacing.ts

// 基于 8px 网格系统
export const spacing = {
  0: '0',
  0.5: '0.125rem',  // 2px
  1: '0.25rem',     // 4px
  1.5: '0.375rem',  // 6px
  2: '0.5rem',      // 8px   ← 基准单位
  3: '0.75rem',     // 12px
  4: '1rem',        // 16px
  5: '1.25rem',     // 20px
  6: '1.5rem',      // 24px
  8: '2rem',        // 32px
  10: '2.5rem',     // 40px
  12: '3rem',       // 48px
  16: '4rem',       // 64px
  20: '5rem',       // 80px
  24: '6rem',       // 96px
  32: '8rem',       // 128px
}

// 组件内边距规范
export const componentPadding = {
  xs: 'px-2 py-1',      // 8px × 4px
  sm: 'px-3 py-1.5',    // 12px × 6px
  md: 'px-4 py-2',      // 16px × 8px
  lg: 'px-6 py-3',      // 24px × 12px
  xl: 'px-8 py-4'       // 32px × 16px
}
```

#### 字体系统

```typescript
// src/design-tokens/typography.ts

// Material Design 3 字阶系统
export const typography = {
  // 字体家族
  fontFamily: {
    sans: [
      'Inter',
      '-apple-system',
      'BlinkMacSystemFont',
      'Segoe UI',
      'Roboto',
      'Helvetica Neue',
      'Arial',
      'sans-serif'
    ],
    mono: [
      'Fira Code',
      'Consolas',
      'Monaco',
      'Courier New',
      'monospace'
    ]
  },
  
  // 字号（9 档）
  fontSize: {
    xs: ['0.75rem', { lineHeight: '1rem' }],      // 12px / 16px
    sm: ['0.875rem', { lineHeight: '1.25rem' }],  // 14px / 20px
    base: ['1rem', { lineHeight: '1.5rem' }],     // 16px / 24px
    lg: ['1.125rem', { lineHeight: '1.75rem' }],  // 18px / 28px
    xl: ['1.25rem', { lineHeight: '1.75rem' }],   // 20px / 28px
    '2xl': ['1.5rem', { lineHeight: '2rem' }],    // 24px / 32px
    '3xl': ['1.875rem', { lineHeight: '2.25rem' }], // 30px / 36px
    '4xl': ['2.25rem', { lineHeight: '2.5rem' }],   // 36px / 40px
    '5xl': ['3rem', { lineHeight: '1' }]            // 48px
  },
  
  // 字重
  fontWeight: {
    light: '300',
    normal: '400',
    medium: '500',
    semibold: '600',
    bold: '700',
    extrabold: '800'
  }
}
```

#### 阴影系统

```typescript
// src/design-tokens/shadows.ts

// Material Design 3 阴影层级
export const shadows = {
  // 层级 1 - 微小提升（卡片）
  sm: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
  
  // 层级 2 - 小提升（悬停卡片）
  DEFAULT: '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1)',
  
  // 层级 3 - 中等提升（下拉菜单）
  md: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1)',
  
  // 层级 4 - 较大提升（模态框）
  lg: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1)',
  
  // 层级 5 - 最大提升（抽屉、浮动按钮）
  xl: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1)',
  
  // 特殊阴影
  '2xl': '0 25px 50px -12px rgba(0, 0, 0, 0.25)',
  inner: 'inset 0 2px 4px 0 rgba(0, 0, 0, 0.05)',
  
  // 彩色阴影（品牌色）
  primary: '0 8px 16px -4px rgba(102, 126, 234, 0.3)',
  secondary: '0 8px 16px -4px rgba(118, 75, 162, 0.3)',
  success: '0 8px 16px -4px rgba(16, 185, 129, 0.3)',
  error: '0 8px 16px -4px rgba(239, 68, 68, 0.3)'
}
```

#### 圆角系统

```typescript
// src/design-tokens/radius.ts

// Material Design 3 圆角规范
export const radius = {
  none: '0',
  sm: '0.25rem',    // 4px - 小元素
  DEFAULT: '0.5rem', // 8px - 基准圆角
  md: '0.75rem',    // 12px - 卡片、按钮
  lg: '1rem',       // 16px - 大卡片
  xl: '1.5rem',     // 24px - 模态框
  '2xl': '2rem',    // 32px - 大型容器
  '3xl': '3rem',    // 48px - 特大容器
  full: '9999px'    // 完全圆形
}
```

### 2. Tailwind 配置

```javascript
// tailwind.config.js
import { colors, spacing, typography, shadows, radius } from './src/design-tokens'

export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}'
  ],
  
  // 启用深色模式
  darkMode: 'class',
  
  theme: {
    extend: {
      // 注入设计令牌
      colors,
      spacing,
      fontSize: typography.fontSize,
      fontFamily: typography.fontFamily,
      fontWeight: typography.fontWeight,
      boxShadow: shadows,
      borderRadius: radius,
      
      // 动画曲线
      transitionTimingFunction: {
        'emphasized': 'cubic-bezier(0.2, 0, 0, 1)',
        'emphasized-decelerate': 'cubic-bezier(0.05, 0.7, 0.1, 1)',
        'emphasized-accelerate': 'cubic-bezier(0.3, 0, 0.8, 0.15)',
        'standard': 'cubic-bezier(0.2, 0, 0, 1)'
      },
      
      // 动画时长
      transitionDuration: {
        short1: '50ms',
        short2: '100ms',
        short3: '150ms',
        short4: '200ms',
        medium1: '250ms',
        medium2: '300ms',
        medium3: '350ms',
        medium4: '400ms',
        long1: '450ms',
        long2: '500ms',
        long3: '550ms',
        long4: '600ms',
        'extra-long1': '700ms',
        'extra-long2': '800ms',
        'extra-long3': '900ms',
        'extra-long4': '1000ms'
      }
    }
  },
  
  plugins: [
    // Tailwind 官方插件
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('@tailwindcss/container-queries'),
    
    // 自定义插件
    function({ addComponents }) {
      addComponents({
        // Material Design 3 表面容器
        '.surface': {
          '@apply bg-surface rounded-lg shadow-sm': {}
        },
        '.surface-container': {
          '@apply bg-surface-container rounded-lg': {}
        },
        
        // 文字省略
        '.text-ellipsis-1': {
          'overflow': 'hidden',
          'text-overflow': 'ellipsis',
          'white-space': 'nowrap'
        },
        '.text-ellipsis-2': {
          'display': '-webkit-box',
          '-webkit-line-clamp': '2',
          '-webkit-box-orient': 'vertical',
          'overflow': 'hidden'
        }
      })
    }
  ]
}
```

---

## 实施计划

### Phase 1: 基础设施搭建（Week 1-2）

#### 1.1 环境配置

**任务清单：**
- [ ] 安装 Tailwind CSS 及相关依赖
- [ ] 配置 PostCSS
- [ ] 创建设计令牌文件
- [ ] 配置 Tailwind 主题
- [ ] 安装 Storybook
- [ ] 配置 Storybook + Tailwind 集成

**依赖安装：**
```bash
# Tailwind CSS 生态
pnpm add -D tailwindcss postcss autoprefixer
pnpm add -D @tailwindcss/forms @tailwindcss/typography @tailwindcss/container-queries

# Storybook
pnpm dlx storybook@latest init

# 工具库
pnpm add clsx tailwind-merge
pnpm add @vueuse/core  # Vue 组合式 API 工具库
```

**验收标准：**
- ✅ Tailwind 类名生效
- ✅ 设计令牌可在组件中使用
- ✅ Storybook 本地运行成功

#### 1.2 项目结构调整

**任务清单：**
- [ ] 创建 `src/components/base/` 目录
- [ ] 创建 `src/components/layout/` 目录
- [ ] 创建 `src/components/business/` 目录
- [ ] 创建 `src/design-tokens/` 目录
- [ ] 创建 `src/composables/` 目录
- [ ] 创建 `src/styles/` 目录

### Phase 2: 基础组件库（Week 3-6）

#### 2.1 第一批组件（P0 优先级）

**Week 3:**
- [ ] GvButton（5 种变体 × 5 种颜色 × 3 种尺寸）
- [ ] GvCard（基础版）
- [ ] GvBadge
- [ ] GvTag

**Week 4:**
- [ ] GvInput
- [ ] GvSelect（封装 ElSelect）
- [ ] GvModal
- [ ] GvAlert

**Week 5:**
- [ ] GvTable（封装 ElTable）
- [ ] GvContainer
- [ ] GvGrid
- [ ] GvFlex

**Week 6:**
- [ ] GvLoading
- [ ] GvEmpty
- [ ] GvDrawer
- [ ] GvDatePicker（封装 ElDatePicker）

**每个组件交付物：**
- ✅ 组件实现（`index.vue`）
- ✅ TypeScript 类型（`types.ts`）
- ✅ Storybook 文档（`stories.ts`）
- ✅ 使用示例（`README.md`）
- ✅ 单元测试（`__tests__/`）

#### 2.2 组件统一导出

```typescript
// src/components/index.ts
export { default as GvButton } from './base/GvButton/index.vue'
export { default as GvCard } from './base/GvCard/index.vue'
export { default as GvBadge } from './base/GvBadge/index.vue'
// ... 其他组件

export type * from './base/GvButton/types'
export type * from './base/GvCard/types'
// ... 其他类型
```

### Phase 3: 业务组件库（Week 7-8）

**Week 7:**
- [ ] AssetCard（资产卡片）
- [ ] TaskCard（任务卡片）
- [ ] StatusBadge（状态徽章）
- [ ] UserAvatar（用户头像）

**Week 8:**
- [ ] FilterBar（筛选栏）
- [ ] DataTable（数据表格）
- [ ] SearchBar（搜索栏）
- [ ] PageHeader（页面头部）

### Phase 4: 页面重构（Week 9-12）

**渐进式重构策略：**

**Week 9: 登录页**
- [ ] 使用新组件重构登录页
- [ ] 应用 Material Design 3 风格
- [ ] 添加深色模式支持

**Week 10: 资产管理页**
- [ ] 使用 AssetCard 组件
- [ ] 使用 FilterBar 组件
- [ ] 使用 DataTable 组件

**Week 11: 任务管理页**
- [ ] 使用 TaskCard 组件
- [ ] 使用新的表单组件

**Week 12: 系统管理页**
- [ ] 用户管理页重构
- [ ] 角色管理页重构
- [ ] 菜单管理页重构

### Phase 5: 优化与文档（Week 13-14）

**Week 13:**
- [ ] 性能优化（Tree Shaking、代码分割）
- [ ] 深色模式完善
- [ ] 响应式优化
- [ ] 可访问性测试

**Week 14:**
- [ ] 组件文档补全
- [ ] Storybook 部署
- [ ] AI 调用规范文档
- [ ] 开发指南文档

---

## 规范文档

### 1. 组件开发规范

#### 文件结构规范

```
GvButton/
├── index.vue              # 组件实现（必需）
├── types.ts               # TypeScript 类型（必需）
├── README.md              # 使用文档（必需）
├── stories.ts             # Storybook 故事（必需）
├── __tests__/             # 单元测试（推荐）
│   └── index.spec.ts
└── styles.css             # 独立样式（可选，优先使用 Tailwind）
```

#### 组件模板规范

```vue
<template>
  <!-- 使用 Tailwind 类名 -->
  <div :class="containerClasses">
    <slot name="prefix" />
    <slot />
    <slot name="suffix" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ComponentProps } from './types'

// Props 定义（使用 TypeScript）
const props = withDefaults(defineProps<ComponentProps>(), {
  // 默认值
})

// Emits 定义
const emit = defineEmits<{
  change: [value: string]
  click: [event: MouseEvent]
}>()

// 计算属性（类名组合）
const containerClasses = computed(() => [
  'gv-component-name',
  props.variant && `gv-component-name--${props.variant}`,
  {
    'gv-component-name--active': props.active
  }
])
</script>

<style scoped>
/* 仅在 Tailwind 无法满足时使用 */
.gv-component-name {
  /* 自定义样式 */
}
</style>
```

### 2. Storybook 故事规范

```typescript
// GvButton.stories.ts
import type { Meta, StoryObj } from '@storybook/vue3'
import GvButton from './index.vue'

const meta: Meta<typeof GvButton> = {
  title: 'Base/GvButton',
  component: GvButton,
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['filled', 'tonal', 'outlined', 'text'],
      description: 'Material Design 3 按钮变体'
    },
    color: {
      control: 'select',
      options: ['primary', 'secondary', 'tertiary', 'error', 'success'],
      description: '按钮颜色'
    },
    size: {
      control: 'select',
      options: ['small', 'medium', 'large'],
      description: '按钮尺寸'
    }
  }
}

export default meta
type Story = StoryObj<typeof GvButton>

// 默认故事
export const Default: Story = {
  args: {
    variant: 'filled',
    color: 'primary',
    size: 'medium'
  },
  render: (args) => ({
    components: { GvButton },
    setup() {
      return { args }
    },
    template: '<GvButton v-bind="args">点击按钮</GvButton>'
  })
}

// 所有变体展示
export const AllVariants: Story = {
  render: () => ({
    components: { GvButton },
    template: `
      <div class="flex gap-4">
        <GvButton variant="filled">Filled</GvButton>
        <GvButton variant="tonal">Tonal</GvButton>
        <GvButton variant="outlined">Outlined</GvButton>
        <GvButton variant="text">Text</GvButton>
      </div>
    `
  })
}

// 所有尺寸展示
export const AllSizes: Story = {
  render: () => ({
    components: { GvButton },
    template: `
      <div class="flex items-center gap-4">
        <GvButton size="small">Small</GvButton>
        <GvButton size="medium">Medium</GvButton>
        <GvButton size="large">Large</GvButton>
      </div>
    `
  })
}

// 加载状态
export const Loading: Story = {
  args: {
    loading: true
  },
  render: (args) => ({
    components: { GvButton },
    setup() {
      return { args }
    },
    template: '<GvButton v-bind="args">加载中...</GvButton>'
  })
}
```

### 3. AI 调用规范文档

创建专门的 AI 规范文档，让 AI 能够准确理解和使用组件。

```markdown
<!-- .cursor/rules/component-usage.mdc -->
# GoyaVision 组件使用规范

## 基础组件

### GvButton - 按钮组件

**路径**: `@/components/base/GvButton`

**基本用法**:
```vue
<GvButton variant="filled" color="primary">
  点击按钮
</GvButton>
```

**Props**:
- `variant`: 'filled' | 'tonal' | 'outlined' | 'text' - 按钮变体
- `color`: 'primary' | 'secondary' | 'tertiary' | 'error' | 'success' - 颜色
- `size`: 'small' | 'medium' | 'large' - 尺寸
- `loading`: boolean - 加载状态
- `disabled`: boolean - 禁用状态
- `icon`: string - 图标名称
- `iconPosition`: 'left' | 'right' - 图标位置

**Events**:
- `@click`: (event: MouseEvent) => void - 点击事件

**示例**:
```vue
<!-- 主要按钮 -->
<GvButton variant="filled" color="primary" @click="handleSubmit">
  提交
</GvButton>

<!-- 次要按钮 -->
<GvButton variant="tonal" color="secondary">
  取消
</GvButton>

<!-- 加载按钮 -->
<GvButton variant="filled" :loading="isLoading">
  保存中...
</GvButton>

<!-- 图标按钮 -->
<GvButton variant="outlined" icon="search" icon-position="left">
  搜索
</GvButton>
```

**使用场景**:
- 表单提交: variant="filled" color="primary"
- 取消操作: variant="tonal" color="secondary"
- 删除操作: variant="outlined" color="error"
- 次要操作: variant="text"

---

### GvCard - 卡片组件

**路径**: `@/components/base/GvCard`

**基本用法**:
```vue
<GvCard>
  <template #header>
    <h3>卡片标题</h3>
  </template>
  <p>卡片内容</p>
  <template #footer>
    <GvButton>操作</GvButton>
  </template>
</GvCard>
```

**Props**:
- `shadow`: 'none' | 'sm' | 'md' | 'lg' | 'xl' - 阴影大小
- `padding`: 'none' | 'sm' | 'md' | 'lg' - 内边距
- `hoverable`: boolean - 是否支持悬停效果

**Slots**:
- `header` - 卡片头部
- `default` - 卡片内容
- `footer` - 卡片底部

**示例**:
```vue
<!-- 基础卡片 -->
<GvCard shadow="md" padding="lg">
  <h3 class="text-lg font-semibold mb-2">资产名称</h3>
  <p class="text-sm text-text-secondary">资产描述信息</p>
</GvCard>

<!-- 可悬停卡片 -->
<GvCard hoverable @click="handleCardClick">
  <AssetInfo :asset="asset" />
</GvCard>
```

---

## 业务组件

### AssetCard - 资产卡片

**路径**: `@/components/business/AssetCard`

**基本用法**:
```vue
<AssetCard :asset="assetData" @click="handleAssetClick" />
```

**Props**:
- `asset`: AssetData - 资产数据对象
- `selectable`: boolean - 是否可选择
- `selected`: boolean - 是否已选择

**Events**:
- `@click`: (asset: AssetData) => void - 点击事件
- `@select`: (selected: boolean) => void - 选择事件

**示例**:
```vue
<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
    <AssetCard
      v-for="asset in assets"
      :key="asset.id"
      :asset="asset"
      @click="handleAssetClick"
    />
  </div>
</template>
```

---

## 布局组件

### GvContainer - 容器组件

**路径**: `@/components/layout/GvContainer`

**基本用法**:
```vue
<GvContainer max-width="xl">
  <p>内容自动居中，最大宽度 1280px</p>
</GvContainer>
```

**Props**:
- `maxWidth`: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | 'full' - 最大宽度
- `padding`: boolean - 是否添加内边距

---

## 使用建议

1. **优先使用组件库**: 不要重复造轮子，先检查组件库是否有现成组件
2. **遵循 Material Design 3**: 使用推荐的变体和配色方案
3. **保持一致性**: 相同功能使用相同组件和样式
4. **响应式设计**: 使用 Tailwind 的响应式类名（sm:、md:、lg:）
5. **可访问性**: 为交互元素添加适当的 aria 属性

## 常见模式

### 表单布局
```vue
<GvCard>
  <template #header>
    <h2>表单标题</h2>
  </template>
  
  <div class="space-y-4">
    <GvInput label="姓名" v-model="form.name" />
    <GvSelect label="类型" v-model="form.type" :options="typeOptions" />
  </div>
  
  <template #footer>
    <div class="flex justify-end gap-2">
      <GvButton variant="tonal" @click="handleCancel">取消</GvButton>
      <GvButton variant="filled" @click="handleSubmit">提交</GvButton>
    </div>
  </template>
</GvCard>
```

### 列表页布局
```vue
<GvContainer>
  <PageHeader title="资产管理">
    <template #actions>
      <GvButton variant="filled" icon="add">新建资产</GvButton>
    </template>
  </PageHeader>
  
  <FilterBar v-model="filters" />
  
  <div class="grid grid-cols-3 gap-4 mt-6">
    <AssetCard v-for="asset in assets" :key="asset.id" :asset="asset" />
  </div>
  
  <GvPagination v-model="page" :total="total" />
</GvContainer>
```
```

---

## 风险评估

### 技术风险

| 风险 | 影响 | 概率 | 应对措施 |
|------|------|------|----------|
| Tailwind 学习曲线 | 中 | 高 | 提供培训文档、代码示例 |
| Element Plus 兼容性 | 中 | 中 | 封装适配层，隔离第三方依赖 |
| 组件库维护成本 | 高 | 中 | 建立完善文档、自动化测试 |
| 性能影响 | 低 | 低 | Tree Shaking、按需加载 |

### 项目风险

| 风险 | 影响 | 概率 | 应对措施 |
|------|------|------|----------|
| 重构周期过长 | 高 | 中 | 渐进式重构、保持旧页面可用 |
| 团队抵触情绪 | 中 | 中 | 充分沟通、展示收益 |
| 设计不统一 | 中 | 低 | 严格遵循 Material Design 3 |
| 文档不完善 | 高 | 中 | 每个组件必须配文档 |

---

## 成功指标

### 量化指标

1. **代码质量**
   - 组件复用率 > 70%
   - TypeScript 类型覆盖率 100%
   - 单元测试覆盖率 > 80%

2. **性能指标**
   - CSS 文件体积减少 30%
   - 首屏渲染时间减少 20%
   - Lighthouse 性能分数 > 90

3. **开发效率**
   - 新页面开发时间减少 50%
   - 样式调整时间减少 70%
   - Bug 修复时间减少 40%

4. **文档完善度**
   - 组件文档覆盖率 100%
   - Storybook 故事数 > 100
   - AI 调用准确率 > 95%

---

## 总结

本重构方案遵循以下原则：

1. ✅ **渐进式重构** - 不影响现有功能，逐步迁移
2. ✅ **标准化架构** - 引入 Tailwind CSS + Material Design 3
3. ✅ **组件化思维** - 封装 30+ 基础组件 + 10+ 业务组件
4. ✅ **规范化文档** - 完善的组件文档和 AI 调用规范
5. ✅ **工程化实践** - Storybook、单元测试、类型安全

**预期收益：**
- 开发效率提升 50%
- 代码可维护性提升 80%
- UI 一致性提升 95%
- AI 辅助开发准确率 > 95%

---

**维护者**: GoyaVision Team  
**最后更新**: 2026-02-03  
**版本**: V1.0  
**状态**: 📋 方案制定完成，待评审
