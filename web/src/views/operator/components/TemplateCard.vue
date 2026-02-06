<template>
  <GvCard shadow="sm" padding="md" class="h-full">
    <div class="flex items-start justify-between gap-3">
      <div>
        <div class="text-base font-semibold text-text-primary">{{ template.name }}</div>
        <div class="mt-1 text-xs text-text-tertiary">{{ template.code }}</div>
      </div>
      <GvTag size="small" variant="tonal">{{ template.exec_mode?.toUpperCase() }}</GvTag>
    </div>

    <div class="mt-3 text-sm text-text-secondary line-clamp-3">
      {{ template.description || '暂无描述' }}
    </div>

    <div class="mt-3 flex items-center gap-2">
      <GvTag size="small" variant="tonal">{{ template.category }}</GvTag>
      <GvTag size="small" variant="tonal">{{ template.type }}</GvTag>
    </div>

    <div class="mt-4 flex justify-end gap-2">
      <GvButton size="small" variant="tonal" @click="emit('preview', template)">预览</GvButton>
      <GvButton size="small" @click="emit('install', template)">安装</GvButton>
    </div>
  </GvCard>
</template>

<script setup lang="ts">
import type { OperatorTemplate } from '@/api/operator'
import GvCard from '@/components/base/GvCard/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'

defineProps<{
  template: OperatorTemplate
}>()

const emit = defineEmits<{
  install: [tpl: OperatorTemplate]
  preview: [tpl: OperatorTemplate]
}>()
</script>
