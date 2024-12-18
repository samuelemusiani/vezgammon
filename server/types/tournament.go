package types

import "time"

const TournamentStatusWaiting = "waiting"
const TournamentStatusInProgress = "in_progress"
const TournamentStatusEnded = "ended"

type Tournament struct {
	ID     int64   `json:"id" example:"1"`
	Name   string  `json:"name" example:"Tournament name"`
	Owner  int64   `json:"owner" example:"1"`
	Status string  `json:"status" example:"open"`
	Users  []int64 `json:"users" example:"1,2,3"`
	// last winner at the end, used to calculate the next round and the leader board, position of users
	Winners      []int64   `json:"winners" example:"0,1,2,3"`
	CreationDate time.Time `json:"creation_date" example:"2021-09-01T00:00:00Z"`
}

type ReturnTournament struct {
	ID     int64    `json:"id" example:"1"`
	Name   string   `json:"name" example:"Tournament name"`
	Owner  string   `json:"owner" example:"Giorgio"`
	Status string   `json:"status" example:"open"`
	Users  []string `json:"users" example:"giorgio,diego,marco"`

	Games        []ReturnGame `json:"games"`
	CreationDate time.Time    `json:"creation_date" example:"2021-09-01T00:00:00Z"`
}

type TournamentInfo struct {
	ID           int64     `json:"id" example:"1"`
	Name         string    `json:"name" example:"Tournament name"`
	Owner        string    `json:"owner" example:"Giorgio"`
	UserNumber   int       `json:"user_number" example:"3"`
	CreationDate time.Time `json:"creation_date" example:"2021-09-01T00:00:00Z"`
	Status       string    `json:"status" example:"open"`
}

type TournamentList []TournamentInfo
