import type {GameState} from "@/utils/game/types";

export interface User {
  id: number
  username: string
  firstname: string
  lastname: string
  mail: string
  is_bot: boolean
  elo: number
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
