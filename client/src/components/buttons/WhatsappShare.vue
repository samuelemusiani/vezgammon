<template>
  <ShareNetwork network="whatsapp" :url="url" :title="title">
    <button class="btn btn-md bg-green-400 text-xl">
      <i class="fab fa-whatsapp"></i>
      Whatsapp
    </button>
  </ShareNetwork>
</template>

<script setup lang="ts">
defineProps<{
  url: string
  title: string
}>()
</script>
