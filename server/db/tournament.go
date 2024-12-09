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
		status BPCHAR DEFAULT 'waiting',
		users INTEGER [],
		winners INTEGER [] DEFAULT '{}'::INTEGER[],
		creation_date TIMESTAMP
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
	INSERT INTO tournaments(name, owner, status, users, creation_date)
	values($1, $2, $3, $4, $5)
	RETURNING id
	`

	res := Conn.QueryRow(
		q,
		t.Name, t.Owner, t.Status, pq.Array(t.Users), t.CreationDate,
	)

	var id int64

	err := res.Scan(&id)
	if err != nil {
		return nil, err
	}

	t.ID = id

	return &t, nil
}

func DeleteTournament(id int64) error {
	q := `
	DELETE FROM tournaments
	WHERE id=$1
	`

	_, err := Conn.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTournament(t *types.Tournament) error {
	q := `
	UPDATE tournaments
	SET status=$1, users=$2, winners=$3
	WHERE id=$4
	`

	_, err := Conn.Exec(
		q,
		t.Status, pq.Array(t.Users), pq.Array(t.Winners), t.ID,
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
	rt.Status = t.Status
	rt.CreationDate = t.CreationDate

	// get usernames
	var owner string
	var users []string

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

	rt.Owner = owner
	rt.Users = users

	// get games
	games, err := GetAllTournamentGames(t.ID)

	for _, g := range games {
		rt.Games = append(rt.Games, *GameToReturnGame(&g))
	}

	return &rt, nil
}

func ReturnTournamentToTournament(rt types.ReturnTournament) (*types.Tournament, error) {
	var t types.Tournament
	t.ID = rt.ID
	t.Name = rt.Name
	t.Status = rt.Status
	t.CreationDate = rt.CreationDate

	var owner int64
	var users []int64

	user, err := GetUserByUsername(rt.Owner)
	if err != nil {
		return nil, err
	}
	owner = user.ID

	for _, u := range rt.Users {
		user, err := GetUserByUsername(u)
		if err != nil {
			return nil, err
		}

		users = append(users, user.ID)
	}

	t.Owner = owner
	t.Users = users

	return &t, nil
}

func GetTournament(id int64) (*types.Tournament, error) {
	q := `
	SELECT *
	FROM tournaments
	WHERE id=$1
	`

	res := Conn.QueryRow(q, id)

	var t types.Tournament

	err := res.Scan(&t.ID, &t.Name, &t.Owner, &t.Status, pq.Array(&t.Users), pq.Array(&t.Winners), &t.CreationDate)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func GetTournamentList() (*types.TournamentList, error) {
	q := `
	SELECT id, name, owner, users, creation_date
	FROM tournaments
	`

	rows, err := Conn.Query(q)
	if err != nil {
		return nil, err
	}

	var list types.TournamentList

	for rows.Next() {
		var entry types.TournamentInfo
		var ownerid int64
		var users []int64
		err := rows.Scan(&entry.ID, &entry.Name, &ownerid, pq.Array(&users), &entry.CreationDate)
		if err != nil {
			return nil, err
		}

		owner, err := GetUser(ownerid)
		if err != nil {
			return nil, err
		}

		entry.Owner = owner.Username
		entry.UserNumber = len(users)

		list = append(list, entry)
	}

	return &list, nil
}

func GetAllTournamentGames(id int64) ([]types.Game, error) {
	q := `
	SELECT id
	FROM games
	WHERE tournament=$1
	ORDER BY id ASC
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
