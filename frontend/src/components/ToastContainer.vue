<template>
  <div class="toast-container">
    <Transition
      v-for="toast in toasts"
      :key="toast.id"
      name="toast-slide"
      @enter="onEnter"
      @leave="onLeave"
    >
      <div :class="['toast', `toast-${toast.type}`]">
        <div class="toast-icon">
          <CheckCircle v-if="toast.type === 'success'" :size="20" />
          <AlertCircle v-else-if="toast.type === 'error'" :size="20" />
          <AlertTriangle v-else-if="toast.type === 'warning'" :size="20" />
          <Info v-else :size="20" />
        </div>
        <div class="toast-message">{{ toast.message }}</div>
        <button class="toast-close" @click="toastStore.removeToast(toast.id)">×</button>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { useToastStore } from '@/stores/toast'
import { computed } from 'vue'
import { CheckCircle, AlertCircle, AlertTriangle, Info } from 'lucide-vue-next'

const toastStore = useToastStore()
const toasts = computed(() => toastStore.toasts)

const onEnter = (el: Element) => {
  const element = el as HTMLElement
  element.style.opacity = '0'
  element.style.transform = 'translateX(400px)'

  requestAnimationFrame(() => {
    element.style.transition = 'all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)'
    element.style.opacity = '1'
    element.style.transform = 'translateX(0)'
  })
}

const onLeave = (el: Element) => {
  const element = el as HTMLElement
  element.style.transition = 'all 0.3s ease-out'
  element.style.opacity = '0'
  element.style.transform = 'translateX(400px)'
}
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 12px;
  pointer-events: none;
}

.toast {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 18px;
  border-radius: 6px;
  backdrop-filter: blur(16px);
  border: 1px solid #E8DDD2;
  font-size: 14px;
  font-weight: 400;
  min-width: 320px;
  max-width: 420px;
  pointer-events: auto;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
  animation: slideIn 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
  font-family: Georgia, serif;
  letter-spacing: 0.3px;
  background: white;
  color: #3D3D3D;
}

/* Success - Green */
.toast-success {
  border-left: 4px solid #D4A574;
}

/* Error - Red */
.toast-error {
  border-left: 4px solid #8B7355;
}

/* Warning - Orange */
.toast-warning {
  border-left: 4px solid #D4A574;
}

/* Info - Blue */
.toast-info {
  border-left: 4px solid #A89968;
}

.toast-icon {
  font-weight: 700;
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}

.toast-success .toast-icon {
  background: #D4A574;
  color: white;
}

.toast-error .toast-icon {
  background: #8B7355;
  color: white;
}

.toast-warning .toast-icon {
  background: #D4A574;
  color: white;
}

.toast-info .toast-icon {
  background: #A89968;
  color: white;
}

.toast-message {
  flex: 1;
  line-height: 1.5;
  letter-spacing: 0.2px;
}

.toast-close {
  flex-shrink: 0;
  background: none;
  border: none;
  color: inherit;
  font-size: 18px;
  cursor: pointer;
  padding: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.5;
  transition: opacity 0.2s ease;
  font-weight: 300;
}

.toast-close:hover {
  opacity: 0.8;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(420px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

:root {
  color-scheme: light dark;
}

@media (prefers-color-scheme: dark) {
  .toast {
    backdrop-filter: blur(20px);
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.3);
    background: #2D2420;
    color: #F5F1E8;
    border-color: #3D3530;
  }

  .toast-success {
    border-left-color: #D4A574;
  }

  .toast-error {
    border-left-color: #B88C6E;
  }

  .toast-warning {
    border-left-color: #D4A574;
  }

  .toast-info {
    border-left-color: #F5DEB3;
  }
}

@media (prefers-reduced-motion: reduce) {
  .toast {
    animation: none;
  }

  .toast-close {
    transition: none;
  }
}

@media (max-width: 640px) {
  .toast-container {
    left: 12px;
    right: 12px;
    top: 12px;
  }

  .toast {
    min-width: auto;
    max-width: 100%;
  }
}
</style>
