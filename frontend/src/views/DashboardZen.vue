<template>
  <div class="dashboard-zen" :class="{ 'dark-mode': themeStore.isDark }">
    <!-- Top Nav -->
    <TopNav />

    <!-- Main content (offset for fixed nav) -->
    <main class="dashboard-content">
      <!-- Welcome section with depth layers -->
      <section class="welcome-section" ref="welcomeSection">
        <div class="welcome-bg-glow" data-depth="0"></div>
        <div class="welcome-inner" data-depth="4">
          <h1 class="welcome-title anim-title">{{ getGreeting }}</h1>
          <p class="welcome-subtitle">{{ $t('app.tagline') }}</p>
        </div>
      </section>

      <!-- Stats cards (Feature 1: task statistics) -->
      <section class="stats-section" ref="statsSection">
        <div v-for="(stat, idx) in statsList" :key="idx" class="stat-card anim-stat-card" :style="{ '--index': idx }">
          <div class="stat-value">{{ stat.value }}</div>
          <div class="stat-label">{{ stat.label }}</div>
        </div>
      </section>

      <!-- Quick actions (Feature 2) -->
      <section class="actions-section">
        <button class="action-btn primary anim-action" @click="handleNewTask">
          <span>+ New Task</span>
        </button>
        <button class="action-btn secondary anim-action" @click="handleViewAll">
          <span>View All</span>
        </button>
      </section>

      <!-- Recent tasks (Feature 3: recent activity) -->
      <section class="recent-section" ref="recentSection" v-if="recentTasks.length > 0">
        <h2 class="section-title anim-title">Recent Tasks</h2>
        <div class="task-list">
          <div v-for="(task, idx) in recentTasks" :key="task.id" class="task-item anim-task-item" :style="{ '--index': idx }" @click="goToTask()">
            <div class="task-status" :class="`status-${task.status}`">
              <span>{{ getStatusIcon(task.status) }}</span>
            </div>
            <div class="task-info">
              <p class="task-name">{{ task.title }}</p>
              <p class="task-desc">{{ task.description || 'No description' }}</p>
              <div class="task-meta-bottom">
                <div v-if="task.tags?.length" class="task-tags">
                  <span v-for="tag in task.tags" :key="tag.id" class="task-tag" :style="{ backgroundColor: tag.color + '20', color: tag.color }">
                    {{ tag.name }}
                  </span>
                </div>
                <span v-if="task.due_date" class="task-due-date">{{ formatDate(task.due_date) }}</span>
              </div>
            </div>
            <div class="task-priority" :class="`priority-${task.priority}`">{{ task.priority }}</div>
          </div>
        </div>
      </section>
      <section v-else class="recent-section empty">
        <p class="empty-message">No tasks yet. Create one to get started!</p>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import TopNav from '@/components/TopNav.vue'
import { useTasksStore } from '@/stores/tasks'
import { useWorkspaceStore } from '@/stores/workspace'
import { useThemeStore } from '@/stores/theme'

const router = useRouter()
const tasksStore = useTasksStore()
const workspaceStore = useWorkspaceStore()
const themeStore = useThemeStore()

// Sync theme changes
watch(() => themeStore.isDark, (isDark) => {
  document.documentElement.setAttribute('data-theme', isDark ? 'dark' : 'light')
})

const statsList = computed(() => [
  { value: tasksStore.stats.total, label: 'Total Tasks' },
  { value: tasksStore.stats.todo, label: 'To Do' },
  { value: tasksStore.stats.doing, label: 'In Progress' },
  { value: tasksStore.stats.done_today, label: 'Done Today' },
])

const recentTasks = computed(() => {
  return tasksStore.tasks
    .filter(t => t.workspace_id === workspaceStore.currentWorkspace?.id)
    .sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
    .slice(0, 5)
})

const getGreeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return 'Good morning'
  if (hour < 18) return 'Good afternoon'
  return 'Good evening'
})

const handleNewTask = () => {
  router.push('/tasks')
}

const handleViewAll = () => {
  router.push('/tasks')
}

const getStatusIcon = (status: string): string => {
  switch (status) {
    case 'todo': return '○'
    case 'doing': return '●'
    case 'done': return '✓'
    default: return '○'
  }
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

const goToTask = () => {
  router.push('/tasks')
}

onMounted(async () => {
  // Load initial data
  await workspaceStore.fetchWorkspaces()
  await tasksStore.fetchTags()
  await tasksStore.fetchTasks()
  await tasksStore.fetchStats()

  // Add subtle animations on scroll
  const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  if (prefersReducedMotion) return

  // Animate stats cards on scroll
  const observer = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        entry.target.classList.add('in-view')
      }
    })
  }, { threshold: 0.1 })

  // Observe all animated elements
  document.querySelectorAll('.anim-title, .anim-subtitle, .anim-stat-card, .anim-action, .anim-task-item').forEach(el => {
    observer.observe(el)
  })
})
</script>

