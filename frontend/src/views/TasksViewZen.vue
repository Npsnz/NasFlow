<template>
  <div class="tasks-zen" :class="{ 'dark-mode': themeStore.isDark }">
    <TopNav />

    <main class="tasks-content">
      <!-- Header -->
      <section class="tasks-header">
        <div class="header-inner">
          <h1>My Tasks</h1>
          <button class="btn-new-task" @click="showModal = true">+ New Task</button>
        </div>

        <!-- Tag Filters -->
        <div v-if="tasksStore.tags.length > 0" class="tag-filters">
          <button
            class="tag-filter"
            :class="{ active: selectedTags.length === 0 }"
            @click="selectedTags = []"
          >
            All Tags
          </button>
          <button
            v-for="tag in tasksStore.tags"
            :key="tag.id"
            class="tag-filter"
            :class="{ active: selectedTags.includes(tag.id) }"
            @click="toggleTag(tag.id)"
            :style="{ '--tag-color': tag.color }"
          >
            {{ tag.name }}
          </button>
        </div>
      </section>

      <!-- Kanban View -->
      <section class="task-board" @touchstart="handleTouchStart" @touchend="handleTouchEnd">
        <div v-for="status in ['todo', 'doing', 'done']" :key="status" class="task-column" @dragover.prevent @drop="handleDrop($event, status)">
          <h2 class="column-title">{{ status === 'todo' ? 'To Do' : status === 'doing' ? 'In Progress' : 'Completed' }}</h2>
          <div class="task-cards" @dragover.prevent @drop.stop="handleDrop($event, status)">
            <div v-if="getFilteredTasks(status).length === 0" class="empty-state">
              No tasks yet
            </div>
            <div
              v-for="task in getFilteredTasks(status)"
              :key="task.id"
              class="task-card"
              @click="openTaskDetail(task)"
              draggable="true"
              @dragstart="handleDragStart($event, task, status)"
              @dragend="handleDragEnd"
            >
              <div class="card-header">
                <span class="card-title">{{ task.title }}</span>
                <span class="priority-badge" :class="task.priority">{{ task.priority }}</span>
              </div>
              <p v-if="task.description" class="card-desc">{{ task.description }}</p>
              <div class="card-footer">
                <div v-if="task.tags?.length" class="tags-row">
                  <span v-for="tag in task.tags" :key="tag.id" class="tag-pill" :style="{ backgroundColor: tag.color + '20', color: tag.color }">
                    {{ tag.name }}
                  </span>
                </div>
                <span v-if="task.due_date" class="due-date">{{ formatDate(task.due_date) }}</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Task Detail Modal -->
      <div v-if="selectedTask" class="modal-overlay" @click.self="selectedTask = null">
        <div class="modal-content task-detail">
          <div class="modal-header">
            <h2>{{ selectedTask.title }}</h2>
            <button class="close-btn" @click="selectedTask = null">×</button>
          </div>

          <div class="task-detail-content">
            <div class="detail-section">
              <label>Title</label>
              <input v-model="selectedTask.title" type="text" />
            </div>

            <div class="detail-row">
              <div class="detail-section">
                <label>Status</label>
                <select v-model="selectedTask.status">
                  <option value="todo">To Do</option>
                  <option value="doing">In Progress</option>
                  <option value="done">Completed</option>
                </select>
              </div>

              <div class="detail-section">
                <label>Priority</label>
                <select v-model="selectedTask.priority">
                  <option value="low">Low</option>
                  <option value="medium">Medium</option>
                  <option value="high">High</option>
                  <option value="urgent">Urgent</option>
                </select>
              </div>
            </div>

            <div class="detail-section">
              <label>Description</label>
              <textarea v-model="selectedTask.description"></textarea>
            </div>

            <div class="detail-section">
              <label>Due Date</label>
              <input v-model="selectedTask.due_date" type="date" />
            </div>

            <div class="detail-section">
              <label>Tags</label>
              <div class="tag-selector">
                <button
                  v-for="tag in tasksStore.tags"
                  :key="tag.id"
                  class="tag-btn"
                  :class="{ active: selectedTask.tags?.some((t: any) => t.id === tag.id) }"
                  @click="toggleTaskTag(tag)"
                  :style="{ '--tag-color': tag.color }"
                >
                  {{ tag.name }}
                </button>
              </div>
            </div>

            <div class="detail-actions">
              <button class="btn-delete" @click="deleteTaskConfirm">Delete Task</button>
              <button class="btn-save-detail" @click="updateTaskAndClose">Save</button>
            </div>
          </div>
        </div>
      </div>

      <!-- New Task Modal -->
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal-content">
          <div class="modal-header">
            <h2>New Task</h2>
            <button class="close-btn" @click="closeModal">×</button>
          </div>

          <form @submit.prevent="handleCreateTask" class="task-form">
            <div class="form-group">
              <label>Title</label>
              <input v-model="newTask.title" type="text" placeholder="Task title" required />
            </div>

            <div class="form-group">
              <label>Description</label>
              <textarea v-model="newTask.description" placeholder="Task description..."></textarea>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label>Priority</label>
                <select v-model="newTask.priority">
                  <option value="low">Low</option>
                  <option value="medium" selected>Medium</option>
                  <option value="high">High</option>
                  <option value="urgent">Urgent</option>
                </select>
              </div>

              <div class="form-group">
                <label>Due Date</label>
                <input v-model="newTask.due_date" type="date" />
              </div>
            </div>

            <div class="form-group">
              <label>Tags</label>
              <div class="tag-selector-modal">
                <button
                  v-for="tag in tasksStore.tags"
                  :key="tag.id"
                  type="button"
                  class="tag-btn-modal"
                  :class="{ active: selectedTaskTags.includes(tag.id) }"
                  @click="toggleNewTaskTag(tag.id)"
                  :style="{ '--tag-color': tag.color }"
                >
                  {{ tag.name }}
                </button>
              </div>
            </div>

            <div class="form-actions">
              <button type="button" class="btn-cancel" @click="closeModal">Cancel</button>
              <button type="submit" class="btn-create" :disabled="creating">
                {{ creating ? 'Creating...' : 'Create Task' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import TopNav from '@/components/TopNav.vue'
import { useTasksStore } from '@/stores/tasks'
import { useWorkspaceStore } from '@/stores/workspace'
import { useToastStore } from '@/stores/toast'
import { useThemeStore } from '@/stores/theme'

const tasksStore = useTasksStore()
const workspaceStore = useWorkspaceStore()
const toastStore = useToastStore()
const themeStore = useThemeStore()

// Sync theme changes
watch(() => themeStore.isDark, (isDark) => {
  document.documentElement.setAttribute('data-theme', isDark ? 'dark' : 'light')
})
const showModal = ref(false)
const creating = ref(false)

const selectedTags = ref<number[]>([])
const selectedTaskTags = ref<number[]>([])
const selectedTask = ref<any>(null)

// Swipe handling for mobile
const touchStart = ref({ x: 0, y: 0 })
const boardScroll = ref({ element: null as HTMLElement | null, scrollLeft: 0 })
const lastSwipeTime = ref(0) // Prevent multiple swipes in quick succession

const newTask = ref({
  title: '',
  description: '',
  priority: 'medium' as 'low' | 'medium' | 'high' | 'urgent',
  status: 'todo' as 'todo' | 'doing' | 'done',
  due_date: '',
  workspace_id: 1, // Default workspace
})

onMounted(async () => {
  try {
    // Get active workspace (first workspace or user's selected)
    if (workspaceStore.workspaces.length === 0) {
      await workspaceStore.fetchWorkspaces()
    }

    const activeWs = workspaceStore.currentWorkspace || workspaceStore.workspaces[0]
    if (activeWs) {
      newTask.value.workspace_id = activeWs.id
      // Fetch tags first so they're available when tasks load
      await tasksStore.fetchTags()
      await tasksStore.fetchTasks({ workspace_id: activeWs.id })
    }
  } catch (err) {
    console.error('Failed to load tasks:', err)
  }
})

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

const toggleTag = (tagId: number) => {
  const idx = selectedTags.value.indexOf(tagId)
  if (idx > -1) {
    selectedTags.value.splice(idx, 1)
  } else {
    selectedTags.value.push(tagId)
  }
}

const getFilteredTasks = (status: string) => {
  let filtered = tasksStore.tasks.filter(t => t.status === status)

  // Filter by selected tags
  if (selectedTags.value.length > 0) {
    filtered = filtered.filter(task =>
      task.tags?.some(tag => selectedTags.value.includes(tag.id))
    )
  }

  return filtered
}

const openTaskDetail = (task: any) => {
  selectedTask.value = { ...task }

  // Format due_date for input type="date" (YYYY-MM-DD)
  if (selectedTask.value.due_date) {
    const date = new Date(selectedTask.value.due_date)
    selectedTask.value.due_date = date.toISOString().split('T')[0]
  }
}

const toggleTaskTag = (tag: any) => {
  if (!selectedTask.value.tags) selectedTask.value.tags = []

  const idx = selectedTask.value.tags.findIndex((t: any) => t.id === tag.id)
  if (idx > -1) {
    selectedTask.value.tags.splice(idx, 1)
  } else {
    selectedTask.value.tags.push(tag)
  }
}

const updateTaskAndClose = async () => {
  if (!selectedTask.value) return

  try {
    let dueDate = selectedTask.value.due_date
    if (dueDate && typeof dueDate === 'string') {
      dueDate = new Date(dueDate).toISOString()
    }

    await tasksStore.updateTask(selectedTask.value.id, {
      title: selectedTask.value.title,
      description: selectedTask.value.description,
      status: selectedTask.value.status,
      priority: selectedTask.value.priority,
      due_date: dueDate || null,
      tag_ids: selectedTask.value.tags?.map((t: any) => t.id) || [],
    })
    toastStore.success('Task updated!')
    selectedTask.value = null
  } catch (err: any) {
    toastStore.error(err.response?.data?.error || 'Failed to update task')
  }
}

const deleteTaskConfirm = async () => {
  if (!selectedTask.value) return
  if (!confirm('Delete this task?')) return

  try {
    await tasksStore.deleteTask(selectedTask.value.id)
    toastStore.success('Task deleted!')
    selectedTask.value = null
  } catch (err: any) {
    toastStore.error(err.response?.data?.error || 'Failed to delete task')
  }
}

const toggleNewTaskTag = (tagId: number) => {
  const idx = selectedTaskTags.value.indexOf(tagId)
  if (idx > -1) {
    selectedTaskTags.value.splice(idx, 1)
  } else {
    selectedTaskTags.value.push(tagId)
  }
}

const handleCreateTask = async () => {
  if (!newTask.value.title.trim()) {
    toastStore.warning('Title is required')
    return
  }

  creating.value = true
  try {
    await tasksStore.createTask({
      title: newTask.value.title,
      description: newTask.value.description,
      priority: newTask.value.priority,
      status: newTask.value.status,
      due_date: newTask.value.due_date || null,
      workspace_id: newTask.value.workspace_id,
      tag_ids: selectedTaskTags.value,
    })

    toastStore.success('Task created!')
    closeModal()
  } catch (err: any) {
    toastStore.error(err.response?.data?.error || 'Failed to create task')
  } finally {
    creating.value = false
  }
}

const closeModal = () => {
  showModal.value = false
  selectedTaskTags.value = []
  newTask.value = {
    title: '',
    description: '',
    priority: 'medium',
    status: 'todo',
    due_date: '',
    workspace_id: newTask.value.workspace_id,
  }
}

let draggedTask: any = null
let draggedFromStatus: string = ''

const handleDragStart = (e: DragEvent, task: any, status: string) => {
  draggedTask = task
  draggedFromStatus = status
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/html', '')
  }
}

const handleDragEnd = () => {
  draggedTask = null
  draggedFromStatus = ''
}

const handleDrop = async (e: DragEvent, targetStatus: string) => {
  e.preventDefault()
  if (!draggedTask) return

  // If dropped in same column, no need to update status
  if (draggedFromStatus === targetStatus) {
    return
  }

  // Update task status to new column
  try {
    await tasksStore.updateTaskStatus(draggedTask.id, targetStatus)
    toastStore.success(`Task moved to ${targetStatus === 'todo' ? 'To Do' : targetStatus === 'doing' ? 'In Progress' : 'Completed'}`)
  } catch (err: any) {
    toastStore.error('Failed to move task')
  }

  draggedTask = null
  draggedFromStatus = ''
}

const handleTouchStart = (e: TouchEvent) => {
  const target = e.currentTarget as HTMLElement
  touchStart.value = {
    x: e.touches[0].clientX,
    y: e.touches[0].clientY,
  }
  boardScroll.value.element = target
  boardScroll.value.scrollLeft = target.scrollLeft
}

const handleTouchEnd = (e: TouchEvent) => {
  if (!boardScroll.value.element) return

  const now = Date.now()
  // Prevent multiple swipes within 600ms to avoid double-column jumps
  if (now - lastSwipeTime.value < 600) return

  const touchEnd = {
    x: e.changedTouches[0].clientX,
    y: e.changedTouches[0].clientY,
  }

  const deltaX = touchStart.value.x - touchEnd.x
  const deltaY = Math.abs(touchStart.value.y - touchEnd.y)

  // Only swipe if horizontal movement is greater than vertical
  // Any swipe moves exactly 1 column
  if (Math.abs(deltaX) > 30 && deltaY < 50) {
    lastSwipeTime.value = now
    const scrollAmount = 320 // Approximate column width + gap
    let newScrollLeft = boardScroll.value.scrollLeft

    if (deltaX > 0) {
      // Swipe left - scroll right to next column only
      newScrollLeft += scrollAmount
    } else {
      // Swipe right - scroll left to previous column only
      newScrollLeft -= scrollAmount
    }

    boardScroll.value.element?.scrollTo({
      left: newScrollLeft,
      behavior: 'smooth',
    })
  }
}
</script>

<style scoped>
.tasks-zen {
  min-height: 100vh;
  background: #FAF8F3;
  padding-top: 80px;
  transition: background 0.3s ease;
}

.tasks-zen.dark-mode {
  background: #1A1410;
}

.tasks-content {
  max-width: 1100px;
  margin: 0 auto;
  padding: 3rem 2rem;
}

/* Header */
.tasks-header {
  margin-bottom: 3rem;
}

.header-inner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.tasks-zen h1 {
  margin: 0;
  font-size: 2rem;
  font-family: Georgia, serif;
  font-weight: 300;
  color: #3D3D3D;
}

.tasks-zen.dark-mode h1 {
  color: #F5F1E8;
}

.btn-new-task {
  padding: 0.75rem 1.5rem;
  background: #A89968;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-new-task:hover {
  background: #8B7355;
}

.tasks-zen.dark-mode .btn-new-task {
  background: #D4A574;
  color: #1A1410;
}

.tasks-zen.dark-mode .btn-new-task:hover {
  background: #F5DEB3;
}

/* Filters */
.filters {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.filter-btn {
  padding: 0.5rem 1rem;
  background: white;
  color: #8B8B8B;
  border: 1px solid #E8DDD2;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  font-size: 0.9rem;
  transition: all 0.3s ease;
}

.tasks-zen.dark-mode .filter-btn {
  background: #2D2420;
  color: #C9C1B8;
  border-color: #3D3530;
}

.filter-btn.active {
  background: #A89968;
  color: white;
  border-color: #A89968;
}

.tasks-zen.dark-mode .filter-btn.active {
  background: #D4A574;
  color: #1A1410;
  border-color: #D4A574;
}

/* Tag filters */
.tag-filters {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
  margin-top: 1rem;
}

.tag-filter {
  padding: 0.5rem 1rem;
  border: 1px solid #E8DDD2;
  border-radius: 20px;
  background: white;
  color: #3D3D3D;
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.tasks-zen.dark-mode .tag-filter {
  background: #2D2420;
  border-color: #3D3530;
  color: #C9C1B8;
}

.tag-filter.active {
  background: var(--tag-color, #A89968);
  color: white;
  border-color: var(--tag-color, #A89968);
}

.tasks-zen.dark-mode .tag-filter.active {
  background: var(--tag-color, #D4A574);
  border-color: var(--tag-color, #D4A574);
  color: #1A1410;
}

/* Tag pills on cards */
.tags-row {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.tag-pill {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
  white-space: nowrap;
}

/* View stub */
.view-stub {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

.stub-content {
  text-align: center;
  color: #8B8B8B;
  font-size: 1.1rem;
}

.tasks-zen.dark-mode .stub-content {
  color: #7A7270;
}

/* Board */
.task-board {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 2rem;
}

.task-column {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  min-height: 500px;
  padding: 1rem;
  border-radius: 8px;
}

.column-title {
  margin: 0 0 1rem;
  font-size: 1.1rem;
  font-weight: 500;
  color: #3D3D3D;
  padding-bottom: 0.75rem;
  border-bottom: 2px solid #E8DDD2;
}

.tasks-zen.dark-mode .column-title {
  color: #F5F1E8;
  border-bottom-color: #3D3530;
}

.task-cards {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.task-card {
  padding: 1.5rem;
  background: white;
  border-radius: 10px;
  border: 1px solid #E8DDD2;
  transition: all 0.3s ease;
  cursor: move;
  user-select: none;
}

.tasks-zen.dark-mode .task-card {
  background: #2D2420;
  border-color: #3D3530;
}

.task-card:hover {
  border-color: #A89968;
  box-shadow: 0 4px 12px rgba(168, 153, 104, 0.1);
}

.task-card:active {
  opacity: 0.7;
  transform: scale(0.98);
}

.tasks-zen.dark-mode .task-card:hover {
  border-color: #F5DEB3;
  box-shadow: 0 4px 12px rgba(245, 222, 179, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 0.75rem;
}

.card-title {
  font-weight: 600;
  color: #3D3D3D;
  font-size: 1rem;
}

.tasks-zen.dark-mode .card-title {
  color: #F5F1E8;
}

.priority-badge {
  padding: 0.25rem 0.6rem;
  background: rgba(168, 153, 104, 0.1);
  color: #A89968;
  font-size: 0.75rem;
  border-radius: 4px;
  font-weight: 500;
}

.tasks-zen.dark-mode .priority-badge {
  background: rgba(212, 165, 116, 0.15);
  color: #F5DEB3;
}

.card-desc {
  margin: 0 0 1rem;
  font-size: 0.9rem;
  color: #8B8B8B;
  line-height: 1.4;
}

.tasks-zen.dark-mode .card-desc {
  color: #C9C1B8;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.75rem;
  font-size: 0.8rem;
}

.tag {
  padding: 0.25rem 0.6rem;
  background: rgba(212, 165, 116, 0.2);
  color: #8B7355;
  border-radius: 4px;
}

.tasks-zen.dark-mode .tag {
  background: rgba(212, 165, 116, 0.25);
  color: #D4A574;
}

.due-date {
  color: #8B8B8B;
  font-weight: 500;
  font-size: 0.85rem;
}

.tasks-zen.dark-mode .due-date {
  color: #C9C1B8;
}

/* Priority badges */
.priority-badge {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.priority-badge.low {
  background: rgba(76, 125, 80, 0.15);
  color: #4C7D50;
}

.priority-badge.medium {
  background: rgba(184, 113, 60, 0.15);
  color: #B8713C;
}

.priority-badge.high {
  background: rgba(139, 69, 61, 0.15);
  color: #8B453D;
}

.priority-badge.urgent {
  background: rgba(139, 69, 61, 0.25);
  color: #8B453D;
  font-weight: 700;
}

/* Empty state */
.empty-state {
  padding: 2rem;
  text-align: center;
  color: #8B8B8B;
  font-style: italic;
}

.tasks-zen.dark-mode .empty-state {
  color: #7A7270;
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.modal-content {
  background: white;
  border-radius: 12px;
  border: 1px solid #E8DDD2;
  padding: 2rem;
  width: 90%;
  max-width: 500px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.tasks-zen.dark-mode .modal-content {
  background: #2D2420;
  border-color: #3D3530;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.5rem;
  font-family: Georgia, serif;
  color: #3D3D3D;
}

.tasks-zen.dark-mode .modal-header h2 {
  color: #F5F1E8;
}

.modal-actions {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.btn-save {
  padding: 0.5rem 1rem;
  background: #A89968;
  color: white;
  border: none;
  border-radius: 6px;
  font-weight: 500;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s ease;
}

.btn-save:hover {
  background: #8B7355;
}

.tasks-zen.dark-mode .btn-save {
  background: #D4A574;
  color: #1A1410;
}

.tasks-zen.dark-mode .btn-save:hover {
  background: #F5DEB3;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #8B8B8B;
  transition: color 0.2s;
}

.close-btn:hover {
  color: #3D3D3D;
}

.tasks-zen.dark-mode .close-btn:hover {
  color: #F5F1E8;
}

/* Task form */
.task-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-weight: 500;
  font-size: 0.9rem;
  color: #3D3D3D;
}

.tasks-zen.dark-mode .form-group label {
  color: #F5F1E8;
}

.form-group input,
.form-group textarea,
.form-group select {
  padding: 0.75rem;
  border: 1px solid #E8DDD2;
  border-radius: 6px;
  font-family: inherit;
  font-size: 0.95rem;
  background: white;
  color: #3D3D3D;
  transition: all 0.2s;
}

.tasks-zen.dark-mode .form-group input,
.tasks-zen.dark-mode .form-group textarea,
.tasks-zen.dark-mode .form-group select {
  background: #1A1410;
  border-color: #3D3530;
  color: #F5F1E8;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: #A89968;
  box-shadow: 0 0 0 3px rgba(168, 153, 104, 0.1);
}

.form-group textarea {
  resize: vertical;
  min-height: 100px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

/* Task detail modal */
.modal-content.task-detail {
  max-width: 600px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

.task-detail-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  max-height: calc(80vh - 100px);
  overflow-y: auto;
  padding-right: 0.5rem;
}

.task-detail-content::-webkit-scrollbar {
  width: 6px;
}

.task-detail-content::-webkit-scrollbar-track {
  background: transparent;
}

.task-detail-content::-webkit-scrollbar-thumb {
  background: #D9CEC0;
  border-radius: 3px;
}

.tasks-zen.dark-mode .task-detail-content::-webkit-scrollbar-thumb {
  background: #5C5450;
}

.detail-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.detail-section {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.detail-section label {
  font-weight: 600;
  font-size: 0.9rem;
  color: #3D3D3D;
}

.tasks-zen.dark-mode .detail-section label {
  color: #F5F1E8;
}

.detail-section select,
.detail-section textarea,
.detail-section input {
  padding: 0.75rem;
  border: 1px solid #E8DDD2;
  border-radius: 6px;
  font-family: inherit;
  font-size: 0.95rem;
  background: white;
  color: #3D3D3D;
}

.detail-section select {
  padding-right: 2.5rem;
  appearance: none;
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%233D3D3D' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right 0.75rem center;
  background-size: 1.5em 1.5em;
}

.tasks-zen.dark-mode .detail-section select {
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23F5F1E8' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
}

.tasks-zen.dark-mode .detail-section select,
.tasks-zen.dark-mode .detail-section textarea,
.tasks-zen.dark-mode .detail-section input {
  background: #1A1410;
  border-color: #3D3530;
  color: #F5F1E8;
}

.detail-section textarea {
  min-height: 120px;
  resize: vertical;
}

.tag-selector,
.tag-selector-modal {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.tag-btn,
.tag-btn-modal {
  padding: 0.5rem 1rem;
  border: 1px solid #E8DDD2;
  border-radius: 20px;
  background: white;
  color: #3D3D3D;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.2s;
}

.tasks-zen.dark-mode .tag-btn,
.tasks-zen.dark-mode .tag-btn-modal {
  background: #2D2420;
  border-color: #3D3530;
  color: #C9C1B8;
}

.tag-btn.active,
.tag-btn-modal.active {
  background: var(--tag-color, #A89968);
  color: white;
  border-color: var(--tag-color, #A89968);
  font-weight: 600;
}

.tasks-zen.dark-mode .tag-btn.active,
.tasks-zen.dark-mode .tag-btn-modal.active {
  background: var(--tag-color, #A89968);
  color: white;
  border-color: var(--tag-color, #A89968);
  box-shadow: 0 0 0 2px rgba(168, 153, 104, 0.3);
}

.detail-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
  border-top: 1px solid #E8DDD2;
  padding-top: 1rem;
}

.tasks-zen.dark-mode .detail-actions {
  border-top-color: #3D3530;
}

.btn-delete {
  padding: 0.75rem 1.5rem;
  background: rgba(139, 69, 61, 0.1);
  color: #8B453D;
  border: 1px solid #E8DDD2;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.tasks-zen.dark-mode .btn-delete {
  border-color: #3D3530;
  background: rgba(200, 120, 100, 0.1);
  color: #C87864;
}

.btn-delete:hover {
  background: rgba(139, 69, 61, 0.2);
}

.btn-save-detail {
  flex: 1;
  padding: 0.75rem 1.5rem;
  background: #A89968;
  color: white;
  border: none;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-save-detail:hover {
  background: #8B7355;
}

.tasks-zen.dark-mode .btn-save-detail {
  background: #D4A574;
  color: #1A1410;
}

.tasks-zen.dark-mode .btn-save-detail:hover {
  background: #F5DEB3;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.btn-cancel,
.btn-create {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  flex: 1;
}

.btn-cancel {
  background: #E8DDD2;
  color: #3D3D3D;
}

.btn-cancel:hover {
  background: #D9CEC0;
}

.tasks-zen.dark-mode .btn-cancel {
  background: #3D3530;
  color: #F5F1E8;
}

.tasks-zen.dark-mode .btn-cancel:hover {
  background: #4D4540;
}

.btn-create {
  background: linear-gradient(135deg, #A89968 0%, #8B7355 100%);
  color: white;
}

.btn-create:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(168, 153, 104, 0.3);
}

.btn-create:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.tasks-zen.dark-mode .btn-create {
  background: linear-gradient(135deg, #D4A574 0%, #A89968 100%);
  color: #1A1410;
}

@media (max-width: 768px) {
  .tasks-content {
    padding: 1.5rem 1rem;
  }

  .task-board {
    display: flex;
    gap: 2rem;
    overflow-x: auto;
    overflow-y: hidden;
    -webkit-overflow-scrolling: touch;
    scroll-behavior: smooth;
    scrollbar-width: none;
    -ms-overflow-style: none;
    padding: 0 calc((100% - 320px) / 2);
    scroll-snap-type: x mandatory;
  }

  .task-board::-webkit-scrollbar {
    display: none;
  }

  .task-column {
    flex: 0 0 100%;
    max-width: 320px;
    scroll-snap-align: center;
  }

  .modal-content {
    width: 95%;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .detail-row {
    grid-template-columns: 1fr;
    gap: 0.75rem;
  }

  .modal-content.task-detail {
    max-height: 85vh;
  }

  .task-detail-content {
    max-height: calc(75vh - 100px);
  }
}
</style>
