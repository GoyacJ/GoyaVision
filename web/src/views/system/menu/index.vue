<template>
  <GvContainer max-width="full">
    <!-- 页面头部 -->
    <PageHeader
      title="菜单管理"
      description="管理系统菜单，配置路由和权限"
    >
      <template #actions>
        <GvButton @click="handleAdd()" v-permission="'menu:create'">
          <template #icon>
            <el-icon><Plus /></el-icon>
          </template>
          新增菜单
        </GvButton>
      </template>
    </PageHeader>

    <!-- 树形表格 -->
    <GvCard>
      <el-table
        :data="tableData"
        v-loading="loading"
        row-key="id"
        :tree-props="{ children: 'children' }"
        stripe
      >
        <el-table-column prop="name" label="菜单名称" min-width="150" />
        <el-table-column prop="code" label="编码" width="150" />
        <el-table-column prop="icon" label="图标" width="80">
          <template #default="{ row }">
            <el-icon v-if="row.icon"><component :is="row.icon" /></el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="80">
          <template #default="{ row }">
            <GvTag :color="getTypeColor(row.type)" size="small">
              {{ getTypeName(row.type) }}
            </GvTag>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路由路径" width="150" />
        <el-table-column prop="component" label="组件" width="150" />
        <el-table-column prop="permission" label="权限标识" width="120" />
        <el-table-column prop="sort" label="排序" width="60" />
        <el-table-column prop="visible" label="显示" width="60">
          <template #default="{ row }">
            {{ row.visible ? '是' : '否' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <StatusBadge :status="row.status === 1 ? 'enabled' : 'disabled'" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <GvSpace size="xs">
              <GvButton size="small" variant="tonal" @click="handleAdd(row.id)" v-permission="'menu:create'">
                新增
              </GvButton>
              <GvButton size="small" @click="handleEdit(row)" v-permission="'menu:update'">
                编辑
              </GvButton>
              <GvButton
                size="small"
                variant="text"
                @click="handleDelete(row)"
                v-permission="'menu:delete'"
              >
                删除
              </GvButton>
            </GvSpace>
          </template>
        </el-table-column>
      </el-table>
    </GvCard>

    <!-- 新增/编辑对话框 -->
    <GvModal
      v-model="dialogVisible"
      :title="dialogTitle"
      size="large"
      :confirm-loading="submitLoading"
      @confirm="handleSubmit"
      @cancel="dialogVisible = false"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="上级菜单" prop="parent_id">
          <el-tree-select
            v-model="form.parent_id"
            :data="menuTreeOptions"
            :props="{ label: 'name', children: 'children', value: 'id' }"
            check-strictly
            clearable
            placeholder="选择上级菜单（可选）"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="菜单类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio :value="1">目录</el-radio>
            <el-radio :value="2">菜单</el-radio>
            <el-radio :value="3">按钮</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="编码" prop="code">
          <GvInput v-model="form.code" :disabled="isEdit" placeholder="唯一标识" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <GvInput v-model="form.name" placeholder="显示名称" />
        </el-form-item>
        <el-form-item v-if="form.type !== 3" label="图标" prop="icon">
          <GvInput v-model="form.icon" placeholder="Element Plus 图标名称" />
        </el-form-item>
        <el-form-item v-if="form.type !== 3" label="路由路径" prop="path">
          <GvInput v-model="form.path" placeholder="/system/user" />
        </el-form-item>
        <el-form-item v-if="form.type === 2" label="组件路径" prop="component">
          <GvInput v-model="form.component" placeholder="system/user/index" />
        </el-form-item>
        <el-form-item label="权限标识" prop="permission">
          <GvInput v-model="form.permission" placeholder="user:list" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item v-if="form.type !== 3" label="是否显示" prop="visible">
          <el-switch v-model="form.visible" />
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
import { menuApi, type Menu } from '../../../api/menu'
import GvContainer from '@/components/layout/GvContainer/index.vue'
import GvCard from '@/components/base/GvCard/index.vue'
import GvModal from '@/components/base/GvModal/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvSpace from '@/components/layout/GvSpace/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'
import GvInput from '@/components/base/GvInput/index.vue'
import PageHeader from '@/components/business/PageHeader/index.vue'
import StatusBadge from '@/components/business/StatusBadge/index.vue'

const loading = ref(false)
const tableData = ref<Menu[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref('')

const form = reactive({
  parent_id: undefined as string | undefined,
  code: '',
  name: '',
  type: 2,
  path: '',
  icon: '',
  component: '',
  permission: '',
  sort: 0,
  visible: true,
  status: 1
})

const rules: FormRules = {
  code: [
    { required: true, message: '请输入菜单编码', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入菜单名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择菜单类型', trigger: 'change' }
  ]
}

const menuTreeOptions = computed(() => {
  return [{ id: '', name: '根目录', children: tableData.value }]
})

function getTypeName(type: number): string {
  switch (type) {
    case 1: return '目录'
    case 2: return '菜单'
    case 3: return '按钮'
    default: return '未知'
  }
}

function getTypeColor(type: number): string {
  switch (type) {
    case 1: return 'neutral'
    case 2: return 'success'
    case 3: return 'warning'
    default: return 'info'
  }
}

async function loadData() {
  loading.value = true
  try {
    const response = await menuApi.listTree()
    tableData.value = response.data || []
  } catch (error: any) {
    ElMessage.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

function handleAdd(parentId?: string) {
  isEdit.value = false
  dialogTitle.value = '新增菜单'
  editingId.value = ''
  Object.assign(form, {
    parent_id: parentId || undefined,
    code: '',
    name: '',
    type: 2,
    path: '',
    icon: '',
    component: '',
    permission: '',
    sort: 0,
    visible: true,
    status: 1
  })
  dialogVisible.value = true
}

function handleEdit(row: Menu) {
  isEdit.value = true
  dialogTitle.value = '编辑菜单'
  editingId.value = row.id
  Object.assign(form, {
    parent_id: row.parent_id || undefined,
    code: row.code,
    name: row.name,
    type: row.type,
    path: row.path,
    icon: row.icon,
    component: row.component,
    permission: row.permission,
    sort: row.sort,
    visible: row.visible,
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
      const data = {
        parent_id: form.parent_id || undefined,
        code: form.code,
        name: form.name,
        type: form.type,
        path: form.path,
        icon: form.icon,
        component: form.component,
        permission: form.permission,
        sort: form.sort,
        visible: form.visible,
        status: form.status
      }

      if (isEdit.value) {
        await menuApi.update(editingId.value, data)
        ElMessage.success('更新成功')
      } else {
        await menuApi.create(data)
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

async function handleDelete(row: Menu) {
  try {
    await ElMessageBox.confirm(`确定删除菜单 "${row.name}" 吗？`, '提示', {
      type: 'warning'
    })
    await menuApi.delete(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch {
    // cancelled
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
</style>
