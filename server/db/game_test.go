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

func TestCreateTurn(t *testing.T) {
	tuser := types.User{
		Username:  "tturn",
		Firstname: "tturn",
		Lastname:  "tturn",
		Mail:      "tturn",
	}

	var err error
	tuser, err = CreateUser(tuser, "tturn")
	assert.NilError(t, err)

	game := types.Game{
		Player1: tuser.ID,
		Elo1:    1000,
		Player2: tuser.ID,
		Elo2:    1000,
		Status:  types.GameStatusOpen,
		Start:   time.Now(),
		End:     time.Now(),
	}

	tgame, err := CreateGame(game)
	assert.NilError(t, err)

	turn := types.Turn{
		GameId: tgame.ID,
		User:   tuser.ID,
		Time:   time.Now(),
		Dices:  types.Dices{6, 6},
		Double: false,
		Moves:  []types.Move{{From: 1, To: 2}, {From: 2, To: 3}},
	}

	tturn, err := CreateTurn(turn)
	assert.NilError(t, err)

	assert.Equal(t, turn.GameId, tturn.GameId)
	assert.Equal(t, turn.User, tturn.User)
	assert.Equal(t, turn.Time, tturn.Time)
	assert.DeepEqual(t, turn.Dices, tturn.Dices)
	assert.Equal(t, turn.Double, tturn.Double)
	assert.DeepEqual(t, turn.Moves, tturn.Moves)
}

func TestGetTurns(t *testing.T) {
	tuser := types.User{
		Username:  "tturngets",
		Firstname: "tturngets",
		Lastname:  "tturngets",
		Mail:      "tturngets",
	}

	var err error
	tuser, err = CreateUser(tuser, "tturngets")
	assert.NilError(t, err)

	game := types.Game{
		Player1: tuser.ID,
		Elo1:    1000,
		Player2: tuser.ID,
		Elo2:    1000,
		Status:  types.GameStatusOpen,
		Start:   time.Now(),
		End:     time.Now(),
	}

	tgame, err := CreateGame(game)
	assert.NilError(t, err)

	turn1 := types.Turn{
		GameId: tgame.ID,
		User:   tuser.ID,
		Time:   time.Now(),
		Dices:  types.Dices{6, 6},
		Double: false,
		Moves:  []types.Move{{From: 4, To: 6}, {From: 2, To: 3}},
	}

	turn2 := types.Turn{
		GameId: tgame.ID,
		User:   tuser.ID,
		Time:   time.Now(),
		Dices:  types.Dices{6, 6},
		Double: false,
		Moves:  []types.Move{{From: 1, To: 3}, {From: 5, To: 7}},
	}

	tturn1, err := CreateTurn(turn1)
	assert.NilError(t, err)

	tturn2, err := CreateTurn(turn2)
	assert.NilError(t, err)

	var retarr []types.Turn
	retarr = append(retarr, *tturn1)
	retarr = append(retarr, *tturn2)

	retturns, err := GetTurns(tgame.ID)
	assert.NilError(t, err)

	for i := range retturns {
		retturns[i].Time = retarr[i].Time
	}
	assert.DeepEqual(t, retturns, retarr)
}

func TestGetLastTurn(t *testing.T) {
	tuser := types.User{
		Username:  "tturngetlast",
		Firstname: "tturngetlast",
		Lastname:  "tturngetlast",
		Mail:      "tturngetlast",
	}

	var err error
	tuser, err = CreateUser(tuser, "tturngetlast")
	assert.NilError(t, err)

	game := types.Game{
		Player1: tuser.ID,
		Elo1:    1000,
		Player2: tuser.ID,
		Elo2:    1000,
		Status:  types.GameStatusOpen,
		Start:   time.Now(),
		End:     time.Now(),
	}

	tgame, err := CreateGame(game)
	assert.NilError(t, err)

	turn1 := types.Turn{
		GameId: tgame.ID,
		User:   tuser.ID,
		Time:   time.Now(),
		Dices:  types.Dices{6, 6},
		Double: false,
		Moves:  []types.Move{{From: 4, To: 6}, {From: 2, To: 3}},
	}

	turn2 := types.Turn{
		GameId: tgame.ID,
		User:   tuser.ID,
		Time:   time.Now().Add(1),
		Dices:  types.Dices{4, 6},
		Double: false,
		Moves:  []types.Move{{From: 6, To: 7}, {From: 5, To: 7}},
	}

	_, err = CreateTurn(turn1)
	assert.NilError(t, err)

	tturn2, err := CreateTurn(turn2)
	assert.NilError(t, err)

	lastturn, err := GetLastTurn(tgame.ID)
	assert.NilError(t, err)

	tturn2.Time = lastturn.Time
	assert.DeepEqual(t, tturn2, lastturn)
}
