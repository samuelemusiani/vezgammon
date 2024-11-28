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
