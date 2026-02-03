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
        primary: {
          DEFAULT: '#667eea',
          50: '#f5f7ff',
          100: '#ebedff',
          200: '#d6dcff',
          300: '#b3bdff',
          400: '#8d9eff',
          500: '#667eea',
          600: '#5568d3',
          700: '#4553bd',
          800: '#3640a6',
          900: '#2a3290',
          950: '#1f2673'
        },
        secondary: {
          DEFAULT: '#764ba2',
          50: '#f9f5fc',
          100: '#f3ebf9',
          200: '#e7d7f3',
          300: '#d4b8e9',
          400: '#b88ed9',
          500: '#9d64c9',
          600: '#764ba2',
          700: '#65408b',
          800: '#543574',
          900: '#432a5d',
          950: '#321f46'
        },
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
        neutral: {
          DEFAULT: '#64748b',
          50: '#f8fafc',
          100: '#f1f5f9',
          200: '#e2e8f0',
          300: '#cbd5e1',
          400: '#94a3b8',
          500: '#64748b',
          600: '#475569',
          700: '#334155',
          800: '#1e293b',
          900: '#0f172a',
          950: '#020617'
        },
        surface: {
          DEFAULT: '#ffffff',
          dark: '#1e293b',
          dim: '#f8fafc',
          bright: '#ffffff',
          container: '#f1f5f9',
          'container-high': '#e2e8f0',
          'container-highest': '#cbd5e1'
        },
        text: {
          primary: '#0f172a',
          secondary: '#475569',
          tertiary: '#64748b',
          disabled: '#cbd5e1',
          inverse: '#ffffff'
        }
      },
      
      boxShadow: {
        sm: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        DEFAULT: '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1)',
        md: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1)',
        lg: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1)',
        xl: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1)',
        '2xl': '0 25px 50px -12px rgba(0, 0, 0, 0.25)',
        inner: 'inset 0 2px 4px 0 rgba(0, 0, 0, 0.05)',
        primary: '0 8px 16px -4px rgba(102, 126, 234, 0.3)',
        secondary: '0 8px 16px -4px rgba(118, 75, 162, 0.3)',
        success: '0 8px 16px -4px rgba(16, 185, 129, 0.3)',
        error: '0 8px 16px -4px rgba(239, 68, 68, 0.3)'
      },
      
      borderRadius: {
        DEFAULT: '0.5rem',
        sm: '0.25rem',
        md: '0.75rem',
        lg: '1rem',
        xl: '1.5rem',
        '2xl': '2rem',
        '3xl': '3rem'
      },
      
      transitionTimingFunction: {
        'emphasized': 'cubic-bezier(0.2, 0, 0, 1)',
        'emphasized-decelerate': 'cubic-bezier(0.05, 0.7, 0.1, 1)',
        'emphasized-accelerate': 'cubic-bezier(0.3, 0, 0.8, 0.15)',
        'standard': 'cubic-bezier(0.2, 0, 0, 1)'
      },
      
      transitionDuration: {
        'short1': '50ms',
        'short2': '100ms',
        'short3': '150ms',
        'short4': '200ms',
        'medium1': '250ms',
        'medium2': '300ms',
        'medium3': '350ms',
        'medium4': '400ms',
        'long1': '450ms',
        'long2': '500ms',
        'long3': '550ms',
        'long4': '600ms',
        'extra-long1': '700ms',
        'extra-long2': '800ms',
        'extra-long3': '900ms',
        'extra-long4': '1000ms'
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
