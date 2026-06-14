<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useTasksStore, Task, Tag } from '@/stores/tasks'
import { useWorkspaceStore } from '@/stores/workspace'
import { useUIStore } from '@/stores/ui'
import {
  FileDown,
  Trash2,
  Tag as TagIcon,
  CheckCircle,
  HelpCircle,
  ChevronDown,
  ChevronRight,
  Sparkles,
  ArrowUpDown,
  FolderClosed,
  Calendar
} from 'lucide-vue-next'

const tasksStore = useTasksStore()
const workspaceStore = useWorkspaceStore()
const uiStore = useUIStore()

const groupBy = ref<'none' | 'status' | 'priority' | 'workspace' | 'due_date'>('status')
const sortBy = ref<'due_date' | 'priority' | 'created' | 'alphabetical'>('due_date')

// Selected tasks for bulk actions
const selectedTaskIds = ref<number[]>([])

// Inline editing states
const editingTaskId = ref<number | null>(null)
const editingField = ref<string | null>(null)

// Pagination state
const currentPage = ref(1)
const itemsPerPage = 20

onMounted(() => {
  tasksStore.fetchTasks()
  tasksStore.fetchTags()
  workspaceStore.fetchWorkspaces()
})

const getWorkspaceName = (wsId: number) => {
  const ws = workspaceStore.workspaces.find(w => w.id === wsId)
  return ws ? ws.name : 'ไม่มีพื้นที่ทำงาน'
}

const getWorkspaceColor = (wsId: number) => {
  const ws = workspaceStore.workspaces.find(w => w.id === wsId)
  return ws ? ws.color : '#cbd5e1'
}

// Reactively filter tasks list
const baseFilteredTasks = computed(() => {
  return tasksStore.tasks.filter((t) => {
    // 1. Search filter
    if (tasksStore.filters.search) {
      const q = tasksStore.filters.search.toLowerCase()
      if (!t.title.toLowerCase().includes(q)) return false
    }

    // 2. Priority Filter
    if (tasksStore.filters.priority) {
      if (t.priority !== tasksStore.filters.priority) return false
    }

    // 3. Tags Filter
    if (tasksStore.filters.tagIds.length > 0) {
      const matches = t.tags.some(tag => tasksStore.filters.tagIds.includes(tag.id))
      if (!matches) return false
    }

    // 4. Due Date Range Preset Filter
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
          const diff = d.getTime() - now.getTime()
          const diffDays = diff / (1000 * 3600 * 24)
          return diffDays >= -1 && diffDays <= 7
        }
        case 'overdue':
          return d < now && t.status !== 'done'
      }
    }

    return true
  })
})

// Sort tasks list
const sortedTasks = computed(() => {
  const tasksCopy = [...baseFilteredTasks.value]
  
  tasksCopy.sort((a, b) => {
    if (sortBy.value === 'due_date') {
      if (!a.due_date) return 1
      if (!b.due_date) return -1
      return new Date(a.due_date).getTime() - new Date(b.due_date).getTime()
    }
    
    if (sortBy.value === 'priority') {
      const weight = { urgent: 4, high: 3, medium: 2, low: 1 }
      return (weight[b.priority] || 0) - (weight[a.priority] || 0)
    }

    if (sortBy.value === 'alphabetical') {
      return a.title.localeCompare(b.title, 'th-TH')
    }

    // Default: 'created' (newest first)
    return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  })

  return tasksCopy
})

// Paginate tasks list (20 items per page)
const totalTasksCount = computed(() => sortedTasks.value.length)
const totalPages = computed(() => Math.ceil(totalTasksCount.value / itemsPerPage))

const paginatedTasks = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage
  return sortedTasks.value.slice(start, start + itemsPerPage)
})

watch([sortBy, groupBy, () => tasksStore.filters], () => {
  currentPage.value = 1 // Reset to page 1 on filter/sort changes
})

