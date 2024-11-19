package db

import (
	"errors"
	"log/slog"
	"math"
	"time"
	"vezgammon/server/types"

	"github.com/lib/pq"
)

// user_id : elo
var Matchmaking = make(map[int64]int64)

func initGame() error {
	q := `
	CREATE TABLE IF NOT EXISTS games(
		id 		  SERIAL PRIMARY KEY,
		p1_id	 	INTEGER REFERENCES users(id),
		p1elo	  INTEGER,
		p2_id		INTEGER REFERENCES users(id),
		p2elo 	INTEGER,

		start 	TIMESTAMP,
		endtime TIMESTAMP,
		status	BPCHAR DEFAULT 'open',

		p1checkers INTEGER [] DEFAULT ARRAY [0, 0, 0, 0, 0, 0, 5, 0, 3, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2],
		P2checkers INTEGER [] DEFAULT ARRAY [0, 0, 0, 0, 0, 0, 5, 0, 3, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2],

		double_value 		INTEGER DEFAULT 1,
		double_owner    BPCHAR DEFAULT 'all',
    want_to_double  BOOL DEFAULT FALSE,

    current_player  BPCHAR DEFAULT 'p1',

		dices INTEGER []
	)
	`
	_, err := Conn.Exec(q)
	if err != nil {
		return err
	}

	q = `
	CREATE TABLE IF NOT EXISTS turns(
		id 			SERIAL PRIMARY KEY,
		game_id 	INTEGER REFERENCES games(id),
		user_id	  	INTEGER REFERENCES users(id),
		time		TIMESTAMP,
		dices		INTEGER [],
		double		BOOL,
		moves	    INTEGER []
	)
	`
	_, err = Conn.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

// This function is called every time a new game is created
func CreateGame(g types.Game) (*types.Game, error) {
	q := `
    INSERT INTO games (
        p1_id, p1elo, p2_id, p2elo,
        start, endtime, current_player, dices
    )
    VALUES (
        $1, $2, $3, $4,
        $5, $6, $7, $8
    )
    RETURNING id
    `
	res := Conn.QueryRow(
		q,
		g.Player1, g.Elo1, g.Player2, g.Elo2,
		g.Start, g.End, g.CurrentPlayer, pq.Array(g.Dices),
	)

	var id int64

	err := res.Scan(&id)
	if err != nil {
		slog.With("creating game: ", err).Debug("error")
		return nil, err
	}

	g.ID = id

	return &g, nil
}

func UpdateGame(g *types.Game) error {
	q := `
		UPDATE games
		SET
      endtime		      = $1,
			status		      = $2,
			p1checkers	    = $3,
			p2checkers	    = $4,
			double_value	  = $5,
			double_owner	  = $6,
			want_to_double	= $7,
			current_player	= $8,
      dices = $9
		WHERE id = $10
		`

	_, err := Conn.Exec(q, g.End, g.Status, pq.Array(g.P1Checkers), pq.Array(g.P2Checkers), g.DoubleValue, g.DoubleOwner, g.WantToDouble, g.CurrentPlayer, pq.Array(g.Dices), g.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetGame(id int64) (*types.Game, error) {
	var g types.Game

	q := "SELECT * FROM games WHERE id = $1"

	row := Conn.QueryRow(q, id)

	// Tmp variables for avoid error type in the db [int8/int]
	var p1CheckersDB pq.Int64Array
	var p2CheckersDB pq.Int64Array
	var dices pq.Int64Array

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
		&g.CurrentPlayer,
		&dices)

	if err != nil {
		return nil, err
	}

	// Conversion of p1Checkers from int8 into int
	for i, v := range p1CheckersDB {
		g.P1Checkers[i] = int8(v)
	}
	for i, v := range p2CheckersDB {
		g.P2Checkers[i] = int8(v)
	}

	if len(dices) < 2 {
		return nil, errors.New("dices are more than 2")
	}

	g.Dices[0] = int(dices[0])
	g.Dices[1] = int(dices[1])

	return &g, nil
}

func GetCurrentGame(user_id int64) (*types.ReturnGame, error) {
	q := `
	SELECT
        g.id,
        u1.username AS p1_username,
        g.p1elo,
        u2.username AS p2_username,
        g.p2elo,
        g.start,
        g.endtime,
        g.status,
        g.p1checkers,
        g.p2checkers,
        g.double_value,
        g.double_owner,
        g.want_to_double,
        g.current_player
    FROM
    	games g
    JOIN
    	users u1 ON g.p1_id = u1.id
    JOIN
    	users u2 ON g.p2_id = u2.id
    WHERE
    	g.status = 'open' AND (g.p1_id = $1 OR g.p2_id = $1)
    LIMIT 1
	`

	row := Conn.QueryRow(q, user_id)

	// Tmp variables for avoid error type in the db [int8/int]
	var p1CheckersDB pq.Int64Array
	var p2CheckersDB pq.Int64Array

	var g types.ReturnGame

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
		&g.CurrentPlayer,
	)

	if err != nil {
		return nil, err
	}

	// Conversion of p1Checkers from int8 into int
	for i, v := range p1CheckersDB {
		g.P1Checkers[i] = int8(v)
	}
	for i, v := range p2CheckersDB {
		g.P2Checkers[i] = int8(v)
	}

	return &g, nil
}

func CreateTurn(t types.Turn) (*types.Turn, error) {
	q := `
	INSERT INTO turns(game_id, user_id, time, dices, double, moves)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`

	res := Conn.QueryRow(q, t.GameId, t.User, t.Time, pq.Array(t.Dices), t.Double, pq.Array(MovesArrayToArray(t.Moves)))
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
	q := "SELECT * FROM turns WHERE game_id = $1 ORDER BY time ASC"
	rows, err := Conn.Query(q, game_id)
	if err != nil {
		return nil, err
	}

	var turns []types.Turn

	var moves pq.Int64Array
	var dices pq.Int64Array

	for rows.Next() {
		var tmp types.Turn
		err = rows.Scan(&tmp.ID, &tmp.GameId, &tmp.User, &tmp.Time, &dices, &tmp.Double, &moves)
		if err != nil {
			return nil, err
		}

		tmp.Dices = types.Dices{int(dices[0]), int(dices[1])}

		tmp.Moves = ArrayToMovesArray([]int64(moves))

		turns = append(turns, tmp)
	}

	return turns, nil
}

func GetLastTurn(game_id int64) (*types.Turn, error) {
	q := "SELECT * FROM turns WHERE game_id = $1 ORDER BY time DESC LIMIT 1"
	row := Conn.QueryRow(q, game_id)

	var turn types.Turn
	var moves pq.Int64Array

	var dices pq.Int64Array
	err := row.Scan(&turn.ID, &turn.GameId, &turn.User, &turn.Time, &dices, &turn.Double, &moves)
	if err != nil {
		return nil, err
	}

	turn.Dices = types.Dices{int(dices[0]), int(dices[1])}

	turn.Moves = ArrayToMovesArray(moves)

	return &turn, nil
}

func SearchGame(uid int64) error {
	u, err := GetUser(uid)
	if err != nil {
		return err
	}
	slog.With("user stats: ", u)

	// start matchmaking
	Matchmaking[uid] = u.Elo

	//cerco l'opponent nel db in base al player e in base a quanto è in queue
	var oppo_id int64
	oppo_id, err = findOpponent(u.Elo, u.ID)
	if err != nil {
		return err
	}

	// remove player from matchmaking map
	delete(Matchmaking, uid)

	oppo, err := GetUser(oppo_id)
	if err != nil {
		return err
	}

	var dices = types.NewDices()
	var CurrentPlayer string
	if dices[0] >= dices[1] {
		CurrentPlayer = types.GameCurrentPlayerP1
	} else {
		CurrentPlayer = types.GameCurrentPlayerP2
	}

	var game types.Game
	game.Player1 = u.ID
	game.Player2 = oppo.ID
	game.Elo1 = u.Elo
	game.Elo2 = oppo.Elo
	game.Start = time.Now()
	game.End = time.Now()
	game.CurrentPlayer = CurrentPlayer
	game.Dices = types.NewDices()

	_, err = CreateGame(game)
	if err != nil {
		return err
	}

	return err
}

func findOpponent(elo int64, uid int64) (int64, error) {
	for key, value := range Matchmaking {
		slog.With("key", key, "uid", uid).Debug("findOpponent")
		if math.Abs(float64(value-elo)) < 200 && key != uid {
			slog.With("player", key).Debug("player found")
			return key, nil
		}
	}

	return 0, errors.New("no opponent found")
}
