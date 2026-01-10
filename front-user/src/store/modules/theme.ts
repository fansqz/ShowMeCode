import { defineStore } from 'pinia'

export type ThemeMode = 'light' | 'dark'

type ThemeState = {
  mode: ThemeMode
  editorTheme: string
}

const useThemeStore = defineStore('theme', {
  state: (): ThemeState => ({
    mode: 'light',
    editorTheme: 'BlulocoLight',
  }),

  getters: {
    isDark: (state): boolean => state.mode === 'dark',
    isLight: (state): boolean => state.mode === 'light',
  },

  actions: {
    setTheme(mode: ThemeMode) {
      this.mode = mode
      localStorage.setItem('theme', mode)

      // 自动调整编辑器主题
      this.autoSetEditorTheme()

      this.applyTheme()
    },

    toggleTheme() {
      this.setTheme(this.mode === 'light' ? 'dark' : 'light')
    },

    setEditorTheme(theme: string) {
      this.editorTheme = theme
      localStorage.setItem('editorTheme', theme)
    },

    // 根据页面主题自动设置编辑器主题
    autoSetEditorTheme() {
      const defaultTheme = this.mode === 'dark' ? 'Monokai' : 'BlulocoLight'
      this.setEditorTheme(defaultTheme)
    },

    applyTheme() {
      const html = document.documentElement
      if (this.mode === 'dark') {
        html.classList.add('dark')
        html.classList.remove('light')
      } else {
        html.classList.remove('dark')
        html.classList.add('light')
      }
    },

    initTheme() {
      // 从localStorage读取保存的主题设置
      const savedTheme = localStorage.getItem('theme') as ThemeMode
      const savedEditorTheme = localStorage.getItem('editorTheme')

      if (savedTheme && ['light', 'dark'].includes(savedTheme)) {
        this.mode = savedTheme
      }

      if (savedEditorTheme) {
        this.editorTheme = savedEditorTheme
      } else {
        // 如果没有保存的编辑器主题，根据页面主题设置默认值
        this.autoSetEditorTheme()
      }

      this.applyTheme()
    },
  },
  persist: true,
})

export default useThemeStore
