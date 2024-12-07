import { ref } from 'vue'
import { useGameEnd } from './useGameEnd'

export function useGameTimer() {
  const timeLeft = ref(60)
  const timerInterval = ref<ReturnType<typeof setTimeout> | null>(null)
  const { handleRetire } = useGameEnd()

  const startTimer = () => {
    timeLeft.value = 60
    if (timerInterval.value) {
      clearInterval(timerInterval.value)
    }
    timerInterval.value = setInterval(() => {
      timeLeft.value--
      if (timeLeft.value <= 0) {
        stopTimer()
        handleRetire()
      }
    }, 1000)
  }

  const stopTimer = () => {
    if (timerInterval.value) {
      clearInterval(timerInterval.value)
      timerInterval.value = null
    }
  }

  return {
    timeLeft,
    startTimer,
    stopTimer,
  }
}
