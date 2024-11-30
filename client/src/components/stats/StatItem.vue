<template>
  <div :class="containerClass || 'stat'">
    <div class="stat">
      <div v-if="icon" class="stat-figure" :class="valueClass">
        <i :class="[icon, 'text-2xl']"></i>
      </div>
      <div class="stat-title">{{ title }}</div>
      <div class="stat-value" :class="valueClass">{{ value }}</div>
      <div v-if="description" class="stat-desc">{{ description }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  title: {
    type: String,
    required: true,
  },
  value: {
    type: [String, Number],
    required: true,
  },
  valueClass: {
    type: String,
    default: '',
  },
  description: {
    type: String,
    default: '',
  },
  icon: {
    type: String,
    default: '',
  },
  containerClass: {
    type: String,
    default: '',
  },
})
</script>
