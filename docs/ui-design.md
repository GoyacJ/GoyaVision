# GoyaVision UI 设计规范

> 参考 ModelScope 等现代化 AI 平台的设计风格，打造具有科技感和专业性的用户界面

## 设计理念

### 核心原则

1. **现代化 (Modern)** - 采用最新的设计趋势和视觉效果
2. **专业性 (Professional)** - 符合 AI 平台的专业定位
3. **易用性 (Usable)** - 清晰的信息层次和流畅的交互
4. **一致性 (Consistent)** - 统一的设计语言和视觉风格
5. **响应式 (Responsive)** - 适配各种屏幕尺寸

### 设计风格

- **Glassmorphism (磨砂玻璃效果)** - 半透明背景 + 模糊效果
- **Gradient (渐变色)** - 大量使用渐变色增强视觉冲击
- **Smooth Animations (流畅动画)** - 微交互和过渡动画
- **Card-based Layout (卡片式布局)** - 清晰的内容分组
- **Soft Shadows (柔和阴影)** - 多层次阴影系统

---

## 配色系统

### 主色调

```css
/* 主渐变色 - 蓝紫色系 */
--primary-gradient: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

/* 主色 */
--primary-color: #667eea;
--primary-dark: #5568d3;
--primary-light: #8798ff;

/* 辅助色 */
--secondary-color: #764ba2;
--secondary-dark: #65408b;
--secondary-light: #8d5cb8;
```

### 功能色

```css
/* 成功 - 青色渐变 */
--success-gradient: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
--success-color: #67c23a;

/* 提示 - 绿色渐变 */
--info-gradient: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
--info-color: #909399;

/* 警告 */
--warning-gradient: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
--warning-color: #e6a23c;

/* 错误/危险 */
--danger-gradient: linear-gradient(135deg, #ff6b6b 0%, #c92a2a 100%);
--danger-color: #f56c6c;
```

### 中性色

```css
/* 背景色 */
--background-primary: #f5f7fa;
--background-secondary: #ffffff;
--background-gradient: linear-gradient(135deg, #f5f7fa 0%, #e8ecf1 100%);

/* 文字色 */
--text-primary: #333333;
--text-secondary: #666666;
--text-tertiary: #999999;
--text-placeholder: #c0c4cc;

/* 边框色 */
--border-color: #dcdfe6;
--border-light: rgba(102, 126, 234, 0.1);
```

---

## 圆角系统

```css
/* 统一圆角规范 */
--radius-sm: 8px;    /* 小按钮、标签 */
--radius-md: 12px;   /* 输入框、普通按钮 */
--radius-lg: 16px;   /* 卡片、对话框 */
--radius-xl: 20px;   /* 登录框、大卡片 */
--radius-full: 50%;  /* 圆形头像 */
```

---

## 阴影系统

```css
/* 阴影层级 */
--shadow-sm: 0 2px 8px rgba(0, 0, 0, 0.05);          /* 输入框 */
--shadow-md: 0 4px 12px rgba(102, 126, 234, 0.15);   /* 悬停 */
--shadow-lg: 0 8px 32px rgba(31, 38, 135, 0.15);     /* 卡片 */
--shadow-xl: 0 20px 60px rgba(31, 38, 135, 0.2);     /* 对话框 */

/* 彩色阴影（按钮、图标等） */
--shadow-primary: 0 8px 24px rgba(102, 126, 234, 0.4);
--shadow-success: 0 8px 24px rgba(103, 194, 58, 0.4);
--shadow-warning: 0 8px 24px rgba(230, 162, 60, 0.4);
```

---

## 动画系统

### 过渡效果

```css
/* 统一过渡 */
--transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

/* 缓动函数 */
--ease-in-out: cubic-bezier(0.4, 0, 0.2, 1);
--ease-out: cubic-bezier(0, 0, 0.2, 1);
--ease-in: cubic-bezier(0.4, 0, 1, 1);
```

### 关键帧动画

#### 淡入动画

```css
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
/* 使用: animation: fadeIn 0.4s ease-out; */
```

#### 右滑入动画

```css
@keyframes slideInRight {
  from {
    opacity: 0;
    transform: translateX(20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}
/* 使用: animation: slideInRight 0.4s ease-out; */
```

#### 脉冲动画

```css
@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
  }
  50% {
    transform: scale(1.05);
    box-shadow: 0 12px 32px rgba(102, 126, 234, 0.6);
  }
}
/* 使用: animation: pulse 2s infinite; */
```

