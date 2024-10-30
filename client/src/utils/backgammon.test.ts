import { expect, test } from 'vitest'
import { newGame, rollDice } from './backgammon'

test('create new game', () => {
  const g = newGame()
  expect(g.starter).toBeNull()
  expect(g.moves.length).toBe(0)
  expect(g.winner).toBeNull()
})

test('dice rolls are always between 0 and 6', () => {
  for (let i = 0; i < 100; i++) {
    const r = rollDice()

    expect(r.first).toBeLessThanOrEqual(6)
    expect(r.first).toBeGreaterThanOrEqual(0)

    expect(r.second).toBeLessThanOrEqual(6)
    expect(r.second).toBeGreaterThanOrEqual(0)
  }
})
