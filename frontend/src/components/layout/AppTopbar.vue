<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace'
import { useTasksStore } from '@/stores/tasks'
import { useUIStore } from '@/stores/ui'
import { useAuthStore } from '@/stores/auth'
import { Menu, Search, X, Filter, Sparkles, LogOut, Settings, User, HelpCircle } from 'lucide-vue-next'

const router = useRouter()
const workspaceStore = useWorkspaceStore()
const tasksStore = useTasksStore()
const uiStore = useUIStore()
const authStore = useAuthStore()

const showUserMenu = ref(false)
const showShortcuts = ref(false)
const user = computed(() => authStore.user)

const currentWorkspace = computed(() => workspaceStore.currentWorkspace)
const activeView = computed(() => uiStore.activeView)
const tags = computed(() => tasksStore.tags)

const showTagDropdown = ref(false)
const searchInput = ref<HTMLInputElement | null>(null)

// Focus search input on Cmd+F / Ctrl+F
const handleGlobalKeydown = (e: KeyboardEvent) => {
  if ((e.metaKey || e.ctrlKey) && e.key === 'f') {
    e.preventDefault()
    searchInput.value?.focus()
  }
  // Open shortcuts on ?
  if (e.key === '?' && !e.ctrlKey && !e.metaKey) {
    e.preventDefault()
    showShortcuts.value = true
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleGlobalKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleGlobalKeydown)
})

const pageTitle = computed(() => {
  switch (activeView.value) {
    case 'calendar':
      return 'ปฏิทินงาน'
    case 'tasks':
      return 'งานทั้งหมด'
    case 'settings':
      return 'ตั้งค่าระบบ'
    case 'board':
    default:
      return currentWorkspace.value?.name || 'แดชบอร์ด'
  }
})

// Debounce search update
let debounceTimeout: any = null
const searchVal = ref(tasksStore.filters.search)

watch(searchVal, (newVal) => {
  if (debounceTimeout) clearTimeout(debounceTimeout)
  debounceTimeout = setTimeout(() => {
    tasksStore.filters.search = newVal
  }, 300)
})

// Sync searchVal back if global filter is reset
watch(() => tasksStore.filters.search, (newVal) => {
  searchVal.value = newVal
})

const hasActiveFilters = computed(() => {
  return (
    tasksStore.filters.search !== '' ||
    tasksStore.filters.priority !== '' ||
    tasksStore.filters.tagIds.length > 0 ||
    tasksStore.filters.dueRange !== ''
  )
})

const clearAllFilters = () => {
  searchVal.value = ''
  tasksStore.filters.search = ''
  tasksStore.filters.priority = ''
  tasksStore.filters.tagIds = []
  tasksStore.filters.dueRange = ''
}

const toggleTagSelection = (tagId: number) => {
  const index = tasksStore.filters.tagIds.indexOf(tagId)
  if (index === -1) {
    tasksStore.filters.tagIds.push(tagId)
  } else {
    tasksStore.filters.tagIds.splice(index, 1)
  }
}

const openQuickAdd = () => {
  if (typeof window !== 'undefined') {
    window.dispatchEvent(new CustomEvent('open-quick-add'))
  }
}

const handleLogout = async () => {
  if (confirm('คุณต้องการออกจากระบบหรือไม่?')) {
    await authStore.logout()
    router.push('/login')
  }
}
</script>

