package handler

import (
	"database/sql"
	"errors"
	"vezgammon/server/matchmaking"
	"vezgammon/server/types"
	"vezgammon/server/ws"
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

func tournamentMatchCreate(user1, user2 int64, tournament sql.NullInt64) error {
	err := matchmaking.CreateGame(user1, user2, tournament)
	if err != nil {
		return err
	}

	err = ws.GameTournamentReady(user1)
	if err != nil {
		return err
	}
	err = ws.GameTournamentReady(user2)
	if err != nil {
		return err
	}

	return nil
}

func tournamentMatchCreator(tournament *types.Tournament) error {
	if len(tournament.Users) != 4 { // only 4 players supported
		return nil
	}

	var err error

	// no games, start the tournament
	if len(tournament.Winners) == 0 {
		err = tournamentMatchCreate(tournament.Users[0], tournament.Users[1], sql.NullInt64{Valid: true, Int64: tournament.ID})
		if err != nil {
			return err
		}

		err = tournamentMatchCreate(tournament.Users[2], tournament.Users[3], sql.NullInt64{Valid: true, Int64: tournament.ID})
		if err != nil {
			return err
		}
	}

	// start finals and third/fourth place match
	if len(tournament.Winners) == 2 {
		err = tournamentMatchCreate(tournament.Winners[0], tournament.Users[1], sql.NullInt64{Valid: true, Int64: tournament.ID})
		if err != nil {
			return err
		}

		// found third/fourth place users
		var losers []int64
		for _, user := range tournament.Users {
			if user != tournament.Winners[0] && user != tournament.Winners[1] {
				losers = append(losers, user)
			}
		}
		err = tournamentMatchCreate(losers[0], losers[1], sql.NullInt64{Valid: true, Int64: tournament.ID})
		if err != nil {
			return err
		}
	}

	return nil
}
