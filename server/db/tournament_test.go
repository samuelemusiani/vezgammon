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

func TestReturnTournamentToTournamentT(t *testing.T) {
	u := types.User{
		Username:  "ttA",
		Firstname: "ttA",
		Lastname:  "ttA",
		Mail:      "ttA@ttA.it",
	}

	password := "2354wdfs"

	retuser, err := CreateUser(u, password)
	assert.NilError(t, err)

	rt := types.ReturnTournament{
		ID:     1,
		Name:   "Tournament1",
		Owner:  retuser.Username,
		Users:  []string{retuser.Username},
		Status: types.TournamentStatusWaiting,
	}

	tournament, err := ReturnTournamentToTournament(rt)
	assert.NilError(t, err)
	assert.Equal(t, tournament.ID, rt.ID)
	assert.Equal(t, tournament.Name, rt.Name)
	assert.Equal(t, tournament.Status, rt.Status)
	assert.Equal(t, tournament.Users[0], retuser.ID)
	assert.Equal(t, tournament.Owner, retuser.ID)
}

func TestDeleteTournament(t *testing.T) {
	u := types.User{
		Username:  "ttB",
		Firstname: "ttB",
		Lastname:  "ttB",
		Mail:      "ttB@ttB.it",
	}

	password := "2354wdfs"

	retuser, err := CreateUser(u, password)
	assert.NilError(t, err)
	tournament = types.Tournament{
		Name:   "Tournament1",
		Owner:  retuser.ID,
		Status: types.TournamentStatusInProgress,
		Users:  []int64{retuser.ID},
	}
	rt, err := CreateTournament(tournament)
	assert.NilError(t, err)

	err = DeleteTournament(rt.ID)
	assert.NilError(t, err)
}

func TestGetTournamentList(t *testing.T) {
	u := types.User{
		Username:  "ttC",
		Firstname: "ttC",
		Lastname:  "ttC",
		Mail:      "ttC@ttC.it",
	}

	password := "2354wdfs"

	retuser, err := CreateUser(u, password)
	assert.NilError(t, err)
	tournament = types.Tournament{
		Name:   "Tournament1",
		Owner:  retuser.ID,
		Status: types.TournamentStatusInProgress,
		Users:  []int64{retuser.ID},
	}
	_, err = CreateTournament(tournament)
	assert.NilError(t, err)

	tlist, err := GetTournamentList()
	assert.NilError(t, err)
	assert.Assert(t, len(*tlist) > 0)
}
