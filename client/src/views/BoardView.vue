<script setup lang="ts">
import { ref } from 'vue'
import router from '@/router'

import type {
  Checker,
  GameState,
  MovesResponse,
  Move,
} from '@/utils/game/types'
import {
  BOARD,
  getTrianglePath,
  getTriangleColor,
  getCheckerX,
  getCheckerY,
  //checkWin,
} from '@/utils/game/game'

import ConfettiExplosion from 'vue-confetti-explosion'
import { useSound } from '@vueuse/sound'
import victorySfx from '@/utils/sounds/victory.mp3'
import diceSfx from '@/utils/sounds/dice.mp3'
import { onMounted } from 'vue'
//import tinSfx from '@/utils/sounds/tintin.mp3'

const gameState = ref<GameState | null>(null)

const selectedChecker = ref<Checker | null>(null)
const availableMoves = ref<MovesResponse | null>(null)
const possibleMoves = ref<number[]>([])
const movesToSubmit = ref<Move[]>([]) // mosse già fatte
const displayedDice = ref<number[]>([])

const isRolling = ref(false)
const diceRolled = ref(false)

const isExploding = ref(false)
const { play: playVictory } = useSound(victorySfx)
const { play: playDice } = useSound(diceSfx)
//const { play: playTin } = useSound(tinSfx)

// Fetch a /api/play on mounted
onMounted(async () => {
  try {
    await fetchGameState()
    await fetchMoves()
  } catch {
    console.error('Error fetching game state')
  }
})

const checkWin = () => {
  if (!gameState.value) return false
  // TODO: when backend is ready, check using gameState.value.status
  if (getOutCheckers(gameState.value.current_player) == 15) {
    return true
  }
  return false
}

const handleWin = () => {
  playVictory()
  isExploding.value = true
  setTimeout(() => {
    isExploding.value = false
  }, 5000)
}

const handleDiceRoll = () => {
  if (diceRolled.value || !availableMoves.value?.dices) return

  isRolling.value = true
  diceRolled.value = true
  playDice()

  // Funzione per generare numeri casuali dei dadi
  const generateRandomDice = () => {
    displayedDice.value = [
      Math.floor(Math.random() * 6) + 1,
      Math.floor(Math.random() * 6) + 1,
    ]
  }

  // Genera numeri casuali ogni 100ms durante l'animazione
  const rollInterval = setInterval(generateRandomDice, 100)

  // Dopo 3 secondi, mostra i veri valori dei dadi
  setTimeout(() => {
    clearInterval(rollInterval)
    isRolling.value = false
    displayedDice.value = availableMoves.value!.dices
  }, 1000)
}

const fetchGameState = async () => {
  try {
    const res = await fetch('/api/play/')
    const data: GameState = await res.json()

    gameState.value = data

    console.log(gameState.value)
  } catch (err) {
    console.error('Error fetching game state:', err)
  }
}

const fetchMoves = async () => {
  try {
    const res = await fetch('/api/play/moves')
    const data: MovesResponse = await res.json()
    availableMoves.value = data
    diceRolled.value = false
    displayedDice.value = []
    console.log(availableMoves.value)
  } catch (err) {
    console.error('Error fetching moves:', err)
  }
}

const isCheckerSelectable = (checker: Checker) => {
  if (!gameState.value || !diceRolled.value) return false

  // Converti il colore della pedina nel formato del player
  const checkerPlayer = checker.color === 'black' ? 'p1' : 'p2'

  // Verifica se la pedina appartiene al giocatore corrente
  if (checkerPlayer !== gameState.value.current_player) return false

  // Check if the player has any checkers in position 0 (the bar)
  let barCheckersCount = 0
  if (checkerPlayer === 'p1') {
    barCheckersCount = gameState.value.p1checkers[0]
  } else {
    barCheckersCount = gameState.value.p2checkers[0]
  }

  // If the player has checkers on the bar (position 0)
  if (barCheckersCount > 0) {
    // Only the checker in position 0 is selectable
    if (checker.position !== 0) return false
    // Only the top checker in position 0 is selectable
    if (checker.stackIndex !== barCheckersCount - 1) return false

    return true
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
    availableMoves.value.possible_moves.every(sequence => sequence.length === 0)
  ) {
    console.log(
      'No possible moves or all sequences are empty, passing the turn',
    )
    // Se non ci sono mosse possibili o tutte le sequenze sono vuote, passa il turno
    fetch('/api/play/moves', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(movesToSubmit.value),
    })
    fetchMoves()
    fetchGameState()
    return
  }

  selectedChecker.value = checker
  possibleMoves.value = [
    ...new Set(
      availableMoves.value.possible_moves
        .map(seq => seq[0]) // Prendo solo la prima mossa di ogni sequenza
        .filter(move => move.from === checker.position)
        .map(move => move.to),
    ),
  ]
  console.log('mosse posibili', possibleMoves.value)
}
// Quando si clicca su un triangolo
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
    availableMoves.value.possible_moves.filter(seq => {
      return seq[0].from === currentMove.from && seq[0].to === currentMove.to
    })

  // rimuovo dalle sequenze possibili la mossa appena giocata
  availableMoves.value.possible_moves = availableMoves.value.possible_moves.map(
    seq => {
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
      const res = await fetch('/api/play/moves', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(movesToSubmit.value),
      })
      console.log('stato POST', res.status)
      // TODO: change with backend check
      if (checkWin()) {
        handleWin()
      }
      movesToSubmit.value = []
      possibleMoves.value = []
      await fetchGameState()
      await fetchMoves()
    } catch (err) {
      console.error('Error submitting moves:', err)
    }
  }
}

