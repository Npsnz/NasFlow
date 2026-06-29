<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import ToastContainer from '@/components/ToastContainer.vue'

const authStore = useAuthStore()
const router = useRouter()

onMounted(async () => {
  // Check if token exists and validate it
  if (authStore.isAuthenticated) {
    await authStore.fetchMe()
  }

  // Listen for logout event (when token refresh fails)
  window.addEventListener('auth-logout', handleAuthLogout)
})

onUnmounted(() => {
  window.removeEventListener('auth-logout', handleAuthLogout)
})

const handleAuthLogout = async () => {
  await authStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="min-h-screen w-screen bg-bg-light dark:bg-bg-dark font-sans select-none">
    <ToastContainer />
    <router-view />
  </div>
</template>
