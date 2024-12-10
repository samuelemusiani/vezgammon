<script setup lang="ts">
import { onMounted, computed, onUnmounted, ref } from 'vue'
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
import TutorialButton from '@/components/buttons/TutorialButton.vue'
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
  isMyTurn,
  avatars,
  fetchGameState,
  fetchMoves,
  fetchSession,
  fetchAvatars,
} = useGameState()
const { timeLeft, startTimer, stopTimer } = useGameTimer()

const {
  diceRolled,
  resetDiceState,
  isRolling,
  displayedDice,
  handleDiceRoll,
  showDiceFromOpponent,
} = useDiceRoll()
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

const isTutorial = ref(false)

// get the variant option from the query
const query = new URLSearchParams(window.location.search)
const variant = query.get('variant')
isTutorial.value = variant === 'tutorial'

onMounted(async () => {
  try {
    await fetchGameState()
    await fetchMoves()
    await fetchSession()
    await fetchAvatars()
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

const handleTurnMade = async () => {
  await fetchGameState()
  await fetchMoves()
  resetDiceState()
  if (gameState.value?.game_type === 'online') startTimer()
  isMyTurn.value = true
}

const handleDiceRolled = (payload: string) => {
  const diceData = JSON.parse(payload)
  console.log('Dice rolled:', diceData)
  showDiceFromOpponent(diceData)
}

const updateCheckerPositionForPlayer1 = (move: Move) => {
  if (gameState.value?.p1checkers[25 - move.to] === 1) {
    gameState.value.p1checkers[25 - move.to] = 0
    gameState.value.p1checkers[0]++
  }
  gameState.value!.p2checkers[move.from]--
  if (move.to !== 25) {
    gameState.value!.p2checkers[move.to]++
  }
}

const updateCheckerPositionForPlayer2 = (move: Move) => {
  if (gameState.value?.p2checkers[25 - move.to] === 1) {
    gameState.value.p2checkers[25 - move.to] = 0
    gameState.value.p2checkers[0]++
  }
  gameState.value!.p1checkers[move.from]--
  if (move.to !== 25) {
    gameState.value!.p1checkers[move.to]++
  }
}

const handleMoveMade = (moveData: { move: Move }) => {
  if (!gameState.value) return

  const move = moveData.move
  if (whichPlayerAmI.value === 'p1') updateCheckerPositionForPlayer1(move)
  else updateCheckerPositionForPlayer2(move)
}

const handleMessage = async (message: WSMessage) => {
  switch (message.type) {
    case 'turn_made':
      await handleTurnMade()
      break
    case 'want_to_double':
      showDoubleModal.value = true
      break
    case 'double_accepted':
      await fetchGameState()
      break
    case 'dice_rolled':
      handleDiceRolled(message.payload)
      break
    case 'game_end':
      await handleEnd(session.value)
      break
    case 'move_made':
      handleMoveMade(JSON.parse(message.payload))
      break
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
  router.push({ name: 'home' })
}

const exitGame = async () => {
  try {
    await handleRetire()
    handleReturnHome()
  } catch (err) {
    console.error('Error exiting game:', err)
  }
}

const handleStopTimer = () => {
  isMyTurn.value = false
  stopTimer()
}

const handleAcceptDouble = async () => {
  showDoubleModal.value = false
  await acceptDouble()
  await fetchGameState()
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
            :avatar="avatars[1]"
          />

          <!-- Double Dice Here -->
          <DoubleDice
            :doubleValue="gameState?.double_value || 1"
            :showDoubleButton="showDoubleButton && isMyTurn"
            @double="handleDouble"
          />
          <!-- Game Timer -->
          <div
            class="my-8 flex flex-col items-center gap-3 border-y border-gray-200 py-4"
          >
            <GameTimer
              :timeLeft="timeLeft"
              :isMyTurn="isMyTurn && gameState?.game_type === 'online'"
            />
            <button class="retro-button" @click="exitGame">Exit Game</button>
          </div>

          <!-- Current Player Info -->
          <PlayerInfo
            :username="gameState?.player1 || ''"
            :elo="gameState?.elo1 || 0"
            :isCurrentTurn="gameState?.current_player === 'p2'"
            :avatar="avatars[0]"
          />
        </div>
      </div>

      <!-- Board Div -->
      <Board
        v-if="gameState"
        :gameState="gameState"
        :availableMoves="availableMoves"
        :isMyTurn="isMyTurn"
        :diceRolled="diceRolled"
        :resetDiceState="resetDiceState"
        :isRolling="isRolling"
        :displayedDice="displayedDice"
        :handleDiceRoll="handleDiceRoll"
        @ws-message="sendWSMessage"
        @fetch-moves="fetchMoves"
        @fetch-game-state="fetchGameState"
        @stop-timer="handleStopTimer"
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
      @confirm="handleAcceptDouble"
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
    <TutorialButton v-if="isTutorial" />
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

@keyframes dice-shake {
  0% {
    transform: rotate(0deg);
  }

  25% {
    transform: rotate(5deg);
  }

  75% {
    transform: rotate(-5deg);
  }

  100% {
    transform: rotate(0deg);
  }
}

.dice-rolling {
  animation: dice-shake 0.3s ease-in-out infinite;
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
