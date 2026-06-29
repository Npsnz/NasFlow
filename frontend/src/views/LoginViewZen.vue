<template>
  <div class="login-zen" :class="{ 'dark-mode': themeStore.isDark }">
    <!-- Background decoration -->
    <div class="bg-decoration"></div>

    <!-- Content -->
    <div class="login-container">
      <div class="login-card">
        <!-- Logo -->
        <div class="logo-section">
          <div class="logo">N</div>
          <h1>NasFlow</h1>
          <p>Personal Task Management</p>
        </div>

        <!-- Tabs -->
        <div class="tabs">
          <button class="tab" :class="{ active: activeTab === 'login' }" @click="activeTab = 'login'">
            Sign In
          </button>
          <button class="tab" :class="{ active: activeTab === 'register' }" @click="activeTab = 'register'">
            Sign Up
          </button>
        </div>

        <!-- Login form -->
        <form v-if="activeTab === 'login'" @submit.prevent="handleLogin" class="form">
          <div class="form-group">
            <label>Email</label>
            <input v-model="form.email" type="email" placeholder="your@email.com" required :disabled="loading" />
          </div>
          <div class="form-group">
            <label>Password</label>
            <input v-model="form.password" type="password" placeholder="••••••••" required :disabled="loading" />
          </div>
          <button type="submit" class="btn-submit" :disabled="loading">{{ loading ? 'ลงชื่นเข้าใช้...' : 'Sign In' }}</button>
        </form>

        <!-- Register form -->
        <form v-if="activeTab === 'register'" @submit.prevent="handleRegister" class="form">
          <div class="form-group">
            <label>Name</label>
            <input v-model="form.name" type="text" placeholder="Your name" required :disabled="loading" />
          </div>
          <div class="form-group">
            <label>Email</label>
            <input v-model="form.email" type="email" placeholder="your@email.com" required :disabled="loading" />
          </div>
          <div class="form-group">
            <label>Password</label>
            <input v-model="form.password" type="password" placeholder="••••••••" required :disabled="loading" />
          </div>
          <button type="submit" class="btn-submit" :disabled="loading">{{ loading ? 'สมัครสมาชิก...' : 'Create Account' }}</button>
        </form>

        <!-- Theme toggle -->
        <div class="theme-toggle">
          <button @click="toggleTheme" class="theme-btn" :title="themeStore.isDark ? 'Light mode' : 'Dark mode'">
            <span v-if="themeStore.isDark">☀️</span>
            <span v-else>🌙</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { useThemeStore } from '@/stores/theme'

const router = useRouter()
const authStore = useAuthStore()
const toastStore = useToastStore()
const themeStore = useThemeStore()
const activeTab = ref<'login' | 'register'>('login')
const loading = ref(false)
const form = ref({
  name: '',
  email: '',
  password: '',
})

const toggleTheme = () => {
  themeStore.toggleTheme()
}

const handleLogin = async () => {
  if (!form.value.email || !form.value.password) {
    toastStore.warning('กรุณากรอกอีเมลและรหัสผ่าน')
    return
  }

  loading.value = true
  try {
    await authStore.login({
      email: form.value.email,
      password: form.value.password,
    })
    toastStore.success('ลงชื่นเข้าใช้สำเร็จ')
    router.push('/dashboard')
  } catch (err: any) {
    const errorMsg = authStore.error || 'ไม่สามารถลงชื่นเข้าใช้ได้'
    toastStore.error(errorMsg)
  } finally {
    loading.value = false
  }
}

