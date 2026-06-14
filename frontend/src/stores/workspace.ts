import { defineStore } from 'pinia'
import client from '@/api/client'

export interface Workspace {
  id: number
  name: string
  slug: string
  color: string
  icon: string
  user_id: number
  sort_order: number
  is_archived: boolean
  created_at: string
  updated_at: string
}

export const useWorkspaceStore = defineStore('workspace', {
  state: () => ({
    workspaces: [] as Workspace[],
    currentWorkspace: null as Workspace | null,
    loading: false,
    error: null as string | null
  }),

  actions: {
    async fetchWorkspaces(showArchived = false) {
      this.loading = true
      try {
        const res = await client.get(`/workspaces?archived=${showArchived}`)
        this.workspaces = res.data.data || []
        
        // Select first active workspace if none selected
        if (!this.currentWorkspace && this.workspaces.length > 0) {
          this.currentWorkspace = this.workspaces[0]
        }
      } catch (err: any) {
        this.error = err.response?.data?.error || 'โหลดพื้นที่งานไม่สำเร็จ'
      } finally {
        this.loading = false
      }
    },

    async createWorkspace(payload: { name: string; color?: string; icon?: string }) {
      try {
        const res = await client.post('/workspaces', payload)
        const newWs = res.data.data
        this.workspaces.push(newWs)
        if (!this.currentWorkspace) {
          this.currentWorkspace = newWs
        }
        return newWs
      } catch (err: any) {
        this.error = err.response?.data?.error || 'สร้างพื้นที่งานไม่สำเร็จ'
        throw err
      }
    },

    async updateWorkspace(id: number, payload: Partial<Workspace>) {
      try {
        const res = await client.put(`/workspaces/${id}`, payload)
        const updated = res.data.data
        
        // Update local list
        const idx = this.workspaces.findIndex(w => w.id === id)
        if (idx !== -1) {
          this.workspaces[idx] = updated
        }
        
        // Update current workspace if it's the updated one
        if (this.currentWorkspace?.id === id) {
          this.currentWorkspace = updated
        }
        return updated
      } catch (err: any) {
        this.error = err.response?.data?.error || 'แก้ไขพื้นที่งานไม่สำเร็จ'
        throw err
      }
    },

    async deleteWorkspace(id: number, hard = false) {
      try {
        await client.delete(`/workspaces/${id}?hard=${hard}`)
        
        // Remove from local list
        this.workspaces = this.workspaces.filter(w => w.id !== id)
        
        // Re-select active workspace if the deleted one was selected
        if (this.currentWorkspace?.id === id) {
          this.currentWorkspace = this.workspaces.length > 0 ? this.workspaces[0] : null
        }
      } catch (err: any) {
        this.error = err.response?.data?.error || 'ลบพื้นที่งานไม่สำเร็จ'
        throw err
      }
    },

    async reorderWorkspaces(ids: number[]) {
      // Optimistic update
      const oldList = [...this.workspaces]
      this.workspaces.sort((a, b) => ids.indexOf(a.id) - ids.indexOf(b.id))
      
      try {
        await client.put('/workspaces/reorder', { ids })
      } catch (err) {
        // Revert on error
        this.workspaces = oldList
        throw err
      }
    },

    setCurrentWorkspace(workspace: Workspace | null) {
      this.currentWorkspace = workspace
    }
  }
})
