<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useTheme } from '@/composables/useTheme'

const { currentTheme, themeOptions, changeTheme } = useTheme()
</script>

<template>
  <div class="bg-neutral bg-retro-pattern border-8 border-primary h-screen w-screen">
    <RouterView />
  </div>
</template>
