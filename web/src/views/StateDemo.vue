<template>
  <div class="state-demo-page">
    <div class="container-gv">
      <h1 class="text-h1 mb-8">状态组件演示</h1>

      <!-- LoadingState -->
      <section class="demo-section">
        <h2 class="text-h2 mb-4">LoadingState - 加载状态</h2>

        <div class="demo-grid">
          <!-- 基础用法 -->
          <div class="demo-card">
            <h3 class="demo-title">基础用法</h3>
            <LoadingState />
          </div>

          <!-- 带提示文本 -->
          <div class="demo-card">
            <h3 class="demo-title">带提示文本</h3>
            <LoadingState message="加载中..." />
          </div>

          <!-- 小尺寸 -->
          <div class="demo-card">
            <h3 class="demo-title">小尺寸</h3>
            <LoadingState size="small" message="加载中..." />
          </div>

          <!-- 大尺寸 -->
          <div class="demo-card">
            <h3 class="demo-title">大尺寸</h3>
            <LoadingState size="large" message="处理中，请稍候..." />
          </div>
        </div>
      </section>

      <!-- ErrorState -->
      <section class="demo-section">
        <h2 class="text-h2 mb-4">ErrorState - 错误状态</h2>

        <div class="demo-grid">
          <!-- 基础用法 -->
          <div class="demo-card">
            <h3 class="demo-title">基础用法</h3>
            <ErrorState @retry="handleRetry" />
          </div>

          <!-- 带错误对象 -->
          <div class="demo-card">
            <h3 class="demo-title">带错误对象</h3>
            <ErrorState :error="sampleError" @retry="handleRetry" />
          </div>

          <!-- 自定义内容 -->
          <div class="demo-card">
            <h3 class="demo-title">自定义内容</h3>
            <ErrorState
              title="网络连接失败"
              message="请检查网络连接后重试"
              retry-text="重新加载"
              @retry="handleRetry"
            />
          </div>

          <!-- 无重试按钮 -->
          <div class="demo-card">
            <h3 class="demo-title">无重试按钮</h3>
            <ErrorState
              title="权限不足"
              message="您没有访问此资源的权限"
              :show-retry="false"
            />
          </div>
        </div>
      </section>

      <!-- EmptyState -->
      <section class="demo-section">
        <h2 class="text-h2 mb-4">EmptyState - 空状态</h2>

        <div class="demo-grid">
          <!-- 基础用法 -->
          <div class="demo-card">
            <h3 class="demo-title">基础用法</h3>
            <EmptyState />
          </div>

          <!-- 自定义内容 -->
          <div class="demo-card">
            <h3 class="demo-title">自定义内容</h3>
            <EmptyState
              icon="🎬"
              title="还没有媒体资产"
              description="开始上传您的第一个视频、图片或音频文件"
            />
          </div>

          <!-- 带操作按钮 -->
          <div class="demo-card">
            <h3 class="demo-title">带操作按钮</h3>
            <EmptyState
              icon="📤"
              title="暂无上传记录"
              description="点击下方按钮开始上传"
              action-text="开始上传"
              show-action
              @action="handleAction"
            />
          </div>

          <!-- 搜索无结果 -->
          <div class="demo-card">
            <h3 class="demo-title">搜索无结果</h3>
            <EmptyState
              icon="🔍"
              title="未找到相关内容"
              description="尝试使用其他关键词搜索"
            />
          </div>
        </div>
      </section>

      <!-- 实际使用示例 -->
      <section class="demo-section">
        <h2 class="text-h2 mb-4">实际使用示例</h2>

        <div class="demo-card">
          <div class="demo-controls">
            <button
              v-for="state in states"
              :key="state"
              :class="['state-button', { active: currentState === state }]"
              @click="currentState = state"
            >
              {{ stateLabels[state] }}
            </button>
          </div>

          <div class="demo-content">
            <LoadingState v-if="currentState === 'loading'" message="加载资产列表..." />

            <ErrorState
              v-else-if="currentState === 'error'"
              :error="sampleError"
              @retry="currentState = 'loading'"
            />

            <EmptyState
              v-else-if="currentState === 'empty'"
              icon="🎬"
              title="还没有媒体资产"
              description="开始上传您的第一个视频、图片或音频文件"
              action-text="上传资产"
              show-action
              @action="handleAction"
            />

            <div v-else class="success-content">
              <div class="success-icon">✅</div>
              <p class="success-text">数据加载成功！</p>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { LoadingState, ErrorState, EmptyState } from '@/components/common'

type State = 'loading' | 'error' | 'empty' | 'success'

const currentState = ref<State>('loading')
const states: State[] = ['loading', 'error', 'empty', 'success']

const stateLabels: Record<State, string> = {
  loading: '加载中',
  error: '错误',
  empty: '空状态',
  success: '成功'
}

const sampleError = new Error('无法连接到服务器，请检查网络连接')

const handleRetry = () => {
  ElMessage.success('点击了重试按钮')
  currentState.value = 'loading'
  setTimeout(() => {
    currentState.value = 'success'
  }, 1000)
}

const handleAction = () => {
  ElMessage.success('点击了操作按钮')
}
</script>

<style scoped>
.state-demo-page {
  padding: 32px 0;
  background: #FAFAFA;
  min-height: 100vh;
}

.demo-section {
  margin-bottom: 48px;
}

.demo-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 24px;
}

.demo-card {
  background: white;
  border: 1px solid #E5E5E5;
  border-radius: 8px;
  overflow: hidden;
}

.demo-title {
  padding: 16px 24px;
  background: #F5F5F5;
  border-bottom: 1px solid #E5E5E5;
  font-size: 14px;
  font-weight: 600;
  color: #262626;
  margin: 0;
}

.demo-controls {
  padding: 16px 24px;
  background: #F5F5F5;
  border-bottom: 1px solid #E5E5E5;
  display: flex;
  gap: 8px;
}

.state-button {
  padding: 8px 16px;
  background: white;
  border: 1px solid #E5E5E5;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  color: #525252;
  cursor: pointer;
  transition: all 150ms;
}

.state-button:hover {
  border-color: #4F5B93;
  color: #4F5B93;
}

.state-button.active {
  background: #4F5B93;
  border-color: #4F5B93;
  color: white;
}

.demo-content {
  min-height: 400px;
}

.success-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

.success-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.success-text {
  font-size: 16px;
  font-weight: 500;
  color: #10B981;
}

@media (max-width: 768px) {
  .demo-grid {
    grid-template-columns: 1fr;
  }
}
</style>
