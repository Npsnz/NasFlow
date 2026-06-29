<template>
  <nav class="top-nav" :class="{ 'dark-mode': themeStore.isDark }">
    <!-- Blur backdrop -->
    <div class="nav-backdrop"></div>

    <!-- Nav content -->
    <div class="nav-container">
      <!-- Logo/Brand -->
      <div class="nav-brand">
        <div class="logo">N</div>
        <div>
          <h1>NasFlow</h1>
          <p>Personal Tasks</p>
        </div>
      </div>

      <!-- Center nav items -->
      <div class="nav-center">
        <router-link to="/dashboard" class="nav-link" :class="{ active: isActive('/dashboard') }">
          {{ $t('nav.dashboard') }}
        </router-link>
        <router-link to="/tasks" class="nav-link" :class="{ active: isActive('/tasks') }">
          {{ $t('nav.tasks') }}
        </router-link>
        <router-link to="/settings" class="nav-link" :class="{ active: isActive('/settings') }">
          {{ $t('nav.settings') }}
        </router-link>
      </div>

      <!-- Right items -->
      <div class="nav-right">
        <button class="icon-btn" @click="toggleTheme" :title="themeStore.isDark ? 'Light mode' : 'Dark mode'">
          <span v-if="themeStore.isDark">☀️</span>
          <span v-else>🌙</span>
        </button>
        <select v-model="locale" class="lang-select">
          <option value="en">EN</option>
          <option value="th">TH</option>
        </select>
        <button class="mobile-menu-btn" @click="mobileMenuOpen = !mobileMenuOpen">
          <span v-if="!mobileMenuOpen">☰</span>
          <span v-else>✕</span>
        </button>
        <button class="logout-btn" @click="logout">{{ $t('nav.logout') }}</button>
      </div>

      <!-- Mobile Menu Dropdown -->
      <div v-if="mobileMenuOpen" class="mobile-menu">
        <router-link to="/dashboard" class="mobile-menu-link" :class="{ active: isActive('/dashboard') }" @click="mobileMenuOpen = false">
          {{ $t('nav.dashboard') }}
        </router-link>
        <router-link to="/tasks" class="mobile-menu-link" :class="{ active: isActive('/tasks') }" @click="mobileMenuOpen = false">
          {{ $t('nav.tasks') }}
        </router-link>
        <router-link to="/settings" class="mobile-menu-link" :class="{ active: isActive('/settings') }" @click="mobileMenuOpen = false">
          {{ $t('nav.settings') }}
        </router-link>
        <button class="mobile-menu-logout" @click="logout">{{ $t('nav.logout') }}</button>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { useThemeStore } from '@/stores/theme'

const router = useRouter()
const route = useRoute()
const { locale } = useI18n()
const authStore = useAuthStore()
const toastStore = useToastStore()
const themeStore = useThemeStore()

const mobileMenuOpen = ref(false)

const toggleTheme = () => {
  themeStore.toggleTheme()
}

const isActive = (path: string) => route.path === path

const logout = async () => {
  await authStore.logout()
  toastStore.success('ออกจากระบบสำเร็จ')
  router.push('/login')
}

onMounted(() => {
  // Initialize theme on mount
  const savedTheme = localStorage.getItem('theme') || 'light'
  document.documentElement.setAttribute('data-theme', savedTheme)
})
</script>

<style scoped>
.top-nav {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 80px;
  z-index: 100;
  transition: all 0.3s ease;
  width: 100%;
}

/* Blur backdrop */
.nav-backdrop {
  position: absolute;
  inset: 0;
  background: rgba(250, 248, 243, 0.7);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(232, 221, 210, 0.3);
}

.top-nav.dark-mode .nav-backdrop {
  background: rgba(26, 20, 16, 0.7);
  border-bottom-color: rgba(61, 53, 48, 0.3);
}

.nav-container {
  position: relative;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 2rem;
  gap: 2rem;
  width: 100%;
  box-sizing: border-box;
}

/* Brand */
.nav-brand {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  min-width: fit-content;
}