<style scoped>
.dashboard-zen {
  min-height: 100vh;
  background: #FAF8F3;
  color: #3D3D3D;
  transition: background 0.3s ease, color 0.3s ease;
  padding-top: 80px;
}

.dashboard-zen.dark-mode {
  background: #1A1410;
  color: #F5F1E8;
}

.dashboard-content {
  max-width: 900px;
  margin: 0 auto;
  padding: 3rem 2rem;
  box-sizing: border-box;
}

/* Welcome section with depth */
.welcome-section {
  position: relative;
  margin-bottom: 4rem;
  padding: 2rem 0;
  border-bottom: 1px solid rgba(168, 153, 104, 0.1);
}

.dashboard-zen.dark-mode .welcome-section {
  border-bottom-color: rgba(212, 165, 116, 0.15);
}

.welcome-bg-glow {
  position: absolute;
  top: -50%;
  right: -20%;
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, rgba(168, 153, 104, 0.08) 0%, transparent 70%);
  filter: blur(40px);
  pointer-events: none;
}

.welcome-inner {
  position: relative;
  z-index: 1;
  text-align: center;
}

.welcome-title {
  margin: 0 0 0.5rem;
  font-size: 2.5rem;
  font-weight: 300;
  font-family: Georgia, serif;
  letter-spacing: -1px;
  color: #3D3D3D;
  opacity: 0;
  transform: translateY(20px);
  animation: fadeInUp 0.8s ease-out 0.2s forwards;
}

.dashboard-zen.dark-mode .welcome-title {
  color: #F5F1E8;
}

.welcome-subtitle {
  margin: 0;
  font-size: 1rem;
  color: #8B8B8B;
  font-weight: 300;
  letter-spacing: 0.3px;
  opacity: 0;
  transform: translateY(20px);
  animation: fadeInUp 0.8s ease-out 0.4s forwards;
}

.dashboard-zen.dark-mode .welcome-subtitle {
  color: #C9C1B8;
}

/* Stats section */
.stats-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
  margin-bottom: 3rem;
}

.stat-card {
  padding: 2rem;
  background: white;
  border-radius: 12px;
  border: 1px solid #E8DDD2;
  text-align: center;
  transition: all 0.3s ease;
}

.dashboard-zen.dark-mode .stat-card {
  background: #2D2420;
  border-color: #3D3530;
}

.stat-card:hover {
  border-color: #A89968;
  box-shadow: 0 4px 12px rgba(168, 153, 104, 0.08);
}

.dashboard-zen.dark-mode .stat-card:hover {
  border-color: #F5DEB3;
  box-shadow: 0 4px 12px rgba(245, 222, 179, 0.08);
}

.stat-value {
  font-size: 2.5rem;
  font-weight: 300;
  font-family: Georgia, serif;
  color: #A89968;
  margin-bottom: 0.5rem;
}

.dashboard-zen.dark-mode .stat-value {
  color: #F5DEB3;
}

