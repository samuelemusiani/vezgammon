<template>
  <div class="stat place-items-center">
    <div class="stat-title">{{ title }}</div>
    <div class="stat-value">{{ value }}</div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  value: string | number
}>()
</script>
