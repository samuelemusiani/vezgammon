<template>
  <div class="card bg-base-200">
    <div class="card-body">
      <h3 class="card-title">Game Performance</h3>
      <div class="grid grid-cols-2 gap-4">
        <div class="grid grid-cols-2 gap-2">
          <StatItem title="Played" :value="0" />
          <StatItem title="Wins" :value="stats.win" />
          <StatItem title="Lost" :value="stats.lost" />
          <StatItem title="Win Rate" :value="`${stats.winrate}%`" />
        </div>
        <div class="grid grid-cols-2 gap-2">
          <StatItem title="CPU" :value="stats.cpu" />
          <StatItem title="Local" :value="stats.local" />
          <StatItem title="Online" :value="stats.online" />
          <StatItem title="Tournament" :value="stats.tournament" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import StatItem from './StatItem.vue'

defineProps<{
  stats: {
    win: number
    lost: number
    winrate: number
    cpu: number
    local: number
    online: number
    tournament: number
  }
}>()
</script>