// Grouping logic
const groupedTasks = computed(() => {
  const list = paginatedTasks.value
  
  if (groupBy.value === 'none') {
    return [{ key: 'งานทั้งหมด', tasks: list }]
  }

  const groups: Record<string, Task[]> = {}

  list.forEach(t => {
    let key = 'ไม่มีข้อมูล'
    
    if (groupBy.value === 'status') {
      const statusMap = { todo: 'รอดำเนินการ', doing: 'กำลังดำเนินการ', done: 'เสร็จแล้ว', archived: 'จัดเก็บแล้ว' }
      key = statusMap[t.status] || t.status
    } else if (groupBy.value === 'priority') {
      const priorityMap = { urgent: 'ด่วนที่สุด (Urgent)', high: 'สูง (High)', medium: 'ปานกลาง (Medium)', low: 'ต่ำ (Low)' }
      key = priorityMap[t.priority] || t.priority
    } else if (groupBy.value === 'workspace') {
      key = getWorkspaceName(t.workspace_id)
    } else if (groupBy.value === 'due_date') {
      if (!t.due_date) {
        key = 'ไม่มีกำหนดส่ง'
      } else {
        const d = new Date(t.due_date)
        key = d.toLocaleDateString('th-TH', { weekday: 'long', day: 'numeric', month: 'short' })
      }
    }
    
    if (!groups[key]) groups[key] = []
    groups[key].push(t)
  })

  return Object.entries(groups).map(([key, tasks]) => ({ key, tasks }))
})

// Checkbox selection
const toggleSelectAll = (e: any) => {
  if (e.target.checked) {
    selectedTaskIds.value = paginatedTasks.value.map(t => t.id)
  } else {
    selectedTaskIds.value = []
  }
}

const toggleSelectTask = (id: number) => {
  const idx = selectedTaskIds.value.indexOf(id)
  if (idx === -1) {
    selectedTaskIds.value.push(id)
  } else {
    selectedTaskIds.value.splice(idx, 1)
  }
}

// Bulk Actions
const bulkMoveStatus = async (status: string) => {
  if (selectedTaskIds.value.length === 0) return
  try {
    const updates = selectedTaskIds.value.map(id => ({ id, sort_order: 0, status }))
    await tasksStore.reorderTasks(updates)
    uiStore.showToast(`ย้ายงานจำนวน ${selectedTaskIds.value.length} งานไปเรียบร้อยแล้ว`, 'success')
    selectedTaskIds.value = []
  } catch (err) {
    uiStore.showToast('ไม่สามารถย้ายงานทั้งหมดได้', 'error')
  }
}

const bulkDelete = async () => {
  if (selectedTaskIds.value.length === 0) return
  if (confirm(`คุณแน่ใจว่าต้องการลบงานที่เลือกจำนวน ${selectedTaskIds.value.length} งานใช่หรือไม่?`)) {
    try {
      for (const id of selectedTaskIds.value) {
        await tasksStore.deleteTask(id)
      }
      uiStore.showToast('ลบงานที่เลือกสำเร็จ', 'success')
      selectedTaskIds.value = []
    } catch (err) {
      uiStore.showToast('การลบงานบางงานล้มเหลว', 'error')
    }
  }
}

const bulkAddTag = async (tagId: number) => {
  if (selectedTaskIds.value.length === 0) return
  try {
    for (const id of selectedTaskIds.value) {
      const task = tasksStore.tasks.find(t => t.id === id)
      if (task) {
        const tagIds = [...task.tags.map(t => t.id), tagId]
        await tasksStore.updateTask(id, { tag_ids: tagIds })
      }
    }
    uiStore.showToast('เพิ่มแท็กให้กับงานที่เลือกสำเร็จ', 'success')
    selectedTaskIds.value = []
  } catch (err) {
    uiStore.showToast('เพิ่มแท็กให้กับงานทั้งหมดล้มเหลว', 'error')
  }
}

