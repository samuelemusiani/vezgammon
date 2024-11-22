import type { BoardDimensions, Checker, GameState } from './types'

export const BOARD: BoardDimensions = {
  width: 800,
  height: 600,
  triangleWidth: 60,
  triangleHeight: 250,
  centerBarWidth: 20,
  padding: 20,
  checkerRadius: 20,
}

// Position is 1 (lower right) to 24 (upper right)
export const getTrianglePath = (position: number): string => {
  const isUpper = position > 12
  let x

  if (isUpper) {
    // 13-24 in alto da sinitra verso destra
    x = BOARD.padding + (position - 13) * BOARD.triangleWidth
  } else {
    // 1-12 in basso da destra verso sinistra
    x = BOARD.padding + (12 - position) * BOARD.triangleWidth
  }

  // Aggiungi lo spazio per la barra centrale dopo i primi 6 triangoli
  if ((isUpper && position >= 19) || (!isUpper && position <= 6)) {
    x += BOARD.centerBarWidth * 2
  }

  const y = isUpper ? BOARD.padding : BOARD.height - BOARD.padding

  return isUpper
    ? `M ${x} ${y}
       L ${x + BOARD.triangleWidth} ${y}
       L ${x + BOARD.triangleWidth / 2} ${y + BOARD.triangleHeight} Z`
    : `M ${x} ${y}
       L ${x + BOARD.triangleWidth} ${y}
       L ${x + BOARD.triangleWidth / 2} ${y - BOARD.triangleHeight} Z`
}

export const getTriangleColor = (position: number): string => {
  return position % 2 === 0 ? '#8B0000' : '#000080'
}

export const getCheckerX = (checker: Checker) => {
  // Hit Checker in the center
  if (checker.position === 0) {
    return 400
  }

  const index =
    checker.position > 12
      ? checker.position - 13 // per posizioni 13-24
      : 12 - checker.position // per posizioni 1-12

  let x = BOARD.padding + index * BOARD.triangleWidth + BOARD.triangleWidth / 2

  if (index >= 6) {
    x += BOARD.centerBarWidth * 2
  }

  return x
}

const normalizePosition = (position: number, color: string) => {
  if (color === 'white' || position === 0) return position
  // Per le pedine nere, inverti la posizione (da 24-1 a 1-24)
  return 25 - position
}

export const getCheckerY = (checker: Checker, gameState: GameState) => {
  const normalizedPosition = normalizePosition(checker.position, checker.color)

  // Hit Checker in the center
  if (checker.position === 0) {
    if (checker.color === 'white') {
      return 100 + checker.stackIndex * BOARD.checkerRadius * 2
    } else {
      return 500 - checker.stackIndex * BOARD.checkerRadius * 2
    }
  }

  let totalCheckers = 0
  if (gameState) {
    if (checker.color === 'black') {
      totalCheckers = gameState.p1checkers[checker.position]
    } else {
      totalCheckers = gameState.p2checkers[checker.position]
    }
  }

  // y-spacing based on number of checkers in each triangle
  const spacing =
    totalCheckers > 4 ? BOARD.checkerRadius * 1.3 : BOARD.checkerRadius * 1.8

  // Posizioni 1-12 in basso, 13-24 in alto
  if (normalizedPosition > 12) {
    return BOARD.padding + BOARD.checkerRadius + checker.stackIndex * spacing
  } else {
    return (
      BOARD.height -
      BOARD.padding -
      BOARD.checkerRadius -
      checker.stackIndex * spacing
    )
  }
}
