<template>
  <div class="settings-zen" :class="{ 'dark-mode': themeStore.isDark }">
    <TopNav />

    <main class="settings-content">
      <h1>Settings</h1>

      <!-- Profile section -->
      <section class="settings-section" v-if="authStore.user">
        <h2>Profile</h2>
        <div class="profile-info">
          <p class="current-user">{{ authStore.user.name }}</p>
          <p class="current-email">{{ authStore.user.email }}</p>
        </div>
        <div class="setting-item">
          <label>Name</label>
          <input v-model="profileForm.name" type="text" placeholder="Your name" />
        </div>
        <div class="setting-item">
          <label>Email</label>
          <input type="email" :value="authStore.user.email" disabled />
        </div>
        <button class="btn-save" @click="handleUpdateProfile" :disabled="profileForm.name === authStore.user.name">
          Save Changes
        </button>
      </section>

      <!-- Appearance section -->
      <section class="settings-section">
        <h2>Appearance</h2>
        <div class="setting-item">
          <label>Theme</label>
          <div class="theme-options">
            <button class="theme-option" :class="{ active: !themeStore.isDark }" @click="handleThemeChange(false)">
              Light
            </button>
            <button class="theme-option" :class="{ active: themeStore.isDark }" @click="handleThemeChange(true)">
              Dark
            </button>
          </div>
        </div>
        <div class="setting-item">
          <label>Color Scheme</label>
          <select>
            <option>Warm Earth</option>
            <option>Cool Minimal</option>
          </select>
        </div>
      </section>

      <!-- Language section -->
      <section class="settings-section">
        <h2>Language</h2>
        <div class="setting-item">
          <label>Preferred Language</label>
          <select v-model="locale">
            <option value="en">English</option>
            <option value="th">ไทย</option>
          </select>
        </div>
      </section>

      <!-- Tags section -->
      <section class="settings-section">
        <h2>Tags</h2>
        <div class="setting-item">
          <label>Create New Tag</label>
          <div class="tag-creator">
            <input v-model="newTag.name" type="text" placeholder="Tag name" />
            <div class="color-picker">
              <input v-model="newTag.color" type="color" />
            </div>
            <button class="btn-create-tag" @click="handleCreateTag" :disabled="!newTag.name">
              Create
            </button>
          </div>
        </div>

        <div class="existing-tags">
          <h3>Existing Tags</h3>
          <div v-if="tasksStore.tags.length === 0" class="no-tags">No tags yet</div>
          <div v-else class="tags-list">
            <div v-for="tag in tasksStore.tags" :key="tag.id" class="tag-item" @click="openEditTag(tag)">
              <div class="tag-preview" :style="{ backgroundColor: tag.color }"></div>
              <span class="tag-name">{{ tag.name }}</span>
              <button class="btn-delete-tag" @click.stop="handleDeleteTag(tag.id)">×</button>
            </div>
          </div>
        </div>
      </section>

      <!-- Danger zone -->
      <section class="settings-section danger">
        <h2>Danger Zone</h2>
        <button class="btn-danger">Delete Account</button>
        <p class="danger-note">This action cannot be undone.</p>
      </section>
    </main>

    <!-- Edit Tag Modal -->
    <div v-if="editingTag" class="modal-overlay" @click.self="editingTag = null">
      <div class="modal-content">
        <div class="modal-header">
          <h2>Edit Tag</h2>
          <button class="close-btn" @click="editingTag = null">×</button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>Tag Name</label>
            <input v-model="editingTag.name" type="text" placeholder="Tag name" />
          </div>

          <div class="form-group">
            <label>Color</label>
            <div class="color-picker">
              <input v-model="editingTag.color" type="color" />
            </div>
          </div>

          <div class="form-actions">
            <button class="btn-cancel" @click="editingTag = null">Cancel</button>
            <button class="btn-save" @click="handleUpdateTag">Save Changes</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useTasksStore } from '@/stores/tasks'
import { useToastStore } from '@/stores/toast'
import { useThemeStore } from '@/stores/theme'
import client from '@/api/client'
import TopNav from '@/components/TopNav.vue'

