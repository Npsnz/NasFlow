import { defineStore } from 'pinia'
import client from '@/api/client'

export interface Tag {
  id: number
  name: string
  color: string
  user_id: number
}

export interface Comment {
  id: number
  content: string
  task_id: number
  user_id: number
  user: {
    id: number
    name: string
    avatar_url?: string
  }
  created_at: string
}

export interface Task {
  id: number
  title: string
  description: string
  status: 'todo' | 'doing' | 'done' | 'archived'
  priority: 'low' | 'medium' | 'high' | 'urgent'
  workspace_id: number
  user_id: number
  due_date: string | null
  completed_at: string | null
  is_recurring: boolean
  recur_rule: string
  parent_task_id: number | null
  sort_order: number
  tags: Tag[]
  comments?: Comment[]
  subtasks?: Task[]
  created_at: string
  updated_at: string
}

export interface Stats {
  total: number
  todo: number
  doing: number
  done_today: number
  overdue: number
  due_this_week: number
}

let sseSource: EventSource | null = null

export const useTasksStore = defineStore('tasks', {
  state: () => ({
    tasks: [] as Task[],
    tags: [] as Tag[],
    stats: {
      total: 0,
      todo: 0,
      doing: 0,
      done_today: 0,
      overdue: 0,
      due_this_week: 0
    } as Stats,
    filters: {
      search: '',
      tagIds: [] as number[],
      priority: '',
      dueRange: '' // 'today' | 'tomorrow' | 'week' | 'overdue' | ''
    },
    activities: [] as string[], // Activity feed events
    loading: false,
    error: null as string | null
  }),

  actions: {
    // Register activity event (Thai relative time will be formatted in view)
    addActivity(actionText: string) {
      this.activities.unshift(`${actionText}|${new Date().toISOString()}`)
      if (this.activities.length > 10) {
        this.activities.pop()
      }
    },

    async fetchStats() {
      try {
        const res = await client.get('/stats')
        this.stats = res.data.data
      } catch (err) {
        console.error('Failed to fetch stats', err)
      }
    },

    async fetchTags() {
      try {
        const res = await client.get('/tags')
        this.tags = res.data.data || []
      } catch (err) {
        console.error('Failed to fetch tags', err)
      }
    },

    async createTag(name: string, color: string) {
      try {
        const res = await client.post('/tags', { name, color })
        this.tags.push(res.data.data)
        return res.data.data
      } catch (err: any) {
        throw err.response?.data?.error || 'สร้างแท็กไม่สำเร็จ'
      }
    },

    async deleteTag(id: number) {
      try {
        await client.delete(`/tags/${id}`)
        this.tags = this.tags.filter(t => t.id !== id)
      } catch (err) {
        console.error('Failed to delete tag', err)
      }
    },

    async fetchTasks(filters: Record<string, any> = {}) {
      this.loading = true
      try {
        // Construct query parameters
        const params = new URLSearchParams()
        Object.entries(filters).forEach(([key, val]) => {
          if (val !== undefined && val !== null && val !== '') {
            params.append(key, String(val))
          }
        })

        const res = await client.get(`/tasks?${params.toString()}`)
        this.tasks = res.data.data || []
      } catch (err: any) {
        this.error = err.response?.data?.error || 'โหลดรายการงานไม่สำเร็จ'
      } finally {
        this.loading = false
      }
    },

    async getTaskDetails(id: number): Promise<Task> {
      try {
        const res = await client.get(`/tasks/${id}`)
        return res.data.data
      } catch (err: any) {
        throw err.response?.data?.error || 'โหลดรายละเอียดงานไม่สำเร็จ'
      }
    },

    async createTask(payload: Partial<Task> & { tag_ids?: number[] }) {
      try {
        const res = await client.post('/tasks', payload)
        const newTask = res.data.data
        
        // SSE will push event, but let's push locally if not SSE connected
        if (!sseSource) {
          this.tasks.push(newTask)
          this.fetchStats()
        }
        return newTask
      } catch (err: any) {
        throw err.response?.data?.error || 'สร้างงานไม่สำเร็จ'
      }
    },

    async updateTask(id: number, payload: Partial<Task> & { tag_ids?: number[] }) {
      // Find task to do optimistic update or cache state
      const taskIndex = this.tasks.findIndex(t => t.id === id)
      const oldTask = taskIndex !== -1 ? { ...this.tasks[taskIndex] } : null

      // If status changed, let's update locally immediately (Optimistic Update)
      if (oldTask && payload.status && payload.status !== oldTask.status) {
        this.tasks[taskIndex].status = payload.status as any
        this.fetchStats()
      }

      try {
        const res = await client.put(`/tasks/${id}`, payload)
        const updated = res.data.data

        if (taskIndex !== -1) {
          this.tasks[taskIndex] = updated
        }
        
        if (!sseSource) {
          this.fetchStats()
        }
        return updated
      } catch (err: any) {
        // Revert optimistic update
        if (oldTask && taskIndex !== -1) {
          this.tasks[taskIndex] = oldTask
          this.fetchStats()
        }
        throw err.response?.data?.error || 'แก้ไขงานไม่สำเร็จ'
      }
    },

    async updateTaskStatus(id: number, status: string) {
      const taskIndex = this.tasks.findIndex(t => t.id === id)
      if (taskIndex === -1) return
      
      const oldStatus = this.tasks[taskIndex].status
      
      // Optimistic update
      this.tasks[taskIndex].status = status as any
      this.fetchStats()

      try {
        const res = await client.put(`/tasks/${id}/status`, { status })
        const updated = res.data.data
        this.tasks[taskIndex] = updated
        return updated
      } catch (err: any) {
        // Revert
        this.tasks[taskIndex].status = oldStatus
        this.fetchStats()
        throw err.response?.data?.error || 'ย้ายงานไม่สำเร็จ'
      }
    },

    async reorderTasks(items: { id: number; sort_order: number; status?: string }[]) {
      // Optimistic sorting updates
      const oldTasks = [...this.tasks]
      
      items.forEach(item => {
        const t = this.tasks.find(x => x.id === item.id)
        if (t) {
          t.sort_order = item.sort_order
          if (item.status) t.status = item.status as any
        }
      })
      this.tasks.sort((a, b) => a.sort_order - b.sort_order)
      this.fetchStats()

      try {
        await client.put('/tasks/reorder', { tasks: items })
      } catch (err) {
        // Revert
        this.tasks = oldTasks
        this.fetchStats()
        throw err
      }
    },

    async completeTask(id: number) {
      const taskIndex = this.tasks.findIndex(t => t.id === id)
      if (taskIndex === -1) return

      const oldStatus = this.tasks[taskIndex].status
      this.tasks[taskIndex].status = 'done'
      this.fetchStats()

      try {
        const res = await client.post(`/tasks/${id}/complete`)
        const updated = res.data.data
        this.tasks[taskIndex] = updated

        // If recurring task created next occurrence, add it to list
        if (res.data.next_occurrence) {
          const nextOccur = res.data.next_occurrence
          this.tasks.push(nextOccur)
        }
        return res.data
      } catch (err: any) {
        this.tasks[taskIndex].status = oldStatus
        this.fetchStats()
        throw err.response?.data?.error || 'ทำเครื่องหมายเสร็จงานไม่สำเร็จ'
      }
    },

    async deleteTask(id: number) {
      try {
        await client.delete(`/tasks/${id}`)
        this.tasks = this.tasks.filter(t => t.id !== id)
        this.fetchStats()
      } catch (err: any) {
        throw err.response?.data?.error || 'ลบงานไม่สำเร็จ'
      }
    },

    // Comments Actions
    async addComment(taskId: number, content: string) {
      try {
        const res = await client.post(`/tasks/${taskId}/comments`, { content })
        return res.data.data
      } catch (err: any) {
        throw err.response?.data?.error || 'แสดงความคิดเห็นไม่สำเร็จ'
      }
    },

    async deleteComment(commentId: number) {
      try {
        await client.delete(`/comments/${commentId}`)
      } catch (err: any) {
        throw err.response?.data?.error || 'ลบความคิดเห็นไม่สำเร็จ'
      }
    },

    // Server-Sent Events Initialization
    initSSE() {
      if (sseSource) return

      // SSE connection
      const proto = window.location.protocol === 'https:' ? 'https:' : 'http:'
      const host = window.location.host
      sseSource = new EventSource(`${proto}//${host}/api/sse`, { withCredentials: true })

      sseSource.addEventListener('task.created', (e: any) => {
        const data = JSON.parse(e.data)
        const idx = this.tasks.findIndex(t => t.id === data.id)
        if (idx === -1) {
          this.tasks.push(data)
          this.addActivity(`สร้างงาน "${data.title}"`)
        }
        this.fetchStats()
      })

      sseSource.addEventListener('task.updated', (e: any) => {
        const data = JSON.parse(e.data)
        // Check if batch reorder
        if (data.reordered) {
          this.fetchTasks() // Full refresh for reorder
          return
        }

        const idx = this.tasks.findIndex(t => t.id === data.id)
        if (idx !== -1) {
          this.tasks[idx] = data
        }
        this.addActivity(`แก้ไขงาน "${data.title}"`)
        this.fetchStats()
      })

      sseSource.addEventListener('task.moved', (e: any) => {
        const payload = JSON.parse(e.data)
        const task = payload.task
        const oldStatus = payload.old_status

        const idx = this.tasks.findIndex(t => t.id === task.id)
        if (idx !== -1) {
          this.tasks[idx] = task
        } else {
          this.tasks.push(task)
        }

        // Map status names for logs
        const statusMap: Record<string, string> = {
          todo: 'รอดำเนินการ',
          doing: 'กำลังดำเนินการ',
          done: 'เสร็จแล้ว',
          archived: 'จัดเก็บแล้ว'
        }
        this.addActivity(`ย้าย "${task.title}" ไป ${statusMap[task.status] || task.status} (เดิม: ${statusMap[oldStatus] || oldStatus})`)
        this.fetchStats()
      })

      sseSource.addEventListener('task.deleted', (e: any) => {
        const data = JSON.parse(e.data)
        const deletedTask = this.tasks.find(t => t.id === data.id)
        const title = deletedTask ? deletedTask.title : `ID ${data.id}`
        this.tasks = this.tasks.filter(t => t.id !== data.id)
        this.addActivity(`ลบงาน "${title}"`)
        this.fetchStats()
      })

      sseSource.addEventListener('tasks.overdue', (e: any) => {
        const payload = JSON.parse(e.data)
        this.stats.overdue = payload.count
        
        // Dispatch window event so views can alert the user
        window.dispatchEvent(new CustomEvent('tasks-overdue-alert', { detail: payload.count }))
      })

      sseSource.onerror = (e) => {
        console.error('SSE connection error, closing stream.', e)
        this.closeSSE()
      }
    },

    closeSSE() {
      if (sseSource) {
        sseSource.close()
        sseSource = null
      }
    }
  }
})
