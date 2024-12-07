<script setup lang="ts">
import { useRoute } from 'vue-router'
import type { GameState } from '@/utils/game/types'
import { ref, onMounted } from 'vue'
import Board from '@/components/game/Board.vue'
import PlayerInfo from '@/components/game/PlayerInfo.vue'
const route = useRoute()

const gameId = Number(route.params.gameId)
const currentMove = ref(0)
const errMessage = ref('')

const gameState = ref<GameState | null>(null)
const avatars = ref<string[]>([])
const dices = ref<number[]>([])

async function getAvatars() {
  try {
    const res1 = await fetch(`/api/player/${gameState.value?.player1}/avatar`)
    const p1 = await res1.json()
    const res2 = await fetch(`/api/player/${gameState.value?.player2}/avatar`)
    const p2 = await res2.json()
    avatars.value = [p1, p2]
  } catch (err) {
    console.error('Error fetching avatars:', err)
  }
}

async function getGame(game_id: number, move: number) {
  try {
    const res = await fetch('/api/replay', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ game_id, move }),
    })
    if (!res.ok) {
      errMessage.value = 'Finished game'
      return false
    }
    errMessage.value = ''
    const data = await res.json()
    gameState.value = data.game
    dices.value = data.dices

    return true
  } catch (error) {
    console.error('Error accepting invite:', error)
    return false
  }
}

async function nextMove() {
  if (gameState.value) {
    if (await getGame(gameId, currentMove.value + 1)) {
      currentMove.value++
    }
  }
}

async function previousMove() {
  if (gameState.value && currentMove.value > 0) {
    if (await getGame(gameId, currentMove.value - 1)) {
      currentMove.value--
    }
  }
}

onMounted(async () => {
  await getGame(gameId, 0)
  await getAvatars()
  console.log(gameState.value)
})
</script>

<template>
  <div class="flex h-full w-full flex-col items-center justify-center p-4">
    <div class="flex h-full w-full gap-6">
      <div class="flex">
        <div
          class="flex w-48 flex-col justify-evenly rounded-lg border-8 border-primary bg-base-100 p-4 shadow-xl"
        >
          <PlayerInfo
            :username="gameState?.player2 || ''"
            :elo="gameState?.elo2 || 0"
            :isCurrentTurn="gameState?.current_player === 'p1'"
            :isOpponent="true"
            :avatar="avatars[1]"
          />

          <!-- Pulsanti -->
          <div class="flex flex-col items-center justify-center">
            <div class="my-4 flex justify-center gap-4">
              <button
                @click="previousMove"
                class="retro-button"
                :disabled="currentMove === 0"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-6 w-6"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M15 19l-7-7 7-7"
                  />
                </svg>
              </button>
              <span class="flex items-center font-bold">{{ currentMove }}</span>

              <button @click="nextMove" class="retro-button">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-6 w-6"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M9 5l7 7-7 7"
                  />
                </svg>
              </button>
            </div>
            <span class="flex items-center font-bold">{{ errMessage }}</span>
          </div>

          <PlayerInfo
            :username="gameState?.player1 || ''"
            :elo="gameState?.elo1 || 0"
            :isCurrentTurn="gameState?.current_player === 'p2'"
            :avatar="avatars[0]"
          />
        </div>
      </div>

      <Board
        v-if="gameState"
        :gameState="gameState"
        :availableMoves="null"
        :dicesReplay="dices"
        :resetDiceState="() => {}"
        :handleDiceRoll="() => {}"
      />
    </div>
  </div>
</template>

<style scoped>
.highlight-container {
  position: relative;
  cursor: pointer;
}

.highlight-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 0, 0.3);
  border: 2px solid yellow;
  border-radius: 0.5rem;
  pointer-events: none;
}

.retro-background {
  @apply bg-base-100;
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
  @apply rounded-lg border-4 border-8 border-primary bg-base-100 shadow-md;
}

.retro-button {
  @apply btn btn-primary rounded-lg border-4 border-primary shadow-md;
  border: 3px solid #8b4513;
  text-transform: uppercase;
  text-shadow: 2px 2px 0 rgba(0, 0, 0, 0.2);
  box-shadow: 0 2px 0 #8b4513;
  font-size: 1.1rem;

  &:hover {
    transform: translateY(2px);
    box-shadow:
      inset 0 0 10px rgba(0, 0, 0, 0.2),
      0 0px 0 #8b4513;
    cursor: url('/tortellino.png'), auto;
  }
}

.selected {
  stroke: yellow !important;
  stroke-width: 3 !important;
}

.confetti-container {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 1000;
  pointer-events: none;
}
</style>
