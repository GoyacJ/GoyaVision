<template>
  <GvContainer max-width="full">
    <!-- 页面头部 -->
    <PageHeader
      title="用户管理"
      description="管理系统用户，分配角色和权限"
    >
      <template #actions>
        <GvButton @click="handleAdd" v-permission="'user:create'">
          <template #icon>
            <el-icon><Plus /></el-icon>
          </template>
          新增用户
        </GvButton>
      </template>
    </PageHeader>

    <!-- 数据表格 -->
    <GvTable
      :data="tableData"
      :columns="columns"
      :loading="loading"
      border
      stripe
      pagination
      :pagination-config="paginationConfig"
      @current-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <template #roles="{ row }">
        <GvSpace wrap size="xs">
          <GvTag v-for="role in row.roles" :key="role.id" size="small">
            {{ role.name }}
          </GvTag>
        </GvSpace>
      </template>
      
      <template #status="{ row }">
        <el-switch
          v-model="row.status"
          :active-value="1"
          :inactive-value="0"
          :loading="row.statusLoading"
          @change="(val) => handleStatusChange(row, val)"
        />
      </template>
      
      <template #created_at="{ row }">
        {{ formatDate(row.created_at) }}
      </template>
      
      <template #actions="{ row }">
        <GvSpace size="xs">
          <GvButton
            size="small"
            @click="handleEdit(row)"
            v-permission="'user:update'"
          >
            编辑
          </GvButton>
          <GvButton
            size="small"
            variant="tonal"
            @click="handleResetPassword(row)"
            v-permission="'user:update'"
          >
            重置密码
          </GvButton>
          <GvButton
            size="small"
            variant="text"
            @click="handleDelete(row)"
            v-permission="'user:delete'"
          >
            删除
          </GvButton>
        </GvSpace>
      </template>
    </GvTable>

    <!-- 新增/编辑对话框 -->
    <GvModal
      v-model="dialogVisible"
      :title="dialogTitle"
      :confirm-loading="submitLoading"
      @confirm="handleSubmit"
      @cancel="dialogVisible = false"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <GvInput v-model="form.username" :disabled="isEdit" />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="密码" prop="password">
          <GvInput v-model="form.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <GvInput v-model="form.nickname" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <GvInput v-model="form.email" type="email" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <GvInput v-model="form.phone" type="tel" />
        </el-form-item>
        <el-form-item label="角色" prop="role_ids">
          <GvSelect
            v-model="form.role_ids"
            :options="roleOptions"
            multiple
            placeholder="选择角色"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
    </GvModal>
  </GvContainer>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { userApi, type User } from '../../../api/user'
import { roleApi, type Role } from '../../../api/role'
import GvContainer from '@/components/layout/GvContainer/index.vue'
import GvTable from '@/components/base/GvTable/index.vue'
import GvModal from '@/components/base/GvModal/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvSpace from '@/components/layout/GvSpace/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'
import GvInput from '@/components/base/GvInput/index.vue'
import GvSelect from '@/components/base/GvSelect/index.vue'
import PageHeader from '@/components/business/PageHeader/index.vue'
import StatusBadge from '@/components/business/StatusBadge/index.vue'
import type { TableColumn } from '@/components/base/GvTable/types'

const loading = ref(false)
const tableData = ref<User[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const roles = ref<Role[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref('')

const form = reactive({
  username: '',
  password: '',
  nickname: '',
  email: '',
  phone: '',
  status: 1,
  role_ids: [] as string[]
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ]
}

const columns: TableColumn[] = [
  { prop: 'username', label: '用户名', width: '120' },
  { prop: 'nickname', label: '昵称', width: '120' },
  { prop: 'email', label: '邮箱', width: '180' },
  { prop: 'phone', label: '手机号', width: '120' },
  { prop: 'roles', label: '角色', minWidth: '150' },
  { prop: 'status', label: '状态', width: '100' },
  { prop: 'created_at', label: '创建时间', width: '180' },
  { prop: 'actions', label: '操作', width: '240', fixed: 'right' }
]

const roleOptions = computed(() => 
  roles.value.map(r => ({ label: r.name, value: r.id }))
)

const paginationConfig = computed(() => ({
  currentPage: currentPage.value,
  pageSize: pageSize.value,
  total: total.value
}))

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleString()
}

async function loadData() {
  loading.value = true
  try {
    const response = await userApi.list({
      limit: pageSize.value,
      offset: (currentPage.value - 1) * pageSize.value
    })
    tableData.value = response.data.items || []
    total.value = response.data.total
  } catch (error: any) {
    ElMessage.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function loadRoles() {
  try {
    const response = await roleApi.list({ status: 1 })
    roles.value = response.data || []
  } catch (error) {
    console.error('加载角色失败', error)
  }
}

function handlePageChange(page: number) {
  currentPage.value = page
  loadData()
}

function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  loadData()
}

async function handleStatusChange(row: any, val: number | string | boolean) {
  if (!row || !row.id) return
  row.statusLoading = true
  try {
    await userApi.update(row.id, {
      status: val as number
    })
    ElMessage.success('状态更新成功')
  } catch (error) {
    row.status = val === 1 ? 0 : 1 // Revert
    ElMessage.error('状态更新失败')
  } finally {
    row.statusLoading = false
  }
}

function handleAdd() {
  isEdit.value = false
  dialogTitle.value = '新增用户'
  editingId.value = ''
  Object.assign(form, {
    username: '',
    password: '',
    nickname: '',
    email: '',
    phone: '',
    status: 1,
    role_ids: []
  })
  dialogVisible.value = true
}

function handleEdit(row: User) {
  isEdit.value = true
  dialogTitle.value = '编辑用户'
  editingId.value = row.id
  Object.assign(form, {
    username: row.username,
    password: '',
    nickname: row.nickname,
    email: row.email,
    phone: row.phone,
    status: row.status,
    role_ids: row.roles?.map(r => r.id) || []
  })
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      if (isEdit.value) {
        await userApi.update(editingId.value, {
          nickname: form.nickname,
          email: form.email,
          phone: form.phone,
          status: form.status,
          role_ids: form.role_ids
        })
        ElMessage.success('更新成功')
      } else {
        await userApi.create({
          username: form.username,
          password: form.password,
          nickname: form.nickname,
          email: form.email,
          phone: form.phone,
          status: form.status,
          role_ids: form.role_ids
        })
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadData()
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || '操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

async function handleResetPassword(row: User) {
  try {
    const { value } = await ElMessageBox.prompt('请输入新密码', '重置密码', {
      inputType: 'password',
      inputValidator: (val) => {
        if (!val || val.length < 6) {
          return '密码长度至少6位'
        }
        return true
      }
    })
    await userApi.resetPassword(row.id, value)
    ElMessage.success('密码重置成功')
  } catch {
    // cancelled
  }
}

async function handleDelete(row: User) {
  try {
    await ElMessageBox.confirm(`确定删除用户 "${row.username}" 吗？`, '提示', {
      type: 'warning'
    })
    await userApi.delete(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch {
    // cancelled
  }
}

onMounted(() => {
  loadData()
  loadRoles()
})
</script>

<style scoped>
</style>
