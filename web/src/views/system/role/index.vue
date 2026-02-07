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
        <el-switch
          v-model="row.status"
          :active-value="1"
          :inactive-value="0"
          :loading="row.statusLoading"
          @change="(val) => handleStatusChange(row, val)"
        />
      </template>

      <template #is_default="{ row }">
        <GvTag v-if="row.is_default" color="primary" size="small" variant="tonal">默认</GvTag>
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
        <el-form-item label="默认角色" prop="is_default">
          <el-switch v-model="form.is_default" />
          <div class="text-xs text-text-tertiary ml-2">新用户注册时将自动获得此角色</div>
        </el-form-item>
        <el-form-item label="自动分配" prop="auto_assign_provider">
          <GvSelect
            v-model="form.auto_assign_provider"
            :options="providerOptions"
            placeholder="选择登录方式（可选）"
            clearable
            class="w-full"
          />
          <div class="text-xs text-text-tertiary ml-2">用户通过此方式登录/绑定时自动获得此角色</div>
        </el-form-item>
      </el-form>
    </GvModal>

    <!-- 权限分配抽屉 -->
    <GvDrawer
      v-model="permissionDialogVisible"
      title="分配权限"
      size="medium"
      :confirm-loading="permissionLoading"
      @confirm="handlePermissionSubmit"
      @cancel="permissionDialogVisible = false"
    >
      <el-tabs v-model="activeTab" class="h-full flex flex-col">
        <el-tab-pane label="菜单权限" name="menus" class="h-full overflow-hidden">
          <div class="h-full overflow-y-auto pr-2">
            <el-tree
              ref="menuTreeRef"
              :data="menuTree"
              show-checkbox
              node-key="id"
              :default-checked-keys="selectedMenuIds"
              :props="{ label: 'name', children: 'children' }"
              highlight-current
            >
              <template #default="{ node, data }">
                <span class="custom-tree-node flex items-center gap-2">
                  <el-icon v-if="data.icon"><component :is="data.icon" /></el-icon>
                  <span>{{ node.label }}</span>
                  <GvTag v-if="data.type === 2" size="small" variant="tonal" class="ml-2">菜单</GvTag>
                  <GvTag v-else size="small" variant="tonal" color="info" class="ml-2">目录</GvTag>
                </span>
              </template>
            </el-tree>
          </div>
        </el-tab-pane>
        <el-tab-pane label="接口权限" name="permissions" class="h-full overflow-hidden">
          <div class="h-full overflow-y-auto pr-2">
            <el-tree
              ref="apiTreeRef"
              :data="apiTreeData"
              show-checkbox
              node-key="id"
              :default-checked-keys="selectedPermissionIds"
              :props="{ label: 'label', children: 'children' }"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </GvDrawer>
  </GvContainer>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import type { ElTree } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { roleApi, permissionApi, type Role, type Permission } from '../../../api/role'
import { menuApi, type Menu } from '../../../api/menu'
import GvContainer from '@/components/layout/GvContainer/index.vue'
import GvTable from '@/components/base/GvTable/index.vue'
import GvModal from '@/components/base/GvModal/index.vue'
import GvDrawer from '@/components/base/GvDrawer/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvSpace from '@/components/layout/GvSpace/index.vue'
import GvInput from '@/components/base/GvInput/index.vue'
import GvSelect from '@/components/base/GvSelect/index.vue'
import GvCard from '@/components/base/GvCard/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'
import PageHeader from '@/components/business/PageHeader/index.vue'
import StatusBadge from '@/components/business/StatusBadge/index.vue'
import type { TableColumn } from '@/components/base/GvTable/types'

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
const apiTreeRef = ref<InstanceType<typeof ElTree>>()

const form = reactive({
  code: '',
  name: '',
  description: '',
  status: 1,
  is_default: false,
  auto_assign_provider: ''
})

const providerOptions = [
  { label: 'Github', value: 'github' },
  { label: '微信', value: 'wechat' },
  { label: '手机号', value: 'phone' }
]

const apiTreeData = computed(() => {
  const root: Record<string, any> = {}
  permissions.value.forEach(p => {
    const parts = p.code.split(':')
    const module = parts[0]
    if (!root[module]) {
      root[module] = {
        id: `module:${module}`,
        label: module.charAt(0).toUpperCase() + module.slice(1) + ' 模块',
        children: []
      }
    }
    root[module].children.push({
      id: p.id,
      label: `${p.name} (${p.code})`
    })
  })
  return Object.values(root)
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
  { prop: 'is_default', label: '默认', width: '80' },
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

async function handleStatusChange(row: any, val: number | string | boolean) {
  if (!row || !row.id) return
  row.statusLoading = true
  try {
    await roleApi.update(row.id, {
      status: val as number
    })
    ElMessage.success('状态更新成功')
  } catch (error) {
    row.status = val === 1 ? 0 : 1
    ElMessage.error('状态更新失败')
  } finally {
    row.statusLoading = false
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
    status: 1,
    is_default: false,
    auto_assign_provider: ''
  })
  dialogVisible.value = true
}

function handleEdit(row: Role) {
  isEdit.value = true
  dialogTitle.value = '编辑角色'
  editingId.value = row.id
  const config = row.auto_assign_config
  const provider = config?.conditions?.provider || ''
  
  Object.assign(form, {
    code: row.code,
    name: row.name,
    description: row.description,
    status: row.status,
    is_default: row.is_default,
    auto_assign_provider: provider
  })
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const autoConfig = form.auto_assign_provider ? {
        trigger: 'login',
        conditions: { provider: form.auto_assign_provider }
      } : undefined

      if (isEdit.value) {
        await roleApi.update(editingId.value, {
          name: form.name,
          description: form.description,
          status: form.status,
          is_default: form.is_default,
          auto_assign_config: autoConfig
        })
        ElMessage.success('更新成功')
      } else {
        await roleApi.create({
          code: form.code,
          name: form.name,
          description: form.description,
          status: form.status,
          is_default: form.is_default,
          auto_assign_config: autoConfig
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
  // 列表数据不包含权限信息，需要单独获取详情
  permissionLoading.value = true
  try {
    const { data } = await roleApi.get(row.id)
    selectedMenuIds.value = data.menus?.map(m => m.id) || []
    selectedPermissionIds.value = data.permissions?.map(p => p.id) || []
    activeTab.value = 'menus'
    permissionDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取角色权限失败')
  } finally {
    permissionLoading.value = false
  }
}

async function handlePermissionSubmit() {
  permissionLoading.value = true
  try {
    const checkedMenuIds = menuTreeRef.value?.getCheckedKeys(false) as string[] || []
    const halfCheckedMenuIds = menuTreeRef.value?.getHalfCheckedKeys() as string[] || []
    const allMenuIds = [...checkedMenuIds, ...halfCheckedMenuIds]

    // 获取选中的接口权限 (排除模块节点)
    const checkedNodes = apiTreeRef.value?.getCheckedNodes(false, false) || []
    const apiIds = checkedNodes.filter((n: any) => !n.id.startsWith('module:')).map((n: any) => n.id)

    await roleApi.update(currentRoleId.value, {
      menu_ids: allMenuIds,
      permission_ids: apiIds
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
