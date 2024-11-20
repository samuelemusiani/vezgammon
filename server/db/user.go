package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"vezgammon/server/types"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func initUser() error {
	q := `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		username BPCHAR UNIQUE NOT NULL,
		password BPCHAR NOT NULL,
		firstname BPCHAR NOT NULL,
		lastname BPCHAR,
		mail BPCHAR UNIQUE,
    elo INTEGER NOT NULL
	)`
	_, err := Conn.Exec(q)
	return err
}

func initCookie() error {
	q := `
    CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL
  )
	`
	_, err := Conn.Exec(q)
	return err
}

func GetUsers() ([]types.User, error) {
	q := "SELECT * FROM users"
	rows, err := Conn.Query(q)
	if err != nil {
		return nil, err
	}

	var users []types.User

	for rows.Next() {
		var tmp types.User
		var pass string
		err = rows.Scan(&tmp.ID, &tmp.Username, &pass, &tmp.Firstname, &tmp.Lastname, &tmp.Mail, &tmp.Elo)
		if err != nil {
			return nil, err
		}

		users = append(users, tmp)
	}

	return users, nil
}

func LoginUser(username string, password string) (*types.User, error) {
	q := "SELECT id, username, firstname, lastname, mail, password, elo FROM users "
	if strings.Contains(username, "@") {
		q = q + "WHERE mail = $1"
	} else {
		q = q + "WHERE username = $1"
	}

	var tmp types.User
	var pass string
	err := Conn.QueryRow(q, username).Scan(
		&tmp.ID,
		&tmp.Username,
		&tmp.Firstname,
		&tmp.Lastname,
		&tmp.Mail,
		&pass,
		&tmp.Elo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("utente non trovato")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
	if err != nil {
		return nil, err
	}

	return &tmp, nil
}

func GenerateSessionToken() string {
	return uuid.NewString()
}

func SaveSessionToken(userID int64, token string) error {
	q := `INSERT INTO sessions (user_id, token, expires_at)
          VALUES ($1, $2, $3)`

	expiresAt := time.Now().Add(1 * time.Hour)
	_, err := Conn.Exec(q, userID, token, expiresAt)
	return err
}

func ValidateSessionToken(token string) (int64, error) {
	q := `SELECT user_id FROM sessions
          WHERE token = $1 AND expires_at > NOW()`

	var userID int64
	err := Conn.QueryRow(q, token).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func CreateUser(u types.User, password string) (types.User, error) {
	q := `INSERT INTO users(username, password, firstname, lastname, mail, elo) VALUES($1, $2, $3, $4, $5, $6) RETURNING id`
	res := Conn.QueryRow(q, u.Username, password, u.Firstname, u.Lastname, u.Mail, types.DefaultElo)

	var id int64
	err := res.Scan(&id)
	if err != nil {
		return u, err
	}

	u.ID = id
	return u, nil
}

func Logout(sessionToken string) error {
	// Rimuovi il token dal database
	q := `DELETE FROM sessions WHERE token = $1`
	_, err := Conn.Exec(q, sessionToken)
	return err
}

func GetUser(user_id any) (*types.User, error) {
	q := `SELECT username, firstname, lastname, mail, elo
          FROM users
          WHERE id = $1`

	var tmp types.User
	err := Conn.QueryRow(q, user_id).Scan(
		&tmp.Username,
		&tmp.Firstname,
		&tmp.Lastname,
		&tmp.Mail,
		&tmp.Elo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("utente non trovato")
		}
		return nil, err
	}

	return &tmp, nil
}
