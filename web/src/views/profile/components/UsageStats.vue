<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
    <!-- 核心统计指标 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div v-for="stat in mainStats" :key="stat.label" class="p-8 rounded-[2.5rem] border border-neutral-100 bg-white hover:shadow-2xl hover:shadow-neutral-100 transition-all duration-500 group">
        <div :class="['w-14 h-14 rounded-2xl flex items-center justify-center text-2xl mb-6 group-hover:scale-110 transition-transform', stat.color]">
          <el-icon><component :is="stat.icon" /></el-icon>
        </div>
        <p class="text-xs text-neutral-400 font-bold uppercase tracking-widest mb-1">{{ stat.label }}</p>
        <div class="flex items-baseline gap-2">
          <h2 class="text-4xl font-black tracking-tight text-neutral-800">{{ stat.value }}</h2>
          <span v-if="stat.unit" class="text-sm text-neutral-400 font-bold">{{ stat.unit }}</span>
        </div>
        <div class="mt-4 flex items-center gap-2">
           <span :class="['text-xs font-bold', stat.trend > 0 ? 'text-emerald-500' : 'text-rose-500']">
             {{ stat.trend > 0 ? '↑' : '↓' }} {{ Math.abs(stat.trend) }}%
           </span>
           <span class="text-[10px] text-neutral-300 font-medium uppercase">对比上月</span>
        </div>
      </div>
    </div>

    <!-- 详细消耗记录 -->
    <GvCard shadow="none" bordered padding="lg" class="!rounded-[2.5rem]">
      <template #header>
        <div class="flex items-center justify-between">
          <h3 class="text-2xl font-black text-neutral-900 tracking-tight">资源消耗明细</h3>
          <div class="flex items-center gap-4">
             <GvTag type="info" size="small" class="rounded-full bg-neutral-100 border-none text-neutral-500 font-bold">最近 30 天</GvTag>
          </div>
        </div>
      </template>

      <div class="mt-4 text-center py-20 border-2 border-dashed border-neutral-100 rounded-3xl">
          <el-icon :size="48" class="text-neutral-200 mb-4"><DataAnalysis /></el-icon>
          <p class="text-neutral-400 font-bold">详细调用日志正在同步中...</p>
      </div>
    </GvCard>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  Operation as OperatorIcon,
  Cpu as CpuIcon,
  DataAnalysis as TokenIcon,
  DataAnalysis
} from '@element-plus/icons-vue'
import GvCard from '@/components/base/GvCard/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'
import { userAssetApi, type UsageStats } from '@/api/user-asset'

const stats = ref<UsageStats>({
  operator_calls: 0,
  ai_model_calls: 0,
  token_usage: 0
})

const mainStats = computed(() => [
  { label: '算子总调用', value: stats.value.operator_calls.toLocaleString(), icon: OperatorIcon, color: 'bg-blue-50 text-blue-600', trend: 0, unit: '次' },
  { label: 'AI 模型调用', value: stats.value.ai_model_calls.toLocaleString(), icon: CpuIcon, color: 'bg-purple-50 text-purple-600', trend: 0, unit: '次' },
  { label: 'Token 消耗', value: (stats.value.token_usage / 1000000).toFixed(2), icon: TokenIcon, color: 'bg-amber-50 text-amber-600', trend: 0, unit: 'M' }
])

onMounted(() => {
  fetchUsage()
})

async function fetchUsage() {
  try {
    const res = await userAssetApi.getUsage()
    if (res.data && res.data.data) {
      stats.value = res.data.data
    }
  } catch (error) {}
}
</script>
