<template>
  <div class="profile-page max-w-7xl mx-auto p-6 md:p-10">
    <!-- 顶部 Banner 重新设计 -->
    <div class="relative mb-12 group">
      <!-- 背景装饰 -->
      <div class="absolute inset-0 bg-gradient-to-r from-primary-600/10 via-purple-600/5 to-transparent rounded-[3rem] -z-10 blur-xl opacity-50 group-hover:opacity-100 transition-opacity duration-700"></div>
      
      <div class="bg-white rounded-[3rem] shadow-2xl shadow-neutral-200/50 border border-neutral-100 overflow-hidden relative">
        <!-- Banner 图 -->
        <div class="h-40 md:h-56 bg-neutral-900 relative overflow-hidden">
          <div class="absolute inset-0 opacity-40">
            <div class="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,_#3b82f6_0%,_transparent_70%)]"></div>
            <div class="absolute inset-0 bg-[url('https://www.transparenttextures.com/patterns/carbon-fibre.png')]"></div>
          </div>
          <!-- 动态光晕 -->
          <div class="absolute -top-24 -right-24 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl animate-pulse"></div>
        </div>

        <!-- 个人信息主体 -->
        <div class="px-8 md:px-12 pb-10 flex flex-col md:flex-row items-center md:items-end gap-8 relative z-10 -mt-20 md:-mt-24">
          <!-- 头像区域 -->
          <div class="relative group/avatar cursor-pointer" @click="handleTriggerUpload">
            <div class="relative p-1.5 bg-white rounded-[2.5rem] shadow-2xl">
              <el-avatar
                :size="isMobile ? 120 : 180"
                :src="userStore.userInfo?.avatar"
                class="!rounded-[2rem] bg-neutral-100 transition-transform duration-700 group-hover/avatar:scale-[1.02]"
              >
                <el-icon :size="60" class="text-neutral-300"><UserFilled /></el-icon>
              </el-avatar>
              <!-- 上传蒙层 -->
              <div class="absolute inset-1.5 rounded-[2rem] bg-neutral-900/40 backdrop-blur-sm flex items-center justify-center opacity-0 group-hover/avatar:opacity-100 transition-all duration-300">
                <el-icon color="white" :size="32"><CameraIcon /></el-icon>
              </div>
            </div>
            <el-upload
              ref="avatarUploadRef"
              class="hidden"
              action=""
              :show-file-list="false"
              :http-request="onUploadAvatar"
              accept="image/*"
            />
          </div>
          
          <!-- 文字信息 -->
          <div class="flex-1 text-center md:text-left mb-2">
            <div class="flex flex-col md:flex-row md:items-center gap-4 mb-3">
              <h1 class="text-4xl md:text-5xl font-black text-neutral-900 tracking-tighter">
                {{ userStore.userInfo?.nickname || userStore.username }}
              </h1>
              <div class="flex gap-2 justify-center">
                <GvTag effect="dark" type="primary" class="rounded-full px-4 font-black text-[10px] uppercase tracking-wider">
                  {{ userStore.roles[0] || 'Member' }}
                </GvTag>
                <GvTag type="success" class="rounded-full px-4 font-black text-[10px] uppercase tracking-wider">
                  Pro Active
                </GvTag>
              </div>
            </div>
            <div class="flex flex-wrap items-center justify-center md:justify-start gap-6 text-neutral-400 font-bold text-sm">
              <span class="flex items-center gap-1.5 hover:text-neutral-600 transition-colors cursor-default">
                <el-icon><UserIcon /></el-icon> @{{ userStore.username }}
              </span>
              <span class="flex items-center gap-1.5 hover:text-neutral-600 transition-colors cursor-default">
                <el-icon><CalendarIcon /></el-icon> 加入于 {{ displayDate(userStore.userInfo?.created_at) }}
              </span>
            </div>
          </div>

          <!-- 快速统计 -->
          <div class="hidden lg:flex gap-10 pb-4">
             <div class="text-center">
               <p class="text-[10px] text-neutral-400 uppercase font-black tracking-widest mb-1">算子调用</p>
               <p class="text-2xl font-black text-neutral-800">12.4k</p>
             </div>
             <div class="w-px h-10 bg-neutral-100"></div>
             <div class="text-center">
               <p class="text-[10px] text-neutral-400 uppercase font-black tracking-widest mb-1">账户积分</p>
               <p class="text-2xl font-black text-neutral-800">5,820</p>
             </div>
          </div>
        </div>
      </div>
    </div>

    <div class="flex flex-col lg:flex-row gap-10">
      <!-- 侧边导航 -->
      <div class="w-full lg:w-80 flex-shrink-0">
        <div class="bg-white rounded-[2.5rem] shadow-xl shadow-neutral-200/40 border border-neutral-100 p-4 sticky top-24">
          <div class="px-4 py-6">
            <p class="text-[10px] text-neutral-300 font-black uppercase tracking-[0.2em] mb-4">账户管理</p>
            <nav class="space-y-1.5">
              <button
                v-for="item in menuTabs.slice(0, 2)"
                :key="item.key"
                class="w-full flex items-center gap-4 px-5 py-4 rounded-2xl cursor-pointer transition-all duration-300 group relative overflow-hidden"
                :class="activeTab === item.key ? 'bg-neutral-900 text-white shadow-xl shadow-neutral-300' : 'hover:bg-neutral-50 text-neutral-500'"
                @click="activeTab = item.key"
              >
                <el-icon :size="20" :class="activeTab === item.key ? 'text-primary-400' : 'group-hover:text-neutral-900'">
                  <component :is="item.icon" />
                </el-icon>
                <span class="font-bold tracking-tight">{{ item.label }}</span>
                <el-icon v-if="activeTab === item.key" class="ml-auto animate-bounce-x"><ArrowRightIcon /></el-icon>
              </button>
            </nav>

            <p class="text-[10px] text-neutral-300 font-black uppercase tracking-[0.2em] mt-10 mb-4">资产与使用</p>
            <nav class="space-y-1.5">
              <button
                v-for="item in menuTabs.slice(2)"
                :key="item.key"
                class="w-full flex items-center gap-4 px-5 py-4 rounded-2xl cursor-pointer transition-all duration-300 group"
                :class="activeTab === item.key ? 'bg-neutral-900 text-white shadow-xl shadow-neutral-300' : 'hover:bg-neutral-50 text-neutral-500'"
                @click="activeTab = item.key"
              >
                <el-icon :size="20" :class="activeTab === item.key ? 'text-primary-400' : 'group-hover:text-neutral-900'">
                  <component :is="item.icon" />
                </el-icon>
                <span class="font-bold tracking-tight">{{ item.label }}</span>
                <el-icon v-if="activeTab === item.key" class="ml-auto animate-bounce-x"><ArrowRightIcon /></el-icon>
              </button>
            </nav>
          </div>
        </div>
      </div>

      <!-- 内容区 -->
      <div class="flex-1 min-w-0">
        <transition name="page-fade" mode="out-in">
          <component 
            :is="currentComponent" 
            :key="activeTab"
            v-bind="componentProps"
            @save="handleSaveProfile"
            @change-password="showPasswordModal = true"
          />
        </transition>
      </div>
    </div>

    <!-- 修改密码弹窗 - 保持原有逻辑但优化样式 -->
    <el-dialog v-model="showPasswordModal" title="更新安全密码" width="480px" class="rounded-custom-dialog" :show-close="false">
      <div class="p-4">
        <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-position="top">
          <GvInput v-model="passwordForm.old_password" type="password" show-password label="当前正在使用的密码" size="large" class="mb-6" />
          <div class="my-8 border-t border-neutral-100 relative">
             <span class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-white px-4 text-[10px] font-black text-neutral-300 uppercase tracking-widest">设置新密码</span>
          </div>
          <GvInput v-model="passwordForm.new_password" type="password" show-password label="设置新登录密码" size="large" class="mb-6" />
          <GvInput v-model="passwordForm.confirm_password" type="password" show-password label="确认新登录密码" size="large" />
        </el-form>
      </div>
      <template #footer>
        <div class="flex gap-4 justify-end p-4">
          <GvButton @click="showPasswordModal = false" size="large" class="!rounded-2xl px-8">取消</GvButton>
          <GvButton type="primary" size="large" class="!rounded-2xl px-10 !bg-neutral-900 font-black border-none" @click="confirmChangePassword">确认更新</GvButton>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch, defineAsyncComponent } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { 
  User as UserIcon, 
  Lock as LockIcon, 
  UserFilled, 
  Camera as CameraIcon, 
  Calendar as CalendarIcon,
  ArrowRight as ArrowRightIcon,
  CreditCard as PaymentIcon,
  Star as PointsIcon,
  Medal as SubscriptionIcon,
  DataLine as UsageIcon
} from '@element-plus/icons-vue'

