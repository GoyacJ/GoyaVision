<template>
  <GvContainer max-width="full">
    <PageHeader
      title="MCP 市场"
      description="浏览并安装 MCP 资源为算子模板"
    >
      <template #actions>
        <GvSpace>
          <SearchBar
            v-model="keyword"
            placeholder="搜索模板"
            class="w-80"
            immediate
            :show-button="false"
            @search="loadTemplates"
          />
          <el-select v-model="selectedServerId" class="w-56" placeholder="选择 MCP Server" clearable>
            <el-option
              v-for="s in mcpServers"
              :key="s.id"
              :label="`${s.name}(${s.id})`"
              :value="s.id"
            />
          </el-select>
          <GvButton variant="tonal" :loading="mcpToolsLoading" @click="loadMCPTools">加载 MCP 工具</GvButton>
          <GvButton :loading="syncing" @click="handleSyncMCP">同步 MCP 模板</GvButton>
        </GvSpace>
      </template>
    </PageHeader>

    <el-row :gutter="16" v-loading="loading">
      <el-col v-for="tpl in templates" :key="tpl.id" :xs="24" :sm="12" :md="8" :lg="6" class="mb-4">
        <TemplateCard :template="tpl" @install="openInstallDialog" @preview="handlePreviewTemplate" />
      </el-col>
    </el-row>

    <EmptyState
      v-if="!loading && templates.length === 0"
      icon="🧩"
      title="暂无模板"
      description="请先同步 MCP 模板或稍后再试"
    />

    <el-card class="mt-4" shadow="never">
      <template #header>
        <div class="font-medium">MCP 工具安装入口</div>
      </template>

      <el-table :data="mcpTools" v-loading="mcpToolsLoading" size="small" border>
        <el-table-column prop="name" label="工具名" min-width="180" />
        <el-table-column prop="version" label="版本" width="120" />
        <el-table-column prop="description" label="描述" min-width="260" />
        <el-table-column label="操作" width="140">
          <template #default="{ row }">
            <GvButton size="small" variant="tonal" @click="openInstallMCPDialog(row)">安装为算子</GvButton>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="!mcpToolsLoading && mcpTools.length === 0" class="mt-2 text-xs text-text-tertiary">
        请先选择 MCP Server 并加载工具
      </div>
    </el-card>

    <GvModal
      v-model="showInstallDialog"
      title="安装模板"
      size="default"
      :show-confirm="false"
      @cancel="showInstallDialog = false"
    >
      <el-form label-width="110px">
        <el-form-item label="模板名称">
          <el-input :model-value="selectedTemplate?.name || '-'" disabled />
        </el-form-item>
        <el-form-item label="算子代码">
          <el-input v-model="installForm.operator_code" placeholder="唯一编码" />
        </el-form-item>
        <el-form-item label="算子名称">
          <el-input v-model="installForm.operator_name" placeholder="显示名称" />
        </el-form-item>
        <el-form-item>
          <div class="flex gap-2">
            <GvButton :loading="installing" @click="handleInstall">确认安装</GvButton>
            <GvButton variant="tonal" @click="showInstallDialog = false">取消</GvButton>
          </div>
        </el-form-item>
      </el-form>
    </GvModal>

    <GvModal
      v-model="showPreviewDialog"
      title="MCP Tool 预览"
      size="large"
      :show-confirm="false"
      cancel-text="关闭"
    >
      <el-descriptions v-if="previewTool" :column="1" border>
        <el-descriptions-item label="工具名称">{{ previewTool.name }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ previewTool.version || '-' }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{ previewTool.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="输入 Schema">
          <pre class="max-h-64 overflow-auto rounded bg-neutral-50 p-3 text-xs">{{ formatJson(previewTool.input_schema) }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="输出 Schema">
          <pre class="max-h-64 overflow-auto rounded bg-neutral-50 p-3 text-xs">{{ formatJson(previewTool.output_schema) }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </GvModal>

    <GvModal
      v-model="showInstallMCPDialog"
      title="安装 MCP 工具为算子"
      size="default"
      :show-confirm="false"
      @cancel="showInstallMCPDialog = false"
    >
      <el-form label-width="110px">
        <el-form-item label="MCP Server">
          <el-input :model-value="selectedServerId || '-'" disabled />
        </el-form-item>
        <el-form-item label="工具名">
          <el-input :model-value="selectedToolName || '-'" disabled />
        </el-form-item>
        <el-form-item label="算子代码">
          <el-input v-model="mcpInstallForm.operator_code" placeholder="唯一编码" />
        </el-form-item>
        <el-form-item label="算子名称">
          <el-input v-model="mcpInstallForm.operator_name" placeholder="显示名称" />
        </el-form-item>
        <el-form-item>
          <div class="flex gap-2">
            <GvButton :loading="installingMCP" @click="handleInstallMCP">确认安装</GvButton>
            <GvButton variant="tonal" @click="showInstallMCPDialog = false">取消</GvButton>
          </div>
        </el-form-item>
      </el-form>
    </GvModal>
  </GvContainer>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { operatorApi, type MCPServer, type MCPTool, type OperatorTemplate } from '@/api/operator'
import GvContainer from '@/components/layout/GvContainer/index.vue'
import GvButton from '@/components/base/GvButton/index.vue'
import GvSpace from '@/components/layout/GvSpace/index.vue'
import GvModal from '@/components/base/GvModal/index.vue'
import PageHeader from '@/components/business/PageHeader/index.vue'
import SearchBar from '@/components/business/SearchBar/index.vue'
import { EmptyState } from '@/components/common'
import TemplateCard from '@/views/operator/components/TemplateCard.vue'

const loading = ref(false)
const syncing = ref(false)
const installing = ref(false)
const installingMCP = ref(false)
const keyword = ref('')
const templates = ref<OperatorTemplate[]>([])
const mcpServers = ref<MCPServer[]>([])
const selectedServerId = ref('')
const mcpTools = ref<MCPTool[]>([])
const mcpToolsLoading = ref(false)
const showInstallDialog = ref(false)
const showInstallMCPDialog = ref(false)
const showPreviewDialog = ref(false)
const selectedTemplate = ref<OperatorTemplate | null>(null)
const selectedToolName = ref('')
const previewTool = ref<MCPTool | null>(null)
const installForm = ref({
  operator_code: '',
  operator_name: ''
})
const mcpInstallForm = ref({
  operator_code: '',
  operator_name: ''
})

async function loadTemplates() {
  loading.value = true
  try {
    const res = await operatorApi.listTemplates({ keyword: keyword.value, page: 1, page_size: 50 })
    templates.value = res.data?.items || []
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '加载模板失败')
  } finally {
    loading.value = false
  }
}

async function loadMCPServers() {
  try {
    const res = await operatorApi.listMCPServers()
    mcpServers.value = res.data || []
    if (!selectedServerId.value && mcpServers.value.length > 0) {
      selectedServerId.value = mcpServers.value[0].id
    }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '加载 MCP Server 失败')
  }
}

async function loadMCPTools() {
  if (!selectedServerId.value) {
    ElMessage.warning('请先选择 MCP Server')
    return
  }
  mcpToolsLoading.value = true
  try {
    const res = await operatorApi.listMCPTools(selectedServerId.value)
    mcpTools.value = res.data || []
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '加载 MCP 工具失败')
  } finally {
    mcpToolsLoading.value = false
  }
}

