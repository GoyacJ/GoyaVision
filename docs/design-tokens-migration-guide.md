# Design Tokens 迁移指南

> 从旧设计系统迁移到新克制设计系统的完整指南

## 概览

本次 Design Tokens 重构的核心目标：
- ✅ 从装饰性设计转向内容优先
- ✅ 降低视觉噪音，提升可读性
- ✅ 移除过度的装饰效果（渐变、彩色阴影、毛玻璃）
- ✅ 采用更专业、更克制的配色方案

---

## 色彩系统变更

### 1. 主色（Primary）

**变更前**：
```css
/* 鲜艳的蓝紫色 */
primary: #667eea
```

**变更后**：
```css
/* 克制的蓝灰色（降低饱和度） */
primary: #4F5B93
```

**迁移策略**：
- 全局搜索 `bg-primary-500` 或 `text-primary-500`
- 检查视觉效果，确保对比度足够
- 按钮、链接、标签等主色元素自动生效

---

### 2. 辅助色（Secondary）- 已弃用

**变更前**：
```css
secondary: #764ba2 /* 紫色，用于渐变 */
```

**变更后**：
```css
/* 已完全移除 */
```

**迁移策略**：

**❌ 不再使用**：
```vue
<!-- 旧代码：渐变背景 -->
<div class="bg-gradient-to-r from-primary-500 to-secondary-600">
```

**✅ 替换为**：
```vue
<!-- 新代码：单色背景 -->
<div class="bg-primary-600">
```

**查找并替换**：
- `from-secondary-*` → 删除
- `to-secondary-*` → 删除
- `bg-gradient-to-*` → 改为 `bg-primary-*`
- `text-secondary-*` → 改为 `text-primary-*` 或 `text-text-secondary`

---

### 3. 中性色（Neutral）

**变更前**：
```css
/* 蓝灰色系（有色相偏向） */
neutral-500: #64748b
```

**变更后**：
```css
/* 纯灰色系（无色相偏向） */
neutral-500: #737373
```

**色彩对照表**：

| 用途 | 旧值 | 新值 |
|------|------|------|
| 页面背景 | `#f8fafc` | `#FAFAFA` |
| 容器背景 | `#f1f5f9` | `#F5F5F5` |
| 边框颜色 | `#e2e8f0` | `#E5E5E5` |
| 主要文本 | `#0f172a` | `#262626` |
| 次要文本 | `#475569` | `#525252` |
| 三级文本 | `#64748b` | `#737373` |

**迁移策略**：
- 旧代码使用的 `neutral-*` 类名保持不变
- 颜色值会自动更新为新的纯灰色系
- 检查深色背景上的文本对比度

---

### 4. 文字色（Text）

**变更前**：
```css
text-primary: #0f172a   /* 深蓝灰 */
text-secondary: #475569 /* 蓝灰 */
text-tertiary: #64748b  /* 浅蓝灰 */
```

**变更后**：
```css
text-primary: #262626   /* 纯黑灰 */
text-secondary: #525252 /* 深灰 */
text-tertiary: #737373  /* 中灰 */
text-placeholder: #A3A3A3 /* 浅灰 */
```

**新增语义化类名**：
```vue
<p class="text-text-primary">主要文本</p>
<p class="text-text-secondary">次要文本</p>
<p class="text-text-tertiary">三级文本</p>
<input placeholder="占位符" class="placeholder:text-text-placeholder" />
```

---

## 阴影系统变更

### 移除彩色阴影

**变更前**：
```css
shadow-primary: 0 8px 16px -4px rgba(102, 126, 234, 0.3)
shadow-secondary: 0 8px 16px -4px rgba(118, 75, 162, 0.3)
```

**变更后**：
```css
/* 已完全移除彩色阴影 */
```

**迁移策略**：

**❌ 删除所有彩色阴影**：
```vue
<!-- 旧代码 -->
<button class="shadow-primary">按钮</button>
<div class="shadow-secondary">卡片</div>
```

**✅ 替换为标准阴影**：
```vue
<!-- 新代码 -->
<button class="shadow-md">按钮</button>
<div class="shadow-sm">卡片</div>
```

---

### 降低阴影透明度

**变更对照表**：

| 层级 | 旧值 | 新值 |
|------|------|------|
| `shadow-sm` | `rgba(0,0,0,0.05)` | `rgba(0,0,0,0.04)` |
| `shadow` | `rgba(0,0,0,0.1)` | `rgba(0,0,0,0.06)` |
| `shadow-md` | `rgba(0,0,0,0.1)` | `rgba(0,0,0,0.07)` |
| `shadow-lg` | `rgba(0,0,0,0.1)` | `rgba(0,0,0,0.08)` |
| `shadow-xl` | `rgba(0,0,0,0.1)` | `rgba(0,0,0,0.10)` |