<template>
  <header class="h-16 border-b border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 flex items-center justify-between px-6 z-30">
    <div class="flex items-center space-x-4">
      <!-- Hamburger Menu on Mobile -->
      <button
        @click="uiStore.toggleSidebar"
        class="lg:hidden p-2 rounded-lg text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800 min-h-[44px]"
      >
        <Menu class="w-5 h-5" />
      </button>

      <!-- Page Title -->
      <div class="flex items-center space-x-2.5">
        <span
          v-if="activeView === 'board' && currentWorkspace"
          class="w-3.5 h-3.5 rounded-full"
          :style="{ backgroundColor: currentWorkspace.color }"
        ></span>
        <h2 class="text-lg font-bold text-slate-800 dark:text-white">{{ pageTitle }}</h2>
      </div>
    </div>

    <!-- Filters Section (Hidden on Settings) -->
    <div v-if="activeView !== 'settings'" class="hidden md:flex items-center space-x-3 flex-grow max-w-2xl justify-end">
      <!-- Search Input -->
      <div class="relative w-48 lg:w-64">
        <span class="absolute inset-y-0 left-3 flex items-center text-slate-400">
          <Search class="w-4 h-4" />
        </span>
        <input
          ref="searchInput"
          type="text"
          v-model="searchVal"
          placeholder="ค้นหาชื่องาน... (Ctrl+F)"
          class="w-full pl-9 pr-3 py-1.5 text-xs bg-slate-50 border border-border-light rounded-lg focus:outline-none focus:ring-1 focus:ring-brand-500 focus:bg-white dark:bg-slate-900/60 dark:border-border-dark dark:focus:bg-slate-900 min-h-[44px]"
        />
        <button
          v-if="searchVal"
          @click="searchVal = ''"
          class="absolute inset-y-0 right-3 flex items-center text-slate-400 hover:text-slate-600"
        >
          <X class="w-3.5 h-3.5" />
        </button>
      </div>

      <!-- Priority Filter -->
      <select
        v-model="tasksStore.filters.priority"
        class="text-xs bg-slate-50 border border-slate-200 rounded-lg px-2.5 py-1.5 focus:outline-none focus:ring-1 focus:ring-brand-500 dark:bg-slate-800/50 dark:border-slate-700 min-h-[44px]"
      >
        <option value="">ทุกความสำคัญ</option>
        <option value="low">ต่ำ (Low)</option>
        <option value="medium">ปานกลาง (Medium)</option>
        <option value="high">สูง (High)</option>
        <option value="urgent">ด่วนที่สุด (Urgent)</option>
      </select>

      <!-- Due Date Preset Filter -->
      <select
        v-model="tasksStore.filters.dueRange"
        class="text-xs bg-slate-50 border border-slate-200 rounded-lg px-2.5 py-1.5 focus:outline-none focus:ring-1 focus:ring-brand-500 dark:bg-slate-800/50 dark:border-slate-700 min-h-[44px]"
      >
        <option value="">ทุกกำหนดส่ง</option>
        <option value="today">วันนี้</option>
        <option value="tomorrow">พรุ่งนี้</option>
        <option value="week">สัปดาห์นี้</option>
        <option value="overdue">เกินกำหนด</option>
      </select>

      <!-- Multi-select Tags Dropdown -->
      <div class="relative">
        <button
          @click="showTagDropdown = !showTagDropdown"
          class="flex items-center space-x-1.5 text-xs bg-slate-50 border border-border-light rounded-lg px-2.5 py-1.5 dark:bg-slate-900/60 dark:border-border-dark hover:bg-slate-100 min-h-[44px]"
        >
          <Filter class="w-3.5 h-3.5" />
          <span>แท็ก ({{ tasksStore.filters.tagIds.length }})</span>
        </button>
        <div
          v-if="showTagDropdown"
          class="absolute right-0 mt-2 w-48 bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-lg shadow-xl py-2 z-50"
        >
          <p class="px-3 py-1 text-[10px] font-semibold text-slate-400 dark:text-slate-500 uppercase tracking-wider">เลือกแท็ก</p>
          <div class="max-h-48 overflow-y-auto px-1">
            <button
              v-for="tag in tags"
              :key="tag.id"
              @click="toggleTagSelection(tag.id)"
              class="w-full flex items-center space-x-2.5 px-2.5 py-2 hover:bg-slate-50 dark:hover:bg-slate-800 rounded text-left text-xs"
            >
              <input
                type="checkbox"
                :checked="tasksStore.filters.tagIds.includes(tag.id)"
                class="rounded text-brand-500 focus:ring-brand-500"
                @click.stop="toggleTagSelection(tag.id)"
              />
              <span class="w-2.5 h-2.5 rounded-full" :style="{ backgroundColor: tag.color }"></span>
              <span class="text-slate-700 dark:text-slate-300 truncate">{{ tag.name }}</span>
            </button>
            <p v-if="tags.length === 0" class="px-3 py-2 text-slate-400 text-xs italic">ไม่มีแท็ก</p>
          </div>
        </div>
      </div>

      <!-- Clear Filters -->
      <button
        v-if="hasActiveFilters"
        @click="clearAllFilters"
        class="text-xs text-red-500 hover:text-red-600 font-semibold px-2 py-1 transition-colors min-h-[44px]"
      >
        ล้างตัวกรอง
      </button>
    </div>

    <!-- Quick Add + Help + User Profile -->
    <div class="flex items-center space-x-3 ml-auto">
      <button
        @click="openQuickAdd"
        class="flex items-center space-x-1 px-3 py-1.5 bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 border border-slate-200 dark:border-slate-700 rounded-lg text-xs hover:bg-slate-200 dark:hover:bg-slate-700 transition-colors font-medium min-h-[44px]"
      >
        <Sparkles class="w-3.5 h-3.5" />
        <span class="hidden sm:inline">ด่วน (Ctrl+K)</span>
      </button>

      <button
        @click="showShortcuts = true"
        class="p-2 rounded-lg text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors min-h-[44px]"
        title="Keyboard Shortcuts (?)"
      >
        <HelpCircle class="w-5 h-5" />
      </button>

      <!-- User Profile Menu -->
      <div class="relative">
        <button
          @click="showUserMenu = !showUserMenu"
          class="flex items-center space-x-2 px-2 py-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800/50 transition-colors min-h-[44px]"
        >
          <img
            v-if="user?.avatar_url"
            :src="user.avatar_url"
            :alt="user.name"
            class="w-8 h-8 rounded-full object-cover border border-slate-200 dark:border-slate-700"
          />
          <div v-else class="w-8 h-8 rounded-full bg-brand-100 dark:bg-brand-500/20 text-brand-500 dark:text-brand-300 flex items-center justify-center font-semibold text-sm">
            {{ user?.name?.charAt(0) || 'U' }}
          </div>
          <div class="hidden md:block text-left">
            <p class="text-xs font-semibold text-slate-900 dark:text-white">{{ user?.name }}</p>
            <p class="text-[10px] text-slate-500 dark:text-slate-400">{{ user?.email }}</p>
          </div>
        </button>

        <!-- Dropdown Menu -->
        <div
          v-if="showUserMenu"
          class="absolute right-0 mt-2 w-56 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg shadow-xl py-1 z-50"
          @click.stop="showUserMenu = false"
        >
          <router-link
            to="/settings"
            class="flex items-center space-x-2 px-4 py-2.5 text-sm text-slate-700 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-700/50"
          >
            <Settings class="w-4 h-4" />
            <span>ตั้งค่าระบบ</span>
          </router-link>
          <button
            @click="handleLogout"
            class="w-full flex items-center space-x-2 px-4 py-2.5 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-950/20"
          >
            <LogOut class="w-4 h-4" />
            <span>ออกจากระบบ</span>
          </button>
        </div>
      </div>
    </div>
  </header>

  <!-- Keyboard Shortcuts Modal -->
  <div v-if="showShortcuts" class="fixed inset-0 z-50 bg-slate-900/50 dark:bg-slate-900/70 backdrop-blur-sm flex items-center justify-center">
    <div class="bg-white dark:bg-slate-800 rounded-xl shadow-xl p-6 w-full max-w-2xl mx-4 max-h-[80vh] overflow-y-auto">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Keyboard Shortcuts</h2>
        <button
          @click="showShortcuts = false"
          class="p-1 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-700"
        >
          <X class="w-6 h-6 text-slate-500" />
        </button>
      </div>

      <div class="space-y-6">
        <!-- Task Management -->
        <div>
          <h3 class="text-sm font-semibold text-slate-900 dark:text-white mb-3 text-brand-500">Task Management</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <div class="flex items-center justify-between p-2.5 bg-slate-50 dark:bg-slate-700/30 rounded-lg">
              <span class="text-sm text-slate-700 dark:text-slate-300">Edit task</span>
              <kbd class="px-2 py-1 bg-slate-200 dark:bg-slate-600 text-slate-900 dark:text-white text-xs font-semibold rounded">E</kbd>
            </div>
            <div class="flex items-center justify-between p-2.5 bg-slate-50 dark:bg-slate-700/30 rounded-lg">
              <span class="text-sm text-slate-700 dark:text-slate-300">Toggle Done</span>
              <kbd class="px-2 py-1 bg-slate-200 dark:bg-slate-600 text-slate-900 dark:text-white text-xs font-semibold rounded">D</kbd>
            </div>
            <div class="flex items-center justify-between p-2.5 bg-slate-50 dark:bg-slate-700/30 rounded-lg">
              <span class="text-sm text-slate-700 dark:text-slate-300">Priority: Low</span>
              <kbd class="px-2 py-1 bg-slate-200 dark:bg-slate-600 text-slate-900 dark:text-white text-xs font-semibold rounded">1</kbd>
            </div>
            <div class="flex items-center justify-between p-2.5 bg-slate-50 dark:bg-slate-700/30 rounded-lg">
              <span class="text-sm text-slate-700 dark:text-slate-300">Priority: Medium</span>
              <kbd class="px-2 py-1 bg-slate-200 dark:bg-slate-600 text-slate-900 dark:text-white text-xs font-semibold rounded">2</kbd>
            </div>
            <div class="flex items-center justify-between p-2.5 bg-slate-50 dark:bg-slate-700/30 rounded-lg">
              <span class="text-sm text-slate-700 dark:text-slate-300">Priority: High</span>
              <kbd class="px-2 py-1 bg-slate-200 dark:bg-slate-600 text-slate-900 dark:text-white text-xs font-semibold rounded">3</kbd>
            </div>
          </div>
        </div>

        <!-- General -->
        <div>
          <h3 class="text-sm font-semibold text-slate-900 dark:text-white mb-3 text-brand-500">General</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <div class="flex items-center justify-between p-2.5 bg-slate-50 dark:bg-slate-700/30 rounded-lg">
              <span class="text-sm text-slate-700 dark:text-slate-300">Quick Add Task</span>
              <div class="flex gap-1">
                <kbd class="px-2 py-1 bg-slate-200 dark:bg-slate-600 text-slate-900 dark:text-white text-xs font-semibold rounded">Ctrl</kbd>
                <kbd class="px-2 py-1 bg-slate-200 dark:bg-slate-600 text-slate-900 dark:text-white text-xs font-semibold rounded">K</kbd>
              </div>
            </div>
            <div class="flex items-center justify-between p-2.5 bg-slate-50 dark:bg-slate-700/30 rounded-lg">
              <span class="text-sm text-slate-700 dark:text-slate-300">Search Tasks</span>
              <div class="flex gap-1">
                <kbd class="px-2 py-1 bg-slate-200 dark:bg-slate-600 text-slate-900 dark:text-white text-xs font-semibold rounded">Ctrl</kbd>
                <kbd class="px-2 py-1 bg-slate-200 dark:bg-slate-600 text-slate-900 dark:text-white text-xs font-semibold rounded">F</kbd>
              </div>
            </div>
            <div class="flex items-center justify-between p-2.5 bg-slate-50 dark:bg-slate-700/30 rounded-lg">
              <span class="text-sm text-slate-700 dark:text-slate-300">Help</span>
              <kbd class="px-2 py-1 bg-slate-200 dark:bg-slate-600 text-slate-900 dark:text-white text-xs font-semibold rounded">?</kbd>
            </div>
          </div>
        </div>

        <div class="pt-4 border-t border-slate-200 dark:border-slate-700">
          <p class="text-xs text-slate-500 dark:text-slate-400">💡 Click on a task to focus it first, then use keyboard shortcuts to manage it</p>
        </div>
      </div>

      <div class="mt-6 flex justify-end">
        <button
          @click="showShortcuts = false"
          class="px-4 py-2 bg-brand-500 text-white rounded-lg hover:bg-brand-600 transition-colors text-sm font-medium"
        >
          Close
        </button>
      </div>
    </div>
  </div>
</template>
