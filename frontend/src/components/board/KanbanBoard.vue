<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useTasksStore, Task } from '@/stores/tasks'
import { useWorkspaceStore } from '@/stores/workspace'
import { useUIStore } from '@/stores/ui'
import KanbanColumn from './KanbanColumn.vue'
import TaskCard from './TaskCard.vue'

const tasksStore = useTasksStore()
const workspaceStore = useWorkspaceStore()
const uiStore = useUIStore()

const handleKey = async (e: KeyboardEvent) => {
  const target = e.target as HTMLElement
  if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) return

  const id = uiStore.focusedTaskId
  if (!id) return

  const task = tasksStore.tasks.find(t => t.id === id)
  if (!task) return

  if (e.key === 'e') {
    e.preventDefault()
    window.dispatchEvent(new CustomEvent('open-task-details', { detail: id }))
  } else if (e.key === 'd') {
    e.preventDefault()
    try {
      if (task.status !== 'done') {
        await tasksStore.completeTask(id)
        uiStore.showToast('ทำเครื่องหมายเสร็จแล้ว', 'success')
      } else {
        await tasksStore.updateTaskStatus(id, 'todo')
        uiStore.showToast('ย้ายกลับรอดำเนินการ', 'info')
      }
    } catch {
      uiStore.showToast('ดำเนินการไม่สำเร็จ', 'error')
    }
  } else if (e.key === '1') {
    e.preventDefault()
    try {
      await tasksStore.updateTask(id, { priority: 'low' })
      uiStore.showToast('ตั้งความสำคัญ: ต่ำ', 'info')
    } catch { /* silent */ }
  } else if (e.key === '2') {
    e.preventDefault()
    try {
      await tasksStore.updateTask(id, { priority: 'medium' })
      uiStore.showToast('ตั้งความสำคัญ: ปานกลาง', 'info')
    } catch { /* silent */ }
  } else if (e.key === '3') {
    e.preventDefault()
    try {
      await tasksStore.updateTask(id, { priority: 'high' })
      uiStore.showToast('ตั้งความสำคัญ: สูง', 'info')
    } catch { /* silent */ }
  }
}

const windowWidth = ref(typeof window !== 'undefined' ? window.innerWidth : 768)

onMounted(() => {
  window.addEventListener('keydown', handleKey)
  const handleResize = () => { windowWidth.value = window.innerWidth }
  window.addEventListener('resize', handleResize)
  return () => {
    window.removeEventListener('keydown', handleKey)
    window.removeEventListener('resize', handleResize)
  }
})
onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKey)
})

const currentWorkspace = computed(() => workspaceStore.currentWorkspace)
const activeMobileColumn = ref<'todo' | 'doing' | 'done'>('todo')

// Dragged task cache
const draggedTask = ref<Task | null>(null)

// Filtering logic
const filteredTasks = computed(() => {
  if (!currentWorkspace.value) return []

  return tasksStore.tasks.filter((t) => {
    // 1. Must belong to the active workspace
    if (t.workspace_id !== currentWorkspace.value?.id) return false
    
    // 2. Filter out subtasks from the main board
    if (t.parent_task_id !== null) return false

    // 3. Title Search
    if (tasksStore.filters.search) {
      const q = tasksStore.filters.search.toLowerCase()
      if (!t.title.toLowerCase().includes(q)) return false
    }

    // 4. Priority Filter
    if (tasksStore.filters.priority) {
      if (t.priority !== tasksStore.filters.priority) return false
    }

    // 5. Multi-select Tags
    if (tasksStore.filters.tagIds.length > 0) {
      const hasMatchingTag = t.tags.some(tag => tasksStore.filters.tagIds.includes(tag.id))
      if (!hasMatchingTag) return false
    }

    // 6. Due Date Presets
    if (tasksStore.filters.dueRange) {
      if (!t.due_date) return false
      const d = new Date(t.due_date)
      const now = new Date()
      
      switch (tasksStore.filters.dueRange) {
        case 'today':
          return d.toDateString() === now.toDateString()
        case 'tomorrow': {
          const tom = new Date()
          tom.setDate(now.getDate() + 1)
          return d.toDateString() === tom.toDateString()
        }
        case 'week': {
          // Check if within next 7 days
          const diff = d.getTime() - now.getTime()
          const diffDays = diff / (1000 * 3600 * 24)
          return diffDays >= -1 && diffDays <= 7
        }
        case 'overdue':
          return d < now && t.status !== 'done'
        default:
          return true
      }
    }

    return true
  })
})

// Columns splitting
const todoTasks = computed(() => filteredTasks.value.filter(t => t.status === 'todo'))
const doingTasks = computed(() => filteredTasks.value.filter(t => t.status === 'doing'))
const doneTasks = computed(() => filteredTasks.value.filter(t => t.status === 'done'))

// HTML5 Drag and Drop Handlers
const onDragStart = (e: DragEvent, task: Task) => {
  draggedTask.value = task
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/plain', String(task.id))
  }
}

const onDragOver = (e: DragEvent) => {
  e.preventDefault()
}

