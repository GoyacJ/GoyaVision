<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
    <!-- 积分概览 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <GvCard shadow="none" bordered padding="lg" class="!rounded-[2.5rem] md:col-span-2 bg-gradient-to-br from-primary-600 to-indigo-700 text-white relative overflow-hidden">
        <div class="absolute top-0 right-0 w-64 h-64 bg-white/10 rounded-full -mr-20 -mt-20 blur-2xl"></div>
        <div class="relative z-10">
          <p class="text-primary-100 font-bold uppercase tracking-widest text-xs mb-3">可用积分</p>
          <div class="flex items-end gap-3 mb-8">
            <h2 class="text-6xl font-black tracking-tighter">{{ points.toLocaleString() }}</h2>
            <span class="text-primary-200 font-bold mb-2">PTS</span>
          </div>
          <div class="flex gap-3">
             <GvTag effect="dark" class="!bg-white/20 !border-none !text-white rounded-full px-4 font-bold">
               等级: {{ memberLevel }}
             </GvTag>
             <GvTag v-if="memberLevel === 'Gold'" effect="dark" class="!bg-white/20 !border-none !text-white rounded-full px-4 font-bold">
               积分倍率: 1.2x
             </GvTag>
          </div>
        </div>
      </GvCard>

      <GvCard shadow="none" bordered padding="lg" class="!rounded-[2.5rem] flex flex-col justify-center items-center text-center group hover:border-primary-200 transition-colors">
        <div class="w-20 h-20 rounded-3xl bg-primary-50 text-primary-600 flex items-center justify-center text-4xl mb-6 group-hover:scale-110 transition-transform">
          <el-icon><StarIcon /></el-icon>
        </div>
        <h4 class="text-xl font-black text-neutral-800 mb-2">每日签到</h4>
        <p class="text-sm text-neutral-400 font-medium mb-6">连续签到可获得额外积分奖励</p>
        <GvButton 
          type="primary" 
          class="!rounded-2xl w-full !h-12 font-black" 
          :loading="submitting" 
          @click="handleCheckIn"
        >
          立即领取
        </GvButton>
      </GvCard>
    </div>

    <!-- 积分明细 -->
    <GvCard shadow="none" bordered padding="lg" class="!rounded-[2.5rem]">
      <template #header>
        <div class="flex items-center justify-between">
          <h3 class="text-2xl font-black text-neutral-900 tracking-tight">积分记录</h3>
          <div class="flex gap-2">
            <GvButton size="small" class="!rounded-full px-4 font-bold" @click="fetchHistory">刷新</GvButton>
          </div>
        </div>
      </template>

      <div class="mt-4">
        <GvTable 
          :columns="columns" 
          :data="history" 
          :loading="loading"
          row-key="id"
        >
          <template #change="{ row }">
            <span :class="['font-black text-lg', row.change > 0 ? 'text-primary-600' : 'text-rose-600']">
              {{ row.change > 0 ? '+' : '' }}{{ row.change }}
            </span>
          </template>
          <template #type="{ row }">
             <div class="flex items-center gap-3">
                <div :class="['w-8 h-8 rounded-lg flex items-center justify-center', row.change > 0 ? 'bg-primary-50 text-primary-600' : 'bg-rose-50 text-rose-600']">
                  <el-icon><component :is="row.change > 0 ? 'CircleCheck' : 'ShoppingCart'" /></el-icon>
                </div>
                <span class="font-bold">{{ row.type }}</span>
             </div>
          </template>
          <template #created_at="{ row }">
            {{ new Date(row.created_at).toLocaleString() }}
          </template>
        </GvTable>
      </div>
    </GvCard>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Star as StarIcon,
  CircleCheck,
  ShoppingCart
} from '@element-plus/icons-vue'
import GvCard from '@/components/base/GvCard/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvTable from '@/components/base/GvTable/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'
import { userAssetApi, type PointRecord } from '@/api/user-asset'

const points = ref(0)
const memberLevel = ref('Free')
const loading = ref(false)
const submitting = ref(false)
const history = ref<PointRecord[]>([])

const columns = [
  { prop: 'type', label: '变动类型', minWidth: '200' },
  { prop: 'change', label: '变动数值', width: '150', align: 'right' },
  { prop: 'balance', label: '变动后余额', width: '150', align: 'right' },
  { prop: 'created_at', label: '操作时间', width: '200' }
]

onMounted(() => {
  fetchSummary()
  fetchHistory()
})

async function fetchSummary() {
  try {
    const res = await userAssetApi.getSummary()
    if (res.data && res.data.data) {
      points.value = res.data.data.points
      memberLevel.value = res.data.data.member_level
    }
  } catch (error) {}
}

async function fetchHistory() {
  loading.value = true
  try {
    const res = await userAssetApi.getPoints({ limit: 10 })
    if (res.data && res.data.data) {
      history.value = res.data.data
    }
  } catch (error) {
    ElMessage.error('获取记录失败')
  } finally {
    loading.value = false
  }
}

async function handleCheckIn() {
  submitting.value = true
  try {
    await userAssetApi.checkIn()
    ElMessage.success('签到成功，获得 50 积分')
    fetchSummary()
    fetchHistory()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '签到失败，请稍后再试')
  } finally {
    submitting.value = false
  }
}
</script>
