<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useWorkspaceStore } from '@/stores/workspace'
import type { Workspace } from '@/stores/workspace'
import { useTasksStore } from '@/stores/tasks'
import { useUIStore } from '@/stores/ui'
import {
  User,
  Key,
  FolderOpen,
  Tag as TagIcon,
  Globe,
  Database,
  Trash2,
  Plus,
  Moon,
  Sun,
  Laptop,
  Check
} from 'lucide-vue-next'

const authStore = useAuthStore()
const workspaceStore = useWorkspaceStore()
const tasksStore = useTasksStore()
const uiStore = useUIStore()

// State tabs for Settings Sections
const activeSection = ref<'profile' | 'workspaces' | 'tags' | 'preferences' | 'data'>('profile')

// Profile State
const name = ref(authStore.user?.name || '')
const email = ref(authStore.user?.email || '')
const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')

// Workspaces editing
const editingWorkspaceId = ref<number | null>(null)
const editWsName = ref('')
const editWsColor = ref('#171717')
const editWsIcon = ref('briefcase')

// Tags editing
const editingTagId = ref<number | null>(null)
const editTagName = ref('')
const editTagColor = ref('#171717')

// New workspace/tag states
const showNewWsForm = ref(false)
const newWsName = ref('')
const newWsColor = ref('#171717')
const newWsIcon = ref('briefcase')

const showNewTagForm = ref(false)
const newTagName = ref('')
const newTagColor = ref('#171717')

const presetColors = ['#171717', '#1D9E75', '#D85A30', '#333333', '#ef4444', '#3b82f6', '#10b981', '#f59e0b']

onMounted(() => {
  uiStore.setActiveView('settings')
  workspaceStore.fetchWorkspaces(true) // Fetch all including archived
  tasksStore.fetchTags()
})

// Avatar image FileReader base64 upload
const handleAvatarUpload = (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = async () => {
    const base64 = reader.result as string
    try {
      await authStore.updateProfile({ avatar_url: base64 })
      uiStore.showToast('อัปโหลดรูปโปรไฟล์สำเร็จแล้ว', 'success')
    } catch (err) {
      uiStore.showToast('ไม่สามารถอัปโหลดรูปโปรไฟล์ได้', 'error')
    }
  }
  reader.readAsDataURL(file)
}

const handleUpdateProfile = async () => {
  if (!name.value || !email.value) {
    uiStore.showToast('กรุณากรอกข้อมูลชื่อและอีเมล', 'error')
    return
  }

  try {
    await authStore.updateProfile({ name: name.value, email: email.value })
    uiStore.showToast('บันทึกข้อมูลส่วนตัวสำเร็จ', 'success')
  } catch (err: any) {
    uiStore.showToast(err || 'ไม่สามารถบันทึกข้อมูลส่วนตัวได้', 'error')
  }
}

const handleChangePassword = async () => {
  if (!currentPassword.value || !newPassword.value) {
    uiStore.showToast('กรุณากรอกรหัสผ่านเพื่อเปลี่ยน', 'error')
    return
  }
  if (newPassword.value.length < 6) {
    uiStore.showToast('รหัสผ่านใหม่ต้องมีความยาวไม่น้อยกว่า 6 ตัวอักษร', 'error')
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    uiStore.showToast('รหัสผ่านใหม่กับรหัสผ่านยืนยันไม่ตรงกัน', 'error')
    return
  }

  try {
    await authStore.updateProfile({
      current_password: currentPassword.value,
      new_password: newPassword.value
    })
    uiStore.showToast('เปลี่ยนรหัสผ่านสำเร็จแล้ว', 'success')
    currentPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
  } catch (err: any) {
    uiStore.showToast(err || 'เปลี่ยนรหัสผ่านล้มเหลว', 'error')
  }
}

// Workspaces logic
const startEditWs = (ws: Workspace) => {
  editingWorkspaceId.value = ws.id
  editWsName.value = ws.name
  editWsColor.value = ws.color
  editWsIcon.value = ws.icon
}

const saveWsEdit = async (id: number) => {
  if (!editWsName.value.trim()) return
  try {
    await workspaceStore.updateWorkspace(id, {
      name: editWsName.value.trim(),
      color: editWsColor.value,
      icon: editWsIcon.value
    })
    editingWorkspaceId.value = null
    uiStore.showToast('แก้ไขพื้นที่งานสำเร็จ', 'success')
  } catch (err) {
    uiStore.showToast('แก้ไขพื้นที่งานล้มเหลว', 'error')
  }
}

