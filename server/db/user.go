package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"math"
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
    avatar BPCHAR NOT NULL,
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
		err = rows.Scan(
            &tmp.ID,
            &tmp.Username,
            &pass,
            &tmp.Firstname,
            &tmp.Lastname,
            &tmp.Mail,
            &tmp.Elo,
            &tmp.Avatar,
            &tmp.IsBot,
        )
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
	q := `INSERT INTO users(username, password, firstname, lastname, mail, elo, avatar, is_bot)
    VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	res := Conn.QueryRow(q, u.Username, password, u.Firstname, u.Lastname, u.Mail, types.DefaultElo, u.Avatar, u.IsBot)

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
	q := `SELECT id, username, firstname, lastname, mail, elo, avatar
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
		&tmp.Avatar,
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
	q := `SELECT username, firstname, lastname, mail, elo, avatar
          FROM users
          WHERE id = $1`

	var tmp types.User
	err := Conn.QueryRow(q, userId).Scan(
		&tmp.Username,
		&tmp.Firstname,
		&tmp.Lastname,
		&tmp.Mail,
		&tmp.Elo,
		&tmp.Avatar,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		return nil, err
	}

	tmp.ID = userId

	return &tmp, nil
}

func insertBotIfNotExists(username, firstname, lastname, mail string, elo int64) error {
	_, err := GetUserByUsername(username)

	slog.With("err", err).Debug("Getting bot user")

	var avatar string
	switch username {
	case "Enzo":
		avatar = "Andrea"
	case "Giovanni":
		avatar = "Luis"

	default:
		avatar = username
	}

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
				Avatar:    "https://api.dicebear.com/6.x/avataaars/svg?seed=" + avatar,
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

func GetStats(userID int64) (*types.Stats, error) {
	stats := new(types.Stats)

	//partite dal db con current state won lost del player(userID)
	u, err := GetUser(userID)
	if err != nil {
		return nil, err
	}

	var gp []types.ReturnGame
	gp, err = GetAllGameFromUser(userID)
	if err != nil {
		return nil, err
	}

	stats.Gameplayed = gp
	stats.Tournament = 0 // not implemented yet
	for _, game := range gp {
		if game.GameType == types.GameTypeBot {
			stats.Cpu++
		} else if game.GameType == types.GameTypeLocal {
			stats.Local++
		} else if game.GameType == types.GameTypeOnline {
			stats.Online++
		} else {
            slog.With("game type", game.GameType).Error("Unknown game type")
        }

		if game.GameType == types.GameTypeOnline { // no sense to count local games in winrate statistics
		    isPlayer1Winner := game.Status == types.GameStatusWinP1 && game.Player1 == u.Username
            isPlayer2Winner := game.Status == types.GameStatusWinP2 && game.Player2 == u.Username

			if isPlayer1Winner || isPlayer2Winner {
				stats.Won++
			} else {
				stats.Lost++
			}

			if game.Player1 == u.Username {
				slog.With("elo", game.Elo1, "game", u.Username, "players", game.Player1, game.Player2).Debug("dio sto elo cazzo 1")
				stats.Elo = append(stats.Elo, game.Elo1)
			} else {
				slog.With("elo", game.Elo2, "game", u.Username).Debug("dio sto elo cazzo 2")
				stats.Elo = append(stats.Elo, game.Elo2)
			}
		}
	}

	// current elo after last game
	stats.Elo = append(stats.Elo, u.Elo)

	if stats.Online == 0 {
		stats.Winrate = 0
	} else {
		// online games only
		stats.Winrate = float32(math.Floor(float64(100*float32(stats.Won)/float32(stats.Online))*100)) / 100
	}

	stats.Leaderboard, err = getLeaderboard()
	if err != nil {
		return nil, err
	}

	slog.With("stats", stats).Debug("Statistiche")
	return stats, nil
}

func UpdateUserElo(userID int64, elo int64) error {
	q := `UPDATE users SET elo = $1 WHERE id = $2`
	_, err := Conn.Exec(q, elo, userID)
	return err
}

func GetBadge(userID int64) (*types.Badge, error) {
	user, err := GetUser(userID)
	if err != nil {
		slog.With("err", err).Debug("Badge")
		return nil, err
	}

	var (
		badge types.Badge

		gp []types.ReturnGame

		gw         int
		gameEnded  int
		homepieces int
	)

	gp, err = GetAllGameFromUser(userID)
	slog.With("gp", gp).Debug("Badge games")

	for _, game := range gp {
		// skip ongoing games
		if game.Status == types.GameStatusOpen {
			continue
		}

		//skip local games
		if game.GameType == types.GameTypeLocal {
			continue
		}

		//bot difficulty
		if game.GameType == types.GameTypeBot {
			// if lost against bot skip game
			isPlayer1Winner := game.Status == types.GameStatusWinP1 && game.Player1 == user.Username
			isPlayer2Winner := game.Status == types.GameStatusWinP2 && game.Player2 == user.Username
			if !(isPlayer1Winner || isPlayer2Winner) {
				continue
			}
			slog.With("game type", game.GameType).Debug("capiamo?")
			p1, e1 := GetUserByUsername(game.Player1)
			if e1 != nil {
				return nil, err
			}
			slog.With("p1", p1).Debug("Badge")

			p2, e2 := GetUserByUsername(game.Player2)
			if e2 != nil {
				return nil, err
			}
			slog.With("p2", p2).Debug("Badge")

			// One return 0 the other is 1/2/3 depends on difficulty
			sum := GetBotLevel(p1.ID) + GetBotLevel(p2.ID)
			switch sum {
			case 1:
				badge.Bot[0] = sum
			case 2:
				badge.Bot[1] = sum
			case 3:
				badge.Bot[2] = sum
			default:
				slog.Debug("2 Humans or 2 Bots")
				err := errors.New("2 Humans or 2 Bots")
				return nil, err
			}
			slog.With("bot", badge.Bot, "sum", sum).Debug("Badge")
			continue
		}

		//only online games
		gameEnded++

		//homepieces
		homepieces += calculateHomePieces(game, user.Username)
		slog.With("home pieces", homepieces).Debug("Badge")

		//game won counter
		isPlayer1Winner := game.Status == types.GameStatusWinP1 && game.Player1 == user.Username
		isPlayer2Winner := game.Status == types.GameStatusWinP2 && game.Player2 == user.Username
		if isPlayer1Winner || isPlayer2Winner {
		    gw++

			//shortest game
			timeDiff := game.End.Sub(game.Start)

			if timeDiff <= 3*time.Minute {
				badge.Wontime[0] = 1
				badge.Wontime[1] = 2
				badge.Wontime[2] = 3
			} else if timeDiff <= 5*time.Minute {
				badge.Wontime[0] = 1
				badge.Wontime[1] = 2
			} else if timeDiff <= 10*time.Minute {
				badge.Wontime[0] = 1
			} else {
                continue
            }
		}
	}

	//homepieces
	if homepieces >= 50 {
		if homepieces <= 100 {
			badge.Homepieces[0] = 1
		} else if homepieces < 200 {
			badge.Homepieces[0] = 1
			badge.Homepieces[1] = 2
		} else {
			badge.Homepieces[0] = 1
			badge.Homepieces[1] = 2
			badge.Homepieces[2] = 3
		}
	}

	//game played
	if gameEnded > 0 {
		if gameEnded <= 10 {
			badge.Gameplayed[0] = 1
		} else if gameEnded <= 100 {
			badge.Gameplayed[0] = 1
			badge.Gameplayed[1] = 2
		} else {
			badge.Gameplayed[0] = 1
			badge.Gameplayed[1] = 2
			badge.Gameplayed[2] = 3
		}
	}

	// game won
	if gw > 0 {
		if gw <= 10 {
			badge.Wongames[0] = 1
		} else if gw <= 50 {
			badge.Wongames[0] = 1
			badge.Wongames[1] = 2
		} else {
			badge.Wongames[0] = 1
			badge.Wongames[1] = 2
			badge.Wongames[2] = 3
		}
	}

	//elo
	if user.Elo > 1000 {
		if user.Elo < 1200 {
			badge.Elo[0] = 1
		} else if user.Elo < 1400 {
			badge.Elo[0] = 1
			badge.Elo[1] = 2
		} else {
			badge.Elo[0] = 1
			badge.Elo[1] = 2
			badge.Elo[2] = 3
		}
	}

	slog.With("badge", badge).Debug("BADGE")
	return &badge, nil
}

func calculateHomePieces(game types.ReturnGame, u string) int {
	var piecesOnBoard int

	if u == game.Player1 {
		for _, t := range game.P1Checkers {
			piecesOnBoard += int(t)
		}
	} else {
		for _, t := range game.P2Checkers {
			piecesOnBoard += int(t)
		}
	}

	return 15 - piecesOnBoard
}

func getLeaderboard() ([]types.LeaderboardUser, error) {
	q := `
	SELECT username, elo
	FROM users
	WHERE is_bot = FALSE
	ORDER BY elo DESC
	`

	rows, err := Conn.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lb []types.LeaderboardUser
	for rows.Next() {
		var user types.LeaderboardUser
		err := rows.Scan(&user.Username, &user.Elo)
		if err != nil {
			return nil, err
		}
		lb = append(lb, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lb, nil
}

func ChangeAvatar(userID int64, avatar string) error {
	slog.With("avatar", avatar).Debug("Avatar")
	q := `
    UPDATE users
    SET avatar = $2
    WHERE id = $1
    `
	_, err := Conn.Exec(q, userID, avatar)
	if err != nil {
		return err
	}

	return nil
}

func ChangePass(username, newPass, oldPass string) error {
	_, err := LoginUser(username, oldPass)
	if err != nil {
		return fmt.Errorf("incorrect old password: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	q := `
    UPDATE users
    SET password = $2
    WHERE username = $1
    `
	_, err = Conn.Exec(q, username, string(hash))
	if err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}

	return nil
}
