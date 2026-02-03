<template>
  <GvCard
    :class="cardClasses"
    hoverable
    @click="handleClick"
  >
    <!-- 选择框 -->
    <div v-if="selectable" class="asset-card__checkbox" @click.stop>
      <el-checkbox
        :model-value="selected"
        @change="handleSelectChange"
      />
    </div>
    
    <!-- 缩略图 -->
    <div class="asset-card__thumbnail">
      <img
        v-if="asset.thumbnail"
        :src="asset.thumbnail"
        :alt="asset.name"
        class="asset-card__image"
      />
      <div v-else class="asset-card__placeholder">
        <el-icon :size="48" class="text-neutral-400">
          <VideoCamera />
        </el-icon>
      </div>
      
      <!-- 状态徽章（覆盖在缩略图上） -->
      <div class="asset-card__status">
        <StatusBadge :status="asset.status as any" />
      </div>
    </div>
    
    <!-- 内容 -->
    <div class="asset-card__content">
      <!-- 标题 -->
      <h3 class="asset-card__title">{{ asset.name }}</h3>
      
      <!-- 描述 -->
      <p v-if="asset.description" class="asset-card__description">
        {{ asset.description }}
      </p>
      
      <!-- 元数据 -->
      <GvFlex class="asset-card__meta" wrap gap="xs">
        <GvTag v-if="asset.type" size="small" variant="tonal">
          {{ asset.type }}
        </GvTag>
        <span v-if="asset.size" class="text-text-tertiary text-sm">
          {{ asset.size }}
        </span>
        <span v-if="asset.duration" class="text-text-tertiary text-sm">
          {{ asset.duration }}
        </span>
      </GvFlex>
      
      <!-- 时间信息 -->
      <div v-if="asset.createdAt" class="asset-card__time">
        <el-icon class="text-text-tertiary"><Clock /></el-icon>
        <span class="text-text-tertiary text-xs ml-1">
          {{ asset.createdAt }}
        </span>
      </div>
    </div>
    
    <!-- 操作按钮（悬停显示） -->
    <div v-if="showActions" class="asset-card__actions">
      <GvSpace size="xs">
        <GvButton
          v-if="showDetail"
          size="small"
          variant="tonal"
          @click.stop="handleDetail"
        >
          详情
        </GvButton>
        <GvButton
          v-if="showEdit"
          size="small"
          @click.stop="handleEdit"
        >
          编辑
        </GvButton>
        <GvButton
          v-if="showDelete"
          size="small"
          variant="text"
          @click.stop="handleDelete"
        >
          删除
        </GvButton>
      </GvSpace>
    </div>
  </GvCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'
import { GvCard, GvFlex, GvTag, GvButton, GvSpace, StatusBadge } from '@/components'
import { VideoCamera, Clock } from '@element-plus/icons-vue'
import type { AssetCardProps, AssetCardEmits } from './types'

const props = withDefaults(defineProps<AssetCardProps>(), {
  selectable: false,
  selected: false,
  showActions: true,
  showDetail: true,
  showEdit: true,
  showDelete: true
})

const emit = defineEmits<AssetCardEmits>()

// 卡片类名
const cardClasses = computed(() => {
  return cn(
    'asset-card',
    'cursor-pointer transition-all duration-200',
    props.selected && 'ring-2 ring-primary-500'
  )
})

// 点击卡片
const handleClick = () => {
  emit('click', props.asset)
}

// 选择变化
const handleSelectChange = (value: boolean) => {
  emit('select', value, props.asset)
}

// 查看详情
const handleDetail = () => {
  emit('detail', props.asset)
}

// 编辑
const handleEdit = () => {
  emit('edit', props.asset)
}

// 删除
const handleDelete = () => {
  emit('delete', props.asset)
}
</script>

<style scoped>
.asset-card {
  @apply relative overflow-hidden;
}

.asset-card__checkbox {
  @apply absolute top-3 left-3 z-10;
}

.asset-card__thumbnail {
  @apply relative w-full aspect-video bg-neutral-100 overflow-hidden;
}

.asset-card__image {
  @apply w-full h-full object-cover;
}

.asset-card__placeholder {
  @apply w-full h-full flex items-center justify-center bg-neutral-100;
}

.asset-card__status {
  @apply absolute top-3 right-3;
}

.asset-card__content {
  @apply p-4 space-y-2;
}

.asset-card__title {
  @apply text-base font-semibold text-text-primary;
  @apply m-0 truncate;
}

.asset-card__description {
  @apply text-sm text-text-secondary;
  @apply m-0 line-clamp-2;
}

.asset-card__meta {
  @apply pt-2;
}

.asset-card__time {
  @apply flex items-center pt-2;
}

.asset-card__actions {
  @apply absolute bottom-4 right-4;
  @apply opacity-0 transition-opacity duration-200;
}

.asset-card:hover .asset-card__actions {
  @apply opacity-100;
}

/* 深色模式 */
.dark .asset-card__placeholder {
  @apply bg-neutral-800;
}

.dark .asset-card__title {
  @apply text-text-inverse;
}

.dark .asset-card__description {
  @apply text-neutral-400;
}
</style>
