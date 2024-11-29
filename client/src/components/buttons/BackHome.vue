<template>
  <button @click="$emit('click')" class="text-xl btn btn-primary">
    Back to Home
  </button>
</template>

<script setup lang="ts">
defineEmits(['click'])
</script>
