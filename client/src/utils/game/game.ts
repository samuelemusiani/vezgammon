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

export const newGame = () => {
  const gameState: GameState = {
    currentPlayer: 'white',
    dice: { value: [0, 0], used: [false, false], double: false },
    board: createDefaultBoard(),
    capturedWhite: [],
    capturedBlack: [],
  }
  return gameState
}

export const createWhiteChecker = (
  position: number,
  stackIndex: number,
): Checker => ({
  color: 'white',
  position,
  stackIndex,
})
export const createBlackChecker = (
  position: number,
  stackIndex: number,
): Checker => ({
  color: 'black',
  position,
  stackIndex,
})

// Create the default board with checkers in their initial positions
export const createDefaultBoard = (): Checker[] => {
  const tmp = []
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
  // Hit checker
  if (position === -1 || position === 24) {
    return 400
  }

  const index = position < 12 ? 11 - position : position - 12
  let x = BOARD.padding + index * BOARD.triangleWidth + BOARD.triangleWidth / 2

  if (index >= 6) {
    x += BOARD.centerBarWidth * 2
  }

  return x
}

export const getCheckerY = (position: number, stackIndex: number) => {
  // Hit checker
  if (position === 24) {
    return 100 + stackIndex * BOARD.checkerRadius * 2
  }
  if (position === -1) {
    return 500 - stackIndex * BOARD.checkerRadius * 2
  }

  const spacing = BOARD.checkerRadius * 1.8

  if (position < 12) {
    return BOARD.padding + BOARD.checkerRadius + stackIndex * spacing
  } else {
    return (
      BOARD.height - BOARD.padding - BOARD.checkerRadius - stackIndex * spacing
    )
  }
}

// Function to update game state after a move has been made
export const updateGameState = (
  gameState: GameState,
  newCheckerPos: number,
  oldCheckerPos: number,
  movesAvailable: number,
) => {
  if (!gameState.dice) return
  if (oldCheckerPos === 24) {
    oldCheckerPos = 23
    newCheckerPos--
  }
  if (oldCheckerPos === -1) {
    oldCheckerPos = 0
    newCheckerPos++
  }
  const move = Math.abs(newCheckerPos - oldCheckerPos)

  const isDouble = gameState.dice.double

  if (isDouble) {
    const n = move / gameState.dice.value[0]

    for (let i = 0; i < n; i++) {
      const unusedIndex = gameState.dice.used.findIndex(used => !used)
      if (unusedIndex !== -1) {
        gameState.dice.used[unusedIndex] = true
        movesAvailable--
      }
    }
  } else {
    if (
      move === gameState.dice.value[0] + gameState.dice.value[1] &&
      !gameState.dice.used[0] &&
      !gameState.dice.used[1]
    ) {
      gameState.dice.used = [true, true]
      movesAvailable = 0
    } else if (move === gameState.dice.value[0] && !gameState.dice.used[0]) {
      gameState.dice.used[0] = true
      movesAvailable--
    } else if (move === gameState.dice.value[1] && !gameState.dice.used[1]) {
      gameState.dice.used[1] = true
      movesAvailable--
    }
  }

  checkWin(gameState)

  if (
    movesAvailable === 0 ||
    (isDouble
      ? gameState.dice.used.every(used => used)
      : gameState.dice.used.every(used => used))
  ) {
    endTurn(gameState)
  }
}

export const endTurn = (gameState: GameState) => {
  gameState.currentPlayer =
    gameState.currentPlayer === 'white' ? 'black' : 'white'
  gameState.dice.value = [0, 0]
}

// temp, here a player wins if all checkers are in position 23 or 0, in reality the game ends when all checkers are off the board
export const checkWin = (gameState: GameState): boolean => {
  const checkers = gameState.board.filter(
    c => c.color === gameState.currentPlayer,
  )
  if (checkers.every(c => c.position === 23 || c.position === 0)) {
    return true
  }
  return false
}
