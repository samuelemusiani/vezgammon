import type { Checker, GameState } from './types'

export const isCheckerMovable = (
  gameState: GameState,
  checker: Checker,
): boolean => {
  if (!gameState.dice.value[0] || !gameState.dice.value[1]) return false
  if (checker.color !== gameState.currentPlayer) return false

  // If there are checkers on the bar, the player must move them first
  const hitCheckers = gameState.board.filter(
    c =>
      c.color === gameState.currentPlayer &&
      (c.position === -1 || c.position === 24),
  )
  if (hitCheckers.length > 0) {
    console.log(hitCheckers)
    console.log(checker)
    return hitCheckers.some(
      c =>
        c.color === checker.color &&
        c.position === checker.position &&
        c.stackIndex === checker.stackIndex,
    )
  }

  const stack = gameState.board.filter(c => c.position === checker.position)
  const topChecker = stack[stack.length - 1]
  return checker === topChecker
}

// Function to calculate possible moves (position reached) for a selected checker
export const calculatePossibleMoves = (
  gameState: GameState,
  checker: Checker,
): number[] => {
  console.log(checker)
  const direction = checker.color === 'white' ? 1 : -1

  const dicesvalue = []
  for (const [i, v] of gameState.dice.used.entries()) {
    if (!v) dicesvalue.push(gameState.dice.value[i % 2])
  }

  let moves = dicesvalue.reduce<number[][]>(
    (subsets, value) =>
      subsets.concat(subsets.map(set => [value, ...set] as number[])),
    [[]],
  )

  moves = moves.filter(v => v.length != 0)

  const possibleMoves = moves.flatMap(
    a => checker.position + a.reduce((sum, v) => sum + v, 0) * direction,
  )

  return possibleMoves.filter(move => isValidMove(gameState, checker, move))
}

// Check if destination position isn't occupied by opponent checkers
export const isValidMove = (
  gameState: GameState,
  checker: Checker,
  targetPosition: number,
) => {
  if (targetPosition < 0 || targetPosition > 23) return false

  const targetCheckers = gameState.board.filter(
    c => c.position === targetPosition,
  )
  if (targetCheckers.length > 1 && targetCheckers[0].color !== checker.color) {
    return false
  }

  return true
}

// Update Stack Indices of checkers in the board
export const updateStackIndices = (gameState: GameState) => {
  const positions: { [key: number]: Checker[] } = {}
  // Group checkers by position (triangle)
  gameState.board.forEach(checker => {
    if (!positions[checker.position]) positions[checker.position] = []
    positions[checker.position].push(checker)
  })

  // Update stack index for each checker in the same position
  for (const stack of Object.values(positions)) {
    stack.forEach((checker, index) => {
      checker.stackIndex = index
    })
  }
}

export const handleHitChecker = (gameState: GameState, position: number) => {
  const targetChecker = gameState.board.find(c => c.position === position)

  // checker is hit
  if (targetChecker && targetChecker.color !== gameState.currentPlayer) {
    targetChecker.position = targetChecker.color === 'white' ? -1 : 24
    updateStackIndices(gameState)
  }
}
