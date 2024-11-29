<template>
  <button @click="$emit('click')" class="btn btn-primary text-xl">
    Back to Home
  </button>
</template>

<script setup lang="ts">
defineEmits(['click'])
</script>
