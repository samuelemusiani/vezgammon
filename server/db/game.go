package db

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
	CREATE TABLE IF NOT EXISTS games(
		id 		SERIAL PRIMARY KEY,
		p1id 	INTEGER REFERENCES users(username),
		p1elo	INTEGER,
		p2id	INTEGER REFERENCES users(username),
		p2elo 	INTEGER,

		start 	TIMESTAMP,
		end 	TIMESTAMP,
		status	GAMESTATE,

		p1checkers INTEGER [],
		P2checkers INTEGER [],

		double 		INTEGER,
		doubleowner INTEGER REFERENCES users(id),

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
