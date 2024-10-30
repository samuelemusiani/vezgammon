interface DiceRoll {
  first: number
  second: number
}

interface Movements {
  first: number | null
  second: number | null
}

interface Move {
  dice: DiceRoll
  movements: Movements
}

interface Game {
  starter: boolean | null
  moves: Move[]
  winner: boolean | null
}

export function newGame(): Game {
  return {
    starter: null,
    moves: [],
    winner: null
  } as Game
}

export function rollDice(): DiceRoll {
  const f = Math.floor(Math.random() * 6 + 1)
  const s = Math.floor(Math.random() * 6 + 1)

  return {
    first: f,
    second: s
  } as DiceRoll
}
