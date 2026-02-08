<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
    <GvCard shadow="none" bordered padding="lg" class="!rounded-[2.5rem]">
      <template #header>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <div class="w-1.5 h-8 bg-primary-600 rounded-full"></div>
            <h3 class="text-2xl font-black text-neutral-900 tracking-tight">基本资料</h3>
          </div>
          <GvTag type="primary" size="small" effect="light" class="rounded-full px-4">
            公开身份
          </GvTag>
        </div>
      </template>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-x-12 gap-y-8 mt-4">
        <GvInput
          v-model="form.nickname"
          label="显示名称"
          placeholder="取一个独特的昵称"
          size="large"
          :disabled="loading"
        >
          <template #prefix>
            <el-icon><UserIcon /></el-icon>
          </template>
        </GvInput>

        <GvInput
          :model-value="username"
          label="用户账号"
          size="large"
          disabled
        >
          <template #prefix>
            <el-icon><MonitorIcon /></el-icon>
          </template>
        </GvInput>

        <GvInput
          v-model="form.email"
          label="电子邮箱"
          placeholder="yourname@domain.com"
          size="large"
          :disabled="loading"
        >
          <template #prefix>
            <el-icon><MessageIcon /></el-icon>
          </template>
        </GvInput>

        <GvInput
          v-model="form.phone"
          label="联系电话"
          placeholder="您的手机号码"
          size="large"
          :disabled="loading"
        >
          <template #prefix>
            <el-icon><PhoneIcon /></el-icon>
          </template>
        </GvInput>
      </div>

      <template #footer>
        <div class="flex flex-col sm:flex-row items-center justify-between gap-6">
          <div class="flex items-center gap-2 text-neutral-400 text-sm font-medium">
            <el-icon class="animate-pulse"><ClockIcon /></el-icon>
            <span>资料最后同步于: {{ lastSyncTime }}</span>
          </div>
          <div class="flex gap-4 items-center">
            <span v-if="hasChanges" class="text-sm text-amber-500 font-bold flex items-center gap-2">
              <span class="w-2 h-2 rounded-full bg-amber-500 animate-ping"></span>
              检测到未保存的更改
            </span>
            <GvButton
              type="primary"
              size="large"
              class="!rounded-2xl px-12 !h-14 font-black shadow-xl shadow-primary-200 transition-all hover:-translate-y-1 active:translate-y-0"
              :loading="loading"
              :disabled="!hasChanges"
              @click="$emit('save', form)"
            >
              更新资料
            </GvButton>
          </div>
        </div>
      </template>
    </GvCard>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed, watch } from 'vue'
import { 
  User as UserIcon, 
  Monitor as MonitorIcon, 
  Message as MessageIcon, 
  Phone as PhoneIcon,
  Clock as ClockIcon
} from '@element-plus/icons-vue'
import GvCard from '@/components/base/GvCard/index.vue'
import GvInput from '@/components/base/GvInput/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'

const props = defineProps<{
  initialData: {
    nickname: string
    email: string
    phone: string
  }
  username: string
  lastSyncTime: string
  loading: boolean
}>()

const emit = defineEmits<{
  (e: 'save', data: any): void
}>()

const form = reactive({ ...props.initialData })

const hasChanges = computed(() => {
  return (
    form.nickname !== props.initialData.nickname ||
    form.email !== props.initialData.email ||
    form.phone !== props.initialData.phone
  )
})

watch(() => props.initialData, (newData) => {
  Object.assign(form, newData)
}, { deep: true })
</script>
