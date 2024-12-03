<script setup lang="ts">
import Die from './Dice.vue'

interface DiceContainerProps {
  diceRolled: boolean
  displayedDice: number[]
  isRolling: boolean
  canRoll: boolean
}

defineProps<DiceContainerProps>()

defineEmits<{
  (e: 'roll'): void
}>()
</script>

<template>
  <div>
    <div v-show="canRoll" class="mb-4 flex justify-center">
      <button @click="$emit('roll')" class="retro-button">Roll Dice</button>
    </div>

    <div class="flex justify-center gap-4">
      <Die
        v-for="(value, index) in diceRolled ? displayedDice : []"
        :key="index"
        :value="value"
        :isRolling="isRolling"
      />
    </div>
  </div>
</template>
