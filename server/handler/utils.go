package handler

import (
	"errors"
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
