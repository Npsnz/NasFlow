<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUIStore } from '@/stores/ui'
import { Eye, EyeOff, User, Mail, Lock, ShieldCheck, Sparkles } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const uiStore = useUIStore()

const activeTab = ref<'login' | 'register'>('login')
const showPassword = ref(false)

// Form fields
const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const rememberMe = ref(true)

const loading = ref(false)
const errors = ref<Record<string, string>>({})

// Real-time password strength check
const passwordStrength = computed(() => {
  const p = password.value
  if (!p) return { score: 0, label: 'ไม่มีความปลอดภัย', color: 'bg-slate-200' }
  let score = 0
  if (p.length >= 6) score++
  if (/[A-Z]/.test(p)) score++
  if (/[0-9]/.test(p)) score++
  if (/[^A-Za-z0-9]/.test(p)) score++

  switch (score) {
    case 4:
      return { score: 4, label: 'แข็งแกร่งมาก', color: 'bg-green-500' }
    case 3:
      return { score: 3, label: 'แข็งแกร่ง', color: 'bg-green-400' }
    case 2:
      return { score: 2, label: 'ปานกลาง', color: 'bg-amber-400' }
    case 1:
    default:
      return { score: 1, label: 'คาดเดาง่าย', color: 'bg-red-500' }
  }
})

const validateForm = () => {
  const errs: Record<string, string> = {}

  if (!email.value) {
    errs.email = 'กรุณาระบุอีเมล'
  } else if (!/\S+@\S+\.\S+/.test(email.value)) {
    errs.email = 'รูปแบบอีเมลไม่ถูกต้อง'
  }

  if (!password.value) {
    errs.password = 'กรุณาระบุรหัสผ่าน'
  }

  if (activeTab.value === 'register') {
    if (!name.value) {
      errs.name = 'กรุณาระบุชื่อแสดงตัวตน'
    }
    if (password.value.length < 6) {
      errs.password = 'รหัสผ่านต้องมีอย่างน้อย 6 ตัวอักษร'
    }
    if (password.value !== confirmPassword.value) {
      errs.confirmPassword = 'รหัสผ่านยืนยันไม่ตรงกัน'
    }
  }

  errors.value = errs
  return Object.keys(errs).length === 0
}

const handleSubmit = async () => {
  if (!validateForm()) return

  loading.value = true
  try {
    if (activeTab.value === 'login') {
      await authStore.login({
        email: email.value,
        password: password.value,
        remember_me: rememberMe.value
      })
      uiStore.showToast('เข้าสู่ระบบสำเร็จ ยินดีต้อนรับ!', 'success')
    } else {
      await authStore.register({
        name: name.value,
        email: email.value,
        password: password.value
      })
      uiStore.showToast('สมัครสมาชิกสำเร็จแล้ว! ระบบสร้างพื้นที่ทำงานตั้งต้นให้คุณเรียบร้อย', 'success')
      // Trigger onboarding tooltip
      localStorage.setItem('tf_onboarding', 'true')
    }
    router.push('/board')
  } catch (err: any) {
    uiStore.showToast(err || 'เกิดข้อผิดพลาดในการตรวจสอบสิทธิ์', 'error')
  } finally {
    loading.value = false
  }
}

const switchTab = (tab: 'login' | 'register') => {
  activeTab.value = tab
  errors.value = {}
  password.value = ''
  confirmPassword.value = ''
}
</script>

