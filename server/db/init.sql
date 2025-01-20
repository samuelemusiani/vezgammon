CREATE TABLE IF NOT EXISTS users(
	id SERIAL PRIMARY KEY,
	username BPCHAR UNIQUE NOT NULL,
	password BPCHAR NOT NULL,
	firstname BPCHAR NOT NULL,
	lastname BPCHAR,
	mail BPCHAR UNIQUE,
	elo INTEGER NOT NULL,
	avatar BPCHAR NOT NULL,
	is_bot BOOL DEFAULT FALSE
)
CREATE TABLE IF NOT EXISTS sessions (
	id SERIAL PRIMARY KEY,
	user_id INTEGER REFERENCES users(id),
	token TEXT UNIQUE NOT NULL,
	expires_at TIMESTAMP NOT NULL
)
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

	dices INTEGER [],
	tournament INTEGER REFERENCES tournaments(id) ON DELETE SET NULL DEFAULT NULL
)
CREATE TABLE IF NOT EXISTS tournaments(
	id SERIAL PRIMARY KEY,
	name BPCHAR,
	owner INTEGER REFERENCES users(id),
	status BPCHAR DEFAULT 'waiting',
	users INTEGER [],
	winners INTEGER [] DEFAULT '{}'::INTEGER[],
	creation_date TIMESTAMP
)
CREATE TABLE IF NOT EXISTS turns(
	id 			SERIAL PRIMARY KEY,
	game_id 	INTEGER REFERENCES games(id),
	user_id	  	INTEGER REFERENCES users(id),
	time		TIMESTAMP,
	dices		INTEGER [],
	double		BOOL,
	moves	    INTEGER []
)
