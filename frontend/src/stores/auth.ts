import { defineStore } from 'pinia'
import client from '@/api/client'

export interface User {
  id: number
  email: string
  name: string
  avatar_url?: string
  created_at: string
  updated_at: string
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('tf_user') || 'null') as User | null,
    token: localStorage.getItem('tf_token') || null as string | null,
    isAuthenticated: !!localStorage.getItem('tf_token'),
    loading: false,
    error: null as string | null
  }),

  actions: {
    persistState(user: User | null, token: string | null) {
      this.user = user
      this.token = token
      this.isAuthenticated = !!token

      if (user) {
        localStorage.setItem('tf_user', JSON.stringify(user))
      } else {
        localStorage.removeItem('tf_user')
      }

      if (token) {
        localStorage.setItem('tf_token', token)
      } else {
        localStorage.removeItem('tf_token')
      }
    },

    async register(payload: any) {
      this.loading = true
      this.error = null
      try {
        const res = await client.post('/auth/register', payload)
        const user = res.data.data
        const token = res.data.token || 'logged_in' // Token is set in cookie, fallback value for state
        this.persistState(user, token)
        return user
      } catch (err: any) {
        this.error = err.response?.data?.error || 'เกิดข้อผิดพลาดในการสมัครสมาชิก'
        throw err
      } finally {
        this.loading = false
      }
    },

    async login(payload: any) {
      this.loading = true
      this.error = null
      try {
        const res = await client.post('/auth/login', payload)
        const user = res.data.data
        const token = res.data.token || 'logged_in'
        this.persistState(user, token)
        return user
      } catch (err: any) {
        this.error = err.response?.data?.error || 'อีเมลหรือรหัสผ่านไม่ถูกต้อง'
        throw err
      } finally {
        this.loading = false
      }
    },

    async logout() {
      this.loading = true
      try {
        await client.post('/auth/logout')
      } catch (err) {
        console.error('Logout request failed', err)
      } finally {
        this.persistState(null, null)
        this.loading = false
      }
    },

    async fetchMe() {
      if (!this.token) return
      try {
        const res = await client.get('/auth/me')
        this.user = res.data.data
        localStorage.setItem('tf_user', JSON.stringify(this.user))
      } catch (err) {
        // If fetchMe fails due to unauthorized, logout
        this.persistState(null, null)
      }
    },

    async updateProfile(payload: any) {
      this.loading = true
      this.error = null
      try {
        const res = await client.put('/auth/profile', payload)
        this.user = res.data.data
        localStorage.setItem('tf_user', JSON.stringify(this.user))
        return this.user
      } catch (err: any) {
        this.error = err.response?.data?.error || 'ไม่สามารถแก้ไขโปรไฟล์ได้'
        throw err
      } finally {
        this.loading = false
      }
    },

    async deleteAccount() {
      this.loading = true
      try {
        await client.delete('/auth/delete')
        this.persistState(null, null)
      } catch (err: any) {
        this.error = err.response?.data?.error || 'ไม่สามารถลบบัญชีผู้ใช้ได้'
        throw err
      } finally {
        this.loading = false
      }
    }
  }
})

// Listen to custom axios logout event
if (typeof window !== 'undefined') {
  window.addEventListener('auth-logout', () => {
    const store = useAuthStore()
    store.persistState(null, null)
  })
}
