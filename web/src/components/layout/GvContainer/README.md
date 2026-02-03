# GvContainer - 容器组件

响应式容器组件，自动居中并限制最大宽度。

## 基本用法

```vue
<GvContainer>
  <p>内容自动居中，最大宽度 1280px</p>
</GvContainer>
```

## 最大宽度

```vue
<GvContainer max-width="sm">640px</GvContainer>
<GvContainer max-width="md">768px</GvContainer>
<GvContainer max-width="lg">1024px</GvContainer>
<GvContainer max-width="xl">1280px</GvContainer>
<GvContainer max-width="2xl">1536px</GvContainer>
<GvContainer max-width="full">100%</GvContainer>
```

## 内边距

```vue
<!-- 带内边距（默认） -->
<GvContainer padding>
  <p>有内边距的容器</p>
</GvContainer>

<!-- 无内边距 -->
<GvContainer :padding="false">
  <p>无内边距的容器</p>
</GvContainer>
```

## 对齐方式

```vue
<!-- 居中对齐（默认） -->
<GvContainer centered>
  <p>居中对齐</p>
</GvContainer>

<!-- 不居中 -->
<GvContainer :centered="false">
  <p>不居中</p>
</GvContainer>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| maxWidth | `'sm' \| 'md' \| 'lg' \| 'xl' \| '2xl' \| 'full'` | `'xl'` | 最大宽度 |
| padding | `boolean` | `true` | 是否添加水平内边距 |
| centered | `boolean` | `true` | 是否居中对齐 |

## 使用场景

### 页面容器

```vue
<template>
  <GvContainer>
    <PageHeader />
    <PageContent />
  </GvContainer>
</template>
```

### 窄容器（表单页面）

```vue
<GvContainer max-width="md">
  <GvCard>
    <el-form>
      <!-- 表单内容 -->
    </el-form>
  </GvCard>
</GvContainer>
```

### 全宽容器

```vue
<GvContainer max-width="full" :padding="false">
  <div class="hero-section">
    <!-- 全宽内容 -->
  </div>
</GvContainer>
```
