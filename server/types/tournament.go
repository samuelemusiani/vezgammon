package types

const TournamentStatusWaiting = "waiting"
const TournamentStatusInProgress = "in_progress"
const TournamentStatusEnded = "ended"

type Tournament struct {
	ID      int64   `json:"id" example:"1"`
	Name    string  `json:"name" example:"Tournament name"`
	Owner   int64   `json:"owner" example:"1"`
	Status  string  `json:"status" example:"open"`
	Users   []int64 `json:"users" example:"1,2,3"`
	Winners []int64 `json:"winners" example:"1,2,3"` // last winner at the end, used to calculate the next round and the leader board
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
	ID         int64  `json:"id" example:"1"`
	Name       string `json:"name" example:"Tournament name"`
	Owner      string `json:"owner" example:"Giorgio"`
	UserNumber int    `json:"user_number" example:"3"`
}

type TournamentList []TournamentInfo
