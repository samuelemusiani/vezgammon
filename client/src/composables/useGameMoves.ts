import { ref } from 'vue'
import type { Move, Checker } from '@/utils/game/types'
import { vfetch } from '@/utils/fetch'

export function useGameMoves() {
  const selectedChecker = ref<Checker | null>(null)
  const possibleMoves = ref<number[]>([])
  const movesToSubmit = ref<Move[]>([])

  const submitMoves = async () => {
    try {
      const res = await vfetch('/api/play/moves', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(movesToSubmit.value),
      })

      // Reset dello stato delle mosse
      movesToSubmit.value = []
      possibleMoves.value = []
      selectedChecker.value = null

      return res.ok
    } catch (err) {
      console.error('Error submitting moves:', err)
      return false
    }
  }

  return {
    selectedChecker,
    possibleMoves,
    movesToSubmit,
    submitMoves,
  }
}
