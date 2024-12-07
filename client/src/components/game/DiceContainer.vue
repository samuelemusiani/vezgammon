<script setup lang="ts">
import Die from './Dice.vue'
import { ref, watchEffect } from 'vue'

interface DiceContainerProps {
  diceRolled: boolean
  displayedDice: number[]
  isRolling: boolean
  canRoll: boolean
  dicesReplay?: number[]
}

const props = defineProps<DiceContainerProps>()
const currentDice = ref<number[]>([])

// Usa watchEffect per aggiornare currentDice quando cambiano le props
watchEffect(() => {
  console.log('Updating dice:', props.dicesReplay || props.displayedDice)
  currentDice.value = props.dicesReplay || props.displayedDice
})

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
        v-for="(value, index) in currentDice"
        :key="index"
        :value="value"
        :isRolling="isRolling"
      />
    </div>
    <div class="text-xs text-gray-500">
      Dice values: {{ currentDice }}
      <br />
      Rolled: {{ diceRolled }}
      <br />
      Rolling: {{ isRolling }}
    </div>
  </div>
</template>
