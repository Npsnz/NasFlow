<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useWorkspaceStore } from '@/stores/workspace'
import { useTasksStore } from '@/stores/tasks'
import type { Task, Tag } from '@/stores/tasks'
import { useUIStore } from '@/stores/ui'
import {
  X,
  Plus,
  Trash2,
  MessageSquare,
  Clock,
  Eye,
  FileText
} from 'lucide-vue-next'

const workspaceStore = useWorkspaceStore()
const tasksStore = useTasksStore()
const uiStore = useUIStore()

const task = ref<Task | null>(null)
const loading = ref(false)
const saving = ref(false)
const showSavedIndicator = ref(false)
const activeTab = ref<'write' | 'preview'>('write')

// New checklist subtask state
const newSubtaskTitle = ref('')

// New comment state
const newCommentContent = ref('')

// Tag autocomplete states
const tagSearch = ref('')
const showTagSelector = ref(false)
const newTagName = ref('')
const newTagColor = ref('#171717')
const presetColors = ['#171717', '#1D9E75', '#D85A30', '#333333', '#ef4444', '#3b82f6', '#10b981', '#f59e0b']

const getPriorityText = (prio: string) => {
  switch (prio) {
    case 'low':
      return 'ต่ำ'
    case 'high':
      return 'สูง'
    case 'urgent':
      return 'ด่วนที่สุด'
    case 'medium':
    default:
      return 'ปานกลาง'
  }
}

const formatRelativeThaiTime = (dateStr: string) => {
  const date = new Date(dateStr)
  const seconds = Math.floor((new Date().getTime() - date.getTime()) / 1000)
  if (seconds < 60) return 'เมื่อครู่'
  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return `${minutes} นาทีที่แล้ว`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours} ชั่วโมงที่แล้ว`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days} วันที่แล้ว`
  return date.toLocaleDateString('th-TH', { day: 'numeric', month: 'short', year: '2-digit' })
}

// Open drawer on receiving event
const openTaskDrawer = async (e: Event) => {
  const taskId = (e as CustomEvent).detail
  loading.value = true
  try {
    const data = await tasksStore.getTaskDetails(taskId)
    task.value = data
    newSubtaskTitle.value = ''
    newCommentContent.value = ''
    tagSearch.value = ''
    showTagSelector.value = false
  } catch (err) {
    uiStore.showToast('ไม่สามารถโหลดข้อมูลรายละเอียดงานได้', 'error')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  window.addEventListener('open-task-details', openTaskDrawer)
})

onUnmounted(() => {
  window.removeEventListener('open-task-details', openTaskDrawer)
})

const closeDrawer = () => {
  task.value = null
}

const triggerAutoSave = async () => {
  if (!task.value) return
  saving.value = true
  try {
    const updated = await tasksStore.updateTask(task.value.id, {
      title: task.value.title,
      description: task.value.description,
      status: task.value.status,
      priority: task.value.priority,
      workspace_id: task.value.workspace_id,
      due_date: task.value.due_date ? new Date(task.value.due_date).toISOString() : null,
      is_recurring: task.value.is_recurring,
      recur_rule: task.value.recur_rule
    })
    // Sync local object
    task.value.completed_at = updated.completed_at

    // Trigger visual save indicator
    showSavedIndicator.value = true
    setTimeout(() => {
      showSavedIndicator.value = false
    }, 1500)
  } catch (err) {
    uiStore.showToast('บันทึกการแก้ไขไม่สำเร็จ', 'error')
  } finally {
    saving.value = false
  }
}

// Priority color circles helper
const priorityColors = {
  low: 'bg-green-500',
  medium: 'bg-amber-500',
  high: 'bg-red-500',
  urgent: 'bg-red-700 animate-pulse'
}

// Subtask checkbox toggle
const toggleSubtask = async (sub: Task) => {
  const nextStatus = sub.status === 'done' ? 'todo' : 'done'
  try {
    await tasksStore.updateTask(sub.id, { status: nextStatus })
    // Refresh task to pull nested list
    if (task.value) {
      const refreshed = await tasksStore.getTaskDetails(task.value.id)
      task.value.subtasks = refreshed.subtasks
    }
  } catch (err) {
    uiStore.showToast('เปลี่ยนสถานะงานย่อยล้มเหลว', 'error')
  }
}

