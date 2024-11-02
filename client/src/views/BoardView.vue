<script setup lang="ts">
import { ref } from 'vue'

interface BoardDimensions {
  width: number
  height: number
  triangleWidth: number
  triangleHeight: number
  centerBarWidth: number
  padding: number
  checkerRadius: number
}

const BOARD: BoardDimensions = {
  width: 800,
  height: 600,
  triangleWidth: 60,
  triangleHeight: 250,
  centerBarWidth: 20,
  padding: 20,
  checkerRadius: 20,
} as const

const selectedChecker = ref(null)
const checkers = ref([
  ...Array(2)
    .fill(null)
    .map((_, i) => ({ color: 'white', position: 0, stackIndex: i })),
  ...Array(5)
    .fill(null)
    .map((_, i) => ({ color: 'gray', position: 5, stackIndex: i })),
  ...Array(3)
    .fill(null)
    .map((_, i) => ({ color: 'gray', position: 7, stackIndex: i })),
  ...Array(5)
    .fill(null)
    .map((_, i) => ({ color: 'white', position: 11, stackIndex: i })),
  ...Array(2)
    .fill(null)
    .map((_, i) => ({ color: 'gray', position: 23, stackIndex: i })),
  ...Array(5)
    .fill(null)
    .map((_, i) => ({ color: 'white', position: 18, stackIndex: i })),
  ...Array(3)
    .fill(null)
    .map((_, i) => ({ color: 'white', position: 16, stackIndex: i })),
  ...Array(5)
    .fill(null)
    .map((_, i) => ({ color: 'gray', position: 12, stackIndex: i })),
])

// Index follows the counter-clockwise direction: upper from right (0) to left (11),
// lower from left (0) to right (11)
const getTrianglePath = (position: number): string => {
  const isUpper = position < 12
  const index = isUpper ? 11 - position : position - 12
  let x = BOARD.padding + index * BOARD.triangleWidth

  if (index >= 6) {
    x += BOARD.centerBarWidth * 2
  }

  const y = isUpper ? BOARD.padding : BOARD.height - BOARD.padding

  return isUpper
    ? `M ${x} ${y}
       L ${x + BOARD.triangleWidth} ${y}
       L ${x + BOARD.triangleWidth / 2} ${y + BOARD.triangleHeight} Z`
    : `M ${x} ${y}
       L ${x + BOARD.triangleWidth} ${y}
       L ${x + BOARD.triangleWidth / 2} ${y - BOARD.triangleHeight} Z`
}

const getTriangleColor = (position: number): string => {
  return position % 2 === 0 ? '#8B0000' : '#000080'
}

// position from 0 (upper right) to 23 (lower right)
const getCheckerX = (position: number) => {
  const index = position < 12 ? 11 - position : position - 12
  let x = BOARD.padding + index * BOARD.triangleWidth + BOARD.triangleWidth / 2

  if (index >= 6) {
    x += BOARD.centerBarWidth * 2
  }

  return x
}

const getCheckerY = (position: number, stackIndex: number) => {
  const spacing = BOARD.checkerRadius * 1.8

  if (position < 12) {
    return BOARD.padding + BOARD.checkerRadius + stackIndex * spacing
  } else {
    return (
      BOARD.height - BOARD.padding - BOARD.checkerRadius - stackIndex * spacing
    )
  }
}
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
