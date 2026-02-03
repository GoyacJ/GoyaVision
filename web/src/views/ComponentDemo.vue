<template>
  <div class="min-h-screen bg-surface p-8">
    <div class="max-w-7xl mx-auto">
      <!-- 页面头部 -->
      <header class="mb-8">
        <h1 class="text-4xl font-bold text-text-primary mb-2">
          GoyaVision 组件库
        </h1>
        <p class="text-text-secondary">
          Material Design 3 风格组件展示
        </p>
      </header>
      
      <!-- GvButton 展示 -->
      <section class="mb-12">
        <h2 class="text-2xl font-semibold text-text-primary mb-4">
          GvButton - 按钮组件
        </h2>
        
        <!-- 变体展示 -->
        <div class="bg-white rounded-xl p-6 shadow-md mb-6">
          <h3 class="text-lg font-medium text-text-primary mb-4">变体（Variants）</h3>
          <div class="flex flex-wrap gap-4">
            <GvButton variant="filled">Filled</GvButton>
            <GvButton variant="tonal">Tonal</GvButton>
            <GvButton variant="outlined">Outlined</GvButton>
            <GvButton variant="text">Text</GvButton>
          </div>
        </div>
        
        <!-- 颜色展示 -->
        <div class="bg-white rounded-xl p-6 shadow-md mb-6">
          <h3 class="text-lg font-medium text-text-primary mb-4">颜色主题</h3>
          <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
            <GvButton color="primary">Primary</GvButton>
            <GvButton color="secondary">Secondary</GvButton>
            <GvButton color="success">Success</GvButton>
            <GvButton color="error">Error</GvButton>
            <GvButton color="warning">Warning</GvButton>
            <GvButton color="info">Info</GvButton>
          </div>
        </div>
        
        <!-- 尺寸展示 -->
        <div class="bg-white rounded-xl p-6 shadow-md mb-6">
          <h3 class="text-lg font-medium text-text-primary mb-4">尺寸（Sizes）</h3>
          <div class="flex items-center gap-4">
            <GvButton size="small">Small</GvButton>
            <GvButton size="medium">Medium</GvButton>
            <GvButton size="large">Large</GvButton>
          </div>
        </div>
        
        <!-- 图标展示 -->
        <div class="bg-white rounded-xl p-6 shadow-md mb-6">
          <h3 class="text-lg font-medium text-text-primary mb-4">图标按钮</h3>
          <div class="flex flex-wrap gap-4">
            <GvButton icon="Plus" icon-position="left">新建</GvButton>
            <GvButton icon="Edit" icon-position="left" variant="tonal">编辑</GvButton>
            <GvButton icon="Delete" icon-position="left" variant="outlined" color="error">删除</GvButton>
            <GvButton icon="Search" variant="outlined" rounded />
            <GvButton icon="Setting" variant="text" rounded />
          </div>
        </div>
        
        <!-- 状态展示 -->
        <div class="bg-white rounded-xl p-6 shadow-md mb-6">
          <h3 class="text-lg font-medium text-text-primary mb-4">状态</h3>
          <div class="flex flex-wrap gap-4">
            <GvButton :loading="true">加载中...</GvButton>
            <GvButton disabled>禁用按钮</GvButton>
            <GvButton variant="tonal" :loading="isLoading" @click="handleLoadingDemo">
              点击测试加载
            </GvButton>
          </div>
        </div>
        
        <!-- 组合示例 -->
        <div class="bg-white rounded-xl p-6 shadow-md">
          <h3 class="text-lg font-medium text-text-primary mb-4">实际应用场景</h3>
          
          <!-- 表单操作 -->
          <div class="mb-6">
            <p class="text-sm text-text-secondary mb-2">表单操作：</p>
            <div class="flex gap-2">
              <GvButton variant="filled" color="primary">提交</GvButton>
              <GvButton variant="tonal" color="secondary">取消</GvButton>
              <GvButton variant="text">重置</GvButton>
            </div>
          </div>
          
          <!-- 数据操作 -->
          <div class="mb-6">
            <p class="text-sm text-text-secondary mb-2">数据操作：</p>
            <div class="flex gap-2">
              <GvButton variant="filled" icon="Plus">新建资产</GvButton>
              <GvButton variant="tonal" icon="Upload">上传</GvButton>
              <GvButton variant="outlined" color="error" icon="Delete">批量删除</GvButton>
            </div>
          </div>
          
          <!-- 块级按钮 -->
          <div>
            <p class="text-sm text-text-secondary mb-2">块级按钮：</p>
            <GvButton block variant="filled" color="primary">立即购买</GvButton>
          </div>
        </div>
      </section>
      
      <!-- 主题切换 -->
      <section class="mb-12">
        <h2 class="text-2xl font-semibold text-text-primary mb-4">
          主题系统
        </h2>
        <div class="bg-white rounded-xl p-6 shadow-md">
          <div class="flex items-center justify-between">
            <div>
              <h3 class="text-lg font-medium text-text-primary mb-1">深色模式</h3>
              <p class="text-sm text-text-secondary">当前主题: {{ actualTheme }}</p>
            </div>
            <GvButton variant="tonal" @click="toggleTheme">
              切换主题
            </GvButton>
          </div>
        </div>
      </section>
      
      <!-- 设计令牌展示 -->
      <section>
        <h2 class="text-2xl font-semibold text-text-primary mb-4">
          设计令牌（Design Tokens）
        </h2>
        
        <!-- 颜色系统 -->
        <div class="bg-white rounded-xl p-6 shadow-md mb-6">
          <h3 class="text-lg font-medium text-text-primary mb-4">颜色系统</h3>
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <div class="h-20 bg-primary-600 rounded-lg mb-2"></div>
              <p class="text-sm text-text-secondary">Primary</p>
            </div>
            <div>
              <div class="h-20 bg-secondary-600 rounded-lg mb-2"></div>
              <p class="text-sm text-text-secondary">Secondary</p>
            </div>
            <div>
              <div class="h-20 bg-success-600 rounded-lg mb-2"></div>
              <p class="text-sm text-text-secondary">Success</p>
            </div>
            <div>
              <div class="h-20 bg-error-600 rounded-lg mb-2"></div>
              <p class="text-sm text-text-secondary">Error</p>
            </div>
          </div>
        </div>
        
        <!-- Tailwind 工具类展示 -->
        <div class="bg-white rounded-xl p-6 shadow-md">
          <h3 class="text-lg font-medium text-text-primary mb-4">Tailwind CSS 已启用</h3>
          <div class="space-y-2 text-sm text-text-secondary">
            <p>✅ 颜色系统（primary, secondary, success, error...）</p>
            <p>✅ 间距系统（基于 8px 网格）</p>
            <p>✅ 阴影系统（Material Design 3 层级）</p>
            <p>✅ 圆角系统（sm, md, lg, xl...）</p>
            <p>✅ 响应式断点（sm, md, lg, xl, 2xl）</p>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { GvButton } from '@/components'
import { useTheme } from '@/composables'

const { actualTheme, toggleTheme } = useTheme()

const isLoading = ref(false)

const handleLoadingDemo = () => {
  isLoading.value = true
  setTimeout(() => {
    isLoading.value = false
  }, 2000)
}
</script>
