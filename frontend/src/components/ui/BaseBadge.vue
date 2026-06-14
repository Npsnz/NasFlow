<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    color?: string
    variant?: 'solid' | 'subtle'
  }>(),
  {
    color: '#171717',
    variant: 'subtle'
  }
)

const styles = computed(() => {
  let hex = props.color.startsWith('#') ? props.color.replace('#', '') : '171717'
  if (hex.length !== 6) {
    hex = '171717'
  }
  const r = parseInt(hex.substring(0, 2), 16)
  const g = parseInt(hex.substring(2, 4), 16)
  const b = parseInt(hex.substring(4, 6), 16)

  if (props.variant === 'subtle') {
    return {
      backgroundColor: `rgba(${r}, ${g}, ${b}, 0.1)`,
      color: `#${hex}`,
      border: `1px solid rgba(${r}, ${g}, ${b}, 0.2)`
    }
  } else {
    return {
      backgroundColor: `#${hex}`,
      color: '#ffffff'
    }
  }
})
</script>

<template>
  <span
    class="inline-flex items-center px-2 py-0.5 text-xs font-medium rounded-full select-none"
    :style="styles"
  >
    <slot />
  </span>
</template>