const handleCreateWorkspace = async () => {
  if (!newWsName.value.trim()) return
  try {
    await workspaceStore.createWorkspace({
      name: newWsName.value.trim(),
      color: newWsColor.value,
      icon: newWsIcon.value
    })
    newWsName.value = ''
    showNewWsForm.value = false
    uiStore.showToast('สร้างพื้นที่งานใหม่สำเร็จ', 'success')
  } catch (err) {
    uiStore.showToast('ไม่สามารถสร้างพื้นที่งานได้', 'error')
  }
}

const handleDeleteWorkspace = async (id: number) => {
  if (confirm('คำเตือน! การลบพื้นที่งานนี้จะลบงานทั้งหมดในพื้นที่นี้อย่างถาวร ยืนยันที่จะลบหรือไม่?')) {
    try {
      await workspaceStore.deleteWorkspace(id, true) // Hard delete
      uiStore.showToast('ลบพื้นที่งานถาวรเรียบร้อยแล้ว', 'success')
    } catch (err) {
      uiStore.showToast('ลบพื้นที่งานไม่สำเร็จ', 'error')
    }
  }
}

const toggleArchiveWorkspace = async (ws: Workspace) => {
  try {
    await workspaceStore.updateWorkspace(ws.id, { is_archived: !ws.is_archived })
    uiStore.showToast(ws.is_archived ? 'เปิดใช้งานพื้นที่งานแล้ว' : 'ปิดใช้งานพื้นที่งานแล้ว', 'success')
  } catch (err) {
    uiStore.showToast('อัปเดตสถานะพื้นที่งานล้มเหลว', 'error')
  }
}

// Workspace order shift reordering
const moveWorkspaceOrder = async (wsId: number, direction: 'up' | 'down') => {
  const wsList = [...workspaceStore.workspaces]
  const idx = wsList.findIndex(w => w.id === wsId)
  if (idx === -1) return
  if (direction === 'up' && idx === 0) return
  if (direction === 'down' && idx === wsList.length - 1) return

  const targetIdx = direction === 'up' ? idx - 1 : idx + 1
  // Swap
  const temp = wsList[idx]
  wsList[idx] = wsList[targetIdx]
  wsList[targetIdx] = temp

  try {
    await workspaceStore.reorderWorkspaces(wsList.map(w => w.id))
    uiStore.showToast('สลับตำแหน่งพื้นที่งานแล้ว', 'success')
  } catch (err) {
    uiStore.showToast('จัดตำแหน่งเรียงลำดับไม่สำเร็จ', 'error')
  }
}

// Tags logic
const startEditTag = (tag: any) => {
  editingTagId.value = tag.id
  editTagName.value = tag.name
  editTagColor.value = tag.color
}

const saveTagEdit = async (id: number) => {
  if (!editTagName.value.trim()) return
  try {
    await tasksStore.createTag(editTagName.value.trim(), editTagColor.value) // Actually we call updateTag
    // Since taskStore manages tags list, let's write a wrapper updateTag inside taskStore:
    // Wait, did we define updateTag in tasksStore? No, but let's check: we can call API directly or we can add it.
    // Let's call API directly to update tag, then refresh tags list!
    await client.put(`/tags/${id}`, { name: editTagName.value.trim(), color: editTagColor.value })
    await tasksStore.fetchTags()
    editingTagId.value = null
    uiStore.showToast('แก้ไขแท็กสำเร็จ', 'success')
  } catch (err) {
    uiStore.showToast('แก้ไขแท็กล้มเหลว', 'error')
  }
}

const handleCreateTag = async () => {
  if (!newTagName.value.trim()) return
  try {
    await tasksStore.createTag(newTagName.value.trim(), newTagColor.value)
    newTagName.value = ''
    showNewTagForm.value = false
    uiStore.showToast('สร้างแท็กใหม่สำเร็จ', 'success')
  } catch (err) {
    uiStore.showToast('สร้างแท็กล้มเหลว', 'error')
  }
}

const handleDeleteTag = async (id: number) => {
  if (confirm('ต้องการลบแท็กนี้หรือไม่? (ไม่มีผลต่อการลบตัวงาน)')) {
    try {
      await tasksStore.deleteTag(id)
      uiStore.showToast('ลบแท็กสำเร็จ', 'success')
    } catch (err) {
      uiStore.showToast('ลบแท็กล้มเหลว', 'error')
    }
  }
}

