<template>
  <ShareNetwork
    network="facebook"
    :url="url"
    :title="title"
    :description="description"
    quote="Check out my Backgammon stats!"
  >
    <button class="text-xl btn btn-success">Share on Facebook</button>
  </ShareNetwork>
</template>

<script setup lang="ts">
import { ShareNetwork } from 'vue-social-sharing'

defineProps<{
  url: string
  title: string
  description?: string
}>()
</script>
