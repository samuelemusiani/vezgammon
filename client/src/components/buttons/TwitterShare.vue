<template>
  <ShareNetwork network="twitter" :url="url" :title="title">
    <button class="btn btn-info btn-md text-xl">
      <i class="fab fa-twitter"></i>
      Tweet
    </button>
  </ShareNetwork>
</template>

<script setup lang="ts">
defineProps<{
  url: string
  title: string
}>()
</script>
