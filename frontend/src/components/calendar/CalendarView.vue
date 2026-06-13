<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useTasksStore, Task } from '@/stores/tasks'
import { useWorkspaceStore } from '@/stores/workspace'
import { useUIStore } from '@/stores/ui'
import { ChevronLeft, ChevronRight, Plus, HelpCircle, User, CalendarDays } from 'lucide-vue-next'

const tasksStore = useTasksStore()
const workspaceStore = useWorkspaceStore()
const uiStore = useUIStore()

const currentWorkspace = computed(() => workspaceStore.currentWorkspace)

// Date states
const today = new Date()
const currentYear = ref(today.getFullYear())
const currentMonth = ref(today.getMonth()) // 0-indexed

// Side panel date selection details
const selectedDayDate = ref<Date | null>(null)

onMounted(() => {
  tasksStore.fetchTasks()
  workspaceStore.fetchWorkspaces()
})

// Thai Month Names
const thaiMonths = [
  'มกราคม', 'กุมภาพันธ์', 'มีนาคม', 'เมษายน', 'พฤษภาคม', 'มิถุนายน',
  'กรกฎาคม', 'สิงหาคม', 'กันยายน', 'ตุลาคม', 'พฤศจิกายน', 'ธันวาคม'
]

// Weekday Headers
const thaiWeekdays = ['จ', 'อ', 'พ', 'พฤ', 'ศ', 'ส', 'อา']

const currentMonthLabel = computed(() => {
  return `${thaiMonths[currentMonth.value]} ${currentYear.value + 543}` // Buddhist era conversion
})

const navigateMonth = (direction: 'prev' | 'next') => {
  if (direction === 'prev') {
    if (currentMonth.value === 0) {
      currentMonth.value = 11
      currentYear.value--
    } else {
      currentMonth.value--
    }
  } else {
    if (currentMonth.value === 11) {
      currentMonth.value = 0
      currentYear.value++
    } else {
      currentMonth.value++
    }
  }
}

const jumpToToday = () => {
  currentYear.value = today.getFullYear()
  currentMonth.value = today.getMonth()
  selectedDayDate.value = null
}

// Generate calendar cells (6 weeks x 7 columns = 42 cells)
const calendarCells = computed(() => {
  const year = currentYear.value
  const month = currentMonth.value

  const firstDay = new Date(year, month, 1)
  // Monday starting: 0=Mon, 6=Sun
  let firstDayIndex = firstDay.getDay() - 1
  if (firstDayIndex === -1) firstDayIndex = 6 // Sunday is index 6

  const totalDays = new Date(year, month + 1, 0).getDate()
  const prevMonthTotalDays = new Date(year, month, 0).getDate()

  const cells = []

  // Trailing days from previous month
  for (let i = firstDayIndex - 1; i >= 0; i--) {
    const d = new Date(year, month - 1, prevMonthTotalDays - i)
    cells.push({
      date: d,
      isCurrentMonth: false,
      dayNum: d.getDate()
    })
  }

  // Current month days
  for (let i = 1; i <= totalDays; i++) {
    const d = new Date(year, month, i)
    cells.push({
      date: d,
      isCurrentMonth: true,
      dayNum: i
    })
  }

  // Leading days from next month
  const remainingCells = 42 - cells.length
  for (let i = 1; i <= remainingCells; i++) {
    const d = new Date(year, month + 1, i)
    cells.push({
      date: d,
      isCurrentMonth: false,
      dayNum: i
    })
  }

  return cells
})

// Match tasks to date cells
const getTasksForDate = (date: Date) => {
  return tasksStore.tasks.filter((t) => {
    if (!t.due_date || t.status === 'archived') return false
    
    // Scoped to current workspace
    if (currentWorkspace.value && t.workspace_id !== currentWorkspace.value.id) return false
    if (t.parent_task_id !== null) return false // No subtasks on calendar

    const taskDate = new Date(t.due_date)
    return taskDate.toDateString() === date.toDateString()
  })
}

// Tasks with no due date
const tasksNoDueDate = computed(() => {
  return tasksStore.tasks.filter((t) => {
    if (t.due_date || t.status === 'archived') return false
    if (currentWorkspace.value && t.workspace_id !== currentWorkspace.value.id) return false
    if (t.parent_task_id !== null) return false
    return true
  })
})

// Verify if date cell has overdue tasks
const isOverdueDay = (date: Date, cellTasks: Task[]) => {
  const now = new Date()
  // Strip time from both for comparison
  const compareDate = new Date(date.getFullYear(), date.getMonth(), date.getDate())
  const compareNow = new Date(now.getFullYear(), now.getMonth(), now.getDate())

  if (compareDate >= compareNow) return false
  return cellTasks.some(t => t.status !== 'done')
}

// Side panel details
const selectedDayTasks = computed(() => {
  if (!selectedDayDate.value) return []
  return getTasksForDate(selectedDayDate.value)
})

