<script setup lang="ts">
interface BoardDimensions {
  width: number
  height: number
  triangleWidth: number
  triangleHeight: number
  centerBarWidth: number
  padding: number
}

const BOARD: BoardDimensions = {
  width: 800,
  height: 600,
  triangleWidth: 60,
  triangleHeight: 250,
  centerBarWidth: 20,
  padding: 20,
} as const

// Index follows the counter-clockwise direction: upper from right (0) to left (11),
// lower from left (0) to right (11)
const getTrianglePath = (index: number, isUpper: boolean): string => {
  let x: number

  if (isUpper) {
    x = BOARD.padding + (11 - index) * BOARD.triangleWidth
    if (index < 6) {
      x += BOARD.centerBarWidth * 2
    }
  } else {
    x = BOARD.padding + index * BOARD.triangleWidth
    if (index >= 6) {
      x += BOARD.centerBarWidth * 2
    }
  }

  const y: number = isUpper ? BOARD.padding : BOARD.height - BOARD.padding

  if (isUpper) {
    return `M ${x} ${y}
            L ${x + BOARD.triangleWidth} ${y}
            L ${x + BOARD.triangleWidth / 2} ${y + BOARD.triangleHeight} Z`
  } else {
    return `M ${x} ${y}
            L ${x + BOARD.triangleWidth} ${y}
            L ${x + BOARD.triangleWidth / 2} ${y - BOARD.triangleHeight} Z`
  }
}

const getTriangleColor = (index: number): string => {
  return index % 2 === 0 ? '#8B0000' : '#000080'
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

            <!-- Upper and Lower triangles -->
            <g v-for="isUpper in [true, false]" :key="String(isUpper)">
              <path
                v-for="index in 12"
                :key="`${isUpper ? 'upper' : 'lower'}-${index}`"
                :d="getTrianglePath(index - 1, isUpper)"
                :fill="getTriangleColor(index - 1)"
                stroke="black"
                stroke-width="1"
              />
            </g>
          </svg>
        </div>
      </div>
    </div>
  </div>
</template>
