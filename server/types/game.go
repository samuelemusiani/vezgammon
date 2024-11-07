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
	User   string    `json:"user"`
	Time   time.Time `json:"time"`
	Moves  [][2]int  `json:"moves"`
}

type Game struct {
	ID      int64  `json:"id"`
	Player1 string `json:"player1"`
	Elo1    int64  `json:"elo1"`
	Player2 string `json:"player2"`
	Elo2    int64  `json:"elo2"`

	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Status string    `json:"status"`

	P1Checkers [25]int `json:"p1checkers"` // arr[0] is bar
	P2Checkers [25]int `json:"p2checkers"` // arr[0] is bar

	Double      uint64 `json:"double"`
	DoubleOwner string `json:"double-owner"`
}
