<script setup lang="ts">
import Die from './Dice.vue'

interface DiceContainerProps {
  diceRolled: boolean
  displayedDice: number[]
  isRolling: boolean
  canRoll: boolean
  dicesReplay?: number[]
}

const props = defineProps<DiceContainerProps>()

function calculateDice() {
  if (props.dicesReplay) {
    return props.dicesReplay
  } else if (props.diceRolled) {
    return props.displayedDice
  }
  return []
}

defineEmits<{
  (e: 'roll'): void
}>()
</script>

<template>
  <div>
    <div
      v-show="canRoll && !dicesReplay && !diceRolled"
      class="mb-4 flex justify-center"
    >
      <button @click="$emit('roll')" class="retro-button">Roll Dice</button>
    </div>

    <div class="flex justify-center gap-4">
      <Die
        v-for="(value, index) in calculateDice()"
        :key="index"
        :value="value"
        :isRolling="isRolling"
      />
    </div>
  </div>
</template>
