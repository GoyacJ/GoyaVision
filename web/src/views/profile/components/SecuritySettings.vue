<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
    <GvCard shadow="none" bordered padding="lg" class="!rounded-[2.5rem]">
      <template #header>
        <div class="flex items-center gap-4">
          <div class="w-1.5 h-8 bg-purple-600 rounded-full"></div>
          <h3 class="text-2xl font-black text-neutral-900 tracking-tight">安全中心</h3>
        </div>
      </template>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-8 mt-4">
        <!-- 密码修改卡片 -->
        <div class="group p-8 rounded-[2rem] border border-neutral-100 bg-neutral-50/50 hover:bg-white hover:shadow-2xl hover:shadow-purple-100 transition-all duration-500">
          <div class="w-16 h-16 rounded-2xl bg-purple-100 text-purple-600 flex items-center justify-center text-3xl mb-8 group-hover:scale-110 group-hover:rotate-6 transition-transform">
            <el-icon><LockIcon /></el-icon>
          </div>
          <h4 class="text-xl font-black text-neutral-800 mb-3">登录密码</h4>
          <p class="text-sm text-neutral-400 mb-8 leading-relaxed font-medium">定期更新您的登录密码，建议使用包含字母、数字与符号的复杂组合。</p>
          <GvButton 
            type="primary" 
            class="!rounded-xl !h-12 w-full font-black !bg-neutral-900 hover:!bg-black !border-none transition-all"
            @click="$emit('change-password')"
          >
            立即修改
          </GvButton>
        </div>

        <!-- MFA 卡片 -->
        <div class="p-8 rounded-[2rem] border border-neutral-100 bg-neutral-50/30 opacity-60">
          <div class="w-16 h-16 rounded-2xl bg-amber-100 text-amber-600 flex items-center justify-center text-3xl mb-8">
            <el-icon><KeyIcon /></el-icon>
          </div>
          <h4 class="text-xl font-black text-neutral-800 mb-3">两步验证 (MFA)</h4>
          <p class="text-sm text-neutral-400 mb-8 leading-relaxed font-medium">在登录时增加额外的动态确认步，为您的账号提供金融级防护。</p>
          <GvButton disabled class="!rounded-xl !h-12 w-full font-black">即将上线</GvButton>
        </div>
      </div>
    </GvCard>

    <!-- 活跃会话 -->
    <GvCard shadow="none" bordered padding="lg" class="!rounded-[2.5rem]">
      <template #header>
        <div class="flex items-center justify-between">
          <h4 class="text-lg font-black text-neutral-800">当前活跃会话</h4>
          <GvTag type="success" size="small" class="rounded-full">运行中</GvTag>
        </div>
      </template>
      <div class="flex items-center gap-6 p-6 bg-primary-50/30 rounded-3xl border border-primary-100/50 mt-2">
        <div class="w-14 h-14 rounded-2xl bg-white shadow-sm flex items-center justify-center text-primary-500 text-2xl">
          <el-icon><MonitorIcon /></el-icon>
        </div>
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-3">
            <span class="text-lg font-black text-neutral-700">Web Dashboard</span>
            <span class="px-2 py-0.5 bg-green-500 text-white text-[10px] font-black rounded uppercase tracking-wider">当前</span>
          </div>
          <div class="text-sm text-neutral-400 mt-1 font-medium">
            {{ locationInfo }} • {{ deviceInfo }}
          </div>
        </div>
      </div>
    </GvCard>
  </div>
</template>

<script setup lang="ts">
import { 
  Lock as LockIcon, 
  Key as KeyIcon, 
  Monitor as MonitorIcon 
} from '@element-plus/icons-vue'
import GvCard from '@/components/base/GvCard/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'

defineProps<{
  locationInfo: string
  deviceInfo: string
}>()

defineEmits<{
  (e: 'change-password'): void
}>()
</script>
