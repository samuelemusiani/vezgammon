package db

import (
	"time"
	"vezgammon/server/types"
)

func initGame() error {

	var q string
	q = `
	CREATE TYPE IF NOT EXISTS gamestate AS
	ENUM('open', 'winp1', 'winp2')
	`
	_, err := conn.Exec(q)
	if err != nil {
		return err
	}

	q = `
	CREATE TYPE IF NOT EXISTS doubleowner AS
	ENUM('all', 'p1', 'p2')
	`
	_, err = conn.Exec(q)
	if err != nil {
		return err
	}

	q = `
	CREATE TABLE IF NOT EXISTS games(
		id 		SERIAL PRIMARY KEY,
		p1	INTEGER REFERENCES users(username),
		p1elo	INTEGER,
		p2	INTEGER REFERENCES users(username),
		p2elo 	INTEGER,

		start 	TIMESTAMP,
		end 	TIMESTAMP,
		status	GAMESTATE DEFAULT 'open',

		p1checkers INTEGER [] DEFAULT {0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 3, 0, 5, 0, 0, 0, 0, 0},
		P2checkers INTEGER [] DEFAULT {0, 0, 0, 0, 0, 0, 5, 0, 3, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},

		double 		INTEGER DEFAULT 1,
		doubleowner INTEGER DOUBLEOWNER DEFAULT 'all',

		nextdices INTEGER []
	)
	`
	_, err = conn.Exec(q)
	if err != nil {
		return err
	}

	q = `
	CREATE TABLE IF NOT EXISTS turns(
		id 		SERIAL PRIMARY KEY,
		game 	INTEGER REFERENCES games(id),
		user	INTEGER REFERENCES users(username),
		time	TIMESTAMP,

		moves	INTEGER [][]
	)
	`
	_, err = conn.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func CreateGame(g types.Game) (*types.Game, error) {
	q := `
	INSERT INTO games (p1, p1elo, p2, p2elo, start, nextdices) 
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`

	start := time.Now()
	dices := types.NewDices()

	res := conn.QueryRow(q, g.Player1, 0, g.Player2, 0, start, dices)
	var id int64
	err := res.Scan(&id)
	if err != nil {
		return nil, err
	}

	g.ID = id
	return &g, nil
}

func UpdateGame(g types.Game) error {
	q := `
	UPDATE games 
	SET	end			= $1, 
		status		= $2, 
		p1checkers	= $3, 
		p2checkers	= $4, 
		double		= $5, 
		doubleowner	= $6, 
		nextdices	= $7 
	WHERE id = $8
	`

	end := time.Now()
	dices := types.NewDices()

	_, err := conn.Exec(q, end, g.Status, g.P1Checkers, g.P2Checkers, g.Double, g.DoubleOwner, dices)
	if err != nil {
		return err
	}

	return nil
}

func GetGame(id int64) (*types.Game, *[2]int, error) {
	var g types.Game

	q := "SELECT * FROM games WHERE id = $1"

	row := conn.QueryRow(q, id)

	var nextdices [2]int
	err := row.Scan(
		&g.ID,
		&g.Player1,
		&g.Elo1,
		&g.Player2,
		&g.Elo2,
		&g.Start,
		&g.End,
		&g.Status,
		&g.P1Checkers,
		&g.P2Checkers,
		&g.Double,
		&g.DoubleOwner,
		&nextdices)
	if err != nil {
		return nil, nil, err
	}

	return &g, &nextdices, nil
}

func CreateTurn(t types.Turn) (*types.Turn, error) {
	q := `
	INSERT INTO turns(id, game, user, time, moves)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	res := conn.QueryRow(q, t.ID, t.GameId, t.User, t.Moves)
	var id int64
	err := res.Scan(&id)
	if err != nil {
		return nil, err
	}
	t.ID = id
	return &t, nil
}

// return all turn for a game last to first
func GetTurns(game_id int64) ([]types.Turn, error) {
	q := "SELECT * FROM turns WHERE game = $1 ORDER BY time ASC"
	rows, err := conn.Query(q, game_id)
	if err != nil {
		return nil, err
	}

	var turns []types.Turn

	for rows.Next() {
		var tmp types.Turn
		err = rows.Scan(&tmp.ID, &tmp.GameId, &tmp.User, &tmp.Time, &tmp.Moves)
		if err != nil {
			return nil, err
		}

		turns = append(turns, tmp)
	}

	return turns, nil
}

func GetLastTurn(game_id int64) (*types.Turn, error) {
	q := "SELECT * FROM turns WHERE game = $1 ORDER BY time ASC LIMIT 1"
	row := conn.QueryRow(q, game_id)

	var turn types.Turn
	err := row.Scan(&turn.ID, &turn.GameId, &turn.User, &turn.Time, &turn.Moves)
	if err != nil {
		return nil, err
	}

	return &turn, nil
}
