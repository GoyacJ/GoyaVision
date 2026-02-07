<template>
  <div>
    <div class="mb-2 flex justify-between items-center">
      <div class="text-sm font-medium text-text-primary">依赖管理</div>
      <el-button size="small" @click="handleAdd">新增依赖</el-button>
    </div>

    <el-table :data="localDeps" v-loading="loading" size="small" border>
      <el-table-column label="依赖算子" min-width="220">
        <template #default="{ row }">
          <GvSelect
            v-model="row.depends_on_id"
            placeholder="搜索算子"
            filterable
            remote
            :remote-method="handleSearchOperators"
            :loading="searching"
            :options="operatorOptions"
            label-key="name"
            value-key="id"
            @change="(val) => handleOperatorChange(row, val)"
          />
        </template>
      </el-table-column>
      <el-table-column label="最低版本" width="140">
        <template #default="{ row }">
          <GvSelect
            v-model="row.min_version"
            placeholder="选择版本"
            :options="versionOptionsMap[row.depends_on_id] || []"
            :disabled="!row.depends_on_id"
          />
        </template>
      </el-table-column>
      <el-table-column label="可选" width="100">
        <template #default="{ row }">
          <el-switch v-model="row.is_optional" />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100">
        <template #default="{ $index }">
          <GvButton variant="text" color="error" size="small" @click="handleRemove($index)">删除</GvButton>
        </template>
      </el-table-column>
    </el-table>

    <div class="mt-3 flex justify-end">
      <GvButton variant="filled" color="primary" size="small" @click="emit('save', localDeps)">保存依赖</GvButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, reactive } from 'vue'
import { operatorApi } from '@/api/operator'
import GvSelect from '@/components/base/GvSelect/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'

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
const operatorOptions = ref<Array<{ id: string; name: string }>>([])
const searching = ref(false)
const versionOptionsMap = reactive<Record<string, Array<{ label: string; value: string }>>>({})

watch(
  () => props.dependencies,
  (value) => {
    localDeps.value = (value || []).map((d) => ({ ...d }))
    // 加载已有依赖的版本列表
    localDeps.value.forEach(d => {
      if (d.depends_on_id) {
        fetchVersions(d.depends_on_id)
      }
    })
  },
  { immediate: true, deep: true }
)

onMounted(() => {
  handleSearchOperators('')
})

async function handleSearchOperators(query: string) {
  searching.value = true
  try {
    const res = await operatorApi.list({ keyword: query, page: 1, page_size: 20 })
    operatorOptions.value = (res.data?.items || []).map((op: any) => ({
      id: op.id,
      name: `${op.name} (${op.code})`
    }))
  } catch (e) {
    console.error(e)
  } finally {
    searching.value = false
  }
}

async function fetchVersions(operatorId: string) {
  if (!operatorId || versionOptionsMap[operatorId]) return
  try {
    const res = await operatorApi.listVersions(operatorId, { page: 1, page_size: 100 })
    versionOptionsMap[operatorId] = (res.data?.items || []).map((v: any) => ({
      label: v.version,
      value: v.version
    }))
  } catch (e) {
    console.error('fetch versions failed', e)
  }
}

function handleOperatorChange(row: DepItem, val: any) {
  row.min_version = '' // 重置版本
  if (val) {
    fetchVersions(val as string)
  }
}

function handleAdd() {
  localDeps.value.push({ depends_on_id: '', min_version: '', is_optional: false })
}

function handleRemove(index: number) {
  localDeps.value.splice(index, 1)
}
</script>