// Theme toggling
const changeTheme = (mode: 'light' | 'dark' | 'system') => {
  uiStore.setTheme(mode)
  uiStore.showToast('อัปเดตการแสดงผลของธีมแล้ว', 'success')
}

// Export raw JSON backup
const exportJSON = () => {
  const dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(tasksStore.tasks, null, 2))
  const link = document.createElement('a')
  link.setAttribute("href", dataStr)
  link.setAttribute("download", "taskflow_backup.json")
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  uiStore.showToast('ดาวน์โหลด JSON สำเร็จ', 'success')
}

// Delete account completely
const handleDeleteAccount = async () => {
  if (confirm('คำเตือนร้ายแรง! การลบบัญชีจะเป็นการทำลายงานทั้งหมดของคุณ และลบข้อมูลออกจากฐานข้อมูลอย่างถาวรโดยไม่สามารถกู้คืนได้ ยืนยันการลบบัญชีหรือไม่?')) {
    try {
      await authStore.deleteAccount()
      uiStore.showToast('ลบบัญชีของคุณสำเร็จแล้ว', 'success')
      window.location.href = '/login'
    } catch (err) {
      uiStore.showToast('ไม่สามารถลบบัญชีได้', 'error')
    }
  }
}

import client from '@/api/client'
</script>