.logo {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: linear-gradient(135deg, #A89968 0%, #8B7355 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 1.25rem;
}

.top-nav.dark-mode .logo {
  background: linear-gradient(135deg, #D4A574 0%, #A89968 100%);
}

.nav-brand h1 {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  font-family: Georgia, serif;
  color: #3D3D3D;
  letter-spacing: 0.5px;
}

.top-nav.dark-mode .nav-brand h1 {
  color: #F5F1E8;
}

.nav-brand p {
  margin: 0;
  font-size: 0.75rem;
  color: #8B8B8B;
  font-weight: 300;
}

.top-nav.dark-mode .nav-brand p {
  color: #C9C1B8;
}

/* Center nav */
.nav-center {
  flex: 1;
  display: flex;
  justify-content: center;
  gap: 2rem;
}

.nav-link {
  position: relative;
  font-size: 0.95rem;
  font-weight: 400;
  color: #8B8B8B;
  text-decoration: none;
  transition: color 0.3s ease;
  font-family: Inter, sans-serif;
  letter-spacing: 0.3px;
}

.top-nav.dark-mode .nav-link {
  color: #A0AEC0;
}

.nav-link:hover {
  color: #A89968;
}

.top-nav.dark-mode .nav-link:hover {
  color: #F5DEB3;
}

.nav-link.active {
  color: #A89968;
  font-weight: 500;
}

.top-nav.dark-mode .nav-link.active {
  color: #F5DEB3;
}

.nav-link.active::after {
  content: '';
  position: absolute;
  bottom: -8px;
  left: 0;
  right: 0;
  height: 2px;
  background: #A89968;
  border-radius: 1px;
}

.top-nav.dark-mode .nav-link.active::after {
  background: #F5DEB3;
}

/* Right items */
.nav-right {
  display: flex;
  align-items: center;
  gap: 1rem;
  min-width: fit-content;
}

.icon-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: none;
  background: rgba(168, 153, 104, 0.1);
  cursor: pointer;
  font-size: 1.2rem;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.top-nav.dark-mode .icon-btn {
  background: rgba(212, 165, 116, 0.15);
}

.icon-btn:hover {
  background: rgba(168, 153, 104, 0.2);
}

.top-nav.dark-mode .icon-btn:hover {
  background: rgba(212, 165, 116, 0.25);
}

.lang-select {
  padding: 0.5rem 0.75rem;
  border-radius: 6px;
  border: 1px solid #E8DDD2;
  background: white;
  color: #3D3D3D;
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.top-nav.dark-mode .lang-select {
  border-color: #3D3530;
  background: #2D2420;
  color: #F5F1E8;
}

.lang-select:hover {
  border-color: #A89968;
}

.top-nav.dark-mode .lang-select:hover {
  border-color: #F5DEB3;
}

.logout-btn {
  padding: 0.5rem 1.25rem;
  border-radius: 6px;
  border: 1px solid #A89968;
  background: white;
  color: #A89968;
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: Inter, sans-serif;
}

.top-nav.dark-mode .logout-btn {
  border-color: #F5DEB3;
  background: transparent;
  color: #F5DEB3;
}

.logout-btn:hover {
  background: #A89968;
  color: white;
}

.top-nav.dark-mode .logout-btn:hover {
  background: #F5DEB3;
  color: #1A1410;
}

.mobile-menu-logout {
  padding: 0.65rem 1.5rem;
  color: #A89968;
  border: none;
  background: none;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: Inter, sans-serif;
  text-align: left;
  border-top: 1px solid rgba(168, 153, 104, 0.2);
  margin-top: 0.5rem;
}

.top-nav.dark-mode .mobile-menu-logout {
  color: #F5DEB3;
  border-top-color: rgba(212, 165, 116, 0.2);
}

.mobile-menu-logout:hover {
  background: rgba(168, 153, 104, 0.06);
  padding-left: 1.75rem;
}

.top-nav.dark-mode .mobile-menu-logout:hover {
  background: rgba(212, 165, 116, 0.06);
}

.mobile-menu-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: none;
  background: rgba(168, 153, 104, 0.1);
  cursor: pointer;
  font-size: 1.2rem;
  transition: all 0.3s ease;
  display: none;
  align-items: center;
  justify-content: center;
}

.top-nav.dark-mode .mobile-menu-btn {
  background: rgba(212, 165, 116, 0.15);
}

.mobile-menu-btn:hover {
  background: rgba(168, 153, 104, 0.2);
}

.top-nav.dark-mode .mobile-menu-btn:hover {
  background: rgba(212, 165, 116, 0.25);
}

.mobile-menu {
  position: absolute;
  top: 88px;
  right: 0;
  left: 0;
  margin-top: 0.5rem;
  background: white;
  border-bottom: 1px solid #E8DDD2;
  display: flex;
  flex-direction: column;
  padding: 0.5rem 0;
  z-index: 99;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  animation: slideDown 0.2s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.top-nav.dark-mode .mobile-menu {
  background: #2D2420;
  border-bottom-color: #3D3530;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.mobile-menu-link {
  padding: 0.65rem 1.5rem;
  color: #8B8B8B;
  text-decoration: none;
  font-size: 0.9rem;
  transition: all 0.2s ease;
  font-weight: 400;
  border-left: 3px solid transparent;
}

.top-nav.dark-mode .mobile-menu-link {
  color: #A0AEC0;
}

.mobile-menu-link:hover {
  background: rgba(168, 153, 104, 0.06);
  color: #A89968;
  padding-left: 1.75rem;
}

.top-nav.dark-mode .mobile-menu-link:hover {
  background: rgba(212, 165, 116, 0.06);
  color: #F5DEB3;
}

.mobile-menu-link.active {
  background: rgba(168, 153, 104, 0.1);
  color: #A89968;
  font-weight: 500;
  border-left-color: #A89968;
}

.top-nav.dark-mode .mobile-menu-link.active {
  background: rgba(212, 165, 116, 0.1);
  color: #F5DEB3;
  border-left-color: #F5DEB3;
}

/* Mobile */
@media (max-width: 768px) {
  .nav-container {
    padding: 0 0.75rem;
    gap: 0.5rem;
  }

  .nav-center {
    display: none !important;
    flex: 0;
  }

  .nav-brand {
    gap: 0.5rem;
  }

  .logo {
    width: 32px;
    height: 32px;
    font-size: 1rem;
  }

  .nav-brand h1 {
    font-size: 0.8rem;
  }

  .nav-brand p {
    font-size: 0.65rem;
  }

  .nav-right {
    gap: 0.25rem;
  }

  .icon-btn {
    width: 32px;
    height: 32px;
    font-size: 1rem;
  }

  .logout-btn {
    display: none;
  }

  .mobile-menu-btn {
    display: flex;
    width: 32px;
    height: 32px;
    font-size: 1rem;
  }

  .lang-select {
    display: none;
  }
}
</style>
