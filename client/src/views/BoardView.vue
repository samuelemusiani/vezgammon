<script setup lang="ts">
import { onMounted, computed, onUnmounted } from 'vue'
import router from '@/router'

import type { Move } from '@/utils/game/types'
import type { WSMessage } from '@/utils/types'

import ConfettiExplosion from 'vue-confetti-explosion'
import Chat from '@/components/ChatContainer.vue'
import GameTimer from '@/components/game/GameTimer.vue'
import PlayerInfo from '@/components/game/PlayerInfo.vue'
import DoubleDice from '@/components/game/DoubleDice.vue'
import Board from '@/components/game/Board.vue'
import Modal from '@/components/Modal.vue'
import { useWebSocketStore } from '@/stores/websocket'
import { useGameState } from '@/composables/useGameState'
import { useGameTimer } from '@/composables/useGameTimer'
import { useGameDouble } from '@/composables/useGameDouble'
import { useDiceRoll } from '@/composables/useDiceRoll'
import { useGameEnd } from '@/composables/useGameEnd'

//import tinSfx from '@/utils/sounds/tintin.mp3'

const {
  gameState,
  availableMoves,
  session,
  fetchGameState,
  fetchMoves,
  fetchSession,
} = useGameState()
const { timeLeft, isMyTurn, startTimer, stopTimer } = useGameTimer()

const { showDiceFromOpponent } = useDiceRoll()
const {
  showResultModal,
  handleLose,
  isWinner,
  isExploding,
  handleEnd,
  handleRetire,
} = useGameEnd()

const { showDoubleModal, handleDouble, acceptDouble, declineDouble } =
  useGameDouble(handleLose)

const webSocketStore = useWebSocketStore()

onMounted(async () => {
  try {
    await fetchGameState()
    await fetchMoves()
    await fetchSession()
    webSocketStore.connect()
    webSocketStore.addMessageHandler(handleMessage)
    if (isMyTurn.value && gameState.value?.game_type === 'online') {
      startTimer()
    }
  } catch {
    console.error('Error fetching game state')
  }
})

onUnmounted(() => {
  stopTimer()
  webSocketStore.removeMessageHandler(handleMessage)
})

const handleMessage = async (message: WSMessage) => {
  if (message.type === 'turn_made') {
    await fetchGameState()
    await fetchMoves()
    startTimer()
    isMyTurn.value = true
  } else if (message.type === 'want_to_double') {
    showDoubleModal.value = true
  } else if (message.type === 'double_accepted') {
    await fetchGameState()
  } else if (message.type === 'dice_rolled') {
    const diceData = JSON.parse(message.payload)
    showDiceFromOpponent(diceData.dices)
  } else if (message.type === 'game_end') {
    await handleEnd(session.value)
  } else if (message.type === 'move_made') {
    if (!gameState.value) return

    const moveData = JSON.parse(message.payload)
    const move = moveData.move as Move

    if (whichPlayerAmI.value === 'p1') {
      if (gameState.value.p1checkers[25 - move.to] === 1) {
        gameState.value.p1checkers[25 - move.to] = 0
        gameState.value.p1checkers[0]++
      }
      gameState.value.p2checkers[move.from]--
      if (move.to !== 25) {
        gameState.value.p2checkers[move.to]++
      }
    } else {
      if (gameState.value.p2checkers[25 - move.to] === 1) {
        gameState.value.p2checkers[25 - move.to] = 0
        gameState.value.p2checkers[0]++
      }
      gameState.value.p1checkers[move.from]--
      if (move.to !== 25) {
        gameState.value.p1checkers[move.to]++
      }
    }
  }
}

const showDoubleButton = computed(() => {
  if (!gameState.value) return false
  if (gameState.value.double_owner === 'all') return true
  if (gameState.value.game_type === 'local') {
    return gameState.value.current_player === gameState.value.double_owner
  }
  return gameState.value.double_owner === whichPlayerAmI.value
})

const handleDoubleWinExit = async () => {
  showResultModal.value = false
  handleReturnHome()
}

const whichPlayerAmI = computed(() => {
  if (gameState.value?.player1 === session.value?.username) {
    return 'p1'
  } else {
    return 'p2'
  }
})

const handleReturnHome = () => {
  router.push('/')
}

const exitGame = async () => {
  try {
    await handleRetire()
    handleReturnHome()
  } catch (err) {
    console.error('Error exiting game:', err)
  }
}

function sendWSMessage(message: WSMessage) {
  webSocketStore.sendMessage(message)
}
</script>

<template>
  <div class="flex h-full w-full flex-col items-center justify-center p-4">
    <div class="fixed top-[25%]">
      <ConfettiExplosion
        v-if="isExploding"
        :stageHeight="1000"
        :particleCount="300"
      />
    </div>
    <div class="flex h-full w-full gap-6">
      <!-- Opponent and Player Info -->
      <div class="flex">
        <div
          class="flex w-48 flex-col justify-evenly rounded-lg border-8 border-primary bg-base-100 p-4 shadow-xl"
        >
          <!-- Opponent Info -->
          <PlayerInfo
            :username="gameState?.player2 || ''"
            :elo="gameState?.elo2 || 0"
            :isCurrentTurn="gameState?.current_player === 'p1'"
            :isOpponent="true"
          />

          <!-- Double Dice Here -->
          <DoubleDice
            :doubleValue="gameState?.double_value || 1"
            :showDoubleButton="showDoubleButton"
            @double="handleDouble"
          />
          <!-- Game Timer -->
          <div
            class="my-8 flex flex-col items-center gap-3 border-y border-gray-200 py-4"
          >
            <GameTimer :timeLeft="timeLeft" :isMyTurn="isMyTurn" />
            <button class="retro-button" @click="exitGame">Exit Game</button>
          </div>

          <!-- Current Player Info -->
          <PlayerInfo
            :username="gameState?.player1 || ''"
            :elo="gameState?.elo1 || 0"
            :isCurrentTurn="gameState?.current_player === 'p2'"
          />
        </div>
      </div>

      <!-- Board Div -->

      <Board
        @ws-message="sendWSMessage"
        v-if="gameState"
        :gameState="gameState"
        :availableMoves="availableMoves"
        @fetch-moves="fetchMoves"
        @fetch-game-state="fetchGameState"
      />
    </div>

    <Chat
      v-if="
        session?.username &&
        (gameState?.game_type === 'online' || gameState?.game_type === 'bot')
      "
      :opponentUsername="
        gameState?.player1 === session?.username
          ? gameState?.player2 || ''
          : gameState?.player1 || ''
      "
      :gameType="gameState?.game_type || ''"
      :myUsername="session?.username || ''"
    />

    <!-- Double Confirmation Modal -->
    <Modal
      :show="showDoubleModal"
      title="Confirm Double"
      confirmText="Confirm"
      cancelText="Cancel"
      confirmVariant="success"
      @confirm="acceptDouble"
      @cancel="declineDouble"
    >
      Your opponent has offered a double. Do you accept?
    </Modal>

    <!-- Game Result Modal -->
    <Modal
      :show="showResultModal"
      :title="isWinner ? 'Victory!' : 'Defeat!'"
      :confirmText="'Return to Menu'"
      :confirmVariant="isWinner ? 'success' : 'danger'"
      @confirm="handleDoubleWinExit"
    >
      <p class="text-lg">
        {{
          isWinner
            ? 'Congratulations! You won the game!'
            : 'Game Over! Better luck next time!'
        }}
      </p>
    </Modal>
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
