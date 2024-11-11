package db

import (
	"log/slog"
	"time"
	"vezgammon/server/types"

	"github.com/lib/pq"
)

func initGame() error {
	q := `
	CREATE TABLE IF NOT EXISTS games(
		id 		SERIAL PRIMARY KEY,
		p1	 	INTEGER REFERENCES users(id),
		p1elo	INTEGER,
		p2		INTEGER REFERENCES users(id),
		p2elo 	INTEGER,

		start 	TIMESTAMP,
		endtime TIMESTAMP,
		status	BPCHAR DEFAULT 'open',

		p1checkers INTEGER [] DEFAULT ARRAY [0, 0, 0, 0, 0, 0, 5, 0, 3, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2],
		P2checkers INTEGER [] DEFAULT ARRAY [0, 0, 0, 0, 0, 0, 5, 0, 3, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2],

		double_value 		INTEGER DEFAULT 1,
		double_owner BPCHAR DEFAULT 'all',
    want_to_double BOOL DEFAULT FALSE,

		dices INTEGER []
	)
	`
	_, err := conn.Exec(q)
	if err != nil {
		return err
	}

	q = `
	CREATE TABLE IF NOT EXISTS turns(
		id 			SERIAL PRIMARY KEY,
		game 		INTEGER REFERENCES games(id),
		playedby	INTEGER REFERENCES users(id),
		time		TIMESTAMP,
		dices		INTEGER [],
		double		BOOL,
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
    INSERT INTO games (
        p1, p1elo, p2, p2elo,
        start, endtime, nextdices,
        p1checkers, p2checkers,
        double, doubleowner, status
    )
    VALUES (
        $1, $2, $3, $4,
        $5, $6, $7,
        $8, $9,
        $10, $11, $12
    )
    RETURNING id, start
    `

	start := time.Now()

	// placeholder for endtime value, can't be nil (i think)
	endtime := time.Now()
	dices := types.NewDices()

	// Default checker positions
	p1Checkers := [25]int8{0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 3, 0, 5, 0, 0, 0, 0, 0}
	p2Checkers := [25]int8{0, 0, 0, 0, 0, 0, 5, 0, 3, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}

	res := conn.QueryRow(
		q,
		g.Player1, g.Elo1, g.Player2, g.Elo2,
		start, endtime, pq.Array(dices),
		pq.Array(p1Checkers), pq.Array(p2Checkers),
		1, "all", "open",
	)

	var id int64
	var startTime time.Time

	err := res.Scan(&id, &startTime)
	if err != nil {
		slog.With("creating game: ", err).Debug("error")
		return nil, err
	}

	g.ID = id
	g.Start = startTime
	g.P1Checkers = p1Checkers
	g.P2Checkers = p2Checkers
	g.DoubleValue = 1
	g.DoubleOwner = "all"
	g.Status = "open"

	return &g, nil
}

func UpdateGame(g types.Game) error {
	q := `
	UPDATE games
	SET	endtime		= $1,
		status		= $2,
		p1checkers	= $3,
		p2checkers	= $4,
		double_value		= $5,
		doubleowner	= $6,
		nextdices	= $7
	WHERE id = $8
	`

	end := time.Now()
	dices := types.NewDices()

	_, err := conn.Exec(q, pq.FormatTimestamp(end), g.Status, pq.Array(g.P1Checkers), pq.Array(g.P2Checkers), g.DoubleValue, g.DoubleOwner, pq.Array(dices))
	if err != nil {
		return err
	}

	return nil
}

func GetGame(id int64) (*types.Game, *[2]int8, error) {
	var g types.Game

	q := "SELECT * FROM games WHERE id = $1"

	row := conn.QueryRow(q, id)

	// Tmp variables for avoid error type in the db [int8/int]
	var p1CheckersDB pq.Int64Array
	var p2CheckersDB pq.Int64Array
	var nextdicesDB pq.Int64Array

	err := row.Scan(
		&g.ID,
		&g.Player1,
		&g.Elo1,
		&g.Player2,
		&g.Elo2,
		&g.Start,
		&g.End,
		&g.Status,
		&p1CheckersDB,
		&p2CheckersDB,
		&g.DoubleValue,
		&g.DoubleOwner,
		&g.WantToDouble,
		&nextdicesDB)

	if err != nil {
		return nil, nil, err
	}

	// Conversion of p1Checkers from int8 into int
	for i, v := range p1CheckersDB {
		g.P1Checkers[i] = int8(v)
	}
	for i, v := range p2CheckersDB {
		g.P2Checkers[i] = int8(v)
	}

	var nextdices [2]int8
	for i, v := range nextdicesDB {
		nextdices[i] = int8(v)
	}

	return &g, &nextdices, nil
}

func CreateTurn(t types.Turn) (*types.Turn, error) {
	q := `
	INSERT INTO turns(id, game, user, time, dices, double, moves)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	res := conn.QueryRow(q, t.ID, t.GameId, t.User, t.Time, t.Dices, t.Double, t.Moves)
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
		err = rows.Scan(&tmp.ID, &tmp.GameId, &tmp.User, &tmp.Time, &tmp.Dices, &tmp.Double, &tmp.Moves)
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
	err := row.Scan(&turn.ID, &turn.GameId, &turn.User, &turn.Time, &turn.Dices, &turn.Double, &turn.Moves)
	if err != nil {
		return nil, err
	}

	return &turn, nil
}