import { useUserStore } from '@/store/user'
import { authApi } from '@/api/auth'
import { fileApi } from '@/api/file'
import { useBreakpoint } from '@/composables/useBreakpoint'

// 基础组件
import GvTag from '@/components/base/GvTag/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvInput from '@/components/base/GvInput/index.vue'

// 异步加载业务组件
const ProfileOverview = defineAsyncComponent(() => import('./components/ProfileOverview.vue'))
const SecuritySettings = defineAsyncComponent(() => import('./components/SecuritySettings.vue'))
const PaymentManager = defineAsyncComponent(() => import('./components/PaymentManager.vue'))
const PointsManager = defineAsyncComponent(() => import('./components/PointsManager.vue'))
const SubscriptionManager = defineAsyncComponent(() => import('./components/SubscriptionManager.vue'))
const UsageStats = defineAsyncComponent(() => import('./components/UsageStats.vue'))

const userStore = useUserStore()
const route = useRoute()
const { isMobile } = useBreakpoint()
const activeTab = ref('profile')

const menuTabs = [
  { key: 'profile', label: '个人资料', icon: UserIcon, component: ProfileOverview },
  { key: 'security', label: '账号安全', icon: LockIcon, component: SecuritySettings },
  { key: 'payment', label: '支付管理', icon: PaymentIcon, component: PaymentManager },
  { key: 'points', label: '积分管理', icon: PointsIcon, component: PointsManager },
  { key: 'subscription', label: '订阅管理', icon: SubscriptionIcon, component: SubscriptionManager },
  { key: 'usage', label: '使用记录', icon: UsageIcon, component: UsageStats }
]

