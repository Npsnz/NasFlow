<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { useWorkspaceStore } from '@/stores/workspace'
import { useTasksStore } from '@/stores/tasks'
import { useUIStore } from '@/stores/ui'
import { Calendar, AlertTriangle, RefreshCw, FolderClosed, X } from 'lucide-vue-next'

const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const workspaceStore = useWorkspaceStore()
const tasksStore = useTasksStore()
const uiStore = useUIStore()

const inputRef = ref<HTMLInputElement | null>(null)
const rawText = ref('')
const loading = ref(false)

// Focus input automatically on mount
watch(() => props.show, (newVal) => {
  if (newVal) {
    nextTick(() => {
      inputRef.value?.focus()
    })
  }
}, { immediate: true })

const parsed = computed(() => {
  const text = rawText.value
  let cleanedTitle = text
  let parsedWorkspaceId = workspaceStore.currentWorkspace?.id || (workspaceStore.workspaces[0]?.id || 0)
  let parsedPriority: 'low' | 'medium' | 'high' | 'urgent' = 'medium'
  let parsedDueDate: Date | null = null
  let parsedIsRecurring = false
  let parsedRecurRule = ''

  // 1. Parse Workspace: #(\S+)
  const wsMatch = text.match(/#([^\s!@]+)/)
  if (wsMatch) {
    const wsName = wsMatch[1]
    const found = workspaceStore.workspaces.find(
      w => w.name.toLowerCase() === wsName.toLowerCase() || w.slug.toLowerCase() === wsName.toLowerCase()
    )
    if (found) {
      parsedWorkspaceId = found.id
      cleanedTitle = cleanedTitle.replace(wsMatch[0], '')
    }
  }

  // 2. Parse Priority: !(\S+)
  const prioMatch = text.match(/!([^\s#@]+)/)
  if (prioMatch) {
    const pVal = prioMatch[1].toLowerCase()
    if (pVal === 'high' || pVal === 'สูง') {
      parsedPriority = 'high'
      cleanedTitle = cleanedTitle.replace(prioMatch[0], '')
    } else if (pVal === 'urgent' || pVal === 'ด่วน') {
      parsedPriority = 'urgent'
      cleanedTitle = cleanedTitle.replace(prioMatch[0], '')
    } else if (pVal === 'low' || pVal === 'ต่ำ') {
      parsedPriority = 'low'
      cleanedTitle = cleanedTitle.replace(prioMatch[0], '')
    } else if (pVal === 'medium' || pVal === 'กลาง') {
      parsedPriority = 'medium'
      cleanedTitle = cleanedTitle.replace(prioMatch[0], '')
    }
  }

  // 3. Parse Due Date: @(\S+)
  const dateMatch = text.match(/@([^\s#!]+)/)
  if (dateMatch) {
    const dVal = dateMatch[1].toLowerCase()
    const now = new Date()
    if (dVal === 'วันนี้' || dVal === 'today') {
      parsedDueDate = now
      cleanedTitle = cleanedTitle.replace(dateMatch[0], '')
    } else if (dVal === 'พรุ่งนี้' || dVal === 'tomorrow') {
      const tomorrow = new Date()
      tomorrow.setDate(now.getDate() + 1)
      parsedDueDate = tomorrow
      cleanedTitle = cleanedTitle.replace(dateMatch[0], '')
    } else if (dVal === 'จันทร์' || dVal === 'monday') {
      const nextMon = new Date()
      nextMon.setDate(now.getDate() + ((1 + 7 - now.getDay()) % 7 || 7))
      parsedDueDate = nextMon
      cleanedTitle = cleanedTitle.replace(dateMatch[0], '')
    }
  }

  // 4. Parse Recurrence: ทุกวัน / every day
  if (text.includes('every day') || text.includes('ทุกวัน')) {
    parsedIsRecurring = true
    parsedRecurRule = 'daily'
    cleanedTitle = cleanedTitle.replace('every day', '').replace('ทุกวัน', '')
  } else if (text.includes('every week') || text.includes('ทุกสัปดาห์')) {
    parsedIsRecurring = true
    parsedRecurRule = 'weekly'
    cleanedTitle = cleanedTitle.replace('every week', '').replace('ทุกสัปดาห์', '')
  } else if (text.includes('every month') || text.includes('ทุกเดือน')) {
    parsedIsRecurring = true
    parsedRecurRule = 'monthly'
    cleanedTitle = cleanedTitle.replace('every month', '').replace('ทุกเดือน', '')
  } else if (text.includes('weekdays') || text.includes('วันทำการ')) {
    parsedIsRecurring = true
    parsedRecurRule = 'weekdays'
    cleanedTitle = cleanedTitle.replace('weekdays', '').replace('วันทำการ', '')
  }

  return {
    title: cleanedTitle.replace(/\s+/g, ' ').trim(),
    workspaceId: parsedWorkspaceId,
    priority: parsedPriority,
    dueDate: parsedDueDate,
    isRecurring: parsedIsRecurring,
    recurRule: parsedRecurRule
  }
})

const activeWorkspace = computed(() => {
  return workspaceStore.workspaces.find(w => w.id === parsed.value.workspaceId)
})

const getPriorityText = (prio: string) => {
  switch (prio) {
    case 'low':
      return 'ต่ำ (Low)'
    case 'high':
      return 'สูง (High)'
    case 'urgent':
      return 'ด่วนที่สุด (Urgent)'
    case 'medium':
    default:
      return 'ปานกลาง (Medium)'
  }
}

const formatDate = (d: Date | null) => {
  if (!d) return ''
  return d.toLocaleDateString('th-TH', { day: 'numeric', month: 'short' })
}

const getRecurText = (rule: string) => {
  switch (rule) {
    case 'daily':
      return 'ทุกวัน'
    case 'weekly':
      return 'ทุกสัปดาห์'
    case 'monthly':
      return 'ทุกเดือน'
    case 'weekdays':
      return 'วันทำการ (จ-ศ)'
    default:
      return ''
  }
}

const handleCreate = async () => {
  const payload = parsed.value
  if (!payload.title) {
    uiStore.showToast('กรุณาระบุหัวข้องาน', 'error')
    return
  }

  loading.value = true
  try {
    const createdTask = await tasksStore.createTask({
      title: payload.title,
      workspace_id: payload.workspaceId,
      priority: payload.priority,
      status: 'todo',
      due_date: payload.dueDate ? payload.dueDate.toISOString() : undefined,
      is_recurring: payload.isRecurring,
      recur_rule: payload.recurRule
    })

    uiStore.showToast('สร้างงานสำเร็จแล้ว ✓', 'success')
    
    // Broadcast notification or trigger drawer open
    // We will show a toast containing a link or trigger an event
    // To implement the "ดูงาน link that opens drawer" action:
    // We dispatch a custom window event that TaskDrawer will capture to slide open!
    nextTick(() => {
      window.dispatchEvent(new CustomEvent('open-task-details', { detail: createdTask.id }))
    })

    rawText.value = ''
    emit('close')
  } catch (err: any) {
    uiStore.showToast(err || 'ไม่สามารถสร้างงานได้', 'error')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <Transition
    enter-active-class="ease-out duration-200"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="ease-in duration-150"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="show"
      class="fixed inset-0 z-50 flex items-start justify-center pt-[15vh] p-4 bg-slate-900/60 dark:bg-slate-950/70 backdrop-blur-sm"
      @click.self="emit('close')"
    >
      <div class="w-full max-w-xl bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl shadow-2xl overflow-hidden flex flex-col">
        <!-- Input area -->
        <div class="p-4 flex items-center space-x-3 border-b border-border-light dark:border-border-dark">
          <input
            ref="inputRef"
            type="text"
            v-model="rawText"
            placeholder="พิมพ์ชื่องานด่วน... เช่น 'ประชุมสไลด์ #งาน !high @วันนี้ ทุกวัน'"
            class="flex-grow bg-transparent text-sm text-slate-800 dark:text-white placeholder-slate-400 focus:outline-none min-h-[44px]"
            @keydown.enter="handleCreate"
            @keydown.esc="emit('close')"
            :disabled="loading"
          />
          <button
            @click="emit('close')"
            class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 min-h-[44px] min-w-[44px]"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Real-time Parsed Pills Section -->
        <div class="px-5 py-3.5 bg-slate-50 dark:bg-slate-900/40 flex flex-wrap gap-2 items-center">
          <span class="text-[10px] font-semibold text-slate-400 uppercase tracking-wide mr-1 select-none">คุณสมบัติ:</span>
          
          <!-- Workspace Pill -->
          <span
            v-if="activeWorkspace"
            class="inline-flex items-center space-x-1.5 px-2 py-0.5 rounded text-[11px] font-medium"
            :style="{ backgroundColor: activeWorkspace.color + '18', color: activeWorkspace.color }"
          >
            <FolderClosed class="w-3.5 h-3.5" />
            <span>{{ activeWorkspace.name }}</span>
          </span>

          <!-- Priority Pill -->
          <span
            class="inline-flex items-center space-x-1.5 px-2 py-0.5 rounded text-[11px] font-medium border"
            :class="{
              'bg-red-50 border-red-200 text-red-600 dark:bg-red-500/10 dark:border-red-500/20': parsed.priority === 'urgent' || parsed.priority === 'high',
              'bg-amber-50 border-amber-200 text-amber-600 dark:bg-amber-500/10 dark:border-amber-500/20': parsed.priority === 'medium',
              'bg-green-50 border-green-200 text-green-600 dark:bg-green-500/10 dark:border-green-500/20': parsed.priority === 'low'
            }"
          >
            <AlertTriangle class="w-3.5 h-3.5" />
            <span>{{ getPriorityText(parsed.priority) }}</span>
          </span>

          <!-- Due Date Pill -->
          <span
            v-if="parsed.dueDate"
            class="inline-flex items-center space-x-1.5 px-2 py-0.5 rounded text-[11px] font-medium bg-blue-50 border border-blue-200 text-blue-600 dark:bg-blue-500/10 dark:border-blue-500/20"
          >
            <Calendar class="w-3.5 h-3.5" />
            <span>ส่ง: {{ formatDate(parsed.dueDate) }}</span>
          </span>

          <!-- Recurrence Pill -->
          <span
            v-if="parsed.isRecurring"
            class="inline-flex items-center space-x-1.5 px-2 py-0.5 rounded text-[11px] font-medium bg-brand-50 border border-brand-200 text-brand-600 dark:bg-brand-600/10 dark:border-brand-500/20"
          >
            <RefreshCw class="w-3.5 h-3.5" />
            <span>ทำซ้ำ: {{ getRecurText(parsed.recurRule) }}</span>
          </span>
        </div>

        <!-- Shortcut guide footer -->
        <div class="px-5 py-2.5 border-t border-border-light dark:border-border-dark flex items-center justify-between text-[10px] text-slate-400 select-none">
          <div class="flex space-x-3">
            <span><strong>#</strong> เพื่อเลือกพื้นที่งาน</span>
            <span><strong>!</strong> เพื่อเลือกความสำคัญ</span>
            <span><strong>@</strong> เพื่อกำหนดวันส่ง</span>
          </div>
          <div>
            <span>กด <strong>Enter</strong> เพื่อสร้างงาน</span>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>
