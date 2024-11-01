package db

import "vezgammon/server/types"

func initUser() error {
	q := `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY, 
		username BPCHAR NOT NULL, 
		password BPCHAR NOT NULL,
		firstname BPCHAR NOT NULL,
		lastname BPCHAR,
		mail BPCHAR UNIQUE
	)`
	_, err := conn.Exec(q)
	return err
}

func GetUsers() ([]types.User, error) {
	q := "SELECT * FROM users"
	rows, err := conn.Query(q)
	if err != nil {
		return nil, err
	}

	var users []types.User

	for rows.Next() {
		var tmp types.User
		var pass string
		err = rows.Scan(&tmp.ID, &tmp.Username, &pass, &tmp.Firstname, &tmp.Lastname, &tmp.Mail)
		if err != nil {
			return nil, err
		}

		users = append(users, tmp)
	}

	return users, nil
}

func LoginUser(username string, password string) (*types.User, error) {
	q := `SELECT id, username, password 
        FROM users 
        WHERE username = $1 AND password = $2`
	row, err := conn.Query(q)
	if err != nil {
		return nil, err
	}

	var tmp types.User
	var pass string
	err = row.Scan(&tmp.ID, &tmp.Username, &pass, &tmp.Firstname, &tmp.Lastname, &tmp.Mail)
	if err != nil {
		return nil, err
	}

	return &tmp, nil
}

func CreateUser(u types.User, password string) (*types.User, error) {
	q := `INSERT INTO users(username, password, firstname, lastname, mail) VALUES($1, $2, $3, $4, $5) RETURNING id`
	res := conn.QueryRow(q, u.Username, password, u.Firstname, u.Lastname, u.Mail)

	var id int64
	err := res.Scan(&id)
	if err != nil {
		return nil, err
	}

	u.ID = id
	return &u, nil
}
