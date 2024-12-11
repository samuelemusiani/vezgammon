package handler

import (
	"database/sql"
	"log/slog"
	"testing"
	"time"
	"vezgammon/server/db"
	"vezgammon/server/matchmaking"
	"vezgammon/server/types"
	"vezgammon/server/ws"

	"gotest.tools/v3/assert"
)

func TestEloCalculation(t *testing.T) {
	elo1, elo2 := calculateElo(1000, 1000, true)
	assert.Equal(t, elo1, int64(1016))
	assert.Equal(t, elo2, int64(984))
}

func TestGetCurrentPlayer(t *testing.T) {
	id, err := getCurrentPlayer("p1", 1, 2)
	assert.NilError(t, err)
	assert.Equal(t, id, int64(1))

	id, err = getCurrentPlayer("p2", 1, 2)
	assert.NilError(t, err)
	assert.Equal(t, id, int64(2))

	id, err = getCurrentPlayer("p3", 1, 2)
	assert.ErrorContains(t, err, "Invalid current player")
}

func TestInvertPlayer(t *testing.T) {
	assert.Equal(t, invertPlayer("p1"), "p2")
	assert.Equal(t, invertPlayer("p2"), "p1")
	assert.Equal(t, invertPlayer("p3"), "")
}

var tournament = types.Tournament{
	ID:           1,
	Name:         "Tournament name",
	Owner:        1,
	Status:       types.TournamentStatusInProgress,
	Users:        []int64{1, 2, 3, 4},
	Winners:      []int64{},
	CreationDate: time.Now(),
}

func TestGetTournamentUserIndex(t *testing.T) {
	index := getTournamentUserIndex(&tournament, int64(2))
	assert.Equal(t, index, int64(1))

	index = getTournamentUserIndex(&tournament, int64(5))
	assert.Equal(t, index, int64(-1))
}

func TestGetTournamentIndexUser(t *testing.T) {
	index := getTournamentIndexUser(&tournament, int64(2))
	assert.Equal(t, index, int64(3))
}

func TestBotVSBotGame(t *testing.T) {
	matchmaking.Init(db.GetDatabase(), ws.GetWebsocket())
	var game *types.Game
	bot := db.GetHardBotID()
	slog.With("bot", bot).Debug("Bot ID")

	game, err := matchmaking.CreateGame(bot, bot, sql.NullInt64{Valid: false})
	assert.NilError(t, err)

	slog.With("game", game).Error("NULLABLE")

	game, err = db.CreateGame(*game)
	assert.NilError(t, err)

	err = botVsBotGame(game)
	assert.NilError(t, err)

	game, err = db.GetGame(game.ID)
	assert.NilError(t, err)
	assert.Assert(t, game.Status != types.GameStatusOpen)
}

func TestIsGameEnded(t *testing.T) {
	game := types.Game{
		Player1:    111,
		Player2:    222,
		Status:     types.GameStatusWinP1,
		P1Checkers: [25]int8{1},
		P2Checkers: [25]int8{0},
	}

	isended, winner := isGameEnded(&game)
	assert.Equal(t, isended, true)
	assert.Equal(t, winner, int64(222))

	game = types.Game{
		Player1:    111,
		Player2:    222,
		Status:     types.GameStatusWinP1,
		P1Checkers: [25]int8{1},
		P2Checkers: [25]int8{1},
	}
	isended, winner = isGameEnded(&game)
	assert.Equal(t, isended, false)
	assert.Equal(t, winner, int64(0))
}
