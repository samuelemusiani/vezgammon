<script setup lang="ts">
import type {
  Checker,
  GameState,
  Move,
  MovesResponse,
} from '@/utils/game/types'
import type { WSMessage } from '@/utils/types'

import DiceContainer from '@/components/game/DiceContainer.vue'
import CapturedCheckers from '@/components/game/CapturedCheckers.vue'

import { useGameMoves } from '@/composables/useGameMoves'
import { useDiceRoll } from '@/composables/useDiceRoll'

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
}>()

const { selectedChecker, possibleMoves, movesToSubmit, submitMoves } =
  useGameMoves()

const {
  diceRolled,
  resetDiceState,
  showDiceFromOpponent,
  isRolling,
  displayedDice,
  handleDiceRoll,
} = useDiceRoll()

const $emits = defineEmits<{
  (e: 'ws-message', payload: WSMessage): void
  (e: 'fetch-moves'): void
  (e: 'fetch-game-state'): void
  (e: 'reset-dice-state'): void
  (e: 'handle-dice-roll'): void
}>()

const handleTriangleClick = async (position: number) => {
  if (
    !selectedChecker.value ||
    !possibleMoves.value.includes(position) ||
    !$props.availableMoves ||
    !$props.gameState
  )
    return

  const currentMove: Move = {
    from: selectedChecker.value.position,
    to: position,
  }

  // Filtra le sequenze di mosse possibili solo quelle che contengono la mossa appena giocata
  $props.availableMoves.possible_moves =
    $props.availableMoves.possible_moves.filter((seq: Move[]) => {
      return seq[0].from === currentMove.from && seq[0].to === currentMove.to
    })

  // rimuovo dalle sequenze possibili la mossa appena giocata
  $props.availableMoves.possible_moves =
    $props.availableMoves.possible_moves.map((seq: Move[]) => {
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
    })
  console.log('sequenze possibili', $props.availableMoves.possible_moves)

  if ($props.gameState.current_player === 'p1') {
    console.log('moving checker')
    if ($props.gameState.p2checkers[25 - position] === 1) {
      // Capture the opponent's checker
      $props.gameState.p2checkers[25 - position] = 0
      $props.gameState.p2checkers[0]++
    }
  } else {
    console.log('moving checker')
    if ($props.gameState.p1checkers[25 - position] === 1) {
      // Capture the opponent's checker
      $props.gameState.p1checkers[25 - position] = 0
      $props.gameState.p1checkers[0]++
    }
  }

  // Aggiorna il $props.gameState per mostrare la mossa appena giocata sulla board
  if ($props.gameState.current_player === 'p1') {
    $props.gameState.p1checkers[currentMove.from]--
    if (currentMove.to !== 25) {
      // Solo se non è una mossa di uscita
      $props.gameState.p1checkers[currentMove.to]++
    }
  } else {
    $props.gameState.p2checkers[currentMove.from]--
    if (currentMove.to !== 25) {
      // Solo se non è una mossa di uscita
      $props.gameState.p2checkers[currentMove.to]++
    }
  }

  if ($props.gameState.game_type === 'online') {
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

  const hasPossibleMoves = $props.availableMoves.possible_moves?.length > 0
  let hasUsedBothDices = movesToSubmit.value.length === 2
  if ($props.availableMoves.dices[0] == $props.availableMoves.dices[1]) {
    hasUsedBothDices = movesToSubmit.value.length === 4
  }

  if (hasUsedBothDices || !hasPossibleMoves) {
    try {
      await submitMoves()
      resetDiceState()
      if ($props.gameState.game_type !== 'online') {
        await $emits('fetch-game-state')
        await $emits('fetch-moves')
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
  console.log('Getting checkers', $props.gameState)
  if (!$props.gameState) return []

  const checkers: Checker[] = []

  // Aggiungi pedine del player 1 (bianche)
  $props.gameState.p1checkers.forEach((count: any, position: any) => {
    for (let i = 0; i < count; i++) {
      checkers.push({
        color: 'black',
        position: position,
        stackIndex: i,
      })
    }
  })

  // Aggiungi pedine del player 2 (nere)
  $props.gameState.p2checkers.forEach((count: any, position: any) => {
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
  console.log('gameState', $props.gameState)
  console.log('diceRolled', diceRolled.value)
  if (!$props.gameState || !diceRolled.value) return false

  const checkerPlayer = checker.color === 'black' ? 'p1' : 'p2'

  console.log('here 1')
  console.log($props.gameState.current_player)
  if (checkerPlayer !== $props.gameState.current_player) return false
  console.log('here 2')

  // Check if the player has any checkers in position 0 (the bar)
  let barCheckersCount = 0
  if (checkerPlayer === 'p1') {
    barCheckersCount = $props.gameState.p1checkers[0]
  } else {
    barCheckersCount = $props.gameState.p2checkers[0]
  }

  // If the player has checkers on the bar Only the top checker in position 0 is selectable
  if (barCheckersCount > 0) {
    return checker.position === 0 && checker.stackIndex === barCheckersCount - 1
  }

  // Ottieni il numero totale di pedine nella posizione della pedina per questo giocatore
  let totalCheckersAtPosition = 0
  if (checkerPlayer === 'p1') {
    totalCheckersAtPosition = $props.gameState.p1checkers[checker.position]
  } else {
    totalCheckersAtPosition = $props.gameState.p2checkers[checker.position]
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

const handleCheckerClick = (checker: Checker) => {
  console.log('handleCheckerClick', checker)
  console.log('available-moves', $props.availableMoves)
  console.log('isCheckerSelectable', isCheckerSelectable(checker))
  if (!$props.availableMoves || !isCheckerSelectable(checker)) return
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

  console.log('mosse possibili', $props.availableMoves.possible_moves.length)

  if (
    $props.availableMoves.possible_moves.length === 0 ||
    $props.availableMoves.possible_moves.every(
      (sequence: Move[]) => sequence.length === 0,
    )
  ) {
    console.log(
      'No possible moves or all sequences are empty, passing the turn',
    )
    submitMoves()
    $emits('fetch-moves')
    $emits('fetch-game-state')
    return
  }

  selectedChecker.value = checker
  possibleMoves.value = [
    ...new Set(
      $props.availableMoves.possible_moves
        .map((seq: Move[]) => seq[0]) // Prendo solo la prima mossa di ogni sequenza
        .filter((move: Move) => move.from === checker.position)
        .map((move: Move) => move.to),
    ),
  ]
  console.log('mosse posibili', possibleMoves.value)
}

const getOutCheckers = (player: 'p1' | 'p2' | string) => {
  if (!$props.gameState) return 0

  // Calcola il numero totale di pedine iniziali (15)
  const initialCheckers = 15

  // Calcola il numero di pedine ancora sulla board
  const remainingCheckers =
    player === 'p1'
      ? $props.gameState.p1checkers.reduce(
          (acc: any, curr: any) => acc + curr,
          0,
        )
      : $props.gameState.p2checkers.reduce(
          (acc: any, curr: any) => acc + curr,
          0,
        )

  // Ritorna la differenza tra le pedine iniziali e quelle rimaste
  return initialCheckers - remainingCheckers
}
</script>

<template>
  <div class="flex flex-1 justify-between gap-6">
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
                $props.gameState?.current_player === 'p2'
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
              $props.gameState?.current_player === 'p2'
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
            :cy="getCheckerY(checker, $props.gameState as GameState)"
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
