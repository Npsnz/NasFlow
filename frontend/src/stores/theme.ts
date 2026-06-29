import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const isDark = ref(localStorage.getItem('theme') === 'dark')

  const setTheme = (dark: boolean) => {
    isDark.value = dark
    localStorage.setItem('theme', dark ? 'dark' : 'light')
    document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light')
  }

  const toggleTheme = () => {
    setTheme(!isDark.value)
  }

  // Initialize theme on store creation
  watch(
    () => isDark.value,
    (dark) => {
      localStorage.setItem('theme', dark ? 'dark' : 'light')
      document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light')
    }
  )

  return {
    isDark,
    setTheme,
    toggleTheme,
  }
})
