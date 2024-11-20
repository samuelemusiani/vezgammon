package db

import "vezgammon/server/types"

func MovesArrayToArray(m []types.Move) []int64 {
	var moves []int64
	for _, move := range m {
		moves = append(moves, move.From, move.To)
	}
	return moves
}

func ArrayToMovesArray(m []int64) []types.Move {
	var moves []types.Move
	for i := 0; i < len(m); i += 2 {
		moves = append(moves, types.Move{From: m[i], To: m[i+1]})
	}
	return moves
}

func ArrayToDices(m []int) types.Dices {
	return types.Dices{m[0], m[1]}
}
