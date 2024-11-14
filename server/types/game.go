package types

import (
	"math/rand"
	"time"
)

func NewDices() Dices {
	var dices Dices
	for i := 0; i < len(dices); i++ {
		dices[i] = rand.Intn(6) + 1
	}
	return dices
}

type Move struct {
	From int64 `json:"from" example:"1"`
	To   int64 `json:"to" example:"2"`
}

type Dices [2]int

// if Double is true it menas that is a fake turn
// where the player doubles the Dices.
// Dices and Moves are empty in this case
type Turn struct {
	ID     int64     `json:"id"`
	GameId int64     `json:"game_id"`
	User   int64     `json:"user"`
	Time   time.Time `json:"time"`
	Dices  Dices     `json:"dices"`
	Double bool      `json:"double"`
	Moves  []Move    `json:"moves"`
}

type FutureTurn struct {
	Dices Dices `json:"dices"`
	// True if the player can double the red dice
	CanDouble     bool     `json:"can_double"`
	PossibleMoves [][]Move `json:"possible_moves"`
}

const GameStatusOpen = "open"
const GameStatusWinP1 = "winp1"
const GameStatusWinP2 = "winp2"

const GameDoubleOwnerAll = "all"
const GameDoubleOwnerP1 = "p1"
const GameDoubleOwnerP2 = "p2"

const GameCurrentPlayerP1 = "p1"
const GameCurrentPlayerP2 = "p2"

type Game struct {
	ID int64 `json:"id"`
	// ID of the player
	Player1 int64 `json:"player1"`
	Elo1    int64 `json:"elo1"`
	Player2 int64 `json:"player2"`
	Elo2    int64 `json:"elo2"`

	Start  time.Time `json:"start"`
	End    time.Time `json:"endtime"`
	Status string    `json:"status"`

	P1Checkers [25]int8 `json:"p1checkers"` // arr[0] is bar
	P2Checkers [25]int8 `json:"p2checkers"` // arr[0] is bar

	// Current value of the red dice
	DoubleValue  uint64 `json:"double_value"`
	DoubleOwner  string `json:"double_owner"`
	WantToDouble bool   `json:"want_to_double"`

	CurrentPlayer string `json:"current_player"`

	Dices Dices `json:"dices"`
}

func (g *Game) PlayMove(moves *[]Move) {
	// if Move.To is 25, it means the checker is out of the board
	// if Move.From is 0, it means the checker is in the bar
	var checkers *[25]int8
	if g.CurrentPlayer == GameCurrentPlayerP1 {
		checkers = &g.P1Checkers

		// if the player is P1, the next player is P2
		g.CurrentPlayer = GameCurrentPlayerP2
	} else {
		checkers = &g.P2Checkers

		// if the player is P2, the next player is P1
		g.CurrentPlayer = GameCurrentPlayerP1
	}

	// estract next dices
	g.Dices = NewDices()

	for _, move := range *moves {
		checkers[move.From]-- // remove checker from the source
		if move.To < 25 {
			checkers[move.To]++ // add checker to the destination
		}
	}
}

type ReturnGame struct {
	ID int64 `json:"id"`
	// Username of the player
	Player1 string `json:"player1"`
	Elo1    int64  `json:"elo1"`
	Player2 string `json:"player2"`
	Elo2    int64  `json:"elo2"`

	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Status string    `json:"status"`

	P1Checkers [25]int8 `json:"p1checkers"` // arr[0] is bar
	P2Checkers [25]int8 `json:"p2checkers"` // arr[0] is bar

	DoubleValue  uint64 `json:"double_value"`
	DoubleOwner  string `json:"double_owner"`
	WantToDouble bool   `json:"want_to_double"`

	CurrentPlayer string `json:"current_player"`
}
