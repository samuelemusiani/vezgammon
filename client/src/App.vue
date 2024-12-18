<script setup lang="ts">
import { RouterView } from 'vue-router'
import Rotate from '@/components/Rotate.vue'
</script>

<template>
  <div
    class="h-screen w-screen border-8 border-primary bg-neutral bg-retro-pattern"
  >
    <RouterView />
    <Rotate />
  </div>
</template>
