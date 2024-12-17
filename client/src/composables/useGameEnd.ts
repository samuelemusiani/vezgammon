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

  const handleWin = (playSound: boolean) => {
    isWinner.value = true
    showResultModal.value = true
    if (playSound) {
      playVictory()
    }
    isExploding.value = true
    setTimeout(() => {
      isExploding.value = false
    }, 5000)
  }

  const handleLose = (playSound: boolean, showResModal: boolean) => {
    isWinner.value = false
    showResultModal.value = showResModal
    if (playSound) {
      playLost()
    }
  }

  const handleEnd = async (
    currentUser: User | undefined,
    localGame: boolean,
  ) => {
    const winner = await fetchWinner()
    if (winner === currentUser?.username) {
      handleWin(!localGame)
    } else {
      handleLose(!localGame, true)
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

      // We do not need to play a sound here, as the server will send game_end
      // ws message and the sound will play
      handleLose(false, false)
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
