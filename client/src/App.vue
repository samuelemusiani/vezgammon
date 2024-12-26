<script setup lang="ts">
import { useTheme } from '@/composables/useTheme'

useTheme()
import { RouterView } from 'vue-router'
import Rotate from '@/components/Rotate.vue'
</script>

<template>
  <div
    class="h-screen w-screen border-2 border-primary bg-neutral bg-retro-pattern lg:border-8"
  >
    <RouterView />
    <Rotate />
  </div>
</template>
