<template>
  <ShareNetwork
    network="facebook"
    :url="url"
    :title="title"
    :description="description"
    quote="Check out my Backgammon stats!"
  >
    <button class="btn btn-success btn-md text-xl">
      <i class="fab fa-facebook"></i>
      Post
    </button>
  </ShareNetwork>
</template>

<script setup lang="ts">
defineProps<{
  url: string
  title: string
  description?: string
}>()
</script>
