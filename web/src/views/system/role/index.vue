<template>
  <GvContainer max-width="full">
    <!-- 页面头部 -->
    <PageHeader
      title="角色管理"
      description="管理系统角色，配置菜单和接口权限"
    >
      <template #actions>
        <GvButton @click="handleAdd" v-permission="'role:create'">
          <template #icon>
            <el-icon><Plus /></el-icon>
          </template>
          新增角色
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
    >
      <template #status="{ row }">
        <StatusBadge :status="row.status === 1 ? 'enabled' : 'disabled'" />
      </template>
      
      <template #created_at="{ row }">
        {{ formatDate(row.created_at) }}
      </template>
      
      <template #actions="{ row }">
        <GvSpace size="xs">
          <GvButton
            size="small"
            @click="handleEdit(row)"
            v-permission="'role:update'"
          >
            编辑
          </GvButton>
          <GvButton
            size="small"
            variant="tonal"
            @click="handlePermission(row)"
            v-permission="'role:update'"
          >
            权限
          </GvButton>
          <GvButton
            size="small"
            variant="text"
            @click="handleDelete(row)"
            v-permission="'role:delete'"
            :disabled="row.code === 'super_admin'"
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
        <el-form-item label="编码" prop="code">
          <GvInput v-model="form.code" :disabled="isEdit" placeholder="唯一标识，如 admin" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <GvInput v-model="form.name" placeholder="显示名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <GvInput v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
    </GvModal>

    <!-- 权限分配对话框 -->
    <GvModal
      v-model="permissionDialogVisible"
      title="分配权限"
      size="large"
      :confirm-loading="permissionLoading"
      @confirm="handlePermissionSubmit"
      @cancel="permissionDialogVisible = false"
    >
      <el-tabs v-model="activeTab">
        <el-tab-pane label="菜单权限" name="menus">
          <el-tree
            ref="menuTreeRef"
            :data="menuTree"
            show-checkbox
            node-key="id"
            :default-checked-keys="selectedMenuIds"
            :props="{ label: 'name', children: 'children' }"
          />
        </el-tab-pane>
        <el-tab-pane label="接口权限" name="permissions">
          <el-checkbox-group v-model="selectedPermissionIds">
            <div v-for="permission in permissions" :key="permission.id" class="permission-item">
              <el-checkbox :value="permission.id">
                {{ permission.name }} ({{ permission.code }})
              </el-checkbox>
            </div>
          </el-checkbox-group>
        </el-tab-pane>
      </el-tabs>
    </GvModal>
  </GvContainer>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import type { ElTree } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { roleApi, permissionApi, type Role, type Permission } from '../../../api/role'
import { menuApi, type Menu } from '../../../api/menu'
import {
  GvContainer,
  GvTable,
  GvModal,
  GvButton,
  GvSpace,
  GvInput,
  PageHeader,
  StatusBadge,
  type TableColumn
} from '@/components'

const loading = ref(false)
const tableData = ref<Role[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref('')

const permissionDialogVisible = ref(false)
const permissionLoading = ref(false)
const activeTab = ref('menus')
const menuTree = ref<Menu[]>([])
const permissions = ref<Permission[]>([])
const selectedMenuIds = ref<string[]>([])
const selectedPermissionIds = ref<string[]>([])
const currentRoleId = ref('')
const menuTreeRef = ref<InstanceType<typeof ElTree>>()

const form = reactive({
  code: '',
  name: '',
  description: '',
  status: 1
})

const rules: FormRules = {
  code: [
    { required: true, message: '请输入角色编码', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' }
  ]
}

const columns: TableColumn[] = [
  { prop: 'code', label: '角色编码', width: '150' },
  { prop: 'name', label: '角色名称', width: '150' },
  { prop: 'description', label: '描述', minWidth: '200' },
  { prop: 'status', label: '状态', width: '100' },
  { prop: 'created_at', label: '创建时间', width: '180' },
  { prop: 'actions', label: '操作', width: '200', fixed: 'right' }
]

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleString()
}

async function loadData() {
  loading.value = true
  try {
    const response = await roleApi.list()
    tableData.value = response.data || []
  } catch (error: any) {
    ElMessage.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function loadMenuTree() {
  try {
    const response = await menuApi.listTree()
    menuTree.value = response.data || []
  } catch (error) {
    console.error('加载菜单失败', error)
  }
}

async function loadPermissions() {
  try {
    const response = await permissionApi.list()
    permissions.value = response.data || []
  } catch (error) {
    console.error('加载权限失败', error)
  }
}

function handleAdd() {
  isEdit.value = false
  dialogTitle.value = '新增角色'
  editingId.value = ''
  Object.assign(form, {
    code: '',
    name: '',
    description: '',
    status: 1
  })
  dialogVisible.value = true
}

function handleEdit(row: Role) {
  isEdit.value = true
  dialogTitle.value = '编辑角色'
  editingId.value = row.id
  Object.assign(form, {
    code: row.code,
    name: row.name,
    description: row.description,
    status: row.status
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
        await roleApi.update(editingId.value, {
          name: form.name,
          description: form.description,
          status: form.status
        })
        ElMessage.success('更新成功')
      } else {
        await roleApi.create({
          code: form.code,
          name: form.name,
          description: form.description,
          status: form.status
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

async function handlePermission(row: Role) {
  currentRoleId.value = row.id
  selectedMenuIds.value = row.menus?.map(m => m.id) || []
  selectedPermissionIds.value = row.permissions?.map(p => p.id) || []
  activeTab.value = 'menus'
  permissionDialogVisible.value = true
}

async function handlePermissionSubmit() {
  permissionLoading.value = true
  try {
    const checkedMenuIds = menuTreeRef.value?.getCheckedKeys(false) as string[] || []
    const halfCheckedMenuIds = menuTreeRef.value?.getHalfCheckedKeys() as string[] || []
    const allMenuIds = [...checkedMenuIds, ...halfCheckedMenuIds]

    await roleApi.update(currentRoleId.value, {
      menu_ids: allMenuIds,
      permission_ids: selectedPermissionIds.value
    })
    ElMessage.success('权限分配成功')
    permissionDialogVisible.value = false
    loadData()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  } finally {
    permissionLoading.value = false
  }
}

async function handleDelete(row: Role) {
  try {
    await ElMessageBox.confirm(`确定删除角色 "${row.name}" 吗？`, '提示', {
      type: 'warning'
    })
    await roleApi.delete(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch {
    // cancelled
  }
}

onMounted(() => {
  loadData()
  loadMenuTree()
  loadPermissions()
})
</script>

<style scoped>
.permission-item {
  @apply py-1;
}
</style>
