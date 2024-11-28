package db

import (
	"cmp"
	"slices"
	"vezgammon/server/types"

	"github.com/lib/pq"
)

func InitTournament() error {
	q := `
	CREATE TABLE IF NOT EXISTS tournaments(
		id SERIAL PRIMARY KEY,
		name BPCHAR,
		owner INTEGER REFERENCES users(id),
		status BPCHAR DEFAULT 'open',
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
	INSERT INTO tournaments(name, owner, status, users)
	values($1, $2, $3, $4)
	RETURNING id
	`

	res := Conn.QueryRow(
		q,
		t.Name, t.Owner, t.Status, pq.Array(t.Users),
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
	SET status=$1, users=$2
	WHERE id=$3
	`

	_, err := Conn.Exec(
		q,
		t.Status, pq.Array(t.Users), t.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func calcLeaderBoard(games []types.ReturnGame) types.LeaderBoard {
	type mape struct {
		win  int
		lose int
	}

	m := make(map[string]mape)

	for _, g := range games {
		switch g.Status {
		case types.GameStatusWinP1:
			m[g.Player1] = mape{m[g.Player1].win + 1, m[g.Player1].lose}
			m[g.Player2] = mape{m[g.Player2].win, m[g.Player2].lose + 1}
		case types.GameStatusWinP2:
			m[g.Player2] = mape{m[g.Player2].win + 1, m[g.Player2].lose}
			m[g.Player1] = mape{m[g.Player1].win, m[g.Player1].lose + 1}
		}
	}

	var list []types.LeaderBoardEntry
	for name, score := range m {
		entry := types.LeaderBoardEntry{
			User: name,
			Win:  score.win,
			Lose: score.lose,
		}

		list = append(list, entry)
	}

	// sort, most wins first
	slices.SortFunc(list, func(i, j types.LeaderBoardEntry) int {
		sum1 := i.Win - i.Lose
		sum2 := j.Win - j.Lose

		return cmp.Compare(sum1, sum2)
	})

	return list
}

func TournamentToReturnTournament(t types.Tournament) (*types.ReturnTournament, error) {
	var rt types.ReturnTournament

	rt.ID = t.ID
	rt.Name = t.Name
	rt.Status = t.Status

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

	// calc leaderboard
	rt.LeaderBoard = calcLeaderBoard(rt.Games)

	return &rt, nil
}

func ReturnTournamentToTournament(rt types.ReturnTournament) (*types.Tournament, error) {
	var t types.Tournament
	t.ID = rt.ID
	t.Name = rt.Name
	t.Status = rt.Status

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

	err := res.Scan(&t.ID, &t.Name, &t.Owner, &t.Status, pq.Array(&t.Users))
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func GetTournamentList() (*types.TournamentList, error) {
	q := `
	SELECT id, name, owner
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
		err := rows.Scan(&entry.ID, &entry.Name, &ownerid)
		if err != nil {
			return nil, err
		}

		owner, err := GetUser(ownerid)
		if err != nil {
			return nil, err
		}

		entry.Owner = owner.Username

		list = append(list, entry)
	}

	return &list, nil
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
