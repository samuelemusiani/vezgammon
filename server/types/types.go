package types

const DefaultElo = 1000

type User struct {
	ID        int64  `json:"id" example:"1"`
	Username  string `json:"username" example:"gio"`
	Firstname string `json:"firstname" example:"giorgio"`
	Lastname  string `json:"lastname" example:"rossi"`
	Mail      string `json:"mail" example:"giorossi@mail.it"`
	Elo       int64  `json:"elo" example:"1000"`
	IsBot     bool   `json:"is_bot" example:"false"`
}

type Stats struct {
	Gameplayed []ReturnGame `json:"games_played"`
	Won        int64        `json:"win"`
	Lost       int64        `json:"lost"`
	Elo        []int64      `json:"elo"`
	Winrate    float32      `json:"winrate"`
	Online     int64        `json:"online"`
	Local      int64        `json:"local"`
	Cpu        int64        `json:"cpu"`
	Tournament int64        `json:"tournament"`
}

type Badge struct {
	Bot        int `json:"bot_icon"`
	Homepieces int `json:"home_pieces"`
	Wongames   int `json:"win_games"`
	Elo        int `json:"elo"`
	Wontime    int `json:"win_time"`
	Gameplayed int `json:"game_played"`
}
