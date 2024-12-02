export interface BoardDimensions {
  width: number
  height: number
  triangleWidth: number
  triangleHeight: number
  centerBarWidth: number
  padding: number
  checkerRadius: number
}

export interface Checker {
  color: string
  position: number
  stackIndex: number
}

export type Board = Checker[]

export interface GameState {
  id: number
  player1: string
  elo1: number
  player2: string
  elo2: number
  start: string
  end: string
  status: 'p1' | 'p2' | 'open' | 'winp1' | 'winp2'
  game_type: string
  p1checkers: number[]
  p2checkers: number[]
  double_value: number
  double_owner: string
  want_to_double: boolean
  current_player: string
}

export interface APIResponse {
  dices_p1: [number, number]
  dices_p2: [number, number]
  game: GameState
}

export interface Move {
  from: number
  to: number
}

export interface MovesResponse {
  dices: [number, number]
  can_double: boolean
  possible_moves: Move[][]
}

export interface DiceRoll {
  value: [number, number]
  used: [boolean, boolean] | [boolean, boolean, boolean, boolean]
  double: boolean
}