<template>
  <div class="min-h-screen bg-gradient-to-tr from-slate-100 via-slate-50 to-brand-100/30 dark:from-bg-dark dark:to-slate-900 flex flex-col justify-center py-12 sm:px-6 lg:px-8 relative select-none">
    <!-- Background blurred decor circles -->
    <div class="absolute w-96 h-96 rounded-full bg-brand-500/10 blur-3xl top-10 left-10 pointer-events-none"></div>
    <div class="absolute w-96 h-96 rounded-full bg-blue-500/10 blur-3xl bottom-10 right-10 pointer-events-none"></div>

    <div class="sm:mx-auto sm:w-full sm:max-w-md text-center z-10">
      <div class="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-brand-500 text-white font-bold text-2xl shadow-xl shadow-brand-500/20 mb-4">
        N
      </div>
      <h2 class="text-3xl font-extrabold tracking-tight text-slate-900 dark:text-white">
        ยินดีต้อนรับสู่ <span class="bg-gradient-to-r from-brand-500 to-brand-700 bg-clip-text text-transparent">NasFlow</span>
      </h2>
      <p class="mt-2 text-xs text-slate-500 dark:text-slate-400">
        จัดการทุกงานสำหรับวันทำงานและชีวิตส่วนตัวของคุณแบบมืออาชีพ
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md z-10 px-4">
      <div class="bg-white/80 dark:bg-surface-dark/80 backdrop-blur-md border border-slate-200 dark:border-border-dark py-8 px-6 shadow-2xl rounded-2xl">
        <!-- Form Tabs Header -->
        <div class="flex border-b border-slate-200 dark:border-border-dark pb-4 mb-6">
          <button
            @click="switchTab('login')"
            class="flex-1 pb-2 text-sm font-semibold transition-all border-b-2 text-center"
            :class="activeTab === 'login'
              ? 'border-brand-500 text-brand-500'
              : 'border-transparent text-slate-400 hover:text-slate-600'"
          >
            เข้าสู่ระบบ (Login)
          </button>
          <button
            @click="switchTab('register')"
            class="flex-1 pb-2 text-sm font-semibold transition-all border-b-2 text-center"
            :class="activeTab === 'register'
              ? 'border-brand-500 text-brand-500'
              : 'border-transparent text-slate-400 hover:text-slate-600'"
          >
            สร้างบัญชี (Register)
          </button>
        </div>

        <!-- Form fields -->
        <form @submit.prevent="handleSubmit" class="space-y-5">
          <!-- Register Name Field -->
          <div v-if="activeTab === 'register'" class="space-y-1">
            <label class="block text-xs font-semibold text-slate-500 dark:text-slate-400">ชื่อผู้ใช้งาน</label>
            <div class="relative">
              <span class="absolute inset-y-0 left-3 flex items-center text-slate-400">
                <User class="w-4.5 h-4.5" />
              </span>
              <input
                type="text"
                v-model="name"
                placeholder="ชื่อของคุณ..."
                class="w-full pl-10 pr-3 py-2 text-sm bg-slate-50 border border-slate-200 rounded-lg focus:outline-none focus:ring-1 focus:ring-brand-500 focus:bg-white dark:bg-slate-900 dark:border-border-dark dark:text-white min-h-[44px]"
                :class="{ 'border-red-400': errors.name }"
              />
            </div>
            <p v-if="errors.name" class="text-[10px] text-red-500 font-semibold">{{ errors.name }}</p>
          </div>

          <!-- Email Field -->
          <div class="space-y-1">
            <label class="block text-xs font-semibold text-slate-500 dark:text-slate-400">อีเมล (Email)</label>
            <div class="relative">
              <span class="absolute inset-y-0 left-3 flex items-center text-slate-400">
                <Mail class="w-4.5 h-4.5" />
              </span>
              <input
                type="email"
                v-model="email"
                placeholder="example@email.com"
                class="w-full pl-10 pr-3 py-2 text-sm bg-slate-50 border border-slate-200 rounded-lg focus:outline-none focus:ring-1 focus:ring-brand-500 focus:bg-white dark:bg-slate-900 dark:border-border-dark dark:text-white min-h-[44px]"
                :class="{ 'border-red-400': errors.email }"
              />
            </div>
            <p v-if="errors.email" class="text-[10px] text-red-500 font-semibold">{{ errors.email }}</p>
          </div>

          <!-- Password Field -->
          <div class="space-y-1">
            <div class="flex justify-between items-center">
              <label class="block text-xs font-semibold text-slate-500 dark:text-slate-400">รหัสผ่าน (Password)</label>
            </div>
            <div class="relative">
              <span class="absolute inset-y-0 left-3 flex items-center text-slate-400">
                <Lock class="w-4.5 h-4.5" />
              </span>
              <input
                :type="showPassword ? 'text' : 'password'"
                v-model="password"
                placeholder="••••••••"
                class="w-full pl-10 pr-10 py-2 text-sm bg-slate-50 border border-slate-200 rounded-lg focus:outline-none focus:ring-1 focus:ring-brand-500 focus:bg-white dark:bg-slate-900 dark:border-border-dark dark:text-white min-h-[44px]"
                :class="{ 'border-red-400': errors.password }"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-3 flex items-center text-slate-400 hover:text-slate-600"
              >
                <EyeOff v-if="showPassword" class="w-4.5 h-4.5" />
                <Eye v-else class="w-4.5 h-4.5" />
              </button>
            </div>
            <p v-if="errors.password" class="text-[10px] text-red-500 font-semibold">{{ errors.password }}</p>
            
            <!-- Real-time Password Strength meter -->
            <div v-if="activeTab === 'register' && password" class="mt-2 space-y-1.5">
              <div class="flex justify-between text-[9px] font-semibold">
                <span class="text-slate-400">ความแข็งแกร่งของรหัสผ่าน:</span>
                <span :class="passwordStrength.score >= 3 ? 'text-green-500' : 'text-amber-500'">{{ passwordStrength.label }}</span>
              </div>
              <div class="h-1.5 w-full bg-slate-100 rounded-full overflow-hidden flex">
                <div
                  v-for="i in 4"
                  :key="i"
                  class="h-full flex-grow border-r border-white last:border-0 transition-all duration-300"
                  :class="i <= passwordStrength.score ? passwordStrength.color : 'bg-transparent'"
                ></div>
              </div>
            </div>
          </div>

          <!-- Confirm Password Field -->
          <div v-if="activeTab === 'register'" class="space-y-1">
            <label class="block text-xs font-semibold text-slate-500 dark:text-slate-400">ยืนยันรหัสผ่าน (Confirm Password)</label>
            <div class="relative">
              <span class="absolute inset-y-0 left-3 flex items-center text-slate-400">
                <ShieldCheck class="w-4.5 h-4.5" />
              </span>
              <input
                type="password"
                v-model="confirmPassword"
                placeholder="••••••••"
                class="w-full pl-10 pr-3 py-2 text-sm bg-slate-50 border border-slate-200 rounded-lg focus:outline-none focus:ring-1 focus:ring-brand-500 focus:bg-white dark:bg-slate-900 dark:border-border-dark dark:text-white min-h-[44px]"
                :class="{ 'border-red-400': errors.confirmPassword }"
              />
            </div>
            <p v-if="errors.confirmPassword" class="text-[10px] text-red-500 font-semibold">{{ errors.confirmPassword }}</p>
          </div>

          <!-- Remember Me Checkbox (only on login) -->
          <div v-if="activeTab === 'login'" class="flex items-center justify-between">
            <div class="flex items-center">
              <input
                id="remember_me"
                type="checkbox"
                v-model="rememberMe"
                class="h-4 w-4 text-brand-500 focus:ring-brand-500 border-slate-300 rounded"
              />
              <label for="remember_me" class="ml-2 block text-xs text-slate-500 dark:text-slate-400 select-none">
                จำฉันไว้ในระบบ (30 วัน)
              </label>
            </div>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full flex justify-center items-center py-2.5 px-4 border border-transparent rounded-lg shadow-md text-sm font-semibold text-white bg-brand-500 hover:bg-brand-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-500 disabled:opacity-50 active:scale-[0.99] transition-transform min-h-[44px]"
          >
            <!-- Spinner -->
            <svg
              v-if="loading"
              class="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <Sparkles v-else class="w-4.5 h-4.5 mr-2 animate-pulse" />
            <span>{{ activeTab === 'login' ? 'เข้าสู่ระบบ (Login)' : 'สร้างบัญชีล็อกอิน' }}</span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>
