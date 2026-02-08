<template>
  <GvContainer>
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-text-primary">系统配置</h1>
    </div>

    <GvCard padding="lg">
      <el-form ref="formRef" :model="form" label-width="120px" class="max-w-2xl">
        <el-form-item label="默认首页" prop="home_path">
          <GvInput v-model="form.home_path" placeholder="/assets" />
          <div class="form-tip">未登录用户访问根路径时将重定向至此路径</div>
        </el-form-item>

        <el-form-item label="公开菜单" prop="public_menus">
          <el-tree-select
            v-model="form.public_menus"
            :data="menuTree"
            multiple
            show-checkbox
            check-strictly
            node-key="id"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            placeholder="选择公开菜单"
            class="w-full"
          />
          <div class="form-tip">选中的菜单将对未登录用户可见（需确保对应页面支持公开访问）</div>
        </el-form-item>

        <el-form-item>
          <GvButton type="primary" :loading="saving" @click="handleSave">保存配置</GvButton>
        </el-form-item>
      </el-form>
    </GvCard>
  </GvContainer>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import GvContainer from '@/components/layout/GvContainer/index.vue'
import GvCard from '@/components/base/GvCard/index.vue'
import GvInput from '@/components/base/GvInput/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import { systemApi } from '@/api/system'
import { menuApi } from '@/api/menu'
import { useAppStore } from '@/store/app'

const form = reactive({
  home_path: '',
  public_menus: [] as string[]
})

const menuTree = ref<any[]>([])
const saving = ref(false)
const appStore = useAppStore()

onMounted(async () => {
  await loadMenus()
  await loadConfig()
})

async function loadMenus() {
  try {
    // Try listTree if available, otherwise list
    // Assuming backend /menus returns tree or flat list that el-tree-select can handle if tree
    // If flat, we might need to build tree, but let's assume /menus returns tree structure for now as it's standard in this project's other parts (e.g. role menu assignment)
    const res = await menuApi.list()
    menuTree.value = res.data
  } catch (error) {
    console.error(error)
  }
}

async function loadConfig() {
  try {
    const res = await systemApi.getPublicConfig()
    form.home_path = res.data.home_path
    
    // Extract IDs from public_menus response
    const ids: string[] = []
    function extractIds(nodes: any[]) {
      for (const node of nodes) {
        ids.push(node.id)
        if (node.children) extractIds(node.children)
      }
    }
    extractIds(res.data.public_menus || [])
    form.public_menus = ids
  } catch (error) {
    console.error(error)
  }
}

async function handleSave() {
  saving.value = true
  try {
    await systemApi.updateConfig({
      'system.home_path': form.home_path,
      'system.public_menus': form.public_menus
    })
    ElMessage.success('保存成功')
    // Refresh local store
    appStore.fetchPublicConfig()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '保存失败')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.form-tip {
  font-size: 12px;
  color: #9ca3af;
  margin-top: 4px;
  line-height: 1.4;
}
</style>
