<template>
  <div class="container py-8 px-4 mx-auto">
    <div class="overflow-auto shadow-xl card h-[90vh] bg-base-100">
      <div class="card-body">
        <h2 class="text-center card-title text-primary">Player Statistics</h2>

        <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <!-- Overall Stats -->
          <div class="card bg-base-200">
            <div class="card-body">
              <h3 class="card-title">Game Performance</h3>
              <div class="grid grid-cols-2 gap-4">
                <div class="grid grid-cols-2 gap-2">
                  <div class="place-items-center stat">
                    <div class="stat-title">Played</div>
                    <div class="stat-value">{{ 0 }}</div>
                  </div>
                  <div class="place-items-center stat">
                    <div class="stat-title">Wins</div>
                    <div class="stat-value">{{ stats.win }}</div>
                  </div>
                  <div class="place-items-center stat">
                    <div class="stat-title">Lost</div>
                    <div class="stat-value">{{ stats.lost }}</div>
                  </div>
                  <div class="place-items-center stat">
                    <div class="stat-title">Win Rate</div>
                    <div class="stat-value">{{ stats.winrate }}%</div>
                  </div>
                </div>
                <div class="grid grid-cols-2 gap-2">
                  <div class="place-items-center stat">
                    <div class="stat-title">CPU</div>
                    <div class="stat-value">{{ stats.cpu }}</div>
                  </div>
                  <div class="place-items-center stat">
                    <div class="stat-title">Local</div>
                    <div class="stat-value">{{ stats.local }}</div>
                  </div>
                  <div class="place-items-center stat">
                    <div class="stat-title">Online</div>
                    <div class="stat-value">{{ stats.online }}</div>
                  </div>
                  <div class="place-items-center stat">
                    <div class="stat-title">Tournament</div>
                    <div class="stat-value">{{ stats.tournament }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Recent Games -->
          <div class="card bg-base-200">
            <div class="card-body">
              <h3 class="card-title">Recent Games</h3>
              <div class="overflow-x-auto">
                <table class="table">
                  <tbody>
                    <tr
                      v-for="game in stats.games_played.slice(0, 5)"
                      :key="game.id"
                      class="text-xl"
                    >
                      <td>{{ game.player1 }} vs {{ game.player2 }}</td>
                      <td>
                        <span
                          :class="
                            game.status === game.current_player
                              ? 'text-success'
                              : 'text-error'
                          "
                        >
                          {{ game.status }}
                        </span>
                      </td>
                      <td class="text-right">
                        {{ new Date(game.start).toDateString() }}
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>

        <div class="divider"></div>

        <!-- ELO Chart -->
        <EloChart :elo="stats.elo" />

        <!-- Back Button -->
        <div class="justify-center mt-4 card-actions">
          <button @click="navigateHome" class="text-xl btn btn-primary">
            Back to Home
          </button>
          <ShareNetwork
            network="facebook"
            :url="gameShareUrl"
            :title="shareTitle"
            :description="shareDescription"
            quote="Check out my Backgammon stats!"
          >
            <button class="text-xl btn btn-success">Share on Facebook</button>
          </ShareNetwork>

          <ShareNetwork
            network="twitter"
            :url="gameShareUrl"
            :title="shareTitle"
          >
            <button class="text-xl btn bg-blue-400">
              <svg
                class="w-10 h-10"
                fill="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  d="M8.29 20.251c7.547 0 11.675-6.253 11.675-11.675 0-.178 0-.355-.012-.53A8.348 8.348 0 0022 5.92a8.19 8.19 0 01-2.357.646 4.118 4.118 0 001.804-2.27 8.224 8.224 0 01-2.605.996 4.107 4.107 0 00-6.993 3.743 11.65 11.65 0 01-8.457-4.287 4.106 4.106 0 001.27 5.477A4.072 4.072 0 012.8 9.713v.052a4.105 4.105 0 003.292 4.022 4.095 4.095 0 01-1.853.07 4.108 4.108 0 003.834 2.85A8.233 8.233 0 012 18.407a11.616 11.616 0 006.29 1.84"
                ></path>
              </svg>
            </button>
          </ShareNetwork>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import router from '@/router'

import EloChart from '@/components/EloChart.vue'

import type { GameState } from '@/utils/game/types'
import type { User } from '@/utils/auth/types'
import { ShareNetwork } from 'vue-social-sharing'

const props = defineProps<{
  playerId?: string
}>()

interface GameStats {
  game_played: GameState[]
  win: number
  lost: number
  winrate: number
  elo: number[]
  cpu: number
  local: number
  online: number
  tournament: number
}

const stats = ref({
  cpu: 0,
  elo: [1280, 1300, 1320, 1340, 1360, 1300, 1280, 1200, 1150, 1000, 700, 1450],
  games_played: [
    {
      current_player: 'p1',
      double_owner: 'all',
      double_value: 1,
      elo1: 1000,
      elo2: 1000,
      end: '2021-01-01T00:00:00Z',
      game_type: 'online',
      id: 0,
      p1checkers: [0],
      p2checkers: [0],
      player1: 'Giorgio',
      player2: 'Mario',
      start: '2021-01-01T00:00:00Z',
      status: 'open',
      want_to_double: false,
    },
    {
      current_player: 'p2',
      double_owner: 'p1',
      double_value: 2,
      elo1: 1050,
      elo2: 1020,
      end: '2021-02-01T00:00:00Z',
      game_type: 'local',
      id: 1,
      p1checkers: [1],
      p2checkers: [1],
      player1: 'Alice',
      player2: 'Bob',
      start: '2021-02-01T00:00:00Z',
      status: 'p1',
      want_to_double: true,
    },
    {
      current_player: 'p1',
      double_owner: 'p2',
      double_value: 3,
      elo1: 1100,
      elo2: 1080,
      end: '2021-03-01T00:00:00Z',
      game_type: 'tournament',
      id: 2,
      p1checkers: [2],
      p2checkers: [2],
      player1: 'Charlie',
      player2: 'Dave',
      start: '2021-03-01T00:00:00Z',
      status: 'open',
      want_to_double: false,
    },
    {
      current_player: 'p2',
      double_owner: 'p1',
      double_value: 4,
      elo1: 1150,
      elo2: 1120,
      end: '2021-04-01T00:00:00Z',
      game_type: 'online',
      id: 3,
      p1checkers: [3],
      p2checkers: [3],
      player1: 'Eve',
      player2: 'Frank',
      start: '2021-04-01T00:00:00Z',
      status: 'p2',
      want_to_double: true,
    },
  ],
  local: 0,
  lost: 0,
  online: 0,
  tournament: 5,
  win: 100,
  winrate: 70,
})

// Reactive variable to store user ID
const currentUserId = ref<string | null>(null)

// Reactive variable for game share URL
const gameShareUrl = ref('')

// Computed properties for share title and description
const shareTitle = computed(() => `Backgammon Player Stats`)
const shareDescription = computed(
  () =>
    `Win Rate: ${stats.value.winrate}% | Games Played: ${stats.value.game_played?.length || 0}`,
)

onMounted(async () => {
  // getting player id from session | props
  const playerId = props.playerId
  try {
    const userResponse = await fetch('/api/session')
    if (!userResponse.ok) {
      throw new Error('Failed to fetch user session')
    }
    const user: User = await userResponse.json()

    // Determine user ID - prioritize prop, then session user, then null
    currentUserId.value = props.playerId || user.id || null

    // Construct share URL
    gameShareUrl.value = `${window.location.origin}/player/${currentUserId.value}`

    let response
    if (!playerId) {
      response = await fetch('/api/stats')
    } else {
      response = await fetch(`/api/stats/${playerId}`)
    }

    if (!response.ok) {
      throw new Error('Failed to fetch stats')
    }
    const tmp: GameStats = await response.json()
    //stats.value = tmp
  } catch (error) {
    console.error('Error fetching stats:', error)
  }
})

const navigateHome = () => {
  router.push('/')
}
</script>
