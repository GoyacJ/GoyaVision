<template>
  <teleport to="body">
    <transition name="gv-drawer-fade" @after-enter="handleOpened" @after-leave="handleClosed">
      <div
        v-show="modelValue"
        class="gv-drawer"
        :style="{ zIndex }"
        @click.self="handleMaskClick"
      >
        <!-- 遮罩层 -->
        <div class="gv-drawer__mask" @click="handleMaskClick"></div>
        
        <!-- 抽屉容器 -->
        <transition :name="`gv-drawer-${direction}`">
          <div
            v-show="modelValue"
            :class="drawerClasses"
            :style="drawerStyle"
            role="dialog"
            aria-modal="true"
            :aria-labelledby="title ? 'gv-drawer-title' : undefined"
          >
            <!-- 头部 -->
            <div v-if="title || $slots.header" class="gv-drawer__header">
              <slot name="header">
                <h3 id="gv-drawer-title" class="gv-drawer__title">
                  {{ title }}
                </h3>
              </slot>
              
              <!-- 关闭按钮 -->
              <button
                v-if="showClose"
                class="gv-drawer__close"
                type="button"
                aria-label="关闭"
                @click="handleClose"
              >
                <el-icon :size="20">
                  <Close />
                </el-icon>
              </button>
            </div>
            
            <!-- 内容 -->
            <div class="gv-drawer__body scrollbar-thin">
              <slot />
            </div>
            
            <!-- 底部 -->
            <div v-if="showFooter || $slots.footer" class="gv-drawer__footer">
              <slot name="footer">
                <div class="flex justify-end gap-3">
                  <GvButton
                    v-if="showCancel"
                    variant="tonal"
                    @click="handleCancel"
                  >
                    {{ cancelText }}
                  </GvButton>
                  <GvButton
                    v-if="showConfirm"
                    variant="filled"
                    :loading="confirmLoading"
                    @click="handleConfirm"
                  >
                    {{ confirmText }}
                  </GvButton>
                </div>
              </slot>
            </div>
          </div>
        </transition>
      </div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { computed, watch, onMounted, onUnmounted } from 'vue'
import { cn } from '@/utils/cn'
import { GvButton } from '@/components'
import type { DrawerProps, DrawerEmits } from './types'

const props = withDefaults(defineProps<DrawerProps>(), {
  direction: 'right',
  size: 'medium',
  showClose: true,
  closeOnClickModal: true,
  closeOnPressEscape: true,
  destroyOnClose: false,
  showFooter: true,
  confirmText: '确定',
  cancelText: '取消',
  showConfirm: true,
  showCancel: true,
  confirmLoading: false,
  zIndex: 1000
})

const emit = defineEmits<DrawerEmits>()

// 抽屉类名
const drawerClasses = computed(() => {
  const base = [
    'gv-drawer__wrapper',
    'fixed bg-white shadow-2xl',
    'flex flex-col',
    'transition-all duration-300'
  ]
  
  // 方向
  const directionClasses = {
    left: 'left-0 top-0 bottom-0',
    right: 'right-0 top-0 bottom-0',
    top: 'top-0 left-0 right-0',
    bottom: 'bottom-0 left-0 right-0'
  }
  
  return cn(base, directionClasses[props.direction], props.customClass)
})

// 抽屉样式
const drawerStyle = computed(() => {
  const isHorizontal = props.direction === 'left' || props.direction === 'right'
  
  // 自定义宽度/高度优先
  if (props.width && isHorizontal) {
    return { width: props.width }
  }
  if (props.height && !isHorizontal) {
    return { height: props.height }
  }
  
  // 预设尺寸
  const sizes = {
    left: {
      small: { width: '300px' },
      medium: { width: '450px' },
      large: { width: '600px' },
      full: { width: '100vw' }
    },
    right: {
      small: { width: '300px' },
      medium: { width: '450px' },
      large: { width: '600px' },
      full: { width: '100vw' }
    },
    top: {
      small: { height: '200px' },
      medium: { height: '300px' },
      large: { height: '400px' },
      full: { height: '100vh' }
    },
    bottom: {
      small: { height: '200px' },
      medium: { height: '300px' },
      large: { height: '400px' },
      full: { height: '100vh' }
    }
  }
  
  return sizes[props.direction][props.size]
})

// 处理遮罩点击
const handleMaskClick = () => {
  if (props.closeOnClickModal) {
    handleClose()
  }
}

