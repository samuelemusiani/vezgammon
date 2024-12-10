<template>
  <div
    class="card glass bg-base-200/80 transition-all duration-300 hover:shadow-lg"
  >
    <div class="card-body">
      <h3 class="card-title text-2xl text-primary">Performance Overview</h3>

      <!-- First Row: Local, CPU, and Tournament Games -->
      <div class="mb-4 grid grid-cols-3 gap-4">
        <StatItem
          title="Local Games"
          :value="stats.local"
          valueClass="text-secondary"
          icon="fas fa-users"
          containerClass="stats bg-base-300/50"
        />
        <StatItem
          title="VS CPU"
          :value="stats.cpu"
          valueClass="text-primary"
          icon="fas fa-robot"
          containerClass="stats bg-base-300/50"
        />
        <StatItem
          title="Tournament"
          :value="stats.tournament"
          valueClass="text-info"
          icon="fas fa-trophy"
          containerClass="stats bg-base-300/50"
        />
      </div>

      <!-- Second Row: Online, Victories, Win Rate, and Elo -->
      <div class="grid grid-cols-2 gap-4">
        <StatItem
          title="Online Games"
          :value="stats.online"
          valueClass="text-accent"
          icon="fas fa-globe"
          containerClass="stats bg-base-300/50"
        />
        <StatItem
          title="Victories"
          :value="stats.win"
          valueClass="text-success"
          icon="fas fa-flag-checkered"
          containerClass="stats bg-base-300/50"
        />
        <StatItem
          title="Win Rate"
          :value="`${stats.winrate}%`"
          valueClass="text-accent"
          icon="fas fa-percentage"
          :description="getTrendDescription"
          containerClass="stats bg-base-300/50"
        />
        <!-- Aggiunto un nuovo StatItem per l'Elo -->
        <StatItem
          title="Elo Rating"
          :value="stats.elo[stats.elo.length - 1]"
          valueClass="text-primary"
          icon="fas fa-chart-line"
          containerClass="stats bg-base-300/50"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import StatItem from './StatItem.vue'
import { computed } from 'vue'

const props = defineProps<{
  stats: {
    win: number
    lost: number
    winrate: number
    cpu: number
    local: number
    online: number
    tournament: number
    elo: number[]
  }
}>()

const getTrendDescription = computed(() => {
  const elo = props.stats.elo
  if (elo.length < 2) return ''

  const lastElo = elo[elo.length - 1]
  const previousElo = elo[elo.length - 2]

  return lastElo > previousElo ? '↗︎ Trending Up' : '↘︎ Trending Down'
})
</script>
