<template>
  <div class="justify-center items-center p-4 h-full md:p-8">
    <div class="mx-auto max-w-7xl">
      <div class="p-4 md:p-8 retro-box card-body">
        <h2 class="mb-4 text-2xl text-center md:text-3xl text-primary">
          Player Statistics
        </h2>

        <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <!-- Overall Stats -->
          <div class="stats-section">
            <h3 class="mb-4 retro-subtitle">Game Performance</h3>
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div class="stats-grid">
                <div class="stat-item">
                  <span class="stat-label">Played</span>
                  <span class="text-lg md:text-xl stat-value">
                    {{ 0 }}
                  </span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">Wins</span>
                  <span class="stat-value">{{ stats.win }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">Lost</span>
                  <span class="stat-value">{{ stats.lost }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">WinRate</span>
                  <span class="stat-value">{{ stats.winrate }}%</span>
                </div>
              </div>
              <div class="stats-grid">
                <div class="stat-item">
                  <span class="stat-label">CPU</span>
                  <span class="stat-value">{{ stats.cpu }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">Local </span>
                  <span class="stat-value">{{ stats.local }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">Online</span>
                  <span class="stat-value">{{ stats.online }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">Tournament</span>
                  <span class="stat-value">{{ stats.tournament }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Recent Games -->
          <div class="stats-section">
            <h3 class="mb-4 retro-subtitle">Recent Games</h3>
            <div class="recent-games">
              <div
                v-for="game in stats.games_played.slice(0, 5)"
                :key="game.id"
                class="mb-2 game-item"
              >
                <span>{{ game.player1 }} vs {{ game.player2 }}</span>
                <span
                  :class="
                    game.status === game.current_player
                      ? 'text-success'
                      : 'text-error'
                  "
                  class="font-bold text-center"
                  >{{ game.status }}</span
                >
                <span class="text-right">{{
                  new Date(game.start).toDateString()
                }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="divider"></div>
        <!-- ELO Chart -->
        <EloChart :elo="stats.elo" />

        <!-- Back Button -->
        <div class="flex justify-center mt-8">
          <button @click="navigateHome" class="retro-button">
            Back to Home
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import router from '@/router'

import EloChart from '@/components/EloChart.vue'

import type { GameState } from '@/utils/game/types'

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

// DUMMY DATA while waiting for API -- TODO: remove
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

onMounted(async () => {
  try {
    const response = await fetch('/api/stats')
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

<style scoped>
.stats-section {
  @apply bg-base-100 p-3 md:p-6;
  border: 3px solid #8b4513;
  border-radius: 8px;
}

.stats-grid {
  @apply grid grid-cols-1 gap-2 sm:grid-cols-2;
}

.stat-item {
  @apply p-2 md:p-4;
  display: flex;
  flex-direction: column;
  align-items: center;
  background: rgba(255, 255, 255, 0.6);
  border: 2px solid #8b4513;
  border-radius: 4px;
  justify-content: center;
  min-height: 70px;
  height: auto;
}

.game-item {
  @apply p-2 md:p-4;
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 0.5rem;
  background: rgba(255, 255, 255, 0.6);
  border: 2px solid #8b4513;
  border-radius: 4px;
}

.stat-label {
  @apply text-sm md:text-base;
  font-family: 'Arial Black', serif;
  color: #8b4513;
}

.stat-value {
  @apply text-base md:text-lg lg:text-xl;
  font-weight: bold;
  color: #d2691e;
}
</style>
