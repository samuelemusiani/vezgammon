package handler

import (
	"errors"
	"vezgammon/server/db"
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

func getBotLevel(id int64) int {
	if id == db.GetEasyBotID() {
		return 1
	} else if id == db.GetMediumBotID() {
		return 2
	} else if id == db.GetHardBotID() {
		return 3
	}

	return 0
}