const currentComponent = computed(() => {
  return menuTabs.find(t => t.key === activeTab.value)?.component || ProfileOverview
})

const componentProps = computed(() => {
  if (activeTab.value === 'profile') {
    return {
      initialData: {
        nickname: userStore.userInfo?.nickname || '',
        email: userStore.userInfo?.email || '',
        phone: userStore.userInfo?.phone || ''
      },
      username: userStore.username,
      lastSyncTime: displayFullTime(userStore.userInfo?.updated_at),
      loading: submitting.value
    }
  }
  if (activeTab.value === 'security') {
    return {
      locationInfo: '中国, 上海市',
      deviceInfo: 'Chrome Browser • macOS Sonoma'
    }
  }
  return {}
})

const avatarUploadRef = ref()
const submitting = ref(false)
const showPasswordModal = ref(false)

// 密码表单逻辑保持不变
const passwordFormRef = ref<FormInstance>()
const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const checkConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (value !== passwordForm.new_password) callback(new Error('两次输入的密码不一致'))
  else callback()
}

const passwordRules: FormRules = {
  old_password: [{ required: true, message: '必须输入原密码', trigger: 'blur' }],
  new_password: [{ required: true, message: '新密码不能为空', trigger: 'blur' }, { min: 6, message: '长度需不少于 6 位', trigger: 'blur' }],
  confirm_password: [{ required: true, message: '请再次输入新密码', trigger: 'blur' }, { validator: checkConfirmPassword, trigger: 'blur' }]
}

onMounted(async () => {
  try {
    await userStore.getProfile()
  } catch (e) {}
  if (route.query.tab) activeTab.value = route.query.tab as string
})

watch(() => route.query.tab, (newTab) => {
  if (newTab) activeTab.value = newTab as string
})

function handleTriggerUpload() {
  const input = avatarUploadRef.value?.$el?.querySelector('input')
  if (input) input.click()
}

async function onUploadAvatar(options: any) {
  try {
    const res = await fileApi.upload(options.file)
    await authApi.updateProfile({ avatar: res.data.url })
    ElMessage.success('头像更新成功')
    await userStore.getProfile()
  } catch (error) {
    ElMessage.error('上传过程中出现错误')
  }
}

async function handleSaveProfile(formData: any) {
  submitting.value = true
  try {
    await authApi.updateProfile(formData)
    ElMessage.success('个人资料已同步')
    await userStore.getProfile()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '更新失败')
  } finally {
    submitting.value = false
  }
}

async function confirmChangePassword() {
  if (!passwordFormRef.value) return
  await passwordFormRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      await authApi.changePassword(passwordForm)
      ElMessage.success('安全密码已更新')
      showPasswordModal.value = false
      Object.keys(passwordForm).forEach(k => (passwordForm as any)[k] = '')
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || '密码验证失败')
    }
  })
}

function displayDate(date: any) {
  if (!date) return '...'
  return new Date(date).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long' })
}

function displayFullTime(date: any) {
  if (!date) return '从未同步'
  return new Date(date).toLocaleString('zh-CN')
}
</script>

<style scoped>
.profile-page {
  animation: page-in 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes page-in {
  from { opacity: 0; transform: translateY(30px); filter: blur(10px); }
  to { opacity: 1; transform: translateY(0); filter: blur(0); }
}

.page-fade-enter-active,
.page-fade-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.page-fade-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.page-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

@keyframes bounce-x {
  0%, 100% { transform: translateX(0); }
  50% { transform: translateX(4px); }
}
.animate-bounce-x {
  animation: bounce-x 1s infinite;
}

:deep(.rounded-custom-dialog) {
  @apply rounded-[3rem] overflow-hidden border-none shadow-2xl !p-6;
}
:deep(.rounded-custom-dialog .el-dialog__header) {
  @apply pb-0 text-2xl font-black px-4 pt-4;
}
</style>
