package types

const TournamentStatusOpen = "open"
const TournamentStatusClosed = "close"

const TournamentVisibilityPublic = "public"
const TournamentVisibilityPrivate = "private"

type Tournament struct {
	ID         int64   `json:"id" example:"1"`
	Name       string  `json:"name" example:"Tournament name"`
	Owner      int64   `json:"owner" example:"1"`
	Start      string  `json:"start" example:"2021-01-01T00:00:00Z"`
	End        string  `json:"end" example:"2021-01-01T00:00:00Z"`
	Status     string  `json:"status" example:"open"`
	Visibility string  `json:"visibility" example:"public"`
	AllowUsers []int64 `json:"allow_users" example:"1,2,3"`
	Users      []int64 `json:"users" example:"1,2,3"`
}

type LeaderBoardEntry struct {
	User string `json:"user" example:"Giorgio"`
	Win  int    `json:"win" example:"1"`
	Lose int    `json:"lose" example:"1"`
}

type ReturnTournament struct {
	ID         int64    `json:"id" example:"1"`
	Name       string   `json:"name" example:"Tournament name"`
	Owner      string   `json:"owner" example:"Giorgio"`
	Start      string   `json:"start" example:"2021-01-01T00:00:00Z"`
	End        string   `json:"end" example:"2021-01-01T00:00:00Z"`
	Status     string   `json:"status" example:"open"`
	Visibility string   `json:"visibility" example:"public"`
	AllowUsers []string `json:"allow_users" example:"giorgio,diego,marco"`
	Users      []string `json:"users" example:"giorgio,diego,marco"`

	LeaderBoard []LeaderBoardEntry `json:"leader_board"`
	Games       []ReturnGame       `json:"games"`
}

type TournamentInfo struct {
	ID    int64  `json:"id" example:"1"`
	Name  string `json:"name" example:"Tournament name"`
	Owner string `json:"owner" example:"Giorgio"`
}

type TournamentList []TournamentInfo
