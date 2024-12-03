<template>
  <ShareNetwork network="reddit" :url="url" :title="title">
    <button class="btn btn-md bg-orange-600 text-xl">
      <i class="fab fa-reddit"></i>
      Reddit
    </button>
  </ShareNetwork>
</template>

<script setup lang="ts">
defineProps<{
  url: string
  title: string
}>()
</script>