// Inline create subtask
const handleAddSubtask = async () => {
  if (!newSubtaskTitle.value.trim() || !task.value) return
  try {
    await tasksStore.createTask({
      title: newSubtaskTitle.value.trim(),
      workspace_id: task.value.workspace_id,
      parent_task_id: task.value.id,
      status: 'todo',
      priority: 'medium'
    })
    newSubtaskTitle.value = ''
    
    // Refresh nested list
    const refreshed = await tasksStore.getTaskDetails(task.value.id)
    task.value.subtasks = refreshed.subtasks
    uiStore.showToast('เพิ่มงานย่อยสำเร็จ', 'success')
  } catch (err) {
    uiStore.showToast('ไม่สามารถสร้างงานย่อยได้', 'error')
  }
}

// Inline delete subtask
const handleDeleteSubtask = async (subId: number) => {
  if (!confirm('ลบงานย่อยนี้หรือไม่?')) return
  try {
    await tasksStore.deleteTask(subId)
    if (task.value) {
      task.value.subtasks = task.value.subtasks?.filter(s => s.id !== subId)
    }
    uiStore.showToast('ลบงานย่อยสำเร็จ', 'success')
  } catch (err) {
    uiStore.showToast('ลบงานย่อยล้มเหลว', 'error')
  }
}

// Comments submit on Ctrl+Enter
const submitComment = async () => {
  if (!newCommentContent.value.trim() || !task.value) return
  try {
    const newComment = await tasksStore.addComment(task.value.id, newCommentContent.value.trim())
    if (!task.value.comments) task.value.comments = []
    task.value.comments.push(newComment)
    newCommentContent.value = ''
    uiStore.showToast('บันทึกความคิดเห็นแล้ว', 'success')
  } catch (err) {
    uiStore.showToast('ไม่สามารถแสดงความคิดเห็นได้', 'error')
  }
}

// Tag filtering and creation in autocomplete
const filteredTags = computed(() => {
  const query = tagSearch.value.toLowerCase().trim()
  if (!query) return tasksStore.tags
  return tasksStore.tags.filter(t => t.name.toLowerCase().includes(query))
})

const handleAddTagToTask = async (tag: Tag) => {
  if (!task.value) return
  const exists = task.value.tags.some(t => t.id === tag.id)
  if (exists) return // Avoid duplicates

  const tagIds = [...task.value.tags.map(t => t.id), tag.id]
  try {
    await tasksStore.updateTask(task.value.id, { tag_ids: tagIds })
    task.value.tags.push(tag)
    showSavedIndicator.value = true
    setTimeout(() => showSavedIndicator.value = false, 1500)
  } catch (err) {
    uiStore.showToast('เพิ่มแท็กไม่สำเร็จ', 'error')
  }
}

const handleRemoveTag = async (tagId: number) => {
  if (!task.value) return
  const tagIds = task.value.tags.filter(t => t.id !== tagId).map(t => t.id)
  try {
    await tasksStore.updateTask(task.value.id, { tag_ids: tagIds })
    task.value.tags = task.value.tags.filter(t => t.id !== tagId)
    showSavedIndicator.value = true
    setTimeout(() => showSavedIndicator.value = false, 1500)
  } catch (err) {
    uiStore.showToast('ลบแท็กออกจากงานไม่สำเร็จ', 'error')
  }
}

const handleCreateTagInline = async () => {
  if (!newTagName.value.trim()) return
  try {
    const newTag = await tasksStore.createTag(newTagName.value.trim(), newTagColor.value)
    newTagName.value = ''
    tagSearch.value = ''
    if (task.value) {
      handleAddTagToTask(newTag)
    }
  } catch (err: any) {
    uiStore.showToast(err || 'ไม่สามารถสร้างแท็กได้', 'error')
  }
}

// Deleting current task
const handleDeleteTask = async () => {
  if (!task.value) return
  if (confirm('ลบงานนี้หรือไม่? ไม่สามารถกู้คืนได้')) {
    try {
      await tasksStore.deleteTask(task.value.id)
      uiStore.showToast('ลบงานสำเร็จ', 'success')
      closeDrawer()
    } catch (err) {
      uiStore.showToast('ลบงานล้มเหลว', 'error')
    }
  }
}

