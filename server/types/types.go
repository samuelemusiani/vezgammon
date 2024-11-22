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
