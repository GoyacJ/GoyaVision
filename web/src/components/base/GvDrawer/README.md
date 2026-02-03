# GvDrawer - 抽屉组件

Material Design 3 风格的抽屉组件，用于从屏幕边缘滑出的面板。

## 基本用法

```vue
<template>
  <div>
    <GvButton @click="visible = true">打开抽屉</GvButton>
    
    <GvDrawer v-model="visible" title="抽屉标题">
      <p>这是抽屉的内容</p>
    </GvDrawer>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { GvDrawer, GvButton } from '@/components'

const visible = ref(false)
</script>
```

## 方向

```vue
<GvDrawer v-model="visible" direction="left" title="从左侧打开">
  <p>内容</p>
</GvDrawer>

<GvDrawer v-model="visible" direction="right" title="从右侧打开">
  <p>内容</p>
</GvDrawer>

<GvDrawer v-model="visible" direction="top" title="从顶部打开">
  <p>内容</p>
</GvDrawer>

<GvDrawer v-model="visible" direction="bottom" title="从底部打开">
  <p>内容</p>
</GvDrawer>
```

## 尺寸

### 预设尺寸

```vue
<GvDrawer v-model="visible" size="small" title="小尺寸">
  <p>300px 宽</p>
</GvDrawer>

<GvDrawer v-model="visible" size="medium" title="中等尺寸">
  <p>450px 宽</p>
</GvDrawer>

<GvDrawer v-model="visible" size="large" title="大尺寸">
  <p>600px 宽</p>
</GvDrawer>

<GvDrawer v-model="visible" size="full" title="全屏">
  <p>100vw 宽</p>
</GvDrawer>
```

### 自定义尺寸

```vue
<!-- 自定义宽度 -->
<GvDrawer v-model="visible" width="800px" title="自定义宽度">
  <p>800px 宽</p>
</GvDrawer>

<!-- 自定义高度（top/bottom 方向） -->
<GvDrawer
  v-model="visible"
  direction="bottom"
  height="500px"
  title="自定义高度"
>
  <p>500px 高</p>
</GvDrawer>
```

## 自定义头部和底部

```vue
<GvDrawer v-model="visible">
  <template #header>
    <div class="flex items-center gap-2">
      <el-icon><InfoFilled /></el-icon>
      <h3>自定义标题</h3>
    </div>
  </template>
  
  <p>抽屉内容</p>
  
  <template #footer>
    <div class="flex justify-between">
      <GvButton variant="text">帮助</GvButton>
      <div class="flex gap-2">
        <GvButton variant="tonal" @click="visible = false">取消</GvButton>
        <GvButton variant="filled" @click="handleSave">保存</GvButton>
      </div>
    </div>
  </template>
</GvDrawer>
```

## 无底部按钮

```vue
<GvDrawer v-model="visible" title="详情" :show-footer="false">
  <p>详情内容</p>
</GvDrawer>
```

## 禁用遮罩点击关闭

```vue
<GvDrawer
  v-model="visible"
  title="重要操作"
  :close-on-click-modal="false"
>
  <p>必须完成操作才能关闭</p>
</GvDrawer>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| modelValue | `boolean` | - | 是否显示（必填） |
| title | `string` | - | 标题 |
| direction | `'left' \| 'right' \| 'top' \| 'bottom'` | `'right'` | 方向 |
| size | `'small' \| 'medium' \| 'large' \| 'full'` | `'medium'` | 尺寸 |
| width | `string` | - | 自定义宽度 |
| height | `string` | - | 自定义高度 |
| showClose | `boolean` | `true` | 是否显示关闭按钮 |
| closeOnClickModal | `boolean` | `true` | 点击遮罩是否关闭 |
| closeOnPressEscape | `boolean` | `true` | 按 ESC 是否关闭 |
| destroyOnClose | `boolean` | `false` | 关闭时销毁子元素 |
| showFooter | `boolean` | `true` | 是否显示底部 |
| confirmText | `string` | `'确定'` | 确认按钮文字 |
| cancelText | `string` | `'取消'` | 取消按钮文字 |
| showConfirm | `boolean` | `true` | 是否显示确认按钮 |
| showCancel | `boolean` | `true` | 是否显示取消按钮 |
| confirmLoading | `boolean` | `false` | 确认按钮加载状态 |
| customClass | `string` | - | 自定义类名 |
| zIndex | `number` | `1000` | z-index |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| update:modelValue | `(value: boolean)` | 更新显示状态 |
| open | `()` | 打开事件 |
| opened | `()` | 打开动画结束 |
| close | `()` | 关闭事件 |
| closed | `()` | 关闭动画结束 |
| confirm | `()` | 确认事件 |
| cancel | `()` | 取消事件 |

## Slots

| 插槽名 | 说明 |
|--------|------|
| header | 头部内容 |
| default | 主体内容 |
| footer | 底部内容 |

## 使用场景

### 详情查看

```vue
<GvDrawer
  v-model="detailVisible"
  title="资产详情"
  :show-confirm="false"
  cancel-text="关闭"
