package types

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Mail      string `json:"mail"`
}

type Stats struct {
	User       User `json:"user"`
	GamePlayed int8 `json:"gameplayed"`
	Elo        int8 `json:"elo"`
}
