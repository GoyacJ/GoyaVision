<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>角色管理</span>
          <el-button type="primary" @click="handleAdd" v-permission="'role:create'">
            <el-icon><Plus /></el-icon>
            新增角色
          </el-button>
        </div>
      </template>

      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="code" label="角色编码" width="150" />
        <el-table-column prop="name" label="角色名称" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleEdit(row)" v-permission="'role:update'">
              编辑
            </el-button>
            <el-button link type="primary" size="small" @click="handlePermission(row)" v-permission="'role:update'">
              权限
            </el-button>
            <el-button
              link
              type="danger"
              size="small"
              @click="handleDelete(row)"
              v-permission="'role:delete'"
              :disabled="row.code === 'super_admin'"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="编码" prop="code">
          <el-input v-model="form.code" :disabled="isEdit" placeholder="唯一标识，如 admin" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="显示名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="permissionDialogVisible" title="分配权限" width="600px">
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
      <template #footer>
        <el-button @click="permissionDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="permissionLoading" @click="handlePermissionSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import type { ElTree } from 'element-plus'
import { roleApi, permissionApi, type Role, type Permission } from '../../../api/role'
import { menuApi, type Menu } from '../../../api/menu'

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
.page-container {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.permission-item {
  padding: 4px 0;
}
</style>
