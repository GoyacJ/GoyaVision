<template>
  <el-container>
    <el-header>
      <div class="header-content">
        <h1>GoyaVision</h1>
        <span class="subtitle">AI 视频流分析处理平台</span>
      </div>
    </el-header>
    <el-container>
      <el-aside width="200px">
        <el-menu
          :default-active="activeMenu"
          router
          mode="vertical"
        >
          <el-menu-item index="/streams">
            <span>视频流</span>
          </el-menu-item>
          <el-menu-item index="/algorithms">
            <span>算法管理</span>
          </el-menu-item>
          <el-menu-item index="/inference-results">
            <span>推理结果</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const activeMenu = computed(() => route.path)
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  height: 100vh;
}

.el-header {
  background-color: #409eff;
  color: white;
  display: flex;
  align-items: center;
  padding: 0 20px;
}

.header-content h1 {
  font-size: 24px;
  margin: 0;
}

.subtitle {
  margin-left: 10px;
  font-size: 14px;
  opacity: 0.9;
}

.el-aside {
  background-color: #f5f5f5;
}

.el-main {
  padding: 20px;
  background-color: #fff;
}
</style>
