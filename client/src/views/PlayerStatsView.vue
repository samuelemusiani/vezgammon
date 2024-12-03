<script setup lang="ts">
import PlayerStats from '@/components/stats/PlayerStats.vue'
import { useRoute } from 'vue-router'

// get player id from url
const route = useRoute()
const username = route.params.username as string
</script>
<template>
  <PlayerStats :sharing-enabled="false" :username="username" />
</template>
