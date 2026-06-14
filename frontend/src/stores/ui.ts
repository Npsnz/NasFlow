import { defineStore } from 'pinia'

export interface Toast {
  id: number
  message: string
  type: 'success' | 'error' | 'info'
}

export const useUIStore = defineStore('ui', {
  state: () => ({
    sidebarOpen: false, // Hidden by default on mobile
    activeView: 'board', // board | calendar | list | settings
    theme: localStorage.getItem('tf_theme') || 'system', // light | dark | system
    language: (localStorage.getItem('tf_lang') || 'th') as 'th' | 'en', // th | en
    toasts: [] as Toast[],
    toastIdCounter: 0,
    focusedTaskId: null as number | null
  }),

  actions: {
    toggleSidebar() {
      this.sidebarOpen = !this.sidebarOpen
    },
    setSidebarOpen(open: boolean) {
      this.sidebarOpen = open
    },
    setActiveView(view: string) {
      this.activeView = view
    },
    setTheme(theme: 'light' | 'dark' | 'system') {
      this.theme = theme
      localStorage.setItem('tf_theme', theme)
      this.applyTheme()
    },
    setLanguage(lang: 'th' | 'en') {
      this.language = lang
      localStorage.setItem('tf_lang', lang)
    },
    showToast(message: string, type: 'success' | 'error' | 'info' = 'success') {
      const id = ++this.toastIdCounter
      this.toasts.push({ id, message, type })
      setTimeout(() => {
        this.toasts = this.toasts.filter(t => t.id !== id)
      }, 3000) // Auto dismiss after 3 seconds
    },
    setFocusedTask(id: number | null) {
      this.focusedTaskId = id
    },
    applyTheme() {
      if (typeof window === 'undefined') return
      const root = document.documentElement
      let isDark = this.theme === 'dark'

      if (this.theme === 'system') {
        isDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      }

      if (isDark) {
        root.classList.add('dark')
      } else {
        root.classList.remove('dark')
      }
    }
  }
})
