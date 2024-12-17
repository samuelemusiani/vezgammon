<script setup lang="ts">
import type {
  Checker,
  GameState,
  Move,
  MovesResponse,
} from '@/utils/game/types'
import type { WSMessage } from '@/utils/types'
import Modal from '@/components/Modal.vue'

import DiceContainer from '@/components/game/DiceContainer.vue'
import CapturedCheckers from '@/components/game/CapturedCheckers.vue'

import { useGameMoves } from '@/composables/useGameMoves'
import { ref, watch } from 'vue'

import {
  BOARD,
  getTrianglePath,
  getTriangleColor,
  getCheckerX,
  getCheckerY,
} from '@/utils/game/game'

const $props = defineProps<{
  gameState: GameState
  availableMoves: MovesResponse | null
  isMyTurn?: boolean
  dicesReplay?: number[]
  diceRolled?: boolean
  isRolling?: boolean
  displayedDice?: number[]
  resetDiceState: () => void
  handleDiceRoll: (
    availableMoves: MovesResponse | null,
    online: boolean,
  ) => void
}>()

const availableMoves = ref($props.availableMoves)
const gameState = ref($props.gameState)
const showNoMovesModal = ref(false)

watch(
  () => $props.availableMoves,
  () => {
    availableMoves.value = $props.availableMoves
  },
)

watch(
  () => $props.gameState,
  () => {
    gameState.value = $props.gameState
  },
)

const { selectedChecker, possibleMoves, movesToSubmit, submitMoves } =
  useGameMoves()

const $emits = defineEmits<{
  (e: 'ws-message', payload: WSMessage): void
  (e: 'fetch-moves'): void
  (e: 'fetch-game-state'): void
  (e: 'stop-timer'): void
}>()

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

  // Aggiorna il gameState.value per mostrare la mossa appena giocata sulla board
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

  if (gameState.value.game_type === 'online') {
    $emits('ws-message', {
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

  const hasNoPossibleMoves =
    availableMoves.value.possible_moves?.length == 0 ||
    availableMoves.value.possible_moves.every(
      (sequence: Move[]) => sequence.length === 0,
    )
  let hasUsedBothDices = movesToSubmit.value.length === 2
  if (availableMoves.value.dices[0] == availableMoves.value.dices[1]) {
    hasUsedBothDices = movesToSubmit.value.length === 4
  }

  if (hasUsedBothDices || hasNoPossibleMoves) {
    try {
      if (!hasUsedBothDices) {
        showNoMovesModal.value = true
        return
      }
      await submitMoves()
      $props.resetDiceState()
      if (gameState.value.game_type !== 'online') {
        $emits('fetch-game-state')
        $emits('fetch-moves')
      } else {
        $emits('stop-timer')
      }
    } catch (err) {
      console.error('Error submitting moves:', err)
    }
  }
}

const getCheckers = () => {
  if (!gameState.value) return []

  const checkers: Checker[] = []

  gameState.value.p1checkers.forEach((count: any, position: any) => {
    for (let i = 0; i < count; i++) {
      checkers.push({
        color: 'black',
        position: position,
        stackIndex: i,
      })
    }
  })

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

const isCheckerSelectable = (checker: Checker) => {
  console.log('isCheckerSelectable', checker)
  console.log('gameState', gameState.value)
  console.log('diceRolled', $props.diceRolled)
  if (!gameState.value || !$props.diceRolled) return false

  const checkerPlayer = checker.color === 'black' ? 'p1' : 'p2'

  console.log('here 1')
  console.log(gameState.value.current_player)
  if (checkerPlayer !== gameState.value.current_player) return false
  console.log('here 2')

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

const isCheckerSelected = (checker: Checker) => {
  return (
    selectedChecker.value &&
    selectedChecker.value.position === checker.position &&
    selectedChecker.value.stackIndex === checker.stackIndex &&
    selectedChecker.value.color === checker.color
  )
}

const handleCheckerClick = async (checker: Checker) => {
  console.log('handleCheckerClick', checker)
  console.log('available-moves', availableMoves.value)
  console.log('isCheckerSelectable', isCheckerSelectable(checker))
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
    showNoMovesModal.value = true
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

const handleNoMovesConfirm = async () => {
  showNoMovesModal.value = false
  await submitMoves()
  $props.resetDiceState()
  if (gameState.value.game_type !== 'online') {
    $emits('fetch-game-state')
    $emits('fetch-moves')
  } else {
    $emits('stop-timer')
  }
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
</script>

<template>
  <div class="flex flex-1 justify-between gap-2 lg:gap-6">
    <div class="flex-1">
      <div class="retro-box h-full w-full rounded-lg p-2 shadow-xl">
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
            v-for="(position, index) in possibleMoves.filter(e => e != 25)"
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
      class="retro-box flex w-1/6 max-w-48 justify-center overflow-y-auto rounded-lg bg-white shadow-xl"
    >
      <div
        class="flex scale-[0.55] flex-col justify-center md:scale-[0.57] lg:scale-[0.80] xl:scale-100"
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
          :diceRolled="diceRolled || false"
          :displayedDice="displayedDice || [0, 0]"
          :isRolling="isRolling || false"
          :canRoll="isMyTurn as boolean"
          :dicesReplay="dicesReplay"
          @roll="
            handleDiceRoll(availableMoves, gameState.game_type === 'online')
          "
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
    <Modal
      :show="showNoMovesModal"
      title="No Moves Available"
      confirm-text="OK"
      confirm-variant="primary"
      @confirm="handleNoMovesConfirm"
    >
      <p>No possible moves left. Your turn is over.</p>
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
  @apply rounded-lg border-4 border-primary bg-base-100 shadow-md lg:border-8;
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
