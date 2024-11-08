package db

import (
	"testing"
	"vezgammon/server/types"

	"gotest.tools/v3/assert"
)

func TestCreateGame(t *testing.T) {
	u := types.User{
		Username:  "aa",
		Firstname: "bb",
		Lastname:  "cc",
		Mail:      "aa.bb@mail.it",
	}

	password := "fgdfdfb"

	retuser, err := CreateUser(u, password)
	if err != nil {
		t.Fatalf("not creating user %s", err)
	}

	g := types.Game{
		Player1: retuser.ID,
		Elo1:    1000,
		Player2: retuser.ID,
		Elo2:    1000,
		Status:  types.GameStatusOpen,
	}

	retgame, err := CreateGame(g)
	assert.NilError(t, err)

	assert.Assert(t, retgame.ID != g.ID)
}

func TestGetGame(t *testing.T) {

	tuser := types.User{
		Username:  "tgame",
		Firstname: "tgame",
		Lastname:  "tame",
		Mail:      "tgame",
	}
	_, err := CreateUser(tuser, "dsdf")
	assert.NilError(t, err)

	g := types.Game{
		Player1: tuser.ID,
		Elo1:    1000,
		Player2: tuser.ID,
		Elo2:    1000,
		Status:  types.GameStatusOpen,
	}

	retgame, err := CreateGame(g)
	assert.NilError(t, err)

	retgame2, _, err := GetGame(retgame.ID)
	assert.NilError(t, err)

	assert.Equal(t, g.Player1, retgame2.Player1)
}
