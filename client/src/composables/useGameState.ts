import { ref } from 'vue'
import type { GameState, MovesResponse } from '@/utils/game/types'
import type { User } from '@/utils/types'

export function useGameState() {
  const gameState = ref<GameState | null>(null)
  const isMyTurn = ref(false)
  const availableMoves = ref<MovesResponse | null>(null)
  const session = ref<User>()
  const avatars = ref<string[]>([])

  const fetchGameState = async () => {
    try {
      if (!session.value) {
        await fetchSession()
      }
      const res = await fetch('/api/play/')
      if (!res.ok) {
        return false
      }
      const data: GameState = await res.json()
      console.log('Game state:', data)
      gameState.value = data
      isMyTurn.value =
        (gameState.value.current_player === 'p1' &&
          gameState.value.player1 === session.value?.username) ||
        (gameState.value.current_player === 'p2' &&
          gameState.value.player2 === session.value?.username)
      return true
    } catch (err) {
      console.error('Error fetching game state:', err)
      return false
    }
  }

  const fetchMoves = async () => {
    try {
      const res = await fetch('/api/play/moves')
      if (!res.ok) return false
      const data: MovesResponse = await res.json()
      console.log('Available moves:', data)
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
      console.log('Session:', data)
      session.value = data
    } catch (err) {
      console.error('Error fetching session:', err)
    }
  }

  const fetchAvatars = async () => {
    try {
      const res1 = await fetch(`/api/player/${gameState.value?.player1}/avatar`)
      const p1 = await res1.json()
      const res2 = await fetch(`/api/player/${gameState.value?.player2}/avatar`)
      const p2 = await res2.json()
      avatars.value = [p1, p2]
    } catch (err) {
      console.error('Error fetching avatars:', err)
    }
  }

  return {
    gameState,
    availableMoves,
    session,
    isMyTurn,
    avatars,
    fetchGameState,
    fetchMoves,
    fetchSession,
    fetchAvatars,
  }
}
