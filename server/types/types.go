package types

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Mail      string `json:"mail"`
}

type Turn struct {
	GameId int64      `json:"game-id"`
	User   string     `json:"user"`
	Moves  [][2]uint8 `json:"moves"`
}

type Game struct {
	ID      int64  `json:"id"`
	Player1 string `json:"player1"`
	Elo1    int64  `json:"elo1"`
	Player2 string `json:"player2"`
	Elo2    int64  `json:"elo2"`

	Status string `json:"status"`

	P1Checkers [25]uint8 `json:"p1checkers"` // arr[0] is bar
	P2Checkers [25]uint8 `json:"p2checkers"` // arr[0] is bar

	Double      uint64 `json:"double"`
	DoubleOwner string `json:"double-owner"`
}
