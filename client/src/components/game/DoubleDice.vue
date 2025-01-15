<script setup lang="ts">
interface DoubleDiceProps {
  doubleValue: number
  showDoubleButton: boolean
}

defineProps<DoubleDiceProps>()

defineEmits<(e: 'double') => void>()
</script>

<template>
  <div class="flex flex-col items-center">
    <div class="flex h-16 w-16 items-center justify-center rounded-lg p-2">
      <svg viewBox="0 0 60 60" width="100%" height="100%">
        <!-- Dice border -->
        <rect
          x="1"
          y="1"
          width="58"
          height="58"
          rx="8"
          class="fill-primary"
          stroke="black"
          stroke-width="2"
        />
        <!-- Number text -->
        <text
          x="30"
          y="40"
          text-anchor="middle"
          font-size="30"
          font-weight="bold"
          fill="white"
        >
          {{ doubleValue }}
        </text>
      </svg>
    </div>
    <button
      @click="$emit('double')"
      class="retro-button mt-2"
      v-show="showDoubleButton"
    >
      Double
    </button>
  </div>
</template>
