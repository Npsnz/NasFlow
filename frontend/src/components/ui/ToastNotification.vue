<script setup lang="ts">
import { useUIStore } from '@/stores/ui'
import { computed } from 'vue'

const uiStore = useUIStore()
const toasts = computed(() => uiStore.toasts)
</script>

<template>
  <div class="fixed top-5 right-5 z-[100] flex flex-col space-y-3 w-full max-w-sm pointer-events-none">
    <TransitionGroup
      enter-active-class="transform ease-out duration-300 transition"
      enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
      enter-to-class="translate-y-0 opacity-100 sm:translate-x-0"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="pointer-events-auto w-full max-w-sm overflow-hidden rounded-lg shadow-lg border p-4 flex items-start space-x-3"
        :class="{
          'bg-white border-green-200 dark:bg-slate-900 dark:border-green-900': toast.type === 'success',
          'bg-white border-red-200 dark:bg-slate-900 dark:border-red-900': toast.type === 'error',
          'bg-white border-blue-200 dark:bg-slate-900 dark:border-blue-900': toast.type === 'info'
        }"
      >
        <!-- Icon -->
        <div class="flex-shrink-0">
          <!-- Success (Green) -->
          <svg
            v-if="toast.type === 'success'"
            class="h-5 w-5 text-green-500"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>

          <!-- Error (Red) -->
          <svg
            v-else-if="toast.type === 'error'"
            class="h-5 w-5 text-red-500"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>

          <!-- Info (Blue) -->
          <svg
            v-else
            class="h-5 w-5 text-blue-500"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>

        <!-- Text -->
        <div class="flex-grow pt-0.5">
          <p
            class="text-sm font-medium"
            :class="{
              'text-slate-900 dark:text-white': true
            }"
          >
            {{ toast.message }}
          </p>
        </div>
      </div>
    </TransitionGroup>
  </div>
</template>
