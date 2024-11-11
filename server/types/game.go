package types

import (
	"math/rand"
	"time"
)

func NewDices() [2]int {
	var dices [2]int
	for i := 0; i < len(dices); i++ {
		dices[i] = rand.Intn(6) + 1
	}
	return dices
}

type Move [2]int
type Dices [2]int

type Turn struct {
	ID     int64     `json:"id"`
	GameId int64     `json:"game_id"`
	User   string    `json:"user"`
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

type Game struct {
	ID      int64 `json:"id"`
	Player1 int64 `json:"player1"`
	Elo1    int64 `json:"elo1"`
	Player2 int64 `json:"player2"`
	Elo2    int64 `json:"elo2"`

	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
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

type ReturnGame struct {
	ID      int64 `json:"id"`
	Player1 int64 `json:"player1"`
	Elo1    int64 `json:"elo1"`
	Player2 int64 `json:"player2"`
	Elo2    int64 `json:"elo2"`

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

func (g *Game) ToReturnGame() ReturnGame {
	return ReturnGame{
		ID:            g.ID,
		Player1:       g.Player1,
		Elo1:          g.Elo1,
		Player2:       g.Player2,
		Elo2:          g.Elo2,
		Start:         g.Start,
		End:           g.End,
		Status:        g.Status,
		P1Checkers:    g.P1Checkers,
		P2Checkers:    g.P2Checkers,
		DoubleValue:   g.DoubleValue,
		DoubleOwner:   g.DoubleOwner,
		WantToDouble:  g.WantToDouble,
		CurrentPlayer: g.CurrentPlayer,
	}
}
