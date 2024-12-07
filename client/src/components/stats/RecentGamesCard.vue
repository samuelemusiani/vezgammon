<template>
  <div
    class="card glass bg-base-200/80 transition-all duration-300 hover:shadow-lg"
  >
    <div class="card-body">
      <h3 class="card-title text-2xl text-primary">Recent Matches</h3>
      <div class="max-h-[350px] overflow-x-auto overflow-y-auto">
        <table class="table table-zebra w-full">
          <thead class="sticky top-0 bg-base-200">
            <tr>
              <th>Players</th>
              <th>Result</th>
              <th>Date</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="game in sortedGames"
              :key="game.id"
              class="cursor-pointer hover:bg-base-300/50"
              @click="navigateToGame(game.id)"
            >
              <td class="font-medium">
                {{ game.player1 }}
                <span class="text-primary">vs</span>
                {{ game.player2 }}
              </td>
              <td>
                <span class="badge" :class="getResultClass(game)">
                  {{ getResultText(game) }}
                </span>
              </td>
              <td>{{ new Date(game.start).toLocaleDateString() }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { GameState } from '@/utils/game/types'
import { useRouter } from 'vue-router'
let router = useRouter()

const props = defineProps<{
  games: GameState[]
  currentUser: string
}>()

const sortedGames = computed(() => {
  const games = props.games
  return games
    .sort((a, b) => new Date(b.start).getTime() - new Date(a.start).getTime())
    .filter(game => game.game_type !== 'local')
})

const getResultClass = (game: GameState) => {
  if (game.status === 'winp1') {
    return game.player1 === props.currentUser ? 'badge-success' : 'badge-error'
  } else if (game.status === 'winp2') {
    return game.player2 === props.currentUser ? 'badge-success' : 'badge-error'
  }
  return 'badge-info'
}

const getResultText = (game: GameState) => {
  if (game.status === 'winp1') {
    return game.player1 === props.currentUser ? 'Win' : 'Lost'
  } else if (game.status === 'winp2') {
    return game.player2 === props.currentUser ? 'Win' : 'Lost'
  }
  return game.status
}

const navigateToGame = (gameId: number) => {
  console.log('here')
  console.log(gameId)
  router.push(`/analysis/${gameId}`)
}
</script>
