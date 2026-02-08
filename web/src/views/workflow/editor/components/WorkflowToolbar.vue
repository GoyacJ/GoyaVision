<template>
  <div class="h-14 border-b bg-white flex items-center px-4 justify-between">
    <div class="flex items-center gap-4">
      <el-button @click="$router.back()">
        <el-icon><Back /></el-icon>
      </el-button>
      <div class="font-medium text-lg">工作流编辑器</div>
    </div>

    <div class="flex items-center gap-2">
      <el-button @click="onAutoLayout">
        自动布局
      </el-button>
      <el-button type="primary" @click="onSave">
        保存
      </el-button>
      <el-button type="success" @click="onRun">
        运行
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Back } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import { useWorkflowEditor } from '../composables/useWorkflowEditor'
import { useAutoLayout } from '../composables/useAutoLayout'

const router = useRouter()
const { saveWorkflow, runWorkflow } = useWorkflowEditor()
const { layout } = useAutoLayout()

async function onSave() {
  await saveWorkflow()
}

async function onRun() {
  await runWorkflow()
}

function onAutoLayout() {
  layout()
}
</script>
