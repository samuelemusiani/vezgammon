import { describe, it, expect } from 'vitest'
import { calculatePossibleMoves } from './moves'
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
