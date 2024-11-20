<script setup lang="ts">
import { ref } from 'vue'
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
//import { useSound } from '@vueuse/sound'
//import victorySfx from '@/utils/sounds/victory.mp3'
//import diceSfx from '@/utils/sounds/dice.mp3'
import { onMounted } from 'vue'
//import tinSfx from '@/utils/sounds/tintin.mp3'

const gameState = ref<GameState>()

const selectedChecker = ref<Checker | null>(null)
const availableMoves = ref<MovesResponse | null>(null)
const possibleMoves = ref<number[]>([])
// Teniamo traccia della sequenza di mosse che stiamo seguendo
const selectedMoveSequence = ref<Move[]>([])
const currentMoveIndex = ref<number>(0)
const movesToSubmit = ref<Move[]>([]) // mosse già fatte
const currentPossibleSequences = ref<Move[][]>([]) // sequenze ancora possibili

const isRolling = ref(false)

const isExploding = ref(false)
//const { play: playVictory } = useSound(victorySfx)
//const { play: playDice } = useSound(diceSfx)
//const { play: playTin } = useSound(tinSfx)

// Fetch a /api/play on mounted
onMounted(async () => {
  await fetchGameState()
  await fetchMoves()
})

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
    console.log(availableMoves.value)
  } catch (err) {
    console.error('Error fetching moves:', err)
  }
}

const isCheckerSelectable = (checker: Checker) => {
  if (!gameState.value) return false
  // Converti il colore della pedina nel formato del player
  const checkerPlayer = checker.color === 'black' ? 'p1' : 'p2'
  console.log(checkerPlayer)
  return checkerPlayer === gameState.value.current_player
}

