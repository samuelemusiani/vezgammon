<template>
  <div class="container mx-auto px-4 py-8">
    <div class="card h-[90vh] overflow-auto bg-base-100 shadow-xl">
      <div class="card-body">
        <h2 class="card-title text-center text-primary">Player Statistics</h2>

        <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <GamePerformanceCard :stats="stats" />
          <RecentGamesCard :games="stats.games_played.slice(0, 5)" />
        </div>

        <div class="divider"></div>

        <EloChart :elo="stats.elo" />

        <div v-if="sharingEnabled" class="card-actions mt-4 justify-center">
          <BackToHomeButton @click="navigateHome" />
          <FacebookShareButton
            :url="gameShareUrl"
            :title="shareTitle"
            :description="shareDescription"
          />
          <TwitterShareButton :url="gameShareUrl" :title="shareTitle" />
        </div>
        <div v-else class="card-actions mt-4 justify-center">
          <BackToHomeButton @click="navigateHome" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, defineProps } from 'vue'
import router from '@/router'

import EloChart from '@/components/stats/EloChart.vue'
import GamePerformanceCard from './GamePerformanceCard.vue'
import RecentGamesCard from './RecentGamesCard.vue'
import BackToHomeButton from '@/components/buttons/BackHome.vue'
import FacebookShareButton from '@/components/buttons/FacebookShare.vue'
import TwitterShareButton from '@/components/buttons/TwitterShare.vue'

import type { GameState } from '@/utils/game/types'
import type { User } from '@/utils/auth/types'

interface GameStats {
  games_played: GameState[]
  win: number
  lost: number
  winrate: number
  elo: number[]
  cpu: number
  local: number
  online: number
  tournament: number
}

const stats = ref<GameStats>({
  games_played: [],
  win: 0,
  lost: 0,
  winrate: 0,
  elo: [],
  cpu: 0,
  local: 0,
  online: 0,
  tournament: 0,
})

const props = withDefaults(
  defineProps<{
    sharingEnabled?: boolean
    username?: string | null
  }>(),
  {
    sharingEnabled: true,
    username: null,
  },
)

const currentUserId = ref<string | null>(null)
const gameShareUrl = ref('')

const shareTitle = computed(() => `Check out my Backgammon stats!`)
const shareDescription = computed(
  () =>
    `Win Rate: ${stats.value.winrate}% | Games Played: ${stats.value.games_played.length || 0}`,
)

async function fetchUserStats() {
  let response
  if (!props.playerId) {
    response = await fetch('/api/stats')
  } else {
    response = await fetch(`/api/player/${props.usernname}`)
  }
  if (!response.ok) {
    throw new Error('Failed to fetch user stats')
  }
  const tmp: GameStats = await response.json()
  stats.value = tmp
}

onMounted(async () => {
  try {
    const userResponse = await fetch('/api/session')
    if (!userResponse.ok) {
      throw new Error('Failed to fetch user session')
    }
    const user: User = await userResponse.json()

    currentUserId.value = user.id || null
    gameShareUrl.value = `${window.location.origin}/player/${currentUserId.value}`
  } catch (error) {
    console.error('Error fetching user info:', error)
  }

  try {
    await fetchUserStats()
  } catch (error) {
    console.error('Error fetching stats:', error)
  }
})

const navigateHome = () => {
  router.push('/')
}
</script>
