import { test, expect } from 'vitest'
import * as game from './game'

test('test getTriangleColor', () => {
  expect(game.getTriangleColor(4)).toBe('#8B0000')
  expect(game.getTriangleColor(16)).toBe('#8B0000')
  expect(game.getTriangleColor(12)).toBe('#8B0000')
  expect(game.getTriangleColor(1)).toBe('#000080')
  expect(game.getTriangleColor(5)).toBe('#000080')
  expect(game.getTriangleColor(7)).toBe('#000080')
})
