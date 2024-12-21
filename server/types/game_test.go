package types

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestNewDices(t *testing.T) {
	dices := NewDices()
	assert.Assert(t, dices[0] >= 1 && dices[0] <= 6)
	assert.Assert(t, dices[1] >= 1 && dices[1] <= 6)
}

func TestSum(t *testing.T) {
	dices := Dices{3, 4}
	assert.Equal(t, dices.Sum(), 7)

	dices = Dices{6, 2}
	assert.Equal(t, dices.Sum(), 8)
}

func TestPlayMove(t *testing.T) {
	dices := Dices{3, 4}
	g := Game{
		ID:            1,
		Player1:       1,
		Elo1:          1000,
		Player2:       2,
		Elo2:          1030,
		Start:         time.Now(),
		End:           time.Now().Add(time.Duration(5) * time.Minute),
		Status:        GameStatusOpen,
		P1Checkers:    [25]int8{},
		P2Checkers:    [25]int8{},
		DoubleValue:   1,
		DoubleOwner:   GameDoubleOwnerAll,
		WantToDouble:  false,
		CurrentPlayer: GameCurrentPlayerP1,
		Dices:         dices,
	}

	// Set up initial board position
	// P1 Checkers (starting from bar - index 0)
	g.P1Checkers = [25]int8{
		0, // bar
		2, // point 1
		0, // point 2
		0, // point 3
		0, // point 4
		0, // point 5
		5, // point 6
		0, // point 7
		3, // point 8
		0, // point 9
		0, // point 10
		0, // point 11
		5, // point 12
		0, // point 13
		0, // point 14
		0, // point 15
		0, // point 16
		0, // point 17
		0, // point 18
		0, // point 19
		0, // point 20
		0, // point 21
		0, // point 22
		0, // point 23
		0, // point 24
	}

	// P2 Checkers (starting from bar - index 0)
	g.P2Checkers = [25]int8{
		0, // bar
		0, // point 24
		0, // point 23
		0, // point 22
		0, // point 21
		0, // point 20
		5, // point 19
		0, // point 18
		3, // point 17
		0, // point 16
		0, // point 15
		0, // point 14
		5, // point 13
		0, // point 12
		0, // point 11
		0, // point 10
		0, // point 9
		0, // point 8
		0, // point 7
		0, // point 6
		0, // point 5
		0, // point 4
		0, // point 3
		0, // point 2
		2, // point 1
	}

	moves := []Move{{From: 6, To: 9}, {From: 6, To: 10}}

	initialP1Checkers := g.P1Checkers

	g.PlayMove(moves)

	// Assertions for regular move
	assert.Equal(t, g.P1Checkers[6], initialP1Checkers[6]-2)
	assert.Equal(t, g.P1Checkers[9], initialP1Checkers[9]+1)
	assert.Equal(t, g.P1Checkers[10], initialP1Checkers[10]+1)

	// Test case 2: Move with hitting opponent's checker
	g = Game{ // Reset game state
		CurrentPlayer: GameCurrentPlayerP2,
		Dices:         Dices{3, 4},
	}

	g.P1Checkers = [25]int8{0}
	g.P2Checkers = [25]int8{0}
	g.P1Checkers[21] = 1
	g.P2Checkers[1] = 1

	moves = []Move{
		{From: 1, To: 4}, // Move that will hit opponent's checker
	}

	g.PlayMove(moves)

	// Assertions for hitting move
	assert.Equal(t, g.P2Checkers[1], int8(0))
	assert.Equal(t, g.P1Checkers[0], int8(1))
	assert.Equal(t, g.P1Checkers[22], int8(0))

	assert.Equal(t, g.CurrentPlayer, GameCurrentPlayerP1)
}

func TestToReturnGame(t *testing.T) {
	// Create test time values
	startTime := time.Now()
	endTime := startTime.Add(5 * time.Minute)

	// Initialize test game
	g := Game{
		ID:            123,
		Player1:       1,
		Elo1:          1500,
		Player2:       2,
		Elo2:          1600,
		Start:         startTime,
		End:           endTime,
		Status:        GameStatusOpen,
		P1Checkers:    [25]int8{1, 2, 3}, // Sample checker positions
		P2Checkers:    [25]int8{4, 5, 6}, // Sample checker positions
		DoubleValue:   2,
		DoubleOwner:   GameDoubleOwnerP1,
		WantToDouble:  true,
		CurrentPlayer: GameCurrentPlayerP1,
		Dices:         Dices{3, 4}, // This field isn't transferred to ReturnGame
	}

	// Test usernames
	username1 := "player_one"
	username2 := "player_two"

	// Convert to ReturnGame
	returnGame := g.ToReturnGame(username1, username2)

	assert.Equal(t, g.ID, returnGame.ID)
	assert.Equal(t, username1, returnGame.Player1)
	assert.Equal(t, g.Elo1, returnGame.Elo1)
	assert.Equal(t, username2, returnGame.Player2)
	assert.Equal(t, g.Elo2, returnGame.Elo2)
	assert.Equal(t, g.Start, returnGame.Start)
	assert.Equal(t, g.End, returnGame.End)
	assert.Equal(t, g.Status, returnGame.Status)
	assert.Equal(t, g.P1Checkers, returnGame.P1Checkers)
	assert.Equal(t, g.P2Checkers, returnGame.P2Checkers)
	assert.Equal(t, g.DoubleValue, returnGame.DoubleValue)
	assert.Equal(t, g.DoubleOwner, returnGame.DoubleOwner)
	assert.Equal(t, g.WantToDouble, returnGame.WantToDouble)
	assert.Equal(t, g.CurrentPlayer, returnGame.CurrentPlayer)
	assert.Equal(t, GameTypeOnline, returnGame.GameType)
}
