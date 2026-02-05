<template>
  <teleport v-if="fullscreen" to="body">
    <transition name="gv-loading-fade">
      <div v-show="loading" :class="fullscreenClasses" :style="fullscreenStyle">
        <div :class="spinnerClasses">
          <component :is="loaderComponent" />
          <p v-if="text" class="gv-loading__text">{{ text }}</p>
        </div>
      </div>
    </transition>
  </teleport>
  
  <div v-else :class="containerClasses">
    <!-- 遮罩层 -->
    <transition name="gv-loading-fade">
      <div v-show="loading" class="gv-loading__mask">
        <div :class="spinnerClasses">
          <component :is="loaderComponent" />
          <p v-if="text" class="gv-loading__text">{{ text }}</p>
        </div>
      </div>
    </transition>
    
    <!-- 内容 -->
    <div :class="{ 'opacity-50 pointer-events-none': loading }">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, watch, onUnmounted } from 'vue'
import { cn } from '@/utils/cn'
import type { LoadingProps } from './types'
import SpinnerLoader from './loaders/SpinnerLoader.vue'
import CircularLoader from './loaders/CircularLoader.vue'
import DotsLoader from './loaders/DotsLoader.vue'
import BarsLoader from './loaders/BarsLoader.vue'

const props = withDefaults(defineProps<LoadingProps>(), {
  loading: true,
  type: 'circular',
  size: 'medium',
  fullscreen: false,
  lock: true,
  color: 'primary'
})

// 加载器组件映射
const loaderComponents = {
  spinner: SpinnerLoader,
  circular: CircularLoader,
  dots: DotsLoader,
  bars: BarsLoader
}

const loaderComponent = computed(() => loaderComponents[props.type])

// 容器类名
const containerClasses = computed(() => {
  return cn('gv-loading', 'relative', props.customClass)
})

// 全屏加载类名
const fullscreenClasses = computed(() => {
  return cn(
    'gv-loading gv-loading--fullscreen',
    'fixed inset-0 z-[2000]',
    'flex items-center justify-center',
    props.customClass
  )
})

// 全屏加载样式
const fullscreenStyle = computed(() => {
  return {
    backgroundColor: props.background || 'rgba(255, 255, 255, 0.9)'
  }
})

// 加载器类名
const spinnerClasses = computed(() => {
  const base = ['gv-loading__spinner', 'flex flex-col items-center justify-center gap-3']
  
  // 尺寸
  const sizeClasses = {
    small: 'w-8 h-8',
    medium: 'w-12 h-12',
    large: 'w-16 h-16'
  }
  
  // 颜色
  const colorClasses = {
    primary: 'text-primary-600',
    secondary: 'text-neutral-600',
    success: 'text-success-600',
    error: 'text-error-600',
    warning: 'text-warning-600',
    info: 'text-info-600',
    white: 'text-white'
  }
  
  return cn(base, sizeClasses[props.size], colorClasses[props.color])
})

// 监听 fullscreen 和 lock
watch(
  () => props.loading && props.fullscreen && props.lock,
  (locked) => {
    if (locked) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }
  },
  { immediate: true }
)

// 组件卸载时恢复滚动
onUnmounted(() => {
  if (props.fullscreen && props.lock) {
    document.body.style.overflow = ''
  }
})
</script>

<style scoped>
.gv-loading__mask {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(255, 255, 255, 0.9);
  border-radius: inherit;
  z-index: 10;
}

.gv-loading__text {
  font-size: 0.875rem;
  color: currentColor;
  margin-top: 8px;
}

/* 淡入淡出动画 */
.gv-loading-fade-enter-active,
.gv-loading-fade-leave-active {
  transition: opacity 0.3s;
}

.gv-loading-fade-enter-from,
.gv-loading-fade-leave-to {
  opacity: 0;
}

/* 深色模式 */
.dark .gv-loading__mask {
  background-color: rgba(30, 41, 59, 0.9);
}

.dark .gv-loading--fullscreen {
  background-color: rgba(30, 41, 59, 0.9) !important;
}
</style>
