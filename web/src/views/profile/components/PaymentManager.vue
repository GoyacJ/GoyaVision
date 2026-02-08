<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
    <!-- 余额概览 -->
    <GvCard shadow="none" bordered padding="lg" class="!rounded-[2.5rem] bg-gradient-to-br from-neutral-900 to-neutral-800 text-white overflow-hidden relative">
      <div class="absolute top-0 right-0 w-64 h-64 bg-primary-500/10 rounded-full -mr-32 -mt-32 blur-3xl"></div>
      <div class="absolute bottom-0 left-0 w-48 h-48 bg-purple-500/10 rounded-full -ml-24 -mb-24 blur-3xl"></div>
      
      <div class="relative z-10">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-8">
          <div>
            <p class="text-neutral-400 font-bold uppercase tracking-widest text-xs mb-2">账户余额</p>
            <h2 class="text-5xl font-black tracking-tighter flex items-baseline gap-2">
              <span class="text-2xl font-medium opacity-50">¥</span>
              {{ balance.toFixed(2) }}
            </h2>
          </div>
          <div class="flex gap-4">
            <GvButton type="primary" size="large" class="!rounded-2xl !h-14 px-8 font-black !bg-white !text-neutral-900 !border-none" @click="showRechargeDialog = true">
              立即充值
            </GvButton>
            <GvButton size="large" class="!rounded-2xl !h-14 px-8 font-black !bg-neutral-700/50 !text-white !border-neutral-600 hover:!bg-neutral-700">
              提现申请
            </GvButton>
          </div>
        </div>
      </div>
    </GvCard>

    <!-- 支付记录 -->
    <GvCard shadow="none" bordered padding="lg" class="!rounded-[2.5rem]">
      <template #header>
        <div class="flex items-center justify-between">
          <h3 class="text-2xl font-black text-neutral-900 tracking-tight">支付记录</h3>
          <GvButton link class="font-bold !text-primary-600" @click="fetchHistory">刷新记录</GvButton>
        </div>
      </template>

      <div class="mt-4">
        <GvTable 
          :columns="columns" 
          :data="history" 
          :loading="loading"
          row-key="id"
        >
          <template #amount="{ row }">
            <span :class="['font-black', row.amount > 0 ? 'text-emerald-600' : 'text-rose-600']">
              {{ row.amount > 0 ? '+' : '' }}{{ row.amount.toFixed(2) }}
            </span>
          </template>
          <template #status="{ row }">
            <GvTag :type="getStatusType(row.status)" size="small" class="rounded-full px-3 font-bold">
              {{ row.status }}
            </GvTag>
          </template>
          <template #method="{ row }">
            <div class="flex items-center gap-2">
              <el-icon :class="getMethodColor(row.method)"><component :is="getMethodIcon(row.method)" /></el-icon>
              <span class="font-medium">{{ row.method }}</span>
            </div>
          </template>
          <template #created_at="{ row }">
            {{ new Date(row.created_at).toLocaleString() }}
          </template>
        </GvTable>
      </div>
    </GvCard>

    <!-- 充值弹窗 -->
    <el-dialog v-model="showRechargeDialog" title="账户充值" width="500px" custom-class="rounded-3xl">
      <div class="p-4">
        <div class="mb-8">
          <p class="text-sm font-bold text-neutral-400 uppercase tracking-widest mb-4">选择充值金额</p>
          <div class="grid grid-cols-3 gap-4">
            <div 
              v-for="amt in [50, 100, 200, 500, 1000, 2000]" 
              :key="amt"
              class="py-4 border-2 rounded-2xl text-center cursor-pointer transition-all"
              :class="rechargeAmount === amt ? 'border-primary-600 bg-primary-50 text-primary-600 font-black' : 'border-neutral-100 hover:border-neutral-200 font-bold text-neutral-600'"
              @click="rechargeAmount = amt"
            >
              ¥{{ amt }}
            </div>
          </div>
        </div>

        <div class="mb-8">
          <p class="text-sm font-bold text-neutral-400 uppercase tracking-widest mb-4">支付方式</p>
          <div class="space-y-3">
            <div 
              v-for="method in paymentMethods" 
              :key="method.id"
              class="flex items-center justify-between p-4 border-2 rounded-2xl cursor-pointer transition-all"
              :class="selectedMethod === method.id ? 'border-primary-600 bg-primary-50' : 'border-neutral-100 hover:border-neutral-200'"
              @click="selectedMethod = method.id"
            >
              <div class="flex items-center gap-4">
                <div :class="['w-10 h-10 rounded-xl flex items-center justify-center text-xl', method.color]">
                  <el-icon><component :is="method.icon" /></el-icon>
                </div>
                <span class="font-bold text-neutral-800">{{ method.name }}</span>
              </div>
              <div v-if="selectedMethod === method.id" class="w-6 h-6 bg-primary-600 rounded-full flex items-center justify-center">
                <el-icon color="white" :size="14"><CheckIcon /></el-icon>
              </div>
            </div>
          </div>
        </div>

        <GvButton 
          type="primary" 
          size="large" 
          class="w-full !rounded-2xl !h-14 font-black shadow-xl shadow-primary-100"
          :loading="submitting"
          @click="handleRecharge"
        >
          确认充值 ¥{{ rechargeAmount }}
        </GvButton>
      </div>
    </el-dialog>

    <!-- 扫码支付弹窗 -->
    <el-dialog v-model="showQRCodeDialog" title="扫码支付" width="400px" center custom-class="rounded-3xl">
      <div class="p-6 text-center">
        <div class="mb-6 bg-neutral-50 p-6 rounded-3xl inline-block border-2 border-neutral-100">
           <!-- 这里实际应用二维码库生成，暂时显示占位 -->
           <div class="w-48 h-48 bg-white flex items-center justify-center border border-dashed border-neutral-200">
              <span class="text-xs text-neutral-400">QR Code: {{ qrcodeUrl }}</span>
           </div>
        </div>
        <p class="text-neutral-500 font-bold mb-2">请使用微信扫码支付</p>
        <p class="text-2xl font-black text-neutral-900 mb-6">¥{{ rechargeAmount.toFixed(2) }}</p>
        <GvButton class="w-full !rounded-2xl" @click="showQRCodeDialog = false">支付完成后点击刷新</GvButton>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  CreditCard as UnionPayIcon, 
  ChatDotRound as WeChatIcon, 
  Wallet as AlipayIcon,
  Check as CheckIcon,
  Timer as ClockIcon
} from '@element-plus/icons-vue'
import GvCard from '@/components/base/GvCard/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvTable from '@/components/base/GvTable/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'
import { userAssetApi, type TransactionRecord } from '@/api/user-asset'

