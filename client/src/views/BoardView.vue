<script setup lang="ts">
import { ref } from 'vue'
import type { BoardDimensions, Checker } from '@/utils/game/types'
import {
  createDefaultBoard,
  BOARD,
  getTrianglePath,
  getTriangleColor,
  getCheckerX,
  getCheckerY,
} from '@/utils/game/game'

// Initialize checkers
const checkers = ref(createDefaultBoard())
</script>

<template>
  <div
    class="flex min-h-screen flex-col items-center justify-center bg-gray-100 p-4"
  >
    <div class="flex w-full max-w-screen-lg">
      <div class="flex-grow">
        <div class="rounded-lg bg-white p-4 shadow-xl">
          <svg
            viewBox="0 0 800 600"
            preserveAspectRatio="xMidYMid meet"
            class="h-auto w-full"
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
              />
            </g>

            <!-- Checkers -->
            <circle
              v-for="(checker, index) in checkers"
              :key="`checker-${index}`"
              :cx="getCheckerX(checker.position)"
              :cy="getCheckerY(checker.position, checker.stackIndex)"
              :r="BOARD.checkerRadius"
              :fill="checker.color"
              stroke="black"
              stroke-width="1.4"
            />
          </svg>
        </div>
      </div>
    </div>
  </div>
</template>
