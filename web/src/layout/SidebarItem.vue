<template>
  <template v-if="item.visible">
    <template v-if="!hasChildren">
      <el-menu-item :index="item.path">
        <el-icon v-if="item.icon">
          <component :is="item.icon" />
        </el-icon>
        <template #title>{{ item.name }}</template>
      </el-menu-item>
    </template>
    <el-sub-menu v-else :index="item.path">
      <template #title>
        <el-icon v-if="item.icon">
          <component :is="item.icon" />
        </el-icon>
        <span>{{ item.name }}</span>
      </template>
      <sidebar-item
        v-for="child in item.children"
        :key="child.id"
        :item="child"
      />
    </el-sub-menu>
  </template>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { MenuInfo } from '../api/auth'

const props = defineProps<{
  item: MenuInfo
}>()

const hasChildren = computed(() => {
  return props.item.children && props.item.children.length > 0
})
</script>