const balance = ref(0)
const loading = ref(false)
const submitting = ref(false)
const showRechargeDialog = ref(false)
const showQRCodeDialog = ref(false)
const rechargeAmount = ref(100)
const selectedMethod = ref('alipay')
const qrcodeUrl = ref('')
const history = ref<TransactionRecord[]>([])

const paymentMethods = [
  { id: 'alipay', name: '支付宝', icon: AlipayIcon, color: 'bg-blue-50 text-blue-600' },
  { id: 'wechat', name: '微信支付', icon: WeChatIcon, color: 'bg-emerald-50 text-emerald-600' },
  { id: 'unionpay', name: '银联支付', icon: UnionPayIcon, color: 'bg-red-50 text-red-600' }
]

const columns = [
  { prop: 'id', label: '订单号', width: '220' },
  { prop: 'type', label: '类型', width: '120' },
  { prop: 'method', label: '支付方式', width: '140' },
  { prop: 'amount', label: '金额 (CNY)', width: '120', align: 'right' },
  { prop: 'status', label: '状态', width: '120', align: 'center' },
  { prop: 'created_at', label: '时间', minWidth: '180' }
]

onMounted(() => {
  fetchSummary()
  fetchHistory()
})

async function fetchSummary() {
  try {
    const res = await userAssetApi.getSummary()
    // res.data 是 AxiosResponse.data (即 ApiResponse)
    // res.data.data 才是 UserAssetSummary
    if (res.data && res.data.data) {
      balance.value = res.data.data.balance
    }
  } catch (error) {}
}

async function fetchHistory() {
  loading.value = true
  try {
    const res = await userAssetApi.getTransactions({ limit: 10 })
    if (res.data && res.data.data) {
      history.value = res.data.data
    }
  } catch (error) {
    ElMessage.error('获取记录失败')
  } finally {
    loading.value = false
  }
}

async function handleRecharge() {
  submitting.value = true
  try {
    const res = await userAssetApi.recharge({
      amount: rechargeAmount.value,
      channel: selectedMethod.value
    })
    
    showRechargeDialog.value = false
    
    const data = res.data.data
    if (selectedMethod.value === 'alipay' && data.pay_url) {
      // 支付宝跳转
      window.location.href = data.pay_url
    } else if (selectedMethod.value === 'wechat' && data.qrcode) {
      // 微信扫码
      qrcodeUrl.value = data.qrcode
      showQRCodeDialog.value = true
    } else if (selectedMethod.value === 'unionpay') {
       ElMessage.warning('银联支付维护中')
    } else {
       ElMessage.success('订单已创建')
       fetchHistory()
    }
  } catch (error) {
    ElMessage.error('充值请求失败')
  } finally {
    submitting.value = false
  }
}

function getStatusType(status: string) {
  if (status === 'Success') return 'success'
  if (status === 'Pending') return 'warning'
  if (status === 'Failed') return 'danger'
  return 'info'
}

function getMethodIcon(method: string) {
  const m = method.toLowerCase()
  if (m.includes('alipay')) return AlipayIcon
  if (m.includes('wechat')) return WeChatIcon
  if (m.includes('union')) return UnionPayIcon
  return ClockIcon
}

function getMethodColor(method: string) {
  const m = method.toLowerCase()
  if (m.includes('alipay')) return 'text-blue-500'
  if (m.includes('wechat')) return 'text-emerald-500'
  if (m.includes('union')) return 'text-red-500'
  return 'text-neutral-400'
}
</script>
