<script setup lang="ts">
import { ref, nextTick } from 'vue'
import type { Task } from '@/stores/tasks'
import { useTasksStore } from '@/stores/tasks'
import { useWorkspaceStore } from '@/stores/workspace'
import { useUIStore } from '@/stores/ui'
import { Plus } from 'lucide-vue-next'

const props = defineProps<{
  status: 'todo' | 'doing' | 'done'
  title: string
  dotColor: string
  tasks: Task[]
}>()

const tasksStore = useTasksStore()
const workspaceStore = useWorkspaceStore()
const uiStore = useUIStore()

const isAdding = ref(false)
const taskTitle = ref('')
const inlineInputRef = ref<HTMLInputElement | null>(null)

const toggleAdding = () => {
  isAdding.value = !isAdding.value
  if (isAdding.value) {
    nextTick(() => {
      inlineInputRef.value?.focus()
    })
  } else {
    taskTitle.value = ''
  }
}

const handleInlineCreate = async () => {
  if (!taskTitle.value.trim()) {
    toggleAdding()
    return
  }

  const workspaceId = workspaceStore.currentWorkspace?.id || 0
  if (!workspaceId) {
    uiStore.showToast('กรุณาสร้างหรือเลือกพื้นที่ทำงานก่อน', 'error')
    return
  }

  try {
    await tasksStore.createTask({
      title: taskTitle.value.trim(),
      status: props.status,
      workspace_id: workspaceId,
      priority: 'medium'
    })
    uiStore.showToast('สร้างงานสำเร็จ', 'success')
    taskTitle.value = ''
    toggleAdding()
  } catch (err: any) {
    uiStore.showToast(err || 'ไม่สามารถสร้างงานได้', 'error')
  }
}
</script>

<template>
  <div
    class="flex-shrink-0 w-full sm:w-72 lg:w-80 rounded-xl p-3 sm:p-4 flex flex-col max-h-[calc(100vh-190px)]"
  >
    <!-- Column Title + Count -->
    <div class="flex items-center justify-between mb-2 sm:mb-3 px-0.5 select-none">
      <div class="flex items-center space-x-2 sm:space-x-3">
        <!-- Status badge pill -->
        <span
          class="inline-flex items-center px-2.5 sm:px-3 py-1 sm:py-1.5 rounded-lg text-xs sm:text-sm font-semibold"
          :class="status === 'todo'
            ? 'bg-amber-100 text-amber-700 dark:bg-amber-950/50 dark:text-amber-400'
            : status === 'doing'
            ? 'bg-blue-100 text-blue-700 dark:bg-blue-950/50 dark:text-blue-400'
            : 'bg-green-100 text-green-700 dark:bg-green-950/50 dark:text-green-400'"
        >
          {{ title }}
        </span>
        <!-- Count badge -->
        <span
          class="inline-flex items-center px-2 sm:px-2.5 py-1 sm:py-1.5 rounded-full text-xs sm:text-sm font-bold bg-white dark:bg-slate-700 border shadow-sm"
          :class="status === 'todo'
            ? 'border-amber-200 text-amber-600 dark:border-amber-800'
            : status === 'doing'
            ? 'border-blue-200 text-blue-600 dark:border-blue-800'
            : 'border-green-200 text-green-600 dark:border-green-800'"
        >
          {{ tasks.length }}
        </span>
      </div>
    </div>

    <!-- Add Button (Top) -->
    <div class="mb-2 sm:mb-3">
      <!-- Input Mode -->
      <div v-if="isAdding" class="bg-white dark:bg-surface-dark border-[0.5px] border-slate-200 dark:border-border-dark rounded-lg p-2 sm:p-3 shadow-sm space-y-2 sm:space-y-3">
        <input
          ref="inlineInputRef"
          type="text"
          v-model="taskTitle"
          placeholder="พิมพ์ชื่องาน..."
          class="w-full text-xs sm:text-sm bg-slate-50 dark:bg-slate-900 border border-border-light dark:border-border-dark rounded px-2 sm:px-3 py-2 sm:py-2.5 focus:outline-none min-h-[36px] sm:min-h-[40px]"
          @keydown.enter="handleInlineCreate"
          @keydown.esc="toggleAdding"
        />
        <div class="flex justify-end gap-1 sm:gap-2">
          <button
            @click="toggleAdding"
            class="px-2 sm:px-3 py-1 sm:py-1.5 text-slate-400 hover:text-slate-600 text-xs sm:text-sm rounded min-h-[32px] sm:min-h-[36px]"
          >
            ยกเลิก
          </button>
          <button
            @click="handleInlineCreate"
            class="px-2 sm:px-3 py-1 sm:py-1.5 bg-brand-500 text-white font-semibold rounded text-xs sm:text-sm hover:bg-brand-600 min-h-[32px] sm:min-h-[36px]"
          >
            เพิ่ม
          </button>
        </div>
      </div>

      <!-- Add Button Mode -->
      <button
        v-else
        @click="toggleAdding"
        class="w-full py-2 sm:py-2.5 border border-dashed border-slate-300 dark:border-border-dark hover:border-slate-400 hover:bg-slate-100/40 dark:hover:bg-slate-800/20 text-slate-500 dark:text-slate-400 hover:text-slate-700 rounded-lg text-xs sm:text-sm font-medium flex items-center justify-center gap-1.5 sm:gap-2 transition-all min-h-[40px] sm:min-h-[44px]"
      >
        <Plus class="w-3.5 h-3.5 sm:w-4 sm:h-4" />
        <span class="hidden sm:inline">เพิ่มงาน</span>
        <span class="sm:hidden">เพิ่ม</span>
      </button>
    </div>

    <!-- Cards Scroll Area / Drag Target -->
    <div
      class="flex-grow overflow-y-auto space-y-3 pb-4 scrollbar-thin drag-zone"
      :data-status="status"
    >
      <slot />
    </div>
  </div>
</template>
