package db

import "vezgammon/server/types"

func initUser() error {
	q := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY, 
		username BPCHAR NOT NULL, 
		password BPCHAR NOT NULL,
		firstname BPCHAR NOT NULL,
		lastname BPCHAR,
		mail BPCHAR UNIQUE,
	)`
	_, err := conn.Exec(q)
	return err
}

func CreateUser(u types.User, password string) (*types.User, error) {
	q := `
	INSERT INTO users (username, password, firstname, lastname, mail) values (?, ?, ?, ?, ?)
	`
	res, err := conn.Exec(q, u.Username, password, u.Firstname, u.Lastname, u.Mail)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	u.ID = id
	return &u, nil
}
