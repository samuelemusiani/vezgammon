package db

import (
	"database/sql"
	"testing"
	"vezgammon/server/types"

	"gotest.tools/v3/assert"
)

var tournament types.Tournament
var u1, u2 types.User

func TestCreateTournament(t *testing.T) {
	u1 = types.User{
		Username:  "tournament1",
		Firstname: "tournament1",
		Lastname:  "tournament1",
		Mail:      "tournament1@mail.it",
	}

	u2 = types.User{
		Username:  "tournament2",
		Firstname: "tournament2",
		Lastname:  "tournament2",
		Mail:      "tournament2@mail.it",
	}

	password := "fgdfdfb"

	retuser1, err := CreateUser(u1, password)
	assert.NilError(t, err)

	retuser2, err := CreateUser(u2, password)
	assert.NilError(t, err)

	tournament = types.Tournament{
		Name:   "Tournament1",
		Owner:  retuser1.ID,
		Status: types.TournamentStatusInProgress,
		Users:  []int64{retuser1.ID, retuser2.ID},
	}

	rettour, err := CreateTournament(tournament)
	assert.NilError(t, err)
	assert.Assert(t, rettour.ID != tournament.ID)
	tournament.ID = rettour.ID
}

func TestUpdateTournament(t *testing.T) {
	tournament.Status = types.TournamentStatusEnded

	err := UpdateTournament(&tournament)
	assert.NilError(t, err)
}

func TestGetTournament(t *testing.T) {
	tour, err := GetTournament(tournament.ID)
	assert.NilError(t, err)
	assert.Equal(t, tour.ID, tournament.ID)
	assert.Equal(t, tour.Name, tournament.Name)
	assert.Equal(t, tour.Status, tournament.Status)
	assert.DeepEqual(t, tour.Users, tournament.Users)
}

func TestCreateGameTournament(t *testing.T) {
	g := types.Game{
		Tournament: sql.NullInt64{Int64: tournament.ID, Valid: true},
		Player1:    tournament.Users[0],
		Elo1:       1000,
		Player2:    tournament.Users[1],
		Elo2:       1000,
		Status:     types.GameStatusOpen,
	}

	_, err := CreateGame(g)
	assert.NilError(t, err)
}

func TestGetAllTournamentGames(t *testing.T) {
	tours, err := GetAllTournamentGames(tournament.ID)
	assert.NilError(t, err)
	assert.Assert(t, len(tours) > 0)
	assert.Equal(t, tours[0].Player1, tournament.Users[0])
}

func TestTournamentToReturnTournament(t *testing.T) {
	rett, err := TournamentToReturnTournament(tournament)
	assert.NilError(t, err)

	assert.Equal(t, rett.Owner, u1.Username)
	assert.Equal(t, rett.Users[0], u1.Username)
	assert.Equal(t, rett.Users[1], u2.Username)
}
