package handler

import (
	"database/sql"
	"errors"
	"log/slog"
	"math"
	"time"
	"vezgammon/server/bgweb"
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

	// if both players are bots
	if (db.GetBotLevel(user1) != 0) && (db.GetBotLevel(user2) != 0) {
		game, err := matchmaking.CreateGame(user1, user2, tournament)
		if err != nil {
			return err
		}

		go botVsBotGame(game) // giusto?
	} else if db.GetBotLevel(user1) != 0 {
		err, _, _ := createBotUserGame(user2, user1, tournament)
		if err != nil {
			return err
		}
		ws.GameTournamentReady(user2)
	} else if db.GetBotLevel(user2) != 0 {
		err, _, _ := createBotUserGame(user1, user2, tournament)
		if err != nil {
			return err
		}
		ws.GameTournamentReady(user1)
	} else {
		_, err := matchmaking.CreateGame(user1, user2, tournament)
		if err != nil {
			return err
		}

		ws.GameTournamentReady(user1)
		ws.GameTournamentReady(user2)
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

func botVsBotGame(g *types.Game) error {
	notabot := errors.New("not a bot")
	levelinvalid := errors.New("Invalid bot level")

	p1, _ := getCurrentPlayer(g.CurrentPlayer, g.Player1, g.Player2)
	p2, _ := getOpponentID(g.CurrentPlayer, g.Player1, g.Player2)

	botLevel1 := db.GetBotLevel(p1)
	var getMovefunc1 func(g *types.Game) ([]types.Move, error)

	botLevel2 := db.GetBotLevel(p2)
	var getMovefunc2 func(g *types.Game) ([]types.Move, error)

	if botLevel1 > 0 {
		switch botLevel1 {
		case 1:
			getMovefunc1 = bgweb.GetEasyMove
		case 2:
			getMovefunc1 = bgweb.GetMediumMove
		case 3:
			getMovefunc1 = bgweb.GetBestMove
		default:
			return levelinvalid
		}
	} else {
		return notabot
	}

	if botLevel2 > 0 {
		switch botLevel2 {
		case 1:
			getMovefunc2 = bgweb.GetEasyMove
		case 2:
			getMovefunc2 = bgweb.GetMediumMove
		case 3:
			getMovefunc2 = bgweb.GetBestMove
		default:
			return levelinvalid
		}
	} else {
		return notabot
	}

	gameEnded := false
	for !gameEnded {
		moves, err := getMovefunc1(g)
		if err != nil {
			return err
		}

		g.PlayMove(moves)
		err = db.UpdateGame(g)
		if err != nil {
			return err
		}

		// save turn
		turn := types.Turn{
			GameId: g.ID,
			User:   p1,
			Time:   time.Now(),
			Dices:  g.Dices,
			Double: false,
			Moves:  moves,
		}

		_, err = db.CreateTurn(turn)
		if err != nil {
			return err
		}

		isended, winner := isGameEnded(g)
		if isended {
			gameEnded = true
			err := endGame(g, winner)
			if err != nil {
				return err
			}
			break
		}

		moves, err = getMovefunc2(g)
		if err != nil {
			return err
		}

		g.PlayMove(moves)
		err = db.UpdateGame(g)
		if err != nil {
			return err
		}

		// save turn
		turn = types.Turn{
			GameId: g.ID,
			User:   p2,
			Time:   time.Now(),
			Dices:  g.Dices,
			Double: false,
			Moves:  moves,
		}

		_, err = db.CreateTurn(turn)
		if err != nil {
			return err
		}

		isended, winner = isGameEnded(g)
		if isended {
			gameEnded = true
			err := endGame(g, winner)
			if err != nil {
				return err
			}
			break
		}

	}

	return nil
}

// True if the game is ended. Return id of the winner
func isGameEnded(g *types.Game) (bool, int64) {
	s := func(a [25]int8) int {
		sum := 0
		for i := 0; i < 25; i++ {
			sum += int(a[i])
		}
		return sum
	}

	if s(g.P1Checkers) == 0 {
		return true, g.Player1
	}

	if s(g.P2Checkers) == 0 {
		return true, g.Player2
	}

	return false, 0
}

func endGame(g *types.Game, winnerID int64) error {
	if winnerID == g.Player1 {
		g.Status = types.GameStatusWinP1
	} else {
		g.Status = types.GameStatusWinP2
	}

	g.End = time.Now()

	err := db.UpdateGame(g)
	if err != nil {
		return err
	}

	rg := db.GameToReturnGame(g)

	if rg.GameType == types.GameTypeOnline {
		u1, err := db.GetUser(g.Player1)
		if err != nil {
			slog.With("error", err).Error("Getting user in endGame")
			return err
		}

		u2, err := db.GetUser(g.Player2)
		if err != nil {
			slog.With("error", err).Error("Getting user in endGame")
			return err
		}

		elo1, elo2 := calculateElo(u1.Elo, u2.Elo, winnerID == g.Player1)

		err = db.UpdateUserElo(g.Player1, elo1)
		if err != nil {
			slog.With("error", err).Error("Updating elo in endGame")
			return err
		}

		err = db.UpdateUserElo(g.Player2, elo2)
		if err != nil {
			slog.With("error", err).Error("Updating elo in endGame")
			return err
		}
	}

	if g.Tournament.Valid {
		err = tournamentGameEndHandler(g.Tournament.Int64, winnerID)
		if err != nil {
			return err
		}
	}

	ws.GameEnd(g.Player1)
	ws.GameEnd(g.Player2)

	return err
}

func createBotUserGame(userId, botId int64, tournament sql.NullInt64) (error, *types.Dices, *types.Dices) {
	var startdicesP1, startdicesP2 types.Dices
	for {
		startdicesP1 = types.NewDices()
		startdicesP2 = types.NewDices()

		if startdicesP1.Sum() != startdicesP2.Sum() {
			if startdicesP1.Sum() < startdicesP2.Sum() {
				startdicesP1, startdicesP2 = startdicesP2, startdicesP1
			}
			break
		}
	}

	// Against a bot the player will always start first
	var startPlayer = types.GameCurrentPlayerP1

	firstdices := types.NewDices()

	g := types.Game{
		Player1:       userId,
		Player2:       botId,
		Start:         time.Now(),
		Status:        types.GameStatusOpen,
		CurrentPlayer: startPlayer,
		Dices:         firstdices,
		Tournament:    tournament,
	}

	slog.With("game", g).Debug("Creating game")

	_, err := db.CreateGame(g)

	return err, &startdicesP1, &startdicesP2
}