<template>
  <div class="p-6 flex flex-col md:flex-row gap-6 select-none max-w-5xl mx-auto h-full overflow-hidden">
    <!-- Left Navigation Menu -->
    <div class="w-full md:w-56 flex-shrink-0 flex flex-row md:flex-col gap-1 md:gap-1.5 border-b md:border-b-0 md:border-r border-border-light dark:border-border-dark pb-4 md:pb-0 md:pr-4">
      <button
        @click="activeSection = 'profile'"
        class="flex-1 md:flex-none flex items-center space-x-2.5 px-3.5 py-2.5 rounded-lg text-xs font-semibold hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors text-left min-h-[44px]"
        :class="activeSection === 'profile' ? 'bg-slate-100 dark:bg-slate-800 text-brand-500 font-bold' : 'text-slate-600 dark:text-slate-400'"
      >
        <User class="w-4 h-4" />
        <span>ข้อมูลส่วนตัว (Profile)</span>
      </button>

      <button
        @click="activeSection = 'workspaces'"
        class="flex-1 md:flex-none flex items-center space-x-2.5 px-3.5 py-2.5 rounded-lg text-xs font-semibold hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors text-left min-h-[44px]"
        :class="activeSection === 'workspaces' ? 'bg-slate-100 dark:bg-slate-800 text-brand-500 font-bold' : 'text-slate-600 dark:text-slate-400'"
      >
        <FolderOpen class="w-4 h-4" />
        <span>จัดการพื้นที่ทำงาน</span>
      </button>

      <button
        @click="activeSection = 'tags'"
        class="flex-1 md:flex-none flex items-center space-x-2.5 px-3.5 py-2.5 rounded-lg text-xs font-semibold hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors text-left min-h-[44px]"
        :class="activeSection === 'tags' ? 'bg-slate-100 dark:bg-slate-800 text-brand-500 font-bold' : 'text-slate-600 dark:text-slate-400'"
      >
        <TagIcon class="w-4 h-4" />
        <span>จัดการแท็ก (Tags)</span>
      </button>

      <button
        @click="activeSection = 'preferences'"
        class="flex-1 md:flex-none flex items-center space-x-2.5 px-3.5 py-2.5 rounded-lg text-xs font-semibold hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors text-left min-h-[44px]"
        :class="activeSection === 'preferences' ? 'bg-slate-100 dark:bg-slate-800 text-brand-500 font-bold' : 'text-slate-600 dark:text-slate-400'"
      >
        <Globe class="w-4 h-4" />
        <span>การแสดงผล (Theme)</span>
      </button>

      <button
        @click="activeSection = 'data'"
        class="flex-1 md:flex-none flex items-center space-x-2.5 px-3.5 py-2.5 rounded-lg text-xs font-semibold hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors text-left min-h-[44px]"
        :class="activeSection === 'data' ? 'bg-slate-100 dark:bg-slate-800 text-brand-500 font-bold' : 'text-slate-600 dark:text-slate-400'"
      >
        <Database class="w-4 h-4" />
        <span>ข้อมูลและการสํารอง</span>
      </button>
    </div>

    <!-- Right Side Config Form Details -->
    <div class="flex-grow overflow-y-auto px-1">
      <!-- 1. Profile Section -->
      <div v-if="activeSection === 'profile'" class="space-y-6">
        <!-- Avatar card -->
        <div class="bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark p-6 rounded-xl space-y-4">
          <h3 class="text-sm font-bold text-slate-800 dark:text-white uppercase tracking-wide">รูปโปรไฟล์ผู้ใช้</h3>
          <div class="flex items-center space-x-5">
            <img
              v-if="authStore.user?.avatar_url"
              :src="authStore.user.avatar_url"
              class="w-16 h-16 rounded-full object-cover border-2 border-brand-500/20"
              alt="Profile"
            />
            <div v-else class="w-16 h-16 rounded-full bg-brand-100 text-brand-500 font-bold text-xl flex items-center justify-center">
              {{ authStore.user?.name.substring(0, 1).toUpperCase() }}
            </div>
            <div>
              <input
                type="file"
                id="avatar-input"
                @change="handleAvatarUpload"
                class="hidden"
                accept="image/*"
              />
              <label
                for="avatar-input"
                class="px-4 py-2 bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 text-xs font-bold rounded-lg cursor-pointer border border-border-light dark:border-border-dark inline-block min-h-[44px]"
              >
                เลือกรูปภาพใหม่
              </label>
              <p class="text-[10px] text-slate-400 mt-1">ประเภทไฟล์ที่รองรับ JPEG, PNG (ขนาดไม่เกิน 2MB)</p>
            </div>
          </div>
        </div>

        <!-- Name & Email update -->
        <div class="bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark p-6 rounded-xl space-y-4">
          <h3 class="text-sm font-bold text-slate-800 dark:text-white uppercase tracking-wide">ข้อมูลส่วนตัว</h3>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-1">
              <label class="text-xs text-slate-400 font-semibold">ชื่อแสดงตัวตน</label>
              <input
                type="text"
                v-model="name"
                class="w-full text-xs bg-slate-50 border border-border-light rounded-lg p-2.5 focus:outline-none dark:bg-slate-900/60 dark:border-border-dark dark:text-white min-h-[44px]"
              />
            </div>
            <div class="space-y-1">
              <label class="text-xs text-slate-400 font-semibold">อีเมลบัญชี</label>
              <input
                type="email"
                v-model="email"
                class="w-full text-xs bg-slate-50 border border-border-light rounded-lg p-2.5 focus:outline-none dark:bg-slate-900/60 dark:border-border-dark dark:text-white min-h-[44px]"
              />
            </div>
          </div>
          <div class="flex justify-end">
            <button
              @click="handleUpdateProfile"
              class="px-4.5 py-2.5 bg-brand-500 hover:bg-brand-600 text-white rounded-lg text-xs font-semibold shadow-md min-h-[44px]"
            >
              บันทึกข้อมูลส่วนตัว
            </button>
          </div>
        </div>

        <!-- Change password -->
        <div class="bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark p-6 rounded-xl space-y-4">
          <h3 class="text-sm font-bold text-slate-800 dark:text-white uppercase tracking-wide flex items-center space-x-1">
            <Key class="w-4.5 h-4.5 text-slate-400" />
            <span>เปลี่ยนรหัสผ่าน</span>
          </h3>
          <div class="space-y-4">
            <div class="space-y-1">
              <label class="text-xs text-slate-400 font-semibold">รหัสผ่านปัจจุบัน</label>
              <input
                type="password"
                v-model="currentPassword"
                placeholder="ป้อนรหัสผ่านเดิม..."
                class="w-full text-xs bg-slate-50 border border-border-light rounded-lg p-2.5 focus:outline-none dark:bg-slate-900/60 dark:border-border-dark dark:text-white min-h-[44px]"
              />
            </div>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div class="space-y-1">
                <label class="text-xs text-slate-400 font-semibold">รหัสผ่านใหม่</label>
                <input
                  type="password"
                  v-model="newPassword"
                  placeholder="อย่างน้อย 6 ตัวอักษร..."
                  class="w-full text-xs bg-slate-50 border border-border-light rounded-lg p-2.5 focus:outline-none dark:bg-slate-900/60 dark:border-border-dark dark:text-white min-h-[44px]"
                />
              </div>
              <div class="space-y-1">
                <label class="text-xs text-slate-400 font-semibold">ยืนยันรหัสผ่านใหม่</label>
                <input
                  type="password"
                  v-model="confirmPassword"
                  placeholder="ป้อนรหัสผ่านใหม่อีกครั้ง..."
                  class="w-full text-xs bg-slate-50 border border-border-light rounded-lg p-2.5 focus:outline-none dark:bg-slate-900/60 dark:border-border-dark dark:text-white min-h-[44px]"
                />
              </div>
            </div>
          </div>
          <div class="flex justify-end">
            <button
              @click="handleChangePassword"
              class="px-4.5 py-2.5 bg-brand-500 hover:bg-brand-600 text-white rounded-lg text-xs font-semibold shadow-md min-h-[44px]"
            >
              เปลี่ยนรหัสผ่านใหม่
            </button>
          </div>
        </div>
      </div>

      <!-- 2. Workspaces Section -->
      <div v-if="activeSection === 'workspaces'" class="space-y-4">
        <!-- New Workspace form toggle -->
        <div class="flex justify-between items-center pb-2 border-b border-border-light dark:border-border-dark">
          <h3 class="text-sm font-bold text-slate-800 dark:text-white uppercase tracking-wide">รายการพื้นที่ทำงาน</h3>
          <button
            @click="showNewWsForm = !showNewWsForm"
            class="px-3.5 py-2 bg-brand-500 hover:bg-brand-600 text-white rounded-lg text-xs font-semibold flex items-center space-x-1 shadow-sm min-h-[44px]"
          >
            <Plus class="w-3.5 h-3.5" />
            <span>สร้างพื้นที่ทำงาน</span>
          </button>
        </div>

        <!-- Create Workspace Form -->
        <div v-if="showNewWsForm" class="bg-slate-50 dark:bg-slate-900/40 p-4 border border-border-light dark:border-border-dark rounded-xl space-y-4">
          <h4 class="text-xs font-bold text-slate-700 dark:text-slate-300">สร้างพื้นที่งานใหม่</h4>
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <input
              type="text"
              v-model="newWsName"
              placeholder="ชื่อพื้นที่งาน... เช่น งานทั่วไป"
              class="col-span-1 sm:col-span-2 text-xs bg-white dark:bg-slate-950 border border-border-light dark:border-border-dark rounded-lg px-2.5 py-1.5 focus:outline-none min-h-[44px]"
            />
            
            <select
              v-model="newWsIcon"
              class="text-xs bg-white dark:bg-slate-950 border border-border-light dark:border-border-dark rounded-lg px-2.5 py-1.5 focus:outline-none min-h-[44px]"
            >
              <option value="briefcase">กระเป๋า (Briefcase)</option>
              <option value="home">บ้าน (Home)</option>
              <option value="heart">หัวใจ (Heart)</option>
            </select>
          </div>

          <!-- Color Swatches -->
          <div class="flex items-center space-x-2">
            <span class="text-[10px] text-slate-400 font-semibold uppercase">สีสัญลักษณ์:</span>
            <div class="flex items-center flex-wrap gap-1.5">
              <button
                v-for="color in presetColors"
                :key="color"
                type="button"
                @click="newWsColor = color"
                class="w-5 h-5 rounded-full border flex items-center justify-center"
                :class="newWsColor === color ? 'border-slate-800 scale-110' : 'border-transparent'"
                :style="{ backgroundColor: color }"
              >
                <Check v-if="newWsColor === color" class="w-3 h-3 text-white" />
              </button>
            </div>
          </div>

          <div class="flex justify-end space-x-2">
            <button
              @click="showNewWsForm = false"
              class="px-3 py-1.5 text-slate-500 text-xs rounded min-h-[44px]"
            >
              ยกเลิก
            </button>
            <button
              @click="handleCreateWorkspace"
              class="px-4.5 py-2.5 bg-brand-500 hover:bg-brand-600 text-white rounded-lg text-xs font-semibold min-h-[44px]"
            >
              สร้างสำเร็จ
            </button>
          </div>
        </div>

        <!-- Workspaces Drag/Shift list -->
        <div class="space-y-3">
          <div
            v-for="ws in workspaceStore.workspaces"
            :key="ws.id"
            class="bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark p-4 rounded-xl flex flex-col sm:flex-row sm:items-center justify-between gap-4"
          >
            <!-- Normal Mode -->
            <div v-if="editingWorkspaceId !== ws.id" class="flex items-center space-x-3 flex-grow">
              <span class="w-3.5 h-3.5 rounded-full" :style="{ backgroundColor: ws.color }"></span>
              <span class="font-semibold text-xs text-slate-900 dark:text-white">{{ ws.name }}</span>
              <span v-if="ws.is_archived" class="text-[9px] font-bold text-slate-400 bg-slate-100 dark:bg-slate-800 px-1.5 py-0.5 rounded">จัดเก็บแล้ว</span>
            </div>

            <!-- Edit Mode -->
            <div v-else class="flex flex-wrap items-center gap-3 flex-grow">
              <input
                type="text"
                v-model="editWsName"
                class="text-xs bg-slate-50 border rounded-lg px-2.5 py-1.5 focus:outline-none dark:bg-slate-900 dark:border-border-dark text-slate-800 dark:text-white min-h-[44px]"
              />
              <select
                v-model="editWsIcon"
                class="text-xs bg-slate-50 border rounded-lg px-2.5 py-1.5 focus:outline-none dark:bg-slate-900 dark:border-border-dark text-slate-800 dark:text-white min-h-[44px]"
              >
                <option value="briefcase">กระเป๋า</option>
                <option value="home">บ้าน</option>
                <option value="heart">หัวใจ</option>
              </select>
              <div class="flex items-center space-x-1.5">
                <button
                  v-for="color in presetColors"
                  :key="color"
                  @click="editWsColor = color"
                  class="w-4.5 h-4.5 rounded-full border"
                  :style="{ backgroundColor: color }"
                ></button>
              </div>
            </div>

            <!-- Action Controls -->
            <div class="flex items-center space-x-2">
              <div v-if="editingWorkspaceId !== ws.id" class="flex items-center space-x-1">
                <!-- Move Up / Down -->
                <button @click="moveWorkspaceOrder(ws.id, 'up')" class="p-2 text-slate-400 hover:text-slate-600 min-h-[36px]">▲</button>
                <button @click="moveWorkspaceOrder(ws.id, 'down')" class="p-2 text-slate-400 hover:text-slate-600 min-h-[36px]">▼</button>
                
                <!-- Edit properties -->
                <button @click="startEditWs(ws)" class="text-xs text-brand-500 hover:text-brand-600 font-bold px-2 py-1 min-h-[44px]">แก้ไข</button>
                
                <!-- Archive/Unarchive -->
                <button @click="toggleArchiveWorkspace(ws)" class="text-xs text-slate-400 hover:text-slate-600 px-2 py-1 min-h-[44px]">
                  {{ ws.is_archived ? 'เลิกจัดเก็บ' : 'จัดเก็บ' }}
                </button>

                <!-- Delete completely -->
                <button @click="handleDeleteWorkspace(ws.id)" class="text-xs text-red-500 hover:text-red-600 font-bold px-2 py-1 min-h-[44px]">ลบ</button>
              </div>
              <div v-else class="flex space-x-1">
                <button @click="editingWorkspaceId = null" class="text-xs text-slate-400 px-2 py-1 min-h-[44px]">ยกเลิก</button>
                <button @click="saveWsEdit(ws.id)" class="text-xs text-brand-500 font-bold px-2 py-1 min-h-[44px]">บันทึก</button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 3. Tags Section -->
      <div v-if="activeSection === 'tags'" class="space-y-4">
        <div class="flex justify-between items-center pb-2 border-b border-border-light dark:border-border-dark">
          <h3 class="text-sm font-bold text-slate-800 dark:text-white uppercase tracking-wide">แท็กทั้งหมด (Tags)</h3>
          <button
            @click="showNewTagForm = !showNewTagForm"
            class="px-3.5 py-2 bg-brand-500 hover:bg-brand-600 text-white rounded-lg text-xs font-semibold flex items-center space-x-1 shadow-sm min-h-[44px]"
          >
            <Plus class="w-3.5 h-3.5" />
            <span>สร้างแท็กใหม่</span>
          </button>
        </div>

        <!-- Create Tag Form -->
        <div v-if="showNewTagForm" class="bg-slate-50 dark:bg-slate-900/40 p-4 border border-border-light dark:border-border-dark rounded-xl space-y-4">
          <h4 class="text-xs font-bold text-slate-700 dark:text-slate-300">สร้างแท็กใหม่</h4>
          <div class="flex items-center space-x-2">
            <input
              type="text"
              v-model="newTagName"
              placeholder="ชื่อแท็ก... เช่น สำคัญด่วน"
              class="flex-grow text-xs bg-white dark:bg-slate-950 border border-border-light dark:border-border-dark rounded-lg px-2.5 py-1.5 focus:outline-none min-h-[44px]"
            />
            
            <div class="flex items-center space-x-1.5 px-2">
              <button
                v-for="color in presetColors.slice(0, 4)"
                :key="color"
                type="button"
                @click="newTagColor = color"
                class="w-5 h-5 rounded-full border"
                :style="{ backgroundColor: color }"
              ></button>
            </div>

            <button
              @click="handleCreateTag"
              class="px-4 py-2 bg-brand-500 hover:bg-brand-600 text-white rounded-lg text-xs font-semibold min-h-[44px]"
            >
              ตกลง
            </button>
          </div>
        </div>

        <!-- Tags List Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div
            v-for="tag in tasksStore.tags"
            :key="tag.id"
            class="bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark p-3.5 rounded-xl flex items-center justify-between gap-3"
          >
            <!-- Normal Mode -->
            <div v-if="editingTagId !== tag.id" class="flex items-center space-x-2.5">
              <span class="w-3 h-3 rounded-full" :style="{ backgroundColor: tag.color }"></span>
              <span class="text-xs font-semibold text-slate-800 dark:text-white">{{ tag.name }}</span>
            </div>

            <!-- Edit Mode -->
            <div v-else class="flex items-center space-x-2 flex-grow">
              <input
                type="text"
                v-model="editTagName"
                class="text-xs bg-slate-50 border rounded px-1.5 py-1 focus:outline-none dark:bg-slate-900 w-24 min-h-[36px]"
              />
              <div class="flex items-center space-x-1">
                <button
                  v-for="color in presetColors.slice(0, 4)"
                  :key="color"
                  @click="editTagColor = color"
                  class="w-3.5 h-3.5 rounded-full"
                  :style="{ backgroundColor: color }"
                ></button>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex items-center space-x-1">
              <div v-if="editingTagId !== tag.id">
                <button @click="startEditTag(tag)" class="text-[11px] text-brand-500 hover:underline px-2 py-1 min-h-[36px]">แก้ไข</button>
                <button @click="handleDeleteTag(tag.id)" class="text-[11px] text-red-500 hover:underline px-2 py-1 min-h-[36px]">ลบ</button>
              </div>
              <div v-else>
                <button @click="editingTagId = null" class="text-[10px] text-slate-400 px-2 py-1 min-h-[36px]">ยกเลิก</button>
                <button @click="saveTagEdit(tag.id)" class="text-[10px] text-brand-500 font-bold px-2 py-1 min-h-[36px]">ตกลง</button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 4. Preferences Section -->
      <div v-if="activeSection === 'preferences'" class="space-y-6">
        <!-- Theme selection -->
        <div class="bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark p-6 rounded-xl space-y-4">
          <h3 class="text-sm font-bold text-slate-800 dark:text-white uppercase tracking-wide">ธีมการแสดงผล (Theme)</h3>
          
          <div class="grid grid-cols-3 gap-3">
            <!-- Light -->
            <button
              @click="changeTheme('light')"
              class="flex flex-col items-center justify-center p-4 rounded-xl border transition-all hover:bg-slate-50 dark:hover:bg-slate-800 min-h-[80px]"
              :class="uiStore.theme === 'light' ? 'border-brand-500 bg-brand-50/20 text-brand-500' : 'border-border-light dark:border-border-dark text-slate-600 dark:text-slate-400'"
            >
              <Sun class="w-5 h-5 mb-1.5" />
              <span class="text-xs font-semibold">สว่าง (Light)</span>
            </button>

            <!-- Dark -->
            <button
              @click="changeTheme('dark')"
              class="flex flex-col items-center justify-center p-4 rounded-xl border transition-all hover:bg-slate-50 dark:hover:bg-slate-800 min-h-[80px]"
              :class="uiStore.theme === 'dark' ? 'border-brand-500 bg-brand-50/20 text-brand-500' : 'border-border-light dark:border-border-dark text-slate-600 dark:text-slate-400'"
            >
              <Moon class="w-5 h-5 mb-1.5" />
              <span class="text-xs font-semibold">มืด (Dark)</span>
            </button>

            <!-- System -->
            <button
              @click="changeTheme('system')"
              class="flex flex-col items-center justify-center p-4 rounded-xl border transition-all hover:bg-slate-50 dark:hover:bg-slate-800 min-h-[80px]"
              :class="uiStore.theme === 'system' ? 'border-brand-500 bg-brand-50/20 text-brand-500' : 'border-border-light dark:border-border-dark text-slate-600 dark:text-slate-400'"
            >
              <Laptop class="w-5 h-5 mb-1.5" />
              <span class="text-xs font-semibold">ระบบ (System)</span>
            </button>
          </div>
        </div>

        <!-- Language toggle -->
        <div class="bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark p-6 rounded-xl space-y-4">
          <h3 class="text-sm font-bold text-slate-800 dark:text-white uppercase tracking-wide">ภาษาของระบบ (Language)</h3>
          <div class="flex space-x-3">
            <button
              @click="uiStore.setLanguage('th')"
              class="px-4 py-2 border rounded-lg text-xs font-semibold min-h-[44px]"
              :class="uiStore.language === 'th' ? 'bg-brand-500 text-white border-transparent' : 'bg-slate-50 border-border-light text-slate-600 dark:bg-slate-900 dark:border-border-dark'"
            >
              ภาษาไทย (Thai)
            </button>
            <button
              @click="uiStore.setLanguage('en')"
              class="px-4 py-2 border rounded-lg text-xs font-semibold min-h-[44px]"
              :class="uiStore.language === 'en' ? 'bg-brand-500 text-white border-transparent' : 'bg-slate-50 border-border-light text-slate-600 dark:bg-slate-900 dark:border-border-dark'"
            >
              English (US)
            </button>
          </div>
        </div>
      </div>

      <!-- 5. Data Section -->
      <div v-if="activeSection === 'data'" class="space-y-6">
        <!-- Backup & Download -->
        <div class="bg-white dark:bg-surface-dark border border-border-light dark:border-border-dark p-6 rounded-xl space-y-4">
          <h3 class="text-sm font-bold text-slate-800 dark:text-white uppercase tracking-wide">การสํารองข้อมูล</h3>
          <p class="text-xs text-slate-500 leading-relaxed">สำรองข้อมูลงานและปฏิทินของคุณเก็บไว้ สามารถดาวน์โหลดเป็นรูปแบบ JSON เพื่อนำไปเก็บสำรองหรือย้ายระบบได้ตลอดเวลา</p>
          <div class="flex items-center space-x-3 pt-2">
            <button
              @click="exportJSON"
              class="flex items-center space-x-1.5 px-4.5 py-2.5 bg-slate-900 text-white dark:bg-slate-800 dark:text-slate-200 border border-transparent rounded-lg text-xs font-semibold hover:bg-slate-800 min-h-[44px]"
            >
              <FileDown class="w-4 h-4" />
              <span>ดาวน์โหลดข้อมูลทั้งหมด (JSON)</span>
            </button>
          </div>
        </div>

        <!-- Danger account wipe -->
        <div class="bg-red-50 border border-red-200 dark:bg-red-950/10 dark:border-red-900/40 p-6 rounded-xl space-y-4">
          <h3 class="text-sm font-bold text-red-600 dark:text-red-400 uppercase tracking-wide">ลบบัญชีและทำลายข้อมูล</h3>
          <p class="text-xs text-red-700 dark:text-red-300 leading-relaxed">คำเตือน! การกดปุ่มด้านล่างจะเริ่มกระบวนการลบชื่อผู้ใช้งาน ข้อมูลงาน ความคิดเห็น แท็ก และพื้นที่งานทั้งหมดของคุณออกจากเซิร์ฟเวอร์แบบถาวร</p>
          <div class="pt-2">
            <button
              @click="handleDeleteAccount"
              class="flex items-center space-x-1.5 px-4.5 py-2.5 bg-red-600 hover:bg-red-700 text-white rounded-lg text-xs font-semibold shadow shadow-red-500/10 min-h-[44px]"
            >
              <Trash2 class="w-4.5 h-4.5" />
              <span>ลบบัญชีของฉันอย่างถาวร</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
