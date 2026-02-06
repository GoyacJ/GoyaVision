<template>
  <div>
    <div class="mb-2 flex justify-between items-center">
      <div class="text-sm font-medium text-text-primary">依赖管理</div>
      <el-button size="small" @click="handleAdd">新增依赖</el-button>
    </div>

    <el-table :data="localDeps" v-loading="loading" size="small" border>
      <el-table-column label="依赖算子ID" min-width="220">
        <template #default="{ row }">
          <el-input v-model="row.depends_on_id" placeholder="请输入依赖算子 ID" />
        </template>
      </el-table-column>
      <el-table-column label="最低版本" width="140">
        <template #default="{ row }">
          <el-input v-model="row.min_version" placeholder="如 v1.2.0" />
        </template>
      </el-table-column>
      <el-table-column label="可选" width="100">
        <template #default="{ row }">
          <el-switch v-model="row.is_optional" />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100">
        <template #default="{ $index }">
          <el-button link type="danger" @click="handleRemove($index)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="mt-3">
      <el-button type="primary" size="small" @click="emit('save', localDeps)">保存依赖</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

type DepItem = {
  depends_on_id: string
  min_version?: string
  is_optional?: boolean
}

const props = defineProps<{
  dependencies: DepItem[]
  loading?: boolean
}>()

const emit = defineEmits<{
  save: [deps: DepItem[]]
}>()

const localDeps = ref<DepItem[]>([])

watch(
  () => props.dependencies,
  (value) => {
    localDeps.value = (value || []).map((d) => ({ ...d }))
  },
  { immediate: true, deep: true }
)

function handleAdd() {
  localDeps.value.push({ depends_on_id: '', min_version: '', is_optional: false })
}

function handleRemove(index: number) {
  localDeps.value.splice(index, 1)
}
</script>
