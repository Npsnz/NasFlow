<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace'
import { useUIStore } from '@/stores/ui'
import {
  LayoutDashboard,
  Calendar,
  ListTodo,
  Settings,
  FolderDot,
  Briefcase,
  Home,
  Heart,
  Edit2,
  X,
  Palette
} from 'lucide-vue-next'

const router = useRouter()
const workspaceStore = useWorkspaceStore()
const uiStore = useUIStore()

const workspaces = computed(() => workspaceStore.workspaces)
const currentWorkspace = computed(() => workspaceStore.currentWorkspace)
const isSidebarOpen = computed(() => uiStore.sidebarOpen)
const activeView = computed(() => uiStore.activeView)

const showIconPicker = ref(false)
const editingWorkspaceId = ref<number | null>(null)
const icons = [
  { name: 'Briefcase', component: Briefcase },
  { name: 'Home', component: Home },
  { name: 'Heart', component: Heart },
  { name: 'FolderDot', component: FolderDot },
  { name: 'LayoutDashboard', component: LayoutDashboard },
  { name: 'Calendar', component: Calendar },
  { name: 'ListTodo', component: ListTodo },
  { name: 'Settings', component: Settings }
]

const selectView = (view: 'board' | 'calendar' | 'tasks' | 'settings') => {
  uiStore.setActiveView(view)
  uiStore.setSidebarOpen(false) // Close sidebar on mobile after selecting
  
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

const selectWorkspace = (ws: any) => {
  workspaceStore.setCurrentWorkspace(ws)
  // Ensure we are on the board view when switching workspaces
  selectView('board')
}


const colors = [
  '#9333ea', '#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#ec4899', '#6366f1', '#14b8a6'
]

const showColorPicker = ref(false)

const openIconPicker = (wsId: number) => {
  editingWorkspaceId.value = wsId
  showIconPicker.value = true
}

const openColorPicker = (wsId: number) => {
  editingWorkspaceId.value = wsId
  showColorPicker.value = true
}

const updateWorkspaceIcon = async (iconName: string) => {
  if (editingWorkspaceId.value) {
    try {
      await workspaceStore.updateWorkspace(editingWorkspaceId.value, { icon: iconName })
    } catch (err) {
      console.error('Failed to update icon:', err)
    }
  }
  showIconPicker.value = false
  editingWorkspaceId.value = null
}

const updateWorkspaceColor = async (color: string) => {
  if (editingWorkspaceId.value) {
    try {
      await workspaceStore.updateWorkspace(editingWorkspaceId.value, { color })
    } catch (err) {
      console.error('Failed to update color:', err)
    }
  }
  showColorPicker.value = false
  editingWorkspaceId.value = null
}
</script>

<template>
  <!-- Sidebar Backdrop on Mobile -->
  <div
    v-if="isSidebarOpen"
    class="fixed inset-0 z-40 bg-slate-900/40 backdrop-blur-sm lg:hidden"
    @click="uiStore.setSidebarOpen(false)"
  ></div>

  <!-- Sidebar Container -->
  <aside
    class="fixed inset-y-0 left-0 z-40 w-64 border-r border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900 flex flex-col justify-between transform transition-transform duration-300 lg:translate-x-0 lg:static lg:h-screen"
    :class="isSidebarOpen ? 'translate-x-0' : '-translate-x-full'"
  >
    <!-- Logo & Brand -->
    <div>
      <div class="h-16 border-b border-slate-200 dark:border-slate-700/50 flex items-center px-6 space-x-3">
        <img src="/flame.png" alt="NasFlow" class="w-8 h-8 rounded-lg shadow-md" />
        <div>
          <span class="text-base font-bold bg-gradient-to-r from-slate-900 to-slate-700 dark:from-white dark:to-slate-300 bg-clip-text text-transparent">NasFlow</span>
          <span class="text-[10px] block text-slate-400 font-medium">Task Management</span>
        </div>
      </div>

      <!-- Main Navigation -->
      <div class="px-3 py-4 space-y-1 border-b border-slate-200 dark:border-slate-700/50">
        <p class="px-3 text-[11px] font-semibold text-slate-400 dark:text-slate-500 uppercase tracking-wider mb-2">เมนูหลัก</p>
        
        <button
          @click="selectView('board')"
          class="w-full flex items-center space-x-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors min-h-[44px]"
          :class="activeView === 'board' ? 'nav-active' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800/50'"
        >
          <LayoutDashboard class="w-4 h-4" />
          <span>กระดานงาน (Board)</span>
        </button>

        <button
          @click="selectView('calendar')"
          class="w-full flex items-center space-x-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors min-h-[44px]"
          :class="activeView === 'calendar' ? 'nav-active' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800/50'"
        >
          <Calendar class="w-4 h-4" />
          <span>ปฏิทินงาน (Calendar)</span>
        </button>

        <button
          @click="selectView('tasks')"
          class="w-full flex items-center space-x-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors min-h-[44px]"
          :class="activeView === 'tasks' ? 'nav-active' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800/50'"
        >
          <ListTodo class="w-4 h-4" />
          <span>งานทั้งหมด (List)</span>
        </button>

        <button
          @click="selectView('settings')"
          class="w-full flex items-center space-x-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors min-h-[44px]"
          :class="activeView === 'settings' ? 'nav-active' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800/50'"
        >
          <Settings class="w-4 h-4" />
          <span>ตั้งค่า (Settings)</span>
        </button>
      </div>

      <!-- Workspaces Section (Bottom) -->
      <div class="px-3 py-4 space-y-1 border-slate-200 dark:border-slate-700/50">
        <p class="px-3 text-[11px] font-semibold text-slate-400 dark:text-slate-500 uppercase tracking-wider mb-2">Workspaces</p>
        <div class="space-y-1">
          <div
            v-for="ws in workspaces"
            :key="ws.id"
            class="group relative flex items-center"
          >
            <button
              @click="selectWorkspace(ws)"
              class="w-full flex items-center px-3 py-2.5 rounded-lg text-sm font-medium transition-colors min-h-[40px]"
              :class="currentWorkspace?.id === ws.id ? 'nav-active' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800/50'"
            >
              <div class="w-6 h-6 rounded-lg flex items-center justify-center flex-shrink-0 mr-2.5 text-white text-sm font-semibold" :style="{ backgroundColor: ws.color }">
                <Briefcase v-if="!ws.icon || ws.icon === 'Briefcase'" class="w-4 h-4" />
                <Home v-else-if="ws.icon === 'Home'" class="w-4 h-4" />
                <Heart v-else-if="ws.icon === 'Heart'" class="w-4 h-4" />
                <FolderDot v-else-if="ws.icon === 'FolderDot'" class="w-4 h-4" />
                <LayoutDashboard v-else-if="ws.icon === 'LayoutDashboard'" class="w-4 h-4" />
                <Calendar v-else-if="ws.icon === 'Calendar'" class="w-4 h-4" />
                <ListTodo v-else-if="ws.icon === 'ListTodo'" class="w-4 h-4" />
                <Settings v-else-if="ws.icon === 'Settings'" class="w-4 h-4" />
                <span v-else>{{ ws.name.charAt(0).toUpperCase() }}</span>
              </div>
              <span class="truncate">{{ ws.name }}</span>
            </button>
            <div class="absolute right-2 opacity-0 group-hover:opacity-100 transition-opacity flex items-center gap-1">
              <button
                @click.stop="openIconPicker(ws.id)"
                class="p-1 rounded bg-slate-200 dark:bg-slate-700 hover:bg-slate-300 dark:hover:bg-slate-600"
                title="เปลี่ยน Icon"
              >
                <Edit2 class="w-3.5 h-3.5 text-slate-700 dark:text-slate-300" />
              </button>
              <button
                @click.stop="openColorPicker(ws.id)"
                class="p-1 rounded bg-slate-200 dark:bg-slate-700 hover:bg-slate-300 dark:hover:bg-slate-600"
                title="เปลี่ยนสี"
              >
                <Palette class="w-3.5 h-3.5 text-slate-700 dark:text-slate-300" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- My Favorites Section -->
      <div class="px-3 py-4 space-y-1">
        <p class="px-3 text-[11px] font-semibold text-slate-400 dark:text-slate-500 uppercase tracking-wider mb-2">⭐ My Favorites</p>
        <button
          v-if="workspaces.length === 0"
          class="w-full px-3 py-2.5 text-xs text-slate-400 italic rounded-lg"
        >
          ไม่มี Favorites
        </button>
      </div>

    </div>
  </aside>

  <!-- Icon Picker Modal -->
  <div v-if="showIconPicker" class="fixed inset-0 z-50 bg-slate-900/50 dark:bg-slate-900/70 backdrop-blur-sm flex items-center justify-center">
    <div class="bg-white dark:bg-slate-800 rounded-xl shadow-xl p-6 w-96 max-w-full mx-4">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-bold text-slate-900 dark:text-white">เลือก Icon</h3>
        <button
          @click="showIconPicker = false"
          class="p-1 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-700"
        >
          <X class="w-5 h-5 text-slate-500" />
        </button>
      </div>

      <div class="grid grid-cols-4 gap-3">
        <button
          v-for="icon in icons"
          :key="icon.name"
          @click="updateWorkspaceIcon(icon.name)"
          class="flex items-center justify-center p-3 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-700 transition-colors"
        >
          <component :is="icon.component" class="w-5 h-5 text-slate-600 dark:text-slate-300" />
        </button>
      </div>
    </div>
  </div>

  <!-- Color Picker Modal -->
  <div v-if="showColorPicker" class="fixed inset-0 z-50 bg-slate-900/50 dark:bg-slate-900/70 backdrop-blur-sm flex items-center justify-center">
    <div class="bg-white dark:bg-slate-800 rounded-xl shadow-xl p-6 w-80 max-w-full mx-4">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-bold text-slate-900 dark:text-white">เลือกสี</h3>
        <button
          @click="showColorPicker = false"
          class="p-1 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-700"
        >
          <X class="w-5 h-5 text-slate-500" />
        </button>
      </div>

      <div class="grid grid-cols-4 gap-3">
        <button
          v-for="color in colors"
          :key="color"
          @click="updateWorkspaceColor(color)"
          class="w-12 h-12 rounded-lg border-2 border-slate-300 dark:border-slate-600 hover:border-slate-500 dark:hover:border-slate-400 transition-colors"
          :style="{ backgroundColor: color }"
        ></button>
      </div>
    </div>
  </div>
</template>
