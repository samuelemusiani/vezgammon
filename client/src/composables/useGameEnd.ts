import { ref } from 'vue'
import { useSound } from '@vueuse/sound'
import victorySfx from '@/utils/sounds/victory.mp3'
import lostSfx from '@/utils/sounds/lostgame.mp3'
import type { User } from '@/utils/types'

export function useGameEnd() {
  const showResultModal = ref(false)
  const isWinner = ref(false)
  const isExploding = ref(false)

  const { play: playVictory } = useSound(victorySfx)
  const { play: playLost } = useSound(lostSfx)

  const fetchWinner = async (): Promise<string | null> => {
    try {
      const res = await fetch('/api/play/last/winner')
      const winner = await res.json()
      return winner
    } catch (err) {
      console.error('Error fetching winner:', err)
      return null
    }
  }

  const handleWin = () => {
    isWinner.value = true
    showResultModal.value = true
    playVictory()
    isExploding.value = true
    setTimeout(() => {
      isExploding.value = false
    }, 5000)
  }

  const handleLose = () => {
    isWinner.value = false
    showResultModal.value = true
    playLost()
  }

  const handleEnd = async (currentUser: User | undefined) => {
    const winner = await fetchWinner()
    if (winner === currentUser?.username) {
      handleWin()
    } else {
      handleLose()
    }
  }

  const handleRetire = async () => {
    try {
      const res = await fetch('/api/play/', {
        method: 'DELETE',
      })

      if (!res.ok) {
        console.error('Error retiring:', res)
      }

      handleLose()
    } catch (err) {
      console.error('Error exiting game:', err)
    }
  }

  return {
    showResultModal,
    isWinner,
    isExploding,
    handleEnd,
    handleLose,
    handleRetire,
  }
}
