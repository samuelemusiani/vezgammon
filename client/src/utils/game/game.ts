import type { BoardDimensions, Checker, GameState } from './types'
import router from '@/router'

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
const checkWin = (gameState: GameState) => {
  // controlla se le pedine sono tutte in posizione 23 o 0
  const checkers = gameState.board.filter(
    c => c.color === gameState.currentPlayer,
  )
  if (checkers.every(c => c.position === 23 || c.position === 0)) {
    alert(`${gameState.currentPlayer} wins!`)
    router.push('/')
  }
}
