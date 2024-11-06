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
