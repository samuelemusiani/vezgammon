<template>
  <div class="flex h-full items-center justify-center overflow-auto">
    <div class="card h-[94%] w-4/5 overflow-y-auto bg-base-100 shadow-xl">
      <div class="card-body w-full">
        <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <GamePerformanceCard :stats="stats" />
          <RecentGamesCard
            :games="stats.games_played"
            :currentUser="currentUsername as string"
            :can-analyze="!generalStats"
          />
        </div>

        <div class="card mt-2 bg-base-200/80 shadow-lg">
          <EloChart :elo="stats.elo" />
        </div>

        <div v-if="!generalStats" class="card-actions mt-2 justify-center">
          <TelegramShareButton :url="gameShareUrl" :title="shareTitle" />
          <RedditShareButton :url="gameShareUrl" :title="shareTitle" />
          <BackToHomeButton @click="navigateHome" />
          <TwitterShareButton :url="gameShareUrl" :title="shareTitle" />
          <WhatsappShareButton :url="gameShareUrl" :title="shareTitle" />
        </div>
        <div v-else class="card-actions mt-2 justify-center">
          <BackToHomeButton @click="navigateHome" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import router from '@/router'

import EloChart from '@/components/stats/EloChart.vue'
import GamePerformanceCard from './GamePerformanceCard.vue'
import RecentGamesCard from './RecentGamesCard.vue'
import BackToHomeButton from '@/components/buttons/BackHome.vue'
import TwitterShareButton from '@/components/buttons/TwitterShare.vue'
import RedditShareButton from '@/components/buttons/RedditShare.vue'
import WhatsappShareButton from '@/components/buttons/WhatsappShare.vue'
import TelegramShareButton from '@/components/buttons/TelegramShare.vue'

import type { GameStats, User } from '@/utils/types'

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
  leaderboard: [],
})

const props = withDefaults(
  defineProps<{
    generalStats?: boolean
    username?: string | null
  }>(),
  {
    generalStats: false,
    username: null,
  },
)

const currentUsername = ref<string | null>(null)
const gameShareUrl = ref('')

const shareTitle = computed(() => `Check out my Backgammon stats!`)

async function fetchUserStats() {
  let response
  if (!props.username) {
    response = await fetch('/api/stats')
  } else {
    response = await fetch(`/api/player/${props.username}`)
  }
  if (!response.ok) {
    throw new Error('Failed to fetch user stats')
  }
  const tmp: GameStats = await response.json()
  stats.value = tmp
}

async function fetchUser() {
  const response = await fetch('/api/session')
  if (!response.ok) {
    throw new Error('Failed to fetch user')
  }
  const user: User = await response.json()
  currentUsername.value = user.username
}

onMounted(async () => {
  try {
    if (props.username) {
      gameShareUrl.value = `${window.location.origin}/player/${props.username}`
    } else {
      await fetchUser()

      gameShareUrl.value = `${window.location.origin}/player/${currentUsername.value}`
    }
  } catch (error) {
    console.error('Error fetching user:', error)
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
