<template>
  <GvCard
    :class="cardClasses"
    shadow="sm"
    padding="none"
    hoverable
    @click="handleClick"
  >
    <!-- 缩略图区域 -->
    <div class="asset-card__thumbnail">
      <!-- 复选框（可选模式） -->
      <div
        v-if="selectable"
        class="asset-card__checkbox"
        @click.stop="handleSelectToggle"
      >
        <el-checkbox :model-value="selected" />
      </div>

      <!-- 缩略图 -->
      <div class="asset-card__image">
        <img
          v-if="thumbnailUrl"
          :src="thumbnailUrl"
          :alt="asset.name"
          class="w-full h-full object-cover"
        />
        <div v-else class="asset-card__placeholder">
          <el-icon :size="48" class="text-neutral-300">
            <component :is="getTypeIcon(asset.type)" />
          </el-icon>
        </div>
      </div>

      <!-- 类型标识（与标签同款 GvTag 样式） -->
      <div class="asset-card__type-badge">
        <GvTag :color="getTypeColor(asset.type)" size="small" variant="tonal">
          <span class="inline-flex items-center gap-1">
            <el-icon :size="14">
              <component :is="getTypeIcon(asset.type)" />
            </el-icon>
            {{ getTypeLabel(asset.type) }}
          </span>
        </GvTag>
      </div>

      <!-- 时长标识（视频/音频） -->
      <div
        v-if="asset.duration && (asset.type === 'video' || asset.type === 'audio')"
        class="asset-card__duration"
      >
        {{ formatDuration(asset.duration) }}
      </div>

      <!-- 悬停操作栏 -->
      <div class="asset-card__actions">
        <GvSpace size="xs">
          <GvButton
            size="small"
            variant="filled"
            color="primary"
            @click.stop="handleView"
          >
            <template #icon>
              <el-icon><View /></el-icon>
            </template>
            查看
          </GvButton>
          <GvButton
            size="small"
            variant="tonal"
            @click.stop="handleEdit"
          >
            <template #icon>
              <el-icon><Edit /></el-icon>
            </template>
            编辑
          </GvButton>
          <GvButton
            size="small"
            variant="text"
            color="error"
            @click.stop="handleDelete"
          >
            <template #icon>
              <el-icon><Delete /></el-icon>
            </template>
            删除
          </GvButton>
        </GvSpace>
      </div>
    </div>

    <!-- 信息区域 -->
    <div class="asset-card__info">
      <!-- 名称 -->
      <div class="asset-card__name" :title="asset.name">
        {{ asset.name }}
      </div>

      <!-- 元信息 -->
      <div class="asset-card__meta">
        <span class="text-text-tertiary text-xs">
          {{ formatSize(asset.size) }}
        </span>
        <span v-if="asset.format" class="text-text-tertiary text-xs">
          {{ asset.format.toUpperCase() }}
        </span>
      </div>

      <!-- 标签 -->
      <div v-if="asset.tags && asset.tags.length > 0" class="asset-card__tags">
        <GvSpace size="xs" wrap>
          <GvTag
            v-for="tag in asset.tags.slice(0, 2)"
            :key="tag"
            size="small"
            color="primary"
            variant="tonal"
          >
            {{ tag }}
          </GvTag>
          <GvTag
            v-if="asset.tags.length > 2"
            size="small"
            color="neutral"
            variant="tonal"
          >
            +{{ asset.tags.length - 2 }}
          </GvTag>
        </GvSpace>
      </div>
    </div>
  </GvCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'
import { VideoCamera, Picture, Headset, View, Edit, Delete } from '@element-plus/icons-vue'
import { GvCard, GvBadge, GvButton, GvSpace, GvTag, StatusBadge } from '@/components'
import type { AssetCardProps, AssetCardEmits } from './types'

const props = withDefaults(defineProps<AssetCardProps>(), {
  selectable: false,
  selected: false
})

const emit = defineEmits<AssetCardEmits>()

const cardClasses = computed(() => {
  return cn(
    'asset-card',
    'overflow-hidden transition-all duration-200',
    props.selected && 'ring-2 ring-primary-500'
  )
})

const thumbnailUrl = computed(() => {
  if (props.asset.type === 'image') {
    return props.asset.path
  }
  return null
})

function getTypeIcon(type: string) {
  const iconMap: Record<string, any> = {
    video: VideoCamera,
    image: Picture,
    audio: Headset
  }
  return iconMap[type] || Picture
}

function getTypeLabel(type: string) {
  const labelMap: Record<string, string> = {
    video: '视频',
    image: '图片',
    audio: '音频'
  }
  return labelMap[type] || type
}

function getTypeColor(type: string) {
  const colorMap: Record<string, string> = {
    video: 'primary',
    image: 'success',
    audio: 'warning'
  }
  return colorMap[type] || 'neutral'
}

function mapStatus(status: string): any {
  const statusMap: Record<string, string> = {
    ready: 'success',
    processing: 'processing',
    pending: 'pending',
    error: 'error'
  }
  return statusMap[status] || 'inactive'
}

function formatSize(size: number): string {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`
  if (size < 1024 * 1024 * 1024) return `${(size / (1024 * 1024)).toFixed(2)} MB`
  return `${(size / (1024 * 1024 * 1024)).toFixed(2)} GB`
}

function formatDuration(seconds: number): string {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = Math.floor(seconds % 60)
  if (h > 0) return `${h}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
  return `${m}:${s.toString().padStart(2, '0')}`
}

function handleClick() {
  emit('click', props.asset)
}

function handleView() {
  emit('view', props.asset)
}

function handleEdit() {
  emit('edit', props.asset)
}

function handleDelete() {
  emit('delete', props.asset)
}

function handleSelectToggle() {
  emit('select', props.asset, !props.selected)
}
</script>

<style scoped>
.asset-card {
  @apply relative cursor-pointer;
}

.asset-card__thumbnail {
  @apply relative w-full aspect-video bg-neutral-100 overflow-hidden;
}

.asset-card__checkbox {
  @apply absolute top-2 left-2 z-10;
}

.asset-card__checkbox :deep(.el-checkbox) {
  @apply bg-white rounded shadow-md;
}

.asset-card__image {
  @apply w-full h-full;
}

.asset-card__placeholder {
  @apply w-full h-full flex items-center justify-center bg-neutral-50;
}

.asset-card__type-badge {
  @apply absolute top-2 right-2 z-10;
}

.asset-card__duration {
  @apply absolute bottom-2 right-2 z-10;
  @apply px-2 py-1 bg-black/70 text-white text-xs rounded;
}

.asset-card__actions {
  @apply absolute inset-0 flex items-center justify-center;
  @apply bg-black/60 opacity-0 transition-opacity duration-200;
}

.asset-card:hover .asset-card__actions {
  @apply opacity-100;
}

.asset-card__info {
  @apply p-3 space-y-2;
}

.asset-card__name {
  @apply text-sm font-medium text-text-primary truncate;
}

.asset-card__meta {
  @apply flex items-center gap-2;
}

.asset-card__tags {
  @apply mt-2;
}

.asset-card__status {
  @apply mt-2;
}

.dark .asset-card__thumbnail {
  @apply bg-neutral-800;
}

.dark .asset-card__placeholder {
  @apply bg-neutral-900;
}

.dark .asset-card__name {
  @apply text-text-inverse;
}
</style>
