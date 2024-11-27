<template>
  <div class="container py-8 px-4 mx-auto ">
    <div class="shadow-xl card bg-base-100 overflow-auto h-[90vh]">
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
          <button @click="navigateHome" class="btn btn-primary">
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
