package db

import "vezgammon/server/types"

func MovesToArray(m types.Move) []int {
	return []int{m.From, m.To}
}

func MovesArrayToArray(m []types.Move) [][]int {
	var moves [][]int
	for _, move := range m {
		moves = append(moves, MovesToArray(move))
	}
	return moves
}

func ArrayToMoves(m []int) types.Move {
	return types.Move{From: m[0], To: m[1]}
}

func ArrayToMovesArray(m [][]int) []types.Move {
	var moves []types.Move
	for _, move := range m {
		moves = append(moves, ArrayToMoves(move))
	}
	return moves
}

func ArrayToDices(m []int) types.Dices {
	return types.Dices{m[0], m[1]}
}