const handleRegister = async () => {
  if (!form.value.name || !form.value.email || !form.value.password) {
    toastStore.warning('กรุณากรอกข้อมูลทั้งหมด')
    return
  }

  if (form.value.password.length < 6) {
    toastStore.warning('รหัสผ่านต้องมีความยาวอย่างน้อย 6 ตัวอักษร')
    return
  }

  loading.value = true
  try {
    await authStore.register({
      name: form.value.name,
      email: form.value.email,
      password: form.value.password,
    })
    toastStore.success('สมัครสมาชิกสำเร็จ ยินดีต้อนรับ!')
    form.value = { name: '', email: '', password: '' }
    router.push('/dashboard')
  } catch (err: any) {
    const errorMsg = authStore.error || 'ไม่สามารถสมัครสมาชิกได้'
    toastStore.error(errorMsg)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-zen {
  min-height: 100vh;
  background: linear-gradient(135deg, #FAF8F3 0%, #F5F1E8 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
  transition: background 0.3s ease;
}

.login-zen.dark-mode {
  background: linear-gradient(135deg, #1A1410 0%, #2D2420 100%);
}

.bg-decoration {
  position: fixed;
  inset: 0;
  opacity: 0.5;
  pointer-events: none;
  z-index: 0;
}

.login-zen.dark-mode .bg-decoration {
  opacity: 0.2;
}

.login-container {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 420px;
}

.login-card {
  background: white;
  border-radius: 16px;
  padding: 3rem;
  border: 1px solid #E8DDD2;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.login-zen.dark-mode .login-card {
  background: #2D2420;
  border-color: #3D3530;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

/* Logo section */
.logo-section {
  text-align: center;
  margin-bottom: 2rem;
}

.logo {
  width: 56px;
  height: 56px;
  margin: 0 auto 1rem;
  background: linear-gradient(135deg, #A89968 0%, #8B7355 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 1.75rem;
  font-weight: 700;
  box-shadow: 0 4px 12px rgba(168, 153, 104, 0.2);
}

.login-zen.dark-mode .logo {
  background: linear-gradient(135deg, #D4A574 0%, #A89968 100%);
  box-shadow: 0 4px 12px rgba(212, 165, 116, 0.2);
}

.logo-section h1 {
  margin: 0 0 0.5rem;
  font-size: 1.75rem;
  font-family: Georgia, serif;
  font-weight: 400;
  color: #3D3D3D;
  letter-spacing: 0.5px;
}

.login-zen.dark-mode .logo-section h1 {
  color: #F5F1E8;
}

.logo-section p {
  margin: 0;
  font-size: 0.85rem;
  color: #8B8B8B;
  font-weight: 300;
}

.login-zen.dark-mode .logo-section p {
  color: #C9C1B8;
}

/* Tabs */
.tabs {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 2rem;
  border-bottom: 1px solid #E8DDD2;
}

.login-zen.dark-mode .tabs {
  border-bottom-color: #3D3530;
}

.tab {
  flex: 1;
  padding: 0.75rem;
  background: none;
  border: none;
  color: #8B8B8B;
  font-weight: 500;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.3s ease;
  font-size: 0.9rem;
}

.login-zen.dark-mode .tab {
  color: #C9C1B8;
}

.tab.active {
  color: #A89968;
  border-bottom-color: #A89968;
}

.login-zen.dark-mode .tab.active {
  color: #F5DEB3;
  border-bottom-color: #F5DEB3;
}

/* Form */
.form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  margin-bottom: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-size: 0.85rem;
  font-weight: 500;
  color: #3D3D3D;
  letter-spacing: 0.3px;
}

.login-zen.dark-mode .form-group label {
  color: #F5F1E8;
}

.form-group input {
  padding: 0.75rem 1rem;
  border: 1px solid #E8DDD2;
  border-radius: 8px;
  font-size: 0.95rem;
  font-family: Inter, sans-serif;
  transition: all 0.3s ease;
  background: white;
  color: #3D3D3D;
}

.login-zen.dark-mode .form-group input {
  background: #1A1410;
  border-color: #3D3530;
  color: #F5F1E8;
}

.form-group input:focus {
  outline: none;
  border-color: #A89968;
  box-shadow: 0 0 0 3px rgba(168, 153, 104, 0.1);
}

.login-zen.dark-mode .form-group input:focus {
  border-color: #F5DEB3;
  box-shadow: 0 0 0 3px rgba(245, 222, 179, 0.1);
}

.form-group input::placeholder {
  color: #D9CEC0;
}

.login-zen.dark-mode .form-group input::placeholder {
  color: #5C5450;
}

/* Submit button */
.btn-submit {
  padding: 0.85rem;
  background: linear-gradient(135deg, #A89968 0%, #8B7355 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 0.95rem;
  letter-spacing: 0.5px;
}

.btn-submit:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(168, 153, 104, 0.3);
}

.btn-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.login-zen.dark-mode .btn-submit {
  background: linear-gradient(135deg, #D4A574 0%, #A89968 100%);
  color: #1A1410;
}

/* Theme toggle */
.theme-toggle {
  text-align: center;
  margin-top: 1.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid #E8DDD2;
}

.login-zen.dark-mode .theme-toggle {
  border-top-color: #3D3530;
}

.theme-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 1px solid #E8DDD2;
  background: white;
  cursor: pointer;
  font-size: 1.2rem;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.login-zen.dark-mode .theme-btn {
  border-color: #3D3530;
  background: #1A1410;
}

.theme-btn:hover {
  border-color: #A89968;
  background: rgba(168, 153, 104, 0.1);
}

.login-zen.dark-mode .theme-btn:hover {
  border-color: #F5DEB3;
  background: rgba(245, 222, 179, 0.1);
}

@media (max-width: 480px) {
  .login-card {
    padding: 2rem;
  }

  .logo-section h1 {
    font-size: 1.5rem;
  }
}
</style>
