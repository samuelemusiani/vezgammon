<template>
  <PlayerStats :sharing-enabled="false" :player-id="playerId" />
</template>

<script setup lang="ts">
import PlayerStats from '@/components/stats/PlayerStats.vue'
import { useRoute } from 'vue-router'

// get player id from url
const route = useRoute()
const playerId = route.params.username as string

</script>
