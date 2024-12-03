<template>
  <ShareNetwork network="telegram" :url="url" :title="title">
    <button class="btn btn-info btn-md text-xl">
      <i class="fab fa-telegram"></i>
      Telegram
    </button>
  </ShareNetwork>
</template>

<script setup lang="ts">
defineProps<{
  url: string
  title: string
}>()
</script>
