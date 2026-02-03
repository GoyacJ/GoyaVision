# 前端重构 - 快速开始指南

> 如何启动和测试重构后的前端项目

---

## 🚀 快速开始

### 1. 安装依赖

```bash
cd web
pnpm install
```

**如果没有 pnpm**，先安装：
```bash
npm install -g pnpm
```

或使用 npm：
```bash
npm install
```

### 2. 启动开发服务器

```bash
pnpm dev
```

服务器将在 `http://localhost:5173` 启动。

### 3. 访问组件展示页面

打开浏览器访问：
```
http://localhost:5173/component-demo
```

这个页面展示了所有已完成的组件。

### 4. 访问主应用

```
http://localhost:5173/login
```

使用默认账号登录：
- 用户名: `admin`
- 密码: `admin123`

---

## 📦 新增依赖说明

### 核心依赖

| 包名 | 版本 | 用途 |
|------|------|------|
| **tailwindcss** | ^3.4.0 | Tailwind CSS 框架 |
| **postcss** | ^8.4.33 | CSS 后处理器 |
| **autoprefixer** | ^10.4.16 | 自动添加浏览器前缀 |
| **clsx** | ^2.1.0 | 类名合并工具 |
| **tailwind-merge** | ^2.2.0 | Tailwind 类名冲突解决 |
| **@vueuse/core** | ^10.7.0 | Vue 组合式 API 工具库 |

### Tailwind 插件

| 包名 | 用途 |
|------|------|
| **@tailwindcss/forms** | 表单样式优化 |
| **@tailwindcss/typography** | 排版样式 |
| **@tailwindcss/container-queries** | 容器查询支持 |

### Storybook（开发依赖）

| 包名 | 版本 | 用途 |
|------|------|------|
| **storybook** | ^7.6.0 | Storybook 核心 |
| **@storybook/vue3** | ^7.6.0 | Vue 3 集成 |
| **@storybook/vue3-vite** | ^7.6.0 | Vite 集成 |
| **@storybook/addon-essentials** | ^7.6.0 | 基础插件 |

---

## 🎨 已完成的组件

### 基础组件（4 个）

1. **GvButton** - 按钮组件
   - 4 种变体 × 6 种颜色 × 3 种尺寸
   - 支持图标、加载状态、圆形/块级按钮

2. **GvCard** - 卡片组件
   - 5 种阴影 × 4 种内边距
   - header/footer 插槽
   - 悬停效果、边框模式

3. **GvBadge** - 徽章组件
   - 7 种颜色 × 3 种变体 × 3 种尺寸
   - 独立徽章 + 角标徽章
   - 数字显示、点状徽章

4. **GvTag** - 标签组件
   - 7 种颜色 × 3 种变体 × 3 种尺寸
   - 图标、可关闭、圆形标签

### 布局组件（1 个）

1. **GvContainer** - 容器组件
   - 6 种最大宽度
   - 响应式内边距
   - 居中对齐

---

## 💻 组件使用示例

### 基本用法

```vue
<template>
  <GvContainer>
    <GvCard shadow="md" padding="lg">
      <template #header>
        <div class="flex justify-between items-center">
          <h3 class="text-lg font-semibold">卡片标题</h3>
          <GvBadge color="success">就绪</GvBadge>
        </div>
      </template>
      
      <p class="text-text-secondary mb-4">卡片内容区域</p>
      
      <div class="flex flex-wrap gap-2">
        <GvTag icon="VideoCamera" color="primary">视频</GvTag>
        <GvTag icon="Check" color="success">已处理</GvTag>
      </div>
      
      <template #footer>
        <div class="flex justify-end gap-2">
          <GvButton variant="tonal">取消</GvButton>
          <GvButton variant="filled">确定</GvButton>
        </div>
      </template>
    </GvCard>
  </GvContainer>
</template>

<script setup lang="ts">
import { GvContainer, GvCard, GvButton, GvBadge, GvTag } from '@/components'
</script>
```