// 处理关闭
const handleClose = () => {
  emit('update:modelValue', false)
  emit('close')
}

// 处理确认
const handleConfirm = () => {
  emit('confirm')
}

// 处理取消
const handleCancel = () => {
  emit('cancel')
  handleClose()
}

// 打开动画结束
const handleOpened = () => {
  emit('opened')
}

// 关闭动画结束
const handleClosed = () => {
  emit('closed')
}

// ESC 键关闭
const handleEscapeKeyDown = (event: KeyboardEvent) => {
  if (event.key === 'Escape' && props.closeOnPressEscape && props.modelValue) {
    handleClose()
  }
}

// 监听显示状态
watch(
  () => props.modelValue,
  (value) => {
    if (value) {
      emit('open')
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }
  }
)

// 生命周期
onMounted(() => {
  document.addEventListener('keydown', handleEscapeKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscapeKeyDown)
  document.body.style.overflow = ''
})
</script>

<style scoped>
.gv-drawer {
  position: fixed;
  inset: 0;
}

.gv-drawer__mask {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.5);
}

.gv-drawer__wrapper {
  z-index: 1;
}

.gv-drawer__header {
  position: relative;
  padding: 24px 24px 16px;
  border-bottom: 1px solid rgb(241 245 249);
  flex-shrink: 0;
}

.gv-drawer__title {
  font-size: 1.25rem;
  font-weight: 600;
  line-height: 1.75rem;
  color: rgb(15 23 42);
  margin: 0;
  padding-right: 32px;
}

.gv-drawer__close {
  position: absolute;
  top: 20px;
  right: 20px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  color: rgb(100 116 139);
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s;
}

.gv-drawer__close:hover {
  background-color: rgb(241 245 249);
  color: rgb(15 23 42);
}

.gv-drawer__body {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  color: rgb(71 85 105);
  line-height: 1.6;
}

.gv-drawer__footer {
  padding: 16px 24px 24px;
  border-top: 1px solid rgb(241 245 249);
  flex-shrink: 0;
}

/* 遮罩淡入淡出 */
.gv-drawer-fade-enter-active,
.gv-drawer-fade-leave-active {
  transition: opacity 0.3s;
}

.gv-drawer-fade-enter-from,
.gv-drawer-fade-leave-to {
  opacity: 0;
}

/* 从左侧滑入 */
.gv-drawer-left-enter-active,
.gv-drawer-left-leave-active {
  transition: transform 0.3s cubic-bezier(0.2, 0, 0, 1);
}

.gv-drawer-left-enter-from {
  transform: translateX(-100%);
}

.gv-drawer-left-leave-to {
  transform: translateX(-100%);
}

/* 从右侧滑入 */
.gv-drawer-right-enter-active,
.gv-drawer-right-leave-active {
  transition: transform 0.3s cubic-bezier(0.2, 0, 0, 1);
}

.gv-drawer-right-enter-from {
  transform: translateX(100%);
}

.gv-drawer-right-leave-to {
  transform: translateX(100%);
}

/* 从顶部滑入 */
.gv-drawer-top-enter-active,
.gv-drawer-top-leave-active {
  transition: transform 0.3s cubic-bezier(0.2, 0, 0, 1);
}

.gv-drawer-top-enter-from {
  transform: translateY(-100%);
}

.gv-drawer-top-leave-to {
  transform: translateY(-100%);
}

/* 从底部滑入 */
.gv-drawer-bottom-enter-active,
.gv-drawer-bottom-leave-active {
  transition: transform 0.3s cubic-bezier(0.2, 0, 0, 1);
}

.gv-drawer-bottom-enter-from {
  transform: translateY(100%);
}

.gv-drawer-bottom-leave-to {
  transform: translateY(100%);
}

/* 深色模式 */
.dark .gv-drawer__wrapper {
  @apply bg-surface-dark;
}

.dark .gv-drawer__header {
  @apply border-neutral-700;
}

.dark .gv-drawer__title {
  @apply text-text-inverse;
}

.dark .gv-drawer__body {
  @apply text-neutral-300;
}

.dark .gv-drawer__footer {
  @apply border-neutral-700;
}

.dark .gv-drawer__close {
  @apply text-neutral-400;
}

.dark .gv-drawer__close:hover {
  @apply bg-neutral-800 text-text-inverse;
}
</style>
