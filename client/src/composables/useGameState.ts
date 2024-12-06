import { ref } from 'vue'
import type { GameState, MovesResponse } from '@/utils/game/types'
import type { User } from '@/utils/types'

export function useGameState() {
  const gameState = ref<GameState | null>(null)
  const availableMoves = ref<MovesResponse | null>(null)
  const session = ref<User>()

  const fetchGameState = async () => {
    try {
      const res = await fetch('/api/play/')
      if (!res.ok) {
        return false
      }
      const data: GameState = await res.json()
      gameState.value = data
      return true
    } catch (err) {
      console.error('Error fetching game state:', err)
      return false
    }
  }

  const fetchMoves = async () => {
    try {
      const res = await fetch('/api/play/moves')
      const data: MovesResponse = await res.json()
      if (!res.ok) return false
      availableMoves.value = data
      return true
    } catch (err) {
      console.error('Error fetching moves:', err)
      return false
    }
  }

  const fetchSession = async () => {
    try {
      const res = await fetch('/api/session')
      const data = await res.json()
      session.value = data
    } catch (err) {
      console.error('Error fetching session:', err)
    }
  }

  return {
    gameState,
    availableMoves,
    session,
    fetchGameState,
    fetchMoves,
    fetchSession,
  }
}
