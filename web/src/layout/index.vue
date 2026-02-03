<template>
  <el-container class="layout-container">
    <el-header class="layout-header">
      <div class="header-left">
        <div class="logo">
          <span class="logo-text">GoyaVision</span>
        </div>
        <el-menu
          :default-active="activeMenu"
          mode="horizontal"
          router
          class="layout-menu"
          :ellipsis="false"
        >
          <sidebar-item
            v-for="menu in menuRoutes"
            :key="menu.id"
            :item="menu"
          />
        </el-menu>
      </div>
      <div class="header-right">
        <el-dropdown @command="handleCommand">
          <span class="user-info">
            <el-avatar :size="32" icon="UserFilled" />
            <span class="username">{{ userStore.nickname }}</span>
            <el-icon><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人信息</el-dropdown-item>
              <el-dropdown-item command="password">修改密码</el-dropdown-item>
              <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>
    <el-main class="layout-main">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </el-main>

    <el-dialog v-model="passwordDialogVisible" title="修改密码" width="400px">
      <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="80px">
        <el-form-item label="原密码" prop="old_password">
          <el-input v-model="passwordForm.old_password" type="password" show-password />
        </el-form-item>
        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="passwordForm.new_password" type="password" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirm_password">
          <el-input v-model="passwordForm.confirm_password" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="passwordDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="passwordLoading" @click="handleChangePassword">
          确定
        </el-button>
      </template>
    </el-dialog>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed, reactive, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { useUserStore } from '../store/user'
import { authApi } from '../api/auth'
import SidebarItem from './SidebarItem.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const passwordDialogVisible = ref(false)
const passwordLoading = ref(false)
const passwordFormRef = ref<FormInstance>()

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const validateConfirmPassword = (rule: any, value: string, callback: any) => {
  if (value !== passwordForm.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const passwordRules: FormRules = {
  old_password: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const menuRoutes = computed(() => {
  return userStore.menus.filter(menu => menu.visible)
})

const activeMenu = computed(() => {
  return route.path
})

function handleCommand(command: string) {
  switch (command) {
    case 'profile':
      ElMessage.info('个人信息功能开发中')
      break
    case 'password':
      passwordDialogVisible.value = true
      passwordForm.old_password = ''
      passwordForm.new_password = ''
      passwordForm.confirm_password = ''
      break
    case 'logout':
      handleLogout()
      break
  }
}

async function handleLogout() {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      type: 'warning'
    })
    await userStore.logout()
    ElMessage.success('已退出登录')
  } catch {
    // cancelled
  }
}

async function handleChangePassword() {
  if (!passwordFormRef.value) return

  await passwordFormRef.value.validate(async (valid) => {
    if (!valid) return

    passwordLoading.value = true
    try {
      await authApi.changePassword({
        old_password: passwordForm.old_password,
        new_password: passwordForm.new_password
      })
      ElMessage.success('密码修改成功')
      passwordDialogVisible.value = false
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || '密码修改失败')
    } finally {
      passwordLoading.value = false
    }
  })
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #ffffff;
}

.layout-header {
  height: 70px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  box-shadow: 0 4px 24px rgba(31, 38, 135, 0.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
  position: sticky;
  top: 0;
  z-index: 1000;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0;
  flex: 1;
  overflow: hidden;
}

.logo {
  display: flex;
  align-items: center;
  margin-right: 48px;
  flex-shrink: 0;
  cursor: pointer;
  transition: transform 0.3s;
}

.logo:hover {
  transform: scale(1.05);
}

.logo-text {
  font-size: 26px;
  font-weight: 800;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.5px;
}

.layout-menu {
  flex: 1;
  border: none;
  background-color: transparent;
}

.layout-menu :deep(.el-menu-item),
.layout-menu :deep(.el-sub-menu__title) {
  border-bottom: none;
  height: 70px;
  line-height: 70px;
  font-weight: 500;
  transition: all 0.3s;
  position: relative;
  margin: 0 4px;
  border-radius: 8px;
  display: flex;
  align-items: center;
}

.layout-menu :deep(.el-menu-item:hover),
.layout-menu :deep(.el-sub-menu__title:hover) {
  background: transparent;
  color: #667eea;
}

.layout-menu :deep(.el-menu-item.is-active) {
  color: #667eea;
  background: transparent;
  border-radius: 8px;
  font-weight: 600;
}

.layout-menu :deep(.el-menu-item.is-active::after) {
  content: '';
  position: absolute;
  bottom: 0;
  left: 20%;
  right: 20%;
  height: 3px;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border-radius: 2px;
}

.header-right {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  padding: 10px 16px;
  border-radius: 12px;
  transition: all 0.3s;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
}

.user-info:hover {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.15) 0%, rgba(118, 75, 162, 0.15) 100%);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

.user-info :deep(.el-avatar) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.username {
  color: #333;
  font-size: 14px;
  font-weight: 600;
}

.layout-main {
  flex: 1;
  background: transparent;
  padding: 32px;
  overflow-y: auto;
  overflow-x: hidden;
}

.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

:deep(.el-dropdown-menu) {
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(31, 38, 135, 0.15);
  border: 1px solid rgba(102, 126, 234, 0.1);
  padding: 8px;
}

:deep(.el-dropdown-menu__item) {
  border-radius: 8px;
  margin: 4px 0;
  transition: all 0.3s;
}

:deep(.el-dropdown-menu__item:hover) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  color: #667eea;
}

:deep(.el-dialog) {
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(31, 38, 135, 0.2);
}

:deep(.el-dialog__header) {
  padding: 24px 24px 16px;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.el-dialog__body) {
  padding: 24px;
}

@media (max-width: 768px) {
  .layout-header {
    padding: 0 16px;
    height: 60px;
  }

  .logo {
    margin-right: 24px;
  }

  .logo-text {
    font-size: 20px;
  }

  .layout-main {
    padding: 16px;
  }

  .username {
    display: none;
  }
}
</style>
