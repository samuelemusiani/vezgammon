import type { GameState } from '@/utils/game/types'

export interface User {
  id: number
  username: string
  firstname: string
  lastname: string
  mail: string
  is_bot: boolean
  elo: number
  avatar: string
}

export interface WSMessage {
  type: string
  payload: string
}

export interface UserStats {
  user: string
  win: number
  lose: number
}

export interface Tournament {
  id: number
  name: string
  owner: string
  status: string
  leader_board: UserStats[]
  games: GameState[]
  users: string[]
  creation_date: string
}

export interface LeaderBoardUser {
  username: string
  elo: number
}

export interface GameStats {
  games_played: GameState[]
  win: number
  lost: number
  winrate: number
  elo: number[]
  cpu: number
  local: number
  online: number
  tournament: number
  leaderboard: LeaderBoardUser[]
}
