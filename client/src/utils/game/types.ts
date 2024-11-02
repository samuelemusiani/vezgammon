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