**迁移策略**：
- 现有 `shadow-*` 类名保持不变
- 视觉效果会自动变得更轻、更克制
- 无需手动修改代码

---

## 圆角系统变更

### 适度减小圆角尺寸

**变更对照表**：

| 类名 | 旧值 | 新值 |
|------|------|------|
| `rounded` | `8px` | `6px` |
| `rounded-md` | `12px` | `8px` |
| `rounded-lg` | `16px` | `12px` |
| `rounded-xl` | `24px` | `16px` |
| `rounded-2xl` | `32px` | `24px` |

**迁移策略**：
- 现有类名保持不变
- 视觉效果会自动变得更克制
- 检查是否有过度圆润的效果

**移除 `rounded-3xl`**：
```vue
<!-- 旧代码 -->
<div class="rounded-3xl"> <!-- 48px，过大 -->

<!-- 新代码 -->
<div class="rounded-2xl"> <!-- 24px，合理 -->
```

---

## 排版系统变更

### 1. 正文字号调整

**变更前**：
```css
text-base: 16px / 24px
```

**变更后**：
```css
text-base: 15px / 24px  /* 参考 Medium */
```

**迁移策略**：
- 无需修改代码，自动生效
- 检查密集文本区域的排版效果
- 确保可读性没有下降

---

### 2. 字距（Letter Spacing）

**新增配置**：
```css
tracking-tighter: -0.02em  /* 超紧凑，用于大标题 */
tracking-tight: -0.01em    /* 紧凑，用于标题 */
tracking-normal: 0         /* 标准，用于正文 */
tracking-wide: 0.025em     /* 宽松，用于标签 */
tracking-wider: 0.05em     /* 超宽松，用于全大写 */
```

**推荐使用**：
```vue
<!-- 页面标题：紧凑字距 -->
<h1 class="text-4xl font-bold tracking-tighter">页面标题</h1>

<!-- 区块标题 -->
<h2 class="text-2xl font-semibold tracking-tight">区块标题</h2>

<!-- 正文 -->
<p class="text-base tracking-normal">正文内容</p>

<!-- 全大写标签 -->
<span class="text-xs uppercase tracking-wider">STATUS</span>
```

---

### 3. 排版工具类

**新增语义化类名**：

```vue
<!-- 标题 -->
<h1 class="text-h1">页面标题</h1>       <!-- 32px, bold, -0.02em -->
<h2 class="text-h2">区块标题</h2>       <!-- 24px, semibold, -0.01em -->
<h3 class="text-h3">卡片标题</h3>       <!-- 18px, semibold -->

<!-- 正文 -->
<p class="text-body">正文内容</p>        <!-- 15px, normal -->
<p class="text-body-small">小正文</p>    <!-- 13px, normal -->
<p class="text-caption">次要文本</p>     <!-- 13px, tertiary -->

<!-- 标签 -->
<span class="text-label">按钮文本</span> <!-- 13px, medium -->
<span class="text-overline">状态</span>  <!-- 12px, uppercase, wider -->
```

---

## 动画系统变更

### 简化动画时长

**变更前**：
```css
duration-short1: 50ms
duration-short2: 100ms
duration-short3: 150ms
duration-short4: 200ms
duration-medium1: 250ms
...
duration-extra-long4: 1000ms
```

**变更后**：
```css
duration-fast: 150ms    /* 快速反馈（hover） */
duration-normal: 200ms  /* 标准过渡（大部分交互） */
duration-slow: 300ms    /* 页面切换（最长） */
```

**迁移策略**：

**❌ 删除旧类名**：
```vue
<div class="transition-all duration-medium2"> <!-- 300ms -->
```

**✅ 使用新类名**：
```vue
<div class="transition-all duration-normal"> <!-- 200ms -->
```

**推荐用法**：
```vue
<!-- 按钮 hover -->
<button class="transition-colors duration-fast">按钮</button>

<!-- 模态框淡入 -->
<div class="transition-opacity duration-normal">模态框</div>

<!-- 页面切换 -->
<div class="transition-transform duration-slow">页面</div>
```

---

## 禁止使用的效果

### 1. 渐变背景

**❌ 禁止**：
```vue
<div class="bg-gradient-to-r from-primary-500 to-secondary-600">
<div class="bg-gradient-to-br from-primary-400 via-secondary-500 to-primary-600">
```

**✅ 替换为**：
```vue
<div class="bg-primary-600">
```

---

### 2. 毛玻璃效果

**❌ 禁止**：
```css
backdrop-filter: blur(20px);
```

**✅ 替换为**：
```css
background: rgba(255, 255, 255, 0.95);
/* 或直接使用纯色 */
background: #FFFFFF;
```