const getCheckers = () => {
  if (!gameState.value) return []

  const checkers: Checker[] = []

  // Aggiungi pedine del player 1 (bianche)
  gameState.value.p1checkers.forEach((count, position) => {
    for (let i = 0; i < count; i++) {
      checkers.push({
        color: 'black',
        position: position,
        stackIndex: i,
      })
    }
  })

  // Aggiungi pedine del player 2 (nere)
  gameState.value.p2checkers.forEach((count, position) => {
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
      ? gameState.value.p1checkers.reduce((acc, curr) => acc + curr, 0)
      : gameState.value.p2checkers.reduce((acc, curr) => acc + curr, 0)

  // Ritorna la differenza tra le pedine iniziali e quelle rimaste
  return initialCheckers - remainingCheckers
}

const exitGame = async () => {
  try {
    const res = await fetch('/api/play/', {
      method: 'DELETE',
    })
    console.log(res.status)
    router.push('/')
  } catch (err) {
    console.error('Error exiting game:', err)
  }
}
</script>

<template>
  <div
    class="retro-background flex min-h-screen flex-col items-center justify-center bg-gray-100 p-4"
  >
    <div class="fixed top-[25%]">
      <ConfettiExplosion
        v-if="isExploding"
        :stageHeight="1000"
        :particleCount="300"
      />
    </div>
    <div class="flex w-full max-w-screen-2xl gap-6">
      <!-- Opponent and Player Info -->
      <div class="flex">
        <div
          class="retro-box flex w-48 flex-col justify-evenly rounded-lg bg-white p-4 shadow-xl"
        >
          <!-- Opponent Info -->
          <div class="mb-8 flex flex-col items-center">
            <div class="relative mb-2">
              <div class="h-16 w-16 overflow-hidden rounded-full bg-gray-200">
                <img
                  src="https://api.dicebear.com/6.x/avataaars/svg?seed=opponent"
                  alt="Opponent avatar"
                  class="h-full w-full object-cover"
                />
              </div>
              <div
                :class="[
                  'absolute -bottom-1 right-0 h-4 w-4 rounded-full border-2 border-white',
                  gameState?.current_player === 'p1'
                    ? 'bg-green-500'
                    : 'bg-gray-300',
                ]"
              ></div>
            </div>
            <h3 class="text-lg font-bold">{{ gameState?.player2 }}</h3>
            <p class="text-gray-600">ELO: {{ gameState?.elo2 }}</p>
          </div>

          <!-- Game Timer -->
          <div
            class="my-8 flex flex-col items-center border-y border-gray-200 py-4"
          >
            <button class="btn btn-primary" @click="exitGame">Exit Game</button>
            <p class="text-sm text-gray-600">Total Time</p>
            <p class="text-2xl font-bold"></p>
          </div>

          <!-- Current Player Info -->
          <div class="mt-8 flex flex-col items-center">
            <div class="relative mb-2">
              <div class="h-16 w-16 overflow-hidden rounded-full bg-gray-200">
                <img
                  src="https://api.dicebear.com/6.x/avataaars/svg?seed=player"
                  alt="Player avatar"
                  class="h-full w-full object-cover"
                />
              </div>
              <div
                :class="[
                  'absolute -bottom-1 right-0 h-4 w-4 rounded-full border-2 border-white',
                  gameState?.current_player === 'p2'
                    ? 'bg-green-500'
                    : 'bg-gray-300',
                ]"
              ></div>
            </div>
            <h3 class="text-lg font-bold">{{ gameState?.player1 }}</h3>
            <p class="text-gray-600">ELO: {{ gameState?.elo1 }}</p>
          </div>
        </div>
      </div>

      <!-- Board Div -->
      <div class="flex-1">
        <div class="retro-box h-full rounded-lg bg-white p-4 shadow-xl">
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

      <!-- Dice Div -->
      <div
        class="retro-box flex w-48 flex-col justify-evenly rounded-lg bg-white p-2 shadow-xl"
      >
        <!-- Opponent's Captured Checkers -->
        <div
          class="captured-checkers-container mb-4 flex flex-col place-items-center"
          :class="{ 'highlight-container': possibleMoves.includes(25) }"
          @click="handleTriangleClick(25)"
        >
          <div
            class="h-64 w-16 overflow-hidden rounded-lg border-2 border-amber-900 p-1"
          >
            <div class="flex flex-col gap-1">
              <div
                v-for="i in getOutCheckers('p1')"
                :key="'black-' + i"
                class="h-3 w-full rounded-full border border-blue-500 bg-black"
              ></div>
            </div>
          </div>
        </div>

        <!-- Roll Dice Button -->
        <div
          v-if="!diceRolled && availableMoves?.dices"
          class="mb-4 flex justify-center"
        >
          <button @click="handleDiceRoll" class="retro-button">
            Roll Dice
          </button>
        </div>

        <div class="flex justify-center gap-4">
          <div
            v-for="(die, index) in diceRolled ? displayedDice : []"
            :key="index"
            class="retro-box flex h-12 w-12 items-center justify-center rounded-lg bg-white p-2 shadow-lg sm:h-16 sm:w-16"
            :class="{ 'dice-rolling': isRolling }"
          >
            <svg viewBox="0 0 60 60">
              <!-- Dice border -->
              <rect
                x="1"
                y="1"
                width="58"
                height="58"
                rx="8"
                fill="white"
                stroke="black"
                stroke-width="2"
              />

              <!-- Number as dots inside the dice -->
              <template v-if="die === 1">
                <circle cx="30" cy="30" r="5" fill="black" />
              </template>

              <template v-if="die === 2">
                <circle cx="20" cy="20" r="5" fill="black" />
                <circle cx="40" cy="40" r="5" fill="black" />
              </template>

              <template v-if="die === 3">
                <circle cx="20" cy="20" r="5" fill="black" />
                <circle cx="30" cy="30" r="5" fill="black" />
                <circle cx="40" cy="40" r="5" fill="black" />
              </template>

              <template v-if="die === 4">
                <circle cx="20" cy="20" r="5" fill="black" />
                <circle cx="40" cy="20" r="5" fill="black" />
                <circle cx="20" cy="40" r="5" fill="black" />
                <circle cx="40" cy="40" r="5" fill="black" />
              </template>

              <template v-if="die === 5">
                <circle cx="20" cy="20" r="5" fill="black" />
                <circle cx="40" cy="20" r="5" fill="black" />
                <circle cx="30" cy="30" r="5" fill="black" />
                <circle cx="20" cy="40" r="5" fill="black" />
                <circle cx="40" cy="40" r="5" fill="black" />
              </template>

              <template v-if="die === 6">
                <circle cx="20" cy="15" r="5" fill="black" />
                <circle cx="40" cy="15" r="5" fill="black" />
                <circle cx="20" cy="30" r="5" fill="black" />
                <circle cx="40" cy="30" r="5" fill="black" />
                <circle cx="20" cy="45" r="5" fill="black" />
                <circle cx="40" cy="45" r="5" fill="black" />
              </template>
            </svg>
          </div>
        </div>

        <div
          class="captured-checkers-container mt-4 flex flex-col place-items-center"
          :class="{ 'highlight-container': possibleMoves.includes(25) }"
          @click="handleTriangleClick(25)"
        >
          <div
            class="h-64 w-16 overflow-hidden rounded-lg border-2 border-amber-900 p-1"
          >
            <div class="flex flex-col gap-1">
              <div
                v-for="i in getOutCheckers('p2')"
                :key="'white-' + i"
                class="h-3 w-full rounded-full border border-black bg-white"
              ></div>
            </div>
          </div>
        </div>
      </div>
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

.retro-button {
  @apply btn;
  background: #d2691e;
  color: white;
  border: 3px solid #8b4513;
  font-family: 'Arial Black', serif;
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

.captured-checkers-container {
  .h-64 {
    background: #f5c27a; /* Un colore leggermente più chiaro del tabellone */
  }
  .w-full {
    transition: all 0.3s ease-out;
  }
  transition: all 0.3s ease;
}
</style>