.stat-label {
  font-size: 0.85rem;
  color: #8B8B8B;
  font-weight: 400;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.dashboard-zen.dark-mode .stat-label {
  color: #A0AEC0;
}

/* Actions section */
.actions-section {
  display: flex;
  gap: 1rem;
  margin-bottom: 3rem;
  justify-content: center;
}

.action-btn {
  padding: 0.75rem 2rem;
  border-radius: 8px;
  border: none;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: Inter, sans-serif;
  font-size: 0.95rem;
  letter-spacing: 0.3px;
}

.action-btn.primary {
  background: #A89968;
  color: white;
}

.action-btn.primary:hover {
  background: #8B7355;
}

.action-btn.secondary {
  background: white;
  color: #A89968;
  border: 1px solid #A89968;
}

.action-btn.secondary:hover {
  background: #A89968;
  color: white;
}

.dashboard-zen.dark-mode .action-btn.primary {
  background: #D4A574;
  color: #1A1410;
}

.dashboard-zen.dark-mode .action-btn.primary:hover {
  background: #F5DEB3;
}

.dashboard-zen.dark-mode .action-btn.secondary {
  background: transparent;
  color: #F5DEB3;
  border-color: #F5DEB3;
}

.dashboard-zen.dark-mode .action-btn.secondary:hover {
  background: #F5DEB3;
  color: #1A1410;
}

/* Recent section */
.section-title {
  margin: 0 0 1.5rem;
  font-size: 1.5rem;
  font-weight: 400;
  font-family: Georgia, serif;
  color: #3D3D3D;
}

.dashboard-zen.dark-mode .section-title {
  color: #F5F1E8;
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.task-item {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem;
  background: white;
  border-radius: 10px;
  border: 1px solid #E8DDD2;
  transition: all 0.3s ease;
  cursor: pointer;
}

.dashboard-zen.dark-mode .task-item {
  background: #2D2420;
  border-color: #3D3530;
}

.task-item:hover {
  border-color: #A89968;
}

.task-status {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #A89968;
  flex-shrink: 0;
  font-size: 0.8rem;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: bold;
}

.task-status.status-todo {
  background: #D4A574;
}

.task-status.status-doing {
  background: #A89968;
}

.task-status.status-done {
  background: #8B7355;
}

.task-status.status-archived {
  background: #C9C1B8;
}

.task-info {
  flex: 1;
}

.task-name {
  margin: 0 0 0.25rem;
  font-weight: 500;
  color: #3D3D3D;
}

.dashboard-zen.dark-mode .task-name {
  color: #F5F1E8;
}

.task-desc {
  margin: 0 0 0.5rem;
  font-size: 0.85rem;
  color: #8B8B8B;
}

.dashboard-zen.dark-mode .task-desc {
  color: #C9C1B8;
}

.task-meta-bottom {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.task-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.task-tag {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.65rem;
  font-weight: 500;
  white-space: nowrap;
}

.task-due-date {
  font-size: 0.75rem;
  color: #A89968;
  font-weight: 500;
}

.dashboard-zen.dark-mode .task-due-date {
  color: #F5DEB3;
}



.task-priority {
  padding: 0.4rem 0.8rem;
  border-radius: 6px;
  background: rgba(168, 153, 104, 0.1);
  color: #A89968;
  font-size: 0.75rem;
  font-weight: 500;
  white-space: nowrap;
  text-transform: capitalize;
}

.task-priority.priority-high {
  background: rgba(212, 165, 116, 0.15);
  color: #D4A574;
}

.task-priority.priority-medium {
  background: rgba(168, 153, 104, 0.1);
  color: #A89968;
}

.task-priority.priority-low {
  background: rgba(139, 115, 85, 0.08);
  color: #8B7355;
}

.dashboard-zen.dark-mode .task-priority {
  background: rgba(212, 165, 116, 0.15);
  color: #F5DEB3;
}

.empty-message {
  text-align: center;
  color: #8B8B8B;
  padding: 3rem;
  font-size: 1rem;
}

.dashboard-zen.dark-mode .empty-message {
  color: #C9C1B8;
}

/* Animations */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes slideInLeft {
  from {
    opacity: 0;
    transform: translateX(-30px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

/* Scroll-triggered animations */
.anim-stat-card {
  opacity: 0;
  transform: translateY(40px) scale(0.95);
  transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
  transition-delay: calc(var(--index, 0) * 0.1s);
}

.anim-stat-card.in-view {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.anim-action {
  opacity: 0;
  transform: translateY(20px);
  animation: fadeInUp 0.6s ease-out forwards;
}

.anim-action:nth-child(1) {
  animation-delay: 0.6s;
}

.anim-action:nth-child(2) {
  animation-delay: 0.8s;
}

.anim-task-item {
  opacity: 0;
  transform: translateX(-20px);
  transition: all 0.5s ease-out;
  transition-delay: calc(var(--index, 0) * 0.1s);
}

.anim-task-item.in-view {
  opacity: 1;
  transform: translateX(0);
}

/* Reduced motion */
@media (prefers-reduced-motion: reduce) {
  .welcome-title,
  .welcome-subtitle,
  .anim-action,
  .anim-stat-card,
  .anim-task-item {
    animation: none;
    transition: none;
    opacity: 1;
    transform: none;
  }
}

/* Mobile */
@media (max-width: 768px) {
  .dashboard-content {
    padding: 1.5rem 0.75rem;
  }

  .welcome-section {
    margin-bottom: 2rem;
    padding: 1rem 0.75rem;
    overflow: hidden;
  }

  .welcome-bg-glow {
    width: 300px;
    height: 300px;
    top: -30%;
    right: -10%;
    filter: blur(30px);
  }

  .welcome-title {
    font-size: 1.5rem;
  }

  .welcome-subtitle {
    font-size: 0.9rem;
  }

  .stats-section {
    grid-template-columns: repeat(2, 1fr);
    gap: 0.75rem;
    margin-bottom: 2rem;
  }

  .stat-card {
    padding: 1rem;
  }

  .stat-value {
    font-size: 1.75rem;
  }

  .stat-label {
    font-size: 0.75rem;
  }

  .actions-section {
    flex-direction: column;
    margin-bottom: 2rem;
  }

  .action-btn {
    padding: 0.6rem 1.5rem;
    font-size: 0.85rem;
  }

  .section-title {
    font-size: 1.25rem;
  }

  .task-item {
    padding: 1rem;
    gap: 0.75rem;
  }

  .task-name {
    font-size: 0.9rem;
  }

  .task-desc {
    font-size: 0.8rem;
  }

  .task-priority {
    font-size: 0.7rem;
    padding: 0.3rem 0.6rem;
  }
}
</style>