// Helper to convert simple Markdown in preview
const renderMarkdown = (text: string) => {
  if (!text) return '<p class="text-slate-400 italic">ไม่มีรายละเอียด</p>'
  let html = text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
    .replace(/`(.*?)`/g, '<code class="bg-slate-100 dark:bg-slate-800 px-1 py-0.5 rounded font-mono text-xs text-red-500">$1</code>')
    .replace(/\n/g, '<br/>')
  return `<p class="text-sm leading-relaxed">${html}</p>`
}

const handleDateChange = (e: any) => {
  if (!task.value) return
  // GORM expects timestamp string
  const val = e.target.value
  task.value.due_date = val ? new Date(val).toISOString() : null
  triggerAutoSave()
}
</script>

<template>
  <!-- Backdrop -->
  <Transition
    enter-active-class="transition ease-out duration-200"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition ease-in duration-150"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="task"
      @click="closeDrawer"
      class="fixed inset-0 z-40 bg-slate-900/40 backdrop-blur-sm hidden md:block"
    ></div>
  </Transition>

  <!-- Drawer Panel -->
  <Transition
    enter-active-class="transform transition ease-out duration-200"
    enter-from-class="translate-x-full"
    enter-to-class="translate-x-0"
    leave-active-class="transform transition ease-in duration-150"
    leave-from-class="translate-x-0"
    leave-to-class="translate-x-full"
  >
    <div
      v-if="task"
      class="fixed inset-y-0 right-0 z-50 w-full md:w-[480px] bg-white dark:bg-surface-dark border-l border-border-light dark:border-border-dark shadow-2xl flex flex-col justify-between"
    >
      <!-- Panel Header -->
      <div class="h-16 border-b border-border-light dark:border-border-dark flex items-center justify-between px-6">
        <div class="flex items-center space-x-2.5">
          <button
            @click="handleDeleteTask"
            class="text-slate-400 hover:text-red-500 transition-colors p-2 rounded-lg"
            title="ลบงานนี้"
          >
            <Trash2 class="w-4 h-4" />
          </button>
          <span
            v-if="showSavedIndicator"
            class="text-xs font-semibold text-green-500 flex items-center space-x-1"
          >
            <span>บันทึกแล้ว ✓</span>
          </span>
        </div>

        <button
          @click="closeDrawer"
          class="w-9 h-9 rounded-full hover:bg-slate-100 dark:hover:bg-slate-800 flex items-center justify-center text-slate-400 hover:text-slate-600 transition-colors"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Panel Body (Scrollable) -->
      <div class="flex-grow overflow-y-auto px-6 py-5 space-y-6">
        <!-- Task Title -->
        <div>
          <input
            type="text"
            v-model="task.title"
            @blur="triggerAutoSave"
            class="w-full text-lg font-bold bg-transparent border-0 focus:ring-0 focus:border-b focus:border-brand-500 pb-1 text-slate-900 dark:text-white"
            placeholder="หัวข้องาน..."
          />
        </div>

        <!-- Task Metadata Grid -->
        <div class="grid grid-cols-3 gap-y-4 gap-x-2 text-sm border-y border-border-light dark:border-border-dark py-4">
          <!-- Status -->
          <div class="text-slate-400 flex items-center space-x-1">
            <span>สถานะ</span>
          </div>
          <div class="col-span-2">
            <select
              v-model="task.status"
              @change="triggerAutoSave"
              class="w-full bg-slate-50 border border-border-light rounded-lg px-3 py-2 focus:outline-none dark:bg-slate-900/60 dark:border-border-dark text-base"
            >
              <option value="todo">รอดำเนินการ</option>
              <option value="doing">กำลังดำเนินการ</option>
              <option value="done">เสร็จแล้ว</option>
              <option value="archived">จัดเก็บ (Archive)</option>
            </select>
          </div>

          <!-- Priority -->
          <div class="text-slate-400 flex items-center space-x-1">
            <span>ความสำคัญ</span>
          </div>
          <div class="col-span-2">
            <div class="flex bg-slate-50 dark:bg-slate-900/60 border border-border-light dark:border-border-dark p-0.5 rounded-lg">
              <button
                v-for="p in ['low', 'medium', 'high', 'urgent']"
                :key="p"
                @click="task.priority = p as any; triggerAutoSave()"
                class="flex-1 text-sm font-medium py-2 px-2.5 rounded transition-all flex items-center justify-center space-x-1"
                :class="task.priority === p ? 'bg-white dark:bg-surface-dark shadow text-slate-900 dark:text-white font-bold' : 'text-slate-500 hover:text-slate-700'"
              >
                <span class="w-1.5 h-1.5 rounded-full" :class="priorityColors[p as 'low'|'medium'|'high'|'urgent']"></span>
                <span>{{ getPriorityText(p) }}</span>
              </button>
            </div>
          </div>

          <!-- Workspace -->
          <div class="text-slate-400 flex items-center space-x-1">
            <span>พื้นที่ทำงาน</span>
          </div>
          <div class="col-span-2">
            <select
              v-model="task.workspace_id"
              @change="triggerAutoSave"
              class="w-full bg-slate-50 border border-border-light rounded-lg px-3 py-2 focus:outline-none dark:bg-slate-900/60 dark:border-border-dark text-base"
            >
              <option
                v-for="ws in workspaceStore.workspaces"
                :key="ws.id"
                :value="ws.id"
              >
                {{ ws.name }}
              </option>
            </select>
          </div>

          <!-- Due Date -->
          <div class="text-slate-400 flex items-center space-x-1">
            <span>วันกำหนดส่ง</span>
          </div>
          <div class="col-span-2 relative">
            <input
              type="date"
              :value="task.due_date ? task.due_date.substring(0, 10) : ''"
              @change="handleDateChange"
              class="w-full bg-slate-50 border border-border-light rounded-lg px-3 py-2 focus:outline-none dark:bg-slate-900/60 dark:border-border-dark text-base"
            />
          </div>

          <!-- Recurring -->
          <div class="text-slate-400 flex items-center space-x-1">
            <span>ทำซ้ำ</span>
          </div>
          <div class="col-span-2 flex items-center space-x-2">
            <input
              type="checkbox"
              v-model="task.is_recurring"
              @change="triggerAutoSave"
              class="rounded text-brand-500 focus:ring-brand-500"
            />
            <select
              v-if="task.is_recurring"
              v-model="task.recur_rule"
              @change="triggerAutoSave"
              class="flex-grow bg-slate-50 border border-border-light rounded-lg px-2 py-0.5 focus:outline-none dark:bg-slate-900/60 dark:border-border-dark text-xs"
            >
              <option value="daily">ทุกวัน (Daily)</option>
              <option value="weekly">ทุกสัปดาห์ (Weekly)</option>
              <option value="monthly">ทุกเดือน (Monthly)</option>
              <option value="weekdays">วันทำการ (จ-ศ)</option>
            </select>
          </div>
        </div>

        <!-- Tags autocomplete section -->
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <h4 class="text-xs font-bold text-slate-400 uppercase tracking-wider">แท็กประจำงาน</h4>
            <button
              @click="showTagSelector = !showTagSelector"
              class="text-xs text-brand-500 hover:text-brand-600 flex items-center space-x-1 min-h-[44px]"
            >
              <Plus class="w-3.5 h-3.5" />
              <span>เพิ่มแท็ก</span>
            </button>
          </div>

          <!-- Tags display list -->
          <div class="flex flex-wrap gap-1.5">
            <span
              v-for="tag in task.tags"
              :key="tag.id"
              class="inline-flex items-center space-x-1 px-2.5 py-0.5 rounded-full text-xs font-semibold select-none border"
              :style="{ backgroundColor: tag.color + '15', color: tag.color, borderColor: tag.color + '30' }"
            >
              <span>{{ tag.name }}</span>
              <button @click="handleRemoveTag(tag.id)" class="hover:text-red-500 font-bold ml-1 text-[10px]">×</button>
            </span>
            <p v-if="!task.tags || task.tags.length === 0" class="text-slate-400 text-xs italic">ไม่มีแท็ก</p>
          </div>

          <!-- Tag Autocomplete Box -->
          <div
            v-if="showTagSelector"
            class="bg-slate-50 dark:bg-slate-900/50 p-3.5 border border-border-light dark:border-border-dark rounded-lg space-y-3 shadow-sm"
          >
            <!-- Search & Add existing -->
            <div>
              <input
                type="text"
                v-model="tagSearch"
                placeholder="ค้นหาแท็กที่มีอยู่..."
                class="w-full text-xs bg-white dark:bg-slate-950 border border-border-light dark:border-border-dark rounded px-2.5 py-1.5 focus:outline-none min-h-[44px]"
              />
              <div v-if="tagSearch" class="max-h-24 overflow-y-auto mt-2 border border-border-light dark:border-border-dark rounded bg-white dark:bg-slate-950 divide-y divide-border-light dark:divide-border-dark">
                <button
                  v-for="tag in filteredTags"
                  :key="tag.id"
                  @click="handleAddTagToTask(tag)"
                  class="w-full px-3 py-1.5 text-xs hover:bg-slate-100 dark:hover:bg-slate-800 flex items-center space-x-2"
                >
                  <span class="w-2.5 h-2.5 rounded-full" :style="{ backgroundColor: tag.color }"></span>
                  <span>{{ tag.name }}</span>
                </button>
                <p v-if="filteredTags.length === 0" class="px-3 py-2 text-slate-400 italic text-[11px]">ไม่พบแท็ก</p>
              </div>
            </div>

            <!-- Create Tag inline -->
            <div class="border-t border-border-light dark:border-border-dark pt-3 space-y-2">
              <p class="text-[10px] font-semibold text-slate-400 uppercase">สร้างแท็กใหม่</p>
              <div class="flex items-center space-x-2">
                <input
                  type="text"
                  v-model="newTagName"
                  placeholder="ชื่อแท็ก..."
                  class="flex-grow text-xs bg-white dark:bg-slate-950 border border-border-light dark:border-border-dark rounded px-2.5 py-1.5 focus:outline-none min-h-[44px]"
                />
                
                <!-- Color picker preset swatches -->
                <div class="flex items-center space-x-1">
                  <button
                    v-for="color in presetColors.slice(0, 4)"
                    :key="color"
                    @click="newTagColor = color"
                    class="w-4.5 h-4.5 rounded-full border"
                    :class="newTagColor === color ? 'border-slate-800 scale-110' : 'border-transparent'"
                    :style="{ backgroundColor: color }"
                  ></button>
                </div>

                <button
                  @click="handleCreateTagInline"
                  class="px-3 py-1.5 bg-brand-500 hover:bg-brand-600 text-white rounded text-xs font-semibold min-h-[44px]"
                >
                  สร้าง
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Description Markdown Editor/Preview -->
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <h4 class="text-xs font-bold text-slate-400 uppercase tracking-wider">รายละเอียด</h4>
            
            <div class="flex space-x-1 bg-slate-100 dark:bg-slate-800 p-0.5 rounded-lg">
              <button
                @click="activeTab = 'write'"
                class="px-2.5 py-1 text-[10px] font-bold rounded"
                :class="activeTab === 'write' ? 'bg-white dark:bg-surface-dark shadow text-slate-900 dark:text-white' : 'text-slate-500 hover:text-slate-700'"
              >
                <FileText class="w-3.5 h-3.5 inline mr-1" />เขียน
              </button>
              <button
                @click="activeTab = 'preview'"
                class="px-2.5 py-1 text-[10px] font-bold rounded"
                :class="activeTab === 'preview' ? 'bg-white dark:bg-surface-dark shadow text-slate-900 dark:text-white' : 'text-slate-500 hover:text-slate-700'"
              >
                <Eye class="w-3.5 h-3.5 inline mr-1" />พรีวิว
              </button>
            </div>
          </div>

          <div v-if="activeTab === 'write'">
            <textarea
              v-model="task.description"
              @blur="triggerAutoSave"
              rows="4"
              class="w-full text-sm bg-slate-50 border border-border-light rounded-lg p-3 focus:outline-none focus:ring-1 focus:ring-brand-500 dark:bg-slate-900/60 dark:border-border-dark dark:focus:bg-slate-900 text-slate-800 dark:text-slate-200"
              placeholder="ใส่รายละเอียดงานที่นี่... รองรับรูปแบบ Markdown แบบย่อ (เช่น **ตัวหนา**, *ตัวเอียง*, `โค้ด`)"
            ></textarea>
          </div>
          <div
            v-else
            class="p-3 bg-slate-50 dark:bg-slate-900/40 border border-border-light dark:border-border-dark rounded-lg min-h-[100px] overflow-x-auto"
            v-html="renderMarkdown(task.description)"
          ></div>
        </div>

        <!-- Subtasks Checklist Section -->
        <div class="space-y-3">
          <h4 class="text-xs font-bold text-slate-400 uppercase tracking-wider">รายการย่อย (Subtasks)</h4>
          
          <div class="space-y-2">
            <div
              v-for="sub in task.subtasks"
              :key="sub.id"
              class="flex items-center justify-between p-2.5 bg-slate-50 dark:bg-slate-900/30 rounded-lg group"
            >
              <div class="flex items-center space-x-3 overflow-hidden flex-grow">
                <input
                  type="checkbox"
                  :checked="sub.status === 'done'"
                  @change="toggleSubtask(sub)"
                  class="rounded text-brand-500 focus:ring-brand-500"
                />
                <span
                  class="text-sm truncate"
                  :class="sub.status === 'done' ? 'line-through text-slate-400 dark:text-slate-500' : 'text-slate-700 dark:text-slate-300'"
                >
                  {{ sub.title }}
                </span>
              </div>
              
              <button
                @click="handleDeleteSubtask(sub.id)"
                class="opacity-0 group-hover:opacity-100 text-slate-400 hover:text-red-500 transition-opacity p-1 min-h-[36px] min-w-[36px]"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </div>

            <!-- Inline add subtask input -->
            <div class="flex items-center space-x-2 mt-2">
              <input
                type="text"
                v-model="newSubtaskTitle"
                placeholder="+ เพิ่มรายการย่อย... (กด Enter)"
                class="flex-grow bg-slate-50 dark:bg-slate-900/40 border border-border-light dark:border-border-dark rounded-lg px-3 py-2 text-xs focus:outline-none min-h-[44px]"
                @keydown.enter="handleAddSubtask"
              />
            </div>
          </div>
        </div>

        <!-- Comments Section -->
        <div class="space-y-4 border-t border-border-light dark:border-border-dark pt-5">
          <h4 class="text-xs font-bold text-slate-400 uppercase tracking-wider flex items-center space-x-1.5">
            <MessageSquare class="w-3.5 h-3.5" />
            <span>ความคิดเห็น ({{ task.comments?.length || 0 }})</span>
          </h4>

          <!-- Comments Feed -->
          <div class="space-y-4.5 max-h-64 overflow-y-auto px-1">
            <div
              v-for="comment in task.comments"
              :key="comment.id"
              class="flex space-x-3 text-xs"
            >
              <img
                v-if="comment.user?.avatar_url"
                :src="comment.user.avatar_url"
                class="w-7 h-7 rounded-full object-cover"
                alt="Avatar"
              />
              <div v-else class="w-7 h-7 rounded-full bg-slate-100 dark:bg-slate-800 text-slate-500 flex items-center justify-center font-bold">
                <User class="w-3 h-3" />
              </div>
              <div class="flex-grow bg-slate-50 dark:bg-slate-900/40 p-2.5 rounded-lg relative group">
                <div class="flex items-center justify-between mb-1">
                  <span class="font-semibold text-slate-800 dark:text-slate-300">{{ comment.user?.name || 'นิรนาม' }}</span>
                  <span class="text-[9px] text-slate-400 flex items-center space-x-0.5">
                    <Clock class="w-2.5 h-2.5" />
                    <span>{{ formatRelativeThaiTime(comment.created_at) }}</span>
                  </span>
                </div>
                <p class="text-slate-600 dark:text-slate-400 leading-relaxed break-words whitespace-pre-wrap">{{ comment.content }}</p>
              </div>
            </div>
            
            <p v-if="!task.comments || task.comments.length === 0" class="text-slate-400 text-xs italic text-center py-4">
              ไม่มีการพูดคุยเกี่ยวกับงานนี้
            </p>
          </div>

          <!-- Add Comment Input -->
          <div class="space-y-2">
            <textarea
              v-model="newCommentContent"
              rows="2"
              class="w-full text-xs bg-slate-50 border border-border-light rounded-lg p-2.5 focus:outline-none focus:ring-1 focus:ring-brand-500 dark:bg-slate-900/60 dark:border-border-dark dark:focus:bg-slate-900 text-slate-800 dark:text-slate-200"
              placeholder="เขียนข้อความของคุณ... (กด Ctrl+Enter เพื่อส่ง)"
              @keydown.ctrl.enter="submitComment"
            ></textarea>
            <div class="flex justify-end">
              <button
                @click="submitComment"
                class="px-4 py-2 bg-brand-500 hover:bg-brand-600 text-white rounded-lg text-xs font-semibold shadow shadow-brand-500/15 min-h-[44px]"
              >
                ส่งความคิดเห็น
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>
