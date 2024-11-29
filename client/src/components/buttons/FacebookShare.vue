<template>
  <ShareNetwork
    network="facebook"
    :url="url"
    :title="title"
    :description="description"
    quote="Check out my Backgammon stats!"
  >
    <button class="btn btn-success text-xl">Share on Facebook</button>
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
