<template>
  <div class="login-container">
    <div class="bg-decoration">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
      <div class="circle circle-3"></div>
    </div>
    <div class="login-box">
      <div class="login-header">
        <div class="logo-icon">
          <el-icon :size="48"><VideoCameraFilled /></el-icon>
        </div>
        <h1>GoyaVision</h1>
        <p>{{ isRegister ? '创建您的账号' : 'AI 多媒体资源分析处理平台' }}</p>
        <div class="header-divider"></div>
      </div>
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        class="login-form"
        @keyup.enter="handleSubmit"
      >
        <el-form-item prop="username">
          <gv-input
            v-model="formData.username"
            placeholder="用户名"
            size="large"
            prefix-icon="User"
            clearable
          />
        </el-form-item>
        <el-form-item prop="password">
          <gv-input
            v-model="formData.password"
            type="password"
            placeholder="密码"
            size="large"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <!-- 注册模式下的额外字段 -->
        <template v-if="isRegister">
          <el-form-item prop="nickname">
            <gv-input
              v-model="formData.nickname"
              placeholder="昵称"
              size="large"
              prefix-icon="Edit"
              clearable
            />
          </el-form-item>
          <el-form-item prop="email">
            <gv-input
              v-model="formData.email"
              placeholder="邮箱 (可选)"
              size="large"
              prefix-icon="Message"
              clearable
            />
          </el-form-item>
        </template>

        <el-form-item>
          <gv-button
            size="large"
            :loading="loading"
            block
            class="submit-btn"
            @click="handleSubmit"
          >
            {{ isRegister ? '立即注册' : '立即登录' }}
          </gv-button>
        </el-form-item>
      </el-form>

      <div class="mode-switch">
        <el-button link type="primary" @click="toggleMode">
          {{ isRegister ? '已有账号？立即登录' : '没有账号？申请注册' }}
        </el-button>
      </div>

      <template v-if="!isRegister">
        <div class="mt-6 mb-4 flex items-center gap-3">
          <div class="h-px bg-gray-200 flex-1"></div>
          <span class="text-xs text-gray-400">其他登录方式</span>
          <div class="h-px bg-gray-200 flex-1"></div>
        </div>
        <div class="flex justify-center gap-6 mb-2">
          <div 
            v-for="item in loginMethods" 
            :key="item.key"
            class="w-10 h-10 rounded-full bg-gray-50 flex items-center justify-center cursor-pointer hover:bg-gray-100 transition-colors border border-gray-200"
            @click="handleOtherLogin(item.key)"
            :title="item.label"
          >
            <span class="text-sm font-bold text-gray-600">{{ item.icon }}</span>
          </div>
        </div>
      </template>
    </div>
    <div class="login-footer-text">
      <p>© 2026 GoyaVision. Powered by AI</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useUserStore } from '@/store/user'
import GvInput from '@/components/base/GvInput/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import { authApi } from '@/api/auth'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const isRegister = ref(false)

const formData = reactive({
  username: '',
  password: '',
  nickname: '',
  email: ''
})

onMounted(async () => {
  const code = route.query.code as string
  if (code) {
    const provider = localStorage.getItem('oauth_provider')
    if (provider) {
      loading.value = true
      try {
        await userStore.loginOAuth({
          provider,
          code
        })
        ElMessage.success('登录成功')
        localStorage.removeItem('oauth_provider')
        
        const redirect = route.query.redirect as string
        const targetPath = redirect || '/assets'
        
        await router.push(targetPath)
      } catch (error: any) {
        ElMessage.error(error.response?.data?.message || '登录失败')
        localStorage.removeItem('oauth_provider')
      } finally {
        loading.value = false
      }
    }
  }
})

