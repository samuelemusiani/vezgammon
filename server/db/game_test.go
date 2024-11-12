package db

import (
	"testing"
	"time"
	"vezgammon/server/types"

	"gotest.tools/v3/assert"
)

func TestCreateGame(t *testing.T) {
	u1 := types.User{
		Username:  "aa",
		Firstname: "aa",
		Lastname:  "aa",
		Mail:      "aa.aa@mail.it",
	}

	u2 := types.User{
		Username:  "bb",
		Firstname: "bb",
		Lastname:  "bb",
		Mail:      "bb.bb@mail.it",
	}

	password1 := "fgdfdfb"
	password2 := "asldfq34n"

	retuser1, err := CreateUser(u1, password1)
	assert.NilError(t, err)

	retuser2, err := CreateUser(u2, password2)
	assert.NilError(t, err)

	g := types.Game{
		Player1: retuser1.ID,
		Elo1:    1000,
		Player2: retuser2.ID,
		Elo2:    1000,
		Status:  types.GameStatusOpen,
	}

	retgame, err := CreateGame(g)
	assert.NilError(t, err)

	assert.Equal(t, g.Player1, retgame.Player1)
	assert.Equal(t, g.Player2, retgame.Player2)
	assert.Equal(t, g.Elo1, retgame.Elo1)
	assert.Equal(t, g.Elo2, retgame.Elo2)
	assert.Equal(t, g.Status, retgame.Status)
}

func TestGetGame(t *testing.T) {

	tuser := types.User{
		Username:  "tgame",
		Firstname: "tgame",
		Lastname:  "tame",
		Mail:      "tgame",
	}

	var err error
	tuser, err = CreateUser(tuser, "asjnoicjeofy23")
	assert.NilError(t, err)

	g := types.Game{
		Player1: tuser.ID,
		Elo1:    1000,
		Player2: tuser.ID,
		Elo2:    1000,
		Status:  types.GameStatusOpen,
		Start:   time.Now(),
		End:     time.Now(),
	}

	retgame, err := CreateGame(g)
	assert.NilError(t, err)

	retgame2, err := GetGame(retgame.ID)
	assert.NilError(t, err)

	assert.Equal(t, retgame.Player1, retgame2.Player1)
	assert.Equal(t, retgame.Player2, retgame2.Player2)
	assert.Equal(t, retgame.Elo1, retgame2.Elo1)
	assert.Equal(t, retgame.Elo2, retgame2.Elo2)
	assert.Equal(t, retgame.Status, retgame2.Status)
}
