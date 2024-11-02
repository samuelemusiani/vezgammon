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

export const createWhiteChecker = (position: number, stackIndex: number): Checker => ({
  color: 'white',
  position,
  stackIndex,
})
export const createBlackChecker = (position: number, stackIndex: number): Checker => ({
  color: 'black',
  position,
  stackIndex,
})

export const createDefaultBoard = (): Checker[] => {
  let tmp = []
  for (let i = 0; i < 5; i++) {
    if (i < 2) {
      tmp.push(createWhiteChecker(0, i))
      tmp.push(createBlackChecker(23, i))
    }
    if (i < 3) {
      tmp.push(createBlackChecker(7, i))
      tmp.push(createWhiteChecker(16, i))
    }
    tmp.push(createBlackChecker(5, i))
    tmp.push(createWhiteChecker(18, i))
    tmp.push(createWhiteChecker(11, i))
    tmp.push(createBlackChecker(12, i))
  }
  return tmp
}

// Index follows the counter-clockwise direction: upper from right (0) to left (11),
// lower from left (0) to right (11)
export const getTrianglePath = (position: number): string => {
  const isUpper = position < 12
  const index = isUpper ? 11 - position : position - 12
  let x = BOARD.padding + index * BOARD.triangleWidth

  if (index >= 6) {
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

// position from 0 (upper right) to 23 (lower right)
export const getCheckerX = (position: number) => {
  const index = position < 12 ? 11 - position : position - 12
  let x = BOARD.padding + index * BOARD.triangleWidth + BOARD.triangleWidth / 2

  if (index >= 6) {
    x += BOARD.centerBarWidth * 2
  }

  return x
}

export const getCheckerY = (position: number, stackIndex: number) => {
  const spacing = BOARD.checkerRadius * 1.8

  if (position < 12) {
    return BOARD.padding + BOARD.checkerRadius + stackIndex * spacing
  } else {
    return (
      BOARD.height - BOARD.padding - BOARD.checkerRadius - stackIndex * spacing
    )
  }
}