### 表单场景

```vue
<template>
  <GvContainer max-width="md">
    <GvCard>
      <template #header>
        <h2 class="text-xl font-semibold">新建资产</h2>
      </template>
      
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.type" class="w-full">
            <el-option label="视频" value="video" />
            <el-option label="图片" value="image" />
          </el-select>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="flex justify-end gap-2">
          <GvButton variant="tonal" @click="handleCancel">
            取消
          </GvButton>
          <GvButton variant="filled" :loading="loading" @click="handleSubmit">
            提交
          </GvButton>
        </div>
      </template>
    </GvCard>
  </GvContainer>
</template>
```

### 列表页场景

```vue
<template>
  <GvContainer>
    <!-- 页面头部 -->
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">资产管理</h1>
      <GvButton variant="filled" icon="Plus">新建资产</GvButton>
    </div>
    
    <!-- 卡片网格 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <GvCard
        v-for="asset in assets"
        :key="asset.id"
        hoverable
        @click="handleAssetClick(asset)"
      >
        <template #header>
          <div class="flex justify-between items-center">
            <h3 class="font-semibold">{{ asset.name }}</h3>
            <GvBadge :color="getStatusColor(asset.status)">
              {{ asset.status }}
            </GvBadge>
          </div>
        </template>
        
        <p class="text-sm text-text-secondary">{{ asset.path }}</p>
        
        <div class="mt-3 flex flex-wrap gap-2">
          <GvTag size="small" :color="getTypeColor(asset.type)">
            {{ asset.type }}
          </GvTag>
        </div>
      </GvCard>
    </div>
  </GvContainer>
</template>
```

---

## 🎨 Tailwind CSS 使用

### 常用工具类

```vue
<!-- 布局 -->
<div class="flex items-center justify-between gap-4">
<div class="grid grid-cols-3 gap-6">

<!-- 间距 -->
<div class="p-4 m-2">              <!-- padding: 16px, margin: 8px -->
<div class="px-6 py-4">            <!-- padding: 24px 16px -->
<div class="space-y-4">            <!-- 子元素垂直间距 16px -->

<!-- 文字 -->
<h1 class="text-2xl font-bold text-text-primary">
<p class="text-sm text-text-secondary">

<!-- 颜色 -->
<div class="bg-primary-600 text-white">
<div class="bg-surface shadow-md rounded-lg">

<!-- 响应式 -->
<div class="w-full md:w-1/2 lg:w-1/3">
<div class="hidden md:block">        <!-- 中屏及以上显示 -->
```

### 设计令牌颜色

```vue
<!-- 主色调 -->
<div class="text-primary-600">      <!-- #667eea -->
<div class="bg-primary-100">        <!-- 浅色背景 -->
<div class="border-primary-600">    <!-- 边框 -->

<!-- 功能色 -->
<div class="text-success-600">      <!-- 成功色 -->
<div class="text-error-600">        <!-- 错误色 -->
<div class="text-warning-600">      <!-- 警告色 -->

<!-- 文字色 -->
<p class="text-text-primary">       <!-- 主要文字 -->
<p class="text-text-secondary">     <!-- 次要文字 -->
<p class="text-text-tertiary">      <!-- 第三级文字 -->
```

---

## 🧪 测试指南

### 测试组件展示页面

1. 启动开发服务器
2. 访问 `http://localhost:5173/component-demo`
3. 测试以下功能：
   - ✅ 所有按钮变体和颜色
   - ✅ 所有按钮尺寸
   - ✅ 按钮图标和加载状态
   - ✅ 卡片悬停效果
   - ✅ 徽章显示
   - ✅ 标签关闭功能
   - ✅ 主题切换（深色/浅色模式）

### 测试主应用

1. 访问 `http://localhost:5173/login`
2. 登录后查看各页面
3. 确认新组件与现有页面兼容

---

## 🔧 开发指南

### 创建新组件

