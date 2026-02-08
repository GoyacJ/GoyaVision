<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
    <!-- 当前订阅卡片 -->
    <GvCard shadow="none" bordered padding="lg" class="!rounded-[2.5rem] bg-white relative overflow-hidden">
      <div class="absolute top-0 right-0 p-8">
         <div class="w-32 h-32 bg-primary-50 rounded-full flex items-center justify-center opacity-20 -mr-16 -mt-16 scale-150">
            <el-icon :size="60" class="text-primary-600"><MedalIcon /></el-icon>
         </div>
      </div>
      
      <div class="relative z-10 flex flex-col md:flex-row md:items-center gap-10">
        <div class="flex-1">
          <div class="flex items-center gap-3 mb-4">
            <h3 class="text-3xl font-black text-neutral-900">{{ currentPlan }}订阅</h3>
            <GvTag type="primary" size="small" class="rounded-full px-4 font-bold">按月计费</GvTag>
          </div>
          <p class="text-neutral-500 font-medium leading-relaxed mb-6 max-w-md">
            {{ getPlanDesc(currentPlan) }}
          </p>
          <div class="flex items-center gap-6">
            <div class="flex flex-col">
              <span class="text-xs text-neutral-400 font-bold uppercase tracking-wider mb-1">当前状态</span>
              <span class="text-neutral-800 font-black">Active</span>
            </div>
            <div class="w-px h-8 bg-neutral-100"></div>
            <div class="flex flex-col">
              <span class="text-xs text-neutral-400 font-bold uppercase tracking-wider mb-1">会员等级</span>
              <span class="text-emerald-600 font-black">{{ memberLevel }}</span>
            </div>
          </div>
        </div>
        <div class="flex flex-col gap-3 min-w-[200px]">
          <GvButton type="primary" class="!rounded-2xl !h-12 font-black w-full shadow-lg shadow-primary-100">管理自动续费</GvButton>
          <GvButton class="!rounded-2xl !h-12 font-black w-full" @click="scrollToPlans">变更订阅计划</GvButton>
        </div>
      </div>
    </GvCard>

    <!-- 计划对比 -->
    <div id="plans-section" class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div 
        v-for="plan in plans" 
        :key="plan.name"
        class="p-8 rounded-[2.5rem] border-2 flex flex-col transition-all duration-300 relative overflow-hidden group"
        :class="currentPlan === plan.name ? 'border-primary-600 bg-primary-50/10' : 'border-neutral-100 bg-white hover:border-neutral-200'"
      >
        <div v-if="plan.popular" class="absolute top-0 right-0 bg-primary-600 text-white text-[10px] font-black px-4 py-1 rounded-bl-xl uppercase tracking-tighter">
          最受欢迎
        </div>
        
        <h4 class="text-xl font-black text-neutral-800 mb-2">{{ plan.name }}</h4>
        <div class="flex items-baseline gap-1 mb-6">
          <span class="text-3xl font-black tracking-tight">¥{{ plan.price }}</span>
          <span class="text-neutral-400 font-medium text-sm">/月</span>
        </div>

        <ul class="space-y-4 mb-10 flex-1">
          <li v-for="feat in plan.features" :key="feat" class="flex items-center gap-3 text-sm font-medium text-neutral-600">
            <el-icon class="text-emerald-500"><CheckIcon /></el-icon>
            {{ feat }}
          </li>
        </ul>

        <GvButton 
          :type="currentPlan === plan.name ? 'default' : 'primary'" 
          :disabled="currentPlan === plan.name"
          :loading="submitting && targetPlan === plan.name"
          class="!rounded-2xl !h-12 font-black w-full"
          @click="handleSubscribe(plan.name)"
        >
          {{ currentPlan === plan.name ? '当前计划' : '立即升级' }}
        </GvButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Medal as MedalIcon,
  Check as CheckIcon
} from '@element-plus/icons-vue'
import GvCard from '@/components/base/GvCard/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'
import { userAssetApi } from '@/api/user-asset'

const currentPlan = ref('Free')
const memberLevel = ref('Free')
const submitting = ref(false)
const targetPlan = ref('')

const plans = [
  {
    name: 'Free',
    price: 0,
    features: ['5GB 存储空间', '每日 100 次算子调用', '基础客服支持', '1 个并发任务'],
    popular: false
  },
  {
    name: 'Pro',
    price: 99,
    features: ['500GB 存储空间', '无限次算子调用', '24/7 优先支持', '10 个并发任务', 'AI 算子优先权'],
    popular: true
  },
  {
    name: 'Enterprise',
    price: 499,
    features: ['无限制存储空间', '定制化算子开发', '独立部署选项', '无限并发任务', '专人技术顾问'],
    popular: false
  }
]

onMounted(() => {
  fetchSummary()
})

async function fetchSummary() {
  try {
    const res = await userAssetApi.getSummary()
    if (res.data && res.data.data) {
      currentPlan.value = res.data.data.subscription || 'Free'
      memberLevel.value = res.data.data.member_level
    }
  } catch (error) {}
}

function getPlanDesc(plan: string) {
  if (plan === 'Pro') return '您当前正在享受专业版带来的全部功能，包括无限次算子调用、优先排队以及 500GB 的云端媒体存储空间。'
  if (plan === 'Enterprise') return '您当前正在享受企业版全量能力，包含定制化算子开发、独占计算资源及专人技术顾问支持。'
  return '您当前处于基础版计划，每日有固定的算子调用配额。升级到专业版可解锁无限次调用及更多存储空间。'
}

function scrollToPlans() {
  document.getElementById('plans-section')?.scrollIntoView({ behavior: 'smooth' })
}

async function handleSubscribe(planName: string) {
  try {
    await ElMessageBox.confirm(`确认将订阅计划变更为 ${planName} 吗？`, '变更确认', {
      confirmButtonText: '确定变更',
      cancelButtonText: '取消',
      type: 'info'
    })
    
    targetPlan.value = planName
    submitting.value = true
    await userAssetApi.subscribe(planName)
    ElMessage.success('订阅计划已更新')
    fetchSummary()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '操作失败')
    }
  } finally {
    submitting.value = false
    targetPlan.value = ''
  }
}
</script>
