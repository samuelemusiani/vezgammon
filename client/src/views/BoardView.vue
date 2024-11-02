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
} from '@/utils/game/game'
import {
  calculatePossibleMoves,
  isCheckerMovable,
  updateStackIndices,
} from '@/utils/game/moves'

const gameState = ref(newGame())
const selectedChecker = ref<Checker | null>(null)
const possibleMoves = ref<number[]>([])
const movesAvailable = ref(2)
const isRolling = ref(false)

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
  moveChecker(selectedChecker.value, position)
  updateGameState(
    gameState.value,
    position,
    oldCheckerPos,
    movesAvailable.value,
  )
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

  isRolling.value = true
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
</script>

<template>
  <div
    class="flex min-h-screen flex-col items-center justify-center bg-gray-100 p-4"
  >
    <div class="flex w-full max-w-screen-lg">
      <!-- Board Div -->
      <div class="flex-grow">
        <div class="rounded-lg bg-white p-4 shadow-xl">
          <svg
            viewBox="0 0 800 600"
            preserveAspectRatio="xMidYMid meet"
            class="h-auto w-full"
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
              stroke="black"
              stroke-width="1.4"
              class="checker-transition"
              :class="{
                'cursor-pointer hover:opacity-80': isCheckerMovable(
                  gameState,
                  checker,
                ),
                selected: selectedChecker === checker,
              }"
              @click.stop="handleCheckerClick(checker)"
            />
          </svg>
        </div>
      </div>

      <!-- Dice Div -->
      <div
        class="flex w-48 flex-col justify-center rounded-lg bg-white p-2 shadow-xl"
      >
        <!-- Dice Section -->
        <div class="mb-4 flex flex-col items-center">
          <button
            @click="rollDice"
            :disabled="
              isRolling ||
              (gameState.dice.value[0] !== 0 && gameState.dice.value[1] !== 0)
            "
            class="mb-4 rounded-lg bg-blue-600 px-4 py-2 font-bold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
          >
            Roll Dice
          </button>

          <div class="mb-4 flex gap-4">
            <div
              v-for="(die, index) in gameState.dice.value"
              :key="index"
              class="flex h-12 w-12 items-center justify-center rounded-lg bg-white p-2 shadow-lg sm:h-16 sm:w-16"
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
