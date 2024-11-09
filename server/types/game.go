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

type Turn struct {
	ID     int64     `json:"id"`
	GameId int64     `json:"game-id"`
	User   int64     `json:"user"`
	Time   time.Time `json:"time"`
	Dices  [2]int    `json:"dices"`
	Double bool      `json:"double"`
	Moves  [][2]int  `json:"moves"`
}

var GameStatusOpen = "open"
var GameStatusWinP1 = "winp1"
var GameStatusWinP2 = "winp2"

var GameDoubleOwnerAll = "all"
var GameDoubleOwnerP1 = "p1"
var GameDoubleOwnerP2 = "p2"

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

	Double      uint64 `json:"double"`
	DoubleOwner string `json:"double-owner"`
}
