/**
 * 主题切换 Composable
 * 支持 light / dark 主题切换
 */

import { ref, computed, watch } from 'vue'
import { useLocalStorage } from '@vueuse/core'

export type Theme = 'light' | 'dark' | 'system'

const THEME_STORAGE_KEY = 'goyavision-theme'

// 全局主题状态
const theme = useLocalStorage<Theme>(THEME_STORAGE_KEY, 'system')

// 计算实际应用的主题
const actualTheme = computed(() => {
  if (theme.value === 'system') {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
  }
  return theme.value
})

// 应用主题到 DOM
function applyTheme() {
  const root = document.documentElement
  
  if (actualTheme.value === 'dark') {
    root.classList.add('dark')
  } else {
    root.classList.remove('dark')
  }
}

// 监听系统主题变化
if (window.matchMedia) {
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
    if (theme.value === 'system') {
      applyTheme()
    }
  })
}

export function useTheme() {
  // 初始化时应用主题
  applyTheme()
  
  // 监听主题变化
  watch(theme, applyTheme)
  
  /**
   * 切换主题
   */
  const toggleTheme = () => {
    theme.value = actualTheme.value === 'dark' ? 'light' : 'dark'
  }
  
  /**
   * 设置主题
   */
  const setTheme = (newTheme: Theme) => {
    theme.value = newTheme
  }
  
  /**
   * 是否为深色模式
   */
  const isDark = computed(() => actualTheme.value === 'dark')
  
  return {
    theme,
    actualTheme,
    isDark,
    toggleTheme,
    setTheme
  }
}
