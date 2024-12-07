import { ref } from 'vue'
import { useSound } from '@vueuse/sound'
import diceSfx from '@/utils/sounds/dice.mp3'
import type { MovesResponse } from '@/utils/game/types'
import { useWebSocketStore } from '@/stores/websocket'

export function useDiceRoll() {
  const isRolling = ref(false)
  const diceRolled = ref(false)
  const displayedDice = ref<number[]>([])
  const webSocketStore = useWebSocketStore()
  const { play: playDice } = useSound(diceSfx)

  const resetDiceState = () => {
    isRolling.value = false
    diceRolled.value = false
    displayedDice.value = []
  }

  const handleDiceRoll = (
    availableMoves: MovesResponse | null,
    online: boolean,
    onRollComplete?: () => void,
  ) => {
    if (diceRolled.value || !availableMoves?.dices) return

    isRolling.value = true
    diceRolled.value = true
    playDice()

    const rollInterval = setInterval(() => {
      displayedDice.value = [
        Math.floor(Math.random() * 6) + 1,
        Math.floor(Math.random() * 6) + 1,
      ]
    }, 100)

    setTimeout(() => {
      clearInterval(rollInterval)
      isRolling.value = false
      displayedDice.value = availableMoves.dices
      if (onRollComplete) {
        onRollComplete()
      }
      if (online) {
        webSocketStore.sendMessage({
          type: 'dice_rolled',
          payload: JSON.stringify(availableMoves.dices),
        })
      }
    }, 1000)
  }

  const showDiceFromOpponent = (dices: number[]) => {
    isRolling.value = true
    diceRolled.value = true
    playDice()

    const rollInterval = setInterval(() => {
      displayedDice.value = [
        Math.floor(Math.random() * 6) + 1,
        Math.floor(Math.random() * 6) + 1,
      ]
    }, 100)

    setTimeout(() => {
      clearInterval(rollInterval)
      isRolling.value = false
      displayedDice.value = dices
      console.log('displayedDice', displayedDice.value, diceRolled.value)
    }, 1000)
  }

  return {
    isRolling,
    diceRolled,
    displayedDice,
    handleDiceRoll,
    resetDiceState,
    showDiceFromOpponent,
  }
}