1. **创建组件目录**
   ```bash
   mkdir -p src/components/base/GvYourComponent
   cd src/components/base/GvYourComponent
   ```

2. **创建必需文件**
   ```bash
   touch index.vue types.ts README.md
   ```

3. **参考现有组件**
   - 查看 `GvButton` 作为模板
   - 遵循 Material Design 3 规范
   - 使用 Tailwind CSS 类名

4. **添加到导出文件**
   ```typescript
   // src/components/index.ts
   export { default as GvYourComponent } from './base/GvYourComponent/index.vue'
   export type * from './base/GvYourComponent/types'
   ```

### 使用组件

1. **导入组件**
   ```typescript
   import { GvButton, GvCard } from '@/components'
   ```

2. **使用组件**
   ```vue
   <template>
     <GvButton variant="filled" color="primary">
       点击按钮
     </GvButton>
   </template>
   ```

---

## 📚 相关文档

| 文档 | 路径 | 用途 |
|------|------|------|
| **快速开始** | `web/REFACTOR-GUIDE.md` | 本文档 |
| **进度追踪** | `docs/REFACTOR-PROGRESS.md` | 实时进度 |
| **详细方案** | `docs/frontend-refactor-plan.md` | 完整方案 |
| **组件规范** | `.cursor/rules/frontend-components.mdc` | AI 开发指南 |
| **UI 设计** | `docs/ui-design.md` | 设计系统 |

---

## ⚠️ 注意事项

### 1. 依赖安装

确保网络畅通，某些依赖可能需要从 npm registry 下载。

如果安装失败，可以尝试：
```bash
# 使用淘宝镜像
pnpm install --registry=https://registry.npmmirror.com

# 或清除缓存后重试
pnpm store prune
pnpm install
```

### 2. Node.js 版本

推荐使用 Node.js 18+ 版本：
```bash
node -v  # 应该 >= 18.0.0
```

### 3. 浏览器兼容性

推荐使用现代浏览器：
- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

### 4. 开发工具

推荐安装以下 VS Code 扩展：
- Volar (Vue 3 支持)
- Tailwind CSS IntelliSense
- PostCSS Language Support
- ESLint
- Prettier

---

## 🎯 下一步

### 继续开发组件

查看 [REFACTOR-PROGRESS.md](../docs/REFACTOR-PROGRESS.md) 了解：
- 下一步要开发的组件
- 当前进度和里程碑
- 详细的实施计划

### 参与贡献

1. 查看组件开发规范
2. 选择一个待开发组件
3. 按照模板开发
4. 提交 Pull Request

---

## 💡 常见问题

### Q: Tailwind CSS 类名不生效？

**A:** 检查以下几点：
1. 确认已在 `main.ts` 中导入 `./styles/tailwind.css`
2. 确认 `tailwind.config.js` 的 `content` 配置正确
3. 重启开发服务器

### Q: 组件导入报错？

**A:** 确认：
1. 组件已在 `src/components/index.ts` 中导出
2. 使用正确的导入路径：`import { GvButton } from '@/components'`
3. TypeScript 配置正确（`tsconfig.json` 中的 `paths`）

### Q: 如何切换深色模式？

**A:** 
```typescript
import { useTheme } from '@/composables'

const { toggleTheme } = useTheme()

// 切换主题
toggleTheme()
```

或在组件展示页面点击"切换主题"按钮。

### Q: 如何查看所有组件文档？

**A:** 
- 访问组件展示页面：`http://localhost:5173/component-demo`
- 查看各组件目录下的 `README.md` 文件
- 将来可以启动 Storybook：`pnpm run storybook`

---

## 📞 获取帮助

- 查看组件 README.md 文档
- 查看 `.cursor/rules/frontend-components.mdc`（AI 调用规范）
- 查看 `docs/frontend-refactor-plan.md`（完整方案）
- 在组件展示页面查看实际效果

---

**祝开发顺利！🎊**
