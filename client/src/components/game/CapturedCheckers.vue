<script setup lang="ts">
interface CapturedCheckersProps {
  player: 'p1' | 'p2'
  checkerCount: number
  isHighlighted: boolean
}

defineProps<CapturedCheckersProps>()

defineEmits<{
  (e: 'click'): void
}>()
</script>

<template>
  <div
    class="captured-checkers-container flex flex-col place-items-center"
    :class="{
      'highlight-container': isHighlighted,
      'mb-4': player === 'p1',
      'mt-4': player === 'p2',
    }"
    @click="$emit('click')"
  >
    <div
      class="h-64 w-16 overflow-hidden rounded-lg border-2 border-amber-900 p-1"
    >
      <div class="flex flex-col gap-1">
        <div
          v-for="i in checkerCount"
          :key="`${player}-${i}`"
          class="h-3 w-full rounded-full"
          :class="{
            'border border-blue-500 bg-black': player === 'p1',
            'border border-black bg-white': player === 'p2',
          }"
        ></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.captured-checkers-container {
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

.captured-checkers-container {
  .h-64 {
    background: #f5c27a;
  }

  .w-full {
    transition: all 0.3s ease-out;
  }

  transition: all 0.3s ease;
}
</style>
