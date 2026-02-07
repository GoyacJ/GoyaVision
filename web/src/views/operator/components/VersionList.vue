<template>
  <el-table :data="versions" v-loading="loading" size="small" border>
    <el-table-column prop="version" label="版本" min-width="120" />
    <el-table-column prop="exec_mode" label="执行模式" width="100" />
    <el-table-column prop="status" label="状态" width="100" />
    <el-table-column label="操作" width="240" fixed="right">
      <template #default="{ row }">
        <span v-if="row.status === 'active'" class="text-success-600 text-sm font-medium px-4">
          当前版本
        </span>
        <template v-else>
          <GvButton
            size="small"
            variant="text"
            @click="emit('activate', row.id)"
          >
            激活
          </GvButton>
          <GvButton
            size="small"
            variant="text"
            @click="emit('rollback', row.id)"
          >
            回滚
          </GvButton>
          <GvButton
            v-if="row.status !== 'archived'"
            size="small"
            variant="text"
            color="warning"
            @click="emit('archive', row.id)"
          >
            归档
          </GvButton>
        </template>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import type { OperatorVersion } from '@/api/operator'
import GvButton from '@/components/base/GvButton/index.vue'

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
