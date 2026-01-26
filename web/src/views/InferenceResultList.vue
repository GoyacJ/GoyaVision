<template>
  <div class="inference-result-list">
    <div class="header">
      <h2>推理结果</h2>
    </div>

    <el-form :inline="true" class="filter-form">
      <el-form-item label="流ID">
        <el-input v-model="query.stream_id" placeholder="可选" clearable />
      </el-form-item>
      <el-form-item label="绑定ID">
        <el-input v-model="query.binding_id" placeholder="可选" clearable />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadResults">查询</el-button>
        <el-button @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="results" v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="300" />
      <el-table-column prop="stream_id" label="流ID" width="300" />
      <el-table-column prop="algorithm_binding_id" label="绑定ID" width="300" />
      <el-table-column prop="ts" label="时间戳" width="200" />
      <el-table-column prop="latency_ms" label="延迟(ms)" width="120" />
      <el-table-column label="输出" min-width="200">
        <template #default="{ row }">
          <el-button size="small" @click="showOutput(row)">查看</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="loadResults"
      @current-change="loadResults"
    />

    <el-dialog v-model="showOutputDialog" title="推理输出" width="600px">
      <pre>{{ outputContent }}</pre>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { inferenceApi, type InferenceResult, type InferenceResultListQuery } from '../api/inference'

const results = ref<InferenceResult[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const showOutputDialog = ref(false)
const outputContent = ref('')

const query = ref<InferenceResultListQuery>({
  stream_id: '',
  binding_id: ''
})

const loadResults = async () => {
  loading.value = true
  try {
    const params: InferenceResultListQuery = {
      ...query.value,
      limit: pageSize.value,
      offset: (currentPage.value - 1) * pageSize.value
    }
    if (!params.stream_id) delete params.stream_id
    if (!params.binding_id) delete params.binding_id

    const res = await inferenceApi.list(params)
    results.value = res.data.items
    total.value = res.data.total
  } catch (error: any) {
    ElMessage.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

const resetQuery = () => {
  query.value = { stream_id: '', binding_id: '' }
  currentPage.value = 1
  loadResults()
}

const showOutput = (row: InferenceResult) => {
  outputContent.value = JSON.stringify(row.output, null, 2)
  showOutputDialog.value = true
}

onMounted(() => {
  loadResults()
})
</script>

<style scoped>
.inference-result-list {
  padding: 20px;
}

.header {
  margin-bottom: 20px;
}

.header h2 {
  margin: 0;
}

.filter-form {
  margin-bottom: 20px;
}

pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  max-height: 400px;
  overflow-y: auto;
}
</style>
