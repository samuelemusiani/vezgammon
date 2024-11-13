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
  currentPlayer: string
  dice: DiceRoll
  board: Board
  time?: string
  capturedWhite: Checker[]
  capturedBlack: Checker[]
}

export interface DiceRoll {
  value: [number, number]
  used: [boolean, boolean] | [boolean, boolean, boolean, boolean]
  double: boolean
}
