package handler

import (
	"errors"
	"math"
	"vezgammon/server/types"
)

func getCurrentPlayer(currentPlayer string, id1, id2 int64) (int64, error) {
	var id int64

	switch currentPlayer {
	case types.GameCurrentPlayerP1:
		id = id1
	case types.GameCurrentPlayerP2:
		id = id2
	default:
		return 0, errors.New("Invalid current player")
	}

	return id, nil
}

func getOpponentID(currentPlayer string, id1, id2 int64) (int64, error) {
	return getCurrentPlayer(currentPlayer, id2, id1)
}

func invertPlayer(currentPlayer string) string {
	switch currentPlayer {
	case types.GameCurrentPlayerP1:
		return types.GameCurrentPlayerP2
	case types.GameCurrentPlayerP2:
		return types.GameCurrentPlayerP1
	default:
		return ""
	}
}

func calculateElo(elo1, elo2 int64, winner1 bool) (int64, int64) {
	diff := float64((elo2 - elo1) / 400)
	pow := math.Pow(10, diff)

	var w1 float64 = 0
	if winner1 {
		w1 = 1
	}

	ea := w1 - 1/(1+pow)

	K := 32
	elo1 += int64(float64(K) * ea)
	elo2 -= int64(float64(K) * ea)

	return elo1, elo2
}
