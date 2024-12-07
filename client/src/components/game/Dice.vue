<script setup lang="ts">
defineProps<{
  value: number
  isRolling: boolean
}>()
</script>

<template>
  <div
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
      <template v-if="value === 1">
        <circle cx="30" cy="30" r="5" fill="black" />
      </template>

      <template v-if="value === 2">
        <circle cx="20" cy="20" r="5" fill="black" />
        <circle cx="40" cy="40" r="5" fill="black" />
      </template>

      <template v-if="value === 3">
        <circle cx="20" cy="20" r="5" fill="black" />
        <circle cx="30" cy="30" r="5" fill="black" />
        <circle cx="40" cy="40" r="5" fill="black" />
      </template>

      <template v-if="value === 4">
        <circle cx="20" cy="20" r="5" fill="black" />
        <circle cx="40" cy="20" r="5" fill="black" />
        <circle cx="20" cy="40" r="5" fill="black" />
        <circle cx="40" cy="40" r="5" fill="black" />
      </template>

      <template v-if="value === 5">
        <circle cx="20" cy="20" r="5" fill="black" />
        <circle cx="40" cy="20" r="5" fill="black" />
        <circle cx="30" cy="30" r="5" fill="black" />
        <circle cx="20" cy="40" r="5" fill="black" />
        <circle cx="40" cy="40" r="5" fill="black" />
      </template>

      <template v-if="value === 6">
        <circle cx="20" cy="15" r="5" fill="black" />
        <circle cx="40" cy="15" r="5" fill="black" />
        <circle cx="20" cy="30" r="5" fill="black" />
        <circle cx="40" cy="30" r="5" fill="black" />
        <circle cx="20" cy="45" r="5" fill="black" />
        <circle cx="40" cy="45" r="5" fill="black" />
      </template>
    </svg>
  </div>
</template>
<style scoped>
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

.retro-box {
  @apply rounded-lg border-4 border-8 border-primary bg-base-100 shadow-md;
  transition: all 0.3s ease;
}
</style>