const handleCheckerClick = (checker: Checker) => {
  if (!availableMoves.value || !isCheckerSelectable(checker)) return
  console.log(checker)
  if (selectedChecker.value === checker) {
    selectedChecker.value = null
    possibleMoves.value = []
    return
  }

  selectedChecker.value = checker

  // Se è la prima mossa
  if (movesToSubmit.value.length === 0) {
    // Trova tutte le possibili destinazioni per questa pedina in tutte le sequenze
    possibleMoves.value = [
      ...new Set(
        availableMoves.value.possible_moves.flatMap(seq =>
          seq
            .filter(move => move.from === checker.position)
            .map(move => move.to),
        ),
      ),
    ]
    console.log('mosse posibili', possibleMoves.value)
  } else {
    // Per le mosse successive, trova tutte le sequenze che contengono le mosse già fatte
    possibleMoves.value = [
      ...new Set(
        currentPossibleSequences.value.flatMap(seq =>
          seq
            .filter(move => move.from === checker.position)
            .map(move => move.to),
        ),
      ),
    ]
  }
}
// Quando si clicca su un triangolo
const handleTriangleClick = async (position: number) => {
  if (
    !selectedChecker.value ||
    !possibleMoves.value.includes(position) ||
    !availableMoves.value
  )
    return

  const currentMove: Move = {
    from: selectedChecker.value.position,
    to: position,
  }

  // Se è la prima mossa
  if (movesToSubmit.value.length === 0) {
    // Filtra le sequenze di mosse possibili solo quelle che contengono la mossa appena giocata
    currentPossibleSequences.value = availableMoves.value.possible_moves.filter(
      seq =>
        seq.some(
          move => move.from === currentMove.from && move.to === currentMove.to,
        ),
    )
  } else {
    // Filtra le sequenze di mosse possibili solo quelle che contengono la mossa appena giocata
    currentPossibleSequences.value = currentPossibleSequences.value.filter(
      seq =>
        seq.some(
          move => move.from === currentMove.from && move.to === currentMove.to,
        ),
    )
  }

  // rimuovo dalle sequenze possibili la mossa appena giocata
  currentPossibleSequences.value = currentPossibleSequences.value.map(seq =>
    seq.filter(
      move => move.from !== currentMove.from || move.to !== currentMove.to,
    ),
  )
  console.log('sequenze possibili', currentPossibleSequences.value)

  // Aggiorna il gameState per mostrare la mossa appena giocata sulla board
  if (gameState.value.current_player === 'p1') {
    gameState.value.p1checkers[currentMove.from]--
    gameState.value.p1checkers[currentMove.to]++
  } else {
    gameState.value.p2checkers[currentMove.from]--
    gameState.value.p2checkers[currentMove.to]++
  }

  // Aggiungi la mossa a quelle fatte
  movesToSubmit.value.push(currentMove)
  console.log('mosse effettuate', movesToSubmit.value)

  // Reset della selezione corrente
  selectedChecker.value = null
  possibleMoves.value = []

  const hasPossibleMoves = currentPossibleSequences.value?.length > 0
  let hasUsedBothDices = movesToSubmit.value.length === 2
  if (availableMoves.value.dices[0] == availableMoves.value.dices[1]) {
    hasUsedBothDices = movesToSubmit.value.length === 4
  }

  if (hasUsedBothDices || !hasPossibleMoves) {
    try {
      // TODO: Remove: mi serve per ordinare le mosse come vuole il backend (altrimenti il backend non me le valida) )
      const matchingSequence = availableMoves.value.possible_moves.find(
        sequence => {
          return movesToSubmit.value.every(move =>
            sequence.some(
              seqMove => seqMove.from === move.from && seqMove.to === move.to,
            ),
          )
        },
      )

      if (!matchingSequence) {
        console.error('No matching sequence found')
        return
      }

      // Ordina movesToSubmit secondo l'ordine in matchingSequence
      const sortedMoves = matchingSequence.filter(seqMove =>
        movesToSubmit.value.some(
          move => move.from === seqMove.from && move.to === seqMove.to,
        ),
      )
      console.log('mosse da inviare', sortedMoves)

      const res = await fetch('/api/play/moves', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(sortedMoves),
      })
      console.log('stato POST', res.status)
      movesToSubmit.value = []
      currentPossibleSequences.value = []
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
  const isP2Turn = gameState.value.current_player === 'p2'

  // Aggiungi pedine del player 1 (bianche)
  gameState.value.p1checkers.forEach((count, position) => {
    // Se è il turno di p2, inverti la posizione (eccetto 0)
    const adjustedPosition =
      position === 0 ? 0 : isP2Turn ? 25 - position : position

    for (let i = 0; i < count; i++) {
      checkers.push({
        color: 'black',
        position: adjustedPosition,
        stackIndex: i,
      })
    }
  })

  // Aggiungi pedine del player 2 (nere)
  gameState.value.p2checkers.forEach((count, position) => {
    // Se è il turno di p2, non invertire la posizione
    // Se è il turno di p1, inverti la posizione (eccetto 0)
    const adjustedPosition =
      position === 0 ? 0 : isP2Turn ? position : 25 - position

    for (let i = 0; i < count; i++) {
      checkers.push({
        color: 'white',
        position: adjustedPosition,
        stackIndex: i,
      })
    }
  })

  return checkers
}

const exitGame = async () => {
  try {
    const res = await fetch('/api/play/', {
      method: 'DELETE',
    })
    console.log(res.status)
    await fetch('/api/play/local')
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
            <h3 class="text-lg font-bold">Opponent</h3>
            <p class="text-gray-600">ELO: 1850</p>
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
            <h3 class="text-lg font-bold">Player</h3>
            <p class="text-gray-600">ELO: 1720</p>
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
                @click="handleTriangleClick(position)"
              />
            </g>

            <!-- Possible moves highlights -->
            <path
              v-for="(position, index) in possibleMoves"
              :key="`highlight-${index}`"
              :d="getTrianglePath(position)"
              fill="yellow"
              opacity="1"
              pointer-events="none"
            />

            <!-- Checkers -->
            <circle
              v-for="(checker, index) in getCheckers()"
              :key="`checker-${index}`"
              :cx="getCheckerX(checker)"
              :cy="getCheckerY(checker)"
              :r="BOARD.checkerRadius"
              :fill="checker.color"
              :stroke="checker.color === 'white' ? 'black' : 'blue'"
              stroke-width="1.4"
              class="checker-transition"
              :class="{ selected: selectedChecker === checker }"
              @click="handleCheckerClick(checker)"
            />
          </svg>
        </div>
      </div>

      <!-- Dice Div -->
      <div
        class="retro-box flex w-48 flex-col justify-center rounded-lg bg-white p-2 shadow-xl"
      >
        <!--<div class="mb-4 flex flex-col items-center">
          <button
            @click="endTurn(gameState)"
            class="retro-button mb-4 rounded-lg bg-blue-600 px-4 py-2 font-bold text-white transition-colors hover:bg-blue-700"
          >
            End Turn
          </button>
        </div>-->

        <!-- Opponent's Captured Checkers -->
        <div
          class="captured-checkers-container mb-4 flex flex-col place-items-center"
        >
          <div
            class="h-64 w-16 overflow-hidden rounded-lg border-2 border-amber-900 p-1"
          >
            <div class="flex flex-col gap-1">
              <div
                v-for="(checker, index) in null"
                :key="'black-' + index"
                class="h-3 w-full rounded-full border border-blue-500 bg-black"
              ></div>
            </div>
          </div>
        </div>

        <!-- Dice Section
        <div class="mb-5 mt-5 flex flex-col items-center">
          <button
            @click="rollDice"
            :disabled="
              isRolling ||
              (gameState.dice.value[0] !== 0 && gameState.dice.value[1] !== 0)
            "
            class="retro-button mb-4 rounded-lg bg-blue-600 px-4 py-2 font-bold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
          >
            Roll Dice
          </button>-->

        <div class="flex gap-4">
          <div
            v-for="(die, index) in availableMoves?.dices"
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
      </div>

      <div
        class="captured-checkers-container mt-4 flex flex-col place-items-center"
      >
        <div
          class="h-64 w-16 overflow-hidden rounded-lg border-2 border-amber-900 p-1"
        >
          <div class="flex flex-col gap-1">
            <!-- Esempio di pedina bianca catturata
              <div
                v-for="(checker, index) in gameState.capturedWhite"
                :key="'white-' + index"
                class="h-3 w-full rounded-full border border-black bg-white"
              ></div>-->
          </div>
        </div>
      </div>
    </div>

    <!-- Game Info -->
    <!--<div class="game-info mb-4 text-center">
           <p class="text-lg font-bold">
            Current Player:
            <span
              :class="
                gameState.currentPlayer === 'white'
                  ? 'text-gray-700'
                  : 'text-gray-900'
              "
            >
              {{ gameState.currentPlayer === 'white' ? 'White' : 'Black' }}
            </span>
          </p>
          <p
            v-if="gameState.dice.value[0] && gameState.dice.value[1]"
            class="text-sm text-gray-600"
          >
            Moves remaining: {{ movesAvailable }}
            <br />
            Available values:
            <span :class="{ 'line-through': gameState.dice.used[0] }">{{
              gameState.dice.value[0]
            }}</span
            >,
            <span :class="{ 'line-through': gameState.dice.used[1] }">{{
              gameState.dice.value[1]
            }}</span>
            <span
              v-if="
                gameState.dice.value[0] !== gameState.dice.value[1] &&
                !gameState.dice.used[0] &&
                !gameState.dice.used[1]
              "
            >
              , or {{ gameState.dice.value[0] + gameState.dice.value[1] }}
            </span>
          </p>
        </div>-->
  </div>
</template>

<style scoped>
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
  stroke-width: 3;
}

.checker-transition {
  transition:
    cx 0.3s ease-out,
    cy 0.3s ease-out;
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
}
</style>
