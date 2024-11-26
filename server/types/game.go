package types

import (
	"log/slog"
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

func (d Dices) Sum() int {
	return d[0] + d[1]
}

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

type NewGame struct {
	DicesP1 Dices      `json:"dices_p1"`
	DicesP2 Dices      `json:"dices_p2"`
	Game    ReturnGame `json:"game"`
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

	Dices      Dices `json:"dices"`
	Tournament int64 `json:"tournament"`
}

func (g *Game) PlayMove(moves []Move) {
	// if Move.To is 25, it means the checker is out of the board
	// if Move.From is 0, it means the checker is in the bar
	var myCheckers *[25]int8
	var opponentCheckers *[25]int8
	if g.CurrentPlayer == GameCurrentPlayerP1 {
		myCheckers = &g.P1Checkers
		opponentCheckers = &g.P2Checkers
		g.CurrentPlayer = GameCurrentPlayerP2
	} else {
		myCheckers = &g.P2Checkers
		opponentCheckers = &g.P1Checkers
		g.CurrentPlayer = GameCurrentPlayerP1
	}

	// estract next dices
	g.Dices = NewDices()

	for _, move := range moves {
		myCheckers[move.From]-- // remove checker from the source
		if move.To < 25 {
			myCheckers[move.To]++ // add checker to the destination

			if opponentCheckers[25-move.To] == 1 {
				opponentCheckers[25-move.To] = 0
				opponentCheckers[0]++
			}
		}
	}

	slog.With("checkers", myCheckers).Debug("Checkers")
	slog.With("checkers", opponentCheckers).Debug("Checkers")
}

const GameTypeLocal = "local"
const GameTypeBot = "bot"
const GameTypeOnline = "online"

type ReturnGame struct {
	ID int64 `json:"id"`
	// Username of the player
	Player1 string `json:"player1" example:"Giorgio"`
	Elo1    int64  `json:"elo1" example:"1000"`
	Player2 string `json:"player2" example:"Mario"`
	Elo2    int64  `json:"elo2" example:"1000"`

	Start  time.Time `json:"start" example:"2021-01-01T00:00:00Z"`
	End    time.Time `json:"end" example:"2021-01-01T00:00:00Z"`
	Status string    `json:"status" example:"open"`

	P1Checkers [25]int8 `json:"p1checkers"` // arr[0] is bar
	P2Checkers [25]int8 `json:"p2checkers"` // arr[0] is bar

	DoubleValue  uint64 `json:"double_value" example:"1"`
	DoubleOwner  string `json:"double_owner" example:"all"`
	WantToDouble bool   `json:"want_to_double" example:"false"`

	CurrentPlayer string `json:"current_player" example:"p1"`

	GameType string `json:"game_type" example:"online"`

	Tournament int64 `json:"tournament" example:"1"`
}
