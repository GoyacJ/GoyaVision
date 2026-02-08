<template>
  <div 
    class="min-w-[120px] bg-white rounded-lg border-2 border-green-500 shadow-sm transition-all p-3"
    :class="{ 'ring-2 ring-green-100': selected }"
  >
    <div class="flex items-center gap-2">
      <div class="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center text-lg text-green-600">
        ▶
      </div>
      <div class="font-bold text-gray-700">开始</div>
    </div>
    <Handle type="source" position="bottom" class="w-3 h-3 !bg-green-500" />
  </div>
</template>

<script setup lang="ts">
import { Handle } from '@vue-flow/core'

defineProps<{
  selected: boolean
}>()
</script>
