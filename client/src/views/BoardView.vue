<script setup lang="ts">
import { ref } from 'vue'
import type { Checker, GameState } from '@/utils/game/types'
import {
  BOARD,
  getTrianglePath,
  getTriangleColor,
  getCheckerX,
  getCheckerY,
  updateGameState,
  newGame,
  endTurn,
  checkWin,
} from '@/utils/game/game'
import {
  calculatePossibleMoves,
  handleHitChecker,
  isCheckerMovable,
  updateStackIndices,
} from '@/utils/game/moves'
import ConfettiExplosion from 'vue-confetti-explosion'
import { useSound } from '@vueuse/sound'
import victorySfx from '@/utils/sounds/victory.mp3'
import diceSfx from '@/utils/sounds/dice.mp3'

const gameState = ref(newGame())
const selectedChecker = ref<Checker | null>(null)
const possibleMoves = ref<number[]>([])
const movesAvailable = ref(2)
const isRolling = ref(false)
const gameStarted = ref(false)
const isExploding = ref(false)
const { play: playVictory } = useSound(victorySfx)
const { play: playDice } = useSound(diceSfx)

const handleCheckerClick = (checker: Checker) => {
  if (!isCheckerMovable(gameState.value, checker)) return

  if (selectedChecker.value === checker) {
    selectedChecker.value = null
    possibleMoves.value = []
  } else {
    selectedChecker.value = checker
    possibleMoves.value = calculatePossibleMoves(gameState.value, checker)
    console.log(possibleMoves.value)
  }
}

const handleTriangleClick = (position: number) => {
  if (!selectedChecker.value || !possibleMoves.value.includes(position)) return
  const oldCheckerPos = selectedChecker.value.position

  // Controlla se c'è una pedina da mangiare
  handleHitChecker(gameState.value, position)

  moveChecker(selectedChecker.value, position)
  updateGameState(
    gameState.value,
    position,
    oldCheckerPos,
    movesAvailable.value,
  )
  if (checkWin(gameState.value)) {
    handleWin()
  }
  possibleMoves.value = []
  selectedChecker.value = null
}

const moveChecker = (checker: Checker, newPosition: number) => {
  checker.position = newPosition
  updateStackIndices(gameState.value)
}

const handleBoardClick = (event: any) => {
  // Deselect checker when clicking on board
  if (event.target.tagName === 'svg' || event.target.tagName === 'rect') {
    selectedChecker.value = null
    possibleMoves.value = []
  }
}

const rollDice = () => {
  if (
    isRolling.value ||
    (gameState.value.dice.value[0] !== 0 && gameState.value.dice.value[1] !== 0)
  )
    return

  if (!gameStarted.value) {
    gameStarted.value = true
    startTimer()
  }
  isRolling.value = true
  playDice()
  const rollAnimation = setInterval(() => {
    gameState.value.dice.value = [
      Math.floor(Math.random() * 6) + 1,
      Math.floor(Math.random() * 6) + 1,
    ]
  }, 50)

  setTimeout(() => {
    clearInterval(rollAnimation)
    isRolling.value = false

    const isDouble =
      gameState.value.dice.value[0] === gameState.value.dice.value[1]

    if (isDouble) {
      gameState.value.dice.used = [false, false, false, false]
      movesAvailable.value = 4
      gameState.value.dice.double = true
    } else {
      gameState.value.dice.used = [false, false]
      movesAvailable.value = 2
      gameState.value.dice.double = false
    }

    selectedChecker.value = null
    possibleMoves.value = []
  }, 500)
}

const handleWin = () => {
  isExploding.value = true
  playVictory()
}

const simulateWin = () => {
  gameState.value.board = gameState.value.board.map(checker => {
    if (checker.color === gameState.value.currentPlayer) {
      return {
        ...checker,
        position: gameState.value.currentPlayer === 'white' ? 0 : 23,
      }
    }
    return checker
  })
  handleWin()
}

let timerInterval: NodeJS.Timeout
const startTimer = () => {
  let seconds = 0
  timerInterval = setInterval(() => {
    seconds++
    const minutes = Math.floor(seconds / 60)
    const remainingSeconds = seconds % 60
    gameState.value.time = `${String(minutes).padStart(2, '0')}:${String(
      remainingSeconds,
    ).padStart(2, '0')}`
  }, 1000)
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
                  gameState.currentPlayer === 'black'
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
            <p class="text-sm text-gray-600">Total Time</p>
            <p class="text-2xl font-bold">{{ gameState.time }}</p>
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
                  gameState.currentPlayer === 'white'
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
            @click="handleBoardClick($event)"
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
                :d="getTrianglePath(position - 1)"
                :fill="getTriangleColor(position - 1)"
                stroke="black"
                stroke-width="1"
                @click="handleTriangleClick(position - 1)"
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
              v-for="(checker, index) in gameState.board"
              :key="`checker-${index}`"
              :cx="getCheckerX(checker.position)"
              :cy="getCheckerY(checker.position, checker.stackIndex)"
              :r="BOARD.checkerRadius"
              :fill="checker.color"
              :stroke="checker.color === 'white' ? 'black' : 'blue'"
              stroke-width="1.4"
              class="checker-transition"
              :class="{
                'cursor-pointer hover:opacity-80': isCheckerMovable(
                  gameState,
                  checker,
                ),
                selected: selectedChecker === checker,
              }"
              :style="
                selectedChecker === checker
                  ? { stroke: checker.color === 'white' ? 'black' : 'white' }
                  : {}
              "
              @click.stop="handleCheckerClick(checker)"
            />
          </svg>
        </div>
      </div>

      <!-- Dice Div -->
      <div
        class="retro-box flex w-48 flex-col justify-center rounded-lg bg-white p-2 shadow-xl"
      >
        <div class="mb-4 flex flex-col items-center">
          <button
            @click="endTurn(gameState)"
            class="retro-button mb-4 rounded-lg bg-blue-600 px-4 py-2 font-bold text-white transition-colors hover:bg-blue-700"
          >
            End Turn
          </button>
          <button
            @click="simulateWin"
            class="mt-4 rounded-lg bg-red-600 px-4 py-2 font-bold text-white transition-colors hover:bg-red-700"
          >
            Debug: Simulate Win
          </button>
        </div>
        <!-- Dice Section -->
        <div class="mb-4 flex flex-col items-center">
          <button
            @click="rollDice"
            :disabled="
              isRolling ||
              (gameState.dice.value[0] !== 0 && gameState.dice.value[1] !== 0)
            "
            class="retro-button mb-4 rounded-lg bg-blue-600 px-4 py-2 font-bold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
          >
            Roll Dice
          </button>

          <div class="mb-4 flex gap-4">
            <div
              v-for="(die, index) in gameState.dice.value"
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

        <!-- Game Info -->
        <div class="game-info mb-4 text-center">
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
        </div>
      </div>
    </div>
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
</style>
