<script setup lang="ts">
import { onMounted, computed, onUnmounted } from 'vue'
import router from '@/router'

import type { Checker, GameState, Move } from '@/utils/game/types'
import type { WSMessage } from '@/utils/types'
import {
  BOARD,
  getTrianglePath,
  getTriangleColor,
  getCheckerX,
  getCheckerY,
} from '@/utils/game/game'

import ConfettiExplosion from 'vue-confetti-explosion'
import Chat from '@/components/ChatContainer.vue'
import GameTimer from '@/components/game/GameTimer.vue'
import PlayerInfo from '@/components/game/PlayerInfo.vue'
import DoubleDice from '@/components/game/DoubleDice.vue'
import DiceContainer from '@/components/game/DiceContainer.vue'
import CapturedCheckers from '@/components/game/CapturedCheckers.vue'
import Modal from '@/components/Modal.vue'
import { useWebSocketStore } from '@/stores/websocket'
import { useGameState } from '@/composables/useGameState'
import { useGameMoves } from '@/composables/useGameMoves'
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
const { selectedChecker, possibleMoves, movesToSubmit, submitMoves } =
  useGameMoves()
const { timeLeft, isMyTurn, startTimer, stopTimer } = useGameTimer()

const {
  isRolling,
  diceRolled,
  displayedDice,
  handleDiceRoll,
  resetDiceState,
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

const isCheckerSelectable = (checker: Checker) => {
  if (!gameState.value || !diceRolled.value) return false

  const checkerPlayer = checker.color === 'black' ? 'p1' : 'p2'

  if (checkerPlayer !== gameState.value.current_player) return false

  // Check if the player has any checkers in position 0 (the bar)
  let barCheckersCount = 0
  if (checkerPlayer === 'p1') {
    barCheckersCount = gameState.value.p1checkers[0]
  } else {
    barCheckersCount = gameState.value.p2checkers[0]
  }

  // If the player has checkers on the bar Only the top checker in position 0 is selectable
  if (barCheckersCount > 0) {
    return checker.position === 0 && checker.stackIndex === barCheckersCount - 1
  }

  // Ottieni il numero totale di pedine nella posizione della pedina per questo giocatore
  let totalCheckersAtPosition = 0
  if (checkerPlayer === 'p1') {
    totalCheckersAtPosition = gameState.value.p1checkers[checker.position]
  } else {
    totalCheckersAtPosition = gameState.value.p2checkers[checker.position]
  }

  // Solo la pedina in cima (con stackIndex più alto) è selezionabile
  if (checker.stackIndex !== totalCheckersAtPosition - 1) return false

  return true
}

const handleCheckerClick = (checker: Checker) => {
  if (!availableMoves.value || !isCheckerSelectable(checker)) return
  console.log(checker)
  if (
    selectedChecker.value &&
    selectedChecker.value.position === checker.position &&
    selectedChecker.value.stackIndex === checker.stackIndex
  ) {
    selectedChecker.value = null
    possibleMoves.value = []
    return
  }

  console.log('mosse possibili', availableMoves.value.possible_moves.length)

  if (
    availableMoves.value.possible_moves.length === 0 ||
    availableMoves.value.possible_moves.every(
      (sequence: Move[]) => sequence.length === 0,
    )
  ) {
    console.log(
      'No possible moves or all sequences are empty, passing the turn',
    )
    submitMoves()
    fetchMoves()
    fetchGameState()
    return
  }

  selectedChecker.value = checker
  possibleMoves.value = [
    ...new Set(
      availableMoves.value.possible_moves
        .map((seq: Move[]) => seq[0]) // Prendo solo la prima mossa di ogni sequenza
        .filter((move: Move) => move.from === checker.position)
        .map((move: Move) => move.to),
    ),
  ]
  console.log('mosse posibili', possibleMoves.value)
}

const handleTriangleClick = async (position: number) => {
  if (
    !selectedChecker.value ||
    !possibleMoves.value.includes(position) ||
    !availableMoves.value ||
    !gameState.value
  )
    return

  const currentMove: Move = {
    from: selectedChecker.value.position,
    to: position,
  }

  // Filtra le sequenze di mosse possibili solo quelle che contengono la mossa appena giocata
  availableMoves.value.possible_moves =
    availableMoves.value.possible_moves.filter((seq: Move[]) => {
      return seq[0].from === currentMove.from && seq[0].to === currentMove.to
    })

  // rimuovo dalle sequenze possibili la mossa appena giocata
  availableMoves.value.possible_moves = availableMoves.value.possible_moves.map(
    (seq: Move[]) => {
      let removed = false
      return seq.filter(move => {
        if (
          !removed &&
          move.from === currentMove.from &&
          move.to === currentMove.to
        ) {
          removed = true
          return false
        }
        return true
      })
    },
  )
  console.log('sequenze possibili', availableMoves.value.possible_moves)

  if (gameState.value.current_player === 'p1') {
    console.log('moving checker')
    if (gameState.value.p2checkers[25 - position] === 1) {
      // Capture the opponent's checker
      gameState.value.p2checkers[25 - position] = 0
      gameState.value.p2checkers[0]++
    }
  } else {
    console.log('moving checker')
    if (gameState.value.p1checkers[25 - position] === 1) {
      // Capture the opponent's checker
      gameState.value.p1checkers[25 - position] = 0
      gameState.value.p1checkers[0]++
    }
  }

  // Aggiorna il gameState per mostrare la mossa appena giocata sulla board
  if (gameState.value.current_player === 'p1') {
    gameState.value.p1checkers[currentMove.from]--
    if (currentMove.to !== 25) {
      // Solo se non è una mossa di uscita
      gameState.value.p1checkers[currentMove.to]++
    }
  } else {
    gameState.value.p2checkers[currentMove.from]--
    if (currentMove.to !== 25) {
      // Solo se non è una mossa di uscita
      gameState.value.p2checkers[currentMove.to]++
    }
  }

  if (gameState.value?.game_type === 'online') {
    webSocketStore.sendMessage({
      type: 'move_made',
      payload: JSON.stringify({
        move: currentMove,
      }),
    })
  }

  // Aggiungi la mossa a quelle fatte
  movesToSubmit.value.push(currentMove)
  console.log('mosse effettuate', movesToSubmit.value)

  // Reset della selezione corrente
  selectedChecker.value = null
  possibleMoves.value = []

  const hasPossibleMoves = availableMoves.value.possible_moves?.length > 0
  let hasUsedBothDices = movesToSubmit.value.length === 2
  if (availableMoves.value.dices[0] == availableMoves.value.dices[1]) {
    hasUsedBothDices = movesToSubmit.value.length === 4
  }

  if (hasUsedBothDices || !hasPossibleMoves) {
    try {
      await submitMoves()
      resetDiceState()
      if (gameState.value.game_type !== 'online') {
        await fetchGameState()
        await fetchMoves()
      } else {
        stopTimer()
        isMyTurn.value = false
      }
    } catch (err) {
      console.error('Error submitting moves:', err)
    }
  }
}

const getCheckers = () => {
  if (!gameState.value) return []

  const checkers: Checker[] = []

  // Aggiungi pedine del player 1 (bianche)
  gameState.value.p1checkers.forEach((count: any, position: any) => {
    for (let i = 0; i < count; i++) {
      checkers.push({
        color: 'black',
        position: position,
        stackIndex: i,
      })
    }
  })

  // Aggiungi pedine del player 2 (nere)
  gameState.value.p2checkers.forEach((count: any, position: any) => {
    for (let i = 0; i < count; i++) {
      checkers.push({
        color: 'white',
        position: position,
        stackIndex: i,
      })
    }
  })

  return checkers
}

const isCheckerSelected = (checker: Checker) => {
  return (
    selectedChecker.value &&
    selectedChecker.value.position === checker.position &&
    selectedChecker.value.stackIndex === checker.stackIndex &&
    selectedChecker.value.color === checker.color
  )
}

const getOutCheckers = (player: 'p1' | 'p2' | string) => {
  if (!gameState.value) return 0

  // Calcola il numero totale di pedine iniziali (15)
  const initialCheckers = 15

  // Calcola il numero di pedine ancora sulla board
  const remainingCheckers =
    player === 'p1'
      ? gameState.value.p1checkers.reduce(
          (acc: any, curr: any) => acc + curr,
          0,
        )
      : gameState.value.p2checkers.reduce(
          (acc: any, curr: any) => acc + curr,
          0,
        )

  // Ritorna la differenza tra le pedine iniziali e quelle rimaste
  return initialCheckers - remainingCheckers
}

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
      <div class="flex-1">
        <div class="retro-box h-full rounded-lg p-4 shadow-xl">
          <svg
            viewBox="0 0 800 600"
            preserveAspectRatio="xMidYMid meet"
            class="h-full w-full"
          >
            <!-- Board background -->
            <rect
              x="0"
              y="0"
              :width="BOARD.width"
              :height="BOARD.height"
              class="fill-[#ebcb97]"
              stroke="brown"
              stroke-width="2"
              rx="6"
            />

            <!-- Center bar -->
            <rect
              :x="BOARD.width / 2 - BOARD.centerBarWidth / 2"
              y="0"
              :width="BOARD.centerBarWidth"
              :height="BOARD.height"
              class="fill-amber-900"
            />

            <!-- Triangles from 0 (upper right) to 23 (lower right) -->
            <g>
              <path
                v-for="position in 24"
                :key="`triangle-${position}`"
                :d="getTrianglePath(position)"
                :fill="getTriangleColor(position)"
                stroke="black"
                stroke-width="1"
                @click="
                  gameState?.current_player === 'p2'
                    ? handleTriangleClick(position)
                    : handleTriangleClick(25 - position)
                "
              />
            </g>

            <!-- Possible moves highlights -->
            <path
              v-for="(position, index) in possibleMoves"
              :key="`highlight-${index}`"
              :d="
                gameState?.current_player === 'p2'
                  ? getTrianglePath(position)
                  : getTrianglePath(25 - position)
              "
              fill="yellow"
              opacity="1"
              pointer-events="none"
            />

            <!-- Checkers -->
            <circle
              v-for="(checker, index) in getCheckers()"
              :key="`checker-${index}`"
              :cx="getCheckerX(checker)"
              :cy="getCheckerY(checker, gameState as GameState)"
              :r="BOARD.checkerRadius"
              :fill="checker.color"
              :stroke="checker.color === 'white' ? 'black' : 'blue'"
              stroke-width="1.4"
              class="checker-transition"
              :class="{ selected: isCheckerSelected(checker) }"
              @click="handleCheckerClick(checker)"
            />
          </svg>
        </div>
      </div>

      <!-- Right Container -->
      <div
        class="retro-box flex w-48 flex-col justify-evenly rounded-lg bg-white p-2 shadow-xl"
      >
        <!-- Captured Checkers -->
        <CapturedCheckers
          player="p1"
          :checkerCount="getOutCheckers('p1')"
          :isHighlighted="possibleMoves.includes(25)"
          @click="handleTriangleClick(25)"
        />

        <!-- Roll Dice Button -->
        <DiceContainer
          :diceRolled="diceRolled"
          :displayedDice="displayedDice"
          :isRolling="isRolling"
          :canRoll="!diceRolled && !!availableMoves?.dices"
          @roll="handleDiceRoll(availableMoves)"
        />

        <!-- Captured Checkers -->
        <CapturedCheckers
          player="p2"
          :checkerCount="getOutCheckers('p2')"
          :isHighlighted="possibleMoves.includes(25)"
          @click="handleTriangleClick(25)"
        />
      </div>
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

.captured-checkers-container {
  .h-64 {
    background: #f5c27a;
    /* Un colore leggermente più chiaro del tabellone */
  }

  .w-full {
    transition: all 0.3s ease-out;
  }

  transition: all 0.3s ease;
}
</style>