#### 浮动动画

```css
@keyframes float {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
  }
}
/* 使用: animation: float 20s infinite ease-in-out; */
```

---

## 组件样式规范

### 按钮 (Button)

#### 主要按钮

```css
.button-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 12px;
  padding: 12px 24px;
  font-weight: 600;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
  transition: all 0.3s;
}

.button-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.5);
}
```

#### 次要按钮

```css
.button-secondary {
  background: transparent;
  color: #667eea;
  border: 1px solid #667eea;
  border-radius: 12px;
  padding: 12px 24px;
  font-weight: 600;
  transition: all 0.3s;
}

.button-secondary:hover {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
}
```

### 输入框 (Input)

```css
.input-wrapper {
  border-radius: 12px;
  padding: 12px 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  transition: all 0.3s;
  border: 1px solid #e4e7ed;
}

.input-wrapper:hover {
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
}

.input-wrapper:focus-within {
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.25);
  border-color: #667eea;
}
```

### 卡片 (Card)

```css
.card {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(31, 38, 135, 0.12);
  border: 1px solid rgba(102, 126, 234, 0.1);
  overflow: hidden;
}

.card-header {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
  padding: 20px 24px;
}

.card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(31, 38, 135, 0.25);
}
```

### 标签 (Tag)

```css
.tag {
  border-radius: 6px;
  padding: 4px 12px;
  font-weight: 500;
  border: none;
  display: inline-flex;
  align-items: center;
}

.tag-primary {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.15) 0%, rgba(118, 75, 162, 0.15) 100%);
  color: #667eea;
}

.tag-success {
  background: linear-gradient(135deg, rgba(103, 194, 58, 0.15) 0%, rgba(56, 249, 215, 0.15) 100%);
  color: #67c23a;
}
```

### 表格 (Table)

```css
.table {
  border-radius: 12px;
  overflow: hidden;
}

.table-header {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  font-weight: 600;
}

.table-row {
  transition: all 0.3s;
}

.table-row:hover {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.03) 0%, rgba(118, 75, 162, 0.03) 100%);
}
```

---

## 页面布局规范

### 页面容器

```css
.page-container {
  padding: 32px;
  animation: fadeIn 0.4s ease-out;
}
```

### 间距系统

```css
/* 统一间距规范 */
--spacing-xs: 8px;
--spacing-sm: 12px;
--spacing-md: 16px;
--spacing-lg: 24px;
--spacing-xl: 32px;
--spacing-2xl: 48px;
```

### 内容宽度

```css
--content-max-width: 1200px;
--content-min-width: 320px;
```

---

## 响应式设计指南

为了确保 GoyaVision 在不同设备（桌面、平板、手机）上的一致体验，采用以下响应式策略。

### 断点系统

沿用 Tailwind CSS 默认断点：
- `sm`: 640px
- `md`: 768px
- `lg`: 1024px
- `xl`: 1280px
- `2xl`: 1536px

### 布局策略

1.  **导航栏**
    - **Desktop (`lg`+)**: 顶部水平导航栏，展示 Logo、菜单、用户信息。
    - **Mobile/Tablet**: 隐藏水平菜单，使用左侧汉堡菜单 (Hamburger Menu) 呼出抽屉式导航 (`el-drawer`)。

2.  **页面布局**
    - **左右布局页 (如资产库)**：
        - **Desktop**: `flex-row`，左侧固定宽筛选栏，右侧自适应内容区。
        - **Mobile**: `flex-col`，筛选栏宽度 `w-full` 并置顶（或折叠），内容区顺延。
    - **工具栏**: 使用 `flex-wrap` 确保搜索框与操作按钮在窄屏下自动换行。

3.  **列表与视图**
    - **网格视图 (Grid)**: 使用响应式 Grid Class，如 `grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5`。
    - **表格视图 (List)**: 移动端体验较差，建议通过 `useBreakpoint` 检测 `isMobile`，强制切换为 **网格视图**。

### 组件适配

- **GvModal / GvDrawer**: 宽度应设为响应式，移动端推荐 `90%` 或 `100%` 宽度。
- **GvSpace**: 开启 `wrap` 属性或手动添加 `flex-wrap` 以防止溢出。

---

## 图标系统

使用 Element Plus Icons 作为主要图标库：