// Client-side CSV Export
const exportCSV = () => {
  const headers = ['ID', 'ชื่องาน', 'รายละเอียด', 'สถานะ', 'ความสำคัญ', 'พื้นที่ทำงาน', 'วันส่งงาน', 'วันที่เสร็จ']
  const rows = sortedTasks.value.map(t => [
    t.id,
    `"${t.title.replace(/"/g, '""')}"`,
    `"${t.description ? t.description.replace(/"/g, '""') : ''}"`,
    t.status,
    t.priority,
    `"${getWorkspaceName(t.workspace_id).replace(/"/g, '""')}"`,
    t.due_date ? t.due_date.substring(0, 10) : 'ไม่มี',
    t.completed_at ? t.completed_at.substring(0, 10) : 'ไม่มี'
  ])

  const csvContent = "data:text/csv;charset=utf-8,\uFEFF" 
    + [headers.join(','), ...rows.map(r => r.join(','))].join('\n')

  const encodedUri = encodeURI(csvContent)
  const link = document.createElement("a")
  link.setAttribute("href", encodedUri)
  link.setAttribute("download", "taskflow_export.csv")
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  uiStore.showToast('ส่งออกไฟล์ CSV สำเร็จ', 'success')
}

// Inline Row Edits
const startEdit = (taskId: number, field: string) => {
  editingTaskId.value = taskId
  editingField.value = field
}

const finishInlineEdit = async (task: Task) => {
  editingTaskId.value = null
  editingField.value = null
  try {
    await tasksStore.updateTask(task.id, {
      title: task.title,
      status: task.status,
      priority: task.priority,
      workspace_id: task.workspace_id,
      due_date: task.due_date ? new Date(task.due_date).toISOString() : null
    })
    uiStore.showToast('บันทึกสำเร็จ ✓', 'success')
  } catch (err) {
    uiStore.showToast('แก้ไขงานไม่สำเร็จ', 'error')
  }
}

const openTaskDrawer = (id: number) => {
  window.dispatchEvent(new CustomEvent('open-task-details', { detail: id }))
}
</script>

