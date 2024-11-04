import { describe, it, expect } from 'vitest'
import { calculatePossibleMoves, isCheckerMovable, isValidMove } from './moves'
import type { GameState, Checker } from './types'
import { createDefaultBoard } from './game'

describe('calculatePossibleMoves', () => {
  it('all dice', () => {
    const gameState: GameState = {
      currentPlayer: 'white',
      dice: {
        value: [3, 5],
        used: [false, false],
        double: false,
      },
      board: createDefaultBoard(),
    }

    const checker: Checker = {
      color: 'white',
      position: 0,
      stackIndex: 0,
    }

    const result = calculatePossibleMoves(gameState, checker)
    expect(result).toEqual([3, 8])
  })

  it('one dice on other player triangle', () => {
    const gameState: GameState = {
      currentPlayer: 'white',
      dice: {
        value: [3, 5],
        used: [true, false],
        double: false,
      },

      board: createDefaultBoard(),
    }

    const checker: Checker = {
      color: 'white',
      position: 0,
      stackIndex: 0,
    }

    const result = calculatePossibleMoves(gameState, checker)
    expect(result).toEqual([])
  })

  it('other direction (opponent)', () => {
    const gameState: GameState = {
      currentPlayer: 'black',
      dice: {
        value: [3, 5],
        used: [false, false],
        double: false,
      },
      board: createDefaultBoard(),
    }

    const checker: Checker = {
      color: 'black',
      position: 23,
      stackIndex: 0,
    }

    const result = calculatePossibleMoves(gameState, checker)
    expect(result).toEqual([20, 15])
  })

  it('no moves when dice are used', () => {
    const gameState: GameState = {
      currentPlayer: 'white',
      dice: {
        value: [3, 5],
        used: [true, true],
        double: false,
      },
      board: createDefaultBoard(),
    }

    const checker: Checker = {
      color: 'white',
      position: 0,
      stackIndex: 0,
    }

    const result = calculatePossibleMoves(gameState, checker)
    expect(result).toEqual([])
  })
})
describe('isValidMove with hitting', () => {
  it('allows hitting single opponent checker', () => {
    const gameState: GameState = {
      currentPlayer: 'white',
      dice: {
        value: [3, 5],
        used: [false, false],
        double: false,
      },
      board: [
        { color: 'white', position: 0, stackIndex: 0 },
        { color: 'black', position: 3, stackIndex: 0 },
      ],
    }

    const checker: Checker = {
      color: 'white',
      position: 0,
      stackIndex: 0,
    }

    expect(isValidMove(gameState, checker, 3)).toBe(true)
  })

  it('no hit if there are two or more opponent checker', () => {
    const gameState: GameState = {
      currentPlayer: 'white',
      dice: {
        value: [3, 5],
        used: [false, false],
        double: false,
      },
      board: [
        { color: 'white', position: 0, stackIndex: 0 },
        { color: 'black', position: 3, stackIndex: 0 },
        { color: 'black', position: 3, stackIndex: 1 },
      ],
    }

    const checker: Checker = {
      color: 'white',
      position: 0,
      stackIndex: 0,
    }

    expect(isValidMove(gameState, checker, 3)).toBe(false)
  })
})

describe('hit checker logic', () => {
  it('only allows moving hit checker first', () => {
    const gameState: GameState = {
      currentPlayer: 'black',
      dice: {
        value: [4, 6],
        used: [false, false],
        double: false,
      },
      board: [
        { color: 'black', position: 24, stackIndex: 0 },
        { color: 'black', position: 23, stackIndex: 0 },
        { color: 'black', position: 0, stackIndex: 0 },
      ],
    }

    const normalChecker: Checker = {
      color: 'black',
      position: 23,
      stackIndex: 0,
    }

    const hitChecker: Checker = {
      color: 'black',
      position: 24,
      stackIndex: 0,
    }

    expect(isCheckerMovable(gameState, normalChecker)).toBe(false)
    expect(isCheckerMovable(gameState, hitChecker)).toBe(true)
  })
})