- **用户相关**: User, UserFilled
- **操作相关**: Edit, Delete, View, Search
- **媒体相关**: VideoCameraFilled, Picture, Microphone
- **状态相关**: Success, Warning, Error, Info
- **导航相关**: ArrowDown, ArrowRight, Menu

---

## 字体系统

### 字体族

```css
font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 
             'Helvetica Neue', Arial, 'Noto Sans', sans-serif,
             'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 
             'Noto Color Emoji';
```

### 字号规范

```css
--font-size-xs: 12px;   /* 辅助文字 */
--font-size-sm: 13px;   /* 提示文字 */
--font-size-base: 14px; /* 正文 */
--font-size-md: 16px;   /* 标题 */
--font-size-lg: 18px;   /* 卡片标题 */
--font-size-xl: 22px;   /* 页面标题 */
--font-size-2xl: 28px;  /* Logo */
--font-size-3xl: 32px;  /* 特殊标题 */
```

### 字重规范

```css
--font-weight-normal: 400;   /* 正文 */
--font-weight-medium: 500;   /* 强调 */
--font-weight-semibold: 600; /* 次标题 */
--font-weight-bold: 700;     /* 标题 */
--font-weight-extrabold: 800; /* Logo */
```

---

## 特殊效果

### 磨砂玻璃效果

```css
.glassmorphism {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}
```

### 渐变文字

```css
.gradient-text {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
```

### 自定义滚动条

```css
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, #5568d3 0%, #65408b 100%);
}
```

---

## 可访问性 (Accessibility)

### 颜色对比度

- 确保文字与背景的对比度至少为 4.5:1（正文）
- 大号文字（18px+ 或 14px+ bold）至少 3:1

### 键盘导航

- 所有交互元素支持键盘访问
- Tab 键导航顺序合理
- 焦点状态清晰可见

### 屏幕阅读器

- 使用语义化 HTML 标签
- 添加必要的 aria 属性
- 提供有意义的 alt 文本

---

## 性能优化

### CSS 优化

1. 使用 CSS 变量减少重复代码
2. 避免过度使用阴影和滤镜
3. 动画使用 transform 和 opacity（GPU 加速）
4. 合理使用 will-change 属性

### 图片优化

1. 使用 WebP 格式
2. 适当的图片尺寸
3. 懒加载非关键图片
4. 使用 SVG 图标

---

## 设计系统维护

### 更新原则

1. **保持一致性** - 任何修改需要全局考虑
2. **渐进增强** - 优先考虑核心功能
3. **向后兼容** - 避免破坏性更改
4. **文档先行** - 更新设计规范文档

### 版本管理

- 主版本号：重大设计语言变更
- 次版本号：新增组件或样式
- 修订号：bug 修复和微调

---

## 实施清单

### ✅ 已完成

- [x] 全局样式系统（CSS 变量）
- [x] 登录页面设计
- [x] 主布局优化
- [x] 资产管理页面优化
- [x] 自定义滚动条
- [x] 动画系统
- [x] 响应式布局重构 (2026-02-08)

### 🚧 进行中

- [ ] 其他页面优化（Workflow、Task、Operator、System）
- [ ] 深色模式支持
- [ ] 多语言支持

### 📋 待实施

- [ ] 组件库文档
- [ ] 设计 Token 导出工具
- [ ] Storybook 组件展示
- [ ] UI 自动化测试
- [ ] 设计系统网站

---

## 参考资源

### 设计灵感

- [ModelScope](https://modelscope.cn/) - AI 模型社区
- [Ant Design](https://ant.design/) - 企业级 UI 设计语言
- [Material Design 3](https://m3.material.io/) - Google 设计系统
- [Fluent Design](https://www.microsoft.com/design/fluent/) - Microsoft 设计语言

### 技术文档

- [Element Plus](https://element-plus.org/) - Vue 3 组件库
- [CSS-Tricks](https://css-tricks.com/) - CSS 技巧
- [MDN Web Docs](https://developer.mozilla.org/) - Web 标准文档

---

## 更新日志

### 2026-02-08

- 📱 新增响应式设计规范
- ✨ 重构全局导航布局 (Drawer/Horizontal)
- ✨ 优化资产库布局 (Stack/Row, Grid/List)

### 2026-02-03

- 🎨 初始化 UI 设计系统
- ✨ 完成登录页面重设计
- ✨ 优化主布局和资产管理页面
- 📝 创建设计规范文档

---

**维护者**: GoyaVision Team  
**最后更新**: 2026-02-08  
**版本**: 1.1.0
