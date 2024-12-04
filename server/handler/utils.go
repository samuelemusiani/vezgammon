package handler

import (
	"database/sql"
	"errors"
	"math"
	"vezgammon/server/db"
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

func tournamentMatchCreate(user1, user2 int64, tournament sql.NullInt64) error {
	err := matchmaking.CreateGame(user1, user2, tournament)
	if err != nil {
		return err
	}

	ws.GameTournamentReady(user1)
	ws.GameTournamentReady(user2)

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

		// start finals last
		err = tournamentMatchCreate(tournament.Winners[0], tournament.Winners[1], sql.NullInt64{Valid: true, Int64: tournament.ID})
		if err != nil {
			return err
		}
	}

	return nil
}

func tournamentGameEndHandler(tournamentId int64, winnerId int64) error {
	tournament, err := db.GetTournament(tournamentId)
	if err != nil {
		return err
	}

	tournament.Winners = append(tournament.Winners, winnerId)

	if len(tournament.Winners) == 4 {
		tournament.Status = types.TournamentStatusEnded
	}

	err = db.UpdateTournament(tournament)
	if err != nil {
		return err
	}

	if len(tournament.Winners) == 4 { // if tournament ended
		for _, user := range tournament.Users {
			ws.TournamentEnded(user)
		}
	} else {
		err = tournamentMatchCreator(tournament) // create next matches if needed
		if err != nil {
			return err
		}
	}
	return nil
}
