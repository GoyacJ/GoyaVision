<template>
  <div class="page-header">
    <!-- 面包屑导航 -->
    <el-breadcrumb v-if="breadcrumb && breadcrumb.length" class="mb-4" separator="/">
      <el-breadcrumb-item
        v-for="(item, index) in breadcrumb"
        :key="index"
        :to="item.to"
      >
        {{ item.label }}
      </el-breadcrumb-item>
    </el-breadcrumb>
    
    <!-- 头部主体 -->
    <GvFlex justify="between" align="start" class="mb-6">
      <!-- 左侧：标题和描述 -->
      <div class="page-header__main">
        <GvFlex align="center" gap="md" class="mb-2">
          <!-- 返回按钮 -->
          <GvButton
            v-if="showBack"
            variant="text"
            size="small"
            @click="handleBack"
          >
            <template #icon>
              <el-icon><ArrowLeft /></el-icon>
            </template>
            {{ backText }}
          </GvButton>
          
          <!-- 标题 -->
          <h1 class="page-header__title">{{ title }}</h1>
        </GvFlex>
        
        <!-- 描述 -->
        <p v-if="description" class="page-header__description">
          {{ description }}
        </p>
      </div>
      
      <!-- 右侧：操作区域 -->
      <div v-if="$slots.actions" class="page-header__actions">
        <slot name="actions" />
      </div>
    </GvFlex>
    
    <!-- 额外内容区域 -->
    <div v-if="$slots.extra" class="page-header__extra">
      <slot name="extra" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { GvFlex, GvButton } from '@/components'
import { ArrowLeft } from '@element-plus/icons-vue'
import type { PageHeaderProps, PageHeaderEmits } from './types'

const props = withDefaults(defineProps<PageHeaderProps>(), {
  showBack: false,
  backText: '返回'
})

const emit = defineEmits<PageHeaderEmits>()

// 返回处理
const handleBack = () => {
  emit('back')
}
</script>

<style scoped>
.page-header {
  @apply mb-6;
}

.page-header__title {
  @apply text-3xl font-bold text-text-primary;
  @apply m-0;
}

.page-header__description {
  @apply text-text-secondary text-base;
  @apply m-0 mt-2;
}

.page-header__actions {
  @apply flex items-center gap-3;
}

/* 深色模式 */
.dark .page-header__title {
  @apply text-text-inverse;
}

.dark .page-header__description {
  @apply text-neutral-400;
}
</style>
