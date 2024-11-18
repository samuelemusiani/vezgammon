<template>
  <div class="retro-background flex h-screen items-center justify-center">
    <div class="card w-3/4">
      <div class="retro-box card-body">
        <h2 class="retro-title mb-4 text-center text-3xl">Player Statistics</h2>

        <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <!-- Overall Stats -->

          <div class="stats-section">
            <h3 class="retro-subtitle mb-4">Game Performance</h3>
            <div className="grid grid-cols-2 gap-1">
              <div class="stats-grid">
                <div class="stat-item">
                  <span class="stat-label">Played</span>
                  <span class="stat-value">{{ stats.gamesPlayed || 0 }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">Wins</span>
                  <span class="stat-value">{{ stats.wins || 0 }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">Lost</span>
                  <span class="stat-value">{{ stats.losses || 0 }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">WinRate</span>
                  <span class="stat-value"
                    >{{
                      calculateWinRate(stats.wins, stats.gamesPlayed)
                    }}%</span
                  >
                </div>
              </div>
              <div class="stats-grid">
                <div class="stat-item">
                  <span class="stat-label">CPU</span>
                  <span class="stat-value">{{ stats.aiGames || 0 }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">Local </span>
                  <span class="stat-value">{{ stats.localGames || 0 }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">Online</span>
                  <span class="stat-value">{{ stats.onlineGames || 0 }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">Tournament</span>
                  <span class="stat-value">{{
                    stats.tournamentGames || 0
                  }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Recent Games -->
          <div class="stats-section">
            <h3 class="retro-subtitle mb-4">Recent Games</h3>
            <div class="recent-games">
              <div
                v-for="game in stats.recentGames"
                :key="game.id"
                class="game-item"
              >
                <span>{{ game.opponent }}</span>
                <span
                  :class="game.result === 'Won' ? 'text-success' : 'text-error'"
                  class="text-center font-bold"
                  >{{ game.result }}</span
                >
                <span class="text-right">{{ game.date }}</span>
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
import { ref } from 'vue'
import router from '@/router'

interface GameStats {
  gamesPlayed: number
  wins: number
  losses: number
  aiGames: number
  localGames: number
  onlineGames: number
  tournamentGames: number
  recentGames: Array<{
    id: number
    opponent: string
    result: 'Won' | 'Lost'
    date: string
  }>
}

// DUMMY DATA while waiting for API -- TODO: remove
const stats = ref<GameStats>({
  gamesPlayed: 0,
  wins: 0,
  losses: 0,
  aiGames: 0,
  localGames: 0,
  onlineGames: 0,
  tournamentGames: 0,
  recentGames: [
    { id: 1, opponent: 'CPU', result: 'Won', date: '2024-09-01' },
    { id: 2, opponent: 'Local', result: 'Lost', date: '2024-09-02' },
    { id: 3, opponent: 'Online', result: 'Won', date: '2024-09-03' },
    { id: 4, opponent: 'Tournament', result: 'Lost', date: '2024-09-04' },
  ],
})

/* TO FETCH DATA FROM API
onMounted(async () => {
  try {
    const response = await fetch('/api/stats')
    if (!response.ok) {
      throw new Error('Failed to fetch stats')
    }
    stats.value = await response.json()
  } catch (error) {
    console.error('Error fetching stats:', error)
  }
})*/

const calculateWinRate = (wins: number, total: number): number => {
  if (total === 0) return 0
  return Math.round((wins / total) * 100)
}

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
  border: 6px solid #d2691e;
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
  background: rgba(210, 105, 30, 0.1);
  padding: 1.5rem;
  border: 3px solid #8b4513;
  border-radius: 8px;
}

.stats-grid {
  @apply grid grid-cols-2 gap-1;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  background: rgba(255, 255, 255, 0.6);
  padding: 1rem;
  border: 2px solid #8b4513;
  border-radius: 4px;
  justify-content: center;
  height: 90px;
}

.stat-label {
  font-family: 'Arial Black', serif;
  color: #8b4513;
  font-size: 0.9rem;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: bold;
  color: #d2691e;
}

.recent-games {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.game-item {
  display: grid;
  grid-template-columns: 1fr 40px 1fr;
  gap: 1rem;
  padding: 0.5rem;
  background: rgba(255, 255, 255, 0.6);
  border: 2px solid #8b4513;
  border-radius: 4px;
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
