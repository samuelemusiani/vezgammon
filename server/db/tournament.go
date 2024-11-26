package db

import (
	"vezgammon/server/types"

	"github.com/lib/pq"
)

func InitTournament() error {
	q := `
	CREATE TABLE IF NOT EXISTS tournaments(
		id SERIAL PRIMARY KEY,
		name BPCHAR,
		owner INTEGER REFERENCES users(id),
		start TIMESTAMP,
		endtime TIMESTAMP,
		status BPCHAR DEFAULT 'open',
		visibility BPCHAR DEFAULT 'public',
		allows_users INTEGER [],
		users INTEGER []
	)
	`

	_, err := Conn.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func CreateTournament(t types.Tournament) (*types.Tournament, error) {
	q := `
	INSERT INTO tournaments(name, owner, start, status, visibility, allows_users, users)
	values($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	res := Conn.QueryRow(
		q,
		t.Name, t.Owner, t.Start, t.Status, t.Visibility, pq.Array(t.AllowUsers), pq.Array(t.Users),
	)

	var id int64

	err := res.Scan(&id)
	if err != nil {
		return nil, err
	}

	t.ID = id

	return &t, nil
}

func UpdateTournament(t *types.Tournament) error {
	q := `
	UPDATE tournaments
	SET status=$1, visibility=$2, allows_users=$3, users=$4
	WHERE id=$5
	`

	_, err := Conn.Exec(
		q,
		t.Status, t.Visibility, pq.Array(t.AllowUsers), pq.Array(t.Users), t.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func TournamentToReturnTournament(t types.Tournament) (*types.ReturnTournament, error) {
	var rt types.ReturnTournament

	rt.ID = t.ID
	rt.Name = t.Name
	rt.Start = t.Start
	rt.End = t.End
	rt.Status = t.Status
	rt.Visibility = t.Visibility

	// get usernames
	var owner string
	var users []string
	var allowUsers []string

	var user *types.User

	user, err := GetUser(t.Owner)
	if err != nil {
		return nil, err
	} else {
		owner = user.Username
	}

	for _, u := range t.Users {
		user, err := GetUser(u)
		if err != nil {
			return nil, err
		} else {
			users = append(users, user.Username)
		}
	}

	for _, u := range t.AllowUsers {
		user, err := GetUser(u)
		if err != nil {
			return nil, err
		} else {
			allowUsers = append(allowUsers, user.Username)
		}
	}

	rt.Owner = owner
	rt.AllowUsers = allowUsers
	rt.Users = users

	// get games
	games, err := GetAllTournamentGames(t.ID)

	for _, g := range games {
		rt.Games = append(rt.Games, *GameToReturnGame(&g))
	}

	// cal leaderboard
	// TODO:

	return &rt, nil
}

func GetTournament(id int64) (*types.ReturnTournament, error) {
	q := `
	SELECT *
	FROM tournaments
	WHERE id=$1
	`

	res := Conn.QueryRow(q, id)

	var t types.Tournament

	err := res.Scan(&t.ID, &t.Name, &t.Owner, &t.Start, &t.End, &t.Status, &t.Visibility, pq.Array(&t.AllowUsers), pq.Array(&t.Users))
	if err != nil {
		return nil, err
	}

	rt, err := TournamentToReturnTournament(t)
	if err != nil {
		return nil, err
	}

	return rt, nil
}

func GetAllTournamentGames(id int64) ([]types.Game, error) {
	q := `
	SELECT id
	FROM games
	WHERE tournament=$1
	`

	rows, err := Conn.Query(q, id)
	if err != nil {
		return nil, err
	}

	var games []types.Game

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		g, err := GetGame(id)
		if err != nil {
			return nil, err
		}

		games = append(games, *g)
	}

	return games, nil
}