const formRules = computed<FormRules>(() => {
  const rules: FormRules = {
    username: [
      { required: true, message: '请输入用户名', trigger: 'blur' },
      { min: 4, message: '用户名至少4位', trigger: 'blur' }
    ],
    password: [
      { required: true, message: '请输入密码', trigger: 'blur' },
      { min: 6, message: '密码长度至少6位', trigger: 'blur' }
    ]
  }

  if (isRegister.value) {
    rules.nickname = [
      { required: true, message: '请输入昵称', trigger: 'blur' }
    ]
    rules.email = [
      { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
    ]
  }

  return rules
})

const loginMethods = [
  { key: 'github', label: 'Github', icon: 'G' },
  { key: 'wechat', label: '微信', icon: 'W' },
  { key: 'phone', label: '手机号', icon: 'P' }
]

function toggleMode() {
  isRegister.value = !isRegister.value
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

function handleOtherLogin(type: string) {
  if (type === 'phone') {
    ElMessage.info('手机号登录功能开发中')
    return
  }
  
  localStorage.setItem('oauth_provider', type)
  // Redirect to OAuth authorize endpoint
  window.location.href = `/api/v1/auth/oauth/login?provider=${type}`
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      if (isRegister.value) {
        await authApi.register({
          username: formData.username,
          password: formData.password,
          nickname: formData.nickname,
          email: formData.email
        })
        ElMessage.success('注册成功，请登录')
        isRegister.value = false
      } else {
        await userStore.login({
          username: formData.username,
          password: formData.password
        })
        ElMessage.success('登录成功')
        
        const redirect = route.query.redirect as string
        const targetPath = redirect || '/assets'
        
        await router.push(targetPath)
      }
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || (isRegister.value ? '注册失败' : '登录失败'))
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
  padding: 20px;
}

.bg-decoration {
  position: absolute;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 0;
}

.circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 20s infinite ease-in-out;
}

.circle-1 {
  width: 400px;
  height: 400px;
  top: -200px;
  left: -200px;
  animation-delay: 0s;
}

.circle-2 {
  width: 300px;
  height: 300px;
  bottom: -150px;
  right: -150px;
  animation-delay: 5s;
}

.circle-3 {
  width: 250px;
  height: 250px;
  top: 50%;
  right: 10%;
  animation-delay: 10s;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
  }
}

.login-box {
  width: 440px;
  max-width: 100%;
  padding: 48px 40px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  position: relative;
  z-index: 1;
  animation: fadeInUp 0.6s ease-out;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
  }
  50% {
    transform: scale(1.05);
    box-shadow: 0 12px 32px rgba(102, 126, 234, 0.6);
  }
}

.login-header h1 {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 12px 0;
  letter-spacing: -0.5px;
}

.login-header p {
  font-size: 14px;
  color: #666;
  margin: 0;
  font-weight: 400;
}

.header-divider {
  width: 60px;
  height: 4px;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border-radius: 2px;
  margin: 20px auto 0;
}

.login-form {
  margin-top: 32px;
}

.login-form :deep(.el-form-item) {
  margin-bottom: 24px;
}

.submit-btn {
  width: 100%;
  height: 48px !important;
  font-size: 16px !important;
  font-weight: 600 !important;
  border-radius: 12px !important;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  border: none !important;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4) !important;
  transition: all 0.3s !important;
  margin-top: 8px;
  color: white !important;
}

.submit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.5) !important;
}

.submit-btn:active {
  transform: translateY(0);
}

.mode-switch {
  text-align: center;
  margin-top: 12px;
}

.login-footer-text {
  margin-top: 32px;
  text-align: center;
  color: rgba(255, 255, 255, 0.8);
  font-size: 13px;
  z-index: 1;
}

.login-footer-text p {
  margin: 0;
}

@media (max-width: 768px) {
  .login-box {
    width: 100%;
    padding: 32px 24px;
  }

  .login-header h1 {
    font-size: 28px;
  }

  .logo-icon {
    width: 64px;
    height: 64px;
  }

  .logo-icon :deep(.el-icon) {
    font-size: 36px;
  }
}
</style>
