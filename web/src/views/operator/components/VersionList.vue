<template>
  <el-table :data="versions" v-loading="loading" size="small" border>
    <el-table-column prop="version" label="版本" min-width="120" />
    <el-table-column prop="exec_mode" label="执行模式" width="100" />
    <el-table-column prop="status" label="状态" width="100" />
    <el-table-column label="操作" width="240" fixed="right">
      <template #default="{ row }">
        <el-button link type="primary" @click="emit('activate', row.id)">激活</el-button>
        <el-button link @click="emit('rollback', row.id)">回滚</el-button>
        <el-button link type="warning" @click="emit('archive', row.id)">归档</el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import type { OperatorVersion } from '@/api/operator'

defineProps<{
  versions: OperatorVersion[]
  loading?: boolean
}>()

const emit = defineEmits<{
  activate: [versionId: string]
  rollback: [versionId: string]
  archive: [versionId: string]
}>()
</script>
