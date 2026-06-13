<script setup lang="ts">
import { computed } from 'vue'
import type { Task } from '@/stores/tasks'
import { useUIStore } from '@/stores/ui'
import { Calendar, CheckSquare } from 'lucide-vue-next'

const props = defineProps<{
  task: Task
  isFocused?: boolean
}>()

const uiStore = useUIStore()

const isOverdue = computed(() => {
  if (!props.task.due_date || props.task.status === 'done') return false
  return new Date(props.task.due_date) < new Date()
})

const formattedDueDate = computed(() => {
  if (!props.task.due_date) return ''
  const d = new Date(props.task.due_date)
  return d.toLocaleDateString('th-TH', { day: 'numeric', month: 'short' })
})

const priorityBadgeClass = {
  low: 'badge-low',
  medium: 'badge-medium',
  high: 'badge-high',
  urgent: 'badge-urgent',
}

const priorityLabel = {
  low: 'ต่ำ',
  medium: 'กลาง',
  high: 'สูง',
  urgent: 'ด่วน',
}

const visibleTags = computed(() => props.task.tags?.slice(0, 2) ?? [])
const extraTagCount = computed(() => Math.max(0, (props.task.tags?.length ?? 0) - 2))

const completedSubtasksCount = computed(() => {
  if (!props.task.subtasks) return 0
  return props.task.subtasks.filter(s => s.status === 'done').length
})

const openDetails = () => {
  uiStore.setFocusedTask(props.task.id)
  window.dispatchEvent(new CustomEvent('open-task-details', { detail: props.task.id }))
}
</script>

<template>
  <div
    @click="openDetails"
    class="group relative bg-white/80 dark:bg-slate-800/60 rounded-xl p-3 sm:p-4 cursor-pointer select-none
           border transition-all duration-200
           hover:-translate-y-0.5 hover:shadow-md
           border-l-2 animate-fade-up"
    :class="isFocused
      ? 'border-brand-400 border-l-brand-500 ring-2 ring-brand-500/20 bg-brand-50/30 dark:bg-brand-500/5 shadow-sm'
      : 'border-slate-200 dark:border-slate-700/60 border-l-transparent hover:border-l-brand-400 hover:border-slate-300 dark:hover:border-slate-600'"
  >
    <!-- Row 1: Priority badge + Tags -->
    <div class="flex items-center justify-between gap-2 mb-2 sm:mb-3">
      <!-- Priority pill -->
      <span
        class="inline-flex items-center px-2 sm:px-2.5 py-0.5 sm:py-1 rounded-md text-[10px] sm:text-xs font-bold tracking-wide uppercase"
        :class="[priorityBadgeClass[task.priority], task.priority === 'urgent' ? 'animate-pulse' : '']"
      >
        {{ priorityLabel[task.priority] }}
      </span>

      <!-- Tags + overflow -->
      <div class="flex items-center gap-1 min-w-0">
        <span
          v-for="tag in visibleTags"
          :key="tag.id"
          class="inline-flex items-center px-1.5 sm:px-2 py-0.5 rounded text-[8px] sm:text-[10px] font-semibold truncate max-w-[50px] sm:max-w-[60px]"
          :style="{ backgroundColor: tag.color + '18', color: tag.color, border: `1px solid ${tag.color}28` }"
        >
          {{ tag.name }}
        </span>
        <span
          v-if="extraTagCount > 0"
          class="text-[9px] font-semibold text-slate-400 dark:text-slate-500"
        >+{{ extraTagCount }}</span>
      </div>
    </div>

    <!-- Row 2: Title -->
    <h4 class="text-sm sm:text-base font-semibold text-slate-800 dark:text-slate-100 leading-snug line-clamp-2 break-words mb-1 sm:mb-2
               group-hover:text-brand-600 dark:group-hover:text-brand-400 transition-colors duration-150">
      {{ task.title }}
    </h4>

    <!-- Row 3: Description preview (if any) -->
    <p
      v-if="task.description"
      class="text-xs sm:text-sm text-slate-400 dark:text-slate-500 line-clamp-1 mb-2"
    >
      {{ task.description }}
    </p>

    <!-- Row 4: Footer meta -->
    <div
      v-if="formattedDueDate || (task.subtasks && task.subtasks.length > 0)"
      class="flex items-center justify-between pt-2 sm:pt-2.5 mt-1 sm:mt-1.5 border-t border-slate-100 dark:border-slate-700/50"
    >
      <!-- Due date -->
      <span
        v-if="formattedDueDate"
        class="inline-flex items-center gap-1 sm:gap-1.5 text-[10px] sm:text-xs font-medium rounded px-1.5 sm:px-2 py-0.5 sm:py-1"
        :class="isOverdue
          ? 'bg-red-50 text-red-600 dark:bg-red-950/30 dark:text-red-400'
          : 'text-slate-400 dark:text-slate-500'"
      >
        <Calendar class="w-3 h-3 sm:w-4 sm:h-4" />
        {{ formattedDueDate }}
      </span>
      <span v-else />

      <!-- Subtasks progress -->
      <span
        v-if="task.subtasks && task.subtasks.length > 0"
        class="inline-flex items-center gap-1 sm:gap-1.5 text-[10px] sm:text-xs text-slate-400 dark:text-slate-500"
      >
        <CheckSquare class="w-3 h-3 sm:w-4 sm:h-4" />
        {{ completedSubtasksCount }}/{{ task.subtasks.length }}
      </span>
    </div>
  </div>
</template>
