<script setup lang="ts">
import { useRoute } from 'vue-router'
import type { GameState } from '@/utils/game/types'
import { ref, onMounted } from 'vue'
import Board from '@/components/game/Board.vue'
import PlayerInfo from '@/components/game/PlayerInfo.vue'
const route = useRoute()

const gameId = Number(route.params.gameId)
const gameState = ref<GameState | null>(null)

async function getGame(game_id: number, move: number) {
  try {
    const res = await fetch('/api/replay', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ game_id, move }),
    })

    gameState.value = await res.json()
  } catch (error) {
    console.error('Error accepting invite:', error)
  }
}

onMounted(async () => {
  await getGame(gameId, 3)
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
          />

          <PlayerInfo
            :username="gameState?.player1 || ''"
            :elo="gameState?.elo1 || 0"
            :isCurrentTurn="gameState?.current_player === 'p2'"
          />
        </div>
      </div>

      <Board v-if="gameState" :gameState="gameState" :availableMoves="null" />
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