const onDrop = async (e: DragEvent, targetStatus: 'todo' | 'doing' | 'done') => {
  e.preventDefault()
  if (!draggedTask.value) return

  const task = draggedTask.value
  const oldStatus = task.status

  if (oldStatus === targetStatus) {
    draggedTask.value = null
    return
  }

  // Calculate new SortOrder LexoRank style
  let targetTasks = [] as Task[]
  if (targetStatus === 'todo') targetTasks = todoTasks.value
  else if (targetStatus === 'doing') targetTasks = doingTasks.value
  else targetTasks = doneTasks.value

  let newSortOrder = 1000.0
  if (targetTasks.length > 0) {
    // Put it at the bottom of the column
    const lastTask = targetTasks[targetTasks.length - 1]
    newSortOrder = lastTask.sort_order + 1000.0
  }

  // Clear dragged task cache
  draggedTask.value = null

  // 1. Optimistic Update in UI
  try {
    await tasksStore.updateTaskStatus(task.id, targetStatus)
    // 2. Save order
    await tasksStore.reorderTasks([{ id: task.id, sort_order: newSortOrder, status: targetStatus }])
  } catch (err) {
    uiStore.showToast('ไม่สามารถย้ายงานได้', 'error')
  }
}
</script>

<template>
  <div class="p-6 space-y-6 bg-slate-50 dark:bg-slate-900/20 min-h-full">
    <!-- Mobile Column Swapper Tabs -->
    <div class="flex md:hidden border-b border-border-light dark:border-border-dark pb-2 space-x-2">
      <button
        v-for="col in [
          { status: 'todo', label: 'รอทำ', count: todoTasks.length, dot: 'bg-slate-400' },
          { status: 'doing', label: 'กำลังทำ', count: doingTasks.length, dot: 'bg-blue-500' },
          { status: 'done', label: 'เสร็จแล้ว', count: doneTasks.length, dot: 'bg-green-500' }
        ]"
        :key="col.status"
        @click="activeMobileColumn = col.status as any"
        class="flex-1 py-2 rounded-lg text-xs font-semibold flex items-center justify-center space-x-1.5 transition-colors border min-h-[44px]"
        :class="activeMobileColumn === col.status
          ? 'bg-slate-900 text-white dark:bg-slate-800 dark:border-slate-700'
          : 'bg-white text-slate-600 dark:bg-surface-dark dark:text-slate-400 dark:border-border-dark'"
      >
        <span class="w-2 h-2 rounded-full" :class="col.dot"></span>
        <span>{{ col.label }}</span>
        <span class="text-[9px] text-slate-400 bg-slate-100 dark:bg-slate-900 px-1 py-0.2 rounded-full">{{ col.count }}</span>
      </button>
    </div>

    <!-- Desktop Kanban Layout (Swipe/scroll-x on small viewports) -->
    <div class="flex space-x-4 overflow-x-auto pb-4 items-start select-none">
      <!-- To Do Column -->
      <KanbanColumn
        v-show="uiStore.sidebarOpen || activeMobileColumn === 'todo' || windowWidth >= 768"
        status="todo"
        title="รอดำเนินการ"
        dotColor="bg-slate-400"
        :tasks="todoTasks"
        @dragover="onDragOver"
        @drop="onDrop($event, 'todo')"
        class="flex-shrink-0"
      >
        <div
          v-for="task in todoTasks"
          :key="task.id"
          draggable="true"
          @dragstart="onDragStart($event, task)"
          class="transform active:scale-[0.98] transition-transform duration-100"
        >
          <TaskCard :task="task" :isFocused="uiStore.focusedTaskId === task.id" />
        </div>
      </KanbanColumn>

      <!-- Doing Column -->
      <KanbanColumn
        v-show="uiStore.sidebarOpen || activeMobileColumn === 'doing' || windowWidth >= 768"
        status="doing"
        title="กำลังดำเนินการ"
        dotColor="bg-blue-500"
        :tasks="doingTasks"
        @dragover="onDragOver"
        @drop="onDrop($event, 'doing')"
        class="flex-shrink-0"
      >
        <div
          v-for="task in doingTasks"
          :key="task.id"
          draggable="true"
          @dragstart="onDragStart($event, task)"
          class="transform active:scale-[0.98] transition-transform duration-100"
        >
          <TaskCard :task="task" :isFocused="uiStore.focusedTaskId === task.id" />
        </div>
      </KanbanColumn>

      <!-- Done Column -->
      <KanbanColumn
        v-show="uiStore.sidebarOpen || activeMobileColumn === 'done' || windowWidth >= 768"
        status="done"
        title="เสร็จแล้ว"
        dotColor="bg-green-500"
        :tasks="doneTasks"
        @dragover="onDragOver"
        @drop="onDrop($event, 'done')"
        class="flex-shrink-0"
      >
        <div
          v-for="task in doneTasks"
          :key="task.id"
          draggable="true"
          @dragstart="onDragStart($event, task)"
          class="transform active:scale-[0.98] transition-transform duration-100"
        >
          <TaskCard :task="task" :isFocused="uiStore.focusedTaskId === task.id" />
        </div>
      </KanbanColumn>
    </div>
  </div>
</template>
