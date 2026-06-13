<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTasksStore } from '@/stores/tasks'
import { useWorkspaceStore } from '@/stores/workspace'
import { useUIStore } from '@/stores/ui'
import AppSidebar from './AppSidebar.vue'
import AppTopbar from './AppTopbar.vue'
import ToastNotification from '@/components/ui/ToastNotification.vue'
import QuickAddModal from '@/components/task/TaskForm.vue' // We will build TaskForm to act as quick-add modal as well
import { LayoutDashboard, Calendar, ListTodo, Settings, Sparkles } from 'lucide-vue-next'

const router = useRouter()
const tasksStore = useTasksStore()
const workspaceStore = useWorkspaceStore()
const uiStore = useUIStore()

const activeView = computed(() => uiStore.activeView)
const showQuickAdd = ref(false)

// Handle Keyboard Shortcuts globally
const handleKeydown = (e: KeyboardEvent) => {
  // Ctrl+K or Cmd+K: Open Quick Add Command Palette
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
    e.preventDefault()
    showQuickAdd.value = true
  }
}

const handleOpenQuickAdd = () => {
  showQuickAdd.value = true
}

const dispatchOpenQuickAdd = () => {
  if (typeof window !== 'undefined') {
    window.dispatchEvent(new CustomEvent('open-quick-add'))
  }
}

const handleOverdueAlert = (e: Event) => {
  const count = (e as CustomEvent).detail
  uiStore.showToast(`คุณมีงานค้างส่งจำนวน ${count} งาน!`, 'error')
}

onMounted(() => {
  // Initialize SSE event streaming
  tasksStore.initSSE()
  // Fetch initial base tags/stats
  tasksStore.fetchStats()
  tasksStore.fetchTags()
  workspaceStore.fetchWorkspaces()

  // Bind key listeners
  window.addEventListener('keydown', handleKeydown)
  window.addEventListener('open-quick-add', handleOpenQuickAdd)
  window.addEventListener('tasks-overdue-alert', handleOverdueAlert)
  
  // Apply theme on load
  uiStore.applyTheme()
})

onUnmounted(() => {
  // Close SSE connection on unmount
  tasksStore.closeSSE()
  window.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('open-quick-add', handleOpenQuickAdd)
  window.removeEventListener('tasks-overdue-alert', handleOverdueAlert)
})

const selectMobileView = (view: 'board' | 'calendar' | 'tasks' | 'settings') => {
  uiStore.setActiveView(view)
  if (view === 'board') {
    router.push('/board')
  } else if (view === 'calendar') {
    router.push('/calendar')
  } else if (view === 'tasks') {
    router.push('/tasks')
  } else if (view === 'settings') {
    router.push('/settings')
  }
}
</script>

<template>
  <div class="h-screen flex overflow-hidden bg-white dark:bg-slate-950">
    <!-- App Sidebar (Collapsible, responsive) -->
    <AppSidebar />

    <!-- Main Content Wrapper -->
    <div class="flex-grow flex flex-col min-w-0 overflow-hidden relative">
      <!-- App Topbar -->
      <AppTopbar />

      <!-- View Area -->
      <main class="flex-grow overflow-y-auto pb-16 lg:pb-0">
        <router-view />
      </main>

      <!-- Mobile Bottom Navigation Bar -->
      <nav class="fixed bottom-0 inset-x-0 h-16 border-t border-border-light dark:border-border-dark bg-white dark:bg-surface-dark flex items-center justify-around px-4 lg:hidden z-30 shadow-lg">
        <button
          @click="selectMobileView('board')"
          class="flex flex-col items-center justify-center space-y-0.5 text-slate-500 hover:text-brand-500 min-h-[44px] min-w-[44px] px-2 py-1 rounded-lg"
          :class="{ 'text-brand-500 font-semibold': activeView === 'board' }"
        >
          <LayoutDashboard class="w-5 h-5" />
          <span class="text-[9px]">กระดาน</span>
        </button>

        <button
          @click="selectMobileView('calendar')"
          class="flex flex-col items-center justify-center space-y-0.5 text-slate-500 hover:text-brand-500 min-h-[44px] min-w-[44px] px-2 py-1 rounded-lg"
          :class="{ 'text-brand-500 font-semibold': activeView === 'calendar' }"
        >
          <Calendar class="w-5 h-5" />
          <span class="text-[9px]">ปฏิทิน</span>
        </button>

        <!-- Floating mobile Quick Add trigger in bottom bar -->
        <button
          @click="dispatchOpenQuickAdd"
          class="flex items-center justify-center w-11 h-11 bg-brand-500 text-white rounded-full shadow-lg shadow-brand-500/30 transform -translate-y-2 select-none active:scale-95"
        >
          <Sparkles class="w-5 h-5 animate-pulse" />
        </button>

        <button
          @click="selectMobileView('tasks')"
          class="flex flex-col items-center justify-center space-y-0.5 text-slate-500 hover:text-brand-500 min-h-[44px] min-w-[44px] px-2 py-1 rounded-lg"
          :class="{ 'text-brand-500 font-semibold': activeView === 'tasks' }"
        >
          <ListTodo class="w-5 h-5" />
          <span class="text-[9px]">งานทั้งหมด</span>
        </button>

        <button
          @click="selectMobileView('settings')"
          class="flex flex-col items-center justify-center space-y-0.5 text-slate-500 hover:text-brand-500 min-h-[44px] min-w-[44px] px-2 py-1 rounded-lg"
          :class="{ 'text-brand-500 font-semibold': activeView === 'settings' }"
        >
          <Settings class="w-5 h-5" />
          <span class="text-[9px]">ตั้งค่า</span>
        </button>
      </nav>
    </div>

    <!-- Toast Container -->
    <ToastNotification />

    <!-- Command Palette / Quick Add Task Modal -->
    <QuickAddModal
      v-if="showQuickAdd"
      :show="showQuickAdd"
      @close="showQuickAdd = false"
    />
  </div>
</template>
