<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useTheme } from '@/composables/useTheme'

const { currentTheme, themeOptions, changeTheme } = useTheme()
</script>

<template>
  <div
    class="h-screen w-screen border-8 border-primary bg-neutral bg-retro-pattern"
  >
    <RouterView />
  </div>
</template>
