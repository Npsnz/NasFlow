<script setup lang="ts">
import { onMounted } from 'vue'
import { useUIStore } from '@/stores/ui'
import TaskListComponent from '@/components/task/TaskList.vue'
import TaskDrawer from '@/components/task/TaskDrawer.vue'

const uiStore = useUIStore()

onMounted(() => {
  uiStore.setActiveView('tasks')
})
</script>

<template>
  <div class="h-full flex flex-col overflow-hidden">
    <TaskListComponent />
    <!-- Global Task Drawer -->
    <TaskDrawer />
  </div>
</template>
