package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"vezgammon/server/types"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var UserNotFound = errors.New("Utente non trovato")

var easyBotID int64 = -1
var mediumBotID int64 = -1
var hardBotID int64 = -1

func GetEasyBotID() int64 {
	if easyBotID == -1 {
		q := `SELECT id FROM users WHERE username = 'Enzo' AND is_bot = TRUE`
		err := Conn.QueryRow(q).Scan(&easyBotID)
		if err != nil {
			slog.With("err", err).Error("Getting easy bot id")
			panic("Cannot get easy bot ID")
		}
	}
	return easyBotID
}

func GetMediumBotID() int64 {
	if mediumBotID == -1 {
		q := `SELECT id FROM users WHERE username = 'Caterina' AND is_bot = TRUE`
		err := Conn.QueryRow(q).Scan(&mediumBotID)
		if err != nil {
			slog.With("err", err).Error("Getting medium bot id")
			panic("Cannot get medium bot ID")
		}
	}
	return mediumBotID
}

func GetHardBotID() int64 {
	if hardBotID == -1 {
		q := `SELECT id FROM users WHERE username = 'Giovanni' AND is_bot = TRUE`
		err := Conn.QueryRow(q).Scan(&hardBotID)
		if err != nil {
			slog.With("err", err).Error("Getting hard bot id")
			panic("Cannot get hard bot ID")
		}
	}

	return hardBotID
}

func GetBotLevel(id int64) int {
	if id == GetEasyBotID() {
		return 1
	} else if id == GetMediumBotID() {
		return 2
	} else if id == GetHardBotID() {
		return 3
	}

	return 0
}

func initUser() error {
	q := `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		username BPCHAR UNIQUE NOT NULL,
		password BPCHAR NOT NULL,
		firstname BPCHAR NOT NULL,
		lastname BPCHAR,
		mail BPCHAR UNIQUE,
    elo INTEGER NOT NULL,
    is_bot BOOL DEFAULT FALSE
	)`
	_, err := Conn.Exec(q)
	if err != nil {
		return err
	}

	// Insert easy bot
	err = insertBotIfNotExists("Enzo", "Re", "Enzo", "enzo@vezgammon.it", 1000)
	if err != nil {
		return err
	}

	err = insertBotIfNotExists("Caterina", "Caterina", "De Vigri", "caterina@vezgammon.it", 2000)
	if err != nil {
		return err
	}

	err = insertBotIfNotExists("Giovanni", "Giovanni", "Bentivoglio", "giovanni@vezgammon.it", 3000)
	if err != nil {
		return err
	}

	return nil
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
	q := "SELECT * FROM users WHERE is_bot = FALSE"
	rows, err := Conn.Query(q)
	if err != nil {
		return nil, err
	}

	var users []types.User

	for rows.Next() {
		var tmp types.User
		var pass string
		err = rows.Scan(&tmp.ID, &tmp.Username, &pass, &tmp.Firstname, &tmp.Lastname, &tmp.Mail, &tmp.Elo, &tmp.IsBot)
		if err != nil {
			return nil, err
		}

		users = append(users, tmp)
	}

	return users, nil
}

func LoginUser(username string, password string) (*types.User, error) {
	q := "SELECT id, username, firstname, lastname, mail, password, elo, is_bot FROM users "
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
		&tmp.IsBot,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		return nil, err
	}

	if tmp.IsBot {
		return nil, fmt.Errorf("User is a bot")
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
	q := `INSERT INTO users(username, password, firstname, lastname, mail, elo, is_bot) 
    VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	res := Conn.QueryRow(q, u.Username, password, u.Firstname, u.Lastname, u.Mail, types.DefaultElo, u.IsBot)

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

func GetUserByUsername(username string) (*types.User, error) {
	q := `SELECT id, username, firstname, lastname, mail, elo
          FROM users
          WHERE username = $1`

	var tmp types.User
	err := Conn.QueryRow(q, username).Scan(
		&tmp.ID,
		&tmp.Username,
		&tmp.Firstname,
		&tmp.Lastname,
		&tmp.Mail,
		&tmp.Elo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		return nil, err
	}

	return &tmp, nil
}

func GetUser(userId int64) (*types.User, error) {
	q := `SELECT username, firstname, lastname, mail, elo
          FROM users
          WHERE id = $1`

	var tmp types.User
	err := Conn.QueryRow(q, userId).Scan(
		&tmp.Username,
		&tmp.Firstname,
		&tmp.Lastname,
		&tmp.Mail,
		&tmp.Elo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		return nil, err
	}

	return &tmp, nil
}

func insertBotIfNotExists(username, firstname, lastname, mail string, elo int64) error {
	_, err := GetUserByUsername(username)

	slog.With("err", err).Debug("Getting bot user")

	if err != nil {
		if errors.Is(err, UserNotFound) {
			// Insert the bot
			_, err := CreateUser(types.User{
				Username:  username,
				Firstname: firstname,
				Lastname:  lastname,
				Mail:      mail,
				Elo:       elo,
				IsBot:     true,
			}, "1234")

			if err != nil {
				slog.With("err", err).Error("Creating bot user")
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func GetStats(user_id int64) (*types.Stats, error) {
	stats := new(types.Stats)

	//partite dal db con current state won lost del player(user_id)
	u, err := GetUser(user_id)
	if err != nil {
		return nil, err
	}

	var gp []types.ReturnGame
	gp, err = GetAllGameFromUser(user_id)
	if err != nil {
		return nil, err
	}

	stats.Gameplayed = gp
	stats.Tournament = 0 // not implemented yet
	for i, game := range gp {
		if game.GameType == types.GameTypeBot {
			stats.Cpu++
		} else if game.GameType == types.GameTypeLocal {
			stats.Local++
		} else if game.GameType == types.GameTypeOnline {
			stats.Online++
		}

		if (game.Status == types.GameStatusWinP1 && game.Player1 == u.Username) || (game.Status == types.GameStatusWinP2 && game.Player2 == u.Username) {
			stats.Won++
		} else {
			stats.Lost++
		}
		if game.Player1 == u.Username {
			stats.Elo[i] = game.Elo1
		} else {
			stats.Elo[i] = game.Elo2
		}
	}
	stats.Elo = append(stats.Elo, u.Elo)

	if len(stats.Gameplayed) == 0 {
		stats.Winrate = 0
	} else {
		stats.Winrate = float32(stats.Won / int64(len(stats.Gameplayed)))
	}
	slog.With("stats", stats).Debug("Statistiche")

	return stats, nil
}
