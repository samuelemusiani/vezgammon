import { ref } from 'vue'

export function useGameTimer() {
  const timeLeft = ref(60)
  const timerInterval = ref<ReturnType<typeof setTimeout> | null>(null)
  const isMyTurn = ref(false)

  const startTimer = () => {
    timeLeft.value = 60
    if (timerInterval.value) {
      clearInterval(timerInterval.value)
    }
    timerInterval.value = setInterval(() => {
      timeLeft.value--
      if (timeLeft.value <= 0) {
        stopTimer()
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
    isMyTurn,
    startTimer,
    stopTimer,
  }
}
