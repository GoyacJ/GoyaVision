/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}'
  ],
  
  darkMode: 'class',
  
  theme: {
    extend: {
      colors: {
        // 主色：克制的蓝灰色
        primary: {
          DEFAULT: '#4F5B93',
          50: '#F5F6FA',
          100: '#EBEDF5',
          200: '#D4D8E8',
          300: '#B3BAD5',
          400: '#8A94B8',
          500: '#4F5B93',
          600: '#3E4A7A',
          700: '#2F3A61',
          800: '#232D4B',
          900: '#1A2238',
          950: '#131A2B'
        },
        // 功能色：保持不变（使用 Tailwind 默认配色）
        success: {
          DEFAULT: '#10b981',
          50: '#f0fdf4',
          100: '#dcfce7',
          200: '#bbf7d0',
          300: '#86efac',
          400: '#4ade80',
          500: '#22c55e',
          600: '#10b981',
          700: '#059669',
          800: '#047857',
          900: '#065f46',
          950: '#064e3b'
        },
        error: {
          DEFAULT: '#ef4444',
          50: '#fef2f2',
          100: '#fee2e2',
          200: '#fecaca',
          300: '#fca5a5',
          400: '#f87171',
          500: '#ef4444',
          600: '#dc2626',
          700: '#b91c1c',
          800: '#991b1b',
          900: '#7f1d1d',
          950: '#450a0a'
        },
        warning: {
          DEFAULT: '#f59e0b',
          50: '#fffbeb',
          100: '#fef3c7',
          200: '#fde68a',
          300: '#fcd34d',
          400: '#fbbf24',
          500: '#f59e0b',
          600: '#d97706',
          700: '#b45309',
          800: '#92400e',
          900: '#78350f',
          950: '#451a03'
        },
        info: {
          DEFAULT: '#3b82f6',
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a',
          950: '#172554'
        },
        // 中性色：纯灰度系统（无色相偏向）
        neutral: {
          DEFAULT: '#737373',
          50: '#FAFAFA',
          100: '#F5F5F5',
          200: '#E5E5E5',
          300: '#D4D4D4',
          400: '#A3A3A3',
          500: '#737373',
          600: '#525252',
          700: '#404040',
          800: '#262626',
          900: '#171717',
          950: '#0A0A0A'
        },
        // 表面色（背景、卡片）
        surface: {
          DEFAULT: '#FFFFFF',
          dim: '#FAFAFA',
          container: '#F5F5F5',
          'container-high': '#E5E5E5',
          dark: '#262626',
          'dark-container': '#171717'
        },
        // 文字色
        text: {
          primary: '#262626',
          secondary: '#525252',
          tertiary: '#737373',
          placeholder: '#A3A3A3',
          disabled: '#D4D4D4',
          inverse: '#FFFFFF'
        }
      },
      
      // 阴影：极简、功能性（移除彩色阴影）
      boxShadow: {
        none: 'none',
        sm: '0 1px 2px rgba(0, 0, 0, 0.04)',
        DEFAULT: '0 1px 3px rgba(0, 0, 0, 0.06)',
        md: '0 4px 6px rgba(0, 0, 0, 0.07)',
        lg: '0 10px 15px rgba(0, 0, 0, 0.08)',
        xl: '0 20px 25px rgba(0, 0, 0, 0.10)',
        '2xl': '0 25px 50px rgba(0, 0, 0, 0.12)',
        inner: 'inset 0 2px 4px rgba(0, 0, 0, 0.04)'
      },
      
      // 圆角：适度减小，更克制
      borderRadius: {
        none: '0',
        sm: '0.25rem',     // 4px
        DEFAULT: '0.375rem',// 6px - 基准
        md: '0.5rem',      // 8px
        lg: '0.75rem',     // 12px
        xl: '1rem',        // 16px
        '2xl': '1.5rem',   // 24px
        full: '9999px'
      },
      
      // 字体家族
      fontFamily: {
        sans: [
          '-apple-system',
          'BlinkMacSystemFont',
          'Segoe UI',
          'Roboto',
          'Helvetica Neue',
          'Arial',
          'Noto Sans SC',
          'sans-serif'
        ],
        mono: [
          'SF Mono',
          'Menlo',
          'Monaco',
          'Consolas',
          'Liberation Mono',
          'Courier New',
          'monospace'
        ]
      },

      // 字号：正文 15px（而非 16px）
      fontSize: {
        xs: ['0.75rem', { lineHeight: '1rem' }],       // 12px
        sm: ['0.8125rem', { lineHeight: '1.25rem' }],  // 13px
        base: ['0.9375rem', { lineHeight: '1.5rem' }], // 15px - 正文
        lg: ['1.125rem', { lineHeight: '1.75rem' }],   // 18px
        xl: ['1.25rem', { lineHeight: '1.75rem' }],    // 20px
        '2xl': ['1.5rem', { lineHeight: '2rem' }],     // 24px
        '3xl': ['1.875rem', { lineHeight: '2.25rem' }],// 30px
        '4xl': ['2rem', { lineHeight: '2.5rem' }],     // 32px
        '5xl': ['2.5rem', { lineHeight: '1.2' }]       // 40px
      },

      // 字距：标题紧凑，标签宽松
      letterSpacing: {
        tighter: '-0.02em',
        tight: '-0.01em',
        normal: '0',
        wide: '0.025em',
        wider: '0.05em'
      },

      // 动画曲线：快速、自然
      transitionTimingFunction: {
        'emphasized': 'cubic-bezier(0.2, 0, 0, 1)',
        'emphasized-decelerate': 'cubic-bezier(0.05, 0.7, 0.1, 1)',
        'emphasized-accelerate': 'cubic-bezier(0.3, 0, 0.8, 0.15)',
        'standard': 'cubic-bezier(0.4, 0, 0.2, 1)'
      },
      
      // 动画时长：更快、不干扰
      transitionDuration: {
        fast: '150ms',     // 快速反馈（hover）
        normal: '200ms',   // 标准过渡（大部分交互）
        slow: '300ms'      // 页面切换（最长）
      }
    }
  },
  
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('@tailwindcss/container-queries'),
    
    function({ addComponents }) {
      addComponents({
        '.surface': {
          '@apply bg-surface rounded-lg shadow-sm': {}
        },
        '.surface-container': {
          '@apply bg-surface-container rounded-lg': {}
        },
        '.text-ellipsis-1': {
          'overflow': 'hidden',
          'text-overflow': 'ellipsis',
          'white-space': 'nowrap'
        },
        '.text-ellipsis-2': {
          'display': '-webkit-box',
          '-webkit-line-clamp': '2',
          '-webkit-box-orient': 'vertical',
          'overflow': 'hidden'
        }
      })
    }
  ]
}
