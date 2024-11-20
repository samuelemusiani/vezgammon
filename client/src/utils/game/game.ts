import type { BoardDimensions, Checker } from './types'

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
  // Pedine mangiate al centro
  if (checker.position === 0) {
    return 400
  }

  // Per posizioni 1-24:
  // 1-12 sono nella metà inferiore (da destra a sinistra)
  // 13-24 sono nella metà superiore (da sinistra a destra)
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

export const getCheckerY = (checker: Checker) => {
  // Pedine mangiate al centro
  if (checker.position === 0) {
    if (checker.color === 'white') {
      return 100 + checker.stackIndex * BOARD.checkerRadius * 2
    } else {
      return 500 - checker.stackIndex * BOARD.checkerRadius * 2
    }
  }

  const spacing = BOARD.checkerRadius * 1.8

  // Posizioni 1-12 in basso, 13-24 in alto
  if (checker.position > 12) {
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
