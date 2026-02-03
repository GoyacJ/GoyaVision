<template>
  <teleport to="body">
    <transition name="gv-modal-fade" @after-enter="handleOpened" @after-leave="handleClosed">
      <div
        v-show="modelValue"
        class="gv-modal"
        :style="{ zIndex }"
        @click.self="handleMaskClick"
      >
        <!-- 遮罩层 -->
        <div class="gv-modal__mask" @click="handleMaskClick"></div>
        
        <!-- 模态框容器 -->
        <transition name="gv-modal-slide">
          <div
            v-show="modelValue"
            :class="modalClasses"
            role="dialog"
            aria-modal="true"
            :aria-labelledby="title ? 'gv-modal-title' : undefined"
          >
            <!-- 头部 -->
            <div v-if="title || $slots.header" class="gv-modal__header">
              <slot name="header">
                <h3 id="gv-modal-title" class="gv-modal__title">
                  {{ title }}
                </h3>
              </slot>
              
              <!-- 关闭按钮 -->
              <button
                v-if="showClose"
                class="gv-modal__close"
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
            <div class="gv-modal__body">
              <slot />
            </div>
            
            <!-- 底部 -->
            <div v-if="showFooter || $slots.footer" class="gv-modal__footer">
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
import type { ModalProps, ModalEmits } from './types'

const props = withDefaults(defineProps<ModalProps>(), {
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
  center: false,
  zIndex: 1000
})

const emit = defineEmits<ModalEmits>()

// 模态框类名
const modalClasses = computed(() => {
  const base = [
    'gv-modal__wrapper',
    'relative bg-white rounded-2xl shadow-2xl',
    'flex flex-col',
    'max-h-[90vh]',
    'transition-all duration-300'
  ]
  
  // 尺寸
  const sizeClasses = {
    small: 'w-[400px]',
    medium: 'w-[600px]',
    large: 'w-[800px]',
    full: 'w-[90vw] h-[90vh]'
  }
  
  // 居中
  const centerClass = props.center ? 'text-center' : ''
  
  return cn(
    base,
    sizeClasses[props.size],
    centerClass,
    props.customClass
  )
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
.gv-modal {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  overflow: auto;
}

.gv-modal__mask {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
}

.gv-modal__wrapper {
  position: relative;
  z-index: 1;
  margin: auto;
}

.gv-modal__header {
  position: relative;
  padding: 24px 24px 16px;
  border-bottom: 1px solid rgb(241 245 249);
}

.gv-modal__title {
  font-size: 1.25rem;
  font-weight: 600;
  line-height: 1.75rem;
  color: rgb(15 23 42);
  margin: 0;
  padding-right: 32px;
}

.gv-modal__close {
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

.gv-modal__close:hover {
  background-color: rgb(241 245 249);
  color: rgb(15 23 42);
}

.gv-modal__body {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  color: rgb(71 85 105);
  line-height: 1.6;
}

.gv-modal__footer {
  padding: 16px 24px 24px;
  border-top: 1px solid rgb(241 245 249);
}

/* 模态框淡入淡出动画 */
.gv-modal-fade-enter-active,
.gv-modal-fade-leave-active {
  transition: opacity 0.3s;
}

.gv-modal-fade-enter-from,
.gv-modal-fade-leave-to {
  opacity: 0;
}

/* 模态框滑入动画 */
.gv-modal-slide-enter-active,
.gv-modal-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.2, 0, 0, 1);
}

.gv-modal-slide-enter-from {
  opacity: 0;
  transform: translateY(-20px) scale(0.95);
}

.gv-modal-slide-leave-to {
  opacity: 0;
  transform: translateY(20px) scale(0.95);
}

/* 深色模式 */
.dark .gv-modal__wrapper {
  @apply bg-surface-dark;
}

.dark .gv-modal__header {
  @apply border-neutral-700;
}

.dark .gv-modal__title {
  @apply text-text-inverse;
}

.dark .gv-modal__body {
  @apply text-neutral-300;
}

.dark .gv-modal__footer {
  @apply border-neutral-700;
}

.dark .gv-modal__close {
  @apply text-neutral-400;
}

.dark .gv-modal__close:hover {
  @apply bg-neutral-800 text-text-inverse;
}
</style>