<template>
  <div class="p-6 space-y-6">
    <!-- Filters and Control row -->
    <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 border-b border-border-light dark:border-border-dark pb-4">
      <div class="flex flex-wrap items-center gap-3">
        <!-- Group By Select -->
        <div class="flex items-center space-x-2">
          <span class="text-xs text-slate-400 font-medium">จัดกลุ่มตาม:</span>
          <select
            v-model="groupBy"
            class="text-xs bg-white border border-border-light rounded-lg px-2.5 py-1.5 focus:outline-none dark:bg-slate-900/60 dark:border-border-dark min-h-[44px]"
          >
            <option value="none">ไม่จัดกลุ่ม</option>
            <option value="status">สถานะ</option>
            <option value="priority">ความสำคัญ</option>
            <option value="workspace">พื้นที่ทำงาน</option>
            <option value="due_date">วันครบกำหนด</option>
          </select>
        </div>

        <!-- Sort By Select -->
        <div class="flex items-center space-x-2">
          <span class="text-xs text-slate-400 font-medium">จัดเรียงตาม:</span>
          <select
            v-model="sortBy"
            class="text-xs bg-white border border-border-light rounded-lg px-2.5 py-1.5 focus:outline-none dark:bg-slate-900/60 dark:border-border-dark min-h-[44px]"
          >
            <option value="due_date">วันกำหนดส่ง</option>
            <option value="priority">ความสำคัญ</option>
            <option value="created">วันสร้างล่าสุด</option>
            <option value="alphabetical">ตามตัวอักษร ก-ฮ</option>
          </select>
        </div>
      </div>

      <!-- Export and Counters -->
      <div class="flex items-center space-x-3 w-full sm:w-auto justify-end">
        <span class="text-xs text-slate-400 font-semibold">{{ sortedTasks.length }} งานที่ค้นพบ</span>
        <button
          @click="exportCSV"
          class="flex items-center space-x-1.5 px-3.5 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 dark:bg-surface-dark dark:text-slate-300 dark:hover:bg-slate-800 border border-border-light dark:border-border-dark rounded-lg text-xs font-semibold min-h-[44px]"
        >
          <FileDown class="w-4 h-4" />
          <span>ส่งออก CSV</span>
        </button>
      </div>
    </div>

    <!-- Bulk Action Overlay Panel (displays when tasks selected) -->
    <div
      v-if="selectedTaskIds.length > 0"
      class="bg-brand-50 border border-brand-200 dark:bg-brand-500/10 dark:border-brand-500/20 p-4 rounded-xl flex flex-col sm:flex-row items-center justify-between gap-4 shadow-lg animate-pulse"
    >
      <div class="flex items-center space-x-2">
        <span class="w-2.5 h-2.5 rounded-full bg-brand-500"></span>
        <span class="text-xs font-bold text-brand-700 dark:text-brand-300">เลือกอยู่ {{ selectedTaskIds.length }} งาน</span>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <!-- Move Status -->
        <select
          @change="bulkMoveStatus(($event.target as HTMLSelectElement).value)"
          class="text-xs bg-white border border-brand-200 rounded-lg px-2.5 py-1.5 focus:outline-none dark:bg-slate-950 dark:border-border-dark min-h-[44px]"
        >
          <option value="">เปลี่ยนสถานะเป็น...</option>
          <option value="todo">รอดำเนินการ</option>
          <option value="doing">กำลังดำเนินการ</option>
          <option value="done">เสร็จแล้ว</option>
        </select>

        <!-- Bulk Tag Selection -->
        <select
          @change="bulkAddTag(Number(($event.target as HTMLSelectElement).value))"
          class="text-xs bg-white border border-brand-200 rounded-lg px-2.5 py-1.5 focus:outline-none dark:bg-slate-950 dark:border-border-dark min-h-[44px]"
        >
          <option value="">เพิ่มแท็ก...</option>
          <option v-for="tag in tasksStore.tags" :key="tag.id" :value="tag.id">{{ tag.name }}</option>
        </select>

        <!-- Bulk Delete -->
        <button
          @click="bulkDelete"
          class="flex items-center space-x-1.5 px-3 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg text-xs font-semibold min-h-[44px]"
        >
          <Trash2 class="w-4 h-4" />
          <span>ลบรายการที่เลือก</span>
        </button>
      </div>
    </div>

    <!-- Grouped Table List -->
    <div class="space-y-8 select-none">
      <div v-for="group in groupedTasks" :key="group.key" class="space-y-2">
        <!-- Group Header Title -->
        <div class="flex items-center space-x-2 px-1">
          <span class="text-xs font-bold text-slate-800 dark:text-slate-200 uppercase tracking-wide">{{ group.key }}</span>
          <span class="text-[10px] text-slate-400 font-bold bg-slate-100 dark:bg-slate-800 px-2 py-0.5 rounded-full">
            {{ group.tasks.length }}
          </span>
        </div>

        <!-- Tasks Table Rows -->
        <div class="bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl overflow-hidden shadow-sm">
          <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse text-xs">
              <thead>
                <tr class="bg-slate-50 dark:bg-slate-900/40 text-slate-400 dark:text-slate-500 border-b border-border-light dark:border-border-dark font-medium uppercase tracking-wider text-[9px] select-none h-11">
                  <th class="p-3 w-10 text-center">
                    <input
                      type="checkbox"
                      :checked="selectedTaskIds.length === paginatedTasks.length && paginatedTasks.length > 0"
                      @change="toggleSelectAll"
                      class="rounded text-brand-500 focus:ring-brand-500"
                    />
                  </th>
                  <th class="p-3 w-16 text-center">ความสำคัญ</th>
                  <th class="p-3">ชื่องาน</th>
                  <th class="p-3 w-32">พื้นที่ทำงาน</th>
                  <th class="p-3 w-36">แท็ก</th>
                  <th class="p-3 w-32">กำหนดส่ง</th>
                  <th class="p-3 w-28">สถานะ</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-border-light dark:divide-border-dark text-slate-700 dark:text-slate-300">
                <tr
                  v-for="task in group.tasks"
                  :key="task.id"
                  class="hover:bg-slate-50/50 dark:hover:bg-slate-800/20 group h-12"
                >
                  <!-- Checkbox Selection -->
                  <td class="p-3 text-center">
                    <input
                      type="checkbox"
                      :checked="selectedTaskIds.includes(task.id)"
                      @change="toggleSelectTask(task.id)"
                      class="rounded text-brand-500 focus:ring-brand-500"
                    />
                  </td>

                  <!-- Priority dot inline edit -->
                  <td class="p-3 text-center">
                    <span v-if="editingTaskId !== task.id || editingField !== 'priority'" @click="startEdit(task.id, 'priority')" class="cursor-pointer inline-flex items-center justify-center">
                      <span
                        class="w-2.5 h-2.5 rounded-full"
                        :class="{
                          'bg-green-500': task.priority === 'low',
                          'bg-amber-500': task.priority === 'medium',
                          'bg-red-500': task.priority === 'high',
                          'bg-red-700 animate-pulse': task.priority === 'urgent'
                        }"
                      ></span>
                    </span>
                    <select
                      v-else
                      v-model="task.priority"
                      @change="finishInlineEdit(task)"
                      @blur="finishInlineEdit(task)"
                      class="text-[10px] bg-white border border-border-light rounded focus:outline-none dark:bg-slate-900"
                    >
                      <option value="low">ต่ำ</option>
                      <option value="medium">กลาง</option>
                      <option value="high">สูง</option>
                      <option value="urgent">ด่วน</option>
                    </select>
                  </td>

                  <!-- Task Title Edit & Click Drawer -->
                  <td class="p-3 font-medium">
                    <div class="flex items-center justify-between">
                      <span
                        v-if="editingTaskId !== task.id || editingField !== 'title'"
                        @click="startEdit(task.id, 'title')"
                        class="cursor-pointer hover:text-brand-500 truncate max-w-xs block"
                      >
                        {{ task.title }}
                      </span>
                      <input
                        v-else
                        type="text"
                        v-model="task.title"
                        @blur="finishInlineEdit(task)"
                        @keydown.enter="finishInlineEdit(task)"
                        class="bg-white border border-brand-500 rounded px-1.5 py-0.5 text-xs focus:outline-none dark:bg-slate-900 text-slate-800 dark:text-white"
                      />
                      <!-- Details button -->
                      <button
                        @click="openTaskDrawer(task.id)"
                        class="opacity-0 group-hover:opacity-100 text-brand-500 font-semibold text-[10px] hover:underline transition-opacity px-2 py-1 min-h-[32px]"
                      >
                        เปิดดูรายละเอียด
                      </button>
                    </div>
                  </td>

                  <!-- Workspace Badge Inline Edit -->
                  <td class="p-3">
                    <span
                      v-if="editingTaskId !== task.id || editingField !== 'workspace'"
                      @click="startEdit(task.id, 'workspace')"
                      class="cursor-pointer inline-flex items-center space-x-1.5 px-2 py-0.5 rounded text-[10px] font-semibold"
                      :style="{ backgroundColor: getWorkspaceColor(task.workspace_id) + '15', color: getWorkspaceColor(task.workspace_id) }"
                    >
                      <FolderClosed class="w-3 h-3" />
                      <span>{{ getWorkspaceName(task.workspace_id) }}</span>
                    </span>
                    <select
                      v-else
                      v-model="task.workspace_id"
                      @change="finishInlineEdit(task)"
                      @blur="finishInlineEdit(task)"
                      class="text-[10px] bg-white border border-border-light rounded focus:outline-none dark:bg-slate-900"
                    >
                      <option v-for="ws in workspaceStore.workspaces" :key="ws.id" :value="ws.id">{{ ws.name }}</option>
                    </select>
                  </td>

                  <!-- Tags -->
                  <td class="p-3">
                    <div class="flex flex-wrap gap-1">
                      <span
                        v-for="tag in task.tags"
                        :key="tag.id"
                        class="text-[9px] font-bold px-1.5 py-0.2 rounded border"
                        :style="{ backgroundColor: tag.color + '15', color: tag.color, borderColor: tag.color + '25' }"
                      >
                        {{ tag.name }}
                      </span>
                      <span v-if="task.tags.length === 0" class="text-slate-400 italic text-[10px]">ไม่มีแท็ก</span>
                    </div>
                  </td>

                  <!-- Due date inline edit -->
                  <td class="p-3">
                    <span
                      v-if="editingTaskId !== task.id || editingField !== 'due_date'"
                      @click="startEdit(task.id, 'due_date')"
                      class="cursor-pointer inline-flex items-center space-x-1 hover:text-brand-500 font-medium"
                    >
                      <Calendar class="w-3 h-3" />
                      <span>{{ task.due_date ? new Date(task.due_date).toLocaleDateString('th-TH', { day: 'numeric', month: 'short' }) : 'ไม่มีกำหนด' }}</span>
                    </span>
                    <input
                      v-else
                      type="date"
                      v-model="task.due_date"
                      @change="finishInlineEdit(task)"
                      @blur="finishInlineEdit(task)"
                      class="text-[10px] bg-white border border-border-light rounded focus:outline-none dark:bg-slate-900"
                    />
                  </td>

                  <!-- Status badge inline edit -->
                  <td class="p-3">
                    <span
                      v-if="editingTaskId !== task.id || editingField !== 'status'"
                      @click="startEdit(task.id, 'status')"
                      class="cursor-pointer inline-flex items-center px-2 py-0.5 rounded text-[10px] font-semibold"
                      :class="{
                        'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300': task.status === 'todo',
                        'bg-blue-50 text-blue-700 dark:bg-blue-950/20 dark:text-blue-400': task.status === 'doing',
                        'bg-green-50 text-green-700 dark:bg-green-950/20 dark:text-green-400': task.status === 'done'
                      }"
                    >
                      {{ task.status === 'todo' ? 'รอทำ' : task.status === 'doing' ? 'กำลังทำ' : 'เสร็จแล้ว' }}
                    </span>
                    <select
                      v-else
                      v-model="task.status"
                      @change="finishInlineEdit(task)"
                      @blur="finishInlineEdit(task)"
                      class="text-[10px] bg-white border border-border-light rounded focus:outline-none dark:bg-slate-900"
                    >
                      <option value="todo">รอดำเนินการ</option>
                      <option value="doing">กำลังดำเนินการ</option>
                      <option value="done">เสร็จแล้ว</option>
                    </select>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
      
      <!-- Empty state inside view -->
      <div
        v-if="sortedTasks.length === 0"
        class="text-center py-16 bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl select-none"
      >
        <ClipboardCopy class="w-12 h-12 text-slate-300 mx-auto mb-3" />
        <h4 class="text-sm font-bold text-slate-800 dark:text-slate-200">ไม่พบงานที่ต้องการ</h4>
        <p class="text-xs text-slate-400 mt-1">ลองล้างตัวกรองที่เลือกไว้ที่แถบด้านบนสุด</p>
        <button
          @click="clearAllFilters"
          class="mt-4 px-4 py-2 bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 text-slate-700 dark:text-slate-300 rounded-lg text-xs font-semibold min-h-[44px]"
        >
          ล้างตัวกรองทั้งหมด
        </button>
      </div>
    </div>

    <!-- Pagination controls (20 per page) -->
    <div v-if="totalPages > 1" class="flex items-center justify-between border-t border-border-light dark:border-border-dark pt-4 select-none">
      <span class="text-xs text-slate-400">หน้า {{ currentPage }} จาก {{ totalPages }} หน้า</span>
      <div class="flex items-center space-x-2">
        <button
          @click="currentPage = Math.max(currentPage - 1, 1)"
          :disabled="currentPage === 1"
          class="px-3.5 py-2 bg-white border border-border-light rounded-lg text-xs font-bold hover:bg-slate-50 disabled:opacity-50 min-h-[44px]"
        >
          ย้อนกลับ
        </button>
        <button
          @click="currentPage = Math.min(currentPage + 1, totalPages)"
          :disabled="currentPage === totalPages"
          class="px-3.5 py-2 bg-white border border-border-light rounded-lg text-xs font-bold hover:bg-slate-50 disabled:opacity-50 min-h-[44px]"
        >
          ถัดไป
        </button>
      </div>
    </div>
  </div>
</template>
