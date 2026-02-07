<template>
  <GvContainer max-width="full">
    <PageHeader
      title="租户管理"
      description="管理多租户数据隔离，创建和维护不同的租户环境"
    >
      <template #actions>
        <GvButton @click="openCreateDialog">
          <template #icon>
            <el-icon><Plus /></el-icon>
          </template>
          添加租户
        </GvButton>
      </template>
    </PageHeader>

    <GvTable
      :data="tenants"
      :columns="columns"
      :loading="loading"
      border
      stripe
    >
      <template #status="{ row }">
        <GvTag :color="row.status === 1 ? 'success' : 'neutral'" size="small">
          {{ row.status === 1 ? '启用' : '禁用' }}
        </GvTag>
      </template>
      <template #created_at="{ row }">
        {{ formatDate(row.created_at) }}
      </template>
      <template #actions="{ row }">
        <GvSpace size="xs">
          <GvButton size="small" variant="text" @click="handleEdit(row)">
            编辑
          </GvButton>
          <GvButton size="small" variant="text" color="error" @click="handleDelete(row)">
            删除
          </GvButton>
        </GvSpace>
      </template>
    </GvTable>

    <GvModal
      v-model="showDialog"
      :title="isEdit ? '编辑租户' : '添加租户'"
      size="medium"
      confirm-text="确定"
      :confirm-loading="submitting"
      @confirm="handleSubmit"
    >
      <div class="space-y-4">
        <GvInput v-model="form.name" label="租户名称" placeholder="请输入租户名称" />
        <GvInput v-model="form.code" label="租户标识" placeholder="请输入唯一标识" :disabled="isEdit" />
        <GvSelect
          v-if="isEdit"
          v-model="form.status"
          label="状态"
          :options="statusOptions"
          class="w-full"
        />
      </div>
    </GvModal>
  </GvContainer>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { tenantApi, type Tenant } from '@/api/tenant'
import { useAsyncData } from '@/composables'
import GvContainer from '@/components/layout/GvContainer/index.vue'
import GvTable from '@/components/base/GvTable/index.vue'
import GvModal from '@/components/base/GvModal/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvInput from '@/components/base/GvInput/index.vue'
import GvSelect from '@/components/base/GvSelect/index.vue'
import GvSpace from '@/components/layout/GvSpace/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'
import PageHeader from '@/components/business/PageHeader/index.vue'

const { data: tenantsData, isLoading: loading, execute: refreshTable } = useAsyncData(
  () => tenantApi.list(),
  { immediate: true }
)

const tenants = computed(() => tenantsData.value?.data || [])

const showDialog = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const editId = ref('')

const form = reactive({
  name: '',
  code: '',
  status: 1
})

const statusOptions = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 }
]

const columns = [
  { prop: 'name', label: '名称', minWidth: '140' },
  { prop: 'code', label: '标识', width: '120' },
  { prop: 'status', label: '状态', width: '100' },
  { prop: 'created_at', label: '创建时间', width: '180' },
  { prop: 'actions', label: '操作', width: '160', fixed: 'right' }
]

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

function openCreateDialog() {
  isEdit.value = false
  form.name = ''
  form.code = ''
  form.status = 1
  showDialog.value = true
}

function handleEdit(row: Tenant) {
  isEdit.value = true
  editId.value = row.id
  form.name = row.name
  form.code = row.code
  form.status = row.status
  showDialog.value = true
}

async function handleSubmit() {
  if (!form.name || !form.code) {
    ElMessage.warning('请填写名称和标识')
    return
  }

  submitting.value = true
  try {
    if (isEdit.value) {
      await tenantApi.update(editId.value, { name: form.name, status: form.status })
      ElMessage.success('更新成功')
    } else {
      await tenantApi.create(form)
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    refreshTable()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(row: Tenant) {
  try {
    await ElMessageBox.confirm('确定要删除此租户吗？', '提示', { type: 'warning' })
    await tenantApi.delete(row.id)
    ElMessage.success('删除成功')
    refreshTable()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

import { computed } from 'vue'
</script>
