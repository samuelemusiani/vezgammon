package types

const TournamentStatusOpen = "open"
const TournamentStatusClosed = "close"

type Tournament struct {
	ID     int64   `json:"id" example:"1"`
	Name   string  `json:"name" example:"Tournament name"`
	Owner  int64   `json:"owner" example:"1"`
	Status string  `json:"status" example:"open"`
	Users  []int64 `json:"users" example:"1,2,3"`
}

type LeaderBoardEntry struct {
	User string `json:"user" example:"Giorgio"`
	Win  int    `json:"win" example:"1"`
	Lose int    `json:"lose" example:"1"`
}

type LeaderBoard []LeaderBoardEntry

type ReturnTournament struct {
	ID     int64    `json:"id" example:"1"`
	Name   string   `json:"name" example:"Tournament name"`
	Owner  string   `json:"owner" example:"Giorgio"`
	Status string   `json:"status" example:"open"`
	Users  []string `json:"users" example:"giorgio,diego,marco"`

	LeaderBoard []LeaderBoardEntry `json:"leader_board"`
	Games       []ReturnGame       `json:"games"`
}

type TournamentInfo struct {
	ID    int64  `json:"id" example:"1"`
	Name  string `json:"name" example:"Tournament name"`
	Owner string `json:"owner" example:"Giorgio"`
}

type TournamentList []TournamentInfo
