<script setup lang="ts">
import { onMounted } from 'vue'
import { useUIStore } from '@/stores/ui'
import CalendarComponent from '@/components/calendar/CalendarView.vue'
import TaskDrawer from '@/components/task/TaskDrawer.vue'

const uiStore = useUIStore()

onMounted(() => {
  uiStore.setActiveView('calendar')
})
</script>

<template>
  <div class="h-full flex flex-col overflow-hidden">
    <CalendarComponent />
    <!-- Global Task Drawer -->
    <TaskDrawer />
  </div>
</template>