function openInstallDialog(tpl: OperatorTemplate) {
  selectedTemplate.value = tpl
  installForm.value = {
    operator_code: `${tpl.code}_copy`,
    operator_name: `${tpl.name}-副本`
  }
  showInstallDialog.value = true
}

async function handleInstall() {
  if (!selectedTemplate.value) return
  if (!installForm.value.operator_code || !installForm.value.operator_name) {
    ElMessage.warning('请填写算子代码与名称')
    return
  }
  installing.value = true
  try {
    await operatorApi.installTemplate({
      template_id: selectedTemplate.value.id,
      operator_code: installForm.value.operator_code,
      operator_name: installForm.value.operator_name
    })
    ElMessage.success('安装模板成功')
    showInstallDialog.value = false
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '安装模板失败')
  } finally {
    installing.value = false
  }
}

async function handlePreviewTemplate(tpl: OperatorTemplate) {
  if (tpl.exec_mode !== 'mcp') {
    ElMessage.info('该模板非 MCP 来源，暂无 Tool 预览信息')
    return
  }

  try {
    const mcpConfig = tpl.exec_config?.mcp as { server_id?: string; tool_name?: string } | undefined
    const servers = await operatorApi.listMCPServers()
    const serverId = mcpConfig?.server_id || servers.data?.[0]?.id
    const toolName = mcpConfig?.tool_name
    if (!serverId) {
      ElMessage.warning('当前无可用 MCP Server')
      return
    }
    if (!toolName) {
      ElMessage.warning('模板缺少 MCP 工具信息，无法预览')
      return
    }
    const res = await operatorApi.previewMCPTool(serverId, toolName)
    previewTool.value = res.data
    showPreviewDialog.value = true
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '预览 MCP Tool 失败')
  }
}

function formatJson(value: any) {
  if (!value) return '{}'
  return JSON.stringify(value, null, 2)
}

async function handleSyncMCP() {
  syncing.value = true
  try {
    const serverId = selectedServerId.value || mcpServers.value[0]?.id
    if (!serverId) {
      ElMessage.warning('当前无可用 MCP Server，请先在后端注册')
      return
    }
    await operatorApi.syncMCPTemplates({ server_id: serverId })
    ElMessage.success('同步 MCP 模板成功')
    await loadTemplates()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '同步 MCP 模板失败')
  } finally {
    syncing.value = false
  }
}

function openInstallMCPDialog(tool: MCPTool) {
  if (!selectedServerId.value) {
    ElMessage.warning('请先选择 MCP Server')
    return
  }
  selectedToolName.value = tool.name
  mcpInstallForm.value = {
    operator_code: `${tool.name.replace(/[^a-zA-Z0-9_-]/g, '_')}_mcp`,
    operator_name: `${tool.name}-MCP算子`
  }
  showInstallMCPDialog.value = true
}

async function handleInstallMCP() {
  if (!selectedServerId.value || !selectedToolName.value) return
  if (!mcpInstallForm.value.operator_code || !mcpInstallForm.value.operator_name) {
    ElMessage.warning('请填写算子代码与名称')
    return
  }
  installingMCP.value = true
  try {
    await operatorApi.installMCPOperator({
      server_id: selectedServerId.value,
      tool_name: selectedToolName.value,
      operator_code: mcpInstallForm.value.operator_code,
      operator_name: mcpInstallForm.value.operator_name
    })
    ElMessage.success('MCP 工具安装成功')
    showInstallMCPDialog.value = false
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '安装 MCP 工具失败')
  } finally {
    installingMCP.value = false
  }
}

loadTemplates()
loadMCPServers()
</script>