const selectDay = (date: Date) => {
  selectedDayDate.value = date
}

// Add task with pre-filled date
const quickAddTaskForDate = (date: Date) => {
  // Sets prefill and triggers event dispatcher
  window.dispatchEvent(new CustomEvent('open-quick-add'))
  // Wait a tick for modal input, then pre-fill date text if input is bound,
  // or we can prefill it in state. In our case, the quick-add command palette
  // parses text, so we can set rawText in palette if needed, or simply pass
  // text preset.
  setTimeout(() => {
    const formatted = `@${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`
    // Query task form input element and insert text
    const input = document.querySelector('input[placeholder*="พิมพ์ชื่องานด่วน"]') as HTMLInputElement
    if (input) {
      input.value = ` ${formatted}`
      input.dispatchEvent(new Event('input'))
    }
  }, 100)
}

const openTaskDrawer = (id: number) => {
  window.dispatchEvent(new CustomEvent('open-task-details', { detail: id }))
}

const getWorkspaceColor = (wsId: number) => {
  const ws = workspaceStore.workspaces.find(w => w.id === wsId)
  return ws ? ws.color : '#cbd5e1'
}
</script>

<template>
  <div class="p-6 flex flex-col lg:flex-row gap-6 select-none">
    <!-- Main Calendar Grid Area -->
    <div class="flex-grow bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-4 shadow-sm flex flex-col">
      <!-- Calendar Controls -->
      <div class="flex items-center justify-between mb-4 pb-2 border-b border-border-light dark:border-border-dark">
        <div class="flex items-center space-x-3">
          <h3 class="text-base font-bold text-slate-800 dark:text-white">{{ currentMonthLabel }}</h3>
        </div>
        
        <div class="flex items-center space-x-2">
          <!-- Back to Today -->
          <button
            @click="jumpToToday"
            class="px-3.5 py-1.5 bg-slate-100 hover:bg-slate-200 text-slate-700 dark:bg-slate-800 dark:text-slate-300 rounded-lg text-xs font-semibold min-h-[44px]"
          >
            วันนี้
          </button>
          
          <!-- Month Switcher arrows -->
          <div class="flex border border-border-light dark:border-border-dark rounded-lg overflow-hidden">
            <button
              @click="navigateMonth('prev')"
              class="p-2 hover:bg-slate-50 dark:hover:bg-slate-800 text-slate-500 min-h-[44px] min-w-[44px]"
            >
              <ChevronLeft class="w-4 h-4" />
            </button>
            <button
              @click="navigateMonth('next')"
              class="p-2 hover:bg-slate-50 dark:hover:bg-slate-800 text-slate-500 border-l border-border-light dark:border-border-dark min-h-[44px] min-w-[44px]"
            >
              <ChevronRight class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- Month Grid Table -->
      <div class="flex-grow grid grid-cols-7 border-t border-l border-border-light dark:border-border-dark text-center">
        <!-- Weekday Headers -->
        <div
          v-for="day in thaiWeekdays"
          :key="day"
          class="py-2.5 bg-slate-50 dark:bg-slate-900/50 border-r border-b border-border-light dark:border-border-dark font-bold text-[10px] text-slate-400 dark:text-slate-500 uppercase tracking-wider"
        >
          {{ day }}
        </div>

        <!-- Date Cells -->
        <div
          v-for="cell in calendarCells"
          :key="cell.date.toISOString()"
          @click="selectDay(cell.date)"
          class="min-h-[90px] p-1.5 border-r border-b border-border-light dark:border-border-dark flex flex-col justify-between cursor-pointer group hover:bg-slate-50/50 dark:hover:bg-slate-800/10 relative transition-colors"
          :class="{
            'bg-slate-50/30 text-slate-400 dark:bg-slate-950/20 dark:text-slate-600': !cell.isCurrentMonth,
            'bg-red-50/50 dark:bg-red-950/10': isOverdueDay(cell.date, getTasksForDate(cell.date)),
            'ring-2 ring-brand-500/50 z-10': cell.date.toDateString() === today.toDateString(),
            'bg-slate-100 dark:bg-slate-800': selectedDayDate?.toDateString() === cell.date.toDateString()
          }"
        >
          <!-- Date label row -->
          <div class="flex items-center justify-between mb-1">
            <span
              class="text-xs font-semibold px-1 py-0.2 rounded"
              :class="{
                'bg-brand-500 text-white': cell.date.toDateString() === today.toDateString()
              }"
            >
              {{ cell.dayNum }}
            </span>

            <!-- Quick Add hover trigger -->
            <button
              @click.stop="quickAddTaskForDate(cell.date)"
              class="opacity-0 group-hover:opacity-100 text-brand-500 hover:text-brand-600 p-0.5 rounded transition-opacity min-h-[32px] min-w-[32px] flex items-center justify-center"
              title="เพิ่มงานในวันนี้"
            >
              <Plus class="w-3.5 h-3.5" />
            </button>
          </div>

          <!-- Tasks list container in cell -->
          <div class="flex-grow flex flex-col space-y-1 overflow-y-auto max-h-[60px] scrollbar-none">
            <div
              v-for="task in getTasksForDate(cell.date).slice(0, 3)"
              :key="task.id"
              @click.stop="openTaskDrawer(task.id)"
              class="text-[9px] font-bold px-1.5 py-0.5 rounded truncate select-none border"
              :style="{
                backgroundColor: getWorkspaceColor(task.workspace_id) + '15',
                color: getWorkspaceColor(task.workspace_id),
                borderColor: getWorkspaceColor(task.workspace_id) + '25'
              }"
              :title="task.title"
            >
              {{ task.title }}
            </div>
            <span
              v-if="getTasksForDate(cell.date).length > 3"
              class="text-[8px] font-semibold text-slate-400 text-left pl-1"
            >
              + อีก {{ getTasksForDate(cell.date).length - 3 }} งาน
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Side details Panel & No Due Date section -->
    <div class="w-full lg:w-72 flex flex-col gap-6">
      <!-- Selected Day details panel -->
      <div
        v-if="selectedDayDate"
        class="bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-4 shadow-sm space-y-3"
      >
        <div class="flex items-center justify-between pb-2 border-b border-border-light dark:border-border-dark">
          <h4 class="text-xs font-bold text-slate-800 dark:text-white uppercase tracking-wide">
            งานวันที่ {{ selectedDayDate.toLocaleDateString('th-TH', { day: 'numeric', month: 'short' }) }}
          </h4>
          <button
            @click="quickAddTaskForDate(selectedDayDate)"
            class="text-xs text-brand-500 hover:text-brand-600 flex items-center space-x-0.5 font-semibold min-h-[36px]"
          >
            <Plus class="w-3.5 h-3.5" />
            <span>เพิ่มงาน</span>
          </button>
        </div>

        <div class="space-y-2 max-h-60 overflow-y-auto pr-1">
          <div
            v-for="task in selectedDayTasks"
            :key="task.id"
            @click="openTaskDrawer(task.id)"
            class="p-2.5 bg-slate-50 dark:bg-slate-900/40 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 cursor-pointer flex items-center justify-between text-xs"
          >
            <span class="font-medium truncate mr-2">{{ task.title }}</span>
            <span
              class="text-[9px] px-1.5 py-0.5 rounded font-bold"
              :class="{
                'bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-400': task.status === 'todo',
                'bg-blue-50 text-blue-600 dark:bg-blue-950/20 dark:text-blue-400': task.status === 'doing',
                'bg-green-50 text-green-600 dark:bg-green-950/20 dark:text-green-400': task.status === 'done'
              }"
            >
              {{ task.status === 'todo' ? 'รอทำ' : task.status === 'doing' ? 'กำลังทำ' : 'เสร็จแล้ว' }}
            </span>
          </div>
          <p v-if="selectedDayTasks.length === 0" class="text-xs text-slate-400 italic text-center py-6">
            ไม่มีกำหนดส่งงานในวันนี้
          </p>
        </div>
      </div>

      <!-- No Due Date section -->
      <div class="bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-4 shadow-sm space-y-3">
        <h4 class="text-xs font-bold text-slate-800 dark:text-white uppercase tracking-wider flex items-center space-x-1.5 pb-2 border-b border-border-light dark:border-border-dark">
          <CalendarDays class="w-4 h-4 text-slate-400" />
          <span>ไม่มีกำหนดการ ({{ tasksNoDueDate.length }})</span>
        </h4>

        <div class="space-y-2 max-h-72 overflow-y-auto pr-1">
          <div
            v-for="task in tasksNoDueDate"
            :key="task.id"
            @click="openTaskDrawer(task.id)"
            class="p-2.5 bg-slate-50 dark:bg-slate-900/40 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 cursor-pointer flex items-center justify-between text-xs"
          >
            <span class="font-medium truncate mr-2">{{ task.title }}</span>
            <span
              class="text-[9px] px-1.5 py-0.5 rounded font-bold"
              :class="{
                'bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-400': task.status === 'todo',
                'bg-blue-50 text-blue-600 dark:bg-blue-950/20 dark:text-blue-400': task.status === 'doing',
                'bg-green-50 text-green-600 dark:bg-green-950/20 dark:text-green-400': task.status === 'done'
              }"
            >
              {{ task.status === 'todo' ? 'รอทำ' : task.status === 'doing' ? 'กำลังทำ' : 'เสร็จแล้ว' }}
            </span>
          </div>
          <p v-if="tasksNoDueDate.length === 0" class="text-xs text-slate-400 italic text-center py-6">
            งานทั้งหมดมีกำหนดส่งครบถ้วน
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