const { locale } = useI18n()
const authStore = useAuthStore()
const tasksStore = useTasksStore()
const toastStore = useToastStore()
const themeStore = useThemeStore()

// Sync theme changes
watch(() => themeStore.isDark, (isDark) => {
  document.documentElement.setAttribute('data-theme', isDark ? 'dark' : 'light')
})

const profileForm = ref({
  name: '',
})

const newTag = ref({
  name: '',
  color: '#A89968',
})

const editingTag = ref<any>(null)

onMounted(async () => {
  try {
    await tasksStore.fetchTags()
    // Initialize profile form with current user data
    if (authStore.user) {
      profileForm.value.name = authStore.user.name
    }
  } catch (err) {
    console.error('Failed to load tags:', err)
  }
})

const handleThemeChange = (dark: boolean) => {
  themeStore.setTheme(dark)
  toastStore.success(dark ? 'Dark mode enabled' : 'Light mode enabled')
}

const handleUpdateProfile = async () => {
  if (!profileForm.value.name.trim()) {
    toastStore.warning('Name is required')
    return
  }

  try {
    const res = await client.put('/auth/profile', {
      name: profileForm.value.name,
    })

    // Update auth store user
    if (authStore.user) {
      authStore.user.name = res.data.data.name
    }

    toastStore.success('Profile updated!')
  } catch (err: any) {
    toastStore.error(err.response?.data?.error || 'Failed to update profile')
  }
}

const handleCreateTag = async () => {
  if (!newTag.value.name.trim()) {
    toastStore.warning('Tag name is required')
    return
  }

  try {
    await tasksStore.createTag(newTag.value.name, newTag.value.color)
    toastStore.success('Tag created!')
    newTag.value = { name: '', color: '#A89968' }
  } catch (err: any) {
    toastStore.error(err || 'Failed to create tag')
  }
}

const handleDeleteTag = async (tagId: number) => {
  if (!confirm('Delete this tag?')) return

  try {
    await tasksStore.deleteTag(tagId)
    toastStore.success('Tag deleted!')
  } catch (err: any) {
    toastStore.error(err || 'Failed to delete tag')
  }
}

const openEditTag = (tag: any) => {
  editingTag.value = { ...tag }
}

const handleUpdateTag = async () => {
  if (!editingTag.value) return
  if (!editingTag.value.name.trim()) {
    toastStore.warning('Tag name is required')
    return
  }

  try {
    const res = await client.put(`/tags/${editingTag.value.id}`, {
      name: editingTag.value.name,
      color: editingTag.value.color,
    })

    // Update local store
    const idx = tasksStore.tags.findIndex(t => t.id === editingTag.value.id)
    if (idx > -1) {
      tasksStore.tags[idx] = res.data.data
    }

    toastStore.success('Tag updated!')
    editingTag.value = null
  } catch (err: any) {
    toastStore.error(err.response?.data?.error || 'Failed to update tag')
  }
}
</script>

<style scoped>
.settings-zen {
  min-height: 100vh;
  background: #FAF8F3;
  padding-top: 80px;
  transition: background 0.3s ease;
}

.settings-zen.dark-mode {
  background: #1A1410;
}

.settings-content {
  max-width: 700px;
  margin: 0 auto;
  padding: 3rem 2rem;
  box-sizing: border-box;
}

.settings-zen h1 {
  margin: 0 0 3rem;
  font-size: 2rem;
  font-family: Georgia, serif;
  font-weight: 300;
  color: #3D3D3D;
}

.settings-zen.dark-mode h1 {
  color: #F5F1E8;
}

/* Sections */
.settings-section {
  margin-bottom: 3rem;
  padding: 2rem;
  background: white;
  border-radius: 12px;
  border: 1px solid #E8DDD2;
  transition: all 0.3s ease;
}

.settings-zen.dark-mode .settings-section {
  background: #2D2420;
  border-color: #3D3530;
}

.settings-section.danger {
  border-color: #EF4444;
}