---

### 3. 过度动画

**❌ 禁止**：
```vue
<div class="transform hover:scale-105 hover:-translate-y-2">
```

**✅ 替换为**：
```vue
<div class="hover:shadow-md transition-shadow duration-fast">
```

---

### 4. 装饰性效果

**❌ 禁止**：
```vue
<!-- 彩色边框 -->
<div class="border-t-4 border-primary-500">

<!-- 脉动动画 -->
<div class="animate-pulse">

<!-- 旋转动画 -->
<div class="animate-spin">
```

**✅ 仅在必要时使用**：
```vue
<!-- Loading 指示器允许使用 animate-spin -->
<div v-if="loading" class="animate-spin">⏳</div>
```

---

## 全局查找与替换清单

### 必须替换

```bash
# 1. 移除 secondary 色系
grep -r "secondary-" web/src/
# 替换为 primary-* 或 neutral-*

# 2. 移除彩色阴影
grep -r "shadow-primary\|shadow-secondary\|shadow-success\|shadow-error" web/src/
# 替换为 shadow-sm / shadow / shadow-md

# 3. 移除渐变背景
grep -r "bg-gradient-to-" web/src/
# 替换为单色背景

# 4. 移除毛玻璃
grep -r "backdrop-filter" web/src/
# 删除该属性

# 5. 移除过度动画
grep -r "hover:scale\|hover:-translate-y" web/src/
# 简化为 hover:shadow-md
```

### 可选优化

```bash
# 1. 使用新的排版工具类
grep -r "text-2xl font-bold" web/src/
# 替换为 text-h2

# 2. 使用新的动画时长
grep -r "duration-\(short\|medium\|long\)" web/src/
# 替换为 duration-fast / duration-normal / duration-slow

# 3. 使用新的字距类名
# 为大标题添加 tracking-tighter
# 为标签添加 tracking-wide
```

---

## 检查清单

完成迁移后，检查以下事项：

### 视觉检查
- [ ] 主色按钮对比度足够
- [ ] 文本可读性良好
- [ ] 卡片阴影不过重
- [ ] 圆角不过度圆润
- [ ] 无渐变色残留
- [ ] 无彩色阴影残留

### 代码检查
- [ ] 无 `secondary-*` 类名
- [ ] 无 `shadow-primary` 等彩色阴影
- [ ] 无 `bg-gradient-to-*` 渐变背景
- [ ] 无 `backdrop-filter` 毛玻璃效果
- [ ] 无 `hover:scale` 等过度动画
- [ ] 无 `duration-short1` 等旧动画时长

### 功能检查
- [ ] 按钮 hover 效果正常
- [ ] 链接颜色清晰可辨
- [ ] 表单输入框 focus 状态明显
- [ ] Loading 状态显示正常
- [ ] 模态框动画流畅

---

## 常见问题

### Q: 主色变暗后，按钮对比度不够怎么办？

**A**: 使用更深的色阶：

```vue
<!-- 旧代码 -->
<button class="bg-primary-500 text-white">

<!-- 新代码：使用 600 色阶 -->
<button class="bg-primary-600 text-white">
```

---

### Q: 移除渐变后，设计感觉太单调？

**A**: 通过以下方式增强设计：

1. 使用微妙的阴影
```vue
<div class="shadow-sm hover:shadow-md transition-shadow">
```

2. 使用边框强调
```vue
<div class="border-l-4 border-primary-600">
```

3. 使用图标和排版
```vue
<div class="flex items-center gap-3">
  <Icon class="text-primary-600" />
  <h3 class="text-h3">标题</h3>
</div>
```

---

### Q: 圆角变小后，感觉太硬怎么办？

**A**: 这是预期效果。克制设计追求专业感而非亲和感。如果确实需要，可使用：

```vue
<!-- 标准圆角 -->
<div class="rounded-md">  <!-- 8px -->

<!-- 大圆角（保留特殊场景） -->
<div class="rounded-lg">  <!-- 12px -->
```

---

### Q: 正文字号从 16px 变为 15px 会影响可读性吗？

**A**: 不会。15px 是 Medium 等内容平台的标准字号。配合 1.6 的行高，可读性实际上更好。

---

## 迁移时间估算

- **小型项目**（< 20 个页面）：2-4 小时
- **中型项目**（20-50 个页面）：1-2 天
- **大型项目**（> 50 个页面）：3-5 天

---

## 支持

如有迁移问题，请查阅：
- `/docs/frontend-refactor-design.md` - 完整设计系统文档
- `/web/src/design-tokens/` - Design Tokens 源代码
- `/web/tailwind.config.js` - Tailwind 配置

---

**版本**：v1.0
**更新日期**：2026-02-05
