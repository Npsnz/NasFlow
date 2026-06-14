<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { useTasksStore } from '@/stores/tasks'
import { useUIStore } from '@/stores/ui'
import KanbanBoard from '@/components/board/KanbanBoard.vue'
import TaskDrawer from '@/components/task/TaskDrawer.vue'
import {
  CalendarDays,
  AlertOctagon,
  CheckCircle,
  Clock,
  Activity
} from 'lucide-vue-next'

const tasksStore = useTasksStore()
const uiStore = useUIStore()

const stats = computed(() => tasksStore.stats)

let statsPollInterval: any = null

onMounted(() => {
  tasksStore.fetchStats()
  tasksStore.fetchTasks()

  // Poll stats every 5 minutes as fallback
  statsPollInterval = setInterval(() => {
    tasksStore.fetchStats()
  }, 5 * 60 * 1000)

  // Listen to onboarding trigger
  const runOnboarding = localStorage.getItem('tf_onboarding')
  if (runOnboarding === 'true') {
    uiStore.showToast('ยินดีต้อนรับสู่ NasFlow! ลองกด Ctrl+K เพื่อทดลองเพิ่มงานด่วนได้เลย ✨', 'info')
    localStorage.removeItem('tf_onboarding')
  }
})

onUnmounted(() => {
  if (statsPollInterval) clearInterval(statsPollInterval)
})

</script>

<template>
  <div class="flex flex-col h-full overflow-hidden bg-gray-50 dark:bg-slate-900/20">
    <!-- Main Workspace Split Area -->
    <div class="flex-grow flex flex-col lg:flex-row overflow-hidden relative">
      <!-- Left side: Kanban Board -->
      <div class="flex-grow overflow-y-auto">
        <KanbanBoard />
      </div>

      <!-- Right side: Stats Panel -->
      <div class="hidden xl:flex flex-col w-80 border-l border-slate-200 dark:border-slate-700/50 bg-white dark:bg-slate-800/50 overflow-hidden flex-shrink-0 select-none">
        <div class="h-14 border-b border-slate-200 dark:border-slate-700/50 flex items-center px-4 space-x-2.5">
          <Activity class="w-4 h-4 text-brand-500" />
          <h3 class="text-sm font-semibold text-slate-900 dark:text-white">Statistics</h3>
        </div>

        <div class="flex-grow overflow-y-auto p-4 space-y-3">
          <!-- Stat 1: Today's Tasks -->
          <div class="bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-700/40 dark:to-slate-800/40 rounded-xl p-4 border border-slate-200 dark:border-slate-700/50">
            <div class="flex items-start justify-between mb-2">
              <div>
                <p class="text-[10px] font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider">Today</p>
                <p class="text-2xl font-bold text-blue-600 dark:text-blue-400 mt-1">{{ stats.todo + stats.doing }}</p>
              </div>
              <CalendarDays class="w-5 h-5 text-blue-500" />
            </div>
            <div class="h-6 flex items-end gap-0.5 bg-white dark:bg-slate-800/50 rounded p-1">
              <div v-for="i in 7" :key="'today-' + i" class="flex-1 bg-blue-300 dark:bg-blue-600 rounded-sm" :style="{ height: Math.random() * 100 + '%' }"></div>
            </div>
          </div>

          <!-- Stat 2: Overdue -->
          <div class="bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-700/40 dark:to-slate-800/40 rounded-xl p-4 border border-slate-200 dark:border-slate-700/50">
            <div class="flex items-start justify-between mb-2">
              <div>
                <p class="text-[10px] font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider">Overdue</p>
                <p class="text-2xl font-bold text-red-600 dark:text-red-400 mt-1">{{ stats.overdue }}</p>
              </div>
              <AlertOctagon class="w-5 h-5 text-red-500" />
            </div>
            <div class="h-6 flex items-end gap-0.5 bg-white dark:bg-slate-800/50 rounded p-1">
              <div v-for="i in 7" :key="'overdue-' + i" class="flex-1 bg-red-300 dark:bg-red-600 rounded-sm" :style="{ height: Math.random() * 100 + '%' }"></div>
            </div>
          </div>

          <!-- Stat 3: Completed -->
          <div class="bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-700/40 dark:to-slate-800/40 rounded-xl p-4 border border-slate-200 dark:border-slate-700/50">
            <div class="flex items-start justify-between mb-2">
              <div>
                <p class="text-[10px] font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider">Completed</p>
                <p class="text-2xl font-bold text-green-600 dark:text-green-400 mt-1">{{ stats.done_today }}</p>
              </div>
              <CheckCircle class="w-5 h-5 text-green-500" />
            </div>
            <div class="h-6 flex items-end gap-0.5 bg-white dark:bg-slate-800/50 rounded p-1">
              <div v-for="i in 7" :key="'completed-' + i" class="flex-1 bg-green-300 dark:bg-green-600 rounded-sm" :style="{ height: Math.random() * 100 + '%' }"></div>
            </div>
          </div>

          <!-- Stat 4: This Week -->
          <div class="bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-700/40 dark:to-slate-800/40 rounded-xl p-4 border border-slate-200 dark:border-slate-700/50">
            <div class="flex items-start justify-between mb-2">
              <div>
                <p class="text-[10px] font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider">This Week</p>
                <p class="text-2xl font-bold text-brand-500 dark:text-brand-400 mt-1">{{ stats.due_this_week }}</p>
              </div>
              <Clock class="w-5 h-5 text-brand-500" />
            </div>
            <div class="h-6 flex items-end gap-0.5 bg-white dark:bg-slate-800/50 rounded p-1">
              <div v-for="i in 7" :key="'week-' + i" class="flex-1 bg-brand-300 dark:bg-brand-600 rounded-sm" :style="{ height: Math.random() * 100 + '%' }"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Slide-in Right Task Drawer Panel -->
    <TaskDrawer />
  </div>
</template>