.settings-zen.dark-mode .settings-section.danger {
  border-color: #DC2626;
}

.settings-section h2 {
  margin: 0 0 1.5rem;
  font-size: 1.25rem;
  font-family: Georgia, serif;
  font-weight: 400;
  color: #3D3D3D;
}

.settings-zen.dark-mode .settings-section h2 {
  color: #F5F1E8;
}

/* Profile info */
.profile-info {
  margin-bottom: 2rem;
  padding: 1rem;
  background: rgba(168, 153, 104, 0.05);
  border-radius: 8px;
  border-left: 3px solid #A89968;
}

.settings-zen.dark-mode .profile-info {
  background: rgba(212, 165, 116, 0.08);
  border-left-color: #F5DEB3;
}

.current-user {
  margin: 0 0 0.25rem;
  font-size: 1rem;
  font-weight: 600;
  color: #3D3D3D;
}

.settings-zen.dark-mode .current-user {
  color: #F5F1E8;
}

.current-email {
  margin: 0;
  font-size: 0.9rem;
  color: #8B8B8B;
}

.settings-zen.dark-mode .current-email {
  color: #C9C1B8;
}

/* Setting items */
.setting-item {
  margin-bottom: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.setting-item label {
  font-weight: 500;
  font-size: 0.9rem;
  color: #3D3D3D;
}

.settings-zen.dark-mode .setting-item label {
  color: #F5F1E8;
}

.setting-item input,
.setting-item select {
  padding: 0.75rem;
  border: 1px solid #E8DDD2;
  border-radius: 8px;
  font-size: 0.95rem;
  font-family: Inter, sans-serif;
  background: white;
  color: #3D3D3D;
  transition: all 0.3s ease;
}

.settings-zen.dark-mode .setting-item input,
.settings-zen.dark-mode .setting-item select {
  background: #1A1410;
  border-color: #3D3530;
  color: #F5F1E8;
}

.setting-item input:focus,
.setting-item select:focus {
  outline: none;
  border-color: #A89968;
  box-shadow: 0 0 0 3px rgba(168, 153, 104, 0.1);
}

.settings-zen.dark-mode .setting-item input:focus,
.settings-zen.dark-mode .setting-item select:focus {
  border-color: #F5DEB3;
  box-shadow: 0 0 0 3px rgba(245, 222, 179, 0.1);
}

/* Theme options */
.theme-options {
  display: flex;
  gap: 1rem;
}

.theme-option {
  flex: 1;
  padding: 0.75rem;
  background: #FAF8F3;
  border: 2px solid #E8DDD2;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
  color: #8B8B8B;
}

.settings-zen.dark-mode .theme-option {
  background: #1A1410;
  border-color: #3D3530;
  color: #C9C1B8;
}

.theme-option.active {
  background: #A89968;
  border-color: #A89968;
  color: white;
}

.settings-zen.dark-mode .theme-option.active {
  background: #D4A574;
  border-color: #D4A574;
  color: #1A1410;
}

/* Buttons */
.btn-save {
  padding: 0.75rem 2rem;
  background: #A89968;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-save:hover {
  background: #8B7355;
}

.settings-zen.dark-mode .btn-save {
  background: #D4A574;
  color: #1A1410;
}

.settings-zen.dark-mode .btn-save:hover {
  background: #F5DEB3;
}

.btn-danger {
  padding: 0.75rem 2rem;
  background: #FEE2E2;
  color: #DC2626;
  border: 1px solid #FECACA;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-danger:hover {
  background: #FCA5A5;
  color: white;
}

.settings-zen.dark-mode .btn-danger {
  background: rgba(220, 38, 38, 0.1);
  border-color: rgba(220, 38, 38, 0.3);
  color: #FCA5A5;
}

.settings-zen.dark-mode .btn-danger:hover {
  background: rgba(220, 38, 38, 0.2);
}

.danger-note {
  margin: 0.5rem 0 0;
  font-size: 0.85rem;
  color: #8B8B8B;
}

.settings-zen.dark-mode .danger-note {
  color: #C9C1B8;
}

/* Tag Creator */
.tag-creator {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}

.tag-creator input {
  flex: 1;
  padding: 0.75rem;
  border: 1px solid #E8DDD2;
  border-radius: 6px;
  font-size: 0.95rem;
}

.settings-zen.dark-mode .tag-creator input {
  background: #1A1410;
  border-color: #3D3530;
  color: #F5F1E8;
}

.color-picker {
  display: flex;
  align-items: center;
}

.color-picker input {
  width: 50px;
  height: 50px;
  border: 1px solid #E8DDD2;
  border-radius: 6px;
  cursor: pointer;
  padding: 0;
  background: white;
}

.settings-zen.dark-mode .color-picker input {
  border-color: #3D3530;
  background: #3D3530;
}

.btn-create-tag {
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #A89968 0%, #8B7355 100%);
  color: white;
  border: none;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-create-tag:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(168, 153, 104, 0.3);
}

.btn-create-tag:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Existing Tags */
.existing-tags {
  margin-top: 2rem;
}

.existing-tags h3 {
  margin: 0 0 1rem;
  font-size: 1rem;
  color: #3D3D3D;
}

.settings-zen.dark-mode .existing-tags h3 {
  color: #F5F1E8;
}

.no-tags {
  color: #8B8B8B;
  font-style: italic;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.tag-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: white;
  border: 1px solid #E8DDD2;
  border-radius: 20px;
  transition: all 0.2s;
}

.settings-zen.dark-mode .tag-item {
  background: #2D2420;
  border-color: #3D3530;
}

.tag-preview {
  width: 16px;
  height: 16px;
  border-radius: 50%;
}

.tag-name {
  font-size: 0.9rem;
  font-weight: 500;
  color: #3D3D3D;
}

.settings-zen.dark-mode .tag-name {
  color: #F5F1E8;
}

.btn-delete-tag {
  background: none;
  border: none;
  color: #8B8B8B;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.2s;
}

.btn-delete-tag:hover {
  color: #8B453D;
}

/* Edit Modal */
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
  max-width: 400px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.settings-zen.dark-mode .modal-content {
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

.settings-zen.dark-mode .modal-header h2 {
  color: #F5F1E8;
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

.settings-zen.dark-mode .close-btn:hover {
  color: #F5F1E8;
}

.modal-body {
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
  font-weight: 600;
  font-size: 0.9rem;
  color: #3D3D3D;
}

.settings-zen.dark-mode .form-group label {
  color: #F5F1E8;
}

.form-group input {
  padding: 0.75rem;
  border: 1px solid #E8DDD2;
  border-radius: 6px;
  font-size: 0.95rem;
  background: white;
  color: #3D3D3D;
}

.settings-zen.dark-mode .form-group input {
  background: #1A1410;
  border-color: #3D3530;
  color: #F5F1E8;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.btn-cancel,
.btn-save {
  flex: 1;
  padding: 0.75rem;
  border: none;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel {
  background: #E8DDD2;
  color: #3D3D3D;
}

.btn-cancel:hover {
  background: #D9CEC0;
}

.settings-zen.dark-mode .btn-cancel {
  background: #3D3530;
  color: #F5F1E8;
}

.settings-zen.dark-mode .btn-cancel:hover {
  background: #4D4540;
}

.btn-save {
  background: linear-gradient(135deg, #A89968 0%, #8B7355 100%);
  color: white;
}

.btn-save:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(168, 153, 104, 0.3);
}

.settings-zen.dark-mode .btn-save {
  background: linear-gradient(135deg, #D4A574 0%, #A89968 100%);
  color: #1A1410;
}

@media (max-width: 768px) {
  .settings-content {
    padding: 1.5rem 1rem;
  }

  .settings-section {
    padding: 1.5rem;
  }

  .tag-creator {
    flex-direction: column;
    gap: 0.5rem;
  }

  .tag-creator input {
    width: 100%;
  }

  .color-picker {
    width: 100%;
  }

  .color-picker input {
    width: 50px;
    height: 50px;
  }

  .btn-create-tag {
    width: 100%;
  }
}
</style>