>
  <div class="space-y-6">
    <div>
      <h4 class="font-semibold mb-3">基本信息</h4>
      <div class="space-y-2">
        <p><span class="text-text-secondary">名称：</span>{{ asset.name }}</p>
        <p><span class="text-text-secondary">类型：</span>{{ asset.type }}</p>
      </div>
    </div>
  </div>
</GvDrawer>
```

### 编辑表单

```vue
<GvDrawer
  v-model="editVisible"
  title="编辑资产"
  :confirm-loading="loading"
  @confirm="handleUpdate"
>
  <el-form :model="form" label-width="100px">
    <el-form-item label="名称">
      <GvInput v-model="form.name" />
    </el-form-item>
    <el-form-item label="描述">
      <el-input v-model="form.description" type="textarea" :rows="4" />
    </el-form-item>
  </el-form>
</GvDrawer>
```

### 筛选面板

```vue
<GvDrawer
  v-model="filterVisible"
  direction="left"
  title="筛选条件"
  size="small"
>
  <div class="space-y-4">
    <div>
      <label class="block mb-2">资产类型</label>
      <GvSelect v-model="filters.type" :options="typeOptions" />
    </div>
    <div>
      <label class="block mb-2">状态</label>
      <GvSelect v-model="filters.status" :options="statusOptions" />
    </div>
  </div>
  
  <template #footer>
    <div class="flex gap-2">
      <GvButton variant="text" @click="resetFilters">重置</GvButton>
      <GvButton variant="filled" class="flex-1" @click="applyFilters">
        应用筛选
      </GvButton>
    </div>
  </template>
</GvDrawer>
```

### 通知面板

```vue
<GvDrawer
  v-model="notificationVisible"
  title="通知"
  direction="right"
  size="small"
  :show-footer="false"
>
  <div class="space-y-4">
    <div
      v-for="notification in notifications"
      :key="notification.id"
      class="p-4 bg-neutral-50 rounded-lg hover:bg-neutral-100 cursor-pointer"
    >
      <h4 class="font-medium mb-1">{{ notification.title }}</h4>
      <p class="text-sm text-text-secondary">{{ notification.message }}</p>
      <span class="text-xs text-text-tertiary">{{ notification.time }}</span>
    </div>
  </div>
</GvDrawer>
```

## 最佳实践

1. **方向选择**：
   - right: 最常用，详情、表单
   - left: 筛选面板、导航
   - bottom: 移动端操作面板
   - top: 通知提示

2. **尺寸选择**：
   - small: 简单信息、筛选
   - medium: 详情、表单
   - large: 复杂内容
   - full: 完整页面

3. **关闭控制**：
   - 编辑表单禁用遮罩点击关闭
   - 多步骤流程禁用 ESC 关闭

4. **用户体验**：
   - 提供清晰的标题
   - 滚动内容使用自定义滚动条
   - 重要操作添加确认

5. **性能优化**：
   - 复杂内容使用 destroyOnClose
   - 避免在抽屉中进行大量计算

## Modal vs Drawer

| 特性 | Modal | Drawer |
|------|-------|--------|
| **位置** | 居中 | 屏幕边缘 |
| **用途** | 重要操作、确认 | 详情、筛选、辅助面板 |
| **视觉影响** | 强，遮盖内容 | 弱，侧边滑出 |
| **适用场景** | 需要用户立即响应 | 查看详情、辅助功能 |

**选择建议**：
- 确认、警告 → Modal
- 详情查看、筛选 → Drawer
