<template>
  <el-container class="layout-container">
    <el-header class="layout-header">
      <div class="header-left">
        <!-- Mobile Toggle -->
        <div v-if="isMobile" class="mobile-toggle" @click="mobileMenuVisible = true">
          <el-icon :size="24"><Menu /></el-icon>
        </div>

        <div class="logo">
          <span class="logo-text">GoyaVision</span>
        </div>
        
        <!-- Desktop Menu -->
        <el-menu
          v-if="!isMobile"
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
        <template v-if="userStore.isLoggedIn">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32" :src="userStore.avatar">
                <el-icon><UserFilled /></el-icon>
              </el-avatar>
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
        </template>
        <template v-else>
          <el-button type="primary" link @click="handleLogin">登录</el-button>
        </template>
      </div>
    </el-header>
    <el-main class="layout-main">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </el-main>

    <!-- Mobile Navigation Drawer -->
    <el-drawer
      v-model="mobileMenuVisible"
      direction="ltr"
      size="240px"
      :with-header="false"
      class="mobile-nav-drawer"
    >
      <div class="mobile-menu-container">
        <div class="mobile-menu-header">
          <span class="logo-text">GoyaVision</span>
        </div>
        <el-menu
          :default-active="activeMenu"
          mode="vertical"
          router
          class="mobile-menu"
        >
          <sidebar-item
            v-for="menu in menuRoutes"
            :key="menu.id"
            :item="menu"
          />
        </el-menu>
      </div>
    </el-drawer>

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
import { useAppStore } from '@/store/app'
import { authApi } from '../api/auth'
import { useBreakpoint } from '@/composables/useBreakpoint'
import SidebarItem from './SidebarItem.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const appStore = useAppStore()
const { isMobile } = useBreakpoint()

const mobileMenuVisible = ref(false)
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
  if (userStore.isLoggedIn) {
    return userStore.menus.filter(menu => menu.visible)
  }
  return appStore.publicMenus.filter(menu => menu.visible)
})

const activeMenu = computed(() => {
  return route.path
})

function handleLogin() {
  router.push(`/login?redirect=${route.path}`)
}

function handleCommand(command: string) {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'password':
      // 跳转到个人中心的安全设置标签
      router.push({ path: '/profile', query: { tab: 'security' } })
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
  background: #ffffff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
  position: sticky;
  top: 0;
  z-index: 1000;
  border-bottom: 1px solid #E5E5E5;
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
}

.logo-text {
  font-size: 26px;
  font-weight: 700;
  color: #262626;
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
  transition: color 150ms;
  position: relative;
  margin: 0 4px;
  border-radius: 6px;
  display: flex;
  align-items: center;
}

.layout-menu :deep(.el-menu-item:hover),
.layout-menu :deep(.el-sub-menu__title:hover) {
  background: transparent;
  color: #4F5B93;
}

.layout-menu :deep(.el-menu-item.is-active) {
  color: #4F5B93;
  background: transparent;
  border-radius: 6px;
  font-weight: 600;
}

.layout-menu :deep(.el-menu-item.is-active::after) {
  content: '';
  position: absolute;
  bottom: 0;
  left: 20%;
  right: 20%;
  height: 2px;
  background: #4F5B93;
  border-radius: 1px;
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
  border-radius: 6px;
  transition: background-color 150ms;
  background: #F5F5F5;
}

.user-info:hover {
  background: #E5E5E5;
}

.user-info :deep(.el-avatar) {
  background: #4F5B93;
}

.username {
  color: #262626;
  font-size: 14px;
  font-weight: 500;
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
  transition: opacity 150ms ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

:deep(.el-dropdown-menu) {
  border-radius: 8px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #E5E5E5;
  padding: 8px;
}

:deep(.el-dropdown-menu__item) {
  border-radius: 6px;
  margin: 4px 0;
  transition: background-color 150ms;
}

:deep(.el-dropdown-menu__item:hover) {
  background: #F5F5F5;
  color: #4F5B93;
}

:deep(.el-dialog) {
  border-radius: 8px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
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

.mobile-toggle {
  display: flex;
  align-items: center;
  margin-right: 16px;
  cursor: pointer;
  color: #606266;
}

.mobile-menu-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.mobile-menu-header {
  height: 60px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  border-bottom: 1px solid #f0f0f0;
}

.mobile-menu {
  flex: 1;
  border-right: none;
  overflow-y: auto;
}
</style>
