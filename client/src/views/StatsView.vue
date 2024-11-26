<template>
  <div
    class="retro-background min-h-screen items-center justify-center p-4 md:p-8"
  >
    <div class="mx-auto max-w-7xl">
      <div class="retro-box card-body p-4 md:p-8">
        <h2 class="retro-title mb-4 text-center text-2xl md:text-3xl">
          Player Statistics
        </h2>

        <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <!-- Overall Stats -->
          <div class="stats-section">
            <h3 class="retro-subtitle mb-4">Game Performance</h3>
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div class="stats-grid">
                <div class="stat-item">
                  <span class="stat-label">Played</span>
                  <span class="stat-value text-lg md:text-xl">
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
            <h3 class="retro-subtitle mb-4">Recent Games</h3>
            <div class="recent-games">
              <div
                v-for="game in stats.game_played"
                :key="game.id"
                class="game-item"
              >
                <!-- TODO: Get opponent name from db and winner -->
                <span>{{ game.player2 }}</span>
                <span
                  :class="game.status === 'Won' ? 'text-success' : 'text-error'"
                  class="text-center font-bold"
                  >{{ game.status }}</span
                >
                <span class="text-right">{{ game.start }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Back Button -->
        <div class="mt-8 flex justify-center">
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
const stats = ref<GameStats>({
  game_played: [],
  win: 0,
  lost: 0,
  winrate: 0,
  elo: [],
  cpu: 0,
  local: 0,
  online: 0,
  tournament: 0,
})

onMounted(async () => {
  try {
    const response = await fetch('/api/stats')
    if (!response.ok) {
      throw new Error('Failed to fetch stats')
    }
    const tmp: GameStats = await response.json()
    console.log(tmp)
    stats.value = tmp
    console.log(stats.value)
  } catch (error) {
    console.error('Error fetching stats:', error)
  }
})

const navigateHome = () => {
  router.push('/')
}
</script>

<style scoped>
.retro-background {
  background: #2c1810;
  background-image: repeating-linear-gradient(
      45deg,
      rgba(139, 69, 19, 0.1) 0px,
      rgba(139, 69, 19, 0.1) 2px,
      transparent 2px,
      transparent 10px
    ),
    repeating-linear-gradient(
      -45deg,
      rgba(139, 69, 19, 0.1) 0px,
      rgba(139, 69, 19, 0.1) 2px,
      transparent 2px,
      transparent 10px
    );
  cursor: url('/tortellino.png'), auto;
}

.retro-box {
  background-color: #ffe5c9;
  border: 5px solid #8b4513;
  box-shadow:
    0 0 0 4px #d2691e,
    inset 0 0 20px rgba(0, 0, 0, 0.2);
}

.retro-title {
  font-family: 'Arial Black', serif;
  color: #8b4513;
  text-shadow: 2px 2px 0 rgba(0, 0, 0, 0.2);
  letter-spacing: 2px;
}

.retro-subtitle {
  font-family: 'Arial Black', serif;
  color: #d2691e;
  font-size: 1.5rem;
  text-shadow: 1px 1px 0 rgba(0, 0, 0, 0.2);
}

.stats-section {
  @apply p-3 md:p-6;
  background: rgba(210, 105, 30, 0.1);
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

.retro-button {
  background: #d2691e;
  color: white;
  border: 3px solid #8b4513;
  font-family: 'Arial Black', serif;
  text-transform: uppercase;
  padding: 0.5rem 2rem;
  text-shadow: 1px 1px 0 rgba(0, 0, 0, 0.2);
  box-shadow: 0 2px 0 #8b4513;
  transition: all 0.2s;

  &:hover {
    transform: translateY(2px);
    box-shadow: none;
    cursor: url('/tortellino.png'), auto;
  }
}
</style>
