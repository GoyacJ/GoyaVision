<template>
  <GvContainer max-width="full">
    <PageHeader
      title="媒体源"
      description="管理拉流/推流媒体源，与 MediaMTX 一一对应，用于流媒体资产接入"
    >
      <template #actions>
        <GvSpace>
          <SearchBar
            v-model="searchName"
            placeholder="搜索媒体源名称"
            class="w-80"
            immediate
            :show-button="false"
            @search="refreshTable"
          />
          <GvButton @click="showCreateDialog = true">
            <template #icon>
              <el-icon><Plus /></el-icon>
            </template>
            添加媒体源
          </GvButton>
        </GvSpace>
      </template>
    </PageHeader>

    <ErrorState
      v-if="error && !loading"
      :error="error"
      title="加载失败"
      @retry="refreshTable"
    />

    <EmptyState
      v-else-if="!loading && sources.length === 0"
      icon="📡"
      title="还没有媒体源"
      description="创建拉流/推流媒体源，用于接入流媒体资产"
      action-text="添加媒体源"
      show-action
      @action="showCreateDialog = true"
    />

    <GvTable
      v-else
      :data="sources"
      :columns="columns"
      :loading="loading"
      border
      stripe
      pagination
      :pagination-config="paginationConfig"
      @current-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <template #type="{ row }">
        <GvTag :color="row.type === 'pull' ? 'primary' : 'success'" size="small" variant="tonal">
          {{ row.type === 'pull' ? '拉流' : '推流' }}
        </GvTag>
      </template>
      <template #enabled="{ row }">
        <GvTag :color="row.enabled ? 'success' : 'neutral'" size="small" variant="tonal">
          {{ row.enabled ? '启用' : '禁用' }}
        </GvTag>
      </template>
      <template #created_at="{ row }">
        {{ formatDate(row.created_at) }}
      </template>
      <template #actions="{ row }">
        <GvSpace size="xs">
          <GvButton size="small" variant="tonal" @click="handlePreview(row)">
            预览
          </GvButton>
          <GvButton size="small" @click="handleView(row)">
            查看
          </GvButton>
          <GvButton size="small" @click="handleEdit(row)">
            编辑
          </GvButton>
          <GvButton size="small" variant="text" @click="handleDelete(row)">
            删除
          </GvButton>
        </GvSpace>
      </template>
    </GvTable>

    <GvModal
      v-model="showCreateDialog"
      title="添加媒体源"
      :confirm-loading="creating"
      @confirm="handleCreate"
      @cancel="showCreateDialog = false"
    >
      <el-form ref="createFormRef" :model="createForm" :rules="createRules" label-width="100px">
        <el-form-item label="名称" prop="name">
          <GvInput v-model="createForm.name" placeholder="显示名称，如：摄像头1" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="createForm.type">
            <el-radio value="pull">拉流</el-radio>
            <el-radio value="push">推流</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="createForm.type === 'pull'" label="流地址" prop="url">
          <GvInput
            v-model="createForm.url"
            placeholder="rtsp://...、rtmp://...、https://.../live.m3u8"
            type="textarea"
            :rows="2"
          />
        </el-form-item>
        <el-form-item v-if="createForm.type === 'push'" label="说明">
          <span class="text-text-tertiary text-sm">推流源创建后，在详情/预览中可获取推流地址，用于 OBS 等配置</span>
        </el-form-item>
        <el-form-item label="可见范围">
          <GvSelect
            v-model="createForm.visibility"
            :options="VISIBILITY_OPTIONS"
            placeholder="请选择可见范围"
          />
        </el-form-item>
        <el-form-item label="启用" prop="enabled">
          <el-switch v-model="createForm.enabled" />
        </el-form-item>
      </el-form>
    </GvModal>

    <GvModal
      v-model="showEditDialog"
      title="编辑媒体源"
      :confirm-loading="updating"
      @confirm="handleUpdate"
      @cancel="showEditDialog = false"
    >
      <el-form ref="editFormRef" :model="editForm" :rules="editRules" label-width="100px">
        <el-form-item label="名称" prop="name">
          <GvInput v-model="editForm.name" />
        </el-form-item>
        <el-form-item v-if="currentSource?.type === 'pull'" label="流地址" prop="url">
          <GvInput v-model="editForm.url" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="可见范围">
          <GvSelect
            v-model="editForm.visibility"
            :options="VISIBILITY_OPTIONS"
            placeholder="请选择可见范围"
          />
        </el-form-item>
        <el-form-item label="启用" prop="enabled">
          <el-switch v-model="editForm.enabled" />
        </el-form-item>
      </el-form>
    </GvModal>

    <GvModal
      v-model="showViewDialog"
      title="媒体源详情"
      size="large"
      :show-confirm="false"
      cancel-text="关闭"
    >
      <template v-if="currentSource">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID" :span="2">{{ currentSource.id }}</el-descriptions-item>
          <el-descriptions-item label="名称">{{ currentSource.name }}</el-descriptions-item>
          <el-descriptions-item label="Path">{{ currentSource.path_name }}</el-descriptions-item>
          <el-descriptions-item label="类型">
            <GvTag :color="currentSource.type === 'pull' ? 'primary' : 'success'" size="small">
              {{ currentSource.type === 'pull' ? '拉流' : '推流' }}
            </GvTag>
          </el-descriptions-item>
          <el-descriptions-item label="启用">
            <GvTag :color="currentSource.enabled ? 'success' : 'neutral'" size="small">
              {{ currentSource.enabled ? '是' : '否' }}
            </GvTag>
          </el-descriptions-item>
          <el-descriptions-item v-if="currentSource.type === 'pull'" label="流地址" :span="2">
            <span class="font-mono text-sm break-all">{{ currentSource.url || '-' }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间" :span="2">{{ formatDate(currentSource.created_at) }}</el-descriptions-item>
        </el-descriptions>
        <div v-if="previewUrls" class="mt-4">
          <h4 class="text-sm font-semibold text-text-primary mb-2">预览与推流地址</h4>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="HLS">{{ previewUrls.hls_url }}</el-descriptions-item>
            <el-descriptions-item label="RTSP">{{ previewUrls.rtsp_url }}</el-descriptions-item>
            <el-descriptions-item label="RTMP">{{ previewUrls.rtmp_url }}</el-descriptions-item>
            <el-descriptions-item v-if="previewUrls.push_url" label="推流地址（OBS）">
              <span class="text-primary-600 font-mono">{{ previewUrls.push_url }}</span>
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </template>
    </GvModal>

    <GvModal
      v-model="showPreviewDialog"
      title="流预览"
      size="large"
      :show-confirm="false"
      cancel-text="关闭"
    >
      <div v-if="previewUrls" class="space-y-2">
        <p class="text-text-secondary text-sm">HLS 预览（若流已就绪可播放）</p>
        <div class="aspect-video bg-black rounded overflow-hidden">
          <video
            v-if="previewHlsUrl"
            ref="previewVideoRef"
            class="w-full h-full"
            controls
            muted
            playsinline
            :src="previewHlsUrl"
          />
        </div>
        <p class="text-text-tertiary text-xs">若无法播放请检查流是否已推/拉成功，或复制下方地址到播放器</p>
        <el-input v-model="previewHlsUrl" readonly class="font-mono text-xs" />
      </div>
    </GvModal>
  </GvContainer>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { sourceApi, type MediaSource, type SourceCreateReq, type SourcePreviewResponse } from '@/api/source'
import { useTable } from '@/composables'
import GvContainer from '@/components/layout/GvContainer/index.vue'
import GvTable from '@/components/base/GvTable/index.vue'
import GvModal from '@/components/base/GvModal/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvSpace from '@/components/layout/GvSpace/index.vue'
import GvTag from '@/components/base/GvTag/index.vue'
import GvInput from '@/components/base/GvInput/index.vue'
import GvSelect from '@/components/base/GvSelect/index.vue'
import PageHeader from '@/components/business/PageHeader/index.vue'
import SearchBar from '@/components/business/SearchBar/index.vue'
import { ErrorState, EmptyState } from '@/components/common'
import { VISIBILITY_OPTIONS } from '@/constants/visibility'

// UI 状态
const searchName = ref('')
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const showViewDialog = ref(false)
const showPreviewDialog = ref(false)
const currentSource = ref<MediaSource | null>(null)
const previewUrls = ref<SourcePreviewResponse | null>(null)
const previewHlsUrl = ref('')
const previewVideoRef = ref<HTMLVideoElement | null>(null)
const creating = ref(false)
const updating = ref(false)
const createFormRef = ref<FormInstance>()
const editFormRef = ref<FormInstance>()

// 使用 useTable 管理媒体源列表
const {
  items: sources,
  isLoading: loading,
  error,
  pagination,
  goToPage,
  changePageSize,
  refreshTable
} = useTable(
  async (params) => {
    // 将 page/page_size 转换为 limit/offset
    const res = await sourceApi.list({
      limit: params.page_size,
      offset: (params.page - 1) * params.page_size
    })
    return {
      items: res.data?.items ?? [],
      total: res.data?.total ?? 0
    }
  },
  {
    immediate: true,
    initialPageSize: 20
  }
)

const paginationConfig = computed(() => ({
  currentPage: pagination.page,
  pageSize: pagination.pageSize,
  total: pagination.total
}))

// 使用 useTable 管理媒体源列表
const createForm = reactive<any>({
  name: '',
  type: 'pull',
  url: '',
  enabled: true,
  visibility: 0
})

const editForm = reactive<any>({
  name: '',
  url: '',
  enabled: true,
  visibility: 0
})

const columns = [
  { prop: 'name', label: '名称', minWidth: 140 },
  { prop: 'path_name', label: 'Path', minWidth: 180 },
  { prop: 'type', label: '类型', width: 90 },
  { prop: 'url', label: '流地址', minWidth: 200 },
  { prop: 'enabled', label: '状态', width: 90 },
  { prop: 'created_at', label: '创建时间', width: 170 },
  { prop: 'actions', label: '操作', width: 220, fixed: 'right' }
]

const createRules: FormRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  url: [
    {
      validator: (_rule, value, callback) => {
        if (createForm.type === 'pull' && (!value || !String(value).trim())) {
          callback(new Error('拉流时请输入流地址'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const editRules: FormRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }]
}

function formatDate(s: string) {
  if (!s) return '-'
  const d = new Date(s)
  return d.toLocaleString('zh-CN')
}

// 直接使用 useTable 提供的方法
const handlePageChange = goToPage
const handleSizeChange = changePageSize

async function handleCreate() {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (!valid) return
    creating.value = true
    try {
      await sourceApi.create({
        name: createForm.name,
        type: createForm.type,
        url: createForm.type === 'pull' ? createForm.url : undefined,
        enabled: createForm.enabled,
        visibility: createForm.visibility ?? 0
      })
      ElMessage.success('创建成功')
      showCreateDialog.value = false
      createForm.name = ''
      createForm.type = 'pull'
      createForm.url = ''
      createForm.enabled = true
      createForm.visibility = 0
      refreshTable()
    } catch (e: any) {
      ElMessage.error(e.response?.data?.message || '创建失败')
    } finally {
      creating.value = false
    }
  })
}

function openView(row: MediaSource) {
  currentSource.value = row
  previewUrls.value = null
  showViewDialog.value = true
  sourceApi.getPreview(row.id).then((res) => {
    previewUrls.value = res.data ?? null
  }).catch(() => {
    previewUrls.value = null
  })
}

function handleView(row: MediaSource) {
  openView(row)
}

async function handlePreview(row: MediaSource) {
  try {
    const res = await sourceApi.getPreview(row.id)
    previewUrls.value = res.data ?? null
    previewHlsUrl.value = res.data?.hls_url ?? ''
    currentSource.value = row
    showPreviewDialog.value = true
  } catch (e: any) {
    ElMessage.error(e.response?.data?.message || '获取预览地址失败')
  }
}

watch(showPreviewDialog, (visible) => {
  if (visible && previewHlsUrl.value && previewVideoRef.value) {
    setTimeout(() => {
      if (previewVideoRef.value) {
        const v = previewVideoRef.value
        if (typeof v.load === 'function') v.load()
      }
    }, 100)
  }
})

function handleEdit(row: MediaSource) {
  currentSource.value = row
  editForm.name = row.name
  editForm.url = row.url ?? ''
  editForm.enabled = row.enabled
  editForm.visibility = row.visibility ?? 0
  showEditDialog.value = true
}

async function handleUpdate() {
  if (!editFormRef.value || !currentSource.value) return
  await editFormRef.value.validate(async (valid) => {
    if (!valid) return
    updating.value = true
    try {
      await sourceApi.update(currentSource.value!.id, {
        name: editForm.name,
        url: currentSource.value!.type === 'pull' ? editForm.url : undefined,
        enabled: editForm.enabled,
        visibility: editForm.visibility ?? 0
      })
      ElMessage.success('更新成功')
      showEditDialog.value = false
      refreshTable()
    } catch (e: any) {
      ElMessage.error(e.response?.data?.message || '更新失败')
    } finally {
      updating.value = false
    }
  })
}

async function handleDelete(row: MediaSource) {
  try {
    await ElMessageBox.confirm('删除媒体源将同时删除 MediaMTX 对应 path，且仅当无关联流媒体资产时允许。确定删除？', '确认删除', {
      type: 'warning'
    })
    await sourceApi.remove(row.id)
    ElMessage.success('删除成功')
    refreshTable()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.response?.data?.message || '删除失败')
    }
  }
}

</script>
